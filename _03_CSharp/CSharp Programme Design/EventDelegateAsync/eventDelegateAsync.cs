using System.ComponentModel;

public abstract partial class EventDelegateAsync {

    private class InternalEventDelegateAsync<Action>(Action action) : BaseEventDelegate<Action, AsyncOperation, AsyncCompletedEventArgs, object>, IEventDelegateAsync where Action : Delegate {
        public bool IsBusy => Interlocked.CompareExchange(ref isBusy, 0, 0) == 1;
        public override void InvokeAsync(dynamic[]? args, object? userState = null) {
            InvokeAsyncBusy(args, userState);
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
    private class InternalEventDelegateAsync<Func, TResult>(Func function) : BaseEventDelegate<Action, AsyncOperation, InvokeCompletedEventArgs<TResult>, TResult>, IEventDelegateAsync<TResult> where Func : Delegate {
        public bool IsBusy => Interlocked.CompareExchange(ref isBusy, 0, 0) == 1;
        public override void InvokeAsync(dynamic[]? args, object? userState = null) {
            InvokeAsyncBusy(args, userState);
        }
        protected override void CallDelegate(dynamic[]? args, out TResult result) {
            result = (TResult)function.DynamicInvoke(args)!;
        }
        protected override void postOperationCompleted(ref AsyncOperation operation, ref Exception error, ref bool cancelled, ref object? userState, ref TResult result) {
            operation.PostOperationCompleted((state) => {
                var (_result, _error, _cancelled, _userState) = ((TResult, Exception, bool, object))state!;
                OnInvokeCompleted(new(_result, _error, _cancelled, _userState));
            }, (result, error, cancelled, userState));
        }
    }
}