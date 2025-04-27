package helper_test

import (
	. "gostd/reflect/helper"
	"testing"
)

func TestSliceType(t *testing.T) {
	testSliceType[int]()
	testSliceType[struct{ V int }]()
	testSliceType[[5]string]()

	if s, _ := SliceOf(nil); s != nil {
		t.Fatal("SliceOf(nil) is not return nil")
	}

	log(TypeSpecifyOf[SliceType]([]int{}))       // []int
	log(TypeSpecifyOf[SliceType](nil))           // <nil>
	log(TypeSpecifyOf[SliceType](10))            // <nil>
	log(TypeSpecifyOf[SliceType]([][][][]int{})) // [][][][]int
}

func testSliceType[T any]() {
	st := SliceFor[T]()
	testTypeCommon(st)
	logf("Elem: %s, Kind: %s", st.Elem().String(), st.Elem().Kind())

}

func TestSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	slice2 := []int{1, 2, 3, 4}
	ps := ValueOf(&slice).(Pointer).Elem().(Slice)

	s := ValueOf(slice2).(Slice)
	log(ps, s)

	pe, _ := ps.Index(0)
	se, _ := s.Index(2)
	pe.(ValueCommon).BaseSetter().Set(ValueOf(10))
	se.(ValueCommon).BaseSetter().Set(ValueOf(99))
	log(slice)
	log(slice2)

	log(ps.Value().Slice(0, 1), s.Value().Slice(0, 1), ps, s)
}
