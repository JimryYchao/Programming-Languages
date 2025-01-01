package gostd

// 包 testing 为 Go 包的自动化测试提供了支持

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

// ! TestMain 为测试主入口 (可选)
// ? go test [-v]
func TestMain(m *testing.M) {
	//? Run 运行测试，返回一个可以传递给 os.Exit 的退出代码
	exitCode := m.Run()

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
	//? 迭代足够次数或运行足够时长来计算单次操作 body 使用的近似时间
	for i := 0; i < b.N; i++ {
		// body
		doSomething(100)
	}
}

// ! FuzzXxx(f *testing.F) 模糊测试
// ? go test -v -run=NONE -fuzz=^FuzzXxx$ [-parallel=8] [-fuzztime=5s] [-short]
func FuzzXxx(f *testing.F) {
	mlimit := debug.SetMemoryLimit(1000 * 100)
	mTreads := debug.SetMaxThreads(1000)
	f.Cleanup(func() {
		debug.SetMemoryLimit(mlimit)
		debug.SetMaxThreads(mTreads)
	})

	f.Log("Fuzzing `FuzzXxx`")
	added := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	for _, v := range added {
		//? Add 添加种子用料
		f.Add(v)
	}
	//? Fuzz 运行模糊目标函数 ff
	f.Fuzz(func(t *testing.T, up int) {
		if up < 2 {
			t.SkipNow()
			return
		}
		// else i > 0, do something
		doSomething(up)
	})
}

// ! ExampleXXX() 示例测试
// ? go test -v -run=Example
func ExampleIdentifier() { // for Identifier
	Println(Identifier)
	// output: Hello World
}
func Example() { // for package
	Println("Testing unordered output")
	for k, v := range M {
		Println(k, v)
	}
	// unordered output:
	// Testing unordered output
	// 4 d
	// 3 c
	// 2 b
	// 1 a
}
func ExampleF() { // for function F
	Println(F())
	// output: Hello World
}
func ExampleT() { // for type T
	var t T = "Testing T"
	Println(t)
	// output: Testing T
}
func ExampleT_M() { // for method T.M
	T.M("Testing T_M")
	// output: Testing T_M
}

// ! AllocsPerRun 函数返回调用 `f` 期间分配的平均次数 (allocs/op)。
// ? go test -v -run=^TestAllocsPerRun$
func TestAllocsPerRun(t *testing.T) {
	var allocsPerRunTests = []struct {
		name   string
		fn     func(temp *any)
		allocs float64
	}{
		{"alloc *byte", func(temp *any) { *temp = new(*byte) }, 1},
		{"alloc complex128", func(temp *any) { *temp = new(complex128) }, 1},
		{"alloc float64", func(temp *any) { *temp = new(float64) }, 1},
		{"alloc int32", func(temp *any) { *temp = new(int32) }, 1},
		{"alloc byte", func(temp *any) { *temp = new(byte) }, 1},
	}
	var temp any
	for _, tt := range allocsPerRunTests {
		if allocs := testing.AllocsPerRun(100, func() { tt.fn(&temp) }); allocs != tt.allocs {
			t.Errorf("AllocsPerRun(100, %s) = %v, want %v.", tt.name, allocs, tt.allocs)
		}
	}
}

// ! CoverMode 报告在测试的覆盖模式: set, count or atomic。
// ! Coverage 报告当前代码的覆盖率, 它不能替代 `go test -cover` 和 `go tool cover` 生成的报告。
// ? go test -v -cover -run=^TestCoverage$ [-coverprofile='coverage.out']
// ? go tool cover -html='coverage.out'
func TestCoverage(t *testing.T) {
	t.Cleanup(func() {
		// 标志 -cover 启用代码覆盖率分析, 默认为 set 模式
		coverMode := testing.CoverMode() // 未设置时返回 ""
		if coverMode == "" {
			t.Log("-cover is not enabled.")
			return
		}
		t.Logf("Cover Mode: %s, Coverage in %s is %v.", coverMode, t.Name(), testing.Coverage())
	})
	// call another file
	Perm(-100)
}

