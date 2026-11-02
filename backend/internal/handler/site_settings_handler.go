package handler

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"webtranshome/internal/middleware"
	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// PostgreSQL site_settings menggunakan pasangan setting_key / setting_value.
// Jangan mengekspos pengaturan internal atau password melalui API publik.
type siteSettingRow struct {
	Key   string         `db:"setting_key"`
	Value sql.NullString `db:"setting_value"`
}

func getSiteSettings(ctx context.Context, s *svc.ServiceContext, includeRatio bool) (map[string]any, error) {
	result := map[string]any{"wa_number": s.Config.WhatsApp.Number, "address": "", "opening_hours": "", "maps_url": ""}
	if includeRatio {
		result["point_ratio_rupiah"] = int64(0)
	}
	var rows []siteSettingRow
	err := s.DB.QueryRowsCtx(ctx, &rows, `SELECT setting_key,setting_value FROM crm_schema.site_settings
        WHERE setting_key IN ('address','opening_hours','maps_url','point_ratio_rupiah')`)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row.Key == "point_ratio_rupiah" {
			if includeRatio {
				v, e := strconv.ParseInt(strings.TrimSpace(row.Value.String), 10, 64)
				if e == nil {
					result[row.Key] = v
				}
			}
		} else if row.Key == "address" || row.Key == "opening_hours" || row.Key == "maps_url" {
			result[row.Key] = row.Value.String
		}
	}
	return result, nil
}

func PublicSiteSettingsHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := getSiteSettings(r.Context(), s, false)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, data)
	}
}
func AdminSiteSettingsHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := getSiteSettings(r.Context(), s, true)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, data)
	}
}

type siteSettingsInput struct {
	Address          string  `json:"address"`
	OpeningHours     string  `json:"opening_hours"`
	MapsURL          *string `json:"maps_url"`
	PointRatioRupiah *int64  `json:"point_ratio_rupiah"`
}

// Admin Konten boleh mengatur info toko; rasio poin dilindungi izin rbac.manage.
// Hanya pengguna dengan izin pengelolaan akses (secara default Super Admin)
// yang boleh mengubah rasio. Perubahan melalui SQL manual di luar cakupan API.
// wa_number tidak diterima: apiDecode menolak field ini untuk semua pengguna.
func UpdateSiteSettingsHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input siteSettingsInput
		if !apiDecode(w, r, &input) {
			return
		}
		input.Address = strings.TrimSpace(input.Address)
		input.OpeningHours = strings.TrimSpace(input.OpeningHours)
		maps := ""
		if input.MapsURL != nil {
			maps = strings.TrimSpace(*input.MapsURL)
		}
		if len(input.Address) == 0 || len(input.Address) > 800 || len(input.OpeningHours) == 0 || len(input.OpeningHours) > 200 || len(maps) > 1000 || (maps != "" && !strings.HasPrefix(maps, "https://")) {
			apiFail(w, r, 422, "Periksa alamat, jam buka, dan link https://.")
			return
		}
		if input.PointRatioRupiah != nil {
			if *input.PointRatioRupiah < 1000 || *input.PointRatioRupiah > 1000000000 {
				apiFail(w, r, 422, "Rasio poin tidak valid.")
				return
			}
			userID, ok := middleware.UserIDFromContext(r.Context())
			if !ok {
				apiFail(w, r, 401, "Sesi tidak valid.")
				return
			}
			var allowed bool
			err := s.DB.QueryRowCtx(r.Context(), &allowed, `SELECT EXISTS (
                SELECT 1 FROM crm_schema.user_roles ur
                JOIN crm_schema.users u ON u.id=ur.user_id
                JOIN crm_schema.roles ro ON ro.id=ur.role_id
                JOIN crm_schema.role_permissions rp ON rp.role_id=ro.id
                JOIN crm_schema.permissions p ON p.id=rp.permission_id
                JOIN crm_schema.menus m ON m.id=p.menu_id
                WHERE ur.user_id=$1 AND u.is_active AND ro.is_active AND m.is_active AND p.code='rbac.manage')`, userID)
			if err != nil {
				apiDBFail(w, r, err)
				return
			}
			if !allowed {
				apiFail(w, r, 403, "Hanya pengelola akses yang boleh mengubah rasio poin.")
				return
			}
		}
		updates := map[string]string{"address": input.Address, "opening_hours": input.OpeningHours, "maps_url": maps}
		if input.PointRatioRupiah != nil {
			updates["point_ratio_rupiah"] = strconv.FormatInt(*input.PointRatioRupiah, 10)
		}
		err := s.DB.TransactCtx(r.Context(), func(ctx context.Context, session sqlx.Session) error {
			for key, value := range updates {
				if _, e := session.ExecCtx(ctx, `INSERT INTO crm_schema.site_settings(setting_key,setting_value,updated_at)
                    VALUES($1,$2,NOW()) ON CONFLICT(setting_key) DO UPDATE
                    SET setting_value=EXCLUDED.setting_value,updated_at=NOW()`, key, value); e != nil {
					return e
				}
			}
			return nil
		})
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		auditAdminAction(s, r, "settings.update", "Settings", map[string]any{"fields": []string{"address", "opening_hours", "maps_url"}, "ratio_changed": input.PointRatioRupiah != nil})
		result, err := getSiteSettings(r.Context(), s, true)
		if err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, result)
	}
}
