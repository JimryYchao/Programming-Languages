#pragma once
typedef  int _defined_;

typedef _defined_ wchar_t;      // 宽字符类型
typedef _defined_ size_t;       
typedef _defined_ mbstate_t;    // 多字节状态类型
typedef _defined_ wint_t;       // 保存宽字符整数类型

typedef struct  tm

#define NULL            // 空指针常量
#define WCHAR_MAX       // 宽字符最大值
#define WCHAR_MIN       // 宽字符最小值
#define WEOF            // 宽字符流结束符

// 宽字符格式化输出
int wprintf(const wchar_t *restrict format, ...);
int fwprintf(FILE *restrict stream, const wchar_t *restrict format, ...);
int swprintf(wchar_t *restrict s, size_t n, const wchar_t *restrict format, ...);
int vfwprintf(FILE *restrict stream, const wchar_t *restrict format, va_list arg);
int vswprintf(wchar_t *restrict s, size_t n, const wchar_t *restrict format, va_list arg);
int vwprintf(const wchar_t *restrict format, va_list arg);
// 宽字符格式化输入
int wscanf(const wchar_t *restrict format, ...);
int fwscanf(FILE *restrict stream, const wchar_t *restrict format, ...);
int swscanf(const wchar_t *restrict s, const wchar_t *restrict format, ...);
int vwscanf(const wchar_t *restrict format, va_list arg);
int vfwscanf(FILE *restrict stream, const wchar_t *restrict format, va_list arg);
int vswscanf(const wchar_t *restrict s, const wchar_t *restrict format, va_list arg);
// 宽字符输入与输出
wint_t fputwc(wchar_t c, FILE *stream);
wint_t putwc(wchar_t c, FILE *stream);
wint_t putwchar(wchar_t c);
wint_t fgetwc(FILE *stream);
wint_t getwc(FILE *stream);
wint_t getwchar(void);
// 宽字符串输入与输出
int fputws(const wchar_t *restrict s, FILE *restrict stream);
wchar_t *fgetws(wchar_t *restrict s, int n, FILE *restrict stream);
// 流操作函数
int fwide(FILE *stream, int mode);
wint_t ungetwc(wint_t c, FILE *stream);

// 宽字符转换数值
double wcstod(const wchar_t *restrict nptr, wchar_t **restrict endptr);             
float wcstof(const wchar_t *restrict nptr, wchar_t **restrict endptr);              
long double wcstold(const wchar_t *restrict nptr, wchar_t **restrict endptr);       
#ifdef __STDC_IEC_60559_DFP__   
_Decimal32 wcstod32(const wchar_t *restrict nptr, wchar_t **restrict endptr);       
_Decimal64 wcstod64(const wchar_t *restrict nptr, wchar_t **restrict endptr);       
_Decimal128 wcstod128(const wchar_t *restrict nptr, wchar_t **restrict endptr);     
#endif
long int wcstol(const wchar_t *restrict nptr, wchar_t **restrict endptr, int base);           
long long int wcstoll(const wchar_t *restrict nptr, wchar_t **restrict endptr, int base);     
unsigned long int wcstoul(const wchar_t *restrict nptr, wchar_t **restrict endptr, int base);        
unsigned long long int wcstoull(const wchar_t *restrict nptr, wchar_t **restrict endptr, int base);  

// 宽字符串操作
size_t wcslen(const wchar_t *s);     // 长度
wchar_t *wcscpy(wchar_t *restrict s1, const wchar_t *restrict s2);              // 字符串复制
wchar_t *wcsncpy(wchar_t *restrict s1, const wchar_t *restrict s2, size_t n);    
wchar_t *wcscat(wchar_t *restrict s1, const wchar_t *restrict s2);              // 追加字符串
wchar_t *wcsncat(wchar_t *restrict s1, const wchar_t *restrict s2, size_t n);  
size_t wcsftime(wchar_t *restrict s, size_t maxsize, const wchar_t *restrict format, const struct tm *restrict timeptr);  // 宽字符时间格式化转换
// 宽字符数组操作
wchar_t *wmemcpy(wchar_t *restrict s1, const wchar_t *restrict s2, size_t n);   // 在不重叠数组间复制
wchar_t *wmemmove(wchar_t *s1, const wchar_t *s2, size_t n);                    // 可能重叠的数组间复制
int wmemcmp(const wchar_t *s1, const wchar_t *s2, size_t n);                    // 数组比较
QWchar_t *wmemchr(QWchar_t *s, wchar_t c, size_t n);                            // 字符检索
wchar_t *wmemset(wchar_t *s, wchar_t c, size_t n);                              // 复制宽字符到数组
// 宽字符串比较
int wcscmp(const wchar_t *s1, const wchar_t *s2);
int wcsncmp(const wchar_t *s1, const wchar_t *s2, size_t n);
int wcscoll(const wchar_t *s1, const wchar_t *s2);                              // 在本地环境下比较
size_t wcsxfrm(wchar_t *restrict s1, const wchar_t *restrict s2, size_t n);     // 变换字符串，是的 wcscmp 产生与 wcscoll 相同的结果 
// 宽字符串检索
QWchar_t *wcschr(QWchar_t *s, wchar_t c);             // 首次出现
QWchar_t *wcsrchr(QWchar_t *s, wchar_t c);            // 最后一次出现
QWchar_t *wcsstr(QWchar_t *s1, const wchar_t *s2);    // 首个子串检索 
QWchar_t *wcspbrk(QWchar_t *s1, const wchar_t *s2);   // 首个任意字符检索
size_t wcsspn(const wchar_t *s1, const wchar_t *s2);  // 连续字符检索
size_t wcscspn(const wchar_t *s1, const wchar_t *s2); // 不连续字符检索
wchar_t *wcstok(wchar_t *restrict s1, const wchar_t *restrict s2, wchar_t **restrict ptr);   // 查找宽字符串中的下一个 token
// 多字节/宽字符转换
wint_t btowc(int c);        // 单字节转宽字符
int wctob(wint_t c);        // 宽字符转单字节
int mbsinit(const mbstate_t *ps);     // 初始状态检查
size_t mbrlen(const char *restrict s, size_t n, mbstate_t *restrict ps);        // 剩余多字节字符长度
size_t mbrtowc(wchar_t *restrict pwc, const char *restrict s, size_t n, mbstate_t *restrict ps);        // 指定状态多字节转宽字符
size_t wcrtomb(char *restrict s, wchar_t wc, mbstate_t *restrict ps);                                   // 指定状态宽字符转多字节
size_t mbsrtowcs(wchar_t *restrict dst, const char **restrict src, size_t len, mbstate_t *restrict ps); // 多字节转换宽字符串
size_t wcsrtombs(char *restrict dst, const wchar_t **restrict src, size_t len, mbstate_t *restrict ps); // 宽字符串转换多字节


