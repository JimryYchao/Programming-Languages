<a id="TOP"></a>

## Package context

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/context_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/context" ><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>


包 `context` 定义了 Context 上下文类型，它在跨 API 边界和进程之间携带 *deadlines*、*cancellation signals* 和其他 *request-scoped values*。 

在 Go 服务器中，每个传入的请求都在自己的 *goroutine* 中处理。请求处理程序经常启动额外的 *goroutine* 来访问后端，如数据库和 RPC 服务。处理请求的 *goroutine* 集合通常需要访问特定于请求的值，例如最终用户的身份、授权令牌和请求的截止日期。当一个请求被取消或超时时，所有处理该请求的 *goroutine* 都应该迅速退出，这样系统就可以回收它们正在使用的任何资源。

对服务器的传入请求应该创建一个 `Context`，对服务器的传出调用应该接受一个 `Context`。它们之间的函数调用链必须传播 `Context`，可以使用 `WithCancel`、`WithDeadline`、`WithDelay` 或 `WithValue` 创建派生 `Context` 进行替换（可选）。当一个 `Context` 被取消时，所有从它派生的 `Context` 也被取消。

`WithCancel`、`WithDeadline` 和 `WithTimeout` 函数接受一个 `Context`（父级），并返回一个派生的 `Context`（子级）和一个 `CancelFunc`。调用 `CancelFunc` 将取消子对象及其所有子对象，删除父对象对该子对象的引用，并停止任何关联的计时器。如果不调用 `CancelFunc`，就会泄漏子对象及其子对象，直到父对象被取消或计时器触发。go vet tool 检查是否在所有控制流路径上使用了 `CancelFunc`s。

`WithCancelCause` 函数返回一个 `CancelCauseFunc`，它接受一个 `error` 并将其记录为取消原因。在取消的上下文或其任何子上下文上调用 `Cause` 将检索取消原因。如果未指定原因，则 `Cause(ctx)` 返回与 `ctx.Err()` 相同的值。

使用 `Context`s 的程序应该遵循这些规则，以保持包之间的接口一致性，并启用静态分析工具来检查上下文传播：不要将 `Context` 存储在结构类型中；相反，将 `Context` 显式传递给需要它的每个函数。`Context` 应该是第一个参数，通常命名为 ctx：

```go
func DoSomething(ctx context.Context, arg Arg) error {
	// ... use ctx ...
}
```

即使函数允许，也不要传递 `nil`。如果不确定使用哪个 `Context`，就传递 `context.TODO`。仅将 *context values* 用于传输进程和 API 请求作用域内的数据，而不是作为可选参数传递给函数。相同的 `Context` 可以传递给运行在不同 goroutine 中的函数；`Context` 可以安全地被多个 goroutine 同时使用。

---
### Context

`Done` 方法返回一个通道，作为代表 `Context` 运行的函数的取消信号：当通道关闭时，函数应该放弃它们的工作并返回。`Err` 方法返回一个错误，指出为什么取消了 `Context`。[Pipelines and Cancellation](https://golang.google.cn/blog/pipelines) 详细地讨论了 `Done` 通道习惯用法。

`Context` 没有 `Cancel` 方法，原因与 `Done` 通道是仅接收的原因相同：接收消除信号的函数通常不是发送信号的函数。特别是，当父操作为子操作启动 *goroutine* 时，这些子操作应该不能取消父操作。相反，`WithCancel` 函数（如下所述）提供了一种取消新的 `Context` 值的方法。

`Context` 可以安全地被多个 *goroutine* 同时使用。代码可以将一个 `Context` 传递给任意数量的 *goroutine*，并取消该 `Context` 以通知所有 goroutine。

`ctx.Deadline` 方法允许函数决定它们是否应该开始工作；如果剩下的时间太少，可能不值得。代码还可以使用一个 `deadline` 来设置 I/O 操作的超时。

`ctx.Value` 允许 `Context` 携带请求作用域的数据。这些数据必须能够安全地被多个 *goroutine* 同时使用。

>---
#### 派生 Context

包 `context` 提供了从现有值派生新的 `Context` 值的函数。这些值形成一个树：当一个 `Context` 被取消时，所有从它派生的 `Contexts` 也被取消。`Background` 是任何 `Context` 树的根；它永远不会被取消。

`WithCancel` 和 `WithTimeout` 返回派生的 `Context` 值，这些值可以比父级 `Context` 更快地被取消。与传入请求相关联的 `Context` 通常在请求处理程序返回时被取消。`WithCancel` 在使用多个副本时也可以用来取消冗余请求。`WithTimeout` 用于设置对后端服务器的请求的截止日期：

```go
type CancelFunc func()
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```

`WithValue` 提供了一种将请求作用域的值与 `Context` 关联的方法

---
#### Examples
<a id="exam" ><a>

- [Google Web Search](./examples/Google%20Web%20Search/Google%20Web%20Search.md)

---