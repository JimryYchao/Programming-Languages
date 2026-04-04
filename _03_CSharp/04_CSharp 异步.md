## CSharp 异步

---

### APM 回调异步模式

APM 基于回调的异步模式，基于 `IAsyncResult` 接口和 `BeginXxx` / `EndXxx` 方法对。
- `BeginXxx` 开始异步操作并返回实现 `IAsyncResult` 对象。
- `EndXxx` 完成异步操作并返回结果。

```csharp
public interface IAsyncResult {
    object? AsyncState { get; }           // 可选特定于程序，包含异步操作的信息
    WaitHandle AsyncWaitHandle { get; }   // 阻塞等待异步完成
    bool CompletedSynchronously { get; }  // 指示异步操作是否在 BeginXxx 线程上完成
    bool IsCompleted { get; }             
}
```

>---
#### APM 示例

> 异步调用委托示例

```csharp
public interface IDelegateAsync {
    IAsyncResult BeginInvoke(dynamic[]? args = null, AsyncCallback? callback = null, object? state = null);
    void EndInvoke(IAsyncResult asyncResult);
}

public interface IDelegateAsync<T> {
    IAsyncResult BeginInvoke(dynamic[]? args = null, AsyncCallback? callback = null, object? state = null);
    T EndInvoke(IAsyncResult asyncResult);
}

public abstract class DelegateAsync {
    public static IDelegateAsync<TResult> Create<Func, TResult>(Func func) where Func : Delegate {
        ArgumentNullException.ThrowIfNull(func, nameof(func));
        return new InternalDelegateAsync<Func, TResult>(func);
    }
    public static IDelegateAsync Create<Action>(Action action) where Action : Delegate {
        ArgumentNullException.ThrowIfNull(action, nameof(action));
        return new InternalDelegateAsync<Action>(action);
    }
    private class AsyncResult(AsyncCallback? callback, object? state) : IAsyncResult {
        public Exception? exception = null;
        private readonly ManualResetEvent waitHandle = new ManualResetEvent(false);
        public object? AsyncState { get; init; } = state;
        public WaitHandle AsyncWaitHandle => waitHandle;
        public bool CompletedSynchronously => false;
        public bool IsCompleted { get; private set; } = false;
        public void Complete() {
            IsCompleted = true;
            waitHandle.Set();
            callback?.Invoke(this);
        }
    }
    private class InternalDelegateAsync<Action>(Action action) : IDelegateAsync where Action : Delegate {
        public IAsyncResult BeginInvoke(dynamic[]? args, AsyncCallback? callback, object? state) {
            var asRt = new AsyncResult(callback, state);
            ThreadPool.QueueUserWorkItem((state) => {
                try {
                    action.DynamicInvoke(args);
                } catch (Exception ex) {
                    asRt.exception = ex;
                } finally {
                    asRt.Complete();
                }
            });
            return asRt;
        }
        public void EndInvoke(IAsyncResult asyncResult) {
            ArgumentNullException.ThrowIfNull(asyncResult);
            if (asyncResult is AsyncResult acResult) {
                if (acResult.exception is not null)
                    throw acResult.exception;
                if (!acResult.IsCompleted)
                    acResult.AsyncWaitHandle.WaitOne();
            } else
                throw new ArgumentException("Invalid async result", nameof(asyncResult));
        }
    }
    private class InternalDelegateAsync<Func, T>(Func func) : IDelegateAsync<T> where Func : Delegate {
        private class AsyncResult(AsyncCallback? callback, object? state) : DelegateAsync.AsyncResult(callback, state) {
            public T result = default!;
        }
        public IAsyncResult BeginInvoke(dynamic[]? args = null, AsyncCallback? callback = null, object? state = null) {
            var funcRt = new AsyncResult(callback, state);
            ThreadPool.QueueUserWorkItem((state) => {
                try {
                    funcRt.result = (T)func.DynamicInvoke(args)!;
                } catch (Exception ex) {
                    funcRt.exception = ex;
                } finally {
                    funcRt.Complete();
                }
            });
            return funcRt;
        }
        public T EndInvoke(IAsyncResult asyncResult) {
            ArgumentNullException.ThrowIfNull(asyncResult);
            if (asyncResult is AsyncResult funcRT) {
                if (funcRT.exception is not null)
                    throw funcRT.exception;
                if (!funcRT.IsCompleted)
                    funcRT.AsyncWaitHandle.WaitOne();
            } else
                throw new ArgumentException("Invalid async result", nameof(asyncResult));
            return funcRT.result;
        }
    }
}
```
```csharp
class Program {
    [ThreadStatic]
    public static int value = 0;
    static void Main(string[] args) {
        List<WaitHandle> handles = new();
        for (int i = 0; i < 10; i++) {
            var ac = DelegateAsync.Create<Func<int, int>, int>((i) => {
                value = 0;
                for (int j = 0; j <= i; j++) {
                    Thread.Sleep(new Random().Next(50, 200));
                    Console.WriteLine($"{Thread.CurrentThread.ManagedThreadId} : {value++}");
                }
                return value;
            });
            handles.Add(ac.BeginInvoke([i]).AsyncWaitHandle);
        }
        // 等待全部异步完成
        WaitHandle.WaitAll([.. handles]);
    }
}
```

