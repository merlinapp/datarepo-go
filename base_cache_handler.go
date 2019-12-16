package datarepo

import (
	"context"
	"github.com/merlinapp/datarepo-go/drreflect"
	"github.com/spf13/cast"
	"reflect"
	"time"
)

type baseCacheHandler struct {
	// The key prefix to use when storing an element in the cache store
	keyPrefix string
	// Name of the field that defines the key that will be used to store elements in the cache
	keyFieldName string
	// expiration time of entries in the cache
	expiration time.Duration
	// handler used for reflection purposes - this represents the type of element to be stored in the cache
	typeHandler drreflect.TypeHandler
}

func (c *baseCacheHandler) CachedType() reflect.Type {
	return c.typeHandler.Type()
}

func (c *baseCacheHandler) CacheKeyPrefix() string {
	return c.keyPrefix
}

func (c *baseCacheHandler) Delete(ctx context.Context, cacheStore CacheStore, key interface{}) error {
	return cacheStore.Delete(ctx, c.cacheKey(key))
}

func (c *baseCacheHandler) Get(ctx context.Context, cacheStore CacheStore, key interface{}, fetcher DataFetcher) (Result, error) {
	strKey := c.cacheKey(key)
	cached := c.typeHandler.NewPtrToElement()
	found, err := cacheStore.Get(ctx, strKey, cached.Ptr())
	if err != nil {
		return nil, err
	}

	if !found {
		result, err := fetcher.FindByKey(ctx, c.keyFieldName, key)
		if err != nil {
			return nil, err
		}
		if !result.IsEmpty() {
			cacheStore.Set(ctx, strKey, result.StoredValue(), c.expiration)
		}
		return result, err
	}

	return ValueResult{Value: cached.Ptr()}, nil
}

func (c *baseCacheHandler) GetMulti(ctx context.Context, cacheStore CacheStore, keys []interface{}, fetcher DataFetcher) ([]Result, error) {
	strKeys := make([]string, len(keys))
	for i, key := range keys {
		strKeys[i] = c.cacheKey(key)
	}
	cached := c.typeHandler.NewPtrToSlice()
	cached.MakeSlice(0, len(keys))
	found, err := cacheStore.GetMulti(ctx, strKeys, cached.Ptr())
	if err != nil {
		return nil, err
	}

	missingKeyMap := make(map[interface{}]int)
	missingKeys := make([]interface{}, 0, len(keys))
	results := make([]Result, len(keys))
	proc := func(i int, handler drreflect.PointerVHandler) {
		if !found[i] {
			if _, ok := missingKeyMap[keys[i]]; !ok {
				missingKeyMap[keys[i]] = len(missingKeys)
				missingKeys = append(missingKeys, keys[i])
			}
			results[i] = EmptyResult{}
		} else {
			results[i] = ValueResult{Value: handler.Element()}
		}
	}
	cached.ForEach(proc)

	if len(missingKeys) > 0 {
		missingResults, err := fetcher.FindByKeys(ctx, c.keyFieldName, missingKeys)
		if err != nil {
			return nil, err
		}
		for i, key := range keys {
			if idx, ok := missingKeyMap[key]; ok {
				results[i] = missingResults[idx]
				if !results[i].IsEmpty() {
					cacheStore.Set(ctx, strKeys[i], results[i].StoredValue(), c.expiration)
				}
			}
		}
	}

	return results, err
}

func (c *baseCacheHandler) cacheKey(keyPart interface{}) string {
	return c.keyPrefix + cast.ToString(keyPart)
}
