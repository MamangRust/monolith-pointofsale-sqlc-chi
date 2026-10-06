package service

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

// EmailEventPublisher publishes email notification events to Kafka topics.
// Implemented by *kafka.Kafka; defined as an interface so it can be faked in
// unit tests (graceful degradation of the email notification path).
// ctx carries the OpenTelemetry trace context, which is injected into the
// Kafka message headers so the consumer continues the same trace.
type EmailEventPublisher interface {
	SendMessage(ctx context.Context, topic string, key string, value []byte) error
}

type TransactionStatsService interface {
	FindMonthlyAmountSuccess(ctx context.Context, req *requests.MonthAmountTransaction) ([]*models.TransactionMonthlyAmountSuccessRow, error)
	FindYearlyAmountSuccess(ctx context.Context, year int) ([]*models.TransactionYearlyAmountSuccessRow, error)
	FindMonthlyAmountFailed(ctx context.Context, req *requests.MonthAmountTransaction) ([]*models.TransactionMonthlyAmountFailedRow, error)
	FindYearlyAmountFailed(ctx context.Context, year int) ([]*models.TransactionYearlyAmountFailedRow, error)

	FindMonthlyMethodSuccess(ctx context.Context, req *requests.MonthMethodTransaction) ([]*models.TransactionMonthlyMethodSuccessRow, error)
	FindYearlyMethodSuccess(ctx context.Context, year int) ([]*models.TransactionYearlyMethodSuccessRow, error)

	FindMonthlyMethodFailed(ctx context.Context, req *requests.MonthMethodTransaction) ([]*models.TransactionMonthlyMethodFailedRow, error)
	FindYearlyMethodFailed(ctx context.Context, year int) ([]*models.TransactionYearlyMethodFailedRow, error)
}

type TransactionStatsByMerchantQueryService interface {
	FindMonthlyAmountSuccessByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountSuccessByMerchantRow, error)
	FindYearlyAmountSuccessByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountSuccessByMerchantRow, error)

	FindMonthlyAmountFailedByMerchant(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountFailedByMerchantRow, error)
	FindYearlyAmountFailedByMerchant(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountFailedByMerchantRow, error)

	FindMonthlyMethodByMerchantSuccess(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantSuccessRow, error)
	FindYearlyMethodByMerchantSuccess(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantSuccessRow, error)

	FindMonthlyMethodByMerchantFailed(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantFailedRow, error)
	FindYearlyMethodByMerchantFailed(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantFailedRow, error)
}

type TransactionQueryService interface {
	FindAllTransactions(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionRow, *int, error)
	FindByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*models.TransactionByMerchantRow, *int, error)
	FindByActive(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionActiveRow, *int, error)
	FindByTrashed(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionTrashedRow, *int, error)
	FindById(ctx context.Context, transactionID int) (*models.Transaction, error)
	FindByOrderId(ctx context.Context, orderID int) (*models.Transaction, error)
}

type TransactionCommandService interface {
	CreateTransaction(ctx context.Context, req *requests.CreateTransactionRequest) (*models.Transaction, error)
	UpdateTransaction(ctx context.Context, req *requests.UpdateTransactionRequest) (*models.Transaction, error)
	TrashedTransaction(ctx context.Context, transaction_id int) (*models.Transaction, error)
	RestoreTransaction(ctx context.Context, transaction_id int) (*models.Transaction, error)
	DeleteTransactionPermanently(ctx context.Context, transactionID int) (bool, error)
	RestoreAllTransactions(ctx context.Context) (bool, error)
	DeleteAllTransactionPermanent(ctx context.Context) (bool, error)
}
