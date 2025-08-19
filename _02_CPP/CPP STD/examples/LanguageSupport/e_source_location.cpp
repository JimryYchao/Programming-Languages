#include <source_location>
#include <iostream>
#include <string>

module LanguageSupport;
using namespace std;

// 在日志系统中使用 source_location
void example_logging_system(const std::string& log_message,
	const std::source_location& loc = std::source_location::current())
{
	std::cout << "[LOG][" << loc.file_name() << ":" << loc.line() << "] "
		<< log_message << "\n";
}


// 在异常中使用 source_location
class LocatedException : public std::exception {
	std::string msg_;
	std::source_location loc_;
public:
	LocatedException(
		const std::string& msg,
		const std::source_location& loc = std::source_location::current())
		: msg_(msg), loc_(loc) {
	}

	const char* what() const noexcept override {
		static std::string full_msg;
		full_msg = msg_ + "\n  at " + loc_.file_name() + ":" + std::to_string(loc_.line());
		return full_msg.c_str();
	}
};
void example_exception_usage()
{
	try {
		throw LocatedException("Something went wrong");
	}
	catch (const std::exception& e) {
		std::cout << "Caught exception: " << e.what() << "\n";
	}
}

void test_source_location(void)
{
	std::cout << "\n=== Logging system ===\n";
	example_logging_system("This is a log");

	std::cout << "\n=== Exception with location ===\n";
	example_exception_usage();
}