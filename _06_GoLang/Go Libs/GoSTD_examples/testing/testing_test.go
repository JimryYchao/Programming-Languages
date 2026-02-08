package gostd_testing

/* 包 testing 为 Go 包的自动化测试提供了支持
! 测试主函数
func TestMain(m *testing.M)
! 常规测试：go test -run=^Test
func TestXxx(t *testing.T)
! 基准测试：go test -run=NONE -bench=Benchmark
func BenchmarkXxx(b *testing.B)
! 模糊测试：go test -run=NONE -fuzz=Fuzz
func FuzzXxx(f *testing.F)
! 运行示例：go test -run=Example
func ExampleIdentifier() { ... }  //? for Identifier
func Example() { ... }      	  //? for package
func ExampleF() { ... }           //? for function F
func ExampleT() { ... }           //? for type T
func ExampleT_M() { ... }         //? for method T.M
func ExampleXxx_suffix() { ... }  //? for Xxx different examples
! 子测试：go test -run TestXxx/subTest
! go test -run /subTest, 测试全部顶级测试的子测试 subTest
func TestXxx(t *testing.T) {
	t.Run("subTest", func(t *testing.T) { ... })
}
*/

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"testing"
	"time"
)

// 辅助声明
func doSomething(up int) []int {
	if up < 0 {
		return nil
	}
	// random generate a slice
	src := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	s := rand.Perm(src.Int() % up)
	// reverse slice
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

var Xxx = "Testing ExampleXxx"

// ! TestMain(m *testing.M) 测试主函数
// ? go test -v -run=^TestMain$
func TestMain(m *testing.M) {
	fmt.Println("Before testing.")
	exitCode := m.Run() // 运行测试
	fmt.Println("After testing.")
	os.Exit(exitCode)
}

// ! TestXxx(t *testing.T) 测试
// ? go test -v -run=Test
func TestXxx(t *testing.T) {
	t.Log("Testing `TestXxx`")
}

// ! BenchmarkXxx(b *testing.B) 基准测试
// ? go test -run=NONE -bench=^BenchmarkXxx$ [-benchtime=3s] [-benchmem]
func BenchmarkXxx(b *testing.B) {
	// 迭代足够次数或运行足够时长来计算单次操作 body 使用的近似时间
	for range b.N {
		// body
		doSomething(1000)
	}

	b.Run("b.Loop", func(b *testing.B) {
		doSomething(1000)
	})

	b.Run("b.RunParallel", func(b *testing.B) {
		templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
		b.RunParallel(func(pb *testing.PB) {
			var buf bytes.Buffer
			for pb.Next() {
				buf.Reset()
				templ.Execute(&buf, "World")
			}
		})
	})
}

// ! FuzzXxx(f *testing.F) 模糊测试
// ? go test -v -run=NONE -fuzz=^FuzzXxx$ [-parallel=8] [-fuzztime=5s] [-short]
func FuzzXxx(f *testing.F) {
	for _, seed := range [][]byte{{}, {0}, {9}, {0xa}, {0xf}, {1, 2, 3, 4}} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		enc := hex.EncodeToString(in)
		out, err := hex.DecodeString(enc)
		if err != nil {
			t.Fatalf("%v: decode: %v", in, err)
		}
		if !bytes.Equal(in, out) {
			t.Fatalf("%v: not equal after round trip: %v", in, out)
		}
	})
}

// ! ExampleXXX() 示例测试
// ? go test -v -run=Example
func ExampleXxx() { // for Identifier Xxx
	fmt.Println(Xxx)
	// output: Testing ExampleXxx
}

func Example() { // for package
	fmt.Println("Testing unordered output")
	for k, v := range map[int]string{1: "a", 2: "b", 3: "c", 4: "d"} {
		fmt.Println(k, v)
	}
	// unordered output:
	// Testing unordered output
	// 4 d
	// 3 c
	// 2 b
	// 1 a
}
