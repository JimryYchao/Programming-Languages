package gostd

import (
	"fmt"
	"math"
	"testing"
)

/*
! cmplx-functions
	[Unary]: Abs, Signbit
	[三角]: Cos, Sin, Sincos, Tan,
	[反三角]: Acos, Asin, Atan, Atan2
	[双曲]: Cosh, Sinh, Tanh
	[反双曲]: Acosh, Asinh, Atanh
	[幂指对]: Sqrt, Cbrt, Exp, Exp2, Expm1, Ilogb, Log, Log10, Log1p, Log2, Logb, Pow, Pow10,
	Hypot 返回 Sqrt(p*p + q*q)
	Frexp 将 f 分解为归一化分数和 2 的整数幂。当满足 f == frac * 2^exp 时返回 frac 和 exp，frac 的绝对值在区间 [1/2，1）内
	Ldexp 是 Frexp 的倒数。它返回 frac * 2^exp。
[其他]: FMA, Inf, IsInf, NaN, IsNaN, Max, Min, Mod, Modf, Remainder
	Copysign 返回值为 f，符号为 sign 的值。
	Ceil 返回大于或等于 x 的最小整数值。
	Floor 返回小于或等于 x 的最大整数值。
	Trunc 截断为整数值。
	Dim 返回 x-y 或 0 的最大值。
	Nextafter32 返回 x 之后的下一个可表示的 float32 值。
	Nextafter 返回 x 之后的下一个可表示的 float64 值。
	Round 返回最接近的整数，向零舍入一半。
	RoundToEven 返回最接近的整数，舍入为偶数。
	Float32bits, Float64bits 返回浮点数 f 的 IEEE 754 二进制表示；
	Float32formbits, Float64formbits 是对应的逆运算
贝塞尔:
	第一类 J0, J1, JN
	第二类 Y0, Y1, YN
Erf:
	Erf 返回 x 的误差函数。
	Erfc 返回 x 的互补误差函数。
	Erfinv 返回 x 的逆误差函数。
	Erfcinv 返回 Erfc(x) 的倒数。
Gamma:
	Gamma 返回 x 的 Gamma 函数。
	Lgamma返回Gamma（x）的自然对数和符号（-1或+1）。
*/

func TestRound(t *testing.T) {
	fmt.Printf("%.1f\n", math.Round(10.5)) // 11

	fmt.Printf("%.1f\n", math.Round(-10.5)) // -11

	fmt.Printf("%.1f\n", math.RoundToEven(11.5)) // 12.0

	fmt.Printf("%.1f\n", math.RoundToEven(12.5)) // 12.0
}

func TestFloat2Bits(t *testing.T) {
	var fs = []float64{1.25, float64(+0), -0, math.Inf(1), math.Inf(-1), math.NaN(), 0.1, 0.1111111111, 0.333333333333, 123456.789, 1.0}
	var bits = []uint64{
		0b0011111111110100000000000000000000000000000000000000000000000000,
		0b0000000000000000000000000000000000000000000000000000000000000000,
		0b1000000000000000000000000000000000000000000000000000000000000000,
		0b0111111111110000000000000000000000000000000000000000000000000000,
		0b1111111111110000000000000000000000000000000000000000000000000000,
		0b0111111111111000000000000000000000000000000000000000000000000001,
		0b0011111110111001100110011001100110011001100110011001100110011010,
		0b0011111110111100011100011100011100011100011001011000111110011101,
		0b0011111111010101010101010101010101010101010101010011110111100001,
		0b0100000011111110001001000000110010011111101111100111011011001001,
		0b0011111111110000000000000000000000000000000000000000000000000000,
	}

	for _, f := range fs {
		logfln("%#15f >> %#064b", f, math.Float64bits(f))
	}

	for _, bit := range bits {
		logfln("%#064b >> %#-15f ", bit, math.Float64frombits(bit))
	}
}
