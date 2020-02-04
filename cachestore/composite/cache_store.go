package composite

import (
	"context"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/drreflect"
	"time"
)

type compositeCacheStore struct {
	delegates []datarepo.CacheStore
}

// Creates a new composite CacheStore backed by the provider CacheStore implementations.
//
// Read operations will be delegated down to the provided caches in the order they're provided.
// If a cache returns a result (key is found), then that result will be used and further caches will not be queried.
//
// Write operations will be propagates to all delegate caches
func NewCompositeCacheStore(delegates ...datarepo.CacheStore) datarepo.CacheStore {
	if len(delegates) == 0 {
		panic("Can't create a composite cache store with no delegate caches")
	}

	store := compositeCacheStore{
		delegates: delegates,
	}
	return &store
}

func (c *compositeCacheStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	for _, cache := range c.delegates {
		cache.Set(ctx, key, value, expiration)
	}
}

func (c *compositeCacheStore) Get(ctx context.Context, key string, out interface{}) (bool, error) {
	var err error
	for _, cache := range c.delegates {
		found, cErr := cache.Get(ctx, key, out)
		if found {
			return true, cErr
		}
		if cErr != nil {
			err = cErr
		}
	}
	return false, err
}

func (c *compositeCacheStore) Delete(ctx context.Context, key string) error {
	var err error
	for _, cache := range c.delegates {
		if cErr := cache.Delete(ctx, key); cErr != nil {
			err = cErr
		}
	}
	return err
}

func (c *compositeCacheStore) GetMulti(ctx context.Context, keys []string, out interface{}) ([]bool, error) {
	found := make([]bool, len(keys))

	// the out interface is expected to be of type: *[]*A assuming this cache stores elements of type A
	sh := drreflect.NewReflectSlicePointerVHandler(out)
	// the type handler represents our type A
	// *[]*A ->  *A               ->        A
	th := sh.ElementTypeHandler().ElementTypeHandler()

	for i, key := range keys {
		value := th.NewPtrToElement()
		found[i], _ = c.Get(ctx, key, value.Ptr())
		sh.Append(value.Ptr())
	}

	return found, nil
}
