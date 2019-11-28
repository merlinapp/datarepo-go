package drreflect

import (
	"reflect"
)

var sliceTypeHandlerCache map[reflect.Type]*reflectSliceTypeHandler

type reflectSliceTypeHandler struct {
	*reflectTypeHandler
}

func NewReflectSliceTypeHandlerFromValue(v interface{}) SliceTypeHandler {
	return NewReflectSliceTypeHandler(reflect.TypeOf(v))
}

func NewReflectSliceTypeHandler(t reflect.Type) *reflectSliceTypeHandler {
	if sliceTypeHandlerCache == nil {
		sliceTypeHandlerCache = make(map[reflect.Type]*reflectSliceTypeHandler)
	}
	if sa, ok := sliceTypeHandlerCache[t]; ok {
		return sa
	}

	if t.Kind() == reflect.Ptr {
		if t.Elem().Kind() != reflect.Slice {
			panic("provided type was a pointer but not to a slice: " + t.String())
		}
		t = t.Elem()
	}
	if t.Kind() != reflect.Slice {
		panic("provided type wasn't a slice or a pointer to a slice: " + t.String())
	}

	sth := reflectSliceTypeHandler{
		NewReflectTypeHandler(t),
	}
	sliceTypeHandlerCache[t] = &sth
	return &sth
}

func (r *reflectSliceTypeHandler) AsInterfaceSlice(value interface{}) []interface{} {
	sliceValue := reflect.ValueOf(value)
	sliceType := sliceValue.Type()
	if sliceType.Kind() == reflect.Ptr {
		sliceType = sliceType.Elem()
		sliceValue = sliceValue.Elem()
	}
	if sliceType != r.t {
		panic("provided value wasn't a slice or a pointer to a slice of the expected type: " + sliceType.String())
	}
	result := make([]interface{}, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		result[i] = sliceValue.Index(i).Interface()
	}
	return result
}
