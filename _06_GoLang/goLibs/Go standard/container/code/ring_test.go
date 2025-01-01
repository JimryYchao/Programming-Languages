package gostd

import (
	"container/ring"
	"math"
	"testing"
)

/*
! New 创建一个 n 元素的环
! ring.Ring 是一个循环列表或环的元素。环没有开始或结束；指向任何环元素的指针用作整个环的引用。空环表示为 nil 环指针。环的零值是一个具有 nil 值的单元素环。
	Len, Next, Prev
	Do 	对环中的每个元素按前向顺序调用函数 f。如果 f 改变 *r，则 Do 的行为未定义
	Move 在环中向后 (n < 0) 或向前 (n >= 0) 移动 n % r.Len() 元素，并返回该环元素
	Link 连接环 r 和环 s，使得 r.Next 变成 s，并返回 r.Next() 的原始值。r 不能为空。
		如果 r 和 s 指向同一个环，连接它们会从环中删除 r 和 s 之间的元素。移除的元素形成一个子环，
		结果是对该子环的引用（如果没有元素被移除，结果仍然是 r.Next() 的原始值，而不是 nil）。
	Unlink 从环 r 中删除 n % r.Len() 元素，从 r.Next() 开始。如果 n % r.Len() == 0，则 r 保持不变。结果是被移除的子环。
*/

func TestRing(t *testing.T) {
	var r = ring.New(10)
	var tmp *ring.Ring = r
	for i := range r.Len() {
		tmp.Value = i
		tmp = tmp.Next()
	}
	sqrt := func(a any) {
		logfln("sqrt of %v is %.4f", a, math.Sqrt(float64(a.(int))))
	}
	iter := func(a any) {
		logfln("%v", a)
	}

	r.Do(sqrt) // 对每个 elem 开平方

	r2 := r.Unlink(5)
	for i := range r2.Len() {
		r2.Value = math.Pow(float64(r2.Value.(int)), (float64)(i))
		r2 = r2.Prev()
	}
	r = r.Link(r2)
	r.Do(iter)
}
