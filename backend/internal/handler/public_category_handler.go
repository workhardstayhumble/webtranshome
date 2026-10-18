package handler

import (
    "net/http"

    "webtranshome/internal/modules/catalog/application/dto"
    "webtranshome/internal/svc"

    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/rest/httpx"
)

// publicCategoryNode is the public API shape expected by categoryService.js,
// KategoriTahap.vue, MegaMenu.vue and Katalog.vue.
type publicCategoryNode struct {
    ID        int64                 `json:"id"`
    Name      string                `json:"name"`
    Slug      string                `json:"slug"`
    ParentID  *int64                `json:"parent_id,omitempty"`
    ImageURL  string                `json:"image_url"`
    SortOrder int64                 `json:"sort_order"`
    Children  []*publicCategoryNode `json:"children"`
}

// PublicCategoriesHandler serves published categories without admin authentication.
// GET /api/v1/categories?tree=1 -> {"data":[{...,"children":[]}]}
// GET /api/v1/categories        -> {"data":[{...,"children":[]}]}, flat.
// Admin CRUD and permission checks are unchanged.
func PublicCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        categories, err := svcCtx.ListCategories.Execute(r.Context())
        if err != nil {
            logx.WithContext(r.Context()).Errorf("load public categories: %v", err)
            httpx.WriteJsonCtx(r.Context(), w, http.StatusInternalServerError, map[string]string{
                "message": "Kategori belum bisa dimuat.",
            })
            return
        }

        var data interface{}
        if r.URL.Query().Get("tree") == "1" {
            data = buildPublicCategoryTree(categories)
        } else {
            data = buildPublicCategoryFlat(categories)
        }
        httpx.OkJsonCtx(r.Context(), w, map[string]interface{}{"data": data})
    }
}

func publicCategoryFromDTO(c dto.Category) *publicCategoryNode {
    return &publicCategoryNode{
        ID: c.ID, Name: c.Name, Slug: c.Slug, ParentID: c.ParentID,
        ImageURL: c.ImageURL, SortOrder: c.SortOrder,
        Children: make([]*publicCategoryNode, 0),
    }
}

func buildPublicCategoryFlat(categories []dto.Category) []*publicCategoryNode {
    out := make([]*publicCategoryNode, 0, len(categories))
    for _, c := range categories {
        if c.IsActive {
            out = append(out, publicCategoryFromDTO(c))
        }
    }
    return out
}

// buildPublicCategoryTree preserves the order established by ListCategories:
// sort_order, name, id. Descendants of inactive/missing parents are not exposed.
func buildPublicCategoryTree(categories []dto.Category) []*publicCategoryNode {
    nodes := make(map[int64]*publicCategoryNode, len(categories))
    for _, c := range categories {
        if c.IsActive {
            nodes[c.ID] = publicCategoryFromDTO(c)
        }
    }

    roots := make([]*publicCategoryNode, 0)
    for _, c := range categories {
        node, ok := nodes[c.ID]
        if !ok { // Inactive categories are never public.
            continue
        }
        if c.ParentID == nil {
            roots = append(roots, node)
            continue
        }
        if parent, found := nodes[*c.ParentID]; found {
            parent.Children = append(parent.Children, node)
        }
        // Missing or inactive parent: keep the category hidden.
    }
    return roots
}