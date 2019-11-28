package datarepo

import (
	"context"
)

type Cache struct {
	Handler     Handler
	Store       CacheStore
	DataFetcher DataFetcher
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	return c.Handler.Delete(ctx, c.Store, key)
}

func (c *Cache) DeleteValue(ctx context.Context, value interface{}) error {
	return c.Handler.DeleteValue(ctx, c.Store, value)
}

func (c *Cache) Get(ctx context.Context, key interface{}) (Result, error) {
	return c.Handler.Get(ctx, c.Store, key, c.DataFetcher)
}

func (c *Cache) GetMulti(ctx context.Context, keys []interface{}) ([]Result, error) {
	return c.Handler.GetMulti(ctx, c.Store, keys, c.DataFetcher)
}

func (c *Cache) Set(ctx context.Context, value interface{}) error {
	return c.Handler.Set(ctx, c.Store, value)
}
