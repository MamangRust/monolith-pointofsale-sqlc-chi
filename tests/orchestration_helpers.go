package tests

import (
	"bytes"
	"context"
	"github.com/MamangRust/monolith-point-of-sale-pb"
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pbmerchants_document "github.com/MamangRust/monolith-point-of-sale-pb/merchants_document"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserrole "github.com/MamangRust/monolith-point-of-sale-pb/user_role"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"mime/multipart"

	"github.com/MamangRust/monolith-point-of-sale-pkg/auth"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"gorm.io/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Role
	role_cache "github.com/MamangRust/monolith-point-of-sale-role/cache"
	role_handler "github.com/MamangRust/monolith-point-of-sale-role/handler"
	role_repo "github.com/MamangRust/monolith-point-of-sale-role/repository"
	role_service "github.com/MamangRust/monolith-point-of-sale-role/service"

	// User
	user_cache "github.com/MamangRust/monolith-point-of-sale-user/cache"
	user_handler "github.com/MamangRust/monolith-point-of-sale-user/handler"
	user_repo "github.com/MamangRust/monolith-point-of-sale-user/repository"
	user_service "github.com/MamangRust/monolith-point-of-sale-user/service"

	// Auth
	auth_cache "github.com/MamangRust/monolith-point-of-sale-auth/cache"
	auth_handler "github.com/MamangRust/monolith-point-of-sale-auth/handler"
	auth_repo "github.com/MamangRust/monolith-point-of-sale-auth/repository"
	auth_service "github.com/MamangRust/monolith-point-of-sale-auth/service"

	// Category
	category_cache "github.com/MamangRust/monolith-point-of-sale-category/cache"
	category_handler "github.com/MamangRust/monolith-point-of-sale-category/handler"
	category_repo "github.com/MamangRust/monolith-point-of-sale-category/repository"
	category_service "github.com/MamangRust/monolith-point-of-sale-category/service"

	// Product
	product_cache "github.com/MamangRust/monolith-point-of-sale-product/cache"
	product_handler "github.com/MamangRust/monolith-point-of-sale-product/handler"
	product_repo "github.com/MamangRust/monolith-point-of-sale-product/repository"
	product_service "github.com/MamangRust/monolith-point-of-sale-product/service"

	// Merchant
	merchant_cache "github.com/MamangRust/monolith-point-of-sale-merchant/cache"
	merchant_handler "github.com/MamangRust/monolith-point-of-sale-merchant/handler"
	merchant_repo "github.com/MamangRust/monolith-point-of-sale-merchant/repository"
	merchant_service "github.com/MamangRust/monolith-point-of-sale-merchant/service"

	// Order
	order_cache "github.com/MamangRust/monolith-point-of-sale-order/cache"
	order_handler "github.com/MamangRust/monolith-point-of-sale-order/handler"
	order_repo "github.com/MamangRust/monolith-point-of-sale-order/repository"
	order_service "github.com/MamangRust/monolith-point-of-sale-order/service"

	// Transaction
	transaction_cache "github.com/MamangRust/monolith-point-of-sale-transacton/cache"
	transaction_handler "github.com/MamangRust/monolith-point-of-sale-transacton/handler"
	transaction_repo "github.com/MamangRust/monolith-point-of-sale-transacton/repository"
	transaction_service "github.com/MamangRust/monolith-point-of-sale-transacton/service"

	// Order Item
	order_item_cache "github.com/MamangRust/monolith-point-of-sale-order-item/cache"
	order_item_handler "github.com/MamangRust/monolith-point-of-sale-order-item/handler"
	order_item_repo "github.com/MamangRust/monolith-point-of-sale-order-item/repository"
	order_item_service "github.com/MamangRust/monolith-point-of-sale-order-item/service"

	// Cashier
	mencache "github.com/MamangRust/monolith-point-of-sale-cashier/cache"
	cashier_handler "github.com/MamangRust/monolith-point-of-sale-cashier/handler"
	cashier_repo "github.com/MamangRust/monolith-point-of-sale-cashier/repository"
	cashier_service "github.com/MamangRust/monolith-point-of-sale-cashier/service"
)

