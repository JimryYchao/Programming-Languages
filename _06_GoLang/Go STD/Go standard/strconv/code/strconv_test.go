package gostd

import (
	"context"
	"math"
	. "strconv"
	"testing"
)

//? go test -v -run=^$

/* to types
! Atoi == ParseInt(s，10，0)
! Parses:
	ParseBool, ParseComplex, ParseFloat, ParseInt, ParseUint
*/

func TestParse(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	getSw := func(t *testing.T, name string) *StringWriter {
		sw, err := NewStringWriter(ctx, name, true)
		if err != nil {
			t.Fatal(err)
		}
		return sw
	}

	t.Run("Atoi", func(t *testing.T) {
		sw := getSw(t, "Atoi")
		sw.WriteLine(">>> Test string to integer <<<")
		atoi := func(s string) {
			if i, err := Atoi(s); err == nil {
				sw.WriteLine("Atoi %s to %d", s, i)
			} else {
				nerr := err.(*NumError)
				sw.WriteError("%s : \"%s\"", nerr.Unwrap(), nerr.Num)
			}
		}
		for _, v := range integerStrings {
			atoi(v)
		}
	})

	t.Run("ParseBool", func(t *testing.T) {
		sw := getSw(t, "ParseBool")
		sw.WriteLine(">>> Test string to bool <<<")
		parseBool := func(s string) {
			if b, err := ParseBool(s); err == nil {
				sw.WriteLine("ParseBool %s to %t", s, b)
			} else {
				nerr := err.(*NumError)
				sw.WriteError("%s : \"%s\"", nerr.Unwrap(), nerr.Num)
			}
		}
		for _, s := range boolStrings {
			parseBool(s)
		}

	})

	t.Run("ParseComplex", func(t *testing.T) {
		sw := getSw(t, "ParseComplex")
		sw.WriteLine(">>> Test string to complex <<<")

		parseComplex := func(s string) {
			if c, err := ParseComplex(s, 128); err == nil {
				sw.WriteLine("ParseComplex %s to %v", s, c)
			} else {
				nerr := err.(*NumError)
				sw.WriteError("%s : \"%s\"", nerr.Unwrap(), nerr.Num)
			}
		}

		for _, s := range cmpStrings {
			parseComplex(s)
		}
		for _, s := range integerStrings {
			parseComplex(s)
		}
	})

	t.Run("ParseFloat", func(t *testing.T) {
		sw := getSw(t, "ParseFloat")
		sw.WriteLine(">>> Test string to float <<<")

		parseFloat := func(s string) {
			if c, err := ParseFloat(s, 64); err == nil {
				sw.WriteLine("ParseFloat %s to %v", s, c)
			} else {
				nerr := err.(*NumError)
				sw.WriteError("%s : \"%s\"", nerr.Unwrap(), nerr.Num)
			}
		}

		for _, s := range cmpStrings {
			parseFloat(s)
		}
		for _, s := range integerStrings {
			parseFloat(s)
		}
		for _, s := range floatStrings {
			parseFloat(s)
		}
	})
	t.Run("ParseInt", func(t *testing.T) {
		sw := getSw(t, "ParseInt")
		sw.WriteLine(">>> Test string to Int <<<")

		parseInt := func(s string, base int) {
			if c, err := ParseInt(s, base, 64); err == nil {
				sw.WriteLine("ParseInt %s to %v", s, c)
			} else {
				nerr := err.(*NumError)
				sw.WriteError("%s in base(%d) : \"%s\"", nerr.Unwrap(), base, nerr.Num)
			}
		}
		for _, base := range bases {
			sw.WriteLine("========= Base %d =========", base)
			for _, s := range integerStrings {
				parseInt(s, base)
			}
			for _, s := range integerStringsOther {
				parseInt(s, base)
			}
		}
	})
}

/* to string
! Itoa == FormatInt(int64(i), 10)
! Formats:
	FormatBool, FormatComplex, FormatFloat, FormatInt, FormatUint
! FormatComplex, FormatFloat 的 fmt 参数可以是 "b", "e", "E", "f", "g", "G", "x", "X"; prec 表示有效位数的最大值（删除尾随零），-1 表示使用最少的数字
*/

