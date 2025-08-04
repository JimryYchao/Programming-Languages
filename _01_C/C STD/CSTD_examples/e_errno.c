#include "test.h"

#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include <math.h>
#include <threads.h> 


// 演示文件操作错误
static void example_file_errors(void)
{
	errno = 0;
	puts("\n[File Operation Errors]");
	FILE* fp = fopen("non_existent_file.txt", "r");
	if (fp == NULL) {
		printf("errorCode [%d]: %s\n", errno, strerror(errno));
		perror("ERROR");
	}
}

// 演示数学运算错误
static void example_math_errors(void)
{
	errno = 0;
	puts("\n[Math Domain Errors]");
	double result = sqrt(-1.0);    // 计算负数的平方根
	if (errno == EDOM)
		printf("Domain error occurred: %s\n", strerror(errno));
	else if (errno != 0)
		perror("ERROR");
}

// 演示内存分配错误
static void example_memory_errors(void)
{
	errno = 0;
	puts("\n[Memory Allocation Errors]");
	void* ptr = malloc(SIZE_MAX);	// 尝试超大内存分配
	if (ptr == NULL) {
		printf("Allocation failed: %s\n", strerror(errno));
		if (errno == ENOMEM)
			puts("System out of memory");
	}
}

// 演示线程局部 errno
static int thread_func(void* arg)
{
	(void)arg;
	errno = EAGAIN;
	printf("Thread errorCode [%d]: %s\n", errno, strerror(errno));
	return 0;
}

void test_errno(void)
{
	example_file_errors();
	example_math_errors();
	example_memory_errors();

	// 线程局部 errno
	puts("\n[Thread-Local errno]");
	thrd_t thr;
	if (thrd_create(&thr, thread_func, NULL) == thrd_success) {
		thrd_join(thr, NULL);
		printf("Main thread errorCode [%d]: %s\n", errno, strerror(errno)); // 应该是 0
	}
}