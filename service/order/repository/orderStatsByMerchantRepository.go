package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/order_errors"
	"gorm.io/gorm"
)

type orderStatsByMerchantRepository struct {
	db *gorm.DB
}

func NewOrderStatsByMerchantRepository(db *gorm.DB) OrderStatByMerchantRepository {
	return &orderStatsByMerchantRepository{db: db}
}

func (r *orderStatsByMerchantRepository) GetMonthlyTotalRevenueByMerchant(ctx context.Context, req *requests.MonthTotalRevenueMerchant) ([]*models.MonthlyTotalRevenueByMerchantRow, error) {
	monthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	monthEnd := monthStart.AddDate(0, 1, 0)

	var result struct {
		Year         string  `gorm:"column:year"`
		Month        string  `gorm:"column:month"`
		TotalRevenue float64 `gorm:"column:total_revenue"`
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			EXTRACT(YEAR FROM created_at)::text AS year,
			EXTRACT(MONTH FROM created_at)::text AS month,
			COALESCE(SUM(total_price), 0)::float8 AS total_revenue
		FROM orders
		WHERE merchant_id = ? AND created_at >= ? AND created_at < ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at), EXTRACT(MONTH FROM created_at)
	`, req.MerchantID, monthStart, monthEnd).Scan(&result).Error
	if err != nil {
		return nil, order_errors.ErrGetMonthlyTotalRevenueByMerchant
	}

	return []*models.MonthlyTotalRevenueByMerchantRow{
		{Year: result.Year, Month: result.Month, TotalRevenue: int32(result.TotalRevenue)},
	}, nil
}

func (r *orderStatsByMerchantRepository) GetYearlyTotalRevenueByMerchant(ctx context.Context, req *requests.YearTotalRevenueMerchant) ([]*models.YearlyTotalRevenueByMerchantRow, error) {
	var result struct {
		Year         string  `gorm:"column:year"`
		TotalRevenue float64 `gorm:"column:total_revenue"`
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			EXTRACT(YEAR FROM created_at)::text AS year,
			COALESCE(SUM(total_price), 0)::float8 AS total_revenue
		FROM orders
		WHERE merchant_id = ? AND EXTRACT(YEAR FROM created_at) = ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at)
	`, req.MerchantID, req.Year).Scan(&result).Error
	if err != nil {
		return nil, order_errors.ErrGetYearlyTotalRevenueByMerchant
	}

	return []*models.YearlyTotalRevenueByMerchantRow{
		{Year: fmt.Sprintf("%d", req.Year), TotalRevenue: int32(result.TotalRevenue)},
	}, nil
}

func (r *orderStatsByMerchantRepository) GetMonthlyOrderByMerchant(ctx context.Context, req *requests.MonthOrderMerchant) ([]*models.MonthlyOrderByMerchantRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)
	yearEnd := yearStart.AddDate(1, 0, 0)

	var result []struct {
		Month          string  `gorm:"column:month"`
		OrderCount     int64   `gorm:"column:order_count"`
		TotalRevenue   float64 `gorm:"column:total_revenue"`
		TotalItemsSold int64   `gorm:"column:total_items_sold"`
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			TO_CHAR(o.created_at, 'YYYY-MM') AS month,
			COUNT(DISTINCT o.order_id) AS order_count,
			COALESCE(SUM(o.total_price), 0)::float8 AS total_revenue,
			COALESCE(SUM(oi.quantity), 0) AS total_items_sold
		FROM orders o
		LEFT JOIN order_items oi ON o.order_id = oi.order_id AND oi.deleted_at IS NULL
		WHERE o.merchant_id = ? AND o.created_at >= ? AND o.created_at < ? AND o.deleted_at IS NULL
		GROUP BY TO_CHAR(o.created_at, 'YYYY-MM')
		ORDER BY month
	`, req.MerchantID, yearStart, yearEnd).Scan(&result).Error
	if err != nil {
		return nil, order_errors.ErrGetMonthlyOrderByMerchant
	}

	var res []*models.MonthlyOrderByMerchantRow
	for _, r := range result {
		res = append(res, &models.MonthlyOrderByMerchantRow{
			Month:          r.Month,
			OrderCount:     r.OrderCount,
			TotalRevenue:   r.TotalRevenue,
			TotalItemsSold: r.TotalItemsSold,
		})
	}
	return res, nil
}

func (r *orderStatsByMerchantRepository) GetYearlyOrderByMerchant(ctx context.Context, req *requests.YearOrderMerchant) ([]*models.YearlyOrderByMerchantRow, error) {
	var result struct {
		Year               string  `gorm:"column:year"`
		OrderCount         int64   `gorm:"column:order_count"`
		TotalRevenue       float64 `gorm:"column:total_revenue"`
		TotalItemsSold     int64   `gorm:"column:total_items_sold"`
		ActiveCashiers     int64   `gorm:"column:active_cashiers"`
		UniqueProductsSold int64   `gorm:"column:unique_products_sold"`
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			EXTRACT(YEAR FROM o.created_at)::text AS year,
			COUNT(DISTINCT o.order_id) AS order_count,
			COALESCE(SUM(o.total_price), 0)::float8 AS total_revenue,
			COALESCE(SUM(oi.quantity), 0) AS total_items_sold,
			COUNT(DISTINCT o.cashier_id) AS active_cashiers,
			COUNT(DISTINCT oi.product_id) AS unique_products_sold
		FROM orders o
		LEFT JOIN order_items oi ON o.order_id = oi.order_id AND oi.deleted_at IS NULL
		WHERE o.merchant_id = ? AND EXTRACT(YEAR FROM o.created_at) = ? AND o.deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM o.created_at)
	`, req.MerchantID, req.Year).Scan(&result).Error
	if err != nil {
		return nil, order_errors.ErrGetYearlyOrderByMerchant
	}

	return []*models.YearlyOrderByMerchantRow{
		{
			Year:               fmt.Sprintf("%d", req.Year),
			OrderCount:         result.OrderCount,
			TotalRevenue:       result.TotalRevenue,
			TotalItemsSold:     result.TotalItemsSold,
			ActiveCashiers:     result.ActiveCashiers,
			UniqueProductsSold: result.UniqueProductsSold,
		},
	}, nil
}
