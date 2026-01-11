#include "test.h"

#include <stdio.h>
#include <stdlib.h>
#include <signal.h>  
#include <threads.h>
#include <windows.h>

// 全局标志位
static volatile sig_atomic_t keep_running = 1;

// 演示基本信号处理
static void handle_sigint(int sig)  // SIGINT 信号处理函数
{
	printf(">>> Received SIGINT (Ctrl+C) <<<\n");
	keep_running = 0;
}
static void example_basic_handling(void)
{
	puts("\n[Basic Signal Handling]");
	printf("Try pressing Ctrl+C...\n");
	signal(SIGINT, handle_sigint);

	while (keep_running) {
		printf("Running... (threadID: %u)\n", thrd_current()._Tid);
		Sleep(1000);
	}
	printf("Clean exit after signal\n");
	//signal(SIGINT, NULL);
}


// 演示信号屏蔽
static volatile int g_signal_mask = 0;      // 信号屏蔽状态
static HANDLE g_signal_event = NULL;		// 用于信号通知的事件对象
static BOOL WINAPI CtrlHandler(DWORD dwCtrlType) { // 控制台控制处理器
	if (dwCtrlType == CTRL_C_EVENT) {       // Ctrl+C
		if (g_signal_mask) {
			printf("Ctrl+C received but blocked!\n");
			SetEvent(g_signal_event);       // 通知有信号到达
			return TRUE;                    // 屏蔽信号
		}
		printf("Ctrl+C received and processed!\n");
	}
	return FALSE;
}
static void init_signal_system() {			// 初始化信号系统
	g_signal_event = CreateEvent(NULL, TRUE, FALSE, NULL);
	SetConsoleCtrlHandler(CtrlHandler, TRUE);
}
static void block_signals() {				// 屏蔽信号
	InterlockedExchange(&g_signal_mask, 1);
	printf("Signals blocked\n");
}
static void unblock_signals() {				// 解除信号屏蔽
	keep_running = 1;
	InterlockedExchange(&g_signal_mask, 0);
	printf("Signals unblocked\n");
	SetConsoleCtrlHandler(NULL, FALSE);
}
static void example_signal_blocking(void) {
	puts("\n[Signal Blocking]");
	init_signal_system();

	printf("Press Ctrl+C to test...\n");
	// 屏蔽信号
	block_signals();
	printf("Signals are blocked for 5 seconds\n");
	for (int i = 5; i > 0; i--) {
		printf("%d...\n", i);
		Sleep(1000);
	}
	// 检查是否有被屏蔽的信号
	if (WaitForSingleObject(g_signal_event, 0) == WAIT_OBJECT_0) {
		printf("A signal was blocked during this time!\n");
		ResetEvent(g_signal_event);
	}
	// 解除屏蔽
	unblock_signals();
	printf("Press Ctrl+C to exit the loop...\n");
	signal(SIGINT, handle_sigint);
	while (keep_running) {
		Sleep(1000);  // Ctrl+C 退出
	}
	printf("Clean exit after signal\n");
}

// 演示错误恢复
static void handle_sigsegv(int sig) // SIGSEGV 信号处理函数
{
	printf("Caught segmentation fault! Signal: %d\n", sig);
	exit(EXIT_FAILURE);
}
static void example_error_recovery(void)
{
	puts("\n[Segmentation Fault Handling]");
	// 设置 SIGSEGV 处理器
	signal(SIGSEGV, handle_sigsegv);
	printf("Before invalid memory access\n");
	// 故意制造段错误
	int* ptr = NULL;
	*ptr = 42;  // 这将触发 SIGSEGV, 并调用 handle_sigsegv()
	printf("This line won't be executed\n");
}

void test_signal(void)
{
	example_basic_handling();
	example_signal_blocking();
	example_error_recovery();
}
/*
* [Basic Signal Handling]
Try pressing Ctrl+C...
Running... (threadID: 49852)
Running... (threadID: 49852)
>>> Received SIGINT (Ctrl+C) <<<
Clean exit after signal

[Signal Blocking]
Press Ctrl+C to test...
Signals blocked
Signals are blocked for 5 seconds
5...
4...
Ctrl+C received but blocked!
3...
2...
1...
A signal was blocked during this time!
Signals unblocked
Press Ctrl+C to exit the loop...
Ctrl+C received and processed!
>>> Received SIGINT (Ctrl+C) <<<
Clean exit after signal

[Segmentation Fault Handling]
Before invalid memory access
Caught segmentation fault! Signal: 11
*/