package helper_test

import (
	. "gostd/reflect/helper"
	"testing"
	"unsafe"
)

func TestMethod(t *testing.T) {
	iterMethods(MethodOf(TypeFor[ArrayType]()))
	iterMethods(MethodOf(TypeFor[*int]()))
	iterMethods(MethodOf(TypeFor[Type]()))
	iterMethods(MethodOf(TypeFor[*Type]()))
	iterMethods(MethodOf(TypeFor[unsafe.Pointer]()))
	iterMethods(MethodOf(TypeOf(unsafe.Pointer(new(int)))))
	iterMethods(MethodOf(TypeFor[ArrayType]()))
	iterMethods(MethodOf(TypeFor[*ArrayType]()))
}

func testMethodSet(set MethodSet) {
	iterMethods(set)
}

func iterMethods(set MethodSet) {
	log(set.Receiver(), set.NumMethod(), set.Receiver().Kind())
	for _, m := range set.Methods() {
		logf("name:%s, type:%s", m.Name, m.Type)
	}
}
