package helper

import (
	r "reflect"
)

//! >>>>>>>>>>>> SliceType <<<<<<<<<<<<

type SliceType = *sliceType

type sliceType struct {
	*tCommon
}

func (t SliceType) typeof(tp r.Type) Type {
	t = &sliceType{newTCommon(tp)}
	return t
}

func (SliceType) Kind() r.Kind { return r.Slice }
func (t SliceType) Elem() Type { return typeFrom(t.t.Elem()) }

// SliceOf
func SliceOf(tp r.Type) (SliceType, error) {
	if tp == nil {
		return nil, newErr("SliceOf type", ErrArgNil)
	}
	return &sliceType{newTCommon(r.SliceOf(tp))}, nil
}

func SliceFor[T any]() SliceType {
	s, _ := SliceOf(r.TypeFor[T]())
	return s
}

// New Slice
// func (t SliceType) new(len, cap int) (Slice, error) {
// 	if len < 0 {
// 		return nil, newErr("len is a negative number")
// 	}
// 	if cap < len {
// 		return nil, newErr("cap is less than len")
// 	}

// 	slice := r.MakeSlice(t.t, len, cap)
// 	return &vSlice{&slice, nil}, nil
// }

// func (t SliceType) New(len int) (Slice, error) {
// 	return t.NewC(len, len)
// }
// func (t SliceType) NewC(len, cap int) (Slice, error) {
// 	if t != nil {
// 		return t.new(len, cap)
// 	}
// 	return nil, newErr("SliceType is invalid")
// }
