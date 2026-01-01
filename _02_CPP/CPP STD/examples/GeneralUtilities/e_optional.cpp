#include <iostream>
#include <optional>
#include <string>
#include <vector>
#include <algorithm>
#include <functional>
#include <memory>
#include <stdexcept>

module GeneralUtilities;
using namespace std;

// std::optional
void example_optional() {
	std::cout << "\n=== std::optional ===\n";

	// Observers
	std::optional<int> opt1 = 42;
	if (opt1.has_value()) {
		std::cout << "opt1.value(): " << opt1.value() << std::endl;
		std::cout << "*opt1: " << *opt1 << std::endl;
	}
	std::optional<std::string> opt2/*= std::nullopt */;  // Create an empty optional
	try {
		opt2.value(); // has no value
	}
	catch (const std::bad_optional_access& e) {
		std::cout << "Exception caught: " << e.what() << std::endl;
	}

	std::cout << "opt1.value_or(0): " << opt1.value_or(0) << std::endl;
	std::cout << "opt2.value_or(\"default\"): " << opt2.value_or("default") << std::endl;

	// Modifiers
	opt1.reset();
	std::cout << "After reset, opt1.has_value(): " << opt1.has_value() << std::endl;
	opt1.emplace(100);
	std::cout << "After emplace, opt1.value(): " << opt1.value() << std::endl;
	std::optional<int> opt3 = 200;
	opt1.swap(opt3);
	std::cout << "After swap, opt1.value(): " << opt1.value() << ", opt3.value(): " << opt3.value() << std::endl;

	// Monadic Ops
	struct S {
		int v;
	} s{ 10086 };

	auto opt = make_optional<S*, 10010>(&s);
	auto and_then_result = opt.and_then([](S* s) -> std::optional<std::string> {
		return nullopt;
		});
	std::cout << "After and_then: " << and_then_result.value_or("empty") << endl;

	std::optional<S> empty_opt;
	auto or_else_result = empty_opt.or_else([]() -> optional<S> {
		return S{ 999 };
		});
	std::cout << "After or_else: " << or_else_result.value_or(S{ 0 }).v << std::endl;

	auto transform_result = opt.transform([](S* s) -> int {
		return s->v * 1000;
		});
	std::cout << "After transform: " << transform_result.value_or(0) << std::endl;

	// make_optional with in-place construction
	auto opt_string = std::make_optional<string>(5, 'a'); // Creates "aaaaa"
	std::cout << "make_optional string with in-place: " << opt_string.value() << std::endl;

	auto opt_vec = std::make_optional<vector<int>>(5, 1);  // {1,1,1,1,1}
	std::cout << "make_optional vector with in-place: {";
	for (size_t i = 0; i < opt_vec->size(); i++)
		cout << opt_vec->operator[](i) << ",";  // (*opt_vec)[i]
	cout << "\b}" << endl;

	// optional of reference
	int value = 42;
	std::optional<int*> opt_ptr = &value;   // 指针
	std::optional<std::reference_wrapper<int>> opt_ref = std::ref(value);  // 引用包装
	cout << (*opt_ptr.value() = 10086, value) << endl;
	cout << (opt_ref->get() = 10010, value) << endl;
}

void test_optional() {
	example_optional();
}