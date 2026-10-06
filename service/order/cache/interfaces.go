package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type OrderStatsCache interface {
	GetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue) ([]*models.MonthlyTotalRevenueRow, bool)
	SetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue, res []*models.MonthlyTotalRevenueRow)

	GetYearlyTotalRevenueCache(ctx context.Context, year int) ([]*models.YearlyTotalRevenueRow, bool)
	SetYearlyTotalRevenueCache(ctx context.Context, year int, res []*models.YearlyTotalRevenueRow)

	GetMonthlyOrderCache(ctx context.Context, year int) ([]*models.MonthlyOrderRow, bool)
	SetMonthlyOrderCache(ctx context.Context, year int, res []*models.MonthlyOrderRow)

	GetYearlyOrderCache(ctx context.Context, year int) ([]*models.YearlyOrderRow, bool)
	SetYearlyOrderCache(ctx context.Context, year int, res []*models.YearlyOrderRow)
}

type OrderStatsByMerchantCache interface {
	GetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*models.MonthlyTotalRevenueByMerchantRow, bool)
	SetMonthlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.MonthTotalRevenueMerchant, res []*models.MonthlyTotalRevenueByMerchantRow)

	GetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*models.YearlyTotalRevenueByMerchantRow, bool)
	SetYearlyTotalRevenueByMerchantCache(ctx context.Context, req *requests.YearTotalRevenueMerchant, res []*models.YearlyTotalRevenueByMerchantRow)

	GetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant) ([]*models.MonthlyOrderByMerchantRow, bool)
	SetMonthlyOrderByMerchantCache(ctx context.Context, req *requests.MonthOrderMerchant, res []*models.MonthlyOrderByMerchantRow)

	GetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant) ([]*models.YearlyOrderByMerchantRow, bool)
	SetYearlyOrderByMerchantCache(ctx context.Context, req *requests.YearOrderMerchant, res []*models.YearlyOrderByMerchantRow)
}

type OrderQueryCache interface {
	GetOrderAllCache(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderRow, *int, bool)
	SetOrderAllCache(ctx context.Context, req *requests.FindAllOrders, data []*models.OrderRow, total *int)

	GetCachedOrderCache(ctx context.Context, orderID int) (*models.Order, bool)
	SetCachedOrderCache(ctx context.Context, data *models.Order)

	GetCachedOrderMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*models.OrderByMerchantRow, *int, bool)
	SetCachedOrderMerchant(ctx context.Context, req *requests.FindAllOrderMerchant, res []*models.OrderByMerchantRow, total *int)

	GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderActiveRow, *int, bool)
	SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrders, data []*models.OrderActiveRow, total *int)

	GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderTrashedRow, *int, bool)
	SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrders, data []*models.OrderTrashedRow, total *int)
}

type OrderCommandCache interface {
	DeleteOrderCache(ctx context.Context, id int)
	DeleteOrderAllCache(ctx context.Context)
}
