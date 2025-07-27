module;

using namespace std;
#include <iostream>

module myModule:subModuleA;


void moduleA::FuncA() {
	cout << "call moduleA::FuncA\n";
}
void moduleA::FuncB() {
	cout << "call moduleA::FuncB\n";
}

