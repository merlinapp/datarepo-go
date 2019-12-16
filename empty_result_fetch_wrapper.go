package datarepo

import (
	"context"
	"github.com/merlinapp/datarepo-go/drreflect"
)

type emptyResultDataFetcherWrapper struct {
	dataType interface{}
	delegate DataFetcher
}

func (w *emptyResultDataFetcherWrapper) FindByKey(ctx context.Context, keyFieldName string, id interface{}) (Result, error) {
	result, err := w.delegate.FindByKey(ctx, keyFieldName, id)
	if err != nil {
		return nil, err
	}
	if result.IsEmpty() {
		th := drreflect.NewReflectStructTypeHandlerFromValue(w.dataType)
		slice := th.NewPtrToSlice()
		slice.MakeSlice(0, 1)
		result = ValueResult{Value: slice.Ptr()}
	}
	return result, err
}

func (w *emptyResultDataFetcherWrapper) FindByKeys(ctx context.Context, keyFieldName string, ids []interface{}) ([]Result, error) {
	result, err := w.delegate.FindByKeys(ctx, keyFieldName, ids)
	if err != nil {
		return nil, err
	}
	arr := result
	for i := 0; i < len(arr); i++ {
		if arr[i].IsEmpty() {
			th := drreflect.NewReflectStructTypeHandlerFromValue(w.dataType)
			slice := th.NewPtrToSlice()
			slice.MakeSlice(0, 1)
			arr[i] = ValueResult{Value: slice.Ptr()}
		}
	}
	return result, err
}
