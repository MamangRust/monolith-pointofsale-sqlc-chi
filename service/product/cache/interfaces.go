package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type ProductQueryCache interface {
	GetCachedProducts(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductRow, *int, bool)
	SetCachedProducts(ctx context.Context, req *requests.FindAllProducts, data []*models.ProductRow, total *int)

	GetCachedProductsByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest) ([]*models.ProductMerchantRow, *int, bool)
	SetCachedProductsByMerchant(ctx context.Context, req *requests.ProductByMerchantRequest, data []*models.ProductMerchantRow, total *int)

	GetCachedProductsByCategory(ctx context.Context, req *requests.ProductByCategoryRequest) ([]*models.ProductCategoryNameRow, *int, bool)
	SetCachedProductsByCategory(ctx context.Context, req *requests.ProductByCategoryRequest, data []*models.ProductCategoryNameRow, total *int)

	GetCachedProductActive(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductActiveRow, *int, bool)
	SetCachedProductActive(ctx context.Context, req *requests.FindAllProducts, data []*models.ProductActiveRow, total *int)

	GetCachedProductTrashed(ctx context.Context, req *requests.FindAllProducts) ([]*models.ProductTrashedRow, *int, bool)
	SetCachedProductTrashed(ctx context.Context, req *requests.FindAllProducts, data []*models.ProductTrashedRow, total *int)

	GetCachedProduct(ctx context.Context, productID int) (*models.Product, bool)
	SetCachedProduct(ctx context.Context, data *models.Product)
}

type ProductCommandCache interface {
	DeleteCachedProduct(ctx context.Context, productID int)
	DeleteCachedProductAllCache(ctx context.Context)
}