---
### EAP 事件异步模式

EAP 基于事件的异步模式。
- `XxxAsync` 同于启动异步操作，可能包含同步版本的镜像。
- `XxxCompleted` 事件在任务完成时触发。
- `AsyncCompletedEventArgs` 事件参数，包含异步操作的结果和异常信息。
- (op)`CancelAsync` 取消异步操作。

```csharp
using System.ComponentModel;

public class AsyncExample {
    // 同步方法
    public int MethodName(string param);
    // 异步方法
    public void MethodNameAsync(string param);   // 不可取消
    public void MethodNameAsync(string param, object userState);  // 可取消
    public class MethodNameCompletedEventArgs : AsyncCompletedEventArgs {
        // MethodNameAsync 具有返回值，void 直接使用 AsyncCompletedEventArgs
        public int /* TResult */ Result { get; }
    }
    public delegate void MethodNameCompletedEventHandler(object sender, MethodNameCompletedEventArgs e);
    public event MethodNameCompletedEventHandler MethodNameCompleted;
    // (可选)取消异步操作
    public void CancelAsync(object userState);            // 支持取消多个异步操作
    public void MethodNameCancelAsync(object userState);  // 取消 MethodNameAsync 操作
    // (可选)类不支持多个并发调用
    public bool IsBusy { get; }  
    // (可选)提供对异步进度报告的支持。范围 0~100
    public event ProgressChangedEventHandler ProgressChanged;            
    public event ProgressChangedEventHandler MethodNameProgressChanged;  
    // (可选)提供对返回增量结果的支持

    // 处理方法中的引用参数 ref,out
    public void RefMethodName(string p1, ref int p2, out ref int p3);
    public class RefMethodNameCompletedEventArgs : System.ComponentModel.AsyncCompletedEventArgs
    {
        public int Result { get; };
        public string Arg2 { get; };
        public string Arg3 { get; };
    }
    public delegate void RefMethodNameCompletedEventHandler(object sender, RefMethodNameCompletedEventArgs e);
    public event RefMethodNameCompletedEventHandler RefMethodNameCompleted;
}
```

>---
#### EAP 示例

