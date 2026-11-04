package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"webtranshome/internal/modules/catalog/domain/product"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ProductRepository struct {
	connection sqlx.SqlConn
}

func NewProductRepository(connection sqlx.SqlConn) *ProductRepository {
	return &ProductRepository{connection: connection}
}

type productRow struct {
	ID              int64           `db:"id"`
	SKU             string          `db:"sku"`
	Name            string          `db:"name"`
	Slug            string          `db:"slug"`
	CategoryID      int64           `db:"category_id"`
	BrandID         sql.NullInt64   `db:"brand_id"`
	Description     sql.NullString  `db:"description"`
	Specifications  sql.NullString  `db:"specifications"`
	DatasheetPDFURL sql.NullString  `db:"datasheet_pdf_url"`
	PriceGeneral    sql.NullFloat64 `db:"price_general"`
	ShowPrice       bool            `db:"show_price"`
	UnitSale        string          `db:"unit_sale"`
	MinOrder        int64           `db:"min_order"`
	StockStatus     string          `db:"stock_status"`
	StockQtyLabel   sql.NullString  `db:"stock_qty_label"`
	IsFeatured      bool            `db:"is_featured"`
	IsActive        bool            `db:"is_active"`
	MetaTitle       sql.NullString  `db:"meta_title"`
	MetaDescription sql.NullString  `db:"meta_description"`
}

func (r *ProductRepository) ListProducts(ctx context.Context) ([]product.Product, error) {
	var rows []productRow
	if err := r.connection.QueryRowsCtx(ctx, &rows,
		`SELECT id, sku, name, slug, category_id, brand_id, description, specifications, datasheet_pdf_url,
			price_general, unit_sale, min_order, stock_status, stock_qty_label, is_featured, is_active,
			meta_title, meta_description, show_price
		FROM crm_schema.products
		ORDER BY name, id`); err != nil {
		return nil, err
	}
	items := make([]product.Product, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapProduct(row))
	}
	return items, nil
}

func (r *ProductRepository) FindProduct(ctx context.Context, productID int64) (product.Product, error) {
	var row productRow
	if err := r.connection.QueryRowCtx(ctx, &row,
		`SELECT id, sku, name, slug, category_id, brand_id, description, specifications, datasheet_pdf_url,
			price_general, unit_sale, min_order, stock_status, stock_qty_label, is_featured, is_active,
			meta_title, meta_description, show_price
		FROM crm_schema.products WHERE id = $1`, productID); err != nil {
		return product.Product{}, mapProductError(err)
	}
	return mapProduct(row), nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, item product.Product) (product.Product, error) {
	var id int64
	err := r.connection.QueryRowCtx(ctx, &id,
		`INSERT INTO crm_schema.products (
			sku, name, slug, category_id, brand_id, description, specifications, datasheet_pdf_url,
			price_general, unit_sale, min_order, stock_status, stock_qty_label, is_featured, is_active,
			meta_title, meta_description, show_price
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id`,
		item.SKU,
		item.Name,
		item.Slug,
		item.CategoryID,
		productNullableInt64(item.BrandID),
		productNullableString(item.Description),
		productNullableString(item.Specifications),
		productNullableString(item.DatasheetPDFURL),
		productNullableFloat64(item.PriceGeneral),
		item.UnitSale,
		item.MinOrder,
		item.StockStatus,
		productNullableString(item.StockQtyLabel),
		item.IsFeatured,
		item.IsActive,
		productNullableString(item.MetaTitle),
		productNullableString(item.MetaDescription),
		item.ShowPrice,
	)
	if err != nil {
		return product.Product{}, product.ErrProductSKUConflict
	}
	item.ID = id
	return item, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, item product.Product) error {
	result, err := r.connection.ExecCtx(ctx,
		`UPDATE crm_schema.products
		SET sku = $1, name = $2, slug = $3, category_id = $4, brand_id = $5,
			description = $6, specifications = $7, datasheet_pdf_url = $8,
			price_general = $9, unit_sale = $10, min_order = $11, stock_status = $12,
			stock_qty_label = $13, is_featured = $14, is_active = $15,
			meta_title = $16, meta_description = $17, show_price = $19, updated_at = CURRENT_TIMESTAMP
		WHERE id = $18`,
		item.SKU,
		item.Name,
		item.Slug,
		item.CategoryID,
		productNullableInt64(item.BrandID),
		productNullableString(item.Description),
		productNullableString(item.Specifications),
		productNullableString(item.DatasheetPDFURL),
		productNullableFloat64(item.PriceGeneral),
		item.UnitSale,
		item.MinOrder,
		item.StockStatus,
		productNullableString(item.StockQtyLabel),
		item.IsFeatured,
		item.IsActive,
		productNullableString(item.MetaTitle),
		productNullableString(item.MetaDescription),
		item.ID,
		item.ShowPrice,
	)
	if err != nil {
		return err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return product.ErrProductNotFound
	}
	return nil
}

