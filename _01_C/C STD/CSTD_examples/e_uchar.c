#include "test.h"

#include <stdio.h>
#include <stdlib.h>
#include <uchar.h>
#include <locale.h>
#include <string.h>

 // 演示 mbrtoc32/c32rtomb 转换
static void example_utf32_conversion(void)
{
	puts("\n[UTF-8 <-> UTF-32 Conversion]");

	char utf8_str[] = u8"🌍地球";  // 窄字符
	char32_t utf32_buf[10] = { 0 };
	char output_buf[20] = { 0 };

	// UTF-8 -> UTF-32 转换
	mbstate_t state = { 0 };
	const char* p = utf8_str;
	char32_t* p32 = utf32_buf;
	char32_t utf32_char = U'😊';        // UTF-32 字符字面量
	size_t step;
	while ((step = mbrtoc32(p32++, p, strlen(p) + 1, &state)) > 0)
		p += step;

	printf("UTF-8: %s\n", utf8_str);
	printf("UTF-32: ");
	for (int i = 0; i < p32 - utf32_buf - 1; i++)
		printf("U+%04X ", utf32_buf[i]);
	putchar('\n');

	// UTF-32 -> UTF-8 转换
	p32 = utf32_buf;
	char* out = output_buf;
	state = (mbstate_t){ 0 };

	while (*p32) {
		step = c32rtomb(out, *p32++, &state);
		out += step;
	}
	printf("Converted UTF-8: %s\n", output_buf);
}

// 演示 Unicode 码点属性
static void example_unicode_properties(void)
{
	puts("\n[Unicode Properties]");

	char32_t chars[] = { U'A', U'你', U'😊', U'𝄞' };
	char c[10] = { 0 };
	for (size_t i = 0; i < sizeof(chars) / sizeof(chars[0]); i++) {
		mbstate_t state = (mbstate_t){ 0 };
		(void)c32rtomb(c, chars[i], &state);
		printf("[U+%04X] %s properties:\n", chars[i], c);
		printf("  Is wide? %s\n", (chars[i] <= 0x7F) ? "true" : "false");
		printf("  UTF-8 length: %d\n", mblen(c, 10));
	}
}

void test_uchar(void)
{
	// 必须设置 UTF-8 本地化环境
	setlocale(LC_ALL, "en_us.UTF8");
	example_utf32_conversion();
	example_unicode_properties();
}