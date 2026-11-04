package dto

type Product struct {
	ID              int64
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
