# A Guide to Writing `slog` Handlers

- [A Guide to Writing `slog` Handlers](#a-guide-to-writing-slog-handlers)
	- [Introduction](#introduction)
	- [Loggers and their handlers](#loggers-and-their-handlers)
	- [Implementing `Handler` methods](#implementing-handler-methods)
		- [The `Enabled` method](#the-enabled-method)
		- [The `Handle` method](#the-handle-method)
		- [The `WithAttrs` method](#the-withattrs-method)
		- [The `WithGroup` method](#the-withgroup-method)
			- [Without pre-formatting](#without-pre-formatting)
			- [Getting the mutex right](#getting-the-mutex-right)
			- [With pre-formatting](#with-pre-formatting)
		- [Testing](#testing)
	- [General considerations](#general-considerations)
		- [Copying records](#copying-records)
		- [Concurrency safety](#concurrency-safety)
		- [Robustness](#robustness)
		- [Speed](#speed)

---
## Introduction

标准库的 `log/slog` 包有两部分设计。由 `Logger` 类型实现的 “前端” 收集结构化的日志信息，如消息、级别和属性，并将它们传递给 “后端”，即 `Handler` 接口的实现。这个包带有两个内置的 *handler*（`TextHandler` 和 `JSONHandler`）。这本指南是来帮助编写自己的 *handler*。

---
## Loggers and their handlers

编写 *handler* 需要理解 `Logger` 和 `Handler` 类型如何一起工作。

每个 *logger* 都包含一个 *handler*。某些 `Logger` 方法会做一些前期工作，比如将键值对收集到 `Attr`s 中，然后调用一个或多个 `Handler` 方法。这些 `Logger` 方法是 `With`、`WithGroup` 和输出方法。

输出方法实现了 *logger* 的主要作用：生成日志输出。下面是对一个输出方法的调用：

    logger.Info("hello", "key", value)

有两种通用的输出方法，`Log` 和 `LogAttrs`。也存在分别用于四个公共级别（`Debug`、`Info`、`Warn` 和 `Error`）对应的输出方法，以及采用 *context* 的对应方法（`DebugContext`、`InfoContext`、`WarnContext` 和 `ErrorContext`）。

每个 `Logger` 输出方法首先调用其 *handler* 的 `Enabled` 方法。如果该调用返回 `true`，则该方法从其参数构造 `Record` 并调用 *handler* 的 `Handle` 方法。

为了方便和优化，属性可以通过调用 `With` 方法添加到 `Logger` 中：

    logger = logger.With("k", v)

这个调用创建了一个新的带有参数属性的 `Logger` 值；原始值保持不变。*logger* 的所有后续输出都将包括这些属性。*logger* 的 `With` 方法调用它的 *handler* 的 `WithAttrs` 方法。

`WithGroup` 方法用于通过建立单独的命名空间来避免大型程序中的键冲突。这个调用创建了一个新的 `Logger` 值，组名为 “g”：

    logger = logger.WithGroup("g")

*logger* 的所有后续键将由组名 “g” 限定。“限定” 的确切含义取决于 *logger* 的 *handler* 如何格式化输出。内置的 `TextHandler` 将组视为键的前缀，例如 g.k。内置的 `JSONHandler` 使用组作为嵌套 JSON 对象的键：

    {"g": {"k": v}}

*logger* 的 `WithGroup` 方法调用它的 *handler* 的 `WithGroup` 方法。

---
## Implementing `Handler` methods

我们现在可以详细讨论这四个 `Handler` 方法。我们将编写一个 *handler*，使用类似 YAML 的格式来格式化日志。它将显示此日志输出调用：

    logger.Info("hello", "key", 23)

输出如:

    time: 2023-05-15T16:29:00
    level: INFO
    message: "hello"
    key: 23
    ---

虽然这个特定的输出是有效的 YAML，但我们的实现没有考虑 YAML 语法的微妙之处，因此有时会产生无效的 YAML。例如，它不引用其中有冒号的键这个 *handler* 命名为 `IndentHandler`。

我们从 `IndentHandler` 类型和使用 `io.Writer` 和 `Options` 构造它的 `New` 函数开始：

```go
type IndentHandler struct {
	opts Options
	// TODO: state for WithGroup and WithAttrs
	mu  *sync.Mutex
	out io.Writer
}

type Options struct {
	// Level reports the minimum level to log.
	// Levels with lower levels are discarded.
	// If nil, the Handler uses [slog.LevelInfo].
	Level slog.Leveler
}

func New(out io.Writer, opts *Options) *IndentHandler {
	h := &IndentHandler{out: out, mu: &sync.Mutex{}}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	return h
}
```
我们将只支持一个选项，即设置最低级别以抑制详细日志输出的能力。*handlers* 应始终将此选项声明为 `slog.Leveler`。`slog.Leveler` 接口由 `Level` 和 `LevelVar` 实现。用户很容易提供 `Level` 值，但是更改多个 *handlers* 的级别需要跟踪它们。如果用户传递了一个 `LevelVar`，那么对 `LevelVar` 的一次修改将改变包含它的所有 *handlers* 的行为。对 `LevelVar` 的修改是 *goroutine-safe* 的。

还可以考虑向 *handler* 添加一个 `ReplaceAttr` 选项，就像 [内置 *handlers*]((https://pkg.go.dev/log/slog#HandlerOptions.ReplaceAttr)) 的选项一样。虽然 `ReplaceAttr` 会使实现复杂化，但它也会使 *handler* 更有用。

互斥锁将用于确保对 `io.Writer` 的写入以原子方式发生。不寻常的是，`IndentHandler` 持有一个指向 `sync.Mutex` 的指针，而不是直接持有一个 `sync.Mutex`。有一个很好的理由，我们 [稍后](#getting-the-mutex-right) 会解释。

我们的 *handler* 需要额外的状态来跟踪对 `WithGroup` 和 `WithAttrs` 的调用。我们将在讨论这些方法时描述这种状态。

### The `Enabled` method

`Enabled` 方法是一种优化，可以避免不必要的工作。一个 `Logger` 输出方法会在处理任何参数之前调用 `Enabled`，看看是否应该继续。

    Enabled(context.Context, Level) bool

*context* 可用于允许基于上下文信息的决策。例如，自定义 HTTP 请求标头可以指定最低级别，服务器将其添加到用于处理该请求的上下文中。*handler* 的 `Enabled` 方法可以报告参数级别是否大于或等于上下文值，从而允许独立控制每个请求所做工作的详细程度。

我们的 `IndentHandler` 不使用上下文。它只是将参数级别与其配置的最小级别进行比较：

```go
func (h *IndentHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}
```

>---
### The `Handle` method

向 `Handle` 方法传递一个 `Record`，其中包含对 `Logger` 输出方法的单个调用要记录的所有信息。`Handle` 方法应该以某种方式处理它。一种方法是以某种格式输出 `Record`，就像 `TextHandler` 和 `JSONHandler` 一样。其他的一些选择是修改 `Record` 并将其传递给另一个 *handler*，或将 `Record` 排队等待以后处理，或者忽略它。

    Handle(context.Context, Record) error

提供上下文是为了支持沿着调用链提供日志记录信息的应用程序。与通常的 Go 实践不同，`Handle` 方法不应该将取消的上下文视为停止工作的信号。

如果 `Handle` 处理它的 `Record`，它应该遵循 [文档](https://pkg.go.dev/log/slog#Handler.Handle) 中的规则。例如，零 `Time` 字段应被忽略，零 `Attr`s 字段也应被忽略。

要生成输出的 `Handle` 方法应该执行以下步骤：

1. 分配一个缓冲区（通常为 `[]byte`）来保存输出。最好先在内存中构造，然后每次调用 `io.Writer.Write` 时将其写入，以最大限度地减少与使用同一个 *writer* 的其他 *goroutine* 的交织。

2. 设置特殊字段的格式：time, level, message, source location (PC)。作为一般规则，这些字段应该首先出现，并且不嵌套在由 `WithGroup` 建立的组中。
3. 格式化 `WithGroup` 和 `WithAttrs` 调用的结果。

4. 格式化 `Record` 中的属性。

5. 输出缓冲区。

以下是 `IndentHandler.Handle` 的结构：

```go
func (h *IndentHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)
	if !r.Time.IsZero() {
		buf = h.appendAttr(buf, slog.Time(slog.TimeKey, r.Time), 0)
	}
	buf = h.appendAttr(buf, slog.Any(slog.LevelKey, r.Level), 0)
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = h.appendAttr(buf, slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", f.File, f.Line)), 0)
	}
	buf = h.appendAttr(buf, slog.String(slog.MessageKey, r.Message), 0)
	indentLevel := 0
	// TODO: output the Attrs and groups from WithAttrs and WithGroup.
	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, a, indentLevel)
		return true
	})
	buf = append(buf, "---\n"...)
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}
```

第一行分配了一个 `[]byte`，对于大多数日志输出来说，它应该足够大。为缓冲区分配一些初始的、相当大的容量是一个简单但重要的优化：它避免了在初始切片为空或很小时发生的重复复制和分配。我们将在关于 [Speed](#speed) 的章节中回到这一行，并展示我们如何做得更好。

我们的 `Handle` 方法的下一部分格式化特殊属性，遵守忽略零 time 和零 PC 的规则。

接下来，该方法处理 `WithAttrs` 和 `WithGroup` 调用的结果。我们先跳过这个。

然后处理参数记录中的属性。我们使用 `Record.Attrs` 方法按照用户将属性传递给 `Logger` 输出方法的顺序来遍历属性。*handler* 可以自由地对属性进行重新排序或重复数据删除，但我们的不能。

最后，在将 `"---"` 行添加到输出中以分隔日志记录之后，我们的 *handler* 使用我们积累的缓冲区对 `h.out.Write` 进行单个调用。我们持有这个写操作的锁，以使它对于可能同时调用 `Handle` 的其他 *goroutine* 来说是原子的。

*handler* 的核心是 `appendAttr` 方法，负责格式化单个属性：

```go
func (h *IndentHandler) appendAttr(buf []byte, a slog.Attr, indentLevel int) []byte {
	// Resolve the Attr's value before doing anything else.
	a.Value = a.Value.Resolve()
	// Ignore empty Attrs.
	if a.Equal(slog.Attr{}) {
		return buf
	}
	// Indent 4 spaces per level.
	buf = fmt.Appendf(buf, "%*s", indentLevel*4, "")
	switch a.Value.Kind() {
	case slog.KindString:
		// Quote string values, to make them easy to parse.
		buf = fmt.Appendf(buf, "%s: %q\n", a.Key, a.Value.String())
	case slog.KindTime:
		// Write times in a standard way, without the monotonic time.
		buf = fmt.Appendf(buf, "%s: %s\n", a.Key, a.Value.Time().Format(time.RFC3339Nano))
	case slog.KindGroup:
		attrs := a.Value.Group()
		// Ignore empty groups.
		if len(attrs) == 0 {
			return buf
		}
		// If the key is non-empty, write it out and indent the rest of the attrs.
		// Otherwise, inline the attrs.
		if a.Key != "" {
			buf = fmt.Appendf(buf, "%s:\n", a.Key)
			indentLevel++
		}
		for _, ga := range attrs {
			buf = h.appendAttr(buf, ga, indentLevel)
		}
	default:
		buf = fmt.Appendf(buf, "%s: %s\n", a.Key, a.Value)
	}
	return buf
}
```

它首先解析属性，如果属性有值，则运行该值的 `LogValuer.LogValue` 方法。所有 *handler* 都应该解析它们处理的每个属性。

接下来，它遵循 *handler* 规则，即应忽略空属性。

然后它 `switch` 属性的 `Kind` 以确定要使用的格式。对于大多数类型（`switch default`），它依赖于 `slog.Value` 的 `String` 方法来产生合理的结果。它专门处理 string 和 time：`KindString` 引用它们的 `value.String()`，`KindTime` 以标准方式 `value.Time().Format(time.RFC3339Nano)` 格式化它们。

当 `appendAttr` 中 *switch* 到 `KindGroup` 时，它会在应用一些*handler*规则后，递归地调用组的属性。首先，不带属性的组被忽略 — 甚至其键也不显示。其次，一个空键的组是内联的：组的边界不会被标记。在我们的例子中，这意味着组的属性没有缩进。

>---
### The `WithAttrs` method

`slog` 的性能优化之一是支持预格式化属性。`Logger.With` 方法将键值对转换为 `Attr`，然后调用 `Handler.WithAttrs`。*handlers* 可以存储这些属性供 `Handle` 方法以后使用，或者它可以利用这个机会现在就格式化这些属性，而不是在每次调用 `Handle` 时重复这样做。

    WithAttrs(attrs []Attr) Handler

属性是传递给 `Logger.With` 的已处理键值对。返回值应该是 *handler* 的一个新实例，它包含属性，可能是预先格式化的。

`WithAttrs` 必须返回一个带有附加属性的新 *handler*，而原始 *handler*（它的接收者）保持不变。例如，这个调用：

    logger2 := logger1.With("k", v)

创建一个新的 *logger* `logger2` ，带有一个附加属性，但对 `logger1` 没有影响。

当我们讨论 `WithGroup` 时，我们将在下面展示 `WithAttrs` 的示例实现。

>---
### The `WithGroup` method

`Logger.WithGroup` 直接调用 `Handler.WithGroup`，使用相同的参数，组名。*handler* 应该记住名称，以便可以使用它来限定所有后续属性。

    WithGroup(name string) Handler

与 `WithAttrs` 类似，`WithGroup` 方法应该返回一个新的 *handler*，而不是修改接收方。

`WithGroup` 和 `WithAttrs` 的实现是交织在一起的。考虑一下这个声明：

    logger = logger.WithGroup("g1").With("k1", 1).WithGroup("g2").With("k2", 2)

随后的 `logger` 输出应该用组 “g1” 限定键 “k1”，用组 “g1” 和 “g2” 限定键 “k2”。`Handler.WithGroup` 和 `Handler.WithAttrs` 的实现必须遵守 `Logger.WithGroup` 和 `Logger.With` 调用的顺序。

我们将研究 `WithGroup` 和 `WithAttrs` 的两种实现，一种是预格式化的，另一种不是。

#### Without pre-formatting

我们的第一个实现将从 `WithGroup` 和 `WithAttrs` 调用中收集信息，以构建一个组名和属性列表切片，并在 `Handle` 中循环该切片。我们从一个结构体开始，它可以保存一个组名或一些属性：

```go
// groupOrAttrs holds either a group name or a list of slog.Attrs.
type groupOrAttrs struct {
	group string      // group name if non-empty
	attrs []slog.Attr // attrs if non-empty
}
```

然后我们将 `groupOrAttrs` 的切片添加到我们的 *handler* 中：

```go
type IndentHandler struct {
	opts Options
	goas []groupOrAttrs
	mu   *sync.Mutex
	out  io.Writer
}
```

如上所述，`WithGroup` 和 `WithAttrs` 方法不应该修改它们的接收器。为此，我们定义了一个方法，它将复制我们的 *handler* 结构并向副本追加一个 `groupOrAttrs`：

```go
func (h *IndentHandler) withGroupOrAttrs(goa groupOrAttrs) *IndentHandler {
	h2 := *h
	h2.goas = make([]groupOrAttrs, len(h.goas)+1)
	copy(h2.goas, h.goas)
	h2.goas[len(h2.goas)-1] = goa
	return &h2
}
```

`IndentHandler` 的大部分字段都可以浅拷贝，但是 `groupOrAttrs` 的切片需要深拷贝，否则克隆和原始将指向相同的底层数组。如果我们使用 `append` 而不是显式复制，我们会引入微妙的别名错误。

使用 `withGroupOrAttrs`，`With` 方法很容易：

```go
func (h *IndentHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}

func (h *IndentHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}
```

`Handle` 方法现在可以在内置属性之后和记录中的属性之前处理 `groupOrAttrs` 切片：

```go
func (h *IndentHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)
	if !r.Time.IsZero() {
		buf = h.appendAttr(buf, slog.Time(slog.TimeKey, r.Time), 0)
	}
	buf = h.appendAttr(buf, slog.Any(slog.LevelKey, r.Level), 0)
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = h.appendAttr(buf, slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", f.File, f.Line)), 0)
	}
	buf = h.appendAttr(buf, slog.String(slog.MessageKey, r.Message), 0)
	indentLevel := 0
	// Handle state from WithGroup and WithAttrs.
	goas := h.goas
	if r.NumAttrs() == 0 {
		// If the record has no Attrs, remove groups at the end of the list; they are empty.
		for len(goas) > 0 && goas[len(goas)-1].group != "" {
			goas = goas[:len(goas)-1]
		}
	}
	for _, goa := range goas {
		if goa.group != "" {
			buf = fmt.Appendf(buf, "%*s%s:\n", indentLevel*4, "", goa.group)
			indentLevel++
		} else {
			for _, a := range goa.attrs {
				buf = h.appendAttr(buf, a, indentLevel)
			}
		}
	}
	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, a, indentLevel)
		return true
	})
	buf = append(buf, "---\n"...)
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}
```

你可能已经注意到，我们记录 `WithGroup` 和 `WithAttrs` 信息的算法在调用这些方法的次数上是二次的，因为重复复制。这在实践中不太可能有什么关系，但如果它困扰你，你可以使用一个链表来代替，`Handle` 将不得不反转或递归访问。有关实现，请参阅 [github.com/jba/slog/withsupport](https://github.com/jba/slog/tree/main/withsupport) 包。

#### Getting the mutex right

让我们重温 `Handle` 的最后几行：

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
    return err

这段代码没有改变，但我们现在可以理解为什么 `h.mu` 是一个指向 `sync.Mutex` 的指针。`WithGroup` 和 `WithAttrs` 都复制 *handler*。两个副本都指向同一个互斥体。如果副本和原始副本使用不同的互斥锁，并且是并发使用的，那么它们的输出可能会交错，或者某些输出可能会丢失。代码如下：

    l2 := l1.With("a", 1)
    go l1.Info("hello")
    l2.Info("goodbye")

可以产生这样的输出：

    hegoollo a=dbye1

请参阅此 [this bug report](https://go.dev/issue/61321) 以了解更多详细信息。

#### With pre-formatting

我们的第二种实现实现了预格式化。这个实现比前一个更复杂。额外的复杂性值得吗？这取决于你的情况，但有一种情况可能是这样。假设你希望服务器记录有关传入请求的大量信息，以及该请求期间发生的每个日志消息。一个典型的 *handler* 可能看起来像这样：

    func (s *Server) handleWidgets(w http.ResponseWriter, r *http.Request) {
        logger := s.logger.With(
            "url", r.URL,
            "traceID": r.Header.Get("X-Cloud-Trace-Context"),
            // many other attributes
            )
        // ...
    }

单个 `handleWidgets` 请求可能会生成数百行日志。例如，它可能包含这样的代码：

    for _, w := range widgets {
        logger.Info("processing widget", "name", w.Name)
        // ...
    }

对于每一行，我们上面写的 `Handle` 方法将格式化所有使用上面的 `With` 添加的属性，以及日志行本身上的属性。

也许所有这些额外的工作并不会显著降低服务器的速度，因为它做了太多的其他工作，花费在日志记录上的时间只是皮毛。但是，也许你的服务器速度足够快，所有额外的格式显示在你的 CPU 配置文件中。也就是说，通过在对 `With` 的调用中格式化属性一次，预格式化可以产生很大的不同。

要将参数预格式化为 `WithAttrs`，我们需要跟踪 `IndentHandler` 结构中的一些额外状态。

```go
type IndentHandler struct {
	opts           Options
	preformatted   []byte   // data from WithGroup and WithAttrs
	unopenedGroups []string // groups from WithGroup that haven't been opened
	indentLevel    int      // same as number of opened groups so far
	mu             *sync.Mutex
	out            io.Writer
}
```
主要是，我们需要一个缓冲区来保存预先格式化的数据。但我们也需要跟踪哪些组我们已经看到，但还没有输出。我们称这些组为 “未打开的”。我们还需要跟踪我们打开了多少个组，我们可以用一个简单的计数器来完成，因为打开的组的唯一效果是改变缩进级别。

`WithGroup` 实现与前一个实现非常相似：只需记住新组，它起初是未打开的。

```go
func (h *IndentHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := *h
	// Add an unopened group to h2 without modifying h.
	h2.unopenedGroups = make([]string, len(h.unopenedGroups)+1)
	copy(h2.unopenedGroups, h.unopenedGroups)
	h2.unopenedGroups[len(h2.unopenedGroups)-1] = name
	return &h2
}
```

`WithAttrs` 执行所有预格式化：

```go
func (h *IndentHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := *h
	// Force an append to copy the underlying array.
	pre := slices.Clip(h.preformatted)
	// Add all groups from WithGroup that haven't already been added.
	h2.preformatted = h2.appendUnopenedGroups(pre, h2.indentLevel)
	// Each of those groups increased the indent level by 1.
	h2.indentLevel += len(h2.unopenedGroups)
	// Now all groups have been opened.
	h2.unopenedGroups = nil
	// Pre-format the attributes.
	for _, a := range attrs {
		h2.preformatted = h2.appendAttr(h2.preformatted, a, h2.indentLevel)
	}
	return &h2
}

func (h *IndentHandler) appendUnopenedGroups(buf []byte, indentLevel int) []byte {
	for _, g := range h.unopenedGroups {
		buf = fmt.Appendf(buf, "%*s%s:\n", indentLevel*4, "", g)
		indentLevel++
	}
	return buf
}
```

它首先打开所有未打开的组。这将处理以下调用：

    logger.WithGroup("g").WithGroup("h").With("a", 1)

这里，`WithAttrs` 必须在 “a” 之前输出 “g” 和 “h”。由于由 `WithGroup` 建立的组对日志行的其余部分有效，因此 `WithAttrs` 为它打开的每个组递增缩进级别。

最后， `WithAttrs` 使用我们上面看到的相同的 `appendAttr` 方法格式化其参数属性。

`Handle` 方法的工作是将预先格式化的材料插入到正确的位置，即在内置属性之后和记录中的属性之前：

```go
func (h *IndentHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)
	if !r.Time.IsZero() {
		buf = h.appendAttr(buf, slog.Time(slog.TimeKey, r.Time), 0)
	}
	buf = h.appendAttr(buf, slog.Any(slog.LevelKey, r.Level), 0)
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = h.appendAttr(buf, slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", f.File, f.Line)), 0)
	}
	buf = h.appendAttr(buf, slog.String(slog.MessageKey, r.Message), 0)
	// Insert preformatted attributes just after built-in ones.
	buf = append(buf, h.preformatted...)
	if r.NumAttrs() > 0 {
		buf = h.appendUnopenedGroups(buf, h.indentLevel)
		r.Attrs(func(a slog.Attr) bool {
			buf = h.appendAttr(buf, a, h.indentLevel+len(h.unopenedGroups))
			return true
		})
	}
	buf = append(buf, "---\n"...)
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}
```

它还必须打开尚未打开的任何组。该逻辑包括如下日志行：

    logger.WithGroup("g").Info("msg", "a", 1)

其中 "g" 在调用 `Handle` 之前未打开，必须写入以产生正确的输出：

    level: INFO
    msg: "msg"
    g:
        a: 1

对 `r.NumAttrs() > 0` 的检查处理这种情况：

    logger.WithGroup("g").Info("msg")

这里没有记录属性，因此没有要打开的组。

>---
### Testing

[`Handler` contract](https://pkg.go.dev/log/slog#Handler) 指定了 *handlers* 的几个约束。要验证你的 *handler* 是否遵循这些规则并通常生成正确的输出，请使用 [testing/slogtest package](https://pkg.go.dev/testing/slogtest)。

该包的 `TestHandler` 函数接受 *handler* 的一个实例和一个函数，该函数返回格式化为映射切片的输出。下面是我们的示例 *handler* 的测试函数：

```go
func TestSlogtest(t *testing.T) {
	var buf bytes.Buffer
	err := slogtest.TestHandler(New(&buf, nil), func() []map[string]any {
		return parseLogEntries(t, buf.Bytes())
	})
	if err != nil {
		t.Error(err)
	}
}
```

调用 `TestHandler` 很简单。困难的部分是解析 *handler* 的输出。`TestHandler` 多次调用你的 *handler*，产生一系列日志条目。你的工作是将每个条目解析为 `map[string]any`。条目中的组应显示为嵌套映射。

如果 *handler* 输出标准格式，则可以使用现有的解析器。例如，如果你的 *handler* 每行输出一个 JSON 对象，那么你可以将输出分成几行，并在每一行上调用 `encoding/json.Unmarshal`。可以开箱即用地使用其他格式的解析器，这些解析器可以解组到映射中。我们的示例输出足够像 YAML，因此我们可以使用 `gopkg.in/yaml.v3` 包来解析它：

```go
func parseLogEntries(t *testing.T, data []byte) []map[string]any {
	entries := bytes.Split(data, []byte("---\n"))
	entries = entries[:len(entries)-1] // last one is empty
	var ms []map[string]any
	for _, e := range entries {
		var m map[string]any
		if err := yaml.Unmarshal([]byte(e), &m); err != nil {
			t.Fatal(err)
		}
		ms = append(ms, m)
	}
	return ms
}
```

如果你必须编写自己的解析器，那么它可能远非完美。`slogtest` 包只使用了少量的简单属性。（它测试 *handler* 的一致性，而不是解析。）解析器可以忽略键和值中的空白和换行符等边缘情况。在切换到 YAML 解析器之前，我们用 65 行代码编写了一个足够的自定义解析器。

---
## General considerations

### Copying records

大多数 *handlers* 不需要复制传递给 `Handle` 方法的 `slog.Record`。那些这样做的人在某些情况下必须特别小心。

*handlers* 可以通过普通的 Go 赋值、通道发送或函数调用来复制 `Record`，如果它不保留原始的话。但是如果它的操作产生多个副本，它应该调用 `Record.Clone` 来创建副本，这样它们就不会共享状态。这个 `Handle` 方法将记录传递给单个 *handler*，所以它不需要 `Clone`：

    type Handler1 struct {
        h slog.Handler
        // ...
    }

    func (h *Handler1) Handle(ctx context.Context, r slog.Record) error {
        return h.h.Handle(ctx, r)
    }

这个 `Handle` 方法可能会将记录传递给多个*handler*，所以它应该使用 `Clone`：

    type Handler2 struct {
        hs []slog.Handler
        // ...
    }

    func (h *Handler2) Handle(ctx context.Context, r slog.Record) error {
        for _, hh := range h.hs {
            if err := hh.Handle(ctx, r.Clone()); err != nil {
                return err
            }
        }
        return nil
    }

>---
### Concurrency safety

当一个 `Logger` 被多个 *goroutine* 共享时，一个 *handler* 必须正常工作。这意味着必须使用锁或其他机制来保护可变状态。实际上，这并不难实现，因为许多 *handlers* 没有任何可变状态。

- `Enabled` 方法通常只查询其参数和配置的级别。该级别通常在初始时设置一次，或者保存在已经是并发安全的 `LevelVar` 中。

- 由于上面讨论的原因，`WithAttrs` 和 `WithGroup` 方法不应该修改接收器。

- `Handle` 方法通常只使用其参数和存储的字段。

对输出方法（如 `io.Writer.Write`）的调用应该被同步，除非可以验证不需要锁定。正如我们在示例中看到的，存储指向互斥锁的指针使 *logger* 及其所有克隆能够彼此同步。要小心像 “Unix 写是原子的” 这样的肤浅的说法；情况要比这微妙得多。

一些 *handlers* 有合法的理由保持状态。例如，一个 *handler* 可能支持 `SetLevel` 方法来动态更改其配置级别。或者，它可能会输出连续调用 `Handle` 之间的时间，这需要一个可变字段来保存最后一次输出时间。同步对这些字段的所有访问，包括读和写。

内置 *handlers* 没有直接可变的状态。它们只使用互斥锁来对它们所包含的 `io.Writer` 进行排序调用。

>---
### Robustness

日志记录通常是调试技术的最后手段。当很难或不可能检查系统时（通常是生产服务器的情况），日志提供了了解其行为的最详细的方法。因此，你的 *handler* 应该对错误输入具有鲁棒性。

例如，当一个函数发现一个问题时，比如一个无效的参数，通常的建议是 *panic* 或返回一个 *error*。内置 *handler* 不遵循该建议。没有什么比无法调试导致日志记录失败的问题更令人沮丧的了；产生一些输出（无论多么不完美）总比什么都不产生好。这就是为什么像 `Logger.Info` 这样的方法会将键值对列表中的编程错误（如丢失的值或格式错误的键）转换为包含尽可能多信息的 `Attr`。

避免 *panics* 的一个地方是处理属性值。想要格式化值的 *handler* 可能会切换值的类型：

    switch attr.Value.Kind() {
    case KindString: ...
    case KindTime: ...
    // all other Kinds
    default: ...
    }

在默认情况下，当 *handler* 遇到它不知道的 `Kind` 时，会发生什么？内置 *handlers* 试图通过使用值的 `String` 方法的结果来蒙混过关，就像我们的示例 *handler* 一样。它们不会 *panic* 或返回 *error*。你自己的 *handler* 可能还希望通过生产监视或错误跟踪遥测系统报告问题。这个问题最有可能的解释是，更新版本的 `slog` 包添加了一个新的 `Kind` — Go 1 兼容性承诺下的向后兼容更改 — *handler* 没有更新。这当然是一个问题，但它不应该阻止 *readers* 看到日志输出的其余部分。

在一种情况下，从 `Handler.Handle` 返回错误是合适的。如果输出操作本身失败，最好的做法是通过返回错误来报告此失败。例如，内置的 `Handle` 方法的最后两行是

    _, err := h.w.Write(*state.buf)
    return err

尽管 `Logger` 的输出方法忽略了错误，但可以编写一个*handler*来处理它，也许可以返回到写入标准错误。

>---
### Speed

大多数程序不需要快速日志记录。在使你的 handler 快速之前，请收集数据 — 最好是生产数据，而不是基准比较 — 证明它需要快速。避免过早优化。

如果你需要一个快速的 *handler*，从预格式化开始。在对 `Logger.With` 的单个调用之后的结果 `logger` 的许多调用的情况下，它可以提供显著的速度提升。

如果日志输出是瓶颈，请考虑使 *handler* 异步。在 *handler* 中进行最少量的处理，然后通过通道发送记录和其他信息。另一个 *goroutine* 可以收集传入的日志条目，并在后台批量写入它们。你可能希望保留同步记录日志的选项，以便可以查看所有日志输出来调试崩溃。

分配往往是系统运行缓慢的主要原因。`slog` 包已经在努力最小化分配。如果你的 *handler* 自己进行分配，并且分析表明这是一个问题，那么看看你是否可以最小化它。

你可以做的一个简单更改是将对 `fmt.Sprintf` 或 `fmt.Appendf` 的调用替换为对缓冲区的直接 *appends*。例如，我们的 `IndentList` 将字符串属性附加到缓冲区，如下所示：

	buf = fmt.Appendf(buf, "%s: %q\n", a.Key, a.Value.String())

从 Go 1.21 开始，这会导致两个分配，每个分配用于传递给 `any` 参数的每个参数。我们可以通过直接使用 `append` 将其降为零：

	buf = append(buf, a.Key...)
	buf = append(buf, ": "...)
	buf = strconv.AppendQuote(buf, a.Value.String())
	buf = append(buf, '\n')

另一个有价值的改变是使用 `sync.Pool` 来管理大多数 *handlers* 需要的一块内存：保存格式化输出的 `[]byte` 缓冲区。

我们的示例 `Handle` 方法以这行开始：

	buf := make([]byte, 0, 1024)

如上所述，提供较大的初始容量可以避免随着切片的增长而重复复制和重新分配切片，从而将分配数量减少到 1。但是我们可以通过保持一个全局 *buffer pool* 来使它在稳定状态下降为零。开始时池将为空，并将分配新的缓冲区。但最终，假设并发日志调用的数量达到稳定的最大值，池中将有足够的缓冲区供所有正在进行的 `Handler` 调用共享。只要没有日志条目增长到超过缓冲区的容量，从垃圾收集器的角度来看，就不会有分配。

我们将把池隐藏在一对函数 `allocBuf` 和 `freeBuf` 后面。在 `Handle` 顶部获取缓冲区的单行变成了两行：

	bufp := allocBuf()
	defer freeBuf(bufp)

创建 `sync.Pool` 切片的一个微妙之处是变量名 `bufp`：你的池必须处理指向切片的指针，而不是切片本身。池 中的值必须始终是指针。如果它们不是，那么 `any` 参数和 `sync.Pool` 方法的返回值本身将导致分配，从而破坏池的目的。

有两种方法可以继续使用切片指针：我们可以在整个函数中将 `buf` 替换为 `*bufp`，或者我们可以解引用它并记住在释放之前重新赋值：

	bufp := allocBuf()
	buf := *bufp
	defer func() {
		*bufp = buf
		freeBuf(bufp)
	}()


以下是我们的池及其相关函数：

```go
var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 1024)
		return &b
	},
}

func allocBuf() *[]byte {
	return bufPool.Get().(*[]byte)
}

func freeBuf(b *[]byte) {
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}
```

池的 `New` 函数做的事情与原始代码相同：创建一个长度为 0 且容量充足的字节片。`allocBuf` 函数只是类型断言池的 `Get` 方法的结果。

`freeBuf` 方法在将缓冲区放回池中之前会截断缓冲区，因此 `allocBuf` 总是返回零长度的切片。它还实现了一个重要的优化：它不向池返回大的缓冲区。要理解这一点的重要性，请考虑一下如果有一个异常大的日志条目（比如格式化时有 1 兆字节）会发生什么。如果将兆字节大小的缓冲区放入池中，它可能会无限期地保留在那里，不断地被重用，但其大部分容量都被浪费了。额外的内存可能永远不会被 *handler* 再次使用，而且由于它在 *handler* 的池中，它可能永远不会被垃圾收集以在其他地方重用。我们可以通过从池中排除大型缓冲区来避免这种情况。

---