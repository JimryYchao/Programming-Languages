#include <iostream>
#include <expected>
#include <string>
#include <vector>
#include <functional>
#include <memory>
#include <stdexcept>

module GeneralUtilities;
using namespace std;

// std::expected observers
void example_expected_observers() {
	std::cout << "\n=== std::expected Observers ===\n";
	// Create an expected object that holds a value
	std::expected<int, std::string> e1 = 42;
	std::cout << "e1.has_value(): " << std::boolalpha << e1.has_value() << std::endl;
	std::cout << "e1.value(): " << e1.value() << std::endl;
	// Create an expected object that holds an error
	std::expected<int, std::string> e2 = std::unexpected("Something went wrong");
	std::cout << "e2.has_value(): " << e2.has_value() << std::endl;
	try {
		e2.value(); // This will throw std::bad_expected_access
	}
	catch (const std::bad_expected_access<std::string>& e) {
		std::cout << "Exception caught: " << e.what() << std::endl;
		std::cout << "Error message: " << e.error() << std::endl;
	}
	// value_or & error_or
	std::cout << "e1.error_or: " << e1.error_or("Unexpected value is inexist") << std::endl;
	std::cout << "e2.value_or: " << e2.value_or(10086) << std::endl;
}

// std::expected monadic operators
void expected_monadic_ops(expected<int, string>&& e) {
	// And_then 
	if (e.has_value()) {
		auto result = e.and_then([](int& val) -> std::expected<std::string, std::string> {
			return std::to_string(val * 2);
			});
		if (result.has_value())
			std::cout << "After and_then: " << result.value() << std::endl;
	}
	// Or_else 
	else {
		auto handled = e.or_else([](const std::string& err) -> std::expected<int, std::string> {
			std::cout << "Handling error: " << err << std::endl;
			return 911; // Return a default value
			});
		std::cout << "After or_else: " << handled.value() << std::endl;
	}
	// transform
	if (e.has_value()) {
		auto result = e.transform([](int& val) -> string {
			return to_string(val * val);
			});
		if (result.has_value())
			std::cout << "After transform: " << result.value() << std::endl;
	}
	// transform_error
	else {
		auto result = e.transform_error([](const std::string& err) {
			cout << "Transform error: " + err << endl;
			return 119;
			});
		std::cout << "After transform_error: " << result.error() << endl;
	}
}
void example_expected_monadic_ops() {
	cout << "\n=== expected monadic operators ===\n";
	expected_monadic_ops(10010);
	cout << endl;
	expected_monadic_ops(unexpected("Something went wrong"));
}

// std::expected modifiers
void example_expected_modifiers() {
	std::cout << "\n=== std::expected modifiers ===\n";

	auto e1 = expected<int, string>(10086);
	auto e2 = expected<int, string>(unexpected<string>("Something went wrong"));

	cout << e1.emplace(10010) << endl;
	e1.swap(e2);

	if (e1.has_value()) cout << e1.value() << endl;
	else cout << e1.error() << endl;
	if (e2.has_value()) cout << e2.value() << endl;
	else cout << e2.error() << endl;
}


// Using std::expected with functions
std::expected<int, std::string> Divide(int a, int b) {
	if (b == 0)
		return std::unexpected("Division by zero");
	return a / b;
}
std::expected<int, std::string> read_positive_number() {
	std::cout << "Enter a positive number: ";
	int num;
	std::cin >> num;
	if (std::cin.fail()) {
		std::cin.clear();
		std::cin.ignore(std::numeric_limits<std::streamsize>::max(), '\n');
		return std::unexpected("Invalid input, not a number");
	}
	if (num <= 0)
		return std::unexpected("Number must be positive");
	return num;
}
void example_expected_with_functions() {
	std::cout << "\n=== std::expected with Functions ===\n";

	// Simple division example
	auto result1 = Divide(10, 2);
	auto result2 = Divide(10, 0);

	std::cout << "10 / 2 = ";
	if (result1.has_value()) std::cout << result1.value() << std::endl;
	else std::cout << "Error: " << result1.error() << std::endl;

	std::cout << "10 / 0 = ";
	if (result2.has_value()) std::cout << result2.value() << std::endl;
	else std::cout << "Error: " << result2.error() << std::endl;

	// Chaining functions
	auto chained = Divide(20, 4).and_then([](int val) {
		return Divide(val, 2);
		});
	std::cout << "Chained division (20 / 4 / 2) = ";
	if (chained.has_value()) std::cout << chained.value() << std::endl;
	else std::cout << "Error: " << chained.error() << std::endl;

	// User input example 
	std::cout << "User input example (demonstration only):" << std::endl;
	auto user_num = read_positive_number();
	if (user_num.has_value()) std::cout << "You entered: " << user_num.value() << std::endl;
	else std::cout << "Error: " << user_num.error() << std::endl;
}



void test_expected() {
	example_expected_observers();
	example_expected_monadic_ops();
	example_expected_modifiers();
	example_expected_with_functions();
}