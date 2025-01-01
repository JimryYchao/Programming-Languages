package gostd

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

/* Context
! context.Context 携带一个 deadline, cancellation signal, 和跨 API 边界的其他值
	Deadline 返回应取消代表此上下文完成工作时的时间，未设置时返回的 ok = false
	Done 返回一个通道，它在代表此上下文完成工作取消上下文时关闭。若 ctx 不能被取消，Done 可能返回 nil; Done 用于 select 语句
		WithCancel, WithDeadline, WithTimeout 安排 Done() 在取消或超时后关闭通道
	Err 在 Done() 关闭时返回一个非 nil err 解释原因：WithCancel, WithDeadline or WithTimeout
	Value 返回与此上下文关联的值；仅将 context values 用于传输进程和 API 边界的 request-scoped data，而不是用于向函数传递可选参数
		key 表示 Context 中的特定值，可以是支持相等的任何类型，应该将键定义为非导出类型
! Background 返回一个非 nil empty 的 Context, 它永远不会被取消，没有 context values，没有 deadline。它通常由主函数、初始化和测试使用，和作为传入请求的顶级上下文
! TODO 返回一个非 nil empty 的 Context。当不清楚要使用哪个 Context 或者它还不可用时（因为周围函数还没有扩展到接受 Context 参数），应使用 context.TODO
! WithValue 返回一个 parent 的副本，其中与 key 关联的值是 val。仅用于传输进程和 API 的请求作用域内的数据；
	提供的键必须是 comparable，并且不应该是 string 类型或任何其他内置类型，以避免使用 Context 的包之间的冲突。WithValue 的用户应该为键定义自己的类型。为了避免在
	分配给 interface{} 时进行分配，上下文键通常是具有具体的类型 struct{}。或者，导出的上下文键变量的静态类型应该是指针或接口。
! WithoutCancel 返回父级的副本，当父级被取消时，该副本不会被取消。返回的上下文不返回 Deadline 或 Err，其 Done 通道为 nil。在返回的上下文上调用 Cause 返回 nil。
*/

func TestCtx_WithValue(t *testing.T) {
	type favContextKey string

	f := func(ctx context.Context, k favContextKey) {
		if v := ctx.Value(k); v != nil {
			fmt.Println("found value:", v)
			return
		}
		fmt.Println("key not found:", k)
	}

	k := favContextKey("language")
	ctx := context.WithValue(context.Background(), k, "Go")

	f(ctx, k)
	f(ctx, favContextKey("color"))
	// found value: Go
	// key not found: color
}

/*
! Cause 返回一个非 nil 错误，解释为什么 c 被取消。第一个 c  或它的父母之一设置 cause。
	如果取消是通过调用 CancelCauseFunc(err) 发生的，那么 Cause 返回 err。否则 Cause(c) 返回与 c.Err() 相同的值。
	如果 c 未被取消，Cause 返回 nil。
! WithCancel 返回带有一个新的 Done() channel 的父级副本和一个 CancelFunc cancel
	调用 cancel 函数被调用时或父上下文的 Done 通道关闭时，返回的 context 的	Done 通道关闭。
	取消此 context 会释放与之关联的资源，因此代码应在此 context 运行的操作完成后立即调用 cancel
! WithCancelCause 与 WithCancel 类型但返回一个 CancelCauseFunc
	调用 cancel 返回一个非 nil 的错误（“cause”）并在 ctx 中记录这个错误，可以使用 Cause(ctx) 检索它
	使用 nil 调用 cancel 会将 cause 设置为 Canceled
! WithDeadline & WithDeadlineCause 返回父上下文的副本，并将终止时间调整为不晚于 d
	在 deadline 或调用 cancel 或父 ctx 的 Done 关闭后，子 ctx 随之关闭
! WithTimeout & WithTimeoutCause 相当于调用 WithDeadline(parent, time.Now().Add(timeout))
*/

func TestWithCancel(t *testing.T) {
	// gen 在单独的例程中生成整数，并将它们发送到返回的通道。
	// gen 的调用方需要在使用完生成的整数后取消上下文，以免泄露 gen 启动的内部例程。
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					checkErr(context.Cause(ctx))
					return // returning not to leak the goroutine
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel() // cancel when we are finished consuming integers

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}

func TestWithDeadline(t *testing.T) {
	var wg sync.WaitGroup
	shortDuration := 3 * time.Second
	d := time.Now().Add(shortDuration)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	// Even though ctx will be expired, it is good practice to call its
	// cancellation function in any case. Failure to do so may keep the
	// context and its parent alive longer than necessary.
	defer cancel()
	wg.Add(1)
	go func() {
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				log("do something")
			case <-ctx.Done():
				checkErr(ctx.Err())
				wg.Done()
				return
			}
		}
	}()
	wg.Wait()
}

