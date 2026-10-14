package types

type AccessControlCatalogResponse struct {
	Roles       []AccessControlRole       `json:"roles"`
	Menus       []AccessControlMenu       `json:"menus"`
	Permissions []AccessControlPermission `json:"permissions"`
}

type AccessControlMenu struct {
	ID        int64  `json:"id"`
	ParentID  int64  `json:"parentId,omitempty"`
	Name      string `json:"name"`
	Code      string `json:"code"`
	Route     string `json:"route,omitempty"`
	Icon      string `json:"icon,omitempty"`
	SortOrder int64  `json:"sortOrder"`
	IsActive  bool   `json:"isActive"`
}

type AccessControlMutationResponse struct {
	Success bool `json:"success"`
}

type AccessControlPermission struct {
	ID          int64  `json:"id"`
	MenuID      int64  `json:"menuId"`
	Code        string `json:"code"`
	Action      string `json:"action"`
	Description string `json:"description,omitempty"`
}

type AccessControlPermissionAssignmentsResponse struct {
	PermissionIDs []int64 `json:"permissionIds"`
}

type AccessControlRole struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"isActive"`
}

type AccessControlRoleAssignmentsResponse struct {
	RoleIDs []int64 `json:"roleIds"`
}

type AccessControlUserResponse struct {
	ID       int64   `json:"id"`
	Email    string  `json:"email"`
	IsActive bool    `json:"isActive"`
	RoleIDs  []int64 `json:"roleIds"`
}

type CatalogBrand struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	LogoURL  string `json:"logoUrl,omitempty"`
	IsActive bool   `json:"isActive"`
}

type CatalogBrandListResponse struct {
	Brands []CatalogBrand `json:"brands"`
}

type CatalogCategory struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	ParentID  *int64 `json:"parentId,omitempty"`
	ImageURL  string `json:"imageUrl,omitempty"`
	SortOrder int64  `json:"sortOrder"`
	IsActive  bool   `json:"isActive"`
}

type CatalogCategoryListResponse struct {
	Categories []CatalogCategory `json:"categories"`
}

type CatalogRoom struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	ImageCover string `json:"imageCover"`
	SortOrder  int64  `json:"sortOrder"`
	IsActive   bool   `json:"isActive"`
}

type CatalogRoomListResponse struct {
	Rooms []CatalogRoom `json:"rooms"`
}

type CatalogMutationResponse struct {
	Success bool `json:"success"`
}

type CatalogProduct struct {
	ID              int64   `json:"id"`
	SKU             string  `json:"sku"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	CategoryID      int64   `json:"categoryId"`
	BrandID         *int64  `json:"brandId,omitempty"`
	Description     string  `json:"description,omitempty"`
	Specifications  string  `json:"specifications,omitempty"`
	DatasheetPDFURL string  `json:"datasheetPdfUrl,omitempty"`
	PriceGeneral    float64 `json:"priceGeneral"`
	UnitSale        string  `json:"unitSale"`
	MinOrder        int64   `json:"minOrder"`
	StockStatus     string  `json:"stockStatus"`
	StockQtyLabel   string  `json:"stockQtyLabel,omitempty"`
	IsFeatured      bool    `json:"isFeatured"`
	IsActive        bool    `json:"isActive"`
	MetaTitle       string  `json:"metaTitle,omitempty"`
	MetaDescription string  `json:"metaDescription,omitempty"`
	Images          []struct {
		ID        int64  `json:"id"`
		ImageURL  string `json:"image_url"`
		IsPrimary bool   `json:"is_primary"`
		SortOrder int64  `json:"sort_order"`
	} `json:"images,omitempty"`
	PrimaryImageURL string `json:"primaryImageUrl,omitempty"`
}

type CatalogProductListResponse struct {
	Products []CatalogProduct `json:"products"`
}

type CreateBrandRequest struct {
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	LogoURL string `json:"logoUrl,omitempty"`
}

type CreateCategoryRequest struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	ParentID  *int64 `json:"parentId,omitempty"`
	ImageURL  string `json:"imageUrl,omitempty"`
	SortOrder int64  `json:"sortOrder"`
}

