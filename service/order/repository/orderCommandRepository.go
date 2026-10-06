package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormtx"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	"gorm.io/gorm"
)

type orderCommandRepository struct {
	db *gorm.DB
}

func NewOrderCommandRepository(db *gorm.DB) OrderCommandRepository {
	return &orderCommandRepository{
		db: db,
	}
}

func (r *orderCommandRepository) DeleteOrder(ctx context.Context, orderID int) error {
	result := r.db.WithContext(ctx).Unscoped().
		Where("order_id = ? AND deleted_at IS NULL", orderID).
		Delete(&models.Order{})
	if result.Error != nil {
		return order_errors.ErrDeleteOrderPermanent
	}
	return nil
}

func (r *orderCommandRepository) CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*models.Order, error) {
	order := &models.Order{
		MerchantID: int32(request.MerchantID),
		CashierID:  int32(request.CashierID),
		TotalPrice: int64(request.TotalPrice),
	}

	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, order_errors.ErrCreateOrder.WithInternal(err)
		}
		return nil, order_errors.ErrCreateOrder.WithInternal(err)
	}

	return order, nil
}

func (r *orderCommandRepository) FindAllTrashed(ctx context.Context) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("deleted_at IS NOT NULL").
		Order("order_id ASC").
		Find(&orders).Error; err != nil {
		return nil, order_errors.ErrFindAllTrashed.WithInternal(err)
	}

	return orders, nil
}

func (r *orderCommandRepository) UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", request.OrderID).
		First(&order).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, order_errors.ErrUpdateOrder.WithInternal(err)
		}
		return nil, order_errors.ErrUpdateOrder.WithInternal(err)
	}

	order.TotalPrice = int64(request.TotalPrice)
	order.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&order).Error; err != nil {
		return nil, order_errors.ErrUpdateOrder.WithInternal(err)
	}

	return &order, nil
}

func (r *orderCommandRepository) TrashedOrder(ctx context.Context, orderID int) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&order).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, order_errors.ErrTrashedOrder.WithInternal(err)
		}
		return nil, order_errors.ErrTrashedOrder.WithInternal(err)
	}

	now := time.Now()
	order.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	order.UpdatedAt = now

	if err := r.db.WithContext(ctx).Save(&order).Error; err != nil {
		return nil, order_errors.ErrTrashedOrder.WithInternal(err)
	}

	return &order, nil
}

func (r *orderCommandRepository) RestoreOrder(ctx context.Context, orderID int) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("order_id = ?", orderID).
		First(&order).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, order_errors.ErrRestoreOrderNotFound.WithInternal(err)
		}
		return nil, order_errors.ErrRestoreOrder.WithInternal(err)
	}

	if !order.DeletedAt.Valid {
		return nil, order_errors.ErrRestoreOrderNotFound
	}

	order.DeletedAt = gorm.DeletedAt{Valid: false}
	order.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Unscoped().Save(&order).Error; err != nil {
		return nil, order_errors.ErrRestoreOrder.WithInternal(err)
	}

	return &order, nil
}

func (r *orderCommandRepository) DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	err := gormtx.RunInTx(ctx, r.db, func(tx *gorm.DB) error {
		if err := tx.Unscoped().
			Where("order_id = ?", orderID).
			Delete(&models.Transaction{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().
			Where("order_id = ?", orderID).
			Delete(&models.OrderItem{}).Error; err != nil {
			return err
		}

		result := tx.Unscoped().
			Where("order_id = ? AND deleted_at IS NOT NULL", orderID).
			Delete(&models.Order{})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return order_errors.ErrDeleteOrderPermanentNotFound
		}

		return nil
	})
	if err != nil {
		if gormerr.IsRecordNotFound(err) || isNotFoundDomainErr(err) {
			return false, order_errors.ErrDeleteOrderPermanentNotFound
		}
		return false, order_errors.ErrDeleteOrderPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *orderCommandRepository) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	var deleted bool

	err := gormtx.RunInTx(ctx, r.db, func(tx *gorm.DB) error {
		targets := tx.Model(&models.Order{}).
			Select("order_id").
			Where("deleted_at IS NOT NULL")

		if err := tx.Unscoped().
			Where("order_id IN (?)", targets).
			Delete(&models.Transaction{}).Error; err != nil {
			return err
		}

		if err := tx.Unscoped().
			Where("order_id IN (?)", targets).
			Delete(&models.OrderItem{}).Error; err != nil {
			return err
		}

		result := tx.Unscoped().
			Where("deleted_at IS NOT NULL").
			Delete(&models.Order{})
		if result.Error != nil {
			return result.Error
		}

		deleted = result.RowsAffected > 0
		return nil
	})
	if err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}

	return deleted, nil
}

// isNotFoundDomainErr reports whether the error chain carries the shared
// not-found AppError used to signal a missing trashed order.
func isNotFoundDomainErr(err error) bool {
	type appErr interface{ GetCode() int }
	return errorsAsNotFound(err)
}

func errorsAsNotFound(err error) bool {
	return gormerr.IsRecordNotFound(err)
}
