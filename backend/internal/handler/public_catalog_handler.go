package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"
)

type publicProductRow struct {
	ID              int64          `db:"id"`
	SKU             string         `db:"sku"`
	Name            string         `db:"name"`
	Slug            string         `db:"slug"`
	CategoryID      int64          `db:"category_id"`
	BrandID         sql.NullInt64  `db:"brand_id"`
	CategoryName    string         `db:"category_name"`
	CategorySlug    string         `db:"category_slug"`
	BrandName       sql.NullString `db:"brand_name"`
	BrandSlug       sql.NullString `db:"brand_slug"`
	Description     sql.NullString `db:"description"`
	Specifications  sql.NullString `db:"specifications"`
	PriceGeneral    float64        `db:"price_general"`
	UnitSale        string         `db:"unit_sale"`
	MinOrder        int64          `db:"min_order"`
	StockStatus     string         `db:"stock_status"`
	StockQtyLabel   sql.NullString `db:"stock_qty_label"`
	IsFeatured      bool           `db:"is_featured"`
	PrimaryImageURL sql.NullString `db:"primary_image_url"`
}

const publicProductColumns = `SELECT p.id,p.sku,p.name,p.slug,p.category_id,p.brand_id,c.name AS category_name,c.slug AS category_slug,b.name AS brand_name,b.slug AS brand_slug,p.description,p.specifications,p.price_general,p.unit_sale,p.min_order,p.stock_status,p.stock_qty_label,p.is_featured, COALESCE((SELECT i.image_url FROM crm_schema.product_images i WHERE i.product_id=p.id ORDER BY i.is_primary DESC,i.sort_order ASC,i.id ASC LIMIT 1),'') AS primary_image_url FROM crm_schema.products p JOIN crm_schema.categories c ON c.id=p.category_id LEFT JOIN crm_schema.brands b ON b.id=p.brand_id WHERE p.is_active=TRUE`

