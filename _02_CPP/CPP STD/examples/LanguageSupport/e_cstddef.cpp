#include <cstddef>
#include <iostream>

module LanguageSupport;
using namespace std;

// 演示 byte 原始内存操作
void example_byte() {
	cout << "\n[Original byte representation]\n";

	std::byte buffer[sizeof(unsigned)];
	unsigned value = 0xDEADBEEF;
	// 将整数转换为字节表示
	for (size_t i = 0; i < sizeof(unsigned); ++i)
		buffer[i] = static_cast<std::byte>((value >> (8 * i)) & 0xFF);
	std::cout << "Value in bytes: ";
	for (std::byte b : buffer)
		std::cout << std::hex << std::to_integer<int>(b) << " ";
	std::cout << "\n";

	// 从字节重建整数
	unsigned reconstructed = 0;
	for (size_t i = 0; i < sizeof(unsigned); ++i)
		reconstructed |= std::to_integer<unsigned>(buffer[i]) << (8 * i);
	std::cout << "Reconstructed value: 0x" << std::hex << reconstructed << "\n";
}

// 演示结构成员的偏移
struct ExampleStruct {
	int x;
	double y;
	char z;
};
void example_offsetof() {
	cout << "\n[offset of struct]\n";
	std::cout << "Offset of x: " << offsetof(ExampleStruct, x) << "\n";
	std::cout << "Offset of y: " << offsetof(ExampleStruct, y) << "\n";
	std::cout << "Offset of z: " << offsetof(ExampleStruct, z) << "\n";
	// 计算结构体总大小
	size_t total_size = sizeof(ExampleStruct);
	std::cout << "Total struct size: " << total_size << " bytes\n";
	std::cout << "Align of struct: " << offsetof(ExampleStruct, z) << "\n";
}

void test_cstddef(void) {
	example_byte();
	example_offsetof();
}