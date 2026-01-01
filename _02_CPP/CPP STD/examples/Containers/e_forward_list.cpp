#include <iostream>
#include <forward_list>
#include <algorithm>
#include <vector>
#include <numeric>

module Containers;
using namespace std;

template<class T>
void print(const forward_list<T>& fl, const string& label) {
	cout << label << ": ";
	for (const auto& num : fl) cout << num << " ";
	cout << endl;
}

// std::forward_list 单向链表
void example_forward_list() {
	cout << "\n=== std::forward_list ===\n";

	// ctor
	forward_list fl_init{ 1, 2, 3, 4, 5 };
	forward_list<int> f_c(5);  // 3 个 10
	forward_list<int> f_cv(3, 10);  // 3 个 10

	// capacity
	forward_list fl(fl_init);
	fl.resize(10, 666);
	cout << "fl empty: " << fl.empty() << ", fl.front: " << fl.front() << endl;
	print(fl, "fl");

	// front
	fl.push_front(10086);
	fl.prepend_range(initializer_list{ 10,20,30 });
	fl.pop_front();
	fl.emplace_front(110);
	print(fl, "fl (front)");

	// after
	fl.insert_after(fl.begin(), { 100,200,300 });
	fl.emplace_after(fl.before_begin(), 999);     // head
	fl.erase_after(fl.begin());
	print(fl, "f (after)");

	// remove and unique
	fl.remove_if([](const auto& v) {return v <= 10; });
	fl.unique();
	print(fl, "fl (remove & unique)");

	// sort & reverse
	fl.sort();
	fl.merge({ 15,155,1555,15555 });
	fl.splice_after(fl.before_begin(), fl_init);	
	fl.reverse();
	print(fl, "fl (reverse)");
}

void test_forward_list() {
	example_forward_list();
}