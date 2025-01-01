package gostd

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"testing"
)

// ! Verbs tests

func TestFormatVerbs(t *testing.T) {
	var compCh = make(chan C)
	h := newCheckVerbs(t, compCh)

	go h.checkVerbs(t, false, "verbs sqxX", compCh, []any{"abcxyz"}, "%v", "%#v", "%s", "%q", "%x", "%X", "% x", "% X", "%#x", "%#X", "%# x", "% #X")

	go h.checkVerbs(t, false, "escape character", compCh, []any{"\a\b\f\b\r\t\v\"\\", "\U0010ffff", string(rune(0x110000)), ""}, "%v", "%#v", "%s", "%x", "%q", "%+q", "%#q", "%#+q")

	// "a", []byte("a"), [1]byte{'a'}, &[1]byte{'a'} 在字符串格式化上等效
	go h.checkVerbs(t, false, "string", compCh, []any{"你a\xffz☺", S("hello"), I(23), IwithString(23), "a", []byte("a"), [1]byte{'a'}, &[1]byte{'a'}},
		"%v", "%#v", "%s", "%2s", "%.2s", "%25s", "%-25s", "%025s", "%25.2s", "%025.2s", "%-25.2s", "%+25.2s",
		"%q", "% q", "%#q", "%+q", "%+#q", "%25q", "%+25q", "%-25q", "%+-25q", "%025q", "%+025q", "%-025q", "%+-025q", "%.4q", "%25.4q", "%025.4q", "%-25.4q",
		"%x", "% x", "%#x", "%.4x", "%.1x", "%25x", "%025x", "%+25x", "%-25x", "% 25.1x", "% X", "% 25.4X")

	go h.checkVerbs(t, false, "integer", compCh, []any{uint16(12345), int16(-12345), 0, 7, 12345, -12345, ^uint16(0), ^uint64(0), 07531, -06420, 0x12af, -0x12af},
		"%v", "%#v", "%d", "%.d", "%8.d", "% d", "%+d", "%-d", "%8d", "%+8d", "%-8d", "%-+8d", "%08d", "%+08d", "%-08d", "%-+08d", "%20.8d", "%020.8d", "%-20.8d",
		"%b", "%#b",
		"%o", "%#o", "%O",
		"%#X", "%x", "%X", "%-#X", "%20.8X", "%-20.8X", "%-#20.8X", "%-#20.8o")

	go h.checkVerbs(t, false, "verb U", compCh, []any{0, -1, rune('\n'), rune('x'), rune('\u263a'), rune('⌘'), rune('我'), rune('\U0001D6C2'), rune('\U0010ffff'), rune(' ')},
		"%U", "%#U", "%+U", "%# U", "%#.2U", "%#14.6U", "%#-14.6U", "%.20U", "%#.20U", "%#014.6U", "%#-014.6U")

	NaN := math.NaN()
	posInf := math.Inf(1)
	negInf := math.Inf(-1)
	fverbs := []string{
		"%v", "%#v", "%f", "%e", "%E", "%x", "%X", "%g", "%G", "%+.3e", "%+.3x", "%+.3f", "%+07.2f", "%-07.2f", "%+-07.2f", "%+10.2f", "% .3f", "% .3e", "% .3x",
		"% .3g", "%+.3g", "%#g", "%#.f", "%#.e", "%#.g", "%#.x", "%#.4f", "%#.4e", "%#.4g", "%#.4x", "%#9.4g", "%b", "%#b", "%#.4b", "%#.68f",
		"%20g", "% 20g", "%020g", "%+20g", "%-20g", "%+-20g", "%+020g", "%-020g",
	}
	go h.checkVerbs(t, false, "float-point", compCh, []any{0.0, 1.0, -1.0, float32(-1.0), float32(1.0), 1e-33323, 123456.0, 1234567.0, 1230000.0, 1100000.0, 1.234, 0.1234, 1.23, 0.123, 12.3, 123.0, 123000.0, posInf, negInf, NaN},
		fverbs...)

	go h.checkVerbs(t, false, "complex", compCh, []any{0i, complex64(0i), -1i, 1i, 1 + 2i, -1 - 2i, complex64(1 + 2i), 123456 + 789012i, 1e-10i, -1e10 - 1.11e100i, 1.23 + 1.0i, 0 + 100000i, 1230000 + 0i, 1230000 - 0i, 1 + 1.23i, 123 + 1.23i,
		complex(posInf, posInf), complex(negInf, negInf), complex(NaN, NaN)}, fverbs...)

	array := [5]int{1, 2, 3, 4, 5}
	iarray := [4]any{1, "hello", 2.5, nil}
	slice := array[:]
	islice := iarray[:]
	type renamedUint8 uint8
	barray := [5]renamedUint8{1, 2, 3, 4, 5}
	bslice := barray[:]
	go h.checkVerbs(t, false, "array and slice", compCh, []any{array, iarray, barray, &array, &iarray, &barray, slice, islice, bslice, &slice, &islice, &bslice}, "%v")

	go h.checkVerbs(t, false, "byte array and slice", compCh, []any{[3]byte{65, 66, 67}, [1]byte{123}, []byte{}, []byte{1, 11, 111}},
		"%v", "%#v", "%b", "%c", "%d", "%o", "%U", "%v", "% v", "%+v", "%012v", "%#012v", "%6v", "%06v", "%-6v", "%-06v", "%#v", "%#6v", "%#06v", "%#-6v", "%#-06v", "%# -6v", "%#+-6v")
}

