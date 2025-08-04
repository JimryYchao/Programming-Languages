#include "test.h"

#include <stdatomic.h>
#include <threads.h>

import std.compat;
using namespace std;

// 演示并发中的原子计数
static atomic_int atomic_counter = { 0 };
static atomic_int t_count = { 0 };
static int atomic_increment(void* arg)
{
	for (int i = 0; i < 100; i++) {
		atomic_fetch_add(&atomic_counter, 1);
		printf("%d\n", atomic_load(&atomic_counter));
	}
	atomic_fetch_sub(&t_count, 1);
	return 0;
}
static void example_atomic_counter() {
	puts("\n[Atomic counter in Concurrency]");
	thrd_t t;
#define T_MAXCOUNT 5

	thrd_create(&t, [](void*)-> int {
		thrd_t ts[T_MAXCOUNT]{ 0 };
		for (int i = 0; i < T_MAXCOUNT; i++)
			if (thrd_create(&ts[i], atomic_increment, NULL) == thrd_success)
				atomic_fetch_add(&t_count, 1);
		while (atomic_load(&t_count) > 0);
		return 0;
		}, NULL);

	thrd_join(t, NULL);
	printf("Atomic counter: %d (expected %d)\n", atomic_load(&atomic_counter), T_MAXCOUNT * 100);
}

// 演示原子标志自旋锁
static int thread_lock_free(void* arg)
{
	static atomic_flag lockFree_flag = ATOMIC_FLAG_INIT;
	const char* name = (const char*)arg;
	struct timespec ts { .tv_nsec = 100000000 };

	// 循环检查是否有其他线程占用自旋锁，直至 flag 被清除
	while (atomic_flag_test_and_set(&lockFree_flag)) {
		printf("[%s] Waiting for lock...\n", name);
		thrd_sleep(&ts, NULL); // 1s
	}
	printf("[%s] Lock acquired\n", name);

	// 占有自旋锁，并执行一些耗时操作
	ts.tv_sec = 1;
	thrd_sleep(&ts, NULL);   // do something

	// 清除自旋锁
	atomic_flag_clear(&lockFree_flag);  // false
	printf("[%s] Lock released\n", name);
	return 0;
}
static void example_lock_free(void) {
	puts("\n[Atomic flag in Concurrency]");
	thrd_t t1, t2;
	// 初始时清除位，ThreadA 获取自旋锁并设置位，ThreadB 等待位清除
	thrd_create(&t1, thread_lock_free, (char[])"ThreadA");
	thrd_create(&t2, thread_lock_free, (char[])"ThreadB");

	thrd_join(t1, NULL);
	thrd_join(t2, NULL);
}

// 演示原子比较交换
static void example_compare_exchange(void)
{
	puts("\n[Compare-and-Swap]");

	atomic_int val = ATOMIC_VAR_INIT(10);
	int expected = 10;

	// val 与 expected 位相等，val = 20
	bool success = atomic_compare_exchange_strong(&val, &expected, 20);
	printf("CAS(10→20): %s (val=%d, expected=%d)\n", success ? "success" : "failed", atomic_load(&val), expected);

	// val 更新为 20，与 expected 不等，expected = val
	success = atomic_compare_exchange_strong(&val, &expected, 30);
	printf("CAS(20→30): %s (val=%d, expected=%d)\n", success ? "success" : "failed ", atomic_load(&val), expected);
}

void test_stdatomic(void)
{
	example_atomic_counter();
	example_lock_free();
	example_compare_exchange();
}