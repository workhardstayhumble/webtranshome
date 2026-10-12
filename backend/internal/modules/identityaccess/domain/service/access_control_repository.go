package service

import (
	"context"
	"errors"

	"webtranshome/internal/modules/identityaccess/domain/menu"
	"webtranshome/internal/modules/identityaccess/domain/permission"
	"webtranshome/internal/modules/identityaccess/domain/role"
)

var (
	ErrUserNotFound       = errors.New("IDENTITY_USER_NOT_FOUND")
	ErrRoleNotFound       = errors.New("IDENTITY_ROLE_NOT_FOUND")
	ErrRoleCodeConflict   = errors.New("IDENTITY_ROLE_CODE_CONFLICT")
	ErrPermissionNotFound = errors.New("IDENTITY_PERMISSION_NOT_FOUND")
	ErrProtectedRole      = role.ErrProtectedRole
)

type AccessControlCatalog struct {
	Roles       []role.Role
	Menus       []menu.Menu
	Permissions []permission.Permission
}

type UserAccess struct {
	UserID      int64
	Email       string
	Roles       []string
	Menus       []menu.Menu
	Permissions []permission.Permission
}

type AccessControlUser struct {
	ID       int64
	Email    string
	IsActive bool
	RoleIDs  []int64
}

type AccessControlCatalogRepository interface {
	GetCatalog(ctx context.Context) (AccessControlCatalog, error)
}

type UserAccessRepository interface {
	GetUserAccess(ctx context.Context, userID int64) (UserAccess, error)
	FindUserByEmail(ctx context.Context, email string) (AccessControlUser, error)
}

type AccessControlQueryRepository interface {
	GetUserRoleIDs(ctx context.Context, userID int64) ([]int64, error)
	GetRolePermissionIDs(ctx context.Context, roleID int64) ([]int64, error)
}

type RoleRepository interface {
	CreateRole(ctx context.Context, item role.Role) (role.Role, error)
	UpdateRole(ctx context.Context, item role.Role) error
}

type AccessAssignmentRepository interface {
	ReplaceUserRoles(ctx context.Context, userID int64, roleIDs []int64) error
	ReplaceRolePermissions(ctx context.Context, roleID int64, permissionIDs []int64) error
}