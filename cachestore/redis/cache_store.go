package redis

import (
	"context"
	"encoding/json"
	redisCache "github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/drreflect"
	"log"
	"time"
)

type redisBasedCacheStore struct {
	redisClient *redis.Client
	cache       *redisCache.Codec
}

// Creates a new CacheStore backed by the provided redis client
//
// Implementation Notes: Currently this implementation serializes the data to JSON
// for storage in Redis
func NewRedisCacheStore(redisClient *redis.Client) datarepo.CacheStore {
	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	codec := &redisCache.Codec{
		Redis:     redisClient,
		Marshal:   cacheMarshal,
		Unmarshal: cacheUnmarshal,
	}

	store := redisBasedCacheStore{
		redisClient: redisClient,
		cache:       codec,
	}
	return &store
}

func (c *redisBasedCacheStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) {
	err := c.cache.Set(&redisCache.Item{
		Key:        key,
		Object:     value,
		Expiration: expiration,
	})
	if err != nil {
		log.Println("Error setting cache value for key: ", key, "-", err)
	}
}

func (c *redisBasedCacheStore) Get(ctx context.Context, key string, out interface{}) (bool, error) {
	if err := c.cache.Get(key, out); err != nil {
		if err == redisCache.ErrCacheMiss {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (c *redisBasedCacheStore) Delete(ctx context.Context, key string) error {
	if err := c.cache.Delete(key); err != nil {
		if err != redisCache.ErrCacheMiss {
			return err
		}
	}
	return nil
}

func (c *redisBasedCacheStore) GetMulti(ctx context.Context, keys []string, out interface{}) ([]bool, error) {
	found := make([]bool, len(keys))
	rawResults, err := c.redisClient.MGet(keys...).Result()
	if err != nil {
		return found, err
	}

	// the out interface is expected to be of type: *[]*A assuming this cache stores elements of type A
	sh := drreflect.NewReflectSlicePointerVHandler(out)
	// the type handler represents our type A
	// *[]*A ->  *A               ->        A
	th := sh.ElementTypeHandler().ElementTypeHandler()

	for i, rawResult := range rawResults {
		value := th.NewPtrToElement()
		if rawResult != nil {
			rawString := rawResult.(string)
			err = cacheUnmarshal([]byte(rawString), value.Ptr())

			if err == nil {
				found[i] = true
			}
		}
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
