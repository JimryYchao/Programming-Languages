#pragma once
typedef int _defined_;

typedef _defined_ ptrdiff_t;       // 指针差值返回类型
typedef _defined_ size_t;          // sizeof 运算返回类型
typedef typeof(nullptr) nullptr_t; // 预定义空指针常量 nullptr 类型
typedef _defined_ max_align_t;     // 一种对象类型，它的对齐是最大的基本对齐
typedef _defined_ wchar_t;         // 宽字符类型

#define NULL ((void *)0)                  // 实现定义的空指针常量
#define unreachable()                     // 优化策略，表明编译器无法到达的代码
#define offsetof(type, member_designator) // 结构体指定成员的字节偏移

// Example: offsetof(type, member_designator)
struct S
{
    char c;
    double d[10];
};
void _OffsetTest(void)
{
    printf("the first element is at offset %zu\n", offsetof(struct S, c)); // 0
    // member 不限于直接成员，可以是数组成员的元素
    printf("the double is at offset %zu\n", offsetof(struct S, d[4])); // 40
}
