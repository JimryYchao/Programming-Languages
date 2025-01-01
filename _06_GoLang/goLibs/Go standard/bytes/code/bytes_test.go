package gostd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
)

/* bytes functions
! Clone, Compare, Join, Repeat, ToValidUTF8
! Contains, ContainsAny, ContainsFunc, ContainsRune
! Cut, CutPrefix, CutSuffix
! HasPrefix, HasSuffix
! Equal, EqualFold
! Fields, FieldsFunc
! Index, IndexAny, IndexByte, IndexFunc, IndexRune
! LastIndex, LastIndexAny, LastIndexByte, LastIndexFunc
! Replace, ReplaceAll
! Split, SplitN, SplitAfter, SplitAfterN
! ToLower, ToLowerSpecial, ToTitle, ToTitleSpecial, ToUpper, ToUpperSpecial
! Trim, TrimFunc, TrimLeft, TrimLeftFunc, TrimRight, TrimRightFunc, TrimSpace, TrimSuffix
! Count 统计 s 中 sep 的非重叠实例的数量, sep 为零长切片则返回 utf8.RuneCount(s) + 1
! Map 对字节切片 s 逐 rune 使用 mapping 函数进行修改，mapping 返回负值则删除该字符
! Rune 将 []byte 解释为 []rune
*/
//? Some functions in `bytes` like `strings` see [package strings]
func TestBytesFunctions(t *testing.T) {
	t.Run("Count", func(t *testing.T) {
		log(bytes.Count([]byte("cheese"), []byte("e")))   // 3
		log(bytes.Count([]byte("cheese"), []byte("ee")))  // 1
		log(bytes.Count([]byte("cheeese"), []byte("ee"))) // 1
		log(bytes.Count([]byte("chees我"), []byte{}))      // 7 = 1+6
		log(bytes.Count([]byte("cheese"), nil))           // 7 = 1+6
	})
	t.Run("Map", func(t *testing.T) {
		rot13 := func(r rune) rune {
			switch {
			case r >= 'A' && r <= 'Z':
				return 'A' + (r-'A'+13)%26
			case r >= 'a' && r <= 'z':
				return 'a' + (r-'a'+13)%26
			}
			return r
		}
		fmt.Printf("%s\n", bytes.Map(rot13, []byte("'Twas brillig and the slithy gopher...")))
		// 'Gjnf oevyyvt naq gur fyvgul tbcure...
	})

	t.Run("Runes", func(t *testing.T) {
		rs := bytes.Runes([]byte("go gopher, 你好"))
		for _, r := range rs {
			logfln("%#U", r)
		}
	})
}

/*
! NewBuffer, NewBufferString
! bytes.Buffer 是一个可变大小的字节缓冲区，具有 Buffer.Read 和 Buffer.Write 方法
	Available, AvailableBuffer
	Cap, Grow, Len, Reset,
	Read, ReadByte, ReadBytes, ReadRune, ReadString, UnreadByte, UnreadRune, WriteTo
	Write, WriteByte, WriteRune, WriteString, ReadFrom
	Bytes 返回缓冲区未读取部分的字节切片，仅在下一次修改缓冲区之前有效
	String 以字符串返回未读取部分
	Next 返回缓冲区中接下来的 n 个字节的切片, 并使缓冲区前进（如同 Buffer.Read）
	Truncate 截断并丢弃 n 个未读字节以外的所有字节
*/

func TestBuffer(t *testing.T) {
	t.Run("Buffer", func(t *testing.T) {
		var b bytes.Buffer // init buffer, == bytes.NewBuffer(nil)
		b.Write([]byte("Hello "))
		fmt.Fprintf(&b, "world!\n")
		b.WriteTo(os.Stdout) // Hello World!

		// A Buffer can turn a string or a []byte into an io.Reader.
		buf := bytes.NewBufferString("R29waGVycyBydWxlIQ==")
		io.Copy(os.Stdout, base64.NewDecoder(base64.StdEncoding, buf)) // Gophers rule!
	})

	t.Run("Buffer.Bytes", func(t *testing.T) {
		buf := bytes.Buffer{}
		buf.Write([]byte{'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd', '\n'})
		os.Stdout.Write(buf.Bytes()) // hello world
		buf.Read(make([]byte, 6))
		os.Stdout.Write(buf.Bytes()) // world
	})

	t.Run("Buffer.Next", func(t *testing.T) {
		var b bytes.Buffer
		b.Write([]byte("abcde"))
		logfln("%s", b.Next(2)) // ab
		logfln("%s", b.Next(2)) // cd
		logfln("%s", b.Next(2)) // e
	})

	t.Run("Buffer.Truncate", func(t *testing.T) {
		b := bytes.NewBufferString("Hello World")
		b.Truncate(6)
		os.Stdout.Write(b.Bytes()) // Hello
		b.WriteString("WORLD")
		os.Stdout.Write(b.Bytes()) // Hello WORLD
	})
}

/*
! NewReader
! bytes.Reader 实现 io.Reader, io.ReaderAt, io.WriterTo, io.Seeker, io.ByteScanner, io.RuneScanner 接口以读取字节片。与 Buffer 不同，Reader 是只读的
	Len, Size, Reset, Seek, WriterTo
	Read, ReadAt, ReadByte, ReadRune, UnReadByte, UnReadRune
*/

func TestReadAtRace(t *testing.T) {
	// Test for the race detector, to verify ReadAt doesn't mutate any state.
	r := bytes.NewReader([]byte("0123456789"))
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			var buf [1]byte
			r.ReadAt(buf[:], int64(i))
			logfln("%c", buf[0])
		}(i)
	}
	wg.Wait()
}
