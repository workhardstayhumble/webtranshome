package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"webtranshome/internal/logic/catalog"
	"webtranshome/internal/middleware"
	"webtranshome/internal/modules/catalog/application/command"
	"webtranshome/internal/modules/catalog/application/query"
	"webtranshome/internal/modules/catalog/domain/product"
	"webtranshome/internal/modules/identityaccess/domain/permission"
	"webtranshome/internal/shared/response"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	catalogProductBadRequestCode = 40025
	catalogProductNotFoundCode   = 40425
	catalogProductConflictCode   = 40925
	catalogProductInternalCode   = 50000
	catalogProductInvalidJSON    = "HTTP_INVALID_JSON"
	catalogProductInternalError  = "INTERNAL_ERROR"
)

func CatalogProductsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := catalog.NewProductLogic(r.Context(), svcCtx).List()
		if err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		for i := range result.Products {
			enrichAdminProduct(r, svcCtx, &result.Products[i])
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func CatalogProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			ProductID int64 `path:"productId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, query.ErrInvalidProductID.Error())
			return
		}
		result, err := catalog.NewProductLogic(r.Context(), svcCtx).Get(path.ProductID)
		if err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		enrichAdminProduct(r, svcCtx, result)
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func CreateCatalogProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, readErr := io.ReadAll(io.LimitReader(r.Body, 2<<20))
		if readErr != nil {
			writeCatalogProductError(w, r, 400, catalogProductBadRequestCode, "Invalid request")
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
		var req types.CreateProductRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, catalogProductInvalidJSON)
			return
		}
		imagePayload, imageErr := decodeImagePayload(body)
		if imageErr == nil {
			imageErr = validateImages(imagePayload)
		}
		if imageErr != nil {
			writeCatalogProductError(w, r, 400, catalogProductBadRequestCode, imageErr.Error())
			return
		}
		result, err := catalog.NewProductLogic(r.Context(), svcCtx).Create(&req)
		if err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		if err := replaceImages(r, svcCtx.DB, result.ID, imagePayload); err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		enrichAdminProduct(r, svcCtx, result)
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func UpdateCatalogProductHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			ProductID int64 `path:"productId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, query.ErrInvalidProductID.Error())
			return
		}
		body, readErr := io.ReadAll(io.LimitReader(r.Body, 2<<20))
		if readErr != nil {
			writeCatalogProductError(w, r, 400, catalogProductBadRequestCode, "Invalid request")
			return
		}
		r.Body = io.NopCloser(bytes.NewReader(body))
		var req types.UpdateProductRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, catalogProductInvalidJSON)
			return
		}
		// Field sensitif membutuhkan izin tambahan meskipun pengguna punya catalog.product.update.
		// Null/omitted tidak dianggap perubahan; pointer dari ParseJsonBody membedakannya.
		if !authorizeProductFields(w, r, svcCtx, &req) {
			return
		}
		imagePayload, imageErr := decodeImagePayload(body)
		if imageErr == nil {
			imageErr = validateImages(imagePayload)
		}
		if imageErr != nil {
			writeCatalogProductError(w, r, 400, catalogProductBadRequestCode, imageErr.Error())
			return
		}
		result, err := catalog.NewProductLogic(r.Context(), svcCtx).Update(&req, path.ProductID)
		if err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		if err := replaceImages(r, svcCtx.DB, result.ID, imagePayload); err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		enrichAdminProduct(r, svcCtx, result)
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

