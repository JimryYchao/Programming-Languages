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

var Log = log
