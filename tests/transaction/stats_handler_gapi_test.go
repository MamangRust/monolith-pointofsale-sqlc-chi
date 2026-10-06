package transaction_test

import (
	"context"
	pbtransactions "github.com/MamangRust/monolith-point-of-sale-pb/transactions"
	"testing"
	"time"

	tests "github.com/MamangRust/monolith-point-of-sale-test"
	"github.com/stretchr/testify/suite"
)

type TransactionStatsGapiTestSuite struct {
	tests.BaseTestSuite
	clientStats pbtransactions.TransactionStatsServiceClient
	userID      int
	merchantID  int
}

func (s *TransactionStatsGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	s.clientStats = pbtransactions.NewTransactionStatsServiceClient(s.Conns["transaction"])

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Create a successful transaction
	gormDB, err := s.GormDB()
	s.Require().NoError(err)
	err = gormDB.WithContext(ctx).Exec(`
		INSERT INTO transactions (order_id, merchant_id, amount, payment_method, payment_status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		orderID, s.merchantID, 10000, "credit_card", "completed", time.Now()).Error
	s.Require().NoError(err)

	err = gormDB.WithContext(ctx).Exec(`
		INSERT INTO transactions (order_id, merchant_id, amount, payment_method, payment_status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		orderID, s.merchantID, 5000, "bank_transfer", "failed", time.Now()).Error
	s.Require().NoError(err)
}

func (s *TransactionStatsGapiTestSuite) TestFindMonthStatusSuccess() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransactions.FindMonthlyTransactionStatus{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.clientStats.FindMonthStatusSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func (s *TransactionStatsGapiTestSuite) TestFindYearStatusSuccess() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransactions.FindYearlyTransactionStatus{
		Year: int32(year),
	}

	res, err := s.clientStats.FindYearStatusSuccess(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func (s *TransactionStatsGapiTestSuite) TestFindMonthStatusFailed() {
	ctx := context.Background()
	now := time.Now()
	req := &pbtransactions.FindMonthlyTransactionStatus{
		Year:  int32(now.Year()),
		Month: int32(now.Month()),
	}

	res, err := s.clientStats.FindMonthStatusFailed(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func (s *TransactionStatsGapiTestSuite) TestFindYearStatusFailed() {
	ctx := context.Background()
	year := time.Now().Year()
	req := &pbtransactions.FindYearlyTransactionStatus{
		Year: int32(year),
	}

	res, err := s.clientStats.FindYearStatusFailed(ctx, req)
	s.NoError(err)
	s.Equal("success", res.Status)
}

func TestTransactionStatsGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsGapiTestSuite))
}
