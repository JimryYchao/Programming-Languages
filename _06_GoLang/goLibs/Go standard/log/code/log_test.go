package gostd

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

/* standard Logger
! New 创建一个新的 `log.Logger` 并关联一个 `io.Writer out`。`prefix` 出现在每个生成日志行的开头。`flag` 定义日志记录的属性。
	Ldata : 2009/01/23
	Ltime : 01:23:23
	Lmicroseconds : 01:23:23.123123
	Llongfile : /a/b/c/d.go:23
	Lshortfile : d.go:23
	LUTC : 使用 UTC 时区
	Lmsgprefix : 将前缀 prefix 移动到 mess 前
	LstdFlags = Ldate | Ltime
! log.Logger 表示一个活动的日志记录对象，它生成到 io.Writer 的输出行。每个日志记录操作都调用 `Writer` 的 Write() 方法。Logger 可以同时在多个 goroutine 中使用；它保证了对 Writer 的序列化访问。
! Default 返回一个包级输出函数使用的标准日志 Logger。
*/

func TestStdLogger(t *testing.T) {
	t.Run("std logger", func(t *testing.T) {
		t.Cleanup(func() {
			// reset std logger
			log.SetFlags(log.LstdFlags)
			log.SetOutput(os.Stdout)
			log.SetPrefix("")
		})
		log.Print("hello", "World")
		log.Print("hello", 1, 2, 3, "World")
		log.Println("hello", 1, 2, 3, "World")

		stdlogger := log.Default()
		fmt.Printf("LOG Prefix is `%s`\n", stdlogger.Prefix())

		stdlogger.SetPrefix("LOG_TEST: ") // 设置日志输出的前缀
		log.Print(time.Now())             // == stdlogger.Print

		log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lmsgprefix) // == stdlogger.SetFlogs
		log.Print("Hello, 世界")

		log.Output(-100, "print")
	})

	t.Run("parallel logger", func(t *testing.T) {
		compCh := make(chan bool)
		var ngorte atomic.Int32
		logInGoroutine := func(logger *log.Logger, n int) {
			ngorte.Add(1)
			for i := range n {
				logger.Printf("line %d", i)
				time.Sleep(200 * time.Millisecond)
			}
			compCh <- true
		}
		go logInGoroutine(log.New(os.Stdout, "[DEBUG]", log.Lmsgprefix), 10)
		go logInGoroutine(log.New(os.Stdout, "[WARNING]", log.LstdFlags|log.Lmsgprefix), 15)
		go logInGoroutine(log.New(os.Stdout, "", log.Ldate|log.LUTC), 5)
		go logInGoroutine(log.New(os.Stdout, "[FATAL]", log.LstdFlags|log.Lmsgprefix|log.Lmicroseconds), 3)
		go logInGoroutine(log.New(os.Stdout, "[INFO]", log.Lmsgprefix), 8)

		var n int32 = 0
		for <-compCh {
			n++
			if ngorte.Load() <= n {
				break
			}
		}
	})
}

func TestLogToFile(t *testing.T) {
	lf, err := os.OpenFile("files/logfile", os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		checkErr(err)
		return
	}
	lf.Seek(0, io.SeekEnd)
	defer lf.Close()
	// 设置 logger 的 log 格式
	setlogfmt := func(lgr *log.Logger) {
		lgr.SetFlags(log.LstdFlags | log.Lmsgprefix)
		lgr.SetPrefix("[MESS] ")
	}

	// 关联 lfw 和 default logger
	lfwr := log.New(lf, "", 0)
	setlogfmt(lfwr)
	setlogfmt(log.Default())
	mulwr := io.MultiWriter(lfwr.Writer(), log.Default().Writer())
	log.SetOutput(mulwr)

	// 测试整合后的 log.Default
	log.Print("Hello World")
	// log.Panic("PANICKING")
}
