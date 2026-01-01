#include <iostream>
#include <vector>
#include <chrono>
#include <map>
#include <random>
#include <boost/container/flat_map.hpp>

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
class Timer_map {
private:
	steady_clock::time_point start_time;
	std::string operation_name;
public:
	Timer_map(const std::string& name) : operation_name(name) {
		start_time = steady_clock::now();
	}

	~Timer_map() {
		auto end_time = steady_clock::now();
		auto duration = 0.001 * duration_cast<microseconds>(end_time - start_time).count();
		std::cout << operation_name << " Time used: " << duration << " ms" << std::endl;
	}
};

// 计算 map 的 heap size
template<class Map>
size_t sizeof_map_heap(const Map& m) {
	if (m.empty()) return 0;

	// 估算每个节点的大小
	// 红黑树节点通常包含：
	// - 3个指针（left, right, parent）
	// - 1个bool（颜色）
	// - 键值对数据
	// - 分配器可能有的开销
	size_t pointer_size = sizeof(void*);
	size_t bool_size = sizeof(bool);
	size_t data_size = sizeof(typename Map::value_type);

	// 典型节点大小估算
	size_t node_size = 3 * pointer_size + bool_size + data_size;
	// 分配器开销（通常 8-16 字节）
	size_t allocation_overhead = 16;
	return m.size() * (node_size + allocation_overhead);
}

// 计算 flat_map 的 heap size
template <class FlatMap>
size_t sizeof_flatmap_heap(const FlatMap& fm) {
	return fm.capacity() * sizeof(typename FlatMap::value_type);
}

// flat_map 和 map 的性能比较
void example_compare_flatmap_map() {
	const size_t DATA_SIZE = 100000;
	const size_t LOOKUP_COUNT = 10000;
	std::cout << "\n=== Compare flat_map & map  (Data size: " << DATA_SIZE << ") ===\n";

	// 生成随机数据
	auto keys = generateRandomNumbers<int>(DATA_SIZE);
	auto lookup_keys = generateRandomNumbers<int>(LOOKUP_COUNT, 0, 1000000);
	std::map<int, int> m;
	flat_map<int, int> fm;

	// 插入性能测试
	// ------------------------------
	std::cout << "\n>>> operator[]: \n";
	// ------------------------------
	{  // map 
		Timer_map t("map");
		for (size_t i = 0; i < DATA_SIZE; ++i)
			m[keys[i]] = i;
	}
	{  // flat_map 
		Timer_map t("flat_map");
		for (size_t i = 0; i < DATA_SIZE; ++i)
			fm[keys[i]] = i;
	}

	m.clear();
	fm.clear();
	// ------------------------------
	std::cout << "\n>>> emplace: \n";
	// ------------------------------
	{ // map 
		Timer_map t("map");
		for (size_t i = 0; i < DATA_SIZE; ++i)
			m.emplace(keys[i], i);
	}
	{ // flat_map 
		Timer_map t("flat_map");
		for (size_t i = 0; i < DATA_SIZE; ++i)
			fm.emplace(keys[i], i);
	}


	// 查找性能测试
	// ------------------------------
	std::cout << "\n>>> find: \n";
	// ------------------------------
	{ // map
		int found_count = 0;
		Timer_map t("map");
		for (size_t i = 0; i < LOOKUP_COUNT; ++i) {
			auto it = m.find(lookup_keys[i]);
			if (it != m.end())
				found_count++;
		}
		std::cout << "map find " << found_count << " elements" << std::endl;
	}
	{ //  flat_map
		int found_count = 0;
		Timer_map t("flat_map");
		for (size_t i = 0; i < LOOKUP_COUNT; ++i) {
			auto it = fm.find(lookup_keys[i]);
			if (it != fm.end())
				found_count++;
		}
		std::cout << "flat_map find " << found_count << " elements" << std::endl;
	}


	// 删除性能比较
	// ------------------------------
	std::cout << "\n>>> erase：\n";
	// ------------------------------
	auto m_copy = m;
	auto fm_copy = fm;
	{ // map 
		Timer_map t("map");
		for (size_t i = 0; i < LOOKUP_COUNT / 2; ++i)
			m_copy.erase(lookup_keys[i]);
	}
	{ // flat_map
		Timer_map t("flat_map");
		for (size_t i = 0; i < LOOKUP_COUNT / 2; ++i) {
			fm_copy.erase(lookup_keys[i]);
		}
	}


	// 内存使用比较
	// ------------------------------
	std::cout << "\ncompare memory used：\n";
	// ------------------------------

	std::cout << "map size: " << m.size() << ", heap: " << sizeof_map_heap(m) << " bytes" << std::endl;
	std::cout << "flat_map size: " << fm.size() << ", heap: " << sizeof_flatmap_heap(fm) << " bytes"  << std::endl;


	// 迭代性能比较
	// ------------------------------
	std::cout << "\n>>> iterate：\n";
	// ------------------------------
	{ // map
		long long sum = 0;
		Timer_map t("map");
		for (const auto& [key, value] : m) 
			sum += key + value;
	}
	{ // flat_map
		long long sum = 0;
		Timer_map t("flat_map");
		for (const auto& [key, value] : fm) 
			sum += key + value;
	}
}

void test_flat_map() {
	example_compare_flatmap_map();
}