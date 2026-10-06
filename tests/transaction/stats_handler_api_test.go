package transaction_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apierror "github.com/MamangRust/monolith-point-of-sale-apigateway/apierror"
	apigw_cache "github.com/MamangRust/monolith-point-of-sale-apigateway/cache/gateway_cache"
	apigateway "github.com/MamangRust/monolith-point-of-sale-apigateway/handler"
	pbcompat "github.com/MamangRust/monolith-point-of-sale-apigateway/pbcompat"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
)

type TransactionStatsApiTestSuite struct {
	tests.BaseTestSuite
	router     chi.Router
	merchantID int
	userID     int
}

func (s *TransactionStatsApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupRoleService()
	s.SetupUserService()
	s.SetupCategoryService()
	s.SetupMerchantService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	s.router = chi.NewRouter()

	apigateway.NewHandlerTransaction(
		s.router,
		pbcompat.NewTransactionServiceClient(s.Conns["transaction"]),
		s.Log,
		response_api.NewTransactionResponseMapper(),
		apierror.NewApiHandler(s.Obs, s.Log),
		apigw_cache.NewGatewayCache(s.GetCacheStore()),
	)

	ctx := context.Background()
	s.userID = s.SeedUser(ctx)
	s.merchantID = s.SeedMerchant(ctx, s.userID)
	catID := s.SeedCategory(ctx)
	prodID := s.SeedProduct(ctx, s.merchantID, catID)
	orderID := s.SeedOrder(ctx, s.userID, s.merchantID, prodID)

	// Seed a successful transaction
	db, err := s.GormDB()
	s.Require().NoError(err)
	err = db.WithContext(ctx).Exec(`
		INSERT INTO transactions (merchant_id, order_id, amount, payment_method, payment_status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		s.merchantID, orderID, 100000, "credit_card", "completed", time.Now()).Error
	s.Require().NoError(err)
}

func (s *TransactionStatsApiTestSuite) TestFindMonthStatusSuccess() {
	now := time.Now()
	url := fmt.Sprintf("/api/transaction/monthly-success?year=%d&month=%d", now.Year(), int(now.Month()))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *TransactionStatsApiTestSuite) TestFindYearStatusSuccess() {
	year := time.Now().Year()
	url := fmt.Sprintf("/api/transaction/yearly-success?year=%d", year)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *TransactionStatsApiTestSuite) TestFindMonthMethodSuccess() {
	now := time.Now()
	url := fmt.Sprintf("/api/transaction/monthly-method-success?year=%d&month=%d", now.Year(), int(now.Month()))
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func (s *TransactionStatsApiTestSuite) TestFindMonthStatusSuccessByMerchant() {
	now := time.Now()
	url := fmt.Sprintf("/api/transaction/merchant/monthly-success?year=%d&month=%d&merchant_id=%d",
		now.Year(), int(now.Month()), s.merchantID)
	req := httptest.NewRequest(http.MethodGet, url, nil)
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	s.Equal("success", res["status"])
	s.NotEmpty(res["data"])
}

func TestTransactionStatsApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(TransactionStatsApiTestSuite))
}
