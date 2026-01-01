#include <iostream>
#include <utility>
#include <string>
#include <functional>
#include <memory>
#include <optional>
#include <tuple>
#include <variant>
#include <any>
#include <expected>

module GeneralUtilities;
using namespace std;

// std::pair 
void example_pair() {
	std::cout << "\n=== std::pair ===\n";

	// std::pair, tuple-like
	std::pair<int, std::string> my_pair(42, "Hello");
	using P = decltype(my_pair);
	std::cout << "Pair size = " << std::tuple_size_v<P> << endl;
	std::cout << "Pair first: " << my_pair.first << ", second: " << my_pair.second << std::endl;
	auto another_pair = std::make_pair(100, "World");
	cout << "Pair(42, Hello) < Pair(100, World) : " << boolalpha << (my_pair < another_pair) << endl;
}

// 移动与交换
void example_move_and_swap_operations() {
	std::cout << "\n=== move & swap ===\n";

	// std::move 
	std::string str1 = "Hello, move semantics!";
	std::string str2 = std::move(str1);
	cout << "str1 after move: " << str1 << endl;
	cout << "str2 after move: " << str2 << std::endl;

	//std::move_if_noexcept
	struct Bad {
		Bad() {};
		Bad(Bad&&) { cout << "Throwing move ctor called\n"; }
		Bad(const Bad&) { cout << "Throwing copy ctor called\n"; }
	};
	struct Good {
		Good() {}
		Good(Good&&) noexcept { cout << "Non-throwing move ctor called\n"; }
		Good(const Good&) noexcept { cout << "Non-throwing copy ctor called\n"; }
	};
	Good _g; Bad _b;
	Good good = std::move_if_noexcept(_g); // 调用移动构造
	Bad bad = std::move_if_noexcept(_b);   // 调用复制构造

	// std::swap 
	int a = 10, b = 20;
	std::cout << "Before swap: a = " << a << ", b = " << b << std::endl;
	std::swap(a, b);
	std::cout << "After swap: a = " << a << ", b = " << b << std::endl;

	// std::exchange - 替换值并返回旧值
	a = 20;
	int old_value = std::exchange(a, 50);
	std::cout << "Old value: " << old_value << ", new value: " << a << std::endl;
}

// 参数转发与引用操作
struct de_Person {
	std::string name;
	template<typename Self>
	auto&& get_name(this Self&& self) {
		return std::forward_like<Self>(self.name);
	}
};
void example_argument_forwarding() {
	std::cout << "\n=== forwarding ===\n";

	auto forward_example = []<typename T>(T && arg) {
		if constexpr (std::is_lvalue_reference_v<decltype(arg)>)
			std::cout << "arg is an lvalue" << std::endl;
		else std::cout << "arg is an rvalue" << std::endl;
	};
	std::string str3 = "Test";
	forward_example(str3);	 // Lvalue
	forward_example(string("Test"));	 // Rvalue
	forward_example(std::move(str3)); // Rvalue

	// std::forward_like - 基于参考类型的转发
	auto forward_like_example = [](auto&& name) {
		if constexpr (std::is_lvalue_reference_v<decltype(name)>)
			std::cout << name << " is an lvalue" << std::endl;
		else std::cout << name << " is an rvalue" << std::endl; };

	de_Person p{ "Alice" };
	auto name1 = p.get_name();  // 返回 std::string&
	auto name2 = de_Person{ "Bob" }.get_name();  // 返回 std::string&& like self(Person{"Bob"})
	forward_like_example(p.get_name());
	forward_like_example(de_Person{ "Bob" }.get_name());
	forward_like_example(std::move(de_Person{ "Tom" }).get_name());

	// std::as_const - 转换为 const 引用
	std::string str = "Test string";
	const std::string& const_ref = std::as_const(str);
	std::cout << "Const reference: " << const_ref << std::endl;

	// std::declval 右值
	static_assert(is_same_v<int&&, decltype(declval<int>())>);
	static_assert(is_same_v<int&, decltype(declval<int&>())>);
	static_assert(is_same_v<int&&, decltype(declval<int&&>())>);
	static_assert(is_same_v<const int&&, decltype(declval<const int>())>);
	static_assert(is_same_v<volatile int&&, decltype(declval<volatile int>())>);

	struct D { int foo() { return 1; }; };
	struct S { S() = delete; int foo() { return 1; }; };

	static_assert(is_same_v<decltype(D().foo()), int>);
	//static_assert(is_same_v<decltype(S().foo()), int>);  // no ctor
	static_assert(is_same_v<decltype(declval<S>().foo()), int>);
}