```csharp

using System.ComponentModel;
public interface IEventActionAsync {
    void InvokeSync(dynamic[]? args);
    void InvokeAsync(dynamic[]? args);
    void InvokeAsync(dynamic[]? args, object? userState = null);
    event EventHandler<AsyncCompletedEventArgs> InvokeCompleted;
    void CancelAsync(object userState);
    bool IsBusy { get; }
}
public abstract class EventActionAsync {
    public static IEventDelegateAsync Create<Action>(Action action) where Action : Delegate {
        ArgumentNullException.ThrowIfNull(action, nameof(action));
        return new InternalEventDelegateAsync<Action>(action);
    }
    private class InternalEventDelegateAsync<Action>(Action action) : IEventDelegateAsync where Action : Delegate {
        private readonly Lock locker = new();
        private int isBusy;
        private readonly Dictionary<object, (CancellationTokenSource CTS, AsyncOperation Operation)> operations = new();
        public event EventHandler<AsyncCompletedEventArgs>? InvokeCompleted;

        public void InvokeSync(dynamic[]? args) {
            action?.DynamicInvoke(args);
        }
        public void InvokeAsync(dynamic[]? args) {
            InvokeAsync(args, null);
        }
        public bool IsBusy => Interlocked.CompareExchange(ref isBusy, 0, 0) == 1;
        public void InvokeAsync(dynamic[]? args, object? userState = null) {
            lock (locker) {
                if (IsBusy && userState == null)
                    throw new InvalidOperationException("An asynchronous operation is already in progress.");
                Interlocked.Exchange(ref isBusy, 1);  // 标记为忙碌
                                                      //var cts = new CancellationTokenSource();  // 创建取消令牌
                var Operation = AsyncOperationManager.CreateOperation(userState);
                var CTS = new CancellationTokenSource();
                if (userState != null)
                    operations[userState] = (CTS, Operation);
                // 在线程池线程上执行操作
                ThreadPool.QueueUserWorkItem(_ => {
                    Exception? error = null;
                    bool cancelled = false;
                    try {
                        // 检查是否已取消
                        if (userState != null && operations.TryGetValue(userState, out var Op))
                            Op.CTS.Token.ThrowIfCancellationRequested();
                        action.DynamicInvoke(args);
                    } catch (OperationCanceledException) {
                        cancelled = true;
                    } catch (Exception ex) {
                        error = ex;
                    } finally {
                        // 清理取消令牌
                        if (userState != null)
                            operations.Remove(userState);
                        Interlocked.Exchange(ref isBusy, 0);
                        Operation.PostOperationCompleted((state) => {
                            var (_error, _cancelled, _userState) = ((Exception, bool, object))state!;
                            OnInvokeCompleted(new(_error, _cancelled, _userState));
                        }, (error, cancelled, userState));
                    }
                });
            }
        }
        public void CancelAsync(object userState) {
            if (userState != null && operations.TryGetValue(userState, out var Op))
                Op.CTS.Cancel();
        }
        protected virtual void OnInvokeCompleted(AsyncCompletedEventArgs e) {
            InvokeCompleted?.Invoke(this, e);
        }
    }
}
```
```csharp
class Program {
    static void Main(string[] args) {

        Console.WriteLine("Testing Action async...");
        var actionAsync = EventActionAsync.Create(() => {
            Console.WriteLine($"Action executing on thread {Thread.CurrentThread.ManagedThreadId}");
            Thread.Sleep(2000); // 模拟耗时操作
            Console.WriteLine("Action completed");
        });
        actionAsync.InvokeCompleted += (sender, e) => {
            if (e.Error != null) {
                Console.WriteLine($"Action error: {e.Error.Message}");
            } else if (e.Cancelled) {
                Console.WriteLine("Action cancelled");
            } else {
                Console.WriteLine("Action completed successfully");
            }
        };
        actionAsync.InvokeAsync(null);
        while (actionAsync.IsBusy)
            continue;   // 阻塞

        // 测试带用户状态
        Console.WriteLine("\nTesting with user state...");
        var stateAsync = EventActionAsync.Create(() => {
            Console.WriteLine($"State action executing on thread {Thread.CurrentThread.ManagedThreadId}");
            Thread.Sleep(1000);
        });
        stateAsync.InvokeCompleted += (sender, e) => {
            Console.WriteLine($"State action completed. User state: {e.UserState}");
        };
        stateAsync.InvokeAsync(null, "TestState123");
        while (stateAsync.IsBusy)
            continue;   // 阻塞

        // 测试取消
        Console.WriteLine("\nTesting with cenceling ...");
        var cancelAsync = EventActionAsync.Create(() => {
            Console.WriteLine($"Cancel action executing on thread {Thread.CurrentThread.ManagedThreadId}");
        });
        cancelAsync.InvokeCompleted += (sender, e) => {
            if (e.Error != null) {
                Console.WriteLine($"Cancel action error: {e.Error.Message}");
            } else if (e.Cancelled) {
                Console.WriteLine("Cancel action cancelled");
            } else {
                Console.WriteLine($"Cancel action completed. User state: {e.UserState}");
            }
        };
        cancelAsync.InvokeAsync(null, "TestCancel");
        cancelAsync.CancelAsync("TestCancel");
        while (cancelAsync.IsBusy)
            continue;   // 阻塞

        // 测试异常
        Console.WriteLine("\nTesting with exception ...");
        var exceptionAsync = EventActionAsync.Create(() => {
            Console.WriteLine($"Exception action executing on thread {Thread.CurrentThread.ManagedThreadId}");
            throw new Exception("Throw an exception");
        });
        exceptionAsync.InvokeCompleted += (sender, e) => {
            if (e.Error != null) {
                Console.WriteLine($"Exception action error: {e.Error.Message}");
            } else if (e.Cancelled) {
                Console.WriteLine("Exception action cancelled");
            } else {
                Console.WriteLine($"Exception action completed. User state: {e.UserState}");
            }
        };
        exceptionAsync.InvokeAsync(null, "TestException");
        while (exceptionAsync.IsBusy)
            continue;   // 阻塞
    }
}
```

