package helper

import (
	r "reflect"
)

//! >>>>>>>>>>>>>> Array <<<<<<<<<<<<<<

type Array = *vArray

type vArray struct {
	*vCommon
	len int
}

func (v Array) valueof(rv r.Value) Value {
	v = &vArray{newValue(rv), rv.Len()}
	return v
}
func (v Array) Kind() r.Kind { return r.Array }
func (v Array) ArrayType() ArrayType {
	t, _ := TypeTo(v.Type()).ArrayType()
	return t
}

func (v Array) Cap() int { return v.len }
func (v Array) Len() int { return v.len }

func (v Array) Index(i int) (Value, error) {
	if i >= v.Len() {
		return nil, newErr("Array.Index", ErrOutOfRange)
	}
	return valueFrom(v.v.Index(i)), nil
}
