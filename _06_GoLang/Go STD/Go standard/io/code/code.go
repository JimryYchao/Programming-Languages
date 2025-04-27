package gostd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	buf         = make([]byte, 512)
	hello       = "Hello World"
	READ        = "Read : `%s`"
	READ_BYTES  = "Read %d Bytes: `%s`"
	WRITE_BYTES = "Write %d Bytes: `%s`"
	EnterTest   = ">>> Enter %s :\n"
	EndTest     = ">>> End   %s\n"
	content     = "some io.Reader stream to be read"
)

func logCase(_case string) {
	logf("case : %s", _case)
}
func logf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}
func checkErr(err error) {
	if err == nil {
		return
	}
	logf("ERROR : %s", err)
}

func readToStdout(r io.Reader) {
	_, err := io.Copy(os.Stdout, r)
	println()
	checkErr(err)
}
func newReader(s string) io.Reader {
	return strings.NewReader(s)
}
func newBytesBuffer(n int) *bytes.Buffer {
	return bytes.NewBuffer(make([]byte, n))
}
func newStdoutWriter() io.Writer {
	return os.Stdout
}
