package order_test

import (
	"context"
	pbcashiers "github.com/MamangRust/monolith-point-of-sale-pb/cashiers"
	pbmerchants "github.com/MamangRust/monolith-point-of-sale-pb/merchants"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	pbproducts "github.com/MamangRust/monolith-point-of-sale-pb/products"
	"testing"

	order_cache "github.com/MamangRust/monolith-point-of-sale-order/cache"
	order_handler "github.com/MamangRust/monolith-point-of-sale-order/handler"
	order_repo "github.com/MamangRust/monolith-point-of-sale-order/repository"
	order_service "github.com/MamangRust/monolith-point-of-sale-order/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderGapiTestSuite struct {
	tests.BaseTestSuite
	client        pborders.OrderQueryServiceClient
	clientCommand pborders.OrderCommandServiceClient
}

func (s *OrderGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Setup dependencies
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupCashierService()
	s.SetupTransactionService()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	// Order dependencies
	mencache := order_cache.NewMencache(cacheStore)
	repos := order_repo.NewRepositories(
		gormDB,
		pbcashiers.NewCashierQueryServiceClient(s.Conns["cashier"]),
		pbmerchants.NewMerchantQueryServiceClient(s.Conns["merchant"]),
		pbproducts.NewProductQueryServiceClient(s.Conns["product"]),
		pbproducts.NewProductCommandServiceClient(s.Conns["product"]),
		pborder_items.NewOrderItemQueryServiceClient(s.Conns["order-item"]),
		pborder_items.NewOrderItemCommandServiceClient(s.Conns["order-item"]),
	)
	svc := order_service.NewService(&order_service.Deps{
		Mencache:      mencache,
		Repositories:  repos,
		Logger:        s.Log,
		Observability: s.Obs,
	})

	// Handler
	handler := order_handler.NewHandler(&order_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pborders.RegisterOrderQueryServiceServer(server, handler.Order)
	pborders.RegisterOrderCommandServiceServer(server, handler.OrderCommand)
	pborders.RegisterOrderStatsServiceServer(server, handler.OrderStats)
	pborders.RegisterOrderStatsByIdServiceServer(server, handler.OrderStatsById)
	pborders.RegisterOrderStatsByMerchantServiceServer(server, handler.OrderStatsByMerchant)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.client = pborders.NewOrderQueryServiceClient(conn)
	s.clientCommand = pborders.NewOrderCommandServiceClient(conn)
}

func (s *OrderGapiTestSuite) TestOrderGapiLifecycle() {
	ctx := context.Background()

	// 1. Seed dependencies
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	// Seed a cashier (orders.cashier_id references cashiers.cashier_id)
	var cashierID int
	gormDB, err := s.GormDB()
	s.Require().NoError(err)
	err = gormDB.WithContext(ctx).Raw(`INSERT INTO cashiers (merchant_id, user_id, name) VALUES (?, ?, 'Order Gapi Cashier') RETURNING cashier_id`,
		merchID, userID,
	).Scan(&cashierID).Error
	s.Require().NoError(err)

	// 2. Create
	createRes, err := s.clientCommand.Create(ctx, &pborders.CreateOrderRequest{
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
	s.Require().NotNil(createRes)
	orderID := createRes.Data.Id

	// 3. FindById
	getRes, err := s.client.FindById(ctx, &pborders.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)
	s.Equal(int32(userID), getRes.Data.CashierId)

	// 4. FindAll
	allRes, err := s.client.FindAll(ctx, &pborders.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 5. FindByActive
	activeRes, err := s.client.FindByActive(ctx, &pborders.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 6. Update
	// Fetch order items first
	itemClient := pborder_items.NewOrderItemQueryServiceClient(s.Conns["order-item"])
	itemsRes, err := itemClient.FindOrderItemByOrder(ctx, &pborder_items.FindByIdOrderItemRequest{Id: orderID})
	s.Require().NoError(err)
	s.NotEmpty(itemsRes.Data)
	orderItemID := itemsRes.Data[0].Id

	_, err = s.clientCommand.Update(ctx, &pborders.UpdateOrderRequest{
		OrderId: orderID,
		Items: []*pborders.UpdateOrderItemRequest{
			{
				OrderItemId: orderItemID,
				ProductId:   int32(prodID),
				Quantity:    1,
			},
		},
	})
	s.Require().NoError(err)

	// 7. Trash
	_, err = s.clientCommand.TrashedOrder(ctx, &pborders.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 8. FindByTrashed
	trashedRes, err := s.client.FindByTrashed(ctx, &pborders.FindAllOrderRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 9. Restore
	_, err = s.clientCommand.RestoreOrder(ctx, &pborders.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 10. DeletePermanent
	_, _ = s.clientCommand.TrashedOrder(ctx, &pborders.FindByIdOrderRequest{Id: orderID})
	_, err = s.clientCommand.DeleteOrderPermanent(ctx, &pborders.FindByIdOrderRequest{Id: orderID})
	s.Require().NoError(err)

	// 11. RestoreAll — create a fresh trashed order first
	createRes2, _ := s.clientCommand.Create(ctx, &pborders.CreateOrderRequest{
		MerchantId: int32(merchID), CashierId: int32(cashierID),
		Items: []*pborders.CreateOrderItemRequest{{ProductId: int32(prodID), Quantity: 1}},
	})
	s.clientCommand.TrashedOrder(ctx, &pborders.FindByIdOrderRequest{Id: createRes2.Data.Id})
	_, err = s.clientCommand.RestoreAllOrder(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 12. DeleteAll — trash again
	createRes3, _ := s.clientCommand.Create(ctx, &pborders.CreateOrderRequest{
		MerchantId: int32(merchID), CashierId: int32(cashierID),
		Items: []*pborders.CreateOrderItemRequest{{ProductId: int32(prodID), Quantity: 1}},
	})
	s.clientCommand.TrashedOrder(ctx, &pborders.FindByIdOrderRequest{Id: createRes3.Data.Id})
	_, err = s.clientCommand.DeleteAllOrderPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestOrderGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderGapiTestSuite))
}
