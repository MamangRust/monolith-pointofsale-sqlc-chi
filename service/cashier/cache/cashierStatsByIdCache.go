package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
)

const (
	cashierStatsMonthTotalSalesByIdCacheKey = "cashier:stats:month:%d:year:%d:id:%d"
	cashierStatsYearTotalSalesByIdCacheKey  = "cashier:stats:year:%d:id:%d"

	cashierStatsMonthSalesByIdCacheKey = "cashier:stats:month:%d:id:%d"
	cashierStatsYearSalesByIdCacheKey  = "cashier:stats:year:%d:id:%d"
)

type cashierStatsByIdCache struct {
	store *cache.CacheStore
}

func NewCashierStatsByIdCache(store *cache.CacheStore) CashierStatsByIdCache {
	return &cashierStatsByIdCache{store: store}
}

func (s *cashierStatsByIdCache) GetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier) ([]*models.CashierMonthlyTotalSalesRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByIdCacheKey, req.Month, req.Year, req.CashierID)
	result, found := cache.GetFromCache[[]*models.CashierMonthlyTotalSalesRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetMonthlyTotalSalesByIdCache(ctx context.Context, req *requests.MonthTotalSalesCashier, res []*models.CashierMonthlyTotalSalesRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthTotalSalesByIdCacheKey, req.Month, req.Year, req.CashierID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier) ([]*models.CashierYearlyTotalSalesRow, bool) {
	key := fmt.Sprintf(cashierStatsYearTotalSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := cache.GetFromCache[[]*models.CashierYearlyTotalSalesRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetYearlyTotalSalesByIdCache(ctx context.Context, req *requests.YearTotalSalesCashier, res []*models.CashierYearlyTotalSalesRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearTotalSalesByIdCacheKey, req.Year, req.CashierID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId) ([]*models.CashierMonthlySaleByIdRow, bool) {
	key := fmt.Sprintf(cashierStatsMonthSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := cache.GetFromCache[[]*models.CashierMonthlySaleByIdRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetMonthlyCashierByIdCache(ctx context.Context, req *requests.MonthCashierId, res []*models.CashierMonthlySaleByIdRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsMonthSalesByIdCacheKey, req.Year, req.CashierID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}

func (s *cashierStatsByIdCache) GetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId) ([]*models.CashierYearlySaleRow, bool) {
	key := fmt.Sprintf(cashierStatsYearSalesByIdCacheKey, req.Year, req.CashierID)
	result, found := cache.GetFromCache[[]*models.CashierYearlySaleRow](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return *result, true
}

func (s *cashierStatsByIdCache) SetYearlyCashierByIdCache(ctx context.Context, req *requests.YearCashierId, res []*models.CashierYearlySaleRow) {
	if res == nil {
		return
	}
	key := fmt.Sprintf(cashierStatsYearSalesByIdCacheKey, req.Year, req.CashierID)
	cache.SetToCache(ctx, s.store, key, &res, ttlDefault)
}
