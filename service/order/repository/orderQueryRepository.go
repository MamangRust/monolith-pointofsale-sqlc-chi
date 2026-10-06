package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	"gorm.io/gorm"
)

type orderQueryRepository struct {
	db *gorm.DB
}

func NewOrderQueryRepository(db *gorm.DB) OrderQueryRepository {
	return &orderQueryRepository{
		db: db,
	}
}

func (r *orderQueryRepository) FindAllOrders(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.OrderRow
	query := r.db.WithContext(ctx).Table("orders").
		Select("order_id, merchant_id, cashier_id, total_price, created_at, updated_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("order_id::TEXT ILIKE ? OR total_price::TEXT ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, order_errors.ErrFindAllOrders.WithInternal(err)
	}

	if len(results) > 0 {
		total := int(results[0].TotalCount)
		return results, &total, nil
	}
	return results, nil, nil
}

func (r *orderQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.OrderActiveRow
	query := r.db.WithContext(ctx).Table("orders").
		Select("order_id, merchant_id, cashier_id, total_price, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("order_id::TEXT ILIKE ? OR total_price::TEXT ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, order_errors.ErrFindByActive.WithInternal(err)
	}

	if len(results) > 0 {
		total := int(results[0].TotalCount)
		return results, &total, nil
	}
	return results, nil, nil
}

func (r *orderQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.OrderTrashedRow
	query := r.db.WithContext(ctx).Table("orders").
		Select("order_id, merchant_id, cashier_id, total_price, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("order_id::TEXT ILIKE ? OR total_price::TEXT ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NOT NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, order_errors.ErrFindByTrashed.WithInternal(err)
	}

	if len(results) > 0 {
		total := int(results[0].TotalCount)
		return results, &total, nil
	}
	return results, nil, nil
}

func (r *orderQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*models.OrderByMerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.OrderByMerchantRow
	query := r.db.WithContext(ctx).Table("orders").
		Select("order_id, merchant_id, cashier_id, total_price, created_at, updated_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("order_id::TEXT ILIKE ? OR total_price::TEXT ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NULL AND merchant_id = ?", req.MerchantID).
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, order_errors.ErrFindByMerchant.WithInternal(err)
	}

	if len(results) > 0 {
		total := int(results[0].TotalCount)
		return results, &total, nil
	}
	return results, nil, nil
}

func (r *orderQueryRepository) FindById(ctx context.Context, orderID int) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&order).Error; err != nil {
		return nil, order_errors.ErrFindById.WithInternal(err)
	}

	return &order, nil
}

func (r *orderQueryRepository) FindByTrashedId(ctx context.Context, orderID int) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("order_id = ?", orderID).
		First(&order).Error; err != nil {
		return nil, order_errors.ErrFindByTrashedId.WithInternal(err)
	}

	return &order, nil
}