/* Printers
! Printf-like 使用格式化 format 输出。
! Print-like 函数在两个非字符串参数之间添加一个空格。对每个参数使用 `%v`。
! Println-like 函数在操作数之间空格，并追加一个换行符。
*/

func TestPrinters(t *testing.T) {
	printErrNew := func(mess string) error {
		return fmt.Errorf("printer error: %s", mess)
	}
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "fmt")
	if err != nil {
		t.Fatal(printErrNew(err.Error()))
	}
	defer file.Close()
	var WriteFormat = "Write: %s"
	var printArgs = []any{1, "H", true, "E", 'A', "L",
		struct{ name string }{"World"}, "L", []byte("!!!"), "O", [...]int{1, 2, 3, 4}, []int{55, 6, 7, 8},
		map[string]int{"a": 1, "b": 2}, &struct{ n int }{100},
		func(a, b int) int { return a + b }, make(chan bool), new(fmt.Formatter)}

	fmt.Printf("\nfmt.Sprintf: \n%s", fmt.Sprintf(WriteFormat, "test Printf\n"))
	fmt.Printf("\nfmt.Sprint: \n%s", fmt.Sprint(printArgs...))
	fmt.Printf("\nfmt.Sprintln: \n%s", fmt.Sprintln(printArgs...))

	fmt.Fprintf(file, WriteFormat, "test Fprintf\n")
	io.Copy(os.Stdout, file)
}

/*
! fmt.GoStringer 提供 `GoString` 方法以定义实现接口值的 Go 语法，GoString 返回值用于传递给 `%#v` 的格式打印。
! fmt.Formatter 提供 `Format` 方法以定义用户格式化输出。实现控制如何解释 `fmt.State` 和 `rune`，并通过 Sprint 或 FPrint(f) 来生成输出。
! fmt.State 表示传递给自定义 formatters 的 printer 状态。它提供对 `io.Writer` 接口的访问，以及有关操作数格式说明符的标志和选项的信息。
	Width 返回宽度值以及是否被设置
	Precision 返回精度值以及是否被设置
	Flag 报告标志 `c` 是否被设置
! fmt.Stringer 提供 `String` 方法以定义实现接口值的 “native” 格式。String 返回值打印到任何接受字符串格式或打印到非 format printer（如 `Print`）。
*/

func TestFmtCustom(t *testing.T) {
	var compCh = make(chan C)
	h := newCheckVerbs(t, compCh)

	vs := []any{F(1), G(2), FG{F(3), G(4)}, SF(5)}

	go h.checkVerbs(t, true, "Formatter", compCh, vs, "%T", "%v", "%.v", "%3.2v", "%.2v", "%3v", "%+v", "%q", "%x", "%a", "%A")

	go h.checkVerbs(t, true, "GoString", compCh, vs, "%#v", "%#q", "%#x")

	go h.checkVerbs(t, true, "Stringer", compCh, vs, "%s", "%#s", "%1.1s")
}

type (
	F  int
	G  int
	FG struct {
		F F
		G G
	}
	SF int
)

