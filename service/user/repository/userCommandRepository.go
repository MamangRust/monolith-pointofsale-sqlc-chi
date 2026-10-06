package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/user_errors"
	"gorm.io/gorm"
)

type userCommandRepository struct {
	db *gorm.DB
}

func NewUserCommandRepository(db *gorm.DB) UserCommandRepository {
	return &userCommandRepository{db: db}
}

func (r *userCommandRepository) CreateUser(ctx context.Context, request *requests.CreateUserRequest) (*models.User, error) {
	user := &models.User{
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, user_errors.ErrUserConflict.WithInternal(err)
		}
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}

	return user, nil
}

// CreateUserRecord persists a user whose password is already hashed by the
// caller (the auth service). The verification fields are auth-managed values.
func (r *userCommandRepository) CreateUserRecord(ctx context.Context, request *requests.RegisterRequest) (*models.User, error) {
	isVerified := request.IsVerified

	user := &models.User{
		Firstname:        request.FirstName,
		Lastname:         request.LastName,
		Email:            request.Email,
		Password:         request.Password,
		VerificationCode: request.VerifiedCode,
		IsVerified:       &isVerified,
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, user_errors.ErrUserConflict.WithInternal(err)
		}
		return nil, user_errors.ErrCreateUser.WithInternal(err)
	}

	return user, nil
}

func (r *userCommandRepository) UpdateUserIsVerified(ctx context.Context, user_id int, is_verified bool) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", user_id).
		First(&user).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, user_errors.ErrUpdateUserVerificationCode.WithInternal(err)
	}

	user.IsVerified = &is_verified
	user.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, user_errors.ErrUpdateUserVerificationCode.WithInternal(err)
	}

	return &user, nil
}

func (r *userCommandRepository) UpdateUserPassword(ctx context.Context, user_id int, password string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", user_id).
		First(&user).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, user_errors.ErrUpdateUserPassword.WithInternal(err)
	}

	user.Password = password
	user.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, user_errors.ErrUpdateUserPassword.WithInternal(err)
	}

	return &user, nil
}

func (r *userCommandRepository) UpdateUser(ctx context.Context, request *requests.UpdateUserRequest) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", *request.UserID).
		First(&user).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}

	user.Firstname = request.FirstName
	user.Lastname = request.LastName
	user.Email = request.Email
	user.Password = request.Password
	user.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, user_errors.ErrUserConflict.WithInternal(err)
		}
		return nil, user_errors.ErrUpdateUser.WithInternal(err)
	}

	return &user, nil
}

func (r *userCommandRepository) TrashedUser(ctx context.Context, user_id int) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", user_id).
		First(&user).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, user_errors.ErrTrashedUser.WithInternal(err)
	}

	now := time.Now()
	user.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	user.UpdatedAt = now

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, user_errors.ErrTrashedUser.WithInternal(err)
	}

	return &user, nil
}

func (r *userCommandRepository) RestoreUser(ctx context.Context, user_id int) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("user_id = ?", user_id).
		First(&user).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, user_errors.ErrUserNotFound.WithInternal(err)
		}
		return nil, user_errors.ErrRestoreUser.WithInternal(err)
	}

	user.DeletedAt = gorm.DeletedAt{Valid: false}
	user.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Unscoped().Save(&user).Error; err != nil {
		return nil, user_errors.ErrRestoreUser.WithInternal(err)
	}

	return &user, nil
}

func (r *userCommandRepository) DeleteUserPermanent(ctx context.Context, user_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("user_id = ? AND deleted_at IS NOT NULL", user_id).
		Delete(&models.User{})

	if result.Error != nil {
		return false, user_errors.ErrDeleteUserPermanent.WithInternal(result.Error)
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *userCommandRepository) RestoreAllUser(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Model(&models.User{}).
		Where("deleted_at IS NOT NULL").
		Update("deleted_at", nil)

	if result.Error != nil {
		return false, user_errors.ErrRestoreAllUsers.WithInternal(result.Error)
	}

	return true, nil
}

func (r *userCommandRepository) DeleteAllUserPermanent(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").
		Delete(&models.User{})

	if result.Error != nil {
		return false, user_errors.ErrDeleteAllUsers.WithInternal(result.Error)
	}

	return true, nil
}
