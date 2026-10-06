package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CashierQueryCache interface {
	GetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierRow, *int, bool)
	SetCachedCashiersCache(ctx context.Context, req *requests.FindAllCashiers, res []*models.CashierRow, total *int)

	GetCachedCashier(ctx context.Context, cashierID int) (*models.Cashier, bool)
	SetCachedCashier(ctx context.Context, res *models.Cashier)

	GetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierActiveRow, *int, bool)
	SetCachedCashiersActive(ctx context.Context, req *requests.FindAllCashiers, res []*models.CashierActiveRow, total *int)

	GetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierTrashedRow, *int, bool)
	SetCachedCashiersTrashed(ctx context.Context, req *requests.FindAllCashiers, res []*models.CashierTrashedRow, total *int)

	GetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*models.CashierByMerchantRow, *int, bool)
	SetCachedCashiersByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant, res []*models.CashierByMerchantRow, total *int)
}

type CashierCommandCache interface {
	DeleteCashierCache(ctx context.Context, id int)
	DeleteCashierListCache(ctx context.Context)
}

type CashierStatsCache interface {
	GetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales) ([]*models.CashierMonthlyTotalSalesRow, bool)
	SetMonthlyTotalSalesCache(ctx context.Context, req *requests.MonthTotalSales, res []*models.CashierMonthlyTotalSalesRow)

	GetYearlyTotalSalesCache(ctx context.Context, year int) ([]*models.CashierYearlyTotalSalesRow, bool)
	SetYearlyTotalSalesCache(ctx context.Context, year int, res []*models.CashierYearlyTotalSalesRow)

	GetMonthlySalesCache(ctx context.Context, year int) ([]*models.CashierMonthlySaleRow, bool)
	SetMonthlySalesCache(ctx context.Context, year int, res []*models.CashierMonthlySaleRow)

	GetYearlySalesCache(ctx context.Context, year int) ([]*models.CashierYearlySaleRow, bool)
	SetYearlySalesCache(ctx context.Context, year int, res []*models.CashierYearlySaleRow)
}

type CashierStatsByIdCache interface {
	GetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*models.CashierMonthlyTotalSalesRow, bool)
	SetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier, res []*models.CashierMonthlyTotalSalesRow)

	GetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*models.CashierYearlyTotalSalesRow, bool)
	SetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier, res []*models.CashierYearlyTotalSalesRow)

	GetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId) ([]*models.CashierMonthlySaleByIdRow, bool)
	SetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId, res []*models.CashierMonthlySaleByIdRow)

	GetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId) ([]*models.CashierYearlySaleRow, bool)
	SetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId, res []*models.CashierYearlySaleRow)
}

type CashierStatsByMerchantCache interface {
	GetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*models.CashierMonthlyTotalSalesRow, bool)
	SetMonthlyTotalSalesByMerchantCache(ctx context.Context, req *requests.MonthTotalSalesMerchant, res []*models.CashierMonthlyTotalSalesRow)

	GetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*models.CashierYearlyTotalSalesRow, bool)
	SetYearlyTotalSalesByMerchantCache(ctx context.Context, req *requests.YearTotalSalesMerchant, res []*models.CashierYearlyTotalSalesRow)

	GetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant) ([]*models.CashierMonthlySaleByMerchantRow, bool)
	SetMonthlyCashierByMerchantCache(ctx context.Context, req *requests.MonthCashierMerchant, res []*models.CashierMonthlySaleByMerchantRow)

	GetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant) ([]*models.CashierYearlySaleRow, bool)
	SetYearlyCashierByMerchantCache(ctx context.Context, req *requests.YearCashierMerchant, res []*models.CashierYearlySaleRow)
}
