#pragma once

#define complex         _Complex        
#define _Complex_I      float _Complex
#define imaginary       _Imaginary
#define _Imaginary_I    float _Imaginary
#define I               _Complex_I   // or _Imaginary_I

#pragma STDC CX_LIMITED_RANGE on-off-switch

/* double: _FUNC_, float: _FUNC_f, long double: _FUNC_l */ 
double complex CMPLX(double x, double y);             // 构建复数
double creal(double complex z);                       // 取实部
double cimag(double complex z);                       // 取虚部     
double cabs(double complex z);                        // 取模    
double carg(double complex z);                        // 取辐角 
double complex conj(double complex z);                // 取共轭复数
double complex cproj(double complex z);               // 取黎曼球投影复数

double complex cexp(double complex z);                   // 计算指数 e^z
double complex clog(double complex z);                   // 计算自然对数 ln(z)
double complex cpow(double complex x, double complex y); // 计算复数幂 x^y
double complex csqrt(double complex z);                  // 计算算数平方根 √z

double complex csin(double complex z);               // 计算正弦 sin(z)
double complex ccos(double complex z);               // 计算余弦 cos(z)    
double complex ctan(double complex z);               // 计算正切 tan(z)
double complex casin(double complex z);              // 计算反正弦 asin(z)    
double complex cacos(double complex z);              // 计算反余弦 acos(z)    
double complex catan(double complex z);              // 计算反正切 atan(z)

double complex csinh(double complex z);              // 计算双曲正弦 sinh(z)
double complex ccosh(double complex z);              // 计算双曲余弦 cosh(z)
double complex ctanh(double complex z);              // 计算双曲正切 tanh(z)
double complex casinh(double complex z);             // 计算反双曲正弦 asinh(z)
double complex cacosh(double complex z);             // 计算反双曲余弦 acosh(z)
double complex catanh(double complex z);             // 计算反双曲正切 atanh(z)