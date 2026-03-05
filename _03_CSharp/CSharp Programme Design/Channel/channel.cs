using System.Runtime.CompilerServices;

public class Channel<T> : IDisposable, IAsyncDisposable {
    private readonly Queue<T>? _bufferQueue;
    private readonly int _capacity;
    private readonly Lock _lockObj = new();
    private readonly Queue<TaskCompletionSource<T>> _receiveWaiters = new Queue<TaskCompletionSource<T>>();
    private readonly Queue<TaskCompletionSource<(T Item, TaskCompletionSource<bool> SendTcs)>> _sendWaiters = new Queue<TaskCompletionSource<(T, TaskCompletionSource<bool>)>>();
    private bool _isClosed;
    private bool _disposed;
    private int _disposeLock;

    public Channel(int capacity = 0) {
        if (capacity < 0)
            throw new ArgumentOutOfRangeException(nameof(capacity), "缓冲容量不能小于 0");
        _capacity = capacity;
        _bufferQueue = capacity > 0 ? new(capacity) : null;
    }

    /// <summary>
    /// 无缓冲通道的发送方法
    /// </summary>
    private async Task<bool> SendAsyncNoBuffer(T value) {
        lock (_lockObj) {
            ObjectDisposedException.ThrowIf(_disposed, nameof(Channel<T>));
            if (_isClosed) return false;
            if (_receiveWaiters.Count > 0) {
                var recvTcs = _receiveWaiters.Dequeue();
                Task.Run(() => recvTcs?.TrySetResult(value));
                return true;
            }
            // 无等待接收方，加入发送等待队列
            var sendTcs = new TaskCompletionSource<bool>(TaskCreationOptions.RunContinuationsAsynchronously);
            var wrapperTcs = new TaskCompletionSource<(T Item, TaskCompletionSource<bool> SendTcs)>(TaskCreationOptions.RunContinuationsAsynchronously);
            wrapperTcs.TrySetResult((value, sendTcs));
            _sendWaiters.Enqueue(wrapperTcs);
        }
        // 等待接收方接收数据
        var sendTask = _sendWaiters.Peek().Task.ContinueWith(t => {
            if (t.IsCompletedSuccessfully) {
                var (_, sendTcs) = t.Result;
                return sendTcs.Task;
            }
            return Task.FromResult(false);
        }, TaskScheduler.Default).Unwrap();
        return await sendTask;
    }
    /// <summary>
    /// 有缓冲通道的发送方法
    /// </summary>
    private async Task<bool> SendAsyncWithBuffer(T value) {
        lock (_lockObj) {
            ObjectDisposedException.ThrowIf(_disposed, nameof(Channel<T>));
            if (_isClosed) return false;
            // 缓冲未满，直接入队
            if (_bufferQueue != null && _bufferQueue.Count < _capacity) {
                _bufferQueue.Enqueue(value);
                WakeUpOneReceiver();
                return true;
            }
            // 缓冲已满，加入发送等待队列
            var sendTcs = new TaskCompletionSource<bool>(TaskCreationOptions.RunContinuationsAsynchronously);
            var wrapperTcs = new TaskCompletionSource<(T Item, TaskCompletionSource<bool> SendTcs)>(TaskCreationOptions.RunContinuationsAsynchronously);
            wrapperTcs.TrySetResult((value, sendTcs));
            _sendWaiters.Enqueue(wrapperTcs);
        }

        // 等待缓冲空间
        var sendTask = _sendWaiters.Peek().Task.ContinueWith(t => {
            if (t.IsCompletedSuccessfully) {
                var (_, sendTcs) = t.Result;
                return sendTcs.Task;
            }
            return Task.FromResult(false);
        }, TaskScheduler.Default).Unwrap();
        return await sendTask;
    }
    public async Task<bool> SendAsync(T value) {
        if (_capacity == 0) {
            return await SendAsyncNoBuffer(value);
        } else {
            return await SendAsyncWithBuffer(value);
        }
    }

    public ChannelAwaiter GetAwaiter() {
        return new ChannelAwaiter(this);
    }
    public struct ChannelAwaiter : INotifyCompletion {
        private readonly Channel<T> _channel;
        private T? _result;
        private bool _isCompleted;
        public ChannelAwaiter(Channel<T> channel) {
            _channel = channel;
            _result = default;
            _isCompleted = false;
        }
        public bool IsCompleted {
            get {
                if (_isCompleted) return true;
                lock (_channel._lockObj) {
                    // 无缓冲通道：有等待发送方则可直接取值
                    if (_channel._capacity == 0)
                        return _channel._disposed || _channel._sendWaiters.Count > 0 || _channel._isClosed;
                    // 有缓冲通道：队列有数据则可取值
                    return _channel._disposed || (_channel._bufferQueue != null && _channel._bufferQueue.Count > 0) || _channel._isClosed;
                }
            }
        }
        public T GetResult() {
            if (_isCompleted && _result != null) return _result;
            lock (_channel._lockObj) {
                ObjectDisposedException.ThrowIf(_channel._disposed, nameof(Channel<T>));
                // 检查通道是否已关闭且无数据可取
                if (_channel._isClosed) {
                    if (_channel._capacity == 0 && _channel._sendWaiters.Count == 0)
                        throw new InvalidOperationException("通道已关闭，无数据可取");
                    if (_channel._capacity > 0 && _channel._bufferQueue != null && _channel._bufferQueue.Count == 0)
                        throw new InvalidOperationException("通道已关闭，无数据可取");
                }
                // 无缓冲通道：从发送等待队列中获取数据
                if (_channel._capacity == 0 && _channel._sendWaiters.Count > 0) {
                    var sendWaitTcs = _channel._sendWaiters.Dequeue();
                    var result = sendWaitTcs.Task.Result;
                    var item = result.Item;
                    var sendTcs = result.SendTcs;
                    _result = item;
                    // 通知发送方发送完成
                    Task.Run(() => sendTcs.TrySetResult(true));
                } else if (_channel._capacity > 0 && _channel._bufferQueue != null && _channel._bufferQueue.Count > 0) {
                    // 有缓冲通道：从缓冲队列中获取数据
                    _result = _channel._bufferQueue.Dequeue();
                    _channel.WakeUpOneSender();
                }
                _isCompleted = true;
                return _result!;
            }
        }

