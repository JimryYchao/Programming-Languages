package gostd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

/* Reader
! io.Reader 包装基本的 `Read` 方法；`Read` 返回读取的字节数
! MultiReader 串联一组 Readers，这些 Readers 在内部按顺序 Read
*/

func TestReader(t *testing.T) {
	r := newReader("some io.Reader stream to be read")
	if c, err := r.Read(buf); err == nil && c > 0 {
		logf(READ_BYTES, c, buf)
	} else {
		// handle error
	}
	// output: some io.Reader stream to be read
}

func TestMultiReader(t *testing.T) {
	t.Run("MultiReader", func(t *testing.T) {
		r1 := newReader("first reader ")
		r2 := newReader("second reader ")
		r3 := newReader("third reader")
		// 按顺序调用 Reader.Read
		mr := io.MultiReader(r1, r2, r3)
		readToStdout(mr)
		// output: first reader second reader third reader
	})
	t.Run("MultiReaderAsWriterTo", func(t *testing.T) {
		mr := io.MultiReader(strings.NewReader("Hello "),
			io.MultiReader(strings.NewReader(""),
				strings.NewReader("World")))
		// MultiReader 内部构造一个 multiReader, 并实现了 io.WriterTo
		if mrAsWriterTo, ok := mr.(io.WriterTo); ok {
			mrAsWriterTo.WriteTo(os.Stdout) // Hello World
		}
	})
}

// ! LimitReader 包装一个限制读取字节数的 Reader *LimitedReader
func TestLimitReader(t *testing.T) {
	lr := io.LimitReader(newReader("some io.Reader stream to be read"), 4) // 限制读取字节数
	readToStdout(lr)
	lr2 := io.LimitReader(newReader(hello), -1)
	readToStdout(lr2) // N < 0, lr2.Read return EOF
	// output: some
}

// ! TeeReader 返回一个关联 w 和 r 的 Reader，从 r 读取的内容会相应的写入 w
func TestTeeReader(t *testing.T) {
	r := newReader("some io.Reader stream to be read\n")
	tr := io.TeeReader(r, os.Stdout) // 关联 tr 到 Stdout

	// 从 tr 的任何读取都会复制到 stdout
	if _, err := io.ReadAll(tr); err != nil {
		t.Fatal(err)
	}
}

/* Writer
! io.Writer 包装基本的 `Write` 方法；`Write` 将最多 `len(p)` 字节写入到底层数据流。
! MultiWriter 串联一组 Writer，这些 Writers 在内部按顺序 Write
*/

func TestWriter(t *testing.T) {
	//? os.Stdout
	var w io.Writer = newStdoutWriter()
	w.Write([]byte("Writing to os.Stdout\n"))

	//? bytes.Buffer
	var bw *bytes.Buffer = newBytesBuffer(128)
	c, err := bw.Write([]byte("Writing to bytes.Buffer"))
	checkErr(err)
	logf(WRITE_BYTES, c, bw.Bytes())

	// output:
	// `Writing to os.Stdout`
	// `Writing to bytes.Buffer`
}

func TestMultiWriter(t *testing.T) {
	t.Run("MultiWriter", func(t *testing.T) {
		w1 := newBytesBuffer(5)
		w2 := &strings.Builder{}
		w3 := os.Stdout

		mw := io.MultiWriter(w1, w2, w3)
		mw.Write([]byte(hello))

		logf("\nbytes.Buffer : %s", w1.Bytes())
		logf("strings.Builder : %s", w2.String())

		// output:
		// Hello World
		// bytes.Buffer : Hello World
		// strings.Builder : Hello World
	})
	t.Run("MultiWriterAsStringWriter", func(t *testing.T) {
		w1 := newBytesBuffer(5)
		w2 := &strings.Builder{}
		w3 := os.Stdout

		if sw, ok := io.MultiWriter(w1, w2, w3).(io.StringWriter); ok {
			sw.WriteString(hello)
		}
		logf("\nbytes.Buffer : %s", w1.Bytes())
		logf("strings.Builder : %s", w2.String())

		// output:
		// Hello World
		// bytes.Buffer : Hello World
		// strings.Builder : Hello World
	})
}

/* Close
! io.Closer 包装基本的 Close 方法。首次调用后的 `Close()` 行为未定义。特定的实现可以记录它们自己的行为。
! NopCloser 返回一个 `io.ReadCloser`，它带有一个无操作的 `Close` 方法，包装了提供的 `Reader r`。如果 `r` 实现了 `WriterTo`，则返回的 `ReadCloser` 将通过转发对 `r` 的调用来实现 `WriterTo`。
*/

