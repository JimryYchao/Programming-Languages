#pragma once
#ifndef __STDC_NO_ATOMICS__
typedef int _defined_;
typedef _defined_ type;

typedef _defined_ A; // 原子类型
typedef _defined_ C; // 对应的非原子类型
typedef _defined_ M; // 原子操作数类型

// 原子无锁宏：0 绝不免锁；1 有时免锁；2 始终免锁
#define ATOMIC_BOOL_LOCK_FREE
#define ATOMIC_CHAR_LOCK_FREE
#define ATOMIC_CHAR8_T_LOCK_FREE
#define ATOMIC_CHAR16_T_LOCK_FREE
#define ATOMIC_CHAR32_T_LOCK_FREE
#define ATOMIC_WCHAR_T_LOCK_FREE
#define ATOMIC_SHORT_LOCK_FREE
#define ATOMIC_INT_LOCK_FREE
#define ATOMIC_LONG_LOCK_FREE
#define ATOMIC_LLONG_LOCK_FREE
#define ATOMIC_POINTER_LOCK_FREE

// 内存同步定序
typedef enum
{
    memory_order_relaxed,                     // 宽松顺序
    memory_order_consume,                     // 消耗
    memory_order_acquire,                     // 获取
    memory_order_release,                     // 释放
    memory_order_acq_rel,                     // memory_order_acquire | memory_order_release
    memory_order_seq_cst                      // 顺序一致性
} memory_order;                               // 定义内存顺序
type kill_dependency(type y);                 // 打破 memory_order_consume 的依赖链
void atomic_thread_fence(memory_order order); // 通用的内存顺序依赖的栅栏同步原语
void atomic_signal_fence(memory_order order); // 线程与执行于同一线程的信号处理函数间的栅栏
// 原子操作函数
void atomic_init(volatile A *obj, C value);       // 初始化原子对象
bool atomic_is_lock_free(const volatile A *obj);  // 是否免锁
void atomic_store(volatile A *object, C desired); // 存储原子对象
void atomic_store_explicit(volatile A *object, C desired, memory_order order);
C atomic_load(const volatile A *object); // 加载原子对象
C atomic_load_explicit(const volatile A *object, memory_order order);
C atomic_exchange(volatile A *object, C desired); // 交换原子对象
C atomic_exchange_explicit(volatile A *object, C desired, memory_order order);
bool atomic_compare_exchange_strong(volatile A *object, C *expected, C desired); // 原子比较交换
bool atomic_compare_exchange_strong_explicit(volatile A *object, C *expected, C desired, memory_order success, memory_order failure);
bool atomic_compare_exchange_weak(volatile A *object, C *expected, C desired); // 原子比较交换
bool atomic_compare_exchange_weak_explicit(volatile A *object, C *expected, C desired, memory_order success, memory_order failure);
C atomic_fetch_add(volatile A *object, M operand); // 原子加法
C atomic_fetch_add_explicit(volatile A *object, M operand, memory_order order);
C atomic_fetch_sub(volatile A *object, M operand); // 原子减法
C atomic_fetch_sub_explicit(volatile A *object, M operand, memory_order order);
C atomic_fetch_or(volatile A *object, M operand); // 原子或操作
C atomic_fetch_or_explicit(volatile A *object, M operand, memory_order order);
C atomic_fetch_xor(volatile A *object, M operand); // 原子异或操作
C atomic_fetch_xor_explicit(volatile A *object, M operand, memory_order order);
C atomic_fetch_and(volatile A *object, M operand); // 原子与操作
C atomic_fetch_and_explicit(volatile A *object, M operand, memory_order order);
// 免锁原子布尔类型
typedef _defined_ atomic_flag;                               // 免锁原子布尔类型
#define ATOMIC_FLAG_INIT ((atomic_flag){0})                  // 原子布尔类型初始化
bool atomic_flag_test_and_set(volatile atomic_flag *object); // 设置 atomic_flag 为 true 并返回旧值
bool atomic_flag_test_and_set_explicit(volatile atomic_flag *object, memory_order order);
void atomic_flag_clear(volatile atomic_flag *object);        // 设置 atomic_flag 为 true 并返回旧值
void atomic_flag_clear_explicit(volatile atomic_flag *object, memory_order order);
#endif