/*
! AfterFunc 安排在 ctx 完成（取消或超时）后在自己的 goroutine 中调用 f。如果 ctx 已经 Done，AfterFunc 会立即在自己的 goroutine 中调用 f。
	调用返回的 stop 函数将停止 ctx 与 f 的关联。
*/

// ? 使用 AfterFunc 定义一个等待 sync.Cond 的函数，当上下文被取消时停止等待
func TestAfterFunc_Cond(t *testing.T) {
	waitOnCond := func(ctx context.Context, cond *sync.Cond, conditionMet func() bool, i int) error {
		stopf := context.AfterFunc(ctx, func() {
			// 我们需要获取 cond。L 在这里确保下面的 Broadcast 不会在调用 Wait 之前发生，否则会导致错过信号(和死锁)。
			cond.L.Lock()
			defer cond.L.Unlock()
			// 如果多个 goroutines 在 cond 的 goroutine 同时等待，我们需要确保我们只唤醒这个 goroutine。
			// 这意味着我们需要广播到所有的例程，这将唤醒它们。如果有 N 个对 waitOnCond 的并发调用，每个 goroutine
			// 将虚假地唤醒 O(N) 个尚未准备好的其他线程，因此这将导致总体 CPU 成本为 O(N²)。
			cond.Broadcast()
			logfln("goroutine %d quit", i)
		})
		defer stopf()

		// 由于唤醒使用的是 Broadcast 而不是 Signal，所以这个调用
		// Wait 可能会解除阻塞，因为其他例程的上下文已经完成，
		// 因此，为了确保 CTX 实际上已经完成，我们需要在循环中检查它。
		for !conditionMet() {
			logfln("goroutine %d is waiting", i)
			cond.Wait()
			if ctx.Err() != nil {
				return ctx.Err()
			}
		}

		return nil
	}

	cond := sync.NewCond(new(sync.Mutex))

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
			defer cancel()

			cond.L.Lock()
			defer cond.L.Unlock()
			err := waitOnCond(ctx, cond, func() bool { return false }, i)
			fmt.Println(err)
		}()
	}
	wg.Wait()

}

func TestAfterFunc_Connection(t *testing.T) {
	// 使用 AfterFunc 定义一个函数，该函数从 net.Conn 读取，当上下文被取消时停止读取。
	readFromConn := func(ctx context.Context, conn net.Conn, b []byte) (n int, err error) {
		stopc := make(chan struct{})
		stop := context.AfterFunc(ctx, func() {
			conn.SetReadDeadline(time.Now())
			close(stopc)
		})
		n, err = conn.Read(b)
		if !stop() {
			// AfterFunc 被启动。等待它完成，并重置 Conn 的截止日期。
			<-stopc
			conn.SetReadDeadline(time.Time{})
			return n, ctx.Err()
		}
		return n, err
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	conn, err := net.Dial(listener.Addr().Network(), listener.Addr().String())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	b := make([]byte, 1024)
	_, err = readFromConn(ctx, conn, b)
	fmt.Println(err)
}

// ? 使用 AfterFunc 定义一个函数, 它合并了多个上下文的取消信号
func TestAfterFunc_Merge(t *testing.T) {
	// mergeCancel 返回包含 ctx 值的上下文，当其中一个 ctx 被取消时，它被取消。
	mergeCancel := func(parent context.Context, ctxs ...context.Context) (context.Context, context.CancelFunc) {
		ctx, cancel := context.WithCancelCause(parent)
		stops := make([]func() bool, len(ctxs))
		for i, cancelCtx := range ctxs {
			stops[i] = context.AfterFunc(cancelCtx, func() {
				cancel(context.Cause(cancelCtx)) // cancelCtx 取消时，调用 f 取消 ctx
			})
		}

		return ctx, func() {
			for _, stop := range stops {
				go stop()
			}
			cancel(context.Canceled)
		}
	}

	// 构造两个 ctx
	ctx1, cancel1 := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel1()
	ctx2, cancel2 := context.WithCancelCause(context.Background())

	// 合并
	mergedCtx, mergedCancel := mergeCancel(context.Background(), ctx1, ctx2)
	defer mergedCancel()

	// time.Sleep(100 * time.Millisecond)   // ctx1 timeout
	cancel2(errors.New("ctx2 canceled")) // ctx2 canceled

	<-mergedCtx.Done()
	fmt.Println(context.Cause(mergedCtx))
}
