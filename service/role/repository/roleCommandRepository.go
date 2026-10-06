package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	role_errors "github.com/MamangRust/monolith-point-of-sale-shared/errors/role_errors"
	"gorm.io/gorm"
)

type roleCommandRepository struct {
	db *gorm.DB
}

func NewRoleCommandRepository(db *gorm.DB) RoleCommandRepository {
	return &roleCommandRepository{db: db}
}

func (r *roleCommandRepository) CreateRole(ctx context.Context, req *requests.CreateRoleRequest) (*models.Role, error) {
	role := &models.Role{
		RoleName: req.Name,
	}

	if err := r.db.WithContext(ctx).Create(role).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, role_errors.ErrRoleConflict.WithInternal(err)
		}
		return nil, role_errors.ErrCreateRole.WithInternal(err)
	}

	return role, nil
}

func (r *roleCommandRepository) UpdateRole(ctx context.Context, req *requests.UpdateRoleRequest) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).
		Where("role_id = ?", *req.ID).
		First(&role).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, role_errors.ErrUpdateRole.WithInternal(err)
	}

	role.RoleName = req.Name
	role.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&role).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, role_errors.ErrRoleConflict.WithInternal(err)
		}
		return nil, role_errors.ErrUpdateRole.WithInternal(err)
	}

	return &role, nil
}

func (r *roleCommandRepository) TrashedRole(ctx context.Context, id int) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).
		Where("role_id = ?", id).
		First(&role).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, role_errors.ErrTrashedRole.WithInternal(err)
	}

	now := time.Now()
	role.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	role.UpdatedAt = now

	if err := r.db.WithContext(ctx).Save(&role).Error; err != nil {
		return nil, role_errors.ErrTrashedRole.WithInternal(err)
	}

	return &role, nil
}

func (r *roleCommandRepository) RestoreRole(ctx context.Context, id int) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("role_id = ?", id).
		First(&role).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, role_errors.ErrRoleNotFound.WithInternal(err)
		}
		return nil, role_errors.ErrRestoreRole.WithInternal(err)
	}

	role.DeletedAt = gorm.DeletedAt{Valid: false}
	role.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Unscoped().Save(&role).Error; err != nil {
		return nil, role_errors.ErrRestoreRole.WithInternal(err)
	}

	return &role, nil
}

func (r *roleCommandRepository) DeleteRolePermanent(ctx context.Context, roleID int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("role_id = ? AND deleted_at IS NOT NULL", roleID).
		Delete(&models.Role{})

	if result.Error != nil {
		return false, role_errors.ErrDeleteRolePermanent.WithInternal(result.Error)
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *roleCommandRepository) RestoreAllRole(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Model(&models.Role{}).
		Where("deleted_at IS NOT NULL").
		Update("deleted_at", nil)

	if result.Error != nil {
		return false, role_errors.ErrRestoreAllRoles.WithInternal(result.Error)
	}

	return true, nil
}

func (r *roleCommandRepository) DeleteAllRolePermanent(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").
		Delete(&models.Role{})

	if result.Error != nil {
		return false, role_errors.ErrDeleteAllRoles.WithInternal(result.Error)
	}

	return true, nil
}
