#include <iostream>
#include <ratio>
#include <chrono>
#include <typeinfo>

module Metaprogramming;
using namespace std;

// Ratio operations
void example_ratio_operations() {
	std::cout << "\n=== Ratio Operations ===\n";

	// Ratio addition
	using one_half = std::ratio<1, 2>;
	using one_third = std::ratio<1, 3>;
	using sum = std::ratio_add<one_half, one_third>;  // 1/2 + 1/3 = 5/6
	std::cout << "1/2 + 1/3 = " << sum::num << "/" << sum::den << std::endl;

	// Ratio subtraction
	using diff = std::ratio_subtract<one_half, one_third>;  // 1/2 - 1/3 = 1/6
	std::cout << "1/2 - 1/3 = " << diff::num << "/" << diff::den << std::endl;

	// Ratio multiplication
	using product = std::ratio_multiply<one_half, one_third>;  // 1/2 * 1/3 = 1/6
	std::cout << "1/2 * 1/3 = " << product::num << "/" << product::den << std::endl;

	// Ratio division
	using quotient = std::ratio_divide<one_half, one_third>;  // 1/2 / 1/3 = 3/2
	std::cout << "1/2 / 1/3 = " << quotient::num << "/" << quotient::den << std::endl;

	// Ratio comparison
	cout << boolalpha;
	std::cout << "1/2 > 1/3: " << std::ratio_greater<one_half, one_third>::value << std::endl;
	std::cout << "1/2 == 1/3: " << std::ratio_equal<one_half, one_third>::value << std::endl;
	std::cout << "1/2 < 1/3: " << std::ratio_less<one_half, one_third>::value << std::endl;
}

// Ratio application in time units
void example_ratio_in_time() {
	std::cout << "\n=== Ratio Application in Time Units ===\n";

	// Use duration types from std::chrono library, which are based on std::ratio
	using seconds = std::chrono::duration<double>;                   // Default is seconds
	using milliseconds = std::chrono::duration<double, std::milli>;  // Milliseconds
	using microseconds = std::chrono::duration<double, std::micro>;  // Microseconds
	using nanoseconds = std::chrono::duration<double, std::nano>;    // Nanoseconds
	using minutes = std::chrono::duration < double, std::ratio<60 >>;// Minutes
	using hours = std::chrono::duration<double, std::ratio<3600>>;   // Hours

	// Conversion examples
	seconds sec(1);
	milliseconds ms = sec;  // 1 second = 1000 milliseconds
	microseconds us = sec;  // 1 second = 1000000 microseconds
	nanoseconds ns = sec;   // 1 second = 1000000000 nanoseconds
	minutes min = sec / 60; // 1 second = 1/60 minute
	hours hr = sec / 3600;  // 1 second = 1/3600 hour

	std::cout << "1 second = " << ms.count() << " milliseconds" << std::endl;
	std::cout << "1 second = " << us.count() << " microseconds" << std::endl;
	std::cout << "1 second = " << ns.count() << " nanoseconds" << std::endl;
	std::cout << "1 second = " << min.count() << " minutes" << std::endl;
	std::cout << "1 second = " << hr.count() << " hours" << std::endl;

	// Operation example
	seconds total = hours(2) + minutes(30) + seconds(15);
	std::cout << "2 hours 30 minutes 15 seconds = " << total.count() << " seconds" << std::endl;
}

// Ratio application in physical units
using meter = std::ratio<1, 1>;				 // Meter
using kilometer = std::ratio<1000, 1>;		 // Kilometer
using centimeter = std::ratio<1, 100>;		 // Centimeter
using millimeter = std::ratio<1, 1000>;		 // Millimeter
using micrometer = std::ratio<1, 1000000>;   // Micrometer
template<typename T, typename R = meter>
struct Length {
	T value;
	explicit Length(T v) : value(v) {}
	template<typename R2>
	Length<T, R2> convert() const {
		using Ratio = std::ratio_divide<R, R2>;
		return Length<T, R2>(static_cast<T>(value * Ratio::num / static_cast<double>(Ratio::den)));
	}
};
void example_ratio_in_physics() {
	std::cout << "\n=== Ratio Application in Physical Units ===\n";
	Length<double, kilometer> distance(1.5);  // 1.5 kilometers
	std::cout << "Distance: " << distance.value << " kilometers" << std::endl;
	// Convert to other units
	auto meters = distance.convert<meter>();
	auto centimeters = distance.convert<centimeter>();

	std::cout << "Converted to meters: " << meters.value << " meters" << std::endl;
	std::cout << "Converted to centimeters: " << centimeters.value << " centimeters" << std::endl;
}

void test_ratio() {
	example_ratio_operations();
	example_ratio_in_time();
	example_ratio_in_physics();
}