package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type OrderStatsService interface {
	FindMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*models.MonthlyTotalRevenueRow, error)
	FindYearlyTotalRevenue(ctx context.Context, year int) ([]*models.YearlyTotalRevenueRow, error)

	FindMonthlyOrder(ctx context.Context, year int) ([]*models.MonthlyOrderRow, error)
	FindYearlyOrder(ctx context.Context, year int) ([]*models.YearlyOrderRow, error)
}

type OrderStatByMerchantQueryService interface {
	FindMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*models.MonthlyTotalRevenueByMerchantRow, error)
	FindYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*models.YearlyTotalRevenueByMerchantRow, error)

	FindMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*models.MonthlyOrderByMerchantRow, error)
	FindYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*models.YearlyOrderByMerchantRow, error)
}

type OrderQueryService interface {
	FindAll(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderRow, *int, error)
	FindById(ctx context.Context, orderID int) (*models.Order, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*models.OrderByMerchantRow, *int, error)
}

type OrderCommandService interface {
	CreateOrder(ctx context.Context, req *requests.CreateOrderRequest) (*models.Order, error)
	UpdateOrder(ctx context.Context, req *requests.UpdateOrderRequest) (*models.Order, error)
	TrashedOrder(ctx context.Context, orderID int) (*models.Order, error)
	RestoreOrder(ctx context.Context, orderID int) (*models.Order, error)
	DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error)
	RestoreAllOrder(ctx context.Context) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}
