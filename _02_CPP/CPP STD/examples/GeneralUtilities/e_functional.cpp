#include <iostream>
#include <functional>
#include <string>
#include <vector>
#include <algorithm>
#include <numeric>
#include <memory>
#include <tuple>
#include <cmath>
#include <any>
#include <bitset>
#include <chrono>

module GeneralUtilities;
using namespace std;

// Basic function objects
void example_functional_basic() {
	std::cout << "\n=== Basic std::functional Usage ===\n";

	// Using function objects with algorithms
	std::vector<int> numbers = { 3, 1, 4, 1, 5, 9, 2, 6 };
	std::sort(numbers.begin(), numbers.end(), std::greater<int>());
	std::cout << "Sorted in descending order: ";
	for (int n : numbers) { std::cout << n << " "; }
	std::cout << std::endl;

	// Using std::function to store callables
	std::function<int(int, int)> operation;
	operation = std::plus<int>();
	std::cout << "Using std::plus: " << operation(10, 20) << std::endl;
	operation = [](int a, int b) { return a - b; };
	std::cout << "Using lambda: " << operation(10, 20) << std::endl;

	// Function as a parameter
	auto apply_operation = [](int a, int b, std::function<int(int, int)> op) {
		return op(a, b);
		};

	std::cout << "Apply multiplication: " << apply_operation(5, 6, std::multiplies()) << std::endl;

}

// functional wrappers
int sum(int x, int y) { return x + y; };
void examples_functional_wrappers() {
	cout << "\nfuntional wrappers\n";
	// 可复制调用对象的可复制包装器
	function<int(int, int)> add = [](int x, int y) {return x + y; };
	cout << add(10, 20) << endl;
	// 仅移动调用对象包装器
	move_only_function<int(int, int)> add_50 = [ptr = std::make_unique<int>(50)](int x, int y) {return *ptr + x + y; };
	cout << add_50(10, 20) << endl;
	// 成员指针
	struct Add_to {
		Add_to(int base) : b(base) {}
		int Add(int v) { return b += v; }
	private:
		int b;
	};
	auto add_to = mem_fn(&Add_to::Add);
	auto u = make_unique<Add_to>(10);
	cout << add_to(u, 1000) << endl; // or add_to(Add_to(10), 1000)
	// 可复制(调用)对象的引用包装，与引用解包装
	using T = reference_wrapper<int>;
	static_assert(is_same_v<unwrap_reference_t<T>, int&>);
	using U = reference_wrapper<int&&>;
	static_assert(is_same_v<unwrap_reference_t<U>, int&>);
	static_assert(is_same_v<unwrap_ref_decay_t<U>, int&>);

	auto fn = sum;
	auto r_fn = std::cref(fn);  // 	reference_wrapper<const T(fn)>(fn);
	cout << r_fn(1, 2) << endl;

	int v = 10086;
	auto r_wrap = std::ref(v);  // reference_wrapper(v)
	cout << (r_wrap.get() += 10010, v) << endl;
	auto& r_unwrap = unwrap_reference_t<decltype(r_wrap)>(r_wrap);
	cout << (r_unwrap = 110, v) << endl;
	auto& r_unwrap_decay = unwrap_ref_decay_t<decltype(r_wrap)>(r_wrap);   // decay<unwrap_reference_t>
	cout << (r_unwrap_decay += 99, v) << endl;
}

// bind
void example_functional_bind() {
	std::cout << "\n=== std::bind Usage ===\n";
	using namespace std::placeholders;
	auto add_three = [](int a, int b, double c) { return  c + a + b; };
	auto fma = [](int a, int b, int c) {return a * b + c; };

	// Bind the first argument to 10
	auto add_to_10 = std::bind(add_three, 10, _1, _2);
	std::cout << "add_to_10(20, 30) = " << add_to_10(20, 30.0) << std::endl;
	// Bind the second argument to 50
	auto add_50 = std::bind(add_three, _1, 50, _2);
	std::cout << "add_50(10, 20) = " << add_50(10, 20) << std::endl;
	// Bind all arguments
	auto always_100 = std::bind(add_three, 25, 50, 25);
	std::cout << "always_100() = " << always_100() << std::endl;

	// Using bind with member functions
	struct Calculator {
		double multiply(double a, double b) { return a * b; }
		static int square(int x) { return x * x; }
	};
	Calculator calc;
	auto calc_multiply = std::bind(&Calculator::multiply, calc, _1, _2);
	std::cout << "calc_multiply(4.5, 2.0) = " << calc_multiply(4.5, 2.0) << std::endl;
	auto calc_square = std::bind(&Calculator::square, _1);
	std::cout << "calc_aquare(666) = " << calc_square(666) << std::endl;

	// bind_front, bind_after
	auto fma_10_at_front = std::bind_front(fma, 10, 20);
	std::cout << "fma[10*20+30] = " << fma_10_at_front(30) << std::endl;
	auto fma_100_at_back = std::bind_back(fma, 100, 200);
	std::cout << "fma[300*100+200]= " << fma_100_at_back(300) << std::endl;
}

