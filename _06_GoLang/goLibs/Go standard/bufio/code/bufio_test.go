package gostd

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

var content = `This is a content for bufio. 
Hello World!
你好，世界
abcdefg123456789
`
var testInput = []byte("012\n345\n678\n9ab\ncde\nfgh\nijk\nlmn\nopq\nrst\nuvw\nxy")
var testInputrn = []byte("012\r\n345\r\n678\r\n9ab\r\ncde\r\nfgh\r\nijk\r\nlmn\r\nopq\r\nrst\r\nuvw\r\nxy\r\n\n\r\n")

/* buffered input
! NewReader, NewReaderSize 构造缓冲 io.Reader 对象。
! bufio.Reader 为 io.Reader 对象实现缓冲。
	Buffered 	返回已读取字节数
	Size 		基础缓冲区大小
	Discard 	丢弃最多 n 字节数
	Peek 		返回接下来 n 个字节但不推进 reader
	Read, ReadByte, ReadRune
	ReadBytes, ReadSlice, ReadString 连续读取直到首次遇到 `delim`
	ReadLine 	返回单行。过长时设置 isPrefix。其余部分在后续调用中返回
	Reset 		重置缓冲区并从 r 读取
	UnreadByte 	取消最近读取操作读取的最后一个字节。Peek、Discard 和 WriteTo 不被视为读操作
	UnreadRune 	取消上一次 `ReadRune` 操作读取的字符
	WriteTo     多次调用底层 Reader 的 Read 或 WriteTo（若实现）
*/

