#include "test.h"

#include <stdio.h>
#include <stdlib.h>
#include <threads.h>
#include <string.h>
#include <stdbool.h>

static volatile int counter = 0;
static mtx_t mutex;

// 演示线程互斥锁同步
static int thread_mtx(void* arg)  // 互斥锁保护操作
{
	const char* name = (const char*)arg;
	struct timespec ts;

	for (int i = 0; i < 3; i++) {
		(void)timespec_get(&ts, TIME_UTC);
		ts.tv_nsec += 150 * 1000 * 1000;
		if (mtx_timedlock(&mutex, &ts) == thrd_success) {
			int val = ++counter;
			printf("[%s] mutex counter: %d\n", name, val);
			mtx_unlock(&mutex);
			thrd_sleep(&(struct timespec) { .tv_nsec = 200 * 1000 * 1000 }, NULL);
			continue;
		}
		else i--;
	}
	return 0;
}
static void example_mutex_protection(void)
{
	puts("\n[Mutex Protection]");
	counter = 0; // 重置计数器
	mtx_init(&mutex, mtx_timed);
	thrd_t thr1, thr2;
	// 创建线程（立即执行）
	thrd_create(&thr1, thread_mtx, "Thread1");
	thrd_create(&thr2, thread_mtx, "Thread2");
	// 阻塞等待线程结束
	thrd_join(thr1, NULL);
	thrd_join(thr2, NULL);
	printf("Final counter value: %d\n", counter);
	mtx_destroy(&mutex);
}

// 演示线程局部存储
static tss_t tss_key;
static thread_local uint32_t tls_id = 0;
static void cleanup_data(void* data)   // 资源清理函数
{
	printf("Cleaning up: %p\n", data);
	free(data);
}
static int thread_lc1(void* arg)  // 设置线程局部数据
{
	tls_id = thrd_current()._Tid;
	// 模拟分配线程局部数据
	int len = strlen((const char*)arg);
	char* data = malloc(sizeof(char) * (len + 1));
	if (!data)
		thrd_exit(0);
	strncpy(data, (const char*)arg, len + 1);
	// 设置 TSS 值
	tss_set(tss_key, data);
	printf("[Thread:%u] Set TSS: %s (%p)\n", tls_id, data, (void*)data);
	// 模拟工作
	thrd_sleep(&(struct timespec) { .tv_sec = 1 }, NULL);
	// 获取并验证数据
	char* retrieved = tss_get(tss_key);
	printf("[Thread:%u] Retrieved: %s (%p)\n", tls_id, retrieved, (void*)retrieved);
	return 0;
}
static int thread_lc2(void* arg)
{
	tls_id = thrd_current()._Tid;
	// 模拟分配线程局部数据
	int* data = malloc(sizeof(int));
	if (!data)
		thrd_exit(0);
	*data = strlen((const char*)arg);

	tss_set(tss_key, data);
	printf("[Thread:%u] Set TSS: %d (%p)\n", tls_id, *data, (void*)data);

	int* retrieved = tss_get(tss_key);
	printf("[Thread:%u] Retrieved: %d (%p)\n", tls_id, *retrieved, retrieved);
	return 0;
}
static void example_local_storage(void)
{
	puts("\n[Thread Local Storage]");
	tls_id = thrd_current()._Tid;
	// 创建 TSS 键（带清理函数）
	if (tss_create(&tss_key, cleanup_data) != thrd_success) {
		perror("tss_create failed");
		exit(EXIT_FAILURE);
	}
	// 创建线程并等待
	thrd_t thr1, thr2;
	thrd_create(&thr1, thread_lc1, "Hello");
	thrd_create(&thr2, thread_lc2, "World");
	thrd_join(thr1, NULL);
	thrd_join(thr2, NULL);

	printf("[Main:%d] TSS value: %p\n", tls_id, tss_get(tss_key));  // NULL
	tss_delete(tss_key);	// 销毁 TSS 键
}

