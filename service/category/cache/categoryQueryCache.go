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
	categoryAllCacheKey     = "category:all:page:%d:pageSize:%d:search:%s"
	categoryByIdCacheKey    = "category:id:%d"
	categoryActiveCacheKey  = "category:active:page:%d:pageSize:%d:search:%s"
	categoryTrashedCacheKey = "category:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type categoryCacheResponse struct {
	Data         []*models.CategoryRow `json:"data"`
	TotalRecords *int                  `json:"totalRecords"`
}

type categoryCacheResponseActive struct {
	Data         []*models.CategoryActiveRow `json:"data"`
	TotalRecords *int                        `json:"totalRecords"`
}

type categoryCacheResponseTrashed struct {
	Data         []*models.CategoryTrashedRow `json:"data"`
	TotalRecords *int                         `json:"totalRecords"`
}

type categoryQueryCache struct {
	store *cache.CacheStore
}

func NewCategoryQueryCache(store *cache.CacheStore) CategoryQueryCache {
	return &categoryQueryCache{store: store}
}

func (s *categoryQueryCache) GetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryRow, *int, bool) {
	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[categoryCacheResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *categoryQueryCache) SetCachedCategoriesCache(ctx context.Context, req *requests.FindAllCategory, data []*models.CategoryRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.CategoryRow{}
	}
	key := fmt.Sprintf(categoryAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &categoryCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryActiveRow, *int, bool) {
	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[categoryCacheResponseActive](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *categoryQueryCache) SetCachedCategoryActiveCache(ctx context.Context, req *requests.FindAllCategory, data []*models.CategoryActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.CategoryActiveRow{}
	}
	key := fmt.Sprintf(categoryActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &categoryCacheResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory) ([]*models.CategoryTrashedRow, *int, bool) {
	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, req.Search)
	result, found := cache.GetFromCache[categoryCacheResponseTrashed](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}
	return result.Data, result.TotalRecords, true
}

func (s *categoryQueryCache) SetCachedCategoryTrashedCache(ctx context.Context, req *requests.FindAllCategory, data []*models.CategoryTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}
	if data == nil {
		data = []*models.CategoryTrashedRow{}
	}
	key := fmt.Sprintf(categoryTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &categoryCacheResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *categoryQueryCache) GetCachedCategoryCache(ctx context.Context, id int) (*models.Category, bool) {
	key := fmt.Sprintf(categoryByIdCacheKey, id)
	result, found := cache.GetFromCache[models.Category](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}
	return result, true
}

func (s *categoryQueryCache) SetCachedCategoryCache(ctx context.Context, data *models.Category) {
	if data == nil {
		return
	}
	key := fmt.Sprintf(categoryByIdCacheKey, data.CategoryID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
