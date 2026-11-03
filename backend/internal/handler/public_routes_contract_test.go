package handler

import (
	"os"
	"strings"
	"testing"
)

// Membantu mencegah regresi 404 saat routes.go di-regenerate oleh goctl.
func TestPublicFrontendRoutesRegistered(t *testing.T) {
	raw, err := os.ReadFile("routes.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	for _, endpoint := range []string{
		`Path: "/categories"`, `Path: "/brands"`, `Path: "/brands/:slug"`,
		`Path: "/rooms"`, `Path: "/rooms/:slug"`,
		`Path: "/settings/public"`, `Path: "/products"`,
		`Path: "/services"`, `Path: "/services/:slug"`,
		`Path: "/admin/dashboard"`,
	} {
		if !strings.Contains(src, endpoint) {
			t.Errorf("route publik belum ada: %s", endpoint)
		}
	}
}

func TestAdminSettingsProtectedByPermission(t *testing.T) {
	raw, err := os.ReadFile("routes.go")
	if err != nil {
		t.Fatal(err)
	}
	src := string(raw)
	for _, endpoint := range []string{`Path: "/admin/settings"`} {
		if strings.Count(src, endpoint) != 2 {
			t.Errorf("butuh GET dan PUT admin settings: %s", endpoint)
		}
	}
	if strings.Count(src, "serverCtx.RequirePermission(permission.ContentPagesManageCode)") != 2 {
		t.Error("GET dan PUT admin settings harus memiliki permission")
	}
}