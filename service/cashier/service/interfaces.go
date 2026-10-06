package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CashierStatsService interface {
	FindMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*models.CashierMonthlyTotalSalesRow, error)
	FindYearlyTotalSales(ctx context.Context, year int) ([]*models.CashierYearlyTotalSalesRow, error)

	FindMonthlySales(ctx context.Context, year int) ([]*models.CashierMonthlySaleRow, error)
	FindYearlySales(ctx context.Context, year int) ([]*models.CashierYearlySaleRow, error)
}

type CashierStatsByIdService interface {
	FindMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*models.CashierMonthlyTotalSalesRow, error)
	FindYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*models.CashierYearlyTotalSalesRow, error)
	FindMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*models.CashierMonthlySaleByIdRow, error)
	FindYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*models.CashierYearlySaleRow, error)
}

type CashierStatsByMerchant interface {
	FindMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*models.CashierMonthlyTotalSalesRow, error)
	FindYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*models.CashierYearlyTotalSalesRow, error)

	FindMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*models.CashierMonthlySaleByMerchantRow, error)
	FindYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*models.CashierYearlySaleRow, error)
}

type CashierQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierRow, *int, error)
	FindById(ctx context.Context, cashierID int) (*models.Cashier, error)
	FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*models.CashierByMerchantRow, *int, error)
}

type CashierCommandService interface {
	CreateCashier(ctx context.Context, req *requests.CreateCashierRequest) (*models.Cashier, error)
	UpdateCashier(ctx context.Context, req *requests.UpdateCashierRequest) (*models.Cashier, error)
	TrashedCashier(ctx context.Context, cashierID int) (*models.Cashier, error)
	RestoreCashier(ctx context.Context, cashierID int) (*models.Cashier, error)
	DeleteCashierPermanent(ctx context.Context, cashierID int) (bool, error)
	RestoreAllCashier(ctx context.Context) (bool, error)
	DeleteAllCashierPermanent(ctx context.Context) (bool, error)
}
