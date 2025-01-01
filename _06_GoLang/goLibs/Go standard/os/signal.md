<a id="TOP"></a>

## Package signal

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/signal_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<!-- <a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	 -->
	<a href="https://pkg.go.dev/os/signal"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `signal` 实现对输入信号的访问。信号主要用于类 Unix 系统

信号 `SIGKILL` 和 `SIGSTOP` 可能未被程序捕获，因此不受此包的影响。

同步信号是由程序执行错误触发的信号：`SIGBUS`、`SIGFPE` 和 `SIGSEGV`。只有当程序执行时，它们才会被视为同步，而不是使用 `os.Process.Kill` 或类 kill 程序机制。一般来说，除了下面讨论的之外，Go 程序会将同步信号转换为运行时崩溃。

其余信号为异步信号。它们不是由程序错误触发的，而是从内核或其他程序发送的。在异步信号中，当程序失去其控制终端时，会发送 `SIGHUP` 信号。当控制终端的用户按下中断字符时，将发送 `SIGINT` 信号，默认为 `^C`（Control-C）。当控制终端的用户按下退出字符时发送 `SIGQUIT` 信号，默认为 `^\` （Control-Backslash）。通常，您可以通过按 `^C` 使程序简单地退出，并且可以通过按 `^\` 使它退出堆栈转储。

默认情况下，同步信号将转换为运行时 panic。`SIGHUP`、`SIGINT` 或 `SIGTERM` 信号会导致程序退出。`SIGQUIT`、`SIGILL`、`SIGTRAP`、`SIGABRT`、`SIGSTKFLT`、`SIGEMT` 或 `SIGSYS` 信号会导致程序退出并出现堆栈转储。`SIGTSTP`、`SIGTTIN` 或 `SIGTTOU` 信号获取系统默认行为（这些信号由 shell 用于作业控制）。`SIGPROF` 信号由 Go 运行时直接处理以实现 `runtime.CPUProfile`。其他信号被捕获时不会采取任何行动。

如果在忽略 `SIGHUP` 或 `SIGINT` 的情况下启动 Go 程序（信号处理程序设置为 `SIG_IGN`），则它们将保持忽略状态。

如果 Go 程序以非空信号掩码启动，则通常会遵循该掩码。但是，有些信号是显式解锁的：同步信号 `SIGILL`、`SIGTRAP`、`SIGSTKFLT`、`SIGCHLD`、`SIGPROF`，以及在 Linux 上，信号 32（`SIGCANCEL`）和 33（`SIGSETXID`）（`SIGCANCEL` 和 `SIGSETXID` 由 glibc 内部使用）。由 `os` 启动的子进程。`exec` 或 `os/exec` 将继承修改后的信号掩码。

> 此包中的函数允许程序更改 Go 程序处理信号的方式。

`Notify` 禁用一组给定异步信号的默认行为，而是通过一个或多个注册通道传递这些信号。

如果程序在启动时忽略了 `SIGHUP` 或 `SIGINT`，并且对任一信号调用了 `Notify`，则将为该信号安装信号处理程序，并且不再忽略该信号。如果稍后对该信号调用了 `Reset` 或 `Ignore`，或者在传递给该信号的 `Notify` 的所有通道上调用了 `Stop`，则该信号将再次被忽略。`Reset` 将恢复信号的系统默认行为，而 `Ignore` 将导致系统完全忽略该信号。

如果程序使用非空信号掩码启动，则如上所述，某些信号将被显式解锁。如果对被阻塞的信号调用 `Notify`，它将被解除阻止。如果稍后对该信号调用 `Reset`，或者在传递给该信号的 `Notify` 的所有通道上调用 `Stop`，则该信号将再次被阻止。

在 Windows 上，`^C`（Control-C）或 `^BREAK`（Control-Break）通常会导致程序退出。如果 `Notify` 调用了 `os.Interrupt`，`^C` 或 `^BREAK` 将导致 `os.Interrupt` 在通道上发送，并且程序不会退出。如果在传递给 `Notify` 的所有通道上调用 `Reset` 或 `Stop`，则将恢复默认行为。

---
<a id="exam" ><a>

