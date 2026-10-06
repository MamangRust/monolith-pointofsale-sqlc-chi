package handler

import (
	"context"
	pborders "github.com/MamangRust/monolith-point-of-sale-pb/orders"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/MamangRust/monolith-point-of-sale-order/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
)

type orderCommandHandleGrpc struct {
	pborders.UnimplementedOrderCommandServiceServer
	orderCommand service.OrderCommandService
	logger       logger.LoggerInterface
}

func NewOrderCommandHandleGrpc(
	service *service.Service,
	logger logger.LoggerInterface,
) pborders.OrderCommandServiceServer {
	return &orderCommandHandleGrpc{
		orderCommand: service.OrderCommand,
		logger:       logger,
	}
}

func (s *orderCommandHandleGrpc) Create(ctx context.Context, request *pborders.CreateOrderRequest) (*pborders.ApiResponseOrder, error) {
	s.logger.Info("Create order called", zap.Int32("merchantId", request.GetMerchantId()))

	req := &requests.CreateOrderRequest{
		MerchantID: int(request.GetMerchantId()),
		CashierID:  int(request.GetCashierId()),
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.CreateOrderItemRequest{
			ProductID: int(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
		})
	}

	if err := req.Validate(); err != nil {
		return nil, order_errors.ErrGrpcValidateCreateOrder
	}

	order, err := s.orderCommand.CreateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Create order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Create order success")

	return &pborders.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully created order",
		Data:    mapResponseOrder(order),
	}, nil
}

func (s *orderCommandHandleGrpc) Update(ctx context.Context, request *pborders.UpdateOrderRequest) (*pborders.ApiResponseOrder, error) {
	s.logger.Info("Update order called", zap.Int32("id", request.GetOrderId()))

	id := int(request.GetOrderId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	req := &requests.UpdateOrderRequest{
		OrderID: &id,
	}

	for _, item := range request.GetItems() {
		req.Items = append(req.Items, requests.UpdateOrderItemRequest{
			OrderItemID: int(item.GetOrderItemId()),
			ProductID:   int(item.GetProductId()),
			Quantity:    int(item.GetQuantity()),
		})
	}

	if err := req.Validate(); err != nil {
		return nil, order_errors.ErrGrpcValidateUpdateOrder
	}

	order, err := s.orderCommand.UpdateOrder(ctx, req)
	if err != nil {
		s.logger.Error("Update order failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("Update order success")

	return &pborders.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully updated order",
		Data:    mapResponseOrder(order),
	}, nil
}

func (s *orderCommandHandleGrpc) TrashedOrder(ctx context.Context, request *pborders.FindByIdOrderRequest) (*pborders.ApiResponseOrderDeleteAt, error) {
	s.logger.Info("TrashedOrder called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderCommand.TrashedOrder(ctx, id)
	if err != nil {
		s.logger.Error("TrashedOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("TrashedOrder success")

	return &pborders.ApiResponseOrderDeleteAt{
		Status:  "success",
		Message: "Successfully trashed order",
		Data:    mapResponseOrderDeleteAt(merchant),
	}, nil
}

func (s *orderCommandHandleGrpc) RestoreOrder(ctx context.Context, request *pborders.FindByIdOrderRequest) (*pborders.ApiResponseOrderDeleteAt, error) {
	s.logger.Info("RestoreOrder called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	merchant, err := s.orderCommand.RestoreOrder(ctx, id)
	if err != nil {
		s.logger.Error("RestoreOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreOrder success")

	return &pborders.ApiResponseOrderDeleteAt{
		Status:  "success",
		Message: "Successfully restored order",
		Data:    mapResponseOrderDeleteAt(merchant),
	}, nil
}

func (s *orderCommandHandleGrpc) DeleteOrderPermanent(ctx context.Context, request *pborders.FindByIdOrderRequest) (*pborders.ApiResponseOrderDelete, error) {
	s.logger.Info("DeleteOrderPermanent called", zap.Int32("id", request.GetId()))

	id := int(request.GetId())
	if id <= 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.orderCommand.DeleteOrderPermanent(ctx, id)
	if err != nil {
		s.logger.Error("DeleteOrderPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteOrderPermanent success")

	return &pborders.ApiResponseOrderDelete{
		Status:  "success",
		Message: "Successfully deleted order permanently",
	}, nil
}

func (s *orderCommandHandleGrpc) RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pborders.ApiResponseOrderAll, error) {
	s.logger.Info("RestoreAllOrder called")

	_, err := s.orderCommand.RestoreAllOrder(ctx)
	if err != nil {
		s.logger.Error("RestoreAllOrder failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("RestoreAllOrder success")

	return &pborders.ApiResponseOrderAll{
		Status:  "success",
		Message: "Successfully restore all order",
	}, nil
}

func (s *orderCommandHandleGrpc) DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pborders.ApiResponseOrderAll, error) {
	s.logger.Info("DeleteAllOrderPermanent called")

	_, err := s.orderCommand.DeleteAllOrderPermanent(ctx)
	if err != nil {
		s.logger.Error("DeleteAllOrderPermanent failed", zap.Error(err))
		return nil, errors.ToGrpcError(err)
	}

	s.logger.Info("DeleteAllOrderPermanent success")

	return &pborders.ApiResponseOrderAll{
		Status:  "success",
		Message: "Successfully delete order permanen",
	}, nil
}

// Map helpers
