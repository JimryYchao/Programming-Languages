#include <coroutine>
#include <iostream>
#include <optional>
#include <memory>
#include <vector>

module LanguageSupport;
using namespace std;

// 用户定义协程
template<typename T>
struct Generator {
	struct promise_type;
	using handle_type = std::coroutine_handle<promise_type>;
	struct promise_type {
		T value_;
		std::exception_ptr eptr_;

		Generator get_return_object() {
			return Generator{ handle_type::from_promise(*this) };
		}
		std::suspend_always initial_suspend() { return {}; }
		std::suspend_always final_suspend() noexcept { return {}; }
		void unhandled_exception() { eptr_ = std::current_exception(); }

		template<std::convertible_to<T> From>
		std::suspend_always yield_value(From&& from) {
			value_ = std::forward<From>(from);
			return {};
		}
		void return_void() {}
	};

	handle_type h_;

	explicit Generator(handle_type h) : h_(h) {}
	~Generator() { if (h_) h_.destroy(); }

	explicit operator bool() {
		fill();
		return !h_.done();
	}
	T operator()() {
		fill();
		full_ = false;
		return std::move(h_.promise().value_);
	}
private:
	bool full_ = false;
	void fill() {
		if (!full_) {
			h_();
			if (h_.promise().eptr_)
				std::rethrow_exception(h_.promise().eptr_);
			full_ = true;
		}
	}
};

Generator<int> example_simple_generator() {
	co_yield 1;
	co_yield 2;
	co_yield 3;
}

void example_simple_coroutine() {
	std::cout << "=== Simple Generator Coroutine ===\n";
	auto gen = example_simple_generator();
	while (gen)
		std::cout << "Generated: " << gen() << "\n";
}

// 异步任务协程
struct AsyncTask {
	struct promise_type;
	using handle_type = std::coroutine_handle<promise_type>;

	struct promise_type {
		int value_;

		AsyncTask get_return_object() {
			return AsyncTask{ handle_type::from_promise(*this) };
		}
		std::suspend_never initial_suspend() { return {}; }
		std::suspend_always final_suspend() noexcept { return {}; }
		void unhandled_exception() { std::terminate(); }
		void return_value(int value) { value_ = value; }
	};

	handle_type h_;

	explicit AsyncTask(handle_type h) : h_(h) {}
	~AsyncTask() { if (h_) h_.destroy(); }

	int get() {
		if (!h_.done())
			h_.resume();
		return h_.promise().value_;
	}
};

AsyncTask example_async_task() {
	std::cout << "Async task started\n";
	co_return 42;
}

void example_async_coroutine() {
	std::cout << "\n=== Async Task Coroutine ===\n";
	auto task = example_async_task();
	std::cout << "Async task result: " << task.get() << "\n";
}

// 协程与 RAII 结合
struct ResourceGuard {
	std::string name;
	explicit ResourceGuard(std::string n) : name(std::move(n)) {
		std::cout << "Acquired resource: " << name << "\n";
	}
	~ResourceGuard() {
		std::cout << "Released resource: " << name << "\n";
	}
};

Generator<int> example_coroutine_with_raii() {
	ResourceGuard guard{ "Coroutine Resource" };
	co_yield 1;
	co_yield 2;
}

void example_raii_coroutine() {
	std::cout << "\n=== Coroutine with RAII ===\n";
	auto gen = example_coroutine_with_raii();
	while (gen)
		std::cout << "Value: " << gen() << "\n";
}

// 协程异常处理
Generator<int> example_coroutine_with_exception() {
	co_yield 1;
	throw std::runtime_error("Coroutine error");
	[[unreachable]] co_yield 2;
}

void example_exception_handling() {
	std::cout << "\n=== Coroutine Exception Handling ===\n";
	auto gen = example_coroutine_with_exception();
	try {
		while (gen)
			std::cout << "Value: " << gen() << "\n";
	}
	catch (const std::exception& e) {
		std::cerr << "Coroutine caught: " << e.what() << "\n";
	}

}

// 协程与智能指针
struct SharedTask {
	struct promise_type;
	using handle_type = std::coroutine_handle<promise_type>;
	struct promise_type {
		std::shared_ptr<int> value_;

		SharedTask get_return_object() {
			return SharedTask{ handle_type::from_promise(*this) };
		}
		std::suspend_always initial_suspend() { return {}; }
		std::suspend_always final_suspend() noexcept { return {}; }
		void unhandled_exception() { std::terminate(); }
		void return_value(int value) {
			value_ = std::make_shared<int>(value);
		}
	};

	handle_type h_;
	explicit SharedTask(handle_type h) : h_(h) {}
	~SharedTask() { if (h_) h_.destroy(); }

	std::shared_ptr<int> get() {
		if (!h_.done())
			h_.resume();
		return h_.promise().value_;
	}
};

SharedTask example_shared_task() {
	co_return 100;
}

void example_shared_state_coroutine() {
	std::cout << "\n=== Shared State Coroutine ===\n";
	auto task = example_shared_task();
	auto result = task.get();
	std::cout << "Shared result: " << *result << "\n";
}

void test_coroutine() {
	example_simple_coroutine();
	example_async_coroutine();
	example_raii_coroutine();
	example_exception_handling();
	example_shared_state_coroutine();
}