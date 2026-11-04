package catalog

import (
	"context"

	"webtranshome/internal/modules/catalog/application/command"
	"webtranshome/internal/modules/catalog/application/dto"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"
)

type ProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProductLogic {
	return &ProductLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *ProductLogic) List() (*types.CatalogProductListResponse, error) {
	items, err := l.svcCtx.ListProducts.Execute(l.ctx)
	if err != nil {
		return nil, err
	}
	result := &types.CatalogProductListResponse{Products: make([]types.CatalogProduct, 0, len(items))}
	for _, item := range items {
		result.Products = append(result.Products, catalogProduct(item))
	}
	return result, nil
}

func (l *ProductLogic) Get(productID int64) (*types.CatalogProduct, error) {
	item, err := l.svcCtx.GetProduct.Execute(l.ctx, productID)
	if err != nil {
		return nil, err
	}
	result := catalogProduct(item)
	return &result, nil
}

func (l *ProductLogic) Create(req *types.CreateProductRequest) (*types.CatalogProduct, error) {
	if req == nil {
		return nil, command.ErrInvalidProduct
	}
	item, err := l.svcCtx.CreateProduct.Execute(l.ctx, command.CreateProductInput{
		ShowPrice:       req.ShowPrice == nil || *req.ShowPrice,
		SKU:             req.SKU,
		Name:            req.Name,
		Slug:            req.Slug,
		CategoryID:      req.CategoryID,
		BrandID:         req.BrandID,
		Description:     req.Description,
		Specifications:  req.Specifications,
		DatasheetPDFURL: req.DatasheetPDFURL,
		PriceGeneral:    req.PriceGeneral,
		UnitSale:        req.UnitSale,
		MinOrder:        req.MinOrder,
		StockStatus:     req.StockStatus,
		StockQtyLabel:   req.StockQtyLabel,
		IsFeatured:      req.IsFeatured,
		IsActive:        req.IsActive,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	})
	if err != nil {
		return nil, err
	}
	result := catalogProduct(item)
	return &result, nil
}

// clearPrice mengosongkan harga: JSON null pada priceGeneral tidak bisa dibedakan dari field yang tidak dikirim di request.
func (l *ProductLogic) Update(req *types.UpdateProductRequest, productID int64, clearPrice bool) (*types.CatalogProduct, error) {
	if req == nil {
		return nil, command.ErrInvalidProduct
	}
	item, err := l.svcCtx.UpdateProduct.Execute(l.ctx, command.UpdateProductInput{
		ShowPrice:       req.ShowPrice,
		ClearPrice:      clearPrice,
		ID:              productID,
		SKU:             req.SKU,
		Name:            req.Name,
		Slug:            req.Slug,
		CategoryID:      req.CategoryID,
		BrandID:         req.BrandID,
		Description:     req.Description,
		Specifications:  req.Specifications,
		DatasheetPDFURL: req.DatasheetPDFURL,
		PriceGeneral:    req.PriceGeneral,
		UnitSale:        req.UnitSale,
		MinOrder:        req.MinOrder,
		StockStatus:     req.StockStatus,
		StockQtyLabel:   req.StockQtyLabel,
		IsFeatured:      req.IsFeatured,
		IsActive:        req.IsActive,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	})
	if err != nil {
		return nil, err
	}
	result := catalogProduct(item)
	return &result, nil
}

func (l *ProductLogic) ReplaceRooms(req *types.ReplaceProductRoomsRequest, productID int64) error {
	if req == nil {
		return command.ErrInvalidProduct
	}
	return l.svcCtx.ReplaceProductRooms.Execute(l.ctx, command.ReplaceProductRoomsInput{
		ProductID: productID,
		RoomIDs:   req.RoomIDs,
	})
}

func catalogProduct(item dto.Product) types.CatalogProduct {
	return types.CatalogProduct{
		ID:              item.ID,
		SKU:             item.SKU,
		Name:            item.Name,
		Slug:            item.Slug,
		CategoryID:      item.CategoryID,
		BrandID:         item.BrandID,
		Description:     item.Description,
		Specifications:  item.Specifications,
		DatasheetPDFURL: item.DatasheetPDFURL,
		PriceGeneral:    item.PriceGeneral,
		ShowPrice:       item.ShowPrice,
		UnitSale:        item.UnitSale,
		MinOrder:        item.MinOrder,
		StockStatus:     item.StockStatus,
		StockQtyLabel:   item.StockQtyLabel,
		IsFeatured:      item.IsFeatured,
		IsActive:        item.IsActive,
		MetaTitle:       item.MetaTitle,
		MetaDescription: item.MetaDescription,
	}
}
