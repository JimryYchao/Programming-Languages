package examples

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s: %d", p.Name, p.Age)
}

// People implements sort.Interface for []Person based on the Age field.
type People []Person

func (a People) Len() int           { return len(a) }
func (a People) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a People) Less(i, j int) bool { return a[i].Age < a[j].Age }

//----------------------------------------------------------------------------
