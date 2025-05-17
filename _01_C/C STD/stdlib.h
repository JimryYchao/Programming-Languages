#pragma once
typedef int _defined_;

#define NULL // 空指针常量
void call_once(once_flag *flag, void (*func)(void)); // 单次调用函数
typedef _defined_ once_flag;                         // 足以保有 call_once 所用标志的完整对象类型
#define ONCE_FLAG_INIT                               // once_flag 初始化

// 字符串转换数值
double atof(const char *nptr);
int atoi(const char *nptr);
long int atol(const char *nptr);
long long int atoll(const char *nptr);
double strtod(const char *restrict nptr, char **restrict endptr);
float strtof(const char *restrict nptr, char **restrict endptr);
long double strtold(const char *restrict nptr, char **restrict endptr);
long int strtol(const char *restrict nptr, char **restrict endptr, int base);
long long int strtoll(const char *restrict nptr, char **restrict endptr, int base);
unsigned long int strtoul(const char *restrict nptr, char **restrict endptr, int base);
unsigned long long int strtoull(const char *restrict nptr, char **restrict endptr, int base);
// 数值转字符串，format：a、A、e、E、f、F、g、G 和可选的精度
int strfromd(char *restrict s, size_t n, const char *restrict format, double fp);
int strfromf(char *restrict s, size_t n, const char *restrict format, float fp);
int strfroml(char *restrict s, size_t n, const char *restrict format, long double fp);
// 随机数生成
int rand(void);
void srand(unsigned int seed);
#define RAND_MAX // rand 返回的最大值
// 内存管理
typedef _defined_ size_t;
void *malloc(size_t size);                                         // 分配未初始化内存
void *calloc(size_t nmemb, size_t size);                           // 分配清零内存
void *realloc(void *ptr, size_t size);                             // 扩充内存
void *aligned_alloc(size_t alignment, size_t size);                // 分配对齐内存
void free(void *ptr);                                              // 释放内存
void free_sized(void *ptr, size_t size);                           // 释放 malloc、realloc、calloc 内存，size 是分配时的大小
void free_aligned_sized(void *ptr, size_t alignment, size_t size); // 释放 aligned_alloc 内存，alignment 是分配时的对齐大小
size_t memalignment(const void *p);                                // 获取地址值所满足的最大对齐
// 程序环境支持
#define EXIT_FAILURE
#define EXIT_SUCCESS
void abort(void) [[noreturn]];            // 中止程序
void _Exit(int status) [[noreturn]];      // 仅正常终止不清理
void exit(int status) [[noreturn]];       // 正常终止并清理
int atexit(void (*func)(void));           // 注册 exit 调用的函数
void quick_exit(int status) [[noreturn]]; // 快速终止并不完全清理
int at_quick_exit(void (*func)(void));    // 注册 quick_exit 调用的函数
char *getenv(const char *name);           // 访问环境变量
int system(const char *string);           // 以参数 string 调用主机命令
// 搜索与排序
QVoid *bsearch(const void *key, // 二分查找升序序列中与 key 匹配的元素位置
               QVoid *base, size_t nmemb, size_t size, int (*compar)(const void *, const void *));
void qsort(void *base, size_t nmemb, // 以 compar 设定的比较规则进行排序
           size_t size, int (*compar)(const void *, const void *));
// 整数算术
typedef _defined_ div_t; // 整数除法的商和余数
typedef _defined_ ldiv_t;
typedef _defined_ lldiv_t;
int abs(int j); // 绝对值
long int labs(long int j);
long long int llabs(long long int j);
div_t div(int numer, int denom); // 商和余数
ldiv_t ldiv(long int numer, long int denom);
lldiv_t lldiv(long long int numer, long long int denom);
// 多字节/宽字符（串）转换
#define MB_CUR_MAX // 当前区域设置的多字节字符中的最大字节数
typedef _defined_ wchar_t;
int wctomb(char *s, wchar_t wc);                                           // 字符，宽字符转多字节
int mbtowc(wchar_t *restrict pwc, const char *restrict s, size_t n);       // 字符，多字节转宽字符
int mblen(const char *s, size_t n);                                        // 确定 s 指向其首字节的多字节字符的字节大小
size_t wcstombs(char *restrict s, const wchar_t *restrict pwcs, size_t n); // 字符串，宽字符转多字节
size_t mbstowcs(wchar_t *restrict pwcs, const char *restrict s, size_t n); // 字符串，多字节转宽字符

#ifdef __STDC_IEC_60559_DFP__
int strfromd32(char *restrict s, size_t n, const char *restrict format, _Decimal32 fp);
int strfromd64(char *restrict s, size_t n, const char *restrict format, _Decimal64 fp);
int strfromd128(char *restrict s, size_t n, const char *restrict format, _Decimal128 fp);
_Decimal32 strtod32(const char *restrict nptr, char **restrict endptr);
_Decimal64 strtod64(const char *restrict nptr, char **restrict endptr);
_Decimal128 strtod128(const char *restrict nptr, char **restrict endptr);
#endif
#define __STDC_LIB_EXT1__
#ifdef(__STDC_LIB_EXT1__ &&__STDC_WANT_LIB_EXT1__)
typedef int errno_t;
typedef size_t rsize_t;
typedef int constraint_handler_t;
constraint_handler_t set_constraint_handler_s(constraint_handler_t handler);
void abort_handler_s(const char *restrict msg, void *restrict ptr, errno_t error);
void ignore_handler_s(const char *restrict msg, void *restrict ptr, errno_t error);
errno_t getenv_s(size_t *restrict len, char *restrict value, rsize_t maxsize, const char *restrict name);
QVoid *bsearch_s(const void *key, QVoid *base, rsize_t nmemb, rsize_t size, int (*compar)(const void *k, const void *y, void *context), void *context);
errno_t qsort_s(void *base, rsize_t nmemb, rsize_t size, int (*compar)(const void *x, const void *y, void *context), void *context);
errno_t wctomb_s(int *restrict status, char *restrict s, rsize_t smax, wchar_t wc);
errno_t mbstowcs_s(size_t *restrict retval, wchar_t *restrict dst, rsize_t dstmax, const char *restrict src, rsize_t len);
errno_t wcstombs_s(size_t *restrict retval, char *restrict dst, rsize_t dstmax, const wchar_t *restrict src, rsize_t len);
#endif