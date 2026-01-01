#include <iostream>
#include <vector>
#include <algorithm>
#include <numeric>
#include <memory_resource>
#include <list>
#include <bitset>

module Containers;
using namespace std;

template<class T, class Alloc>
void print(const vector<T, Alloc> vec, string label) {
	std::cout << label << ": ";
	for (const auto& e : vec) std::cout << e << " ";
	std::cout << std::endl;
}

// std::vector 动态数组
void example_vector() {
	std::cout << "\n=== std::vector ===\n";

	// ctor
	vector<int> v_empty;				  // Empty 
	vector<int> v_init{ 1, 2, 3, 4, 5 };  // Initializer list
	vector<int> v_count(10);		      // size = 10	
	vector<int> v_cv(5, 10);              // 5 elements, each initialized to 10
	vector<int> v_it(v_init.begin(), v_init.end());   // iterator range

	// assign
	auto v_assign(v_init);
	print(v_assign, "v before assign");
	v_assign.assign(5, 10);
	print(v_assign, "v after assign");

	// capacity
	pmr::monotonic_buffer_resource res(1024);
	std::vector<int, pmr::polymorphic_allocator<int>> vec{ std::pmr::polymorphic_allocator<int>(&res) };  // alloc
	std::cout << "\nvec empty: " << std::boolalpha << vec.empty() << ", size: " << vec.size() << ", capacity: " << vec.capacity() << std::endl;
	cout << "vec capaciry after reserve(128) : " << (vec.reserve(vec.capacity() + 128), vec.capacity())
		<< ", size after resize(10, 1) : " << (vec.resize(5, 1), vec.size())
		<< ", vec.first : " << vec.front() << ", vec.last : " << vec.back() << endl;
	cout << "vec capacity after shrink_to_fit : " << (vec.shrink_to_fit(), vec.capacity()) << endl;

	// add to back
	vec.push_back(10086);
	vec.emplace_back(10010);
	vec.append_range(v_init);
	vec.pop_back();   // vec.erase(vec.end());
	print(vec, "\nv (add to back)");

	// insert into front
	vec.insert(vec.begin(), 911);
	vec.insert(vec.begin(), { 6,7,8,9,10 });
	vec.insert_range(vec.begin(), vector<int>(2, 666));
	print(vec, "v (insert into begin)");

	// vector.erase
	std::erase_if(vec, [](auto& e) { return e < 100; });
	print(vec, "v (erase if less 100)");
}

// vector<bool> 动态位集
void example_vector_bool() {
	std::cout << "\n=== std::vector<bool> ===\n";
	std::vector<bool> bitvec(10, false);  
	std::cout << "bitvec size: " << bitvec.size() << ", capacity in bits: " << bitvec.capacity() << std::endl;

	auto printbit = [&bitvec] {
		for (bool b : bitvec) std::cout << (b ? '1' : '0');
		return "";
		};

	// 设置位值
	bitvec[0] = true;
	bitvec[2] = true;
	bitvec[5] = true;
	bitvec[8] = true;
	std::cout << "bitvec: " << printbit() << endl;

	// 翻转所有位
	bitvec.flip();  
	std::cout << "bitvec after flip           : " << printbit() << endl;

	// vector<bool> 的引用行为特殊
	std::vector<bool>::reference ref = bitvec[1];  // 获取位的引用
	ref = false;  // 修改位值
	std::cout << "After setting bit 1 to false: " << printbit() << std::endl;
}


void test_vector() {
	example_vector();
	example_vector_bool();
}