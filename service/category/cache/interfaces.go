package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CategoryQueryCache interface {
	GetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryRow, *int, bool)
	SetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory, data []*models.CategoryRow, total *int)

	GetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryActiveRow, *int, bool)
	SetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory, data []*models.CategoryActiveRow, total *int)

	GetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryTrashedRow, *int, bool)
	SetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory, data []*models.CategoryTrashedRow, total *int)

	GetCachedCategoryCache(ctx context.Context, id int) (*models.Category, bool)
	SetCachedCategoryCache(ctx context.Context, data *models.Category)
}

type CategoryCommandCache interface {
	DeleteCachedCategoryCache(ctx context.Context, id int)
	DeleteCachedCategoryAllCache(ctx context.Context)
}

type CategoryStatsCache interface {
	GetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice) ([]*models.CategoryMonthlyTotalPriceRow, bool)
	SetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice, data []*models.CategoryMonthlyTotalPriceRow)

	GetCachedYearTotalPriceCache(ctx context.Context, year int) ([]*models.CategoryYearlyTotalPriceRow, bool)
	SetCachedYearTotalPriceCache(ctx context.Context, year int, data []*models.CategoryYearlyTotalPriceRow)

	GetCachedMonthPriceCache(ctx context.Context, year int) ([]*models.CategoryMonthlyRow, bool)
	SetCachedMonthPriceCache(ctx context.Context, year int, data []*models.CategoryMonthlyRow)

	GetCachedYearPriceCache(ctx context.Context, year int) ([]*models.CategoryYearlyRow, bool)
	SetCachedYearPriceCache(ctx context.Context, year int, data []*models.CategoryYearlyRow)
}

type CategoryStatsByIdCache interface {
	GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*models.CategoryMonthlyTotalPriceByIdRow, bool)
	SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory, data []*models.CategoryMonthlyTotalPriceByIdRow)

	GetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*models.CategoryYearlyTotalPriceByIdRow, bool)
	SetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory, data []*models.CategoryYearlyTotalPriceByIdRow)

	GetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId) ([]*models.CategoryMonthlyByIdRow, bool)
	SetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId, data []*models.CategoryMonthlyByIdRow)

	GetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId) ([]*models.CategoryYearlyByIdRow, bool)
	SetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId, data []*models.CategoryYearlyByIdRow)
}

type CategoryStatsByMerchantCache interface {
	GetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant) ([]*models.CategoryMonthlyTotalPriceByMerchantRow, bool)
	SetCachedMonthTotalPriceByMerchantCache(ctx context.Context, req *requests.MonthTotalPriceMerchant, data []*models.CategoryMonthlyTotalPriceByMerchantRow)

	GetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant) ([]*models.CategoryYearlyTotalPriceByMerchantRow, bool)
	SetCachedYearTotalPriceByMerchantCache(ctx context.Context, req *requests.YearTotalPriceMerchant, data []*models.CategoryYearlyTotalPriceByMerchantRow)

	GetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant) ([]*models.CategoryMonthlyByMerchantRow, bool)
	SetCachedMonthPriceByMerchantCache(ctx context.Context, req *requests.MonthPriceMerchant, data []*models.CategoryMonthlyByMerchantRow)

	GetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant) ([]*models.CategoryYearlyByMerchantRow, bool)
	SetCachedYearPriceByMerchantCache(ctx context.Context, req *requests.YearPriceMerchant, data []*models.CategoryYearlyByMerchantRow)
}