// ! Short 报告在测试中是否启用 `-test.short`。
// ! Verbose 报告在测试中是否启用 `-test.v`。
// ! Testing 报告当前代码是否在测试中运行。
// ? go test -v -short -run=^TestTesting$
func TestTesting(t *testing.T) {
	if !testing.Testing() {
		t.Errorf("%s runs only in programs created by `go test`.", t.Name())
	}
	if testing.Verbose() {
		t.Log("Enable -v mode.")
	}
	if !testing.Short() {
		t.Skipf("Skipping %s in non-short mode.", t.Name())
	}
}

/*
! Init 初始化注册测试标志，通常由 “go test” 命令自动注册。不使用指令单独调用 Benchmark 等函数时才需要 Init。
! Benchmark 基准测试对单个函数进行基准测试（不依赖 `go test`）。
! BenchmarkResult 包含基准测试运行的结果。
	N          			迭代总次数
	T          			运行基准测试的总时间
	Bytes      			单次迭代处理的字节数
	MemAllocs  			内存分配总次数
	MemBytes   			分配的总字节数
	Extra   			记录由 ReportMetric 添加的用户度量标准

	AllocedBytesPerOp	返回 B/op 指标，其计算公式为 r.MemBytes/r.N
	AllocsPerOp			返回 allocs/op 指标，计算公式为 r.MemAllocs/r.N
	MemString			以与 `go test` 相同的格式返回 r.AllocedBytesPerOp 和 r.AllocsPerOp
	NsPerOp				返回 ns/op 度量
	String				返回基准测试结果的摘要。额外指标覆盖相同名称的内置指标; String 不包括 allocs/op 或 B/op
*/
//? go test -v -run=^TestBenchmark$
func TestBenchmark(t *testing.T) {
	testBenchmark()
}
func testBenchmark() {
	testing.Init() // Init before any Benchmark without `go test`
	var br testing.BenchmarkResult = testing.Benchmark(func(b *testing.B) {
		b.ReportAllocs()
		var s []int64
		var compares int64
		for range b.N {
			s = []int64{5, 4, 3, 2, 1}
			sort.Slice(s, func(i, j int) bool {
				compares++
				return s[i] < s[j]
			})
		}
		b.SetBytes(int64(len(s)) * int64(unsafe.Sizeof(int64(0))))
		b.ReportMetric(float64(compares)/float64(b.N), "compares/op")
		b.ReportMetric(float64(compares)/float64(b.Elapsed().Milliseconds()), "compares/ms")
	})

	fmt.Println("testBenchmark\t" + br.String() + "\t" + br.MemString())
}

/*
! testing.TB 是 `testing.T`、`testing.B` 和 `testing.F` 的公共接口
	Log         Println；对于测试，仅当测试失败或设置了 `—test. v` 标志时，才会打印文本。对于基准测试，总是打印文本
	Logf   		Printf，自动换行
	Fail		将函数标记为失败，但继续执行
	Error		相当于Log + Fail
	Errorf		相当于 Logf + Fail
	FailNow		将函数标记为失败，并停止其执行
	Fatal       相当于 Log + FailNow
	Fatalf		相当于 Logf + FailNow
	Failed		报告函数是否失败

	SkipNow		将测试函数标记为已跳过，并停止其执行
	Skip		相当于 Log + SkipNow
	Skipf		相当于 Logf + SkipNow
	Skipped     报告是否跳过测试

	Cleanup		注册一个在（子）测试完成后要调用的函数。测试结束时，函数按注册逆序调用
	Helper 		将调用方标记为测试帮助函数。当打印文件和行信息时，该函数将被跳过
	Name		返回正在运行的（子）测试或基准的名称
	Setenv		调用 os.Setenv(key，value)，并在测试后使用 Cleanup 还原。不能用于并行测试
	TempDir		返回一个临时目录供测试使用。（子）测试完成后自动删除
*/
// ? go test -v -run=^TestTBFCommonFunctions$
func TestTBFCommonFunctions(t *testing.T) {
	beforeTest := func(t *testing.T) {
		t.Helper() // 标记测试帮助函数
		t.Logf("Testing >>> %s", t.Name())
		// 测试结束时按注册逆序依次调用
		t.Cleanup(func() {
			t.Logf("End Test >>> %s", t.Name())
		})
	}
	t.Run("Fail", func(t *testing.T) {
		beforeTest(t)
		t.Cleanup(func() {
			if t.Skipped() {
				t.Log("Please use command: `go test -v -run=^TestTB$`.")
			}
		})
		var _sub *testing.T
		if runflag := flag.Lookup("test.run"); runflag != nil && runflag.Value.String() == "^TestTB$" {
			// 仅运行当前测试时故意失败
			t.Run("-run=^TestTB$", func(t *testing.T) {
				_sub = t
				t.Fatal("TestTB/Fail fails deliberately.")
			})
		} else {
			t.Skip("TestTB/Fail runs only in `-run=^TestTB$` mode.")
		}
		if _sub != nil && _sub.Failed() {
			t.Logf("%s Pass", t.Name())
		}
	})
}

