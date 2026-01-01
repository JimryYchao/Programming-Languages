#include <iostream>
#include <array>
#include <algorithm>
#include <numeric>

module Containers;
using namespace std;

// std::array 固定长度数组
void example_array() {
	std::cout << "\n=== std::array ===\n";

	// ctor
	array<int, 5> arr_init = { 1, 2, 3, 4, 5 };
	array<int, 5> arr_fill; 
	arr_fill.fill(42); // Fill all elements with 42
	int o_arr[]{ 1,2,3,4,5 };
	array arr_ori = std::to_array(o_arr); 
	std::array<int, 10> arr_iota;
	ranges::iota(arr_iota, 1);

	// capacity
	auto arr = array(arr_init);
	std::cout << "arr size: " << arr.size() << ", max_size: " << arr.max_size() << ", empty: " << std::boolalpha << arr.empty() << std::endl;
	std::cout << "arr.front: " << arr.front() << ", arr.back: " << arr.back() << ", arr[2]: " << arr[2] << ", arr.at(3): " << arr.at(3) << std::endl;
}

void test_array() {
	example_array();
}
