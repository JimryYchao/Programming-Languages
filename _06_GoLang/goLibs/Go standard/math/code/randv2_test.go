package gostd

import (
	"fmt"
	"math"
	r "math/rand"
	"math/rand/v2"
	"os"
	"testing"
	"text/tabwriter"
	"time"
)

type Src struct {
	r.Source
}

func (s *Src) Uint64() uint64 {
	return uint64(s.Source.Int63() + int64((math.MaxInt64+time.Now().Nanosecond()*1357986543*time.Now().Nanosecond())%math.MaxInt32))

}
func NewSource() *Src {
	return &Src{r.NewSource(int64(time.Now().Nanosecond()))}
}

/*
! rand.Source 表示在范围 [0，1<<64) 中的均匀分布的伪随机 uint64 值的源。NewSource 构造一个 Source
! rand.Rand 表示为随机数的来源。New 构造一个 Rand
	Seed 使用提供的 seed 将生成器初始化为确定状态, 它不应调用其他 Rand 函数
	Float32, Float64 以相应浮点数类型的形式返回半开区间 [0.0，1.0) 内的伪随机数
	Int, Int32, Int64, IntN, Int32N, Int64N, Uint32, Uint64 返回一个非负的伪随机整数。
	ExpFloat64 返回一个在 (0，+[math.MaxFloat64]] 范围内呈指数分布的 float64，其指数分布的速率参数 lambda 为 1，平均值为 1/lambda(1)
	NormFloat64 返回一个在 [-math.MaxFloat64, +[math.MaxFloat64]] 范围内的正态分布 float64，具有标准正态分布 (mean = 0, stddev = 1)。
	Perm 返回 n 个整数的切片，半开区间 [0，n) 中整数的伪随机序列
	Read 生成 len(p) 随机字节并将它们写入 p。它总是返回 len(p) 和 nil 错误。Read 不应与任何其他 Rand 方法并发调用。
	Shuffle 伪随机化元素的顺序。n 是元素的个数。如果 n < 0，则 Panic。swap 交换索引为 i 和 j 的元素。
*/

func TestRandv2Functions(t *testing.T) {
	// Create and seed the generator.
	// Typically a non-fixed seed should be used, such as time.Now().UnixNano().
	// Using a fixed seed will produce the same output on every run.
	r := rand.New(NewSource())

	// The tabwriter here helps us generate aligned output.
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	defer w.Flush()
	show := func(name string, v1, v2, v3 any) {
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
	show("Int32", r.Int32(), r.Int32(), r.Int32())
	show("Int64", r.Int64(), r.Int64(), r.Int64())
	show("Uint32", r.Uint32(), r.Uint32(), r.Uint32())
	show("Uint64", r.Uint64(), r.Uint64(), r.Uint64())

	// Intn, Int31n, and Int63n limit their output to be < n.
	// They do so more carefully than using r.Int()%n.
	show("Intn(10)", r.IntN(10), r.IntN(10), r.IntN(10))
	show("Int32N", r.Int32N(1<<16), r.Int32N(1<<16), r.Int32N(1<<16))
	show("Int64N", r.Int64N(1<<32), r.Int64N(1<<32), r.Int64N(1<<32))

	// Perm generates a random permutation of the numbers [0, n).
	show("Perm", r.Perm(5), r.Perm(5), r.Perm(5))
}

/*
! rand.ChaCh8 是一个基于 ChaCha8 的强随机数生成器。NewChaCha8 返回一个新的 ChaCha8
	marshal: MarshalBinary & UnmarshalBinary
	Seed 重置 ChaCha8, 等效于 NewChaCha8(seed)
	Uint64 返回均匀分布的随机 uint64 值
*/

func TestChaCha8(t *testing.T) {
	cc8 := rand.NewChaCha8([32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 21, 22, 23, 24, 25, 36, 37, 38, 39, 40, 11, 13, 15, 17, 19, 20, 30, 40, 50, 66, 67, 68})
	for range 10 {
		log(cc8.Uint64())
	}
}

/*
! rand.PCG 是具有 128 位内部状态的 PCG 生成器。零 PCG 等于 NewPCG(0, 0)。NewPCG 返回以给定值作为种子的新 PCG
	MarshalBinary & UnmarshalBinary
	Seed 重置 PCG, 等效于 NewPCG(seed1, seed2)
	Uint64 返回均匀分布的随机 uint64 值。
*/

func TestPCG(t *testing.T) {
	pcg := rand.NewPCG(0, math.MaxUint64)
	for range 10 {
		log(pcg.Uint64())
	}
}

/*
! rand.Zipf 生成 Zipf 分布变量。
	Uint64 返回一个从 Zipf 对象描述的 Zipf 分布中提取的值。
! NewZipf 返回 Zipf 变量生成器。生成器生成值 k ∈ [0, imax]，使得 P(k) 与 (v + k)^(-s) 成比例。要求：s > 1 且 v >= 1
*/

func TestV2_Zipf(t *testing.T) {
	pcg := rand.NewPCG(0, math.MaxUint64)
	for range 10 {
		log(pcg.Uint64())
	}
	z := rand.NewZipf(rand.New(NewSource()), 1.25, 10, 1<<8)
	for range 10 {
		log(z.Uint64())
	}
}
