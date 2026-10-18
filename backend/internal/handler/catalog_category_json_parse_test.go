package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"webtranshome/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// Regression: go-zero needs `,optional` (not `,omitempty`) for request fields.
func TestCreateCategoryJsonAllowsMissingParentID(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantParent bool
	}{
		{
			name: "kategori utama tanpa parentId",
			body: `{"name":"Keramik & Granit","slug":"keramik-granit","imageUrl":"/uploads/keramik.jpg","sortOrder":1}`,
		},
		{
			name: "kategori utama tanpa parentId dan tanpa imageUrl",
			body: `{"name":"Keramik & Granit","slug":"keramik-granit","sortOrder":1}`,
		},
		{
			name:       "subkategori dengan parentId",
			body:       `{"name":"Keramik Lantai","slug":"keramik-lantai","parentId":1,"sortOrder":2}`,
			wantParent: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/catalog/categories", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			var payload types.CreateCategoryRequest
			if err := httpx.ParseJsonBody(req, &payload); err != nil {
				t.Fatalf("ParseJsonBody gagal: %v", err)
			}
			if (payload.ParentID != nil) != tt.wantParent {
				t.Fatalf("ParentID=%v, ingin memiliki parent=%v", payload.ParentID, tt.wantParent)
			}
		})
	}
}

func TestUpdateCategoryJsonAllowsPartialPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/catalog/categories/1", strings.NewReader(`{"isActive":false}`))
	req.Header.Set("Content-Type", "application/json")
	var payload types.UpdateCategoryRequest
	if err := httpx.ParseJsonBody(req, &payload); err != nil {
		t.Fatalf("ParseJsonBody PATCH parsial gagal: %v", err)
	}
	if payload.IsActive == nil || *payload.IsActive {
		t.Fatalf("isActive = %v, ingin false", payload.IsActive)
	}
}
