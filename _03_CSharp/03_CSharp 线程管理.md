## CSharp 线程管理

---

### 1. 线程管理

所有的线程管理由 `Thread` 完成。托管线程分为后台线程和前台线程，后台线程不会终止进程停止。`thread.IsBackground` 设置线程状态。托管线程池的线程是后台线程，创建或启动的 `Thread` 对象默认为前台线程。
- `thread.Start` 启动线程。
- `thread.Join` 等待线程终止。
- `thread.Interrupt` 中止线程。
- `thread.Priority` 线程优先级。
- `Thread.Sleep` 线程休眠。

`ThreadPool` 提供辅助线程池，`Task` 在线程池上执行异步任务。线程池线程是后台线程。`ThreadPool.QueueUserWorkItem` 压入一个方法，在线程池线程可用时调用。

```csharp
ThreadPool.QueueUserWorkItem((_) => Console.WriteLine("In thread {0}", Thread.CurrentThread.ManagedThreadId));
```

利用 `CancellationTokenSource` 和 `CancellationToken` 管理取消线程或任务。

```csharp
class Sample {
    static void Main(string[] args) {
        using CancellationTokenSource CTS = new CancellationTokenSource();
        CTS.Token.Register(() => {
            Console.WriteLine("Request cancelled");
        });
        ThreadPool.QueueUserWorkItem(new WaitCallback((state) => {
            if (state is null)
                return;
            var token = (CancellationToken)state;
            for (int i = 0; i < 100; i++) {
                if (token.IsCancellationRequested) {
                    Console.WriteLine("In iteration {0}, Break by cancellation...", i + 1);
                    break;
                }
                Thread.Sleep(new Random().Next(10, 50));
            }
        }), CTS.Token);
        Thread.Sleep(1000);
        CTS.Cancel();
        Thread.Sleep(1000);
    }
}
```

`WaitHandle` 封装本机等待句柄，使用信号机制进行线程交互。利用 `WaitHandle.WaitAny`（阻塞）监听 `WaitHandle` 事件完成或 `CancellationToken` 取消请求。

```csharp
ManualResetEvent mre = new ManualResetEvent(false);
using CancellationTokenSource CTS = new CancellationTokenSource();
// do something
var signal = WaitHandle.WaitAny(waitHandles: [mre, CTS.Token.WaitHandle], timeout: new(0, 0, 20));
// 0 >> mre.Set() 
// 1 >> CTS.Cancel()
// ...
// 258 >> time out
```

`ThreadPool.RegisterWaitForSingleObject` 注册一个等待句柄，句柄完成或超时时触发回调。

```csharp
class Sample {
    class TaskInfo {
        public RegisteredWaitHandle handle;
    }
    static void Main(string[] args) {
        using AutoResetEvent ev = new(false);
        var ti = new TaskInfo();
        ti.handle = ThreadPool.RegisterWaitForSingleObject(ev, callBack: (state, timeout) => {
            var ti = state as TaskInfo;
            var cause = "TIMED OUT";
            if (!timeout) {
                cause = "SIGNALED";
                if (ti.handle is not null)
                    ti.handle.Unregister(null);  // 注销回调
            }
            Console.WriteLine("CallBack executes on thread {0}; cause = {1}.",
         Thread.CurrentThread.ManagedThreadId, cause);
        }, ti, 2000, false);
        ev.Set();
        Thread.Sleep(1000);
    }
}
```

>---
#### 1.1. 线程存储

`ThreadLocal<T>` 创建延迟初始化的线程本地对象。`[ThreadStatic]` 创建线程静态字段。`Thread.AllocateDataSlot` 方法创建数据槽，每个线程都有一个独立的槽值。`Thread.AllocateNamedDataSlot` / `Thread.GetNamedDataSlot` 方法创建命名数据槽，所有线程共享。


