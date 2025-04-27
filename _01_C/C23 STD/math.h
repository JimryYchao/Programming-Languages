#pragma once

typedef float float_t;   // 至少为 float
typedef double double_t; // 至少为 double

#define HUGE_VAL // 指示过大而无法表示的 float 值，无穷大
#define INFINITY // 无穷大
#define NAN      // 非数

// 可选定义，表示对应函数的执行速度与操作数的直接算术操作一样快
#define FP_FAST_FMA
#define FP_FAST_FADD
#define FP_FAST_FSUB
#define FP_FAST_FMUL
#define FP_FAST_FDIV
#define FP_FAST_FSQRT
#define FP_FAST_FFMA

#define FP_ILOGB0   // ilogb(x) 在 x 为零时返回
#define FP_ILOGBNAN // ilogb(x) 在 x 为非数时返回
#define FP_LLOGB0   // llogb(x) 在 x 为零时返回
#define FP_LLOGBNAN // llogb(x) 在 x 为非数时返回

// 数学函数的错误处理机制
#define MATH_ERRNO 1       
#define MATH_ERREXCEPT 2
#define math_errhandling (MATH_ERRNO | MATH_ERREXCEPT)

#pragma STDC FP_CONTRACT on_off
// 分类
int fpclassify(real_floating x);  // 分类为以下常量之一
#define FP_INFINITE               // 无穷大
#define FP_NAN                    // 非数
#define FP_NORMAL                 // 规范化数
#define FP_SUBNORMAL              // 非规范化数
#define FP_ZERO                   // 零
int signbit(real_floating x);     // 是否为负数
int iscanonical(real_floating x); // 是否为规范化数
int isfinite(real_floating x);    // 是否有限
int isinf(real_floating x);       // 是否为无穷大
int isnan(real_floating x);       // 是否为非数
int isnormal(real_floating x);    // 是否为正规
int issubnormal(real_floating x); // 是否为非正规
int issignaling(real_floating x); // 是否 signaling NaN
int iszero(real_floating x);      // 是否为零
// 比较
int isgreater(real_floating x, real_floating y);      // >
int isgreaterequal(real_floating x, real_floating y); // >=
int isless(real_floating x, real_floating y);         // <
int islessequal(real_floating x, real_floating y);    // <=
int islessgreater(real_floating x, real_floating y);  // < or >
int isunordered(real_floating x, real_floating y);    // 是否无序
int iseqsig(real_floating x, real_floating y);        // 是否相等
// 三角函数
double acos(double x);              // 反余弦
double asin(double x);              // 反正弦
double atan(double x);              // 反正切
double atan2(double y, double x);   // 反正切，y/x
double cos(double x);               // 余弦
double sin(double x);               // 正弦
double tan(double x);               // 正切
double acospi(double x);            // 反余弦，pi 进制
double asinpi(double x);            // 反正弦，pi 进制
double atanpi(double x);            // 反正切，pi 进制
double atan2pi(double y, double x); // 反正切，y/x，pi 进制
double cospi(double x);             // 余弦，pi 进制
double sinpi(double x);             // 正弦，pi 进制
double tanpi(double x);             // 正切，pi 进制
// 双曲函数
double acosh(double x); // 反双曲余弦
double asinh(double x); // 反双曲正弦
double atanh(double x); // 反双曲正切
double cosh(double x);  // 双曲余弦
double sinh(double x);  // 双曲正弦
double tanh(double x);  // 双曲正切
// 幂函数
double exp(double x);                        // e^x
double expm1(double x);                      // e^x - 1
double exp10(double x);                      // 10^x
double exp10m1(double x);                    // 10^x - 1
double exp2(double x);                       // 2^x
double exp2m1(double x);                     // 2^x - 1
double frexp(double value, int *p);          // 分解为 m * 2^p
double ldexp(double x, int p);               // x * 2^p
double log(double x);                        // log_e x
double log10(double x);                      // log_10 x
double log10p1(double x);                    // log_10 (x + 1)
double log1p(double x);                      // log_e (x + 1)
double logp1(double x);                      // log_e (1 + x)
double log2(double x);                       // log_2 x
double log2p1(double x);                     // log_2 (x + 1)
double logb(double x);                       // 将 x 的指数提取为 double 值
double scalbn(double x, int n);              // x * FLT_RADIX^n
double cbrt(double x);                       // x^(1/3)
double compoundn(double x, long long int n); // (1+x)^n
double hypot(double x, double y);            // sqrt(x^2 + y^2)
double pow(double x, double y);              // x^y
double pown(double x, long long int n);      // x^n
double powr(double y, double x);             // e^(y * log_e(x))
double rootn(double x, long long int n);     // x^(1/n)
double rsqrt(double x);                      // 1/sqrt(x)
double sqrt(double x);                       // x^(1/2)
// 功能函数
double fabs(double x);                       // |x|
double fmod(double x, double y);             // 取模，x - n*y，n 为整数
double modf(double value, double *iptr);     // 分解为整数和小数部分
double remquo(double x, double y, int *quo); // 余数和商
double remainder(double x, double y);        // x - n*y，n 是最接近 x/y 准确值的整数值
// 误差与伽马函数
double erf(double x);                        // 误差函数
double erfc(double x);                       // 互补误差函数
double lgamma(double x);                     // 对数伽马函数
double tgamma(double x);                     // 伽马函数
// 舍入函数
double ceil(double x);      // 向上取整
double floor(double x);     // 向下取整
double nearbyint(double x); // 使用当前舍入方向，进行取整
double rint(double x);      // 使用当前舍入方向，进行取整
double round(double x);     // 四舍五入，不受当前舍入模式影响
double roundeven(double x); // 向最接近的偶数舍入
double trunc(double x);     // 向零舍入
// 整数转换
#define FP_INT_UPWARD                                   // 向上舍入
#define FP_INT_DOWNWARD                                 // 向下舍入
#define FP_INT_TOWARDZERO                               // 向零舍入
#define FP_INT_TONEARESTFROMZERO                        // 从零开始的四舍五入
#define FP_INT_TONEAREST                                // 四舍五入
double fromfp(double x, int rnd, unsigned int width);   // 使用指定的舍入方向舍入为有符号整数
double ufromfp(double x, int rnd, unsigned int width);  // 使用指定的舍入方向舍入为无符号整数
double fromfpx(double x, int rnd, unsigned int width);  // 使用指定的舍入方向舍入为有符号整数，报告不精确性
double ufromfpx(double x, int rnd, unsigned int width); // 使用指定的舍入方向舍入为无符号整数，报告不精确性
// 浮点操作
double copysign(double x, double y);           // 复制符号位
double nan(const char *tagp);                  // 转换为 quiet NAN
double nextafter(double x, double y);          // x -> y 的下一个有限浮点数
double nexttoward(double x, long double y);    // x -> y 的下一个有限浮点数
double nextup(double x);                       // 确定下一个大于 x 的可表示的浮点值
double nextdown(double x);                     // 确定下一个小于 x 的可表示的浮点值
int canonicalize(double *cx, const double *x); // 规范化浮点数
double fma(double x, double y, double z);      // x*y + z
float fadd(double x, double y);                // 加法
float fsub(double x, double y);                // 减法
float fmul(double x, double y);                // 乘法
float fdiv(double x, double y);                // 除法
float ffma(double x, double y, double z);      // x*y + z 舍入窄
float fsqrt(double x);                         // 平方根舍入窄
// 比较和正差函数
double fdim(double x, double y);     // x - y，若 x < y，则返回 0
double fmax(double x, double y);     // x > y ? x : y
double fmin(double x, double y);     // x < y ? x : y
double fmaximum(double x, double y); // 两数之间的最大值
double fmaximum_num(double x, double y);
double fmaximum_mag(double x, double y); // 两数的模之间的最大值
double fmaximum_mag_num(double x, double y);
double fminimum(double x, double y); // 两数之间的最小值
double fminimum_num(double x, double y);
double fminimum_mag(double x, double y); // 两数的模之间的最小值
double fminimum_mag_num(double x, double y);

