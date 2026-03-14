## C STD (C23)

---

| Header                                                               | Description                                                      | Example                                            |
| :------------------------------------------------------------------- | :--------------------------------------------------------------- | :------------------------------------------------- |
| [assert.h](https://zh.cppreference.com/w/c/header/assert.html)       | 程序断言，`NDEBUG` 禁用断言。                                    |                                                    |
| [complex.h](https://zh.cppreference.com/w/c/header/complex.html)     | 复数支持。                                                       |                                                    |
| [ctype.h](https://zh.cppreference.com/w/c/header/ctype.html)         | 字符处理，提供分类（`is*`）和映射（`to*`）函数。                 |
| [errno.h](https://zh.cppreference.com/w/c/header/errno.html)         | 错误报告，定义错误条件宏。                                       | [[↗]](./C%20Libs/CSTD_examples/errno_test.c)       |
| [fenv.h](https://zh.cppreference.com/w/c/header/fenv.html)           | 浮点环境，提供对线程局域浮点环境的访问和控制。                   | [[↗]](./C%20Libs/CSTD_examples/fenv_test.c)        |
| [float.h](https://zh.cppreference.com/w/c/header/float.html)         | 浮点类型标准，提供常规浮点类型和十进制浮点类型的限制。           |
| [inttypes.h](https://zh.cppreference.com/w/c/header/inttypes.html)   | 整型格式转换，提供转换函数和格式化输入输出预定义宏。             |                                                    |
| [iso646.h](https://zh.cppreference.com/w/c/header/iso646.html)       | 替用拼写。                                                       |                                                    |
| [limits.h](https://zh.cppreference.com/w/c/header/limits.html)       | 整数类型限制。                                                   |                                                    |
| [locale.h](https://zh.cppreference.com/w/c/header/locale.html)       | 本地化支持，提供查询和设置区域支持的方法。                       | [[↗]](./C%20Libs/CSTD_examples/locale_test.c)      |
| [math.h](https://zh.cppreference.com/w/c/header/math.html)           | 数学库。                                                         |                                                    |
| [setjmp.h](https://zh.cppreference.com/w/c/header/setjmp.html)       | 非局部跳转。                                                     | [[↗]](./C%20Libs/CSTD_examples/setjmp_test.c)      |
| [signal.h](https://zh.cppreference.com/w/c/header/signal.html)       | 信号处理。                                                       | [[↗]](./C%20Libs/CSTD_examples/signal_test.c)      |
| [stdarg.h](https://zh.cppreference.com/w/c/header/stdarg.html)       | 可变参数支持。                                                   |                                                    |
| [stdatomic.h](https://zh.cppreference.com/w/c/header/stdatomic.html) | 原子操作支持。                                                   | [[↗]](./C%20Libs/CSTD_examples/stdatomic_test.cpp) |
| [stdbit.h](https://zh.cppreference.com/w/c/header/stdbit.html)       | 位与字节工具。                                                   |                                                    |
| [stdckdint.h](https://zh.cppreference.com/w/c/header/stdckdint.html) | 校验整数算术。                                                   |                                                    |
| [stddef.h](https://zh.cppreference.com/w/c/header/stddef.html)       | 常用定义。                                                       |                                                    |
| [stdint.h](https://zh.cppreference.com/w/c/header/stdint.html)       | 定宽整数类型支持。                                               |                                                    |
| [stdio.h](https://zh.cppreference.com/w/c/header/stdio.html)         | 输入与输出，通用文件操作支持，窄字符 IO 支持。                   | [[↗]](./C%20Libs/CSTD_examples/stdio_test.c)       |
| [stdlib.h](https://zh.cppreference.com/w/c/header/stdlib.html)       | 基础工具库，例如内存管理、程序工具、字符串转换、随机数、算法等。 | [[↗]](./C%20Libs/CSTD_examples/stdlib_test.c)      |
| [string.h](https://zh.cppreference.com/w/c/header/string.html)       | 字符串处理。                                                     | [[↗]](./C%20Libs/CSTD_examples/string_test.c)      |
| [tgmath.h](https://zh.cppreference.com/w/c/header/tgmath.html)       | 泛型数学，整合 math.h 和 complex.h。                             |                                                    |
| [threads.h](https://zh.cppreference.com/w/c/header/threads.html)     | 线程库，提供了对线程、互斥、条件变量和线程专有存储的支持。       | [[↗]](./C%20Libs/CSTD_examples/threads_test.c)     |
| [time.h](https://zh.cppreference.com/w/c/header/time.html)           | 时间与日期工具。                                                 | [[↗]](./C%20Libs/CSTD_examples/time_test.c)        |
| [uchar.h](https://zh.cppreference.com/w/c/header/uchar.html)         | Unicode 字符工具，提供多字节和 UTF-16、UTF-32 之间的转换函数。   | [[↗]](./C%20Libs/CSTD_examples/uchar_test.c)       |
| [wchar.h](https://zh.cppreference.com/w/c/header/wchar.html)         | 扩展多字节和宽字符工具。                                         |                                                    |
| [wctype.h](https://zh.cppreference.com/w/c/header/wctype.html)       | 宽字符处理，当前区域设置 `LC_CTYPE` 影响。                       |

> 弃用 
     
- [stdalign.h](https://zh.cppreference.com/w/c/header/stdalign.html)                                                                                      
- [stdbool.h](https://zh.cppreference.com/w/c/header/stdbool.html)                                                                                                                    
- [stdnoreturn.h](https://zh.cppreference.com/w/c/header/stdnoreturn.html)                                                                                                                 

<!-- 
- stdmchar C29 ??
-->

---