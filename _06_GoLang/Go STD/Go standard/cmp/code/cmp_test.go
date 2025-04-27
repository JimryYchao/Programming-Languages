package gostd

import (
	"cmp"
	"math"
	"testing"
	"unsafe"
)

//? go test -v -run=^$

/*
! Compare(x, y) x<y -> -1; x==y -> 0; x>y -> 1
! Less return x > y
! Or 返回它的第一个不等于零值的参数
*/

type rt interface {
	~int | ~bool
}

func TestCmp(t *testing.T) {
	for _, test := range tests {
		switch test.x.(type) {
		case int:
			T := struct {
				x int
				y int
			}{test.x.(int), test.y.(int)}
			f("Compare", T, cmp.Compare)
			f("   Less", T, cmp.Less)
		case string:
			T := struct {
				x string
				y string
			}{test.x.(string), test.y.(string)}
			f("Compare", T, cmp.Compare)
			f("   Less", T, cmp.Less)
		case float64:
			T := struct {
				x float64
				y float64
			}{test.x.(float64), test.y.(float64)}
			f("Compare", T, cmp.Compare)
			f("   Less", T, cmp.Less)
		case uintptr:
			T := struct {
				x uintptr
				y uintptr
			}{test.x.(uintptr), test.y.(uintptr)}
			f("Compare", T, cmp.Compare)
			f("   Less", T, cmp.Less)
		}
	}
}

var negzero = math.Copysign(0, -1)
var nonnilptr uintptr = uintptr(unsafe.Pointer(&negzero))
var nilptr uintptr = uintptr(unsafe.Pointer(nil))

func f[T cmp.Ordered, R rt](_case string, t struct{ x, y T }, fn func(x, y T) R) {
	logfln("%s(%v, %v) = %v", _case, t.x, t.y, fn(t.x, t.y))
}

var tests = []struct {
	x, y any
}{
	{1, 2},
	{1, 1},
	{2, 1},
	{"a", "aa"},
	{"a", "a"},
	{"aa", "a"},
	{1.0, 1.1},
	{1.1, 1.1},
	{1.1, 1.0},
	{math.Inf(1), math.Inf(1)},
	{math.Inf(-1), math.Inf(-1)},
	{math.Inf(-1), 1.0},
	{1.0, math.Inf(-1)},
	{math.Inf(1), 1.0},
	{1.0, math.Inf(1)},
	{math.NaN(), math.NaN()},
	{0.0, math.NaN()},
	{math.NaN(), 0.0},
	{math.NaN(), math.Inf(-1)},
	{math.Inf(-1), math.NaN()},
	{0.0, 0.0},
	{negzero, negzero},
	{negzero, 0.0},
	{0.0, negzero},
	{negzero, 1.0},
	{negzero, -1.0},
	{nilptr, nonnilptr},
	{nonnilptr, nilptr},
	{nonnilptr, nonnilptr},
}
