#include <iostream>
#include <string>
#include <vector>
#include <algorithm>
#include <sstream>

module Strings;
using namespace std;

// Basic usage of std::string
void example_string_basic() {
	std::cout << "\n=== Basic std::string Usage ===\n";

	// String construction
	std::string str = "hello";
	std::u8string u8str = u8"hello";
	std::u16string u16str = u"hello";
	std::u32string u32str = U"hello";
	std::wstring wstr = L"hello";

	// String construction
	std::string s1;                       // Empty string
	std::string s2(5, 'a');               // 5 'a' characters: "aaaaa"
	std::string s3("Hello");              // From C-string
	std::string s4(s3, 1, 3);             // Substring: "ell"
	std::string s5({ 'w', 'o', 'r', 'l', 'd' }); // Initializer list
	std::string s6(str.begin(), str.end() - 4);

	// String size and capacity
	auto hi = s3 + "世界";  // 5 + 6
	std::cout << "hi empty: " << std::boolalpha << hi.empty() << std::endl;
	std::cout << "hi length: " << hi.length() << std::endl;  // == size
	std::cout << "hi capacity: " << hi.capacity() << std::endl;
	std::cout << "hi append '132456789', size: " << (hi.append("123456789"), hi.size()) << endl;
	std::cout << "hi capacity: " << hi.capacity() << std::endl;
	std::cout << "hi size after clear: " << (hi.clear(), hi.size()) << endl;
	std::cout << "hi capacity after clear: " << (hi.capacity()) << endl;
	std::cout << "hi capacity after shrink: " << (hi.shrink_to_fit(), hi.capacity()) << endl;
	std::cout << "hi capacity after reserve(2*cap): " << (hi.reserve(2 * hi.capacity()), hi.capacity()) << endl;
	std::cout << "hi max_size: " << hi.max_size() << std::endl;

	// Accessing characters
	hi = "Hello World";
	std::cout << "First character of hi: " << hi[0] << std::endl;
	std::cout << "Third character of hi: " << hi.at(2) << std::endl; // With bounds checking
	std::cout << "Front character of hi: " << hi.front() << std::endl;
	std::cout << "Back character of hi: " << hi.back() << std::endl;
	const char* c_str = hi.c_str();   // c-style-string
	std::cout << "hi as C-string: " << c_str << std::endl;
	auto cpp_str = "Hello WOLRD"s;

	// Assign
	cout << R"(assign("hi"): )" << hi.assign("hi") << endl;
	cout << R"(assign("This is a String", 4): )" << hi.assign("This is a String", 4) << endl;
	cout << R"(assign("This is a String", 4, 6): )" << hi.assign("This is a String", 4, 6) << endl;
	cout << R"(assign(5, 'a'): )" << hi.assign(5, 'a') << endl;
	cout << R"(assign({ 'w', 'o', 'r', 'l', 'd' }): )" << hi.assign({ 'w', 'o', 'r', 'l', 'd' }) << endl;
	string_view sv(cpp_str);
	cout << R"(assign(sv, 6): )" << hi.assign(sv, 6) << endl;
	cout << R"(assign(sv, 0, 5): )" << hi.assign(sv, 0, 5) << endl;
	cout << R"(assign(inputFst, inputLst): )" << hi.assign(begin(sv), begin(sv) + 2) << endl;

	// Modifier
	cout << "cpp_str: " << cpp_str << endl;
	cpp_str.resize(cpp_str.size() + 5, '!');  // fill with '!'
	cout << "cpp_str after resize(size+5, '!'): " << cpp_str << std::endl;
	cpp_str.resize(cpp_str.size() - 10);
	cout << "cpp_str after resize(size-10): " << cpp_str << endl;
	auto str_append = "resize_and_overwrite"s;
	cpp_str.resize_and_overwrite(cpp_str.size() + str_append.size(), [sz = cpp_str.size(), append = string_view(str_append)](char* p, size_t nsz) -> size_t {
		std::memcpy(p + sz, append.data(), append.size());
		return sz + append.size();
		});
	cout << "cpp_str after resize overwrite: " << cpp_str << endl;
	cout << "cpp_str insert ':' : " << cpp_str.insert(5, ":") << endl;
	cout << "cpp_str erase string behind ':' : " << cpp_str.erase(cpp_str.find(':') + 1) << endl;
	cout << "cpp_str pop_back : " << (cpp_str.pop_back(), cpp_str) << endl;
	cout << "cpp_str push_back ',' : " << (cpp_str.push_back(','), cpp_str) << endl;
	cout << "cpp_str append 'World' : " << (cpp_str.append("World"), cpp_str) << endl;
	cout << "cpp_str replace ',' to ' ' : " << (cpp_str.replace(cpp_str.find(','), 1, " "), cpp_str) << endl;

	// Find
	std::string str_find = "  \t Hello World!   \n\n";
	size_t start = str_find.find_first_not_of(" \t\v");
	size_t end = str_find.find_last_not_of(" \t\v\r\n");
	if (start != std::string::npos && end != std::string::npos) {
		str_find = str_find.substr(start, end - start + 1);
		std::cout << "str_find trimmed: " << str_find << std::endl;
	}
	if (str_find.contains("World"))
		std::cout << "Found \"World\" at position: " << str_find.find("World") << std::endl;
}

void test_string() {
	example_string_basic();
}