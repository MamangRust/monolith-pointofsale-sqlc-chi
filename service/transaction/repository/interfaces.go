package repository

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	"gorm.io/gorm"
)

type TransactionStatsRepository interface {
	GetMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*models.TransactionMonthlyAmountSuccessRow, error)
	GetYearlyAmountSuccess(ctx context.Context, year int) ([]*models.TransactionYearlyAmountSuccessRow, error)
	GetMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*models.TransactionMonthlyAmountFailedRow, error)
	GetYearlyAmountFailed(ctx context.Context, year int) ([]*models.TransactionYearlyAmountFailedRow, error)

	GetMonthlyTransactionMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*models.TransactionMonthlyMethodSuccessRow, error)
	GetYearlyTransactionMethodSuccess(ctx context.Context, year int) ([]*models.TransactionYearlyMethodSuccessRow, error)
	GetMonthlyTransactionMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*models.TransactionMonthlyMethodFailedRow, error)
	GetYearlyTransactionMethodFailed(ctx context.Context, year int) ([]*models.TransactionYearlyMethodFailedRow, error)
}

type TransactionStatsByMerchantRepository interface {
	GetMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountSuccessByMerchantRow, error)
	GetYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountSuccessByMerchantRow, error)
	GetMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountFailedByMerchantRow, error)
	GetYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountFailedByMerchantRow, error)

	GetMonthlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantSuccessRow, error)
	GetYearlyTransactionMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantSuccessRow, error)
	GetMonthlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantFailedRow, error)
	GetYearlyTransactionMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantFailedRow, error)
}

type TransactionQueryRepository interface {
	FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionTrashedRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*models.TransactionByMerchantRow, *int, error)
	FindById(ctx context.Context, transaction_id int) (*models.Transaction, error)
	FindByOrderId(ctx context.Context, order_id int) (*models.Transaction, error)
}

type TransactionCommandRepository interface {
	CreateTransaction(ctx context.Context, request *requests.CreateTransactionRequest) (*models.Transaction, error)
	// CreateTransactionInTx persists the transaction inside the given GORM
	// transaction so the caller can commit the business write and its outbox
	// event atomically (Phase 6 — transactional outbox).
	CreateTransactionInTx(ctx context.Context, tx *gorm.DB, request *requests.CreateTransactionRequest) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, request *requests.UpdateTransactionRequest) (*models.Transaction, error)
	TrashTransaction(ctx context.Context, transaction_id int) (*models.Transaction, error)
	RestoreTransaction(ctx context.Context, transaction_id int) (*models.Transaction, error)
	DeleteTransactionPermanently(ctx context.Context, transaction_id int) (bool, error)
	RestoreAllTransactions(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}
