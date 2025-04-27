package helper_test

import (
	. "gostd/reflect/helper"
	"reflect"
	"testing"
)

func TestMapType(t *testing.T) {

	testMapType(MapFor[string, int]())
	testMapType(MapFor[string, []int]())
	testMapType(MapFor[[5]int, []int]())
	testMapType(MapFor[any, []int]())

	if t, ok := TypeSpecifyOf[MapType](map[any]int{}); ok {
		testMapType(t)
	}
	log(TypeSpecifyOf[MapType](10086))
	log(MapOf(reflect.TypeFor[any](), reflect.TypeFor[any]()))
	log(MapOf(reflect.TypeFor[[]int](), reflect.TypeFor[any]()))
	log(MapOf(reflect.FuncOf(nil, nil, false), reflect.TypeFor[any]()))
	// m := make(map[any]int)
}

func testMapType(m MapType) {
	if m == nil {
		log("m is nil")
		return
	}
	testTypeCommon(m)
	logf("Key:%s, Elem:%s", m.Key(), m.Elem())
}

func TestMap(t *testing.T) {
	m := map[int]int{
		1: 1, 2: 2, 3: 3, 4: 4,
	}
	p := ValueOf(&m).(Pointer)
	setter := p.Elem().(Map)
	log(setter.MapType().Elem())
	// MapKeys
	for _, v := range setter.MapKeys() {
		r := v
		log(r)
	}
	iter := setter.MapRange()

	// MapSetter
	mr, _ := SetterFor[MapSetter](p)
	mr.Value().SetMapIndex(reflect.ValueOf(5), reflect.ValueOf(5))
	log(mr.SetMapIndex(ValueOf(1.1), ValueOf(6)))
	log(mr.SetMapIndex(ValueOf(1), nil))
	for iter.Next() {
		log(iter.Value())
	}
}