func (r *ProductRepository) ReplaceProductRooms(ctx context.Context, productID int64, roomIDs []int64) error {
	return r.connection.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, `DELETE FROM crm_schema.product_rooms WHERE product_id = $1`, productID); err != nil {
			return err
		}
		if len(roomIDs) == 0 {
			return nil
		}
		query := `INSERT INTO crm_schema.product_rooms (product_id, room_id) VALUES `
		args := make([]interface{}, 0, len(roomIDs)*2)
		placeholders := make([]string, 0, len(roomIDs))
		argIndex := 1
		for _, roomID := range roomIDs {
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", argIndex, argIndex+1))
			args = append(args, productID, roomID)
			argIndex += 2
		}
		query += joinPlaceholders(placeholders)
		_, err := session.ExecCtx(ctx, query, args...)
		return err
	})
}

func joinPlaceholders(values []string) string {
	if len(values) == 0 {
		return ""
	}
	result := ""
	for i, value := range values {
		if i > 0 {
			result += ", "
		}
		result += value
	}
	return result
}

func mapProduct(row productRow) product.Product {
	value := product.Product{
		ID:              row.ID,
		SKU:             row.SKU,
		Name:            row.Name,
		Slug:            row.Slug,
		CategoryID:      row.CategoryID,
		BrandID:         nullableInt64Ptr(row.BrandID),
		Description:     nullableStringPtr(row.Description),
		Specifications:  nullableStringPtr(row.Specifications),
		DatasheetPDFURL: nullableStringPtr(row.DatasheetPDFURL),
		PriceGeneral:    nullableFloat64Ptr(row.PriceGeneral),
		ShowPrice:       row.ShowPrice,
		UnitSale:        row.UnitSale,
		MinOrder:        row.MinOrder,
		StockStatus:     row.StockStatus,
		StockQtyLabel:   nullableStringPtr(row.StockQtyLabel),
		IsFeatured:      row.IsFeatured,
		IsActive:        row.IsActive,
		MetaTitle:       nullableStringPtr(row.MetaTitle),
		MetaDescription: nullableStringPtr(row.MetaDescription),
	}
	return value
}

func mapProductError(err error) error {
	if errors.Is(err, sqlx.ErrNotFound) {
		return product.ErrProductNotFound
	}
	return err
}

func productNullableString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}

func productNullableInt64(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *value, Valid: true}
}

func productNullableFloat64(value *float64) sql.NullFloat64 {
	if value == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *value, Valid: true}
}

func nullableFloat64Ptr(value sql.NullFloat64) *float64 {
	if !value.Valid {
		return nil
	}
	copy := value.Float64
	return &copy
}

func nullableStringPtr(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	copy := value.String
	return &copy
}

func nullableInt64Ptr(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	copy := value.Int64
	return &copy
}

var (
	_ product.Reader = (*ProductRepository)(nil)
	_ product.Writer = (*ProductRepository)(nil)
)
