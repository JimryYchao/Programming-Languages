package gostd

import (
	"bytes"
	"fmt"
	"gostd/examples"
	"io"
	"net/http"
	"os"
	. "sync"
	"testing"
	"time"
)

/*
! sync.Locker 表示一个可以锁定 Lock 和解锁 Unlock 的对象。
! sync.Mutex 是一种互斥锁。Mutex 的零值是未锁定的 Mutex。Mutex 在第一次使用后不能复制。
	Lock 锁定 m。如果锁已经在使用中，调用的 goroutine 阻塞直到互斥锁可用
	TryLock 尝试锁定 m 并报告是否成功。
	Unlock 解锁 m。如果 m 没有在解锁时被未锁定，则是一个运行时错误。
! sync.RWMutex 是一种读写互斥锁。锁可以由任意数量的 reader 或单个 writer 持有。RWMutex 的零值是一个未锁定的互斥锁。RWMutex 在第一次使用后不能复制
	如果任何一个 goroutine 调用 Lock，而锁已经被一个或多个读取器持有，那么对 RLock 的并发调用将阻塞，直到 writer 获得（并释放）锁，以确保锁最终对 writer 可用。
	Lock 锁定 rw 进行写入。如果锁已锁定以进行读取或写入，则 Lock 将阻塞，直到锁可用为止。
	TryLock 尝试锁定 rw 进行写入，并报告是否成功
	Unlock 解锁 rw。如果 rw 未被 Lock 时调用，则这是一个运行时 panic
	RLock 锁定 rw 以进行读取。
	TryRLock 尝试锁定 rw 以进行读取，并报告是否成功。
	RUnlock 撤消单个 RLock 调用; 它不影响其他并发 reader。如果 rw 没有被 RLock 时调用，则这是一个运行时 panic
	RLocker 返回一个 Locker 接口值，该接口对象通过调用 rw.RLock 和 rw.RUnlock 实现 Lock 和 Unlock
*/

func TestLocker(t *testing.T) {
	hammerMutex := func(m *Mutex, loops int, cdone chan bool, n *int) {
		for i := 0; i < loops; i++ {
			go func(i int) {
				select {
				case <-cdone:
					return
				default:
					if i%2 == 0 {
						if m.TryLock() {
							logfln("loop %d: tryLock", i)
							*n++
							m.Unlock()
							logfln("loop %d: unLock try", i)
						}
					}
				}
			}(i)
			m.Lock()
			logfln("loop %d: lock", i)
			m.Unlock()
			logfln("loop %d: unlock", i)
		}
		cdone <- true
	}

	t.Run("Mutex", func(t *testing.T) {
		m := new(Mutex)
		m.Lock()
		logfln("call TryLock = %t while Lock", m.TryLock())
		m.Unlock()
		logfln("call TryLock = %t while UnLock", m.TryLock())
		m.Unlock()
		n := 0
		c := make(chan bool)
		go hammerMutex(m, 20, c, &n)
		for range 1 + n {
			<-c
		}
	})

	t.Run("RWMutex", func(t *testing.T) {
		rwm := new(RWMutex)
		n := new(int)
		done := make(chan bool)
		reader := func(l *RWMutex) {
			logfln("read %d", *n)
			l.RUnlock()
		}

		go func() {
			for {
				select {
				case <-done:
					return
				case <-time.After(500 * 60 * 2000):
					if rwm.TryRLock() {
						go reader(rwm)
					}
				}
			}
		}()

		go func() {
			for range 100 {
				rwm.Lock()
				*n++
				rwm.Unlock()
				<-time.After(1)
			}
			done <- true
		}()
		<-done
		log(*n) // 100
	})

}

/*
! OnceFunc 返回一个只调用 f 一次的函数。返回的函数可以并发调用。
! OnceValue, OnceValues 返回一个只调用 f 一次的函数，并返回 f 返回的值。返回的函数可以并发调用。每次调用时都出现相同的返回值
! sync.Once 是一个只执行一个操作的对象；Do 之后不能复制 Once
*/

