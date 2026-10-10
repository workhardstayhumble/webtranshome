package product

import "context"

// ImageReader menyediakan operasi baca untuk use case gambar produk.
type ImageReader interface {
	ProductExists(ctx context.Context, productID int64) (bool, error)
	ListProductImages(ctx context.Context, productID int64) ([]Image, error)
}

// ImageWriter menyediakan operasi tulis yang diperlukan use case tambah gambar.
type ImageWriter interface {
	CreateProductImage(ctx context.Context, image Image) (Image, error)
}
