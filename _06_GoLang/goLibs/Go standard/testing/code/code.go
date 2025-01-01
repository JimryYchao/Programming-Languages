package gostd

import (
	"errors"
	"fmt"
	"math/rand"
	"unicode/utf8"
)

func logCase(_case string) {
	logfln("case : %s", _case)
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

var Identifier string = "Hello World"

type T string

func (t T) M() {
	fmt.Println(t)
}

func F() string {
	return Identifier
}

var M = map[int]string{
	1: "a", 2: "b", 3: "c", 4: "d",
}

func Abs(i int) int {
	if i < 0 {
		i = -i
	}
	return i
}

func Rand() int {
	return rand.Int()
}

type Big struct {
	r int
}

func (b *Big) Do() {
	b.r = rand.Int()
}

func NewBig() Big {
	return Big{r: 10086}
}

func Println(s ...any) {
	fmt.Println(s...)
}

func Perm(n int) []int {
	if n < 0 {
		return nil
	}
	s := make([]int, n+1)
	for i := range len(s) {
		s[i] = i
	}
	return s
}

func checkErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("LOG ERROR: %s\n", err)
}

func Reverse(s string) (string, error) {
	if !utf8.ValidString(s) {
		return s, errors.New("input is not valid UTF-8")
	}
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r), nil
}
