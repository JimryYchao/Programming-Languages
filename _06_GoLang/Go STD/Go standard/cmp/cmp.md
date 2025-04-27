<a id="TOP"></a>

## Package cmp

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/cmp_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<!-- <a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	 -->
	<a href="https://pkg.go.dev/cmp"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>


包 `cmp` 提供了与比较有序值相关的类型和函数。

```go
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}
```

---
<a id="exam" ><a>