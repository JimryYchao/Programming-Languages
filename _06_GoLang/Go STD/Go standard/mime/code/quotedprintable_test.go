package gostd

import (
	"fmt"
	"io"
	"mime/quotedprintable"
	"os"
	"strings"
	"testing"
)

/*
! Reader, NewReader 是 quoted-printable 解码器
! Writer, NewWriter 是 quoted-printable 编码器
*/

func TestReader(t *testing.T) {
	for _, s := range []string{
		`=48=65=6C=6C=6F=2C=20=47=6F=70=68=65=72=73=21`,
		`invalid escape: <b style="font-size: 200%">hello</b>`,
		"Hello, Gophers! This symbol will be unescaped: =3D and this will be written in =\r\none line.",
	} {
		b, err := io.ReadAll(quotedprintable.NewReader(strings.NewReader(s)))
		fmt.Printf("%s %v\n", b, err)
	}

	// Hello, Gophers! <nil>
	// invalid escape: <b style="font-size: 200%">hello</b> <nil>
	// Hello, Gophers! This symbol will be unescaped: = and this will be written in one line. <nil>
}

func TestWriter(t *testing.T) {
	w := quotedprintable.NewWriter(os.Stdout)
	w.Write([]byte("These symbols will be escaped: =\a\f\t"))
	w.Close()
}
