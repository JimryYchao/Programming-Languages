package helper

import (
	r "reflect"
)

//! >>>>>>>>>>>> PointerType <<<<<<<<<<<<

type PointerType = *pointerType

type pointerType struct {
	*tCommon
}

func (t PointerType) typeof(tp r.Type) Type {
	t = &pointerType{newTCommon(tp)}
	return t
}

func (PointerType) Kind() r.Kind { return r.Pointer }
func (t PointerType) Elem() Type { return typeFrom(t.t.Elem()) }

// PointerTo

func PointerTo(tp r.Type) (PointerType, error) {
	if tp == nil {
		return nil, newErr("PointerTo type", ErrArgNil)
	}
	return &pointerType{newTCommon(r.PointerTo(tp))}, nil
}

func PointerFor[T any]() PointerType {
	s, _ := PointerTo(r.TypeFor[T]())
	return s
}
