package handler

import (
	"context"

	pborder_items "github.com/MamangRust/monolith-point-of-sale-pb/order_items"
	"go.uber.org/zap"

	"github.com/MamangRust/monolith-point-of-sale-order-item/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
)

type orderItemCommandHandleGrpc struct {
	pborder_items.UnimplementedOrderItemCommandServiceServer
	orderItemCommandService service.OrderItemCommandService
	logger                  logger.LoggerInterface
}

func NewOrderItemCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pborder_items.OrderItemCommandServiceServer {
	return &orderItemCommandHandleGrpc{
		orderItemCommandService: service.OrderItemCommand,
		logger:                  logger,
	}
}

func (s *orderItemCommandHandleGrpc) CreateOrderItem(ctx context.Context, request *pborder_items.CreateOrderItemRecordRequest) (*pborder_items.ApiResponseOrderItem, error) {
	s.logger.Info("CreateOrderItem called", zap.Int32("order_id", request.GetOrderId()))

	if request.GetOrderId() <= 0 || request.GetProductId() <= 0 || request.GetQuantity() <= 0 {
		return nil, errors.ToGrpcError(errors.ErrBadRequest.WithMessage("Invalid order item payload"))
	}

	req := &requests.CreateOrderItemRecordRequest{
		OrderID:   int(request.GetOrderId()),
		ProductID: int(request.GetProductId()),
		Quantity:  int(request.GetQuantity()),
		Price:     int(request.GetPrice()),
	}

	orderItem, err := s.orderItemCommandService.CreateOrderItem(ctx, req)
	if err != nil {
		s.logger.Error("CreateOrderItem failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("CreateOrderItem success")

	return &pborder_items.ApiResponseOrderItem{
		Status:  "success",
		Message: "Successfully created order item",
		Data:    mapOrderItemToProto(orderItem),
	}, nil
}

func (s *orderItemCommandHandleGrpc) UpdateOrderItem(ctx context.Context, request *pborder_items.UpdateOrderItemRecordRequest) (*pborder_items.ApiResponseOrderItem, error) {
	s.logger.Info("UpdateOrderItem called", zap.Int32("order_item_id", request.GetOrderItemId()))

	if request.GetOrderItemId() <= 0 {
		return nil, errors.ToGrpcError(errors.ErrBadRequest.WithMessage("Invalid order item ID"))
	}

	req := &requests.UpdateOrderItemRecordRequest{
		OrderItemID: int(request.GetOrderItemId()),
		OrderID:     int(request.GetOrderId()),
		ProductID:   int(request.GetProductId()),
		Quantity:    int(request.GetQuantity()),
		Price:       int(request.GetPrice()),
	}

	orderItem, err := s.orderItemCommandService.UpdateOrderItem(ctx, req)
	if err != nil {
		s.logger.Error("UpdateOrderItem failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("UpdateOrderItem success")

	return &pborder_items.ApiResponseOrderItem{
		Status:  "success",
		Message: "Successfully updated order item",
		Data:    mapOrderItemToProto(orderItem),
	}, nil
}

func (s *orderItemCommandHandleGrpc) DeleteOrderItem(ctx context.Context, request *pborder_items.DeleteOrderItemRecordRequest) (*pborder_items.ApiResponseOrderItemDelete, error) {
	s.logger.Info("DeleteOrderItem called", zap.Int32("order_item_id", request.GetOrderItemId()))

	if request.GetOrderItemId() <= 0 {
		return nil, errors.ToGrpcError(errors.ErrBadRequest.WithMessage("Invalid order item ID"))
	}

	if err := s.orderItemCommandService.DeleteOrderItem(ctx, int(request.GetOrderItemId())); err != nil {
		s.logger.Error("DeleteOrderItem failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteOrderItem success")

	return &pborder_items.ApiResponseOrderItemDelete{
		Status:  "success",
		Message: "Successfully deleted order item",
	}, nil
}