// DependencyClients bundles generated gRPC clients for the identity and
// catalog services that consumers now reach through the shared adapters.
type DependencyClients struct {
	UserQuery     pbusers.UserQueryServiceClient
	UserCommand   pbusers.UserCommandServiceClient
	RoleQuery     pbroles.RoleQueryServiceClient
	RoleCommand   pbroles.RoleCommandServiceClient
	UserRole      pbuserrole.UserRoleServiceClient
	CategoryQuery pbcategories.CategoryQueryServiceClient
	MerchantQuery pbmerchants.MerchantQueryServiceClient
}

// StartDependencyClients boots user, role, category and merchant services
// in-process against the given database and returns clients for them plus a
// cleanup function. It is used by the plain tests.TestSuite suites that build
// adapters directly instead of going through BaseTestSuite orchestration.
// dependencyCleanups holds the in-process servers started by
// StartDependencyClients so a suite can release them in TearDownSuite.
var dependencyCleanups []func()

// StopDependencyClients shuts down every server started via
// StartDependencyClients. Safe to call when nothing was started.
func StopDependencyClients() {
	for _, cleanup := range dependencyCleanups {
		cleanup()
	}
	dependencyCleanups = nil
}

func StartDependencyClients(
	gormDB *gorm.DB,
	cacheStore *cache.CacheStore,
	log logger.LoggerInterface,
	obs observability.TraceLoggerObservability,
) *DependencyClients {
	var servers []*grpc.Server
	var conns []*grpc.ClientConn

	start := func(register func(*grpc.Server)) *grpc.ClientConn {
		server := grpc.NewServer()
		register(server)
		addr, err := RunGRPCServer(server)
		if err != nil {
			panic(err)
		}
		servers = append(servers, server)
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		conns = append(conns, conn)
		return conn
	}

	hasher := hash.NewHashingPassword()

	roleConn := start(func(server *grpc.Server) {
		roleRepos := role_repo.NewRepositories(gormDB)
		roleSvc := role_service.NewService(&role_service.Deps{
			Repositories:  roleRepos,
			Logger:        log,
			Mencache:      role_cache.NewMencache(cacheStore),
			Observability: obs,
		})
		roleGapi := role_handler.NewHandler(&role_handler.Deps{Service: roleSvc, Logger: log})
		pbroles.RegisterRoleQueryServiceServer(server, roleGapi.Role)
		pbroles.RegisterRoleCommandServiceServer(server, roleGapi.RoleCommand)
		pbuserrole.RegisterUserRoleServiceServer(server, roleGapi.UserRole)
	})

	userConn := start(func(server *grpc.Server) {
		userRepos := user_repo.NewRepositories(&user_repo.Deps{
			Db:       gormDB,
			Role:     pbroles.NewRoleQueryServiceClient(roleConn),
			UserRole: pbuserrole.NewUserRoleServiceClient(roleConn),
		})
		userSvc := user_service.NewService(&user_service.Deps{
			Repositories:  userRepos,
			Logger:        log,
			Hash:          hasher,
			Mencache:      user_cache.NewMencache(cacheStore),
			Observability: obs,
		})
		userGapi := user_handler.NewHandler(&user_handler.Deps{Service: userSvc, Logger: log})
		pbusers.RegisterUserQueryServiceServer(server, userGapi.User)
		pbusers.RegisterUserCommandServiceServer(server, userGapi.UserCommand)
	})

	categoryConn := start(func(server *grpc.Server) {
		catRepos := category_repo.NewRepositories(gormDB)
		catSvc := category_service.NewService(&category_service.Deps{
			Mencache:      category_cache.NewMencache(cacheStore),
			Repositories:  catRepos,
			Logger:        log,
			Observability: obs,
		})
		catGapi := category_handler.NewHandler(&category_handler.Deps{Service: catSvc, Logger: log})
		pbcategories.RegisterCategoryQueryServiceServer(server, catGapi.Category)
		pbcategories.RegisterCategoryCommandServiceServer(server, catGapi.CategoryCommand)
	})

	merchantConn := start(func(server *grpc.Server) {
		merchantRepos := merchant_repo.NewRepositories(gormDB, pbusers.NewUserQueryServiceClient(userConn))
		merchantSvc := merchant_service.NewService(&merchant_service.Deps{
			Mencache:      merchant_cache.NewMencache(cacheStore),
			Repositories:  merchantRepos,
			Logger:        log,
			Observability: obs,
			Kafka:         nil,
		})
		merchantGapi := merchant_handler.NewHandler(&merchant_handler.Deps{Service: merchantSvc, Logger: log})
		pbmerchants.RegisterMerchantQueryServiceServer(server, merchantGapi.Merchant)
		pbmerchants.RegisterMerchantCommandServiceServer(server, merchantGapi.MerchantCommand)
	})

	clients := &DependencyClients{
		UserQuery:     pbusers.NewUserQueryServiceClient(userConn),
		UserCommand:   pbusers.NewUserCommandServiceClient(userConn),
		RoleQuery:     pbroles.NewRoleQueryServiceClient(roleConn),
		RoleCommand:   pbroles.NewRoleCommandServiceClient(roleConn),
		UserRole:      pbuserrole.NewUserRoleServiceClient(roleConn),
		CategoryQuery: pbcategories.NewCategoryQueryServiceClient(categoryConn),
		MerchantQuery: pbmerchants.NewMerchantQueryServiceClient(merchantConn),
	}

	dependencyCleanups = append(dependencyCleanups, func() {
		for _, conn := range conns {
			conn.Close()
		}
		for _, server := range servers {
			server.GracefulStop()
		}
	})

	return clients
}

