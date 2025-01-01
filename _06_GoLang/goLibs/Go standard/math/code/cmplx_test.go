package gostd

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"
)

/*
! cmplx-functions
	[Unary]: Abs,
	[三角]: Cos, Sin, Tan, Cot
	[反三角]: Acos, Asin, Atan
	[双曲]: Cosh, Sinh, Tanh
	[反双曲]: Acosh, Asinh, Atanh
	[幂指对]: Exp, Log, Log10, Pow, Sqrt
	[其他]: Conj 复共轭, Inf, IsInf, NaN, IsNaN, Phase 相位
	Polar 返回 c 的极坐标和相位角
	Rect 构造具有极坐标 r,θ 的复数
*/

func TestPolar(t *testing.T) {
	r, theta := cmplx.Polar(2i)
	fmt.Printf("r: %f, θ: %f*π", r, theta/math.Pi)
}