// 十进制浮点数
#ifdef __STDC_IEC_60559_DFP__
#ifdef __STDC_WANT_IEC_60559_EXT__
typedef _Decimal32 _Decimal32_t;
typedef _Decimal64 _Decimal64_t;
#define HUGE_VAL_D32
#define HUGE_VAL_D64
#define HUGE_VAL_D128
#endif
#define DEC_INFINITY
#define DEC_NAN
#define FP_FAST_FMAD32
#define FP_FAST_D32ADDD64
#define FP_FAST_D32SUBD64
#define FP_FAST_D32MULD64
#define FP_FAST_D32DIVD64
#define FP_FAST_D32FMAD64
#define FP_FAST_D32SQRTD64

_Decimal32 acosd32(_Decimal32 x);
_Decimal32 asind32(_Decimal32 x);
_Decimal32 atand32(_Decimal32 x);
_Decimal32 atan2d32(_Decimal32 y, _Decimal32 x);
_Decimal32 cosd32(_Decimal32 x);
_Decimal32 sind32(_Decimal32 x);
_Decimal32 tand32(_Decimal32 x);
_Decimal32 acospid32(_Decimal32 x);
_Decimal32 asinpid32(_Decimal32 x);
_Decimal32 atanpid32(_Decimal32 x);
_Decimal32 atan2pid32(_Decimal32 y, _Decimal32 x);
_Decimal32 cospid32(_Decimal32 x);
_Decimal32 sinpid32(_Decimal32 x);
_Decimal32 tanpid32(_Decimal32 x);
_Decimal32 acoshd32(_Decimal32 x);
_Decimal32 asinhd32(_Decimal32 x);
_Decimal32 atanhd32(_Decimal32 x);
_Decimal32 coshd32(_Decimal32 x);
_Decimal32 sinhd32(_Decimal32 x);
_Decimal32 tanhd32(_Decimal32 x);
_Decimal32 expd32(_Decimal32 x);
_Decimal32 exp10d32(_Decimal32 x);
_Decimal32 exp10m1d32(_Decimal32 x);
_Decimal32 exp2d32(_Decimal32 x);
_Decimal32 exp2m1d32(_Decimal32 x);
_Decimal32 expm1d32(_Decimal32 x);
_Decimal32 frexpd32(_Decimal32 value, int *p);
int ilogbd32(_Decimal32 x);
_Decimal32 ldexpd32(_Decimal32 x, int p);
long int llogbd32(_Decimal32 x);
_Decimal32 logd32(_Decimal32 x);
_Decimal32 log10d32(_Decimal32 x);
_Decimal32 log10p1d32(_Decimal32 x);
_Decimal32 log1pd32(_Decimal32 x);
_Decimal32 logp1d32(_Decimal32 x);
_Decimal32 log2d32(_Decimal32 x);
_Decimal32 log2p1d32(_Decimal32 x);
_Decimal32 logbd32(_Decimal32 x);
_Decimal32 modfd32(_Decimal32 x, _Decimal32 *iptr);
_Decimal32 scalbnd32(_Decimal32 x, int n);
_Decimal32 scalblnd32(_Decimal32 x, long int n);
_Decimal32 cbrtd32(_Decimal32 x);
_Decimal32 compoundnd32(_Decimal32 x, long long int n);
_Decimal32 fabsd32(_Decimal32 x);
_Decimal32 hypotd32(_Decimal32 x, _Decimal32 y);
_Decimal32 powd32(_Decimal32 x, _Decimal32 y);
_Decimal32 pownd32(_Decimal32 x, long long int n);
_Decimal32 powrd32(_Decimal32 y, _Decimal32 x);
_Decimal32 rootnd32(_Decimal32 x, long long int n);
_Decimal32 rsqrtd32(_Decimal32 x);
_Decimal32 sqrtd32(_Decimal32 x);
_Decimal32 erfd32(_Decimal32 x);
_Decimal32 erfcd32(_Decimal32 x);
_Decimal32 lgammad32(_Decimal32 x);
_Decimal32 tgammad32(_Decimal32 x);
_Decimal32 ceild32(_Decimal32 x);
_Decimal32 floord32(_Decimal32 x);
_Decimal32 nearbyintd32(_Decimal32 x);
_Decimal32 rintd32(_Decimal32 x);
long int lrintd32(_Decimal32 x);
long long int llrintd32(_Decimal32 x);
_Decimal32 roundd32(_Decimal32 x);
long int lroundd32(_Decimal32 x);
long long int llroundd32(_Decimal32 x);
_Decimal32 roundevend32(_Decimal32 x);
_Decimal32 truncd32(_Decimal32 x);
_Decimal32 fromfpd32(_Decimal32 x, int rnd, unsigned int width);
_Decimal32 ufromfpd32(_Decimal32 x, int rnd, unsigned int width);
_Decimal32 fromfpxd32(_Decimal32 x, int rnd, unsigned int width);
_Decimal32 ufromfpxd32(_Decimal32 x, int rnd, unsigned int width);
_Decimal32 fmodd32(_Decimal32 x, _Decimal32 y);
_Decimal32 remainderd32(_Decimal32 x, _Decimal32 y);
_Decimal32 copysignd32(_Decimal32 x, _Decimal32 y);
_Decimal32 nand32(const char *tagp);
_Decimal32 nextafterd32(_Decimal32 x, _Decimal32 y);
_Decimal32 nexttowardd32(_Decimal32 x, _Decimal128 y);
_Decimal32 nextupd32(_Decimal32 x);
_Decimal32 nextdownd32(_Decimal32 x);
int canonicalized32(_Decimal32 *cx, const _Decimal32 *x);
_Decimal32 fdimd32(_Decimal32 x, _Decimal32 y);
_Decimal32 fmaxd32(_Decimal32 x, _Decimal32 y);
_Decimal32 fmind32(_Decimal32 x, _Decimal32 y);
_Decimal32 fmaximumd32(_Decimal32 x, _Decimal32 y);
_Decimal32 fminimumd32(_Decimal32 x, _Decimal32 y);
_Decimal32 fmaximum_magd32(_Decimal32 x, _Decimal32 y);
_Decimal32 fminimum_magd32(_Decimal32 x, _Decimal32 y);
_Decimal32 fmaximum_numd32(_Decimal32 x, _Decimal32 y);
_Decimal32 fminimum_numd32(_Decimal32 x, _Decimal32 y);
_Decimal32 fmaximum_mag_numd32(_Decimal32 x, _Decimal32 y);
_Decimal32 fminimum_mag_numd32(_Decimal32 x, _Decimal32 y);
_Decimal32 fmad32(_Decimal32 x, _Decimal32 y, _Decimal32 z);
_Decimal32 d32addd64(_Decimal64 x, _Decimal64 y);
_Decimal32 d32subd64(_Decimal64 x, _Decimal64 y);
_Decimal32 d32muld64(_Decimal64 x, _Decimal64 y);
_Decimal32 d32divd64(_Decimal64 x, _Decimal64 y);
_Decimal32 d32fmad64(_Decimal64 x, _Decimal64 y, _Decimal64 z);
_Decimal32 d32sqrtd64(_Decimal64 x);

