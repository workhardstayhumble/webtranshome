package handler

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"webtranshome/internal/middleware"

	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"webtranshome/internal/svc"
)

type memberRow struct {
	ID                   int64          `db:"id"`
	ProfileID            int64          `db:"profile_id"`
	MemberCode           string         `db:"member_code"`
	FullName             string         `db:"full_name"`
	Email                string         `db:"email"`
	PhoneNumber          string         `db:"phone_number"`
	Tier                 string         `db:"tier"`
	IsActive             bool           `db:"is_active"`
	EmailVerified        bool           `db:"email_verified"`
	CommunicationConsent bool           `db:"communication_consent"`
	Balance              int64          `db:"balance"`
	SyncStatus           sql.NullString `db:"sync_status"`
	ConsentAt            sql.NullTime   `db:"consent_at"`
	GoogleConnected      bool           `db:"google_connected"`
	LastLoginAt          sql.NullTime   `db:"last_login_at"`
	CreatedAt            sql.NullTime   `db:"created_at"`
	ContactID            sql.NullInt64  `db:"contact_id"`
}

const memberFrom = ` FROM crm_schema.member_profiles mp JOIN crm_schema.users u ON u.id=mp.user_id LEFT JOIN crm_schema.qontak_contacts qc ON qc.user_id=u.id LEFT JOIN (SELECT member_id,SUM(points)::bigint AS balance FROM crm_schema.point_transactions GROUP BY member_id) pt ON pt.member_id=mp.id `
const memberColumns = `u.id,mp.id AS profile_id,mp.member_code,mp.full_name,u.email,u.phone_number,mp.tier,u.is_active,(u.email_verified_at IS NOT NULL) AS email_verified,mp.communication_consent,COALESCE(pt.balance,0) AS balance,qc.sync_status,mp.consent_at,(u.google_id IS NOT NULL) AS google_connected,u.last_login_at,mp.created_at,qc.id AS contact_id`

func (m memberRow) value(detail bool) map[string]any {
	out := map[string]any{"id": m.ID, "member_code": m.MemberCode, "full_name": m.FullName, "email": m.Email, "phone_number": m.PhoneNumber, "tier": m.Tier, "is_active": m.IsActive, "email_verified": m.EmailVerified, "communication_consent": m.CommunicationConsent, "balance": m.Balance, "sync_status": nullableStr(m.SyncStatus), "created_at": nullableTime(m.CreatedAt)}
	if detail {
		out["consent_at"] = nullableTime(m.ConsentAt)
		out["google_connected"] = m.GoogleConnected
		out["last_login_at"] = nullableTime(m.LastLoginAt)
		if m.ContactID.Valid {
			out["qontak"] = map[string]any{"id": m.ContactID.Int64, "user_id": m.ID, "name": m.FullName, "sync_status": nullableStr(m.SyncStatus)}
		} else {
			out["qontak"] = nil
		}
	}
	return out
}
func memberWhere() string {
	return ` WHERE ($1='' OR u.email ILIKE '%'||$1||'%' OR u.phone_number ILIKE '%'||$1||'%' OR mp.full_name ILIKE '%'||$1||'%' OR mp.member_code ILIKE '%'||$1||'%') AND ($2='' OR mp.communication_consent=($2='1')) AND ($3='' OR (u.phone_number<>'')=($3='1')) AND ($4='' OR qc.sync_status=$4) AND ($5='' OR u.is_active=($5='1')) `
}
func memberFilters(r *http.Request) []any {
	q := r.URL.Query()
	return []any{strings.TrimSpace(q.Get("q")), q.Get("consent"), q.Get("hp"), q.Get("sync"), q.Get("active")}
}
func memberListHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, per, offset := pageParams(r)
		args := memberFilters(r)
		where := memberWhere()
		var total int64
		err := s.DB.QueryRowCtx(r.Context(), &total, `SELECT COUNT(*)`+memberFrom+where, args...)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		var rows []memberRow
		err = s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT `+memberColumns+memberFrom+where+` ORDER BY mp.created_at DESC,mp.id DESC LIMIT $6 OFFSET $7`, append(args, per, offset)...)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		data := make([]any, 0, len(rows))
		for _, x := range rows {
			data = append(data, x.value(false))
		}
		apiOK(w, r, map[string]any{"data": data, "meta": pageMeta(page, per, total)})
	}
}
func memberSearchHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		if len([]rune(q)) < 3 {
			apiFail(w, r, 422, "Pencarian minimal 3 karakter.")
			return
		}
		var rows []memberRow
		err := s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT `+memberColumns+memberFrom+` WHERE u.email ILIKE '%'||$1||'%' OR u.phone_number LIKE '%'||$1||'%' OR mp.full_name ILIKE '%'||$1||'%' OR mp.member_code ILIKE '%'||$1||'%' ORDER BY CASE WHEN lower(mp.member_code)=lower($1) OR lower(u.email)=lower($1) OR u.phone_number=$1 THEN 0 ELSE 1 END,mp.full_name LIMIT 10`, q)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		data := make([]any, 0, len(rows))
		for _, x := range rows {
			data = append(data, x.value(false))
		}
		apiData(w, r, data)
	}
}
func readMember(s *svc.ServiceContext, r *http.Request, id int64) (map[string]any, error) {
	var m memberRow
	err := s.DB.QueryRowCtx(r.Context(), &m, `SELECT `+memberColumns+memberFrom+` WHERE u.id=$1`, id)
	return m.value(true), err
}
func memberDetailHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := apiID(w, r)
		if !ok {
			return
		}
		v, err := readMember(s, r, id)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, v)
	}
}
func memberPatchHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := apiID(w, r)
		if !ok {
			return
		}
		var input struct {
			IsActive *bool `json:"is_active"`
		}
		if !apiDecode(w, r, &input) {
			return
		}
		if input.IsActive == nil {
			apiFail(w, r, 422, "is_active wajib boolean.")
			return
		}
		var userID int64
		err := s.DB.QueryRowCtx(r.Context(), &userID, `UPDATE crm_schema.users SET is_active=$2 WHERE id=$1 AND EXISTS(SELECT 1 FROM crm_schema.member_profiles WHERE user_id=$1) RETURNING id`, id, *input.IsActive)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		auditAdminAction(s, r, "member.update", "Member", map[string]any{"id": id, "is_active": *input.IsActive})
		v, err := readMember(s, r, id)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, v)
	}
}

