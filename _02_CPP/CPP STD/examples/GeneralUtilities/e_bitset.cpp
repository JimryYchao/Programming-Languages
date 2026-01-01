#include <iostream>
#include <bitset>
#include <string>
#include <algorithm>
#include <numeric>

module GeneralUtilities;
using namespace std;

// Basic bitset operations
void example_bitset_basic() {
	std::cout << "\n=== Basic std::bitset Operations ===\n";

	// Create bitsets
	std::bitset<8> b1;                  // Default constructor (all zeros)
	std::bitset<8> b2(42);              // From unsigned long long
	std::bitset<8> b3(std::string("10101010")); // From string

	std::cout << "b1: " << b1 << std::endl;
	std::cout << "b2 (42): " << b2 << std::endl;
	std::cout << "b3 (\"10101010\"): " << b3 << std::endl;

	// Bit access
	std::cout << "b3[0]: " << b3[0] << std::endl; 
	std::cout << "b3.test(7): " << boolalpha << b3.test(7) << dec << std::endl;

	// Set, reset, flip
	std::cout << "b1 after set(): " << b1.set() << std::endl;
	std::cout << "b1 after reset(): " << b1.reset() << std::endl;
	std::cout << "b1 after set(3): " << b1.set(3) << std::endl;
	std::cout << "b1 after flip(): " << b1.flip() << std::endl;
	std::cout << "b1 after flip(3): " << b1.flip(3) << std::endl;

	// Count bits
	std::cout << "b2 has " << b2.size() << " bits in total" << std::endl;
	std::cout << "b2 has " << b2.count() << " bits set" << std::endl;
	std::cout << "b2 is " << (b2.none() ? "all zeros" : "not all zeros") << std::endl;
	std::cout << "b2 is " << (b2.all() ? "all ones" : "not all ones") << std::endl;
	std::cout << "b2 is " << (b2.any() ? "has at least one 1" : "all zeros") << std::endl;
}

// Bitset bitwise operations
void example_bitset_bitwise() {
	std::cout << "\n=== std::bitset Bitwise Operations ===\n";

	std::bitset<8> b1(0b10101010); // 170
	std::bitset<8> b2(0b00111100); // 60

	std::cout << "b1: " << b1 << std::endl;
	std::cout << "b2: " << b2 << std::endl;

	// Bitwise 
	std::cout << "b1 & b2: " << (b1 & b2) << std::endl;
	std::cout << "b1 | b2: " << (b1 | b2) << std::endl;
	std::cout << "b1 ^ b2: " << (b1 ^ b2) << std::endl;
	std::cout << "~b1: " << (~b1) << std::endl;
	std::cout << "b1 << 2: " << (b1 << 2) << std::endl;
	std::cout << "b1 >> 2: " << (b1 >> 2) << std::endl;
	// Compound assignments
	b1 |= b2;
	std::cout << "b1 after |= b2: " << b1 << std::endl;
}

void test_bitset() {
	example_bitset_basic();
	example_bitset_bitwise();
}