---
### TAP 任务异步模式

TAP 基于任务的异步模式，基于 `System.Threading.Tasks` 的 `Task` 和 `Task<TResult>`、`ValueTask` 和 `ValueTask<TResult>`。

`async` / `await` 底层依靠异步状态机（AsyncStateMachine）、Task Awaiter 模式、同步上下文（SynchronizationContext）三大核心机制实现。`await` 挂起（非阻塞）当前方法执行线程并释放线程回底层线程池；异步操作完成后通过状态机恢复执行。底层执行过程：

- 编译器为 `async` 方法生成一个异步状态机（AsyncStateMachine）：合成一个 `IAsyncStateMachine` 嵌套结构体。同时把 `await` 断点拆分为多个执行片段，对应状态机的每个执行状态（state）。
- 初始化并启动状态机：调用异步方法，首先初始化状态机实例，设置初始状态、捕获当前同步上下文（SynchronizationContext）、保存方法参数与局部变量，随后调用状态机的 `MoveNext()` 方法，执行同步代码片段，直到遇见第一个 `await`。
- `await` 断点处判断：检查异步任务是否完成。首先编译器获取对应 `Task` 的 `Awaiter` 对象（`INotifyCompletion` / `ICriticalNotifyCompletion` 实现），调用 `Awaiter.IsCompleted` 检查异步任务是否完成：若任务已同步完成，跳过挂起逻辑，直接执行后续；若未完成，则将状态机当前状态保存，调用 `Awaiter.OnCompleted()` 方法，把状态机的 `MoveNext` 方法注册为回调并挂起状态机，释放当前线程，方法立即返回一个未完成的 `Task` 对象给调用方。
- 异步完成后触发状态机恢复：底层异步操作（IO、网络、数据库等）完成后，操作系统或底层驱动会通知 CLR，触发之前注册的 `OnCompleted` 回调，回调内部会再次调用状态机的 `MoveNext()` 方法，唤醒状态机。
- 恢复执行后续同步逻辑：状态机根据保存的 `state` 字段，定位到 `await` 断点后的代码片段，通过同步上下文切换到对应执行线程（UI 线程 / 线程池线程），执行后续代码；若有返回值，设置 `Task` 结果；若有异常，捕获并封装到 `Task` 中。

