package gostd

import (
	"container/heap"
	"testing"
)

/*
! heap.Interface: sort.Interface, Push(x any), Pop() any
! Init 建立供包中其他例程所需的 heap 不变量。Init 对于堆不变量是幂等的，并且可以在堆不变量失效时调用
! Fix 在索引 i 处的元素更改其值后重新建立堆排序。更改 i 处的值，然后调用 Fix，这相当于调用 Remove(h，i)，然后调用 Push(h, x) 新值
! Pop 从堆中移除并返回最小元素（根据 Less）
! Push 将元素 x 推到堆上。
! Remove 从堆中移除并返回索引 i 处的元素
*/

// An IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length, not just its contents.
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func TestIntHeap(t *testing.T) {
	h := &IntHeap{64, 5, 8, 9, 74, 5, 8, 7, 8, 1, 89, 2}
	heap.Init(h)

	heap.Push(h, 50)
	logfln("min num : %d", (*h)[0])

	(*h)[5] = 10086
	// heap.Fix(h, 5)
	log(heap.Remove(h, 1))
	log(heap.Remove(h, 1))
	log(heap.Remove(h, 1))

	for i := range len(*h) {
		logfln("i : %d, v : %v", i, heap.Pop(h))
	}
}
