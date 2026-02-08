package gostd_testing

import (
	"sync/atomic"
	"testing"
	"testing/synctest"
)

/*
! synctest 包提供并发代码测试
- Test 创建一个测试气泡，在其中运行测试函数
- Wait 等待所有其他 goroutine 都永久阻塞
- 提供确定性的并发测试环境
*/

// TestSynctest 测试使用 WaitGroup 的情况
// go test -v -run=TestSynctest
func TestSynctest(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		counter := atomic.Int32{}
		// 启动多个 goroutine
		for range 100 {
			go func() {
				for range 1000000 {
					counter.Add(1)
				}
			}()
		}
		// 等待所有 goroutine 完成
		synctest.Wait()
		if counter.Load() != 100*1000000 {
			t.Fatalf("Expected counter to be %d, got %d", 100*1000000, counter.Load())
		}
	})
}