func TestCloser(t *testing.T) {
	t.Run("Close", func(t *testing.T) {
		tmpfile, err := os.CreateTemp(t.TempDir(), "tmpFile")
		if err != nil || tmpfile == nil {
			t.Fatalf("CreateTemp(%s) failed: %v", "tmpFile", err)
		}
		defer func() {
			checkErr(tmpfile.Close()) // 在首次调用 Close 之后都会返回 err 或其他实现定义的行为
		}()
		tmpfile.Close() // 首次调用
	})

	t.Run("NopCloser", func(t *testing.T) {
		tmpfile, _ := os.CreateTemp(t.TempDir(), "tmpFile")
		readCloser := io.NopCloser(tmpfile)
		checkErr(tmpfile.Close())    // Closer 正常关闭
		checkErr(readCloser.Close()) // 转发一个无操作的 Closer, 永远不会发生 err
	})
}

/* Seek
! io.Seeker 包装基本的 Seek 方法；Seek 将下一次读取或写入的偏移量依照 `whence` 设置为 `offset`; `whence` 解释为：
		SeekStart		相对于开始
		SeekEnd 		相对于末尾
		SeekCurrent 	相对于当前偏移量
*/

func TestSeek(t *testing.T) {
	file, err := os.OpenFile("./seek.file", os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	file.Write([]byte(content + "\n"))
	t.Cleanup(func() {
		file.Close()
		os.Remove("./seek.file")
	})
	readFile("./seek.file") // 更改之前的文件内容

	file.Seek(8, io.SeekStart) // 跳转至开头的 +8 offset
	file.WriteString("Writer")
	readFile("./seek.file") // 检查更改后的内容

	file.Seek(-5, io.SeekEnd) // 跳转至末尾的 -5 offset
	file.WriteString("write\n")
	readFile("./seek.file") // 检查更改后的内容

	// output:
	// some io.Reader stream to be read
	// some io.Writer stream to be read
	// some io.Writer stream to be write
}

/* Copy Functions
! Copy 将副本从 `src` 复制到 `dst`。它返回复制的字节数和第一个错误（如果有）。
! CopyBuffer 等效于 Copy。不能提供 0 长度的 `buf`，传递 `nil` 时将内部创建一个 `buf`。
! CopyN 将最多 `n` 个字节从 `src` 复制到 `dst`）。
*/

func TestCopyFunctions(t *testing.T) {
	var wb *bytes.Buffer = newBytesBuffer(128)
	t.Run("Copy", func(t *testing.T) {
		//? Copy
		wb.Reset() // 从 Reader 到 Writer 复制
		if c, err := io.Copy(wb, newReader(hello)); err != nil {
			t.Fatal(err)
		} else if c > 0 {
			logf(READ_BYTES, c, wb.Bytes())
		}

		//? Copy with negative LimitedReader
		wb.Reset() // N < 0, 将返回 ""
		c, _ := io.Copy(wb, &io.LimitedReader{R: newReader(hello), N: -1})
		logf(READ_BYTES, c, wb.Bytes())

		// output:
		// `Hello World`
		// ``
	})

	t.Run("CopyBuffer", func(t *testing.T) {
		//? CopyBuffer
		wb.Reset() // `CopyBuffer(dst, src, nil)` same as `Copy(dst, src)`
		if c, err := io.CopyBuffer(wb, newReader(hello), buf); err != nil {
			t.Fatal(err)
		} else {
			logf(READ_BYTES, c, wb.Bytes())
		}

		//? CopyBuffer with empty buffer
		defer func() {
			if err := recover(); err != nil {
				logf("Panicking : %s", err)
			}
		}()
		wb.Reset() // panicking with empty buf
		io.CopyBuffer(wb, newReader(hello), []byte{})

		// output:
		// `Hello World`
		// Panicking : empty buffer in CopyBuffer
	})

	t.Run("CopyN", func(t *testing.T) {
		copyN := func(_case string, n int64) {
			wb.Reset()
			c, err := io.CopyN(wb, newReader(hello), n)
			logCase(_case)
			checkErr(err)
			logf(READ_BYTES, c, wb.Bytes())
		}

		copyN("CopyN with small N: len(hello) > 5, return (5, nil)", 5)                 // `Hello`
		copyN("CopyN with negative N: N < 0, return (0, nil)", -1)                      // ``
		copyN("CopyN with large N: len(hello) < 100, return (len(hello), io.EOF)", 100) // `Hello World`, EOF
	})
}

/* PipeReader & PipeWriter
! Pipe 创建一组同步内存管道。它用于连接一组 `io.Reader` 和 `io. Writer`。管道上的读取和写入是一（多）对一匹配的。
! io.PipeReader 是 Pipe() 的读取端。
	Close			关闭读取端；管道的后续写入将返回错误 `ErrClosedPipe`。
	CloseWithError	关闭读取端；管道的后续写入将返回错误 `err`。
	Read			从管道中读取数据，阻塞直到写入端末尾或写入端关闭。如果写入端因错误而关闭，则返回该错误。
! io.PipeWriter 是 Pipe() 的写入端。
	Close 			关闭写入端；管道的后续读取 Read 将返回 `(0, EOF)`。
	CloseWithError	关闭写入端；管道的后续读取将返回错误 `(0, err)` 或 `(0, EOF) // err == nil`。
	Write 			将数据写入管道，阻塞直到一个或多个读取端消耗了所有数据或读取端关闭（返回 `ErrClosedPipe`）。如果读取端因错误而关闭，则返回该错误。
*/

func TestPipe(t *testing.T) {
	checkWrite := func(w io.Writer, data []byte, c chan int) {
		if _, err := w.Write(data); err != nil {
			checkErr(err)
		}
		c <- 0
	}
	reader := func(r io.Reader, c chan int) {
		for {
			n, err := r.Read(make([]byte, 64))
			if err != nil {
				if err != io.EOF {
					checkErr(err)
				}
				break
			}
			c <- n
		}
	}

	t.Run("A single r/w pair", func(t *testing.T) {
		c := make(chan int)
		pr, pw := io.Pipe()
		var buf = make([]byte, 64)

		go checkWrite(pw, []byte(content), c) // write part
		n, err := pr.Read(buf)                // read part
		checkErr(err)
		logf(READ_BYTES, n, buf)
		<-c // 等待读取结束

		pr.Close()
		pw.Close()
	})

	t.Run("A sequence of r/w pairs", func(t *testing.T) {
		c := make(chan int)
		pr, pw := io.Pipe()

		go reader(pr, c) // continuous read
		var buf = make([]byte, 64)
		for i := 0; i < 5; i++ {
			p := buf[0 : 5+i*10]
			if _, err := pw.Write(p); err != nil {
				checkErr(err)
				return
			}
			logf("Reader read %d bytes", <-c)
		}
		pw.Close()
	})

	t.Run("A large write and multiple reads", func(t *testing.T) {
		c := make(chan struct {
			n   int
			err error
		})
		pr, pw := io.Pipe()
		go func() { // write a large date
			n, err := pw.Write(make([]byte, 1024))
			pw.Close()
			c <- struct {
				n   int
				err error
			}{n, err}
		}()

		tot := 0
		rdat := make([]byte, 64)
		for {
			n, err := pr.Read(rdat)
			if err != nil {
				tot += n
				logf("Read %d bytes and read %s", n, err)
				break
			}
			tot += n
			logf("Read %d bytes", n)
		}
		rt := <-c
		checkErr(rt.err)
		if rt.n != tot {
			t.Fatalf("total read %d != 1024", tot)
		}
	})

	t.Run("close writer", func(t *testing.T) {
		pr, pw := io.Pipe()
		// PipeWriter
		go func() {
			defer pw.Close() // 后续返回 EOF
			for range 3 {
				pw.Write([]byte(hello))
			}
		}()
		// PipeReader
		defer pr.Close()
		for {
			if _, err := pr.Read(buf); err != nil {
				checkErr(err) // EOF
				break
			}
		}
	})

	t.Run("close reader", func(t *testing.T) {
		pr, pw := io.Pipe()
		// PipeReader
		go func() {
			defer pr.Close()
			for range 3 {
				pr.Read(buf)
				time.Sleep(500 * time.Millisecond)
			}
		}()
		// PipeWriter
		defer pw.Close()
		for {
			if _, err := pw.Write([]byte(hello)); err != nil {
				checkErr(err)
				break
			}
		}
	})

	t.Run("close with error", func(t *testing.T) {
		pr, pw := io.Pipe()
		var uerr = errors.New("user error")
		go func() {
			pw.CloseWithError(uerr)
		}()
		if _, err := pr.Read(buf); err != nil {
			checkErr(err)
		}
	})
}

/* ReadAt
! io.ReaderAt 包装 `ReadAt` 方法。`ReadAt` 从底层输入源中的偏移 `off` 开始将最多 `len(p)` 字节读入 `p`。
! io.SectionReader 在底层 `ReaderAt` 的片段上实现 `Read`、`Seek` 和 `ReadAt`。
	Outer 返回底层 `(ReaderAt, off, n)`。是创建它的 `NewSectionReader` 的逆运算。
	Size 返回片段的字节大小。
! NewSectionReader 包装一个 `ReaderAt` 并返回一个 `io.SectionReader`，它从 `off` 偏移开始读取，并在最多 `n` 个字节处停止。
*/

func TestReadAt(t *testing.T) {
	t.Helper()
	t.Run("ReadAt", func(t *testing.T) {
		var ra io.ReaderAt = strings.NewReader(content)
		readAt := func(_case string, lenbuf int64, offset int64) {
			buf := make([]byte, lenbuf)
			n, err := ra.ReadAt(buf, offset)
			logCase(_case)
			logf(READ_BYTES, n, buf)
			checkErr(err)
		}
		readAt("lenbuf(50) > len(content) - offset(10)", 50, 10)   // ok, EOF
		readAt("len(content) - offset(10) > len(content)", 15, 10) // ok
		readAt("offset(40) > len(content)", 10, 50)                // read 0, EOF
		readAt("negative offset(-1)", 10, -1)                      // read 0, Err : negative offset
	})

	t.Run("SectionReader", func(t *testing.T) {
		readSec := func(_case string, offset, n int64) {
			sr := io.NewSectionReader(strings.NewReader(content), offset, n)
			logCase(_case)
			readToStdout(sr)
		}

		readSec("offset + n < len(content)", 10, 10)                        // ok
		readSec("offset + n > len(content); offset < len(content)", 10, 50) // ok, full read
		readSec("offset > len(content)", 50, 10)                            // ""
		readSec("negative offset", -10, 20)                                 // "", Err : negative offset
		readSec("negative n", 10, -1)                                       // ok, full read
	})
}

func readFile(fname string) {
	if f, err := os.Open(fname); err == nil {
		io.Copy(os.Stdout, f)
		f.Close()
	}
}

/* WriteAt
! io.WriterAt 包装 `WriteAt` 方法。`WriteAt` 将 `len(p)` 个字节从 `p` 写入偏移量为 `off` 的底层数据流。
! io.OffsetWriter 将基础偏移量处的写入映射到基础写入器中的偏移量 base+off。
! NewOffsetWriter 返回一个 `OffsetWriter`，它从 `off` 偏移开始 `WriterAt` 写入。
*/

func TestWriteAt(t *testing.T) {
	fname := "./hello.file"
	if f, err := os.Open(fname); err == nil {
		f.Close()
		os.Remove(fname)
	}
	t.Cleanup(func() {
		os.Remove(fname)
	})

	f, _ := os.Create(fname)
	f.WriteString(content + "\n")
	f.Close()

	t.Run("WriteAt", func(t *testing.T) {
		if fwrAt, err := os.OpenFile(fname, os.O_RDWR, 0644); err != nil {
			t.Fatal(err)
		} else {
			io.Copy(os.Stdout, fwrAt)
			// read file: some io.Reader stream to be read
			fwrAt.WriteAt([]byte(hello), 10)
			fwrAt.Close()
			readFile(fname)
			// read again: some io.ReHello World to be read
		}
	})

	t.Run("OffsetWriter", func(t *testing.T) {
		if f, err := os.OpenFile(fname, os.O_RDWR, 0644); err != nil {
			t.Fatal(err)
		} else {
			offwr := io.NewOffsetWriter(f, 10)
			// before: some io.ReHello World to be read
			_, err := offwr.WriteAt([]byte("ader stream"), 0)
			// now: some io.Reader stream to be read
			checkErr(err)
			if _, err := offwr.Seek(5, io.SeekStart); err == nil { // 移动到末尾，在换行符之前
				offwr.Write([]byte("STREAM")) // expect: some io.Reader STREAM to be read
			}
			f.Close()
			readFile(fname)
		}
	})
}

/* ReaderFrom & WriterTo
! io.ReaderFrom 包装 `ReadFrom` 方法。`ReadFrom` 从 `r` 读取数据并返回值读取的字节数 `n`。
! io.WriterTo 包装 `WriteTo` 方法。`WriteTo` 将数据写入 `Writer`。
*/

func TestReadFromAndWriteTo(t *testing.T) {
	file, err := os.OpenFile("./writeTo.file", os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		file.Close()
		os.Remove("./writeTo.file")
	})

	writeTo := func(wt io.WriterTo) {
		wt.WriteTo(file)
	}

	rf := newBytesBuffer(16)
	for i := range 5 {
		rf.ReadFrom(strings.NewReader(fmt.Sprintf("line%d : %s\n", i+1, content)))
		writeTo(rf)
	}
	readFile("./writeTo.file") // content
}

