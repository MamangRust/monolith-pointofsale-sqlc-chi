package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

const (
	monthlyTotalRevenueCacheKey = "order:monthly:totalRevenue:month:%d:year:%d"
	yearlyTotalRevenueCacheKey  = "order:yearly:totalRevenue:year:%d"

	monthlyOrderCacheKey = "order:monthly:order:month:%d"
	yearlyOrderCacheKey  = "order:yearly:order:year:%d"
)

type orderStatsCache struct {
	store *cache.CacheStore
}

func NewOrderStatsCache(store *cache.CacheStore) OrderStatsCache {
	return &orderStatsCache{store: store}
}

func (s *orderStatsCache) GetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue) ([]*models.MonthlyTotalRevenueRow, bool) {
	key := fmt.Sprintf(monthlyTotalRevenueCacheKey, req.Month, req.Year)
	result, found := cache.GetFromCache[[]*models.MonthlyTotalRevenueRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *orderStatsCache) SetMonthlyTotalRevenueCache(ctx context.Context, req *requests.MonthTotalRevenue, data []*models.MonthlyTotalRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyTotalRevenueCacheKey, req.Month, req.Year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetYearlyTotalRevenueCache(ctx context.Context, year int) ([]*models.YearlyTotalRevenueRow, bool) {
	key := fmt.Sprintf(yearlyTotalRevenueCacheKey, year)
	result, found := cache.GetFromCache[[]*models.YearlyTotalRevenueRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *orderStatsCache) SetYearlyTotalRevenueCache(ctx context.Context, year int, data []*models.YearlyTotalRevenueRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyTotalRevenueCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetMonthlyOrderCache(ctx context.Context, year int) ([]*models.MonthlyOrderRow, bool) {
	key := fmt.Sprintf(monthlyOrderCacheKey, year)
	result, found := cache.GetFromCache[[]*models.MonthlyOrderRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *orderStatsCache) SetMonthlyOrderCache(ctx context.Context, year int, data []*models.MonthlyOrderRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(monthlyOrderCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}

func (s *orderStatsCache) GetYearlyOrderCache(ctx context.Context, year int) ([]*models.YearlyOrderRow, bool) {
	key := fmt.Sprintf(yearlyOrderCacheKey, year)
	result, found := cache.GetFromCache[[]*models.YearlyOrderRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *orderStatsCache) SetYearlyOrderCache(ctx context.Context, year int, data []*models.YearlyOrderRow) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(yearlyOrderCacheKey, year)
	cache.SetToCache(ctx, s.store, key, &data, ttlDefault)
}
