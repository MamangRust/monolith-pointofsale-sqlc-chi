package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{
		db: db,
	}
}

func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.WithContext(ctx).
		Where("token = ?", token).
		First(&refreshToken).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) FindByUserId(ctx context.Context, user_id int) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", user_id).
		Order("created_at DESC").
		First(&refreshToken).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) CreateRefreshToken(ctx context.Context, req *requests.CreateRefreshToken) (*models.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	refreshToken := &models.RefreshToken{
		UserID:     int32(req.UserId),
		Token:      req.Token,
		Expiration: expirationTime,
	}

	if err := r.db.WithContext(ctx).Create(refreshToken).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return refreshToken, nil
}

func (r *refreshTokenRepository) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*models.RefreshToken, error) {
	layout := "2006-01-02 15:04:05"
	expirationTime, err := time.Parse(layout, req.ExpiresAt)
	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	result := r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("user_id = ? AND deleted_at IS NULL", req.UserId).
		Updates(map[string]interface{}{
			"token":      req.Token,
			"expiration": expirationTime,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(result.Error)
	}
	// The sqlc :one query returned pgx.ErrNoRows when no active token row was
	// found for the user; surface the equivalent not-found error.
	if result.RowsAffected == 0 {
		return nil, sharedErrors.ErrInternal.WithInternal(gorm.ErrRecordNotFound)
	}

	var refreshToken models.RefreshToken
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", req.UserId).
		Order("created_at DESC").
		First(&refreshToken).Error; err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	// The original sqlc query hard-deletes the row, so the soft-delete scope
	// of the model must be bypassed here.
	if err := r.db.WithContext(ctx).Unscoped().
		Where("token = ?", token).
		Delete(&models.RefreshToken{}).Error; err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}

func (r *refreshTokenRepository) DeleteRefreshTokenByUserId(ctx context.Context, user_id int) error {
	if err := r.db.WithContext(ctx).Unscoped().
		Where("user_id = ?", user_id).
		Delete(&models.RefreshToken{}).Error; err != nil {
		return sharedErrors.ErrInternal.WithInternal(err)
	}
	return nil
}
