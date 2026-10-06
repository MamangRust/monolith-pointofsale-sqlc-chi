package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"gorm.io/gorm"
)

type orderItemQueryRepository struct {
	db *gorm.DB
}

func NewOrderItemQueryRepository(db *gorm.DB) OrderItemQueryRepository {
	return &orderItemQueryRepository{db: db}
}

func (r *orderItemQueryRepository) FindAllOrderItems(ctx context.Context, req *requests.FindAllOrderItems) ([]*models.OrderItemRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.OrderItemRow
	query := r.db.WithContext(ctx).Table("order_items").
		Select("order_item_id, order_id, product_id, quantity, price, created_at, updated_at, COUNT(*) OVER() AS total_count").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		query = query.Where("order_item_id::TEXT ILIKE ?", "%"+req.Search+"%")
	}

	query = query.Order("created_at DESC").Limit(int(req.PageSize)).Offset(int(offset))
	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}
	return results, &total, nil
}

func (r *orderItemQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*models.OrderItemActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.OrderItemActiveRow
	query := r.db.WithContext(ctx).Table("order_items").
		Select("order_item_id, order_id, product_id, quantity, price, created_at, updated_at, deleted_at, COUNT(*) OVER() AS total_count").
		Where("deleted_at IS NULL")

	if req.Search != "" {
		query = query.Where("order_item_id::TEXT ILIKE ?", "%"+req.Search+"%")
	}

	query = query.Order("created_at DESC").Limit(int(req.PageSize)).Offset(int(offset))
	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}
	return results, &total, nil
}

func (r *orderItemQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*models.OrderItemTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize
	var results []*models.OrderItemTrashedRow
	query := r.db.WithContext(ctx).Table("order_items").
		Select("order_item_id, order_id, product_id, quantity, price, created_at, updated_at, deleted_at, COUNT(*) OVER() AS total_count").
		Where("deleted_at IS NOT NULL")

	if req.Search != "" {
		query = query.Where("order_item_id::TEXT ILIKE ?", "%"+req.Search+"%")
	}

	query = query.Order("created_at DESC").Limit(int(req.PageSize)).Offset(int(offset))
	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}
	return results, &total, nil
}

func (r *orderItemQueryRepository) FindOrderItemByOrder(ctx context.Context, order_id int) ([]*models.OrderItem, error) {
	var items []*models.OrderItem
	if err := r.db.WithContext(ctx).
		Where("order_id = ? AND deleted_at IS NULL", order_id).
		Find(&items).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, nil
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return items, nil
}
