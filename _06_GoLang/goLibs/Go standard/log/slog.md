<a id="TOP"></a>

## Package slog

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/slog_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/log/slog" ><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>


包 `slog` 提供结构化日志记录，其中日志记录包括消息、严重性级别和以键值对表示的各种其他属性。

`slog` 定义了 Debug、Info、Warn、Error 四种常用级别。在调用 `SetDefault` 之前，slog.defaultLogger 关联使用 log.defaultLogger 进行包级日志函数的输出，`SetLogLoggerLevel` 同时控制 `slog.defaultLogger` 和默认的 log.Logger 的日志级别（默认为 Info）。调用 `SetLogLoggerLevel(slog.LevelDebug)` 后，对 `slog.Debug` 的调用将传递给 log.defaultLogger。

```go
func TestSetLoggerLevel(t *testing.T) {
	beforeTest(t)
	slog.Info("test slog info")				// output
	slog.Debug("this mess will not output") // not output
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Debug("set logger level to Debug") // output
}
```

每一个 `Logger` 关联一个 `Handler`，`Logger` 的日志方法从参数中创建一个 `Record`，并将其传递给 `Handler` 进行处理。日志记录由 time、level、massage 和一组键值对组成。`Handler` 提供了默认的 `TextHandler` 和 `JSONHandler`。

```go
slog.						// time of call 
	Info(					// level of Info
	"Hello", 				// msg of "Hello"
	"count", 3)				// kv-pair of ("count", 3)
// 2024/01/01 01:23:45 INFO Hello count=3

txtSlgr := slog.New(slog.NewTextHandler(os.Stdout, nil))
txtSlgr.Info("Hello Text")
// time=2024-05-25T16:40:25.921+08:00 level=INFO msg="Hello Text"

jsonSlgr := slog.New(slog.NewJSONHandler(os.Stdout, nil))
jsonSlgr.Info("Hello Json")
// {"time":"2024-05-25T16:14:02.5273149+08:00","level":"INFO","msg":"Hello Json"}
```


`SetDefault` 用来设置默认的 slogLogger 的 `Handler`。在调用 `SetDefault` 之后，`log.Print` 及相关函数的日志记录也会发送到日志 Handler 处理程序。

```go
slog.Info("test SetDefault")
// 2024/05/25 16:46:24 INFO test SetDefault
slog.SetDefault(jsonSlgr)
slog.Info("output json format")
// {"time":"2024-05-25T16:46:24.9309425+08:00","level":"INFO","msg":"output json format"}

log.Print("Call log.Print")
// {"time":"2024-05-25T16:58:25.3666103+08:00","level":"INFO","msg":"Call log.Print"}
//! log.Fatal("Call log.Fatal")
// {"time":"2024-05-25T17:02:48.5319853+08:00","level":"INFO","msg":"Call log.Fatal"}
//! log.Panic("Call log.Panic")
// {"time":"2024-05-25T17:03:07.7079422+08:00","level":"INFO","msg":"Call log.Panic"}
```

`TextHandler` 和 `JSONHandler` 都可以使用 `HandlerOptions` 进行配置。例如设置最低级别、显示日志调用的源文件和行，以及在记录属性之前修改属性。

```go
le := &slog.LevelVar{} // 0 : Info
lgr := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: le}))
lgr.Info("print test")
// {"time":"2024-05-26T22:41:48.3181011+08:00","level":"INFO","msg":"print test"}
le.Set(slog.LevelError)
lgr.Info("not printed")
lgr.Error("with attrs", slog.Group("G", slog.String("k1", "v1"), slog.String("k2", "v2")))
// {"time":"2024-05-26T22:41:48.3278887+08:00","level":"ERROR","msg":"with attrs","G":{"k1":"v1","k2":"v2"}}
```

---
<a id="exam" ><a>
### Examples

- [LevelVarHandler](./examples/LevelVarHandler.go)

- [LevelVarSLogger](./examples/LevelVarSLogger.go)

- [slog-handler-guide](./examples/slog-handler-guide/README.md) <a href="https://github.com/golang/example/blob/master/slog-handler-guide/README.md"><img src="../_rsc/link-src.drawio.png" id="other"/>
---