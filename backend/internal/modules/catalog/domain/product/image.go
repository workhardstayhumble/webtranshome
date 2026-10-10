package product

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	ErrInvalidProductImage  = errors.New("CATALOG_INVALID_PRODUCT_IMAGE")
	ErrImageProductNotFound = errors.New("CATALOG_PRODUCT_NOT_FOUND")
)

type Image struct {
	ID        int64
	ProductID int64
	ImageURL  string
	IsPrimary bool
	SortOrder int64
	CreatedAt *time.Time
}

// NewImage menerapkan aturan bisnis dasar sebelum penyimpanan.
func NewImage(productID int64, imageURL string, isPrimary bool, sortOrder int64) (Image, error) {
	imageURL = strings.TrimSpace(imageURL)
	if productID <= 0 || imageURL == "" || utf8.RuneCountInString(imageURL) > 255 ||
		sortOrder < 0 || sortOrder > 1<<31-1 {
		return Image{}, ErrInvalidProductImage
	}
	return Image{
		ProductID: productID,
		ImageURL:  imageURL,
		IsPrimary: isPrimary,
		SortOrder: sortOrder,
	}, nil
}