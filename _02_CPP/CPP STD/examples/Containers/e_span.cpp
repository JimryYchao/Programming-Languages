#include <iostream>
#include <span>
#include <vector>
#include <array>
#include <algorithm>
#include <numeric>
#include <iomanip>

module Containers;
using namespace std;

// std::span 
void print(const auto& sp, const std::string& label = "") {
	std::cout << label << ": ";
	for (const auto& e : sp) std::cout << e << " ";
	std::cout << std::endl;
}
void prints(const auto& sp, const std::string& label = "") {
	std::cout << label << ": ";
	for (const auto& e : sp) std::cout << e;
	std::cout << std::endl;
}
void example_span() {
	std::cout << "\n=== std::span ===\n";

	// 基本构造与初始化
	int arr[] = { 1, 2, 3, 4, 5 };
	std::array<int, 5> array = { 6, 7, 8, 9, 10 };
	std::vector<int> vec(5);
	std::iota(vec.begin(), vec.end(), 11);  // 填充 11,..., 15
	std::string str = "HELLO WORLD"s;
	std::string_view sv = "Hello string_view"sv;
	const char* c_str = "Hello C";

	// 从不同容器创建span
	std::span sp_arr(arr);
	std::span sp_array(array);
	std::span sp_vec(vec);
	std::span sp_str(str); // char
	std::span sp_sv(sv);   // const char
	std::span<const char> sp_cstr(c_str, strlen(c_str));
	std::span sp_literal("Hello World");   // const char

	print(sp_arr, "from arr[]");
	print(sp_array, "from array");
	print(sp_vec, "from vec");
	prints(sp_str, "from str");
	prints(sp_sv, "from sv");
	prints(sp_cstr, "from c_str");
	prints(sp_literal, "from literal");

	std::span sp_vec_part(vec.begin() + 2, 2);
	std::span<int> sp_arr_part(arr, 3);
	print(sp_vec_part, "vec[2,3]");
	print(sp_arr_part, "arr[0,3)");

	// size
	std::cout << "sp_arr.size(): " << sp_arr.size() << std::endl;
	std::cout << "sp_arr.size_bytes(): " << sp_arr.size_bytes() << std::endl;

	// modifier
	sp_str[0] = 'h';
	prints(sp_str, "H to h");
	std::cout << "str: " << str << std::endl;

	// sub span
	std::span hello = sp_str.first(5);
	std::span world = sp_str.last(5);
	prints(hello, "first 5");
	prints(world, "last 5");
	std::span sub_sp_literal = sp_literal.subspan(6);
	prints(sub_sp_literal, "sub_sp_literal");
}

void example_as_bytes() {
	std::cout << "\n=== span as_bytes ===\n";

	auto printbytes = [](float const x, std::span<const std::byte> const bytes)
		{
			std::cout << std::setprecision(6) << std::setw(8) << x << " = { "
				<< std::hex << std::uppercase << std::setfill('0');
			for (auto const b : bytes)
				std::cout << std::setw(2) << std::to_integer<int>(b) << ' ';
			std::cout << std::dec << "}\n";
		};

	float data[1]{ 3.141592f };
	auto const const_bytes = std::as_bytes(std::span{ data });
	printbytes(data[0], const_bytes);

	auto const writable_bytes = std::as_writable_bytes(std::span{ data });
	// Change the sign bit that is the MSB (IEEE 754 Floating-Point Standard).
	writable_bytes[3] |= std::byte{ 0B1000'0000 };
	printbytes(data[0], const_bytes);
}

void test_span() {
	example_span();
	example_as_bytes();
}