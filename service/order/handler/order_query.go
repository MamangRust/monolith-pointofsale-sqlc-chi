package handler

import (
	"context"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"go.uber.org/zap"
	"math"

	"github.com/MamangRust/monolith-point-of-sale-order/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
)

type orderQueryHandleGrpc struct {
	pborders.UnimplementedOrderQueryServiceServer
	orderQuery service.OrderQueryService
	logger     logger.LoggerInterface
}

func NewOrderQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pborders.OrderQueryServiceServer {
	return &orderQueryHandleGrpc{
		orderQuery: service.OrderQuery,
		logger:     logger,
	}
}

func (s *orderQueryHandleGrpc) FindAll(ctx context.Context, request *pborders.FindAllOrderRequest) (*pborders.ApiResponsePaginationOrder, error) {
	s.logger.Info("FindAll orders called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderQuery.FindAll(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll orders success")

	return &pborders.ApiResponsePaginationOrder{
		Status:     "success",
		Message:    "Successfully fetched order",
		Data:       mapResponsesOrder(merchant),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderQueryHandleGrpc) FindById(ctx context.Context, request *pborders.FindByIdOrderRequest) (*pborders.ApiResponseOrder, error) {
	s.logger.Info("FindById order called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderQuery.FindById(ctx, id)
	if err != nil {
		s.logger.Error("FindById order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindById order success")

	return &pborders.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully fetched order",
		Data:    mapResponseOrder(merchant),
	}, nil
}

func (s *orderQueryHandleGrpc) FindByActive(ctx context.Context, request *pborders.FindAllOrderRequest) (*pborders.ApiResponsePaginationOrderDeleteAt, error) {
	s.logger.Info("FindByActive orders called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	merchant, totalRecords, err := s.orderQuery.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive orders success")

	return &pborders.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active order",
		Data:       mapResponsesOrderActive(merchant),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderQueryHandleGrpc) FindByTrashed(ctx context.Context, request *pborders.FindAllOrderRequest) (*pborders.ApiResponsePaginationOrderDeleteAt, error) {
	s.logger.Info("FindByTrashed orders called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrders{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	users, totalRecords, err := s.orderQuery.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed orders failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed orders success")

	return &pborders.ApiResponsePaginationOrderDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed order",
		Data:       mapResponsesOrderTrashed(users),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}
