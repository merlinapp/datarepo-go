package drreflect

type SlicePointerHandler interface {
	PointerVHandler
	// Returns a TypeHandler for the elements of this slice. For example, if this slice is of type
	// []A, then the returned TypeHandler will be for type A
	ElementTypeHandler() TypeHandler
	Append(v interface{})
	MakeSlice(len, cap int)
	// Copies the values from the provided array into this slice using the
	// mapFunction to convert from one type to another. The mapFunction is called
	// once per element in the input array
	CopyFrom(arr interface{}, mapFunction ElementMapper)
	// Executes the provided processingFunction for each element in the slice
	ForEach(processingFunction ElementProcessor)
	// Returns the length of this slice
	Len() int
	// Returns this slice as a slice of interface{} instances.
	AsInterfaceSlice() []interface{}
}

// An ElementMapper receives a pointer to an element and an input object for processing
//
// The ElementMapper can be used to update the value of an element based on the input object
type ElementMapper func(idx int, h PointerVHandler, in interface{})

// An ElementProcessor receives a pointer to an element and can be used to update the value of an
// element based on external logic
type ElementProcessor func(idx int, h PointerVHandler)
