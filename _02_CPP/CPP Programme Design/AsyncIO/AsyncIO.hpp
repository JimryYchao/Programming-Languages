#pragma once
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
        return Awaiter{ *this };
    }
    void complete(T value) {
        completed = true;
        result = std::move(value);
        if (callback)
            EventLoop::instance().post(callback);
    }
};