package gostd

import (
	"slices"
	"testing"
)

var orderedS = []int{0, 1, 2, 3, 4, 8, 10, 56, 87, 89, 99, 100}
var unorderedS = []int{10, 5, -5, -87, 88, -1, -10, 10, -33, 89, -1, 748, 25, 1, 25, 33, 48}

func foreach(print bool, s []int) {
	if print {
		logfln("len(%d), cap(%d): %v", len(s), cap(s), s)
	}
}

// ! BinarySearch 在一个已正向排序的切片中搜索 target，并返回找到 target 的位置
// ! BinarySearchFunc 类似于 BinarySearch，但使用自定义比较函数 cmp, 匹配时返回 0；小于返回负数，大于返回正数
func TestBinarySearch(t *testing.T) {
	search := func(target int) {
		if i, ok := slices.BinarySearch(orderedS, target); ok {
			logfln("target(%d) is in Slice[%d]:%d", target, i, orderedS[i])

		} else {
			logfln("target(%d) is not in Slice", target)
		}
	}
	searchFunc := func(target int) {
		if i, ok := slices.BinarySearchFunc(orderedS, target, func(E, T int) int {
			if E*E > T {
				return 1
			} else if E*E == T {
				return 0
			}
			return -1
		}); ok {
			logfln("sqrt(target(%d)) is in Slice[%d]:%d", target, i, orderedS[i])
		} else {
			logfln("target(%d) is not in Slice", target)
		}
	}

	search(10)
	search(5)

	searchFunc(16)
	searchFunc(100)
	searchFunc(25)
}

// ! Clip 从切片中删除未使用的容量，返回 s[0:len(s):len(s)]
// ! Clone 浅克隆返回切片的副本
// ! Concat 合并多个切片
// ! Grow 会在必要时增加切片的容量，以保证额外 n 个元素的空间
// ! Insert 在索引 i 处插入元素
// ! Replace 将元素 s[i:j] 替换为给定的 v...; 当 len(v)<(j-i) 时，其余的元素归零
// ! Reverse 将切片的元素反转
// ? go test -v -run=^TestSlicesFunctions$
func TestSlicesFunctions(t *testing.T) {
	testSlicesFunctions(true)
}

func testSlicesFunctions(print bool) []int {
	s := make([]int, len(unorderedS), 2*len(unorderedS))
	copy(s, unorderedS)
	foreach(print, s)

	cop := slices.Clone(s)
	cop = slices.Clip(cop)
	foreach(print, cop)

	cop = slices.Grow(cop, len(orderedS))
	cop = slices.Concat(cop, orderedS)
	foreach(print, cop)

	cop = slices.Replace(cop, 0, 3, []int{10, -10, 10}...)
	foreach(print, cop)

	slices.Reverse(cop)
	foreach(print, cop)

	cop = slices.Insert(cop, 10, 10, 10, 10, 10)
	foreach(print, cop)

	return cop
}

// ! Compact 修改切片并用单个副本替换连续的相等元素；返回的切片可能更短
// ! CompactFunc 类似于 Compact，但使用 eq 函数来比较元素；CompactFunc 保留第一个
func TestCompact(t *testing.T) {
	s := testSlicesFunctions(false)
	foreach(true, s)
	s1 := slices.Compact(slices.Clone(s))
	foreach(true, s1)

	s1 = slices.CompactFunc(slices.Clone(s), func(e1, e2 int) bool {
		if e1+e2 == 0 || e1 == e2 {
			return true
		}
		return false
	})
	foreach(true, s1)
}

//! Compare 使用 cmp.Compare 比较两个切片的元素，元素都相等时比较切片的长度
//! CompareFunc 类似于 Compare，但使用自定义的比较函数
//! Contains 报告 v 是否存在于 s 中
//! ContainsFunc 报告 s 中是否至少有一个元素 e 满足 f(e)
//! Delete 从 s 中删除元素 s[i:j]；j 大于 len(s) 或 s[i:j] 不是有效的切片则 panic
//! DeleteFunc 从 s 中删除满足函数 del 的所有元素
//! Equal 报告两个切片是否相等：长度相同且所有元素都相等
//! EqualFunc 类似于 Equal，对每对元素使用 eq 函数。如果长度不同，则返回 false
//! Index 返回 v 在 s 中第一次出现的索引
//! IndexFunc 返回第一个满足 f(s[i]) 的索引 i
//! IsSorted 报告 s 是否按升序排序。
//! IsSortedFunc 报告 x 是否按升序排序，将 cmp 作为 SortFunc 比较函数
//! Max 返回 s 中的最大值
//! MaxFunc 返回 s 中的最大值，使用 cmp 比较元素
//! Min 返回 s 中的最小值
//! MinFunc 返回 s 中的最小值，使用 cmp 比较元素
//! Sort 按升序对任何有序类型的切片进行排序
//! SortFunc 按照 cmp 函数确定的升序对切片 s 进行排序
//! SortStableFunc 对切片 s 进行排序，同时保持相等元素的原始顺序
