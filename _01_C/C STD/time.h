#pragma once
typedef int _defined_;

typedef _defined_ clock_t; // 程序运行时间
typedef _defined_ time_t;  // 从纪元开始的日历时间
struct timespec            // 保存以秒和纳秒指定的间隔
{
    time_t tv_sec;
    signed int tv_nsec;
};
struct tm // 分解时间，包含年、月、日、时、分、秒等信息
{
    int tm_sec;   // 秒数 [0, 60]，允许闰秒
    int tm_min;   // 分钟数 [0, 59]
    int tm_hour;  // 小时数 [0, 23]
    int tm_mday;  // 一个月中的天数 [1, 31]
    int tm_mon;   // 月份 [0, 11]，0 表示 1 月
    int tm_year;  // 年份，从 1900 年开始的偏移量
    int tm_wday;  // 星期几 [0, 6]，0 表示星期天
    int tm_yday;  // 一年中的天数 [0, 365]，0表示1月1日
    int tm_isdst; // 夏令时标志，正值表示夏令时，负值表示非夏令时，零表示未知
};

#define NULL ((void *)0)
#define CLOCKS_PER_SEC ((clock_t)1000)
#define TIME_UTC           // 指定日历时基
#define TIME_MONOTONIC     // 指定单调时基
#define TIME_ACTIVE        // 指定活动时基
#define TIME_THREAD_ACTIVE // 指定线程活动时基

// 时间操作
clock_t clock(void);                                // 进程运行时间
time_t time(time_t *timer);                         // 获取当前日历时间
double difftime(time_t time1, time_t time0);        // 日历时间差，秒
time_t mktime(struct tm *timeptr);                  // 将 tm 转换从纪元开始的日历时间
time_t timegm(struct tm *timeptr);                  // 将 tm 转换为 UTC 日历时间
int timespec_get(struct timespec *ts, int base);    // 保存以时间基 base 表示的当前日历时间
int timespec_getres(struct timespec *ts, int base); // 获取时间基 base 的分辨率
// 时间格式化
struct tm *gmtime(const time_t *timer); // 将日历时间转换为 UTC 时间
struct tm *gmtime_r(const time_t *timer, struct tm *buf);
struct tm *localtime(const time_t *timer); // 将日历时间转换为本地时间
struct tm *localtime_r(const time_t *timer, struct tm *buf);
size_t strftime(char *restrict s, size_t maxsize, const char *restrict format, const struct tm *restrict timeptr); // 格式化时间

[[deprecated]] char *asctime(const struct tm *timeptr); // tm 转换文本 Www Mmm dd hh:mm:ss yyyy\n
errno_t asctime_s(char *buf, rsize_t bufsz, const struct tm *time_ptr);
[[deprecated]] char *ctime(const time_t *timer); // time_t 转换文本 Www Mmm dd hh:mm:ss yyyy\n
errno_t ctime_s(char *buf, rsize_t bufsz, const time_t *timer);

// Example: 获取本地当前时间并格式化输出
void get_local_time()
{
    time_t now = time(NULL);             // 获取当前时间
    struct tm *tm_now = localtime(&now); // 将时间转换为本地时间结构体
    char buff[100] = {0};
    strftime(buff, sizeof buff, "%D %T", tm_now); // 格式化时间
    printf("当前时间: %s UTC+8\n", buff);         // 当前时间: 01/01/25 00:00:00 UTC+8
}