// 类型与比较操作
void example_type_and_comparison_operations() {
	std::cout << "\n=== range ===\n";

	// std::to_underlying - 将枚举转换为其底层类型
	enum class Color : short { Red = 1, Green = 2, Blue = 3 };
	static_assert (is_same_v<decltype(to_underlying(Color::Red)), short>);
	std::cout << "Color::Green underlying value: " << dec << std::to_underlying(Color::Green) << std::endl;

	// std::in_range - 检查值是否在类型范围内
	std::cout << "10000000000 is in range of int: " << std::in_range<int>(10000000000LL) << std::endl;
	std::cout << "-1 is in range of unsigned char: " << std::in_range<unsigned char>(-1) << std::endl;

	// cmp
	cout << "-1 > 0u : " << boolalpha << (-1 > 0u) << endl;
	std::cout << "cmp_greater(-1 > 0u): " << std::cmp_greater(-1, 0u) << std::endl;
}

// 分段构造, 消除接受两个元组参数的不同函数之间的歧义
static struct Point2D {
	int a, b;
	Point2D(int a, int b) { std::cout << "ctor with (" << a << ", " << b << ")\n"; }
	Point2D(const Point2D& other) { std::cout << "copy Point\n"; };
	Point2D(Point2D&& other) { std::cout << "move Point\n"; };
};
void example_piecewise_construct() {
	std::cout << "\n=== piecewise_construct ===\n";

	// 使用 piecewise_construct 构造包含 Complex 对象的 pair
	cout << "Creating p1 : \n"; // 构造+移动
	std::pair<Point2D, Point2D> p1(Point2D(1, 2), Point2D(3, 4));
	cout << "Creating p2: \n";  // 就地构造
	std::pair<Point2D, Point2D> p2(std::piecewise_construct, forward_as_tuple(1, 2), forward_as_tuple(3, 4));
}

// 就地构造, 消歧义标签
void example_in_place_operations() {
	std::cout << "\n=== in_place  ===\n";

	// in_place_type
	std::any a(Point2D(1, 2)); // 构造+移动
	a.emplace<Point2D>(5, 6);
	
	std::any a_p(std::in_place_type<Point2D>, 3, 4);  // 就地构造
	
	std::optional<Point2D> op_p(std::in_place, 7, 8);
	
	std::variant<Point2D> var_p(std::in_place_type<Point2D>, 9, 10);
	std::variant<Point2D, int> var_int(std::in_place_type<int>, 10086);
	std::variant<int, monostate, Point2D> var_mono(std::in_place_index<1>);
	std::variant<int, monostate, Point2D> var_i0(std::in_place_index<0>, 10010);
	std::variant<int, monostate, Point2D> var_i2(std::in_place_index<2>, 11, 12);

	std::expected<Point2D, string> ex_p(std::in_place, 13, 14);
	std::expected<Point2D, string> ex_err(std::unexpect, "Something went wrong");
}


void test_utility() {
	example_pair();
	example_move_and_swap_operations();
	example_argument_forwarding();
	example_type_and_comparison_operations();
	example_piecewise_construct();
	example_in_place_operations();

	return;
	std::unreachable();
}