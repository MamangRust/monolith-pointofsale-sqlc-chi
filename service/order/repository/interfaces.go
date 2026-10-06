package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type OrderStatsRepository interface {
	GetMonthlyTotalRevenue(ctx context.Context, req *requests.MonthTotalRevenue) ([]*models.MonthlyTotalRevenueRow, error)
	GetYearlyTotalRevenue(ctx context.Context, year int) ([]*models.YearlyTotalRevenueRow, error)
	GetMonthlyOrder(ctx context.Context, year int) ([]*models.MonthlyOrderRow, error)
	GetYearlyOrder(ctx context.Context, year int) ([]*models.YearlyOrderRow, error)
}

type OrderStatByMerchantRepository interface {
	GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*models.MonthlyTotalRevenueByMerchantRow, error)
	GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*models.YearlyTotalRevenueByMerchantRow, error)
	GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*models.MonthlyOrderByMerchantRow, error)
	GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*models.YearlyOrderByMerchantRow, error)
}

type OrderQueryRepository interface {
	FindAllOrders(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*models.OrderByMerchantRow, *int, error)
	FindById(ctx context.Context, orderID int) (*models.Order, error)
	FindByTrashedId(ctx context.Context, orderID int) (*models.Order, error)
}

type OrderCommandRepository interface {
	DeleteOrder(ctx context.Context, orderID int) error
	CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*models.Order, error)
	UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*models.Order, error)
	FindAllTrashed(ctx context.Context) ([]*models.Order, error)
	TrashedOrder(ctx context.Context, orderID int) (*models.Order, error)
	RestoreOrder(ctx context.Context, orderID int) (*models.Order, error)
	DeleteOrderPermanent(ctx context.Context, orderID int) (bool, error)
	DeleteAllOrderPermanent(ctx context.Context) (bool, error)
}

