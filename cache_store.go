package datarepo

import (
	"context"
	"time"
)

type CacheStore interface {
	// Deletes the provided key from the cache
	Delete(ctx context.Context, key string) error
	// Retrieves the provided key from the cache and places the output value in the out variable
	//
	// The type of out is expected to be a pointer to the element being stored, for example, if we're
	// storing elements of type A, then out is expected to be of type *A
	Get(ctx context.Context, key string, out interface{}) (bool, error)
	// Retrieves the provided key from the cache and places the output value in the out variable
	//
	// The type of out is expected to be a pointer to a slice of pointers of the elements being stored,
	// for example, if we're storing elements of type A, then out is expected to be of type *[]*A
	GetMulti(ctx context.Context, keys []string, out interface{}) ([]bool, error)
	// Sets the key in the cache with the provided value
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration)
}
