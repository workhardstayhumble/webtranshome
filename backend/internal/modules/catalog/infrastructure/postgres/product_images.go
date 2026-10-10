package postgres

import (
	"context"
	"database/sql"
	"errors"

	"webtranshome/internal/modules/catalog/domain/product"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Repository konkret memakai sqlx, konsisten dengan repository Category/Brand/Room.
type ProductImageRepository struct {
	connection sqlx.SqlConn
}

func NewProductImageRepository(connection sqlx.SqlConn) *ProductImageRepository {
	return &ProductImageRepository{connection: connection}
}

type productImageRow struct {
	ID        int64        `db:"id"`
	ProductID int64        `db:"product_id"`
	ImageURL  string       `db:"image_url"`
	IsPrimary bool         `db:"is_primary"`
	SortOrder int64        `db:"sort_order"`
	CreatedAt sql.NullTime `db:"created_at"`
}

func (r *ProductImageRepository) ProductExists(ctx context.Context, productID int64) (bool, error) {
	var exists bool
	err := r.connection.QueryRowCtx(ctx, &exists,
		`SELECT EXISTS(SELECT 1 FROM crm_schema.products WHERE id = $1)`, productID)
	return exists, err
}

func (r *ProductImageRepository) ListProductImages(ctx context.Context, productID int64) ([]product.Image, error) {
	var rows []productImageRow
	err := r.connection.QueryRowsCtx(ctx, &rows,
		`SELECT id, product_id, image_url, is_primary, sort_order, created_at
         FROM crm_schema.product_images
         WHERE product_id = $1
         ORDER BY is_primary DESC, sort_order ASC, id ASC`, productID)
	if err != nil {
		return nil, err
	}

	images := make([]product.Image, 0, len(rows))
	for _, row := range rows {
		images = append(images, mapProductImage(row))
	}
	return images, nil
}

// Lock baris product agar pergantian gambar utama terserialisasi
// antar request yang menggunakan repository ini.
func (r *ProductImageRepository) CreateProductImage(ctx context.Context, image product.Image) (product.Image, error) {
	var inserted productImageRow

	err := r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		var id int64
		if err := session.QueryRowCtx(ctx, &id,
			`SELECT id FROM crm_schema.products WHERE id = $1 FOR UPDATE`, image.ProductID); err != nil {
			if errors.Is(err, sqlx.ErrNotFound) {
				return product.ErrImageProductNotFound
			}
			return err
		}

		if image.IsPrimary {
			if _, err := session.ExecCtx(ctx,
				`UPDATE crm_schema.product_images SET is_primary = FALSE
                 WHERE product_id = $1 AND is_primary = TRUE`, image.ProductID); err != nil {
				return err
			}
		}

		return session.QueryRowCtx(ctx, &inserted,
			`INSERT INTO crm_schema.product_images
                 (product_id, image_url, is_primary, sort_order)
             VALUES ($1, $2, $3, $4)
             RETURNING id, product_id, image_url, is_primary, sort_order, created_at`,
			image.ProductID, image.ImageURL, image.IsPrimary, image.SortOrder)
	})
	if err != nil {
		return product.Image{}, err
	}
	return mapProductImage(inserted), nil
}

func mapProductImage(row productImageRow) product.Image {
	image := product.Image{
		ID:        row.ID,
		ProductID: row.ProductID,
		ImageURL:  row.ImageURL,
		IsPrimary: row.IsPrimary,
		SortOrder: row.SortOrder,
	}
	if row.CreatedAt.Valid {
		createdAt := row.CreatedAt.Time
		image.CreatedAt = &createdAt
	}
	return image
}

var (
	_ product.ImageReader = (*ProductImageRepository)(nil)
	_ product.ImageWriter = (*ProductImageRepository)(nil)
)
