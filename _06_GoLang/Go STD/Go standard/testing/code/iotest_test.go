package gostd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"testing/iotest"
)

var content = "this is a content for iotest."

/*	testing readers
! DataErrReader 更改 Reader 处理错误的方式：将错误和数据一起返回
! ErrReader 从所有 Read 调用中返回 0，err
! HalfReader 从 r 读取一半的请求字节来实现 Read
! NewReadLogger （using log.Printf）读取到标准错误, 打印前缀和读取的十六进制数据
! OneByteReader 从 r 中读取一个字节来实现每个非空 Read
! TestReader 测试从 r 中读取返回预期的 content 内容
! TimeoutReader 在第二次读取时返回无数据的错误 ErrTimeout。后续读取调用成功。
*/

func TestReaders(t *testing.T) {
	buf := make([]byte, 128)
	read := func(mess string, reader io.Reader) {
		reader = io.TeeReader(reader, os.Stdout)
		n, err := reader.Read(buf)
		logfln(fmt.Sprintf("\n%s ", mess)+"total read %d bytes: err : %v", n, err)
	}

	read("DataErrReader", iotest.DataErrReader(strings.NewReader(content)))
	read("ErrReader", iotest.ErrReader(fmt.Errorf("this is a error")))
	read("HalfReader", iotest.HalfReader(strings.NewReader(content)))
	read("NewReaderLogger", iotest.NewReadLogger("TEST", strings.NewReader(content)))
	read("OneByteReader", iotest.OneByteReader(strings.NewReader(content)))
	read("TestReader", iotest.ErrReader(iotest.TestReader(strings.NewReader(content), []byte(content+"\n"))))

	bufwr := &bytes.Buffer{}
	tr := iotest.TimeoutReader(bufwr)
	tr = io.TeeReader(tr, os.Stdout)

	checkTr := func() {
		bufwr.Reset()
		bufwr.Write([]byte(content))
		n, err := tr.Read(buf)
		logfln("\ntotal read %d bytes: err : %v", n, err)
	}
	checkTr()
	checkTr() // return ErrTimeout
	checkTr()
}

/* testing writers
! NewWriteLogger （using log.Printf）写入到标准错误，打印前缀和写入的十六进制数据。
! TruncateWriter 写入 w，但在 n 字节后停止。后续调入不再写入。
*/

func TestWriters(t *testing.T) {
	iotest.NewWriteLogger("TEST", &bytes.Buffer{}).Write([]byte(content))

	w := iotest.TruncateWriter(os.Stdout, 10)
	w.Write([]byte(content))
	w.Write([]byte(content)) // 不再写入
}
