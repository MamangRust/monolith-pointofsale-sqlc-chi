package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-point-of-sale-order-item/cache"
	"github.com/MamangRust/monolith-point-of-sale-order-item/repository"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharederrorhandler "github.com/MamangRust/monolith-point-of-sale-shared/errorhandler"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderItemCommandDeps struct {
	Cache         mencache.OrderItemCommandCache
	Repo          repository.OrderItemCommandRepository
	Logger        logger.LoggerInterface
	Observability observability.TraceLoggerObservability
}

type orderItemCommandService struct {
	mencache      mencache.OrderItemCommandCache
	repo          repository.OrderItemCommandRepository
	logger        logger.LoggerInterface
	observability observability.TraceLoggerObservability
}

func NewOrderItemCommandService(params *orderItemCommandDeps) OrderItemCommandService {
	return &orderItemCommandService{
		mencache:      params.Cache,
		repo:          params.Repo,
		logger:        params.Logger,
		observability: params.Observability,
	}
}

func (s *orderItemCommandService) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*models.OrderItem, error) {
	const method = "CreateOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order.id", req.OrderID),
		attribute.Int("product.id", req.ProductID),
	)
	defer func() {
		end(status)
	}()

	orderItem, err := s.repo.CreateOrderItem(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.OrderItem](
			s.logger,
			orderitem_errors.ErrFailedCreateOrderItem.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedOrderItems(ctx, req.OrderID)
	logSuccess("Successfully created order item", zap.Int32("orderItem.id", orderItem.OrderItemID))
	return orderItem, nil
}

func (s *orderItemCommandService) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*models.OrderItem, error) {
	const method = "UpdateOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderItem.id", req.OrderItemID),
		attribute.Int("order.id", req.OrderID),
	)
	defer func() {
		end(status)
	}()

	orderItem, err := s.repo.UpdateOrderItem(ctx, req)
	if err != nil {
		status = "error"
		return sharederrorhandler.HandleError[*models.OrderItem](
			s.logger,
			orderitem_errors.ErrFailedUpdateOrderItem.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
	}

	s.mencache.DeleteCachedOrderItems(ctx, req.OrderID)
	logSuccess("Successfully updated order item", zap.Int32("orderItem.id", orderItem.OrderItemID))
	return orderItem, nil
}

func (s *orderItemCommandService) DeleteOrderItem(ctx context.Context, orderItemID int) error {
	const method = "DeleteOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderItem.id", orderItemID),
	)
	defer func() {
		end(status)
	}()

	if err := s.repo.DeleteOrderItem(ctx, orderItemID); err != nil {
		status = "error"
		_, mappedErr := sharederrorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedDeleteOrderItem.WithInternal(err),
			method,
			span,
			zap.Error(err),
		)
		return mappedErr
	}

	s.mencache.DeleteCachedOrderItemsAll(ctx)
	logSuccess("Successfully deleted order item", zap.Int("orderItem.id", orderItemID))
	return nil
}
