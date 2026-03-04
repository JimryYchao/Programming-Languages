#include "test.h"

#include <stdio.h>
#include <stdlib.h>  
#include <errno.h>   
#include <string.h>  
#include <locale.h>
#include <time.h>

// 演示环境变量操作
static void example_environment(void)
{
	puts("\n[Environment Variables]");
	// 获取环境变量
	const char* path = getenv("PATH");
	printf("PATH: %.40s...\n", path ? path : "(not found)");
	// 设置环境变量
	if (!_putenv("TEST_VAR=hello;"))
		printf("Set TEST_VAR: %s\n", getenv("TEST_VAR"));
}

// 演示伪随机数生成
static void example_random100(void)
{
	puts("\n[Random Numbers]");
	//srand(42);  // 固定种子保证可重复结果
	srand((unsigned)time(NULL));
	printf("Random values: ");
	for (int i = 0; i < 5; i++)
		printf("%d ", rand() % 100);
	putchar('\n');
}

// 演示程序终止
inline void atexit_func(void) {
	printf(">>> Cleanup handler executed <<<\n");
}
static void example_program_termination(void)
{
	puts("\n[Program Termination]");
	atexit(atexit_func); // 注册退出处理函数
	printf("Normal exit in progress...\n");
}

// 演示多字节/宽字符转换
static void example_mb_wc_Conversion(void)
{
	puts("\n[Multibyte/Wide Char Conversion]");
	// 必须设置本地化环境才能正确处理多字节字符
	setlocale(LC_ALL, "");

	/* wcstombs 宽字符串转多字节 */
	const wchar_t* wstr = L"中文ABC";
	char mbstr[100] = { 0 };
	size_t conv = wcstombs(mbstr, wstr, sizeof(mbstr));
	printf("wcs2mbs:\n\tWide: %ls\n\tMulti: %s (%zu bytes)\n", wstr, mbstr, conv);

	/* mbstowcs 多字节字符串转宽字符 */
	const char* mbstr2 = "日本語123";
	wchar_t wstr2[100] = { 0 };
	conv = mbstowcs(wstr2, mbstr2, sizeof(wstr2) / sizeof(wchar_t));
	printf("mbs2wcs: \n\tMulti: %s\n\tWide: %ls (%zu chars)\n", mbstr2, wstr2, conv);

	/* mbtowc 单个多字节字符转宽字符 */
	char mbchar[MB_LEN_MAX] = "€"; // 欧元符号(3字节UTF-8)
	wchar_t wc;
	int bytes = mbtowc(&wc, mbchar, MB_CUR_MAX);
	printf("mb2wc:\n\tMulti: %s → U+%04X (%lc), %d bytes\n", mbchar, wc, wc, bytes);

	/* wctomb 单个宽字符转多字节 */
	wchar_t wc2 = L'ß'; // 德语sharp-s
	char mbchar2[MB_LEN_MAX] = { 0 };
	bytes = wctomb(mbchar2, wc2);
	printf("wc2mb:\n\tWide: U+%04X (%lc) → Multi: ", wc2, wc2);
	for (int i = 0; i < bytes; i++)
		printf("%02X ", (unsigned char)mbchar2[i]);
	putchar('\n');
}

// 演示整数数组查找
static int compare_ints(const void* a, const void* b)  // 整数比较函数（升序）
{
	int arg1 = *(const int*)a;
	int arg2 = *(const int*)b;
	return (arg1 > arg2) - (arg1 < arg2);
}

static void example_bsearch(void)
{
	puts("\n[Integer Array Search]");
	int arr[] = { 10, 20, 30, 40, 50, 60 };
	int count = sizeof(arr) / sizeof(int);
	int key = 40;
	void* result = bsearch(&key, arr, count, sizeof(int), compare_ints);
	if (result)
		printf("Found %d at index %lld\n", key, (int*)result - arr);
	else printf("%d not found\n", key);

	// 字符串查找
	puts("\n[String Array Search]");
	const char* strs[] = { "apple", "banana", "grape", "orange", "pear" };
	int count = sizeof(strs) / sizeof(strs[0]);
	const char* key = "grape";

	// 按字典序查找
	void* result = bsearch(&key, strs, count, sizeof(char*), strcmp);
	if (result) printf("Found '%s' at index %lld\n", *(char**)result, (char**)result - strs);
	else printf("'%s' not found\n", key);
}

// 演示结构体查找
typedef struct {
	int id;
	char name[20];
} Person;
static int compare_struct_id(const void* a, const void* b) {  // 按 Person.id 查找
	int id1 = ((const Person*)a)->id;
	int id2 = ((const Person*)b)->id;
	return (id1 > id2) - (id1 < id2);
};
static void example_struct_bsearch(void)
{
	puts("\n[Structure Array Search]");
	Person people[] = {
		{101, "Alice"},
		{203, "Bob"},
		{305, "Charlie"}
	};
	int count = sizeof(people) / sizeof(people[0]);
	Person key = { .id = 203 };

	void* result = bsearch(&key, people, count, sizeof(Person), compare_struct_id);
	if (result)
		printf("Found ID %d: %s\n", ((Person*)result)->id, ((Person*)result)->name);
}

void test_stdlib(void)
{
	example_environment();
	example_random100();
	example_program_termination();
	example_mb_wc_Conversion();
	// bearch
	example_bsearch();
	example_struct_bsearch();
}
/*

*/