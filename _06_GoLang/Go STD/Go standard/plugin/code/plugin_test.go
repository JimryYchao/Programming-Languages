package gostd

import (
	"plugin"
	"testing"
)

/*
! Open 打开一个 Go 插件。如果路径已经打开，则返回现有的 *Plugin。它对于多个 goroutine 并发使用是安全的。
! Plugin 是一个加载的 Go 插件。
	Lookup 在 plugin p 中搜索名为 symName 的 Symbol; Symbol 是任何导出的变量或函数。找不到时会报告一个错误。并发安全
! Symbol 是指向变量或函数的指针。
*/

func TestPluginWIN(t *testing.T) {
	p, err := plugin.Open("plugin_name.so")
	if err != nil {
		t.Fatal(err)
	}

	v, err := p.Lookup("V")
	if err != nil {
		t.Fatal(err)
	}
	f, err := p.Lookup("F")
	if err != nil {
		t.Fatal(err)
	}
	*v.(*int) = 7
	f.(func())() // prints "Hello, number 7"
}
