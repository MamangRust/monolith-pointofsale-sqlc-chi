package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	merchant_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/merchant_errors"
	"gorm.io/gorm"
)

type merchantQueryRepository struct {
	db *gorm.DB
}

func NewMerchantQueryRepository(db *gorm.DB) MerchantQueryRepository {
	return &merchantQueryRepository{db: db}
}

func (r *merchantQueryRepository) FindAllMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.MerchantRow
	query := r.db.WithContext(ctx).Table("merchants").
		Select("merchant_id, user_id, name, description, address, contact_email, contact_phone, status, created_at, updated_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("name ILIKE ? OR contact_email ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, merchant_errors.ErrFindAllMerchants.WithInternal(err)
	}

	if len(results) > 0 {
		total := int(results[0].TotalCount)
		return results, &total, nil
	}
	return results, nil, nil
}

func (r *merchantQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantActiveRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.MerchantActiveRow
	query := r.db.WithContext(ctx).Table("merchants").
		Select("merchant_id, user_id, name, description, address, contact_email, contact_phone, status, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("name ILIKE ? OR contact_email ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, merchant_errors.ErrFindByActive.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}

	return results, &total, nil
}

func (r *merchantQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantTrashedRow, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.MerchantTrashedRow
	query := r.db.WithContext(ctx).Table("merchants").
		Select("merchant_id, user_id, name, description, address, contact_email, contact_phone, status, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("name ILIKE ? OR contact_email ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	query = query.Where("deleted_at IS NOT NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, nil, merchant_errors.ErrFindByTrashed.WithInternal(err)
	}

	var total int
	if len(results) > 0 {
		total = int(results[0].TotalCount)
	}

	return results, &total, nil
}

func (r *merchantQueryRepository) FindById(ctx context.Context, merchantID int) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.db.WithContext(ctx).
		Where("merchant_id = ?", merchantID).
		First(&merchant).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, merchant_errors.ErrFindById.WithInternal(err)
		}
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return &merchant, nil
}
