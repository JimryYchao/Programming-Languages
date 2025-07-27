#pragma once
typedef  int _defined_;
typedef struct _defined_ imaxdiv_t;
typedef _defined_ intmax_t;		// 最大整数类型
typedef _defined_ uintmax_t;	// 最大无符号整数类型
typedef _defined_ wchar_t;	    // 宽字符类型

// 格式化输出说明符，PRI_M_N: M = b,B,d,i,o,u,x,X; N = 8,16,32,64
#define PRIi8       // 有符号整数
#define PRIb8       // 无符号二进制
#define PRIB8       // 无符号二进制（可选）
#define PRIo8       // 无符号八进制整数
#define PRId8       // 有符号十进制整数
#define PRIu8       // 无符号十进制整数
#define PRIx8       // 无符号十六进制整数（小写）
#define PRIX8       // 无符号十六进制整数（大写）
// 格式化输入说明符，SCN_M_N: M = b,d,i,o,u,x,; N = 8,16,32,64
#define SCNi8       // 有符号整数
#define SCNu8       // 无符号整数
#define SCNb8       // 无符号二进制整数
#define SCNo8       // 无符号八进制整数
#define SCNd8       // 有符号十进制整数
#define SCNx8       // 无符号十六进制整数

// 其他类型：PRI_M_T_N: M = b,B,d,i,o,u,x,X; N = 8,16,32,64
#define PRIdLEAST8  // int_least8_t
#define PRIdFAST8   // int_fast8_t
#define PRIdMAX     // intmax_t
#define PRIdPTR     // intptr_t
#define SCNdLEAST8  // int_least8_t
#define SCNdFAST8   // int_fast8_t
#define SCNdMAX     // intmax_t
#define SCNdPTR     // intptr_t

intmax_t imaxabs(intmax_t j);                        //  绝对值
imaxdiv_t imaxdiv(intmax_t numer, intmax_t denom);   // 计算商余
intmax_t strtoimax(const char * restrict nptr, char ** restrict endptr, int base);      // 字符串转整数
intmax_t wcstoimax(const wchar_t *restrict nptr, wchar_t **restrict endptr, int base);  // 宽字符串转整数
uintmax_t strtoumax(const char * restrict nptr, char ** restrict endptr, int base);     // 字符串转无符号整数
uintmax_t wcstoumax(const wchar_t *restrict nptr, wchar_t **restrict endptr, int base); // 宽字符串转无符号整数