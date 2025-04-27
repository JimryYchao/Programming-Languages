#pragma once
typedef int Tres;
typedef int T1;
typedef int T2;

// 泛型宏尝试对 a,b 进行加减乘法，并将正确结果存储到 *result, 运算失败则返回 false
bool ckd_add(Tres *result, T1 a, T2 b);    
bool ckd_sub(Tres *result, T1 a, T2 b);
bool ckd_mul(Tres *result, T1 a, T2 b);