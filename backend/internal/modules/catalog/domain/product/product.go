package product

import (
	"errors"
	"strings"
	"unicode/utf8"

	"webtranshome/internal/modules/catalog/domain/slug"
)

var (
	ErrInvalidProduct      = errors.New("CATALOG_INVALID_PRODUCT")
	ErrProductNotFound     = errors.New("CATALOG_PRODUCT_NOT_FOUND")
	ErrProductSKUConflict   = errors.New("CATALOG_PRODUCT_SKU_CONFLICT")
	ErrProductSlugConflict  = errors.New("CATALOG_PRODUCT_SLUG_CONFLICT")
	ErrInvalidCategoryID   = errors.New("CATALOG_INVALID_CATEGORY_ID")
	ErrInvalidBrandID      = errors.New("CATALOG_INVALID_BRAND_ID")
	ErrInvalidRoomID       = errors.New("CATALOG_INVALID_ROOM_ID")
	ErrInvalidStockStatus  = errors.New("CATALOG_INVALID_STOCK_STATUS")
)

var validStockStatuses = map[string]struct{}{
	"TERSEDIA": {},
	"SISA_STOK": {},
	"PRE_ORDER": {},
	"HABIS": {},
}

type Product struct {
	ID              int64
	SKU             string
	Name            string
	Slug            string
	CategoryID      int64
	BrandID         *int64
	Description     *string
	Specifications  *string
	DatasheetPDFURL *string
	PriceGeneral    *float64
	ShowPrice       bool
	UnitSale        string
	MinOrder        int64
	StockStatus     string
	StockQtyLabel   *string
	IsFeatured      bool
	IsActive        bool
	MetaTitle       *string
	MetaDescription *string
}

type Changes struct {
	SKU             *string
	Name            *string
	Slug            *string
	CategoryID      *int64
	BrandID         *int64
	Description     *string
	Specifications  *string
	DatasheetPDFURL *string
	PriceGeneral    *float64
	ShowPrice       *bool
	// ClearPrice mengosongkan harga; hanya sah bila harga tidak ditampilkan.
	ClearPrice      bool
	UnitSale        *string
	MinOrder        *int64
	StockStatus     *string
	StockQtyLabel   *string
	IsFeatured      *bool
	IsActive        *bool
	MetaTitle       *string
	MetaDescription *string
}

func New(sku, name, slug string, categoryID int64, brandID *int64, description, specifications, datasheetPDFURL string,
	priceGeneral *float64, showPrice bool, unitSale string, minOrder int64, stockStatus string, stockQtyLabel string,
	isFeatured, isActive bool, metaTitle, metaDescription string) (Product, error) {
	item := Product{
		SKU:             strings.TrimSpace(sku),
		Name:            strings.TrimSpace(name),
		Slug:            strings.TrimSpace(slug),
		CategoryID:      categoryID,
		BrandID:         copyInt64(brandID),
		Description:     optionalString(description),
		Specifications:  optionalString(specifications),
		DatasheetPDFURL: optionalString(datasheetPDFURL),
		PriceGeneral:    copyFloat64(priceGeneral),
		ShowPrice:       showPrice,
		UnitSale:        strings.TrimSpace(unitSale),
		MinOrder:        minOrder,
		StockStatus:     strings.TrimSpace(stockStatus),
		StockQtyLabel:   optionalString(stockQtyLabel),
		IsFeatured:      isFeatured,
		IsActive:        isActive,
		MetaTitle:       optionalString(metaTitle),
		MetaDescription: optionalString(metaDescription),
	}
	if err := item.Validate(); err != nil {
		return Product{}, err
	}
	return item, nil
}

