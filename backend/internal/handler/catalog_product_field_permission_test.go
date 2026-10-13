package handler

import (
	"reflect"
	"testing"
	"webtranshome/internal/modules/identityaccess/domain/permission"
	"webtranshome/internal/types"
)

func TestProductFieldPermissions(t *testing.T) {
	tests := []struct {
		name string
		req  *types.UpdateProductRequest
		want []string
	}{
		{"empty", &types.UpdateProductRequest{}, nil},
		{"price", &types.UpdateProductRequest{PriceGeneral: new(float64)}, []string{permission.CatalogProductPriceUpdateCode}},
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
