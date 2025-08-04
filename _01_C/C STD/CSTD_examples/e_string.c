#include "test.h"

#include <stdio.h>
#include <string.h>
#include <stdlib.h>

// 演示基础字符串操作
static void example_string_basics(void)
{
	puts("\n[Basic String Operations]");

	char src[50] = "Hello";
	char dest[50] = { 0 };

	// 字符串复制
	strcpy_s(dest, 50, src);
	printf("strcpy_s: %s\n", dest);

	// 字符串复制
	const char* source = " World!";
	char* end = &dest[strlen(dest)];
	(void)strncpy_s(end, sizeof(dest) - (end - dest) * sizeof(char), source, sizeof(source));
	printf("strncpy_s: %s\n", dest);

	// 字符串连接
	(void)strcat_s(dest, sizeof(dest), " Welcome");
	printf("strcat: %s\n", dest);

	// 字符串比较
	int cmp = strcmp("apple", "banana");
	printf("strcmp: 'apple' %s 'banana'\n", cmp < 0 ? "<" : ">");
}

// 演示内存操作函数
static void example_memory_operations(void)
{
	puts("\n[Memory Operations]");

	int arr1[5] = { 1, 2, 3, 4, 5 };
	int arr2[5] = { 0 };

	// 内存复制
	memcpy(arr2, arr1, sizeof(arr1));
	printf("memcpy: arr2[2] = %d\n", arr2[2]);

	// 内存移动（处理重叠区域）
	memmove(arr1 + 2, arr1, 3 * sizeof(int));
	printf("memmove: arr1[3] = %d\n", arr1[3]);

	// 内存设置
	memset(arr1, 0, sizeof(arr1));
	printf("memset: arr1[3] = %d\n", arr1[0]);
}

// 演示字符串搜索与分割
static void example_string_search(void)
{
	puts("\n[String Searching]");

	char str[] = "The quick brown fox jumps over the lazy dog";
	printf("Parsing the str: '%s'\n", str);

	// 字符搜索
	char* p = strchr(str, 'q');
	printf("strchr: 'q' found at position %zd\n", p - str);

	// 字符串搜索
	p = strstr(str, "fox");
	printf("strstr: 'fox' found at position %zd\n", p - str);

	// 内存内容搜索
	int numbers[] = { 10, 20, 30, 40, 50 };
	int* num = memchr(numbers, 30, sizeof(numbers));
	printf("memchr: Found %d at offset %zd\n", *num, (char*)num - (char*)numbers);

	//char buf[16];
	char* next_token = NULL;
	char* token = strtok_s(str, " ", &next_token);
	while (token) {
		puts(token);
		token = strtok_s(NULL, " ", &next_token);
	}
	printf("Contents of the str now: '");
	for (size_t n = 0; n < sizeof str; ++n)
		str[n] ? putchar(str[n]) : fputs("\\0", stdout);
	puts("'");
}


// 演示 memcmp 和排序
typedef struct {
	int id;
	char name[20];
}Item;
static int sort_f(const void* a, const void* b) {
	return ((const Item*)a)->id - ((const Item*)b)->id;
}
static void example_sorting(void)
{
	puts("\n[Memory Comparison]");

	Item items[] = {
			{3, "Apple"},
			{1, "Banana"},
			{4, "Orange"},
			{2, "Cherry"},
	};

	// 排序比较函数
	int eleCount = sizeof(items) / sizeof(Item);
	qsort(items, eleCount, sizeof(Item), sort_f);
	printf("Sorted items:\n");
	for (int i = 0; i < eleCount; i++)
		printf("%d: %s\n", items[i].id, items[i].name);
}

void test_string(void)
{
	example_string_basics();
	example_memory_operations();
	example_string_search();
	example_sorting();
}