type CreateProductRequest struct {
	SKU             string  `json:"sku"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	CategoryID      int64   `json:"categoryId"`
	BrandID         *int64  `json:"brandId,omitempty"`
	Description     string  `json:"description,omitempty"`
	Specifications  string  `json:"specifications,omitempty"`
	DatasheetPDFURL string  `json:"datasheetPdfUrl,omitempty"`
	PriceGeneral    float64 `json:"priceGeneral"`
	UnitSale        string  `json:"unitSale"`
	MinOrder        int64   `json:"minOrder"`
	StockStatus     string  `json:"stockStatus"`
	StockQtyLabel   string  `json:"stockQtyLabel,omitempty"`
	IsFeatured      bool    `json:"isFeatured"`
	IsActive        bool    `json:"isActive"`
	MetaTitle       string  `json:"metaTitle,omitempty"`
	MetaDescription string  `json:"metaDescription,omitempty"`
}

type CreateRoleRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CreateRoomRequest struct {
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	ImageCover string `json:"imageCover"`
	SortOrder  int64  `json:"sortOrder"`
}

type HealthResp struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type ReplaceRolePermissionsRequest struct {
	PermissionIDs []int64 `json:"permissionIds"`
}

type ReplaceUserRolesRequest struct {
	RoleIDs []int64 `json:"roleIds"`
}

type SearchAccessControlUsersRequest struct {
	Email string `form:"email"`
}

type UpdateBrandRequest struct {
	Name     *string `json:"name,omitempty"`
	Slug     *string `json:"slug,omitempty"`
	LogoURL  *string `json:"logoUrl,omitempty"`
	IsActive *bool   `json:"isActive,omitempty"`
}

type UpdateCategoryRequest struct {
	Name      *string `json:"name,omitempty"`
	Slug      *string `json:"slug,omitempty"`
	ParentID  *int64  `json:"parentId,omitempty"`
	ImageURL  *string `json:"imageUrl,omitempty"`
	SortOrder *int64  `json:"sortOrder,omitempty"`
	IsActive  *bool   `json:"isActive,omitempty"`
}

type UpdateProductRequest struct {
	SKU             *string  `json:"sku,omitempty"`
	Name            *string  `json:"name,omitempty"`
	Slug            *string  `json:"slug,omitempty"`
	CategoryID      *int64   `json:"categoryId,omitempty"`
	BrandID         *int64   `json:"brandId,omitempty"`
	Description     *string  `json:"description,omitempty"`
	Specifications  *string  `json:"specifications,omitempty"`
	DatasheetPDFURL *string  `json:"datasheetPdfUrl,omitempty"`
	PriceGeneral    *float64 `json:"priceGeneral,omitempty"`
	UnitSale        *string  `json:"unitSale,omitempty"`
	MinOrder        *int64   `json:"minOrder,omitempty"`
	StockStatus     *string  `json:"stockStatus,omitempty"`
	StockQtyLabel   *string  `json:"stockQtyLabel,omitempty"`
	IsFeatured      *bool    `json:"isFeatured,omitempty"`
	IsActive        *bool    `json:"isActive,omitempty"`
	MetaTitle       *string  `json:"metaTitle,omitempty"`
	MetaDescription *string  `json:"metaDescription,omitempty"`
}

type ReplaceProductRoomsRequest struct {
	RoomIDs []int64 `json:"roomIds"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsActive    bool   `json:"isActive"`
}

type UpdateRoomRequest struct {
	Name       *string `json:"name,omitempty"`
	Slug       *string `json:"slug,omitempty"`
	ImageCover *string `json:"imageCover,omitempty"`
	SortOrder  *int64  `json:"sortOrder,omitempty"`
	IsActive   *bool   `json:"isActive,omitempty"`
}

type CurrentAccessUser struct {
	ID    int64    `json:"id"`
	Email string   `json:"email"`
	Roles []string `json:"roles"`
}

type UserAccessResponse struct {
	User            CurrentAccessUser   `json:"user"`
	Menus           []AccessControlMenu `json:"menus"`
	PermissionCodes []string            `json:"permissionCodes"`
}