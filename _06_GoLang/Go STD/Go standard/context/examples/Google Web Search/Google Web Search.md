### Google Web Search


示例是一个 HTTP 服务器，它通过将查询 “golang” 转发到 Google Web Search API 并呈现结果来处理像 "/search?q=golang&timeout=1s" 这样的 URL。`timeout` 参数告诉服务器在超时后取消请求。

代码分为三个包：
- `server` 提供 `main` 函数和 `/search` 的处理程序。
- `userip` 提供了从请求中提取用户 IP 地址并将其与 `Context` 关联的函数。
- `google` 提供了 `Search` 功能，用于向 `Google` 发送查询。

>---
#### [server](./server.go)

`server` 程序处理类似于 "/search?q=golang" 的请求，通过提供 `golang` 的前几个 Google 搜索结果。它注册 `handleSearch` 以处理 "/search" 端点。处理程序创建一个名为 `ctx` 的初始 `Context`，并安排在处理程序返回时取消它。如果请求中包含 `timeout` URL 参数，则在超时后，`Context` 会自动取消：

```go
func handleSearch(w http.ResponseWriter, req *http.Request) {
    // ctx is the Context for this handler. Calling cancel closes the
    // ctx.Done channel, which is the cancellation signal for requests
    // started by this handler.
    var (
        ctx    context.Context
        cancel context.CancelFunc
    )
    timeout, err := time.ParseDuration(req.FormValue("timeout"))
    if err == nil {
        // The request has a timeout, so create a context that is
        // canceled automatically when the timeout expires.
        ctx, cancel = context.WithTimeout(context.Background(), timeout)
    } else {
        ctx, cancel = context.WithCancel(context.Background())
    }
    defer cancel() // Cancel ctx as soon as handleSearch returns.
    ...
}
```

处理程序从请求中提取查询，并通过调用 `userip` 包提取客户端的 IP 地址。后端请求需要客户端的 IP 地址，因此 `handleSearch` 将其附加到 `ctx`：

```go
    ...
    // Check the search query.
    query := req.FormValue("q")
    if query == "" {
        http.Error(w, "no query", http.StatusBadRequest)
        return
    }

    // Store the user IP in ctx for use by code in other packages.
    userIP, err := userip.FromRequest(req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    ctx = userip.NewContext(ctx, userIP)
    ...
```

处理程序使用 `ctx` 和 `query` 调用 `google.Search`：

```go
    ...
    // Run the Google search and print the results.
    start := time.Now()
    results, err := google.Search(ctx, query)
    elapsed := time.Since(start)
    ...
```

如果搜索成功，处理程序将呈现结果：

```go
    ...
    if err := resultsTemplate.Execute(w, struct {
        Results          google.Results
        Timeout, Elapsed time.Duration
    }{
        Results: results,
        Timeout: timeout,
        Elapsed: elapsed,
    }); err != nil {
        log.Print(err)
        return
    }
}
```

>---
#### [userip](./userip/userip.go) 

`userip` 包提供了从请求中提取用户 IP 地址并将其与 `Context` 关联的函数。`Context` 提供了一个键值映射，其中键和值都是类型 `interface{}`。键类型必须支持相等，值必须并发安全（可供多个 *goroutine* 同时使用）。像 `userip` 这样的包隐藏了这个映射的细节，并提供了对特定 `Context` 值的强类型访问。

为了避免键冲突，`userip` 定义了一个未导出的类型 `key`，并使用该类型的值作为上下文键：

```go
// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int
// userIPkey is the context key for the user IP address.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const userIPKey key = 0
```

`FromRequest` 从 `http.Request` 中提取 `userIP` 值：

```go
func FromRequest(req *http.Request) (net.IP, error) {
    ip, _, err := net.SplitHostPort(req.RemoteAddr)
    if err != nil {
        return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
    }
    ...
```

`NewContext` 返回一个新的 `Context`，它携带一个提供的 `userIP` 值：

```go
func NewContext(ctx context.Context, userIP net.IP) context.Context {
    return context.WithValue(ctx, userIPKey, userIP)
}
```

`FromContext` 从 `Context` 中提取 `userIP`：

```go
func FromContext(ctx context.Context) (net.IP, bool) {
    // ctx.Value returns nil if ctx has no value for the key;
    // the net.IP type assertion returns ok=false for nil.
    userIP, ok := ctx.Value(userIPKey).(net.IP)
    return userIP, ok
}
```

>---
#### [google](./google/google.go)

`google.Search` 函数向 Google Web Search API 发出 HTTP 请求，并解析 JSON 编码的结果。它接受一个 `Context` 参数 `ctx`，如果在请求进行中关闭了 `ctx.Done`，则立即返回。

Google Web Search API 请求包括搜索查询和用户 IP 作为查询参数：

```go
func Search(ctx context.Context, query string) (Results, error) {
    // Prepare the Google Search API request.
    req, err := http.NewRequest("GET", "https://ajax.googleapis.com/ajax/services/search/web?v=1.0", nil)
    if err != nil {
        return nil, err
    }
    q := req.URL.Query()
    q.Set("q", query)

    // If ctx is carrying the user IP address, forward it to the server.
    // Google APIs use the user IP to distinguish server-initiated requests
    // from end-user requests.
    if userIP, ok := userip.FromContext(ctx); ok {
        q.Set("userip", userIP.String())
    }
    req.URL.RawQuery = q.Encode()
    ...
```

`Search` 使用辅助函数 `httpDo` 来发出 HTTP 请求，如果在处理请求或响应时 `ctx.Done` 关闭，则取消该请求。`Search` 传递一个闭包给 `httpDo` 处理 HTTP 响应：

```go
    ...
    var results Results
    err = httpDo(ctx, req, func(resp *http.Response, err error) error {
        if err != nil {
            return err
        }
        defer resp.Body.Close()

        // Parse the JSON search result.
        // https://developers.google.com/web-search/docs/#fonje
        var data struct {
            ResponseData struct {
                Results []struct {
                    TitleNoFormatting string
                    URL               string
                }
            }
        }
        if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
            return err
        }
        for _, res := range data.ResponseData.Results {
            results = append(results, Result{Title: res.TitleNoFormatting, URL: res.URL})
        }
        return nil
    })
    // httpDo waits for the closure we provided to return, so it's safe to
    // read results here.
    return results, err
}
```

`httpDo` 函数运行 HTTP 请求，并在一个新的 *goroutine* 中处理其响应。如果 `ctx.Done` 在 *goroutine* 退出之前关闭，它会取消请求：

```go
func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
    // Run the HTTP request in a goroutine and pass the response to f.
    c := make(chan error, 1)
    req = req.WithContext(ctx)
    go func() { c <- f(http.DefaultClient.Do(req)) }()
    select {
    case <-ctx.Done():
        <-c // Wait for f to return.
        return ctx.Err()
    case err := <-c:
        return err
    }
}
```

>---
#### Conclusion

想要在 `Context` 上构建的服务器框架应该提供 `Context` 的实现，以在它们的包和那些需要 `Context` 参数的包之间架起桥梁。然后，它们的客户端库将从调用代码中接受 `Context`。通过为请求范围的数据和取消建立一个公共接口，`Context`  使包开发人员更容易共享创建可伸缩服务的代码。

---