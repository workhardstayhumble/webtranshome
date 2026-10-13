package permission

// Permission selalu diverifikasi terhadap role_permissions di PostgreSQL.
// Role tidak otomatis menjadi bypass; SUPER_ADMIN menerima izin lewat migration.
const (
	AccessControlManageCode   = "rbac.manage"
	CatalogCategoryManageCode = "catalog.category.manage" // legacy
	CatalogBrandManageCode    = "catalog.brand.manage"    // legacy
	CatalogRoomManageCode     = "catalog.room.manage"     // legacy

	CatalogCategoryReadCode   = "catalog.category.read"
	CatalogCategoryCreateCode = "catalog.category.create"
	CatalogCategoryUpdateCode = "catalog.category.update"
	CatalogBrandReadCode      = "catalog.brand.read"
	CatalogBrandCreateCode    = "catalog.brand.create"
	CatalogBrandUpdateCode    = "catalog.brand.update"
	CatalogRoomReadCode       = "catalog.room.read"
	CatalogRoomCreateCode     = "catalog.room.create"
	CatalogRoomUpdateCode     = "catalog.room.update"

	CatalogProductReadCode        = "catalog.product.read"
	CatalogProductCreateCode      = "catalog.product.create"
	CatalogProductUpdateCode      = "catalog.product.update"
	CatalogProductPriceUpdateCode = "catalog.product.price.update"
	CatalogProductStockUpdateCode = "catalog.product.stock.update"
	CatalogProductFeaturedCode    = "catalog.featured.manage"
	ContentBannerReadCode         = "content.banner.read"
	ContentBannerManageCode       = "content.banner.manage"
	ContentArticleReadCode        = "content.article.read"
	ContentArticleManageCode      = "content.article.manage"
	MembershipMemberReadCode      = "membership.member.read"
	MembershipMemberUpdateCode    = "membership.member.update"
	MembershipPointsReadCode      = "membership.points.read"
	MembershipPointsInputCode     = "membership.points.input"
	IntegrationQontakReadCode     = "integration.qontak.read"
	IntegrationQontakRetryCode    = "integration.qontak.retry"
	AuditReadCode                 = "audit.read"
)

type Permission struct {
	ID          int64
	MenuID      int64
	Code        string
	Action      string
	Description *string
}