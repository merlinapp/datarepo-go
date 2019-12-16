package drreflect

import (
	"reflect"
)

var typeHandlerCache map[reflect.Type]*reflectTypeHandler

type reflectTypeHandler struct {
	t        reflect.Type
	ptr      reflect.Type
	slice    reflect.Type
	slicePtr reflect.Type
}

func NewReflectTypeHandlerFromValue(v interface{}) TypeHandler {
	return NewReflectTypeHandler(reflect.TypeOf(v))
}

func NewReflectTypeHandler(t reflect.Type) *reflectTypeHandler {
	if typeHandlerCache == nil {
		typeHandlerCache = make(map[reflect.Type]*reflectTypeHandler)
	}

	if th, ok := typeHandlerCache[t]; ok {
		return th
	}

	ptr := reflect.PtrTo(t)
	th := reflectTypeHandler{
		t:        t,
		ptr:      ptr,
		slice:    reflect.SliceOf(t),
		slicePtr: reflect.SliceOf(ptr),
	}
	typeHandlerCache[t] = &th
	return &th
}

func (r *reflectTypeHandler) Type() reflect.Type {
	return r.t
}

func (r *reflectTypeHandler) NewPtrToElement() PointerVHandler {
	return NewReflectPointerVHandler(reflect.New(r.t).Interface())
}

func (r *reflectTypeHandler) NewPtrToSlice() SlicePointerHandler {
	return NewReflectSlicePointerVHandler(reflect.New(r.slicePtr).Interface())
}

func (r *reflectTypeHandler) IsOfType(input interface{}) bool {
	return reflect.TypeOf(input) == r.t
}

func (r *reflectTypeHandler) IsOfPtrType(input interface{}) bool {
	return reflect.TypeOf(input) == r.ptr
}

func (r *reflectTypeHandler) ElementTypeHandler() TypeHandler {
	return NewReflectTypeHandler(r.t.Elem())
}

func (r *reflectTypeHandler) SliceTypeHandler() TypeHandler {
	return NewReflectTypeHandler(r.slice)
}

func (r *reflectTypeHandler) SlicePtrTypeHandler() TypeHandler {
	return NewReflectTypeHandler(r.slicePtr)
}