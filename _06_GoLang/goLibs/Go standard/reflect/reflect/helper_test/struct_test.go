package helper_test

import (
	. "gostd/reflect/helper"
	"testing"
)

func TestStruct(t *testing.T) {
	type s1 struct {
		v1 int `tag:"????"`
		v2 int `tag:"????" tag2:"tag2"`
		E1 int `mess tag:"????" tag2:"tag2"`
		E2 int `tag:'????'`
		int
	}

	log(TypeTo(TypeOf([]int{})).SliceType())
	if s, ok := TypeTo(TypeOf(s1{})).StructType(); ok {
		fields := VisibleFields(s)
		for _, f := range fields {
			logf("name:%s, tag:%s, tag2:%s", f.Name, f.Tag.Get("tag"), f.Tag.Get("tag2"))
			logf(f.Get("tag"))
		}
	}

	type s2 struct {
		s1
		E1 int
		v1 string
	}
	type s3 struct {
		E1 string
		s2
	}

	testStruct(TypeSpecifyOf[StructType](s1{}))
	testStruct(TypeSpecifyOf[StructType](s2{}))
	testStruct(TypeSpecifyOf[StructType](s3{}))

	if f, ok := TypeSpecifyOf[StructType](s3{}); ok {
		log(f.FieldByIndex([]int{1, 2}))
	}
}

func testStruct(s StructType, ok bool) {
	if s == nil || !ok {
		log("struct is a nil")
		return
	}
	// testTypeCommon(s)
	logf("num: %d", s.NumField())
	fs := VisibleFields(s)
	for _, f := range fs {
		logf("field: %s, Kind:%s, Index:%v", f.Name, f.Type().Kind(), f.Index)
	}
	iterMethods(MethodOf(s))
}
