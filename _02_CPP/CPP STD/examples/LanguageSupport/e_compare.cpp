#include <compare>
#include <iostream>
#include <algorithm>
#include <vector>
#include <cassert>
#include <ranges>

module LanguageSupport;
using namespace std;

// 演示三路比较
static void example_compare_three_way()
{
	std::cout << "\n[Compare three way]\n";

	int a = 5, b = 3;
	double x = 3.14, y = 2.71;

	// 整数比较（返回 strong_ordering）
	strong_ordering int_res = a <=> b;
	std::cout << "5 <=> 3: ";
	if (int_res > 0)       std::cout << "greater\n";
	else if (int_res < 0)  std::cout << "less\n";
	else                   std::cout << "equal\n";

	// 浮点数比较（返回 partial_ordering）
	partial_ordering float_res = x <=> y;
	std::cout << "3.14 <=> 2.71: ";
	if (float_res > 0)      std::cout << "greater\n";
	else if (float_res < 0) std::cout << "less\n";
	else                    std::cout << "equal\n";

	// 处理 NaN
	double nan = std::numeric_limits<double>::quiet_NaN();
	auto nan_res = x <=> nan;
	std::cout << "3.14 <=> NaN: ";
	if (nan_res == std::partial_ordering::unordered)
		std::cout << "unordered (NaN)\n";

	// 自定义类型的默认三路比较
	struct Point {
		int x, y;
		auto operator<=>(const Point&) const = default;
	};
	Point p1{ 1, 2 }, p2{ 1, 3 };
	std::cout << "p1 < p2: " << (p1 < p2) << "\n";   // true 
	std::cout << "p1 == p2: " << (p1 == p2) << "\n"; // false

	// 手动实现三路比较
	struct Book {
		std::string title;
		int year;
		auto operator<=>(const Book& other) const {
			return  title <=> other.title;	// 只比较 title
		}
	};
	Book b1{ "C++ Primer", 2020 }, b2{ "The C++ Programming Language", 2013 };
	std::cout << "b1 < b2: " << (b1 < b2) << "\n";  // true ("C++" < "The")
}

// 演示标准比较函数
static void example_compare_functions() {
	cout << "\n[base compare functions]\n";
	// 基础类型比较
	double a = 3.14, b = 2.71;
	assert(std::strong_order(a, b) == std::strong_ordering::greater);

	// 处理 NaN
	double nan = std::numeric_limits<double>::quiet_NaN();
	auto nan_res = std::partial_order(a, nan);
	assert(nan_res == std::partial_ordering::unordered);

	// 自定义类型回退比较，当类型没有 <=> 时，使用 < 和 == 生成比较结果
	struct Data {
		double value;
		bool operator ==(const Data& other) const {
			return this->value == other.value;
		}
		bool operator < (const Data& other) const {
			return this->value < other.value;
		}
		Data(double l) {
			value = l;
		}
	} l1 = 3.14, l2 = 9.18;
	assert(std::compare_strong_order_fallback(l1, l2) == strong_ordering::less);
}

// 演示比较函数算法应用
static struct Person {
	string name;
	int age;
	auto operator <=> (const Person& other) const {
		if (auto cmp = this->name <=> other.name; cmp != 0)
			return cmp;
		else return this->age <=> other.age;
	};
	friend ostream& operator <<(std::ostream& os, const Person& p) {
		os << "{\"" + p.name + "\", " << p.age << "}";
		return os;
	}
};
static void example_compare_Algorithms()
{
	std::cout << "\n[compare with Algorithms]\n";

	// 使用三路比较排序
	std::vector<Person> p = { {"Tom", 12}, {"Anne",15}, {"Bob", 16}, {"Tom", 17} };
	std::ranges::sort(p, [](auto a, auto b) {
		return (a <=> b) < 0;  // 正序
		});
	std::cout << "Sorted by compare_three_way : ";
	for (auto n : p) std::cout << n << " ";
	std::cout << "\n";

	// 使用自定义比较函数
	std::vector<int> iv = { 5, 3, 1, 4, 2 };
	auto cmp = [](int a, int b)-> int {
		return a >= b ? 1 : 0;  // 倒序
		};
	std::sort(iv.begin(), iv.end(), cmp);
	std::cout << "Sorted by custom cmp: ";
	for (int n : iv) std::cout << n << " ";  // 5 4 3 2 1
	std::cout << "\n";

	// 标准比较函数
	std::vector<double> dv = { 5.18, 3.14, 1.123, 4.880,  numeric_limits<double>::quiet_NaN(), 2.354, numeric_limits<double>::infinity() };
	ranges::sort(dv, [](auto a, auto  b) {
		return std::partial_order(a, b) == partial_ordering::unordered || (a <=> b) < 0;
		});
	std::cout << "Sorted by compareFunc: ";
	for (auto n : dv) std::cout << n << " ";
	std::cout << "\n";
}

// 验证三路比较
template<class T, class U = T>
struct Compare {
	static_assert(three_way_comparable_with<T, U>);
	using cmp = compare_three_way_result<T, U>::type;

	constexpr static int With(const T&& t, const U&& u) {
		if constexpr (is_same_v<cmp, partial_ordering>) {
			if ((t <=> u) == partial_ordering::unordered)
				return 1;
		}
		else {
			auto order = t <=> u;
			return order > 0 ? 1 : (order < 0 ? -1 : 0);
		}
	}
};
static void example_verify_compare()
{
	std::cout << "\n[verify compare_three_way]\n";

	// 验证三路比较
	cout << "int <=> float : " << three_way_comparable_with<int, float> << "\n";
	cout << "void* <=> int*: " << three_way_comparable_with<void*, int*> << "\n";
	cout << "const char* <=> char*: " << three_way_comparable_with<const char*, char*> << "\n";
	cout << "const char* <=> string: " << three_way_comparable_with<const char*, string> << "\n";
	cout << "char* <=> string: " << three_way_comparable_with<char*, string> << "\n";
	cout << "int* <=> int[]: " << three_way_comparable_with<int*, int[]> << "\n";

	// 用户定义类型 
	// 若满足概念 three_way_comparable_with<T, U> 要求，需定义 T(U), operator==, operator<=>, operator T <=> U
	struct Duration {
		int seconds;
		constexpr Duration(int s) : seconds{ s } {
		}   // T(U)
		constexpr bool operator==(const Duration& other) const {
			return this->seconds == other.seconds;
		};  // operator== 
		constexpr auto operator<=>(const Duration& other) const {
			return seconds <=> other.seconds;
		};  // operator<=>
		constexpr auto operator<=>(int n) const {
			return seconds <=> n;
		}   // operator T <=> U
	};
	static_assert(3 <=> Duration{ 3 } == strong_ordering::equivalent);

	// 测试 Compare
	cout << "Compare 3 with Duration{4} = " << Compare<int, Duration>::With(3, Duration{ 4 }) << "\n";
	cout << "Compare 3 with Duration{3} = " << Compare<int, Duration>::With(3, Duration{ 3 }) << "\n";
	cout << "Compare 3 with 3.0 = " << Compare<int, double>::With((int)3, (double)3.0) << "\n";
	cout << "Compare 3.0 with nan = " << Compare<float, float>::With(3.0, std::numeric_limits<float>::quiet_NaN()) << "\n";
}

void test_compare(void) {
	std::cout << std::boolalpha;
	example_compare_three_way();
	example_compare_functions();
	example_compare_Algorithms();
	example_verify_compare();
}