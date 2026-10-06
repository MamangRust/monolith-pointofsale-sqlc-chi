package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type CashierStatsRepository interface {
	GetMonthlyTotalSales(ctx context.Context, req *requests.MonthTotalSales) ([]*models.CashierMonthlyTotalSalesRow, error)
	GetYearlyTotalSales(ctx context.Context, year int) ([]*models.CashierYearlyTotalSalesRow, error)

	GetMonthyCashier(ctx context.Context, year int) ([]*models.CashierMonthlySaleRow, error)
	GetYearlyCashier(ctx context.Context, year int) ([]*models.CashierYearlySaleRow, error)
}

type CashierStatByIdRepository interface {
	GetMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*models.CashierMonthlyTotalSalesRow, error)
	GetYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*models.CashierYearlyTotalSalesRow, error)

	GetMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*models.CashierMonthlySaleByIdRow, error)
	GetYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*models.CashierYearlySaleRow, error)
}

type CashierStatByMerchantRepository interface {
	GetMonthlyTotalSalesByMerchant(ctx context.Context, req *requests.MonthTotalSalesMerchant) ([]*models.CashierMonthlyTotalSalesRow, error)
	GetYearlyTotalSalesByMerchant(ctx context.Context, req *requests.YearTotalSalesMerchant) ([]*models.CashierYearlyTotalSalesRow, error)

	GetMonthlyCashierByMerchant(ctx context.Context, req *requests.MonthCashierMerchant) ([]*models.CashierMonthlySaleByMerchantRow, error)
	GetYearlyCashierByMerchant(ctx context.Context, req *requests.YearCashierMerchant) ([]*models.CashierYearlySaleRow, error)
}

type CashierQueryRepository interface {
	FindAllCashiers(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierRow, *int, error)
	FindById(ctx context.Context, cashier_id int) (*models.Cashier, error)
	FindByActive(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllCashiers) ([]*models.CashierTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllCashierMerchant) ([]*models.CashierByMerchantRow, *int, error)
}

type CashierCommandRepository interface {
	CreateCashier(ctx context.Context, request *requests.CreateCashierRequest) (*models.Cashier, error)
	UpdateCashier(ctx context.Context, request *requests.UpdateCashierRequest) (*models.Cashier, error)
	TrashedCashier(ctx context.Context, cashier_id int) (*models.Cashier, error)
	RestoreCashier(ctx context.Context, cashier_id int) (*models.Cashier, error)
	DeleteCashierPermanent(ctx context.Context, cashier_id int) (bool, error)
	RestoreAllCashier(ctx context.Context) (bool, error)
	DeleteAllCashierPermanent(ctx context.Context) (bool, error)
}