func TestOnce(t *testing.T) {
	t.Run("OnceFunc", func(t *testing.T) {

		f := func() {
			log("call f once")
		}
		fvalue := func() int {
			log("call fvalue once")
			return 0
		}
		fvalues := func() (int, error) {
			log("call fvalues once")
			return 0, nil
		}

		oncef := OnceFunc(f)
		oncefvalue := OnceValue(fvalue)
		oncefvalues := OnceValues(fvalues)

		for range 10 {
			oncef()
			oncefvalue()
			oncefvalues()
		}
	})

	t.Run("sync.Once", func(t *testing.T) {
		var once Once
		c := once
		onceBody := func() {
			fmt.Println("Only once")
		}
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				once.Do(onceBody)
				done <- true
			}()
		}
		for i := 0; i < 10; i++ {
			<-done
		}

		c.Do(onceBody) // copy before once.Do; ok
	})
}

/*
! sync.Cond 实现了一个条件变量，一个等待或宣布事件发生的 goroutine 的集合点。NewCond 构造一个 Cond
	每个 Cond 都有一个关联的 Locker（通常是 *Mutex 或 *RWMutex），在更改条件和调用 Wait 方法时必须持有；
	首次使用后不得复制 Cond
	在 Go 内存模型的术语中，Cond 安排对 Broadcast 或 Signal 的调用 “优先于” 它解除阻塞的任何 Wait 调用。
	对于许多简单的用例，用户使用通道比使用 Cond 更好（Broadcast 对应于关闭通道，Signal 对应于在通道上发送）。
Cond.methods
	Broadcast 唤醒所有等待 cond 的 goroutine
	Singal 唤醒一个等待 cond 的 goroutine，如果有的话。允许但不要求主叫方在调用过程中持有 C.L。
		Signal() 不影响 goroutine 的调度优先级; 如果其他 goroutine 试图锁定 c.L，它们可能会在一个 “等待” 的 goroutine 之前被唤醒。
	Wait 自动解锁 c.L 并挂起调用 goroutine 的执行。在稍后恢复执行之后，Wait 在返回之前锁定 c.L。
		与其他系统不同，Wait 不能返回，除非被 Broadcast 或 Signal 唤醒。
		由于 c.L 在 Wait 等待时未被锁定，因此调用方通常不能假定 Wait 返回时条件为 true。相反，调用者应该在循环中等待：
*/

func TestCondSignal(t *testing.T) {
	m := Mutex{}
	c := NewCond(&m)
	n := 3
	wait := func() {
		log("cond waiting")
		c.Wait()
	}
	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			m.Lock()
			running <- true
			wait()
			awake <- true
			m.Unlock()
			log("cond awake")
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}

	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}
		m.Lock()
		c.Signal() // 唤醒一个 goroutine
		m.Unlock()
		<-awake // Will deadlock if no goroutine wakes up
		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}
	c.Signal()
}

func TestCondBroadcast(t *testing.T) {
	m := Mutex{}
	c := NewCond(&m)
	n := 3
	running := make(chan bool, n+1)
	awake := make(chan int, n)
	for i := 0; i < n; i++ {
		go func(i int) {
			m.Lock()
			running <- true
			logfln("wait %d", i)
			c.Wait()
			awake <- i
			m.Unlock()
		}(i)
		<-running
	}
	running <- true
	log("Broadcast all")
	c.Broadcast()
	<-running
	for range n {
		a := <-awake
		logfln("awake %d", a)
	}
}

/*
! sync.WaitGroup 等待一个 goroutine 集合的完成。
	主 goroutine 调用 `Add`` 来设置要等待的 goroutine 的数量。然后每个 goroutine 运行并在完成时调用 `Done`。
	同时，`Wait` 可以用来阻塞主 goroutine，直到所有的 goroutine 都完成。首次使用后不得复制 WaitGroup
*/

func TestWaitGroup(t *testing.T) {
	var wg WaitGroup
	var urls = []string{
		"https://pkg.go.dev/sync",
		"http://www.bing.com",
		"http://www.example.com/",
	}
	for _, url := range urls {
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go func(url string) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Fetch the URL.
			resp, err := http.Get(url)
			checkErr(err)
			if err == nil {
				log(resp.Request.URL)
			}
		}(url)
	}
	// Wait for all HTTP fetches to complete.
	wg.Wait()
}

