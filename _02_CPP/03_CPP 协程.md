## CPP 协程

---

### 1. 协程概述

协程（Coroutine）是一种可以暂停和恢复执行的函数。与普通函数不同，协程可以在执行过程中暂停，保存当前状态，然后在稍后恢复执行。这种特性使得协程非常适合处理异步操作、生成器、惰性求值等场景。

C++ 协程设计为无栈协程，允许顺序代码异步执行。协程函数不使用可变参数、占位符返回类型、`return` 语句。`constexpr` / `consteval`、构造函数、析构函数不能是协程函数。

C++20 引入了标准化的协程支持，主要特性包括：
- `co_await`：等待异步操作完成
- `co_yield`：生成值并暂停
- `co_return`：返回结果并结束协程
- 编译器自动生成状态机

---
### 2. 协程的底层机制

#### 2.1. 协程的生命周期

```
┌───────────────────────────────────────┐
│              协程生命周期              │  
└───────────────────┬───────────────────┘
                    │
        ┌───────────▼──────────────┐
        │  创建协程 (调用协程函数)   │
        └───────────┬──────────────┘
                    │
        ┌───────────▼───────────────────────────────┐
        │  初始挂起 (执行到第一个 co_await/co_yield)  │
        └───────────┬───────────────────────────────┘
                    │
        ┌───────────▼──────────┐     ┌───────────────────────────┐
        │  运行中 (执行协程体)   │────>│  恢复执行 (co_await 完成)  │
        └───────────┬──────────┘     └──────┬────────────────────┘
                    │                       │
        ┌───────────▼─────────────┐         │
        │  挂起状态 (等待异步操作)  │<────────┘
        └───────────┬─────────────┘
                    │
        ┌───────────▼─────────────────┐
        │  完成/异常 (co_return/异常)  │
        └───────────┬─────────────────┘
                    │
        ┌───────────▼─────────┐
        │  销毁协程 (释放资源)  │
        └─────────────────────┘
```

>---
#### 2.2. 协程的内存布局

C++20 协程是无栈协程，其内存布局主要包括：

```
┌─────────────────────────────────────┐
│ 协程帧 (Coroutine Frame)             │
├─────────────────────────────────────┤
│ 承诺对象 (Promise)                   │
│ - 协程状态管理                       │
│ - 返回值存储                         │
│ - 异常处理                           │
├─────────────────────────────────────┤
│ 参数拷贝                             │
│ - 按值传递的参数                     │
│ - 按引用传递的参数（需特别注意生命周期）│
├─────────────────────────────────────┤
│ 局部变量                             │
│ - 在挂起点之间需要保持的变量           │
├─────────────────────────────────────┤
│ 恢复点/挂起点信息                    │
│ - 当前执行位置                       │
│ - 寄存器状态                         │
└─────────────────────────────────────┘
```

>---
#### 2.3. 编译器生成的状态机

> 原始协程函数

```cpp
Task<int> example(int x) {
    int a = co_await async_op1();
    int b = co_await async_op2();
    co_return a + b + x;
}
```

> 编译器生成的状态机（概念性表示）

```cpp
struct __example_coroutine_frame {
    enum class State { initial, after_async_op1, after_async_op2, done };
    Task<int> promise;
    State current_state = State::initial;
    int x;           // 参数
    int a, b;        // 局部变量
    static std::coroutine_handle<__example_coroutine_frame> from_promise(Task<int>::promise_type& promise){
        return reinterpret_cast<__example_coroutine_frame*>(reinterpret_cast<char*>(&promise)) 
                - offsetof(__example_coroutine_frame, promise);
    }
    void resume() {
        switch (current_state) {
            case State::initial:
                // 执行到第一个 co_await
                current_state = State::after_async_op1;
                {
                    auto awaitable = async_op1();
                    auto awaiter = awaitable.operator co_await();
                    if (!awaiter.await_ready()) {   
                        // 挂起协程
                        awaiter.await_suspend(from_promise(promise));
                        return; // 挂起
                    }
                    // 如果 awaiter 不需要挂起，直接获取结果
                    a = awaiter.await_resume();
                }
            case State::after_async_op1:
                // 执行到第二个 co_await
                current_state = State::after_async_op2;
                {
                    auto awaitable = async_op2();
                    auto awaiter = awaitable.operator co_await();
                    if (!awaiter.await_ready()) {   
                        // 挂起协程
                        awaiter.await_suspend(from_promise(promise));
                        return; // 挂起
                    }
                    // 如果 awaiter 不需要挂起，直接获取结果
                    b = awaiter.await_resume();
                }
            case State::after_async_op2:
                // 执行到 co_return
                current_state = State::done;
                {
                    int result = a + b + x;
                    promise.return_value(result);
                    // 调用 final_suspend
                    auto awaiter = promise.final_suspend();
                    if (!awaiter.await_ready()) {
                        awaiter.await_suspend(from_promise(promise));
                        return;
                    }
                }
            default:
                return;
        }
    }
};
```

