package gostd

import (
	"fmt"
)

func logCase(_case string) {
	logfln("case : %s", _case)
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("LOG ERROR: %s\n", err)
}

func log(a ...any) {
	s := fmt.Sprintln(a...)
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	fmt.Print(s)
}
func logfln(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	fmt.Print(s)
}

var integerStrings = []string{
	"", "0", "123", "-456", "+789", "0777", "0o777",
	"18446744073709551615", "18446744073709551616",
	"1_2_3", "_123", "1__23", "123_", "0_123",
	"0b110", "0b_110", "0b1_10",
	"0o123", "0o_777", "0o7_77",
	"0123", "0_777", "7_77",
	"0x123", "0x_fff", "0xf_ff",
	"0x123", "0x_FFF", "0xF_FF",
	"0X123", "0X_fff", "0Xf_ff",
	"0X123", "0X_FFF", "0XF_FF",
	"0110", "0_110", "01_10", "abc", "a_bc", "777", "7_77", "fff", "f_ff",
}
var integerStringsOther = []string{
	"abcdefg", "10", "holycow",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
}

var boolStrings = []string{
	"1", "t", "T", "true", "TRUE", "True",
	"tRue", "trUE", "TRuE",
	"0", "f", "F", "false", "FALSE", "False",
	"fAlse", "FALse", "fALsE",
}

var cmpStrings = []string{
	"i", "+i", "-i", "0i", "-0i", "+0i", "I", "-I", "(1i)", "-(1i)", "(3.0+5.5i)", "3.0+5.5i",
	"NaN", "NANi", "nan+nAni", "+NaN", "-NaN", "NaN-NaNi",
	"Inf", "+inf", "-inf", "Infinity", "+INFINITY", "-infinity", "+infi", "0-infinityi", "Inf+Infi", "+Inf-Infi", "-Infinity+Infi", "inf-inf",
	"0", "0i", "-0.0i", "0+0.0i", "0e+0i", "0e-0+0i", "-0.0-0.0i", "0e+012345", "0x0p+012345i", "0x0.00p-012345i", "+0e-0+0e-0i", "0e+0+0e+0i", "-0e+0-0e+0i",
	"0.1", "0.1i", "0.123", "0.123i", "0.123+0.123i", "99", "+99", "-99", "+1i", "-1i", "+3+1i", "30+3i", "+3e+3-3e+3i", "+3e+3+3e+3i", "+3e+3+3e+3i+",
	"0.1_2_3", "+0x_3p3i", "0_0+0x_0p0i", "0x_10.3p-8+0x3p3i", "+0x_1_0.3p-8+0x_3_0p3i", "0x1_0.3p+8-0x_3p3i",
	"0x10.3p-8+0x3p3i", "+0x10.3p-8+0x3p3i", "0x10.3p+8-0x3p3i", "0x1p0", "0x1p1", "0x1p-1", "0x1ep-1", "-0x1ep-1", "-0x2p3", "0x1e2", "1p2", "0x1e2i",
}

var floatStrings = []string{
	"",
	"1", "+1", "1x", "1.1.", "1e23", "1E23", "100000000000000000000000", "1e-100", "123456700", "99999999999999974834176", "100000000000000000000001", "100000000000000008388608", "100000000000000016777215", "100000000000000016777216",
	"-1", "-0.1", "-0", "1e-20", "625e-3", "0x1p0", "0x1p1", "0x1p-1", "0x1ep-1", "-0x1ep-1", "-0x1_ep-1", "0x1p-200", "0x1p200", "0x1fFe2.p0", "0x1fFe2.P0", "-0x2p3", "0x0.fp4", "0x0.fp0", "0x1e2", "1p2",
	"0", "0e0", "-0e0", "+0e0", "0e-0", "-0e-0", "+0e-0", "0e+0", "-0e+0", "+0e+0", "0e+01234567890123456789", "0.00e-01234567890123456789", "-0e+01234567890123456789", "-0.00e-01234567890123456789", "0x0p+01234567890123456789", "0x0.00p-01234567890123456789", "-0x0p+01234567890123456789", "-0x0.00p-01234567890123456789",
	"nan", "NaN", "NAN", "inf", "-Inf", "+INF", "-Infinity", "+INFINITY", "Infinity",
	"1.7976931348623157e308", "-1.7976931348623157e308", "0x1.fffffffffffffp1023", "-0x1.fffffffffffffp1023", "0x1fffffffffffffp+971", "-0x1fffffffffffffp+971", "0x.1fffffffffffffp1027", "-0x.1fffffffffffffp1027",
	"1.7976931348623158e308", "-1.7976931348623158e308", "0x1.fffffffffffff7fffp1023", "-0x1.fffffffffffff7fffp1023", "1.797693134862315808e308", "-1.797693134862315808e308", "0x1.fffffffffffff8p1023", "-0x1.fffffffffffff8p1023", "0x1fffffffffffff.8p+971", "-0x1fffffffffffff8p+967", "0x.1fffffffffffff8p1027",
	"-0x.1fffffffffffff9p1027", "1e", "1e-", ".e-1", "1\x00.2", "0x", "0x.", "0x1", "0x.1", "0x1p", "0x.1p", "0x1p+", "0x.1p+", "0x1p-", "0x.1p-", "0x1p+2", "0x.1p+2", "0x1p-2", "0x.1p-2",
}

var bases = []int{
	0, 2, 8, 10, 16, 3, 17, 20, 25, 29, 35, 36,
}

var quotes = []string{
	"quoted \u263a string, 你好",
	`quoted \u263a string, 你好`,
	"\"quoted \u263a string, 你好\"",
	"`quoted \u263a string, 你好`",
	`"quoted \u263a string, 你好"`,
	`\"quoted \u263a string, 你好\"`,
	`\'quoted \u263a string, 你好\'`,
	"'\u263a'",
	`'\u263a'`,
	"'quoted \u263a string, 你好'",
}
