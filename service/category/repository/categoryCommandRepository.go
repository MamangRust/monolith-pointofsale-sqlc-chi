package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
	"gorm.io/gorm"
)

type categoryCommandRepository struct {
	db *gorm.DB
}

func NewCategoryCommandRepository(db *gorm.DB) CategoryCommandRepository {
	return &categoryCommandRepository{db: db}
}

func (r *categoryCommandRepository) CreateCategory(ctx context.Context, request *requests.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:         request.Name,
		Description:  &request.Description,
		SlugCategory: request.SlugCategory,
	}

	// Note: category_errors has no dedicated conflict error, so duplicate-key
	// violations (e.g. unique slug_category) surface as the create failure.
	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		return nil, category_errors.ErrCreateCategory.WithInternal(err)
	}

	return category, nil
}

func (r *categoryCommandRepository) UpdateCategory(ctx context.Context, request *requests.UpdateCategoryRequest) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).
		Where("category_id = ?", *request.CategoryID).
		First(&category).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, category_errors.ErrUpdateCategory.WithInternal(err)
		}
		return nil, category_errors.ErrUpdateCategory.WithInternal(err)
	}

	category.Name = request.Name
	category.Description = &request.Description
	category.SlugCategory = request.SlugCategory
	category.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&category).Error; err != nil {
		return nil, category_errors.ErrUpdateCategory.WithInternal(err)
	}

	return &category, nil
}

func (r *categoryCommandRepository) TrashedCategory(ctx context.Context, category_id int) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).
		Where("category_id = ?", category_id).
		First(&category).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, category_errors.ErrTrashedCategory.WithInternal(err)
		}
		return nil, category_errors.ErrTrashedCategory.WithInternal(err)
	}

	now := time.Now()
	category.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	category.UpdatedAt = now

	if err := r.db.WithContext(ctx).Save(&category).Error; err != nil {
		return nil, category_errors.ErrTrashedCategory.WithInternal(err)
	}

	return &category, nil
}

func (r *categoryCommandRepository) RestoreCategory(ctx context.Context, category_id int) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("category_id = ?", category_id).
		First(&category).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, category_errors.ErrRestoreCategory.WithInternal(err)
		}
		return nil, category_errors.ErrRestoreCategory.WithInternal(err)
	}

	category.DeletedAt = gorm.DeletedAt{Valid: false}
	category.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Unscoped().Save(&category).Error; err != nil {
		return nil, category_errors.ErrRestoreCategory.WithInternal(err)
	}

	return &category, nil
}

func (r *categoryCommandRepository) DeleteCategoryPermanently(ctx context.Context, category_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("category_id = ? AND deleted_at IS NOT NULL", category_id).
		Delete(&models.Category{})

	if result.Error != nil {
		return false, category_errors.ErrDeleteCategoryPermanently.WithInternal(result.Error)
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *categoryCommandRepository) RestoreAllCategories(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Model(&models.Category{}).
		Where("deleted_at IS NOT NULL").
		Update("deleted_at", nil)

	if result.Error != nil {
		return false, category_errors.ErrRestoreAllCategories.WithInternal(result.Error)
	}

	return true, nil
}

func (r *categoryCommandRepository) DeleteAllPermanentCategories(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").
		Delete(&models.Category{})

	if result.Error != nil {
		return false, category_errors.ErrDeleteAllPermanentCategories.WithInternal(result.Error)
	}

	return true, nil
}
