package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-apigateway/apierror"
	gateway_cache "github.com/MamangRust/monolith-point-of-sale-apigateway/cache/gateway_cache"
	"github.com/MamangRust/monolith-point-of-sale-apigateway/httpx"
	pb "github.com/MamangRust/monolith-point-of-sale-apigateway/pbcompat"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/response"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	"github.com/go-chi/chi/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

type cashierHandleApi struct {
	client     pb.CashierServiceClient
	logger     logger.LoggerInterface
	mapping    response_api.CashierResponseMapper
	apiHandler apierror.ApiHandler
	cache      *gateway_cache.GatewayCache
}

func NewHandlerCashier(
	router chi.Router,
	client pb.CashierServiceClient,
	logger logger.LoggerInterface,
	mapping response_api.CashierResponseMapper,
	apiHandler apierror.ApiHandler,
	cache *gateway_cache.GatewayCache,
) *cashierHandleApi {
	cashierHandler := &cashierHandleApi{
		client:     client,
		logger:     logger,
		mapping:    mapping,
		apiHandler: apiHandler,
		cache:      cache,
	}

	router.Route("/api/cashier", func(routerCashier chi.Router) {

		routerCashier.Get("/", apiHandler.Handle("get-cashier-findallcashier", cashierHandler.FindAllCashier))
		routerCashier.Get("/{id}", apiHandler.Handle("get-cashier-findbyid", cashierHandler.FindById))
		routerCashier.Get("/active", apiHandler.Handle("get-cashier-findbyactive", cashierHandler.FindByActive))
		routerCashier.Get("/trashed", apiHandler.Handle("get-cashier-findbytrashed", cashierHandler.FindByTrashed))

		routerCashier.Get("/monthly-total-sales", apiHandler.Handle("get-cashier-findmonthlytotalsales", cashierHandler.FindMonthlyTotalSales))
		routerCashier.Get("/yearly-total-sales", apiHandler.Handle("get-cashier-findyeartotalsales", cashierHandler.FindYearTotalSales))

		routerCashier.Get("/merchant/monthly-total-sales", apiHandler.Handle("get-cashier-findmonthlytotalsalesbymerchant", cashierHandler.FindMonthlyTotalSalesByMerchant))
		routerCashier.Get("/merchant/yearly-total-sales", apiHandler.Handle("get-cashier-findyeartotalsalesbymerchant", cashierHandler.FindYearTotalSalesByMerchant))

		routerCashier.Get("/mycashier/monthly-total-sales", apiHandler.Handle("get-cashier-findmonthlytotalsalesbyid", cashierHandler.FindMonthlyTotalSalesById))
		routerCashier.Get("/mycashier/yearly-total-sales", apiHandler.Handle("get-cashier-findyeartotalsalesbyid", cashierHandler.FindYearTotalSalesById))

		routerCashier.Get("/monthly-sales", apiHandler.Handle("get-cashier-findmonthsales", cashierHandler.FindMonthSales))
		routerCashier.Get("/yearly-sales", apiHandler.Handle("get-cashier-findyearsales", cashierHandler.FindYearSales))
		routerCashier.Get("/merchant/monthly-sales", apiHandler.Handle("get-cashier-findmonthsalesbymerchant", cashierHandler.FindMonthSalesByMerchant))
		routerCashier.Get("/merchant/yearly-sales", apiHandler.Handle("get-cashier-findyearsalesbymerchant", cashierHandler.FindYearSalesByMerchant))
		routerCashier.Get("/mycashier/monthly-sales", apiHandler.Handle("get-cashier-findmonthsalesbyid", cashierHandler.FindMonthSalesById))
		routerCashier.Get("/mycashier/yearly-sales", apiHandler.Handle("get-cashier-findyearsalesbyid", cashierHandler.FindYearSalesById))

		routerCashier.Post("/create", apiHandler.Handle("post-cashier-createcashier", cashierHandler.CreateCashier))
		routerCashier.Post("/update/{id}", apiHandler.Handle("post-cashier-updatecashier", cashierHandler.UpdateCashier))

		routerCashier.Post("/trashed/{id}", apiHandler.Handle("post-cashier-trashedcashier", cashierHandler.TrashedCashier))
		routerCashier.Post("/restore/{id}", apiHandler.Handle("post-cashier-restorecashier", cashierHandler.RestoreCashier))
		routerCashier.Delete("/permanent/{id}", apiHandler.Handle("delete-cashier-deletecashierpermanent", cashierHandler.DeleteCashierPermanent))

		routerCashier.Post("/restore/all", apiHandler.Handle("post-cashier-restoreallcashier", cashierHandler.RestoreAllCashier))
		routerCashier.Post("/permanent/all", apiHandler.Handle("post-cashier-deleteallcashierpermanent", cashierHandler.DeleteAllCashierPermanent))

	})
	return cashierHandler
}

