package mencache

import (
	"context"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

const (
	transactionAllCacheKey  = "transaction:all:page:%d:pageSize:%d:search:%s"
	transactionByIdCacheKey = "transaction:id:%d"

	transactionByMerchantCacheKey = "transaction:merchant:%d:page:%d:pageSize:%d:search:%s"

	transactionActiveCacheKey  = "transaction:active:page:%d:pageSize:%d:search:%s"
	transactionTrashedCacheKey = "transaction:trashed:page:%d:pageSize:%d:search:%s"

	transactionByOrderCacheKey = "transaction:order:%d"

	ttlDefault = 5 * time.Minute
)

type transactionCacheResponse struct {
	Data         []*models.TransactionRow `json:"data"`
	TotalRecords *int                     `json:"totalRecords"`
}

type transactionMerchantCacheResponse struct {
	Data         []*models.TransactionByMerchantRow `json:"data"`
	TotalRecords *int                               `json:"totalRecords"`
}

type transactionCacheResponseActive struct {
	Data         []*models.TransactionActiveRow `json:"data"`
	TotalRecords *int                           `json:"totalRecords"`
}

type transactionCacheResponseTrashed struct {
	Data         []*models.TransactionTrashedRow `json:"data"`
	TotalRecords *int                            `json:"totalRecords"`
}

type transactionQueryCache struct {
	store *cache.CacheStore
}

func NewTransactionQueryCache(store *cache.CacheStore) TransactionQueryCache {
	return &transactionQueryCache{store: store}
}

func (t *transactionQueryCache) GetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionRow, *int, bool) {
	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[transactionCacheResponse](ctx, t.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionsCache(ctx context.Context, req *requests.FindAllTransaction, data []*models.TransactionRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.TransactionRow{}
	}

	key := fmt.Sprintf(transactionAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant) ([]*models.TransactionByMerchantRow, *int, bool) {
	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[transactionMerchantCacheResponse](ctx, t.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionByMerchant(ctx context.Context, req *requests.FindAllTransactionByMerchant, data []*models.TransactionByMerchantRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.TransactionByMerchantRow{}
	}

	key := fmt.Sprintf(transactionByMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &transactionMerchantCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionActiveRow, *int, bool) {
	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[transactionCacheResponseActive](ctx, t.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionActiveCache(ctx context.Context, req *requests.FindAllTransaction, data []*models.TransactionActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.TransactionActiveRow{}
	}

	key := fmt.Sprintf(transactionActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCacheResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction) ([]*models.TransactionTrashedRow, *int, bool) {
	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[transactionCacheResponseTrashed](ctx, t.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (t *transactionQueryCache) SetCachedTransactionTrashedCache(ctx context.Context, req *requests.FindAllTransaction, data []*models.TransactionTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.TransactionTrashedRow{}
	}

	key := fmt.Sprintf(transactionTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &transactionCacheResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, t.store, key, payload, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionCache(ctx context.Context, id int) (*models.Transaction, bool) {
	key := fmt.Sprintf(transactionByIdCacheKey, id)

	result, found := cache.GetFromCache[models.Transaction](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionCache(ctx context.Context, data *models.Transaction) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transactionByIdCacheKey, data.TransactionID)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}

func (t *transactionQueryCache) GetCachedTransactionByOrderId(ctx context.Context, orderID int) (*models.Transaction, bool) {
	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)

	result, found := cache.GetFromCache[models.Transaction](ctx, t.store, key)
	if !found || result == nil {
		return nil, false
	}

	return result, true
}

func (t *transactionQueryCache) SetCachedTransactionByOrderId(ctx context.Context, orderID int, data *models.Transaction) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(transactionByOrderCacheKey, orderID)
	cache.SetToCache(ctx, t.store, key, data, ttlDefault)
}
