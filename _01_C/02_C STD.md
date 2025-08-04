### C Standard Libraries

- [assert.h](https://zh.cppreference.com/w/c/header/assert.html)           : 程序断言，定义了宏 `assert` 并引用了 `NDEBUG` 用于取消断言。                                   
- [complex.h](https://zh.cppreference.com/w/c/header/complex.html)         : 复数类型和运算支持。                                                         
- [ctype.h](https://zh.cppreference.com/w/c/header/ctype.html)             : 字符处理，提供用于分类（`is*`）和映射字符（`to*`）的函数。                                     
- [errno.h](https://zh.cppreference.com/w/c/header/errno.html)             : 错误报告，定义错误条件宏。                                                                     
- [fenv.h](https://zh.cppreference.com/w/c/header/fenv.html)               : 浮点环境，提供对线程局域浮点环境的访问和控制。                                                 
- [float.h](https://zh.cppreference.com/w/c/header/float.html)             : 浮点类型标准，提供常规浮点类型和十进制浮点类型的限制。                                         
- [inttypes.h](https://zh.cppreference.com/w/c/header/inttypes.html)       : 整型格式转换，提供转换函数和格式化输入输出预定义宏。                                           
- [iso646.h](https://zh.cppreference.com/w/c/header/iso646.html)           : 替用拼写。                                                                                     
- [limits.h](https://zh.cppreference.com/w/c/header/limits.html)           : 整数类型限制，定义整数类型的宽度和范围。                                                       
- [locale.h](https://zh.cppreference.com/w/c/header/locale.html)           : 本地化支持，提供查询和设置区域支持的方法。                                                     
- [math.h](https://zh.cppreference.com/w/c/header/math.html)               : 数学库。                                                                                       
- [setjmp.h](https://zh.cppreference.com/w/c/header/setjmp.html)           : 非局部跳转，用于绕过正常的函数调用和返回规则。                                                 
- [signal.h](https://zh.cppreference.com/w/c/header/signal.html)           : 信号处理，提供用于信号管理的函数和宏常量。                                                     
- [stdarg.h](https://zh.cppreference.com/w/c/header/stdarg.html)           : 可变参数，提供对变参函数的支持。                                                               
- [stdatomic.h](https://zh.cppreference.com/w/c/header/stdatomic.html)     : 原子操作，提供原子操作、互斥、条件变量的内建支持。                                             
- [stdbit.h](https://zh.cppreference.com/w/c/header/stdbit.html)           : 位与字节工具，提供处理各类型的字节和位表示的宏函数。                                           
- [stdckdint.h](https://zh.cppreference.com/w/c/header/stdckdint.html)     : 校验整数算术，提供用于执行校验整数算术的宏函数。                                               
- [stddef.h](https://zh.cppreference.com/w/c/header/stddef.html)           : 常用定义。                                                                                     
- [stdint.h](https://zh.cppreference.com/w/c/header/stdint.html)           : 定宽整型，指定宽度的整数类型和对应的整数限制。                                                 
- [stdio.h](https://zh.cppreference.com/w/c/header/stdio.html)             : 输入与输出，提供了通用的文件操作支持，并提供了具有窄字符输入和输出功能的函数。                 
- [stdlib.h](https://zh.cppreference.com/w/c/header/stdlib.html)           : 基础工具库，例如内存管理、程序工具、字符串转换、随机数、算法等。                               
- [string.h](https://zh.cppreference.com/w/c/header/string.html)           : 字符串处理，提供一些操作字符类型的数组和字符串的函数等。                                       
- [tgmath.h](https://zh.cppreference.com/w/c/header/tgmath.html)           : 泛型数学，整合 math.h 和 complex.h。                                                           
- [threads.h](https://zh.cppreference.com/w/c/header/threads.html)         : 线程库，提供了对线程、互斥、条件变量和线程专有存储的支持。                                     
- [time.h](https://zh.cppreference.com/w/c/header/time.html)               : 时间与日期，提供与时间相关的类型和函数。                                                       
- [uchar.h](https://zh.cppreference.com/w/c/header/uchar.html)             : Unicode 字符工具，提供多字节和 UTF-16、UTF-32 之间的转换函数。                                 
- [wchar.h](https://zh.cppreference.com/w/c/header/wchar.html)             : 扩展多字节和宽字符工具。                                                                       
- [wctype.h](https://zh.cppreference.com/w/c/header/wctype.html)           : 宽字符处理，提供用于宽字符分类和字符映射的函数和类型，这些函数受当前区域设置 `LC_CTYPE` 影响。 
- [stdaligh.h](https://zh.cppreference.com/w/c/header/stdalign.html)       : 对齐，C23 起弃用                                                                               
- [stdbool.h](https://zh.cppreference.com/w/c/header/stdbool.html)         : 布尔定义，C23 起弃用。                                                                         
- [stdnoreturn.h](https://zh.cppreference.com/w/c/header/stdnoreturn.html) : _Noreturn，C23 起弃用。                                                                        

<!-- 
- stdmchar C29 ??
-->

>---
#### Examples

- [examples](./C%20STD/README.md)

---