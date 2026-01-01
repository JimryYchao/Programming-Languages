#include <iostream>
#include <execution>
#include <vector>
#include <algorithm>
#include <numeric>
#include <chrono>
#include <string>
#include <future>
#include <random>

module GeneralUtilities;
using namespace std;

// execution policies
void execution_policies(int vecSize) {
	std::vector<int> v(vecSize);
	std::iota(v.begin(), v.end(), 1);  // Fill with 1, 2, 3, ..., vecSize
	auto ope = [](int& n) { n = n * n; };
	cout << "vector size: " << vecSize << endl;

	// Sequential 顺序执行策略
	auto v_seq = v;
	{
		auto start = std::chrono::high_resolution_clock::now();
		std::for_each(std::execution::seq, v_seq.begin(), v_seq.end(), ope);
		auto end = std::chrono::high_resolution_clock::now();
		auto seq_time = std::chrono::duration_cast<std::chrono::microseconds>(end - start).count();
		std::cout << "Sequential execution time: " << (double)seq_time / 1000 << " ms" << std::endl;
	}

	// Parallel 并行执行策略
	auto v_par = v;
	{
		auto start = std::chrono::high_resolution_clock::now();
		std::for_each(std::execution::par, v_par.begin(), v_par.end(),
			ope);
		auto end = std::chrono::high_resolution_clock::now();
		auto par_time = std::chrono::duration_cast<std::chrono::microseconds>(end - start).count();
		std::cout << "Parallel execution time: " << (double)par_time / 1000 << " ms" << std::endl;
	}

	// 并行无序执行策略
	auto v_par_unseq = v;
	{
		auto start = std::chrono::high_resolution_clock::now();
		std::for_each(std::execution::par_unseq, v_par_unseq.begin(), v_par_unseq.end(), ope);
		auto end = std::chrono::high_resolution_clock::now();
		auto par_unseq_time = std::chrono::duration_cast<std::chrono::microseconds>(end - start).count();
		std::cout << "Parallel and unsequenced execution time: " << (double)par_unseq_time / 1000 << " ms" << std::endl;
	}

	// 无序执行策略
	auto v_unseq = v;
	{
		auto start = std::chrono::high_resolution_clock::now();
		std::for_each(std::execution::unseq, v_unseq.begin(), v_unseq.end(),
			[](int& n) { n = n * n; });
		auto end = std::chrono::high_resolution_clock::now();
		auto unseq_time = std::chrono::duration_cast<std::chrono::microseconds>(end - start).count();
		std::cout << "Unsequenced execution time: " << (double)unseq_time / 1000 << " ms" << std::endl;
	}

	// Verify results are the same
	bool all_same = (v_seq == v_par) && (v_par == v_par_unseq) && (v_par_unseq == v_unseq);
	std::cout << "All results are the same: " << std::boolalpha << all_same << std::endl;
}

void example_execution_policies() {
	std::cout << "\n=== execution Policies ===\n";
	execution_policies(100);
	execution_policies(10000);
	execution_policies(1000000);
	execution_policies(100000000);
}

void test_execution() {
	example_execution_policies();
}