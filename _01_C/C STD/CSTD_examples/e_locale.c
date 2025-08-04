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
	strcpy(old_lc, setlocale(LC_NUMERIC, NULL));
	char* old_locale = setlocale(LC_NUMERIC, NULL);
	printf("Current numeric locale: %s\n", old_lc);

	// 测试不同区域的数字格式
	setlocale(LC_NUMERIC, "C");
	printf("[C locale] Decimal point: '%s'\n", localeconv()->decimal_point);
	setlocale(LC_NUMERIC, "fr_FR.UTF-8");
	printf("[French locale] Decimal point: '%s'\n", localeconv()->decimal_point);

	// 恢复原始设置
	setlocale(LC_NUMERIC, old_lc);
}

// 演示货币格式的区域化特性
static void example_currency_format(void)
{
	char old_lc[128] = { 0 };
	strcpy(old_lc, setlocale(LC_MONETARY, NULL));
	struct lconv* lc = NULL; 

#define currency(loc)   \
	setlocale(LC_MONETARY, (loc));	\
	lc = localeconv();				\
	printf("\nCurrent monetary locale: %s\n", setlocale(LC_MONETARY, NULL));\
	printf("Currency symbols:\n");\
	printf("International: %s\n", lc->int_curr_symbol);\
	printf("Local: %s\n", lc->currency_symbol);\
	printf("Positive format: %d\n", lc->p_sign_posn);

	currency(old_lc);
	currency("C");
	currency("fr_FR.UTF-8");
	currency("cu-RU.UTF8");
	currency("en_us.UTF8");
	currency("ja_jp.utf8");
	setlocale(LC_MONETARY, old_lc);
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
	setlocale(LC_ALL, "");  // 初始化为本地区域设置

	example_numeric_format();
	example_currency_format();
	example_time_format();
}