// ? go test -v -run=^TestTBSetenv&
func TestTB_Setenv(t *testing.T) {
	env := struct {
		key string
		val string
	}{
		key: "Setenv", val: "Hello World",
	}
	if v, ok := os.LookupEnv(env.key); ok {
		t.Logf("The value of Env `Setenv` is %s", v)
	} else {
		t.Setenv(env.key, env.val)
		if os.Getenv(env.key) == "Hello World" {
			t.Log("Set environment variable `Setenv:Hello World` successfully.")
		}
	}

	panickingRecover := func(name string) {
		if got := recover(); got != nil {
			t.Logf("panicking in %s: \n%#v.", name, got)
		} else {
			t.Logf("Test: %s PASS.", name)
		}
	}
	//? 测试完成后，将恢复 Setenv 进入（子）测试之前的值
	t.Run("RestoreWhileTestCompleted", func(t *testing.T) {
		t.Setenv(env.key, "dlroW olleH")
		t.Logf("Change env %s to %s", env.key, os.Getenv(env.key))
	})
	if os.Getenv(env.key) == env.val {
		t.Log("The environment variable `Setenv:Hello World` restores.")
	}

	//? Setenv 不能用于并行测试或具有并行祖先的测试
	t.Run("ParallelAfterSetenv", func(t *testing.T) {
		defer panickingRecover(t.Name())
		t.Setenv("Setenv", "Hello")
		t.Parallel()
	})

	t.Run("ParallelBeforeSetenv", func(t *testing.T) {
		defer panickingRecover(t.Name())
		t.Parallel()
		t.Setenv("Setenv", "Hello")
	})

	t.Run("ParallelParentBeforeSetenv", func(t *testing.T) {
		t.Parallel()
		t.Run("child", func(t *testing.T) {
			defer panickingRecover(t.Name())
			t.Setenv("Setenv", "Hello")
		})
	})
}

// ? go test -v -run=^TestTBTempDir$
func TestTB_TempDir(t *testing.T) {
	var dir string
	fn_checkDir := func(dir string) {
		fi, err := os.Stat(dir)
		if fi != nil {
			t.Fatalf("Directory %q from user Cleanup still exists", dir)
		}
		if !os.IsNotExist(err) {
			t.Fatalf("Unexpected error: %v", err)
		}
	}

	t.Run("CheckExist", func(t *testing.T) {
		dirCh := make(chan string, 1)
		t.Cleanup(func() {
			// 验证目录 directory 已经在测试完成时被删除
			select {
			case dir := <-dirCh:
				fn_checkDir(dir)
			default:
				if !t.Failed() {
					t.Fatal("never received dir channel")
				}
			}
		})
		dir := t.TempDir()
		t.Logf("create a tempDir=%v\n", dir)
		dirCh <- dir // 传递 tempDir
	})

	t.Run("InCleanup", func(t *testing.T) {
		t.Helper()
		t.Run("test", func(t *testing.T) {
			t.Cleanup(func() {
				dir = t.TempDir()
			})
			_ = t.TempDir()
		})
		fn_checkDir(dir)
	})

	t.Run("InBenchmark", func(t *testing.T) {
		testing.Benchmark(func(b *testing.B) {
			if !b.Run("test", func(b *testing.B) {
				// Add a loop so that the test won't fail.
				for i := 0; i < b.N; i++ {
					_ = b.TempDir()
				}
			}) {
				t.Fatal("Sub test failure in a benchmark")
			}
		})
	})
}

