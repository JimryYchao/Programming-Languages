<a id="TOP"></a>

## Package bufio

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/bufio_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<!-- <a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a> -->	
	<a href="https://pkg.go.dev/bufio" ><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `bufio` 实现了缓冲 I/O。它包装一个 `io.Reader` 或 `io.Writer` 对象，创建了另一个对象（`Reader` 或 `Writer`），该对象也实现了该接口，但提供了缓冲和一些文本 I/O 帮助。
   
---
<a id="exam" ><a>

### Examples

- [Custom Scan File](examples/customScanFile_test.go)

---