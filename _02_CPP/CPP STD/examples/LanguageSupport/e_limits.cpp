#include <limits>
#include <iostream>

module LanguageSupport;
using namespace std;

// 示例1: 基本数值类型极限值查询
void example_numeric_limits() {
	std::cout << "=== Numeric Limits ===\n";

	// 整数类型
	std::cout << "int min: " << std::numeric_limits<int>::min() << "\n";
	std::cout << "int max: " << std::numeric_limits<int>::max() << "\n";
	std::cout << "unsigned int max: " << std::numeric_limits<unsigned int>::max() << "\n";

	// 浮点类型
	std::cout << "float min: " << std::numeric_limits<float>::min() << "\n";
	std::cout << "float max: " << std::numeric_limits<float>::max() << "\n";
	std::cout << "float epsilon: " << std::numeric_limits<float>::epsilon() << "\n";
	std::cout << "float digits10: " << std::numeric_limits<float>::digits10 << "\n";

	// 特殊值
	std::cout << "float has infinity: " << std::boolalpha
		<< std::numeric_limits<float>::has_infinity << "\n";
	std::cout << "float infinity: " << std::numeric_limits<float>::infinity() << "\n";
	std::cout << "float quiet NaN: " << std::numeric_limits<float>::quiet_NaN() << "\n";
}

// 类型特性查询
void example_type_properties() {
	std::cout << "\n=== Type Properties ===\n";

	std::cout << "int is signed: " << std::numeric_limits<int>::is_signed << "\n";
	std::cout << "int is integer: " << std::numeric_limits<int>::is_integer << "\n";
	std::cout << "float is exact: " << std::numeric_limits<float>::is_exact << "\n";
	std::cout << "int is modulo: " << std::numeric_limits<int>::is_modulo << "\n";

	// 浮点类型精度比较
	std::cout << "float digits: " << std::numeric_limits<float>::digits << "\n";
	std::cout << "double digits: " << std::numeric_limits<double>::digits << "\n";
	std::cout << "long double digits: " << std::numeric_limits<long double>::digits << "\n";
}

// 安全边界检查
void example_boundary_check() {
	std::cout << "\n=== Boundary Check ===\n";

	int max_int = std::numeric_limits<int>::max();
	int min_int = std::numeric_limits<int>::min();

	std::cout << "Max int + 1 would be: " << max_int + 1 << " (undefined behavior)\n";
	std::cout << "Min int - 1 would be: " << min_int - 1 << " (undefined behavior)\n";

	// 安全边界检查示例
	auto safe_add = [](int a, int b) {
		if ((b > 0) && (a > std::numeric_limits<int>::max() - b))
			throw std::overflow_error("Integer overflow");
		if ((b < 0) && (a < std::numeric_limits<int>::min() - b)) 
			throw std::underflow_error("Integer underflow");
		return a + b;
		};
	try {
		std::cout << "Safe add: 100 + 200 = " << safe_add(100, 200) << "\n";
		std::cout << "Safe add: max_int + 1 = " << safe_add(max_int, 1) << "\n";  // 这会抛出异常
	}
	catch (const std::exception& e) {
		std::cout << "Error: " << e.what() << "\n";
	}
}

// 特殊浮点值处理
void example_special_floats() {
	std::cout << "\n=== Special Float Values ===\n";

	float inf = std::numeric_limits<float>::infinity();
	float nan = std::numeric_limits<float>::quiet_NaN();
	float normal = 1.0f;
	auto check_float = [](float f) {
		std::cout << "Value: " << f << " - ";
		if (std::isinf(f))
			std::cout << "Infinity";
		else if (std::isnan(f))
			std::cout << "NaN";
		else std::cout << "Normal";
		std::cout << "\n";
		};

	check_float(inf);
	check_float(nan);
	check_float(normal);
	check_float(inf / inf);		 // 产生 NaN
	check_float(normal / 0.0f);  // 产生 Infinity
}

void test_limits(void) {
	std::cout << std::boolalpha;
	example_numeric_limits();
	example_type_properties();
	example_boundary_check();
	example_special_floats();
}