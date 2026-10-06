package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type OrderItemQueryCache interface {
	GetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*models.OrderItemRow, *int, bool)
	SetCachedOrderItemsAll(ctx context.Context, req *requests.FindAllOrderItems, data []*models.OrderItemRow, total *int)

	GetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*models.OrderItemActiveRow, *int, bool)
	SetCachedOrderItemActive(ctx context.Context, req *requests.FindAllOrderItems, data []*models.OrderItemActiveRow, total *int)

	GetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*models.OrderItemTrashedRow, *int, bool)
	SetCachedOrderItemTrashed(ctx context.Context, req *requests.FindAllOrderItems, data []*models.OrderItemTrashedRow, total *int)

	GetCachedOrderItems(ctx context.Context, orderID int) ([]*models.OrderItem, bool)
	SetCachedOrderItems(ctx context.Context, data []*models.OrderItem)
}

type OrderItemCommandCache interface {
	DeleteCachedOrderItems(ctx context.Context, orderID int)
	DeleteCachedOrderItemsAll(ctx context.Context)
}
