
#include <iostream>
#include <stacktrace>
#include <stdexcept>
#include <string>

module Diagnostics;
using namespace std;

// 基本栈跟踪捕获与输出
void example_basic_stacktrace() {
	// 获取当前栈跟踪
	auto st = std::stacktrace::current(0, 4);

	std::cout << "=== Basic Stack Trace ===\n";
	std::cout << st << "\n\n";
}

// 带帧跳过的栈跟踪
void helper_function(int depth) { // 辅助函数: 用于创建调用链
	if (depth > 0) {
		helper_function(depth - 1);  // 递归 depth +1 
	}
	else {
		// 获取当前栈跟踪并跳过2个帧
		std::stacktrace st = std::stacktrace::current(2,3);
		std::cout << st << "\n\n";
	}
}
void example_stacktrace_skipping() {
	std::cout << "=== Stack Trace with Frame Skipping ===\n";
	helper_function(3);
}

// 异常中的栈跟踪
void example_exception_stacktrace() {
	try {
		throw std::runtime_error("Sample exception with stack trace");
	}
	catch (const std::exception& e) {
		std::cout << "=== Exception with Stack Trace ===\n";
		std::cout << "Exception: " << e.what() << "\n";
		std::cout << "Stack trace:\n" << std::stacktrace::current(1) << "\n\n";
	}
}

// 栈跟踪条目详细信息
void example_stacktrace_entries() {
	std::stacktrace st = std::stacktrace::current();

	std::cout << "=== Stack Trace Entry Details ===\n";
	std::cout << "Stack trace size: " << st.size() << " frames\n\n";

	for (size_t i = 0; i < st.size() && i < 4; ++i) {
		const std::stacktrace_entry& entry = st[i];
		std::cout << "Frame " << i << ":\n";
		std::cout << "  Function: " << entry.description() << "\n";
		std::cout << "  File:     " << entry.source_file() << "\n";
		std::cout << "  Line:     " << entry.source_line() << "\n\n";
	}
}

void test_stacktrace(void) {
	example_basic_stacktrace();
	example_stacktrace_skipping();
	example_exception_stacktrace();
	example_stacktrace_entries();
}