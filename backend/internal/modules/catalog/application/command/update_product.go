package command

import (
	"context"
	"errors"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/product"
)

var ErrInvalidProductID = errors.New("CATALOG_INVALID_PRODUCT_ID")

type UpdateProductInput struct {
	ID              int64
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

type UpdateProduct struct {
	writer product.Writer
	reader product.Reader
}

func NewUpdateProduct(reader product.Reader, writer product.Writer) *UpdateProduct {
	return &UpdateProduct{reader: reader, writer: writer}
}

func (c *UpdateProduct) Execute(ctx context.Context, input UpdateProductInput) (dto.Product, error) {
	if input.ID <= 0 {
		return dto.Product{}, ErrInvalidProductID
	}
	item, err := c.reader.FindProduct(ctx, input.ID)
	if err != nil {
		return dto.Product{}, err
	}
	if err := item.Apply(product.Changes{
		SKU:             input.SKU,
		Name:            input.Name,
		Slug:            input.Slug,
		CategoryID:      input.CategoryID,
		BrandID:         input.BrandID,
		Description:     input.Description,
		Specifications:  input.Specifications,
		DatasheetPDFURL: input.DatasheetPDFURL,
		PriceGeneral:    input.PriceGeneral,
		ShowPrice:       input.ShowPrice,
		ClearPrice:      input.ClearPrice,
		UnitSale:        input.UnitSale,
		MinOrder:        input.MinOrder,
		StockStatus:     input.StockStatus,
		StockQtyLabel:   input.StockQtyLabel,
		IsFeatured:      input.IsFeatured,
		IsActive:        input.IsActive,
		MetaTitle:       input.MetaTitle,
		MetaDescription: input.MetaDescription,
	}); err != nil {
		return dto.Product{}, err
	}
	if err := c.writer.UpdateProduct(ctx, item); err != nil {
		return dto.Product{}, err
	}
	return toProductDTO(item), nil
}

func toProductDTO(item product.Product) dto.Product {
	var description, specifications, datasheetPDFURL, stockQtyLabel, metaTitle, metaDescription string
	if item.Description != nil {
		description = *item.Description
	}
	if item.Specifications != nil {
		specifications = *item.Specifications
	}
	if item.DatasheetPDFURL != nil {
		datasheetPDFURL = *item.DatasheetPDFURL
	}
	if item.StockQtyLabel != nil {
		stockQtyLabel = *item.StockQtyLabel
	}
	if item.MetaTitle != nil {
		metaTitle = *item.MetaTitle
	}
	if item.MetaDescription != nil {
		metaDescription = *item.MetaDescription
	}
	return dto.Product{
		ID:              item.ID,
		SKU:             item.SKU,
		Name:            item.Name,
		Slug:            item.Slug,
		CategoryID:      item.CategoryID,
		BrandID:         item.BrandID,
		Description:     description,
		Specifications:  specifications,
		DatasheetPDFURL: datasheetPDFURL,
		PriceGeneral:    item.PriceGeneral,
		ShowPrice:       item.ShowPrice,
		UnitSale:        item.UnitSale,
		MinOrder:        item.MinOrder,
		StockStatus:     item.StockStatus,
		StockQtyLabel:   stockQtyLabel,
		IsFeatured:      item.IsFeatured,
		IsActive:        item.IsActive,
		MetaTitle:       metaTitle,
		MetaDescription: metaDescription,
	}
}