/*
! T.non-common Functions
	Deadline	报告测试的二进制文件的运行时间将超过 `-timeout` 指定的时间（默认为 10m）。`-timeout=0s` 表示无超时，`ok` 始终返回 false
	Parallel	标记并行信号，表示此测试将与其他的具有并行信号的测试并行运行
	Run			在一个单独的 goroutine 中运行子测试 f(t *T)
*/
//? go test -v -run=^TestT_Functions$ [-timeout=0s]
func TestT_Functions(t *testing.T) {
	t.Run("NoParallel", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 2; i++ {
			i := i
			t.Run(fmt.Sprintf("outer%d", i), func(t *testing.T) {
				for j := 0; j < 2; j++ {
					j := j
					t.Run(fmt.Sprintf("inner%d", j), func(t *testing.T) {
						time.Sleep(2 * time.Second)
						t.Logf("End NoParallel/outer%d/inner%d", i, j)
					})
				}
				t.Logf("End NoParallel/outer%d", i)
			})
		}
	})

	t.Run("Parallel", func(t *testing.T) {
		t.Parallel()
		for i := 0; i < 2; i++ {
			i := i
			t.Run(fmt.Sprintf("outer%d", i), func(t *testing.T) {
				t.Parallel()
				for j := 0; j < 2; j++ {
					j := j
					t.Run(fmt.Sprintf("inner%d", j), func(t *testing.T) {
						t.Parallel()
						time.Sleep(2 * time.Second)
						t.Logf("End Parallel/outer%d/inner%d", i, j)
					})
				}
				t.Logf("End Parallel/outer%d", i)
			})
		}
	})

	t.Run("Deadline", func(t *testing.T) {
		t.Parallel()
		t.Logf("Running Time : %#v", time.Now())
		//? 获取测试当前函数的距离截止时间, false 表示无超时 (-timeout=0s)
		if dl, ok := t.Deadline(); !ok {
			t.Skipf("Testing deadline not set.")
		} else {
			remain := time.Until(dl) // 计算距离 deadline 的剩余时间
			if remain < 0 {          // 超时终止测试
				t.Fatal("Test has exceeded the deadline.")
			}

			t.Logf("Test sleeps for 5 seconds.")
			for i := range 6 {
				t.Logf("Test resumes after %d seconds.", 5-i)
				time.Sleep(1 * time.Second)
			}
		}
	})
}

