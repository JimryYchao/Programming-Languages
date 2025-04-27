package helper

import (
	r "reflect"
)

//! >>>>>>>>>>>>>> Pointer <<<<<<<<<<<<<<

type Pointer = *vPointer
type PointerToArray = *vPointerToArray

type vPointer struct {
	*vCommon
}

func (v Pointer) valueof(rv r.Value) Value {
	v = &vPointer{newValue(rv)}
	return v
}
func (v Pointer) Kind() r.Kind { return r.Pointer }
func (v Pointer) PointerType() PointerType {
	t, _ := TypeTo(v.Type()).PointerType()
	return t
}

func (v Pointer) Elem() Value    { return valueFrom(v.v.Elem()) }
func (v Pointer) Interface() any { return v.v.Interface() }

// func (v Pointer)

// ! >>>>>>>>>>>>>> Pointer Elems <<<<<<<<<<<<<<
type vPointerToArray struct {
	Pointer
}

func (v Pointer) ToArrayPointer() PointerToArray {
	if v.Elem().Kind() == r.Array {
		return &vPointerToArray{v}
	}
	return nil
}
func (p PointerToArray) Cap() int { return p.v.Cap() }
func (p PointerToArray) Len() int { return p.v.Len() }
