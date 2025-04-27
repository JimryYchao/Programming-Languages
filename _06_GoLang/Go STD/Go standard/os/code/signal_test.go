package gostd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"testing"
	"time"
)

/*
! Ignore, Ignored 忽略提供的信号；Ignore 也将撤消之前对所提供信号的 `Notify` 的任何调用的影响
! Notify 通知 package signal 将给定信号类别传递给通道 c；若未提供则所有的输入信号都将中继到 c
	package signal 不会阻塞发送到 c，调用方必须确保 c 有足够的缓冲空间来跟上预期的信号速率；
	对于仅用于通知一个信号值的通道，大小为 1 的缓冲区就足够了；允许使用同一通道多次调用 Notify
	每次调用都会扩展发送到该通道的信号集。从集合中删除信号的唯一方法是调用 Stop
	允许使用不同的通道和相同的信号多次调用 Notify：每个通道独立接收传入信号的副本。
! NotifyContext 返回一个与 signal 关联的 context，当调用 stop, 或 parent.Done, 或接收到信号时，该 context 取消
	stop 函数取消注册信号行为；调用 NotifyContext(parent, os.Interrupt) 将更改中断行为为取消返回的上下文。
	stop 函数释放与上下文关联的资源；因此，一旦在此上下文中运行的操作完成或不再需要时，主动调用 stop
! Reset 将撤消之前对所提供信号 signal 的 Notify 调用的任何效果
! Stop 使 signal 停止将输入信号到 c。并撤销已设置的 signals；当 Stop 返回时，可以保证 c 不会再接收到信号。
*/

// ? go test -v -run=^TestNotifyContext$
func TestNotifyContext(t *testing.T) {
	errCh := make(chan error)

	log("please press ^c to send os.Interrupt to sigCtx")
	hello := func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				errCh <- context.Cause(ctx)
				return
			default:
				fmt.Println("Hello World")
				time.Sleep(250 * time.Millisecond)
			}
		}
	}
	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second) // 执行 15s
	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	defer signal.Reset()

	go hello(sigCtx)
	log(<-errCh)
}