> 协程函数的编译器转换

```cpp
Task<int> example(int x) {
    // 分配协程帧
    void* frame_memory = ::operator new(sizeof(__example_coroutine_frame));
    __example_coroutine_frame* frame = new (frame_memory) __example_coroutine_frame();
    // 初始化协程帧
    frame->current_state = __example_coroutine_frame::State::initial;
    frame->x = x;
    // 创建 Promise 对象
    auto& promise = frame->promise;
    // 获取返回对象
    Task<int> result = promise.get_return_object();
    // 调用 initial_suspend
    auto initial_awaiter = promise.initial_suspend();
    if (!initial_awaiter.await_ready()) {
        // 挂起协程
        initial_awaiter.await_suspend(std::coroutine_handle<__example_coroutine_frame>::from_promise(promise));
        return result;
    }
    // 如果 initial_suspend 不挂起，直接开始执行
    frame->resume();
    return result;
}
```

---
### 3. 核心概念

#### 3.1. Promise 类型

Promise 类型是协程与调用者之间的契约，定义了协程的行为：

```cpp
struct Promise {
    // 返回包装 coroutine_handle 的对象
    ReturnType get_return_object();     
    // 协程开始的挂起，通常返回 std::suspend_always()/std::suspend_never
    Awaitable initial_suspend();        
    // 协程结束的挂起
    Awaitable final_suspend() noexcept;
    // 处理协程中未捕获的异常
    void unhandled_exception();
    // void or value
    void return_void();            // co_return;
    void return_value(T value);    // co_return value;
// 可选 ======================
    // co_yield value;
    Awaitable yield_value(T value);      
    // co_await 自定义转换行为;
    Awaitable await_transform(T value);  
    // 当协程帧分配失败时返回的替代 get_return_object()
    ReturnType get_return_object_on_allocation_failure();   
    // 自定义分配协程帧
    void* operator new(size_t size);
    void operator delete(void* ptr, size_t size)
    // 自定义分配协程帧数组
    void* operator new[](size_t size);
    void operator delete[](void* ptr, size_t size);
};
```

>---
#### 3.2. Awaitable 与 Awaiter

`co_await` 表达式需要一个 Awaitable 对象，Awaitable 必须提供 `operator co_await` 或本身就是 Awaiter：

```cpp
struct Awaiter {
    // 是否立即完成（不挂起）
    bool await_ready() const noexcept;
    // 挂起时的操作，返回 void 或 awaitable
    void await_suspend(std::coroutine_handle<> handle);   // 协程挂起返回调用方
    bool await_suspend(std::coroutine_handle<> handle);   // true 挂起；false 继续执行
    std::coroutine_handle<> await_suspend(std::coroutine_handle<> handle);  // 对称转移，恢复目标协程
    // 作为 co_await 表达式的结果
    void await_resume();   // 无返回值
    T await_resume();
};

struct Awaitable {
    // 提供 operator co_await，返回 Awaiter 对象
    auto operator co_await() {
        return Awaiter{/* ... */};
    }
};
```

>---
#### 3.3. 标准库 Awaitable

C++20 标准库提供了几个基本的 Awaitable：

| 类型                      | 描述     | 用途              |
| ------------------------- | -------- | ----------------- |
| `std::suspend_always`     | 总是挂起 | 初始/最终挂起策略 |
| `std::suspend_never`      | 从不挂起 | 立即执行的协程    |
| `std::coroutine_handle<>` | 协程句柄 | 控制协程的执行    |

```cpp
// std::suspend_always 实现
struct suspend_always {
    constexpr bool await_ready() const noexcept { return false; }
    constexpr void await_suspend(coroutine_handle<>) const noexcept {}
    constexpr void await_resume() const noexcept {}
};
// std::suspend_never 实现
struct suspend_never {
    constexpr bool await_ready() const noexcept { return true; }
    constexpr void await_suspend(coroutine_handle<>) const noexcept {}
    constexpr void await_resume() const noexcept {}
};
```

