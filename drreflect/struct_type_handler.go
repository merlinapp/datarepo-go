package drreflect

type StructTypeHandler interface {
	TypeHandler
	// Returns the value of the specified field
	//
	// If the struct is of type A, the `input` to this function is expected to be of type *A or of type A
	//
	// This function panics if the input is not of the expected type
	GetFieldValue(input interface{}, fieldName string) interface{}
}

