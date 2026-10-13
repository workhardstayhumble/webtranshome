package svc

import (
	"webtranshome/internal/config"
	"webtranshome/internal/middleware"
	"webtranshome/internal/model"
	catalogcommand "webtranshome/internal/modules/catalog/application/command"
	catalogquery "webtranshome/internal/modules/catalog/application/query"
	catalogpostgres "webtranshome/internal/modules/catalog/infrastructure/postgres"
	"webtranshome/internal/modules/identityaccess/application/command"
	"webtranshome/internal/modules/identityaccess/application/query"
	passwordbcrypt "webtranshome/internal/modules/identityaccess/infrastructure/password/bcrypt"
	authpostgres "webtranshome/internal/modules/identityaccess/infrastructure/postgres"
	jwtissuer "webtranshome/internal/modules/identityaccess/infrastructure/token/jwtissuer"

	"github.com/zeromicro/go-zero/core/stores/postgres"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	DB                       sqlx.SqlConn
	Config                   config.Config
	UserModel                model.UsersModel
	LoginAdmin               *command.Login
	AccessControlQueries     *query.AccessControl
	ReplaceAccessAssignments *command.ReplaceAccessAssignments
	ManageRoles              *command.ManageRoles
	ListCategories           *catalogquery.ListCategories
	GetCategory              *catalogquery.GetCategory
	CreateCategory           *catalogcommand.CreateCategory
	UpdateCategory           *catalogcommand.UpdateCategory
	ListBrands               *catalogquery.ListBrands
	GetBrand                 *catalogquery.GetBrand
	CreateBrand              *catalogcommand.CreateBrand
	UpdateBrand              *catalogcommand.UpdateBrand
	ListRooms                *catalogquery.ListRooms
	GetRoom                  *catalogquery.GetRoom
	CreateRoom               *catalogcommand.CreateRoom
	UpdateRoom               *catalogcommand.UpdateRoom
	ListProducts             *catalogquery.ListProducts
	GetProduct               *catalogquery.GetProduct
	CreateProduct            *catalogcommand.CreateProduct
	UpdateProduct            *catalogcommand.UpdateProduct
	ReplaceProductRooms      *catalogcommand.ReplaceProductRooms
	Authenticate             rest.Middleware
	RequirePermission        func(string) rest.Middleware
	CheckPermission          *query.CheckPermission
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := postgres.New(c.Database.DataSource)
	userModel := model.NewUsersModel(conn)
	userRepository := authpostgres.NewUserRepository(conn, userModel)
	permissionRepository := authpostgres.NewPermissionRepository(conn)
	accessControlRepository := authpostgres.NewAccessControlRepository(conn)
	categoryRepository := catalogpostgres.NewCategoryRepository(conn)
	loginAdmin := command.NewLogin(
		userRepository,
		passwordbcrypt.NewVerifier(),
		jwtissuer.NewIssuer(c.Auth.AccessSecret, c.Auth.AccessExpire),
	)
	checkPermission := query.NewCheckPermission(permissionRepository)
	tokenVerifier := jwtissuer.NewVerifier(c.Auth.AccessSecret)
	accessControlQueries := query.NewAccessControl(
		accessControlRepository,
		accessControlRepository,
		accessControlRepository,
	)
	replaceAccessAssignments := command.NewReplaceAccessAssignments(accessControlRepository)
	manageRoles := command.NewManageRoles(accessControlRepository)
	listCategories := catalogquery.NewListCategories(categoryRepository)
	getCategory := catalogquery.NewGetCategory(categoryRepository)
	createCategory := catalogcommand.NewCreateCategory(categoryRepository, categoryRepository)
	updateCategory := catalogcommand.NewUpdateCategory(categoryRepository, categoryRepository)
	brandRepository := catalogpostgres.NewBrandRepository(conn)
	listBrands := catalogquery.NewListBrands(brandRepository)
	getBrand := catalogquery.NewGetBrand(brandRepository)
	createBrand := catalogcommand.NewCreateBrand(brandRepository)
	updateBrand := catalogcommand.NewUpdateBrand(brandRepository, brandRepository)
	roomRepository := catalogpostgres.NewRoomRepository(conn)
	listRooms := catalogquery.NewListRooms(roomRepository)
	getRoom := catalogquery.NewGetRoom(roomRepository)
	createRoom := catalogcommand.NewCreateRoom(roomRepository)
	updateRoom := catalogcommand.NewUpdateRoom(roomRepository, roomRepository)
	productRepository := catalogpostgres.NewProductRepository(conn)
	listProducts := catalogquery.NewListProducts(productRepository)
	getProduct := catalogquery.NewGetProduct(productRepository)
	createProduct := catalogcommand.NewCreateProduct(productRepository, productRepository)
	updateProduct := catalogcommand.NewUpdateProduct(productRepository, productRepository)
	replaceProductRooms := catalogcommand.NewReplaceProductRooms(productRepository, roomRepository, productRepository)

	return &ServiceContext{
		DB:                       conn,
		Config:                   c,
		UserModel:                userModel,
		LoginAdmin:               loginAdmin,
		AccessControlQueries:     accessControlQueries,
		ReplaceAccessAssignments: replaceAccessAssignments,
		ManageRoles:              manageRoles,
		ListCategories:           listCategories,
		GetCategory:              getCategory,
		CreateCategory:           createCategory,
		UpdateCategory:           updateCategory,
		ListBrands:               listBrands,
		GetBrand:                 getBrand,
		CreateBrand:              createBrand,
		UpdateBrand:              updateBrand,
		ListRooms:                listRooms,
		GetRoom:                  getRoom,
		CreateRoom:               createRoom,
		UpdateRoom:               updateRoom,
		ListProducts:             listProducts,
		GetProduct:               getProduct,
		CreateProduct:            createProduct,
		UpdateProduct:            updateProduct,
		ReplaceProductRooms:      replaceProductRooms,
		CheckPermission:          checkPermission,
		Authenticate:             middleware.Authenticate(tokenVerifier),
		RequirePermission: func(code string) rest.Middleware {
			return middleware.RequirePermission(checkPermission, code)
		},
	}
}