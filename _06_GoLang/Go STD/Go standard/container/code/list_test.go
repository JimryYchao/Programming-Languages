package gostd

import (
	"container/list"
	"testing"
)

/*
! New 范围一个初始化的 List
! list.List 表示一个双向链表，它的零值是一个空列表；元素为 Element
	Init 	初始化并清空列表
	Len 	返回链表的长度
	Back 	返回最后一个元素；列表为空返回 nil
	Front 	返回第一个元素；列表为空返回 nil
	Remove 	移除元素 e 并返回 e.Value
	InsertAfter, InsertBefore 	在 mark 元素后 (前) 插入一个值为 v 的新元素 e 并返回；mark 非 l 的元素无修改，不能为 nil
	MoveAfter, MoveBefore 		将元素 e 移动到 mark 之前（后）的新位置
	MoveToBack, MoveToFront 	移动元素到 l 的后面（前面）
	PushBack, PushFront 		在链表后 (前) 插入一个值为 v 的新元素
	PushBackList, PushFrontList 在链表后 (前) 插入另一个链表的副本
! Element 为链表的元素类型
	Value, Next, Prev
*/

func TestList(t *testing.T) {
	l1 := list.New()
	lter := func(l *list.List) {
		e := l.Front()
		values := make([]any, 0, l.Len())
		for e != nil {
			values = append(values, e.Value)
			e = e.Next()
		}
		log(values...)
	}

	for i := range 5 {
		l1.PushBack(i + 1) // 1,2,3,4,5
	}

	l2 := list.New()
	e := l1.Front()
	for e != nil {
		l2.PushBack(e.Value.(int) * 10)
		e = e.Next()
	}
	lter(l2) // 10,20,30,40,50

	l2.MoveBefore(l2.Back(), l2.Front())
	lter(l2) // 50,10,20,30,40

	l1.PushBackList(l2)
	lter(l1) // 1 2 3 4 5 50 10 20 30 40

	l1.Init()
	lter(l1) // clear
	lter(l2) // 50 10 20 30 40
}
