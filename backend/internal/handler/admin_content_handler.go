package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"webtranshome/internal/middleware"
	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// The database dates are local wall-clock timestamps, matching the YYYY-MM-DD
// fields of the existing Vue banner form; no timezone conversion is performed.
type bannerRow struct {
	ID        int64          `db:"id"`
	Title     string         `db:"title"`
	ImageURL  string         `db:"image_url"`
	LinkURL   sql.NullString `db:"link_url"`
	Position  string         `db:"position"`
	StartAt   sql.NullString `db:"start_at"`
	EndAt     sql.NullString `db:"end_at"`
	SortOrder int64          `db:"sort_order"`
	IsActive  bool           `db:"is_active"`
	Status    string         `db:"status"`
}

const bannerColumns = `id,title,image_url,link_url,position,
 to_char(start_at,'YYYY-MM-DD"T"HH24:MI:SS') AS start_at,
 to_char(end_at,'YYYY-MM-DD"T"HH24:MI:SS') AS end_at,sort_order,is_active,
 CASE WHEN NOT is_active THEN 'NONAKTIF' WHEN start_at > timezone('Asia/Jakarta',now()) THEN 'TERJADWAL'
 WHEN end_at < timezone('Asia/Jakarta',now()) THEN 'BERAKHIR' ELSE 'TAYANG' END AS status`

func (b bannerRow) value() map[string]any {
	return map[string]any{"id": b.ID, "title": b.Title, "image_url": b.ImageURL, "link_url": nullableStr(b.LinkURL), "position": b.Position, "start_at": nullableStr(b.StartAt), "end_at": nullableStr(b.EndAt), "start_date": datePart(b.StartAt), "end_date": datePart(b.EndAt), "sort_order": b.SortOrder, "is_active": b.IsActive, "status": b.Status}
}
func datePart(t sql.NullString) string {
	if !t.Valid || len(t.String) < 10 {
		return ""
	}
	return t.String[:10]
}
func readBanner(s *svc.ServiceContext, r *http.Request, id int64) (map[string]any, error) {
	var row bannerRow
	err := s.DB.QueryRowCtx(r.Context(), &row, `SELECT `+bannerColumns+` FROM crm_schema.banners WHERE id=$1`, id)
	return row.value(), err
}
func bannersHandler(s *svc.ServiceContext, admin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rows []bannerRow
		var err error
		if admin {
			err = s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT `+bannerColumns+` FROM crm_schema.banners ORDER BY CASE position WHEN 'HOME_SLIDER' THEN 0 ELSE 1 END,sort_order,id`)
		} else {
			err = s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT `+bannerColumns+` FROM crm_schema.banners WHERE is_active=TRUE AND (start_at IS NULL OR start_at<=timezone('Asia/Jakarta',now())) AND (end_at IS NULL OR end_at>=timezone('Asia/Jakarta',now())) AND ($1='' OR position=$1) ORDER BY CASE position WHEN 'HOME_SLIDER' THEN 0 ELSE 1 END,sort_order,id`, r.URL.Query().Get("position"))
		}
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		items := make([]any, 0, len(rows))
		for _, b := range rows {
			items = append(items, b.value())
		}
		apiData(w, r, items)
	}
}
func bannerDetailHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := apiID(w, r)
		if !ok {
			return
		}
		v, err := readBanner(s, r, id)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, v)
	}
}

