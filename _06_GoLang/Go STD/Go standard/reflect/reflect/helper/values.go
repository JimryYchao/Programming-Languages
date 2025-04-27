package helper

import r "reflect"

// type Int
type (
	Int  = *vInt
	Uint = *vUint
	Bool = *vBool
)

type vInt struct {
	*vCommon
}

func (i Int) valueof(x r.Value) Value {
	i = &vInt{newValue(x)}
	return i
}
func (i Int) CanInt() bool             { return i.v.CanInt() }
func (i Int) Int() int64               { return i.v.Int() }
func (i Int) OverflowInt(x int64) bool { return i.v.OverflowInt(x) }

type IntSetter = *vIntSetter
type vIntSetter struct {
	*vSetter
}

func (s IntSetter) SetInt(x int64) {
	s.v.SetInt(x)
}

type vUint struct{}
type vBool struct{}

// type vUint struct {}
// type vUint struct {}
// type vUint struct {}
// type vUint struct {}
