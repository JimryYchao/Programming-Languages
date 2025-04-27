#pragma once
typedef int _defined_;

typedef _defined_ jmp_buf;                       // 保存上下文的缓冲区类型
int setjmp(jmp_buf env);                         // 保存当前上下文，调用 longjmp 恢复当前上下文
[[noreturn]] void longjmp(jmp_buf env, int val); // 跳转至指定位置，并将 val 传递给 setjmp 返回

// Example: longjmp
[[noreturn]] void _JumpTest(jmp_buf env, int count)
{
    printf("Jump count = %d\n", count);
    longjmp(env, count + 1); // 跳转到 setjmp 的调用位置
}
void _SetjmpTest(void)
{
	static int count = 0;
	jmp_buf env;
    if (setjmp(env) < 10)        // 保存当前上下文
        _JumpTest(env, ++count); // 保存当前上下文并跳转到 _JumpTest
}
