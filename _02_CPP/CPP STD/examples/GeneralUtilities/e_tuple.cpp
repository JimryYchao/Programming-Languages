#include <iostream>
#include <tuple>
#include <string>
#include <array>
#include <algorithm>
#include <functional>
#include <memory>

module GeneralUtilities;
using namespace std;

// std::tuple 
void example_tuple() {
	std::cout << "\n=== std::tuple Usage ===\n";
	// make tuple
	std::tuple<int, std::string, double> my_tuple(42, "Hello", 3.14);
	using T = decltype(my_tuple);
	auto another_tuple = std::make_tuple(5, "World", 2.71);
	// access elem
	std::cout << "Elements: " << std::get<0>(my_tuple) << ", " << std::get<1>(my_tuple) << ", " << std::get<2>(my_tuple) << std::endl;
	std::cout << "Tuple size: " << tuple_size<T>::value << std::endl;
	tuple_element_t<0, T> e0;       // int 
	tuple_element_t<1, T> e1;       // string
	std::tie(e0, e1, ignore) = my_tuple;
	cout << "tie tuple: " << e0 << ", " << e1 << endl;
	std::get<0>(my_tuple) = 10086;
	std::cout << "Modified first element: " << std::get<0>(my_tuple) << e0 << std::endl;

	// tuple cat
	auto new_tuple = tuple_cat(my_tuple, another_tuple);
	std::cout << "New tuple size: " << std::tuple_size_v<decltype(new_tuple)> << std::endl;

	// forward_as_tuple
	auto ref_tuple = forward_as_tuple(e0, e1);
	std::get<0>(ref_tuple) = 10010;
	cout << "ref tuple<0>" << e0 << endl;

	// make obj from tuple
	struct S {
		S(int fst, string& sec) { cout << "make S : " << (fst += 100) << ", " << (sec += " World") << endl; }
	};
	auto s = make_from_tuple<S>(std::move(ref_tuple));
	cout << e1 << endl;

	// common_type_t of tuples
	using Tuple = common_type_t<tuple<int, char*, float>, tuple<long, string, double>>;
	auto com_tuple = Tuple(1, "Hello", 3.1415);
	std::cout << "Elements: " << std::get<0>(com_tuple) << ", " << std::get<1>(com_tuple) << ", " << std::get<2>(com_tuple) << std::endl;

}

// Advanced tuple operations
void example_tuple_operations() {
	std::cout << "\n=== tuple Operations ===\n";

	auto my_tuple = std::make_tuple(1, "apple", 3.14);
	auto [x, y, z] = my_tuple;
	std::cout << "Structured bindings: x = " << x << ", y = " << y << ", z = " << z << std::endl;

	// Tie values
	int a = 10; std::string b = "banana";
	auto tied_tuple = std::tie(a, b);
	std::get<0>(tied_tuple) = 100;
	std::cout << "After tie modification: a = " << a << std::endl;

	// Apply function to tuple elements
	auto sum_tuple = std::make_tuple(1, 2, 3, 4, 5);
	auto sum = std::apply([](auto&&... args) { return (args + ...); }, sum_tuple);
	std::cout << "Sum of tuple elements: " << sum << std::endl;
}


void test_tuple() {
	example_tuple();
	example_tuple_operations();
}