// List all cashier (paginated)
// List all cashier (paginated)
// List all cashier (paginated)
// @Summary List all cashier (paginated)
// @Tags Cashier
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationCashier
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindAllCashier(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page := parseQueryInt(r, "page", 1)
	pageSize := parseQueryInt(r, "page_size", 10)
	search := r.URL.Query().Get("search")

	cacheKey := fmt.Sprintf("cashier:findallcashier:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationCashier](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindAll(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationCashier(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier by ID
// Get cashier by ID
// Get cashier by ID
// @Summary Get cashier by ID
// @Tags Cashier
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} response.ApiResponseCashier
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/:id [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return errors.NewBadRequestError("invalid cashier ID")
	}

	cacheKey := fmt.Sprintf("cashier:findbyid:id_%d", id)
	if cached, found := gateway_cache.Get[response.ApiResponseCashier](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}
	res, err := h.client.FindById(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashier(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// List active cashier
// List active cashier
// List active cashier
// @Summary List active cashier
// @Tags Cashier
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationCashierDeleteAt
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/active [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindByActive(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page := parseQueryInt(r, "page", 1)
	pageSize := parseQueryInt(r, "page_size", 10)
	search := r.URL.Query().Get("search")

	cacheKey := fmt.Sprintf("cashier:findbyactive:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationCashierDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByActive(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationCashierDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// List trashed cashier
// List trashed cashier
// List trashed cashier
// @Summary List trashed cashier
// @Tags Cashier
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param search query string false "Search keyword"
// @Success 200 {object} response.ApiResponsePaginationCashierDeleteAt
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/trashed [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindByTrashed(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	page := parseQueryInt(r, "page", 1)
	pageSize := parseQueryInt(r, "page_size", 10)
	search := r.URL.Query().Get("search")

	cacheKey := fmt.Sprintf("cashier:findbytrashed:page_%d:size_%d:search_%s", page, pageSize, search)
	if cached, found := gateway_cache.Get[response.ApiResponsePaginationCashierDeleteAt](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	req := &pb.FindAllCashierRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
		Search:   search,
	}

	res, err := h.client.FindByTrashed(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponsePaginationCashierDeleteAt(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashier
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/monthly-total-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindMonthlyTotalSales(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(r, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthlytotalsales:year_%d:month_%d", year, month)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindMonthlyTotalSales(ctx, &pb.FindYearMonthTotalSales{
		Year:  int32(year),
		Month: int32(month),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashierYearlyTotalSales
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/yearly-total-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindYearTotalSales(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("cashier:findyeartotalsales:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindYearlyTotalSales(ctx, &pb.FindYearTotalSales{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashier
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/mycashier/monthly-total-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindMonthlyTotalSalesById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(r, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}
	cashierID, err := parseQueryIntWithValidation(r, "cashier_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid cashier_id")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthlytotalsalesbyid:year_%d:month_%d:cashierID_%d", year, month, cashierID)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindMonthlyTotalSalesById(ctx, &pb.FindYearMonthTotalSalesById{
		Year:      int32(year),
		Month:     int32(month),
		CashierId: int32(cashierID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashierYearlyTotalSales
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/mycashier/yearly-total-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindYearTotalSalesById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	cashierID, err := parseQueryIntWithValidation(r, "cashier_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid cashierID")
	}

	cacheKey := fmt.Sprintf("cashier:findyeartotalsalesbyid:year_%d:cashierID_%d", year, cashierID)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindYearlyTotalSalesById(ctx, &pb.FindYearTotalSalesById{
		Year:      int32(year),
		CashierId: int32(cashierID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashier
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/merchant/monthly-total-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindMonthlyTotalSalesByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	month, err := parseQueryIntWithValidation(r, "month", 1, 12)
	if err != nil {
		return errors.NewBadRequestError("invalid month")
	}
	merchantID, err := parseQueryIntWithValidation(r, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchantID")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthlytotalsalesbymerchant:year_%d:month_%d:merchantID_%d", year, month, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindMonthlyTotalSalesByMerchant(ctx, &pb.FindYearMonthTotalSalesByMerchant{
		Year:       int32(year),
		Month:      int32(month),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseMonthlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashierYearlyTotalSales
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/merchant/yearly-total-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindYearTotalSalesByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchantID, err := parseQueryIntWithValidation(r, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchantID")
	}

	cacheKey := fmt.Sprintf("cashier:findyeartotalsalesbymerchant:year_%d:merchantID_%d", year, merchantID)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearlyTotalSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindYearlyTotalSalesByMerchant(ctx, &pb.FindYearTotalSalesByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchantID),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseYearlyTotalSales(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashierMonthSales
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/monthly-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindMonthSales(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthsales:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindMonthSales(ctx, &pb.FindYearCashier{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashierYearSales
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/yearly-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindYearSales(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}

	cacheKey := fmt.Sprintf("cashier:findyearsales:year_%d", year)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindYearSales(ctx, &pb.FindYearCashier{
		Year: int32(year),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashierMonthSales
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/merchant/monthly-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindMonthSalesByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchant, err := parseQueryIntWithValidation(r, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchant_id")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthsalesbymerchant:year_%d:merchantID_%d", year, merchant)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindMonthSalesByMerchant(ctx, &pb.FindYearCashierByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashierYearSales
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/merchant/yearly-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindYearSalesByMerchant(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	merchant, err := parseQueryIntWithValidation(r, "merchant_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid merchant_id")
	}

	cacheKey := fmt.Sprintf("cashier:findyearsalesbymerchant:year_%d:merchantID_%d", year, merchant)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindYearSalesByMerchant(ctx, &pb.FindYearCashierByMerchant{
		Year:       int32(year),
		MerchantId: int32(merchant),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashierMonthSales
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/mycashier/monthly-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindMonthSalesById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	cashier, err := parseQueryIntWithValidation(r, "cashier_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid cashier_id")
	}

	cacheKey := fmt.Sprintf("cashier:findmonthsalesbyid:year_%d:cashierID_%d", year, cashier)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierMonthSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindMonthSalesById(ctx, &pb.FindYearCashierById{
		Year:      int32(year),
		CashierId: int32(cashier),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierMonthlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Get cashier statistics
// Get cashier statistics
// Get cashier statistics
// @Summary Get cashier statistics
// @Tags Cashier
// @Accept json
// @Produce json
// @Param year query int false "Year (e.g. 2026)"
// @Param month query int false "Month (1-12)"
// @Success 200 {object} response.ApiResponseCashierYearSales
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/mycashier/yearly-sales [get]
// @Security BearerAuth
func (h *cashierHandleApi) FindYearSalesById(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	year, err := parseQueryIntWithValidation(r, "year", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid year")
	}
	cashier_id, err := parseQueryIntWithValidation(r, "cashier_id", 1, 9999)
	if err != nil {
		return errors.NewBadRequestError("invalid cashier_id")
	}

	cacheKey := fmt.Sprintf("cashier:findyearsalesbyid:year_%d:cashierID_%d", year, cashier_id)
	if cached, found := gateway_cache.Get[response.ApiResponseCashierYearSales](ctx, h.cache, cacheKey); found && cached != nil {
		return httpx.JSON(w, http.StatusOK, cached)
	}

	res, err := h.client.FindYearSalesById(ctx, &pb.FindYearCashierById{
		Year:      int32(year),
		CashierId: int32(cashier_id),
	})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierYearlySale(res)

	gateway_cache.Set(ctx, h.cache, cacheKey, so, 5*time.Minute)
	return httpx.JSON(w, http.StatusOK, so)
}

// Create cashier
// Create cashier
// Create cashier
// @Summary Create cashier
// @Tags Cashier
// @Accept json
// @Produce json
// @Param request body requests.CreateCashierRequest true "Request body"
// @Success 200 {object} response.ApiResponseCashier
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/create [post]
// @Security BearerAuth
func (h *cashierHandleApi) CreateCashier(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	var body requests.CreateCashierRequest
	if err := httpx.Bind(r, &body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.CreateCashierRequest{
		MerchantId: int32(body.MerchantID),
		UserId:     int32(body.UserID),
		Name:       body.Name,
	}

	res, err := h.client.CreateCashier(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashier(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return httpx.JSON(w, http.StatusCreated, so)
}

// Update cashier
// Update cashier
// Update cashier
// @Summary Update cashier
// @Tags Cashier
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Param request body requests.UpdateCashierRequest true "Request body"
// @Success 200 {object} response.ApiResponseCashier
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/update/:id [post]
// @Security BearerAuth
func (h *cashierHandleApi) UpdateCashier(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	idStr, err := strconv.Atoi(id)
	if err != nil {
		return errors.NewBadRequestError("invalid cashier ID")
	}

	var body requests.UpdateCashierRequest
	if err := httpx.Bind(r, &body); err != nil {
		return errors.NewBadRequestError("bind failed")
	}
	if err := body.Validate(); err != nil {
		return errors.NewBadRequestError("validation failed")
	}

	req := &pb.UpdateCashierRequest{
		CashierId: int32(idStr),
		Name:      body.Name,
	}

	res, err := h.client.UpdateCashier(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashier(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return httpx.JSON(w, http.StatusOK, so)
}

// Trash cashier
// Trash cashier
// Trash cashier
// @Summary Trash cashier
// @Tags Cashier
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} response.ApiResponseCashierDeleteAt
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/trashed/:id [post]
// @Security BearerAuth
func (h *cashierHandleApi) TrashedCashier(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return errors.NewBadRequestError("invalid cashier ID")
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}
	res, err := h.client.TrashedCashier(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return httpx.JSON(w, http.StatusOK, so)
}

// Restore cashier
// Restore cashier
// Restore cashier
// @Summary Restore cashier
// @Tags Cashier
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} response.ApiResponseCashierDeleteAt
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/restore/:id [post]
// @Security BearerAuth
func (h *cashierHandleApi) RestoreCashier(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return errors.NewBadRequestError("invalid cashier ID")
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}
	res, err := h.client.RestoreCashier(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierDeleteAt(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return httpx.JSON(w, http.StatusOK, so)
}

// Delete cashier permanently
// Delete cashier permanently
// Delete cashier permanently
// @Summary Delete cashier permanently
// @Tags Cashier
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} response.ApiResponseCashierDelete
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/permanent/:id [delete]
// @Security BearerAuth
func (h *cashierHandleApi) DeleteCashierPermanent(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return errors.NewBadRequestError("invalid cashier ID")
	}

	req := &pb.FindByIdCashierRequest{Id: int32(id)}
	res, err := h.client.DeleteCashierPermanent(ctx, req)
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierDelete(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return httpx.JSON(w, http.StatusOK, so)
}

// Restore all trashed cashier
// Restore all trashed cashier
// Restore all trashed cashier
// @Summary Restore all trashed cashier
// @Tags Cashier
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCashierAll
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/restore/all [post]
// @Security BearerAuth
func (h *cashierHandleApi) RestoreAllCashier(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	res, err := h.client.RestoreAllCashier(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return httpx.JSON(w, http.StatusOK, so)
}

// Delete all trashed cashier permanently
// Delete all trashed cashier permanently
// Delete all trashed cashier permanently
// @Summary Delete all trashed cashier permanently
// @Tags Cashier
// @Accept json
// @Produce json
// @Success 200 {object} response.ApiResponseCashierAll
// @Failure 400 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/cashier/permanent/all [post]
// @Security BearerAuth
func (h *cashierHandleApi) DeleteAllCashierPermanent(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	res, err := h.client.DeleteAllCashierPermanent(ctx, &emptypb.Empty{})
	if err != nil {
		return errors.ParseGrpcError(err)
	}

	so := h.mapping.ToApiResponseCashierAll(res)

	gateway_cache.InvalidatePattern(ctx, h.cache, "cashier:*")
	return httpx.JSON(w, http.StatusOK, so)
}
