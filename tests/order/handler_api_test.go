package order_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	apierror "github.com/MamangRust/monolith-point-of-sale-apigateway/apierror"
	"github.com/MamangRust/monolith-point-of-sale-apigateway/cache/gateway_cache"
	"github.com/MamangRust/monolith-point-of-sale-apigateway/handler"
	pbcompat "github.com/MamangRust/monolith-point-of-sale-apigateway/pbcompat"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/suite"
)

type OrderApiTestSuite struct {
	tests.BaseTestSuite
	router    chi.Router
	orderID   int
	userID    int
	cashierID int
	merchID   int
	prodID    int
}

func (s *OrderApiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()
	s.SetupUserService()
	s.SetupMerchantService()
	s.SetupCategoryService()
	s.SetupProductService()
	s.SetupOrderItemService()
	s.SetupTransactionService()
	s.SetupOrderService()

	// Seed dependencies
	ctx := context.Background()
	userID := s.SeedUser(ctx)
	catID := s.SeedCategory(ctx)
	merchID := s.SeedMerchant(ctx, userID)
	prodID := s.SeedProduct(ctx, merchID, catID)

	// Seed a cashier (orders.cashier_id references cashiers.cashier_id)
	var cashierID int
	gormDB, err := s.GormDB()
	s.Require().NoError(err)
	err = gormDB.WithContext(ctx).Raw(`INSERT INTO cashiers (merchant_id, user_id, name) VALUES (?, ?, 'Order Api Cashier') RETURNING cashier_id`,
		merchID, userID,
	).Scan(&cashierID).Error
	s.Require().NoError(err)

	s.userID = userID
	s.cashierID = cashierID
	s.merchID = merchID
	s.prodID = prodID

	s.router = chi.NewRouter()

	orderConn := s.Conns["order"]

	handler.NewHandlerOrder(
		s.router,
		pbcompat.NewOrderServiceClient(orderConn),
		s.Log,
		response_api.NewOrderResponseMapper(),
		apierror.NewApiHandler(s.Obs, s.Log),
		gateway_cache.NewGatewayCache(s.GetCacheStore()),
	)
}

func (s *OrderApiTestSuite) TestOrderApiLifecycle() {
	// 1. Create
	reqBody := requests.CreateOrderRequest{
		MerchantID: s.merchID,
		CashierID:  s.cashierID,
		Items: []requests.CreateOrderItemRequest{
			{ProductID: s.prodID, Quantity: 1},
		},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/order/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.orderID = int(data["id"].(float64))

	// 2. FindById
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/order/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(float64(s.orderID), data["id"])

	// 3. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/order", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/order/active", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. Update
	updateBody := requests.UpdateOrderRequest{
		Items: []requests.UpdateOrderItemRequest{
			{ProductID: s.prodID, Quantity: 1},
		},
	}
	body, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order/update/%d", s.orderID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. Trash
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order/trashed/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/order/trashed", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order/restore/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. DeletePermanent
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order/trashed/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/order/permanent/%d", s.orderID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. RestoreAll — create a fresh trashed order first
	createBody, _ := json.Marshal(requests.CreateOrderRequest{
		MerchantID: s.merchID, CashierID: s.cashierID,
		Items: []requests.CreateOrderItemRequest{{ProductID: s.prodID, Quantity: 1}},
	})
	createReq := httptest.NewRequest(http.MethodPost, "/api/order/create", bytes.NewBuffer(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	s.router.ServeHTTP(createRec, createReq)
	var createRes map[string]interface{}
	json.Unmarshal(createRec.Body.Bytes(), &createRes)
	newOrderID := int(createRes["data"].(map[string]interface{})["id"].(float64))
	// Trash the new order
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order/trashed/%d", newOrderID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	// RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/order/restore/all", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. DeleteAll — trash again
	createBody2, _ := json.Marshal(requests.CreateOrderRequest{
		MerchantID: s.merchID, CashierID: s.cashierID,
		Items: []requests.CreateOrderItemRequest{{ProductID: s.prodID, Quantity: 1}},
	})
	createReq2 := httptest.NewRequest(http.MethodPost, "/api/order/create", bytes.NewBuffer(createBody2))
	createReq2.Header.Set("Content-Type", "application/json")
	createRec2 := httptest.NewRecorder()
	s.router.ServeHTTP(createRec2, createReq2)
	var createRes2 map[string]interface{}
	json.Unmarshal(createRec2.Body.Bytes(), &createRes2)
	newOrderID2 := int(createRes2["data"].(map[string]interface{})["id"].(float64))
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/order/trashed/%d", newOrderID2), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	// DeleteAll
	req = httptest.NewRequest(http.MethodPost, "/api/order/permanent/all", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func (s *OrderApiTestSuite) Test12_GetMonthlyTotalRevenue() {
	req := httptest.NewRequest(http.MethodGet, "/api/order/monthly-total-revenue?year=2024&month=1", nil)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func (s *OrderApiTestSuite) Test13_GetYearlyTotalRevenue() {
	req := httptest.NewRequest(http.MethodGet, "/api/order/yearly-total-revenue?year=2024", nil)
	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code)
}

func TestOrderApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(OrderApiTestSuite))
}
