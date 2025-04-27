<a id="TOP"></a>

## Package exec

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/exec_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/os/exec"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `exec` 运行外部命令。它包装了 `os.StartProcess` 用以重新映射 `stdin` 和 `stdout`、将 I/O 与 Pipes 连接、以及执行其他调整。

`os/exec` 包有意不调用系统 shell，也不扩展任何 glob 模式或处理通常由 shell 完成的其他扩展、管道或重定向

要扩展 glob 模式，请直接调用 shell，注意转义任何危险输入，或者使用包 `path/filepath`的 `Glob` 函数。若要扩展环境变量，请使用包 `os` 的 `ExpandEnv`。

函数 `Command` 和 `LookPath` 按照主机操作系统的约定在当前路径列出的目录中查找程序；此包不再使用相对于当前目录的隐式或显式路径条目来解析程序。

---
<a id="exam" ><a>