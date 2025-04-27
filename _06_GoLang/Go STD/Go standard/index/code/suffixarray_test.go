package gostd

import (
	"bufio"
	"index/suffixarray"
	"os"
	"regexp"
	"testing"
)

/*
! New 为 data 创建新索引。对于 N=len(data)，索引创建时间是 O(N)
! Index 实现了一个后缀数组，用于快速子字符串搜索。
	Bytes 创建索引时所依据的数据。不得修改
	FindAllIndex 返回正则表达式 r 的非重叠匹配的排序列表，最多匹配 n 个；-1 表示 All
	Lookup 返回最多 n 个 Index 的未排序列表，其中字节串 s 出现在索引数据中；O(log(N)*len(s) + len(result))；-1 表示 All
	Read 从 r 读入 Index; Index 不能为 nil
	Write 将 Index 写入 w。
*/

func TestIndex(t *testing.T) {
	context := []byte("this is a content for suffixarray test")
	a := suffixarray.New(context)

	blankOff := a.Lookup([]byte{' '}, -1)
	log(blankOff)

	if r, err := regexp.Compile("[^ ]*s[^ ]*"); err == nil {
		indexex := a.FindAllIndex(r, -1)
		for _, i := range indexex {
			logfln("%s", a.Bytes()[i[0]:i[1]])
		}
	}
}

func TestIndexBytes(t *testing.T) {
	f, err := os.Open("suffixarray_test.go")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	bfr := bufio.NewReader(f)
	r, _ := regexp.Compile("^func Test.*")
	for {
		datas, err := bfr.ReadBytes('\n')
		if err != nil {
			break
		}
		a := suffixarray.New(datas)
		s := a.FindAllIndex(r, 1)
		if len(s) > 0 {
			logfln("%s", a.Bytes()[s[0][0]:s[0][1]-2])
		}
	}
}
