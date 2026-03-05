#include "Generator.hpp"

#include <iostream>
// 使用示例：斐波那契数列生成器
Generator<long long> fibonacci(int n) {
	long long a = 0, b = 1;
	for (int i = 0; i < n; ++i) {
		co_yield a;
		auto next = a + b;
		a = b;
		b = next;
	}
}
// 无限序列生成器
Generator<long long> infinite_fibonacci() {
	long long a = 0, b = 1;
	while (true) {
		co_yield a;
		auto next = a + b;
		a = b;
		b = next;
	}
}
int main() {
	std::cout << "前10个斐波那契数:" << std::endl;
	for (auto num : fibonacci(10)) {
		std::cout << num << " ";
	}
	std::cout << std::endl;
	std::cout << "\n无限序列的前15个:" << std::endl;
	auto gen = infinite_fibonacci();
	gen.next(); // 先让生成器产生第一个值
	for (int i = 0; i < 15; ++i) {
		std::cout << gen.current() << " ";
		gen.next();
	}
	std::cout << std::endl;
}