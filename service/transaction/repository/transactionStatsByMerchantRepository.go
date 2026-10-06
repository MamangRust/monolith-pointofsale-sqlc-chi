package repository

import (
	"context"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"github.com/MamangRust/monolith-point-of-sale-shared/errors/transaction_errors"
	"gorm.io/gorm"
)

type transactionStatsByMerchantRepository struct {
	db *gorm.DB
}

func NewTransactionStatsByMerchantRepository(db *gorm.DB) TransactionStatsByMerchantRepository {
	return &transactionStatsByMerchantRepository{db: db}
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountSuccessByMerchantRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)

	var result struct {
		Year         string `gorm:"column:year"`
		Month        string `gorm:"column:month"`
		TotalAmount  int32  `gorm:"column:total_amount"`
		TotalSuccess int64  `gorm:"column:success_count"`
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(YEAR FROM created_at)::text AS year,
			EXTRACT(MONTH FROM created_at)::text AS month,
			COALESCE(SUM(CASE WHEN payment_status = 'completed' THEN amount ELSE 0 END), 0) AS total_amount,
			COUNT(CASE WHEN payment_status = 'completed' THEN 1 END) AS success_count
		FROM transactions
		WHERE merchant_id = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at), EXTRACT(MONTH FROM created_at)
	`, req.MerchantID, currentDate, lastDayCurrentMonth).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountSuccessByMerchant
	}

	return []*models.TransactionMonthlyAmountSuccessByMerchantRow{
		{Year: result.Year, Month: result.Month, TotalAmount: result.TotalAmount, TotalSuccess: result.TotalSuccess},
	}, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountSuccessByMerchantRow, error) {
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
		WHERE merchant_id = ? AND EXTRACT(YEAR FROM created_at) = ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at)
	`, req.MerchantID, req.Year).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountSuccessByMerchant
	}

	return []*models.TransactionYearlyAmountSuccessByMerchantRow{
		{Year: result.Year, TotalAmount: result.TotalAmount, TotalSuccess: int32(result.TotalSuccess)},
	}, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountFailedByMerchantRow, error) {
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
		WHERE merchant_id = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at), EXTRACT(MONTH FROM created_at)
	`, req.MerchantID, currentDate, lastDayCurrentMonth).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyAmountFailedByMerchant
	}

	return []*models.TransactionMonthlyAmountFailedByMerchantRow{
		{Year: result.Year, Month: result.Month, TotalFailed: result.TotalFailed},
	}, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountFailedByMerchantRow, error) {
	var result struct {
		Year        string `gorm:"column:year"`
		TotalFailed int64  `gorm:"column:failed_count"`
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(YEAR FROM created_at)::text AS year,
			COUNT(CASE WHEN payment_status = 'failed' THEN 1 END) AS failed_count
		FROM transactions
		WHERE merchant_id = ? AND EXTRACT(YEAR FROM created_at) = ? AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at)
	`, req.MerchantID, req.Year).Scan(&result).Error
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyAmountFailedByMerchant
	}

	return []*models.TransactionYearlyAmountFailedByMerchantRow{
		{Year: result.Year, TotalFailed: int32(result.TotalFailed)},
	}, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantSuccessRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)

	var results []models.TransactionMonthlyMethodByMerchantSuccessRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(MONTH FROM created_at)::text AS month,
			payment_method,
			COUNT(*) AS total_transactions,
			COALESCE(SUM(amount), 0) AS total_amount
		FROM transactions
		WHERE merchant_id = ? AND created_at BETWEEN ? AND ? AND payment_status = 'completed' AND deleted_at IS NULL
		GROUP BY EXTRACT(MONTH FROM created_at), payment_method
		ORDER BY payment_method
	`, req.MerchantID, currentDate, lastDayCurrentMonth).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant
	}
	return []*models.TransactionMonthlyMethodByMerchantSuccessRow{
		&results[0],
	}, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantSuccessRow, error) {

	var results []models.TransactionYearlyMethodByMerchantSuccessRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(YEAR FROM created_at)::text AS year,
			payment_method,
			COUNT(*) AS total_transactions,
			COALESCE(SUM(amount), 0) AS total_amount
		FROM transactions
		WHERE merchant_id = ? AND EXTRACT(YEAR FROM created_at) = ? AND payment_status = 'completed' AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at), payment_method
		ORDER BY payment_method
	`, req.MerchantID, req.Year).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant
	}
	var res []*models.TransactionYearlyMethodByMerchantSuccessRow
	for _, r := range results {
		res = append(res, &models.TransactionYearlyMethodByMerchantSuccessRow{
			Year:              r.Year,
			PaymentMethod:     r.PaymentMethod,
			TotalTransactions: r.TotalTransactions,
			TotalAmount:       r.TotalAmount,
		})
	}
	return res, nil
}

func (r *transactionStatsByMerchantRepository) GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantFailedRow, error) {
	currentDate := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	lastDayCurrentMonth := currentDate.AddDate(0, 1, -1)

	var results []models.TransactionMonthlyMethodByMerchantFailedRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(MONTH FROM created_at)::text AS month,
			payment_method,
			COUNT(*) AS total_transactions,
			COALESCE(SUM(amount), 0) AS total_amount
		FROM transactions
		WHERE merchant_id = ? AND created_at BETWEEN ? AND ? AND payment_status = 'failed' AND deleted_at IS NULL
		GROUP BY EXTRACT(MONTH FROM created_at), payment_method
		ORDER BY payment_method
	`, req.MerchantID, currentDate, lastDayCurrentMonth).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrGetMonthlyTransactionMethodByMerchant
	}
	var res []*models.TransactionMonthlyMethodByMerchantFailedRow
	for _, r := range results {
		res = append(res, &models.TransactionMonthlyMethodByMerchantFailedRow{
			Month:             r.Month,
			PaymentMethod:     r.PaymentMethod,
			TotalTransactions: r.TotalTransactions,
			TotalAmount:       r.TotalAmount,
		})
	}
	return res, nil
}

func (r *transactionStatsByMerchantRepository) GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantFailedRow, error) {

	var results []models.TransactionYearlyMethodByMerchantFailedRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			EXTRACT(YEAR FROM created_at)::text AS year,
			payment_method,
			COUNT(*) AS total_transactions,
			COALESCE(SUM(amount), 0) AS total_amount
		FROM transactions
		WHERE merchant_id = ? AND EXTRACT(YEAR FROM created_at) = ? AND payment_status = 'failed' AND deleted_at IS NULL
		GROUP BY EXTRACT(YEAR FROM created_at), payment_method
		ORDER BY payment_method
	`, req.MerchantID, req.Year).Scan(&results).Error
	if err != nil {
		return nil, transaction_errors.ErrGetYearlyTransactionMethodByMerchant
	}
	var res []*models.TransactionYearlyMethodByMerchantFailedRow
	for _, r := range results {
		res = append(res, &models.TransactionYearlyMethodByMerchantFailedRow{
			Year:              r.Year,
			PaymentMethod:     r.PaymentMethod,
			TotalTransactions: r.TotalTransactions,
			TotalAmount:       r.TotalAmount,
		})
	}
	return res, nil
}
