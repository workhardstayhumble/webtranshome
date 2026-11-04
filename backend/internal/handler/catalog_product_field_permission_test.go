package handler

import (
	"reflect"
	"testing"
	"webtranshome/internal/modules/identityaccess/domain/permission"
	"webtranshome/internal/types"
)

func TestRequestsPriceRemoval(t *testing.T) {
	tampil, sembunyi, harga := true, false, 1000.0
	tests := []struct {
		name string
		body string
		req  types.UpdateProductRequest
		want bool
	}{
		{"hidden with null price", `{"showPrice":false,"priceGeneral":null}`, types.UpdateProductRequest{ShowPrice: &sembunyi}, true},
		{"hidden without price field keeps price", `{"showPrice":false}`, types.UpdateProductRequest{ShowPrice: &sembunyi}, false},
		{"hidden with price keeps price", `{"showPrice":false,"priceGeneral":1000}`, types.UpdateProductRequest{ShowPrice: &sembunyi, PriceGeneral: &harga}, false},
		{"shown with null price", `{"showPrice":true,"priceGeneral":null}`, types.UpdateProductRequest{ShowPrice: &tampil}, false},
		{"null price without showPrice", `{"priceGeneral":null}`, types.UpdateProductRequest{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requestsPriceRemoval([]byte(tt.body), &tt.req); got != tt.want {
				t.Fatalf("got %v want %v", got, tt.want)
			}
		})
	}
}

func TestProductFieldPermissions(t *testing.T) {
	tests := []struct {
		name string
		req  *types.UpdateProductRequest
		want []string
	}{
		{"empty", &types.UpdateProductRequest{}, nil},
		{"price", &types.UpdateProductRequest{PriceGeneral: new(float64)}, []string{permission.CatalogProductPriceUpdateCode}},
		{"show price", &types.UpdateProductRequest{ShowPrice: new(bool)}, []string{permission.CatalogProductPriceUpdateCode}},
		{"stock", &types.UpdateProductRequest{StockStatus: new(string)}, []string{permission.CatalogProductStockUpdateCode}},
		{"stock label", &types.UpdateProductRequest{StockQtyLabel: new(string)}, []string{permission.CatalogProductStockUpdateCode}},
		{"featured", &types.UpdateProductRequest{IsFeatured: new(bool)}, []string{permission.CatalogProductFeaturedCode}},
		{"all", &types.UpdateProductRequest{PriceGeneral: new(float64), StockStatus: new(string), IsFeatured: new(bool)}, []string{permission.CatalogProductPriceUpdateCode, permission.CatalogProductStockUpdateCode, permission.CatalogProductFeaturedCode}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := requiredProductFieldPermissions(tt.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v want %v", got, tt.want)
			}
		})
	}
}
