package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

const (
	transactonMonthAmountSuccessByMerchantKey = "transaction:month:amount:success:merchant:%d:month:%d:year:%d"
	transactonMonthAmountFailedByMerchantKey  = "transaction:month:amount:failed:merchant:%d:month:%d:year:%d"

	transactonYearAmountSuccessByMerchantKey = "transaction:year:amount:success:merchant:%d:year:%d"
	transactonYearAmountFailedByMerchantKey  = "transaction:year:amount:failed:merchant:%d:year:%d"

	transactonMonthMethodSuccessByMerchantKey = "transaction:month:method:success:merchant:%d:month:%d:year:%d"
	transactonMonthMethodFailedByMerchantKey  = "transaction:month:method:failed:merchant:%d:month:%d:year:%d"

	transactonYearMethodSuccessByMerchantKey = "transaction:year:method:success:merchant:%d:year:%d"
	transactonYearMethodFailedByMerchantKey  = "transaction:year:method:failed:merchant:%d:year:%d"
)

type transactionStatsByMerchantCache struct {
	store *cache.CacheStore
}

func NewTransactionStatsByMerchantCache(store *cache.CacheStore) TransactionStatsByMerchantCache {
	return &transactionStatsByMerchantCache{store: store}
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*models.TransactionMonthlyAmountSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountSuccessByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*models.TransactionMonthlyAmountSuccessByMerchantRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthAmountSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant) ([]*models.TransactionMonthlyAmountFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*models.TransactionMonthlyAmountFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthAmountFailedByMerchantCached(ctx context.Context, req *requests.MonthAmountTransactionMerchant, res []*models.TransactionMonthlyAmountFailedByMerchantRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthAmountFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountFailedByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*models.TransactionYearlyAmountFailedByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountFailedByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*models.TransactionYearlyAmountFailedByMerchantRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearAmountFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant) ([]*models.TransactionYearlyAmountSuccessByMerchantRow, bool) {
	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*models.TransactionYearlyAmountSuccessByMerchantRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearAmountSuccessByMerchantCached(ctx context.Context, req *requests.YearAmountTransactionMerchant, res []*models.TransactionYearlyAmountSuccessByMerchantRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearAmountSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantSuccessRow, bool) {
	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*models.TransactionMonthlyMethodByMerchantSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodSuccessByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*models.TransactionMonthlyMethodByMerchantSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthMethodSuccessByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantSuccessRow, bool) {
	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*models.TransactionYearlyMethodByMerchantSuccessRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodSuccessByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*models.TransactionYearlyMethodByMerchantSuccessRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearMethodSuccessByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant) ([]*models.TransactionMonthlyMethodByMerchantFailedRow, bool) {
	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)

	result, found := cache.GetFromCache[[]*models.TransactionMonthlyMethodByMerchantFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedMonthMethodFailedByMerchantCached(ctx context.Context, req *requests.MonthMethodTransactionMerchant, res []*models.TransactionMonthlyMethodByMerchantFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonMonthMethodFailedByMerchantKey, req.MerchantID, req.Month, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}

func (t *transactionStatsByMerchantCache) GetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant) ([]*models.TransactionYearlyMethodByMerchantFailedRow, bool) {
	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)

	result, found := cache.GetFromCache[[]*models.TransactionYearlyMethodByMerchantFailedRow](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (t *transactionStatsByMerchantCache) SetCachedYearMethodFailedByMerchantCached(ctx context.Context, req *requests.YearMethodTransactionMerchant, res []*models.TransactionYearlyMethodByMerchantFailedRow) {
	if res == nil {
		return
	}

	key := fmt.Sprintf(transactonYearMethodFailedByMerchantKey, req.MerchantID, req.Year)
	cache.SetToCache(ctx, t.store, key, &res, ttlDefault)
}
