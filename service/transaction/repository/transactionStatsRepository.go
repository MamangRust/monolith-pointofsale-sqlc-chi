package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"gorm.io/gorm"
)

type transactonStatsRepository struct {
	db *gorm.DB
}

func NewTransactionStatsRepository(db *gorm.DB) TransactionStatsRepository {
	return &transactonStatsRepository{db: db}
}

func (r *transactonStatsRepository) GetMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*models.TransactionMonthlyAmountSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)

	var result []struct {
		Year         string `gorm:"column:year"`
		Month        string `gorm:"column:month"`
		TotalSuccess int64  `gorm:"column:success_count"`
		TotalAmount  int32  `gorm:"column:total_amount"`
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(YEAR FROM created_at)::text AS year,
			EXTRACT(MONTH FROM created_at)::text AS month,
			COALESCE(SUM(CASE WHEN payment_status = 'completed' THEN amount ELSE 0 END), 0)::float8 AS total_amount,
			COUNT(CASE WHEN payment_status = 'completed' THEN 1 END) AS success_count,
			COALESCE(SUM(CASE WHEN payment_status = 'completed' AND created_at BETWEEN ? AND ? THEN amount ELSE 0 END), 0)::float8 AS prev_amount,
			COUNT(CASE WHEN payment_status = 'completed' AND created_at BETWEEN ? AND ? THEN 1 END) AS prev_count
		FROM transactions
		WHERE created_at BETWEEN ? AND ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at), EXTRACT(MONTH FROM created_at)
	`, lastDayCurrentMonth.AddDate(0, -1, 0), lastDayCurrentMonth, lastDayCurrentMonth.AddDate(0, -1, 0), lastDayCurrentMonth, currentDate, lastDayCurrentMonth).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccess
	}

	var res []*models.TransactionMonthlyAmountSuccessRow
	for _, r := range result {
		res = append(res, &models.TransactionMonthlyAmountSuccessRow{
			Year:         r.Year,
			Month:        r.Month,
			TotalAmount:  r.TotalAmount,
			TotalSuccess: r.TotalSuccess,
		})
	}
	return res, nil
}

func (r *transactonStatsRepository) GetYearlyAmountSuccess(ctx context.Context, year int) ([]*models.TransactionYearlyAmountSuccessRow, error) {
	var result struct {
		Year         string `gorm:"column:year"`
		TotalAmount  int32  `gorm:"column:total_amount"`
		TotalSuccess int64  `gorm:"column:success_count"`
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(YEAR FROM created_at)::text AS year,
			COALESCE(SUM(CASE WHEN payment_status = 'completed' THEN amount ELSE 0 END), 0) AS total_amount,
			COUNT(CASE WHEN payment_status = 'completed' THEN 1 END) AS success_count
		FROM transactions
		WHERE EXTRACT(YEAR FROM created_at) = ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at)
	`, year).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccess
	}

	return []*models.TransactionYearlyAmountSuccessRow{
		{Year: strconv.Itoa(year), TotalAmount: result.TotalAmount, TotalSuccess: int32(result.TotalSuccess)},
	}, nil
}

func (r *transactonStatsRepository) GetMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*models.TransactionMonthlyAmountFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)

	var result struct {
		Year        string `gorm:"column:year"`
		Month       string `gorm:"column:month"`
		TotalFailed int64  `gorm:"column:failed_count"`
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(YEAR FROM created_at)::text AS year,
			EXTRACT(MONTH FROM created_at)::text AS month,
			COUNT(CASE WHEN payment_status = 'failed' THEN 1 END) AS failed_count
		FROM transactions
		WHERE created_at BETWEEN ? AND ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at), EXTRACT(MONTH FROM created_at)
	`, currentDate, lastDayCurrentMonth).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailed
	}

	return []*models.TransactionMonthlyAmountFailedRow{
		{Year: result.Year, Month: result.Month, TotalFailed: int64(result.TotalFailed)},
	}, nil
}

func (r *transactonStatsRepository) GetYearlyAmountFailed(ctx context.Context, year int) ([]*models.TransactionYearlyAmountFailedRow, error) {
	var result struct {
		Year        string `gorm:"column:year"`
		TotalFailed int64  `gorm:"column:failed_count"`
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(YEAR FROM created_at)::text AS year,
			COUNT(CASE WHEN payment_status = 'failed' THEN 1 END) AS failed_count
		FROM transactions
		WHERE EXTRACT(YEAR FROM created_at) = ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at)
	`, year).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailed
	}

	return []*models.TransactionYearlyAmountFailedRow{
		{Year: strconv.Itoa(year), TotalFailed: int32(result.TotalFailed)},
	}, nil
}

