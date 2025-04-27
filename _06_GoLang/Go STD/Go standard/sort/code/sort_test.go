package gostd

import (
	"fmt"
	. "gostd/examples"
	"slices"
	"sort"
	"testing"
)

/*
! sort.Interface 定义一个可排序的接口对象，sort 的方法通过整数索引引用其基础集合的元素
	Len 返回集合中的元素个数
	Less 报告 i 是否必须在 j 之前排序
	Swap 交换两个元素
! Sort 对 Interface 进行升序排序
! IsSorted 报告 Interface 是否已排序
! Reverse 返回 Interface 的降序, 调用 Sort(Reverse(Interface))
*/

func TestSort(t *testing.T) {
	log(people)

	sort.Sort(people)
	log(people)

	sort.SliceStable(people, func(i, j int) bool {
		return people[i].Age > people[j].Age
	})
	log(people) // reverse
}

var people = People{
	{"Bob", 31},
	{"John", 42},
	{"Michael", 17},
	{"Jenny", 26},
}

func TestSortKeys(t *testing.T) {
	var planets = []Planet{
		{"Mercury", 0.055, 0.4},
		{"Venus", 0.815, 0.7},
		{"Earth", 1.0, 1.0},
		{"Mars", 0.107, 1.5},
	}

	name := func(p1, p2 *Planet) bool {
		return p1.Name < p2.Name
	}
	mass := func(p1, p2 *Planet) bool {
		return p1.Mass < p2.Mass
	}
	distance := func(p1, p2 *Planet) bool {
		return p1.Distance < p2.Distance
	}
	decreasingDistance := func(p1, p2 *Planet) bool {
		return distance(p2, p1)
	}

	// Sort the planets by the various criteria.
	By(name).Sort(planets)
	fmt.Println("By name:", planets)

	By(mass).Sort(planets)
	fmt.Println("By mass:", planets)

	By(distance).Sort(planets)
	fmt.Println("By distance:", planets)

	By(decreasingDistance).Sort(planets)
	fmt.Println("By decreasing distance:", planets)
}

func TestSortMulKeys(t *testing.T) {
	var changes = []Change{
		{"gri", "Go", 100},
		{"ken", "C", 150},
		{"glenda", "Go", 200},
		{"rsc", "Go", 200},
		{"r", "Go", 100},
		{"ken", "Go", 200},
		{"dmr", "C", 100},
		{"r", "C", 150},
		{"gri", "Smalltalk", 80},
	}

	user := func(c1, c2 *Change) bool {
		return c1.User < c2.User
	}
	language := func(c1, c2 *Change) bool {
		return c1.Language < c2.Language
	}
	increasingLines := func(c1, c2 *Change) bool {
		return c1.Lines < c2.Lines
	}
	decreasingLines := func(c1, c2 *Change) bool {
		return c1.Lines > c2.Lines // Note: > orders downwards.
	}

	// Simple use: Sort by user.
	OrderedBy(user).Sort(changes)
	fmt.Println("By user:", changes)

	// More examples.
	OrderedBy(user, increasingLines).Sort(changes)
	fmt.Println("By user,<lines:", changes)

	OrderedBy(user, decreasingLines).Sort(changes)
	fmt.Println("By user,>lines:", changes)

	OrderedBy(language, increasingLines).Sort(changes)
	fmt.Println("By language,<lines:", changes)

	sorter := OrderedBy(language, increasingLines, user)
	sorter.Sort(changes)
	fmt.Println("By language,<lines,user:", changes)

	sort.Sort(sort.Reverse(sorter))
	fmt.Println("reverse:", changes)
}

func TestSortWrapper(t *testing.T) {
	s := Organs{
		{"brain", 1340},
		{"heart", 290},
		{"liver", 1494},
		{"pancreas", 131},
		{"prostate", 62},
		{"spleen", 162},
	}
	printOrgans := func(s Organs) {
		for _, o := range s {
			fmt.Printf("%-8s (%v)\n", o.Name, o.Weight)
		}
	}

	sort.Sort(ByWeight{s})
	fmt.Println("Organs by weight:")
	printOrgans(s)

	sort.Sort(ByName{s})
	fmt.Println("Organs by name:")
	printOrgans(s)

	fmt.Println("reverse Organs by name:")
	if sort.IsSorted(ByName{s}) {
		sort.Sort(sort.Reverse(ByName{s}))
		printOrgans(s)
	}
}

/*
! Find 在 [0,n) 中二分查找使得 cmp(i)<=0 的最小索引，没有时返回 n；i < n 且 cmp(i) == 0 则 found = true
	cmp 要求对 i 之前满足 cmp(i0) > 0, 之后的元素 <= 0
! Search 在 [0,n) 中二分查找使得 f(i) 为真的最小索引；f 要求对 i 之前的元素返回 false，之后的元素都返回 true
*/

func TestSeachers(t *testing.T) {
	people := people
	sort.Sort(people)
	t.Run("Find", func(t *testing.T) {
		target := 27
		i, found := sort.Find(people.Len(), func(i int) int {
			return target - people[i].Age
		})
		if found {
			fmt.Printf("found age %d at entry %d\n", target, i)
		} else {
			fmt.Printf("age %d not found, would insert at %d", target, i)
			people = slices.Insert(people, i, Person{Name: "Ychao", Age: target})
			log(people)
		}
	})

	t.Run("Search", func(t *testing.T) {
		i := sort.Search(len(people), func(i int) bool {
			return people[i].Name == "Ychao"
		})
		if i > 0 {
			log(people[i])
		}
	})
}

/*
! Sort: Float64s, Ints, Strings, Slice, SliceStable 对切片的进行升序排序
	SliceStable 可以保证相等元素保持原始位置
! IsSorted: Float64sAreSorted, IntsAreSorted, StringsAreSorted, SliceIsSorted 报告切片是否是升序排序
! sort.Float64Slice, IntSlice, StringSlice 实现 Interface 接口
	x.Sort 使用 Sort(x)
	x.SearchT 使用 Search(x)
! SearchT: SearchFloat64s, SearchInts, SearchStrings 查找升序切片中元素 x 索引；不存在时返回 x 可以插入位置的索引
*/

func TestSortSlices(t *testing.T) {
	ints := []int{3, 4, 5, 8, 7, 1, 7, 4, 8, 8, 41, 7, 48}
	sort.Ints(ints)
	find := func(target int) {
		n := sort.Search(len(ints), func(i int) bool {
			return ints[i] >= target
		})
		if n < len(ints) && ints[n] == target {
			logfln("found %d in index %d", ints[n], n)
		} else {
			logfln("not found %d", target)
		}
	}
	find(10)
	find(8)
	find(41)
}

func TestFloat64Slice(t *testing.T) {
	floats := sort.Float64Slice([]float64{0, 1, -10, 1.1, .2, 12.3, 6.5, 3.3})
	sort.Sort(sort.Reverse(floats))
	log(floats)

	floats.Sort()
	i := floats.Search(10)
	if floats[i] != 10 {
		floats = slices.Insert(floats, i, 10)
	}
	log(floats)
}
