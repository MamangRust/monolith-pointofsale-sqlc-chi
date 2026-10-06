package tests

import (
	"context"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"reflect"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type BaseTestSuite struct {
	suite.Suite
	ts      *TestSuite
	Log     logger.LoggerInterface
	Obs     observability.TraceLoggerObservability
	Conns   map[string]*grpc.ClientConn
	Servers []*grpc.Server
	Ctx     context.Context
	Cancel  context.CancelFunc
}

func (s *BaseTestSuite) SetupSuite() {
	s.Ctx, s.Cancel = context.WithCancel(context.Background())

	ts, err := SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	s.Log, _ = logger.NewLogger("test", nil)

	if s.Log == nil || (reflect.ValueOf(s.Log).Kind() == reflect.Ptr && reflect.ValueOf(s.Log).IsNil()) {
		z, _ := zap.NewDevelopment()
		s.Log = &logger.Logger{Log: z}
	}

	s.Obs, err = observability.NewObservability("test", s.Log)
	s.Require().NoError(err)
	s.Require().NotNil(s.Obs)
	s.Conns = make(map[string]*grpc.ClientConn)
}

func (s *BaseTestSuite) TearDownSuite() {
	for _, conn := range s.Conns {
		conn.Close()
	}
	for _, server := range s.Servers {
		server.GracefulStop()
	}
	if s.ts != nil {
		s.ts.Teardown()
	}
	if s.Cancel != nil {
		s.Cancel()
	}
}

func (s *BaseTestSuite) GormDB() (*gorm.DB, error) {
	return s.ts.GormDB()
}

func (s *BaseTestSuite) RedisClient() *goredis.Client {
	return s.ts.RedisClient()
}

func (s *BaseTestSuite) RegisterServer(server *grpc.Server) string {
	addr, err := RunGRPCServer(server)
	s.Require().NoError(err)
	s.Servers = append(s.Servers, server)
	return addr
}

func (s *BaseTestSuite) GetConnection(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	s.Require().NoError(err)
	return conn
}

func (s *BaseTestSuite) SeedUser(ctx context.Context) int {
	db, err := s.GormDB()
	s.Require().NoError(err)
	err = db.WithContext(ctx).Exec(`
		INSERT INTO roles (role_name, created_at, updated_at) 
		VALUES ('Admin Access 1', current_timestamp, current_timestamp),
		       ('ROLE_ADMIN', current_timestamp, current_timestamp)
		ON CONFLICT (role_name) DO NOTHING
	`).Error
	s.Require().NoError(err)

	res, err := pbusers.NewUserCommandServiceClient(s.Conns["user"]).Create(ctx, &pbusers.CreateUserRequest{
		Firstname:       "Seed",
		Lastname:        "User",
		Email:           "seed.user@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedCategory(ctx context.Context) int {
	res, err := pbcategories.NewCategoryCommandServiceClient(s.Conns["category"]).Create(ctx, &pbcategories.CreateCategoryRequest{
		Name:        "Seed Category",
		Description: "Seed Description",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedMerchant(ctx context.Context, userID int) int {
	res, err := pbmerchants.NewMerchantCommandServiceClient(s.Conns["merchant"]).Create(ctx, &pbmerchants.CreateMerchantRequest{
		UserId:       int32(userID),
		Name:         "Seed Merchant",
		Description:  "Seed Description",
		Address:      "Seed Address",
		ContactEmail: "merchant@example.com",
		ContactPhone: "08123456789",
		Status:       "active",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedProduct(ctx context.Context, merchantID int, categoryID int) int {
	res, err := pbproducts.NewProductCommandServiceClient(s.Conns["product"]).Create(ctx, &pbproducts.CreateProductRequest{
		MerchantId:   int32(merchantID),
		CategoryId:   int32(categoryID),
		Name:         "Seed Product",
		Description:  "Seed Description",
		Price:        10000,
		CountInStock: 100,
		Brand:        "Seed Brand",
		Weight:       1000,
		ImageProduct: "seed.jpg",
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedOrder(ctx context.Context, userID int, merchID int, prodID int) int {
	db, err := s.GormDB()
	s.Require().NoError(err)

	var cashierID int32
	err = db.WithContext(ctx).
		Raw(`SELECT cashier_id FROM cashiers WHERE user_id = ? AND merchant_id = ? AND deleted_at IS NULL LIMIT 1`, userID, merchID).
		Scan(&cashierID).Error
	if err != nil || cashierID == 0 {
		err = db.WithContext(ctx).
			Raw(`INSERT INTO cashiers (merchant_id, user_id, name, created_at, updated_at)
			VALUES (?, ?, 'Seed Cashier', current_timestamp, current_timestamp)
			RETURNING cashier_id`, merchID, userID).
			Scan(&cashierID).Error
		s.Require().NoError(err)
	}

	res, err := pborders.NewOrderCommandServiceClient(s.Conns["order"]).Create(ctx, &pborders.CreateOrderRequest{
		MerchantId: int32(merchID),
		CashierId:  int32(cashierID),
		Items: []*pborders.CreateOrderItemRequest{
			{
				ProductId: int32(prodID),
				Quantity:  1,
			},
		},
	})
	s.Require().NoError(err)
	return int(res.Data.Id)
}

func (s *BaseTestSuite) SeedOrderItem(ctx context.Context, orderID int, productID int) int {
	db, err := s.GormDB()
	s.Require().NoError(err)

	var id int32
	err = db.WithContext(ctx).
		Raw(`INSERT INTO "order_items" (order_id, product_id, quantity, price) VALUES (?, ?, 1, 1000) RETURNING order_item_id`, orderID, productID).
		Scan(&id).Error
	s.Require().NoError(err)
	return int(id)
}
