package helper_test

import (
	. "gostd/reflect/helper"
	"testing"
)

func TestPointerType(t *testing.T) {
	testPointerType[int]()
	testSliceType[struct{ V int }]()
	testPointerType[[5]string]()

	if s, _ := PointerTo(nil); s != nil {
		t.Fatal("SliceOf(nil) is not return nil")
	}

	log(TypeSpecifyOf[PointerType]([]int{}))          // []int
	log(TypeSpecifyOf[PointerType](nil))              // <nil>
	log(TypeSpecifyOf[PointerType](new(PointerType))) // <nil>
	log(TypeSpecifyOf[PointerType](new(int)))         //
}

func testPointerType[T any]() {
	st := PointerFor[T]()
	testTypeCommon(st)
	logf("Elem: %s, Kind: %s", st.Elem().String(), st.Elem().Kind())
}
