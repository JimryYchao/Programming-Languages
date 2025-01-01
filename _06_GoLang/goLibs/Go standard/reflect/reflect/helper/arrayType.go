package helper

import (
	r "reflect"
)

//! >>>>>>>>>>>> ArrayType <<<<<<<<<<<<

type ArrayType = *arrayType

type arrayType struct {
	*tCommon
	len int
}

func (t ArrayType) typeof(tp r.Type) Type {
	t = &arrayType{newTCommon(tp), tp.Len()}
	return t
}
func (ArrayType) Kind() r.Kind { return r.Array }
func (t ArrayType) Elem() Type { return typeFrom(t.t.Elem()) }
func (t ArrayType) Len() int   { return t.len }

// ArrayOf
func ArrayOf(length int, tp r.Type) (ArrayType, error) {
	if tp == nil {
		return nil, newErr("ArrayOf type", ErrArgNil)
	}
	if length < 0 {
		return nil, newErr("ArrayOf length", ErrNegative)
	}
	return &arrayType{newTCommon(r.ArrayOf(int(length), tp)), int(length)}, nil
}

func ArrayFor[T any](length int) ArrayType {
	a, _ := ArrayOf(length, r.TypeFor[T]())
	return a
}
