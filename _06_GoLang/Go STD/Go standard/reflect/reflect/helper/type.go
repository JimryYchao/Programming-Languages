package helper

import (
	"errors"
	"fmt"
	r "reflect"
	"unsafe"
)

// ! >>>>>>>>>>>>>> Type <<<<<<<<<<<<<<<<
type Type interface {
	typeof(r.Type) Type
	Kind() r.Kind
	Type() r.Type
	Name() string
	String() string
	// common() TypeCommon
}
type tCommon struct {
	t r.Type
}

func (b *tCommon) typeof(r.Type) Type { return b }
func (b *tCommon) Kind() r.Kind       { return b.t.Kind() }
func (b *tCommon) Type() r.Type       { return b.t }
func (b *tCommon) Name() string       { return b.t.Name() }
func (b *tCommon) String() string     { return b.t.String() }
func (c *tCommon) Size() uintptr      { return c.t.Size() }
func (c *tCommon) Align() int         { return c.t.Align() }
func (c *tCommon) PkgPath() string    { return c.t.PkgPath() }
func (c *tCommon) Comparable() bool   { return c.t.Comparable() }
func (c *tCommon) FieldAlign() int    { return c.t.FieldAlign() }
func (c *tCommon) AssignableTo(u r.Type) bool {
	if u == nil {
		return false
	}
	return c.t.AssignableTo(u)
}
func (c *tCommon) ConvertibleTo(u r.Type) bool {
	if u == nil {
		return false
	}
	return c.t.ConvertibleTo(u)
}

func newTCommon(tp r.Type) *tCommon {
	return &tCommon{tp}
}

func getType[T Type](tp r.Type) Type {
	return T.typeof(*(*T)(unsafe.Pointer(new(T))), tp)
}

func typeFrom(tp r.Type) Type {
	if tp == nil || tp.Kind() == r.Invalid {
		return Nil{}
	}
	switch tp.Kind() {
	case r.Slice:
		return getType[SliceType](tp)
	case r.Map:
		return getType[MapType](tp)
	case r.Array:
		return getType[ArrayType](tp)
	case r.Func:
		return getType[FuncType](tp)
	case r.Chan:
		return getType[ChanType](tp)
	case r.Struct:
		return getType[StructType](tp)
	case r.Pointer:
		return getType[PointerType](tp)
	case r.Interface:
		return getType[InterfaceType](tp)
	default:
		return newTCommon(tp) // 返回常规 reflect.Type
	}
}

// 从类型构造一个 Type
func TypeFor[T any]() Type {
	return typeFrom(r.TypeOf((*T)(nil)).Elem())
}

// 从 v 提取一个 Type
func TypeOf(i any) Type {
	return typeFrom(r.TypeOf(i))
}

// 包装一个 reflect.Type
func TypeFrom(tp r.Type) Type {
	return typeFrom(tp)
}

// 检查包装类型
func IsType[T Type](t Type) bool {
	if t == nil {
		return false
	}
	_, ok := t.(T)
	return ok
}

// 尝试包装为特定的 Type
func TypeSpecify[T Type](tp r.Type) (T, bool) {
	return toT[T](typeFrom(tp))
}

func TypeSpecifyOf[T Type](i any) (T, bool) {
	return toT[T](TypeOf(i))
}

//! >>>>>>>>>>>>>> totype <<<<<<<<<<<<<<<<

type totype struct {
	t Type
}

func (t totype) ArrayType() (ArrayType, bool)         { return toT[ArrayType](t.t) }
func (t totype) ChanType() (ChanType, bool)           { return toT[ChanType](t.t) }
func (t totype) PointerType() (PointerType, bool)     { return toT[PointerType](t.t) }
func (t totype) FuncType() (FuncType, bool)           { return toT[FuncType](t.t) }
func (t totype) MapType() (MapType, bool)             { return toT[MapType](t.t) }
func (t totype) StructType() (StructType, bool)       { return toT[StructType](t.t) }
func (t totype) SliceType() (SliceType, bool)         { return toT[SliceType](t.t) }
func (t totype) InterfaceType() (InterfaceType, bool) { return toT[InterfaceType](t.t) }

func toT[T Type](t Type) (T, bool) {
	out, ok := t.(T)
	return out, ok
}

func TypeTo(t Type) totype {
	if t == nil || t.Kind() == r.Invalid {
		return totype{Nil{}}
	}
	return totype{t}
}

// ! >>>>>>>>>>>>>> TypeCommon <<<<<<<<<<<<<<<<
type tCommonBase interface {
	Type
	Size() uintptr
	Align() int
	PkgPath() string
	AssignableTo(r.Type) bool
	ConvertibleTo(r.Type) bool
	Comparable() bool
	FieldAlign() int
}
type TypeCommon interface {
	tCommonBase
	Implements(u r.Type) bool
}

type typeBase1 struct {
	*tCommon
}

func (c *tCommon) Implements(u r.Type) bool { return c.t.Implements(u) }

//	func TypeCom(c Type) (typeCommon, error) {
//		if !IsNilType(c) {
//			return &typeCom{&typeBase{c.Type()}}, nil
//		}
//		return nil, newErr("TypeCom", ErrArgNil)
//	}

//! >>>>>>>>>>>>>> TypeProperty <<<<<<<<<<<<<<<<

type typePropConstraint interface {
	TypeCommon
	Type
}
type typeProperty interface {
	IsDefined() bool
	IsBuildIn() bool
	IsAnonymous() bool
}

type typeProper struct {
	com typePropConstraint
}

func PropFor(t typePropConstraint) typeProperty {
	return typeProper{t}
}

func (c typeProper) IsDefined() bool   { return (c.com).Name() != "" }
func (c typeProper) IsBuildIn() bool   { return (c.com).Name() != "" && (c.com).PkgPath() == "" }
func (c typeProper) IsAnonymous() bool { return (c.com).Name() == "" && (c.com).PkgPath() == "" }

//! >>>>>>>>>>>>>> Nil <<<<<<<<<<<<<<<<

type Nil struct{}

func (n Nil) typeof(r.Type) Type { return n }
func (n Nil) Type() r.Type       { return nil }
func (n Nil) Kind() r.Kind       { return r.Invalid }
func (n Nil) String() string     { return "<nil>" }
func (n Nil) Name() string       { return "nil" }
func (n Nil) common() TypeCommon { return nil }
func IsNilType(t Type) bool {
	_, ok := t.(Nil)
	return ok
}

//! >>>>>>>>>>>>>> Err <<<<<<<<<<<<<<<<

type HelperErr struct {
	op  string
	err error
}

func (e *HelperErr) Error() string {
	return e.op + ": " + e.err.Error()
}

func (e *HelperErr) UnWrap() error {
	return e.err
}

func newErr(op string, err error) error {
	return &HelperErr{op, err}
}
func newErrorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

var ErrArgNil = errors.New("arg is nil")
var ErrOutOfRange = errors.New("index is out of range")
var ErrChanElemSize = errors.New("element size too large")
var ErrNegative = errors.New("arg is not negative")
var ErrTooManyArgs = errors.New("too many args")
var ErrVaNotSlice = errors.New("last arg of variadic func must be slice")

// var ErrTypeNil = newErr("reflect.Type is nil")
