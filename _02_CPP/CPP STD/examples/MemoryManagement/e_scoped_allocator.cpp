#include <iostream>
#include <vector>
#include <string>
#include <scoped_allocator>
#include <memory_resource>
#include <chrono>
#include <thread>
#include <mutex>
#include <string_view>

module MemoryManagement;
using namespace std;
using namespace std::pmr;

// 自定义分配器与 scoped_allocator 结合使用
template <typename T>
class CustomAllocator {
public:
	static int count;
	using value_type = T;
	using pointer = T*;
	using const_pointer = const T*;
	using reference = T&;
	using const_reference = const T&;
	using size_type = std::size_t;
	using difference_type = std::ptrdiff_t;
	CustomAllocator() noexcept {}
	template <typename U>
	CustomAllocator(const CustomAllocator<U>&) noexcept {}
	~CustomAllocator() {}
	T* allocate(size_t n) {
		if (n > std::numeric_limits<size_t>::max() / sizeof(T))
			throw std::bad_alloc();
		std::cout << "CustomAllocator: Allocated " << n * sizeof(T) << " bytes\n";
		return static_cast<T*>(::operator new(n * sizeof(T)));
	}
	void deallocate(T* p, size_t n) {
		std::cout << "CustomAllocator: Deallocated " << n * sizeof(T) << " bytes\n";
		::operator delete(p);
	}
	bool operator==(const CustomAllocator&) const noexcept {
		return true;
	}
	bool operator!=(const CustomAllocator& other) const noexcept {
		return false;
	}
};
void example_custom_allocator_with_scoped() {
	std::cout << "\n=== custom allocator with scoped ===\n";
	using outerVec = std::vector<std::string, CustomAllocator<std::string>>;
	using ScopedCustomAlloc = scoped_allocator_adaptor<CustomAllocator<outerVec>, CustomAllocator<std::string>>;
	std::vector<outerVec, ScopedCustomAlloc> nested_vec;

	nested_vec.resize(2);  // 创建两个内层 vector
	nested_vec[0].push_back("Hello");
	nested_vec[0].push_back("World");
	nested_vec[1].push_back("C++");
	nested_vec[1].push_back("Scoped Allocator");

	// 遍历嵌套容器并显示内容
	std::cout << "Nested container contents:\n";
	for (size_t i = 0; i < nested_vec.size(); ++i) {
		std::cout << "Inner vector " << i << ": ";
		for (const auto& str : nested_vec[i]) {
			std::cout << str << " ";
		}
		std::cout << "\n";
	}
	std::cout << "Nested container size: " << nested_vec.size() << "\n";
}

// 在复杂数据结构中使用 polymorphic_allocator（不需要 scoped_allocator_adaptor）
void example_scoped_allocator_in_complex_structures() {
	std::cout << "\n=== polymorphic_allocator in Complex Data Structures ===\n";
	// 定义一个复杂结构体
	struct Person {
		std::pmr::string name;  // 使用pmr::string
		int age;
		std::pmr::vector<std::string> hobbies;  // 使用 pmr::vector 和 pmr::string
		// 使用 pmr::polymorphic_allocator 的构造函数
		Person(polymorphic_allocator<char> alloc) : name(alloc), age(0), hobbies(alloc) {}
		Person(std::string n, int a, std::pmr::vector<std::string> h, polymorphic_allocator<char> alloc)
			: name(std::move(n), alloc), age(a), hobbies(std::move(h)) {
		}
	};

	// 创建内存资源
	synchronized_pool_resource pool;
	polymorphic_allocator<char> char_alloc(&pool);
	// polymorphic_allocator 会自动传播到嵌套容器
	std::pmr::vector<Person> people(char_alloc);

	// 添加元素到容器
	people.emplace_back("Zhang San", 30, std::pmr::vector<std::string>({ "Reading", "Traveling" }, char_alloc), char_alloc);
	people.emplace_back("Li Si", 25, std::pmr::vector<std::string>({ "Sports", "Music" }, char_alloc), char_alloc);

	// 显示人员信息
	std::cout << "Person information:\n";
	for (const auto& person : people) {
		std::cout << "Name: " << person.name << ", Age: " << person.age << ", Hobbies: ";
		for (const auto& hobby : person.hobbies)
			std::cout << hobby << " ";
		std::cout << "\n";
	}
}

