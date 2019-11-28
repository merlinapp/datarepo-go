package drreflect

import (
	"reflect"
)

type reflectPointerVHandler struct {
	t     reflect.Type
	ptr   reflect.Type
	value reflect.Value
}

func NewReflectPointerVHandler(value interface{}) PointerVHandler {
	return newReflectPointerVHandler(reflect.ValueOf(value))
}

func newReflectPointerVHandler(value reflect.Value) *reflectPointerVHandler {
	if value.Type().Kind() != reflect.Ptr {
		panic("value is not a pointer")
	}

	return &reflectPointerVHandler{
		t:     value.Elem().Type(),
		ptr:   value.Type(),
		value: value,
	}
}

func (r *reflectPointerVHandler) Type() reflect.Type {
	return r.ptr
}

func (r *reflectPointerVHandler) ElementType() reflect.Type {
	return r.t
}

func (r *reflectPointerVHandler) Ptr() interface{} {
	return r.value.Interface()
}

func (r *reflectPointerVHandler) Element() interface{} {
	return r.value.Elem().Interface()
}

func (r *reflectPointerVHandler) SetZeroElement() {
	r.value.Elem().Set(reflect.Zero(r.t))
}

func (r *reflectPointerVHandler) SetElement(value interface{}) {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	if t == r.ptr {
		v = v.Elem()
	}
	r.value.Elem().Set(v)
}
