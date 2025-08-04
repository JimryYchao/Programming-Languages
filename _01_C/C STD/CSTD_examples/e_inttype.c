#include "test.h"
#include <stddef.h>
#include <inttypes.h>
#include <stdio.h>

void example_str2Int() {
	const char* num = "123456789abcdefgh";
	char buf[20] = { 0 };
	char* endptr = NULL;
	intmax_t i = 0;
#define _strtoimax(base) \
i = strtoimax(num, &endptr, base);  \
printf("str(%.*s) toimax by base(%d) = %lld\n", endptr - num, num, base, i);
	_strtoimax(0);
	_strtoimax(2);
	_strtoimax(4);
	_strtoimax(8);
	_strtoimax(10);
	_strtoimax(16);
	_strtoimax(32);

}

void test_inttype(void) {
	example_str2Int();
}