func TestFormat(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	t.Run("FormatFloat", func(t *testing.T) {
		sw, err := NewStringWriter(ctx, "FormatFloat", false)
		if err != nil {
			t.Fatal(err)
		}
		sw.WriteLine(">>> Test float to string <<<")
		fmt := []byte{'b', 'e', 'E', 'f', 'g', 'G', 'x', 'X'}
		prec := []int{-1, 0, 1, 3, 6, 10, 100}
		for _, flt := range []float64{
			0, -0, 1, 1.2, 1.000, 1.23000,
			math.Pi, math.Inf(1), math.Inf(-1), math.Inf(0), math.NaN(),
		} {
			for _, f := range fmt {
				for _, p := range prec {
					sw.WriteLine("%10f in (%c,%3d): %s", flt, f, p, FormatFloat(flt, f, p, 64))
				}
			}
		}
	})
}

/*
! Quotes:
	Quote, QuoteToASCII, QuoteToGraphic
	QuoteRune, QuoteRuneToASCII, QuoteRuneToGraphic
! QuotedPrefix 返回 s 中以带引号（",`）的字符串为前缀的片段, 单引号视为字符
! IsPrint 报告 r 是否被 Go 定义为可打印的，定义与 unicode.IsPrint 相同。IsPrint：字母，数字，标点符号，符号和 ASCII 空格
! IsGraphic 报告 r 是否被 Unicode 定义为图形。这些字符包括字母、标记、数字、标点符号、符号和空格，来自类别 L、M、N、P、S 和 Z。
! CanBackquote 报告字符串 s 是否可以不加修改地表示为不包含除 \t 以外的控制字符的单行反引号字符串
*/

func TestQuote(t *testing.T) {
	t.Run("Quote", func(t *testing.T) {
		for _, q := range quotes {
			logfln("(%s) Quote   = %s", q, Quote(q))
			logfln("(%s) ToASCII = %s", q, QuoteToASCII(q))
			logfln("(%s) Graphic = %s", q, QuoteToGraphic(q))
		}
	})

	t.Run("QuotedPrefix", func(t *testing.T) {
		quotedPrefix := func(s string) {
			s, err := QuotedPrefix(s)
			logfln("%q, %v\n", s, err)
		}

		quotedPrefix("“not” a quoted string")
		quotedPrefix("\"double-quoted string\" with trailing text")
		quotedPrefix("`or backquoted` with more trailing text")
		quotedPrefix("'\u263a' is also okay")

		// "", invalid syntax
		// "\"double-quoted string\"", <nil>
		// "`or backquoted`", <nil>
		// "'☺'", <nil>
	})

	t.Run("CanBackquote", func(t *testing.T) {
		quote := func(str string) {
			if CanBackquote(str) {
				logfln(Quote(str))
			}
		}
		quote("\tFran & Freddie's Diner ☺")
		quote("\t\rcan't backquote this")
	})
}

/*
! Unquote 将 s 解释为单引号、双引号或反引号的 Go 字符串字面量
! UnquoteChar 解释字符串的第一个字符或字节
	若 `quote` 为 '\'' 或 '\"' 则允许解释 \' 或 \"
	设置为其他值，则表示不解释任何转义字符，仅解释 ' 或 "
*/

func TestUnquote(t *testing.T) {
	t.Run("Unquote", func(t *testing.T) {
		for _, q := range quotes {
			if s, err := Unquote(q); err == nil {
				logfln("(%s) Unquote = %s", q, s)
			} else {
				logfln("(%s) Error: %s", q, err)
			}
		}
	})

	t.Run("UnquoteChar", func(t *testing.T) {
		for _, s := range quotes {
			logfln(">>> %s <<<", s)
			for _, q := range []byte{0, '`', '\'', '"', '\n', 'q', '\\', '\t'} {
				if v, m, t, err := UnquoteChar(s, q); err == nil {
					logfln("(%30s) Unquote(%#6q) = value: %c, mul: %t, tail: %s", s, q, v, m, t)
				} else {
					logfln("(%30s)   Error(%#6q) = %s", s, q, err)
				}
			}
			log(" ")
		}
	})
}

/*
! Appends:
	AppendBool, AppendFloat, AppendInt, AppendUint
	AppendQuote, AppendQuoteToASCII, AppendQuoteToGraphic
	AppendQuoteRune, AppendQuoteRuneToASCII, AppendQuoteRuneToGraphic
Appends 和 Quote 类似，但是附加到目标切片上
*/

func Test(t *testing.T) {
	b := []byte("quote (ascii):")
	b = AppendQuoteToASCII(b, `"Fran & Freddie's Diner, 你好"`)
	log(string(b))
}
