package gostd

import (
	"fmt"
	"strings"
	"testing"
	"unicode"
)

/*
! In 报告 r 是否是给定的 `RangeTable` ranges 之一的成员;
! Is 报告 r 是否在指定的 rangeTab 表中: IsXxx(r) ==> Is(tab, r)
	IsControl : Is(C, r), Is(Other, r)
	IsDigit
	IsGraphic : L、M、N、P、S、Zs
	IsLetter  : L
	IsLower, IsUpper
	IsMark	  : M
	IsNumber  : N
	IsOneOf   : 等效于 In
	IsPrint   : L, M, N, P, S, ASCII space ' '
	IsPunct   : P
	IsSpace   : '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP) ...
	IsSymbol  : symbolic character
	IsTitle   : title case letter

! To 将 r 映射到指定的大小写: UpperCase、LowerCase, TitleCase
	ToLower, ToTitle, ToUpper

! SimpleFold 迭代 Unicode 定义的简单大小写折叠下的 Unicode 码位。在相当于 rune 的码位中（包括 rune 本身），
	SimpleFold 返回最小的 rune > r（如果存在），否则返回最小的 rune >= 0。
	如果 r 不是有效的 Unicode 码位，SimpleFold(r) 返回r。
! unicode.CaseRange 表示用于简单大小写转换的 Unicode 码位范围 [Lo, Hi], 固定步长为 1
! unicode.SpecialCase 表示语言特定的大小写映射 []CaseRange
	ToLower, ToTitle, ToUpper
*/

func TestIsXxx(tt *testing.T) {
	// constant with mixed type runes
	const mixed = "\b5Ὂg̀9! ℃ᾭG৩"
	for _, c := range mixed {
		fmt.Printf("For %q:\n", c)
		if unicode.In(c, unicode.C, unicode.ASCII_Hex_Digit, unicode.Georgian) {
			fmt.Printf("\tis in rangeTable\n")
		}
		if unicode.IsControl(c) {
			fmt.Println("\tis control rune")
		}
		if unicode.IsDigit(c) {
			fmt.Println("\tis digit rune")
		}
		if unicode.IsGraphic(c) {
			fmt.Println("\tis graphic rune")
		}
		if unicode.IsLetter(c) {
			fmt.Println("\tis letter rune")
		}
		if unicode.IsLower(c) {
			fmt.Println("\tis lower case rune")
		}
		if unicode.IsMark(c) {
			fmt.Println("\tis mark rune")
		}
		if unicode.IsNumber(c) {
			fmt.Println("\tis number rune")
		}
		if unicode.IsPrint(c) {
			fmt.Println("\tis printable rune")
		}
		if !unicode.IsPrint(c) {
			fmt.Println("\tis not printable rune")
		}
		if unicode.IsPunct(c) {
			fmt.Println("\tis punct rune")
		}
		if unicode.IsSpace(c) {
			fmt.Println("\tis space rune")
		}
		if unicode.IsSymbol(c) {
			fmt.Println("\tis symbol rune")
		}
		if unicode.IsTitle(c) {
			fmt.Println("\tis title case rune")
		}
		if unicode.IsUpper(c) {
			fmt.Println("\tis upper case rune")
		}
	}
}

var Chars string = "ǅ5Ὂg̀9! ℃ᾭG৩gAaK\u212A1G"

func TestToXxx(t *testing.T) {
	toXxx := func(r rune) {
		fmt.Printf("%#U\nto upper: %#U\nto lower: %#U\nto title: %#U\n", r,
			unicode.To(unicode.UpperCase, r),
			unicode.To(unicode.LowerCase, r),
			unicode.To(unicode.TitleCase, r))
	}
	var err error
	var r rune
	var sr = strings.NewReader(Chars)
	for err := err; err == nil; {
		r, _, err = sr.ReadRune()
		toXxx(r)
	}
}

func TestSimpleFold(t *testing.T) {
	SimpleFold := func(r rune) {
		fmt.Printf("SimpleFold(%#U): %#U\n", r, unicode.SimpleFold(r))
	}
	var err error
	var r rune
	var sr = strings.NewReader(Chars)
	for err := err; err == nil; {
		r, _, err = sr.ReadRune()
		SimpleFold(r)
	}
}

func TestSpecialCase(t *testing.T) {
	SpecialCase := func(_case string, cs unicode.SpecialCase) {
		var sr = strings.NewReader(Chars)
		var err error
		var r rune
		for err := err; err == nil; {
			r, _, err = sr.ReadRune()
			fmt.Printf("%#U in case %s\nto upper: %#U\nto lower: %#U\nto title: %#U\n",
				r, _case, cs.ToUpper(r), cs.ToLower(r), cs.ToTitle(r))
		}
	}
	SpecialCase("CaseRanges", unicode.CaseRanges)
	SpecialCase("AzeriCase", unicode.AzeriCase)
	SpecialCase("TurkishCase", unicode.TurkishCase)
}