// 量子函数与量子指数函数
_Decimal32 quantized32(_Decimal32 x, _Decimal32 y); // 函数计算具有 x 和 y 的量子指数的值
bool samequantumd32(_Decimal32 x, _Decimal32 y);    // 函数检查 x 和 y 的量子指数是否相同
_Decimal32 quantumd32(_Decimal32 x);                // 函数计算有限参数 x 的量子
long long int llquantexpd32(_Decimal32 x);          // 函数计算有限参数 x 的量子指数
// 十进制重编码函数
void encodedecd32(unsigned char encptr[restrict static 4], const _Decimal32 *restrict xptr); // 转换 *xptr 为十进制编码
void decodedecd32(_Decimal32 *restrict xptr, const unsigned char encptr[restrict static 4]); // 解码十进制数
void encodebind32(unsigned char encptr[restrict static 4], const _Decimal32 *restrict xptr); // 转换 *xptr 为二进制编码
void decodebind32(_Decimal32 *restrict xptr, const unsigned char encptr[restrict static 4]); // 解码二进制数

int totalorderd32(const _Decimal32 *x, const _Decimal32 *y);
int totalordermagd32(const _Decimal32 *x, const _Decimal32 *y);
_Decimal32 getpayloadd32(const _Decimal32 *x);
int setpayloadd32(_Decimal32 *res, _Decimal32 pl);
int setpayloadsigd32(_Decimal32 *res, _Decimal32 pl);
#endif