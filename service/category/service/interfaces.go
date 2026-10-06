package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CategoryStatsService interface {
	FindMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*models.CategoryMonthlyTotalPriceRow, error)
	FindYearlyTotalPrice(ctx context.Context, year int) ([]*models.CategoryYearlyTotalPriceRow, error)

	FindMonthPrice(ctx context.Context, year int) ([]*models.CategoryMonthlyRow, error)
	FindYearPrice(ctx context.Context, year int) ([]*models.CategoryYearlyRow, error)
}

type CategoryStatsByIdService interface {
	FindMonthlyTotalPriceById(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*models.CategoryMonthlyTotalPriceByIdRow, error)
	FindYearlyTotalPriceById(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*models.CategoryYearlyTotalPriceByIdRow, error)

	FindMonthPriceById(ctx context.Context, req *requests.MonthPriceId) ([]*models.CategoryMonthlyByIdRow, error)
	FindYearPriceById(ctx context.Context, req *requests.YearPriceId) ([]*models.CategoryYearlyByIdRow, error)
}

type CategoryStatsByMerchantQueryService interface {
	FindMonthlyTotalPriceByMerchant(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*models.CategoryMonthlyTotalPriceByMerchantRow, error)
	FindYearlyTotalPriceByMerchant(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*models.CategoryYearlyTotalPriceByMerchantRow, error)

	FindMonthPriceByMerchant(ctx context.Context, req *requests.MonthPriceMerchant) ([]*models.CategoryMonthlyByMerchantRow, error)
	FindYearPriceByMerchant(ctx context.Context, req *requests.YearPriceMerchant) ([]*models.CategoryYearlyByMerchantRow, error)
}

type CategoryQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryRow, *int, error)
	FindById(ctx context.Context, category_id int) (*models.Category, error)
	FindByName(ctx context.Context, name string) (*models.Category, error)
	FindByIds(ctx context.Context, category_ids []int) ([]*models.Category, error)
	FindByActive(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryTrashedRow, *int, error)
}

type CategoryCommandService interface {
	CreateCategory(ctx context.Context, req *requests.CreateCategoryRequest) (*models.Category, error)
	UpdateCategory(ctx context.Context, req *requests.UpdateCategoryRequest) (*models.Category, error)
	TrashedCategory(ctx context.Context, category_id int) (*models.Category, error)
	RestoreCategory(ctx context.Context, categoryID int) (*models.Category, error)
	DeleteCategoryPermanent(ctx context.Context, categoryID int) (bool, error)
	RestoreAllCategories(ctx context.Context) (bool, error)
	DeleteAllCategoriesPermanent(ctx context.Context) (bool, error)
}
