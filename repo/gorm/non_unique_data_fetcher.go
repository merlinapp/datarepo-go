package gorm

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/merlinapp/datarepo-go"
	"github.com/merlinapp/datarepo-go/drreflect"
)

type nonUniqueDataFetcher struct {
	db                *gorm.DB
	typeHandler       drreflect.StructTypeHandler
	fieldToColumnName map[string]string
}

func (u *nonUniqueDataFetcher) FindByKey(ctx context.Context, keyFieldName string, id interface{}) (datarepo.Result, error) {
	result, err := u.FindByKeys(ctx, keyFieldName, []interface{}{id})
	if err != nil {
		return nil, err
	}
	return result[0], err
}

func (u *nonUniqueDataFetcher) FindByKeys(ctx context.Context, keyFieldName string, ids []interface{}) ([]datarepo.Result, error) {
	dataSlice := u.typeHandler.NewPtrToSlice()
	columnName, ok := u.fieldToColumnName[keyFieldName]
	if !ok {
		return nil, errors.New("column name not defined for: " + keyFieldName)
	}

	err := u.db.Find(dataSlice.Ptr(), columnName+" IN (?)", ids).Error
	if err != nil {
		return nil, err
	}

	result := make([]datarepo.Result, len(ids))
	resultsPerId := make(map[interface{}]drreflect.SlicePointerHandler)
	proc := func(_ int, handler drreflect.PointerVHandler) {
		keyValue := u.typeHandler.GetFieldValue(handler.Element(), keyFieldName)
		if _, ok := resultsPerId[keyValue]; !ok {
			resultsPerId[keyValue] = u.typeHandler.NewPtrToSlice()
		}
		resultsPerId[keyValue].Append(handler.Element())
	}
	dataSlice.ForEach(proc)
	for i, id := range ids {
		if value, ok := resultsPerId[id]; ok {
			result[i] = datarepo.ValueResult{Value: value.Ptr()}
		} else {
			result[i] = datarepo.EmptyResult{}
		}
	}

	return result, err
}
