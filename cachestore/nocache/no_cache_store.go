package nocache

import (
	"context"
	"time"
)

// An implementation of an empty CacheStore that doesn't store any values.
type Store struct {
}

func (c *Store) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) {
}

func (c *Store) Delete(ctx context.Context, key string) error {
	return nil
}

func (c *Store) Get(ctx context.Context, key string, out interface{}) (bool, error) {
	return false, nil
}

func (c *Store) GetMulti(ctx context.Context, keys []string, out interface{}) ([]bool, error) {
	found := make([]bool, len(keys))
	return found, nil
}
