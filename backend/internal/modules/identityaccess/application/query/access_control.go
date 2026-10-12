package query

import (
	"context"
	"errors"
	"strings"

	"webtranshome/internal/modules/identityaccess/application/dto"
	"webtranshome/internal/modules/identityaccess/domain/service"
)

var ErrInvalidAccessControlQuery = errors.New("IDENTITY_INVALID_ACCESS_CONTROL_QUERY")

type AccessControl struct {
	catalogRepository    service.AccessControlCatalogRepository
	userAccessRepository service.UserAccessRepository
	queryRepository      service.AccessControlQueryRepository
}

func NewAccessControl(
	catalogRepository service.AccessControlCatalogRepository,
	userAccessRepository service.UserAccessRepository,
	queryRepository service.AccessControlQueryRepository,
) *AccessControl {
	return &AccessControl{
		catalogRepository:    catalogRepository,
		userAccessRepository: userAccessRepository,
		queryRepository:      queryRepository,
	}
}

func (q *AccessControl) Catalog(ctx context.Context) (dto.AccessControlCatalog, error) {
	catalog, err := q.catalogRepository.GetCatalog(ctx)
	if err != nil {
		return dto.AccessControlCatalog{}, err
	}

	result := dto.AccessControlCatalog{
		Roles:       make([]dto.AccessControlRole, 0, len(catalog.Roles)),
		Menus:       make([]dto.AccessControlMenu, 0, len(catalog.Menus)),
		Permissions: make([]dto.AccessControlPermission, 0, len(catalog.Permissions)),
	}
	for _, item := range catalog.Roles {
		result.Roles = append(result.Roles, dto.AccessControlRole{
			ID:          item.ID,
			Code:        item.Code,
			Name:        item.Name,
			Description: stringValue(item.Description),
			IsActive:    item.IsActive,
		})
	}
	for _, item := range catalog.Menus {
		result.Menus = append(result.Menus, dto.AccessControlMenu{
			ID:        item.ID,
			ParentID:  int64Value(item.ParentID),
			Name:      item.Name,
			Code:      item.Code,
			Route:     stringValue(item.Route),
			Icon:      stringValue(item.Icon),
			SortOrder: item.SortOrder,
			IsActive:  item.IsActive,
		})
	}
	for _, item := range catalog.Permissions {
		result.Permissions = append(result.Permissions, dto.AccessControlPermission{
			ID:          item.ID,
			MenuID:      item.MenuID,
			Code:        item.Code,
			Action:      item.Action,
			Description: stringValue(item.Description),
		})
	}

	return result, nil
}

func (q *AccessControl) UserAccess(ctx context.Context, userID int64) (dto.EffectiveAccess, error) {
	if userID <= 0 {
		return dto.EffectiveAccess{}, ErrInvalidAccessControlQuery
	}

	access, err := q.userAccessRepository.GetUserAccess(ctx, userID)
	if err != nil {
		return dto.EffectiveAccess{}, err
	}

	result := dto.EffectiveAccess{
		UserID:          access.UserID,
		Email:           access.Email,
		Roles:           append([]string{}, access.Roles...),
		Menus:           make([]dto.AccessControlMenu, 0, len(access.Menus)),
		PermissionCodes: make([]string, 0, len(access.Permissions)),
	}
	for _, item := range access.Menus {
		result.Menus = append(result.Menus, dto.AccessControlMenu{
			ID:        item.ID,
			ParentID:  int64Value(item.ParentID),
			Name:      item.Name,
			Code:      item.Code,
			Route:     stringValue(item.Route),
			Icon:      stringValue(item.Icon),
			SortOrder: item.SortOrder,
			IsActive:  item.IsActive,
		})
	}
	for _, item := range access.Permissions {
		result.PermissionCodes = append(result.PermissionCodes, item.Code)
	}

	return result, nil
}

func (q *AccessControl) FindUserByEmail(ctx context.Context, email string) (dto.AccessControlUser, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return dto.AccessControlUser{}, ErrInvalidAccessControlQuery
	}

	user, err := q.userAccessRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return dto.AccessControlUser{}, err
	}
	return dto.AccessControlUser{
		ID:       user.ID,
		Email:    user.Email,
		IsActive: user.IsActive,
		RoleIDs:  user.RoleIDs,
	}, nil
}

func (q *AccessControl) UserRoleIDs(ctx context.Context, userID int64) ([]int64, error) {
	if userID <= 0 {
		return nil, ErrInvalidAccessControlQuery
	}

	return q.queryRepository.GetUserRoleIDs(ctx, userID)
}

func (q *AccessControl) RolePermissionIDs(ctx context.Context, roleID int64) ([]int64, error) {
	if roleID <= 0 {
		return nil, ErrInvalidAccessControlQuery
	}

	return q.queryRepository.GetRolePermissionIDs(ctx, roleID)
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func int64Value(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}