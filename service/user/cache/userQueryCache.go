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
	userAllCacheKey     = "user:all:page:%d:pageSize:%d:search:%s"
	userByIdCacheKey    = "user:id:%d"
	userActiveCacheKey  = "user:active:page:%d:pageSize:%d:search:%s"
	userTrashedCacheKey = "user:trashed:page:%d:pageSize:%d:search:%s"

	ttlDefault = 5 * time.Minute
)

type userCacheResponse struct {
	Data         []*models.UserRow `json:"data"`
	TotalRecords *int              `json:"total_records"`
}

type userCacheResponseActive struct {
	Data         []*models.UserActiveRow `json:"data"`
	TotalRecords *int                    `json:"total_records"`
}

type userCacheResponseTrashed struct {
	Data         []*models.UserTrashedRow `json:"data"`
	TotalRecords *int                     `json:"total_records"`
}

type userQueryCache struct {
	store *cache.CacheStore
}

func NewUserQueryCache(store *cache.CacheStore) UserQueryCache {
	return &userQueryCache{store: store}
}

func (s *userQueryCache) GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserRow, *int, bool) {
	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[userCacheResponse](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data []*models.UserRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.UserRow{}
	}

	key := fmt.Sprintf(userAllCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponse{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserActiveRow, *int, bool) {
	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[userCacheResponseActive](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data []*models.UserActiveRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.UserActiveRow{}
	}

	key := fmt.Sprintf(userActiveCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseActive{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) ([]*models.UserTrashedRow, *int, bool) {
	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)

	result, found := cache.GetFromCache[userCacheResponseTrashed](ctx, s.store, key)
	if !found || result == nil {
		return nil, nil, false
	}

	return result.Data, result.TotalRecords, true
}

func (s *userQueryCache) SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data []*models.UserTrashedRow, total *int) {
	if total == nil {
		zero := 0
		total = &zero
	}

	if data == nil {
		data = []*models.UserTrashedRow{}
	}

	key := fmt.Sprintf(userTrashedCacheKey, req.Page, req.PageSize, req.Search)
	payload := &userCacheResponseTrashed{Data: data, TotalRecords: total}
	cache.SetToCache(ctx, s.store, key, payload, ttlDefault)
}

func (s *userQueryCache) GetCachedUserCache(ctx context.Context, id int) (*models.User, bool) {
	key := fmt.Sprintf(userByIdCacheKey, id)

	result, found := cache.GetFromCache[*models.User](ctx, s.store, key)
	if !found || result == nil {
		return nil, false
	}

	return *result, true
}

func (s *userQueryCache) SetCachedUserCache(ctx context.Context, data *models.User) {
	if data == nil {
		return
	}

	key := fmt.Sprintf(userByIdCacheKey, data.UserID)
	cache.SetToCache(ctx, s.store, key, data, ttlDefault)
}
