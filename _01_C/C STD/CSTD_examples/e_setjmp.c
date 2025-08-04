#include "test.h"

#include <stdio.h>
#include <setjmp.h>  
#include <stdlib.h>  

// 全局跳转缓冲区
static jmp_buf jump_buffer;

// 模拟从深度嵌套函数跳出
static void nested_function(int depth, int maxDepth)
{
	printf("Entering nested function, depth: %d\n", depth);

	if (depth > maxDepth) {  // 达成跳出条件
		printf("longJmp condition triggered at depth %d\n", depth);
		longjmp(jump_buffer, depth);  // 跳回 setjmp 位置
	}

	nested_function(depth + 1, maxDepth);

	printf("This line won't be executed\n");  // 演示控制流中断
}
static void example_jump_nested_function(int maxDepth)
{
	volatile int should_retry = 1;  // 需要 volatile 防止优化
	int retry_count = 0;
	int jp_state = 0;

	while (should_retry) {
		// 设置跳转点
		jp_state = setjmp(jump_buffer);
		if (jp_state == 0) {

			printf("\nAttempt %d: Calling nested functions\n", ++retry_count);
			nested_function(1, maxDepth);
			should_retry = 0;  // 正常完成
		}
		else if (jp_state > maxDepth) {

			printf("Recovered from depth %d\n", jp_state);
			should_retry = 0;
		}
		else {
			printf("Unknown condition\n");
			exit(EXIT_FAILURE);
		}
	}
}

// 模拟错误处理
typedef void (*error_handler_t)(int);
static jmp_buf err_jmp;
static void (default_errHandler)(int code) {
	exit(code);
}
void error_handle(int code, error_handler_t f) {
	if (f)
		f(code);
	else default_errHandler(code);
}
#define  Error(code, format, ...) \
fprintf(stderr, "ERROR[%d] : " format "\n", code,__VA_ARGS__); \
longjmp(err_jmp, code)

static void raise_error(int code) {
	Error(code, "error testing in %s() at line %d", __func__, __LINE__);
}


void test_setjmp(void)
{
	puts("\n[Jump out nested function Example]");
	example_jump_nested_function(4);

	puts("\n[Error Handling Example]");
	int code = 0;
	if (code = setjmp(err_jmp))
		error_handle(code, NULL);
	raise_error(10086);
}