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
	merchantDocumentAllCacheKey     = "merchant_document:all:page:%d:pageSize:%d:search:%s"
	merchantDocumentByIdCacheKey    = "merchant_document:id:%d"
	merchantDocumentActiveCacheKey  = "merchant_document:active:page:%d:pageSize:%d:search:%s"
	merchantDocumentTrashedCacheKey = "merchant_document:trashed:page:%d:pageSize:%d:search:%s"

	ttlDocumentDefault = 5 * time.Minute
)

type merchantDocumentQueryCachedResponse struct {
	Data         []*models.MerchantDocumentRow `json:"data"`
	TotalRecords *int                          `json:"total_records"`
}

type merchantDocumentQueryCachedResponseActive struct {
	Data         []*models.MerchantDocumentActiveRow `json:"data"`
	TotalRecords *int                                `json:"total_records"`
}

type merchantDocumentQueryCachedResponseTrashed struct {
	Data         []*models.MerchantDocumentTrashedRow `json:"data"`
	TotalRecords *int                                 `json:"total_records"`
}

type merchantDocumentQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantDocumentQueryCache(store *cache.CacheStore) MerchantDocumentQueryCache {
	return &merchantDocumentQueryCache{store: store}
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentRow, *int, bool) {
	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDocumentQueryCachedResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocuments(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.MerchantDocumentRow{}
	}

	key := fmt.Sprintf(merchantDocumentAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDocumentDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentActiveRow, *int, bool) {
	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDocumentQueryCachedResponseActive](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsActive(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.MerchantDocumentActiveRow{}
	}

	key := fmt.Sprintf(merchantDocumentActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDocumentDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments) ([]*models.MerchantDocumentTrashedRow, *int, bool) {
	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantDocumentQueryCachedResponseTrashed](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocumentsTrashed(ctx context.Context, req *requests.FindAllMerchantDocuments, data []*models.MerchantDocumentTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.MerchantDocumentTrashedRow{}
	}

	key := fmt.Sprintf(merchantDocumentTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantDocumentQueryCachedResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDocumentDefault)
}

func (s *merchantDocumentQueryCache) GetCachedMerchantDocument(ctx context.Context, id int) (*models.MerchantDocument, bool) {
	key := fmt.Sprintf(merchantDocumentByIdCacheKey, id)

	result, found := cache.GetFromCache[*models.MerchantDocument](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *merchantDocumentQueryCache) SetCachedMerchantDocument(ctx context.Context, data *models.MerchantDocument) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantDocumentByIdCacheKey, data.DocumentID)

	cache.SetToCache(ctx, s.store, key, data, ttlDocumentDefault)
}
