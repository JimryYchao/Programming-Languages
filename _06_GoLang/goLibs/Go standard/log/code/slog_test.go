package gostd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"
)

func logAll(l *slog.Logger, msg string, args ...any) {
	l.Debug(msg, args...)
	l.Info(msg, args...)
	l.Warn(msg, args...)
	l.Error(msg, args...)
	fmt.Println()
}
func logAllContext(l *slog.Logger, ctx context.Context, msg string, args ...any) {
	l.DebugContext(ctx, msg, args...)
	l.InfoContext(ctx, msg, args...)
	l.WarnContext(ctx, msg, args...)
	l.ErrorContext(ctx, msg, args...)
	fmt.Println()
}

var content = "this is a message"
var contentWith = content + " with "

// Logger, Handler, Leveler, Level, LevelVar, Attr, Value, Kind

/*
! slog.Logger 记录有关对其 `Log`,`Debug`,`Info`,`Warn`,`Error` 方法每次调用的结构化信息，相应的包级方法调用 slog 默认的 logger。对它的日志方法的每次调用都会创建一个 `slog.Record` 对象并传递给与之关联 `slog.Handler`。
	Enabled 报告 logger 在给定上下文 ctx 和级别 level 上发出日志
	Handler 返回与 logger 关联的日志处理程序
	Log & LogAttr 发出一条包含当前时间、给定级别和消息的日志记录
	Debug & DebugContext 在 LevelDebug 上发出日志；DebugContext 使用给定的上下文 ctx 记录日志
	Info & InfoContext 在 LevelInfo 上发出日志
	Warn & WarnContext 在 LevelWarn 上发出日志
	Error & ErrorContext 在 LevelContext 上发出日志
	With 返回一个 logger，它将在每个日志输出操作中包含给定的属性
	WithGroup 返回一个 logger，它的输出方法的所有属性的键都由给定的 `group` 限定
! New 构造一个新的关联给定 handler 的 slog.Logger。
! Default 返回 slog 包级使用的默认 *slog.Logger。slog 的顶级日志函数使用其进行日志事件的输出。未调用 SetDefault 之前使用 log 包的默认 logger。
! SetDefault 重新关联一个新的 logger 作为 slog 包的 defaultLogger。log 包的输出函数将使用 logger 的 handler 进行日志处理和输出。
! SetLogLoggerLevel 设置 slog.defaultLogger 和桥接到 log.defaultLogger 的日志级别。默认级别为 slog.LevelInfo。
! slog.Level 表示日志事件的重要程度或严重级别，默认为 Info:0
*/

func TestSlogTopFunctions(t *testing.T) {
	logAll(slog.Default(), "before set level")
	old := slog.SetLogLoggerLevel(slog.LevelWarn)
	t.Cleanup(func() {
		slog.SetLogLoggerLevel(old) // 恢复原始 level 设置
	})
	logAll(slog.Default(), "SetLogLoggerLevel to WARN")

	slog.SetLogLoggerLevel(slog.LevelDebug)
	logAll(slog.Default(), "SetLogLoggerLevel to DEBUG")
	logAllContext(slog.Default(), nil, "nil context")
	logAll(slog.Default(), "log with attrs", "a1", "v1", "a2", "v2")
}

func TestNewLogger(t *testing.T) {
	lgr := slog.New(&wrappingHandler{slog.Default().Handler(), slog.LevelInfo})
	lgr.Debug("not printed")
	lgr.Info("new a logger with a wrappingHandler")
	lgr.Error(contentWith, "err", errors.New("logging error"))

	if !lgr.Enabled(context.Background(), slog.LevelDebug) {
		lgr.Handler().(*wrappingHandler).Set(slog.LevelDebug)
		lgr.Debug("set handler-level to DEBUG")
	}
	// 2024/05/26 15:46:02 INFO new a logger with a wrappingHandler
	// 2024/05/26 15:46:02 ERROR this is a message with  err="logging error"
	// 2024/05/26 15:46:02 DEBUG set handler-level to DEBUG
}

func TestWithLogger(t *testing.T) {
	lgr := slog.With("author", "Ychao")
	lgr.Info("create a logger by slog.With")
	args := []any{"a", 1, slog.Int64("b", 2), slog.Group("group", "c", 3, "d", 4)}
	// 2024/05/26 03:53:40 INFO this is a message with  author=Ychao a=1 b=2 group.c=3 group.d=4
	lgr.Info(contentWith, args...)

	lgr = lgr.WithGroup("class")
	lgr.Info("create a logger by lgr.WithGroup")
	lgr.Info(contentWith, args...)
	// 2024/05/26 03:53:40 INFO this is a message with  author=Ychao class.a=1 class.b=2 class.group.c=3 class.group.d=4
}

