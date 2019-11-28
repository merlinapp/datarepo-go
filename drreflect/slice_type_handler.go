package drreflect

type SliceTypeHandler interface {
	TypeHandler
	// Returns the given slice as a slice of interface{} instances.
	AsInterfaceSlice(value interface{}) []interface{}
}
