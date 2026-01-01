#include <iostream>
#include <any>
#include <string>
#include <vector>
#include <typeinfo>

module GeneralUtilities;
using namespace std;

// Basic usage of std::any
void example_any_basic() {
	std::cout << "\n=== Basic std::any Usage ===\n";

	std::any a1 = 42;                   // int
	std::any a2 = 3.14;                 // double
	std::any a3 = std::string("Hello"); // string
	std::any a4;                        // Empty any

	// Check if any object contains a value
	std::cout << "a1 has value: " << std::boolalpha << a1.has_value() << std::endl;
	std::cout << "a4 has value: " << a4.has_value() << std::endl;
	if (a1.has_value())
		std::cout << "a1 type name: " << a1.type().name() << std::endl;
	a1.swap(a3);
	std::cout << "a1 type name: " << a1.type().name() << std::endl;
	a4 = make_any<vector<int>>({ 1,2,3,4,5,6 });
	if (a4.has_value())
		cout << "a4 type name: " << a4.type().name() << std::endl;

	// Try to get value with type
	try {
		auto val = std::any_cast<double>(a2);
		std::cout << "a2 value: " << val << std::endl;
		// This will throw bad_any_cast
		a1.emplace<int>(std::any_cast<int>(a2));
	}
	catch (const std::bad_any_cast& e) {
		std::cout << "Exception: " << e.what() << std::endl;
	}
	a1.reset();
	std::cout << "After reset, a1 has value: " << a1.has_value() << std::endl;

	// std::any supports in_place_type_t constructor
	std::any a5(std::in_place_type<short>, 100);
	std::cout << "a5 with in_place_type: " << a5.type().name() << ", value: " << std::any_cast<short>(a5) << std::endl;

	// std::any_cast with pointer to get a nullable result instead of exception
	if (const auto* p = std::any_cast<short>(&a5); p)
		std::cout << "Value via any_cast pointer: " << *p << std::endl;

	// Trying to get a different type will return nullptr
	const auto* p_str = std::any_cast<std::string>(&a5);
	std::cout << "Pointer to string in a5 is " << (p_str ? "not null" : "nullptr") << std::endl;
}

void test_any() {
	example_any_basic();
}