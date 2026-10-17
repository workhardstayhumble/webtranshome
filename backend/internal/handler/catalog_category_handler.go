package handler

import (
	"errors"
	"fmt"
	"net/http"

	"webtranshome/internal/logic/catalog"
	"webtranshome/internal/modules/catalog/application/command"
	"webtranshome/internal/modules/catalog/application/query"
	"webtranshome/internal/modules/catalog/domain/category"
	"webtranshome/internal/shared/response"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	catalogCategoryBadRequestCode = 40020
	catalogCategoryNotFoundCode   = 40420
	catalogCategoryConflictCode   = 40920
	catalogCategoryInternalCode   = 50000
	catalogCategoryInvalidJSON    = "HTTP_INVALID_JSON"
	catalogCategoryInternalError  = "INTERNAL_ERROR"
)

func CatalogCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := catalog.NewCategoryLogic(r.Context(), svcCtx).List()
		if err != nil {
			writeCatalogCategoryFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func CatalogCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			CategoryID int64 `path:"categoryId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogCategoryError(w, r, http.StatusBadRequest, catalogCategoryBadRequestCode, query.ErrInvalidCategoryID.Error())
			return
		}
		result, err := catalog.NewCategoryLogic(r.Context(), svcCtx).Get(path.CategoryID)
		if err != nil {
			writeCatalogCategoryFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func CreateCatalogCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateCategoryRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			// Jangan log body atau header Authorization: keduanya bisa mengandung data sensitif.
			logx.WithContext(r.Context()).Error(fmt.Errorf("CreateCatalogCategory invalid JSON (content-type=%q): %w", r.Header.Get("Content-Type"), err))
			writeCatalogCategoryError(w, r, http.StatusBadRequest, catalogCategoryBadRequestCode, catalogCategoryInvalidJSON)
			return
		}
		result, err := catalog.NewCategoryLogic(r.Context(), svcCtx).Create(&req)
		if err != nil {
			writeCatalogCategoryFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func UpdateCatalogCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var path struct {
			CategoryID int64 `path:"categoryId"`
		}
		if err := httpx.ParsePath(r, &path); err != nil {
			writeCatalogCategoryError(w, r, http.StatusBadRequest, catalogCategoryBadRequestCode, query.ErrInvalidCategoryID.Error())
			return
		}
		var req types.UpdateCategoryRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			writeCatalogCategoryError(w, r, http.StatusBadRequest, catalogCategoryBadRequestCode, catalogCategoryInvalidJSON)
			return
		}
		result, err := catalog.NewCategoryLogic(r.Context(), svcCtx).Update(&req, path.CategoryID)
		if err != nil {
			writeCatalogCategoryFailure(w, r, err)
			return
		}
		httpx.OkJsonCtx(r.Context(), w, result)
	}
}

func writeCatalogCategoryFailure(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, category.ErrInvalidCategory), errors.Is(err, category.ErrInvalidParent),
		errors.Is(err, query.ErrInvalidCategoryID), errors.Is(err, command.ErrInvalidCategoryID):
		writeCatalogCategoryError(w, r, http.StatusBadRequest, catalogCategoryBadRequestCode, err.Error())
	case errors.Is(err, category.ErrCategoryNotFound):
		writeCatalogCategoryError(w, r, http.StatusNotFound, catalogCategoryNotFoundCode, err.Error())
	case errors.Is(err, category.ErrCategorySlugConflict):
		writeCatalogCategoryError(w, r, http.StatusConflict, catalogCategoryConflictCode, err.Error())
	default:
		logx.WithContext(r.Context()).Error(err)
		writeCatalogCategoryError(w, r, http.StatusInternalServerError, catalogCategoryInternalCode, catalogCategoryInternalError)
	}
}

func writeCatalogCategoryError(w http.ResponseWriter, r *http.Request, status, code int, message string) {
	httpx.WriteJsonCtx(r.Context(), w, status, response.Body{Code: code, Message: message})
}
