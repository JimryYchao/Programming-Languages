package gostd

import (
	"testing"
	"unicode/utf16"
)

/*
encode functions:
! AppendRune 将 Unicode 码位 r 的 UTF-16 编码附加到 []uint16 p 的末尾
! Encode 返回 Unicode 码位编码的 UTF-16 编码
! EncodeRune 返回代理项对的 UTF-16 编码
! IsSurrogate 报告指定的 Unicode 码位是否可以出现在代理项对中。
decode functions:
! Decode 返回由 UTF-16 编码表示的 Unicode 代码点序列
! DecodeRune 返回代理项对 r1, r2 的 UTF-16 解码
*/
//? go test -v -run=^$
func Test(t *testing.T) {
	content := []rune("Hello, 世界")
	var u16 []uint16 = make([]uint16, len(content)*2)
	n := 0

	for _, r := range content {
		if utf16.IsSurrogate(r) {
			r1, r2 := utf16.EncodeRune(r)
			u16[n] = uint16(r1)
			u16[n+1] = uint16(r2)
			n += 2
		} else {
			u16[n] = uint16(r)
			n++
		}
	}
	logfln("%s", string(utf16.Decode(u16)))

	// 等效于
	u16 = utf16.Encode(content)
	logfln("%s", string(utf16.Decode(u16)))
}
