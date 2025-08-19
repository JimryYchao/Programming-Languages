### CPP Standard Libraries

- [Examples]()

> 语言支持

| header                                                                                       | description                                            | examples                                                           |
| :------------------------------------------------------------------------------------------- | :----------------------------------------------------- | :----------------------------------------------------------------- |
| [cfloat](https://zh.cppreference.com/w/cpp/types/climits.html)                               | 浮点类型限制                                           |
| [climits](https://zh.cppreference.com/w/cpp/header/climits.html)                             | 整数类型限制                                           |
| TODO [compare](https://zh.cppreference.com/w/cpp/header/compare.html)                        | 三路比较运算符支持                                     | [[↗]](./CPP%20STD/examples/LanguageSupport/e_compare.cpp)          |
| [coroutine](https://zh.cppreference.com/w/cpp/header/coroutine.html)                         | 协程支持                                               | [[↗]](./CPP%20STD/examples/LanguageSupport/e_coroutine.cpp)        |
| [csetjmp](https://zh.cppreference.com/w/cpp/header/csetjmp.html)                             | 非局部跳转                                             |
| [csignal](https://zh.cppreference.com/w/cpp/header/csignal.html)                             | 信号处理                                               |
| [cstdarg](https://zh.cppreference.com/w/cpp/header/cstdarg.html)                             | 可变参数处理                                           |
| [cstddef](https://zh.cppreference.com/w/cpp/header/cstddef.html)                             | 常用宏和类型定义                                       | [[↗]](./CPP%20STD/examples/LanguageSupport/e_cstddef.cpp)          |
| [cstdlib](https://cppreference.cn/w/cpp/header/cstdlib)                                      | 通用工具库：程序控制、动态内存分配、随机数、排序和搜索 |
| [cstdint](https://zh.cppreference.com/w/cpp/header/cstdint.html)                             | 定宽整数及其限制                                       |
| [exception](https://zh.cppreference.com/w/cpp/header/exception.html)                         | 异常处理工具                                           | [[↗]](./CPP%20STD/examples/LanguageSupport/e_exception.cpp)        |
| [initializer_list](https://zh.cppreference.com/w/cpp/header/initializer_list.html)           | 列表初始化器模板化支持                                 | [[↗]](./CPP%20STD/examples/LanguageSupport/e_initializer_list.cpp) |
| [limits](https://zh.cppreference.com/w/cpp/header/limits.html)                               | 查询算术类型的属性                                     | [[↗]](./CPP%20STD/examples/LanguageSupport/e_limits.cpp)           |
| [new](https://zh.cppreference.com/w/cpp/header/new.html)                                     | 底层内存管理工具                                       | [[↗]](./CPP%20STD/examples/LanguageSupport/e_new.cpp)              |
| [source_location](https://zh.cppreference.com/w/cpp/header/source_location.html)             | 获取源代码位置                                         | [[↗]](./CPP%20STD/examples/LanguageSupport/e_source_location.cpp)  |
| [stdfloat](https://zh.cppreference.com/w/cpp/header/stdfloat.html)                           | 定宽浮点类型                                           |
| [typeindex](https://zh.cppreference.com/w/cpp/header/typeindex.html)                         | `std::type_info` 对象的包装器                          | [[↗]](./CPP%20STD/examples/LanguageSupport/e_typeindex.cpp)        |
| [typeinfo](https://zh.cppreference.com/w/cpp/header/typeinfo.html)                           | 运行时（typeid-expr）类型信息工具                      |                                                                    |
| [version](https://zh.cppreference.com/w/cpp/experimental/feature_test.html#Library_features) | 提供用于验证库实现状态的宏                             | [[↗]](./CPP%20STD/examples/LanguageSupport/e_version.cpp)          |
| [concepts](https://zh.cppreference.com/w/cpp/header/concepts.html)                           | 预定义概念库                                           |                                                                    |

<!-- [contracts]  契约支持 -->

> 诊断库

| header                                                                     | description                        | examples                                                   |
| :------------------------------------------------------------------------- | :--------------------------------- | :--------------------------------------------------------- |
| [cassert](https://zh.cppreference.com/w/cpp/header/cassert.html)           | 运行时断言                         |
| [cerrno](https://zh.cppreference.com/w/cpp/header/cerrno.html)             | 错误指示                           |
| [stacktrace](https://zh.cppreference.com/w/cpp/header/stacktrace.html)     | 堆栈追踪                           | [[↗]](./CPP%20STD/examples/Diagnostics/e_stacktrace.cpp)   |
| [stdexcept](https://zh.cppreference.com/w/cpp/header/stdexcept.html)       | 标准异常类型                       |
| [system_error](https://zh.cppreference.com/w/cpp/header/system_error.html) | 平台相关错误代码 `std::error_code` | [[↗]](./CPP%20STD/examples/Diagnostics/e_system_error.cpp) |

<!-- [debugging] 调试库 -->

> 内存管理库

| header                                                                             | description          | examples                                                            |
| :--------------------------------------------------------------------------------- | :------------------- | :------------------------------------------------------------------ |
| [memory](https://zh.cppreference.com/w/cpp/header/memory.html)                     | 内存管理工具         | [[↗]](./CPP%20STD/examples/MemoryManagement/e_memory.cpp)           |
| [memory_resource](https://zh.cppreference.com/w/cpp/header/memory_resource.html)   | 多态分配器与内存资源 | [[↗]](./CPP%20STD/examples/MemoryManagement/e_memory_resource.cpp)  |
| [scoped_allocator](https://zh.cppreference.com/w/cpp/header/scoped_allocator.html) | 嵌套分配器类         | [[↗]](./CPP%20STD/examples/MemoryManagement/e_scoped_allocator.cpp) |

> 元编程库

| header                                                                   | description    | examples                                                      |
| :----------------------------------------------------------------------- | :------------- | :------------------------------------------------------------ |
| [ratio](https://zh.cppreference.com/w/cpp/header/ratio.html)             | 编译时有理算术 | [[↗]](./CPP%20STD/examples/Metaprogramming/e_ratio.cpp)       |
| [type_traits](https://zh.cppreference.com/w/cpp/header/type_traits.html) | 编译时类型特征 | [[↗]](./CPP%20STD/examples/Metaprogramming/e_type_traits.cpp) |

> 通用工具库



> 字符串库

> 容器库

> 迭代器库

> 范围库
> 算法库
> 数值库
> 时间库

> 本地化库
> 输入/输出库
> 正则表达式库
> 并发支持库

>---
### Examples

- [examples](./CPP%20STD/README.md)


---