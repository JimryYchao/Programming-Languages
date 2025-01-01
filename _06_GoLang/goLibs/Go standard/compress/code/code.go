package gostd

import (
	"fmt"
)

func _logCase(_case string) {
	_logfln("case : %s", _case)
}

var fileNames = []string{
	"e.txt",
	"compress.pdf",
	"Isaac.Newton-Opticks.txt",
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("LOG ERROR: %s\n", err)
}

func _log(a ...any) {
	s := fmt.Sprintln(a...)
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	fmt.Print(s)
}
func _logfln(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	fmt.Print(s)
}
