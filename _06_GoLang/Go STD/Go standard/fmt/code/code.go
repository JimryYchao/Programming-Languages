package gostd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync/atomic"
	"testing"
)

func logCase(_case string) {
	logfln("case : %s", _case)
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("LOG ERROR: \n%s", err)
}

func log(s any) {
	fmt.Println(s)
}
func logfln(format string, args ...any) {
	fmt.Printf(format+"\n", args...)

}

// check verbs

func checkVerbsHelpers(w io.Writer, name string, vs []any, verbs ...string) error {
	if len(vs) == 0 {
		return errors.New("no vs")
	}
	for _, v := range vs {

		var rt []byte
		if _, ok := v.(string); ok {
			rt = []byte(fmt.Sprintf("(T:%T, V:%s)", v, fmt.Sprintf("%q", v)))
		} else {
			rt = []byte(fmt.Sprintf("(T:%T, V:%v)", v, v))
		}

		if len(verbs) == 0 {
			return errors.New("no verbs")
		}
		fmts, okCh := make([][]byte, len(verbs)), make(chan bool)

		for i, verb := range verbs {
			i, verb := i, verb
			go func() {
				fmts[i] = []byte(fmt.Sprintf("%-10s", verb) + "|" + fmt.Sprintf(verb, v) + "|")
				okCh <- true
			}()
		}
		var c = 0
		for {
			<-okCh
			c++
			if c >= len(verbs) {
				break
			}
		}
		var sb strings.Builder
		sb.Write(rt)
		for _, format := range fmts {
			sb.WriteString("\n    " + string(format))
		}
		_, err := io.WriteString(w, sb.String()+"\n")
		if err != nil {
			return errors.Join(fmt.Errorf("write to file %s failed", name), err)
		}
	}
	return nil
}

type C struct {
	c   bool
	err error
}

type nhelper struct {
	n atomic.Uint32
}

func newCheckVerbs(t *testing.T, ch chan C) *nhelper {
	t.Helper()
	h := nhelper{atomic.Uint32{}}
	t.Cleanup(func() {
		var n = 0
		for {
			rt := <-ch
			if !rt.c {
				checkErr(rt.err)
			}
			n++
			if n >= int(h.n.Load()) {
				break
			}
		}
	})
	return &h
}

func (helper *nhelper) writer(t *testing.T, file string, stdout bool) io.Writer {
	if !stdout {
		f, err := os.OpenFile("files\\"+file, os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil || f == nil {
			return nil
		}
		t.Cleanup(func() {
			f.Close()
		})
		return f
	}
	return os.Stdout
}

func (helper *nhelper) checkVerbs(t *testing.T, stdout bool, file string, compCh chan C, vs []any, verbs ...string) {
	helper.n.Add(1)
	f := helper.writer(t, file, stdout)
	if f == nil {
		compCh <- C{false, nil}
		return
	}

	if err := checkVerbsHelpers(f, file, vs, verbs...); err != nil {
		io.WriteString(f, "\n\n ========== ERROR =========\n"+err.Error())
		compCh <- C{false, err}
		return
	}
	compCh <- C{true, nil}
}

type S string
type I int
type IwithString int

func (i IwithString) String() string {
	return fmt.Sprintf("<%d>", int(i))
}
