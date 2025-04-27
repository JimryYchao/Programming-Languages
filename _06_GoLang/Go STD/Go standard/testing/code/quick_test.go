package gostd

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
	"time"
)

/*
! Check 查找 f 的输入，任何返回 bool 的函数，使得 f 返回 false。它反复调用 f，每个参数都有任意值。
	如果 f 在给定的输入上返回 false，则 Check 将该输入作为 *CheckError 返回。
! CheckEqual 查找 f 和 g 返回不同结果的输入。它反复调用 f 和 g，每个参数都有任意值。
	如果 f 和 g 返回不同的答案，则 CheckEqual 返回一个描述输入和输出的 *CheckEqualError。
! quick.Config 包含用于运行测试函数 Check, CheckEqual 的选项。
	MaxCount 设置最大的迭代次数
	MaxCountScale 是应用于默认最大值的非负比例因子，通常为 100，可以使用 -quickchecks 标志设置
	Rand 指定随机数的来源
	Values 指定一个函数生成任意反射的切片，与被测试函数的参数一致；未设置时使用 `Value` 函数
! Value 返回给定类型的任意值。如果类型实现了 Generator 接口，则将使用该接口。注意：要为结构体创建任意值，必须导出所有字段。
! quick.Generator 包装 Generate，表示可以生成自己类型的随机值。
*/

func TestQuickCheck(t *testing.T) {
	c := &quick.Config{MaxCount: 100000, Rand: rand.New(rand.NewSource(int64(time.Now().Nanosecond())))}
	logfln("Check: %s", quick.Check(IsNonInf, c))
	logfln("CheckEqual: %s", (quick.CheckEqual(fSqrtMInt, fSqrt[mInt], c)))
}

func IsNonInf(num mInt, denom mInt) bool {
	return !math.IsInf(float64(num)/float64(denom), 0)
}

type Ints interface {
	~int
}

func fSqrtMInt(m mInt) float64 {
	return math.Sqrt(float64(m))
}
func fSqrt[I Ints](m I) float64 {
	return math.Sqrt(float64((m)))
}

type mInt int

func (m mInt) Generate(r *rand.Rand, _ int) reflect.Value {
	return reflect.ValueOf(mInt(int(r.Int31n(1000)))) // -500 ~ 500
}
