#include <iostream>
#include <stack>
#include <vector>
#include <string>
#include <forward_list>
#include <span>

module Containers;
using namespace std;

// 栈格式化打印工具
template<typename T, typename Container>
void print(const stack<T, Container>& s, const string& label) {
	cout << label << ": ";
	auto tmp = s;
	while (!tmp.empty()) {
		cout << tmp.top() << " ";
		tmp.pop();
	}
	cout << endl;
}

void example_stack() {
	cout << "\n=== std::stack ===\n";

	// container adapt
	stack s_dq(deque<int>{});
	stack s_vec(vector<int>{});
	stack s_str("Hello"sv);
	s_str.push(' ');

	stack<int>    ({ 1,2,3,4,5 });
	stack<string> strStack({ "Hello", "World", "C++" });
	stack s_string("string"s);
	stack s_span(span("span string_view"sv));

	stack s(s_vec);
	cout << "s size: " << s.size() << ", empty: " << boolalpha << s.empty() << endl;
	print(s, "s");

	s.top() = 999;
	print(s, "top = 999");

	s.push_range(vector{ 10,20,30,40,50 });
	print(s, "push range");

	s.emplace(666);
	print(s, "emplace 666");
}

// 平衡括号检查
void example_balanced_parentheses() {
	cout << "\n=== Balanced Parentheses Check ===\n";
	auto checkBalance = [](const string& expr) {
		stack<char> brackets;
		for (char c : expr) {
			if (c == '(' || c == '{' || c == '[') brackets.push(c);
			else if (c == ')' || c == '}' || c == ']') {
				if (brackets.empty()) return false;
				char top = brackets.top();
				brackets.pop();
				if ((c == ')' && top != '(') || (c == '}' && top != '{') || (c == ']' && top != '['))
					return false;
			}
		}
		return brackets.empty();
		};

	vector<string> exprs = { "()", "()[]{}", "(]", "([)]", "{[]}", "((())" };
	for (const auto& expr : exprs)
		cout << "\"" << expr << "\" is " << (checkBalance(expr) ? "balanced" : "unbalanced") << endl;
}

// 使用栈反转序列
void example_reverse_string() {
	cout << "\n=== reverse string by stack ===" << endl;
	string original = "Hello, World!";
	stack charStack(original);
	string reversed;
	reversed.reserve(original.length());
	while (!charStack.empty()) {
		reversed.push_back(charStack.top());
		charStack.pop();
	}
	cout << "original: " << original << endl;
	cout << "reversed: " << reversed << endl;
}

void test_stack() {
	example_stack();
	example_balanced_parentheses();
	example_reverse_string();
}