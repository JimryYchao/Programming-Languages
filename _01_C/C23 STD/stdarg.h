#pragma once

typedef char *  va_list;     // 保存宏 va_start、va_arg、va_end、va_copy 所需的信息

void va_start(va_list ap, ...);   // 初始化 va_list
type va_arg(va_list ap, type);    // 访问下一个参数 
void va_copy(va_list dest, va_list src);   // 复制 va_list 状态 
void va_end(va_list ap);          // 结束 va_list 访问

// Example: va_list
int add_nums_c23(int count, ...) {
    int result = 0;
    va_list args;
    va_start(args);
    for (int i = 0; i < count; ++i) 
        result += va_arg(args, int);
    va_end(args);
}  // sum = add_nums(4, 25, 25, 50, 50); 