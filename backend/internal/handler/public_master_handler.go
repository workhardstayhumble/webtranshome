package handler

import (
	"net/http"
	"strings"

	"webtranshome/internal/svc"
)

// Bentuk respons sengaja memakai snake_case sesuai catalogService.js.
type publicBrandRow struct {
	ID           int64  `db:"id"`
	Name         string `db:"name"`
	Slug         string `db:"slug"`
	LogoURL      string `db:"logo_url"`
	ProductCount int64  `db:"product_count"`
}

func (v publicBrandRow) value() map[string]any {
	return map[string]any{"id": v.ID, "name": v.Name, "slug": v.Slug, "logo_url": v.LogoURL, "product_count": v.ProductCount}
}

const publicBrandColumns = `SELECT b.id,b.name,b.slug,COALESCE(b.logo_url,'') AS logo_url,
 (SELECT COUNT(*) FROM crm_schema.products p WHERE p.brand_id=b.id AND p.is_active=TRUE) AS product_count
 FROM crm_schema.brands b WHERE b.is_active=TRUE`

func PublicBrandsHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rows []publicBrandRow
		if err := s.DB.QueryRowsCtx(r.Context(), &rows, publicBrandColumns+` ORDER BY b.name,b.id`); err != nil {
			apiDBFail(w, r, err)
			return
		}
		data := make([]any, 0, len(rows))
		for _, v := range rows {
			data = append(data, v.value())
		}
		apiData(w, r, data)
	}
}
func PublicBrandDetailHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/api/v1/brands/")
		if slug == "" || strings.Contains(slug, "/") {
			http.NotFound(w, r)
			return
		}
		var v publicBrandRow
		if err := s.DB.QueryRowCtx(r.Context(), &v, publicBrandColumns+` AND b.slug=$1`, slug); err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, v.value())
	}
}

type publicRoomRow struct {
	ID           int64  `db:"id"`
	Name         string `db:"name"`
	Slug         string `db:"slug"`
	ImageCover   string `db:"image_cover"`
	SortOrder    int64  `db:"sort_order"`
	ProductCount int64  `db:"product_count"`
}

func (v publicRoomRow) value() map[string]any {
	return map[string]any{"id": v.ID, "name": v.Name, "slug": v.Slug, "image_cover": v.ImageCover, "sort_order": v.SortOrder, "product_count": v.ProductCount}
}

const publicRoomColumns = `SELECT room.id,room.name,room.slug,room.image_cover,room.sort_order,
 (SELECT COUNT(*) FROM crm_schema.product_rooms pr JOIN crm_schema.products p ON p.id=pr.product_id
  WHERE pr.room_id=room.id AND p.is_active=TRUE) AS product_count
 FROM crm_schema.rooms room WHERE room.is_active=TRUE`

func PublicRoomsHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var rows []publicRoomRow
		if err := s.DB.QueryRowsCtx(r.Context(), &rows, publicRoomColumns+` ORDER BY room.sort_order,room.name,room.id`); err != nil {
			apiDBFail(w, r, err)
			return
		}
		data := make([]any, 0, len(rows))
		for _, v := range rows {
			data = append(data, v.value())
		}
		apiData(w, r, data)
	}
}
func PublicRoomDetailHandler(s *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := strings.TrimPrefix(r.URL.Path, "/api/v1/rooms/")
		if slug == "" || strings.Contains(slug, "/") {
			http.NotFound(w, r)
			return
		}
		var v publicRoomRow
		if err := s.DB.QueryRowCtx(r.Context(), &v, publicRoomColumns+` AND room.slug=$1`, slug); err != nil {
			apiDBFail(w, r, err)
			return
		}
		apiData(w, r, v.value())
	}
}