// invoke
void example_functional_invoker() {
	std::cout << "\n=== functional invoker ===\n";

	struct Calculator {
		double multiply(double a, double b) { return a * b; }
		static int square(int x) { return x * x; }
	};
	auto lambda = [](int a, int b) { return a * b; };

	std::cout << std::invoke(&Calculator::multiply, Calculator{}, 10086, 10010) << std::endl;
	std::cout << std::invoke_r<int>(&Calculator::square, 99) << std::endl;
	std::cout << std::invoke(lambda, 5, 6) << std::endl;
}

// searchers
void example_functional_searchers() {
	std::cout << "\n=== std::functional String Searchers ===\n";

	// 测试文本和模式
	constexpr std::string_view text =
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed "
		"do eiusmod tempor incididunt ut labore et dolore magna aliqua";
	std::string_view pattern = "pisci";

	// default_searcher 是一个简单的搜索算法（通常是暴力匹配）
	std::cout << "\nUsing default_searcher: " << std::endl;
	auto start = std::chrono::high_resolution_clock::now();
	auto default_result = std::search(
		text.begin(), text.end(),
		std::default_searcher(pattern.begin(), pattern.end())
	);
	auto end = std::chrono::high_resolution_clock::now();
	auto time = std::chrono::duration_cast<std::chrono::microseconds>(end - start).count();
	std::cout << "default_searcher use " << (double)time / 1000 << " ms" << std::endl;
	if (default_result != text.end())
		std::cout << "Pattern found at position: " << std::distance(text.begin(), default_result) << std::endl;
	else std::cout << "Pattern not found" << std::endl;


	// boyer_moore_searcher 对较长模式有更好的性能
	std::cout << "\nUsing boyer_moore_searcher: " << std::endl;
	start = std::chrono::high_resolution_clock::now();
	auto bm_result = std::search(
		text.begin(), text.end(),
		std::boyer_moore_searcher(pattern.begin(), pattern.end())
	);
	end = std::chrono::high_resolution_clock::now();
	time = std::chrono::duration_cast<std::chrono::microseconds>(end - start).count();
	std::cout << "boyer_moore_searcher use " << (double)time / 1000 << " ms" << std::endl;
	if (bm_result != text.end())
		std::cout << "Pattern found at position: " << std::distance(text.begin(), bm_result) << std::endl;
	else std::cout << "Pattern not found" << std::endl;

	// boyer_moore_horspool_searcher 是 Boyer-Moore 的简化版本
	std::cout << "\nUsing boyer_moore_horspool_searcher: " << std::endl;
	start = std::chrono::high_resolution_clock::now();
	auto bmh_result = std::search(
		text.begin(), text.end(),
		std::boyer_moore_horspool_searcher(pattern.begin(), pattern.end())
	);
	end = std::chrono::high_resolution_clock::now();
	time = std::chrono::duration_cast<std::chrono::microseconds>(end - start).count();
	std::cout << "boyer_moore_horspool_searcher use " << (double)time / 1000 << " ms" << std::endl;
	if (bmh_result != text.end())
		std::cout << "Pattern found at position: " << std::distance(text.begin(), bmh_result) << std::endl;
	else std::cout << "Pattern not found" << std::endl;

	// 使用自定义字符类型和不同的模式
	std::wstring wtext = L"这是一个用于字符串搜索算法的示例文本。";
	std::wstring wpattern = L"搜索算法";

	std::cout << "\nSearching in wide string with boyer_moore_horspool_searcher: " << std::endl;
	auto wresult = std::search(
		wtext.begin(), wtext.end(),
		std::boyer_moore_horspool_searcher(wpattern.begin(), wpattern.end())
	);
	if (wresult != wtext.end())
		std::cout << "Wide string pattern found at position: " << std::distance(wtext.begin(), wresult) << std::endl;
	else std::cout << "Wide string pattern not found" << std::endl;

	// 在 vector 中搜索子序列
	std::vector<int> data = { 1, 2, 3, 4, 5, 6, 7, 8, 9 };
	std::vector<int> subseq = { 4, 5, 6 };

	std::cout << "\nSearching for subsequence in vector with default_searcher: " << std::endl;
	auto vec_result = std::search(
		data.begin(), data.end(),
		std::default_searcher(subseq.begin(), subseq.end())
	);
	if (vec_result != data.end()) 
		std::cout << "Subsequence found starting at index: " << std::distance(data.begin(), vec_result) << std::endl;
	else std::cout << "Subsequence not found" << std::endl;
}


void test_functional() {
	example_functional_basic();
	examples_functional_wrappers();
	example_functional_bind();
	example_functional_invoker();
	example_functional_searchers();
}