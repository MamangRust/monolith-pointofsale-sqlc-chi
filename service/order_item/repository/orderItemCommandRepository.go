package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/order_item_errors"
	"gorm.io/gorm"
)

type orderItemCommandRepository struct {
	db *gorm.DB
}

func NewOrderItemCommandRepository(db *gorm.DB) OrderItemCommandRepository {
	return &orderItemCommandRepository{db: db}
}

func (r *orderItemCommandRepository) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*models.OrderItem, error) {
	orderItem := &models.OrderItem{
		OrderID:   int32(req.OrderID),
		ProductID: int32(req.ProductID),
		Quantity:  int32(req.Quantity),
		Price:     int32(req.Price),
	}

	if err := r.db.WithContext(ctx).Create(orderItem).Error; err != nil {
		return nil, orderitem_errors.ErrCreateOrderItem.WithInternal(err)
	}

	return orderItem, nil
}

func (r *orderItemCommandRepository) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*models.OrderItem, error) {
	var orderItem models.OrderItem
	if err := r.db.WithContext(ctx).
		Where("order_item_id = ?", req.OrderItemID).
		First(&orderItem).Error; err != nil {
		return nil, orderitem_errors.ErrUpdateOrderItem.WithInternal(err)
	}

	orderItem.OrderID = int32(req.OrderID)
	orderItem.ProductID = int32(req.ProductID)
	orderItem.Quantity = int32(req.Quantity)
	orderItem.Price = int32(req.Price)
	orderItem.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&orderItem).Error; err != nil {
		return nil, orderitem_errors.ErrUpdateOrderItem.WithInternal(err)
	}

	return &orderItem, nil
}

// DeleteOrderItem hard-deletes a single order item by its primary key. It is
// used as a compensating action when an order creation is rolled back, so
// deleting an already-absent row is treated as success.
func (r *orderItemCommandRepository) DeleteOrderItem(ctx context.Context, orderItemID int) error {
	result := r.db.WithContext(ctx).Unscoped().
		Where("order_item_id = ?", orderItemID).
		Delete(&models.OrderItem{})

	if result.Error != nil {
		return orderitem_errors.ErrDeleteOrderItemPermanent.WithInternal(result.Error)
	}

	return nil
}