/* Bytes Readers & Writers
! io.ByteReader 包装 `ReadByte` 方法。`ReadByte` 读取并返回输入中的下一个字节或错误。
! io.ByteScanner 将 `UnreadByte` 方法添加到 `ByteReader`。它导致下次调用 `ReadByte` 将返回最后读取的字节。
! io.ByteWriter 包装 `WriteByte` 方法。
*/

func TestBytesReadWriter(t *testing.T) {
	brw := bytes.NewBuffer([]byte(hello))
	// read bytes
	for {
		if b, err := brw.ReadByte(); err != nil {
			if err == io.EOF {
				logf("Read EOF")
				brw.UnreadByte() // 对 brw 的下次 `ReadByte` 将返回 b
				logf("Last byte in buf is %#v %c", b, b)
				err = nil
				break // 忽略 EOF 并停止读取字节
			}
			t.Fatal(err)
		} else {
			logf("Read byte : %#v %c", b, b)
		}
	}

	// write bytes
	brw.Reset() // 清空 buf
	sr := strings.NewReader(content)
	for b, err := sr.ReadByte(); err == nil; b, err = sr.ReadByte() {
		brw.WriteByte(b)
		logf("Write byte : %#v %c", b, b)
	}
	readToStdout(brw) // some io.Reader stream to be read
}

