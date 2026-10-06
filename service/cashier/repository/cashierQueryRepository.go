package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"gorm.io/gorm"
)

type cashierQueryRepository struct {
	db *gorm.DB
}

func NewCashierQueryRepository(db *gorm.DB) CashierQueryRepository {
	return &cashierQueryRepository{
		db: db,
	}
}

func (r *cashierQueryRepository) FindAllCashiers(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.CashierRow
	query := r.db.WithContext(ctx).Table("cashiers").
		Select("cashier_id, merchant_id, user_id, name, created_at, updated_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, cashier_errors.ErrFindAllCashiers.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}

	return results, &total, nil
}

func (r *cashierQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.CashierActiveRow
	query := r.db.WithContext(ctx).Table("cashiers").
		Select("cashier_id, merchant_id, user_id, name, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, cashier_errors.ErrFindActiveCashiers.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}

	return results, &total, nil
}

func (r *cashierQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.CashierTrashedRow
	query := r.db.WithContext(ctx).Table("cashiers").
		Select("cashier_id, merchant_id, user_id, name, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NOT NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, cashier_errors.ErrFindTrashedCashiers.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}

	return results, &total, nil
}

func (r *cashierQueryRepository) FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*models.CashierByMerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.CashierByMerchantRow
	query := r.db.WithContext(ctx).Table("cashiers").
		Select("cashier_id, merchant_id, user_id, name, created_at, updated_at, COUNT(*) OVER () AS total_count").
		Where("merchant_id = ?", req.MerchantID)

	if req.Search != "" {
		query = query.Where("name ILIKE ?", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, cashier_errors.ErrFindCashiersByMerchant.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}

	return results, &total, nil
}

func (r *cashierQueryRepository) FindById(ctx context.Context, cashier_id int) (*models.Cashier, error) {
	var cashier models.Cashier
	if err := r.db.WithContext(ctx).
		Where("cashier_id = ?", cashier_id).
		First(&cashier).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, cashier_errors.ErrFindCashierById.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	return &cashier, nil
}
