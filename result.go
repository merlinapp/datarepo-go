package datarepo

import (
	"github.com/merlinapp/datarepo-go/drreflect"
)

type Result interface {
	IsEmpty() bool
	InjectResult(out interface{})
	StoredValue() interface{}
}

type EmptyResult struct {
}

func (e EmptyResult) IsEmpty() bool {
	return true
}

func (e EmptyResult) InjectResult(out interface{}) {
	panic("cannot inject an empty result")
}

func (e EmptyResult) StoredValue() interface{} {
	panic("cannot get the value of an empty result")
}

type ValueResult struct {
	Value interface{}
}

func (r ValueResult) IsEmpty() bool {
	return false
}

func (r ValueResult) InjectResult(out interface{}) {
	ph := drreflect.NewReflectPointerVHandler(out)
	r.injectResult(ph)
}

func (r ValueResult) StoredValue() interface{} {
	return r.Value
}

func (r ValueResult) injectResult(out drreflect.PointerVHandler) {
	out.SetElement(r.Value)
}

func InjectResults(r []Result, out interface{}) {
	handler := drreflect.NewReflectSlicePointerVHandler(out)
	mappingFunction := func(_ int, pointerHandler drreflect.PointerVHandler, in interface{}) {
		r := in.(Result)
		if r.IsEmpty() {
			pointerHandler.SetZeroElement()
		} else {
			r.InjectResult(pointerHandler.Ptr())
		}
	}
	handler.CopyFrom(r, mappingFunction)
}
