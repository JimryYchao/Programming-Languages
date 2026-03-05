#pragma once
#include <coroutine>
#include <optional>
#include <iterator>
#include <iostream>

template<typename T>
class Generator {
public:
    struct promise_type {
        std::optional<T> current_value;
        Generator get_return_object() {
            return Generator{ std::coroutine_handle<promise_type>::from_promise(*this) };
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
        return Iterator{ handle.done() ? nullptr : handle };
    }
    Iterator end() {
        return Iterator{ nullptr };
    }
    bool done() const { return !handle || handle.done(); }
    T current() const {
        if (!handle || handle.done() || !handle.promise().current_value) {
            throw std::runtime_error("Generator is not ready");
        }
        return *handle.promise().current_value;
    }
    void next() {
        if (handle && !handle.done()) {
            handle.resume();
        }
    }
};