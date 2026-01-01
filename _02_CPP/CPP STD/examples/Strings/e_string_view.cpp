#include <iostream>
#include <string>
#include <string_view>
#include <vector>
#include <algorithm>

module Strings;
using namespace std;

// std::string_view
void example_string_view() {
	std::cout << "\n=== std::string_view ===\n";
	auto str = "     HELLO, WORLD!     "s;
	auto sv = string_view(str);
	cout << "str: " << str << endl;
	cout << "sv : " << sv << endl;
	cout << "sv remove_prefix(5): " << (sv.remove_prefix(5), sv) << endl;
	cout << "sv remove_suffix(5): " << (sv.remove_suffix(5), sv) << endl;
	for (auto& c : str)
		c = std::tolower(c);

	cout << "str tolower: " << str << endl;
	cout << "sv: " << str << endl;
	cout << sv << endl;
}

void test_string_view() {
	example_string_view();
}