func (n F) Format(s fmt.State, c rune) {
	var format, flag string
	if w, ok := s.Width(); ok {
		format += fmt.Sprintf(" w:%d", w)
	}
	if p, ok := s.Precision(); ok {
		format += fmt.Sprintf(" p:%d", p)
	}

	if s.Flag('+') {
		flag = "+"
	}
	format = fmt.Sprintf("<%s=%sF(%d)", fmt.FormatString(s, c), flag, int(n)) + format + ">"
	fmt.Fprint(s, format)
}

func (n G) GoString() string {
	return fmt.Sprintf("GoString(%d)", int(n))
}

func (n SF) Format(s fmt.State, c rune) {
	switch c {
	case 'a':
		fmt.Fprintf(s, "<%c=custom(%d)>", c, int(n))
	case 'A':
		fmt.Fprintf(s, "<%c=CUSTOM(%d)>", c, int(n))
	default:
		fmt.Fprintf(s, "<%c=F(%d)>", c, int(n))
	}
}
func (n SF) String() string {
	return fmt.Sprintf("String(%d)", int(n))
}

/* Scanners
! Scanf-like 根据格式字符串解析参数。
! Scan-like 将输入的换行符视为空格。
! Scanln-like 在扫描到换行符或 EOF 后停止。
	Scanners 的 verbs 中不包括 %T、%p、#、+；除了 %c 其他 verbs 的实现将忽略任何的前导空格；%s
*/

func TestScanners(t *testing.T) {
	var text string = "5 true gophers\naaaaa"
	t.Run("Scan-like", func(t *testing.T) {
		var (
			b     bool
			n     int
			s, s0 string
		)
		if c, err := fmt.Sscan(text, &n, &b, &s, &s0); err == nil {
			fmt.Printf("Read %d : %s", c, fmt.Sprintln(n, b, s, s0))
		}
	})
	t.Run("Scanln-like", func(t *testing.T) {
		var (
			b bool
			n int
			s string
		)
		if c, err := fmt.Sscanln(text, &n, &b, &s); err == nil {
			fmt.Printf("Read %d : %s", c, fmt.Sprintln(n, b, s))
		}
	})

	t.Run("Scanf-like", func(t *testing.T) {
		var (
			r rune
			b bool
			s string
			h int
		)
		if c, err := fmt.Sscanf(text, "%c %t %s\n%x", &r, &b, &s, &h); err == nil {
			fmt.Printf("Read %d : %s", c, fmt.Sprintln(r, b, s, h))
		}
	})
}

// ! fmt.Scanner 提供 `Scan` 方法，该方法扫描输入以获得值的表示并将结果存储在接收器（有效指针）中。对于实现 `fmt.Scanner` 的任何参数 arg，`Scan`，`Scanf` 或 `Scanln` 方法都会调用它的 arg.Scan。
// ! fmt.ScanState 表示传递给自定义 Scanner.Scan 的参数 state。Scanners 可以执行一次运行扫描或要求 ScanState 发现下一个空格分隔的令牌。

func TestScanInts(t *testing.T) {
	var r RecursiveInt
	if c, err := fmt.Sscan(string(makeInts(100)), &r); err == nil {
		fmt.Print(c, &r)
	}
}

// RecursiveInt 接受一个形如 %d.%d.%d.... 的字符串并将其解析为一个链表。
type RecursiveInt struct {
	i    int
	next *RecursiveInt
	len  int
}

func (r *RecursiveInt) String() string {
	var data []byte = []byte(fmt.Sprint(r.i))

	for r = r.next; r != nil; r = r.next {
		data = fmt.Appendf(data, ".%d", r.i)
	}
	return fmt.Sprintf("%s", data)
}

func makeInts(n int) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "1")
	for i := 1; i < n; i++ {
		fmt.Fprintf(&buf, ".%d", i+1)
	}
	return buf.Bytes()
}

func (r *RecursiveInt) Scan(state fmt.ScanState, verb rune) (err error) {
	_, err = fmt.Fscan(state, &r.i)
	if err != nil {
		return
	}
	r.len++
	next := new(RecursiveInt)
	_, err = fmt.Fscanf(state, ".%v", next)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			err = nil
		}
		return
	}
	r.next = next
	return
}
