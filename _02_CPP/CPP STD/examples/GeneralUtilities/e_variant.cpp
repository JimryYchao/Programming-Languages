#include <iostream>
#include <variant>
#include <string>
#include <vector>
#include <functional>
#include <memory>
#include <execution>

module GeneralUtilities;
using namespace std;

// std::variant 
void example_variant() {
	std::cout << "\n=== std::variant ===\n";
	// access
	std::variant<int, double, std::string> var = forward<double>(3.14);
	std::cout << "Var size" << variant_size_v<decltype(var)> << endl;
	static_assert(is_same_v<variant_alternative_t<0, decltype(var)>, int>);  // index 0 is int
	std::cout << "Current hold index: " << var.index() << std::endl;
	std::cout << "Holds double: " << std::boolalpha << std::holds_alternative<double>(var) << std::endl;
	std::cout << "Holds int: " << std::boolalpha << std::holds_alternative<int>(var) << std::endl;

	var = std::string("Hello, variant!");
	std::cout << "Holds string after assignment: " << std::holds_alternative<std::string>(var) << std::endl;
	std::get<string>(var) = "Hello";   // ref lvalue
	std::cout << "Value as string: " << std::get<string>(var) << std::endl;  // or get<2>(var)
	*get_if<string>(&var) = "World";
	std::cout << "Value as string: " << std::get<string>(var) << std::endl;  // or get<2>(var)
	// exception
	try {
		std::cout << "Value as int: " << std::get<int>(var) << std::endl;
	}
	catch (bad_variant_access& e) {
		cout << e.what() << endl;
	}
	// valueless
	struct S {
		S(int) {};
		S(const S&) { throw "error"; }
	};
	std::variant<S, int> var1 = 0;
	try {
		var1 = S{ 10086 };  // throw
	}
	catch (...) {
		std::cout << "var1 is valueless: " << var1.valueless_by_exception()
			<< ", index == variant_npos: " << boolalpha << (var1.index() == variant_npos) << std::endl;
	}
	variant<int, monostate> var2;
	std::cout << "var2 is valueless: " << var2.valueless_by_exception()
		<< ", index: " << var2.index() << std::endl;


}

// Variant visitors and custom types
void example_variant_visitors() {
	std::cout << "\n=== std::variant Visitors and Custom Types ===\n";
	// Variant with custom types
	struct Circle { double radius; };
	struct Rectangle { double width; double height; };
	struct Triangle { double base; double height; };
	auto area_calculator = [](auto&& shape) {
		using T = std::decay_t<decltype(shape)>;
		if constexpr (std::is_same_v<T, Circle>) 
			return 3.14159 * shape.radius * shape.radius;
		else if constexpr (std::is_same_v<T, Rectangle>) 
			return shape.width * shape.height;
		else if constexpr (std::is_same_v<T, Triangle>) 
			return 0.5 * shape.base * shape.height;
		return 0.0;
		};
	std::variant<Circle, Rectangle, Triangle> shape = Circle{ 5.0 };
	std::cout << "Area of shape: " << std::visit(area_calculator, shape) << std::endl;

	// Multiple variant visit (visit multiple variants simultaneously)
	using Var_t = std::variant<int, double, std::string>;
	Var_t var1 = 10, var2 = 3.14;

	// Generic lambda that accepts different types
	auto binary_visitor = [](auto&& a, auto&& b) {
		using T1 = std::decay_t<decltype(a)>;
		using T2 = std::decay_t<decltype(b)>;

		if constexpr (std::is_arithmetic_v<T1> && std::is_arithmetic_v<T2>) {
			std::cout << "Sum: " << (a + b) << std::endl;
			std::cout << "Product: " << (a * b) << std::endl;
		}
		else {
			std::cout << "Cannot perform arithmetic on non-arithmetic types" << std::endl;
		}};
	std::visit(binary_visitor, var1, var2);
}

// variant hold reference
void example_variant_reference() {
	std::cout << "\n=== Advanced: how variant hold reference ===\n";
	// ===========================================================
	// 重要注意事项：
	// C++ 标准明确禁止 std::variant 直接持有引用类型（如 T&）、数组类型或 void 类型
	// 如需在 variant 中使用引用语义，请使用 std::reference_wrapper<T>
	// ===========================================================
	int value = 42;
	std::variant<std::reference_wrapper<int>, std::reference_wrapper<std::string>> ref_wrapper_var = std::ref(value);
	std::cout << "ref_wrapper_var references value: " << std::get<std::reference_wrapper<int>>(ref_wrapper_var).get() << std::endl;
	std::get<std::reference_wrapper<int>>(ref_wrapper_var).get() = 10086;
	std::cout << "Original value after modification through variant: " << value << std::endl;
}

void test_variant() {
	example_variant();
	example_variant_visitors();
	example_variant_reference();
}