package gostd

import (
	"fmt"
	"time"
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

func log(a any) {
	s := fmt.Sprint(a)
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

var dates = []time.Time{
	time.Date(2008, 9, 17, 20, 4, 26, 0, time.UTC),
	time.Date(1994, 9, 17, 20, 4, 26, 0, time.FixedZone("EST", -18000)),
	time.Date(2000, 12, 26, 1, 15, 6, 0, time.FixedZone("OTO", 15600)),
}

var layouts = []struct {
	name   string
	layout string
}{
	{"Layout", time.Layout},
	{"ANSIC", time.ANSIC},
	{"UnixDate", time.UnixDate},
	{"RubyDate", time.RubyDate},
	{"RFC822", time.RFC822},
	{"RFC822Z", time.RFC822Z},
	{"RFC850", time.RFC850},
	{"RFC1123", time.RFC1123},
	{"RFC1123Z", time.RFC1123Z},
	{"RFC3339", time.RFC3339},
	{"RFC3339Nano", time.RFC3339Nano},
	{"Kitchen", time.Kitchen},
	{"Stamp", time.Stamp},
	{"StampMilli", time.StampMilli},
	{"StampMicro", time.StampMicro},
	{"StampNano", time.StampNano},
	{"DateTime", time.DateTime},
	{"DateOnly", time.DateOnly},
	{"TimeOnly", time.TimeOnly},
	{"YearDay", "Jan  2 002 __2 2"},
	{"Year", "2006 6 06 _6 __6 ___6"},
	{"Month", "Jan January 1 01 _1"},
	{"DayOfMonth", "2 02 _2 __2"},
	{"DayOfWeek", "Mon Monday"},
	{"Hour", "15 3 03 _3"},
	{"Minute", "4 04 _4"},
	{"Second", "5 05 _5"},
}

var tms = []string{
	"06/01 12:04:57AM '24 +0800",
	"Sat Jun  1 00:04:57 2024",
	"Sat Jun  1 00:04:57 CST 2024",
	"Sat Jun 01 00:04:57 +0800 2024",
	"01 Jun 24 00:04 CST",
	"01 Jun 24 00:04 +0800",
	"Saturday, 01-Jun-24 00:04:57 CST",
	"Sat, 01 Jun 2024 00:04:57 CST",
	"Sat, 01 Jun 2024 00:04:57 +0800",
	"2024-06-01T00:04:57+08:00",
	"2024-06-01T00:04:57.9328875+08:00",
	"12:04AM",
	"Jun  1 00:04:57",
	"Jun  1 00:04:57.932",
	"Jun  1 00:04:57.932887",
	"Jun  1 00:04:57.932887500",
	"2024-06-01 00:04:57",
	"2024-06-01",
	"00:04:57",
	"Jun  1 153 153 1",
	"2024 6 24 _6 __6 ___6",
	"Jun June 6 06 _6",
	"1 01  1 153",
	"Sat Saturday",
	"00 12 12 _12",
	"4 04 _4",
	"57 57 _57",
}

var locations = []string{
	"Asia/Baghdad",
	"America/Blanc-Sablon",
	"Australia/Sydney",
}

var parseDuration = []string{
	"0",
	"5s",
	"30s",
	"1478s",
	"-5s",
	"+5s",
	"-0",
	"+0",
	"5.0s",
	"5.6s",
	"5.s",
	".5s",
	"1.0s",
	"1.00s",
	"1.004s",
	"1.0040s",
	"100.00100s",
	"10ns",
	"11us",
	"12µs",
	"12μs",
	"13ms",
	"14s",
	"15m",
	"16h",
	"3h30m",
	"10.5s4m",
	"-2m3.4s",
	"1h2m3s4ms5us6ns",
	"39h9m14.425s",
	"52763797000ns",
	"0.3333333333333333333h",
	"9007199254740993ns",
	"9223372036854775807ns",
	"9223372036854775.807us",
	"9223372036s854ms775us807ns",
	"-9223372036854775808ns",
	"-9223372036854775.808us",
	"-9223372036s854ms775us808ns",
	"-9223372036854775808ns",
	"-2562047h47m16.854775808s",
	"0.100000000000000000000h",
	"0.830103483285477580700h",
}
var ds = []time.Duration{
	time.Nanosecond,
	time.Microsecond,
	time.Millisecond,
	time.Second,
	2 * time.Second,
	time.Minute,
	10 * time.Minute,
	time.Hour,
}