```csharp
// 原始代码
internal class AsyncMethods
{
    internal static async Task<int> MultiCallMethodAsync(int arg0, int arg1, int arg2, int arg3)
    {
        HelperMethods.Before();
        int resultOfAwait1 = await MethodAsync(arg0, arg1);
        HelperMethods.Continuation1(resultOfAwait1);
        int resultOfAwait2 = await MethodAsync(arg2, arg3);
        HelperMethods.Continuation2(resultOfAwait2);
        int resultToReturn = resultOfAwait1 + resultOfAwait2;
        return resultToReturn;
    }
}
// 编译后
internal class CompiledAsyncMethods
{
    [DebuggerStepThrough]
    [AsyncStateMachine(typeof(MultiCallMethodAsyncStateMachine))] // async
    internal static /*async*/ Task<int> MultiCallMethodAsync_(int arg0, int arg1, int arg2, int arg3)
    {
        MultiCallMethodAsyncStateMachine multiCallMethodAsyncStateMachine = new MultiCallMethodAsyncStateMachine()
            {
                Arg0 = arg0,
                Arg1 = arg1,
                Arg2 = arg2,
                Arg3 = arg3,
                ResultToReturn = new TaskCompletionSource<int>(),
                // -1: Begin
                //  0: 1st await is done
                //  1: 2nd await is done
                //     ...
                // -2: End
                State = -1
            };
        (multiCallMethodAsyncStateMachine as IAsyncStateMachine).MoveNext(); // Original code are in this method.
        return multiCallMethodAsyncStateMachine.ResultToReturn.Task;
    }
}
[CompilerGenerated]
[StructLayout(LayoutKind.Auto)]
internal struct MultiCallMethodAsyncStateMachine : IAsyncStateMachine
{
    // State:
    // -1: Begin
    //  0: 1st await is done
    //  1: 2nd await is done
    //     ...
    // -2: End
    public int State;
    public TaskCompletionSource<int> ResultToReturn; // int resultToReturn ...
    public int Arg0; // int Arg0
    public int Arg1; // int arg1
    public int Arg2; // int arg2
    public int Arg3; // int arg3
    public int ResultOfAwait1; // int resultOfAwait1 ...
    public int ResultOfAwait2; // int resultOfAwait2 ...
    private Task<int> currentTaskToAwait;

    /// <summary>
    /// Moves the state machine to its next state.
    /// </summary>
    void IAsyncStateMachine.MoveNext()
    {
        try
        {
            switch (this.State)
            {
                IAsyncStateMachine that = this; // Cannot use "this" in lambda so create a local copy. 
                // Orginal code is splitted by "case"s:
                // case -1:
                //      HelperMethods.Before();
                //      MethodAsync(Arg0, arg1);
                // case 0:
                //      int resultOfAwait1 = await ...
                //      HelperMethods.Continuation1(resultOfAwait1);
                //      MethodAsync(arg2, arg3);
                // case 1:
                //      int resultOfAwait2 = await ...
                //      HelperMethods.Continuation2(resultOfAwait2);
                //      int resultToReturn = resultOfAwait1 + resultOfAwait2;
                //      return resultToReturn;
                case -1: // -1 is begin.
                    HelperMethods.Before(); // Code before 1st await.
                    this.currentTaskToAwait = AsyncMethods.MethodAsync(this.Arg0, this.Arg1); // 1st task to await
                    // When this.currentTaskToAwait is done, run this.MoveNext() and go to case 0.
                    this.State = 0;
                    this.currentTaskToAwait.ContinueWith(_ => that.MoveNext()); // Callback
                    break;
                case 0: // Now 1st await is done.
                    this.ResultOfAwait1 = this.currentTaskToAwait.Result; // Get 1st await's result.
                    HelperMethods.Continuation1(this.ResultOfAwait1); // Code after 1st await and before 2nd await.
                    this.currentTaskToAwait = AsyncMethods.MethodAsync(this.Arg2, this.Arg3); // 2nd task to await
                    // When this.currentTaskToAwait is done, run this.MoveNext() and go to case 1.
                    this.State = 1;
                    this.currentTaskToAwait.ContinueWith(_ => that.MoveNext()); // Callback
                    break;
                case 1: // Now 2nd await is done.
                    this.ResultOfAwait2 = this.currentTaskToAwait.Result; // Get 2nd await's result.
                    HelperMethods.Continuation2(this.ResultOfAwait2); // Code after 2nd await.
                    int resultToReturn = this.ResultOfAwait1 + this.ResultOfAwait2; // Code after 2nd await.
                    // End with resultToReturn. No more invocation of MoveNext().
                    this.State = -2; // -2 is end.
                    this.ResultToReturn.SetResult(resultToReturn);
                    break;
            }
        }
        catch (Exception exception)
        {
            // End with exception.
            this.State = -2; // -2 is end. Exception will also when the execution of state machine.
            this.ResultToReturn.SetException(exception);
        }
    }

    /// <summary>
    /// Configures the state machine with a heap-allocated replica.
    /// </summary>
    /// <param name="stateMachine">The heap-allocated replica.</param>
    [DebuggerHidden]
    void IAsyncStateMachine.SetStateMachine(IAsyncStateMachine stateMachine)
    {
        // No core logic.
    }
}
```

>---
#### TAP 示例

