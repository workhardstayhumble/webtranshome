package handler

import (
	"os"
	"strings"
	"testing"
)

// Mencegah goctl generate ulang sehingga route admin kembali tanpa RBAC.
func TestAdminRoutesHavePermissionMiddleware(t *testing.T) {
	source, err := os.ReadFile("routes.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	totalAdmin := strings.Count(text, `Path: "/admin/`)
	protected := strings.Count(text, `serverCtx.RequirePermission(permission.`)
	// /admin/dashboard memakai middleware dashboardAccessMiddleware (izin
	// diverifikasi berdasarkan role aktif dan data disaring per permission).
	dashboard := strings.Count(text, `dashboardAccessMiddleware(serverCtx)`)
	if dashboard != 1 {
		t.Fatalf("dashboard guard=%d; wajib tepat satu", dashboard)
	}
	if totalAdmin == 0 || totalAdmin != protected+dashboard {
		t.Fatalf("admin routes=%d permission checks=%d: periksa routes.go setelah regenerasi", totalAdmin, protected)
	}
	if strings.Contains(text, `serverCtx.RequirePermission(permission.AccessControlManageCode)},
			rest.Route{Method: http.MethodPatch, Path: "/admin/catalog/products`) {
		t.Fatal("PATCH Product tidak boleh memakai rbac.manage")
	}
}