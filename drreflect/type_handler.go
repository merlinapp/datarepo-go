package drreflect

import (
	"reflect"
)

type TypeHandler interface {
	// The type handled by this instance
	Type() reflect.Type
	// Creates a new pointer
	//
	// For example, if the struct is of type A, this will return an instance of type *A pointing to a zero-valued A
	NewPtrToElement() PointerVHandler
	// Creates a new slice of pointers
	//
	// For example, if Type() returns type A, this will return an instance of type *[]A
	NewPtrToSlice() SlicePointerHandler
	// Returns a TypeHandler of the element hold by this type.
	//
	// For example, if Type() returns type *A, then this will return a TypeHandler for type A
	// This method will panic if the type doesn't hold an element.
	ElementTypeHandler() TypeHandler
	// Returns a TypeHandler of a slice of pointers to this type.
	//
	// For example, if Type() returns type A, then this will return a TypeHandler for type []*A
	SlicePtrTypeHandler() TypeHandler
	// Checks if the given input is of the type handled by this handler
	//
	// For example, if Type() returns type A, then this will return true if input is of type A, false otherwise
	IsOfType(input interface{}) bool
	// Checks if the given input is of the pointer type handled by this handler
	//
	// For example, if Type() returns type A, then this will return true if input is of type *A, false otherwise
	IsOfPtrType(input interface{}) bool
}
