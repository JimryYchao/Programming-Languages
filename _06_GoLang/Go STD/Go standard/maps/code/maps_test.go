package gostd

import (
	"maps"
	"strings"
	"testing"
)

var one, ten, thousand string = "one", "ten", "thousand"
var One, Ten, Thousand string = "One", "Ten", "Thousand"

var m1 = map[int](*string){
	1:    &one,
	10:   &ten,
	1000: &thousand,
}
var m2 = map[int]*string{
	1:    &One,
	10:   &Ten,
	1000: &Thousand,
}
var m3 = map[int]string{
	1:    One,
	10:   Ten,
	1000: Thousand,
}

// ! Clone 返回 m 的副本，这是一个浅克隆：键和值是使用普通赋值设置的
// ! DeleteFunc 从 m 中删除任何满足给定函数 del 的键/值对
// ! Equal 报告两个映射是否包含相同的键/值对。使用 == 比较值。
// ! EqualFunc 函数类似于 Equal，但使用给定函数 eq 比较值。键仍然以 == 进行比较
func TestEqual(t *testing.T) {
	m := maps.Clone(m1)
	log(maps.Equal(m1, m2)) // false
	log(maps.Equal(m1, m))  // true
	log(maps.EqualFunc(m1, m2, func(v1 *string, v2 *string) bool {
		return strings.EqualFold(*v1, *v2)
	})) // true
	log(maps.EqualFunc(m2, m3, func(v1 *string, v2 string) bool {
		return *v1 == v2
	})) // true

	maps.DeleteFunc(m1, func(k int, v *string) bool { return k < 10 })
	for k, v := range m1 {
		logfln("k,v = %d, %s", k, *v)
	}
	log(maps.EqualFunc(m1, m2, func(v1 *string, v2 *string) bool {
		return strings.EqualFold(*v1, *v2)
	})) // false
}
