package handler

import (
	"context"
	"net/http"
	"time"

	"webtranshome/internal/middleware"
	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

// Semua angka dan log di dashboard difilter menurut izin dari database.
// Middleware ini MEMBATASI akses ke akun admin aktif dan setidaknya satu izin admin.
// Tidak memakai bypass berdasarkan nama SUPER_ADMIN.
type dashboardPermissionsKey struct{}

func dashboardAccessMiddleware(s *svc.ServiceContext) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			userID, ok := middleware.UserIDFromContext(r.Context())
			if !ok {
				apiFail(w, r, 401, "Login admin diperlukan.")
				return
			}
			var rows []struct {
				Code string `db:"code"`
			}
			err := s.DB.QueryRowsCtx(r.Context(), &rows, `SELECT DISTINCT p.code
                FROM crm_schema.users u
                JOIN crm_schema.user_roles ur ON ur.user_id=u.id
                JOIN crm_schema.roles ro ON ro.id=ur.role_id
                JOIN crm_schema.role_permissions rp ON rp.role_id=ro.id
                JOIN crm_schema.permissions p ON p.id=rp.permission_id
                JOIN crm_schema.menus m ON m.id=p.menu_id
                WHERE u.id=$1 AND u.is_active=TRUE AND ro.is_active=TRUE AND m.is_active=TRUE
                AND ro.code IN ('SUPER_ADMIN','ADMIN_KATALOG','ADMIN_MEMBERSHIP','ADMIN_KONTEN')`, userID)
			if err != nil {
				apiDBFail(w, r, err)
				return
			}
			allowed := map[string]bool{}
			for _, row := range rows {
				allowed[row.Code] = true
			}
			// Mencegah CUSTOMER yang ber-token sah melihat statistik internal.
			permitted := false
			for _, code := range []string{"rbac.manage", "catalog.product.read", "catalog.category.read", "membership.member.read", "content.banner.read", "content.article.read", "integration.qontak.read", "audit.read"} {
				if allowed[code] {
					permitted = true
					break
				}
			}
			if !permitted {
				apiFail(w, r, 403, "Tidak mempunyai izin melihat dashboard.")
				return
			}
			ctx := context.WithValue(r.Context(), dashboardPermissionsKey{}, allowed)
			next(w, r.WithContext(ctx))
		}
	}
}

type dashboardItemRow struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
	SKU   string `db:"sku"`
}
type dashboardSyncRow struct {
	ID         int64  `db:"id"`
	Title      string `db:"title"`
	RetryCount int64  `db:"retry_count"`
	StatusCode *int64 `db:"status_code"`
}
type dashboardBannerRow struct {
	ID    int64     `db:"id"`
	Title string    `db:"title"`
	EndAt time.Time `db:"end_at"`
}
type dashboardLogRow struct {
	ID        int64     `db:"id"`
	AdminName string    `db:"admin_name"`
	Action    string    `db:"action"`
	Module    string    `db:"module"`
	CreatedAt time.Time `db:"created_at"`
}

func AdminDashboardHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		permissions, ok := r.Context().Value(dashboardPermissionsKey{}).(map[string]bool)
		if !ok {
			apiFail(w, r, 403, "Izin dashboard belum diverifikasi.")
			return
		}
		type stat struct{ Code, Key, Label, Link, SQL string }
		stats := []stat{
			{"catalog.product.read", "products", "Produk aktif", "/admin/produk", `SELECT COUNT(*) FROM crm_schema.products WHERE is_active=TRUE`},
			{"catalog.category.read", "categories", "Kategori aktif", "/admin/kategori", `SELECT COUNT(*) FROM crm_schema.categories WHERE is_active=TRUE`},
			{"membership.member.read", "members", "Member", "/admin/member", `SELECT COUNT(*) FROM crm_schema.member_profiles`},
			{"content.banner.read", "banners", "Banner aktif", "/admin/banner", `SELECT COUNT(*) FROM crm_schema.banners WHERE is_active=TRUE`},
			{"content.article.read", "articles", "Artikel terbit", "/admin/artikel", `SELECT COUNT(*) FROM crm_schema.articles WHERE status='PUBLISHED'`},
			{"integration.qontak.read", "qontak_failed", "Sync Qontak gagal", "/admin/sync-qontak?status=FAILED", `SELECT COUNT(*) FROM crm_schema.qontak_contacts WHERE sync_status='FAILED'`},
		}
		summary := make([]any, 0, len(stats))
		for _, item := range stats {
			if !permissions[item.Code] {
				continue
			}
			var count int64
			if err := s.DB.QueryRowCtx(r.Context(), &count, item.SQL); err != nil {
				apiDBFail(w, r, err)
				return
			}
			summary = append(summary, map[string]any{"key": item.Key, "label": item.Label, "value": count, "link": item.Link})
		}
		attention := make([]any, 0)
		attentionTotal := map[string]int64{}
		if permissions["catalog.product.read"] {
			for _, t := range []struct{ Name, Where string }{
				{"stok_habis", `p.is_active=TRUE AND p.stock_status='HABIS'`},
				{"tanpa_foto", `p.is_active=TRUE AND NOT EXISTS (SELECT 1 FROM crm_schema.product_images pi WHERE pi.product_id=p.id)`},
			} {
				var count int64
				if err := s.DB.QueryRowCtx(r.Context(), &count, `SELECT COUNT(*) FROM crm_schema.products p WHERE `+t.Where); err != nil {
					apiDBFail(w, r, err)
					return
				}
				attentionTotal[t.Name] = count
				var items []dashboardItemRow
				if err := s.DB.QueryRowsCtx(r.Context(), &items, `SELECT p.id,p.name AS title,p.sku FROM crm_schema.products p WHERE `+t.Where+` ORDER BY p.id DESC LIMIT 3`); err != nil {
					apiDBFail(w, r, err)
					return
				}
				for _, v := range items {
					attention = append(attention, map[string]any{"type": t.Name, "id": v.ID, "title": v.Title, "sku": v.SKU})
				}
			}
		}
		if permissions["integration.qontak.read"] {
			var count int64
			if err := s.DB.QueryRowCtx(r.Context(), &count, `SELECT COUNT(*) FROM crm_schema.qontak_contacts WHERE sync_status='FAILED'`); err != nil {
				apiDBFail(w, r, err)
				return
			}
			attentionTotal["sync_gagal"] = count
			var items []dashboardSyncRow
			err := s.DB.QueryRowsCtx(r.Context(), &items, `SELECT c.id,COALESCE(mp.full_name,u.email) AS title,
                COALESCE(sl.retry_count,0) AS retry_count,sl.status_code
                FROM crm_schema.qontak_contacts c
                JOIN crm_schema.users u ON u.id=c.user_id
                LEFT JOIN crm_schema.member_profiles mp ON mp.user_id=u.id
                LEFT JOIN LATERAL (SELECT sl.retry_count,sl.status_code FROM crm_schema.sync_logs sl
                    WHERE sl.user_id=c.user_id ORDER BY sl.created_at DESC,sl.id DESC LIMIT 1) sl ON TRUE
                WHERE c.sync_status='FAILED' ORDER BY c.updated_at DESC,c.id DESC LIMIT 3`)
			if err != nil {
				apiDBFail(w, r, err)
				return
			}
			for _, v := range items {
				attention = append(attention, map[string]any{"type": "sync_gagal", "id": v.ID, "title": v.Title, "sync_status": "FAILED", "retry_count": v.RetryCount, "status_code": v.StatusCode})
			}
		}
		if permissions["content.banner.read"] {
			constQuery := ` FROM crm_schema.banners b WHERE b.is_active=TRUE AND b.end_at IS NOT NULL
                AND b.end_at >= NOW() AND b.end_at < NOW()+interval '7 days'`
			var count int64
			if err := s.DB.QueryRowCtx(r.Context(), &count, `SELECT COUNT(*)`+constQuery); err != nil {
				apiDBFail(w, r, err)
				return
			}
			attentionTotal["banner_berakhir"] = count
			var items []dashboardBannerRow
			if err := s.DB.QueryRowsCtx(r.Context(), &items, `SELECT b.id,b.title,b.end_at`+constQuery+` ORDER BY b.end_at LIMIT 3`); err != nil {
				apiDBFail(w, r, err)
				return
			}
			for _, v := range items {
				attention = append(attention, map[string]any{"type": "banner_berakhir", "id": v.ID, "title": v.Title, "end_at": v.EndAt.Format(time.RFC3339)})
			}
		}
		userID, _ := middleware.UserIDFromContext(r.Context())
		scope := "mine"
		if permissions["audit.read"] {
			scope = "all"
		}
		var logs []dashboardLogRow
		query := `SELECT l.id,COALESCE(mp.full_name,u.email) AS admin_name,l.action,l.module,l.created_at
            FROM crm_schema.activity_logs l JOIN crm_schema.users u ON u.id=l.admin_id
            LEFT JOIN crm_schema.member_profiles mp ON mp.user_id=u.id`
		var err error
		if scope == "all" {
			err = s.DB.QueryRowsCtx(r.Context(), &logs, query+` ORDER BY l.created_at DESC,l.id DESC LIMIT 8`)
		}
		if scope == "mine" {
			err = s.DB.QueryRowsCtx(r.Context(), &logs, query+` WHERE l.admin_id=$1 ORDER BY l.created_at DESC,l.id DESC LIMIT 8`, userID)
		}
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		recent := make([]any, 0, len(logs))
		for _, v := range logs {
			recent = append(recent, map[string]any{"id": v.ID, "admin_name": v.AdminName, "action": v.Action, "module": v.Module, "created_at": v.CreatedAt.Format(time.RFC3339)})
		}
		apiData(w, r, map[string]any{"summary": summary, "attention": attention, "attention_total": attentionTotal, "logs": recent, "logs_scope": scope})
	}
}