---
### 4. 具体实现

#### 4.1. 异步任务实现

```cpp
#include <coroutine>
#include <optional>
#include <exception>
#include <iostream>

// 简单的异步任务实现
template<typename T = void>
class Task {
public:
    struct promise_type {
        std::optional<T> value;
        std::exception_ptr exception;
        std::coroutine_handle<> continuation;
        Task get_return_object() {
            return Task{std::coroutine_handle<promise_type>::from_promise(*this)};
        }
        std::suspend_always initial_suspend() { return {}; }
        
        // 最终挂起，等待结果被获取
        auto final_suspend() noexcept {
            struct FinalAwaiter {
                bool await_ready() noexcept { return false; }
                void await_suspend(std::coroutine_handle<promise_type> h) noexcept {
                    if (h.promise().continuation) 
                        h.promise().continuation.resume();
                }
                void await_resume() noexcept {}
            };
            return FinalAwaiter{};
        }
        void return_value(T v) {
            value = std::move(v);
        }
        void unhandled_exception() {
            exception = std::current_exception();
        }
    };
    
    using Handle = std::coroutine_handle<promise_type>;
    Handle handle;
    explicit Task(Handle h) : handle(h) {}
    Task(Task&& other) noexcept : handle(other.handle) {
        other.handle = nullptr;
    }
    Task& operator=(Task&& other) noexcept {
        if (this != &other) {
            if (handle) handle.destroy();
            handle = other.handle;
            other.handle = nullptr;
        }
        return *this;
    }
    ~Task() {
        if (handle) handle.destroy();
    }
    
    // 使 Task 可等待
    auto operator co_await() {
        struct TaskAwaiter {
            Handle handle;
            bool await_ready() const {
                return handle.done();
            }
            void await_suspend(std::coroutine_handle<> h) {
                handle.promise().continuation = h;
                handle.resume();
            }
            T await_resume() {
                if (handle.promise().exception) 
                    std::rethrow_exception(handle.promise().exception);
                return std::move(*handle.promise().value);
            }
        };
        return TaskAwaiter{handle};
    }
    
    T get() {
        if (!handle.done()) 
            handle.resume();
        if (handle.promise().exception) 
            std::rethrow_exception(handle.promise().exception);
        return std::move(*handle.promise().value);
    }
    bool done() const { return handle.done(); }
};

// void 特化
template<>
class Task<void> {
public:
    struct promise_type {
        std::exception_ptr exception;
        std::coroutine_handle<> continuation;
        Task get_return_object() {
            return Task{std::coroutine_handle<promise_type>::from_promise(*this)};
        }
        std::suspend_always initial_suspend() { return {}; }
        auto final_suspend() noexcept {
            struct FinalAwaiter {
                bool await_ready() noexcept { return false; }
                void await_suspend(std::coroutine_handle<promise_type> h) noexcept {
                    if (h.promise().continuation) 
                        h.promise().continuation.resume();
                }
                void await_resume() noexcept {}
            };
            return FinalAwaiter{};
        }
        void return_void() {}  // void
        void unhandled_exception() {
            exception = std::current_exception();
        }
    };
    using Handle = std::coroutine_handle<promise_type>;
    Handle handle;
    explicit Task(Handle h) : handle(h) {}
    Task(Task&& other) noexcept : handle(other.handle) {
        other.handle = nullptr;
    }
    Task& operator=(Task&& other) noexcept {
        if (this != &other) {
            if (handle) handle.destroy();
            handle = other.handle;
            other.handle = nullptr;
        }
        return *this;
    }
    ~Task() {
        if (handle) handle.destroy();
    }
    
    auto operator co_await() {
        struct TaskAwaiter {
            Handle handle;
            bool await_ready() const {
                return handle.done();
            }
            void await_suspend(std::coroutine_handle<> h) {
                handle.promise().continuation = h;
                handle.resume();
            }
            void await_resume() {
                if (handle.promise().exception)
                    std::rethrow_exception(handle.promise().exception);
            }
        };
        return TaskAwaiter{handle};
    }
    
    void get() {
        if (!handle.done()) 
            handle.resume();
        if (handle.promise().exception)
            std::rethrow_exception(handle.promise().exception);
    }
};
```
```cpp
// 使用示例
Task<int> async_add(int a, int b) {
    std::cout << "开始计算 " << a << " + " << b << std::endl;
    co_return a + b;
}

Task<int> async_computation() {
    int x = co_await async_add(1, 2);
    int y = co_await async_add(x, 3);
    co_return y * 2;
}

int main() {
    auto task = async_computation();
    int result = task.get();
    std::cout << "结果: " << result << std::endl;  // 输出: 12
    return 0;
}
```