#define __STDC_LIB_EXT1__    // 线程安全扩展
#ifdef(__STDC_LIB_EXT1__ &&__STDC_WANT_LIB_EXT1__)
typedef __errno_t errno_t;
typedef size_t rsize_t;
int fwprintf_s(FILE *restrict stream, const wchar_t *restrict format, ...);
int fwscanf_s(FILE *restrict stream, const wchar_t *restrict format, ...);
int snwprintf_s(wchar_t *restrict s, rsize_t n, const wchar_t *restrict format, ...);
int swprintf_s(wchar_t *restrict s, rsize_t n, const wchar_t *restrict format, ...);
int swscanf_s(const wchar_t *restrict s, const wchar_t *restrict format, ...);
int vfwprintf_s(FILE *restrict stream, const wchar_t *restrict format, va_list arg);
int vfwscanf_s(FILE *restrict stream, const wchar_t *restrict format, va_list arg);
int vsnwprintf_s(wchar_t *restrict s, rsize_t n, const wchar_t *restrict format, va_list arg);
int vswprintf_s(wchar_t *restrict s, rsize_t n, const wchar_t *restrict format, va_list arg);
int vswscanf_s(const wchar_t *restrict s, const wchar_t *restrict format, va_list arg);
int vwprintf_s(const wchar_t *restrict format, va_list arg);
int vwscanf_s(const wchar_t *restrict format, va_list arg);
int wprintf_s(const wchar_t *restrict format, ...);
int wscanf_s(const wchar_t *restrict format, ...);
errno_t wcscpy_s(wchar_t *restrict s1, rsize_t s1max, const wchar_t *restrict s2);
errno_t wcsncpy_s(wchar_t *restrict s1, rsize_t s1max, const wchar_t *restrict s2, rsize_t n);
errno_t wmemcpy_s(wchar_t *restrict s1, rsize_t s1max, const wchar_t *restrict s2, rsize_t n);
errno_t wmemmove_s(wchar_t *s1, rsize_t s1max, const wchar_t *s2, rsize_t n);
errno_t wcscat_s(wchar_t *restrict s1, rsize_t s1max, const wchar_t *restrict s2);
errno_t wcsncat_s(wchar_t *restrict s1, rsize_t s1max, const wchar_t *restrict s2, rsize_t n);
wchar_t *wcstok_s(wchar_t *restrict s1, rsize_t *restrict s1max, const wchar_t *restrict s2, wchar_t **restrict ptr);
size_t wcsnlen_s(const wchar_t *s, size_t maxsize);
errno_t wcrtomb_s(size_t *restrict retval, char *restrict s, rsize_t smax, wchar_t wc, mbstate_t *restrict ps);
errno_t mbsrtowcs_s(size_t *restrict retval, wchar_t *restrict dst, rsize_t dstmax, const char **restrict src, rsize_t len, mbstate_t *restrict ps);
errno_t wcsrtombs_s(size_t *restrict retval, char *restrict dst, rsize_t dstmax, const wchar_t **restrict src, rsize_t len, mbstate_t *restrict ps);
#endif

// Example: wcsftime
void _wcsftime(void)
{
	wchar_t buff[40];
	struct tm my_time = {};
	time_t curtime = time(NULL);
	gmtime_s(&my_time, &curtime);
	wcsftime(buff, sizeof buff, L"%A %c", &my_time);
	printf("UTC: %ls\n", buff);   // UTC: Tuesday 01/01/25 00:00:00
}