func (p Product) Validate() error {
	if p.SKU == "" || utf8.RuneCountInString(p.SKU) > 100 ||
		p.Name == "" || utf8.RuneCountInString(p.Name) > 255 ||
		!slug.IsValid(p.Slug, 255) ||
		p.CategoryID <= 0 ||
		(p.BrandID != nil && *p.BrandID <= 0) ||
		(p.PriceGeneral != nil && *p.PriceGeneral < 0) ||
		(p.ShowPrice && p.PriceGeneral == nil) ||
		p.UnitSale == "" || utf8.RuneCountInString(p.UnitSale) > 50 ||
		p.MinOrder < 1 ||
		p.StockStatus == "" ||
		!isValidStockStatus(p.StockStatus) ||
		(p.Description != nil && utf8.RuneCountInString(*p.Description) > 10000) ||
		(p.Specifications != nil && utf8.RuneCountInString(*p.Specifications) > 5000) ||
		(p.DatasheetPDFURL != nil && utf8.RuneCountInString(*p.DatasheetPDFURL) > 255) ||
		(p.StockQtyLabel != nil && utf8.RuneCountInString(*p.StockQtyLabel) > 50) ||
		(p.MetaTitle != nil && utf8.RuneCountInString(*p.MetaTitle) > 255) ||
		(p.MetaDescription != nil && utf8.RuneCountInString(*p.MetaDescription) > 500) {
		return ErrInvalidProduct
	}
	return nil
}

func (p *Product) Apply(changes Changes) error {
	if changes.Empty() {
		return ErrInvalidProduct
	}
	candidate := p.applyChanges(changes)
	if err := candidate.Validate(); err != nil {
		return err
	}
	*p = candidate
	return nil
}

func (c Changes) Empty() bool {
	return c.SKU == nil && c.Name == nil && c.Slug == nil && c.CategoryID == nil &&
		c.BrandID == nil && c.Description == nil && c.Specifications == nil &&
		c.DatasheetPDFURL == nil && c.PriceGeneral == nil && c.ShowPrice == nil && !c.ClearPrice && c.UnitSale == nil &&
		c.MinOrder == nil && c.StockStatus == nil && c.StockQtyLabel == nil &&
		c.IsFeatured == nil && c.IsActive == nil && c.MetaTitle == nil && c.MetaDescription == nil
}

func (p Product) applyChanges(changes Changes) Product {
	candidate := p
	if changes.SKU != nil {
		candidate.SKU = strings.TrimSpace(*changes.SKU)
	}
	if changes.Name != nil {
		candidate.Name = strings.TrimSpace(*changes.Name)
	}
	if changes.Slug != nil {
		candidate.Slug = strings.TrimSpace(*changes.Slug)
	}
	if changes.CategoryID != nil {
		candidate.CategoryID = *changes.CategoryID
	}
	if changes.BrandID != nil {
		candidate.BrandID = copyInt64(changes.BrandID)
	}
	if changes.Description != nil {
		candidate.Description = optionalString(*changes.Description)
	}
	if changes.Specifications != nil {
		candidate.Specifications = optionalString(*changes.Specifications)
	}
	if changes.DatasheetPDFURL != nil {
		candidate.DatasheetPDFURL = optionalString(*changes.DatasheetPDFURL)
	}
	if changes.ClearPrice {
		candidate.PriceGeneral = nil
	}
	if changes.PriceGeneral != nil {
		candidate.PriceGeneral = copyFloat64(changes.PriceGeneral)
	}
	if changes.ShowPrice != nil {
		candidate.ShowPrice = *changes.ShowPrice
	}
	if changes.UnitSale != nil {
		candidate.UnitSale = strings.TrimSpace(*changes.UnitSale)
	}
	if changes.MinOrder != nil {
		candidate.MinOrder = *changes.MinOrder
	}
	if changes.StockStatus != nil {
		candidate.StockStatus = strings.TrimSpace(*changes.StockStatus)
	}
	if changes.StockQtyLabel != nil {
		candidate.StockQtyLabel = optionalString(*changes.StockQtyLabel)
	}
	if changes.IsFeatured != nil {
		candidate.IsFeatured = *changes.IsFeatured
	}
	if changes.IsActive != nil {
		candidate.IsActive = *changes.IsActive
	}
	if changes.MetaTitle != nil {
		candidate.MetaTitle = optionalString(*changes.MetaTitle)
	}
	if changes.MetaDescription != nil {
		candidate.MetaDescription = optionalString(*changes.MetaDescription)
	}
	return candidate
}

func copyInt64(value *int64) *int64 {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func copyFloat64(value *float64) *float64 {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func optionalString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func isValidStockStatus(value string) bool {
	_, ok := validStockStatuses[strings.TrimSpace(value)]
	return ok
}
