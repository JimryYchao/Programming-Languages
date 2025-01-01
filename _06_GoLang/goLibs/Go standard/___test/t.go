package main

import (
	"testing"
)

// func main() {
// 	flag.Int("n", 1234, "help message for flag n")
// 	flag.String("s", "flag", "help message for flag s")
// 	flag.Usage()
// }

/*
Usage of test.exe:
  -n int
        help message for flag n (default 1234)
  -s string
        help message for flag s (default "flag")
*/

func TestS(t *testing.T) {
	// r := reflect.Zero(reflect.TypeFor[io.Reader]())
	// fmt.Print(r, r.Type(), r.Kind())

	// r = reflect.ValueOf(r)
	// fmt.Print(r, r.Type(), r.Kind(), r.Interface().(reflect.Value).Type())

}

func main() {
	const (
		a  = 1 << iota //  1 << 0 (iota = 0)
		a1 = 2         //  (iota == 1, unused)
		b  = 2 * iota  //  2 * 2  (iota = 2)
		b1             //  2 * 3  (iota = 3)
		b2             //  2 * 4  (iota = 4)
		b3             //  2 * 5  (iota = 5)
		c  = 1 + iota  //  1 + 6  (iota = 6)
		c1             //  1 + 7  (iota = 7)
		c2             //  1 + 8  (iota = 8)
	)
	print(b2)

	var s []rune // nil
	var s1 []int = make([]int, 1024)
	var a = [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9} // array
	s2 := a[0:5]
}
