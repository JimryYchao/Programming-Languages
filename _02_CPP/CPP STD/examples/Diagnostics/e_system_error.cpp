#include <system_error>
#include <iostream>
#include <string>
#include <cassert>

module Diagnostics;
using namespace std;

// error_code 
void example_error_code() {
	// 创建一个表示文件未找到的错误码
	std::error_code ec = make_error_code(errc::no_such_file_or_directory);

	std::cout << "=== Example Error Code ===\n";
	std::cout << "Error message: " << ec.message() << "\n";
	std::cout << "Error code: " << ec.value() << "\n";
	std::cout << "Error category: " << ec.category().name() << "\n\n";
}

// error_condition 
void example_error_condition() {
	std::error_code ec = make_error_code(std::errc::permission_denied);
	std::error_condition ecnd = make_error_condition(std::errc::permission_denied);

	std::cout << "=== Example Error Condition ===\n";
	std::cout << "Error code matches condition: " << std::boolalpha << (ec == ecnd) << "\n";
	std::cout << "Condition message: " << ecnd.message() << "\n\n";
}

// system_error 异常
void example_system_error_exception() {
	std::cout << "=== Example System Error Exception ===\n";
	try {
		// 抛出一个 system_error 异常
		throw std::system_error(make_error_code(errc::invalid_argument), "Invalid argument error");
	}
	catch (const std::system_error& e) {
		std::cout << "Caught system_error: " << e.what() << "\n";
		std::cout << "Error code: " << e.code().value() << "\n";
		std::cout << "Error category: " << e.code().category().name() << "\n\n";
	}
}

// 自定义错误类别
class custom_error_category : public std::error_category {
public:
	const char* name() const noexcept override { return "custom_category"; }
	std::string message(int ev) const override {
		switch (ev) {
		case 1: return "Custom error: invalid parameter";
		case 2: return "Custom error: operation failed";
		default: return "Unknown custom error";
		}
	}
};
void example_custom_error() {
	std::error_code ec(1, custom_error_category{});
	std::cout << "=== Example Custom Error ===\n";
	std::cout << "Custom error message: " << ec.message() << "\n";
	std::cout << "Custom error code: " << ec.value() << "\n";
	std::cout << "Custom category name: " << ec.category().name() << "\n\n";
}

void test_system_error() {
	example_error_code();
	example_error_condition();
	example_system_error_exception();
	example_custom_error();
}