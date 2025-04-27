package helper

import (
	r "reflect"
)

//! >>>>>>>>>>>> FuncType <<<<<<<<<<<<

type FuncType = *funcType

type funcType struct {
	*tCommon
	in         uint
	out        uint
	isVariadic bool
}

func (t FuncType) typeof(tp r.Type) Type {
	t = &funcType{newTCommon(tp), (uint)(tp.NumIn()), (uint)(tp.NumOut()), tp.IsVariadic()}
	return t
}

func (FuncType) Kind() r.Kind       { return r.Func }
func (t FuncType) IsVariadic() bool { return t.isVariadic }
func (t FuncType) NumIn() int       { return int(t.in) }
func (t FuncType) NumOut() int      { return int(t.out) }
func (t FuncType) In(i uint) (Type, bool) {
	if i < t.in && t.in > 0 {
		return typeFrom(t.t.In(int(i))), true
	}
	return nil, false
}
func (t FuncType) Out(i uint) (Type, bool) {
	if i < t.out && t.out > 0 {
		return typeFrom(t.t.In(int(i))), true
	}
	return nil, false
}

// 补充
func (t FuncType) TryIn(i uint, out *Type) bool {
	o, ok := t.In(i)
	*out = o
	return ok
}
func (t FuncType) TryOut(i uint, out *Type) bool {
	o, ok := t.Out(i)
	*out = o
	return ok
}

func (t FuncType) Ins() []Type {
	tps := make([]Type, t.in)
	for i := range t.in {
		tps[i], _ = t.In(i)
	}
	return tps
}

func (t FuncType) Outs() []Type {
	tps := make([]Type, t.out)
	for i := range t.out {
		tps[i], _ = t.Out(i)
	}
	return tps
}

// FuncOf
// va == nil 时，表示无可变参
func FuncOf(in []r.Type, out []r.Type, va r.Type) (FuncType, error) {
	if n := len(in) + len(out); n > 128 {
		return nil, newErr("FuncOf", ErrTooManyArgs)
	}
	if va == nil {
		return &funcType{newTCommon(r.FuncOf(in, out, false)), uint(len(in)), uint(len(out)), false}, nil
	} else {
		if va.Kind() != r.Slice {
			return nil, newErr("FuncOf", ErrVaNotSlice)
		}
		in = append(in, va)
		return &funcType{newTCommon(r.FuncOf(in, out, true)), uint(len(in)), uint(len(out)), true}, nil
	}
}

// New Func

// func (t *FuncType) New() {

// }
