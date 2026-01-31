#include "test.h"

#include <fenv.h>
#include <stdio.h>
#include <math.h>
#include <stdbool.h>

// 启用浮点环境访问（C23）
#if __STDC_VERSION__ >= 202311L
#pragma STDC FENV_ACCESS ON
#else
#pragma fenv_access (on)  
#endif

// 测试浮点异常
void show_fe_exceptions(bool isclear) {
	printf("current floating-point exception raised: ");
	if (fetestexcept(FE_DIVBYZERO)) printf("  FE_DIVBYZERO");
	if (fetestexcept(FE_INVALID))   printf("  FE_INVALID");
	if (fetestexcept(FE_OVERFLOW))  printf("  FE_OVERFLOW");
	if (fetestexcept(FE_UNDERFLOW)) printf("  FE_UNDERFLOW");
	if (fetestexcept(FE_INEXACT))   printf("  FE_INEXACT");
	if (fetestexcept(FE_ALL_EXCEPT) == 0) printf("  NONE");
	printf("\n");
	if (isclear)
		feclearexcept(FE_ALL_EXCEPT);    // 清除所有异常标志
}
void example_fe_exceptions(void) {
	printf("\n>>> Floating-Point Exception Testing\n");
	feclearexcept(FE_ALL_EXCEPT);
	double a = 1.0, b = 0.0;
	printf("1.0 / 0.0 = %f\n", a / b);		       // FE_DIVBYZERO
	show_fe_exceptions(true);

	printf("sqrt(-1.0) = %f\n", sqrt(-1.0));    // FE_INVALID
	show_fe_exceptions(true);

	printf("exp(1000.0) = %e\n", exp(1000));    // FE_OVERFLOW 
	show_fe_exceptions(true);

	printf("exp(-1000.0) = %e\n", exp(-1000));  // FE_UNDERFLOW 
	show_fe_exceptions(true);

	printf("1.0 / 3.0 = %.20f\n", (1.0 / 3.0));	   // FE_INEXACT
	show_fe_exceptions(true);
}

// 测试舍入方向
void show_fe_rounding(void)
{
	printf("current rounding direction:  ");
	switch (fegetround())
	{
	case FE_TONEAREST:  printf("FE_TONEAREST");  break;
	case FE_DOWNWARD:   printf("FE_DOWNWARD");   break;
	case FE_UPWARD:     printf("FE_UPWARD");     break;
	case FE_TOWARDZERO: printf("FE_TOWARDZERO"); break;
	// case FE_TONEARESTFROMZERO: printf("FE_TONEARESTFROMZERO"); break;   // C23
	default:            printf("unknown");
	};
	printf("\n");
}
void example_fe_rounding(void) {
	printf("\n>>> Rounding Direction Testing\n");
	fenv_t env;
	fegetenv(&env); // 保存当前环境

	volatile double value = 1.0 / 3.0; // 0.333...
	double result;

	// 测试所有舍入方向
	const struct {
		int mode;
		const char* name;
	} modes[] = {
		{FE_TONEAREST,  "FE_TONEAREST (Round to nearest)"},
		{FE_DOWNWARD,   "FE_DOWNWARD (Round toward -∞)"},
		{FE_UPWARD,     "FE_UPWARD (Round toward +∞)"},
		{FE_TOWARDZERO, "FE_TOWARDZERO (Round toward zero)"}
	};

	for (size_t i = 0; i < sizeof(modes) / sizeof(modes[0]); i++) {
		fesetround(modes[i].mode);	// 设置舍入方向
		result = rint(value * 10);		// 将 3.333... 舍入到整数
		printf("%-35s: %.20f → %.0f\n",
			modes[i].name, value * 10, result);
	}
	fesetenv(&env); // 恢复原始环境
}

// 测试浮点环境控制
void example_fe_environment(void) {
	printf("\n>>> Environment Control\n");
	fenv_t env1, env2;
	fexcept_t ex;
	fegetenv(&env1);  // 保存初始浮点环境，初始为 FE_TONEAREST

	feraiseexcept(FE_DIVBYZERO | FE_INVALID);
	fesetround(FE_DOWNWARD);
	fegetexceptflag(&ex, FE_INVALID);  // 提取 flag
	show_fe_exceptions(false);
	show_fe_rounding();
	feholdexcept(&env2);  // 保存浮点环境和异常，并清除异常标志
	show_fe_exceptions(false);
	show_fe_rounding();

	// 引发 inexact，overflow
	fesetround(FE_UPWARD);
	double x = DBL_MAX; x *= 2;
	show_fe_exceptions(false);
	show_fe_rounding();

	// 恢复 env2 并引发已保存的浮点异常
	feupdateenv(&env2);
	show_fe_exceptions(false);
	show_fe_rounding();

	// 复制 flag
	fesetexceptflag(&ex, FE_ALL_EXCEPT);
	show_fe_exceptions(false);
	show_fe_rounding();

	// 恢复初始浮点环境
	fesetenv(&env1);
	show_fe_exceptions(false);
	show_fe_rounding();
}

void test_fenv(void) {
	example_fe_exceptions();
	example_fe_rounding();
	example_fe_environment();
}
/*
>>> Floating-Point Exception Testing
1.0 / 0.0 = inf
current floating-point exception raised:   FE_DIVBYZERO
sqrt(-1.0) = -nan(ind)
current floating-point exception raised:   FE_INVALID
exp(1000.0) = inf
current floating-point exception raised:   FE_OVERFLOW  FE_INEXACT
exp(-1000.0) = 0.000000e+00
current floating-point exception raised:   FE_UNDERFLOW  FE_INEXACT
1.0 / 3.0 = 0.33333333333333331483
current floating-point exception raised:   FE_INEXACT

>>> Rounding Direction Testing
FE_TONEAREST (Round to nearest)    : 3.33333333333333303727 → 3
FE_DOWNWARD (Round toward -∞)    : 3.33333333333333303727 → 3
FE_UPWARD (Round toward +∞)      : 3.33333333333333348137 → 4
FE_TOWARDZERO (Round toward zero)  : 3.33333333333333303727 → 3

>>> Environment Control
current floating-point exception raised:   FE_DIVBYZERO  FE_INVALID
current rounding direction:  FE_DOWNWARD
current floating-point exception raised:   NONE
current rounding direction:  FE_DOWNWARD
current floating-point exception raised:   FE_OVERFLOW  FE_INEXACT
current rounding direction:  FE_UPWARD
current floating-point exception raised:   FE_DIVBYZERO  FE_INVALID  FE_OVERFLOW  FE_INEXACT
current rounding direction:  FE_DOWNWARD
current floating-point exception raised:   FE_INVALID
current rounding direction:  FE_DOWNWARD
current floating-point exception raised:   NONE
current rounding direction:  FE_TONEAREST
*/