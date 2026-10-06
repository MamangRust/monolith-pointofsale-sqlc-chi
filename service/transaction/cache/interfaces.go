package mencache

import (
	"context"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

type TransactionStatsCache interface {
	GetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*models.TransactionMonthlyAmountSuccessRow, bool)
	SetCachedMonthAmountSuccessCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*models.TransactionMonthlyAmountSuccessRow)

	GetCachedYearAmountSuccessCached(ctx context.Context, year int) ([]*models.TransactionYearlyAmountSuccessRow, bool)
	SetCachedYearAmountSuccessCached(ctx context.Context, year int, res []*models.TransactionYearlyAmountSuccessRow)

	GetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction) ([]*models.TransactionMonthlyAmountFailedRow, bool)
	SetCachedMonthAmountFailedCached(ctx context.Context, req *requests.MonthAmountTransaction, res []*models.TransactionMonthlyAmountFailedRow)

	GetCachedYearAmountFailedCached(ctx context.Context, year int) ([]*models.TransactionYearlyAmountFailedRow, bool)
	SetCachedYearAmountFailedCached(ctx context.Context, year int, res []*models.TransactionYearlyAmountFailedRow)

	GetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*models.TransactionMonthlyMethodSuccessRow, bool)
	SetCachedMonthMethodSuccessCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*models.TransactionMonthlyMethodSuccessRow)

	GetCachedYearMethodSuccessCached(ctx context.Context, year int) ([]*models.TransactionYearlyMethodSuccessRow, bool)
	SetCachedYearMethodSuccessCached(ctx context.Context, year int, res []*models.TransactionYearlyMethodSuccessRow)

	GetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction) ([]*models.TransactionMonthlyMethodFailedRow, bool)
	SetCachedMonthMethodFailedCached(ctx context.Context, req *requests.MonthMethodTransaction, res []*models.TransactionMonthlyMethodFailedRow)

	GetCachedYearMethodFailedCached(ctx context.Context, year int) ([]*models.TransactionYearlyMethodFailedRow, bool)
	SetCachedYearMethodFailedCached(ctx context.Context, year int, res []*models.TransactionYearlyMethodFailedRow)
}

type TransactionStatsByMerchantCache interface {
	GetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountSuccessByMerchantRow, bool)
	SetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*models.TransactionMonthlyAmountSuccessByMerchantRow)

	GetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountSuccessByMerchantRow, bool)
	SetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*models.TransactionYearlyAmountSuccessByMerchantRow)

	GetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountFailedByMerchantRow, bool)
	SetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*models.TransactionMonthlyAmountFailedByMerchantRow)

	GetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountFailedByMerchantRow, bool)
	SetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*models.TransactionYearlyAmountFailedByMerchantRow)

	GetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantSuccessRow, bool)
	SetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*models.TransactionMonthlyMethodByMerchantSuccessRow)

	GetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantSuccessRow, bool)
	SetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*models.TransactionYearlyMethodByMerchantSuccessRow)

	GetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantFailedRow, bool)
	SetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*models.TransactionMonthlyMethodByMerchantFailedRow)

	GetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantFailedRow, bool)
	SetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*models.TransactionYearlyMethodByMerchantFailedRow)
}

type TransactionQueryCache interface {
	GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionRow, *int, bool)
	SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data []*models.TransactionRow, total *int)

	GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*models.TransactionByMerchantRow, *int, bool)
	SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data []*models.TransactionByMerchantRow, total *int)

	GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionActiveRow, *int, bool)
	SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data []*models.TransactionActiveRow, total *int)

	GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionTrashedRow, *int, bool)
	SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data []*models.TransactionTrashedRow, total *int)

	GetCachedTransactionCache(ctx context.Context, id int) (*models.Transaction, bool)
	SetCachedTransactionCache(ctx context.Context, data *models.Transaction)

	GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*models.Transaction, bool)
	SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *models.Transaction)
}

type TransactionCommandCache interface {
	DeleteTransactionCache(ctx context.Context, transactionID int)
	DeleteTransactionAllCache(ctx context.Context)
}
