package gostd

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"text/tabwriter"
	"time"
)

func TestRand(t *testing.T) {
	answers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}
	for range 10 {
		fmt.Println("Magic 8-Ball says:", answers[rand.Intn(time.Now().Nanosecond())%len(answers)])
	}
}

/*
! rand.Source 表示在范围 [0，1<<63) 中的均匀分布的伪随机 int64 值的源。NewSource 构造一个 Source
! rand.Rand 表示为随机数的来源。New 构造一个 Rand
	Seed 使用提供的 seed 将生成器初始化为确定状态, 它不应调用其他 Rand 函数
	Float32, Float64 以相应浮点数类型的形式返回半开区间 [0.0，1.0) 内的伪随机数
	Int, Int31, Int63, Intn, Int31n, Int63n, Uint32, Uint64 返回一个非负的伪随机整数。
	ExpFloat64 返回一个在 (0，+[math.MaxFloat64]] 范围内呈指数分布的 float64，其指数分布的速率参数 lambda 为 1，平均值为 1/lambda(1)
	NormFloat64 返回一个在 [-math.MaxFloat64, +[math.MaxFloat64]] 范围内的正态分布 float64，具有标准正态分布 (mean = 0, stddev = 1)。
	Perm 返回 n 个整数的切片，半开区间 [0，n) 中整数的伪随机序列
	Read 生成 len(p) 随机字节并将它们写入 p。它总是返回 len(p) 和 nil 错误。Read 不应与任何其他 Rand 方法并发调用。
	Shuffle 伪随机化元素的顺序。n 是元素的个数。如果 n < 0，则 Panic。swap 交换索引为 i 和 j 的元素。
*/

func TestRandFunctions(t *testing.T) {
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	r := rand.New(rand.NewSource(-99))

	// The tabwriter here helps us generate aligned output.
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	defer w.Flush()
	show := func(name string, v1, v2, v3 any) {
		r.Seed(int64(time.Now().Nanosecond()))
		fmt.Fprintf(w, "%s\t%v\t%v\t%v\n", name, v1, v2, v3)
	}

	// Float32 and Float64 values are in [0, 1).
	show("Float32", r.Float32(), r.Float32(), r.Float32())
	show("Float64", r.Float64(), r.Float64(), r.Float64())

	// ExpFloat64 values have an average of 1 but decay exponentially.
	show("ExpFloat64", r.ExpFloat64(), r.ExpFloat64(), r.ExpFloat64())

	// NormFloat64 values have an average of 0 and a standard deviation of 1.
	show("NormFloat64", r.NormFloat64(), r.NormFloat64(), r.NormFloat64())

	// Int31, Int63, and Uint32 generate values of the given width.
	// The Int method (not shown) is like either Int31 or Int63
	// depending on the size of 'int'.
	show("Int", r.Int(), r.Int(), r.Int())
	show("Int31", r.Int31(), r.Int31(), r.Int31())
	show("Int63", r.Int63(), r.Int63(), r.Int63())
	show("Uint32", r.Uint32(), r.Uint32(), r.Uint32())
	show("Uint64", r.Uint64(), r.Uint64(), r.Uint64())

	// Intn, Int31n, and Int63n limit their output to be < n.
	// They do so more carefully than using r.Int()%n.
	show("Intn(10)", r.Intn(10), r.Intn(10), r.Intn(10))
	show("Int31n", r.Int31n(1<<16), r.Int31n(1<<16), r.Int31n(1<<16))
	show("Int63n", r.Int63n(1<<32), r.Int63n(1<<32), r.Int63n(1<<32))

	// Perm generates a random permutation of the numbers [0, n).
	show("Perm", r.Perm(5), r.Perm(5), r.Perm(5))
}

/*
! rand.Zipf 生成 Zipf 分布变量。
	Uint64 返回一个从 Zipf 对象描述的 Zipf 分布中提取的值。
! NewZipf 返回 Zipf 变量生成器。生成器生成值 k ∈ [0, imax]，使得 P(k) 与 (v + k)^(-s) 成比例。要求：s > 1 且 v >= 1
*/

func TestZipf(t *testing.T) {
	z := rand.NewZipf(rand.New(rand.NewSource(int64(time.Now().Nanosecond()))), 1.25, 10, 1<<8)
	for range 10 {
		log(z.Uint64())
	}
}
