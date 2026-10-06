package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"gorm.io/gorm"
)

type resetTokenRepository struct {
	db *gorm.DB
}

func NewResetTokenRepository(db *gorm.DB) ResetTokenRepository {
	return &resetTokenRepository{
		db: db,
	}
}

func (r *resetTokenRepository) FindByToken(ctx context.Context, code string) (*models.ResetToken, error) {
	var resetToken models.ResetToken
	if err := r.db.WithContext(ctx).
		Where("token = ?", code).
		First(&resetToken).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &resetToken, nil
}

func (r *resetTokenRepository) CreateResetToken(ctx context.Context, req *requests.CreateResetTokenRequest) (*models.ResetToken, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	resetToken := &models.ResetToken{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	}

	if err := r.db.WithContext(ctx).Create(resetToken).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return resetToken, nil
}

// CreateResetTokenInTx persists the reset token inside the given database
// transaction so the caller can commit the token write and its outbox event
// atomically (Phase 6).
func (r *resetTokenRepository) CreateResetTokenInTx(ctx context.Context, tx *gorm.DB, req *requests.CreateResetTokenRequest) (*models.ResetToken, error) {
	expiryDate, err := time.Parse("2006-01-02 15:04:05", req.ExpiredAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	resetToken := &models.ResetToken{
		UserID:     int64(req.UserID),
		Token:      req.ResetToken,
		ExpiryDate: expiryDate,
	}

	if err := tx.WithContext(ctx).Create(resetToken).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return resetToken, nil
}

func (r *resetTokenRepository) DeleteResetToken(ctx context.Context, user_id int) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", user_id).
		Delete(&models.ResetToken{}).Error; err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}
