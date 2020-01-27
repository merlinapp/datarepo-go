package memory

import (
	"context"
	"encoding/json"
	"github.com/coocood/freecache"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/drreflect"
	"log"
	"time"
)

type memoryBasedCacheStore struct {
	cache *freecache.Cache
}

// Creates a new BookCacheStore backed by a freecache Cache (github.com/coocood/freecache)
//
// Implementation Notes: Currently this implementation serializes the data to JSON
// for storage in the memory cache
func NewFreeCacheInMemoryStore(cacheSize int) datarepo.CacheStore {
	cache := freecache.NewCache(cacheSize)

	store := memoryBasedCacheStore{
		cache: cache,
	}
	return &store
}

func (c *memoryBasedCacheStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	bytesToCache, err := cacheMarshal(value)
	if err != nil {
		log.Println("Error setting cache value for key: ", key, "-", err)
		return
	}

	if err = c.cache.Set([]byte(key), bytesToCache, int(expiration.Seconds())); err != nil {
		log.Println("Error setting cache value for key: ", key, "-", err)
	}
}

func (c *memoryBasedCacheStore) Get(ctx context.Context, key string, out interface{}) (bool, error) {
	cachedBytes, err := c.cache.Get([]byte(key))
	if err != nil {
		if err == freecache.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	if err = cacheUnmarshal(cachedBytes, out); err != nil {
		return false, err
	}
	return true, nil
}

func (c *memoryBasedCacheStore) Delete(ctx context.Context, key string) error {
	c.cache.Del([]byte(key))
	return nil
}

func (c *memoryBasedCacheStore) GetMulti(ctx context.Context, keys []string, out interface{}) ([]bool, error) {
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

func cacheMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func cacheUnmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, v)
}
