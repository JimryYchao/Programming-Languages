#include <iostream>
#include <unordered_map>
#include <string>

module Containers;
using namespace std;

template<typename T>
void print(const T& container, const std::string& name) {
	std::cout << name << ": ";
	for (const auto& [key, value] : container)
		std::cout << "{" << key << "," << value << "} ";
	std::cout << std::endl;
}

static struct Point_um {
	int x, y;
	Point_um(int x_, int y_) : x(x_), y(y_) {}
	bool operator==(const Point_um& other) const { return x == other.x && y == other.y; }
};
struct PointHash {
	std::size_t operator()(const Point_um& p) const {
		std::size_t h1 = std::hash<int>{}(p.x);
		std::size_t h2 = std::hash<int>{}(p.y);
		return h1 ^ (h2 << 1); // 哈希值组合
	}
};
// 不区分大小写的哈希函数 
struct CaseInsensitiveHash {
	std::size_t operator()(const std::string& s) const {
		std::string lower;
		lower.reserve(s.size());
		for (char c : s)
			lower.push_back(std::tolower(c));
		return std::hash<std::string>{}(lower);
	}
};

// 不区分大小写的比较器
struct CaseInsensitiveEqual {
	bool operator()(const std::string& a, const std::string& b) const {
		if (a.length() != b.length()) return false;
		for (size_t i = 0; i < a.length(); ++i)
			if (std::tolower(a[i]) != std::tolower(b[i]))
				return false;
		return true;
	}
};

// std::unordered_map 
void example_unordered_map() {
	std::cout << "\n=== std::unordered_map ===\n";

	std::unordered_map<std::string, int> um = { {"three", 3}, {"one", 1}, {"two", 2}, {"four", 4}, {"six", 6} };
	cout << "um size: " << um.size() << ", empty" << boolalpha << um.empty() << endl;
	print(um, "um");

	// search
	if (um.contains("one"))
		um.at("one") = 10086;
	if (auto p = um.find("five"); p == um.end())
		um["five"] = 5;
	print(um, "um");

	std::cout << "\num bucketIndex of key[two]: " << um.bucket("two") << std::endl;
	std::cout << "In um bucket[two]: ";
	auto b2 = um.bucket("two");
	if (um.bucket_size(b2) > 0)
		for (auto it = um.begin(b2); it != um.end(b2); ++it)
			std::cout << "{" << it->first << "," << it->second << "} " << std::endl;

	// bucket operator
	std::cout << "\num size: " << um.size() << std::endl;
	std::cout << "um bucket count: " << um.bucket_count() << std::endl;
	std::cout << "um loadfactor: " << um.load_factor() << std::endl;

	// 性能调优 
	um.reserve(128); // 预留空间避免多次重哈希
	std::cout << "um bucket count after reserve: " << um.bucket_count() << std::endl;
	um.max_load_factor(0.5); // 设置最大装载因子
	std::cout << "um set max loadfactor: " << um.max_load_factor() << std::endl;
	um.insert_range(initializer_list < pair<string, int>>{
		{"two2", 2 << 2}, { "two3", 2 << 3 }, { "two5", 2 << 5 },
		{ "three1",3 << 1 }, { "three3", 3 << 3 },
		{ "five8", 5 << 8 }});
	cout << "um size after insert_range: " << um.size() << ", bucket count: " << um.bucket_count() << ", loadfactor: " << um.load_factor() << endl;

	// 自定义哈希函数和比较器 
	std::unordered_map<std::string, int, CaseInsensitiveHash, CaseInsensitiveEqual> caseInsensitiveMap;
	caseInsensitiveMap["Hello"] = 1;
	caseInsensitiveMap["WORLD"] = 2;
	caseInsensitiveMap["hello"] = 3; // 覆盖 "Hello"，因为不区分大小写
	print(caseInsensitiveMap, "\ncaseInsensitive um");

	// 自定义键类型 
	std::unordered_map<Point_um, std::string, PointHash> pointMap;
	pointMap[{1, 2}] = "p1";
	pointMap.insert_range(initializer_list<pair<Point_um, string>>{
		{{3, 4}, "p2"}, { {0,0}, "ori" }, { {10,10} , "end" }
	});
	for (const auto& [pt, desc] : pointMap)
		std::cout << "Point{" << pt.x << ", " << pt.y << "}: " << desc << std::endl;
}

// std::unordered_multimap 特有功能示例
void example_unordered_multimap() {
	std::cout << "\n=== std::unordered_multimap ===\n";

	std::unordered_multimap<std::string, int> umm = { {"two", 2}, {"one", 1}, {"two", 22}, {"three", 3} };
	cout << "umm size: " << umm.size() << ", empty: " << boolalpha << umm.empty() << ", bucket_count: " << umm.bucket_count() << endl;

	umm.reserve(128);
	umm.max_load_factor(0.5);
	umm.insert_range(initializer_list<pair<string, int>>{
		{"two", 22222}, { "two",2222 }, { "one", 11111 }, { "three", 333 }, { "one",111 }, { "two",22222222 }
	});


	for (size_t i = 0; i < umm.bucket_count(); i++)
	{
		if (umm.bucket_size(i) > 0) {
			cout << "in bucket " << i << ":";
			for (auto it = umm.begin(i); it != umm.end(i); ++it)
				std::cout << "{" << it->first << "," << it->second << "} ";
			cout << endl;
		}
	}


	// equal_range
	auto twoRange = umm.equal_range("two");
	cout << "the values of key=two: ";
	for (auto it = twoRange.first; it != twoRange.second; ++it)
		std::cout << it->second << " ";
	std::cout << std::endl;

	// count
	std::cout << "the values count of key=two : " << umm.count("two") << std::endl;
}

void test_unordered_map() {
	example_unordered_map();
	example_unordered_multimap();
}