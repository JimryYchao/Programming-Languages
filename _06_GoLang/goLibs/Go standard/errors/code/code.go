package gostd

import (
	"fmt"
)

func checkErr(err error) {
	if err == nil {
		return
	}
	logf("LOG ERROR: \n%s", err)
}
func logCase(_case string) {
	logf("case : %s", _case)
}
func logf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}
