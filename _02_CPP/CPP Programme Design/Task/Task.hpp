#pragma once
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
            return Task{ std::coroutine_handle<promise_type>::from_promise(*this) };
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
        return TaskAwaiter{ handle };
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
            return Task{ std::coroutine_handle<promise_type>::from_promise(*this) };
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
        return TaskAwaiter{ handle };
    }

    void get() {
        if (!handle.done())
            handle.resume();
        if (handle.promise().exception)
            std::rethrow_exception(handle.promise().exception);
    }
};