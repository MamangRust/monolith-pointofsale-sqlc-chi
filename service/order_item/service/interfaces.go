package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type OrderItemQueryService interface {
	FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*models.OrderItemRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*models.OrderItemActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*models.OrderItemTrashedRow, *int, error)
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*models.OrderItem, error)
}

type OrderItemCommandService interface {
	CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*models.OrderItem, error)
	UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*models.OrderItem, error)
	DeleteOrderItem(ctx context.Context, orderItemID int) error
}
