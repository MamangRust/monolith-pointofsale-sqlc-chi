package order_test

import (
	"context"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"testing"
	"time"

	tests "github.com/MamangRust/monolith-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type OrderStatsServiceTestSuite struct {
	tests.BaseTestSuite
	clientStats  pborders.OrderStatsServiceClient
	clientStatsM pborders.OrderStatsByMerchantServiceClient
	merchantID   int
	userID       int
}

func (s *OrderStatsServiceTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	s.clientStats = pborders.NewOrderStatsServiceClient(s.Conns["order"])
	s.clientStatsM = pborders.NewOrderStatsByMerchantServiceClient(s.Conns["order"])

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Ensure created_at is set to current time to be picked up by stats
	gormDB, err := s.GormDB()
	s.Require().NoError(err)
	err = gormDB.WithContext(ctx).Exec("UPDATE orders SET created_at = ? WHERE order_id = ?",
		time.Now(), orderID).Error
	s.Require().NoError(err)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyTotalRevenue() {
	ctx := context.Background()
	now := time.Now()
	req := &pborders.FindYearMonthTotalRevenue{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.clientStats.FindMonthlyTotalRevenue(ctx, req)
	s.NoError(err)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyTotalRevenue() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.clientStats.FindYearlyTotalRevenue(ctx, &pborders.FindYearTotalRevenue{Year: int32(year)})
	s.NoError(err)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyRevenue() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.clientStats.FindMonthlyRevenue(ctx, &pborders.FindYearOrder{Year: int32(year)})
	s.NoError(err)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyRevenue() {
	ctx := context.Background()
	year := time.Now().Year()

	res, err := s.clientStats.FindYearlyRevenue(ctx, &pborders.FindYearOrder{Year: int32(year)})
	s.NoError(err)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyTotalRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pborders.FindYearMonthTotalRevenueByMerchant{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientStatsM.FindMonthlyTotalRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyTotalRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pborders.FindYearTotalRevenueByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientStatsM.FindYearlyTotalRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsServiceTestSuite) TestFindMonthlyRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pborders.FindYearOrderByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientStatsM.FindMonthlyRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res.Data)
}

func (s *OrderStatsServiceTestSuite) TestFindYearlyRevenueByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pborders.FindYearOrderByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientStatsM.FindYearlyRevenueByMerchant(ctx, req)
	s.NoError(err)
	s.NotEmpty(res.Data)
}

func TestOrderStatsServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderStatsServiceTestSuite))
}