```csharp
class Sample {
    public static readonly Sample Instance = new Sample();
    ThreadLocal<int> local = new();
    public static void Method() {
        Instance.local.Value = Thread.CurrentThread.ManagedThreadId;
        Console.WriteLine(Instance.local.Value);
    }
    static void Main(string[] args) {
        Thread t = new Thread(Method);
        t.Start(); 
        t.Join();   // maybe 11
        Task.Run(async () => { 
            await Task.Delay(1000);
            Method();
        }).Wait();   // maybe 10
        Console.WriteLine(Instance.local.Value);   // 0
    }
}
```

>---
#### 1.2. 计时器

`Threading.Timer` 在 `ThreadPool` 上执行定时回调。
`Timers.Timer` 在 `ThreadPool` 上执行定时任务。
`Threading.PeriodicTimer` 在调用方等待滴答后执行工作。

```csharp
class Sample {
    static System.Threading.Timer threadTimer;
    static System.Timers.Timer timersTimer;
    static System.Threading.PeriodicTimer periodicTimer;
    static void Main(string[] args) {
        CancellationTokenSource CTS = new();
        // Threading.Timer
        threadTimer = new(
            callback: (state) => {
                var token = (CancellationToken)state;
                if (token.IsCancellationRequested) 
                    threadTimer?.Dispose();
                else Console.WriteLine($"Threading.Timer SignalTime: {DateTime.Now:yyyy/MM/dd hh:mm:ss}.");
            },
            state: CTS.Token, dueTime: 0, period: 1000);
        // Timers.Timer
        timersTimer = new(TimeSpan.FromSeconds(1));
        timersTimer.Elapsed += (sender, e) => {
            if (CTS.IsCancellationRequested) 
                timersTimer?.Dispose();
            else Console.WriteLine($"Timers.Timer SignalTime: {e.SignalTime}.");
        };
        timersTimer.Start();
        // Threading.PeriodicTimer
        periodicTimer = new(TimeSpan.FromSeconds(1));
        var task = Task.Run(async () => {
            while (await periodicTimer.WaitForNextTickAsync()) {
                if (CTS.IsCancellationRequested) 
                    periodicTimer?.Dispose();
                else Console.WriteLine($"Threading.PeriodicTimer SignalTime: {DateTime.Now:yyyy/MM/dd hh:mm:ss}");
            }
        });

        for (int i = 0; i < 3; i++)
            Task.Delay(1000).Wait();
        CTS.Cancel();
        Task.Delay(2000).Wait();
        Console.WriteLine($"{DateTime.Now:HH:mm:ss}: done.");
    }
}
```

>---
#### 1.3. 异步通信 Channel 

`Channel` 提供线程安全的异步通信机制，提供一组在生产者和消费者之间传递数据的同步数据结构。

```csharp
using System.Threading.Channels;
class Sample {
    static int counter = 0;
    static async Task Main(string[] args) {
        var channel = Channel.CreateUnbounded<int>(new() { SingleWriter = true });   // 无界通道
        var consumerTask = Task.Run(async () => {
            await foreach (var item in channel.Reader.ReadAllAsync()) {
                Console.CursorLeft = 0;
                Console.Write($"Consumed: {item}");
            }
        });
        var producerTask = async () => {
            for (int i = 0; i < 100; i++) {
                Interlocked.Increment(ref counter);
                await channel.Writer.WriteAsync(counter);
                await Task.Delay(new Random().Next(10, 50));
            }
        };
        Parallel.Invoke(producerTask().Wait, producerTask().Wait, producerTask().Wait);
        channel.Writer.Complete();
        await consumerTask;
        Console.WriteLine("\r\nCompleted!");  // 300
    }
}
```



---
### 2. 线程同步

#### 2.1. Monitor, lock

`Monitor` 获取或释放锁，提供对共享资源的互斥访问，可用于线程间同步。
`lock` 使用 `Monitor.Enter` 和 `Monitor.Exit` 实现，对象是 `Lock` 时使用 `Lock.EnterScope` 方法创建同步区域。

```csharp

```

>---
#### 2.2. SpinLock

