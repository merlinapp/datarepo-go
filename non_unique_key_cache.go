package datarepo

import (
	"context"
	"github.com/merlinapp/datarepo-go/drreflect"
	"github.com/spf13/cast"
)

type nonUniqueKeyCacheHandler struct {
	baseCacheHandler
	// Name of the field that defines the subkey that will be used to compare and store elements that
	// belong to the same key
	subKeyFieldName string
	subTypeHandler  drreflect.StructTypeHandler
}

func NonUniqueKeyCache(v interface{}, cacheDef NonUniqueKeyCacheDefinition) Handler {
	th := drreflect.NewReflectStructTypeHandlerFromValue(v)
	definition := nonUniqueKeyCacheHandler{
		baseCacheHandler{
			keyPrefix:    cacheDef.KeyPrefix,
			keyFieldName: cacheDef.KeyFieldName,
			expiration:   cacheDef.Expiration,
			typeHandler:  th.SlicePtrTypeHandler(),
		},
		cacheDef.SubKeyFieldName,
		th,
	}
	definition.validateConfiguration()
	return &definition
}

func (c *nonUniqueKeyCacheHandler) SingleResultPerKey() bool {
	return false
}

func (c *nonUniqueKeyCacheHandler) DeleteValue(ctx context.Context, cacheStore CacheStore, value interface{}) error {
	return cacheStore.Delete(ctx, c.cacheKeyFromValue(value))
}

func (c *nonUniqueKeyCacheHandler) Set(ctx context.Context, cacheStore CacheStore, value interface{}) error {
	key := c.cacheKeyFromValue(value)
	if key == c.keyPrefix {
		return nil
	}
	cached := c.typeHandler.NewPtrToElement()
	found, err := cacheStore.Get(ctx, key, cached.Ptr())
	if !found {
		return nil
	}
	if err != nil {
		return err
	}

	return c.setInCache(ctx, cacheStore, key, value, cached.Ptr())
}

func (c *nonUniqueKeyCacheHandler) setInCache(ctx context.Context, cacheStore CacheStore, key string, value interface{}, existent interface{}) error {
	subKey := c.cacheSubKey(value)
	var found bool
	var values interface{}
	if existent != nil {
		sliceHandler := drreflect.NewReflectSlicePointerVHandler(existent)

		found = false
		procFunction := func(_ int, ph drreflect.PointerVHandler) {
			valSubKey := c.cacheSubKey(ph.Element())
			if subKey == valSubKey {
				ph.SetElement(value)
				found = true
			}
		}
		sliceHandler.ForEach(procFunction)
		values = sliceHandler.Element()
	}
	if !found {
		var sliceHandler drreflect.SlicePointerHandler
		if existent == nil {
			sliceHandler = c.typeHandler.NewPtrToSlice()
			sliceHandler.MakeSlice(0, 1)
		} else {
			sliceHandler = drreflect.NewReflectSlicePointerVHandler(existent)
		}
		sliceHandler.Append(value)
		values = sliceHandler.Element()
	}
	cacheStore.Set(ctx, key, values, c.expiration)
	return nil
}

func (c *nonUniqueKeyCacheHandler) cacheKeyFromValue(value interface{}) string {
	return c.cacheKey(cast.ToString(c.getFieldValue(value, c.keyFieldName)))
}

func (c *nonUniqueKeyCacheHandler) cacheSubKey(value interface{}) string {
	return cast.ToString(c.getFieldValue(value, c.subKeyFieldName))
}

func (c *nonUniqueKeyCacheHandler) getFieldValue(value interface{}, fieldName string) interface{} {
	return c.subTypeHandler.GetFieldValue(value, fieldName)
}

func (c *nonUniqueKeyCacheHandler) validateConfiguration() {
	if c.keyFieldName == "" {
		panic("A keyFieldName must be defined")
	}
	if c.subKeyFieldName == "" {
		panic("A subKeyFieldName must be defined for caches of type OneToMany")
	}
}
