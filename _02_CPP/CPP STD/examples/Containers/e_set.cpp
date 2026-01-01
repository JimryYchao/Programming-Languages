#include <iostream>
#include <set>
#include <string>
#include <algorithm>
#include <iterator>
#include <vector>
#include <random>

module Containers;
using namespace std;

void print(const auto& s, const std::string& name) {
	std::cout << name << ": ";
	for (const auto& num : s) {
		std::cout << num << " ";
	}
	std::cout << std::endl;
};

// std::set - 有序且唯一的元素集合
void example_set() {
	std::cout << "\n=== std::set ===\n";

	std::set<int> s1;
	std::set<int> s2 = { 3, 1, 4, 1, 5, 9 };
	std::set<int> s3(s2.begin(), s2.end());
	std::set<int> s4(s2);

	set<int, greater<int>> s(greater<int>{});  // 倒序
	s.insert_range(vector<int>{9, 5, 6, 8, 7, 1, 5});
	s.merge(s2);
	std::cout << "s size: " << s.size() << ", empty: " << std::boolalpha << s.empty() << std::endl;
	print(s, "s");


	s.insert(10086);
	s.insert(911);
	s.insert({ 100,111,222,333,444,555 });
	auto lg100 = s.lower_bound(100);
	set s_less100(std::next(lg100), s.end()); // 小于 100 的元素
	set s_gte100(s.begin(), next(lg100)); // 大于等于 100 的元素
	s_gte100.erase(s_gte100.upper_bound(500), s_gte100.end()); // 删除大于 500 的元素
	print(s_less100, "s less 100");
	print(s_gte100, "s greater 100");
}

// std::multiset, 允许重复元素的有序集合
void example_multiset() {
	std::cout << "\n=== std::multiset ===\n";

	std::set<int> s = { 3, 1, 4, 1, 5, 9 };  // 1,3,4,5,9

	multiset<int> ms({ 9, 5, 6, 8, 7, 1, 5 });
	ms.merge(s);
	std::cout << "ms size: " << ms.size() << ", empty: " << std::boolalpha << ms.empty() << std::endl;
	std::cout << "s size after merge: " << s.size() << ", empty: " << std::boolalpha << s.empty() << std::endl;
	print(ms, "ms");

	ms.erase(1);
	cout << "the count_1 after erasing k1: " << ms.count(1) << ", ms size: " << ms.size() << endl;

	auto randi = [](int min, int max) {
		static std::random_device rd;
		static std::mt19937 gen(rd());
		std::uniform_int_distribution<int> dis(min, max);
		return dis(gen);
		};
	for (int i = 0; i < 1000; ++i)
		ms.emplace(randi(0, 10));
	cout << "the count of 1 after emplace randi 1000 times: " << ms.count(1) << endl;
}

void test_set() {
	example_set();
	example_multiset();
}