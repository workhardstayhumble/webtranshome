package query

import (
	"context"

	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/modules/catalog/domain/product"
)

type ListProductImages struct {
	reader product.ImageReader
}

func NewListProductImages(reader product.ImageReader) *ListProductImages {
	return &ListProductImages{reader: reader}
}

func (q *ListProductImages) Execute(ctx context.Context, productID int64) ([]dto.ProductImage, error) {
	if productID <= 0 {
		return nil, product.ErrInvalidProductImage
	}

	exists, err := q.reader.ProductExists(ctx, productID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, product.ErrImageProductNotFound
	}

	images, err := q.reader.ListProductImages(ctx, productID)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ProductImage, 0, len(images))
	for _, image := range images {
		result = append(result, dto.ProductImageFromDomain(image))
	}
	return result, nil
}
