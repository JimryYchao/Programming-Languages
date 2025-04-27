// 此示例演示了使用堆接口构建的优先级队列。
package examples

import (
	"container/heap"
	"fmt"
	"testing"
)

// Item 是我们在优先队列中管理的东西。
type Item struct {
	value    string // Item 的值；任意
	priority int    // 优先级
	// 更新需要索引，并由 heap.Interface 方法维护
	index int // Item 在堆中的索引
}

// PriorityQueue实现 heap.Interface 并保存 Item。
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	// 希望 Pop 给我们的优先级是最高的，而不是最低的，所以使用 >。
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // 避免内存泄漏
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update 修改队列中 Item 的优先级和值。
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}

// ! Test 创建了一个 PriorityQueue，添加和操作一个 Item，然后按优先级顺序移除这些 Items。
func TestPriorityQueue(*testing.T) {
	items := map[string]int{
		"banana": 3, "apple": 2, "pear": 4,
	}

	// 创建一个优先级队列，将 Item 放入其中，并建立优先级队列 (堆) 不变量。
	pq := make(PriorityQueue, len(items))
	i := 0
	for value, priority := range items {
		pq[i] = &Item{
			value:    value,
			priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// 插入新 Item，然后修改其优先级。
	item := &Item{
		value:    "orange",
		priority: 1,
	}
	heap.Push(&pq, item)
	pq.update(item, item.value, 5)

	// Pop Item, 它们按优先级递减的顺序弹出。
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		fmt.Printf("%.2d:%s ", item.priority, item.value)
	}
}
