#include "test.h"

#include <stdalign.h>
#include <stddef.h>
#include <stdio.h>
#include <stdbool.h>


typedef struct
{
	int i;
	short s;
	bool b;
	long double ld;
} S1;
typedef struct
{
	int i;
	alignas(alignof(int) * 2) short s;
	bool b;
	long double ld;
} S2;
#pragma pack(push, 1)
typedef struct
{
	int i;
	alignas(alignof(int) * 2) short s;
	bool b;
	long double ld;
} S3;
#pragma pack(pop)

void test_stdalign(void)
{
#define check_mem_offset(T)            \
	printf("In alignas(%zu) " #T ":\n" \
		   "	sizeof = %zu\n"           \
		   "	offset(i) = %zu\n"        \
		   "	offset(s) = %zu\n"        \
		   "	offset(b) = %zu\n"        \
		   "	offset(ld) = %zu\n",   \
		   alignof(T), sizeof(T), offsetof(T, i), offsetof(T, s), offsetof(T, b), offsetof(T, ld))

	check_mem_offset(S1);
	check_mem_offset(S2);
	check_mem_offset(S3);
}