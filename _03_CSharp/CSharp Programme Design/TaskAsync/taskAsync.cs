using System.Runtime.CompilerServices;


public interface IAwaitable {
    IAwaiter GetAwaiter();
}
public interface IAwaiter : INotifyCompletion {
    bool IsCompleted { get; }
    void GetResult();
}
public interface IAwaitable<out TResult> {
    IAwaiter<TResult> GetAwaiter();
}

public interface IAwaiter<out TResult> : INotifyCompletion // or ICriticalNotifyCompletion
{
    bool IsCompleted { get; }

    TResult GetResult();
}
public class TaskAsync : IAwaitable {
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
    public IAwaiter GetAwaiter() {
        return new TaskAsyncAwaiter(this);
    }
}
public class TaskAsync<TResult> : TaskAsync, IAwaitable<TResult> {
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
    public new IAwaiter<TResult> GetAwaiter() {
        return new TaskAsyncAwaiter<TResult>(this);
    }
}
// Awaiter 实现
public class TaskAsyncAwaiter(TaskAsync task) : IAwaiter {
    public bool IsCompleted => task.IsCompleted;
    public void GetResult() => task.Wait();
    public void OnCompleted(Action continuation) => task.ContinueWith(continuation);
}
public class TaskAsyncAwaiter<TResult>(TaskAsync<TResult> task) : IAwaiter<TResult> {
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
    private class TaskAsyncImpl<TResult>(Func<TResult> func) : TaskAsync<TResult> {
        protected override void Execute() => Complete(func());
    }
}