func TestBufReader(t *testing.T) {
	checkErr := func(t *testing.T, err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Run("Read", func(t *testing.T) {
		r := bufio.NewReaderSize(strings.NewReader(content), 64)
		size, buffered := r.Size(), r.Buffered()
		logfln("Buffer Reader: Size=%d, Buffered=%d", size, buffered)

		//>>>>> ReadByte
		for {
			if b, err := r.ReadByte(); err != nil || b == '\n' {
				checkErr(t, err)
				break
			} else {
				logfln("read byte: %#U", b)
			}
		}

		//>>>>> ReadBytes
		if bs, err := r.ReadBytes('\n'); err != nil {
			checkErr(t, err)
		} else {
			logfln("read %d bytes, the string is: %s", len(bs), bs)
		}

		//>>>>> ReadRune
		for {
			c, size, err := r.ReadRune()
			if err != nil || c == '\n' {
				checkErr(t, err)
				break
			}
			logfln("read rune: %#U, size : %d", c, size)
		}

		//>>>>> Discard, Peek
		if n, err := r.Discard(7); err != nil {
			checkErr(t, err)
		} else {
			bs, err := r.Peek(10)
			if err != nil {
				checkErr(t, err)
			}
			logfln("discard %d bytes, peek %d bytes and to string: %s", n, len(bs), bs)
		}
	})

	t.Run("ReadLine", func(t *testing.T) {
		var err error
		var line []byte
		testReadLine := func(input []byte) {
			r := bufio.NewReader(bytes.NewReader(input))
			for err == nil {
				line, _, err = r.ReadLine()
				if len(line) > 0 {
					logfln("read line: %s", line)
				}
			}
			if err != io.EOF {
				checkErr(t, err)
			}
			err = nil
		}
		testReadLine(testInput)
		testReadLine(testInputrn)
	})

	t.Run("read slice, string", func(t *testing.T) {
		r := bufio.NewReader(strings.NewReader(content))
		log(r.ReadString('\n'))

		slice, _ := r.ReadSlice('\n')
		logfln("%s", slice)

		for r.Buffered() > 0 {
			slice, _ = r.ReadSlice('\n')
			logfln("read slice: %s", slice[0:len(slice)-1])
		}
	})

	t.Run("WriteTo", func(t *testing.T) {
		r := bufio.NewReader(strings.NewReader("This is a content for WriteTo test\n"))
		r.WriteTo(os.Stdout)
	})
}

/* buffered output
! bufio.Writer 为 io.Writer 对象实现缓冲。
	Available 缓冲区可用字节数
	AvailableBuffer 返回 Available() 容量的空缓冲区。此缓冲区旨在追加并传递给紧接其后的 Writer.Write 调用
	Reset, Flush, Size, Buffered
	Write, WriteByte, WriteRune, WriteString
	ReadFrom 实现了 io.ReaderFrom。它调用底层 writer.ReadFrom（若实现）, 存在缓冲数据，则调用之前填充缓冲区并写入缓冲区
! NewWriterSize & NewWriter 构造缓冲 io Writer 对象。
*/

func TestBufWriter(t *testing.T) {
	// 写入到标准输出, 缓冲区大小为 10
	w := bufio.NewWriterSize(os.Stdout, 10)
	bs := bytes.NewBuffer([]byte(content))
	getBytes := func() []byte {
		if bytes, err := bs.ReadBytes('\n'); err == nil {
			return bytes
		}
		return nil
	}
	t.Run("WriteByte", func(t *testing.T) {
		defer w.Flush()
		for _, b := range getBytes() {
			if w.Available() < 1 {
				w.Flush()
			}
			w.WriteByte(b)
		}
	})

	t.Run("Write&AvailableBuffer", func(t *testing.T) {
		defer w.Flush()
		for _, b := range getBytes() {
			buf := w.AvailableBuffer()
			buf = append(buf, b, ',')
			w.Write(buf)
		}
	})

	t.Run("WriteRune", func(t *testing.T) {
		defer w.Flush()
		for _, r := range string(getBytes()) {
			w.WriteRune(r)
		}
	})

	t.Run("ReadFrom", func(t *testing.T) {
		defer w.Flush()
		if n, err := w.ReadFrom(bs); err == nil {
			logfln("total read %d bytes from r", n)
		}
	})
}

// ! ReadWriter 保存一组 Reader 和 Writer 的指针
func TestReadWriter(t *testing.T) {
	wr := bufio.NewReadWriter(bufio.NewReader(strings.NewReader(content)), bufio.NewWriter(os.Stdout))
	for {
		line, _, err := wr.ReadLine()
		if len(line) > 0 {
			wr.Write(line)
			wr.Flush()
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
	}
}

/* Scanner
! NewScanner 从 io.Reader 构造一个 Scanner
! bufio.Scanner 根据给定的拆分函数 split 读取数据并依次返回拆分后的令牌（数据）
	Buffer 	设置 Scan 使用的初始缓冲区以及最大缓冲区大小, 用以保存 Scan 读取的令牌
	Err 	返回 Scanner 遇到的第一个非 EOF 错误
	Scan 	使 Scanner 前进到下一个令牌，然后通过 Scanner.Bytes 或 Scanner.Text 获取该令牌。
				如果 split 函数返回太多的空令牌而没有推进输入则 panic
	Text, Bytes 读取调用 Scan 之后生成的最新令牌
	Split 	设置 split 拆分函数。默认为 ScanLines。如果在 Scan 后调用 Split，则会 panic
! ScanBytes, ScanLines, ScanRunes, ScanWords 为内置的拆分函数
! bufio.SplitFunc 表示 split 函数的签名
*/

func TestScanner(t *testing.T) {
	scan := func(fn bufio.SplitFunc, kind string) {
		scr := bufio.NewScanner(strings.NewReader(content))
		//? 调用 Scan 之前调用
		scr.Split(fn)
		scr.Buffer(make([]byte, 10), 128)
		for scr.Scan() {
			logfln("scan %s : %s", kind, scr.Text())
		}
		checkErr(scr.Err())
		log("")
	}

	scan(bufio.ScanLines, "line")
	scan(bufio.ScanRunes, "rune")
	scan(bufio.ScanWords, "word")
	scan(bufio.ScanBytes, "byte")
}

// ? 自定义拆分函数
func TestCustomSpiltFunc(t *testing.T) {
	const input = "1234 5678 1234567901234567890"
	scr := bufio.NewScanner(strings.NewReader(input))
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanWords(data, atEOF)
		if err == nil && token != nil {
			_, err = strconv.ParseInt(string(token), 10, 32)
		}
		return
	}
	scr.Split(split)
	for scr.Scan() {
		logfln("%s\n", scr.Text())
	}
	if err := scr.Err(); err != nil {
		logfln("Invalid input: %s", err)
	}
}
