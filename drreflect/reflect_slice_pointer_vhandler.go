package drreflect

import (
	"reflect"
)

type reflectSlicePointerVHandler struct {
	*reflectPointerVHandler
}

func NewReflectSlicePointerVHandler(value interface{}) SlicePointerHandler {
	return newReflectSlicePointerVHandler(reflect.ValueOf(value))
}

func newReflectSlicePointerVHandler(value reflect.Value) *reflectSlicePointerVHandler {
	ptrHandler := newReflectPointerVHandler(value)
	if ptrHandler.ElementType().Kind() != reflect.Slice {
		panic("The provided instance is not a Slice, got " + ptrHandler.Type().String())
	}
	return &reflectSlicePointerVHandler{ptrHandler}
}

func (s *reflectSlicePointerVHandler) ElementTypeHandler() TypeHandler {
	e := s.ElementType().Elem()
	return NewReflectTypeHandler(e)
}

func (s *reflectSlicePointerVHandler) MakeSlice(len, cap int) {
	sliceValue := reflect.MakeSlice(s.ElementType(), len, cap)
	s.SetElement(sliceValue.Interface())
}

func (s *reflectSlicePointerVHandler) Append(v interface{}) {
	newSlice := reflect.Append(reflect.ValueOf(s.Element()), reflect.ValueOf(v))
	s.SetElement(newSlice.Interface())
}

func (s *reflectSlicePointerVHandler) CopyFrom(arr interface{}, mapper ElementMapper) {
	arrValue := reflect.ValueOf(arr)
	if arrValue.Type().Kind() == reflect.Ptr {
		arrValue = arrValue.Elem()
	}
	if arrValue.Type().Kind() != reflect.Slice {
		panic("input isn't a Slice: " + arrValue.Type().String())
	}

	sliceValue := reflect.MakeSlice(s.ElementType(), arrValue.Len(), arrValue.Len())
	s.SetElement(sliceValue.Interface())
	for i := 0; i < arrValue.Len(); i++ {
		ph := newReflectPointerVHandler(sliceValue.Index(i).Addr())
		mapper(i, ph, arrValue.Index(i).Interface())
	}
}

func (s *reflectSlicePointerVHandler) AsInterfaceSlice() []interface{} {
	sliceValue := reflect.ValueOf(s.Element())
	result := make([]interface{}, sliceValue.Len())
	for i := 0; i < sliceValue.Len(); i++ {
		result[i] = sliceValue.Index(i).Interface()
	}
	return result
}

func (s *reflectSlicePointerVHandler) ForEach(procFunction ElementProcessor) {
	sliceValue := reflect.ValueOf(s.Element())
	for i := 0; i < sliceValue.Len(); i++ {
		ph := newReflectPointerVHandler(sliceValue.Index(i).Addr())
		procFunction(i, ph)
	}
}

func (s *reflectSlicePointerVHandler) Len() int {
	sliceValue := reflect.ValueOf(s.Element())
	return sliceValue.Len()
}
