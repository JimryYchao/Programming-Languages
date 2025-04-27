<a id="TOP"></a>

## Package log

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/log_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/log" ><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>


包 `log` 实现了一个简单的日志记录包。每个日志消息都在单独的一行中输出。`Fatal` 函数在写入日志消息后调用 `os.exit(1)`。`Panic` 函数在写入日志消息后调用 `panic()`。



---