func fetchPublicProducts(ctx context.Context, s *svc.ServiceContext) ([]publicProductRow, error) {
	var rows []publicProductRow
	err := s.DB.QueryRowsCtx(ctx, &rows, publicProductColumns+` ORDER BY p.id DESC`)
	return rows, err
}
func publicItem(row publicProductRow) map[string]interface{} {
	var specs interface{} = map[string]interface{}{}
	if row.Specifications.Valid {
		_ = json.Unmarshal([]byte(row.Specifications.String), &specs)
	}
	var brand interface{} = nil
	if row.BrandID.Valid {
		brand = map[string]interface{}{"id": row.BrandID.Int64, "name": row.BrandName.String, "slug": row.BrandSlug.String}
	}
	return map[string]interface{}{"id": row.ID, "sku": row.SKU, "name": row.Name, "slug": row.Slug, "category_id": row.CategoryID, "brand_id": row.BrandID.Int64, "description": row.Description.String, "specifications": specs, "price_general": row.PriceGeneral, "unit_sale": row.UnitSale, "min_order": row.MinOrder, "stock_status": row.StockStatus, "stock_qty_label": row.StockQtyLabel.String, "is_featured": row.IsFeatured, "category": map[string]interface{}{"id": row.CategoryID, "name": row.CategoryName, "slug": row.CategorySlug}, "brand": brand, "primary_image_url": row.PrimaryImageURL.String}
}
func writePublicJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func PublicProductsHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := fetchPublicProducts(r.Context(), s)
		if err != nil {
			http.Error(w, "Database error", 500)
			return
		}
		q := r.URL.Query()
		items := make([]map[string]interface{}, 0)
		for _, row := range rows {
			if q.Get("q") != "" && !strings.Contains(strings.ToLower(row.Name+" "+row.SKU), strings.ToLower(q.Get("q"))) {
				continue
			}
			if q.Get("category") != "" && row.CategorySlug != q.Get("category") {
				continue
			}
			if q.Get("brand") != "" && !strings.Contains(","+q.Get("brand")+",", ","+row.BrandSlug.String+",") {
				continue
			}
			if q.Get("stock") != "" && !strings.Contains(","+q.Get("stock")+",", ","+row.StockStatus+",") {
				continue
			}
			if q.Get("featured") == "1" && !row.IsFeatured {
				continue
			}
			if q.Get("exclude") == strconv.FormatInt(row.ID, 10) {
				continue
			}
			if v, e := strconv.ParseFloat(q.Get("min"), 64); e == nil && row.PriceGeneral < v {
				continue
			}
			if v, e := strconv.ParseFloat(q.Get("max"), 64); e == nil && row.PriceGeneral > v {
				continue
			}
			items = append(items, publicItem(row))
		}
		switch q.Get("sort") {
		case "nama":
			sortProducts(items, "name", false)
		case "harga_terendah":
			sortProducts(items, "price_general", false)
		case "harga_tertinggi":
			sortProducts(items, "price_general", true)
		}
		page, _ := strconv.Atoi(q.Get("page"))
		if page < 1 {
			page = 1
		}
		per, _ := strconv.Atoi(q.Get("per_page"))
		if per < 1 {
			per = 24
		}
		if per > 60 {
			per = 60
		}
		total := len(items)
		pages := (total + per - 1) / per
		if pages < 1 {
			pages = 1
		}
		if page > pages {
			page = pages
		}
		start := (page - 1) * per
		end := start + per
		if end > total {
			end = total
		}
		writePublicJSON(w, 200, map[string]interface{}{"data": items[start:end], "meta": map[string]interface{}{"page": page, "per_page": per, "total": total, "total_pages": pages, "facets": map[string]interface{}{}}})
	}
}
func PublicProductDetailHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
		if slug == "" {
			http.NotFound(w, r)
			return
		}
		rows, err := fetchPublicProducts(r.Context(), s)
		if err != nil {
			http.Error(w, "Database error", 500)
			return
		}
		for _, row := range rows {
			if row.Slug != slug {
				continue
			}
			v := publicItem(row)
			images, err := readImages(r.Context(), s, row.ID)
			if err != nil {
				http.Error(w, "Database error", 500)
				return
			}
			v["images"] = images
			v["rooms"] = []interface{}{}
			writePublicJSON(w, 200, map[string]interface{}{"data": v})
			return
		}
		http.NotFound(w, r)
	}
}
func PublicProductSuggestHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := fetchPublicProducts(r.Context(), s)
		if err != nil {
			http.Error(w, "Database error", 500)
			return
		}
		items := make([]interface{}, 0)
		q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		if len(q) >= 2 {
			for _, row := range rows {
				if strings.Contains(strings.ToLower(row.Name+" "+row.SKU), q) {
					items = append(items, map[string]interface{}{"id": row.ID, "name": row.Name, "slug": row.Slug, "sku": row.SKU})
					if len(items) == 5 {
						break
					}
				}
			}
		}
		writePublicJSON(w, 200, map[string]interface{}{"data": map[string]interface{}{"products": items, "categories": []interface{}{}, "brands": []interface{}{}}})
	}
}
func readImages(ctx context.Context, s *svc.ServiceContext, id int64) ([]imageOutput, error) {
	var rows []struct {
		ID        int64  `db:"id"`
		ImageURL  string `db:"image_url"`
		IsPrimary bool   `db:"is_primary"`
		SortOrder int64  `db:"sort_order"`
	}
	err := s.DB.QueryRowsCtx(ctx, &rows, `SELECT id,image_url,is_primary,sort_order FROM crm_schema.product_images WHERE product_id=$1 ORDER BY is_primary DESC,sort_order,id`, id)
	output := make([]imageOutput, 0, len(rows))
	for _, row := range rows {
		output = append(output, imageOutput{row.ID, row.ImageURL, row.IsPrimary, row.SortOrder})
	}
	return output, err
}
func enrichAdminProduct(r *http.Request, s *svc.ServiceContext, v *types.CatalogProduct) {
	if v == nil {
		return
	}
	images, err := readImages(r.Context(), s, v.ID)
	if err != nil {
		return
	}
	for _, im := range images {
		v.Images = append(v.Images, struct {
			ID        int64  `json:"id"`
			ImageURL  string `json:"image_url"`
			IsPrimary bool   `json:"is_primary"`
			SortOrder int64  `json:"sort_order"`
		}{im.ID, im.ImageURL, im.IsPrimary, im.SortOrder})
		if v.PrimaryImageURL == "" {
			v.PrimaryImageURL = im.ImageURL
		}
	}
}

func sortProducts(items []map[string]interface{}, key string, reverse bool) {
	sort.SliceStable(items, func(i, j int) bool {
		if key == "name" {
			a, b := items[i][key].(string), items[j][key].(string)
			if reverse {
				return a > b
			}
			return a < b
		}
		a, b := items[i][key].(float64), items[j][key].(float64)
		if reverse {
			return a > b
		}
		return a < b
	})
}