#pragma once
typedef int _defined_;

typedef _defined_ size_t;
#define NULL

// 文件操作
typedef _defined_ FILE;                       // 文件，一种对象类型，能够记录控制流所需的所有信息
int remove(const char *filename);             // 删除文件
int rename(const char *old, const char *new); // 重命名文件
char *tmpnam(char *s);                        // 生成唯一的临时文件名
#define L_tmpnam                              // 临时文件名的最大长度
#define FILENAME_MAX                          // 文件名的最大长度
FILE *tmpfile(void);                          // 创建临时文件并打开，关闭或程序终止时自动删除
#define TMP_MAX                               // 临时文件的最大数量
// 文件访问
#define FOPEN_MAX // 同时打开的最大文件数
#define EOF       // 文件结束标志
#define stdin
#define stdout
#define stderr
FILE *fopen(const char *restrict filename, const char *restrict mode); // 以 mode 打开文件
int fclose(FILE *stream);                                              // 关闭文件，接触流与文件关联
int fflush(FILE *stream);                                              // 刷新文件流，强制将缓冲区的内容写入文件
FILE *freopen(const char *restrict filename,                           // 尝试关闭并以指定模式重新打开文件
              const char *restrict mode, FILE *restrict stream);
void setbuf(FILE *restrict stream, char *restrict buf); // 设置流操作内部缓冲区，buf = NULL 无缓冲，buf != NULL 全缓冲
#define BUFSIZ                                          // setbuf 使用的缓冲区大小
int setvbuf(FILE *restrict stream,                      // 设置缓冲区
            char *restrict buf, int mode, size_t size);
#define _IOFBF // 全缓冲
#define _IOLBF // 行缓冲
#define _IONBF // 无缓冲
size_t fread(void *restrict ptr,
             size_t size, size_t nmemb, FILE *restrict stream); // 从文件读取
size_t fwrite(const void *restrict ptr,
              size_t size, size_t nmemb, FILE *restrict stream); // 向文件写入
// 格式化输入与输出
#define _PRINTF_NAN_LEN_MAX                                                 // printf 输出 NaN 的最大长度
int printf(const char *restrict format, ...);                               // 格式化输出到 stdout
int sprintf(char *restrict s, const char *restrict format, ...);            // 格式化输出到字符串
int snprintf(char *restrict s, size_t n, const char *restrict format, ...); // 格式化输出到字符串，最多 n 个字符
int fprintf(FILE *restrict stream, const char *restrict format, ...);       // 格式化输出到指定文件流
int vprintf(const char *restrict format, va_list arg);
int vsprintf(char *restrict s, const char *restrict format, va_list arg);
int vsnprintf(char *restrict s, size_t n, const char *restrict format, va_list arg);
int vfprintf(FILE *restrict stream, const char *restrict format, va_list arg);
int scanf(const char *restrict format, ...);                          // 格式化输入
int sscanf(const char *restrict s, const char *restrict format, ...); // 格式化输入到字符串
int fscanf(FILE *restrict stream, const char *restrict format, ...);  // 格式化输入到指定文件流
int vscanf(const char *restrict format, va_list arg);
int vsscanf(const char *restrict s, const char *restrict format, va_list arg);
int vfscanf(FILE *restrict stream, const char *restrict format, va_list arg);
// 字符输入与输出
int putchar(int c);             // 输出字符到 stdout
int puts(const char *s);        // 输出字符串到 stdout
int fputc(int c, FILE *stream); // 输出字符到指定文件流
int putc(int c, FILE *stream);
int fputs(const char *restrict s, FILE *restrict stream); // 输出字符串到指定文件流

int ungetc(int c, FILE *stream);  // 将字符放回到指定文件流
int getchar(void);                // 从 stdin 读取字符
char *gets_s(char *s, rsize_t n); // 从 stdin 读取字符串，最多 n 个字符
int fgetc(FILE *stream);          // 从指定文件流读取字符
int getc(FILE *stream);
char *fgets(char *restrict s, int n, FILE *restrict stream); // 从指定文件流读取字符串，最多 n 个字符
// 文件定位
long int ftell(FILE *stream);                             // 获取文件流的当前位置
int fseek(FILE *stream, long int offset, int whence);     // 设置文件流的位置
#define SEEK_SET                                          // 从文件开头偏移
#define SEEK_CUR                                          // 从当前位置偏移
#define SEEK_END                                          // 从文件末尾偏移
void rewind(FILE *stream);                                // 将文件流位置设置到开头
typedef _defined_ fpos_t;                                 // 记录指定 FILE 的位置和多字节剖析状态
int fgetpos(FILE *restrict stream, fpos_t *restrict pos); // 获取文件流的位置
int fsetpos(FILE *stream, const fpos_t *pos);             // 恢复文件流的位置
// 错误处理
int feof(FILE *stream);      // 检查是否到达给定文件流的结尾
void clearerr(FILE *stream); // 清除文件流的错误标志
int ferror(FILE *stream);    // 检查文件流的错误标志
void perror(const char *s);  // 打印 errno 错误信息，s 为前缀字符串


#define __STDC_LIB_EXT1__
#ifdef(__STDC_LIB_EXT1__ &&__STDC_WANT_LIB_EXT1__)
#define L_tmpnam_s L_tmpnam
#define TMP_MAX_S TMP_MAX
typedef int errno_t;
typedef size_t rsize_t;
errno_t tmpfile_s(FILE *restrict *restrict streamptr);
errno_t tmpnam_s(char *s, rsize_t maxsize);
errno_t fopen_s(FILE *restrict *restrict streamptr, const char *restrict filename, const char *restrict mode);
errno_t freopen_s(FILE *restrict *restrict newstreamptr, const char *restrict filename, const char *restrict mode, FILE *restrict stream);
int fprintf_s(FILE *restrict stream, const char *restrict format, ...);
int fscanf_s(FILE *restrict stream, const char *restrict format, ...);
int printf_s(const char *restrict format, ...);
int scanf_s(const char *restrict format, ...);
int snprintf_s(char *restrict s, rsize_t n, const char *restrict format, ...);
int sprintf_s(char *restrict s, rsize_t n, const char *restrict format, ...);
int sscanf_s(const char *restrict s, const char *restrict format, ...);
int vfprintf_s(FILE *restrict stream, const char *restrict format, va_list arg);
int vfscanf_s(FILE *restrict stream, const char *restrict format, va_list arg);
int vprintf_s(const char *restrict format, va_list arg);
int vscanf_s(const char *restrict format, va_list arg);
int vsnprintf_s(char *restrict s, rsize_t n, const char *restrict format, va_list arg);
int vsprintf_s(char *restrict s, rsize_t n, const char *restrict format, va_list arg);
int vsscanf_s(const char *restrict s, const char *restrict format, va_list arg);
char *gets_s(char *s, rsize_t n);
#endif