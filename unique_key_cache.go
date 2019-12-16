package datarepo

import (
	"context"
	"github.com/merlinapp/datarepo-go/drreflect"
	"github.com/spf13/cast"
)

type uniqueKeyCacheHandler struct {
	baseCacheHandler
	subTypeHandler drreflect.StructTypeHandler
}

func UniqueKeyCache(v interface{}, cacheDefinition UniqueKeyCacheDefinition) Handler {
	th := drreflect.NewReflectStructTypeHandlerFromValue(v)
	definition := uniqueKeyCacheHandler{
		baseCacheHandler{
			keyPrefix:    cacheDefinition.KeyPrefix,
			keyFieldName: cacheDefinition.KeyFieldName,
			expiration:   cacheDefinition.Expiration,
			typeHandler:  th,
		},
		th,
	}
	definition.validateConfiguration()
	return &definition
}

func (c *uniqueKeyCacheHandler) SingleResultPerKey() bool {
	return true
}

func (c *uniqueKeyCacheHandler) DeleteValue(ctx context.Context, cacheStore CacheStore, value interface{}) error {
	return cacheStore.Delete(ctx, c.cacheKeyFromValue(value))
}

func (c *uniqueKeyCacheHandler) Set(ctx context.Context, cacheStore CacheStore, value interface{}) error {
	cacheStore.Set(ctx, c.cacheKeyFromValue(value), value, c.expiration)
	return nil
}

func (c *uniqueKeyCacheHandler) cacheKeyFromValue(value interface{}) string {
	return c.cacheKey(cast.ToString(c.getFieldValue(value, c.keyFieldName)))
}

func (c *uniqueKeyCacheHandler) getFieldValue(value interface{}, fieldName string) interface{} {
	return c.subTypeHandler.GetFieldValue(value, fieldName)
}

func (c *uniqueKeyCacheHandler) validateConfiguration() {
	if c.keyFieldName == "" {
		panic("A keyFieldName must be defined")
	}
}
