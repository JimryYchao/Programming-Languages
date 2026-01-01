#include <iostream>
#include <deque>
#include <string>
#include <algorithm>
#include <iterator>
#include <numeric>
#include <vector>
#include <span>

module Containers;
using namespace std;

template<typename T>
void print(const deque<T>& dq, const string& label) {
	cout << label << ": ";
	for (const auto& item : dq)
		cout << item << " ";
	cout << endl;
}

static struct Person_dq {
	string name;
	int age;
	Person_dq(const string& n, int a) : name(n), age(a) {}
};

static ostream& operator<<(ostream& os, const Person_dq& p) {
	return os << "{" << p.name << ", " << p.age << "}";
}

// std::deque 双端队列
void example_deque() {
	cout << "\n=== std::deque ===\n";

	// ctor
	deque<int> dq_init{ 1, 2, 3, 4, 5 };
	deque dq_fill(5, 10); // 5 elements initialized to 10
	deque dq_vec(from_range, initializer_list<int>{1, 2, 3, 4, 5});
	deque dq_str(from_range, "Hello World");  // char
	deque<int> de_iota(10);   // size = 10
	ranges::iota(de_iota, 1); // Fills 1, 2, 3, ..., 10

	// capacity
	deque dq(dq_vec);
	cout << "dq size: " << dq.size() << ", empty: " << boolalpha << dq.empty() << endl;
	print(dq, "dq");
	dq.resize(20, 666);
	print(dq, "after resize");
	dq.clear(); de_iota.shrink_to_fit();
	cout << "after clear and shrink_to_fit, size remains: " << dq.size() << endl;

	// access
	dq.swap(dq_init);
	dq.front() = 999;
	dq.back() = 888;
	dq[1] = 777;
	dq.at(dq.size() - 2) = 666;
	print(dq, "dq (modified)");

	// front
	dq.prepend_range(vector{ 10,20,30 });  // 10,20,30,...
	dq.pop_front();								 // 20, ...
	dq.push_front(10010);					 // 10010,20, ...
	print(dq, "dq (front)");

	// back
	dq.append_range(vector{ 40,50,60 });  // ..., 40,50,60
	dq.pop_back();								// ..., 40,50
	dq.push_back(10086);					// ..., 50,10086
	print(dq, "dq (back)");

	// Insert 
	auto midpos = dq.begin() + dq.size() / 2;
	dq.insert_range(dq.begin() + 1, vector{ 11,22,33 });
	dq.insert(dq.end() - 1, { 44,55,66 });
	print(dq, "dq (insert)");

	// erase 
	erase_if(dq, [](const auto& v) {return v < 100; });
	print(dq, "dq (erase_if e < 100)");
}

void example_customtype_in_deque()
{
	cout << "\n=== Custom type in deque ===\n";
	deque<Person_dq> people;

	people.emplace_back("Alice", 30);
	people.emplace_back("Bob", 25);
	people.emplace_front("Charlie", 35);
	people.emplace_front("David", 28);
	people.emplace(people.begin() + people.size() / 2 , "Mid", 99);
	print(people, "people deque");
	// Sort by age ascending
	ranges::sort(people, [](const Person_dq& a, const Person_dq& b) {
		return a.age < b.age;
		});
	print(people, "people sorted by age");
	// Reverse
	ranges::reverse(people);
	print(people, "people after reverse");
}

void test_deque() {
	example_deque();
	example_customtype_in_deque();
}