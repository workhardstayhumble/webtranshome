package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"webtranshome/internal/middleware"
	"webtranshome/internal/modules/identityaccess/domain/permission"
	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type siteSettingsDBStub struct {
	sqlx.SqlConn
	values       map[string]string
	reads        []string
	writes       []string
	transactions int
	allowRatio   bool
	audit        map[string]any
}

func (s *siteSettingsDBStub) QueryRowsCtx(_ context.Context, target any, query string, _ ...any) error {
	s.reads = append(s.reads, query)
	rows := make([]siteSettingRow, 0, len(s.values))
	for key, value := range s.values {
		rows = append(rows, siteSettingRow{Key: key, Value: sql.NullString{String: value, Valid: true}})
	}
	// Include the legacy DB number to also verify it cannot override server config.
	*target.(*[]siteSettingRow) = rows
	return nil
}

func (s *siteSettingsDBStub) QueryRowCtx(_ context.Context, target any, query string, _ ...any) error {
	s.reads = append(s.reads, query)
	*target.(*bool) = s.allowRatio
	return nil
}

func (s *siteSettingsDBStub) ExecCtx(_ context.Context, query string, args ...any) (sql.Result, error) {
	if strings.Contains(query, "INSERT INTO crm_schema.site_settings") {
		key := args[0].(string)
		s.writes = append(s.writes, key)
		s.values[key] = args[1].(string)
	} else if strings.Contains(query, "INSERT INTO crm_schema.activity_logs") {
		if err := json.Unmarshal([]byte(args[3].(string)), &s.audit); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (s *siteSettingsDBStub) TransactCtx(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	s.transactions++
	return fn(ctx, s)
}

type siteSettingsTokenVerifier struct{}

func (siteSettingsTokenVerifier) Verify(context.Context, string) (int64, error) { return 17, nil }

type siteSettingsPermissionChecker struct{ granted bool }

func (c siteSettingsPermissionChecker) Execute(_ context.Context, userID int64, code string) (bool, error) {
	return c.granted && userID == 17 && code == permission.ContentPagesManageCode, nil
}

func protectedSettingsHandler(s *svc.ServiceContext, allowed bool) http.HandlerFunc {
	return middleware.Authenticate(siteSettingsTokenVerifier{})(
		middleware.RequirePermission(siteSettingsPermissionChecker{granted: allowed}, permission.ContentPagesManageCode)(
			UpdateSiteSettingsHandler(s)))
}

func newSettingsService() (*svc.ServiceContext, *siteSettingsDBStub) {
	db := &siteSettingsDBStub{values: map[string]string{
		"wa_number": "6289999999999", "address": "Alamat lama", "opening_hours": "08.00-17.00",
		"maps_url": "", "point_ratio_rupiah": "1000000",
	}}
	s := &svc.ServiceContext{DB: db}
	s.Config.WhatsApp.Number = "6281234567890"
	return s, db
}

func TestSiteSettingsUsesServerWhatsAppNumber(t *testing.T) {
	for _, number := range []string{"6281234567890", ""} {
		for _, admin := range []bool{false, true} {
			t.Run(fmt.Sprintf("number=%q/admin=%t", number, admin), func(t *testing.T) {
				s, db := newSettingsService()
				s.Config.WhatsApp.Number = number
				handler := PublicSiteSettingsHandler(s)
				if admin {
					handler = AdminSiteSettingsHandler(s)
				}
				w := httptest.NewRecorder()
				handler(w, httptest.NewRequest(http.MethodGet, "/api/v1/settings/public", nil))
				var response struct{ Data map[string]any }
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatal(err)
				}
				if w.Code != http.StatusOK || response.Data["wa_number"] != number {
					t.Fatalf("status=%d body=%s, want server number %q", w.Code, w.Body, number)
				}
				if _, exists := response.Data["point_ratio_rupiah"]; exists != admin {
					t.Fatalf("ratio exposed=%t, admin=%t", exists, admin)
				}
				if len(db.reads) != 1 || strings.Contains(db.reads[0], "'wa_number'") {
					t.Fatal("settings must not read the WhatsApp number from the database")
				}
			})
		}
	}
}

func TestAdminSettingsRejectsWhatsAppWritesForEveryRole(t *testing.T) {
	for _, role := range []struct {
		name   string
		logged bool
		allow  bool
	}{
		{"SUPER_ADMIN", true, true}, {"ADMIN_KONTEN", true, true},
		{"ADMIN_KATALOG", true, false}, {"ADMIN_MEMBERSHIP", true, false},
		{"CUSTOMER", true, false}, {"anonymous", false, false},
	} {
		for _, field := range []string{"wa_number", "WA_NUMBER", "waNumber"} {
			for _, value := range []string{`"6281111111111"`, `"6281234567890"`, `""`, `null`, `123`, `{}`} {
				t.Run(role.name+"/"+field+"/"+value, func(t *testing.T) {
					s, db := newSettingsService()
					body := fmt.Sprintf(`{"address":"Alamat baru","opening_hours":"09.00-18.00",%q:%s}`, field, value)
					r := httptest.NewRequest(http.MethodPut, "/api/v1/admin/settings", strings.NewReader(body))
					if role.logged {
						r.Header.Set("Authorization", "Bearer test-token")
					}
					w := httptest.NewRecorder()
					protectedSettingsHandler(s, role.allow)(w, r)
					want := http.StatusForbidden
					if !role.logged {
						want = http.StatusUnauthorized
					} else if role.allow {
						want = http.StatusUnprocessableEntity
					}
					if w.Code != want || len(db.writes) != 0 || len(db.reads) != 0 || db.transactions != 0 {
						t.Fatalf("status=%d want=%d body=%s; rejected request must not touch DB", w.Code, want, w.Body)
					}
					if s.Config.WhatsApp.Number != "6281234567890" {
						t.Fatal("server WhatsApp number changed")
					}
				})
			}
		}
	}
}

func TestAdminSettingsUpdatesOtherFieldsWithoutChangingWhatsApp(t *testing.T) {
	for _, ratio := range []bool{false, true} {
		t.Run(fmt.Sprintf("withRatio=%t", ratio), func(t *testing.T) {
			s, db := newSettingsService()
			db.allowRatio = ratio
			body := `{"address":"Alamat baru","opening_hours":"09.00-18.00","maps_url":"https://maps.example.test/"}`
			if ratio {
				body = strings.TrimSuffix(body, "}") + `,"point_ratio_rupiah":2000000}`
			}
			r := httptest.NewRequest(http.MethodPut, "/api/v1/admin/settings", strings.NewReader(body))
			r.Header.Set("Authorization", "Bearer test-token")
			w := httptest.NewRecorder()
			protectedSettingsHandler(s, true)(w, r)
			if w.Code != http.StatusOK || db.values["address"] != "Alamat baru" || db.transactions != 1 {
				t.Fatalf("settings update failed: status=%d body=%s", w.Code, w.Body)
			}
			var response struct{ Data map[string]any }
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatal(err)
			}
			if response.Data["wa_number"] != s.Config.WhatsApp.Number || db.values["wa_number"] != "6289999999999" {
				t.Fatal("update must retain server number and leave legacy DB number untouched")
			}
			for _, key := range db.writes {
				if key == "wa_number" {
					t.Fatal("WhatsApp number written to database")
				}
			}
			if db.audit == nil || db.audit["ratio_changed"] != ratio {
				t.Fatal("settings update must still be audited")
			}
			for _, field := range db.audit["fields"].([]any) {
				if field == "wa_number" {
					t.Fatal("audit incorrectly reports WhatsApp number as editable")
				}
			}
		})
	}
}

func TestAdminSettingsRatioStillRequiresAdditionalPermission(t *testing.T) {
	s, db := newSettingsService()
	r := httptest.NewRequest(http.MethodPut, "/api/v1/admin/settings", strings.NewReader(
		`{"address":"Alamat baru","opening_hours":"09.00-18.00","point_ratio_rupiah":2000000}`))
	r.Header.Set("Authorization", "Bearer test-token")
	w := httptest.NewRecorder()
	protectedSettingsHandler(s, true)(w, r)
	if w.Code != http.StatusForbidden || db.transactions != 0 || len(db.writes) != 0 {
		t.Fatalf("ratio must remain protected: status=%d body=%s", w.Code, w.Body)
	}
}
