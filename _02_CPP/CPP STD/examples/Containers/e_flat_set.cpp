#include <iostream>
#include <vector>
#include <chrono>
#include <set>
#include <random>
#include <boost/container/flat_set.hpp>

module Containers;
using namespace std;
using namespace boost::container;
using namespace std::chrono;

// 用于生成随机数据的辅助函数
template<typename T>
std::vector<T> generateRandomNumbers(size_t count, T min = 0, T max = 1000000) {
	std::vector<T> numbers;
	numbers.reserve(count);

	random_device rd;
	mt19937 gen(rd());
	uniform_int_distribution<T> dis(min, max);
	for (size_t i = 0; i < count; ++i)
		numbers.push_back(dis(gen));

	return numbers;
}

// 计时器类，用于测量代码执行时间
class Timer_set {
private:
	steady_clock::time_point start_time;
	std::string operation_name;
public:
	Timer_set(const std::string& name) : operation_name(name) {
		start_time = steady_clock::now();
	}

	~Timer_set() {
		auto end_time = steady_clock::now();
		auto duration = 0.001 * duration_cast<microseconds>(end_time - start_time).count();
		std::cout << operation_name << " Time used: " << duration << " ms" << std::endl;
	}
};


// flat_set 和 set 的性能比较
void example_compare_flatset_set() {
	const size_t DATA_SIZE = 100000;
	const size_t LOOKUP_COUNT = 10000;
	std::cout << "\n=== Compare flat_set & set  (Data size: " << DATA_SIZE << ") ===\n";

	// 生成随机数据
	auto keys = generateRandomNumbers<int>(DATA_SIZE);
	auto lookup_keys = generateRandomNumbers<int>(LOOKUP_COUNT, 0, 1000000);
	std::set<int> s;
	flat_set<int> fs;

	// 插入性能测试
	// ------------------------------
	std::cout << "\n>>> insert single element: \n";
	// ------------------------------
	{  // set 
		Timer_set t("set");
		for (size_t i = 0; i < DATA_SIZE; ++i)
			s.insert(keys[i]);
	}
	{  // flat_set 
		Timer_set t("flat_set");
		for (size_t i = 0; i < DATA_SIZE; ++i)
			fs.insert(keys[i]);
	}

	s.clear();
	fs.clear();
	// ------------------------------
	std::cout << "\n>>> insert range: \n";
	// ------------------------------
	{
		// 先打乱数据
		random_device rd;
		mt19937 g(rd());
		std::shuffle(keys.begin(), keys.end(), g);
	}
	{ // set 
		Timer_set t("set");
		s.insert(keys.begin(), keys.end());
	}
	{ // flat_set 
		Timer_set t("flat_set");
		fs.insert(keys.begin(), keys.end());
	}


	// 查找性能测试
	// ------------------------------
	std::cout << "\n>>> find: \n";
	// ------------------------------
	{ // set
		int found_count = 0;
		Timer_set t("set");
		for (size_t i = 0; i < LOOKUP_COUNT; ++i) {
			auto it = s.find(lookup_keys[i]);
			if (it != s.end())
				found_count++;
		}
		std::cout << "set find " << found_count << " elements" << std::endl;
	}
	{ //  flat_set
		int found_count = 0;
		Timer_set t("flat_set");
		for (size_t i = 0; i < LOOKUP_COUNT; ++i) {
			auto it = fs.find(lookup_keys[i]);
			if (it != fs.end())
				found_count++;
		}
		std::cout << "flat_set find " << found_count << " elements" << std::endl;
	}


	// 删除性能比较
	// ------------------------------
	std::cout << "\n>>> erase：\n";
	// ------------------------------
	auto s_copy = s;
	auto fs_copy = fs;
	{ // set 
		Timer_set t("set");
		for (size_t i = 0; i < LOOKUP_COUNT / 2; ++i)
			s_copy.erase(lookup_keys[i]);
	}
	{ // flat_set
		Timer_set t("flat_set");
		for (size_t i = 0; i < LOOKUP_COUNT / 2; ++i) {
			fs_copy.erase(lookup_keys[i]);
		}
	}


	// 迭代性能比较
	// ------------------------------
	std::cout << "\n>>> iterate：\n";
	// ------------------------------
	{ // set
		long long sum = 0;
		Timer_set t("set");
		for (const auto& elem : s) 
			sum += elem;
	}
	{ // flat_set
		long long sum = 0;
		Timer_set t("flat_set");
		for (const auto& elem : fs) 
			sum += elem;
	}
}

void test_flat_set() {
	example_compare_flatset_set();
}