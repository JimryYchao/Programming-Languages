package helper_test

import (
	. "gostd/reflect/helper"
	"testing"
)

func TestArrayType(t *testing.T) {
	testArrayType[int](10)
	testArrayType[struct{ V int }](0)
	testArrayType[[5]string](15)

	if a, _ := ArrayOf(10, nil); a != nil {
		t.Fatal("ArrayOf(nil) is not return nil")
	}

	log(TypeSpecifyOf[ArrayType]([20]int{}))       // []int
	log(TypeSpecifyOf[ArrayType](nil))             // <nil>
	log(TypeSpecifyOf[ArrayType](10))              // <nil>
	log(TypeSpecifyOf[ArrayType]([10][][][]int{})) // [10][][][]int
}

func testArrayType[T any](len int) {
	at := ArrayFor[T](len)
	testTypeCommon(at)
	logf("Len: %d, Elem: %s", at.Len(), at.Elem().String())
}
