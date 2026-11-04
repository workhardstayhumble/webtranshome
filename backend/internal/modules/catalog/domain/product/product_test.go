package product

import (
	"errors"
	"testing"
)

func TestNewValidatesAndNormalizes(t *testing.T) {
	brandID := int64(7)
	item, err := New(
		" SKU-1 ",
		"  Smart Lamp ",
		" smart-lamp ",
		5,
		&brandID,
		"  detail  ",
		" {\"power\":\"12W\"} ",
		"/datasheet.pdf",
		float64Pointer(150000),
		true,
		" pcs ",
		2,
		"TERSEDIA",
		"2 pcs",
		true,
		true,
		"Smart Lamp",
		"A smart home lamp",
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if item.SKU != "SKU-1" || item.Name != "Smart Lamp" || item.Slug != "smart-lamp" || item.CategoryID != 5 ||
		item.BrandID == nil || *item.BrandID != brandID || item.Description == nil || *item.Description != "detail" ||
		item.Specifications == nil || *item.Specifications != "{\"power\":\"12W\"}" || item.DatasheetPDFURL == nil ||
		*item.DatasheetPDFURL != "/datasheet.pdf" || item.PriceGeneral == nil || *item.PriceGeneral != 150000 || !item.ShowPrice || item.UnitSale != "pcs" ||
		item.MinOrder != 2 || item.StockStatus != "TERSEDIA" || item.StockQtyLabel == nil || *item.StockQtyLabel != "2 pcs" ||
		!item.IsFeatured || !item.IsActive || item.MetaTitle == nil || *item.MetaTitle != "Smart Lamp" ||
		item.MetaDescription == nil || *item.MetaDescription != "A smart home lamp" {
		t.Fatalf("New() = %+v", item)
	}
}

func TestNewRejectsInvalidValues(t *testing.T) {
	brandID := int64(7)
	tests := []struct {
		name    string
		item    Product
		wantErr error
	}{
		{name: "empty sku", item: Product{SKU: "", Name: "Lamp", Slug: "lamp", CategoryID: 1, PriceGeneral: float64Pointer(1000), UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA"}, wantErr: ErrInvalidProduct},
		{name: "invalid stock status", item: Product{SKU: "SKU-1", Name: "Lamp", Slug: "lamp", CategoryID: 1, PriceGeneral: float64Pointer(1000), UnitSale: "pcs", MinOrder: 1, StockStatus: "SAATINI"}, wantErr: ErrInvalidProduct},
		{name: "negative price", item: Product{SKU: "SKU-1", Name: "Lamp", Slug: "lamp", CategoryID: 1, PriceGeneral: float64Pointer(-1), UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA"}, wantErr: ErrInvalidProduct},
		{name: "invalid brand id", item: Product{SKU: "SKU-1", Name: "Lamp", Slug: "lamp", CategoryID: 1, BrandID: &[]int64{0}[0], PriceGeneral: float64Pointer(1000), UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA"}, wantErr: ErrInvalidProduct},
		{name: "invalid slug", item: Product{SKU: "SKU-1", Name: "Lamp", Slug: "bad slug", CategoryID: 1, BrandID: &brandID, PriceGeneral: float64Pointer(1000), UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA"}, wantErr: ErrInvalidProduct},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.item.Validate(); !errors.Is(err, test.wantErr) {
				t.Fatalf("Validate() error = %v, want %v", err, test.wantErr)
			}
		})
	}
}

func TestApplyChanges(t *testing.T) {
	item := Product{ID: 8, SKU: "SKU-1", Name: "Lamp", Slug: "lamp", CategoryID: 1, PriceGeneral: float64Pointer(1000), UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA", IsActive: true}
	if err := item.Apply(Changes{PriceGeneral: float64Pointer(1500), StockStatus: stringPointer("HABIS"), IsFeatured: boolPointer(true)}); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if *item.PriceGeneral != 1500 || item.StockStatus != "HABIS" || !item.IsFeatured {
		t.Fatalf("Apply() = %+v", item)
	}
}

func TestApplyChangesRejectsEmptyPatch(t *testing.T) {
	item := Product{ID: 8, SKU: "SKU-1", Name: "Lamp", Slug: "lamp", CategoryID: 1, PriceGeneral: float64Pointer(1000), UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA", IsActive: true}
	if err := item.Apply(Changes{}); !errors.Is(err, ErrInvalidProduct) {
		t.Fatalf("Apply() error = %v, want %v", err, ErrInvalidProduct)
	}
}

func TestHiddenPriceMayBeEmpty(t *testing.T) {
	base := Product{SKU: "SKU-1", Name: "Lamp", Slug: "lamp", CategoryID: 1, UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA"}
	hidden := base
	if err := hidden.Validate(); err != nil {
		t.Fatalf("harga disembunyikan boleh kosong, error = %v", err)
	}
	shown := base
	shown.ShowPrice = true
	if err := shown.Validate(); !errors.Is(err, ErrInvalidProduct) {
		t.Fatalf("harga yang ditampilkan wajib ada, error = %v", err)
	}
}

func TestClearPrice(t *testing.T) {
	item := Product{ID: 8, SKU: "SKU-1", Name: "Lamp", Slug: "lamp", CategoryID: 1, PriceGeneral: float64Pointer(1000), ShowPrice: true, UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA"}
	if err := item.Apply(Changes{ClearPrice: true}); !errors.Is(err, ErrInvalidProduct) {
		t.Fatalf("harga yang masih ditampilkan tidak boleh dikosongkan, error = %v", err)
	}
	if err := item.Apply(Changes{ClearPrice: true, ShowPrice: boolPointer(false)}); err != nil {
		t.Fatalf("Apply() error = %v", err)
	}
	if item.PriceGeneral != nil || item.ShowPrice {
		t.Fatalf("Apply() = %+v", item)
	}
}

func float64Pointer(value float64) *float64 { return &value }
func stringPointer(value string) *string     { return &value }
func boolPointer(value bool) *bool           { return &value }
