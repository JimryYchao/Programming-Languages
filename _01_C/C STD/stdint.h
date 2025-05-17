#pragma once
typedef int _defined_;

// N = 8, 16, 32, 64
typedef _defined_ intN_t;        // 定宽精确整数
typedef _defined_ uintN_t;       
typedef _defined_ int_least8_t;  // 最小整数
typedef _defined_ uint_least8_t; 
typedef _defined_ int_fast8_t;   // 最快整数
typedef _defined_ uint_fast8_t;  
typedef _defined_ intptr_t;      // 保存对象指针整数类型
typedef _defined_ uintptr_t;     
typedef _defined_ intmax_t;      // 最大整数类型
typedef _defined_ uintmax_t;
// 整数范围，N = 8, 16, 32, 64
#define INT8_MIN 
#define INT8_MAX 
#define UINT8_MAX 
#define INT_LEAST8_MIN 
#define INT_LEAST8_MAX 
#define UINT_LEAST8_MAX 
#define INT_FAST8_MIN 
#define INT_FAST8_MAX 
#define UINT_FAST8_MAX 
#define INTPTR_MIN 
#define INTPTR_MAX 
#define UINTPTR_MAX 
#define INTMAX_MIN 
#define INTMAX_MAX 
#define UINTMAX_MAX 
// 整数宽度，N = 8, 16, 32, 64
#define INT8_WIDTH 
#define UINT8_WIDTH 
#define INT_LEAST8_WIDTH 
#define UINT_LEAST8_WIDTH 
#define INT_FAST8_WIDTH 
#define UINT_FAST8_WIDTH 
#define INTPTR_WIDTH 
#define UINTPTR_WIDTH 
#define INTMAX_WIDTH
#define UINTMAX_WIDTH 
// 其他整数类型
#define PTRDIFF_MIN       // ptrdiff_t
#define PTRDIFF_MAX 
#define SIG_ATOMIC_WIDTH  // sig_atomic_t
#define SIG_ATOMIC_MIN 
#define SIG_ATOMIC_MAX 
#define SIZE_WIDTH        // size_t 
#define SIZE_MAX 
#define WCHAR_WIDTH       // wchar_t
#define WCHAR_MIN
#define WCHAR_MAX
#define WINT_WIDTH        // wint_t 
#define WINT_MIN
#define WINT_MAX
#define __STDC_LIB_EXT1__
#ifdef(__STDC_WANT_LIB_EXT1__ &&__STDC_LIB_EXT1__)
#define RSIZE_MAX SIZE_MAX      // rsize_t
#endif

// 构造整数常量的函数宏，N = 8, 16, 32, 64
#define INT8_C(x)     // 定宽整数常量
#define UINT8_C(x) 
#define INTMAX_C(x)   // intmax_t 常量
#define UINTMAX_C(x)