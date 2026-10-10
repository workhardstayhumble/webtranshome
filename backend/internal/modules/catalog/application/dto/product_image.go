package dto

import (
	"time"

	"webtranshome/internal/modules/catalog/domain/product"
)

// Input use case, bukan struct tabel database ataupun HTTP generated type.
type AddProductImageInput struct {
	ProductID int64
	ImageURL  string
	IsPrimary bool
	SortOrder int64
}

type ProductImage struct {
	ID        int64
	ProductID int64
	ImageURL  string
	IsPrimary bool
	SortOrder int64
	CreatedAt *time.Time
}

func ProductImageFromDomain(image product.Image) ProductImage {
	return ProductImage{
		ID:        image.ID,
		ProductID: image.ProductID,
		ImageURL:  image.ImageURL,
		IsPrimary: image.IsPrimary,
		SortOrder: image.SortOrder,
		CreatedAt: image.CreatedAt,
	}
}
