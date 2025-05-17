#pragma once

typedef int sig_atomic_t;
// 信号处理函数
#define SIG_DFL
#define SIG_IGN
#define SIG_ERR    // 指示发生错误信号的返回类型
// 定义信号类型
#define SIGABRT     // 进程中止信号
#define SIGFPE      // 算术异常信号，如除零错误
#define SIGILL      // 非法程序映像，如非法指令
#define SIGINT      // 中断信号，如用户发动
#define SIGSEGV     // 非法内存访问
#define SIGTERM     // 发送终止信号

void (*signal(int sig, void (*func)(int)))(int);   // 为特定信号设置处理函数
int raise(int sig);         // 运行特定信号的处理函数

// Example: signal
void signal_Handler(int sig) // 信号处理函数的示例实现
{
	// 处理信号的代码
	printf("Received signal: %d\n", sig);
	if (sig == SIGABRT)
		printf("SIGABRT");
}
int _SignalTest(void)
{
	// 设置信号处理函数
	signal(SIGILL, signal_Handler);
	signal(SIGABRT, signal_Handler);
	// 触发 SIGILL 信号
	raise(SIGILL); 
	
	abort(); // 触发 SIGABRT 信号
}
