package command

import (
	"context"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/product"
)

var ErrInvalidProduct = product.ErrInvalidProduct

type CreateProductInput struct {
	SKU             string
	Name            string
	Slug            string
	CategoryID      int64
	BrandID         *int64
	Description     string
	Specifications  string
	DatasheetPDFURL string
	PriceGeneral    *float64
	ShowPrice       bool
	UnitSale        string
	MinOrder        int64
	StockStatus     string
	StockQtyLabel   string
	IsFeatured      bool
	IsActive        bool
	MetaTitle       string
	MetaDescription string
}

type CreateProduct struct {
	reader product.Reader
	writer product.Writer
}

func NewCreateProduct(reader product.Reader, writer product.Writer) *CreateProduct {
	return &CreateProduct{reader: reader, writer: writer}
}

func (c *CreateProduct) Execute(ctx context.Context, input CreateProductInput) (dto.Product, error) {
	item, err := product.New(
		input.SKU,
		input.Name,
		input.Slug,
		input.CategoryID,
		input.BrandID,
		input.Description,
		input.Specifications,
		input.DatasheetPDFURL,
		input.PriceGeneral,
		input.ShowPrice,
		input.UnitSale,
		input.MinOrder,
		input.StockStatus,
		input.StockQtyLabel,
		input.IsFeatured,
		input.IsActive,
		input.MetaTitle,
		input.MetaDescription,
	)
	if err != nil {
		return dto.Product{}, err
	}
	created, err := c.writer.CreateProduct(ctx, item)
	if err != nil {
		return dto.Product{}, err
	}
	return toProductDTO(created), nil
}
