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
/*
str(123456789) toimax by base(0) = 123456789
str(1) toimax by base(2) = 1
str(123) toimax by base(4) = 27
str(1234567) toimax by base(8) = 342391
str(123456789) toimax by base(10) = 123456789
str(123456789abcdef) toimax by base(16) = 81985529216486895
str(123456789abcdefgh) toimax by base(32) = 9223372036854775807
*/