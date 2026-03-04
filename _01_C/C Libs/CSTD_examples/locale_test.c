#include "test.h"

#include <stdio.h>
#include <locale.h>  
#include <time.h>    
#include <wchar.h>   
#include <string.h>

// 演示数字格式的区域化差异
static void example_numeric_format(void)
{
	// 获取当前区域设置
	char old_lc[128] = { 0 };
	strcpy_s(old_lc, sizeof(old_lc), setlocale(LC_NUMERIC, NULL));
	char* old = setlocale(LC_NUMERIC, NULL);
	printf("Current numeric locale: %s\n", old_lc);

	// 测试不同区域的数字格式
	setlocale(LC_NUMERIC, "C");
	printf("[C locale] Decimal point: '%s'\n", localeconv()->decimal_point);
	setlocale(LC_NUMERIC, "fr_FR.UTF-8");
	printf("[French locale] Decimal point: '%s'\n", localeconv()->decimal_point);

	// 恢复原始设置
	setlocale(LC_NUMERIC, old_lc);
}

// 演示时间格式的区域化展示
static void example_time_format(void)
{
	time_t now = time(NULL);
	struct tm* local = localtime(&now);
	char buffer[100];

	setlocale(LC_TIME, "en_US.UTF-8");
	strftime(buffer, sizeof(buffer), "%A %c", local);
	printf("\nEnglish time format: %s\n", buffer);

	setlocale(LC_TIME, "zh_CN.UTF-8");
	strftime(buffer, sizeof(buffer), "%A %c", local);
	printf("Chinese time format: %s\n", buffer);

	setlocale(LC_TIME, "ja_jp.utf8");
	strftime(buffer, sizeof(buffer), "%A %c", local);
	printf("Japanese time format: %s\n", buffer);
}

void test_locale(void)
{
	// 初始化为本地区域设置
	setlocale(LC_ALL, "");  // "C" 为默认行为
	example_numeric_format();
	example_time_format();
}
/*
Current numeric locale: Chinese (Simplified)_China.utf8
[C locale] Decimal point: '.'
[French locale] Decimal point: ','

English time format: Sunday 1/11/2026 3:52:10 PM
Chinese time format: 星期日 2026/1/11 15:52:10
Japanese time format: 日曜日 2026/01/11 15:52:10
*/