        /// <summary>
        /// 设置操作完成后的回调
        /// </summary>
        public void OnCompleted(Action continuation) {
            // 如果通道已释放、已关闭，或者有数据可取，直接执行回调
            if (_channel._disposed || _channel._isClosed ||
                (_channel._capacity == 0 && _channel._sendWaiters.Count > 0) ||
                (_channel._capacity > 0 && _channel._bufferQueue != null && _channel._bufferQueue.Count > 0)) {
                continuation();
                return;
            }
            // 创建一个 TaskCompletionSource 用于等待数据
            var recvWaitTcs = new TaskCompletionSource<T>(TaskCreationOptions.RunContinuationsAsynchronously);
            // 将其加入接收等待队列
            _channel._receiveWaiters.Enqueue(recvWaitTcs);
            // 异步等待完成后执行回调，锁外执行，无阻塞
            recvWaitTcs.Task.ContinueWith(t => {
                if (t.IsCompletedSuccessfully) continuation();
                else Task.Run(continuation);
            }, TaskScheduler.Default);
        }
    }

    private void WakeUpOneReceiver() {
        lock (_lockObj) {
            // 如果有等待接收的任务且缓冲队列有数据
            if (_receiveWaiters.Count > 0 && _bufferQueue != null && _bufferQueue.Count > 0) {
                var recvTcs = _receiveWaiters.Dequeue();
                var data = _bufferQueue.Dequeue();
                // 异步设置结果，避免阻塞
                Task.Run(() => recvTcs?.TrySetResult(data));
                WakeUpOneSender();
            }
        }
    }
    private void WakeUpOneSender() {
        lock (_lockObj) {
            // 如果有等待发送的任务且缓冲队列有空间
            if (_sendWaiters.Count > 0 && _bufferQueue != null && _bufferQueue.Count < _capacity) {
                var sendTcsWrapper = _sendWaiters.Dequeue();
                var (_, sendTcs) = sendTcsWrapper.Task.Result;
                // 异步设置结果，避免阻塞
                Task.Run(() => sendTcs.TrySetResult(true));
            }
        }
    }

    /// <summary>
    /// 唤醒所有等待的任务，通知通道已释放或关闭
    /// </summary>
    private void WakeUpAllWaiters() {
        lock (_lockObj) {
            // 唤醒所有接收等待方，设置异常
            while (_receiveWaiters.Count > 0) {
                var tcs = _receiveWaiters.Dequeue();
                Task.Run(() => tcs.TrySetException(new ObjectDisposedException(nameof(Channel<T>), "通道已释放或关闭")));
            }
            // 唤醒所有发送等待方，设置异常
            while (_sendWaiters.Count > 0) {
                var wrapperTcs = _sendWaiters.Dequeue();
                if (wrapperTcs.Task.IsCompletedSuccessfully) {
                    var result = wrapperTcs.Task.Result;
                    if (result.SendTcs != null) {
                        Task.Run(() => result.SendTcs.TrySetException(new ObjectDisposedException(nameof(Channel<T>), "通道已释放或关闭")));
                    }
                }
            }
        }
    }

    public void Close() {
        ObjectDisposedException.ThrowIf(_disposed, nameof(Channel<T>));
        lock (_lockObj) {
            if (_isClosed) return;
            _isClosed = true;
            WakeUpAllWaiters();
        }
    }
    ~Channel() => Dispose(false);
    
    public void Dispose() {
        // 使用 Interlocked.CompareExchange 确保 Dispose 只被调用一次
        if (Interlocked.CompareExchange(ref _disposeLock, 1, 0) != 0)
            return;
        Dispose(true);
        // 抑制垃圾回收器对析构函数的调用
        GC.SuppressFinalize(this);
    }
    
    /// <param name="disposing">是否是手动调用</param>
    private void Dispose(bool disposing) {
        lock (_lockObj) {
            if (_disposed) return;
            _disposed = true;
            _isClosed = true;
            if (disposing) {
                _bufferQueue?.Clear();
                WakeUpAllWaiters();
            }
        }
    }

    public ValueTask DisposeAsync() {
        // 使用 Interlocked.CompareExchange 确保 DisposeAsync 只被调用一次
        if (Interlocked.CompareExchange(ref _disposeLock, 1, 0) != 0)
            return ValueTask.CompletedTask;

        lock (_lockObj) {
            if (_disposed) return ValueTask.CompletedTask;
            _disposed = true;
            _isClosed = true;
            _bufferQueue?.Clear();
            WakeUpAllWaiters();
        }
        GC.SuppressFinalize(this);
        return ValueTask.CompletedTask;
    }
}