type pointRow struct {
	ID             int64          `db:"id"`
	UserID         int64          `db:"user_id"`
	MemberCode     string         `db:"member_code"`
	Type           string         `db:"type"`
	Points         int64          `db:"points"`
	ReferenceNo    sql.NullString `db:"reference_no"`
	PurchaseAmount sql.NullString `db:"purchase_amount"`
	Note           sql.NullString `db:"note"`
	CreatedByName  sql.NullString `db:"created_by_name"`
	CreatedAt      time.Time      `db:"created_at"`
}

func (p pointRow) value() map[string]any {
	return map[string]any{"id": p.ID, "user_id": p.UserID, "member_code": p.MemberCode, "type": p.Type, "points": p.Points, "reference_no": nullableStr(p.ReferenceNo), "purchase_amount": nullableStr(p.PurchaseAmount), "note": nullableStr(p.Note), "created_by_name": nullableStr(p.CreatedByName), "created_at": p.CreatedAt.Format(time.RFC3339)}
}
func pointsListHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := parseIDPart(r.URL.Query().Get("user_id"))
		if !ok {
			apiFail(w, r, 422, "user_id wajib diisi.")
			return
		}
		page, per, offset := pageParams(r)
		typeFilter := r.URL.Query().Get("type")
		var member struct {
			ID int64 `db:"id"`
		}
		if err := s.DB.QueryRowCtx(r.Context(), &member, `SELECT id FROM crm_schema.member_profiles WHERE user_id=$1`, userID); err != nil {
			apiDBFail(w, r, err)
			return
		}
		var total, balance int64
		if err := s.DB.QueryRowCtx(r.Context(), &total, `SELECT COUNT(*) FROM crm_schema.point_transactions WHERE member_id=$1 AND ($2='' OR type=$2)`, member.ID, typeFilter); err != nil {
			apiDBFail(w, r, err)
			return
		}
		if err := s.DB.QueryRowCtx(r.Context(), &balance, `SELECT COALESCE(SUM(points),0)::bigint FROM crm_schema.point_transactions WHERE member_id=$1`, member.ID); err != nil {
			apiDBFail(w, r, err)
			return
		}
		var rows []pointRow
		err := s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT pt.id,mp.user_id,mp.member_code,pt.type,pt.points,pt.reference_no,pt.purchase_amount::text AS purchase_amount,pt.note,COALESCE(adminp.full_name,admin.email) AS created_by_name,pt.created_at FROM crm_schema.point_transactions pt JOIN crm_schema.member_profiles mp ON mp.id=pt.member_id LEFT JOIN crm_schema.users admin ON admin.id=pt.created_by LEFT JOIN crm_schema.member_profiles adminp ON adminp.user_id=admin.id WHERE pt.member_id=$1 AND ($2='' OR pt.type=$2) ORDER BY pt.created_at DESC,pt.id DESC LIMIT $3 OFFSET $4`, member.ID, typeFilter, per, offset)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		data := make([]any, 0, len(rows))
		for _, x := range rows {
			data = append(data, x.value())
		}
		meta := pageMeta(page, per, total)
		meta["balance"] = balance
		apiOK(w, r, map[string]any{"data": data, "meta": meta})
	}
}
func pointsCheckReferenceHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reference := strings.TrimSpace(r.URL.Query().Get("ref"))
		if reference == "" {
			apiFail(w, r, 422, "Nomor nota wajib diisi.")
			return
		}
		var row struct {
			CreatedAt  time.Time `db:"created_at"`
			MemberCode string    `db:"member_code"`
		}
		err := s.DB.QueryRowCtx(r.Context(), &row, `SELECT pt.created_at,mp.member_code FROM crm_schema.point_transactions pt JOIN crm_schema.member_profiles mp ON mp.id=pt.member_id WHERE pt.type='EARN' AND lower(pt.reference_no)=lower($1) LIMIT 1`, reference)
		if err == nil {
			apiData(w, r, map[string]any{"available": false, "used_at": row.CreatedAt.Format(time.RFC3339), "member_code": row.MemberCode})
			return
		}
		if err == sql.ErrNoRows {
			apiData(w, r, map[string]any{"available": true})
			return
		}
		// sqlx uses a distinct not-found sentinel.
		if isRowMissing(err) {
			apiData(w, r, map[string]any{"available": true})
			return
		}
		apiDBFail(w, r, err)
	}
}

// Real point transactions use a member-row lock so simultaneous redemptions
// cannot overdraw the same account. The ratio must be configured explicitly.
type pointInput struct {
	UserID         int64   `json:"user_id"`
	Type           string  `json:"type"`
	ReferenceNo    string  `json:"reference_no"`
	PurchaseAmount float64 `json:"purchase_amount"`
	Points         int64   `json:"points"`
	Note           string  `json:"note"`
}

func createPointHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in pointInput
		if !apiDecode(w, r, &in) {
			return
		}
		if in.UserID <= 0 || (in.Type != "EARN" && in.Type != "REDEEM" && in.Type != "ADJUST") {
			apiFail(w, r, 422, "Member dan jenis transaksi wajib diisi.")
			return
		}
		if in.Type == "ADJUST" {
			apiFail(w, r, 403, "Penyesuaian saldo belum dibuka melalui API ini; memerlukan alur persetujuan khusus.")
			return
		}
		if in.Type == "REDEEM" {
			adminID, ok := middleware.UserIDFromContext(r.Context())
			if !ok {
				apiFail(w, r, 401, "Sesi tidak valid.")
				return
			}
			granted, err := s.CheckPermission.Execute(r.Context(), adminID, "membership.points.redeem")
			if err != nil {
				apiDBFail(w, r, err)
				return
			}
			if !granted {
				apiFail(w, r, 403, "Tidak memiliki izin penukaran poin.")
				return
			}
		}
		if in.Type == "EARN" && (strings.TrimSpace(in.ReferenceNo) == "" || in.PurchaseAmount <= 0) {
			apiFail(w, r, 422, "Nomor nota dan nilai pembelian wajib diisi.")
			return
		}
		if in.Type == "REDEEM" && (in.Points <= 0 || strings.TrimSpace(in.Note) == "") {
			apiFail(w, r, 422, "Jumlah poin dan alasan penukaran wajib diisi.")
			return
		}
		if in.Type == "EARN" {
			var ratio struct {
				Value string `db:"setting_value"`
			}
			err := s.DB.QueryRowCtx(r.Context(), &ratio, `SELECT setting_value FROM crm_schema.site_settings WHERE setting_key='point_ratio_rupiah'`)
			if err != nil {
				apiFail(w, r, 422, "Rasio poin belum diatur. Hubungi Super Admin sebelum input nota.")
				return
			}
			val, e := strconv.ParseFloat(ratio.Value, 64)
			if e != nil || val < 1000 {
				apiFail(w, r, 422, "Konfigurasi rasio poin tidak valid.")
				return
			}
			in.Points = int64(math.Floor(in.PurchaseAmount / val))
			if in.Points < 1 {
				apiFail(w, r, 422, "Total pembelian belum mencapai batas minimal poin.")
				return
			}
		}
		if in.Type == "REDEEM" {
			in.Points = -in.Points
		}
		adminID, _ := middleware.UserIDFromContext(r.Context())
		var transactionID, balance int64
		var failureStatus int
		var failureMessage string
		err := s.DB.TransactCtx(r.Context(), func(ctx context.Context, tx sqlx.Session) error {
			var member struct {
				ID     int64 `db:"id"`
				Active bool  `db:"active"`
			}
			if e := tx.QueryRowCtx(ctx, &member, `SELECT mp.id,u.is_active AS active FROM crm_schema.member_profiles mp JOIN crm_schema.users u ON u.id=mp.user_id WHERE mp.user_id=$1 FOR UPDATE OF mp`, in.UserID); e != nil {
				return e
			}
			if !member.Active {
				failureStatus = 422
				failureMessage = "Akun member nonaktif."
				return errPointRule
			}
			if e := tx.QueryRowCtx(ctx, &balance, `SELECT COALESCE(SUM(points),0)::bigint FROM crm_schema.point_transactions WHERE member_id=$1`, member.ID); e != nil {
				return e
			}
			if balance+in.Points < 0 {
				failureStatus = 422
				failureMessage = "Saldo poin tidak mencukupi."
				return errPointRule
			}
			var reference any
			if in.Type == "EARN" {
				ref := strings.ToUpper(strings.TrimSpace(in.ReferenceNo))
				reference = ref
				var used bool
				if e := tx.QueryRowCtx(ctx, &used, `SELECT EXISTS(SELECT 1 FROM crm_schema.point_transactions WHERE type='EARN' AND upper(reference_no)=$1)`, ref); e != nil {
					return e
				}
				if used {
					failureStatus = 409
					failureMessage = "Nomor nota sudah pernah digunakan."
					return errPointRule
				}
			}
			var amount any
			if in.Type == "EARN" {
				amount = in.PurchaseAmount
			}
			return tx.QueryRowCtx(ctx, &transactionID, `INSERT INTO crm_schema.point_transactions(member_id,type,points,reference_no,purchase_amount,note,created_by) VALUES($1,$2,$3,$4,$5,$6,$7) ON CONFLICT (type,reference_no) DO NOTHING RETURNING id`, member.ID, in.Type, in.Points, reference, amount, nilIfEmpty(in.Note), adminID)
		})
		if errors.Is(err, errPointRule) {
			apiFail(w, r, failureStatus, failureMessage)
			return
		}
		if in.Type == "EARN" && isRowMissing(err) {
			apiFail(w, r, 409, "Nomor nota sudah pernah digunakan.")
			return
		}
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		auditAdminAction(s, r, "point.create", "Poin", map[string]any{"transaction_id": transactionID, "user_id": in.UserID, "type": in.Type, "points": in.Points})
		var row pointRow
		e := s.DB.QueryRowCtx(r.Context(), &row, `SELECT pt.id,mp.user_id,mp.member_code,pt.type,pt.points,pt.reference_no,pt.purchase_amount::text AS purchase_amount,pt.note,COALESCE(adminp.full_name,admin.email) AS created_by_name,pt.created_at FROM crm_schema.point_transactions pt JOIN crm_schema.member_profiles mp ON mp.id=pt.member_id LEFT JOIN crm_schema.users admin ON admin.id=pt.created_by LEFT JOIN crm_schema.member_profiles adminp ON adminp.user_id=admin.id WHERE pt.id=$1`, transactionID)
		if e != nil {
			apiDBFail(w, r, e)
			return
		}
		apiData(w, r, map[string]any{"transaction": row.value(), "balance": balance + in.Points})
	}
}

var errPointRule = errors.New("point validation failed")