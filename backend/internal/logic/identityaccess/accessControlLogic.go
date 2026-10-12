package identityaccess

import (
	"context"

	"webtranshome/internal/modules/identityaccess/application/command"
	"webtranshome/internal/svc"
	"webtranshome/internal/types"
)

type AccessControlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAccessControlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccessControlLogic {
	return &AccessControlLogic{ctx: ctx, svcCtx: svcCtx}
}

func (l *AccessControlLogic) Catalog() (*types.AccessControlCatalogResponse, error) {
	catalog, err := l.svcCtx.AccessControlQueries.Catalog(l.ctx)
	if err != nil {
		return nil, err
	}

	response := &types.AccessControlCatalogResponse{
		Roles:       make([]types.AccessControlRole, 0, len(catalog.Roles)),
		Menus:       make([]types.AccessControlMenu, 0, len(catalog.Menus)),
		Permissions: make([]types.AccessControlPermission, 0, len(catalog.Permissions)),
	}
	for _, item := range catalog.Roles {
		response.Roles = append(response.Roles, types.AccessControlRole{
			ID:          item.ID,
			Code:        item.Code,
			Name:        item.Name,
			Description: item.Description,
			IsActive:    item.IsActive,
		})
	}
	for _, item := range catalog.Menus {
		response.Menus = append(response.Menus, types.AccessControlMenu{
			ID:        item.ID,
			ParentID:  item.ParentID,
			Name:      item.Name,
			Code:      item.Code,
			Route:     item.Route,
			Icon:      item.Icon,
			SortOrder: int64(item.SortOrder),
			IsActive:  item.IsActive,
		})
	}
	for _, item := range catalog.Permissions {
		response.Permissions = append(response.Permissions, types.AccessControlPermission{
			ID:          item.ID,
			MenuID:      item.MenuID,
			Code:        item.Code,
			Action:      item.Action,
			Description: item.Description,
		})
	}

	return response, nil
}

func (l *AccessControlLogic) UserAccess(userID int64) (*types.UserAccessResponse, error) {
	access, err := l.svcCtx.AccessControlQueries.UserAccess(l.ctx, userID)
	if err != nil {
		return nil, err
	}

	response := &types.UserAccessResponse{
		User: types.CurrentAccessUser{
			ID:    access.UserID,
			Email: access.Email,
			Roles: access.Roles,
		},
		Menus:           make([]types.AccessControlMenu, 0, len(access.Menus)),
		PermissionCodes: access.PermissionCodes,
	}
	for _, item := range access.Menus {
		response.Menus = append(response.Menus, types.AccessControlMenu{
			ID:        item.ID,
			ParentID:  item.ParentID,
			Name:      item.Name,
			Code:      item.Code,
			Route:     item.Route,
			Icon:      item.Icon,
			SortOrder: int64(item.SortOrder),
			IsActive:  item.IsActive,
		})
	}

	return response, nil
}

func (l *AccessControlLogic) UserRoleIDs(userID int64) (*types.AccessControlRoleAssignmentsResponse, error) {
	ids, err := l.svcCtx.AccessControlQueries.UserRoleIDs(l.ctx, userID)
	if err != nil {
		return nil, err
	}
	return &types.AccessControlRoleAssignmentsResponse{RoleIDs: ids}, nil
}

func (l *AccessControlLogic) FindUserByEmail(email string) (*types.AccessControlUserResponse, error) {
	result, err := l.svcCtx.AccessControlQueries.FindUserByEmail(l.ctx, email)
	if err != nil {
		return nil, err
	}
	return &types.AccessControlUserResponse{
		ID:       result.ID,
		Email:    result.Email,
		IsActive: result.IsActive,
		RoleIDs:  result.RoleIDs,
	}, nil
}

func (l *AccessControlLogic) RolePermissionIDs(roleID int64) (*types.AccessControlPermissionAssignmentsResponse, error) {
	ids, err := l.svcCtx.AccessControlQueries.RolePermissionIDs(l.ctx, roleID)
	if err != nil {
		return nil, err
	}
	return &types.AccessControlPermissionAssignmentsResponse{PermissionIDs: ids}, nil
}

func (l *AccessControlLogic) ReplaceUserRoles(req *types.ReplaceUserRolesRequest, userID int64) error {
	if req == nil {
		return command.ErrInvalidAccessAssignment
	}
	return l.svcCtx.ReplaceAccessAssignments.UserRoles(l.ctx, userID, req.RoleIDs)
}

func (l *AccessControlLogic) ReplaceRolePermissions(req *types.ReplaceRolePermissionsRequest, roleID int64) error {
	if req == nil {
		return command.ErrInvalidAccessAssignment
	}
	return l.svcCtx.ReplaceAccessAssignments.RolePermissions(l.ctx, roleID, req.PermissionIDs)
}

func (l *AccessControlLogic) CreateRole(req *types.CreateRoleRequest) (*types.AccessControlRole, error) {
	if req == nil {
		return nil, command.ErrInvalidRole
	}
	result, err := l.svcCtx.ManageRoles.Create(l.ctx, command.CreateRoleCommand{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}

	return &types.AccessControlRole{
		ID:          result.ID,
		Code:        result.Code,
		Name:        result.Name,
		Description: stringValue(result.Description),
		IsActive:    result.IsActive,
	}, nil
}

func (l *AccessControlLogic) UpdateRole(req *types.UpdateRoleRequest, roleID int64) error {
	if req == nil {
		return command.ErrInvalidRole
	}
	return l.svcCtx.ManageRoles.Update(l.ctx, command.UpdateRoleCommand{
		ID:          roleID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	})
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}