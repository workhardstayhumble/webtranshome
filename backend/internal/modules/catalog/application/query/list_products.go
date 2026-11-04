package query

import (
	"context"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/product"
)

type ListProducts struct {
	reader product.Reader
}

func NewListProducts(reader product.Reader) *ListProducts {
	return &ListProducts{reader: reader}
}

func (q *ListProducts) Execute(ctx context.Context) ([]dto.Product, error) {
	items, err := q.reader.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Product, 0, len(items))
	for _, item := range items {
		result = append(result, toProductDTO(item))
	}
	return result, nil
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
