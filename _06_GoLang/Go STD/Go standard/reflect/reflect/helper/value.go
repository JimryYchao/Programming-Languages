package helper

import (
	"fmt"
	r "reflect"
	"unsafe"
)

// ! >>>>>>>>>>>> Value <<<<<<<<<<<<

type Value interface {
	valueof(r.Value) Value
	IsValid() bool
	Kind() r.Kind
	Type() Type
	Value() r.Value
	String() string
	IsExported() bool
}

type ValueCommon interface {
	Value
	Comparable() bool
	IsZero() bool
	CanSet() bool
	CanConvert(u r.Type) bool
	Convert(u r.Type) (Value, bool)

	BaseSetter() ValueSetter
	// Equal() bool
}

type vCommon struct {
	v          r.Value
	isExported bool
}

func (v *vCommon) valueof(x r.Value) Value { // 派生实现
	rt := vCommon{x, true}
	v = &rt
	return v
}
func (v *vCommon) IsValid() bool    { return true }
func (v *vCommon) Kind() r.Kind     { return v.v.Kind() }
func (v *vCommon) Type() Type       { return typeFrom(v.v.Type()) }
func (v *vCommon) Value() r.Value   { return v.v }
func (v *vCommon) String() string   { return fmt.Sprint(v.v.Interface()) }
func (v *vCommon) IsExported() bool { return v.isExported }
func (v *vCommon) IsZero() bool     { return v.v.IsZero() }
func (v *vCommon) CanSet() bool     { return v.v.CanSet() }
func (v *vCommon) Comparable() bool { return v.v.Comparable() }
func (v *vCommon) CanConvert(u r.Type) bool {
	if u == nil {
		return false
	}
	return v.v.CanConvert(u)
}
func (v *vCommon) Convert(u r.Type) (Value, bool) {
	if v.CanConvert(u) {
		return valueFrom(v.v.Convert(u)), true
	}
	return nil, false
}
func (v *vCommon) BaseSetter() ValueSetter {
	if v.CanSet() {
		return struct {
			*vCommon
			*vSetterPitch
		}{v, &vSetterPitch{vSetter{&v.v}}}
	}
	return nil
}

// func (v *valueBase) ValueSetter() ValueSetter { return &vSetter{v} }

// func (v *valueBase) To() toValue           { return tovalue{valueFrom(v.v)} }

// todo
func newValue(v r.Value) *vCommon {
	return &vCommon{v, true}
}

func getValue[V Value](v r.Value) Value {
	return V.valueof(*(*V)(unsafe.Pointer(new(V))), v)
}

func valueFrom(v r.Value) Value {
	if v.Kind() == r.Invalid {
		return NilValue{v}
	}
	switch v.Kind() {
	case r.Slice:
		return getValue[Slice](v)
	case r.Array:
		return getValue[Array](v)
	case r.Func:
		return getValue[Func](v)
	case r.Map:
		return getValue[Map](v)
	case r.Chan:
		return getValue[Chan](v)
	case r.Pointer:
		return getValue[Pointer](v)
	case r.Struct:
		return getValue[Struct](v)
	case r.String:
		return getValue[String](v)
	case r.Interface:
		return getValue[Interface](v)
	case r.Int, r.Int8, r.Int16, r.Int32, r.Int64:
		return getValue[Int](v)
	// case r.Uint, r.Uint8, r.Uint16, r.Uint32, r.Uint64:
	// 	return getValue[Uint](v)
	default:
		return &vCommon{v, true}
	}
}

// 构造 Type 零值
func Zero(t Type) Value {
	if IsNilType(t) {
		return nil
	}
	return valueFrom(r.Zero(t.Type()))
}

func ValueOf(i any) Value {
	return valueFrom(r.ValueOf(i))
}

func ValueFrom(v r.Value) Value {
	return valueFrom(v)
}

func IsValue[V Value](v Value) bool {
	if v == nil {
		return false
	}
	return (*(*V)(unsafe.Pointer(new(V)))).Kind() == v.Kind()
}

func ValueSpecify[V Value](v r.Value) (V, bool) {
	return toV[V](valueFrom(v))
}

func ValueSpecifyOf[V Value](i any) (V, bool) {
	return toV[V](ValueOf(i))
}

// ! >>>>>>>>>>>> toValue <<<<<<<<<<<<
type tovalue struct {
	Value
}

func toV[V Value](v Value) (V, bool) {
	out, ok := v.(V)
	return out, ok
}

func (t tovalue) Array() (Array, bool)         { return toV[Array](t.Value) }
func (t tovalue) Chan() (Chan, bool)           { return toV[Chan](t.Value) }
func (t tovalue) Pointer() (Pointer, bool)     { return toV[Pointer](t.Value) }
func (t tovalue) Func() (Func, bool)           { return toV[Func](t.Value) }
func (t tovalue) Map() (Map, bool)             { return toV[Map](t.Value) }
func (t tovalue) Struct() (Struct, bool)       { return toV[Struct](t.Value) }
func (t tovalue) Slice() (Slice, bool)         { return toV[Slice](t.Value) }
func (t tovalue) String() (String, bool)       { return toV[String](t.Value) }
func (t tovalue) Interface() (Interface, bool) { return toV[Interface](t.Value) }

// ! >>>>>>>>>>>> ValueSetter <<<<<<<<<<<<
type ValueSetter interface {
	ValueCommon
	SetZero()
	Set(Value) error
	init(Pointer) ValueSetter
	kind() r.Kind
	SetIterKey(*MapIter) error
	SetIterValue(*MapIter) error
}

type vSetter struct {
	v *r.Value
}
type vSetterPitch struct {
	vSetter
}

func (v *vSetterPitch) kind() r.Kind             { return v.v.Kind() }
func (v *vSetterPitch) init(Pointer) ValueSetter { return nil }

func (v *vSetter) SetZero() {
	v.v.SetZero()
}
func (v *vSetter) Set(x Value) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s", e)
		}
	}()
	if x == nil {
		return ErrArgNil
	}
	v.v.Set(x.Value())
	return
}
func (v *vSetter) SetIterKey(iter *MapIter) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s", e)
		}
	}()
	v.v.SetIterKey(iter.iter)
	return nil
}
func (v *vSetter) SetIterValue(iter *MapIter) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s", e)
		}
	}()
	v.v.SetIterValue(iter.iter)
	return nil
}
func SetterFor[T ValueSetter](p Pointer) (T, bool) {
	var t T
	if !p.Elem().Value().CanSet() {
		return t, false
	}
	if p.Elem().Kind() == (*(*T)(unsafe.Pointer(new(T)))).kind() {
		t = t.init(p).(T)
		return t, true
	}
	return t, false
}

// // func (c *valueCom) IsZero() bool     { return c.v.IsZero() }

// func ValueCom(v Value) (ValueCommon, error) {
// 	// if !IsNilType()
// 	return nil, nil
// }

// func FromValue[V Value](v r.Value) V {
// 	// TODO
// 	return *(*V)(nil)
// }

//! >>>>>>>>>>>>>> Nil Value <<<<<<<<<<<<<<<<

type NilValue struct{ v r.Value }

func (n NilValue) valueof(r.Value) Value { return nil }
func (n NilValue) IsValid() bool         { return false }
func (n NilValue) Kind() r.Kind          { return r.Invalid }
func (n NilValue) Type() Type            { return nil }
func (n NilValue) String() string        { return "<nil value>" }
func (n NilValue) Value() r.Value        { return n.v }
func (n NilValue) IsExported() bool      { return false }
