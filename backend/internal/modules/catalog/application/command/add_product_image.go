package command

import (
	"context"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/product"
)

// AddProductImage adalah use case untuk menambahkan referensi gambar produk.
type AddProductImage struct {
	reader product.ImageReader
	writer product.ImageWriter
}

func NewAddProductImage(reader product.ImageReader, writer product.ImageWriter) *AddProductImage {
	return &AddProductImage{reader: reader, writer: writer}
}

func (c *AddProductImage) Execute(ctx context.Context, input dto.AddProductImageInput) (dto.ProductImage, error) {
	image, err := product.NewImage(input.ProductID, input.ImageURL, input.IsPrimary, input.SortOrder)
	if err != nil {
		return dto.ProductImage{}, err
	}

	exists, err := c.reader.ProductExists(ctx, image.ProductID)
	if err != nil {
		return dto.ProductImage{}, err
	}
	if !exists {
		return dto.ProductImage{}, product.ErrImageProductNotFound
	}

	created, err := c.writer.CreateProductImage(ctx, image)
	if err != nil {
		return dto.ProductImage{}, err
	}
	return dto.ProductImageFromDomain(created), nil
}
