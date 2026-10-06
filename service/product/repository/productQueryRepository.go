package repository

import (
	"context"
	"strings"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
	"gorm.io/gorm"
)

type productQueryRepository struct {
	db *gorm.DB
}

func NewProductQueryRepository(db *gorm.DB) ProductQueryRepository {
	return &productQueryRepository{db: db}
}

func (r *productQueryRepository) FindAllProducts(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.ProductRow
	query := r.db.WithContext(ctx).Table("products AS p").
		Select("p.product_id, p.merchant_id, p.category_id, p.name, p.description, p.price, p.count_in_stock, p.brand, p.weight, p.slug_product, p.image_product, p.barcode, p.created_at, p.updated_at, p.deleted_at, COUNT(*) OVER () AS total_count").
		Where("p.deleted_at IS NULL")

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("(p.name ILIKE ? OR p.description ILIKE ? OR p.brand ILIKE ? OR p.slug_product ILIKE ? OR p.barcode ILIKE ?)",
			search, search, search, search, search)
	}

	if err := query.Order("p.created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Scan(&results).Error; err != nil {
		return nil, product_errors.ErrFindAllProducts.WithInternal(err)
	}

	return results, nil
}

func (r *productQueryRepository) FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductActiveRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.ProductActiveRow
	query := r.db.WithContext(ctx).Table("products AS p").
		Select("p.product_id, p.merchant_id, p.category_id, p.name, p.description, p.price, p.count_in_stock, p.brand, p.weight, p.slug_product, p.image_product, p.barcode, p.created_at, p.updated_at, p.deleted_at, COUNT(*) OVER () AS total_count").
		Where("p.deleted_at IS NULL")

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("(p.name ILIKE ? OR p.description ILIKE ? OR p.brand ILIKE ? OR p.slug_product ILIKE ? OR p.barcode ILIKE ?)",
			search, search, search, search, search)
	}

	if err := query.Order("p.created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Scan(&results).Error; err != nil {
		return nil, product_errors.ErrFindByActive.WithInternal(err)
	}

	return results, nil
}

func (r *productQueryRepository) FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductTrashedRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.ProductTrashedRow
	query := r.db.WithContext(ctx).Table("products AS p").
		Select("p.product_id, p.merchant_id, p.category_id, p.name, p.description, p.price, p.count_in_stock, p.brand, p.weight, p.slug_product, p.image_product, p.barcode, p.created_at, p.updated_at, p.deleted_at, COUNT(*) OVER () AS total_count").
		Where("p.deleted_at IS NOT NULL")

	if req.Search != "" {
		search := "%" + strings.TrimSpace(req.Search) + "%"
		query = query.Where("(p.name ILIKE ? OR p.description ILIKE ? OR p.brand ILIKE ? OR p.slug_product ILIKE ? OR p.barcode ILIKE ?)",
			search, search, search, search, search)
	}

	if err := query.Order("p.created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Scan(&results).Error; err != nil {
		return nil, product_errors.ErrFindByTrashed.WithInternal(err)
	}

	return results, nil
}

func (r *productQueryRepository) FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*models.ProductMerchantRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.ProductMerchantRow
	// category_name is left empty here and resolved by the service layer via
	// the Category adapter, so this repository never reads the categories table.
	query := r.db.WithContext(ctx).Table("products AS p").
		Select("p.product_id, p.category_id, p.name, p.description, p.price, p.count_in_stock, p.brand, p.image_product, p.created_at, COUNT(*) OVER () AS total_count").
		Where("p.deleted_at IS NULL").
		Where("p.merchant_id = ?", req.MerchantID)

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("(p.name ILIKE ? OR p.description ILIKE ?)", search, search)
	}

	if req.CategoryID != nil && *req.CategoryID != 0 {
		query = query.Where("p.category_id = ?", *req.CategoryID)
	}

	query = query.Where("p.price >= ? AND p.price <= ?", minPriceFilter(req.MinPrice), maxPriceFilter(req.MaxPrice))

	if err := query.Order("p.created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Scan(&results).Error; err != nil {
		return nil, product_errors.ErrFindByMerchant.WithInternal(err)
	}

	return results, nil
}

// FindByCategoryID returns products under a category id. The caller resolves
// the category name through the Category adapter first, so this repository
// never reads the categories table.
func (r *productQueryRepository) FindByCategoryID(ctx context.Context, categoryID int32, req *requests.ProductByCategoryRequest) ([]*models.ProductCategoryNameRow, error) {
	offset := (req.Page - 1) * req.PageSize

	var results []*models.ProductCategoryNameRow
	query := r.db.WithContext(ctx).Table("products AS p").
		Select("p.product_id, p.merchant_id, p.category_id, p.slug_product, p.weight, p.name, p.description, p.price, p.count_in_stock, p.brand, p.image_product, p.barcode, p.created_at, p.updated_at, p.deleted_at, COUNT(*) OVER () AS total_count").
		Where("p.deleted_at IS NULL").
		Where("p.category_id = ?", categoryID)

	if req.Search != "" {
		search := "%" + req.Search + "%"
		query = query.Where("(p.name ILIKE ? OR p.description ILIKE ?)", search, search)
	}

	query = query.Where("p.price >= ? AND p.price <= ?", minPriceFilter(req.MinPrice), maxPriceFilter(req.MaxPrice))

	if err := query.Order("p.created_at DESC").
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Scan(&results).Error; err != nil {
		return nil, product_errors.ErrFindByCategory.WithInternal(err)
	}

	return results, nil
}

// minPriceFilter maps an optional minimum-price filter to the lower bound used
// by the original sqlc queries (COALESCE(NULLIF($min, 0), 0)). nil and
// non-positive values mean "no lower bound".
func minPriceFilter(p *int) int {
	if p == nil || *p <= 0 {
		return 0
	}
	return *p
}

// maxPriceFilter maps an optional maximum-price filter to the upper bound used
// by the original sqlc queries (COALESCE(NULLIF($max, 0), 999999999)). nil and
// non-positive values mean "no upper bound".
func maxPriceFilter(p *int) int {
	if p == nil || *p <= 0 {
		return 999999999
	}
	return *p
}

func (r *productQueryRepository) FindById(ctx context.Context, product_id int) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", product_id).
		First(&product).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, product_errors.ErrFindById.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	return &product, nil
}

func (r *productQueryRepository) FindByIdTrashed(ctx context.Context, id int) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("product_id = ?", id).
		First(&product).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, product_errors.ErrFindByIdTrashed.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	return &product, nil
}
