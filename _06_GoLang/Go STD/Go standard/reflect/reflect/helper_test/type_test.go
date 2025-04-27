package helper_test

import (
	"fmt"
	. "gostd/reflect/helper"
	"reflect"
	"testing"
)

type mInt int
type MInt int

var log = fmt.Println

func logf(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func TestGetType(tt *testing.T) {
	var tp = reflect.TypeFor[[]int]()

	if t := TypeFrom(tp); IsType[SliceType](t) {
		fmt.Println(t.String())
	}

	if n := TypeOf(nil); IsType[Nil](n) {
		fmt.Println(n.String())
	}

	if t := TypeFor[[]int](); IsType[SliceType](t) {
		fmt.Println(t.String())
	}
}

func TestTypeOf(t *testing.T) {
	type IntInline int
	if st := TypeFrom(reflect.TypeOf(any([]MInt{}))); IsType[SliceType](st) {
		testTypeCommon(st)
	}
	testTypeCommon(TypeOf(mInt(0)))

	testTypeCommon(TypeOf([]MInt{}).(SliceType))
	testTypeCommon(TypeFrom(reflect.TypeOf(int(0))))
	testTypeCommon(TypeFrom(reflect.TypeOf(IntInline(0))))
	testTypeCommon(TypeFor[struct{ anom int }]())
	testTypeCommon(TypeFor[[]struct{ anom int }]().(SliceType))
	testTypeCommon(TypeOf(nil))
}

func testTypeCommon(t Type) {
	if t == nil {
		log("t is nil")
	}
	logf("\n>>>>>  t : %s  <<<<<", t.Name())
	logf("Type : %s", t.Type())
	logf("String : %s", t.String())
	logf("Kind : %s", t.Kind())

	if !IsNilType(t) {
		if t, ok := t.(TypeCommon); ok {
			logf("Size : %d", t.Size())
			logf("Align : %d", t.Align())
			logf("PkgPath : %s", t.PkgPath())
			logf("AssignableTo any : %t", t.AssignableTo(reflect.TypeFor[any]()))
			logf("ConvertibleTo any : %t", t.ConvertibleTo(reflect.TypeFor[any]()))
			logf("Comparable : %t", t.Comparable())
			logf("FieldAlign : %d", t.FieldAlign())
		}
	}
}
