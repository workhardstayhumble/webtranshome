package postgres

import (
	"context"
	"database/sql"
	"errors"

	"webtranshome/internal/modules/identityaccess/domain/menu"
	"webtranshome/internal/modules/identityaccess/domain/permission"
	"webtranshome/internal/modules/identityaccess/domain/role"
	"webtranshome/internal/modules/identityaccess/domain/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

func (r *AccessControlRepository) GetCatalog(ctx context.Context) (service.AccessControlCatalog, error) {
	var roleRows []struct {
		ID          int64          `db:"id"`
		Code        string         `db:"code"`
		Name        string         `db:"name"`
		Description sql.NullString `db:"description"`
		IsActive    bool           `db:"is_active"`
	}
	if err := r.connection.QueryRowsCtx(ctx, &roleRows,
		`SELECT id, code, name, description, is_active
		FROM crm_schema.roles
		ORDER BY id`); err != nil {
		return service.AccessControlCatalog{}, err
	}

	var menuRows []struct {
		ID        int64          `db:"id"`
		ParentID  sql.NullInt64  `db:"parent_id"`
		Name      string         `db:"name"`
		Code      string         `db:"code"`
		Route     sql.NullString `db:"route"`
		Icon      sql.NullString `db:"icon"`
		SortOrder int            `db:"sort_order"`
		IsActive  bool           `db:"is_active"`
	}
	if err := r.connection.QueryRowsCtx(ctx, &menuRows,
		`SELECT id, parent_id, name, code, route, icon, sort_order, is_active
		FROM crm_schema.menus
		ORDER BY sort_order, id`); err != nil {
		return service.AccessControlCatalog{}, err
	}

	var permissionRows []struct {
		ID          int64          `db:"id"`
		MenuID      int64          `db:"menu_id"`
		Code        string         `db:"code"`
		Action      string         `db:"action"`
		Description sql.NullString `db:"description"`
	}
	if err := r.connection.QueryRowsCtx(ctx, &permissionRows,
		`SELECT id, menu_id, code, action, description
		FROM crm_schema.permissions
		ORDER BY menu_id, id`); err != nil {
		return service.AccessControlCatalog{}, err
	}

	catalog := service.AccessControlCatalog{
		Roles:       make([]role.Role, 0, len(roleRows)),
		Menus:       make([]menu.Menu, 0, len(menuRows)),
		Permissions: make([]permission.Permission, 0, len(permissionRows)),
	}
	for _, row := range roleRows {
		catalog.Roles = append(catalog.Roles, role.Role{
			ID:          row.ID,
			Code:        row.Code,
			Name:        row.Name,
			Description: stringPointer(row.Description),
			IsActive:    row.IsActive,
		})
	}
	for _, row := range menuRows {
		catalog.Menus = append(catalog.Menus, menu.Menu{
			ID:        row.ID,
			ParentID:  int64Pointer(row.ParentID),
			Name:      row.Name,
			Code:      row.Code,
			Route:     stringPointer(row.Route),
			Icon:      stringPointer(row.Icon),
			SortOrder: row.SortOrder,
			IsActive:  row.IsActive,
		})
	}
	for _, row := range permissionRows {
		catalog.Permissions = append(catalog.Permissions, permission.Permission{
			ID:          row.ID,
			MenuID:      row.MenuID,
			Code:        row.Code,
			Action:      row.Action,
			Description: stringPointer(row.Description),
		})
	}

	return catalog, nil
}

func (r *AccessControlRepository) GetUserAccess(ctx context.Context, userID int64) (service.UserAccess, error) {
	var account struct {
		Email    string `db:"email"`
		IsActive bool   `db:"is_active"`
	}
	if err := r.connection.QueryRowCtx(ctx, &account,
		`SELECT email, is_active FROM crm_schema.users WHERE id = $1`, userID); err != nil {
		if errors.Is(err, sqlx.ErrNotFound) {
			return service.UserAccess{}, service.ErrUserNotFound
		}
		return service.UserAccess{}, err
	}
	if !account.IsActive {
		return service.UserAccess{}, service.ErrUserNotFound
	}

	// Hanya role aktif yang masuk ke sesi frontend. Tidak pernah mengambil role dari input browser.
	var roleRows []struct {
		Code string `db:"code"`
	}
	if err := r.connection.QueryRowsCtx(ctx, &roleRows,
		`SELECT DISTINCT r.code
		 FROM crm_schema.user_roles ur
		 JOIN crm_schema.roles r ON r.id = ur.role_id AND r.is_active = TRUE
		 WHERE ur.user_id = $1
		 ORDER BY r.code`, userID); err != nil {
		return service.UserAccess{}, err
	}

	var menuRows []struct {
		ID        int64          `db:"id"`
		ParentID  sql.NullInt64  `db:"parent_id"`
		Name      string         `db:"name"`
		Code      string         `db:"code"`
		Route     sql.NullString `db:"route"`
		Icon      sql.NullString `db:"icon"`
		SortOrder int            `db:"sort_order"`
		IsActive  bool           `db:"is_active"`
	}
	const menusQuery = `WITH RECURSIVE visible_menus AS (
		SELECT m.id, m.parent_id, m.name, m.code, m.route, m.icon, m.sort_order, m.is_active
		FROM crm_schema.user_roles ur
		JOIN crm_schema.roles r ON r.id = ur.role_id AND r.is_active = TRUE
		JOIN crm_schema.role_permissions rp ON rp.role_id = r.id
		JOIN crm_schema.permissions p ON p.id = rp.permission_id
		JOIN crm_schema.menus m ON m.id = p.menu_id AND m.is_active = TRUE
		WHERE ur.user_id = $1
		UNION
		SELECT parent.id, parent.parent_id, parent.name, parent.code, parent.route,
			parent.icon, parent.sort_order, parent.is_active
		FROM crm_schema.menus parent
		JOIN visible_menus child ON child.parent_id = parent.id
		WHERE parent.is_active = TRUE
	)
	SELECT DISTINCT id, parent_id, name, code, route, icon, sort_order, is_active
	FROM visible_menus
	ORDER BY sort_order, id`
	if err := r.connection.QueryRowsCtx(ctx, &menuRows, menusQuery, userID); err != nil {
		return service.UserAccess{}, err
	}

	var permissionRows []struct {
		ID          int64          `db:"id"`
		MenuID      int64          `db:"menu_id"`
		Code        string         `db:"code"`
		Action      string         `db:"action"`
		Description sql.NullString `db:"description"`
	}
	const permissionsQuery = `SELECT DISTINCT p.id, p.menu_id, p.code, p.action, p.description
		FROM crm_schema.user_roles ur
		JOIN crm_schema.roles r ON r.id = ur.role_id AND r.is_active = TRUE
		JOIN crm_schema.role_permissions rp ON rp.role_id = r.id
		JOIN crm_schema.permissions p ON p.id = rp.permission_id
		JOIN crm_schema.menus m ON m.id = p.menu_id AND m.is_active = TRUE
		WHERE ur.user_id = $1
		ORDER BY p.code`
	if err := r.connection.QueryRowsCtx(ctx, &permissionRows, permissionsQuery, userID); err != nil {
		return service.UserAccess{}, err
	}

	access := service.UserAccess{
		UserID:      userID,
		Email:       account.Email,
		Roles:       make([]string, 0, len(roleRows)),
		Menus:       make([]menu.Menu, 0, len(menuRows)),
		Permissions: make([]permission.Permission, 0, len(permissionRows)),
	}
	for _, row := range roleRows {
		access.Roles = append(access.Roles, row.Code)
	}
	for _, row := range menuRows {
		access.Menus = append(access.Menus, menu.Menu{
			ID:        row.ID,
			ParentID:  int64Pointer(row.ParentID),
			Name:      row.Name,
			Code:      row.Code,
			Route:     stringPointer(row.Route),
			Icon:      stringPointer(row.Icon),
			SortOrder: row.SortOrder,
			IsActive:  row.IsActive,
		})
	}
	for _, row := range permissionRows {
		access.Permissions = append(access.Permissions, permission.Permission{
			ID:          row.ID,
			MenuID:      row.MenuID,
			Code:        row.Code,
			Action:      row.Action,
			Description: stringPointer(row.Description),
		})
	}

	return access, nil
}

func stringPointer(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

func int64Pointer(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	return &value.Int64
}