`SpinLock` 旋转锁，在等待事件较短时性能优于其他锁类型。开发阶段可以启用线程追踪模式，锁不可重入时抛出异常。

```csharp
class Sample {
    static int counter;
    private readonly static object locker = new();
    public static void Increment() {
        for (int i = 0; i < 100000; i++) {
            Monitor.Enter(locker);
            counter++;
            Monitor.Exit(locker);
        }
    }
    static void Main() {
        Parallel.Invoke(Increment, Increment);
        Console.WriteLine(counter);  // 100000 * 2
    }
}
```

>---
#### 2.3. ReaderWriterLockSlim

`ReaderWriterLockSlim` 提供读写锁，支持读多写少场景。

```csharp

class Sample {
    private static readonly ReaderWriterLockSlim _lock = new ReaderWriterLockSlim(LockRecursionPolicy.SupportsRecursion);
    private static List<int> _sharedList = new List<int> { 1, 2, 3, 4, 5 };
    static void ReadAndUpdateList() {
        for (int i = 0; i < 3; i++) {
            try {
                _lock.EnterUpgradeableReadLock(); // 升级读锁独占模式，可升级为写锁
                Console.WriteLine($"Thread {Thread.CurrentThread.ManagedThreadId} reading list: [{string.Join(", ", _sharedList)}]");
                // 检查是否需要更新
                if (_sharedList.Count < 10) {
                    _lock.EnterWriteLock();
                    try {
                        int newValue = _sharedList.Count + 1;
                        _sharedList.Add(newValue);
                        Console.WriteLine($"Thread {Thread.CurrentThread.ManagedThreadId} added value: {newValue}");
                    } finally {
                        _lock.ExitWriteLock();
                    }
                }
            } finally {
                _lock.ExitUpgradeableReadLock();
            }
            Thread.Sleep(100);
        }
    }
    static void Main(string[] args) {
        Task[] tasks = [Task.Run(ReadAndUpdateList), Task.Run(ReadAndUpdateList), Task.Run(ReadAndUpdateList)];
        Task.WaitAll(tasks);
    }
}
```

>---
#### 2.4. Semaphore, SemaphoreSlim

`Semaphore` 表示一个系统范围内的命名信号量或本地信号量，限制可同时访问共享资源或资源池的线程数，没有线程相关性。`SemaphoreSlim` 用于单个进程边界内的同步。

```csharp
class Sample {
    const string name = "Global\\MySemaphore";
    static void Main(string[] args) {
        // 创建一个系统级命名信号量，初始计数为 2，最大计数为 3
        using Semaphore semaphore = new Semaphore(2, 3, name);
        // 模拟 5 个线程尝试访问资源
        var tasks = new Task[5];
        for (int i = 0; i < 5; i++) {
            var taskId = i;
            tasks[taskId] = Task.Run(() => {
                Console.WriteLine($"Task {taskId} waiting to access resource");
                var s = Semaphore.OpenExisting(name);
                s?.WaitOne();
                try {
                    Console.WriteLine($"Task {taskId} accessing resource");
                    // 模拟资源访问
                    Thread.Sleep(2000);
                } finally {
                    Console.WriteLine($"Task {taskId} releasing resource");
                    s?.Release();
                }
            });
        }
        Task.WaitAll(tasks);
        Console.WriteLine("All tasks completed");
    }
}
```

>---
#### 2.5. EventWaitHandle

`EventWaitHandle` 允许线程通过信号和等待信号相互通信，可自动或手动重置事件，命名句柄进程间可见。`AutoResetEvent` 为自动版本，`ManualResetEvent` 为手动版本，仅为本地等待句柄。

```csharp
class Sample {
    static void Main(string[] args) {
        using EventWaitHandle ew = new EventWaitHandle(true, EventResetMode.AutoReset, "AUTO");
        Task[] tasks = new Task[5];
        for (int i = 0; i < 5; i++) {
            var taskId = i;
            tasks[taskId] = Task.Run(() => {
                Console.WriteLine($"Task {taskId} waiting to access resource");
                var ew = EventWaitHandle.OpenExisting("AUTO");
                ew.WaitOne();
                try {
                    Console.WriteLine($"Task {taskId} accessing resource");
                    // 模拟资源访问
                    Thread.Sleep(200);
                } finally {
                    Console.WriteLine($"Task {taskId} releasing resource");
                    ew.Set();  // auto reset
                }
            });
        }
        Task.WaitAll(tasks);
    }
}
```

