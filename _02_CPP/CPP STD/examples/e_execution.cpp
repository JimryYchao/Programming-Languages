//module;

#include <execution>
#include <algorithm>
#include <iostream>
#include <vector>
#include <chrono>
#include <numeric>

using namespace std;

// 示例1: 演示并行执行策略
void example_parallel_policy() {
	std::vector<int> data(10000000);
	std::iota(data.begin(), data.end(), 0);

	// 顺序执行
	auto start_seq = std::chrono::high_resolution_clock::now();
	std::sort(std::execution::seq, data.begin(), data.end());
	auto end_seq = std::chrono::high_resolution_clock::now();

	// 并行执行
	auto start_par = std::chrono::high_resolution_clock::now();
	std::sort(std::execution::par, data.begin(), data.end());
	auto end_par = std::chrono::high_resolution_clock::now();

	std::chrono::duration<double> seq_time = end_seq - start_seq;
	std::chrono::duration<double> par_time = end_par - start_par;

	std::cout << "Sequential sort time: " << seq_time.count() << " seconds\n";
	std::cout << "Parallel sort time: " << par_time.count() << " seconds\n";
	std::cout << "Speedup: " << seq_time.count() / par_time.count() << "x\n";
}

// 示例2: 演示向量化执行策略
void example_vectorized_policy() {
	std::vector<double> a(1000000, 1.0);
	std::vector<double> b(1000000, 2.0);
	std::vector<double> c(1000000);

	// 使用并行+向量化策略执行元素级加法
	std::transform(std::execution::par_unseq,
		a.begin(), a.end(), b.begin(), c.begin(),
		[](double a, double b) { return a + b; });

	// 验证结果
	bool all_correct = std::all_of(std::execution::par,
		c.begin(), c.end(),
		[](double val) { return val == 3.0; });

	std::cout << "Vector addition result is "
		<< (all_correct ? "correct" : "incorrect") << "\n";
}

// 示例3: 并行累加
void example_parallel_reduce() {
	std::vector<int> nums(1000000);
	std::iota(nums.begin(), nums.end(), 1);

	// 并行累加
	int sum = std::reduce(std::execution::par,
		nums.begin(), nums.end());

	std::cout << "Sum of first " << nums.size()
		<< " integers: " << sum << "\n";
}

// 示例4: 并行查找
void example_parallel_find() {
	std::vector<int> data(10000000);
	std::iota(data.begin(), data.end(), 0);
	int target = 9999999;

	// 并行查找
	auto it = std::find(std::execution::par,
		data.begin(), data.end(),
		target);

	if (it != data.end()) {
		std::cout << "Found target " << target
			<< " at position " << (it - data.begin()) << "\n";
	}
	else {
		std::cout << "Target not found\n";
	}
}

// 示例5: 并行for_each
void example_parallel_for_each() {
	std::vector<int> data(20);

	std::for_each(std::execution::par,
		data.begin(), data.end(),
		[](int& n) {
			n = std::rand() % 100;
			std::cout << "Thread " << std::this_thread::get_id()
				<< " processing element\n";
		});

	std::cout << "Random numbers: ";
	for (int n : data) {
		std::cout << n << " ";
	}
	std::cout << "\n";
}

void test_execution(void) {
	std::cout << "=== Example 1: Parallel vs Sequential Sort ===\n";
	example_parallel_policy();

	std::cout << "\n=== Example 2: Vectorized Operations ===\n";
	example_vectorized_policy();

	std::cout << "\n=== Example 3: Parallel Reduce ===\n";
	example_parallel_reduce();

	std::cout << "\n=== Example 4: Parallel Find ===\n";
	example_parallel_find();

	std::cout << "\n=== Example 5: Parallel For Each ===\n";
	example_parallel_for_each();
}