/*
! B.non-common Functions
	Elapsed			返回基准测试的目前的运行时间
	StartTimer		开始计时测试。此函数在基准测试开始之前自动调用，但也可用于在调用 B.StopTimer 之后恢复计时
	StopTimer		停止测试计时。例如在执行部分无需测量的复杂初始化时，暂停计时器
	ResetTimer		将基准运行时间和内存分配计数器归零，并删除用户报告的度量。它不会影响计时器是否正在运行

	Run				运行基准测试 f 子基准测试
	RunParallel		并行运行基准测试。它创建多个 goroutine 并在它们之间分配 b.N 迭代
!   PB.Next			PB 由 RunParallel 用于运行并行基准测试。Next 报告是否有更多的迭代要执行。
	SetMaxelism 	将 B.RunParallel 使用的 goroutine 数量设置为 p * GOMAXPROCS。对于 CPU 限制的基准测试，通常不需要调用 SetValuelism

	SetData			记录单个操作中处理的字节数。如果调用此函数，基准测试将报告 ns/op 和 MB/s
	ReportAllocs	为此基准测试启用 malloc 统计信息。它等效于设置 `-test.benchmem`，但它只影响调用 ReportAllocs 的基准函数

	ReportMetric	将用户基准测试度量标准 “n 单位” 添加到报告的基准测试结果中。
					如果度量是 per/op，调用者应该除以 b.N；如果度量是 per/ns，应除以 b.Elapsed；
					ReportMetric 覆盖同一单位以前报告的任何值。如果 unit 为空字符串或 unit 包含任何空格，ReportMetric 将 panic。
					如果 unit 是一个通常由基准框架本身报告的单位（例如 “allocs/op”），ReportMetric 将覆盖该度量。
					将 “ns/op” 设置为 0 将抑制该内置度量。
*/
//? go test -v -run=NONE -bench=^Benchmark_Functions$ [-cpu='1,2,4,8']
func Benchmark_Functions(b *testing.B) {
	b.Helper()

	b.Run("BenchmarkTimer", func(b *testing.B) {
		// Benchmark sleeps for 1 seconds
		time.Sleep(1 * time.Second)
		b.ResetTimer()

		// Benchmark stop Timer
		b.StopTimer()
		t1 := b.Elapsed()

		// Benchmark sleeps for 1 seconds
		time.Sleep(1 * time.Second)
		if b.Elapsed() != t1 {
			b.Fatalf("Timer is not stopped.")
		}
		// Benchmark resume Timer
		b.StartTimer()
		time.Sleep(1 * time.Second) // ≈ 1E9 ns/op
	})

	b.Run("ReportAllocs", func(b *testing.B) {
		b.ReportAllocs() // 对该子测试报告 `-benchmem`
		for range b.N {
			doSomething(100)
		}
	})

	b.Run("RunParallel", func(b *testing.B) {
		if testing.Short() {
			b.Skip("Skipping in short mode")
		}
		var procs atomic.Uint32
		var iters atomic.Uint64
		b.SetParallelism(3) // 3 * GOMAXPROCS 个 goroutine 数
		b.RunParallel(func(pb *testing.PB) {
			procs.Add(1)
			for pb.Next() { // 循环体在所有 goroutine 中总共执行 b.N 次
				iters.Add(1)
			}
		})
		// 校验是否创建了 3 * GOMAXPROCS 个 goroutine 数
		if want := uint32(3 * runtime.GOMAXPROCS(0)); procs.Load() != want {
			b.Errorf("got %v procs, want %v", procs.Load(), want)
		}
		// 校验原子计数是否和 b.N 相等
		if iters.Load() != uint64(b.N) {
			b.Errorf("got %v iters, want %v", iters.Load(), b.N)
		}
	})

	b.Run("ReportMetric", func(b *testing.B) {
		var compares int64
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			s := doSomething(100)
			sort.Slice(s, func(i, j int) bool {
				compares++
				return s[i] < s[j]
			})
		}
		// 度量 `metric` 是 per-operation, 应除以 b.N, 单位为 metric/op
		b.ReportMetric(float64(compares)/float64(b.N), "compares/op")
		// 度量 `metric` 是 per-time, 应除以 b.Elapsed, 单位为 metric/ns
		b.ReportMetric(float64(compares)/float64(b.Elapsed().Nanoseconds()), "compares/ns")
	})

	b.Run("ReportMetricInParallel", func(b *testing.B) {
		var compares atomic.Int64
		s := doSomething(100)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sort.Slice(s, func(i, j int) bool {
					compares.Add(1)
					return s[i] < s[j]
				})
			}
		})
		b.SetBytes(int64(len(s)) * int64(unsafe.Sizeof(int(1))))
		b.ReportMetric(float64(compares.Load())/float64(b.N), "compares/op")
		b.ReportMetric(float64(compares.Load())/float64(b.Elapsed().Nanoseconds()), "compares/ns")
	})
}

/*
! F.non-common Functions
	Add			将把参数添加到模糊测试的种子语料库中。args 必须匹配 Fuzz 目标的参数
	Fuzz   		运行模糊目标函数 ff 进行开始模糊测试。如果 ff 对一组模糊参数失败，这些参数将被添加到种子语料库；
				ff 不能调用任何 *F 方法，例如 (*F).Log、(*F).Error、(*F).Skip。使用相应的 *T 方法。
				在 ff 中允许的唯一 *F 方法是 (*F).Failed 和（*F).Name
*/
//? go test -v -run=NONE -fuzz=^Fuzz_Functions$ [-fuzztime=3s] [-parallel=8] [-short]
func Fuzz_Functions(f *testing.F) {
	//? 防止异常崩溃
	mlimit := debug.SetMemoryLimit(1000 * 100)
	mTreads := debug.SetMaxThreads(1000)
	f.Cleanup(func() {
		debug.SetMemoryLimit(mlimit)
		debug.SetMaxThreads(mTreads)
	})
	for _, v := range []int{10, 20, 30, 50, 100} {
		f.Add(v)
	}
	f.Fuzz(func(t *testing.T, a int) {
		if f.Failed() {
			t.Skipf("%s Failed", f.Name())
		}
		if a < 2 {
			t.Skip()
		}
		doSomething(a)
	})
}

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
