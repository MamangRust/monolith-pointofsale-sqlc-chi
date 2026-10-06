package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type ProductQueryRepository interface {
	FindAllProducts(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductRow, error)
	FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductTrashedRow, error)
	FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*models.ProductMerchantRow, error)
	FindByCategoryID(ctx context.Context, categoryID int32, req *requests.ProductByCategoryRequest) ([]*models.ProductCategoryNameRow, error)
	FindById(ctx context.Context, product_id int) (*models.Product, error)
	FindByIdTrashed(ctx context.Context, id int) (*models.Product, error)
}

type ProductCommandRepository interface {
	CreateProduct(ctx context.Context, request *requests.CreateProductRequest) (*models.Product, error)
	UpdateProduct(ctx context.Context, request *requests.UpdateProductRequest) (*models.Product, error)
	UpdateProductCountStock(ctx context.Context, product_id int, stock int) (*models.Product, error)
	AdjustProductCountStock(ctx context.Context, product_id int, delta int) (*models.Product, error)
	TrashedProduct(ctx context.Context, product_id int) (*models.Product, error)
	RestoreProduct(ctx context.Context, product_id int) (*models.Product, error)
	DeleteProductPermanent(ctx context.Context, product_id int) (bool, error)
	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductPermanent(ctx context.Context) (bool, error)
}
