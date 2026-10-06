package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"gorm.io/gorm"
)

type cashierCommandRepository struct {
	db *gorm.DB
}

func NewCashierCommandRepository(db *gorm.DB) CashierCommandRepository {
	return &cashierCommandRepository{
		db: db,
	}
}

func (r *cashierCommandRepository) CreateCashier(ctx context.Context, request *requests.CreateCashierRequest) (*models.Cashier, error) {
	cashier := &models.Cashier{
		MerchantID: int32(request.MerchantID),
		UserID:     int32(request.UserID),
		Name:       request.Name,
	}

	if err := r.db.WithContext(ctx).Create(cashier).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, cashier_errors.ErrCreateCashier.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	return cashier, nil
}

func (r *cashierCommandRepository) UpdateCashier(ctx context.Context, request *requests.UpdateCashierRequest) (*models.Cashier, error) {
	var cashier models.Cashier
	if err := r.db.WithContext(ctx).
		Where("cashier_id = ?", *request.CashierID).
		First(&cashier).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, cashier_errors.ErrUpdateCashier.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	cashier.Name = request.Name
	cashier.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&cashier).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, cashier_errors.ErrUpdateCashier.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	return &cashier, nil
}

func (r *cashierCommandRepository) TrashedCashier(ctx context.Context, cashier_id int) (*models.Cashier, error) {
	var cashier models.Cashier
	if err := r.db.WithContext(ctx).
		Where("cashier_id = ?", cashier_id).
		First(&cashier).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, cashier_errors.ErrTrashedCashier.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	now := time.Now()
	cashier.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	cashier.UpdatedAt = now

	if err := r.db.WithContext(ctx).Save(&cashier).Error; err != nil {
		return nil, cashier_errors.ErrTrashedCashier.WithInternal(err)
	}

	return &cashier, nil
}

func (r *cashierCommandRepository) RestoreCashier(ctx context.Context, cashier_id int) (*models.Cashier, error) {
	var cashier models.Cashier
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("cashier_id = ?", cashier_id).
		First(&cashier).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, cashier_errors.ErrRestoreCashier.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	cashier.DeletedAt = gorm.DeletedAt{Valid: false}
	cashier.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Unscoped().Save(&cashier).Error; err != nil {
		return nil, cashier_errors.ErrRestoreCashier.WithInternal(err)
	}

	return &cashier, nil
}

func (r *cashierCommandRepository) DeleteCashierPermanent(ctx context.Context, cashier_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("cashier_id = ? AND deleted_at IS NOT NULL", cashier_id).
		Delete(&models.Cashier{})

	if result.Error != nil {
		return false, cashier_errors.ErrDeleteCashierPermanent.WithInternal(result.Error)
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *cashierCommandRepository) RestoreAllCashier(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Model(&models.Cashier{}).
		Where("deleted_at IS NOT NULL").
		Update("deleted_at", nil)

	if result.Error != nil {
		return false, cashier_errors.ErrRestoreAllCashiers.WithInternal(result.Error)
	}

	return true, nil
}

func (r *cashierCommandRepository) DeleteAllCashierPermanent(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").
		Delete(&models.Cashier{})

	if result.Error != nil {
		return false, cashier_errors.ErrDeleteAllCashiersPermanent.WithInternal(result.Error)
	}

	return true, nil
}
