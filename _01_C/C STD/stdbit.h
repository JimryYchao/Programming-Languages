#pragma once

#define __STDC_ENDIAN_BIG__     // 字节顺序存储方法，小端排序
#define __STDC_ENDIAN_LITTLE__  // 字节顺序存储方法，大端排序
#define __STDC_ENDIAN_NATIVE__  // 表示执行环境的字节顺序存储方法, __STDC_ENDIAN_BIG__ or __STDC_ENDIAN_LITTLE__ or other

// function_XX: uc, us, ui, ul, ull
unsigned int stdc_leading_zeros_uc(unsigned char value) [[unsequenced]];       // 前导 0 计数
unsigned int stdc_leading_ones_uc(unsigned char value) [[unsequenced]];        // 前导 1 计数
unsigned int stdc_trailing_zeros_uc(unsigned char value) [[unsequenced]];      // 尾随 0 计数
unsigned int stdc_trailing_ones_uc(unsigned char value) [[unsequenced]];       // 尾随 1 计数
unsigned int stdc_first_leading_zero_uc(unsigned char value) [[unsequenced]];  // 首个前导 0 位索引
unsigned int stdc_first_leading_one_uc(unsigned char value) [[unsequenced]];   // 首个前导 1 位索引
unsigned int stdc_first_trailing_zero_uc(unsigned char value) [[unsequenced]]; // 首个尾随 0 位索引
unsigned int stdc_first_trailing_one_uc(unsigned char value) [[unsequenced]];  // 首个尾随 1 位索引
unsigned int stdc_count_zeros_uc(unsigned char value) [[unsequenced]];         // 0 位计数
unsigned int stdc_count_ones_uc(unsigned char value) [[unsequenced]];          // 1 位计数  
bool stdc_has_single_bit_uc(unsigned char value) [[unsequenced]];              // 是否为 2 的整数次幂
unsigned int stdc_bit_width_uc(unsigned char value) [[unsequenced]];           // 存储该值的最小位数
unsigned char stdc_bit_floor_uc(unsigned char value) [[unsequenced]];          // 不大于值的 2 的最大整数次幂
unsigned char stdc_bit_ceil_uc(unsigned char value) [[unsequenced]];           // 不小于值的 2 的最小整数次幂