// 多线程环境中使用 polymorphic_allocator
void example_scoped_allocator_in_multithreading() {
	std::cout << "\n=== polymorphic_allocator in Multithreading ===\n";

	// 使用 synchronized_pool_resource，它是线程安全的
	synchronized_pool_resource shared_pool;
	polymorphic_allocator<char> char_alloc(&shared_pool);

	// 共享的容器，使用互斥锁保护
	// pmr::vector会自动将分配器传播到嵌套容器
	std::pmr::vector<std::pmr::vector<int>> shared_data(char_alloc);
	std::mutex data_mutex;
	// 创建多个线程向共享容器添加数据
	std::vector<std::thread> threads;
	for (int t = 0; t < 3; ++t) {
		threads.emplace_back([t, &shared_data, &data_mutex, &char_alloc]() {
			// 创建线程局部数据，使用相同的分配器
			std::pmr::vector<int> local_data(char_alloc);
			for (int i = 0; i < 100; ++i)
				local_data.push_back(t * 1000 + i);
			// 线程安全地添加到共享容器
			{
				std::lock_guard<std::mutex> lock(data_mutex);
				shared_data.push_back(local_data);
			}
			});
	}
	// 等待所有线程完成
	for (auto& thread : threads)
		thread.join();
	// 验证结果
	std::cout << "Shared container size after multithreading: " << shared_data.size() << "\n";
	std::cout << "First 5 elements of the first inner container: ";
	for (int i = 0; i < 5 && i < shared_data[0].size(); ++i)
		std::cout << shared_data[0][i] << " ";
	std::cout << "\n";
}

// scoped_allocator 的性能对比示例
void example_scoped_allocator_performance() {
	std::cout << "\n=== scoped_allocator Performance Comparison ===\n";

	const int OUTER_LOOPS = 100;
	const int INNER_LOOPS = 100;

	// 使用标准分配器
	auto start = std::chrono::high_resolution_clock::now();
	{
		std::vector<std::vector<std::string>> standard_nested;
		for (int i = 0; i < OUTER_LOOPS; ++i) {
			standard_nested.emplace_back();
			for (int j = 0; j < INNER_LOOPS; ++j)
				standard_nested.back().push_back("S" + to_string(i * OUTER_LOOPS + j));
		}
	}
	auto end = std::chrono::high_resolution_clock::now();
	auto standard_time = std::chrono::duration_cast<std::chrono::microseconds>(end - start).count();
	std::cout << "Standard allocator time: " << standard_time << " microseconds\n";

	// 使用 scoped_allocator 与内存池
	start = std::chrono::high_resolution_clock::now();
	{
		using outerVec = std::vector<std::string, allocator<std::string>>;
		using ScopedCustomAlloc = scoped_allocator_adaptor<allocator<outerVec>, allocator<std::string>>;
		std::vector<outerVec, ScopedCustomAlloc> scoped_nested;
		for (int i = 0; i < OUTER_LOOPS; ++i) {
			scoped_nested.emplace_back();
			for (int j = 0; j < INNER_LOOPS; ++j)
				scoped_nested.back().push_back("S" + to_string(i * OUTER_LOOPS + j));
		}
	}
	end = std::chrono::high_resolution_clock::now();
	auto scoped_time = std::chrono::duration_cast<std::chrono::microseconds>(end - start).count();
	std::cout << "Scoped allocator time: " << scoped_time << " microseconds\n";
	// 计算性能提升百分比
	double improvement = (standard_time - scoped_time) * 100.0 / standard_time;
	std::cout << "Performance improvement: " << improvement << "%\n";
}

void test_scoped_allocator() {
	example_custom_allocator_with_scoped();
	example_scoped_allocator_in_complex_structures();
	example_scoped_allocator_in_multithreading();
	example_scoped_allocator_performance();
}