type bannerInput struct {
	Title     string `json:"title"`
	ImageURL  string `json:"image_url"`
	LinkURL   string `json:"link_url"`
	Position  string `json:"position"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	IsActive  *bool  `json:"is_active"`
}

func bannerDates(p bannerInput) (any, any, bool) {
	start, err := time.Parse("2006-01-02", p.StartDate)
	if err != nil {
		return nil, nil, false
	}
	var endValue any
	if p.EndDate != "" {
		end, err := time.Parse("2006-01-02", p.EndDate)
		if err != nil || end.Before(start) {
			return nil, nil, false
		}
		endValue = end.Add(24*time.Hour - time.Second).Format("2006-01-02 15:04:05")
	}
	return start.Format("2006-01-02 15:04:05"), endValue, true
}
func validateBanner(p bannerInput) bool {
	if strings.TrimSpace(p.Title) == "" || strings.TrimSpace(p.ImageURL) == "" || (p.Position != "HOME_SLIDER" && p.Position != "HOME_SIDE") {
		return false
	}
	if p.LinkURL != "" && !strings.HasPrefix(p.LinkURL, "/") && !strings.HasPrefix(p.LinkURL, "https://") {
		return false
	}
	_, _, ok := bannerDates(p)
	return ok
}
func bannerSaveHandler(s *svc.ServiceContext, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input bannerInput
		if !apiDecode(w, r, &input) {
			return
		}
		if !validateBanner(input) {
			apiFail(w, r, 422, "Isi judul, gambar, posisi, tanggal, dan tautan banner dengan benar.")
			return
		}
		start, end, _ := bannerDates(input)
		active := true
		if input.IsActive != nil {
			active = *input.IsActive
		}
		var id int64
		var err error
		if create {
			err = s.DB.QueryRowCtx(r.Context(), &id, `INSERT INTO crm_schema.banners(title,image_url,link_url,position,start_at,end_at,sort_order,is_active) VALUES($1,$2,$3,$4,$5,$6,(SELECT COALESCE(MAX(sort_order),0)+1 FROM crm_schema.banners WHERE position=$4),$7) RETURNING id`, strings.TrimSpace(input.Title), strings.TrimSpace(input.ImageURL), nilIfEmpty(input.LinkURL), input.Position, start, end, active)
		} else {
			var ok bool
			id, ok = apiID(w, r)
			if !ok {
				return
			}
			err = s.DB.QueryRowCtx(r.Context(), &id, `UPDATE crm_schema.banners SET title=$2,image_url=$3,link_url=$4,position=$5,start_at=$6,end_at=$7,is_active=$8,updated_at=NOW() WHERE id=$1 RETURNING id`, id, strings.TrimSpace(input.Title), strings.TrimSpace(input.ImageURL), nilIfEmpty(input.LinkURL), input.Position, start, end, active)
		}
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		auditAdminAction(s, r, map[bool]string{true: "banner.create", false: "banner.update"}[create], "Banner", map[string]any{"id": id, "title": input.Title})
		result, err := readBanner(s, r, id)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, result)
	}
}
func bannerPatchHandler(s *svc.ServiceContext) http.HandlerFunc {
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
		var newID int64
		err := s.DB.QueryRowCtx(r.Context(), &newID, `UPDATE crm_schema.banners SET is_active=COALESCE($2,is_active),sort_order=COALESCE($3,sort_order) WHERE id=$1 RETURNING id`, id, change.IsActive, change.SortOrder)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		auditAdminAction(s, r, "banner.update", "Banner", map[string]any{"id": id})
		item, err := readBanner(s, r, id)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, item)
	}
}

var articleSlug = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type articleRow struct {
	ID           int64          `db:"id"`
	Type         string         `db:"type"`
	Title        string         `db:"title"`
	Slug         string         `db:"slug"`
	ThumbnailURL sql.NullString `db:"thumbnail_url"`
	Content      string         `db:"content"`
	Status       string         `db:"status"`
	PublishedAt  sql.NullTime   `db:"published_at"`
	UpdatedAt    sql.NullTime   `db:"updated_at"`
	AuthorName   sql.NullString `db:"author_name"`
}

const articleColumns = `a.id,a.type,a.title,a.slug,a.thumbnail_url,a.content,a.status,a.published_at,a.updated_at,mp.full_name AS author_name`
const articleFrom = ` FROM crm_schema.articles a LEFT JOIN crm_schema.member_profiles mp ON mp.user_id=a.author_id `

func (a articleRow) value(detail bool) map[string]any {
	v := map[string]any{"id": a.ID, "type": a.Type, "title": a.Title, "slug": a.Slug, "thumbnail_url": nullableStr(a.ThumbnailURL), "status": a.Status, "published_at": nullableTime(a.PublishedAt), "updated_at": nullableTime(a.UpdatedAt), "author_name": nullableStr(a.AuthorName)}
	if detail {
		v["content"] = a.Content
	}
	return v
}
func readArticle(s *svc.ServiceContext, r *http.Request, id int64) (map[string]any, error) {
	var a articleRow
	err := s.DB.QueryRowCtx(r.Context(), &a, `SELECT `+articleColumns+articleFrom+` WHERE a.id=$1`, id)
	return a.value(true), err
}
func articlesHandler(s *svc.ServiceContext, public bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, per, offset := pageParams(r)
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		typeFilter := r.URL.Query().Get("type")
		status := r.URL.Query().Get("status")
		if public {
			status = "PUBLISHED"
		}
		where := ` WHERE ($1='' OR a.title ILIKE '%'||$1||'%' OR a.slug ILIKE '%'||$1||'%') AND ($2='' OR a.type=$2) AND ($3='' OR a.status=$3) `
		if public && typeFilter == "" {
			where += ` AND a.type<>'HALAMAN' `
		}
		var count int64
		err := s.DB.QueryRowCtx(r.Context(), &count, `SELECT COUNT(*)`+articleFrom+where, q, typeFilter, status)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		var rows []articleRow
		err = s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT `+articleColumns+articleFrom+where+` ORDER BY a.updated_at DESC NULLS LAST,a.id DESC LIMIT $4 OFFSET $5`, q, typeFilter, status, per, offset)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		out := make([]any, 0, len(rows))
		for _, a := range rows {
			v := a.value(false)
			if public {
				excerpt := strings.TrimSpace(a.Content)
				if len(excerpt) > 160 {
					excerpt = string([]rune(excerpt)[:min(len([]rune(excerpt)), 160)])
				}
				v["excerpt"] = excerpt
			}
			out = append(out, v)
		}
		apiOK(w, r, map[string]any{"data": out, "meta": pageMeta(page, per, count)})
	}
}
func articleDetailHandler(s *svc.ServiceContext, public bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if public {
			slug := strings.TrimPrefix(r.URL.Path, "/api/v1/articles/")
			var a articleRow
			err := s.DB.QueryRowCtx(r.Context(), &a, `SELECT `+articleColumns+articleFrom+` WHERE a.slug=$1 AND a.status='PUBLISHED'`, slug)
			if err != nil {
				apiDBFail(w, r, err)
				return
			}
			apiData(w, r, a.value(true))
			return
		}
		id, ok := apiID(w, r)
		if !ok {
			return
		}
		v, err := readArticle(s, r, id)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, v)
	}
}

