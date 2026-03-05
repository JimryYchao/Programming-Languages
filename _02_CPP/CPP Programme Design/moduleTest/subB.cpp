module;
using namespace std;
#include <iostream>

module myModule:subModuleB;


void moduleB::FuncA() {
	cout << "call moduleB::FuncA\n";
}
void moduleB::FuncB() {
	cout << "call moduleB::FuncB\n";
}

