#include <iostream>
#include <map>
#include <string>
#include <algorithm>
#include <iterator>
#include <vector>

module Containers;
using namespace std;

void print(auto& m, const std::string& name) {
	std::cout << name << ": ";
	for (const auto& pair : m)
		std::cout << "{" << pair.first << "," << pair.second << "} ";
	std::cout << std::endl;
}
template<typename T>
struct Person_m {
	std::string name;
	int age;
	Person_m(const std::string& n, int a) : name(n), age(a) {}
	auto operator <=> (const Person_m& other) const = default;
};
static struct CaseInsensitiveLess { // 不区分大小写的比较器
	bool operator()(const std::string& a, const std::string& b) const {
		std::string a_lower, b_lower;
		std::transform(a.begin(), a.end(), std::back_inserter(a_lower), ::tolower);
		std::transform(b.begin(), b.end(), std::back_inserter(b_lower), ::tolower);
		return a_lower < b_lower;
	}
};

// std::map 
void example_map() {
	std::cout << "\n=== std::map ===\n";

	std::map<int, std::string> m1;
	std::map<int, std::string, greater<int>> m2({ {1, "one"}, {2, "two"}, {3, "three"} }, greater<int>{});

	std::map<int, std::string> m3(m2.begin(), m2.end()); // range
	auto m(m2);
	std::cout << "m size: " << m.size() << ", empty: " << std::boolalpha << m.empty() << std::endl;

	// access
	std::cout << "Value for key 2: " << m2[2] << std::endl;
	std::cout << "Value for key 3: " << m2.at(3) << std::endl;
	std::cout << "Accessing non-existent m[10]: " << m2[10] << ", size after access m[10]: " << m.size() << std::endl;

	// search
	auto it = m.find(4);
	if (it != m.end())
		std::cout << "Found key 4 with value: " << it->second << std::endl;
	else std::cout << "set m[4] :" << (m[4] = "four") << std::endl;

	// insert 
	m.insert({ {5, "five"}, { 7, "seven" }, { 9, "nine" } });
	print(m, "m after insert");

	// custom key 
	std::map<Person_m<int>, std::string> people;
	people.insert({ {"Alice", 30}, "Engineer" });
	people.insert({ {"Bob", 25}, "Designer" });
	people.insert({ {"Charlie", 35}, "Manager" });
	people.insert({ {"Alice", 28}, "Developer" }); // 不同年龄，相同名称是允许的

	std::cout << "\nPeople map: " << std::endl;
	for (const auto& entry : people) {
		std::cout << "{Name: " << entry.first.name << ", Age: " << entry.first.age
			<< "}, Occupation: " << entry.second << std::endl;
	}

	// 使用自定义比较器
	std::map<std::string, int, CaseInsensitiveLess> caseInsensitiveMap;
	caseInsensitiveMap["Hello"] = 1;
	caseInsensitiveMap["WORLD"] = 2;
	caseInsensitiveMap["hello"] = 3; // 这会覆盖"Hello"条目
	std::cout << "\nCase-insensitive map: " << std::endl;
	print(caseInsensitiveMap, "caseInsensitiveMap");

	// 合并maps
	m.clear();
	m.insert_range(vector<pair<int, string>>{ {1, "one"}, {2, "two"}, {3, "three"} });
	std::map<int, std::string> m_merged = { {1, "ONE"}, {3, "THREE"}, {4, "four"}, {5, "five"} };

	print(m, "\nm1 before merge");
	print(m_merged, "m2 before merge");
	m.merge(m_merged);
	print(m, "m1 after merge");
	print(m_merged, "m2 after merge");
}

// std::multimap 可重复键映射
void example_multimap() {
	std::cout << "\n=== std::multimap ===\n";

	std::multimap<int, std::string> mm = {
		{1, "one"},
		{2, "two"},
		{2, "TWO"},
		{3, "three"},
		{3, "THREE"}
	};

	mm.insert({ { 2, "2" }, {2, "II"}, {1, "1"}});
	cout << "mm size: " << mm.size() << ", empty: " << std::boolalpha << mm.empty() << std::endl;
	print(mm, "mm");

	// count
	std::cout << "Count of key 2: " << mm.count(2) << std::endl;
	std::cout << "Count of key 3: " << mm.count(3) << std::endl;
	auto range = mm.equal_range(2);
	std::cout << "All values for key 2: ";
	for (auto it = range.first; it != range.second; ++it) 
		std::cout << it->second << " ";
	std::cout << std::endl;

	// erase all k2
	mm.erase(2);
	print(mm, "mm after erasing k2");
}

void test_map() {
	example_map();
	example_multimap();
}