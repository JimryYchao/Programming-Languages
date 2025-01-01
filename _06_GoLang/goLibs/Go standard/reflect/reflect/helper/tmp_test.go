package helper

import (
	"fmt"
	"reflect"
	"testing"
)

var log = fmt.Println
var logf = fmt.Printf

func Test(t *testing.T) {
	type S struct {
		v int
	}
	s := reflect.ValueOf([]int{1, 2, 3, 4, 5, 6, 7, 8})

	s.SetLen(5) // err

}
