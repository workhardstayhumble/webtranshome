package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Jasa (KEPUTUSAN #20 frontend): layanan pemasangan/perbaikan. Harga tidak disimpan;
// konsultasi lewat WhatsApp. Kontrak: frontend/docs/FRONTEND_BUILD_PLAN.md "Detail kontrak jasa".
type serviceRow struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	Slug        string         `db:"slug"`
	Summary     string         `db:"summary"`
	Description sql.NullString `db:"description"`
	ImageURL    sql.NullString `db:"image_url"`
	ServiceArea sql.NullString `db:"service_area"`
	SortOrder   int64          `db:"sort_order"`
	IsActive    bool           `db:"is_active"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
	Categories  string         `db:"categories"`
}

// categories = JSON [{id,name,slug,parent_id,is_active}] kategori terkait.
const serviceColumns = `s.id,s.name,s.slug,s.summary,s.description,s.image_url,s.service_area,s.sort_order,s.is_active,s.created_at,s.updated_at,
 COALESCE((SELECT json_agg(json_build_object('id',c.id,'name',c.name,'slug',c.slug,'parent_id',c.parent_id,'is_active',c.is_active) ORDER BY c.sort_order,c.name)
  FROM crm_schema.service_categories sc JOIN crm_schema.categories c ON c.id=sc.category_id WHERE sc.service_id=s.id),'[]')::text AS categories`

type serviceCategory struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	ParentID *int64 `json:"parent_id"`
	IsActive bool   `json:"is_active"`
}

func (s serviceRow) categories() []serviceCategory {
	var out []serviceCategory
	_ = json.Unmarshal([]byte(s.Categories), &out)
	return out
}
func (s serviceRow) publicValue() map[string]any {
	return map[string]any{"id": s.ID, "name": s.Name, "slug": s.Slug, "summary": s.Summary, "image_url": nullableStr(s.ImageURL), "service_area": nullableStr(s.ServiceArea)}
}
func (s serviceRow) publicDetail() map[string]any {
	v := s.publicValue()
	v["description"] = s.Description.String
	cats := make([]any, 0)
	for _, c := range s.categories() {
		if c.IsActive {
			cats = append(cats, map[string]any{"name": c.Name, "slug": c.Slug})
		}
	}
	v["categories"] = cats
	return v
}
func (s serviceRow) adminValue() map[string]any {
	ids := make([]int64, 0)
	cats := make([]any, 0)
	for _, c := range s.categories() {
		ids = append(ids, c.ID)
		cats = append(cats, map[string]any{"id": c.ID, "name": c.Name, "slug": c.Slug, "parent_id": c.ParentID})
	}
	return map[string]any{"id": s.ID, "name": s.Name, "slug": s.Slug, "summary": s.Summary, "description": s.Description.String,
		"image_url": nullableStr(s.ImageURL), "service_area": nullableStr(s.ServiceArea), "sort_order": s.SortOrder, "is_active": s.IsActive,
		"category_ids": ids, "categories": cats, "created_at": nullableTime(s.CreatedAt), "updated_at": nullableTime(s.UpdatedAt)}
}
func readService(s *svc.ServiceContext, r *http.Request, id int64) (map[string]any, error) {
	var row serviceRow
	err := s.DB.QueryRowCtx(r.Context(), &row, `SELECT `+serviceColumns+` FROM crm_schema.services s WHERE s.id=$1`, id)
	return row.adminValue(), err
}

// GET /services?category=<slug>&limit=n: jasa aktif yang terhubung ke kategori itu atau induknya.
func publicServicesHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		category := strings.TrimSpace(r.URL.Query().Get("category"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 0 || limit > 100 {
			limit = 100
		}
		var rows []serviceRow
		err := s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT `+serviceColumns+` FROM crm_schema.services s
			WHERE s.is_active=TRUE AND ($1='' OR EXISTS(SELECT 1 FROM crm_schema.service_categories sc JOIN crm_schema.categories c ON c.slug=$1
				WHERE sc.service_id=s.id AND (sc.category_id=c.id OR sc.category_id=c.parent_id)))
			ORDER BY s.sort_order,s.id LIMIT NULLIF($2,0)`, category, limit)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		items := make([]any, 0, len(rows))
		for _, row := range rows {
			items = append(items, row.publicValue())
		}
		apiData(w, r, items)
	}
}
func publicServiceDetailHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/api/v1/services/")
		var row serviceRow
		err := s.DB.QueryRowCtx(r.Context(), &row, `SELECT `+serviceColumns+` FROM crm_schema.services s WHERE s.slug=$1 AND s.is_active=TRUE`, slug)
		if isRowMissing(err) {
			apiFail(w, r, 404, "Jasa tidak ditemukan.")
			return
		}
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, row.publicDetail())
	}
}
func adminServicesHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rows []serviceRow
		if err := s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT `+serviceColumns+` FROM crm_schema.services s ORDER BY s.sort_order,s.id`); err != nil {
			apiDBFail(w, r, err)
			return
		}
		items := make([]any, 0, len(rows))
		for _, row := range rows {
			items = append(items, row.adminValue())
		}
		apiData(w, r, items)
	}
}
func adminServiceDetailHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := apiID(w, r)
		if !ok {
			return
		}
		v, err := readService(s, r, id)
		if isRowMissing(err) {
			apiFail(w, r, 404, "Jasa tidak ditemukan.")
			return
		}
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, v)
	}
}

type serviceInput struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Summary     string  `json:"summary"`
	Description string  `json:"description"`
	ImageURL    *string `json:"image_url"`
	ServiceArea *string `json:"service_area"`
	CategoryIDs []int64 `json:"category_ids"`
	IsActive    *bool   `json:"is_active"`
}

// validateService memeriksa isian per kolom; pesan mengikuti mock frontend.
func validateService(p *serviceInput) map[string][]string {
	errs := map[string][]string{}
	p.Name = strings.TrimSpace(p.Name)
	p.Summary = strings.TrimSpace(p.Summary)
	p.Description = strings.TrimSpace(p.Description)
	switch {
	case p.Name == "":
		errs["name"] = []string{"Nama jasa wajib diisi"}
	case utf8.RuneCountInString(p.Name) > 100:
		errs["name"] = []string{"Nama jasa maksimal 100 karakter"}
	}
	switch {
	case p.Slug == "":
		errs["slug"] = []string{"Slug wajib diisi"}
	case len(p.Slug) > 150 || !articleSlug.MatchString(p.Slug):
		errs["slug"] = []string{"Slug hanya boleh huruf kecil, angka, dan tanda hubung"}
	}
	switch {
	case p.Summary == "":
		errs["summary"] = []string{"Ringkasan wajib diisi"}
	case utf8.RuneCountInString(p.Summary) > 160:
		errs["summary"] = []string{"Ringkasan maksimal 160 karakter"}
	}
	if p.ServiceArea != nil && utf8.RuneCountInString(strings.TrimSpace(*p.ServiceArea)) > 100 {
		errs["service_area"] = []string{"Area layanan maksimal 100 karakter"}
	}
	if p.ImageURL != nil && len(strings.TrimSpace(*p.ImageURL)) > 2048 {
		errs["image_url"] = []string{"Alamat gambar terlalu panjang"}
	}
	seen := map[int64]bool{}
	ids := make([]int64, 0, len(p.CategoryIDs))
	for _, id := range p.CategoryIDs {
		if id <= 0 {
			errs["category_ids"] = []string{"Kategori tidak valid"}
		} else if !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	p.CategoryIDs = ids
	return errs
}
func joinIDs(ids []int64) string {
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = strconv.FormatInt(id, 10)
	}
	return strings.Join(parts, ",")
}
func apiFieldFail(w http.ResponseWriter, r *http.Request, status int, msg string, errs map[string][]string) {
	httpx.WriteJsonCtx(r.Context(), w, status, map[string]any{"message": msg, "errors": errs})
}
func optionalText(v *string) any {
	if v == nil {
		return nil
	}
	return nilIfEmpty(*v)
}

// POST /admin/services (create=true) dan PUT /admin/services/:id.
func serviceSaveHandler(s *svc.ServiceContext, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int64
		if !create {
			var ok bool
			if id, ok = apiID(w, r); !ok {
				return
			}
		}
		var input serviceInput
		if !apiDecode(w, r, &input) {
			return
		}
		if errs := validateService(&input); len(errs) > 0 {
			apiFieldFail(w, r, 422, "Periksa kembali isian jasa.", errs)
			return
		}
		ids := joinIDs(input.CategoryIDs)
		var found int64
		if err := s.DB.QueryRowCtx(r.Context(), &found, `SELECT COUNT(*) FROM crm_schema.categories WHERE id = ANY(string_to_array($1,',')::bigint[])`, ids); err != nil {
			apiDBFail(w, r, err)
			return
		}
		if found != int64(len(input.CategoryIDs)) {
			apiFieldFail(w, r, 422, "Periksa kembali isian jasa.", map[string][]string{"category_ids": {"Kategori tidak valid"}})
			return
		}
		var taken bool
		if err := s.DB.QueryRowCtx(r.Context(), &taken, `SELECT EXISTS(SELECT 1 FROM crm_schema.services WHERE slug=$1 AND id<>$2)`, input.Slug, id); err != nil {
			apiDBFail(w, r, err)
			return
		}
		if taken {
			apiFieldFail(w, r, 409, "Slug sudah dipakai jasa lain.", map[string][]string{"slug": {"Slug sudah dipakai jasa lain"}})
			return
		}
		active := input.IsActive == nil || *input.IsActive
		var before serviceRow
		if !create {
			err := s.DB.QueryRowCtx(r.Context(), &before, `SELECT `+serviceColumns+` FROM crm_schema.services s WHERE s.id=$1`, id)
			if isRowMissing(err) {
				apiFail(w, r, 404, "Jasa tidak ditemukan.")
				return
			}
			if err != nil {
				apiDBFail(w, r, err)
				return
			}
		}
		err := s.DB.TransactCtx(r.Context(), func(ctx context.Context, tx sqlx.Session) error {
			var e error
			if create {
				// Jasa baru masuk urutan terakhir.
				e = tx.QueryRowCtx(ctx, &id, `INSERT INTO crm_schema.services(name,slug,summary,description,image_url,service_area,is_active,sort_order)
					VALUES($1,$2,$3,$4,$5,$6,$7,(SELECT COALESCE(MAX(sort_order),0)+1 FROM crm_schema.services)) RETURNING id`,
					input.Name, input.Slug, input.Summary, input.Description, optionalText(input.ImageURL), optionalText(input.ServiceArea), active)
			} else {
				_, e = tx.ExecCtx(ctx, `UPDATE crm_schema.services SET name=$2,slug=$3,summary=$4,description=$5,image_url=$6,service_area=$7,is_active=$8 WHERE id=$1`,
					id, input.Name, input.Slug, input.Summary, input.Description, optionalText(input.ImageURL), optionalText(input.ServiceArea), active)
			}
			if e != nil {
				return e
			}
			if _, e = tx.ExecCtx(ctx, `DELETE FROM crm_schema.service_categories WHERE service_id=$1`, id); e != nil {
				return e
			}
			_, e = tx.ExecCtx(ctx, `INSERT INTO crm_schema.service_categories(service_id,category_id)
				SELECT $1::bigint, unnest(string_to_array($2,',')::bigint[]) ON CONFLICT DO NOTHING`, id, ids)
			return e
		})
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		result, err := readService(s, r, id)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		if create {
			auditAdminAction(s, r, "service.create", "Jasa", map[string]any{"id": id, "name": input.Name, "category_ids": input.CategoryIDs})
			httpx.WriteJsonCtx(r.Context(), w, http.StatusCreated, map[string]any{"data": result})
			return
		}
		auditAdminAction(s, r, "service.update", "Jasa", serviceChanges(before.adminValue(), result))
		apiData(w, r, result)
	}
}

// PATCH /admin/services/:id: { is_active } atau { sort_order }. Tidak ada DELETE; jasa dinonaktifkan.
func servicePatchHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := apiID(w, r)
		if !ok {
			return
		}
		var change struct {
			IsActive  *bool `json:"is_active"`
			SortOrder *int  `json:"sort_order"`
		}
		if !apiDecode(w, r, &change) {
			return
		}
		if change.IsActive == nil && change.SortOrder == nil {
			apiFail(w, r, 422, "Isi is_active atau sort_order.")
			return
		}
		before, err := readService(s, r, id)
		if isRowMissing(err) {
			apiFail(w, r, 404, "Jasa tidak ditemukan.")
			return
		}
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		if _, err = s.DB.ExecCtx(r.Context(), `UPDATE crm_schema.services SET is_active=COALESCE($2,is_active),sort_order=COALESCE($3,sort_order) WHERE id=$1`, id, change.IsActive, change.SortOrder); err != nil {
			apiDBFail(w, r, err)
			return
		}
		item, err := readService(s, r, id)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		auditAdminAction(s, r, "service.update", "Jasa", serviceChanges(before, item))
		apiData(w, r, item)
	}
}

// serviceChanges mencatat kolom yang berubah, dengan nama jasa sebagai penanda (seperti mock frontend).
func serviceChanges(before, after map[string]any) map[string]any {
	changed := map[string]any{"id": after["id"], "name": after["name"]}
	for _, key := range []string{"name", "slug", "summary", "description", "image_url", "service_area", "sort_order", "is_active", "category_ids"} {
		b, _ := json.Marshal(before[key])
		a, _ := json.Marshal(after[key])
		if string(b) != string(a) {
			changed[key] = map[string]any{"before": before[key], "after": after[key]}
		}
	}
	return changed
}
