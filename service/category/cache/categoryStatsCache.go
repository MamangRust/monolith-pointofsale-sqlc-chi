package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

const (
	categoryStatsMonthTotalPriceCacheKey = "category:stats:month:%d:year:%d"
	categoryStatsYearTotalPriceCacheKey  = "category:stats:year:%d"

	categoryStatsMonthPriceCacheKey = "category:stats:month:%d"
	categoryStatsYearPriceCacheKey  = "category:stats:year:%d"
)

type categoryStatsCache struct {
	store *cache.CacheStore
}

func NewCategoryStatsCache(store *cache.CacheStore) CategoryStatsCache {
	return &categoryStatsCache{store: store}
}

func (s *categoryStatsCache) GetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice) ([]*models.CategoryMonthlyTotalPriceRow, bool) {
	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*models.CategoryMonthlyTotalPriceRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *categoryStatsCache) SetCachedMonthTotalPriceCache(ctx context.Context, req *requests.MonthTotalPrice, data []*models.CategoryMonthlyTotalPriceRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsMonthTotalPriceCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearTotalPriceCache(ctx context.Context, year int) ([]*models.CategoryYearlyTotalPriceRow, bool) {
	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	result, found := cache.GetFromCache[[]*models.CategoryYearlyTotalPriceRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *categoryStatsCache) SetCachedYearTotalPriceCache(ctx context.Context, year int, data []*models.CategoryYearlyTotalPriceRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsYearTotalPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedMonthPriceCache(ctx context.Context, year int) ([]*models.CategoryMonthlyRow, bool) {
	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	result, found := cache.GetFromCache[[]*models.CategoryMonthlyRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *categoryStatsCache) SetCachedMonthPriceCache(ctx context.Context, year int, data []*models.CategoryMonthlyRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsMonthPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *categoryStatsCache) GetCachedYearPriceCache(ctx context.Context, year int) ([]*models.CategoryYearlyRow, bool) {
	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	result, found := cache.GetFromCache[[]*models.CategoryYearlyRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *categoryStatsCache) SetCachedYearPriceCache(ctx context.Context, year int, data []*models.CategoryYearlyRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryStatsYearPriceCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
