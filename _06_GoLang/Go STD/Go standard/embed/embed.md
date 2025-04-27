<a id="TOP"></a>

## Package embed

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/embed_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/embed" ><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `embed` 提供了对嵌入在运行的 Go 程序中的文件的访问。`import "embed"` 的 Go 源文件可以使用 `//go:embed` 指令初始化 `string`、`[]byte` 或 `embed.FS` 类型的包级变量，并使用编译时从包目录或子目录读取的文件内容。如果任何模式无效或具有无效匹配，则构建将失败。

有三种方式可以嵌入 `hello.txt` 的文件，然后在运行时打印其内容。

```go
import _ "embed"

//go:embed hello.txt
var s string
print(s)

//go:embed hello.txt
var b []byte
print(string(b))

//go:embed hello.txt
var f embed.FS
data, _ := f.ReadFile("hello.txt")
print(string(data))
```

`//go:embed` 指令使用一个或多个路径匹配模式指定要嵌入的文件。多文件嵌入仅对 `embed.FS` 有效。

```go
// content holds our static web server content.
//go:embed image/* template/*
//go:embed html/index.html
var content embed.FS
```

`FS` 实现了 `io/fs` 包的 `FS` 接口，它可以与任何理解文件系统的包一起使用，例如 `net/http`，`text/template` 和 `html/template`。



---
<a id="exam" ><a>