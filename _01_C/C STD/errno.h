#pragma once

#define errno      // 展开为 POSIX 兼容的线程错误编号变量, 一些标准库函数通过写入正整数到 errno 指定错误。
// 错误编号宏
#define EDOM        // 数学参数在函数定义域外
#define EILSEQ      // 非法字节序列
#define ERANGE      // 数学结果超出范围

#define __STDC_LIB_EXT1__
#ifdef __STDC_WANT_LIB_EXT1__
typedef __errno_t   errno_t;
#endif


void Example(){
    errno = 0;
    printf("log(-1.0) = %f\n", log(-1.0));  // log(-1.0) = nan
    printf("%s\n", strerror(errno));        // Numerical argument out of domain
}