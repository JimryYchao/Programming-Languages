package gostd

import (
	"math"
	"testing"
	. "unicode/utf8"
)

/*
! AppendRune 将 r 的 UTF-8 编码附加到 p 的末尾，并返回扩展缓冲区。如果 rune 超出范围，它会附加 RuneError(U+FFFD) 的编码
! FullRune, FullRuneInString 报告 p 是否以 rune 的完整 UTF-8 编码开始
! DecodeRune, DecodeRuneInString 解包 p 中的第一个 UTF-8 编码, 并返回 rune 及其字节宽度
! DecodeLastRune, DecodeLastRuneInString 解包 p 中的最后一个 UTF-8 编码，并返回 rune 及其字节宽度
! EncodeRune 将 r 的 UTF-8 编码写入 p（必须足够大，>= 3）
! RuneCount, RuneCountInString 返回 p 中的 rune 数。错误和短编码被视为宽度为 1 字节的单个 rune
! RuneLen 返回编码 rune 所需的字节数。如果 rune 不是一个可用 UTF-8 编码的有效值，则返回 -1
! RuneStart 报告该字节是否可能是编码的、可能无效的 rune 的第一个字节
! Valid, ValidString 报告 p 是否完全由有效的 UTF-8 编码 rune 组成
! ValidRune 报告 r 是否可以合法地编码为 UTF-8
*/

func TestUtf8(t *testing.T) {
	t.Run("AppendRune", func(t *testing.T) {
		appendrune := func(r rune) {
			logfln("append %#U as %x`%[2]s`", r, AppendRune([]byte{}, r))
		}
		appendrune(rune(math.MaxInt32))
		appendrune(rune(RuneError))
		appendrune(rune(0))
		appendrune('H')
	})

	t.Run("DecodeRune", func(t *testing.T) {
		p := []byte("Hello, 世界")
		for FullRune(p) {
			r, n := DecodeRune(p)
			p = p[n:]
			logfln("read rune: %#U, %d", r, n)
		}
	})

	t.Run("DecodeLastRune", func(t *testing.T) {
		p := []byte("Hello, 世界")
		for {
			if r, n := DecodeLastRune(p); r != RuneError {
				logfln("%#U", r)
				p = p[:len(p)-n]
			} else {
				return
			}
		}
	})

	t.Run("EncodeRune", func(t *testing.T) {
		content := "Hello, 世界"
		p, l := make([]byte, 3*RuneCountInString(content)), 0
		for {
			if r, n := DecodeLastRuneInString(content); r != RuneError {
				content = content[0 : len(content)-n]
				if cap(p)-l < 3 {
					p = append(p, 0, 0, 0, 0)
				}
				n := EncodeRune(p[l:], r)
				l += n
			} else {
				p = p[:l]
				break
			}
		}
		logfln("%s", p)
	})

	t.Run("Valid", func(t *testing.T) {
		content := "Hello世界"[0:10]
		for !ValidString(content) {
			for _, b := range []byte(content) {
				if RuneStart(b) {
					logfln("%#x, %[1]c", b)
				} else {
					logfln("%#x", b)
					return
				}
			}
		}
	})
}
