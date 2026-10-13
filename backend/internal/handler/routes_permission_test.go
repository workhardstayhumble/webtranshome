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
	if totalAdmin == 0 || totalAdmin != protected {
		t.Fatalf("admin routes=%d permission checks=%d: periksa routes.go setelah regenerasi", totalAdmin, protected)
	}
	if strings.Contains(text, `serverCtx.RequirePermission(permission.AccessControlManageCode)},
			rest.Route{Method: http.MethodPatch, Path: "/admin/catalog/products`) {
		t.Fatal("PATCH Product tidak boleh memakai rbac.manage")
	}
}