>---
#### 2.6. CountdownEvent

`CountdownEvent` 在收到特定次数的信号（`Signal()`）后，取消阻塞等待线程的同步基元。

```csharp
class Sample {
    static void Main(string[] args) {
        using CountdownEvent cte = new CountdownEvent(5);
        for (int i = 0; i < 5; i++) {
            int taskId = i;
            Task.Run(() => {
                Task.Delay(2000).Wait();
                Console.WriteLine($"Task {taskId} completed");
                cte.Signal();
            });
        }
        cte.Wait();
    }
}
```

>---
#### 2.7. Barrier

`Barrier` 使多个线程（参与者）可以分阶段并发执行，直至到达某个屏障点被阻塞。当所有参与者到达屏障点后，执行一次回调（*PostAction*），并释放所有参与者继续执行。

```csharp
class Sample {
    const int N = 2;
    static bool isCompleted = false;
    static Barrier barrier = new Barrier(N, postPhaseAction: (b) => {
        Console.WriteLine("All participant arrive at barrier");
        if (b.CurrentPhaseNumber >= 2) {  // 0，1，2
            isCompleted = true;
            Console.WriteLine("\r\nTasks completed");
        }
    }); // b.CurrentPhaseNumber ++
    static void Main(string[] args) {
        Task[] tasks = new Task[N];
        for (int i = 0; i < N; i++) {
            var taskId = i;
            tasks[taskId] = Task.Run(() => {
                while (!isCompleted) {
                    Console.WriteLine($"Do some work on Task {taskId}");
                    barrier.SignalAndWait();
                }
            });
        }
        Task.WaitAll(tasks);
    }
}
```

>---
#### 2.8. Interlocked

`Interlocked` 为多个线程共享的变量提供原子操作。

```csharp
class Sample {
    static int Counter = 0;
    public static void Increment() {
        for (int j = 0; j < 500000; j++)
            Interlocked.Increment(ref Counter);
    }
    static void Main(string[] args) {
        Parallel.Invoke(Increment, Increment, Increment);
        Console.WriteLine(Counter);  // 3* 500000
    }
}
```

>---
#### 2.9. SpinWait

`SpinWait` 自旋等待，避免线程切换开销，用于短时间等待的情况。

```csharp
// 无锁 SpinWait
public class LightLock {
    private int state = 0;
    public void Enter() {
        SpinWait spinner = new SpinWait();
        while (Interlocked.CompareExchange(ref state, 1, 0) != 0) {
            spinner.SpinOnce();
        }
    }
    public void Exit() => Interlocked.Exchange(ref state, 0);
}
class Sample {
    static int counter = 0;
    static LightLock lightLock = new LightLock();
    static void Increment() {
        for (int j = 0; j < 500000; j++) {
            lightLock.Enter();
            counter++;
            lightLock.Exit();
        }
    }
    static void Main(string[] args) {
        Parallel.Invoke(Increment, Increment, Increment);
        Console.WriteLine(counter);  // 3 * 500000
    }
}
```

>---
#### 2.10. MethodImpAttribute

`[MethodImp(MethodImplOptions.Synchronized)]` 修饰方法为同步方法。   

```csharp
using System.Runtime.CompilerServices;
class Sample {
    static int counter;
    [MethodImpl(MethodImplOptions.Synchronized)]
    static void Increment() => ++counter;
    static void Main(string[] args) {
        var increment = () => {
            for (int i = 0; i < 500000; ++i)
                Increment();
        };
        Parallel.Invoke(increment, increment, increment);
        Console.WriteLine(counter);  // 3 * 500000
    }
}
```

---