// PATCH /admin/catalog/products/:productId/featured
// Endpoint khusus Admin Konten dan Admin Katalog: tidak menerima mutasi selain isFeatured.
func UpdateCatalogProductFeaturedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			ProductID int64 `path:"productId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil || path.ProductID <= 0 {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, query.ErrInvalidProductID.Error())
			return
		}
		var req struct {
			IsFeatured *bool `json:"isFeatured"`
		}
		if err := httpx.ParseJsonBody(r, &req); err != nil || req.IsFeatured == nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, catalogProductInvalidJSON)
			return
		}
		result, err := catalog.NewProductLogic(r.Context(), svcCtx).Update(&types.UpdateProductRequest{IsFeatured: req.IsFeatured}, path.ProductID)
		if err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func ReplaceCatalogProductRoomsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			ProductID int64 `path:"productId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, query.ErrInvalidProductID.Error())
			return
		}
		var req types.ReplaceProductRoomsRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, catalogProductInvalidJSON)
			return
		}
		if err := catalog.NewProductLogic(r.Context(), svcCtx).ReplaceRooms(&req, path.ProductID); err != nil {
			writeCatalogProductFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, types.CatalogMutationResponse{Success: true})
	}
}

func writeCatalogProductFailure(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, product.ErrInvalidProduct), errors.Is(err, product.ErrInvalidCategoryID),
		errors.Is(err, product.ErrInvalidBrandID), errors.Is(err, product.ErrInvalidRoomID),
		errors.Is(err, product.ErrInvalidStockStatus), errors.Is(err, query.ErrInvalidProductID),
		errors.Is(err, command.ErrInvalidProductID):
		writeCatalogProductError(w, r, http.StatusBadRequest, catalogProductBadRequestCode, err.Error())
	case errors.Is(err, product.ErrProductNotFound):
		writeCatalogProductError(w, r, http.StatusNotFound, catalogProductNotFoundCode, err.Error())
	case errors.Is(err, product.ErrProductSKUConflict), errors.Is(err, product.ErrProductSlugConflict):
		writeCatalogProductError(w, r, http.StatusConflict, catalogProductConflictCode, err.Error())
	default:
		logx.WithContext(r.Context()).Error(err)
		writeCatalogProductError(w, r, http.StatusInternalServerError, catalogProductInternalCode, catalogProductInternalError)
	}
}

func writeCatalogProductError(w http.ResponseWriter, r *http.Request, status, code int, message string) {
	httpx.WriteJsonCtx(r.Context(), w, status, response.Body{Code: code, Message: message})
}

// authorizeProductFields adalah defense-in-depth untuk operasi field-level.
// Backend selalu menjadi sumber kebenaran RBAC, bukan tombol di browser.
func authorizeProductFields(w http.ResponseWriter, r *http.Request, svcCtx *svc.ServiceContext, req *types.UpdateProductRequest) bool {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeCatalogProductError(w, r, http.StatusUnauthorized, 40110, "AUTHENTICATION_REQUIRED")
		return false
	}
	checks := requiredProductFieldPermissions(req)
	for _, code := range checks {
		granted, err := svcCtx.CheckPermission.Execute(r.Context(), userID, code)
		if err != nil {
			logx.WithContext(r.Context()).Error(err)
			writeCatalogProductError(w, r, http.StatusInternalServerError, 50000, "INTERNAL_ERROR")
			return false
		}
		if !granted {
			writeCatalogProductError(w, r, http.StatusForbidden, 40310, "AUTHORIZATION_DENIED")
			return false
		}
	}
	return true
}

// Nilai dikirim tanpa perubahan pun membutuhkan izin, agar tidak ada jalur bypass.
func requiredProductFieldPermissions(req *types.UpdateProductRequest) []string {
	if req == nil {
		return nil
	}
	var codes []string
	if req.PriceGeneral != nil {
		codes = append(codes, permission.CatalogProductPriceUpdateCode)
	}
	if req.StockStatus != nil || req.StockQtyLabel != nil {
		codes = append(codes, permission.CatalogProductStockUpdateCode)
	}
	if req.IsFeatured != nil {
		codes = append(codes, permission.CatalogProductFeaturedCode)
	}
	return codes
}

func decodeImagePayload(b []byte) ([]imageInput, error) {
	var payload struct {
		Images json.RawMessage `json:"images"`
	}
	if err := json.Unmarshal(b, &payload); err != nil {
		return nil, err
	}
	if len(payload.Images) == 0 {
		return nil, nil
	}
	var images []imageInput
	err := json.Unmarshal(payload.Images, &images)
	return images, err
}