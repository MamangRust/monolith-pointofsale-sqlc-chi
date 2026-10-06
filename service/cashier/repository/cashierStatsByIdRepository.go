package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/cashier_errors"
	"gorm.io/gorm"
)

type cashierStatsByIdRepository struct {
	db *gorm.DB
}

func NewCashierStatsByIdRepository(db *gorm.DB) CashierStatByIdRepository {
	return &cashierStatsByIdRepository{
		db: db,
	}
}

const monthlyTotalSalesByIdQuery = `
WITH monthly_totals AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::TEXT AS year,
        EXTRACT(MONTH FROM o.created_at)::integer AS month,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND (
            (o.created_at >= @current_month_start AND o.created_at <= @current_month_end)
            OR (o.created_at >= @prev_month_start AND o.created_at <= @prev_month_end)
        )
        AND c.cashier_id = @cashier_id
    GROUP BY
        EXTRACT(YEAR FROM o.created_at),
        EXTRACT(MONTH FROM o.created_at)
),
all_months AS (
    SELECT
        EXTRACT(YEAR FROM @current_month_start)::TEXT AS year,
        EXTRACT(MONTH FROM @current_month_start)::integer AS month,
        TO_CHAR(@current_month_start, 'FMMonth') AS month_name

    UNION

    SELECT
        EXTRACT(YEAR FROM @prev_month_start)::TEXT AS year,
        EXTRACT(MONTH FROM @prev_month_start)::integer AS month,
        TO_CHAR(@prev_month_start, 'FMMonth') AS month_name
)
SELECT
    COALESCE(am.year, EXTRACT(YEAR FROM @current_month_start)::TEXT) AS year,
    COALESCE(am.month_name, TO_CHAR(@current_month_start, 'FMMonth')) AS month,
    COALESCE(mt.total_sales, 0) AS total_sales
FROM
    all_months am
LEFT JOIN
    monthly_totals mt ON am.year = mt.year AND am.month = mt.month
ORDER BY
    am.year::INT DESC,
    am.month DESC`

const yearlyTotalSalesByIdQuery = `
WITH yearly_data AS (
    SELECT
        EXTRACT(YEAR FROM o.created_at)::integer AS year,
        COALESCE(SUM(o.total_price), 0)::INTEGER AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND (
            EXTRACT(YEAR FROM o.created_at) = @year::integer
            OR EXTRACT(YEAR FROM o.created_at) = @year::integer - 1
        )
        AND c.cashier_id = @cashier_id
    GROUP BY
        EXTRACT(YEAR FROM o.created_at)
),
all_years AS (
    SELECT @year::integer AS year
    UNION
    SELECT @year::integer - 1 AS year
)
SELECT
    a.year::text AS year,
    COALESCE(yd.total_sales, 0) AS total_sales
FROM
    all_years a
LEFT JOIN
    yearly_data yd ON a.year = yd.year
ORDER BY
    a.year DESC`

const monthlyCashierByIdQuery = `
WITH date_range AS (
    SELECT
        date_trunc('month', @ref_date::timestamp) AS start_date,
        date_trunc('month', @ref_date::timestamp) + interval '1 year' - interval '1 day' AS end_date
),
cashier_activity AS (
    SELECT
        c.cashier_id,
        c.name AS cashier_name,
        date_trunc('month', o.created_at) AS activity_month,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price) AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND c.cashier_id = @cashier_id
        AND o.created_at BETWEEN (SELECT start_date FROM date_range)
                             AND (SELECT end_date FROM date_range)
    GROUP BY
        c.cashier_id, c.name, activity_month
)
SELECT
    ca.cashier_id,
    ca.cashier_name,
    TO_CHAR(ca.activity_month, 'Mon') AS month,
    ca.order_count,
    ca.total_sales
FROM
    cashier_activity ca
ORDER BY
    ca.activity_month, ca.cashier_id`

const yearlyCashierByIdQuery = `
WITH last_five_years AS (
    SELECT
        c.cashier_id,
        c.name AS cashier_name,
        EXTRACT(YEAR FROM o.created_at)::text AS year,
        COUNT(o.order_id) AS order_count,
        SUM(o.total_price) AS total_sales
    FROM
        orders o
    JOIN
        cashiers c ON o.cashier_id = c.cashier_id
    WHERE
        o.deleted_at IS NULL
        AND c.deleted_at IS NULL
        AND c.cashier_id = @cashier_id
        AND EXTRACT(YEAR FROM o.created_at) BETWEEN (EXTRACT(YEAR FROM @ref_date::timestamp) - 4) AND EXTRACT(YEAR FROM @ref_date::timestamp)
    GROUP BY
        c.cashier_id, c.name, EXTRACT(YEAR FROM o.created_at)
)
SELECT
    year,
    cashier_id,
    cashier_name,
    order_count,
    total_sales
FROM
    last_five_years
ORDER BY
    year, cashier_id`

func (r *cashierStatsByIdRepository) GetMonthlyTotalSalesById(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*models.CashierMonthlyTotalSalesRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)

	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	var results []*models.CashierMonthlyTotalSalesRow
	if err := r.db.WithContext(ctx).Raw(monthlyTotalSalesByIdQuery, sql.Named("current_month_start", currentMonthStart),
		sql.Named("current_month_end", currentMonthEnd),
		sql.Named("prev_month_start", prevMonthStart),
		sql.Named("prev_month_end", prevMonthEnd),
		sql.Named("cashier_id", int32(req.CashierID)),
	).Scan(&results).Error; err != nil {
		return nil, cashier_errors.ErrGetMonthlyTotalSalesById.WithInternal(err)
	}

	return results, nil
}

func (r *cashierStatsByIdRepository) GetYearlyTotalSalesById(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*models.CashierYearlyTotalSalesRow, error) {
	var results []*models.CashierYearlyTotalSalesRow
	if err := r.db.WithContext(ctx).Raw(yearlyTotalSalesByIdQuery, sql.Named("year", int32(req.Year)),
		sql.Named("cashier_id", int32(req.CashierID)),
	).Scan(&results).Error; err != nil {
		return nil, cashier_errors.ErrGetYearlyTotalSalesById.WithInternal(err)
	}

	return results, nil
}

func (r *cashierStatsByIdRepository) GetMonthlyCashierById(ctx context.Context, req *requests.MonthCashierId) ([]*models.CashierMonthlySaleByIdRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	var results []*models.CashierMonthlySaleByIdRow
	if err := r.db.WithContext(ctx).Raw(monthlyCashierByIdQuery, sql.Named("ref_date", yearStart),
		sql.Named("cashier_id", int32(req.CashierID)),
	).Scan(&results).Error; err != nil {
		return nil, cashier_errors.ErrGetMonthlyCashierById.WithInternal(err)
	}

	return results, nil
}

func (r *cashierStatsByIdRepository) GetYearlyCashierById(ctx context.Context, req *requests.YearCashierId) ([]*models.CashierYearlySaleRow, error) {
	yearStart := time.Date(req.Year, 1, 1, 0, 0, 0, 0, time.UTC)

	var results []*models.CashierYearlySaleRow
	if err := r.db.WithContext(ctx).Raw(yearlyCashierByIdQuery, sql.Named("ref_date", yearStart),
		sql.Named("cashier_id", int32(req.CashierID)),
	).Scan(&results).Error; err != nil {
		return nil, cashier_errors.ErrGetYearlyCashierById.WithInternal(err)
	}

	return results, nil
}
