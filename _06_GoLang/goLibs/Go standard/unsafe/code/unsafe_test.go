package gostd

import (
	"bytes"
	"io"
	"os"
	"testing"
	. "unsafe"
)

/*
! Alignof 返回类型的对齐方式, 和 reflect.TypeOf(x).Align() 值相同；
	若 s 是结构类型，则字段 f 的 Alignof(s.f) 返回结构中字段的所需对齐方式，与 reflect.TypeOf(s.f).FieldAlign() 相同
! Offsetof 返回由 x 表示的字段 f 在结构内的偏移量; Offsetof(x.f)
! Sizeof 接受类型 x 的大小（以字节为单位）; 若 x 是切片, 字符串, 则返回其描述符的大小而非引用的内存大小; map 返回指针大小
! String 返回一个字符串值，其基础字节从 ptr 开始，长度为 len;
	由于字符串是不可变的，因此传递给 String 的字节不能在之后被修改
! Slice 返回一个切片，其底层数组以 ptr 开始，长度和容量为 len。Slice(ptr,len) 等价于
	(*[len]ArbitraryType)(unsafe.Pointer(ptr))[:]
! SliceData 返回指向参数 slice 的底层数组的指针

! unsafe.Pointer 表示指向任意类型的指针
	有四种特殊操作:
		任何类型的指针值都可以转换 Pointer
		Pointer 可以转换为任何类型的指针值
		uintptr 类型可以转换为 Pointer 类型
		Pointer 可以转换为 uintptr 类型
! Add 将 len 添加到 ptr 并返回更新后的指针 `Pointer(uintptr(ptr) + uintptr(len))`
*/

func TestAlignof(t *testing.T) {
	logfln("bool:%d, uintptr:%d, ptr:%d, byte:%d, rune:%d, any:%d, interface:%d, func:%d, A1:%d, A2:%d, A3:%d",
		Alignof(true), Alignof(uintptr(0)), Alignof(&c), Alignof(byte(0)), Alignof(rune(0)), Alignof(any(0)), Alignof(io.Reader(&bytes.Reader{})), Alignof(func() {}),
		Alignof(A1{}), Alignof(A2{}), Alignof(A3{}))
	logfln("int8:%d, int16:%d, int32:%d, int64:%d, float32:%d, float64:%d, complex64:%d, complex128:%d",
		Alignof(int8(0)), Alignof(int16(0)), Alignof(int32(0)), Alignof(int64(0)), Alignof(float32(0)),
		Alignof(float64(0)), Alignof(complex64(0)), Alignof(complex128(0)))
	logfln("slice:%d, array:%d, map:%d, chan:%d, string:%d, S:%d", Alignof([]byte{}), Alignof([4]int{}),
		Alignof(map[int]A3{}), Alignof(make(chan int, 10)), Alignof("Hello"), Alignof(S{}))

	// bool:1, uintptr:8, ptr:8, byte:1, rune:4, any:8, interface:8, func:8, A1:2, A2:8, A3:8
	// int8:1, int16:2, int32:4, int64:8, float32:4, float64:8, complex64:4, complex128:8
	// slice:8, array:8, map:8, chan:8, string:8, S:8
}

func TestOffsetof(t *testing.T) {
	s := S{}
	a1, a2, a3 := A1{}, A2{}, A3{}
	logfln("offset: V1:%d, V2:%d, V3:%d", Offsetof(s.V1), Offsetof(s.V2), Offsetof(s.V3))
	logfln("offset: V1:%d, V2:%d", Offsetof(a1.V1), Offsetof(a1.V2))
	logfln("offset: V1:%d, V2:%d, V3:%d", Offsetof(a2.V1), Offsetof(a2.V2), Offsetof(a2.V3))
	logfln("offset: V1:%d, V2:%d, V3:%d", Offsetof(a3.V1), Offsetof(a3.V2), Offsetof(a3.V3))

	// offset: V1:0, V2:8, V3:32
	// offset: V1:0, V2:2
	// offset: V1:0, V2:8, V3:16
	// offset: V1:0, V2:4, V3:8
}

func TestSizeof(t *testing.T) {
	logfln("bool:%d, uintptr:%d, ptr:%d, byte:%d, rune:%d, any:%d, interface:%d, func:%d, A1:%d, A2:%d, A3:%d",
		Sizeof(true), Sizeof(uintptr(0)), Sizeof(&c), Sizeof(byte(0)), Sizeof(rune(0)), Sizeof(any(0)), Sizeof(io.Reader(&bytes.Reader{})), Alignof(func() {}),
		Sizeof(A1{}), Sizeof(A2{}), Sizeof(A3{}))
	logfln("int8:%d, int16:%d, int32:%d, int64:%d, float32:%d, float64:%d, complex64:%d, complex128:%d",
		Sizeof(int8(0)), Sizeof(int16(0)), Sizeof(int32(0)), Sizeof(int64(0)), Sizeof(float32(0)),
		Sizeof(float64(0)), Sizeof(complex64(0)), Sizeof(complex128(0)))
	logfln("slice[0:4]:%d, [4]array:%d, map[int]A3{}:%d, chan 10:%d, string:%d, S:%d", Sizeof([]byte{0, 1, 2, 3, 4}), Sizeof([4]int{}),
		Sizeof(map[int]A3{1: A3{}, 2: A3{}}), Sizeof(make(chan int, 10)), Sizeof("Hello World"), Sizeof(S{}))

	// bool:1, uintptr:8, ptr:8, byte:1, rune:4, any:16, interface:16, func:8, A1:4, A2:24, A3:16
	// int8:1, int16:2, int32:4, int64:8, float32:4, float64:8, complex64:8, complex128:16
	// slice[0:4]:24, [4]array:32, map[int]A3{}:8, chan 10:8, string:16, S:48
}

var s = S{V1: int(rune('A')), V2: [16]byte([]byte("Hello World!!!!!!!!")), V3: nil}

func TestString(t *testing.T) {
	r := String((*byte)(Pointer(&s)), 8+16+8)
	logfln(`"%s"`, r)
}

func TestSlice(t *testing.T) {
	slice := Slice((*byte)(Pointer(&s)), Sizeof(s))

	s1 := (*S)(Pointer(SliceData(slice)))
	logfln("%#v", slice)
	for i := range len(slice) - int(Sizeof(s1.V3)) {
		slice[i] = slice[0]
	}
	logfln("%s", slice)
	logfln("%v, %v, %v", s1.V1, s1.V2, s1.V3)

	pr := (*io.Writer)(Add(Pointer(s1), Offsetof(s1.V3)))

	*pr = os.Stdout

	s1.V3.Write([]byte("Hello World\n"))

	logfln("%#v", slice)
}

var c int

type S struct {
	V1 int
	_  [8]byte // 填充字段
	V2 [16]byte
	V3 io.Writer
}

func (s S) Func() {}

type (
	A1 struct {
		V1 int16
		V2 bool
	}
	A2 struct {
		V1 bool
		V2 int64
		V3 int32
	}
	A3 struct {
		V1 bool
		V2 int32
		V3 int64
	}
)
