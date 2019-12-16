package datarepo

import (
	"context"
	"reflect"
)

// this interface is for internal use only - users of the library shouldn't need to implement this interface
type Handler interface {
	Delete(ctx context.Context, cacheStore CacheStore, key interface{}) error
	DeleteValue(ctx context.Context, cacheStore CacheStore, value interface{}) error
	Get(ctx context.Context, cacheStore CacheStore, key interface{}, fetcher DataFetcher) (Result, error)
	GetMulti(ctx context.Context, cacheStore CacheStore, keys []interface{}, fetcher DataFetcher) ([]Result, error)
	Set(ctx context.Context, cacheStore CacheStore, value interface{}) error

	CachedType() reflect.Type
	CacheKeyPrefix() string
	SingleResultPerKey() bool
}
