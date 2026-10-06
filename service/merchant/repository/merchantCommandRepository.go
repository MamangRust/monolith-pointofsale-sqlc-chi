package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	merchant_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
	"gorm.io/gorm"
)

type merchantCommandRepository struct {
	db *gorm.DB
}

func NewMerchantCommandRepository(db *gorm.DB) MerchantCommandRepository {
	return &merchantCommandRepository{db: db}
}

func (r *merchantCommandRepository) CreateMerchant(ctx context.Context, request *requests.CreateMerchantRequest) (*models.Merchant, error) {
	return r.createMerchant(ctx, r.db, request)
}

// CreateMerchantInTx persists the merchant inside the given GORM transaction
// so the caller can commit the business write and its outbox event atomically
// (Phase 6 — transactional outbox).
func (r *merchantCommandRepository) CreateMerchantInTx(ctx context.Context, tx *gorm.DB, request *requests.CreateMerchantRequest) (*models.Merchant, error) {
	return r.createMerchant(ctx, tx, request)
}

func (r *merchantCommandRepository) createMerchant(ctx context.Context, db *gorm.DB, request *requests.CreateMerchantRequest) (*models.Merchant, error) {
	merchant := &models.Merchant{
		UserID:       int32(request.UserID),
		Name:         request.Name,
		Description:  &request.Description,
		Address:      &request.Address,
		ContactEmail: &request.ContactEmail,
		ContactPhone: &request.ContactPhone,
		Status:       "inactive",
	}

	if err := r.db.WithContext(ctx).Create(merchant).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, merchant_errors.ErrCreateMerchant.WithInternal(err)
		}
		return nil, merchant_errors.ErrCreateMerchant.WithInternal(err)
	}

	return merchant, nil
}

func (r *merchantCommandRepository) UpdateMerchant(ctx context.Context, request *requests.UpdateMerchantRequest) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).
		Where("merchant_id = ?", *request.MerchantID).
		First(&merchant).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, merchant_errors.ErrUpdateMerchant.WithInternal(err)
		}
		return nil, merchant_errors.ErrUpdateMerchant.WithInternal(err)
	}

	merchant.Name = request.Name
	merchant.Description = &request.Description
	merchant.Address = &request.Address
	merchant.ContactEmail = &request.ContactEmail
	merchant.ContactPhone = &request.ContactPhone
	merchant.Status = request.Status
	merchant.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&merchant).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, merchant_errors.ErrUpdateMerchant.WithInternal(err)
		}
		return nil, merchant_errors.ErrUpdateMerchant.WithInternal(err)
	}

	return &merchant, nil
}

func (r *merchantCommandRepository) UpdateMerchantStatus(ctx context.Context, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error) {
	return r.updateMerchantStatus(ctx, r.db, request)
}

// UpdateMerchantStatusInTx updates the merchant status inside the given GORM
// transaction so the caller can commit the business write and its outbox event
// atomically (Phase 6 — transactional outbox).
func (r *merchantCommandRepository) UpdateMerchantStatusInTx(ctx context.Context, tx *gorm.DB, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error) {
	return r.updateMerchantStatus(ctx, tx, request)
}

func (r *merchantCommandRepository) updateMerchantStatus(ctx context.Context, db *gorm.DB, request *requests.UpdateMerchantStatusRequest) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).
		Where("merchant_id = ?", *request.MerchantID).
		First(&merchant).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, merchant_errors.ErrUpdateMerchantStatusFailed.WithInternal(err)
		}
		return nil, merchant_errors.ErrUpdateMerchantStatusFailed.WithInternal(err)
	}

	merchant.Status = request.Status
	merchant.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&merchant).Error; err != nil {
		return nil, merchant_errors.ErrUpdateMerchantStatusFailed.WithInternal(err)
	}

	return &merchant, nil
}

func (r *merchantCommandRepository) TrashedMerchant(ctx context.Context, merchantID int) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		First(&merchant).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, merchant_errors.ErrTrashedMerchant.WithInternal(err)
		}
		return nil, merchant_errors.ErrTrashedMerchant.WithInternal(err)
	}

	now := time.Now()
	merchant.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	merchant.UpdatedAt = now

	if err := r.db.WithContext(ctx).Save(&merchant).Error; err != nil {
		return nil, merchant_errors.ErrTrashedMerchant.WithInternal(err)
	}

	return &merchant, nil
}

func (r *merchantCommandRepository) RestoreMerchant(ctx context.Context, merchantID int) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("merchant_id = ?", merchantID).
		First(&merchant).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, merchant_errors.ErrRestoreMerchant.WithInternal(err)
		}
		return nil, merchant_errors.ErrRestoreMerchant.WithInternal(err)
	}

	merchant.DeletedAt = gorm.DeletedAt{Valid: false}
	merchant.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Unscoped().Save(&merchant).Error; err != nil {
		return nil, merchant_errors.ErrRestoreMerchant.WithInternal(err)
	}

	return &merchant, nil
}

func (r *merchantCommandRepository) DeleteMerchantPermanent(ctx context.Context, merchantID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("merchant_id = ? AND deleted_at IS NOT NULL", merchantID).
		Delete(&models.Merchant{})

	if result.Error != nil {
		return false, merchant_errors.ErrDeleteMerchantPermanent.WithInternal(result.Error)
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *merchantCommandRepository) RestoreAllMerchant(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Model(&models.Merchant{}).
		Where("deleted_at IS NOT NULL").
		Update("deleted_at", nil)

	if result.Error != nil {
		return false, merchant_errors.ErrRestoreAllMerchant.WithInternal(result.Error)
	}

	return true, nil
}

func (r *merchantCommandRepository) DeleteAllMerchantPermanent(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").
		Delete(&models.Merchant{})

	if result.Error != nil {
		return false, merchant_errors.ErrDeleteAllMerchantPermanent.WithInternal(result.Error)
	}

	return true, nil
}
