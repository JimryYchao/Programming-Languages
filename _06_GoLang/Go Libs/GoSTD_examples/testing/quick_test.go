package gostd_testing

/* 包 quick 提供了快速测试功能
! quick 包主要用于生成随机测试数据
! 核心功能：
! - Check: 检查函数是否对所有随机输入都返回 true
! - CheckEqual: 检查两个函数是否对所有随机输入都返回相同结果
*/

import (
	"testing"
	"testing/quick"
)

// ! 测试函数：检查整数乘法的结合律
// ? go test -v -run=TestQuickCheckMultiply
func TestQuickCheckMultiply(t *testing.T) {
	F := func(a, b, c int) bool {
		return (a*b)*c == a*(b*c)
	}
	if err := quick.Check(F, nil); err != nil {
		t.Fatalf("Multiplication is not associative: %v", err)
	}
}
