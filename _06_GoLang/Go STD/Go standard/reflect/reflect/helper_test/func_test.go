package helper_test

import (
	. "gostd/reflect/helper"
	"reflect"
	"testing"
)

func TestFuncType(t *testing.T) {

	i := reflect.TypeFor[int]()
	s := reflect.TypeFor[string]()
	b := reflect.TypeFor[bool]()
	e := reflect.TypeFor[error]()
	f := reflect.TypeFor[float64]()
	slice := reflect.TypeFor[[]any]()

	testFuncType([]reflect.Type{}, []reflect.Type{b}, b)
	testFuncType([]reflect.Type{i, s}, []reflect.Type{}, f)
	testFuncType([]reflect.Type{f, s, slice}, nil, nil)
	testFuncType(nil, []reflect.Type{e}, SliceFor[int]().Type())
	testFuncType(nil, nil, nil)
}

func testFuncType(in []reflect.Type, out []reflect.Type, va reflect.Type) {

	ft, err := FuncOf(in, out, va)
	if err != nil {
		log(err)
		return
	}
	testTypeCommon(ft)
	logf("IsVariadic: %t", ft.IsVariadic())
	logf("NumIn: %d", ft.NumIn())
	logf("Ins: %s", ft.Ins())
	logf("NumOut: %d", ft.NumOut())
	logf("Outs: %s", ft.Outs())
}

// func fmtTypes(tps []Type) string {
// 	var sb strings.Builder
// 	for _, v := range tps {
// 		sb.WriteString(v.String() + ", ")
// 	}
// 	return
// }