/* Rune Reader
! io.RuneReader 包装 `ReadRune` 方法。`ReadRune` 读取单个编码的 Unicode 字符，并返回该字符及其字节大小。
! io.RuneScanner 将 `UnreadRune` 方法添加到 `RuneReader`。`UnreadRune` 导致下一次调用 `ReadRune` 返回最后读取的字符。
*/

func TestRuneRead(t *testing.T) {
	sr := strings.NewReader("Hello, 你好")
	for r, size, err := sr.ReadRune(); err == nil; r, size, err = sr.ReadRune() {
		logf("Read rune : %c of %d bytes", r, size)
	}
}

/* WriteString
! io.StringWriter 包装 `WriteString` 方法。
! WriteString 将字符串 `s` 的内容写入 `w`。如果 `w` 实现了 `StringWriter`，则直接调用 `StringWriter.WriteString`。
*/

func TestStringWriter(t *testing.T) {
	var sw io.StringWriter = os.Stdout
	sw.WriteString(hello + "\n")                           // Hello World
	io.WriteString(os.Stdout, strings.ToLower(hello)+"\n") // hello world
}

/* Read Functions
! ReadAll 从 `Reader` 开始读取并返回读取的数据，直到出现错误或 EOF。
! ReadAtLeast 从 `Reader` 读取到 `buf`，直到它至少读取了 `min` 字节。
! ReadFull 将最多 `len(buf)` 字节从 `Reader` 精确读取到 `buf`。
*/

