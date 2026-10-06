package handler

import (
	"context"
	pbcommon "github.com/MamangRust/monolith-point-of-sale-pb/common"
	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"go.uber.org/zap"
	"math"

	"github.com/MamangRust/monolith-point-of-sale-order-item/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
)

type orderItemQueryHandleGrpc struct {
	pborder_items.UnimplementedOrderItemQueryServiceServer
	orderItemService service.OrderItemQueryService
	logger           logger.LoggerInterface
}

func NewOrderItemQueryHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pborder_items.OrderItemQueryServiceServer {
	return &orderItemQueryHandleGrpc{
		orderItemService: service.OrderItemQuery,
		logger:           logger,
	}
}

func (s *orderItemQueryHandleGrpc) FindAll(ctx context.Context, request *pborder_items.FindAllOrderItemRequest) (*pborder_items.ApiResponsePaginationOrderItem, error) {
	s.logger.Info("FindAll order items called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindAllOrderItems(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindAll order items failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindAll order items success")

	return &pborder_items.ApiResponsePaginationOrderItem{
		Status:     "success",
		Message:    "Successfully fetched order items",
		Data:       mapResponsesOrderItem(orderItems),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderItemQueryHandleGrpc) FindByActive(ctx context.Context, request *pborder_items.FindAllOrderItemRequest) (*pborder_items.ApiResponsePaginationOrderItemDeleteAt, error) {
	s.logger.Info("FindByActive order items called", zap.Int32("page", request.GetPage()))

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByActive(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByActive order items failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByActive order items success")

	return &pborder_items.ApiResponsePaginationOrderItemDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched active order items",
		Data:       mapResponsesOrderItemActive(orderItems),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderItemQueryHandleGrpc) FindByTrashed(ctx context.Context, request *pborder_items.FindAllOrderItemRequest) (*pborder_items.ApiResponsePaginationOrderItemDeleteAt, error) {
	s.logger.Info("FindByTrashed order items called")

	page := int(request.GetPage())
	pageSize := int(request.GetPageSize())
	search := request.GetSearch()

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	reqService := requests.FindAllOrderItems{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}

	orderItems, totalRecords, err := s.orderItemService.FindByTrashed(ctx, &reqService)
	if err != nil {
		s.logger.Error("FindByTrashed order items failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	totalPages := int(math.Ceil(float64(*totalRecords) / float64(pageSize)))

	paginationMeta := &pbcommon.PaginationMeta{
		CurrentPage:  int32(page),
		PageSize:     int32(pageSize),
		TotalPages:   int32(totalPages),
		TotalRecords: int32(*totalRecords),
	}

	s.logger.Info("FindByTrashed order items success")

	return &pborder_items.ApiResponsePaginationOrderItemDeleteAt{
		Status:     "success",
		Message:    "Successfully fetched trashed order items",
		Data:       mapResponsesOrderItemTrashed(orderItems),
		Pagination: mapPaginationMeta(paginationMeta),
	}, nil
}

func (s *orderItemQueryHandleGrpc) FindOrderItemByOrder(ctx context.Context, request *pborder_items.FindByIdOrderItemRequest) (*pborder_items.ApiResponsesOrderItem, error) {
	s.logger.Info("FindOrderItemByOrder called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	orderItems, err := s.orderItemService.FindOrderItemByOrder(ctx, id)
	if err != nil {
		s.logger.Error("FindOrderItemByOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("FindOrderItemByOrder success")

	return &pborder_items.ApiResponsesOrderItem{
		Status:  "success",
		Message: "Successfully fetched order items by order",
		Data:    mapResponsesOrderItemFromModel(orderItems),
	}, nil
}

// Map helpers
