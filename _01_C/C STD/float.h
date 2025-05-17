#pragma once

#define FLT_RADIX           // 浮点类型基数   
#define FLT_ROUNDS          // 浮点舍入模式的特征值，-1 表示不支持
#define INFINITY            // 正无穷大
#define NAN                 // 非数值
#define FLT_SNAN            // （DBL，LDBL）扩展为表示 signaling NaN 的相应类型的常量表达式
#define FLT_EVAL_METHOD	    // 浮点中间运算的浮点格式
#define DECIMAL_DIG         // 默认十进制位数的舍入精度，C23 起弃用

// FLT，DBL，LDBL
#define FLT_IS_IEC_60559    // 是否符合IEC 60559标准
#define FLT_DECIMAL_DIG     // 十进制数字位数
#define FLT_MANT_DIG        // 以 FLT_RADIX 位基数浮点有效位数
#define FLT_DIG             // 相应类型的精度的小数位数
#define FLT_EPSILON         // 使 x + 1.0 不等于 1.0 的最小正数 x
#define FLT_MIN             // 最小规范化正浮点数
#define FLT_MAX             // 可表示的最大浮点数
#define FLT_MIN_EXP         // 最小二进制指数 
#define FLT_MAX_EXP         // 最大二进制指数
#define FLT_MIN_10_EXP      // 最小十进制指数
#define FLT_MAX_10_EXP      // 最大十进制指数
#define FLT_TRUE_MIN        // 最小可表示的正浮点数
#define FLT_NORM_MAX        // 最大归一化浮点数
#define FLT_HAS_SUBNORM	    // 是否有非规格化数，C23 起弃用
// DEC32，64，128
#ifdef __STDC_IEC_60559_DFP__   
#define DEC_EVAL_METHOD
#define DEC_INFINITY
#define DEC_NAN
#define DEC32_SNAN
#define DEC32_MANT_DIG
#define DEC32_EPSILON
#define DEC32_MIN_EXP
#define DEC32_MAX_EXP
#define DEC32_MIN
#define DEC32_MAX
#define DEC32_TRUE_MIN
#endif