func TestReadFunctions(t *testing.T) {
	content := "some io.Reader stream to be read"
	logBuf := func(buf []byte, err error, n int) {
		if n < 0 {
			logf(READ, buf)
		} else {
			logf(READ_BYTES, n, buf)
		}
		checkErr(err)
	}

	t.Run("ReadAll", func(t *testing.T) {
		buff, err := io.ReadAll(newReader(content))
		logBuf(buff, err, -1)
		// output :
		// `some io.Reader stream to be read`
	})

	t.Run("ReadAtLeast", func(t *testing.T) {
		readAtLeast := func(_case string, min, lenbuf int) {
			buff := make([]byte, lenbuf)
			logCase(_case)
			n, err := io.ReadAtLeast(newReader(content), buff, min)
			logBuf(buff, err, n)
		}
		readAtLeast("min(10) < len(content) < lenbuf(60))", 10, 60)
		readAtLeast("min(10) < lenbuf(15) < len(content)", 10, 15)
		readAtLeast("len(content) < min(50) < lenbuf(60)", 50, 60)
		readAtLeast("len(content) < lenbuf(50) < min(60)", 60, 50) // io.ErrUnexpectedEOF: unexpected EOF
		readAtLeast("lenbuf(10) < len(content) < min", 50, 10)     // io.ErrShortBuffer: short buffer
		readAtLeast("lenbuf(10) < min(15) < len(content)", 15, 10) // io.ErrShortBuffer: short buffer
	})

	t.Run("ReadFull", func(t *testing.T) {
		readFull := func(_case string, lenbuf int) {
			buff := make([]byte, lenbuf)
			logCase(_case)
			n, err := io.ReadFull(newReader(content), buff)
			logBuf(buff, err, n)
		}
		readFull("lenbuf < len(content)", 10)
		readFull("lenbuf > len(content)", 50) // ERROR : unexpected EOF
		readFull("lenbuf = len(content)", len(content))
	})
}
