package datarepo

import (
	"context"
)

type DataFetcher interface {
	FindByKey(ctx context.Context, keyFieldName string, id interface{}) (Result, error)
	FindByKeys(ctx context.Context, keyFieldName string, ids []interface{}) ([]Result, error)
}
