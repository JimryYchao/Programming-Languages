using System.ComponentModel;

// 基于 EAP 模式的异步委托接口
public class InvokeCompletedEventArgs<TResult> : AsyncCompletedEventArgs {
    public TResult Result { get; }
    public InvokeCompletedEventArgs(TResult result, Exception error, bool cancelled, object userState)
        : base(error, cancelled, userState) {
        Result = result;
    }
}

public interface IEventDelegateCommon<EventArgs> where EventArgs : System.EventArgs {
    void InvokeAsync(dynamic[]? args);
    void InvokeAsync(dynamic[]? args, object? userState = null);
    event EventHandler<EventArgs> InvokeCompleted;
    void CancelAsync(object userState);
}
public interface IEventDelegateAsync : IEventDelegateCommon<AsyncCompletedEventArgs> {
    bool IsBusy { get; }
}
public interface IEventDelegateAsync<T> : IEventDelegateCommon<InvokeCompletedEventArgs<T>> {
    bool IsBusy { get; }
}
public interface IParallelEventDelegateAsync : IEventDelegateCommon<AsyncCompletedEventArgs> {
    bool Reset(object? userState);
}
public interface IParallelEventDelegateAsync<T> : IEventDelegateCommon<InvokeCompletedEventArgs<T>>;


public abstract partial class EventDelegateAsync {
    public static IEventDelegateAsync Create<Action>(Action action) where Action : Delegate {
        ArgumentNullException.ThrowIfNull(action, nameof(action));
        return new InternalEventDelegateAsync<Action>(action);
    }
    public static IEventDelegateAsync<TResult> Create<Func, TResult>(Func func) where Func : Delegate {
        ArgumentNullException.ThrowIfNull(func, nameof(func));
        return new InternalEventDelegateAsync<Func, TResult>(func);
    }
    public static IParallelEventDelegateAsync CreateParallel<Action>(Action action) where Action : Delegate {
        ArgumentNullException.ThrowIfNull(action, nameof(action));
        return new InternalEventDelegateAsyncParallel<Action>(action);
    }

    public static IParallelEventDelegateAsync<TResult> CreateParallel<Func, TResult>(Func func) where Func : Delegate {
        ArgumentNullException.ThrowIfNull(func, nameof(func));
        return new InternalEventDelegateAsyncParallel<Func, TResult>(func);
    }
}

/// <summary>
/// 基于事件的异步委托调用实现
/// </summary>
public abstract partial class EventDelegateAsync {
    private abstract class BaseEventDelegate<Dele, OptType, EventArgs, TResult> : IEventDelegateCommon<EventArgs>
   where EventArgs : System.EventArgs
   where Dele : Delegate
   where OptType : class {
        protected readonly Lock locker = new();
        protected readonly Dictionary<object, (CancellationTokenSource CTS, OptType Ops)> operations = new();
        protected int isBusy;
        public event EventHandler<EventArgs>? InvokeCompleted;
        public virtual void CancelAsync(object userState) {
            if (userState != null && operations.TryGetValue(userState, out var Op))
                Op.CTS.Cancel();
        }
        public virtual void InvokeAsync(dynamic[]? args) {
            InvokeAsync(args, null);
        }
        protected virtual void OnInvokeCompleted(EventArgs e) {
            InvokeCompleted?.Invoke(this, e);
        }
        public abstract void InvokeAsync(dynamic[]? args, object? userState = null);

        protected abstract void CallDelegate(dynamic[]? args, out TResult result);
        protected abstract void postOperationCompleted(ref AsyncOperation operation, ref Exception error, ref bool cancelled, ref object? userState, ref TResult result);
        protected virtual void InvokeAsyncBusy(dynamic[]? args, object? userState = null) {
            CancellationTokenSource? CTS = null;
            AsyncOperation? _operation = null;
            OptType? Operation = default;
            TResult? result = default;
            lock (locker) {
                if (isBusy == 1 && userState == null)
                    throw new InvalidOperationException("An asynchronous operation is already in progress.");
                Interlocked.Exchange(ref isBusy, 1);  // 标记为忙碌
                                                      //var cts = new CancellationTokenSource();  // 创建取消令牌
                _operation = AsyncOperationManager.CreateOperation(userState);
                Operation = _operation as OptType;
                CTS = new CancellationTokenSource();
                if (userState != null)
                    operations[userState] = (CTS, Operation!);
                // 在线程池线程上执行操作
                ThreadPool.QueueUserWorkItem(_ => {
                    Exception? error = null;
                    bool cancelled = false;
                    try {
                        // 检查是否已取消
                        if (userState != null && operations.TryGetValue(userState, out var Op))
                            Op.CTS.Token.ThrowIfCancellationRequested();
                        CallDelegate(args, out result);
                    } catch (OperationCanceledException) {
                        cancelled = true;
                    } catch (Exception ex) {
                        error = ex;
                    } finally {
                        // 清理取消令牌
                        if (userState != null)
                            operations.Remove(userState);
                        Interlocked.Exchange(ref isBusy, 0);
                        postOperationCompleted(ref _operation!, ref error!, ref cancelled, ref userState, ref result!);
                    }
                });
            }
        }
        protected virtual void InvokeAsyncParallel<OpList>(dynamic[]? args, object? userState = null) where OpList : OptType, IList<AsyncOperation>, new() {
            object operationState = userState ?? new object();
            int op_UID = -1;
            CancellationTokenSource? cts = null;
            AsyncOperation operation = null!;
            Exception? error = null;
            bool cancelled = false;
            TResult result = default!;
            lock (locker) {
                operation = AsyncOperationManager.CreateOperation(operationState);
                // 包含 userState
                if (operations.TryGetValue(operationState, out var Op)) {
                    OpList oplist = (OpList)Op.Ops;
                    op_UID = oplist?.Count ?? 0;
                    cts = Op.CTS;
                    // 已取消的情况下，不添加任何异步操作
                    if (cts.IsCancellationRequested) {
                        error = new OperationCanceledException();
                        cancelled = true;
                        postOperationCompleted(ref operation, ref error!, ref cancelled, ref operationState!, ref result);
                        return;
                    }
                    ((OpList)Op.Ops).Add(operation);
                } else {
                    op_UID = 0;
                    cts = new CancellationTokenSource();
                    OpList oplist = new();
                    oplist.Add(operation);
                    operations.Add(operationState, (cts, oplist));
                }
            }
            ThreadPool.QueueUserWorkItem(_ => {
                try {
                    // 检查是否已取消
                    if (operationState != null) {
                        lock (locker) {
                            if (operations.TryGetValue(operationState, out var op))
                                op.CTS.Token.ThrowIfCancellationRequested();
                        }
                    }
                    CallDelegate(args, out result);
                } catch (OperationCanceledException) {
                    cancelled = true;
                } catch (Exception ex) {
                    error = ex;
                } finally {
                    AsyncOperation operation;
                    (CancellationTokenSource CTS, OptType Ops) Op;
                    lock (locker) {
                        // 并发池清空，移除状态对
                        Op = operations[operationState];
                        operation = ((OpList)Op.Ops)[op_UID];
                        ((OpList)Op.Ops).RemoveAt(op_UID);
                        if (userState is null) {  // 移除非用户定义 userState
                            operations.Remove(operationState);
                            Console.WriteLine("Clear");
                        }
                        postOperationCompleted(ref operation, ref error!, ref cancelled, ref operationState!, ref result);
                    }
                }
            });
        }
    }
}