// 借助条件变量 演示基本生产者-消费者模型
typedef struct {
	mtx_t mutex;
	cnd_t cond;
	volatile int data;
	volatile bool ready;
} SharedData;
static volatile int current_consumer_count = 0;
static int producer(void* arg)      // 生产者线程函数
{
	SharedData* sd = (SharedData*)arg;
	{
		mtx_lock(&sd->mutex);
		// 生产数据
		sd->data = 10086;
		sd->ready = true;
		thrd_sleep(&(struct timespec) { .tv_sec = 2 }, NULL);
		printf("[Producer] Data ready: %d\n", sd->data);
		// 通知一个消费者
		cnd_signal(&sd->cond);
		// cnd_broadcast(&sd->cond);    // 或者通知所有消费者: 
		mtx_unlock(&sd->mutex);
	}
	while (current_consumer_count > 0) {
		thrd_sleep(&(struct timespec) { .tv_sec = 1 }, NULL);
	}
	return 0;
}
static int consumer(void* arg)      // 消费者线程函数
{
	mtx_lock(&mutex);
	current_consumer_count++;
	mtx_unlock(&mutex);

	SharedData* sd = (SharedData*)arg;
	{
		mtx_lock(&sd->mutex);
		struct timespec ts;
		(void)timespec_get(&ts, TIME_UTC);
		ts.tv_sec += 4;
		// 等待条件满足
		printf("[Consumer:%d] Waiting for data...\n", thrd_current()._Tid);
		while (!sd->ready) {
			if (cnd_timedwait(&sd->cond, &sd->mutex, &ts) == thrd_timedout) {
				printf("[Consumer:%d] Waiting timeout\n", thrd_current()._Tid);
				mtx_unlock(&sd->mutex);

				mtx_lock(&mutex);
				current_consumer_count--;
				mtx_unlock(&mutex);
				thrd_exit(EXIT_FAILURE);
			}
		}
		// 消费数据
		printf("[Consumer:%d] Consumed: %d\n", thrd_current()._Tid, sd->data);
		sd->ready = false;
		mtx_unlock(&sd->mutex);
	}

	mtx_lock(&mutex);
	current_consumer_count--;
	mtx_unlock(&mutex);
	return 0;
}
static void example_producer_consumer(void)
{
	puts("\n[Producer-Consumer Example]");
	SharedData sd = {
		.data = 0,
		.ready = false
	};
	// 初始化同步原语
	mtx_init(&sd.mutex, mtx_plain);
	cnd_init(&sd.cond);
	thrd_t prod, cons1, cons2;

	// 创建消费者线程（先启动）
	thrd_create(&cons1, consumer, &sd);
	thrd_create(&cons2, consumer, &sd);
	// 创建生产者线程
	thrd_create(&prod, producer, &sd);

	// 等待所有线程完成
	thrd_join(prod, NULL);

	// 清理资源
	mtx_destroy(&sd.mutex);
	cnd_destroy(&sd.cond);
}


// 演示广播唤醒
typedef struct { mtx_t* m; cnd_t* c; int* volatile cnt; } Mct;
static int thread_worker(void* arg) {
	mtx_t* m = ((Mct*)arg)->m;
	cnd_t* c = ((Mct*)arg)->c;
	int* volatile cnt = ((Mct*)arg)->cnt;

	mtx_lock(m);
	// 等待条件
	while (*cnt < 5) 
		cnd_wait(c, m);
	printf("[Worker:%d] Woken up!\n", thrd_current()._Tid);
	mtx_unlock(m);
	return 0;
}
static void example_broadcast(void)
{
	puts("\n[Broadcast Wakeup Example]");
	mtx_t mtx;
	cnd_t cond;
	volatile int count = 0;
	mtx_init(&mtx, mtx_plain);
	cnd_init(&cond);

	// 创建多个工作线程
	Mct arg = { &mtx, &cond, &count };
	thrd_t workers[3];
	for (int i = 0; i < 3; i++) 
		thrd_create(&workers[i], thread_worker, &arg);
	// 主线程改变条件
	mtx_lock(&mtx);
	count = 5;
	puts("[Main] Broadcasting wakeup...");
	cnd_broadcast(&cond);
	mtx_unlock(&mtx);
	// 等待工作线程
	for (int i = 0; i < 3; i++) 
		thrd_join(workers[i], NULL);

	mtx_destroy(&mtx);
	cnd_destroy(&cond);
}
void test_threads(void)
{
	example_mutex_protection();
	example_local_storage();
	example_producer_consumer();
	example_broadcast();
}
/*
[Mutex Protection]
[Thread1] mutex counter: 1
[Thread2] mutex counter: 2
[Thread2] mutex counter: 3
[Thread1] mutex counter: 4
[Thread2] mutex counter: 5
[Thread1] mutex counter: 6
Final counter value: 6

[Thread Local Storage]
[Thread:50020] Set TSS: 5 (00000255F8C58A30)
[Thread:52660] Set TSS: Hello (00000255F8C6BB10)
[Thread:50020] Retrieved: 5 (00000255F8C58A30)
Cleaning up: 00000255F8C58A30
[Thread:52660] Retrieved: Hello (00000255F8C6BB10)
Cleaning up: 00000255F8C6BB10
[Main:19596] TSS value: 0000000000000000

[Producer-Consumer Example]
[Consumer:42256] Waiting for data...
[Producer] Data ready: 10086
[Consumer:42256] Consumed: 10086
[Consumer:2196] Waiting for data...
[Consumer:2196] Waiting timeout

[Broadcast Wakeup Example]
[Main] Broadcasting wakeup...
[Worker:43272] Woken up!
[Worker:21408] Woken up!
[Worker:34956] Woken up!
*/
