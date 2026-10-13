package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"webtranshome/internal/middleware"
	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Helpers for API contracts currently used by the Vue administration screens.
func apiOK(w http.ResponseWriter, r *http.Request, data any) { httpx.OkJsonCtx(r.Context(), w, data) }
func apiData(w http.ResponseWriter, r *http.Request, data any) {
	apiOK(w, r, map[string]any{"data": data})
}
func apiFail(w http.ResponseWriter, r *http.Request, status int, msg string) {
	httpx.WriteJsonCtx(r.Context(), w, status, map[string]any{"message": msg})
}
func apiDBFail(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, sqlx.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
		apiFail(w, r, 404, "Data tidak ditemukan.")
		return
	}
	logx.WithContext(r.Context()).Error(err)
	apiFail(w, r, 500, "Kesalahan database. Silakan coba lagi.")
}
func apiDecode(w http.ResponseWriter, r *http.Request, target any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		apiFail(w, r, 422, "Format data tidak valid: "+err.Error())
		return false
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		apiFail(w, r, 422, "Kirim satu objek JSON saja.")
		return false
	}
	return true
}
func apiID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	// Only called for paths with exactly one numeric :id parameter.
	raw := r.URL.Path
	idx := strings.LastIndexByte(raw, '/')
	if idx < 0 {
		apiFail(w, r, 400, "ID tidak valid")
		return 0, false
	}
	id, err := strconv.ParseInt(raw[idx+1:], 10, 64)
	if err != nil || id <= 0 {
		apiFail(w, r, 400, "ID tidak valid")
		return 0, false
	}
	return id, true
}
func pageParams(r *http.Request) (page, per, offset int) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	if page > 100000 {
		page = 100000
	}
	per, _ = strconv.Atoi(r.URL.Query().Get("per_page"))
	if per < 1 {
		per = 20
	}
	if per > 100 {
		per = 100
	}
	return page, per, (page - 1) * per
}
func pageMeta(page, per int, total int64) map[string]any {
	pages := int((total + int64(per) - 1) / int64(per))
	if pages < 1 {
		pages = 1
	}
	return map[string]any{"page": page, "per_page": per, "total": total, "total_pages": pages}
}
func nullableStr(s sql.NullString) any {
	if s.Valid {
		return s.String
	}
	return nil
}
func nullableTime(t sql.NullTime) any {
	if t.Valid {
		return t.Time.Format(time.RFC3339)
	}
	return nil
}
func nullableID(t sql.NullInt64) any {
	if t.Valid {
		return t.Int64
	}
	return nil
}
func nilIfEmpty(v string) any {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return v
}
func auditAdminAction(s *svc.ServiceContext, r *http.Request, action, module string, after any) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		return
	}
	b, _ := json.Marshal(after)
	_, err := s.DB.ExecCtx(r.Context(), `INSERT INTO crm_schema.activity_logs (admin_id,action,module,after_data,ip_address) VALUES($1,$2,$3,$4,$5)`, userID, action, module, string(b), strings.Split(r.RemoteAddr, ":")[0])
	if err != nil {
		logx.WithContext(r.Context()).Errorf("audit insert failed: %v", err)
	}
}

func parseIDPart(raw string) (int64, bool) {
	v, e := strconv.ParseInt(raw, 10, 64)
	return v, e == nil && v > 0
}
func rawJSON(s sql.NullString) any {
	if !s.Valid {
		return nil
	}
	var data any
	if json.Unmarshal([]byte(s.String), &data) != nil {
		return nil
	}
	return data
}

func isRowMissing(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || errors.Is(err, sqlx.ErrNotFound)
}