#include <exception>
#include <iostream>

module LanguageSupport;
using namespace std;

// 用户定义异常
void example_custom_exception() try {
	cout << "\n[Custom exception]\n";
	class MyException : public std::exception {
		std::string msg;
	public:
		MyException(const std::string& message) : msg(message) {}
		const char* what() const noexcept override {
			return msg.c_str();
		}
	};
	// 抛出异常
	throw MyException("Custom exception occurred");
}
catch (exception& e) {
	std::cout << "Caught custom exception: " << e.what() << "\n";
}

// 嵌套异常处理
void example_nested_exceptions() {
	cout << "\n[Nested exceptions]\n";
	try {
		try {
			throw std::runtime_error("Inner exception");
		}
		catch (...) {
			// 混合 runtime_error 并抛出 参数
			std::throw_with_nested(std::logic_error("Outer exception"));
		}
	}
	catch (const std::exception& e) {
		std::cout << "Caught exception: " << e.what() << "\n";
		try {
			// 抛出内嵌异常
			std::rethrow_if_nested(e);
		}
		catch (const std::exception& nested) {
			std::cout << "Nested exception: " << nested.what() << "\n";
		}
	}
}

// 异常对象指针
exception_ptr pcall(void(*f)()) noexcept {
	try {
		f();
	}
	catch (...) {
		return std::current_exception();
	}
	return nullptr;
}
void handle_exception_ptr(exception_ptr eptr) {
	try {
		if (eptr)
			rethrow_exception(eptr);
	}
	catch (bad_exception&) {
		cout << "Bad exception\n";
	}
	catch (exception& e) {
		std::cout << "Rethrown exception: " << e.what() << "\n";
	}
	catch (...) {
		cout << "Not exception type" << "\n";
	}
}
void example_exception_ptr() {
	cout << "\n[Exception pointer]\n";

	handle_exception_ptr(pcall([] {
		throw std::range_error("Range error example");
		}));

	handle_exception_ptr(pcall([] {
		throw "Throw a non-exception type";
		}));
}

// 使用 std::terminate 和设置终止处理器
void example_terminate_handler() {
	cout << "\n[Terminate handler]\n";

	static exception_ptr ptr;

	auto old_handler = std::get_terminate();
	std::set_terminate([]() {
		cout << "Terminate handler called!\n";

		// 捕获调用 terminate 是否引发异常
		auto eptr = current_exception();
		if (eptr) try {
			rethrow_exception(eptr);
		}
		catch (bad_exception&) {
			cout << "Bad exception\n";
		}
		catch (exception& e) {
			cout << e.what() << endl;
		}
		catch (...) {
			cout << "Unknown exception\n";
		}
		});

	// std::terminate();
	try {
		std::cout << "About to throw an uncaught exception...\n";
		throw "This is a string exception, not derived from std::exception";
	}
	catch (std::exception& e) {
		std::cout << "This won't be reached\n";
	}
}

void test_exception(void) {

	example_custom_exception();
	example_nested_exceptions();
	example_exception_ptr();
	// 注释掉这行以避免程序终止
	//example_terminate_handler();
}