type articleInput struct {
	Type         string `json:"type"`
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	ThumbnailURL string `json:"thumbnail_url"`
	Content      string `json:"content"`
}

func validateArticle(v articleInput) bool {
	return (v.Type == "ARTIKEL" || v.Type == "PROMO" || v.Type == "EVENT" || v.Type == "HALAMAN") && strings.TrimSpace(v.Title) != "" && articleSlug.MatchString(v.Slug) && strings.TrimSpace(v.Content) != ""
}
func articleSaveHandler(s *svc.ServiceContext, create bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input articleInput
		if !apiDecode(w, r, &input) {
			return
		}
		if !validateArticle(input) {
			apiFail(w, r, 422, "Tipe, judul, slug, dan isi artikel wajib valid.")
			return
		}
		var id int64
		var err error
		if create {
			userID, _ := middleware.UserIDFromContext(r.Context())
			err = s.DB.QueryRowCtx(r.Context(), &id, `INSERT INTO crm_schema.articles(type,title,slug,thumbnail_url,content,status,author_id) VALUES($1,$2,$3,$4,$5,'DRAFT',$6) ON CONFLICT(slug) DO NOTHING RETURNING id`, input.Type, strings.TrimSpace(input.Title), input.Slug, nilIfEmpty(input.ThumbnailURL), input.Content, userID)
		} else {
			id, _ = apiID(w, r)
			if id == 0 {
				return
			}
			err = s.DB.QueryRowCtx(r.Context(), &id, `UPDATE crm_schema.articles SET type=$2,title=$3,slug=$4,thumbnail_url=$5,content=$6 WHERE id=$1 AND NOT EXISTS(SELECT 1 FROM crm_schema.articles WHERE slug=$4 AND id<>$1) RETURNING id`, id, input.Type, strings.TrimSpace(input.Title), input.Slug, nilIfEmpty(input.ThumbnailURL), input.Content)
		}
		if errors.Is(err, sqlx.ErrNotFound) {
			apiFail(w, r, 409, "Artikel tidak ditemukan atau slug sudah digunakan.")
			return
		}
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		auditAdminAction(s, r, map[bool]string{true: "article.create", false: "article.update"}[create], "Artikel", map[string]any{"id": id, "title": input.Title})
		v, err := readArticle(s, r, id)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, v)
	}
}
func articleStatusHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := apiID(w, r)
		if !ok {
			return
		}
		var in struct {
			Status string `json:"status"`
		}
		if !apiDecode(w, r, &in) {
			return
		}
		if in.Status != "DRAFT" && in.Status != "PUBLISHED" {
			apiFail(w, r, 422, "Status harus DRAFT atau PUBLISHED.")
			return
		}
		var ret int64
		err := s.DB.QueryRowCtx(r.Context(), &ret, `UPDATE crm_schema.articles SET status=$2,published_at=CASE WHEN $2='PUBLISHED' THEN COALESCE(published_at,NOW()) ELSE NULL END WHERE id=$1 RETURNING id`, id, in.Status)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		auditAdminAction(s, r, "article.publish", "Artikel", map[string]any{"id": id, "status": in.Status})
		v, err := readArticle(s, r, id)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, v)
	}
}