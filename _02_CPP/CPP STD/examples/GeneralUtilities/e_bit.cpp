#include <iostream>
#include <bit>
#include <cstdint>
#include <bitset>
#include <string>

module GeneralUtilities;
using namespace std;

// 二进制位 0 1 计数
void example_bit_count() {
	std::cout << "\n=== count zero or one ===\n";

	// Count leading zeros
	uint8_t value8 = 0b110000;  // 48 in decimal
	uint16_t value16 = 0b110000;  // 48 in decimal

	std::cout << "value8: 0b" << std::bitset<8>(value8) << std::endl;
	std::cout << "Count leading zeros: " << std::countl_zero(value8) << std::endl;

	std::cout << "value16: 0b" << std::bitset<16>(value16) << std::endl;
	std::cout << "Count leading zeros: " << std::countl_zero(value16) << std::endl;

	// Count trailing zeros
	uint8_t value8_t = 0b00110000;  // 48 in decimal
	std::cout << "value8_t: 0b" << std::bitset<8>(value8_t) << std::endl;
	std::cout << "Count trailing zeros: " << std::countr_zero(value8_t) << std::endl;

	// Count set bits
	uint8_t value8_c = 0b00110101;  // 53 in decimal
	std::cout << "value8_c: 0b" << std::bitset<8>(value8_c) << std::endl;
	std::cout << "Count set bits: " << std::popcount(value8_c) << std::endl;

	// Find first set bit
	uint8_t value8_f = 0b00110000;  // 48 in decimal
	std::cout << "value8_f: 0b" << std::bitset<8>(value8_f) << std::endl;
	std::cout << "Find first set bit: " << std::countr_zero(value8_f) << " (from least significant bit)" << std::endl;
}

// 按位旋转
void example_bit_rotate() {
	std::cout << "\n=== bit rotate ===\n";

	uint8_t value = 0b10110001;
	std::cout << "Original value: 0b" << std::bitset<8>(value) << std::endl;

	// Rotate left
	auto rotated_left = std::rotl(value, 1);
	std::cout << "After rotate left by 1: 0b" << std::bitset<8>(rotated_left) << std::endl;
	rotated_left = std::rotl(value, 3);
	std::cout << "After rotate left by 3: 0b" << std::bitset<8>(rotated_left) << std::endl;

	// Rotate right
	auto rotated_right = std::rotr(value, 1);
	std::cout << "After rotate right by 1: 0b" << std::bitset<8>(rotated_right) << std::endl;
	rotated_right = std::rotr(value, 3);
	std::cout << "After rotate right by 3: 0b" << std::bitset<8>(rotated_right) << std::endl;

	// 按字节反转
	uint32_t value32 = 0x12345678;
	std::cout << "Original value32: 0x" << std::hex << value32 << std::dec << std::endl;
	auto swapped = std::byteswap(value32);   // 0x78563412
	std::cout << "After byteswap: 0x" << std::hex << swapped << std::dec << std::endl;
}

// 2 的整数幂
void example_bit_power_of_two() {
	std::cout << "\n=== power of two ===\n";

	std::cout << "Is 1024 a power of two? " << std::boolalpha << has_single_bit(1024u) << std::endl;
	std::cout << "Is 9999 a power of two? " << std::has_single_bit(9999u) << std::endl;

	// 查找最接近的 2 的整数幂
	std::cout << "Next power of two after 9999: " << std::bit_ceil(9999u) << std::endl;
	std::cout << "Before power of two after 1111: " << std::bit_floor(1111u) << std::endl;

	// 存储值至少需要的位数
	std::cout << "the number of bits at least needed to store 10086: " << std::bit_width(10086u) << std::endl;
	std::cout << "value(10086): 0b" << std::bitset<bit_ceil((unsigned)bit_width(10086u))>(10086u) << endl;
}

void test_bit() {
	example_bit_count();
	example_bit_rotate();
	example_bit_power_of_two();
}