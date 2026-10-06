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
	orderAllCacheKey      = "order:all:page:%d:pageSize:%d:search:%s"
	orderByIdCacheKey     = "order:id:%d"
	orderMerchantCacheKey = "order:merchant:%d:page:%d:pageSize:%d:search:%s"
	orderActiveCacheKey   = "order:active:page:%d:pageSize:%d:search:%s"
	orderTrashedCacheKey  = "order:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type orderCacheResponse struct {
	Data         []*models.OrderRow `json:"data"`
	TotalRecords *int               `json:"total_records"`
}

type orderCacheResponseByMerchant struct {
	Data         []*models.OrderByMerchantRow `json:"data"`
	TotalRecords *int                         `json:"total_records"`
}

type orderCacheResponseActive struct {
	Data         []*models.OrderActiveRow `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type orderCacheResponseTrashed struct {
	Data         []*models.OrderTrashedRow `json:"data"`
	TotalRecords *int                      `json:"total_records"`
}

type orderQueryCache struct {
	store *cache.CacheStore
}

func NewOrderQueryCache(store *cache.CacheStore) OrderQueryCache {
	return &orderQueryCache{store: store}
}

func (s *orderQueryCache) GetOrderAllCache(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderRow, *int, bool) {
	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[orderCacheResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderAllCache(ctx context.Context, req *requests.FindAllOrders, data []*models.OrderRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.OrderRow{}
	}

	key := fmt.Sprintf(orderAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetCachedOrderMerchant(ctx context.Context, req *requests.FindAllOrderMerchant) ([]*models.OrderByMerchantRow, *int, bool) {
	key := fmt.Sprintf(orderMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[orderCacheResponseByMerchant](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetCachedOrderMerchant(ctx context.Context, req *requests.FindAllOrderMerchant, res []*models.OrderByMerchantRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if res == nil {
		res = []*models.OrderByMerchantRow{}
	}

	key := fmt.Sprintf(orderMerchantCacheKey, req.MerchantID, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponseByMerchant{Data: res, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetOrderActiveCache(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderActiveRow, *int, bool) {
	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[orderCacheResponseActive](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderActiveCache(ctx context.Context, req *requests.FindAllOrders, data []*models.OrderActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.OrderActiveRow{}
	}

	key := fmt.Sprintf(orderActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrders) ([]*models.OrderTrashedRow, *int, bool) {
	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[orderCacheResponseTrashed](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *orderQueryCache) SetOrderTrashedCache(ctx context.Context, req *requests.FindAllOrders, data []*models.OrderTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.OrderTrashedRow{}
	}

	key := fmt.Sprintf(orderTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &orderCacheResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *orderQueryCache) GetCachedOrderCache(ctx context.Context, orderID int) (*models.Order, bool) {
	key := fmt.Sprintf(orderByIdCacheKey, orderID)
	result, found := cache.GetFromCache[models.Order](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *orderQueryCache) SetCachedOrderCache(ctx context.Context, data *models.Order) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(orderByIdCacheKey, data.OrderID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
