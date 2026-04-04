// 基于 APM 模式的异步委托接口
public interface IDelegateAsync {
    IAsyncResult BeginInvoke(dynamic[]? args = null, AsyncCallback? callback = null, object? state = null, CancellationToken token = default);
    void EndInvoke(IAsyncResult asyncResult);
}

public interface IDelegateAsync<T> : IDelegateAsync {
    new T EndInvoke(IAsyncResult asyncResult);
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
    private class AsyncResult(AsyncCallback? callback, object? state, CancellationToken token) : IAsyncResult {
        public Exception? exception = null;
        private readonly ManualResetEvent waitHandle = new ManualResetEvent(false);
        public CancellationToken? cancellationToken { get; init; } = token;
        public WaitHandle AsyncWaitHandle => waitHandle;
        public bool CompletedSynchronously => false;
        public bool IsCompleted { get; private set; } = false;
        public object? AsyncState => state ?? new();
        public void Complete() {
            IsCompleted = true;
            waitHandle.Set();
            callback?.Invoke(this);
        }
    }
    private class AsyncResult<T>(AsyncCallback? callback, object? state, CancellationToken token) : AsyncResult(callback, state, token) {
        public T Result = default!;
    }

    private abstract class BaseDelegateAsync<Dele, TResult>(Dele dele) : IDelegateAsync<TResult> where Dele : Delegate {
        protected abstract void BaseInvoke(dynamic[]? args, out TResult result);
        public IAsyncResult BeginInvoke(dynamic[]? args = null, AsyncCallback? callback = null, object? state = null, CancellationToken token = default) {
            var asRt = new AsyncResult<TResult>(callback, state, token);
            if (token.IsCancellationRequested) {
                asRt.exception = new OperationCanceledException(token);
                asRt.Complete();
                return asRt;
            }
            ThreadPool.QueueUserWorkItem((state) => {
                try {
                    BaseInvoke(args, out asRt.Result);
                } catch (Exception ex) {
                    asRt.exception = ex;
                } finally {
                    asRt.Complete();
                }
            });
            return asRt;
        }
        TResult IDelegateAsync<TResult>.EndInvoke(IAsyncResult asyncResult) {
            ArgumentNullException.ThrowIfNull(asyncResult);
            if (asyncResult is AsyncResult<TResult> asRes) {
                if (asRes.exception is not null)
                    throw asRes.exception;
                if (!asRes.IsCompleted)
                    asRes.AsyncWaitHandle.WaitOne();
            } else
                throw new ArgumentException("Invalid async result", nameof(asyncResult));
            return asRes.Result;
        }
        void IDelegateAsync.EndInvoke(IAsyncResult asyncResult) {
            ArgumentNullException.ThrowIfNull(asyncResult);
            if (asyncResult is AsyncResult asRes) {
                if (asRes.exception is not null)
                    throw asRes.exception;
                if (!asRes.IsCompleted)
                    asRes.AsyncWaitHandle.WaitOne();
            } else
                throw new ArgumentException("Invalid async result", nameof(asyncResult));
        }
    }

    private class InternalDelegateAsync<Action>(Action action) : BaseDelegateAsync<Action, object>(action), IDelegateAsync where Action : Delegate {
        protected override void BaseInvoke(dynamic[]? args, out object result) {
            action.DynamicInvoke(args);
            result = default!;
        }
    }
    private class InternalDelegateAsync<Func, T>(Func func) : BaseDelegateAsync<Func, T>(func), IDelegateAsync<T> where Func : Delegate {
        protected override void BaseInvoke(dynamic[]? args, out T result) {
            result = (T)func.DynamicInvoke(args)!;
        }
    }
}