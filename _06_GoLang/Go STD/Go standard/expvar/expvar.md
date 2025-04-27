<a id="TOP"></a>

## Package expvar

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/expvar_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<!-- <a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	 -->
	<a href="https://pkg.go.dev/expvar"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `expvar` 提供了一个公共变量的标准化接口，比如服务器中的操作计数器。它通过 HTTP/debug/vars 以 JSON 格式公开这些变量。

设置或修改这些公共变量的操作是原子操作。除了添加 HTTP 处理程序外，此包还注册以下变量：

```go
cmdline   os.Args
memstats  runtime.Memstats
```

有时候导入包只是为了注册它的 HTTP 处理程序和上述变量的副作用：

```go
import _ "expvar"
```

---
<a id="exam" ><a>