package command

import (
	"context"
	"errors"
	"testing"

	"webtranshome/internal/modules/catalog/domain/product"
)

type productRepositoryStub struct {
	found      product.Product
	created    product.Product
	updated    product.Product
	findCalls  int
	createCalls int
	updateCalls int
	findErr    error
	createErr  error
	updateErr  error
}

func (r *productRepositoryStub) ListProducts(context.Context) ([]product.Product, error) {
	return nil, nil
}

func (r *productRepositoryStub) FindProduct(context.Context, int64) (product.Product, error) {
	r.findCalls++
	return r.found, r.findErr
}

func (r *productRepositoryStub) CreateProduct(_ context.Context, item product.Product) (product.Product, error) {
	r.createCalls++
	r.created = item
	item.ID = 21
	return item, r.createErr
}

func (r *productRepositoryStub) UpdateProduct(_ context.Context, item product.Product) error {
	r.updateCalls++
	r.updated = item
	return r.updateErr
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name      string
		input     CreateProductInput
		wantErr   error
		wantCalls int
	}{
		{name: "creates valid product", input: CreateProductInput{SKU: "SKU-1", Name: "Lamp", Slug: "lamp", CategoryID: 5, PriceGeneral: productFloat64Pointer(120000), UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA"}, wantCalls: 1},
		{name: "rejects invalid product", input: CreateProductInput{SKU: "", Name: "Lamp", Slug: "lamp", CategoryID: 5, PriceGeneral: productFloat64Pointer(120000), UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA"}, wantErr: product.ErrInvalidProduct},
		{name: "rejects invalid stock status", input: CreateProductInput{SKU: "SKU-2", Name: "Lamp", Slug: "lamp-2", CategoryID: 5, PriceGeneral: productFloat64Pointer(120000), UnitSale: "pcs", MinOrder: 1, StockStatus: "N/A"}, wantErr: product.ErrInvalidProduct},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &productRepositoryStub{}
			created, err := NewCreateProduct(repository, repository).Execute(context.Background(), test.input)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Execute() error = %v, want %v", err, test.wantErr)
			}
			if repository.createCalls != test.wantCalls {
				t.Fatalf("CreateProduct() calls = %d, want %d", repository.createCalls, test.wantCalls)
			}
			if test.wantErr == nil && (created.ID != 21 || created.SKU != "SKU-1" || created.Name != "Lamp") {
				t.Fatalf("created product = %+v", created)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	repository := &productRepositoryStub{found: product.Product{ID: 9, SKU: "SKU-1", Name: "Lamp", Slug: "lamp", CategoryID: 5, PriceGeneral: productFloat64Pointer(1000), UnitSale: "pcs", MinOrder: 1, StockStatus: "TERSEDIA", IsActive: true}}
	result, err := NewUpdateProduct(repository, repository).Execute(context.Background(), UpdateProductInput{
		ID: 9,
		Name: productStringPointer("Lamp Pro"),
		PriceGeneral: productFloat64Pointer(1500),
		StockStatus: productStringPointer("HABIS"),
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if repository.updateCalls != 1 || result.Name != "Lamp Pro" || *result.PriceGeneral != 1500 || result.StockStatus != "HABIS" {
		t.Fatalf("updated product = %+v, updates = %d", result, repository.updateCalls)
	}
}

func TestUpdateProductRejectsInvalidID(t *testing.T) {
	_, err := NewUpdateProduct(&productRepositoryStub{}, &productRepositoryStub{}).Execute(context.Background(), UpdateProductInput{ID: 0})
	if !errors.Is(err, ErrInvalidProductID) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrInvalidProductID)
	}
}

func productStringPointer(value string) *string { return &value }
func productFloat64Pointer(value float64) *float64 { return &value }

var (
	_ product.Reader = (*productRepositoryStub)(nil)
	_ product.Writer = (*productRepositoryStub)(nil)
)
