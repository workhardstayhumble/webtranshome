package handler

import (
	"database/sql"
	"net/http"
	"strings"
	"time"

	"webtranshome/internal/svc"
)

type contactRow struct {
	ID              int64          `db:"id"`
	UserID          int64          `db:"user_id"`
	Name            string         `db:"name"`
	Email           string         `db:"email"`
	PhoneNumber     string         `db:"phone_number"`
	MemberCode      sql.NullString `db:"member_code"`
	QontakContactID sql.NullString `db:"qontak_contact_id"`
	SyncStatus      string         `db:"sync_status"`
	RetryCount      int64          `db:"retry_count"`
	StatusCode      sql.NullInt64  `db:"status_code"`
	LastSyncedAt    sql.NullTime   `db:"last_synced_at"`
	CreatedAt       sql.NullTime   `db:"created_at"`
	UpdatedAt       sql.NullTime   `db:"updated_at"`
}

func (c contactRow) value() map[string]any {
	return map[string]any{"id": c.ID, "user_id": c.UserID, "name": c.Name, "email": c.Email, "phone_number": c.PhoneNumber, "member_code": nullableStr(c.MemberCode), "qontak_contact_id": nullableStr(c.QontakContactID), "sync_status": c.SyncStatus, "retry_count": c.RetryCount, "status_code": nullableID(c.StatusCode), "last_synced_at": nullableTime(c.LastSyncedAt), "created_at": nullableTime(c.CreatedAt), "updated_at": nullableTime(c.UpdatedAt)}
}

const contactQuery = ` FROM crm_schema.qontak_contacts c JOIN crm_schema.users u ON u.id=c.user_id LEFT JOIN crm_schema.member_profiles mp ON mp.user_id=u.id LEFT JOIN LATERAL (SELECT sl.status_code,sl.retry_count FROM crm_schema.sync_logs sl WHERE sl.user_id=c.user_id ORDER BY sl.created_at DESC,sl.id DESC LIMIT 1) sl ON TRUE `
const contactColumns = `c.id,c.user_id,COALESCE(mp.full_name,u.email) AS name,u.email,u.phone_number,mp.member_code,c.qontak_contact_id,c.sync_status,COALESCE(sl.retry_count,0) AS retry_count,sl.status_code,c.last_synced_at,c.created_at,c.updated_at`

func readContact(s *svc.ServiceContext, r *http.Request, id int64) (map[string]any, error) {
	var row contactRow
	err := s.DB.QueryRowCtx(r.Context(), &row, `SELECT `+contactColumns+contactQuery+` WHERE c.id=$1`, id)
	return row.value(), err
}
func qontakListHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, per, offset := pageParams(r)
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		status := r.URL.Query().Get("status")
		where := ` WHERE ($1='' OR u.email ILIKE '%'||$1||'%' OR u.phone_number ILIKE '%'||$1||'%' OR mp.full_name ILIKE '%'||$1||'%' OR mp.member_code ILIKE '%'||$1||'%') `
		var counts []struct {
			Status string `db:"status"`
			Count  int64  `db:"count"`
		}
		err := s.DB.QueryRowsCtx(r.Context(), &counts, `SELECT c.sync_status AS status,COUNT(*) AS count`+contactQuery+where+` GROUP BY c.sync_status`, q)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		countsMap := map[string]int64{"ALL": 0, "PENDING": 0, "SYNCED": 0, "FAILED": 0}
		for _, x := range counts {
			countsMap[x.Status] = x.Count
			countsMap["ALL"] += x.Count
		}
		var rows []contactRow
		err = s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT `+contactColumns+contactQuery+where+` AND ($2='' OR c.sync_status=$2) ORDER BY c.created_at DESC,c.id DESC LIMIT $3 OFFSET $4`, q, status, per, offset)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		total := countsMap["ALL"]
		if status != "" {
			total = countsMap[status]
		}
		items := make([]any, 0, len(rows))
		for _, c := range rows {
			items = append(items, c.value())
		}
		meta := pageMeta(page, per, total)
		meta["counts"] = countsMap
		apiOK(w, r, map[string]any{"data": items, "meta": meta})
	}
}
func qontakDetailHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := apiID(w, r)
		if !ok {
			return
		}
		v, err := readContact(s, r, id)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, v)
	}
}

type contactLogRow struct {
	ID              int64          `db:"id"`
	EventType       string         `db:"event_type"`
	StatusCode      sql.NullInt64  `db:"status_code"`
	RetryCount      int64          `db:"retry_count"`
	RequestPayload  sql.NullString `db:"request_payload"`
	ResponsePayload sql.NullString `db:"response_payload"`
	CreatedAt       time.Time      `db:"created_at"`
}

func qontakLogsHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// /admin/qontak/contacts/:id/logs => second-last path segment is ID.
		pieces := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pieces) < 2 {
			apiFail(w, r, 400, "ID tidak valid")
			return
		}
		id, ok := parseIDPart(pieces[len(pieces)-2])
		if !ok {
			apiFail(w, r, 400, "ID tidak valid")
			return
		}
		var user struct {
			ID int64 `db:"user_id"`
		}
		if err := s.DB.QueryRowCtx(r.Context(), &user, `SELECT user_id FROM crm_schema.qontak_contacts WHERE id=$1`, id); err != nil {
			apiDBFail(w, r, err)
			return
		}
		var rows []contactLogRow
		err := s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT id,event_type,status_code,retry_count,request_payload::text AS request_payload,response_payload::text AS response_payload,created_at FROM crm_schema.sync_logs WHERE user_id=$1 ORDER BY created_at DESC,id DESC LIMIT 100`, user.ID)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		out := make([]any, 0, len(rows))
		for _, x := range rows {
			out = append(out, map[string]any{"id": x.ID, "event_type": x.EventType, "status_code": nullableID(x.StatusCode), "retry_count": x.RetryCount, "request_payload": rawJSON(x.RequestPayload), "response_payload": rawJSON(x.ResponsePayload), "created_at": x.CreatedAt.Format(time.RFC3339)})
		}
		apiData(w, r, out)
	}
}
func qontakRetryUnavailableHandler(w http.ResponseWriter, r *http.Request) {
	apiFail(w, r, http.StatusConflict, "Retry Qontak belum tersedia: worker dan koneksi Mekari Qontak belum diimplementasikan. Tidak ada kontak yang diubah.")
}