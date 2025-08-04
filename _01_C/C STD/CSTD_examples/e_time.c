#include "test.h"

#include <stdio.h>
#include <time.h>
#include <stdint.h>
#include <threads.h>

static char time_buf[100];

// 演示基础时间获取
static void example_basic_time(void)
{
	puts("\n[Basic Time Operations]");
	time_t now = time(NULL);
	ctime_s(time_buf, sizeof(time_buf), &now);
	printf("Epoch seconds: %lld\n", (long long)now);
	printf("Local time: %s\n", time_buf);

	// 转换为 UTC 时间
	struct tm* utc = gmtime(&now);
	asctime_s(time_buf, sizeof(time_buf), utc);
	printf("UTC time: %s\n", time_buf);

	// 计算时间差
	struct tm target = {
		.tm_year = 125, // 2025 - 1900
		.tm_mon = 11,   // December
		.tm_mday = 31
	};
	time_t new_year = mktime(&target);
	double diff = difftime(new_year, now);
	printf("Seconds until 2026: %.0f\n", diff);
}

// 演示自定义时间格式
static void example_custom_format(void)
{
	puts("\n[Custom Time Formatting]");

	time_t now = time(NULL);
	struct tm* tm = localtime(&now);   // 转换本地区域 tm 

	// 自定义格式
	strftime(time_buf, sizeof(time_buf), "now: %Y-%m-%d %H:%M:%S", tm);
	puts(time_buf);

	// ISO 8601 格式
	strftime(time_buf, sizeof(time_buf), "ISO 8601: %FT%T%z", tm);
	puts(time_buf);
}

// 测量代码执行时间
static void example_exec_time(void)
{
	puts("\n[Measure code execution time]");
	struct timespec start, end;
	(void)timespec_get(&start, TIME_UTC);

	// 计算操作耗时
	volatile double x = 1.0;
	for (int i = 0; i < 1000000; i++)
		x *= 1.000001;
	(void)timespec_get(&end, TIME_UTC);
	long duration = (end.tv_sec - start.tv_sec) * 1000000000L +
		(end.tv_nsec - start.tv_nsec);
	printf("Loop took: %ld ns (x=%.2f)\n", duration, x);
}

void test_time(void)
{
	// 处理器当前执行时间
	clock_t t1 = clock();

	example_basic_time();
	example_custom_format();
	example_exec_time();

	thrd_sleep(&(struct timespec) { .tv_sec = 1 }, NULL);
	clock_t t2 = clock();
	printf("CPU time used: %0.3f ms\n", 1000.0 * (t2 - t1) / CLOCKS_PER_SEC);
}