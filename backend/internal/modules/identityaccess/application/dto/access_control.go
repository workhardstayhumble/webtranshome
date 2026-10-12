package dto

type AccessControlRole struct {
	ID          int64
	Code        string
	Name        string
	Description string
	IsActive    bool
}

type AccessControlMenu struct {
	ID        int64
	ParentID  int64
	Name      string
	Code      string
	Route     string
	Icon      string
	SortOrder int
	IsActive  bool
}

type AccessControlPermission struct {
	ID          int64
	MenuID      int64
	Code        string
	Action      string
	Description string
}

type AccessControlCatalog struct {
	Roles       []AccessControlRole
	Menus       []AccessControlMenu
	Permissions []AccessControlPermission
}

type EffectiveAccess struct {
	UserID          int64
	Email           string
	Roles           []string
	Menus           []AccessControlMenu
	PermissionCodes []string
}

type AccessControlUser struct {
	ID       int64
	Email    string
	IsActive bool
	RoleIDs  []int64
}