func TestSetDefault(t *testing.T) {
	log.Print("log.Print before SetDefault")
	slog.Debug("not printed")
	slog.Info("printed")
	lgr := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(lgr)
	slog.Debug("call SetDefault")
	log.Print("log.Print after SetDefault")
	// 2024/05/26 16:03:02 log.Print before SetDefault
	// 2024/05/26 16:03:02 INFO printed
	// {"time":"2024-05-26T16:03:02.1233064+08:00","level":"DEBUG","msg":"call SetDefault"}
	// {"time":"2024-05-26T16:03:02.1233064+08:00","level":"INFO","msg":"log.Print after SetDefault"}
}

func TestSLogToFile(t *testing.T) {
	logfile, _ := os.OpenFile("files/logfile.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	defer logfile.Close()
	mLogger := slog.NewLogLogger(slog.Default().Handler(), slog.LevelDebug) // 关联一个 handler 和 level
	mLogger.Print("create a new log.Logger with a slog.defaultHandler")
	mLogger.SetOutput(io.MultiWriter(os.Stdout, logfile)) // 关联 logger 的输出到 logfile 和 stdout
	mLogger.Print("Hello logfile")

	log.SetOutput(mLogger.Writer())

	slog.Debug("test log a debug")
	slog.SetLogLoggerLevel(slog.LevelDebug)
	// 未调用 SetDefault 之前 slog 使用 log.defaultLogger
	slog.Debug("test log a debug again") // 同时发出到 logfile 和 stdout
}

/*
! slog.Attr 使用 (string, Value) 表示一组 key-value 对，键为 string，值为 `slog.Value`。Group 函数打包一组参数值为一个 Attr。Attr 可以作为日志输出函数的属性参数。
! slog.Value 用来表示任何 Go 值。与 any 类型不同的是，它可以表示大多数没有分配的小值。零值对应于 nil。slog.Kind 表示 Value 的类别。
! AnyValue 返回 v 对应的 Value。
	- 预声明的 string, bool, 非复数值分别返回 KindString, KindBool, KindUInt64, KindInt64, KindFloat64
	- time.Time, time.Duration 返回 KindTime, KindDuration
	- 对于 nil 和其他类型，包括命名类型，返回 KindAny
! GroupValue 用于组合一组 Attrs 为一个 Value
*/

func TestLogAttrs(t *testing.T) {
	slog.Log(context.Background(), slog.LevelInfo+1, "hello", "a1", "v1", "a2", "v2", slog.Int("a", 1), slog.String("b", "two"), slog.Group("group", "c", 3, "d", "four"))
	// 2024/05/26 02:43:30 INFO+1 hello a=1 b=two group.c=3 group.d=four
	slog.LogAttrs(context.Background(), slog.LevelDebug+3, "world", slog.Any("any", time.Now()), slog.String("Hello", "World"))
	// 2024/05/26 02:48:27 DEBUG+3 world any=2024-05-26T02:48:27.490+08:00 Hello=World
}

// ! slog.LogValuer 是可以将自己转换为用于 logging 值的任何 Go 值。这种机制可以用来推迟昂贵的操作。
func TestLogValuer(t *testing.T) {
	n := Name{"Jimry", "Ychao"}
	lgr := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	lgr.Info("author", "name", n)

	// json output:
	/*
		{
			"time":"2024-05-26T21:45:47.8978386+08:00",
			"level":"INFO",
			"msg":"author",
			"name":{
				"first":"Jimry",
				"last":"Ychao"
			}
		}
	*/
}

type Name struct {
	First, Last string
}

func (n Name) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("first", n.First),
		slog.String("last", n.Last))
}

/*
! slog.Handler 处理 logger 生成的 records。用户不应该直接调用 handler 的方法，应调用其关联的 logger
	Enabled 报告 handler 是否处理给定 ctx 和 level 的记录，它在任何的记录参数处理之前调用，低级别将被忽略
	Handle 在 ctx 下处理 record，仅当在 Enabled 为 true 时被调用
	WithAttrs 返回一个 handler，它包含其接收方的属性和参数组成
	WithGroup 返回一个 handler，它将给定的 group 附加到接收方的现有组中，后续的所有键都由 group 限定。
! JSONHandler 将 record 作为行形式的 JSON 对象写入 io.Writer。NewJSONHandler 创建 JSONHandler。
! TextHandler 将 record 作为空格分隔的 key=value 序列写入 io.Writer。NewTextHandler 创建 TextHandler。
*/

