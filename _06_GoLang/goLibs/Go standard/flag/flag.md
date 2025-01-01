<a id="TOP"></a>

## Package XXX

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/XXX_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/XXX"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `flag` 实现命令行标志解析。使用 `flag.String, Bool, Int, ...` 等定义标记。或 `VarT` 将标志绑定到变量；

```go
var ip = flag.Int("n", 1234, "help message for flag n")

flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
// or
flag.Var(&flagVal, "name", "help message for flagname")  // flagVal 实现 flag.Value 接口
```

调用 `flag.Parse` 将命令行解析为定义的标志。

```go
flag.Parse()
// use flags
fmt.Println("ip has value ", *ip)
fmt.Println("flagvar has value ", flagvar)
```

标志后面的参数可以作为切片 `flag.Args` 或单独作为 `flag.Arg(i)`。参数的索引范围是从 0 到 `flag.NArg-1`。

命令行标志可以使用以下形式：

```shell
$ cmd -flag
$ cmd --flag   // double dashes are also permitted
$ cmd -flag=x
$ cmd -flag x  // non-boolean flags only

$ go test -v -run=Test
```

整数标志可以是 `1234`, `077`, `-0x1234` 的形式；布尔标志可以是 `1`, `0`, `t`, `f`, `T`, `F`, `true`, `false`, `TRUE`, `FALSE`, `True`, `False`；

`Duration` 标志接受任何的对 `time.ParseDuration` 有效的输入。

默认的命令行标志集由包级函数控制；`FlagSet` 允许定义独立的标志集，例如在命令行界面中实现子命令。`FlagSet` 的方法类似于命令行标志集的顶级函数。 

`flag.Usage()` 将所有定义的命令行标志的使用情况消息打印到 CommandLine's output：

```go
package main
import "flag"

func main() {
	flag.Int("n", 1234, "help message for flag n")
	flag.String("s", "flag", "help message for flag s")
	flag.Usage()
}
/*
Usage of test.exe:
  -n int
        help message for flag n (default 1234)
  -s string
        help message for flag s (default "flag")
*/
```

---
<a id="exam" ><a>