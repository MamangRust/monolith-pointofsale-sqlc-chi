package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	sharedErrors "github.com/MamangRust/monolith-point-of-sale-shared/errors"
	"gorm.io/gorm"
)

type categoryStatsRepository struct {
	db *gorm.DB
}

func NewCategoryStatsRepository(db *gorm.DB) CategoryStatsRepository {
	return &categoryStatsRepository{
		db: db,
	}
}

func (r *categoryStatsRepository) GetMonthlyTotalPrice(ctx context.Context, req *requests.MonthTotalPrice) ([]*models.CategoryMonthlyTotalPriceRow, error) {
	currentMonthStart := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	currentMonthEnd := currentMonthStart.AddDate(0, 1, -1)
	prevMonthStart := currentMonthStart.AddDate(0, -1, 0)
	prevMonthEnd := prevMonthStart.AddDate(0, 1, -1)

	var results []*models.CategoryMonthlyTotalPriceRow

	err := r.db.WithContext(ctx).Raw(`
		WITH monthly_totals AS (
			SELECT
				EXTRACT(YEAR FROM o.created_at)::TEXT AS year,
				EXTRACT(MONTH FROM o.created_at)::integer AS month,
				COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
			FROM
				orders o
			JOIN
				order_items oi ON o.order_id = oi.order_id
			JOIN
				products p ON oi.product_id = p.product_id
			JOIN
				categories c ON p.category_id = c.category_id
			WHERE
				o.deleted_at IS NULL
				AND oi.deleted_at IS NULL
				AND (
					(o.created_at >= ? AND o.created_at <= ?)
					OR (o.created_at >= ? AND o.created_at <= ?)
				)
			GROUP BY
				EXTRACT(YEAR FROM o.created_at),
				EXTRACT(MONTH FROM o.created_at)
		),
		all_months AS (
			SELECT
				EXTRACT(YEAR FROM ?::timestamp)::TEXT AS year,
				EXTRACT(MONTH FROM ?::timestamp)::integer AS month,
				TO_CHAR(?::timestamp, 'FMMonth') AS month_name

			UNION

			SELECT
				EXTRACT(YEAR FROM ?::timestamp)::TEXT AS year,
				EXTRACT(MONTH FROM ?::timestamp)::integer AS month,
				TO_CHAR(?::timestamp, 'FMMonth') AS month_name
		)
		SELECT
			COALESCE(am.year, EXTRACT(YEAR FROM ?::timestamp)::TEXT) AS year,
			COALESCE(am.month_name, TO_CHAR(?::timestamp, 'FMMonth')) AS month,
			COALESCE(mt.total_revenue, 0) AS total_revenue
		FROM
			all_months am
		LEFT JOIN
			monthly_totals mt ON am.year = mt.year AND am.month = mt.month
		ORDER BY
			am.year::INT DESC,
			am.month DESC`,
		currentMonthStart, currentMonthEnd, prevMonthStart, prevMonthEnd,
		currentMonthStart, currentMonthStart, currentMonthStart,
		prevMonthStart, prevMonthStart, prevMonthStart,
		currentMonthStart, currentMonthStart,
	).Scan(&results).Error

	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return results, nil
}

func (r *categoryStatsRepository) GetYearlyTotalPrices(ctx context.Context, year int) ([]*models.CategoryYearlyTotalPriceRow, error) {
	var results []*models.CategoryYearlyTotalPriceRow

	err := r.db.WithContext(ctx).Raw(`
		WITH yearly_data AS (
			SELECT
				EXTRACT(YEAR FROM o.created_at)::integer AS year,
				COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
			FROM
				orders o
			JOIN
				order_items oi ON o.order_id = oi.order_id
			JOIN
				products p ON oi.product_id = p.product_id
			JOIN
				categories c ON p.category_id = c.category_id
			WHERE
				o.deleted_at IS NULL
				AND oi.deleted_at IS NULL
				AND p.deleted_at IS NULL
				AND c.deleted_at IS NULL
				AND (
					EXTRACT(YEAR FROM o.created_at) = ?::integer
					OR EXTRACT(YEAR FROM o.created_at) = ?::integer - 1
				)
			GROUP BY
				EXTRACT(YEAR FROM o.created_at)
		),
		all_years AS (
			SELECT ?::integer AS year
			UNION
			SELECT ?::integer - 1 AS year
		)
		SELECT
			a.year::text AS year,
			COALESCE(yd.total_revenue, 0) AS total_revenue
		FROM
			all_years a
		LEFT JOIN
			yearly_data yd ON a.year = yd.year
		ORDER BY
			a.year DESC`,
		year, year, year, year,
	).Scan(&results).Error

	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return results, nil
}