>---
#### 4.2. 生成器实现

```cpp
#include <coroutine>
#include <optional>
#include <iterator>

template<typename T>
class Generator {
public:
    struct promise_type {
        std::optional<T> current_value;
        Generator get_return_object() {
            return Generator{std::coroutine_handle<promise_type>::from_promise(*this)};
        }
        std::suspend_always initial_suspend() { return {}; }
        std::suspend_always final_suspend() noexcept { return {}; }
        std::suspend_always yield_value(T value) {
            current_value = std::move(value);
            return {};
        }
        void return_void() {}
        void unhandled_exception() { std::terminate(); }
    };
    using Handle = std::coroutine_handle<promise_type>;
    Handle handle;
    explicit Generator(Handle h) : handle(h) {}
    Generator(Generator&& other) noexcept : handle(other.handle) {
        other.handle = nullptr;
    }
    ~Generator() {
        if (handle) handle.destroy();
    }
    
    // 迭代器支持
    struct Iterator {
        Handle handle;
        Iterator& operator++() {
            handle.resume();
            if (handle.done()) handle = nullptr;
            return *this;
        }
        T operator*() const {
            return *handle.promise().current_value;
        }
        bool operator!=(const Iterator& other) const {
            return handle != other.handle;
        }
    };
    Iterator begin() {
        if (handle && !handle.done()) {
            handle.resume();
        }
        return Iterator{handle.done() ? nullptr : handle};
    }
    Iterator end() {
        return Iterator{nullptr};
    }
    bool done() const { return !handle || handle.done(); }
    T current() const {
        return *handle.promise().current_value;
    }
    void next() {
        if (handle && !handle.done()) {
            handle.resume();
        }
    }
};
```
```cpp
#include <iostream>
// 使用示例：斐波那契数列生成器
Generator<long long> fibonacci(int n) {
    long long a = 0, b = 1;
    for (int i = 0; i < n; ++i) {
        co_yield a;
        auto next = a + b;
        a = b;
        b = next;
    }
}
// 无限序列生成器
Generator<long long> infinite_fibonacci() {
    long long a = 0, b = 1;
    while (true) {
        co_yield a;
        auto next = a + b;
        a = b;
        b = next;
    }
}
int main() {
    std::cout << "前10个斐波那契数:" << std::endl;
    for (auto num : fibonacci(10)) {
        std::cout << num << " ";
    }
    std::cout << std::endl;
    std::cout << "\n无限序列的前15个:" << std::endl;
    auto gen = infinite_fibonacci();
    for (int i = 0; i < 15; ++i) {
        std::cout << gen.current() << " ";
        gen.next();
    }
    std::cout << std::endl;
    return 0;
}
```

>---
#### 4.3. 异步 I/O 实现

```cpp
#include <coroutine>
#include <functional>
#include <vector>
#include <queue>
#include <iostream>

// 简化的异步 I/O 系统
class EventLoop {
public:
    using Callback = std::function<void()>;
    void post(Callback cb) {
        tasks.push(std::move(cb));
    }
    void run() {
        while (!tasks.empty()) {
            auto task = std::move(tasks.front());
            tasks.pop();
            task();
        }
    }
    static EventLoop& instance() {
        static EventLoop loop;
        return loop;
    }
private:
    std::queue<Callback> tasks;
};

// 异步 I/O 操作
template<typename T>
class AsyncIO {
private:
    bool completed = false;
    T result;
    std::function<void()> callback;
public:
    struct Awaiter {
        AsyncIO& io;
        bool await_ready() const {
            return io.completed;
        }
        void await_suspend(std::coroutine_handle<> h) {
            io.callback = [h]() mutable { h.resume(); };
        }
        T await_resume() const {
            return io.result;
        }
    };
    auto operator co_await() {
        return Awaiter{*this};
    }
    void complete(T value) {
        completed = true;
        result = std::move(value);
        if (callback) 
            EventLoop::instance().post(callback);
    }
};
```
```cpp
// 模拟异步读取
AsyncIO<std::string> async_read_file(const std::string& filename) {
    AsyncIO<std::string> op;
    // 模拟异步操作
    EventLoop::instance().post([&op, filename]() {
        // 模拟文件读取延迟
        std::cout << "读取文件: " << filename << std::endl;
        op.complete("文件内容: " + filename);
    });
    return op;
}

// 模拟异步写入
AsyncIO<bool> async_write_file(const std::string& filename, const std::string& content) {
    AsyncIO<bool> op;
    EventLoop::instance().post([&op, filename, content]() {
        std::cout << "写入文件: " << filename << std::endl;
        std::cout << "内容: " << content << std::endl;
        op.complete(true);
    });
    return op;
}

// 使用协程进行异步文件操作
Task<std::string> process_files() {
    auto content1 = co_await async_read_file("input1.txt");
    auto content2 = co_await async_read_file("input2.txt");
    std::string combined = content1 + " + " + content2;
    co_await async_write_file("output.txt", combined);
    co_return combined;
}

int main() {
    auto task = process_files();
    // 启动协程
    if (!task.done()) {
        task.handle.resume();
    }
    // 运行事件循环
    EventLoop::instance().run();
    std::cout << "最终结果: " << task.get() << std::endl;
    return 0;
}
```

