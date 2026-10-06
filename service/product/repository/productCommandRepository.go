package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/gormerr"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/product_errors"
	"gorm.io/gorm"
)

type productCommandRepository struct {
	db *gorm.DB
}

func NewProductCommandRepository(db *gorm.DB) ProductCommandRepository {
	return &productCommandRepository{db: db}
}

func (r *productCommandRepository) CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*models.Product, error) {
	weight := int32(request.Weight)

	product := &models.Product{
		MerchantID:   int32(request.MerchantID),
		CategoryID:   int32(request.CategoryID),
		Name:         request.Name,
		Description:  &request.Description,
		Price:        int32(request.Price),
		CountInStock: int32(request.CountInStock),
		Brand:        &request.Brand,
		Weight:       &weight,
		SlugProduct:  request.SlugProduct,
		ImageProduct: &request.ImageProduct,
		Barcode:      request.Barcode,
	}

	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, product_errors.ErrCreateProduct.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	return product, nil
}

func (r *productCommandRepository) UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", *request.ProductID).
		First(&product).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, product_errors.ErrUpdateProduct.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	weight := int32(request.Weight)

	product.CategoryID = int32(request.CategoryID)
	product.Name = request.Name
	product.Description = &request.Description
	product.Price = int32(request.Price)
	product.CountInStock = int32(request.CountInStock)
	product.Brand = &request.Brand
	product.Weight = &weight
	product.ImageProduct = &request.ImageProduct
	product.Barcode = request.Barcode
	product.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&product).Error; err != nil {
		if gormerr.IsDuplicateKey(err) {
			return nil, product_errors.ErrUpdateProduct.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	return &product, nil
}

func (r *productCommandRepository) UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", product_id).
		First(&product).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	product.CountInStock = int32(stock)
	product.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Save(&product).Error; err != nil {
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
	}

	return &product, nil
}

// AdjustProductCountStock applies a relative stock change (delta) in a single
// conditional UPDATE so concurrent adjustments cannot oversell. A decrement
// that would drive the stock negative reports ErrInsufficientStock.
func (r *productCommandRepository) AdjustProductCountStock(ctx context.Context, product_id int, delta int) (*models.Product, error) {
	query := r.db.WithContext(ctx).
		Model(&models.Product{}).
		Where("product_id = ?", product_id)

	if delta < 0 {
		query = query.Where("count_in_stock >= ?", -delta)
	}

	result := query.UpdateColumn("count_in_stock", gorm.Expr("count_in_stock + ?", delta))
	if result.Error != nil {
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(result.Error)
	}

	if result.RowsAffected == 0 {
		var existing models.Product
		if err := r.db.WithContext(ctx).Where("product_id = ?", product_id).First(&existing).Error; err != nil {
			if gormerr.IsRecordNotFound(err) {
				return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
			}
			return nil, errors.ErrInternal.WithInternal(err)
		}
		return nil, product_errors.ErrInsufficientStock
	}

	var product models.Product
	if err := r.db.WithContext(ctx).Where("product_id = ?", product_id).First(&product).Error; err != nil {
		return nil, product_errors.ErrUpdateProductCountStock.WithInternal(err)
	}

	return &product, nil
}

func (r *productCommandRepository) TrashedProduct(ctx context.Context, product_id int) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", product_id).
		First(&product).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, product_errors.ErrTrashedProduct.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	now := time.Now()
	product.DeletedAt = gorm.DeletedAt{Time: now, Valid: true}
	product.UpdatedAt = now

	if err := r.db.WithContext(ctx).Save(&product).Error; err != nil {
		return nil, product_errors.ErrTrashedProduct.WithInternal(err)
	}

	return &product, nil
}

func (r *productCommandRepository) RestoreProduct(ctx context.Context, product_id int) (*models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).
		Unscoped().
		Where("product_id = ?", product_id).
		First(&product).Error; err != nil {
		if gormerr.IsRecordNotFound(err) {
			return nil, product_errors.ErrRestoreProduct.WithInternal(err)
		}
		return nil, errors.ErrInternal.WithInternal(err)
	}

	product.DeletedAt = gorm.DeletedAt{Valid: false}
	product.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Unscoped().Save(&product).Error; err != nil {
		return nil, product_errors.ErrRestoreProduct.WithInternal(err)
	}

	return &product, nil
}

func (r *productCommandRepository) DeleteProductPermanent(ctx context.Context, product_id int) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("product_id = ? AND deleted_at IS NOT NULL", product_id).
		Delete(&models.Product{})

	if result.Error != nil {
		return false, product_errors.ErrDeleteProductPermanent.WithInternal(result.Error)
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (r *productCommandRepository) RestoreAllProducts(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Model(&models.Product{}).
		Where("deleted_at IS NOT NULL").
		Update("deleted_at", nil)

	if result.Error != nil {
		return false, product_errors.ErrRestoreAllProducts.WithInternal(result.Error)
	}

	return true, nil
}

func (r *productCommandRepository) DeleteAllProductPermanent(ctx context.Context) (bool, error) {
	result := r.db.WithContext(ctx).Unscoped().
		Where("deleted_at IS NOT NULL").
		Delete(&models.Product{})

	if result.Error != nil {
		return false, product_errors.ErrDeleteAllProductPermanent.WithInternal(result.Error)
	}

	return true, nil
}