func (r *transactonStatsRepository) GetMonthlyTransactionMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*models.TransactionMonthlyMethodSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)

	var results []models.TransactionMonthlyMethodSuccessRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(MONTH FROM created_at)::text AS month,
			payment_method,
			COUNT(*) AS total_transactions,
			COALESCE(SUM(amount), 0) AS total_amount
		FROM transactions
		WHERE created_at BETWEEN ? AND ? AND payment_status = 'completed' AND deleted_at IS NULL
		GROUP BY EXTRACT(MONTH FROM created_at), payment_method
		ORDER BY payment_method
	`, currentDate, lastDayCurrentMonth).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod
	}

	var res []*models.TransactionMonthlyMethodSuccessRow
	for _, r := range results {
		res = append(res, &models.TransactionMonthlyMethodSuccessRow{
			Month:             r.Month,
			PaymentMethod:     r.PaymentMethod,
			TotalTransactions: r.TotalTransactions,
			TotalAmount:       r.TotalAmount,
		})
	}
	return res, nil
}

func (r *transactonStatsRepository) GetYearlyTransactionMethodSuccess(ctx context.Context, year int) ([]*models.TransactionYearlyMethodSuccessRow, error) {
	var results []models.TransactionYearlyMethodSuccessRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(YEAR FROM created_at)::text AS year,
			payment_method,
			COUNT(*) AS total_transactions,
			COALESCE(SUM(amount), 0) AS total_amount
		FROM transactions
		WHERE EXTRACT(YEAR FROM created_at) = ? AND payment_status = 'completed' AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at), payment_method
		ORDER BY payment_method
	`, year).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod
	}

	var res []*models.TransactionYearlyMethodSuccessRow
	for _, r := range results {
		res = append(res, &models.TransactionYearlyMethodSuccessRow{
			Year:              r.Year,
			PaymentMethod:     r.PaymentMethod,
			TotalTransactions: r.TotalTransactions,
			TotalAmount:       r.TotalAmount,
		})
	}
	return res, nil
}

func (r *transactonStatsRepository) GetMonthlyTransactionMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*models.TransactionMonthlyMethodFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)

	var results []models.TransactionMonthlyMethodFailedRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(MONTH FROM created_at)::text AS month,
			payment_method,
			COUNT(*) AS total_transactions,
			COALESCE(SUM(amount), 0) AS total_amount
		FROM transactions
		WHERE created_at BETWEEN ? AND ? AND payment_status = 'failed' AND deleted_at IS NULL
		GROUP BY EXTRACT(MONTH FROM created_at), payment_method
		ORDER BY payment_method
	`, currentDate, lastDayCurrentMonth).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethod
	}

	var res []*models.TransactionMonthlyMethodFailedRow
	for _, r := range results {
		res = append(res, &models.TransactionMonthlyMethodFailedRow{
			Month:             r.Month,
			PaymentMethod:     r.PaymentMethod,
			TotalTransactions: r.TotalTransactions,
			TotalAmount:       r.TotalAmount,
		})
	}
	return res, nil
}

func (r *transactonStatsRepository) GetYearlyTransactionMethodFailed(ctx context.Context, year int) ([]*models.TransactionYearlyMethodFailedRow, error) {
	var results []models.TransactionYearlyMethodFailedRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(YEAR FROM created_at)::text AS year,
			payment_method,
			COUNT(*) AS total_transactions,
			COALESCE(SUM(amount), 0) AS total_amount
		FROM transactions
		WHERE EXTRACT(YEAR FROM created_at) = ? AND payment_status = 'failed' AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at), payment_method
		ORDER BY payment_method
	`, year).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethod
	}

	var res []*models.TransactionYearlyMethodFailedRow
	for _, r := range results {
		res = append(res, &models.TransactionYearlyMethodFailedRow{
			Year:              r.Year,
			PaymentMethod:     r.PaymentMethod,
			TotalTransactions: r.TotalTransactions,
			TotalAmount:       r.TotalAmount,
		})
	}
	return res, nil
}
