package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"webtranshome/internal/svc"
)

type logRow struct {
	ID         int64          `db:"id"`
	AdminID    int64          `db:"admin_id"`
	AdminName  string         `db:"admin_name"`
	Action     string         `db:"action"`
	Module     string         `db:"module"`
	BeforeData sql.NullString `db:"before_data"`
	AfterData  sql.NullString `db:"after_data"`
	IPAddress  sql.NullString `db:"ip_address"`
	CreatedAt  time.Time      `db:"created_at"`
}

func (l logRow) value() map[string]any {
	var before, after any
	if l.BeforeData.Valid {
		var v any
		if json.Unmarshal([]byte(l.BeforeData.String), &v) == nil {
			before = v
		}
	}
	if l.AfterData.Valid {
		var v any
		if json.Unmarshal([]byte(l.AfterData.String), &v) == nil {
			after = v
		}
	}
	return map[string]any{"id": l.ID, "admin_id": l.AdminID, "admin_name": l.AdminName, "action": l.Action, "module": l.Module, "before_data": before, "after_data": after, "ip_address": nullableStr(l.IPAddress), "created_at": l.CreatedAt.Format(time.RFC3339)}
}
func adminAuditListHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, per, offset := pageParams(r)
		q := r.URL.Query()
		admin := strings.TrimSpace(q.Get("admin_id"))
		module := q.Get("module")
		start := q.Get("from")
		end := q.Get("to")
		for _, d := range []string{start, end} {
			if d != "" {
				if _, err := time.Parse("2006-01-02", d); err != nil {
					apiFail(w, r, 422, "Format tanggal harus YYYY-MM-DD.")
					return
				}
			}
		}
		where := ` WHERE ($1='' OR l.admin_id::text=$1) AND ($2='' OR l.module=$2) AND ($3='' OR l.created_at >= NULLIF($3,'')::date) AND ($4='' OR l.created_at < (NULLIF($4,'')::date + interval '1 day')) `
		var total int64
		err := s.DB.QueryRowCtx(r.Context(), &total, `SELECT COUNT(*) FROM crm_schema.activity_logs l`+where, admin, module, start, end)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		var rows []logRow
		err = s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT l.id,l.admin_id,COALESCE(mp.full_name,u.email) AS admin_name,l.action,l.module,l.before_data::text AS before_data,l.after_data::text AS after_data,l.ip_address,l.created_at FROM crm_schema.activity_logs l JOIN crm_schema.users u ON u.id=l.admin_id LEFT JOIN crm_schema.member_profiles mp ON mp.user_id=u.id `+where+` ORDER BY l.created_at DESC,l.id DESC LIMIT $5 OFFSET $6`, admin, module, start, end, per, offset)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		items := make([]any, 0, len(rows))
		for _, item := range rows {
			items = append(items, item.value())
		}
		apiOK(w, r, map[string]any{"data": items, "meta": pageMeta(page, per, total)})
	}
}

type auditOptionAdmin struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
type auditOptionModule struct {
	Module string `db:"module"`
}

func adminAuditOptionsHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var admins []auditOptionAdmin
		err := s.DB.QueryRowsCtx(r.Context(), &admins, `SELECT DISTINCT u.id,COALESCE(mp.full_name,u.email) AS name FROM crm_schema.users u JOIN crm_schema.user_roles ur ON ur.user_id=u.id JOIN crm_schema.roles role ON role.id=ur.role_id LEFT JOIN crm_schema.member_profiles mp ON mp.user_id=u.id WHERE role.code IN ('SUPER_ADMIN','ADMIN_KATALOG','ADMIN_MEMBERSHIP','ADMIN_KONTEN') ORDER BY name`)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		var modules []auditOptionModule
		err = s.DB.QueryRowsCtx(r.Context(), &modules, `SELECT DISTINCT module FROM crm_schema.activity_logs ORDER BY module`)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		aa := make([]any, 0, len(admins))
		for _, a := range admins {
			aa = append(aa, map[string]any{"id": a.ID, "name": a.Name})
		}
		mm := make([]string, 0, len(modules))
		for _, m := range modules {
			mm = append(mm, m.Module)
		}
		apiData(w, r, map[string]any{"admins": aa, "modules": mm})
	}
}