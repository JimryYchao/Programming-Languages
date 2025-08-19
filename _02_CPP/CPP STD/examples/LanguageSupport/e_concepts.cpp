#include <concepts>
#include <iostream>
#include <vector>
#include <string>
#include <type_traits>
#include <ranges>
#include <algorithm>

module LanguageSupport;
using namespace std;

// 示例1: 基本概念使用
template <class T>
concept Numeric = integral<T> || floating_point<T>;
struct S_default {};
struct S_delete {
	S_delete(const S_delete&) = delete;  // 复制构造
	S_delete(const S_delete&&) = delete; // 移动构造
};
void example_basic_concepts() {
	std::cout << "=== Basic Concepts ===\n";

	static_assert(default_initializable<S_default> && move_constructible<S_default> && copy_constructible<S_default>);
	static_assert(!default_initializable<S_delete>);
	static_assert(!move_constructible<S_delete>);
	static_assert(!copy_constructible<S_delete>);

	// type_traits
	std::cout << "int is integral: " << std::boolalpha
		<< std::integral<int> << "\n";
	std::cout << "double is floating_point: "
		<< std::floating_point<double> << "\n";
	std::cout << "string is derived_from<string>: "
		<< std::derived_from<std::string, std::string> << "\n";
}

// 自定义概念
template<typename T>
concept Addable = requires(T a, T b) {
	{ a + b } -> std::same_as<T>;
};
template<Addable T>
T add(T a, T b) {
	return a + b;
}
void example_custom_concept() {
	std::cout << "\n=== Custom Addable Concept ===\n";

	std::cout << "Sum of integers: " << add(5, 3) << "\n";
	std::cout << "Sum of doubles: " << add(3.14, 2.71) << "\n";
	std::cout << "Sum of strings: " << add(std::string("Hello"), std::string("World"));
}

// 概念与迭代器
template<typename T>
concept RandomAccessIterator = requires(T it) {
	requires std::input_iterator<T>;
{ it[0] } -> std::same_as<typename std::iterator_traits<T>::reference>;
{ it + 1 } -> std::same_as<T>;
};
template<RandomAccessIterator Iter>
void sort_and_print(Iter begin, Iter end) {
	std::sort(begin, end);
	std::cout << "Sorted range: ";
	for (auto it = begin; it != end; ++it) 
		std::cout << *it << " ";
	std::cout << "\n";
}
void example_iterator_concept() {
	std::cout << "\n=== Iterator Concept ===\n";
	std::vector<int> nums = { 5, 2, 8, 1, 9 };
	sort_and_print(nums.begin(), nums.end());
}

// 概念与 Ranges
template<typename R>
concept PrintableRange = std::ranges::range<R> &&
	requires(std::ranges::range_value_t<R> v) {
	std::cout << v;
};
void print_range(const PrintableRange auto& r) {
	std::cout << "Range contents: ";
	for (const auto& item : r) 
		std::cout << item << " ";
	std::cout << "\n";
}
void example_range_concept() {
	std::cout << "\n=== Range Concept ===\n";

	std::vector<int> nums = { 1, 2, 3 };
	print_range(nums);

	std::array<std::string, 2> strs = { "Hello", "World" };
	print_range(strs);
}

void test_concepts(void) {
	example_basic_concepts();
	example_custom_concept();
	example_iterator_concept();
	example_range_concept();
}