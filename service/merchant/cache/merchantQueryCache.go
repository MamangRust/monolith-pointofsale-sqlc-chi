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
	merchantAllCacheKey      = "merchant:all:page:%d:pageSize:%d:search:%s"
	merchantByIdCacheKey     = "merchant:id:%d"
	merchantActiveCacheKey   = "merchant:active:page:%d:pageSize:%d:search:%s"
	merchantTrashedCacheKey  = "merchant:trashed:page:%d:pageSize:%d:search:%s"
	merchantByUserIdCacheKey = "merchant:user_id:%d"

	ttlDefault = 5 * time.Minute
)

type merchantCachedResponse struct {
	Data         []*models.MerchantRow `json:"data"`
	TotalRecords *int                  `json:"total_records"`
}

type merchantCachedResponseActive struct {
	Data         []*models.MerchantActiveRow `json:"data"`
	TotalRecords *int                        `json:"total_records"`
}

type merchantCachedResponseTrashed struct {
	Data         []*models.MerchantTrashedRow `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

type merchantQueryCache struct {
	store *cache.CacheStore
}

func NewMerchantQueryCache(store *cache.CacheStore) MerchantQueryCache {
	return &merchantQueryCache{store: store}
}

func (m *merchantQueryCache) GetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantRow, *int, bool) {
	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantCachedResponse](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetCachedMerchants(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.MerchantRow{}
	}

	key := fmt.Sprintf(merchantAllCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantActiveRow, *int, bool) {
	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantCachedResponseActive](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetCachedMerchantActive(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.MerchantActiveRow{}
	}

	key := fmt.Sprintf(merchantActiveCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants) ([]*models.MerchantTrashedRow, *int, bool) {
	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[merchantCachedResponseTrashed](ctx, m.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (m *merchantQueryCache) SetCachedMerchantTrashed(ctx context.Context, req *requests.FindAllMerchants, data []*models.MerchantTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.MerchantTrashedRow{}
	}

	key := fmt.Sprintf(merchantTrashedCacheKey, req.Page, req.PageSize, req.Search)

	payload := &merchantCachedResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, m.store, key, payload, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchant(ctx context.Context, id int) (*models.Merchant, bool) {
	key := fmt.Sprintf(merchantByIdCacheKey, id)

	result, found := cache.GetFromCache[*models.Merchant](ctx, m.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantQueryCache) SetCachedMerchant(ctx context.Context, data *models.Merchant) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(merchantByIdCacheKey, data.MerchantID)

	cache.SetToCache(ctx, m.store, key, data, ttlDefault)
}

func (m *merchantQueryCache) GetCachedMerchantsByUserId(ctx context.Context, id int) ([]*models.Merchant, bool) {
	key := fmt.Sprintf(merchantByUserIdCacheKey, id)

	result, found := cache.GetFromCache[[]*models.Merchant](ctx, m.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (m *merchantQueryCache) SetCachedMerchantsByUserId(ctx context.Context, userId int, data []*models.Merchant) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(merchantByUserIdCacheKey, userId)

	cache.SetToCache(ctx, m.store, key, &data, ttlDefault)
}
