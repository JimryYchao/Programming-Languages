<a id="TOP"></a>

## Package log

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/log_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/log" ><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>


包 `syslog` 为系统日志服务提供了一个简单的接口。它可以使用 UNIX 域套接字、UDP 或 TCP 向系统日志守护程序发送消息。仅需要一个调用一次 `Dial` 进行连接。写入失败时，系统日志客户端将尝试重新连接到服务器并再次写入。

TODO：包未在 Windows 上实现。


---