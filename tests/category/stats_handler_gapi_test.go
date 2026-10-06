package category_test

import (
	"context"
	pbcategories "github.com/MamangRust/monolith-point-of-sale-pb/categories"
	"testing"
	"time"

	tests "github.com/MamangRust/monolith-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type CategoryStatsGapiTestSuite struct {
	tests.BaseTestSuite
	clientStats   pbcategories.CategoryStatsServiceClient
	clientStatsID pbcategories.CategoryStatsByIdServiceClient
	clientStatsM  pbcategories.CategoryStatsByMerchantServiceClient
	categoryID    int
	merchantID    int
	userID        int
}

func (s *CategoryStatsGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	s.clientStats = pbcategories.NewCategoryStatsServiceClient(s.Conns["category"])
	s.clientStatsID = pbcategories.NewCategoryStatsByIdServiceClient(s.Conns["category"])
	s.clientStatsM = pbcategories.NewCategoryStatsByMerchantServiceClient(s.Conns["category"])

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	s.categoryID = s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, s.categoryID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Ensure created_at is set to current time to be picked up by stats
	db, err := s.GormDB()
	s.Require().NoError(err)
	err = db.WithContext(ctx).Exec("UPDATE orders SET created_at = ? WHERE order_id = ?",
		time.Now(), orderID).Error
	s.Require().NoError(err)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthlyTotalPrices() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearMonthTotalPrices{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.clientStats.FindMonthlyTotalPrices(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearlyTotalPrices() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbcategories.FindYearTotalPrices{
		Year: int32(year),
	}

	res, err := s.clientStats.FindYearlyTotalPrices(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthPrice() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbcategories.FindYearCategory{
		Year: int32(year),
	}

	res, err := s.clientStats.FindMonthPrice(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearPrice() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbcategories.FindYearCategory{
		Year: int32(year),
	}

	res, err := s.clientStats.FindYearPrice(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthlyTotalPricesById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearMonthTotalPriceById{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientStatsID.FindMonthlyTotalPricesById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearlyTotalPricesById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearTotalPriceById{
		Year:       int32(now.Year()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientStatsID.FindYearlyTotalPricesById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthPriceById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearCategoryById{
		Year:       int32(now.Year()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientStatsID.FindMonthPriceById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearPriceById() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearCategoryById{
		Year:       int32(now.Year()),
		CategoryId: int32(s.categoryID),
	}

	res, err := s.clientStatsID.FindYearPriceById(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthlyTotalPricesByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearMonthTotalPriceByMerchant{
		Year:       int32(now.Year()),
		Month:      int32(now.Month()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientStatsM.FindMonthlyTotalPricesByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearlyTotalPricesByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearTotalPriceByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientStatsM.FindYearlyTotalPricesByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindMonthPriceByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearCategoryByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientStatsM.FindMonthPriceByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func (s *CategoryStatsGapiTestSuite) TestFindYearPriceByMerchant() {
	ctx := context.Background()
	now := time.Now()
	req := &pbcategories.FindYearCategoryByMerchant{
		Year:       int32(now.Year()),
		MerchantId: int32(s.merchantID),
	}

	res, err := s.clientStatsM.FindYearPriceByMerchant(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
	s.NotEmpty(res.Data)
}

func TestCategoryStatsGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryStatsGapiTestSuite))
}