```csharp
// 一个简易的任务类
using System.Runtime.CompilerServices;
public class TaskAsync {
    // 同步原语
    private readonly ManualResetEvent completionEvent = new(false);
    private readonly Lock locker = new();
    // 任务状态
    public enum TaskStatus {
        Created,     // 已创建但未开始
        Running,     // 正在执行
        RanToCompletion, // 成功完成
        Faulted,     // 执行过程中发生异常
        Canceled     // 被取消
    }
    private readonly List<Action> continuations = new();
    public TaskStatus Status { get; private set; } = TaskStatus.Created;
    public Exception Exception { get; private set; } = null!;
    public bool IsCompleted => Status == TaskStatus.RanToCompletion ||
                              Status == TaskStatus.Faulted ||
                              Status == TaskStatus.Canceled;
    protected virtual void Execute() { }
    protected void Complete() {
        lock (locker) {
            if (IsCompleted) return;
            Status = TaskStatus.RanToCompletion;
            completionEvent.Set();  // 完成时取消 this.Wait 堵塞
            // 执行所有延续任务
            foreach (var continuation in continuations) 
                ThreadPool.QueueUserWorkItem(_ => continuation());
        }
    }
    protected void Fault(Exception ex) {
        lock (locker) {
            if (IsCompleted) return;
            Status = TaskStatus.Faulted;
            Exception = ex;
            completionEvent.Set();
            // 执行所有延续任务
            foreach (var continuation in continuations) 
                ThreadPool.QueueUserWorkItem(_ => continuation());
        }
    }
    public bool Start() {
        lock (locker) {
            if (Status != TaskStatus.Created)
                return false;
            Status = TaskStatus.Running;
        }
        // 在线程池上执行任务
        ThreadPool.QueueUserWorkItem(_ => {
            try {
                Execute();
                Complete();
            } catch (Exception ex) {
                Fault(ex);
            }
        });
        return true;
    }
    public void Wait() {
        completionEvent.WaitOne();
        if (Status == TaskStatus.Faulted && Exception != null) 
            throw Exception;
    }
    public bool Wait(int millisecondsTimeout) {
        bool completed = completionEvent.WaitOne(millisecondsTimeout);
        if (completed && Status == TaskStatus.Faulted && Exception != null) 
            throw Exception;
        return completed;
    }
    public TaskAsync ContinueWith(Action action) {
        var continuationTask = new TaskAsync();
        lock (locker) {
            if (IsCompleted) {
                ThreadPool.QueueUserWorkItem(_ => {
                    try {
                        action();
                        continuationTask.Complete();
                    } catch (Exception ex) {
                        continuationTask.Fault(ex);
                    }
                });
            } else {
                // 添加到延续列表
                continuations.Add(() => {
                    try {
                        action();
                        continuationTask.Complete();
                    } catch (Exception ex) {
                        continuationTask.Fault(ex);
                    }
                });
            }
        }
        return continuationTask;
    }
    public TaskAsyncAwaiter GetAwaiter() {
        return new TaskAsyncAwaiter(this);
    }
}
public class TaskAsync<TResult> : TaskAsync {
    private TResult result;
    public TResult Result {
        get {
            Wait();
            if (Status == TaskStatus.Faulted && Exception != null) 
                throw Exception;
            return result;
        }
    }
    public void Complete(TResult result) {
        this.result = result;
        base.Complete();
    }
    public new TaskAsyncAwaiter<TResult> GetAwaiter() {
        return new TaskAsyncAwaiter<TResult>(this);
    }
}
// Awaiter 实现
public class TaskAsyncAwaiter(TaskAsync task) : INotifyCompletion {
    public bool IsCompleted => task.IsCompleted;
    public void GetResult() => task.Wait();
    public void OnCompleted(Action continuation) => task.ContinueWith(continuation);
}
public class TaskAsyncAwaiter<TResult> (TaskAsync<TResult> task) : INotifyCompletion {
    public bool IsCompleted => task.IsCompleted;
    public TResult GetResult() => task.Result;
    public void OnCompleted(Action continuation) => task.ContinueWith(continuation);
}

// 任务工厂
public static class TaskAsyncFactory {
    public static TaskAsync Run(Action action) {
        var task = new TaskAsyncImpl(action);
        task.Start();
        return task;
    }
    public static TaskAsync<TResult> Run<TResult>(Func<TResult> function) {
        var task = new TaskAsyncImpl<TResult>(function);
        task.Start();
        return task;
    }
    public static void WaitAll(params TaskAsync[] tasks) {
        foreach (var task in tasks) {
            task.Wait();
        }
    }
    private class TaskAsyncImpl(Action action) : TaskAsync {
        protected override void Execute() => action();
    }
    private class TaskAsyncImpl<TResult> (Func<TResult> func): TaskAsync<TResult> {
        protected override void Execute() => Complete(func());
    }
}
```
```csharp
public class TaskAsyncTest {
    static IEnumerable<TaskAsync<Task<int>>> Counter() {
        int v = 0;
        for (int i = 0; i < 10; i++)
            yield return TaskAsyncFactory.Run(async () => {
                await Task.Delay(new Random().Next(50,200));
                return v++;
            });
    }
    public static async Task Test() {
        foreach (var t in Counter())
            Console.WriteLine(t.Result.Result);
    }
}
```

---