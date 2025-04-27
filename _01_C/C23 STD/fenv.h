#pragma once

#pragma STDC FENV_ACCESS   // on-off-switch, 开启对浮点环境的访问
// #pragma fenv_access(on | off) 
typedef int _defined_;    

// 浮点异常控制  
#define FE_INVALID               // 无效操作的定义域错误。
#define FE_DIVBYZERO             // 除数为零的定义域错误。 
#define FE_OVERFLOW              // 溢出错误。
#define FE_UNDERFLOW             // 下溢错误。
#define FE_INEXACT               // 不精确错误。
#define FE_ALL_EXCEPT       (FE_INVALID | FE_DIVBYZERO | FE_OVERFLOW | FE_UNDERFLOW | FE_INEXACT)
typedef _defined_ fexcept_t;        // 浮点异常类型
int fesetexcept(int excepts);       // 设置指定异常标志
int feraiseexcept(int excepts);     // 触发指定异常标志
int feclearexcept(int excepts);     // 清除指定异常标志
int fetestexcept(int excepts);      // 测试指定异常标志
int fegetexceptflag(fexcept_t *flagp, int excepts);          // 在 flagp 中获取指定异常标志
int fesetexceptflag(const fexcept_t *flagp, int excepts);    // 在 flagp 中设置指定异常标志
int fetestexceptflag(const fexcept_t *flagp, int excepts);   // 在 flagp 中测试指定异常标志

// 浮点舍入方向控制
#define FE_TONEAREST	         // 四舍五入到最接近的值。
#define FE_DOWNWARD	             // 向下舍入到最接近的值。
#define FE_UPWARD	             // 向上舍入到最接近的值。
#define FE_TOWARDZERO	         // 向零舍入到最接近的值。
#define FE_TONEARESTFROMZERO     // 从零开始的四舍五入到最接近的值。
int fegetround(void);            // 获取当前舍入方向
int fesetround(int rnd);         // 设置舍入方向
#pragma STDC FENV_ROUND      direction    // 为翻译单元或复合语句中设置当前浮点舍入方向
#pragma STDC FENV_ROUND      FE_DYNAMIC   // 设置为动态浮点环境
#ifdef __STDC_IEC_60559_DFP__    // 对应的十进制浮点数舍入方向
#define FE_DEC_TONEAREST
#define FE_DEC_DOWNWARD
#define FE_DEC_UPWARD
#define FE_DEC_TOWARDZERO
#define FE_DEC_TONEARESTFROMZERO
#pragma STDC FENV_DEC_ROUND dec-direction
int fe_dec_getround(void);
int fe_dec_setround(int rnd);
#endif

// 动态浮点环境控制
typedef _defined_ femode_t;     // 浮点模式类型
typedef _defined_ fenv_t;       // 浮点环境类型
#define FE_DFL_ENV    ((const fenv_t *) 0)    // 默认浮点环境
#define FE_DFL_MODE   ((const femode_t *) 0)  // 默认浮点模式
int fegetenv(fenv_t *envp);             // 获取当前浮点环境
int fesetenv(const fenv_t *envp);       // 设置当前浮点环境
int feholdexcept(fenv_t *envp);         // 保存当前浮点环境并清除异常标志
int feupdateenv(const fenv_t *envp);    // 更新当前浮点环境并清除异常标志
int fegetmode(femode_t *modep);         // 获取当前浮点模式
int fesetmode(const femode_t *modep);   // 设置当前浮点模式