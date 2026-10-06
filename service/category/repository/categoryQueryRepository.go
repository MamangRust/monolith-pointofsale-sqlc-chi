package repository

import (
	"context"
	"strings"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/category_errors"
	"gorm.io/gorm"
)

type categoryQueryRepository struct {
	db *gorm.DB
}

func NewCategoryQueryRepository(db *gorm.DB) CategoryQueryRepository {
	return &categoryQueryRepository{db: db}
}

func (r *categoryQueryRepository) FindAllCategory(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.CategoryRow
	query := r.db.WithContext(ctx).Table("categories").
		Select("category_id, name, description, slug_category, created_at, updated_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("name ILIKE ? OR slug_category ILIKE ?", "%"+strings.TrimSpace(req.Search)+"%", "%"+strings.TrimSpace(req.Search)+"%")
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, category_errors.ErrFindAllCategory.WithInternal(err)
	}

	return results, nil
}

func (r *categoryQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.CategoryActiveRow
	query := r.db.WithContext(ctx).Table("categories").
		Select("category_id, name, description, slug_category, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("name ILIKE ? OR slug_category ILIKE ?", "%"+strings.TrimSpace(req.Search)+"%", "%"+strings.TrimSpace(req.Search)+"%")
	}

	query = query.Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, category_errors.ErrFindByActive.WithInternal(err)
	}

	return results, nil
}

func (r *categoryQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.CategoryTrashedRow
	query := r.db.WithContext(ctx).Table("categories").
		Select("category_id, name, description, slug_category, created_at, updated_at, deleted_at, COUNT(*) OVER () AS total_count")

	if req.Search != "" {
		query = query.Where("name ILIKE ? OR slug_category ILIKE ?", "%"+strings.TrimSpace(req.Search)+"%", "%"+strings.TrimSpace(req.Search)+"%")
	}

	query = query.Where("deleted_at IS NOT NULL").
		Order("created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset))

	if err := query.Scan(&results).Error; err != nil {
		return nil, category_errors.ErrFindByTrashed.WithInternal(err)
	}

	return results, nil
}

func (r *categoryQueryRepository) FindById(ctx context.Context, category_id int) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).
		Where("category_id = ?", category_id).
		First(&category).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, category_errors.ErrFindById.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	return &category, nil
}

func (r *categoryQueryRepository) FindByIdTrashed(ctx context.Context, category_id int) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("category_id = ? AND deleted_at IS NOT NULL", category_id).
		First(&category).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, category_errors.ErrFindById.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	return &category, nil
}

func (r *categoryQueryRepository) FindByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&category).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, category_errors.ErrFindByName.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	return &category, nil
}

func (r *categoryQueryRepository) FindByIds(ctx context.Context, category_ids []int) ([]*models.Category, error) {
	if len(category_ids) == 0 {
		return nil, nil
	}

	var categories []*models.Category
	if err := r.db.WithContext(ctx).
		Where("category_id IN ?", category_ids).
		Where("deleted_at IS NULL").
		Find(&categories).Error; err != nil {
		return nil, category_errors.ErrFindByIds.WithInternal(err)
	}

	return categories, nil
}

func (r *categoryQueryRepository) FindByNameAndId(ctx context.Context, req *requests.CategoryNameAndId) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).
		Where("name = ? AND category_id = ?", req.Name, req.CategoryID).
		First(&category).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, category_errors.ErrFindByNameAndId.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	return &category, nil
}
