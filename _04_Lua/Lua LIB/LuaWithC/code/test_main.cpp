extern "C" {
#include "test.h"
}

int main() {
	test_BitArray();
	test_C_Call_Lua();
	test_Lua_Call_C();
	test_lproc();
	test_Continue_pcall();
	test_Continue_CFunc();
	test_memlimit();
}