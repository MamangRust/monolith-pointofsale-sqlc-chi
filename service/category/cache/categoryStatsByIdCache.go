package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

const (
	categoryStatsByIdMonthTotalPriceCacheKey = "category:stats:byid:%d:month:%d:year:%d"
	categoryStatsByIdYearTotalPriceCacheKey  = "category:stats:byid:%d:year:%d"

	categoryStatsByIdMonthPriceCacheKey = "category:stats:byid:%d:month:%d"
	categoryStatsByIdYearPriceCacheKey  = "category:stats:byid:%d:year:%d"
)

type categoryStatsByIdCache struct {
	store *cache.CacheStore
}

func NewCategoryStatsByIdCache(store *cache.CacheStore) CategoryStatsByIdCache {
	return &categoryStatsByIdCache{store: store}
}

func (s *categoryStatsByIdCache) GetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory) ([]*models.CategoryMonthlyTotalPriceByIdRow, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*models.CategoryMonthlyTotalPriceByIdRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthTotalPriceByIdCache(ctx context.Context, req *requests.MonthTotalPriceCategory, data []*models.CategoryMonthlyTotalPriceByIdRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByIdMonthTotalPriceCacheKey, req.CategoryID, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory) ([]*models.CategoryYearlyTotalPriceByIdRow, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)
	result, found := cache.GetFromCache[[]*models.CategoryYearlyTotalPriceByIdRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedYearTotalPriceByIdCache(ctx context.Context, req *requests.YearTotalPriceCategory, data []*models.CategoryYearlyTotalPriceByIdRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByIdYearTotalPriceCacheKey, req.CategoryID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId) ([]*models.CategoryMonthlyByIdRow, bool) {
	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)
	result, found := cache.GetFromCache[[]*models.CategoryMonthlyByIdRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedMonthPriceByIdCache(ctx context.Context, req *requests.MonthPriceId, data []*models.CategoryMonthlyByIdRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsByIdMonthPriceCacheKey, req.CategoryID, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsByIdCache) GetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId) ([]*models.CategoryYearlyByIdRow, bool) {
	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)
	result, found := cache.GetFromCache[[]*models.CategoryYearlyByIdRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *categoryStatsByIdCache) SetCachedYearPriceByIdCache(ctx context.Context, req *requests.YearPriceId, data []*models.CategoryYearlyByIdRow) {
	key := fmt.Sprintf(categoryStatsByIdYearPriceCacheKey, req.CategoryID, req.Year)
	if data == nil {
		return
	}
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