>---
#### 4.4. 数据处理管道

```cpp
// 使用协程构建数据处理管道
#include <coroutine>
#include <vector>
#include <functional>
// 异步数据流
template<typename T>
class AsyncStream {
public:
    struct promise_type {
        T current_value;
        bool has_value = false;
        
        AsyncStream get_return_object() {
            return AsyncStream{std::coroutine_handle<promise_type>::from_promise(*this)};
        }
        
        std::suspend_always initial_suspend() { return {}; }
        std::suspend_always final_suspend() noexcept { return {}; }
        
        std::suspend_always yield_value(T value) {
            current_value = std::move(value);
            has_value = true;
            return {};
        }
        
        void return_void() {}
        void unhandled_exception() { std::terminate(); }
    };
    
    using Handle = std::coroutine_handle<promise_type>;
    Handle handle;
    
    explicit AsyncStream(Handle h) : handle(h) {}
    
    AsyncStream(AsyncStream&& other) noexcept : handle(other.handle) {
        other.handle = nullptr;
    }
    ~AsyncStream() {
        if (handle) handle.destroy();
    }
    struct Iterator {
        Handle handle;
        Iterator& operator++() {
            handle.promise().has_value = false;
            handle.resume();
            if (handle.done()) handle = nullptr;
            return *this;
        }
        T operator*() const {
            return handle.promise().current_value;
        }
        bool operator!=(const Iterator& other) const {
            return handle != other.handle;
        }
    };
    Iterator begin() {
        if (handle && !handle.done()) {
            handle.resume();
        }
        return Iterator{handle.done() ? nullptr : handle};
    }
    Iterator end() {
        return Iterator{nullptr};
    }
};
```
```cpp
// 数据管道操作
AsyncStream<int> read_sensor_data() {
    for (int i = 0; i < 100; ++i) {
        co_await async_sleep(std::chrono::milliseconds(100));
        co_yield read_sensor();  // 模拟读取传感器数据
    }
}
AsyncStream<double> process_data(AsyncStream<int> input) {
    for (auto value : input) {
        // 数据过滤和转换
        if (value > threshold) {
            co_yield transform(value);
        }
    }
}
AsyncStream<ProcessedData> aggregate_data(AsyncStream<double> input) {
    std::vector<double> batch;
    
    for (auto value : input) {
        batch.push_back(value);
        
        if (batch.size() >= 10) {
            co_await async_save_to_database(batch);
            co_yield calculate_statistics(batch);
            batch.clear();
        }
    }
}
// 使用管道
void run_pipeline() {
    auto sensor_stream = read_sensor_data();
    auto processed_stream = process_data(std::move(sensor_stream));
    // auto aggregated_stream = aggregate_data(std::move(processed_stream));
    
    for (auto result : aggregated_stream) {
        std::cout << "统计结果: " << result << std::endl;
    }
}
```

---
### 5. 参考资源

- [C++20 Coroutines Proposal](http://www.open-std.org/jtc1/sc22/wg21/docs/papers/2018/n4775.pdf)
- [cppreference: Coroutines](https://en.cppreference.com/w/cpp/language/coroutines)
- [Lewis Baker's Blog on Coroutines](https://lewissbaker.github.io/)
- [C++ Coroutines: Understanding the Compiler Transform](https://lewissbaker.github.io/2022/08/27/understanding-the-compiler-transform)

---