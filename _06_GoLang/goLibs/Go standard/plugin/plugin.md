<a id="TOP"></a>

## Package plugin

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/plugin_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/plugin"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `plugin` 实现 Go 插件的加载和 symbol 符号解析。

插件是一个 Go 主包，包含导出的函数和变量；构建时使用 `go build -buildmode=plugin`。当一个插件第一次打开时，所有嵌入的包的 `init` 函数都会被调用。`main` 函数没有运行。插件只初始化一次，不能关闭。

`Open` 加载一个路径下的插件 `Plugin`，`plugin.Lookup` 查找特定名称的 `Symbol`；`Symbol` 是指向变量或函数的指针。

```go
// 例如定义
package main
import "fmt"
var V int
func F() { fmt.Printf("Hello, number %d\n", V) }
```
```go
// Open
p, err := plugin.Open("plugin_name.so")
if err != nil {
	panic(err)
}
v, err := p.Lookup("V")
if err != nil {
	panic(err)
}
f, err := p.Lookup("F")
if err != nil {
	panic(err)
}
*v.(*int) = 7
f.(func())() // prints "Hello, number 7"
```

---
<a id="exam" ><a>	