func (s *BaseTestSuite) SetupRoleService() {
	cacheStore := s.GetCacheStore()
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(gormDB)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repositories:  roleRepos,
		Logger:        s.Log,
		Mencache:      roleMencache,
		Observability: s.Obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbroles.RegisterRoleQueryServiceServer(server, roleGapi.Role)
	pbroles.RegisterRoleCommandServiceServer(server, roleGapi.RoleCommand)
	pbuserrole.RegisterUserRoleServiceServer(server, roleGapi.UserRole)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["role"] = conn
	s.Servers = append(s.Servers, server)
}

// UserRepositories builds the user repositories with the shared role and
// user-role adapters pointed at the in-process role service, starting it if it
// is not already running.
func (s *BaseTestSuite) UserRepositories() user_repo.Repositories {
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	if _, ok := s.Conns["role"]; !ok {
		s.SetupRoleService()
	}

	return user_repo.NewRepositories(&user_repo.Deps{
		Db:       gormDB,
		Role:     pbroles.NewRoleQueryServiceClient(s.Conns["role"]),
		UserRole: pbuserrole.NewUserRoleServiceClient(s.Conns["role"]),
	})
}

func (s *BaseTestSuite) SetupUserService() {
	cacheStore := s.GetCacheStore()
	hasher := hash.NewHashingPassword()

	userMencache := user_cache.NewMencache(cacheStore)
	userRepos := s.UserRepositories()
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        s.Log,
		Hash:          hasher,
		Mencache:      userMencache,
		Observability: s.Obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbusers.RegisterUserQueryServiceServer(server, userGapi.User)
	pbusers.RegisterUserCommandServiceServer(server, userGapi.UserCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["user"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupAuthService() {
	cacheStore := s.GetCacheStore()
	gormDB, err := s.GormDB()
	s.Require().NoError(err)
	hasher := hash.NewHashingPassword()
	tokenManager, _ := auth.NewManager("mysecret")

	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}
	if _, ok := s.Conns["role"]; !ok {
		s.SetupRoleService()
	}

	userQueryClient := pbusers.NewUserQueryServiceClient(s.Conns["user"])
	userCommandClient := pbusers.NewUserCommandServiceClient(s.Conns["user"])
	roleQueryClient := pbroles.NewRoleQueryServiceClient(s.Conns["role"])
	roleCommandClient := pbroles.NewRoleCommandServiceClient(s.Conns["role"])

	authRepos := auth_repo.NewRepositories(gormDB, userQueryClient, userCommandClient, roleQueryClient, roleCommandClient, pbuserrole.NewUserRoleServiceClient(s.Conns["role"]))
	authMencache := auth_cache.NewMencache(cacheStore)
	authSvc := auth_service.NewService(&auth_service.Deps{
		Repositories:  authRepos,
		Logger:        s.Log,
		Mencache:      authMencache,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: s.Obs,
	})
	authGapi := auth_handler.NewAuthHandleGrpc(authSvc, s.Log)
	server := grpc.NewServer()
	pb.RegisterAuthServiceServer(server, authGapi)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.Conns["auth"] = conn
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCategoryService() {
	cacheStore := s.GetCacheStore()
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	catMencache := category_cache.NewMencache(cacheStore)
	catRepos := category_repo.NewRepositories(gormDB)
	catSvc := category_service.NewService(&category_service.Deps{
		Mencache:      catMencache,
		Repositories:  catRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	catGapi := category_handler.NewHandler(&category_handler.Deps{
		Service: catSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbcategories.RegisterCategoryQueryServiceServer(server, catGapi.Category)
	pbcategories.RegisterCategoryCommandServiceServer(server, catGapi.CategoryCommand)
	pbcategories.RegisterCategoryStatsServiceServer(server, catGapi.CategoryStats)
	pbcategories.RegisterCategoryStatsByIdServiceServer(server, catGapi.CategoryStatsById)
	pbcategories.RegisterCategoryStatsByMerchantServiceServer(server, catGapi.CategoryStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["category"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupProductService() {
	cacheStore := s.GetCacheStore()
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	if _, ok := s.Conns["category"]; !ok {
		s.SetupCategoryService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}

	categoryClient := pbcategories.NewCategoryQueryServiceClient(s.Conns["category"])
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"])

	prodMencache := product_cache.NewMencache(cacheStore)
	prodRepos := product_repo.NewRepositories(gormDB, categoryClient, merchantClient)
	prodSvc := product_service.NewService(&product_service.Deps{
		Mencache:      prodMencache,
		Repositories:  prodRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	prodGapi := product_handler.NewHandler(&product_handler.Deps{
		Service: prodSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbproducts.RegisterProductQueryServiceServer(server, prodGapi.Product)
	pbproducts.RegisterProductCommandServiceServer(server, prodGapi.ProductCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["product"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupMerchantService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}

	cacheStore := s.GetCacheStore()
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	merchantMencache := merchant_cache.NewMencache(cacheStore)
	userQueryClient := pbusers.NewUserQueryServiceClient(s.Conns["user"])
	merchantRepos := merchant_repo.NewRepositories(gormDB, userQueryClient)
	merchantSvc := merchant_service.NewService(&merchant_service.Deps{
		Mencache:      merchantMencache,
		Repositories:  merchantRepos,
		Logger:        s.Log,
		Observability: s.Obs,
		Kafka:         nil,
	})
	merchantGapi := merchant_handler.NewHandler(&merchant_handler.Deps{
		Service: merchantSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbmerchants.RegisterMerchantQueryServiceServer(server, merchantGapi.Merchant)
	pbmerchants.RegisterMerchantCommandServiceServer(server, merchantGapi.MerchantCommand)
	pbmerchants_document.RegisterMerchantDocumentQueryServiceServer(server, merchantGapi.MerchantDocument)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["merchant"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderService() {
	if _, ok := s.Conns["cashier"]; !ok {
		s.SetupCashierService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}
	if _, ok := s.Conns["product"]; !ok {
		s.SetupProductService()
	}
	if _, ok := s.Conns["order-item"]; !ok {
		s.SetupOrderItemService()
	}

	cacheStore := s.GetCacheStore()
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	cashierClient := pbcashiers.NewCashierQueryServiceClient(s.Conns["cashier"])
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"])
	productClient := pbproducts.NewProductQueryServiceClient(s.Conns["product"])
	productCommandClient := pbproducts.NewProductCommandServiceClient(s.Conns["product"])
	orderItemClient := pborder_items.NewOrderItemQueryServiceClient(s.Conns["order-item"])
	orderItemCommandClient := pborder_items.NewOrderItemCommandServiceClient(s.Conns["order-item"])

	orderMencache := order_cache.NewMencache(cacheStore)
	orderRepos := order_repo.NewRepositories(gormDB, cashierClient, merchantClient, productClient, productCommandClient, orderItemClient, orderItemCommandClient)
	orderSvc := order_service.NewService(&order_service.Deps{
		Mencache:      orderMencache,
		Repositories:  orderRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	orderGapi := order_handler.NewHandler(&order_handler.Deps{
		Service: orderSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pborders.RegisterOrderQueryServiceServer(server, orderGapi.Order)
	pborders.RegisterOrderCommandServiceServer(server, orderGapi.OrderCommand)
	pborders.RegisterOrderStatsServiceServer(server, orderGapi.OrderStats)
	pborders.RegisterOrderStatsByIdServiceServer(server, orderGapi.OrderStatsById)
	pborders.RegisterOrderStatsByMerchantServiceServer(server, orderGapi.OrderStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupTransactionService() {
	if _, ok := s.Conns["cashier"]; !ok {
		s.SetupCashierService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}
	if _, ok := s.Conns["order"]; !ok {
		s.SetupOrderService()
	}
	if _, ok := s.Conns["order-item"]; !ok {
		s.SetupOrderItemService()
	}

	cacheStore := s.GetCacheStore()
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	cashierClient := pbcashiers.NewCashierQueryServiceClient(s.Conns["cashier"])
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"])
	orderClient := pborders.NewOrderQueryServiceClient(s.Conns["order"])
	orderItemClient := pborder_items.NewOrderItemQueryServiceClient(s.Conns["order-item"])

	transactionMencache := transaction_cache.NewMencache(cacheStore)
	transactionRepos := transaction_repo.NewRepositories(gormDB, cashierClient, merchantClient, orderClient, orderItemClient)
	transactionSvc := transaction_service.NewService(&transaction_service.Deps{
		Mencache:      transactionMencache,
		Repositories:  transactionRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	transactionGapi := transaction_handler.NewHandler(&transaction_handler.Deps{
		Service: transactionSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbtransactions.RegisterTransactionQueryServiceServer(server, transactionGapi.Transaction)
	pbtransactions.RegisterTransactionCommandServiceServer(server, transactionGapi.TransactionCommand)
	pbtransactions.RegisterTransactionStatsServiceServer(server, transactionGapi.TransactionStats)
	pbtransactions.RegisterTransactionStatsByMerchantServiceServer(server, transactionGapi.TransactionStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["transaction"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupOrderItemService() {
	cacheStore := s.GetCacheStore()
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	itemMencache := order_item_cache.NewMencache(cacheStore)
	itemRepos := order_item_repo.NewRepositories(gormDB)
	itemSvc := order_item_service.NewService(&order_item_service.Deps{
		Mencache:      itemMencache,
		Repositories:  itemRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	itemGapi := order_item_handler.NewHandler(&order_item_handler.Deps{
		Service: itemSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pborder_items.RegisterOrderItemQueryServiceServer(server, itemGapi.OrderItem)
	pborder_items.RegisterOrderItemCommandServiceServer(server, itemGapi.OrderItemCommand)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["order-item"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) SetupCashierService() {
	if _, ok := s.Conns["user"]; !ok {
		s.SetupUserService()
	}
	if _, ok := s.Conns["merchant"]; !ok {
		s.SetupMerchantService()
	}

	cacheStore := s.GetCacheStore()
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	userClient := pbusers.NewUserQueryServiceClient(s.Conns["user"])
	merchantClient := pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"])

	cashierMencache := mencache.NewMencache(cacheStore)
	cashierRepos := cashier_repo.NewRepositories(gormDB, userClient, merchantClient)
	cashierSvc := cashier_service.NewService(&cashier_service.Deps{
		Ctx:           context.Background(),
		Mencache:      cashierMencache,
		Repositories:  cashierRepos,
		Logger:        s.Log,
		Observability: s.Obs,
	})
	cashierGapi := cashier_handler.NewHandler(&cashier_handler.Deps{
		Service: cashierSvc,
		Logger:  s.Log,
	})
	server := grpc.NewServer()
	pbcashiers.RegisterCashierQueryServiceServer(server, cashierGapi.Cashier)
	pbcashiers.RegisterCashierCommandServiceServer(server, cashierGapi.CashierCommand)
	pbcashiers.RegisterCashierStatsServiceServer(server, cashierGapi.CashierStats)
	pbcashiers.RegisterCashierStatsByIdServiceServer(server, cashierGapi.CashierStatsById)
	pbcashiers.RegisterCashierStatsByMerchantServiceServer(server, cashierGapi.CashierStatsByMerchant)
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Conns["cashier"] = s.dial(addr)
	s.Servers = append(s.Servers, server)
}

func (s *BaseTestSuite) dial(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	return conn
}

func (s *BaseTestSuite) GetCacheStore() *cache.CacheStore {
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	return cache.NewCacheStore(s.ts.RedisClient(), s.Log, cacheMetrics)
}

func (s *BaseTestSuite) BuildMultipartRequestBody(fields map[string]string, fieldName, fileName string) ([]byte, string) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range fields {
		fw, _ := w.CreateFormField(key)
		fw.Write([]byte(r))
	}
	fw, _ := w.CreateFormFile(fieldName, fileName)
	fw.Write([]byte("dummy image content"))
	w.Close()
	return b.Bytes(), w.FormDataContentType()
}
