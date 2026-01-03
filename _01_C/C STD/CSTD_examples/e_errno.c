#include "test.h"

#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>
#include <math.h>
#include <threads.h> 


static void example_errors(void)
{
	errno = 0;
	// file operation errors
	FILE* fp = fopen("non_existent_file.txt", "r");
	if (fp == NULL)
		printf("ERRNO[%d]: %s\n", errno, strerror(errno));
	// math domain errors
	errno = 0;
	double result = sqrt(-1.0);
	if (errno == EDOM)
		printf("ERRNO[%d]: %s\n", errno, strerror(errno)); 
	// memory allocation errors
	errno = 0;
	void* ptr = malloc(SIZE_MAX * 2);
	if (ptr == NULL && errno == ENOMEM)
		printf("ERRNO[%d]: %s\n", errno, strerror(errno));
}

// 演示线程局部 errno
static int thread_func(void* arg)
{
	(void)arg;
	errno = EAGAIN;
	printf("Thread ERRNO[%d]: %s\n", errno, strerror(errno));
	return 0;
}
static void example_thread_errors() {
	thrd_t thr;
	if (thrd_create(&thr, thread_func, NULL) == thrd_success) {
		thrd_join(thr, NULL);
		printf("Main ERRNO[%d]: %s\n", errno, strerror(errno)); // 应该是 0
	}
}

void test_errno(void)
{
	example_errors();
	example_thread_errors();
}