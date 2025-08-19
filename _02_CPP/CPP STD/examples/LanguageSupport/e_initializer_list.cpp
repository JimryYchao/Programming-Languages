#include <initializer_list>
#include <iostream>
#include <vector>
#include <algorithm>

module LanguageSupport;
using namespace std;

// 自定义类使用 initializer_list
void example_custom_class() {
	cout << "\n[Custom class with initializer_list]\n";

	class ShoppingList {
		std::vector<std::string> items;
	public:
		ShoppingList(std::initializer_list<std::string> init) : items(init) {}
		void print() const {
			std::cout << "Shopping List:\n";
			for (const auto& item : items) {
				std::cout << "- " << item << "\n";
			}
		}
		void addList(std::initializer_list<std::string> list) {
			items.insert(items.end(), list);
		}
	};
	ShoppingList myList{ "Milk", "Eggs", "Bread", "Cheese" };  // 调用 initializer_list 构造函数
	myList.addList({ "Oranges", "Bananas", "Fish" });
	myList.print();
}

// 使用 initializer_list 与算法
void example_with_algorithms() {
	cout << "\n[initializer_list with algorithms]\n";

	// 使用 initializer_list 直接调用算法
	auto max_num = std::max({ 10, 30, 20, 50, 40 });
	std::cout << "Max number: " << max_num << "\n";

	// 排序 initializer_list 内容
	std::initializer_list<std::string> words = { "banana", "apple", "cherry" };
	std::vector sort_words(words);
	std::sort(sort_words.begin(), sort_words.end());
	std::cout << "Sorted words: ";
	for (const auto& word : sort_words)
		std::cout << word << " ";
	std::cout << "\n";
}

void test_initializer_list(void) {
	example_custom_class();
	example_with_algorithms();
}