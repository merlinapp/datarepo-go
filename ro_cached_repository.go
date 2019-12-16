package datarepo

import (
	"context"
	"errors"
	"github.com/merlinapp/datarepo-go/drreflect"
)

type ReadOnlyCachedRepository interface {
	// Retrieves the data from the repository using the given keyFieldName for the given id value.
	//
	// - If the keyFieldName is a Unique Key of the entity, then the result will be a single element.
	//
	// - If the keyFieldName is not a Unique Key of the entity, then the result will be a slice of elements.
	FindByKey(ctx context.Context, keyFieldName string, id interface{}) (Result, error)
	// Retrieves the data from the repository using the given keyFieldName for the given ids.
	//
	// - If the keyFieldName is a Unique Key of the entity, then the result will be a single element per id.
	//
	// - If the keyFieldName is not a Unique Key of the entity, then the result will be a slice of elements per id.
	//
	// Each element in the returned slice corresponds to an id, that is, the result in position 0 corresponds to
	// the id in position 0
	//
	// Ids is expected to be a pointer to a slice or a slice of the corresponding type stored in the keyFieldName,
	// for example, if the keyFieldName stores strings, then ids is expected to be of type *[]string or []string
	FindByKeys(ctx context.Context, keyFieldName string, ids interface{}) ([]Result, error)
}

type readOnlyCachedRepository struct {
	caches map[string]Cache
}

func (r *readOnlyCachedRepository) FindByKey(ctx context.Context, keyFieldName string, id interface{}) (Result, error) {
	return r.FetchSingleFromCache(ctx, keyFieldName, id)
}

func (r *readOnlyCachedRepository) FindByKeys(ctx context.Context, keyFieldName string, ids interface{}) ([]Result, error) {
	sh := drreflect.NewReflectSliceTypeHandlerFromValue(ids)
	return r.FetchMultiFromCache(ctx, keyFieldName, sh.AsInterfaceSlice(ids))
}

func (r *readOnlyCachedRepository) FetchSingleFromCache(ctx context.Context, keyFieldName string, id interface{}) (Result, error) {
	if cacheConfig, ok := r.caches[keyFieldName]; ok {
		return cacheConfig.Get(ctx, id)
	} else {
		return nil, errors.New("Undefined cache for: " + keyFieldName)
	}
}

func (r *readOnlyCachedRepository) FetchMultiFromCache(ctx context.Context, keyFieldName string, ids []interface{}) ([]Result, error) {
	if cacheConfig, ok := r.caches[keyFieldName]; ok {
		return cacheConfig.GetMulti(ctx, ids)
	} else {
		return nil, errors.New("Undefined cache for: " + keyFieldName)
	}
}
