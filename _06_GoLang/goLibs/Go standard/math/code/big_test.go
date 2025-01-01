package gostd

import (
	"math"
	"math/big"
	"testing"
)

/*
! NewInt 构造一个 *big.Int
! big.Word 表示多精度无符号整数的单个数字。
! big.Int 表示有符号多精度整数
	Sign 符号位
	[Set]: Set, SetInt64, SetUint64, SetBytes, SetBits
	SetBit 设置 z 为 x 并将 x 的第 i 位设置位 0 或 1
	Bit 返回 x 的第 i 位
	[To]: Int64, Uint64, Float64, Bytes, FillBytes
	[Unary]: Abs 绝对值, Neg 负,
	[Binary]: Add 加法, Sub 减法, Mul 乘法, Div 欧几里得除法, Mod 欧几里得取模, Quo 截断除法, Rem 截断取余, Sqrt 平方根
	[Logic] Lsh <<, Rsh >>, And &, AndNot &^, Or |, Xor ^, Not ^x,
	[Pred]: Cmp 比较, CmpAbs 绝对值比较, IsInt64, IsUint64
	QuoRem 设置 z 为 x Quo y, 设置 r 为 x Rem y
	DivMod 设置 z 为 x Div y, 设置 r 为 x Mod y
	MulRange 设置 z 为 [a,b] 内所有整数的乘积
	Binomial 设置 z 为二项分布系数 C(n,k)
	BitLen 返回 x 绝对值的长度
	TrailingZeroBits 返回 x 绝对值位中连续最少有效零位的个数
	Exp 设置 z 为 x^y mod |m|
	GCD 设置 z 为 a,b 的最大公约数; 若 x,y 不为 nil 则设置 x,y 以满足 z == a*x+b*y
	Rand 将 z 设置为 [0,n) 中的伪随机数
	ModInverse 设置 z 为环 Z/nZ 中 g 的乘法逆
	ModSqrt 设置 z 为 x Mod p 的平方根
! Int.conv
	Text 返回给定基数(2~63)的字符串形式
	Append 对 x 的字符串表示形式追加到 buf
	String 返回 Text(10)
	SetString 给定基数和字符串转换为 Int
	Format 实现 fmt.Formatter
	Scan 对 fmt.Scanner 进行支持
! Int.marsh
	GobEncode & GobDecode 实现 encoding/gob.GobEncoder & GobDecoder
	MarshalText & UnmarshalText 实现 encoding.TextMarshaler & TextUnmarshaler
	MarshalJSON & UnmarshalJSON 实现 encoding/json.Marshaler & Unmarshaler
! Jacobi 返回 Jacobi 符号 (x/y)，+1、-1 或 0。y 参数必须是奇数
*/

func TestBigInt(t *testing.T) {
	var i = big.NewInt(0)
	i = set(i.Set, big.NewInt(1))
	i = set(i.SetInt64, math.MinInt64)
	i = set(i.SetUint64, math.MaxInt64)
	i = set(i.SetBits, []big.Word{math.MaxInt, math.MaxUint64})
	i = set(i.SetBytes, []byte{12, 34, 56, 78, 90})
	if i, ok := i.SetString("1234567890ABCDEF", 16); ok {
		logfln("set `1234567890ABCDEF` to %#X", i)
	}
	for n := range 20 {
		i = i.SetBit(i, n, 1)
	}
	logfln("set i(bit[0,20): %b) to 1", i)
	logfln("TrailingZeroBits(%v) : %d", i, i.TrailingZeroBits())
	i1 := big.NewInt(0).Set(i)
	logfln("Binomial(%v, n:10, k:2) : %d", i1, i.Binomial(10, 2))

	i2 := big.NewInt(0)
	fn := func(x, y int64) {
		R := big.NewInt(0)
		M := big.NewInt(0)
		d, _ := i2.SetInt64(x).DivMod(i2, big.NewInt(y), M)
		logfln("%d DivMod %d = %d .... %d", x, y, d, M)
		q, _ := i2.SetInt64(x).QuoRem(i2, big.NewInt(y), R)
		logfln("%d QuoRem %d = %d .... %d", x, y, q, R)
	}
	fn(11, 3)
	fn(11, -3)
	fn(-11, 3)
	fn(-11, -3)

	i.SetInt64(10086)
}
func set[T, U any](fn func(t T) *U, t T) *U {
	r := fn(t)
	logfln("set %v to %v", t, r)
	return r
}

/*
! NewFloat 构造一个 *big.Float
! big.RoundingMode 表示如何舍入 Float
! big.Accuracy 描述最近一次产生的浮点舍入误差
! big.Float 表示非零有限浮点数表示的多精度浮点数
	SetPrec & Prec 设置（可能四舍五入）或返回 z 的尾数精度
	MinPrec 返回精确表示 x 所需的最小精度
	SetMode & Mode 设置或获取 z 的舍入模式
	Acc 返回最近一次操作产生 x 的精度
	MantExp 分解 x 为尾数和指数并返回指数，非空 mant 参数设置为具有相同精度的尾数
	SetMantExp 将 z 设置为 mant*2^exp
	[Is]: IsInf, IsInt
	Sign 报告 x < 0 or =±0 or >0
	SignBit 报告符号位，是正 0 还是负
	[Set]: Set, SetInt64, SetUint64, SetFloat64, SetInt, SetRat, SetInf
	[To]: Copy, Uint64, Int64, Float32, Float64, Int, Rat
	[Unary]: Abs, Neg, Sqrt
	[Binary]: Add, Sub, Mul, Quo,
	[Pred]: Cmp
! Float.conv
	SetString, Parse, Scan
	Text, String, Append, Format
! Float.marshal
	GobEncode & GobDecode
	MarshalText & UnmarshalText
! ParseFloat 相当于 `Float.Parse`
*/

func TestBigFloat(t *testing.T) {
	f := big.NewFloat(1.234567891011121314)
	logfln("prec: %d, value: %x", f.Prec(), f)
	f = f.SetPrec(5)
	logfln("set prec: %d, value: %x", f.Prec(), f)
	mant := new(big.Float)
	exp := f.MantExp(mant)
	logfln("%x: mant= %x, exp= %d", f, mant, exp)
	f.SetPrec(53)
}

/*
! NewRat 构造一个 *big.Rat
! big.Rat 表示任意精度的商 a/B
	[Unary]: Abs, Denom 分母, Neg, Num 分子
	[Set]: Set, SetFloat64, SetFrac, SetFrac64, SetInt, SetUint64
	[Binary]: Add, Mul, Quo, Sub
	[Pred]: Cmp,
	[To]: Float32, Float64, FloatPrec 指定精度,
	[Is]: IsInt
	Sign 报告 x < 0 or =±0 or >0
	Inv 设置 1/x
! big.conv
	FloatString, RatString, Scan
	SetString, String
! big.marshal
	GobDecode & GobEncode
	MarshalText % UnmarshalText
*/

func TestBigRat(t *testing.T) {
	r := big.NewRat(1, 3)
	log(r.FloatString(53), r.RatString(), r.String())
	log(r.SetString("1/6"))
	log(r.SetString("0o11/6"))
	log(r.SetString("0x11/6"))
	log(r.SetString("0b11/6"))
}
