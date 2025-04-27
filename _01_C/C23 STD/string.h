#pragma once
typedef unsigned char QChar;
typedef unsigned char QVoid;

// 字符串复制
char *strcpy(char *restrict s1, const char *restrict s2);            // 复制字符串
char *strncpy(char *restrict s1, const char *restrict s2, size_t n); // 复制字符串，最多 n 个字符
char *strcat(char *restrict s1, const char *restrict s2);            // 连接字符串
char *strncat(char *restrict s1, const char *restrict s2, size_t n); // 连接字符串，最多 n 个字符
char *strdup(const char *s);                                         // 分配新内存复制字符串，需要使用 free 释放内存
char *strndup(const char *s, size_t n);                              // 分配新内存复制字符串，最多 n 个字符，需要使用 free 释放内存
// 字符串校验
size_t strlen(const char *s);                          // 长度
char *strerror(int errnum);                            // 获取错误码对应的字符串描述
int strcmp(const char *s1, const char *s2);            // 比较字符串
int strncmp(const char *s1, const char *s2, size_t n); // 比较字符串，最多 n 个字符
int strcoll(const char *s1, const char *s2);           // 比较字符串，使用当前区域设置
size_t strxfrm(char *restrict s1,                      // 转换字符串，最多 n 个字符，使得 strcmp 会产生与 strcoll 相同的结果
               const char *restrict s2, size_t n);
QChar *strchr(QChar *s, int c);                           // 查找字符首次出现的位置
QChar *strrchr(QChar *s, int c);                          // 查找字符最后一次出现的位置
size_t strspn(const char *s1, const char *s2);            // 从起始检索字符集连续包含的字符数量
size_t strcspn(const char *s1, const char *s2);           // 从起始检索字符集连续不包含的字符数量
QChar *strpbrk(QChar *s1, const char *s2);                // 定位字符集中任何字符在 dst 字节串中的第一个匹配项
QChar *strstr(QChar *s1, const char *s2);                 // 查找子字符串首次出现的位置
char *strtok(char *restrict s1, const char *restrict s2); // 查找字节字符串中的下一个 token
// 字符组数操作
QVoid *memchr(QVoid *s, int c, size_t n);                                   // 查找字符首次出现的位置
int memcmp(const void *s1, const void *s2, size_t n);                       // 比较内存区域，最多 n 个字节
void *memcpy(void *restrict s1, const void *restrict s2, size_t n);         // 复制内存区域，最多 n 个字节
void *memccpy(void *restrict s1, const void *restrict s2, int c, size_t n); // 复制内存区域，最多 n 个字节，或直到遇到字符 c
void *memmove(void *s1, const void *s2, size_t n);                          // 复制内存区域，最多 n 个字节，允许重叠
void *memset(void *s, int c, size_t n);                                     // 设置内存区域，最多 n 个字节为字符 c
void *memset_explicit(void *s, int c, size_t n);                            // 设置内存区域，最多 n 个字节为字符 c，但对敏感信息是安全的

#define __STDC_LIB_EXT1__
#ifdef(__STDC_LIB_EXT1__ &&__STDC_WANT_LIB_EXT1__)
typedef int errno_t;
typedef size_t rsize_t;
errno_t memcpy_s(void *restrict s1, rsize_t s1max, const void *restrict s2, rsize_t n);
errno_t memmove_s(void *s1, rsize_t s1max, const void *s2, rsize_t n);
errno_t strcpy_s(char *restrict s1, rsize_t s1max, const char *restrict s2);
errno_t strncpy_s(char *restrict s1, rsize_t s1max, const char *restrict s2, rsize_t n);
errno_t strcat_s(char *restrict s1, rsize_t s1max, const char *restrict s2);
errno_t strncat_s(char *restrict s1, rsize_t s1max, const char *restrict s2, rsize_t n);
char *strtok_s(char *restrict s1, rsize_t *restrict s1max, const char *restrict s2, char **restrict ptr);
errno_t memset_s(void *s, rsize_t smax, int c, rsize_t n);
errno_t strerror_s(char *s, rsize_t maxsize, errno_t errnum);
size_t strerrorlen_s(errno_t errnum);
size_t strnlen_s(const char *s, size_t maxsize);
#endif