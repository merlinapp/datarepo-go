package drreflect

import "reflect"

type PointerVHandler interface {
	// type represented by this pointer handler, for example *A
	Type() reflect.Type
	// type represented by the element this pointer points to,
	// for example if the pointer is of type *A, then this returns A
	ElementType() reflect.Type
	// the value handled by this instance
	// this value is of the same type as returned by Type()
	Ptr() interface{}
	// the value pointed by this instance
	// this value is of the same type as returned by Elem()
	Element() interface{}
	// changes the element to be the zero value for the respective ElementType
	SetZeroElement()
	// changes the element to be the provided value
	SetElement(value interface{})
}