func (r *categoryStatsRepository) GetMonthPrice(ctx context.Context, year int) ([]*models.CategoryMonthlyRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	var results []*models.CategoryMonthlyRow

	err := r.db.WithContext(ctx).Raw(`
		WITH date_range AS (
			SELECT
				date_trunc('month', ?::timestamp) AS start_date,
				date_trunc('month', ?::timestamp) + interval '1 year' - interval '1 day' AS end_date
		),
		monthly_category_stats AS (
			SELECT
				c.category_id,
				c.name AS category_name,
				date_trunc('month', o.created_at) AS activity_month,
				COUNT(DISTINCT o.order_id) AS order_count,
				SUM(oi.quantity) AS items_sold,
				COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue
			FROM
				orders o
			JOIN
				order_items oi ON o.order_id = oi.order_id
			JOIN
				products p ON oi.product_id = p.product_id
			JOIN
				categories c ON p.category_id = c.category_id
			WHERE
				o.deleted_at IS NULL
				AND oi.deleted_at IS NULL
				AND p.deleted_at IS NULL
				AND c.deleted_at IS NULL
				AND o.created_at BETWEEN (SELECT start_date FROM date_range)
									 AND (SELECT end_date FROM date_range)
			GROUP BY
				c.category_id, c.name, activity_month
		)
		SELECT
			TO_CHAR(mcs.activity_month, 'Mon') AS month,
			mcs.category_id,
			mcs.category_name,
			mcs.order_count,
			mcs.items_sold,
			mcs.total_revenue
		FROM
			monthly_category_stats mcs
		ORDER BY
			mcs.activity_month, mcs.total_revenue DESC`,
		yearStart, yearStart,
	).Scan(&results).Error

	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return results, nil
}

func (r *categoryStatsRepository) GetYearPrice(ctx context.Context, year int) ([]*models.CategoryYearlyRow, error) {
	yearStart := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	var results []*models.CategoryYearlyRow

	err := r.db.WithContext(ctx).Raw(`
		WITH last_five_years AS (
			SELECT
				c.category_id,
				c.name AS category_name,
				EXTRACT(YEAR FROM o.created_at)::text AS year,
				COUNT(DISTINCT o.order_id) AS order_count,
				SUM(oi.quantity) AS items_sold,
				COALESCE(SUM(o.total_price), 0)::INTEGER AS total_revenue,
				COUNT(DISTINCT oi.product_id) AS unique_products_sold
			FROM
				orders o
			JOIN
				order_items oi ON o.order_id = oi.order_id
			JOIN
				products p ON oi.product_id = p.product_id
			JOIN
				categories c ON p.category_id = c.category_id
			WHERE
				o.deleted_at IS NULL
				AND oi.deleted_at IS NULL
				AND p.deleted_at IS NULL
				AND c.deleted_at IS NULL
				AND EXTRACT(YEAR FROM o.created_at) BETWEEN (EXTRACT(YEAR FROM ?::timestamp) - 4) AND EXTRACT(YEAR FROM ?::timestamp)
			GROUP BY
				c.category_id, c.name, EXTRACT(YEAR FROM o.created_at)
		)
		SELECT
			year,
			category_id,
			category_name,
			order_count,
			items_sold,
			total_revenue,
			unique_products_sold
		FROM
			last_five_years
		ORDER BY
			year, total_revenue DESC`,
		yearStart, yearStart,
	).Scan(&results).Error

	if err != nil {
		return nil, sharedErrors.ErrInternal.WithInternal(err)
	}

	return results, nil
}
