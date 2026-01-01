#include <iostream>
#include <list>
#include <algorithm>
#include <numeric>
#include <vector>

module Containers;
using namespace std;


template <class T, class Alloc>
void print(const std::list<T, Alloc>& list, const std::string& label) {
	std::cout << label << ": ";
	for (const auto& num : list) std::cout << num << " ";
	std::cout << std::endl;
}

// std::list 双向链表
void example_list() {
	std::cout << "\n=== std::list ===\n";

	// ctor
	list<int> l_init{ 1, 2, 3, 4, 5 };
	list<int> l_c(3);
	list<int> l_cv(5, 1);
	list l_lt(l_init.begin(), l_init.end());  // 填充 1,..., 10

	// capacity
	list l(l_init);
	cout << "l size: " << l.size() << ", empty: " << l.empty() << endl;
	cout << "l.front: " << l.front() << ", l.back: " << l.back() << endl;
	l.resize(10, 666);
	print(l, "l resize");

	// front
	l.push_front(10010);
	l.prepend_range(initializer_list{ 10,20,30 });
	l.pop_front();
	l.emplace_front(666);
	print(l, "l (front)");

	// back
	l.push_back(10086);
	l.append_range(initializer_list{ 90,80,70 });
	l.pop_back();
	l.emplace_back(999);
	print(l, "l (back)");

	// insert
	l.insert(prev(l.end()), { 9,99 });
	l.insert_range(next(l.begin()), initializer_list{ 66,6 });
	auto mid = l.begin();
	advance(mid, l.size() / 2);
	l.emplace(mid, 555);
	print(l, "l (insert)");

	// remove
	l.remove(555);
	l.remove_if([](const auto& e) {return e <= 10; });
	print(l, "l (remove)");

	// sort and merge
	list l_m({ 95,15,85,25,75,35,65,45,55 });
	l.sort(greater<>{});
	l_m.sort(greater<>{});
	l.merge(l_m, greater<>{});
	cout << "l size: " << l.size() << ", l_m size: " << l_m.size() << endl;
	print(l, "l (merge)");

	// reverse and splice
	l.reverse();
	l.splice(l.begin(), l_init);
	cout << "l size: " << l.size() << ", l_init size: " << l_init.size() << endl;
	print(l, "l (reverse)");

	// erase_if and unique
	l.unique();
	erase_if(l, [](auto& v) {return  v < 50 || v > 100; });
	print(l, "l (unique)");
}

void test_list() {
	example_list();
}