func TestSlogHandlers(t *testing.T) {
	t.Run("JSONHandler", func(t *testing.T) {
		jsonSlgr := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		jsonSlgr.Info("create a slogger with a JSONHandler")
		jsonSlgr.Info(contentWith, "a", 1, slog.Int64("b", 2),
			slog.Group("map", slog.String("k1", "v1"), "k2", "v2"))

		jSlgrWith := jsonSlgr.With("author", "Ychao")
		jSlgrWith.Info(content)

		jSlgrWithGroup := jSlgrWith.WithGroup("G")
		jSlgrWithGroup.Info(content)
		jSlgrWithGroup.Info(contentWith, "a", 1, slog.Int64("b", 2),
			slog.Group("map", slog.String("k1", "v1"), "k2", "v2"))
		// output:
		// {"time":"2024-05-26T16:19:41.1795976+08:00","level":"INFO","msg":"create a slogger with a JSONHandler"}
		// {"time":"2024-05-26T16:19:41.1853126+08:00","level":"INFO","msg":"this is a message with ","a":1,"b":2,"map":{"k1":"v1","k2":"v2"}}
		// {"time":"2024-05-26T16:19:41.1853126+08:00","level":"INFO","msg":"this is a message","author":"Ychao"}
		// {"time":"2024-05-26T16:19:41.1853126+08:00","level":"INFO","msg":"this is a message","author":"Ychao"}
		// {"time":"2024-05-26T16:19:41.1853126+08:00","level":"INFO","msg":"this is a message with ","author":"Ychao","G":{"a":1,"b":2,"map":{"k1":"v1","k2":"v2"}}}
	})

	t.Run("TextHandler", func(t *testing.T) {
		txtSlgr := slog.New(slog.NewTextHandler(os.Stdout, nil))
		txtSlgr.Info("create a slogger with a JSONHandler")
		txtSlgr.Info(contentWith, "a", 1, slog.Int64("b", 2),
			slog.Group("map", slog.String("k1", "v1"), "k2", "v2"))

		tSlgrWith := txtSlgr.With("author", "Ychao")
		tSlgrWith.Info(content)

		tSlgrWithGroup := tSlgrWith.WithGroup("G")
		tSlgrWithGroup.Info(content)
		tSlgrWithGroup.Info(contentWith, "a", 1, slog.Int64("b", 2),
			slog.Group("map", slog.String("k1", "v1"), "k2", "v2"))
		// output:
		// time=2024-05-26T16:21:52.631+08:00 level=INFO msg="create a slogger with a JSONHandler"
		// time=2024-05-26T16:21:52.631+08:00 level=INFO msg="this is a message with " a=1 b=2 map.k1=v1 map.k2=v2
		// time=2024-05-26T16:21:52.631+08:00 level=INFO msg="this is a message" author=Ychao
		// time=2024-05-26T16:21:52.631+08:00 level=INFO msg="this is a message" author=Ychao
		// time=2024-05-26T16:21:52.631+08:00 level=INFO msg="this is a message with " author=Ychao G.a=1 G.b=2 G.map.k1=v1 G.map.k2=v2
	})
}

/*
! slog.HandlerOptions 用于构造 TextHandler 或 JSONHandler 时的选项。
	`AddSource bool` 表示在是否在输出中添加源信息
	`Level Leveler` 表示 handler 处理 record 的最小等级，若要动态更改，可使用 slog.LevelVar
	`ReplaceAttr fn` 在重写 record 的每个非 group 属性之前被调用，内置键 time,level,source,msg 默认被传递。
! slog.Leveler 提供一个 Level 值。Level 和 LevelVar 实现了该接口，可以通过 `LevelVar.Set` 方法动态更改 [Handler] 的级别。
*/

func TestLevelVar(t *testing.T) {
	le := &slog.LevelVar{} // 0 : Info
	lgr := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: le}))
	logAll(lgr, "print test")
	le.Set(slog.LevelError)
	logAll(lgr, content)
	le.Set(5)
	logAll(lgr, "set to 5", "a", 1)
	logAll(lgr, contentWith, slog.Group("G", slog.String("k1", "v1"), slog.String("k2", "v2")))
}

// ! NewLogLogger 构造一个关联给定 handler 和 level 的 log.Logger
func TestLogLogger(t *testing.T) {
	lgr := slog.NewLogLogger(slog.Default().Handler(), slog.LevelDebug)

	lgr.Print("not printed")
	slog.SetLogLoggerLevel(slog.LevelDebug)
	lgr.Print("printed")
	// 2024/05/26 17:51:49 DEBUG printed

	lgr = slog.NewLogLogger(slog.Default().Handler(), slog.LevelInfo)
	lgr.SetPrefix("MESS: ")
	lgr.Print("printed")
	// 2024/05/26 17:51:49 INFO MESS: printed
}

/*
! slog.Record 保存有关日志事件的信息。复制 record 将在副本之间共享状态。
	Time 		调用 output methods (Log, Info ...) 的调用时间
	Message 	log mess
	Level 		日志事件级别
	PC			在构造 record 时的程序计数器，0 表示不可用；仅作为 [runtime.CallersFrames] 参数。
! slog.Record Methods
	Add & AddAttrs 添加属性到 record 的属性列表
	Attrs 在 record 的每个属性上调用 f，f 返回 false 时迭代停止
	Clone 克隆一个与原始 record 没有共享状态的副本。
	NumAttrs 返回 record 中的属性数目
! NewRecord 从给定的参数创建一个 Record。使用 Record.AddAttrs 向记录添加属性。NewRecord 用于记录期望支持 Handler 作为后端的 API。
*/
