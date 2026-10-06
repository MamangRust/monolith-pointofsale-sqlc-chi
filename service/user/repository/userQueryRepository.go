package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	"gorm.io/gorm"
)

type userQueryRepository struct {
	db *gorm.DB
}

func NewUserQueryRepository(db *gorm.DB) UserQueryRepository {
	return &userQueryRepository{db: db}
}

func (r *userQueryRepository) FindAllUsers(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.UserRow
	query := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, verification_code, is_verified, created_at, updated_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("firstname ILIKE ? OR lastname ILIKE ? OR email ILIKE ?", search, search, search)
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, user_errors.ErrFindAllUsers.WithInternal(err)
	}

	return results, nil
}

func (r *userQueryRepository) FindById(ctx context.Context, user_id int) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", user_id).
		First(&user).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, user_errors.ErrInternalServerError.WithInternal(err)
	}

	return &user, nil
}

func (r *userQueryRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, nil
		}
		return nil, user_errors.ErrInternalServerError.WithInternal(err)
	}

	return &user, nil
}

// FindByEmailWithPassword returns the full user row (including the password
// hash) for credential checks. It is only used by the auth service over gRPC.
func (r *userQueryRepository) FindByEmailWithPassword(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, nil
		}
		return nil, user_errors.ErrInternalServerError.WithInternal(err)
	}

	return &user, nil
}

// FindByEmailAndVerify resolves a verified user by email for login.
func (r *userQueryRepository) FindByEmailAndVerify(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("email = ? AND is_verified = true", email).
		First(&user).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, nil
		}
		return nil, user_errors.ErrInternalServerError.WithInternal(err)
	}

	return &user, nil
}

// FindByVerificationCode resolves a user by their verification code.
func (r *userQueryRepository) FindByVerificationCode(ctx context.Context, verification_code string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("verification_code = ?", verification_code).
		First(&user).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, nil
		}
		return nil, user_errors.ErrInternalServerError.WithInternal(err)
	}

	return &user, nil
}

func (r *userQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.UserActiveRow
	query := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, verification_code, is_verified, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("firstname ILIKE ? OR lastname ILIKE ? OR email ILIKE ?", search, search, search)
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, user_errors.ErrFindActiveUsers.WithInternal(err)
	}

	return results, nil
}

func (r *userQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.UserTrashedRow
	query := r.db.WithContext(ctx).Table("users").
		Select("user_id, firstname, lastname, email, password, verification_code, is_verified, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("firstname ILIKE ? OR lastname ILIKE ? OR email ILIKE ?", search, search, search)
	}

	query = query.Where("deleted_at IS NOT NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, user_errors.ErrFindTrashedUsers.WithInternal(err)
	}

	return results, nil
}