/*
! sync.Map 类似于 Go map[any]any，但对于多个 goroutine 并发使用是安全的，无需额外的锁定或协调；Map 类型是专用的。应使用普通的 map，并带有单独的锁定或协调
	CompareAndDelete 在 may[key] == old 时删除 key 项；在 map 中没有 key 的当前值，CompareAndDelete 返回 false
	CompareAndSwap 在 map[key] == old 时，交换 key 的 old, new 值
	Delete 删除 key 的值。
	Load 返回 map 中存储的 key 的值; ok 描述是否找到
	LoadAndDelete 删除键的值，（如果有）则返回以前的值。loaded 报告 key 是否存在。
	LoadOrStore 返回键的值（如果有）。否则，它存储并返回给定的 value。如果值被加载，则 loaded 为 true，如果 value 被存储，则为 false。
	Range 为 map 中的每个键和值顺序调用 f。如果 f 返回 false，则 Range 停止迭代。
	Store 存储设置 key 的 value。
	Swap 交换键的值，并返回先前的值（如果有），没有键时则是 Store。loaded 的结果报告 key 是否存在。
*/

func TestMap(t *testing.T) {
	var m Map
	var wg WaitGroup
	wg.Add(1)
	Range := func() {
		m.Range(func(k, v any) bool {
			log(k, v)
			return true
		})
	}

	go func() {
		defer wg.Done()
		for k, v := range map[int]string{
			1: "a", 2: "b", 3: "c", 5: "e",
		} {
			m.Store(k, v)
		}
	}()

	if v, ok := m.Load(2); ok {
		log(v)
	} else {
		m.Store(2, "B")
	}
	wg.Wait()

	if m.CompareAndDelete(1, "a") {
		_, ok := m.Load(1)
		logfln("load k1 %t", ok)
	}

	if v, ok := m.LoadOrStore(4, "d"); ok { // store
		log(v)
	} else {
		log("load k4 failed")
		Range()
	}

	if vold, ok := m.Swap(1, "A"); ok {
		vnew, _ := m.Load(1)
		logfln("set m[1]:%v to %v", vold, vnew)
	} else {
		log("swap k1 failed")
		Range()
	}
	if v, ok := m.LoadAndDelete(4); ok {
		logfln("4,%v is deleted", v)
	}

	if m.CompareAndSwap(1, "A", "a") {
		if v, ok := m.Load(1); ok {
			logfln("load k1 %t", v)
		}
	}
}

/*
! sync.Pool 是一组可以单独保存和检索的临时对象。可以安全地被多个 goroutine 同时使用
	存储在 Pool 中的任何项可能会随时自动删除，且不会通知。如果发生这种情况时 pool 持有唯一的引用，则可能会释放该项。
	Pool 的目的是缓存已分配但未使用的项以供以后重用，从而减轻垃圾收集器的压力。
	Pool 的一个适当用法是管理一组临时项，这些临时项在包的并发独立客户端之间静默共享，并可能被包的并发独立客户端重用。Pool 提供了一种在多个客户端之间分摊分配开销的方法。
	一个使用 Pool 的很好的例子是 fmt 包
Pool-functions
	Pool.New 可选的指定一个函数来生成一个值，否则 Get 返回 nil；不能在调用 Get 的同时修改 New 的值
	Get 从池中选择任意项，将其从池中删除，并将其返回给调用方。如果 Get 返回 nil，而 p.New 不为 nil，那么 Get 返回 p.New 调用的结果。
	Put 将 x 添加到池中。
*/

func TestPool(t *testing.T) {
	var bufPool = examples.NewSyncPool(func() any {
		log("create a pool")
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	})
	timeNow := func() time.Time {
		return time.Now()
	}

	log := func(w io.Writer, key, val string) {
		b := bufPool.Get().(*bytes.Buffer)
		b.Reset()
		// Replace this with time.Now() in a real logger.
		b.WriteString(timeNow().UTC().Format(time.RFC3339))
		b.WriteByte(' ')
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteString(val + "\n")
		w.Write(b.Bytes())
		bufPool.Put(b)
	}

	n := 100
	c := make(chan bool, n)
	for i := 1; i < n+1; i++ {
		i := i
		if i == n/2 {
			bufPool.Clear()
		}
		go func() {
			log(os.Stdout, fmt.Sprint(i), fmt.Sprintf("%#U", i+64))
			c <- true
		}()
	}

	for range n {
		<-c
	}
}
