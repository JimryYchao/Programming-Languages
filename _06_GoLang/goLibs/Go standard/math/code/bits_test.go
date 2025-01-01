package gostd

import (
	"fmt"
	"math/bits"
	"testing"
)

/*
! Add, Add32, Add64 返回 x，y 和 carry 的和：sum = x + y + carry。
	carry 进位输入必须为 0 或 1。carryOut 输出保证为 0 或 1
*/

func TestAdd(t *testing.T) {
	add := func(n1, n2 []uint32) {
		d1, carry := bits.Add32(n1[1], n2[1], 0)
		d0, _ := bits.Add32(n1[0], n2[0], carry)
		nsum := []uint32{d0, d1}
		fmt.Printf("%v + %v = %v (carry bit was %v)\n", n1, n2, nsum, carry)
	}
	n1 := []uint32{33, 12} // 33<<32 + 12
	n2 := []uint32{21, 23} // 21<<32 + 23
	add(n1, n2)            // [33 12] + [21 23] = [54 35] (carry bit was 0)

	n1 = []uint32{1, 0x80000000} // 1<<32 + 2147483648
	n2 = []uint32{1, 0x80000000} // 1<<32 + 2147483648
	add(n1, n2)                  // [1 2147483648] + [1 2147483648] = [3 0] (carry bit was 1)
}

/*
! Sub, Sub32, Sub64 返回 x,y 和 borrow 的差值：diff = x - y - borrow。
	borrow 输入必须为 0 或 1; 否则行为未定义。borrowOut 为 0 或 1。
*/

func TestSub(t *testing.T) {
	sub := func(n1, n2 []uint32) {
		d1, carry := bits.Sub32(n1[1], n2[1], 0)
		d0, _ := bits.Sub32(n1[0], n2[0], carry)
		nsum := []uint32{d0, d1}
		fmt.Printf("%v - %v = %v (carry bit was %v)\n", n1, n2, nsum, carry)
	}
	n1 := []uint32{33, 23} // 33<<32 + 23
	n2 := []uint32{21, 12} // 21<<32 + 12
	sub(n1, n2)            // [33 23] - [21 12] = [12 11] (carry bit was 0)

	n1 = []uint32{3, 0x80000000 - 1} // 3<<32 + 2147483647
	n2 = []uint32{1, 0x80000000}     // 1<<32 + 2147483648
	sub(n1, n2)                      // [3 2147483647] - [1 2147483648] = [1 4294967295] (carry bit was 1)
}

/*
! Div, Div32, Div64 返回 (hi，lo)/y 的商和余数：quo =(hi，lo)/y，rem =(hi，lo)%y
	被除数位的上半部分在参数 hi 中，下半部分在参数 lo 中。
	Div 在 y == 0（被零除）或 y <= hi（商溢出）时 panic。
*/

func TestDiv(t *testing.T) {
	div := func(n1, n2 []uint32) {
		quo, rem := bits.Div32(n1[0], n1[1], n2[1])
		nsum := []uint32{quo, rem}
		fmt.Printf("[%v %v] / %v = %v\n", n1[0], n1[1], n2[1], nsum)
	}
	n1 := []uint32{0, 6} // 0<<32 + 6
	n2 := []uint32{0, 3} // 0<<32 + 3
	div(n1, n2)          // [0 6] / 3 = [2 0]

	n1 = []uint32{2, 0x80000000} // 2<<32 + 2147483648
	n2 = []uint32{0, 0x80000000} // 0<<32 + 2147483648
	div(n1, n2)                  // [2 2147483648] / 2147483648 = [5 0]
}

/*
! Mul, Mul32, Mul64 返回 x 和 y 的全宽乘积：(hi，lo) = x * y
	其中乘积位的上半部分返回到 hi，下半部分返回到 lo。
*/

func TestMul(t *testing.T) {
	mul := func(n1, n2 []uint32) {
		hi, lo := bits.Mul32(n1[1], n2[1])
		nsum := []uint32{hi, lo}
		fmt.Printf("%v * %v = %v\n", n1[1], n2[1], nsum)
	}
	n1 := []uint32{0, 12} // 0<<32 + 12
	n2 := []uint32{0, 12} // 0<<32 + 12
	mul(n1, n2)           // 12 * 12 = [0 144]

	n1 = []uint32{0, 0x80000000} // 0<<32 + 2147483648
	n2 = []uint32{0, 2}          // 0<<32 + 2
	mul(n1, n2)                  // 2147483648 * 2 = [1 0]
}

/*
! Rem, Rem32, Rem64 返回 (hi，lo)/y 的余数。
	y == 0（被零除）时会 panic，但与 Div 不同的是，它不会在商溢出时 panic
*/

func TestRem(t *testing.T) {
	rem := func(n1, n2 []uint32) {
		rem := bits.Rem32(n1[0], n1[1], n2[1])
		fmt.Printf("[%v %v] / %v ... %v\n", n1[0], n1[1], n2[1], rem)
	}
	n1 := []uint32{0, 7} // 0<<32 + 7
	n2 := []uint32{0, 3} // 0<<32 + 3
	rem(n1, n2)          // [0 7] / 3 ... 1

	n1 = []uint32{2, 0x80000000}     // 2<<32 + 2147483648
	n2 = []uint32{0, 0x80000000 - 1} // 0<<32 + 2147483648-1
	rem(n1, n2)                      // [2 2147483648] / 2147483647 ... 5
}

/*
! LeadingZeros, LeadingZeros8, LeadingZeros16, LeadingZeros32, LeadingZeros64 返回 x 前导零的个数
! TrailingZeros, TrailingZeros8, TrailingZeros16, TrailingZeros32, TrailingZeros64, 返回 x 中尾随零的个数
! Len, Len8, Len16, Len32, Len64 返回表示 x 所需的最小位数
! OnesCount, OnesCount8, OnesCount16, OnesCount32, OnesCount64 返回表示 x 中位 1 的个数
! Reverse, Reverse8, Reverse16, Reverse32, Reverse64 返回 x 的反转位值
! ReverseBytes, ReverseBytes16, ReverseBytes32, ReverseBytes64 返回 x 的反转字节值
! RotateLeft, RotateLeft8, RotateLeft16, RotateLeft32, RotateLeft64 返回左旋转 (k mod UintSize) 位的 x 的值。要将 x 向右旋转 k 位，调用 RotateLeft(x，-k)
*/

func TestBitsFunctions(t *testing.T) {
	logfln("LeadingZeros8(%08b) = %d", 1<<4, bits.LeadingZeros8(1<<4))     // 3
	logfln("TrailingZeros8(%08b) = %d", 1<<4, bits.TrailingZeros8(1<<4))   // 4
	logfln("Len8(%08b) = %d", 8, bits.Len8(8))                             // 4
	logfln("OnesCount8(%08b) = %d", 123, bits.OnesCount32(123))            // 6
	logfln("Reverse8(%08b) = %08b", 123, bits.Reverse8(123))               // 01111011 => 11011110
	logfln("ReverseBytes16(%016b) = %016b", 123, bits.ReverseBytes16(123)) // 0000000001111011) = 0111101100000000
	logfln("RotateLeft8(%08b, -1) = %08b", 123, bits.RotateLeft8(123, -1)) // 0111101_1 => 1_0111101
	logfln("RotateLeft8(%08b, 1) = %08b", 123, bits.RotateLeft8(123, 1))   // 0_1111011 => 1111011_0

}
