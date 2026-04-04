

using System.ComponentModel;

public abstract partial class EventDelegateAsync {
    private class InternalEventDelegateAsyncParallel<Action>(Action action) : BaseEventDelegate<Action, List<AsyncOperation>, AsyncCompletedEventArgs, object>, IParallelEventDelegateAsync where Action : Delegate {
        public override void InvokeAsync(dynamic[]? args, object? userState = null) {
            InvokeAsyncParallel<List<AsyncOperation>>(args, userState);
        }
        public bool Reset(object? userState) {
            if (userState is not null) {
                lock (locker)
                    return operations.Remove(userState);
            }
            return false;
        }
        protected override void CallDelegate(dynamic[]? args, out object result) {
            action?.DynamicInvoke(args);
            result = default!;
        }
        protected override void postOperationCompleted(ref AsyncOperation operation, ref Exception error, ref bool cancelled, ref object? userState, ref object result) {
            operation.PostOperationCompleted((state) => {
                var (_error, _cancelled, _userState) = ((Exception, bool, object))state!;
                OnInvokeCompleted(new(_error, _cancelled, _userState));
            }, (error, cancelled, userState));
        }
    }
    private class InternalEventDelegateAsyncParallel<Func, TResult>(Func function) : BaseEventDelegate<Action, List<AsyncOperation>, InvokeCompletedEventArgs<TResult>, TResult>, IParallelEventDelegateAsync<TResult> where Func : Delegate {
        public override void InvokeAsync(dynamic[]? args, object? userState = null) {
            InvokeAsyncParallel<List<AsyncOperation>>(args, userState);
        }
        public bool Reset(object? userState) {
            if (userState is not null) {
                lock (locker)
                    return operations.Remove(userState);
            }
            return false;
        }
        protected override void CallDelegate(dynamic[]? args, out TResult result) {
            result = (TResult)function?.DynamicInvoke(args)!;
        }
        protected override void postOperationCompleted(ref AsyncOperation operation, ref Exception error, ref bool cancelled, ref object? userState, ref TResult result) {
            operation.PostOperationCompleted((state) => {
                var (_result, _error, _cancelled, _userState) = ((TResult, Exception, bool, object))state!;
                OnInvokeCompleted(new(_result, _error, _cancelled, _userState));
            }, (result, error, cancelled, userState));
        }
    }
}