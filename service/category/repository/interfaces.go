package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CategoryStatsRepository interface {
	GetMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*models.CategoryMonthlyTotalPriceRow, error)
	GetYearlyTotalPrices(ctx context.Context, year int) ([]*models.CategoryYearlyTotalPriceRow, error)

	GetMonthPrice(ctx context.Context, year int) ([]*models.CategoryMonthlyRow, error)
	GetYearPrice(ctx context.Context, year int) ([]*models.CategoryYearlyRow, error)
}

type CategoryStatsByIdRepository interface {
	GetMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*models.CategoryMonthlyTotalPriceByIdRow, error)
	GetYearlyTotalPricesById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*models.CategoryYearlyTotalPriceByIdRow, error)

	GetMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*models.CategoryMonthlyByIdRow, error)
	GetYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*models.CategoryYearlyByIdRow, error)
}

type CategoryStatsByMerchantRepository interface {
	GetMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*models.CategoryMonthlyTotalPriceByMerchantRow, error)
	GetYearlyTotalPricesByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*models.CategoryYearlyTotalPriceByMerchantRow, error)

	GetMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*models.CategoryMonthlyByMerchantRow, error)
	GetYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*models.CategoryYearlyByMerchantRow, error)
}

type CategoryQueryRepository interface {
	FindAllCategory(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryRow, error)
	FindById(ctx context.Context, category_id int) (*models.Category, error)
	FindByNameAndId(ctx context.Context, req *requests.CategoryNameAndId) (*models.Category, error)
	FindByName(ctx context.Context, name string) (*models.Category, error)
	FindByIds(ctx context.Context, category_ids []int) ([]*models.Category, error)

	FindByIdTrashed(ctx context.Context, category_id int) (*models.Category, error)

	FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryActiveRow, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryTrashedRow, error)
}

type CategoryCommandRepository interface {
	CreateCategory(ctx context.Context, request *requests.CreateCategoryRequest) (*models.Category, error)
	UpdateCategory(ctx context.Context, request *requests.UpdateCategoryRequest) (*models.Category, error)
	TrashedCategory(ctx context.Context, category_id int) (*models.Category, error)
	RestoreCategory(ctx context.Context, category_id int) (*models.Category, error)
	DeleteCategoryPermanently(ctx context.Context, category_id int) (bool, error)
	RestoreAllCategories(ctx context.Context) (bool, error)
	DeleteAllPermanentCategories(ctx context.Context) (bool, error)
}
