package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type ProductQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*models.ProductMerchantRow, *int, error)
	FindByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*models.ProductCategoryNameRow, *int, error)
	FindById(ctx context.Context, productID int) (*models.Product, error)
	FindByActive(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductTrashedRow, *int, error)
}

type ProductCommandService interface {
	CreateProduct(ctx context.Context, req *requests.CreateProductRequest) (*models.Product, error)
	UpdateProduct(ctx context.Context, req *requests.UpdateProductRequest) (*models.Product, error)
	DecrementStock(ctx context.Context, productID int, quantity int) (*models.Product, error)
	IncrementStock(ctx context.Context, productID int, quantity int) (*models.Product, error)
	TrashProduct(ctx context.Context, productID int) (*models.Product, error)
	RestoreProduct(ctx context.Context, productID int) (*models.Product, error)
	DeleteProductPermanent(ctx context.Context, productID int) (bool, error)
	RestoreAllProducts(ctx context.Context) (bool, error)
	DeleteAllProductsPermanent(ctx context.Context) (bool, error)
}
