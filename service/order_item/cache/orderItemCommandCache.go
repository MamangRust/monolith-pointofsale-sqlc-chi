package mencache

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
)

type orderItemCommandCache struct {
	store *cache.CacheStore
}

func NewOrderItemCommandCache(store *cache.CacheStore) OrderItemCommandCache {
	return &orderItemCommandCache{store: store}
}

// DeleteCachedOrderItems invalidates the cached order items for a single order.
func (o *orderItemCommandCache) DeleteCachedOrderItems(ctx context.Context, orderID int) {
	key := fmt.Sprintf(orderItemByIdCacheKey, orderID)
	cache.DeleteFromCache(ctx, o.store, key)
}

// DeleteCachedOrderItemsAll clears every order item cache entry. It is used by
// deletes, where the owning order is not known from the request.
func (o *orderItemCommandCache) DeleteCachedOrderItemsAll(ctx context.Context) {
	if o.store == nil {
		return
	}
	_, _ = o.store.InvalidateCache(ctx, "order_item:*")
}
