package gostd

import (
	"expvar"
	"math"
	"testing"
)

/*
! Do 对每个导出的变量调用 f
! Handle 返回 expvar HTTP。这只需要在非标准位置安装处理程序。
! Publish 声明一个命名的导出变量。这应该在包创建其 Vars 时从包的 init 函数调用。如果该名称已经注册，则将记录。
! Get 检索命名的导出变量
! expvar.Var 所有导出变量的抽象类型
! expvar.Float, NewFloat 是一个满足 Var 接口的 64 位浮点变量
	Add, Set, String, Value
! expvar.Func 通过调用函数并使用 JSON 格式化返回值来实现 Var
	String, Value
! expvar.Int, NewInt 是一个满足 Var 接口的 64 位整数变量
	Add, Set, String, Value
! expvar.Map, KeyValue, NewMap 是满足 Var 接口的字符串到 Var 映射变量
	Add, AddInt, AddFloat, Delete, Do, Get, Init, Set, String
! expvar.String, NewString 是一个满足 Var 接口的字符串变量
	Set, String, Value
*/

func TestExpVar(t *testing.T) {
	e := expvar.Get("e")
	log(e)
	if e != nil {
		e.(*expvar.Float).Set(1.234)
		log(e)
	}

	author := expvar.Get("Author")
	if author != nil {
		logfln("Author : %s", author)
	}
}

func init() {
	expvar.NewFloat("PI").Set(3.141592653)
	expvar.NewFloat("e").Set(math.E)
	expvar.NewInt("N").Set(10010)
	expvar.Publish("Author", &Name{"Jimry", "Ychao"})
}

type Name struct {
	First string
	Last  string
}

func (n *Name) String() string {
	return n.First + " " + n.Last
}
