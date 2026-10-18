package handler

import (
	"testing"

	"webtranshome/internal/modules/catalog/application/dto"
)

func TestBuildPublicCategoryTree(t *testing.T) {
    rootID := int64(1)
    inactiveID := int64(9)
    data := []dto.Category{
        {ID: 1, Name: "Keramik & Granit", Slug: "keramik-granit", ImageURL: "/api/v1/uploads/category.jpg", IsActive: true},
        {ID: 2, Name: "Keramik Lantai", Slug: "keramik-lantai", ParentID: &rootID, IsActive: true},
        {ID: 9, Name: "Tidak Aktif", Slug: "nonaktif", IsActive: false},
        {ID: 10, Name: "Anak Tidak Aktif", Slug: "anak-nonaktif", ParentID: &inactiveID, IsActive: true},
    }
    roots := buildPublicCategoryTree(data)
    if len(roots) != 1 { t.Fatalf("want one visible root; got %d", len(roots)) }
    if roots[0].Slug != "keramik-granit" || roots[0].ImageURL != "/api/v1/uploads/category.jpg" {
        t.Fatalf("unexpected root: %+v", roots[0])
    }
    if len(roots[0].Children) != 1 || roots[0].Children[0].Slug != "keramik-lantai" {
        t.Fatalf("unexpected children: %+v", roots[0].Children)
    }
    if roots[0].Children[0].Children == nil {
        t.Fatal("children must be [] rather than null")
    }
}

func TestPublicCategoryEmptyIsArray(t *testing.T) {
    roots := buildPublicCategoryTree(nil)
    if roots == nil || len(roots) != 0 {
        t.Fatalf("empty category response must be [], got %#v", roots)
    }
    flat := buildPublicCategoryFlat(nil)
    if flat == nil || len(flat) != 0 {
        t.Fatalf("empty flat list must be [], got %#v", flat)
    }
}