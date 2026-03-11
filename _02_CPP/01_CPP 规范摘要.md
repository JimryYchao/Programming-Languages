## CPP 规范摘要

---
### 1. 基本概念

#### 1.1. 类型系统

C++ 类型分为基础类型和复合类型。变量和成员函数声明可以包含 CV 限定。
- 基础类型：布尔、字符、整数、浮点数等类型；
- 复合类型：引用（左值 `&`、右值 `&&`）、指针、成员指针、数组、函数、枚举、类类型（`class`，`struct`，`union`）。	

***隐式生存期类型*** 包括标量类型（算术、指针、成员指针、枚举、`nullptr_t`）、隐式生存期类类型、数组等。隐式生存期类或聚合体需满足：平凡可构造、平凡可析构、无虚基类、非静态成员隐式生存期。

***平凡可复制类型*** 包括标量类型和平凡可复制类（数组）。平凡可复制类满足：平凡复制/移动构造、平凡复制/移动赋值运算符（delete 或 default）、平凡可析构。

***标准布局类型*** 包括标量类型和标准布局类类型（数组）。标准布局类满足：无虚函数或虚基类、非静态成员和基类均满足标准布局；非静态数据成员具有相同访问控制；继承层级中仅有一个类具有非静态数据成员。

***结构化类型*** 包括有算术类型、左值引用、指针或成员指针、枚举、`nullptr_t`、无闭包捕获 Lambda、全部公开非 `mutable` 非静态数据成员（包括基类）且成员类型也是结构化类型或数组的字面类类型。

> 布局与 POD

当类包含如虚基类、虚拟函数、或具有不同访问控制的成员时，编译器可以自由选择布局，因此存储对象可能不会使用连续的内存区域。例如某个类具有虚函数，则该类的所有实例可能会共享单个虚拟函数表。由于布局未定义，无法将这类对象传递到其他语言，因为它们可能是非连续的。POD（简单旧数据）类型同时为平凡和标准布局，它的内存布局是连续的，可以逐字节复制和二进制 I/O。标量类型是 POD 类型。

```cpp
#include <type_traits>
#include <iostream>
using namespace std;

struct B {
protected:
	virtual void Foo() {}
};
// 非平凡非标准布局
struct A : B {
	int a;
	int b;
	void Foo() override {}
};
// 平凡非标准布局
struct C {
	int a;
private:
	int b;   // 具有不同访问的非静态数据成员
};
// 标准布局非平凡
struct D {
	int a;
	int b;
	D() {} // 具有用户定义构造函数
};
struct POD {
	int a;
	int b;
};

int main() {
	cout << boolalpha;
	cout << "A is trivial is " << is_trivial<A>() << endl;					// false
	cout << "A is standard-layout is " << is_standard_layout<A>() << endl;  // false

	cout << "C is trivial is " << is_trivial<C>() << endl;					// true
	cout << "C is standard-layout is " << is_standard_layout<C>() << endl;  // false

	cout << "D is trivial is " << is_trivial<D>() << endl;					// false
	cout << "D is standard-layout is " << is_standard_layout<D>() << endl;	// true

	cout << "POD is trivial is " << is_trivial<POD>() << endl;				   // true
	cout << "POD is standard-layout is " << is_standard_layout<POD>() << endl; // true
}
```

>---
#### 1.2. 作用域和存储期

**作用域** 包含有全局范围、文件范围、命名空间范围、类范围和枚举范围、函数范围、块范围等。**链接** 包括外部链接（`extern`）、内部链接（`static`）、模块链接或无链接（块范围）。顶级命名空间中非 `extern` 的 `const` 名称具有内部链接。

**存储期** 包括有静态存储期（`static`,`extern`）、线程存储期（`thread_local`）、自动存储期（局部非静态变量）、动态存储期（`new` 或显式分配的对象，异常对象）。`thread_local` 变量在线程重用时可能被污染（线程池），因此可以在异步任务开始之前手动重置线程变量的状态。也可以利用线程指针对象（`new`）手动控制线程对象的创建和销毁，而不是在线程结束时自动调用线程对象的析构。

>---

#### 1.3. 对象资源管理

C++ 没有自动回收垃圾，程序负责资源的释放。当对象初始化时，它会获取它拥有的资源，并且该对象负责在其析构函数中释放资源。堆栈对象直接包含资源本身，这类原则称为 “资源获取即初始化” (RAII)。当堆栈对象超出范围时，会自动调用其析构函数。智能指针自动处理资源的释放。

```c++
#include <iostream>
struct buffer {
	size_t size;
	char* data;
	buffer(size_t size) :size{ size }, data(new char[size]) {}
	~buffer() { delete data; std::cout << "delete buffer" << std::endl; }
};
int main()
{
	buffer b(512);
	int	n = sprintf_s(b.data, b.size, "Hello World");
	std::cout << b.data;
}
// Hello World
// delete buffer
```

>---
#### 1.4. 程序启动与终止

```cpp
int main(void);
int main(int argc, char *argv[]);    // argv[0] 为程序名称
```
```cpp
#include <cstdlib>
int main() {
    if (cond1){
        exit(EXIT_SUCCESS);
    }else if (cond2) {
        abort();        // 跳过 atexit
    }else 
        return value;  // 相当于 exit(value) 
}
```

`main` 以返回值作为参数调用 `exit` 并清理自动变量。直接调用 `exit` 时不会销毁自动变量，仅销毁与当前线程关联的线程对象，然后调用注册的 `atexit` 函数（如果有）之后销毁静态对象。

>---
#### 1.5. 词法元素

> **标准关键字**

| description    | keywords                                                                                   |
| :------------- | :----------------------------------------------------------------------------------------- |
| 无类型         | `void`                                                                                     |
| 布尔           | `bool`,`false`,`true`                                                                      |
| 数值类型       | `signed`,`unsigned`,`int`,`short`,`long`,`long long`,`char`,`float`,`double`,`long double` |
| 字符类型       | `char`,`char8_t`,`char16_t`,`char32_t`,`wchar_t`                                           |
| 用户类型       | `class`,`struct`,`union`,`enum`                                                            |
| 命名空间与模块 | `namespace`,`module`,`using`,`export`,`import`                                             |
| 别名声明       | `typedef`,`using`                                                                          |
| 成员访问性     | `public`,`protected`,`private`                                                             |
| 继承与封装     | `virtual`,`final`,`override`                                                               |
| 类型说明符     | `const`,`volatile`,`register`,`extern`,`static`,`thread_local`,`alignas`                   |
| 参数修饰符     | `register`(弃用),`restrict`,                                                               |
| 说明符         | `friend`,`inline`,`mutable`,`operator`,`explicit`,`decltype`,`default`,`delete`            |
| 模板           | `template`,`tempname`,`class`,`concept`,`requires`                                         |
| 内存管理       | `new`,`delete`                                                                             |
| 编译时计算     | `consteval`,`constexpr`,`constinit`                                                        |
| 类型转换       | `const_cast`,`dynamic_cast`,`reinterpret_cast`,`static_cast`                               |
| 表达式         | `sizeof`,`typeid`,`auto`,`this`,`nullptr`                                                  |
| 协程           | `co_await`,`co_return`,`co_yield`                                                          |
| 迭代语句       | `do`,`while`,`for`                                                                         |
| 条件分支语句   | `if-else`,`switch-case-default`                                                            |
| 跳转语句       | `break`,`return`,`continue`,`goto`                                                         |
| 异常处理       | `try`,`catch`,`finally`,`noexcept`                                                         |
| 静态断言       | `static_assert`                                                                            |
| 汇编           | `asm`                                                                                      |


> **操作符**

```cpp
{ } [ ] ( )
<: :> <% %> ; : ...
? :: . .* -> ->* ~
! + - * / % ^ & |
= += -= *= /= %= ^= &= |=
== != < > <= >= <=> && ||
<< >> <<= >>= ++ -- ,
and or xor not bitand bitor compl
and_eq or_eq xor_eq not_eq
```

> **替用标志**

```cpp
and          &&    
or           ||  
not          !  
xor          ^     
bitand       &    
bitor        | 
compl        ~     
not_eq       !=  
xor_eq       ^= 
and_eq       &= 
or_eq        |=  
```


>---
#### 1.6. 模块

模块用于在翻译单元间共享声明和定义，可以在某些地方替代头文件的使用。一个完整的模块包含有一个主模块接口单元（*.ixx）和若干模块分区接口单元，以及模块实现单元。模块可以导出顶级类型声明、模块分区、头文件或其他模块的内容。不会导出任何宏定义。

| 语法                                                   | 描述                            |
| :----------------------------------------------------- | :------------------------------ |
| `export module <Name> [:<Partition>] [<Property>]`     | 声明模块接口单元                |
| `module <Name> [:<Partition>] [<Property>] `           | 声明模块实现单元                |
| `module;`                                              | 开始一个全局模块片段            |
| `module : private;`                                    | 开始一个私有模块片段            |
| `export <declaration>`,`export { <declaration-list> }` | 导出声明                        |
| `[export] import <Name>`                               | 导入（再 `export`）一个模块单元 |
| `[export] import <: Partition>`                        | 导入（再 `export`）一个模块分区 |
| `[export] import <Header>`                             | 导入（再 `export`）一个头文件   |

```cpp
export module moName;   // 主模块声明
namespace A{            // 隐式导出
    export void foo();  // 导出
    void bar();         // 未导出
}
export namespace B {    // 导出所有成员
    void foo();   
    void bar();   
}
```

> **模块组成的基本大纲**

主模块接口定义文件：

```c++
module; // 可选。全局模块起始位置
// #include 指令应放置于此处，但仅适用于本文件，不会与其他模块实现文件共享。
// 宏定义在本文件之外以及导入文件中均不可见。此处不使用导入语句。它们应放在模块的前置部分。
export module <Name>; // 主模块的起始部分
// import ...; 对属于该模块的所有文件均有效
// export import ...; 导出子模块,头文件
// export ...; 导出函数、类型、模板
module :private; // 可选。私有模块分区起始位置。
// 之后所有内容仅在本文件中可见
```

模块分区文件：

```c++
module; // 可选。全局模块起始位置
// #include 指令应放置于此处，但仅适用于本文件，不会与其他模块实现文件共享。
// 宏定义在本文件之外以及导入文件中均不可见。此处不使用导入语句。它们应放在模块的前置部分。
export module <Name>:<Partition>;  // 分区模块的起始位置
// import ...; 对属于该模块分区的所有文件均有效
// import :OtherPartition; 导入其他分区
// export import ...; 
// export ...; 
module :private; // 可选。私有模块分区起始位置。
// 之后所有内容仅在本文件中可见
```

模块实现单元：

```c++
// #include 或 import 语句; 可选。导入的内容对本文件可用
module <Name>[:<Partition>];    // 声明所属命名模块或模块分区
// 以下为模块实现
```

> **标准库模块**

C++23 标准库引入了两个命名模块：

- `import std` 导入 `std` 中定义的声明和名称，它还会导出 C 包装器标头的内容，例如 `<cstdio>` 和 `<cstdlib>`，提供类似 `std::printf()` 函数的内容。
- `import std.compat` 导入 `std` 中的所有内容，并添加 C 运行时全局命名空间，例如 `::printf`、`::fopen`、`::size_t`、`::strlen` 等。

```c++
import std.compat;
int main() {
    printf("Hello World! %s", "CPP!");
}
```

>---
#### 1.7. 命名空间

命名空间是一个声明性区域，`using` 导入命名空间或类型。匿名命名空间的成员仅对当前翻译单元可见（相当于内部链接）。

```c++
namespace CustomData {
	class ObjectManager {
	public:
		void DoSomething() {}
	};
	void Func(ObjectManager) {}
}
int main() {
	// 使用完全限定名称访问：
	CustomData::ObjectManager mgr;
	mgr.DoSomething();
	CustomData::Func(mgr);
	// using 导入 CustomData
	using namespace CustomData;        
	ObjectManager mgr;
	mgr.DoSomething();
	Func(mgr);
	// using 导入 ObjectManager
	using CustomData::ObjectManager;  
	ObjectManager mgr;
	mgr.DoSomething();
}
```

> *内联命名空间*

内联命名空间（`inline namespace`）的成员被视为父级空间的成员。

```c++
namespace Parent
{
    inline namespace new_ns
    {
        template <typename T>
        struct C
        {
            T member;
        };
        void Func();
    }
    template<> class C<int> {};  // 专用化
    void new_ns::Func() {		 // 提供实现
        cout << 1;
    }
}
int main() {
    Parent::Func();   // 相当于 Parent::new_ns::Func();
}
```

---
### 2. 声明

#### 2.1. 别名：typedef, using, namespace

`typedef` 在任意范围内创建任意既存类型的别名，具有外部链接。

```cpp
typedef struct {/* ... */ } S,*pS;  // 匿名类
typedef void Action();			// 函数
typedef int IntArray[];			// 数组
typedef int(&rFunc)();          // 函数引用
typedef void(*pFunc)(int);		// 函数指针
typedef int int_t, *intp_t, (&fp)(int, ulong), arr_t[10];  
```

`using` 创建定义类型或模板的别名。

```cpp
using S = struct { /* ... */ };
using Action = void();
using IntArray = int[];
using rFunc = int(&)();
using pFunc = int(*)();
// 别名模板
template <typename T,typename U>
using Vec = vector<T,U>;  
```

命名空间别名：

```c++
namespace a_very_long_namespace_name { class Foo {}; }
namespace AVLNN = a_very_long_namespace_name;
void Bar(AVLNN::Foo foo){ }
```

>---
#### 2.2. 对齐：alignas, alignof

`alignof(Type|Expr)` 获取对齐要求，`alignas(Type|Expr)` 设置（非位域、函数）对齐要求，对齐要求 = 0,1,2,4,8,16,... 不低于原类型自然对齐。

```cpp
struct alignas(8) S {
	char alignas(4) v1; 
	double v2;
};
int main(){
    bool alignas(8) b;
    cout << alignof(S);   // 8
}
```

预处理指令 `#pragma pack` 设置当前翻译单元的对齐方式：

```cpp
#pragma pack(2)  
struct T {
	int i;      // size 4
	short j;    // size 2
	double k;   // size 8
};
#pragma pack(0)  // 恢复默认
int main() {
	cout << alignof(T) << endl;      // 2
	cout << sizeof(T) << endl;       // 14
	cout << offsetof(T, i) << endl;  // 0
	cout << offsetof(T, j) << endl;  // 4
	cout << offsetof(T, k) << endl;  // 6
}
```

>---

#### 2.3. 类型推断：auto, decltype

`auto` 类型推断或类型占位符。

```cpp
auto Hi = "Hello World";    // const char*
auto f(int v) { return v * v; }   // int
auto sum(auto typename v1, auto typename v2) { return v + v2; }   // template<class>
```

`decltype` 生成指定表达式的类型，包含限定符。

```cpp
template<class T> requires !(std::is_reference_v<T>)
using V = decltype(std::forward<T>(T{}));

V<const double> v = 10;   // const double&&
```

`decltype(auto)` 可用于函数占位符的类型推断返回。

```c++
template<typename T, typename U>
auto Func(T&& t, U&& u) -> decltype(auto)  // 或 decltype(auto) Func;
    { return forward<T>(t) + forward<U>(u); };
```

>---
#### 2.4. 语言连接： extern "C", extern"C++"

语言规范提供在不同编程的程序单元之间进行连接的功能。*语言连接* 封装了与另一种语言程序单元链接所需的要求集：调用约定、名称修饰等。

```cpp
extern "C"|"C++" Decl | { Decl-List }

extern "C" void Func();
extern "C" {
#include "CHeader.h"   // 可能由 CHeader.c 提供实现
}
```

`"C++"` 为默认语言连接；`"C"` 允许 C++ 程序连接在 C 翻译单元定义的函数，并遵循 C 调用约定。当类成员、带有尾随 `requires` 友元函数、非静态成员函数声明 `"C"` 连接时，其类型的连接仍然是 `"C++"` 但参数类型（如果有）仍然是 `"C"`。

```cpp
extern "C" int open(const char *path_name, int flags); // C function declaration
int main() {
    int fd = open("test.txt", 0); // calls a C function from a C++ program
}

// This C++ function can be called from C code
extern "C" void handler(int) {
    std::cout << "Callback invoked\n"; // It can use C++ features
}
```

>---
#### 2.5. 编译时：constexpr, consteval, constinit

`constexpr` 编译时求值。`consteval` 指定立即函数，对函数的每次调用都必须产生编译时常量。`constinit` 断言变量具有静态初始化；变量是引用时，`constinit` 等价于 `constexpr`；`extern constinit` 引用外部已初始化的变量。

```cpp
consteval int InitEval() {
	if consteval 
		return 10086;   // 编译时立即函数，只会返回 10086
	else 
		return 10010;
}
constexpr int InitExpr() {
	if consteval 
		return 10086;   // 编译时返回 10086
	else 
		return 10010;   // 运行时返回 10010
}
int main() {
	constinit static int init_v = InitExpr();   // 编译时求值，10086
	constexpr int expr_v = InitExpr();			// 编译时求值，10086
	int x = InitExpr();			// 运行时求值，10010
	int v = InitEval();			// 编译时求值，10086
}
```

>---
#### 2.6. 属性声明

```cpp
[[attr-list]]
[[using attrNamespace : attr-list]]   // 指定的命名空间适用于所有这些属性, 例如 [[using CC: opt(1), debug]]

attr-list = attrId [(params)] 
          | attrNamespace : attrId [(params)]
```

标准属性定义：
- `[[noreturn]]` 指定函数永不返回；它始终引发异常或退出。
- `[[deprecated]]`,`[[deprecated("reason")]]` 指定声明可能出于某些原因不再使用；应用于类、*typedef*、变量、非静态数据成员、函数、命名空间、枚举、枚举器或模板专用化的声明。
- `[[fallthrough]]` 指示从前一个 `case` 贯穿是有意的，无需发出编译器诊断。
- `[[maybe_unused]]` 如果存在未使用的实体，则抑制编译器警告。
- `[[nodiscard]]`,`[[nodiscard("reason")]]` 鼓励编译器在返回值被丢弃时发出警告。
- `[[likely]]`,`[[unlikely]]` 指示编译器应针对语句的执行路径比任何其他执行路径或多或少可能的情况进行优化
- `[[assume(expression)]]` 指定 *expression* 在给定点始终求值为 `true`。

```c++
template<typename T1, typename T2>
[[nodiscard]] auto Plus(T1&& t1, T2&& t2) -> decltype(auto) {
    return forward<T1>(t1) + forward<T2>(t2);
}
int main() {
    std::cout << Plus(1, 2) << std::endl;
}
[[noreturn]] void Foo(int exitcode) {
    _Exit(exitcode);
}
```

---

### 3. 基础类型

| type       | keywords                                         |
| :--------- | :----------------------------------------------- |
| 无类型     | `void`                                           |
| 布尔       | `bool`,`false`,`true`                            |
| 空指针     | `std::nullptr_t`,`nullptr`                       |
| 字符       | `char`,`wchar_t`,`char8_t`,`char16_t`,`char32_t` |
| 标准浮点   | `float`,`double`,`long double`                   |
| 有符号整数 | `[signed]? char/short/int/long/long long`        |
| 无符号整数 | `unsigned char/short/int/long/long long`         |

> **整数字面值**

```cpp
0b1010'0101U;            // 二进制, u,U
013245670;               // 八进制
1234567890uLL;           // 十进制，ull
0x12345'67890'abcdefLL;  // 十六进制，ll
size_t z = -10010z;
ssize_t uz = 10086uz;
```

> **浮点数字面值**

```c++
float f         = 3.1415f;	   // F
double d        = 3.1415;
long double ld  = 3.1415L;     
float16_t f16   = 3.1415f16;   // F16
float32_t f32   = 3.1415f32;   // F32
float64_t d2    = 3.1415f64;   
float128_t f64  = 3.1415f128;
bfloat16_t bf16 = 3.1415bf16;  // BF16
// E-十进制表示法
double e1 = 1.234e-2; 
double e2 = 1.234E2;
// P-十六进制表示
double p1 = 0x1p-2;	  // == 0.25
double p2 = 0x2.p10;  // == 2048.0
```

> **字符字面值**

```cpp
// 字面值
char c       =   'A';     
wchar_t w    =  L'A';  
char8_t c8   = u8'A';  
char16_t c16 =  u'A';  
char32_t c32 =  U'A';       
int m        =   'abcd';       

// 转义序列
char u1 = 'A';          // 'A'
char u2 = '\101';       // octal; 1-3 位八进制，最大 \377
char u3 = '\x41';       // hexadecimal; 
char u4 = '\u0041';     // \u hhhh
char u5 = '\U00000041'; // \U HHHHHHHH

// 转义字符
'\'' '\"' '\?' '\\' '\a' '\b' '\f' '\n' '\r' '\t' '\v'
```

> 字符串字面值

常规字符串和 UTF-8 字符串称为窄字符串。

```c++
auto s0 =   "hello";   // const char* 
auto s1 = u8"hello";   // const char8_t* 
auto s2 =  L"hello";   // const wchar_t*
auto s3 =  u"hello";   // const char16_t*, UTF-16
auto s4 =  U"hello";   // const char32_t*, UTF-32
```

原始字符串以 `R"delimiter( char-sequence )delimiter"` 作为引导序列，例如 `R"==(Hello)=="` 为 `"Hello"`，`R"((a|b))"` 等价于 `"(a|b)"`。

```c++
auto R0 =   R"("Hello \ world")";   // const char*
auto R1 = u8R"("Hello \ world")";   // const char8_t*
auto R2 =  LR"("Hello \ world")";   // const wchar_t*
auto R3 =  uR"("Hello \ world")";   // const char16_t*, UTF-16
auto R4 =  UR"("Hello \ world")";   // const char32_t*, UTF-32
```

`s` 后缀表示映射到 `std::string`。

```c++
auto S   =   "hello"s; // std::string
auto S8  = u8"hello"s; // std::u8string
auto SW  =  L"hello"s; // std::wstring
auto S16 =  u"hello"s; // std::u16string
auto S32 =  U"hello"s; // std::u32string
auto R   =   R"("Hello \ world")"s;  // std::string
```

字符串拼接要求编码相同或其中一个无前缀：

```c++
char str[] = "12" "34";  // "1234"
auto hi = u8"hello" " " "world"s;  // "hello world"
```

---
### 4. 枚举

枚举包含一组命名的关联常数项，默认为 `int`。非作用域枚举项直接访问，可隐式转换为整数；作用域枚举项由枚举名称限定访问。由 `static_cast` 转换为整数。

```c++
enum Identifier [: integerType ] { ... }                 // 非作用域枚举
enum class | struct  Identifier [: integerType ] { ... }  // 作用域枚举

enum Week : unsigned char { 
    Monday /*0*/, Tuesday /*1*/, 
    Wednesday = 10, Thursday /*11*/, Friday, Saturday, Sunday };
enum class Suit { Diamonds, Hearts, Clubs, Spades };

int main() {
    Week day = Monday;
    Suit kind = Suit::Hearts;

    int tue = Tuesday;
    int clubs = static_cast<int>(Suit::Clubs);
}
```

>--- 
#### 4.1. 笼统枚举

没有枚举项的枚举称为笼统枚举。可以利用笼统枚举来声明一种新整数类型，例如 `std::byte`。

```cpp
enum E : <base-intType>;
enum class CE [: <base-intType>];
```



```c++
// std::byte
enum class byte : unsigned char;
byte operator|(const byte _Left, const byte _Right) noexcept {
    return static_cast<byte>(
        static_cast<unsigned char>(static_cast<unsigned int>(_Left) | static_cast<unsigned int>(_Right)));
}
byte operator&(const byte _Left, const byte _Right) noexcept {
    return static_cast<byte>(
        static_cast<unsigned char>(static_cast<unsigned int>(_Left) & static_cast<unsigned int>(_Right)));
}
byte operator^(const byte _Left, const byte _Right) noexcept {
    return static_cast<byte>(
        static_cast<unsigned char>(static_cast<unsigned int>(_Left) ^ static_cast<unsigned int>(_Right)));
}
byte operator~(const byte _Arg) noexcept {
    return static_cast<byte>(static_cast<unsigned char>(~static_cast<unsigned int>(_Arg)));
}
byte& operator|=(byte& _Left, const byte _Right) noexcept {
    return _Left = _Left | _Right;
}
byte& operator&=(byte& _Left, const byte _Right) noexcept {
    return _Left = _Left & _Right;
}
byte& operator^=(byte& _Left, const byte _Right) noexcept {
    return _Left = _Left ^ _Right;
}
int main() {
    byte b1{ 0 }, b2{ 1 };
    byte rt = b1 | b2;
    rt = b1 & b2;
    rt = b1 ^ b2;
    rt &= b1;
    rt |= b1;
    rt ^= b1;
}
```

>---
#### 4.2. using enum

`using enum E` 声明将 `E` 的枚举项作为成员引入到声明范围内。

```cpp
enum class E { a, b, c, d };
using enum E;       // ::a, ::b, ::c, ::d
struct S {
	using enum E;  // S::a, S::b, S::c, S::d
};
namespace N {
	using enum E;  // N::a, N::b, N::c, N::d
}
int main() {
	using enum E;
	auto ea = a;
	auto eb = ::b;  
	auto sc = S::c;
	auto nc = N::c;
}
```

---
### 5. 引用

`T&`,`T&&` 声明 `T` 类型的引用。编译器将已具名的右值引用视为左值。引用可以转换为兼容指针或其他左值引用。**可以为临时右值创建 `const T&` 左值引用，但是要避免脱离函数域传递临时右值构造的 `const T&` 而出现悬挂引用。**

```C++
// 直接绑定
int n = 10010;
int& lr = n;
int* pl = &lr; assert(&lr == pl);  // okay
int&& rrn = static_cast<int&&>(n); // rrn 直接绑定到 n 
int&& rrm = std::move(lr);         // rrm 直接绑定到 n

// 间接绑定
int&& rr = 10086;
const int& cr = 110;             // 指代 110.0 的临时量
const int& lm = std::move(911);  // 指代 991.0 的临时量
double&& rrd = rr;               // 指代值 (double)rr 的临时量
rrd = 1000000000000;
cout << rr;   // 10086
```

`&&` 保留对右值表达式的引用；*Rvalue* 引用支持 “移动语义”。利用移动语义，可以将资源（如动态分配的内存）从一个对象转移到另一个对象，例如利用 `std::move` 转换为右值引用。

```cpp
string a = "hello";
string&& b = std::move(a);
cout << b;   // "hello"
cout << a;   // "hello"
```

实现移动语义，通常需要提供 “移动构造函数”（`class(class&&)`），或移动赋值运算符（`operator=(class&&)`）。利用移动语义，可以直接移动对象而不必执行内存分配和复制操作。例如 `std::string`：

```cpp
string s = string("h") + "e" + "ll" + "o";
```

>---
#### 5.1. 临时量的生存期

当引用绑定到临时对象或它的子对象，临时对象的生存期被延长以匹配引用的生存期。例外情景：
- `return` 语句中绑定到函数返回值的临时量不会被延续。`return` 始终返回悬垂引用 (C++26 前)。
    ```cpp
    const std::string& getString() {
        return "hello";  // 临时字符串在 return 结束时销毁
    }
    int main() {
        const std::string& s = getString(); // s 是悬垂引用
        std::cout << s;                     // 未定义行为
    }
    ```
+ 在函数调用中绑定到函数形参的临时量，生存期仅当前函数域：

    ```cpp
    const string& forwardRef(const string& param) {
        return param;  // 返回绑定到临时对象的引用
    }
    int main() {
        const string& ref = forwardRef("hello");  // 临时对象，超出范围，悬垂引用
        std::cout << ref;                         // 未定义行为
    }
    ```

- 绑定到 `new` 引用的临时量，仅存在到表达式结尾：

    ```cpp
    struct Test {
        const string& ref;  // 引用成员
        Test(const string& r) : ref(r) {}  // 绑定临时对象
    };
    int main() {
        // 直接绑定临时对象（悬垂引用）
        Test* t1 = new Test("Hello");  // 临时对象，超出范围，悬垂引用
        std::cout << t1->ref;          // 未定义行为：访问已释放的内存

        // 通过函数返回临时对象（同样悬垂）
        auto create_temp = []() { return "10086"; };
        Test* t2 = new Test(create_temp());  // 临时对象在表达式 create_temp 后销毁
        std::cout << t2->ref;                // 未定义行为
    }
    ```


+ 绑定到用直接初始化语法（`T(value)`）引用的临时量，存在到初始化器末尾。

    ```cpp
    struct Test {
    	int v;
    	const string& ref;
    };
    int main() {
    	Test t1(1, "hello");		// 直接初始化，临时对象存在到初始化器末尾
    	Test t2{ 2, "Hello" };      // 列表初始化

    	cout << t1.ref;   // 悬挂引用
    	cout << t2.ref;   // Hello
    }
    ```

---
### 6. 指针

指针可以指向类型化对象、类型成员、函数或其他指针，无法指向引用或位域。原始指针生存期不受封装对象控制。

```cpp
struct S {
	static void Foo() {}
	void foo(const string& mess) { cout << mess; };
};
int main() {
	S s{};
	struct S* ps = &s;		     // 对象指针
	void (*pF)() = S::Foo;       // 函数指针，静态函数
	void (S::*pf)(const string&) = &S::foo;   // 成员指针
	(s.*pf)("Hello");            // 从对象调用成员指针
	(ps->*pf)(" World");         // 从对象指针调用成员函数
}
```

指针的 CV 限定语法（`*` 后随 CV 限定）：

```c
int num = 0;
// 指针 CV 限定
int* const cp = &num;    
int* volatile vp = &num;
int* const volatile cvp = &num;
// CV 限定类型
const int* pc = &num;           // 基础类型为 const int
volatile int* pv = &num;		// 基础类型为 volatile int
const volatile int* pcv = &num; // 基础类型为 const volatile int
// 组合
volatile const int* const volatile cvpcv = &num;   // 指向 CV 对象的 CV 指针
```

`new` 运算返回一个堆分配对象的指针，需手动 `delete` 释放，未释放的内存将导致内存泄漏。

```c++
MyClass* mc = new MyClass; // allocate object on the heap
delete mc; 	  // delete object

int *arr = new int[size];
delete[] arr;
```

> **指针运算**

```c++
// address-of
int * pi = &i;
const type cp = &const_type;
volatile type vp = &vola_type;
// 算数
p = p + integer;  // -
p += integer;     // -=
p++;              // --
// 指针差
ptrdiff_t diff = p1 - p2;
```

>---
#### 6.1. 函数指针

函数指针可以指向非成员函数或静态成员函数。

```cpp
void foo(int v) { cout << v; }
template<typename T>
void foo(T v) { cout << v; }

int main() {
	using Foo = void(int);
	Foo* pf = foo;
	Foo& rf = foo;    // 函数引用

	pf(1);			 // 指针调用
	rf(2);			 // 引用调用
	(*pf)(3);		 // 解引用调用

	void (&Tfoo)(double) = foo;   // foo<double>
	void (&Tfoo)(int) = foo;      // foo(int)  重载决策优先级更高
}
```

>---
#### 6.2. 成员指针

成员指针指向类类型的非静态成员。

```cpp
struct S {
	int M;
	void F(int v) { cout << v; }
};

int main() {
	int S::* pM = &S::M;      // 数据成员指针
	void (S::* pF)(int) = &S::F; // 函数成员指针
	// 通过对象访问
	S s{10086};
	(s.*pF)(s.*pM);
	// 通过指针对象访问
	S* ps = &s;
	(ps->*pF)(ps->*pM);
}
```

指向一个可访问且无歧义（非虚）基类数据成员的指针，可以隐式转换为指向派生同一数据成员的指针。或由 `static_cast` 和显式转换，即使基类没有该成员。

```cpp

struct B {   // 非虚基类
	int M;
};
struct D : B {
	double N;
	void F() { cout << "F\n"; }
	D(int m, double n) : N{ n } { M = m; }
};
int main() {
	D d{ 10086, 3.1415 };
	B b;  

	int B::* bM = &B::M;
	int D::* dM = bM;        // 隐式转换, int B::* -> int D::*
	cout << d.*dM << endl;   // 10086

	double D::* dN = &D::N;
	double B::* bN = static_cast<double B::*>(dN);  // 显式转换, double D::* -> double B::*
	cout << d.*bN << endl;   // 3.1415
	cout << b.*bN << endl;   // 未定义行为

	void (D::* dF)() = &D::F;
	void (B::* bF)() = static_cast<void (B::*)()>(dF);
	(d.*bF)();    // F
	(b.*bF)();    // 未定义行为
}
```

指向基类成员函数的指针可以隐式转换为指向派生同一成员函数的指针。

```cpp
struct B { virtual void F() { cout << "B\n"; } };
struct D : B { void F() override final { cout << "D\n"; } };
int main() {
	D d;
	void (B:: * bF)() = &B::F;
	void (D:: * dF)() = bF;      // 隐式转换
	(d.*bF)();   // D, 重写
	(d.*dF)();   // D
}
```

>---
#### 6.3. 智能指针

智能指针（`<memory>`）包装原始指针，自动管理资源释放。

- `unique_ptr` 独占指针，可以移动（`std::move`）到新所有者，但不会复制或共享。支持 *rvalue* 引用。
- `shared_ptr` 共享指针，采用引用计数。全部所有者超出了范围或放弃所有权，才会释放原始指针。
- `weak_ptr` 解决 `shared_ptr` 循环引用导致的内存泄漏问题，不增加对象的引用计数。

```c++
Point *p = new Point{3, 4};
unique_ptr<Point> p1(p);
unique_ptr<Point> p2 = make_unique<Point>(*p);
unique_ptr<Point> p3 = make_unique<Point>(3, 4);
```

> **unique_ptr**

`unique_ptr` 不共享它的指针，只能移动 `unique_ptr`。创建 `unique_ptr` 实例并在函数之间传递实例：

```c++
struct Song {
    string artist;
    string title;
    Song(string artist, string title) : artist{ artist }, title{ title } {};
};
unique_ptr<Song> SongFactory(const std::string& artist, const std::string& title) {
    // Implicit move operation into the variable that stores the result.
    return make_unique<Song>(artist, title);
}
void MakeSongs() {
    auto song = make_unique<Song>("Mr. Children", "Namonaki Uta");
    unique_ptr<Song> song2 = std::move(song);
    auto song3 = SongFactory("Michael Jackson", "Beat It");
}
```

使用 `make_unique` 将 `unique_ptr` 创建到数组，但无法使用 `make_unique` 初始化数组元素：

```c++
// Create a unique_ptr to an array of 5 integers.
auto p = make_unique<int[]>(5);
// Initialize the array.
for (int i = 0; i < 5; ++i) 
    p[i] = i;
```

> **shared_ptr**

`shared_ptr` 所有实例均指向同一个对象，并共享对一个 “控制块” 的访问权限。当引用计数达到零时，控制块将删除内存资源和自身。

```c++
struct MediaAsset {
    virtual ~MediaAsset() = default; // make it polymorphic
};
struct Song : public MediaAsset {
    std::string artist;
    std::string title;
    Song(const std::string& artist_, const std::string& title_) :
        artist{ artist_ }, title{ title_ } {}
};
struct Photo : public MediaAsset {
    std::string date;
    std::string location;
    std::string subject;
    Photo(const std::string& date_,
          const std::string& location_,
          const std::string& subject_) 
        : date{ date_ }, location{ location_ }, subject{ subject_ } {}
};

int main() {
    std::shared_ptr<Song> sp1 = make_shared<Song>("The Beatles", "Im Happy Just to Dance With You");
    std::shared_ptr<Song> sp2(new Song("Lady Gaga", "Just Dance"));
    std::shared_ptr<Song> sp3 = sp2;
}
```

使用 `dynamic_pointer_cast`、`static_pointer_cast` 和 `const_pointer_cast` 来转换 `shared_ptr`：

```c++
vector<shared_ptr<MediaAsset>> assets {
    make_shared<Song>("Himesh Reshammiya", "Tera Surroor"),
    make_shared<Song>("Penaz Masani", "Tu Dil De De"),
    make_shared<Photo>("2011-04-06", "Redmond, WA", "Soccer field at Microsoft.")
};
vector<shared_ptr<MediaAsset>> photos;
copy_if(assets.begin(), assets.end(), back_inserter(photos), [] (shared_ptr<MediaAsset> p) -> bool
{
    // Use dynamic_pointer_cast to test whether element is a shared_ptr<Photo>.
    shared_ptr<Photo> temp = dynamic_pointer_cast<Photo>(p);
    return temp.get() != nullptr;
});

for (const auto&  p : photos)
    // the photos vector contains only shared_ptr<Photo> objects, so use static_cast.
    cout << "Photo location: " << (static_pointer_cast<Photo>(p))->location << endl;
```

按值传递 `shared_ptr` 将调用复制构造并增加引用计数。按引用传递 `shared_ptr`，引用计数不会增加。传递原始指针或对原始对象的引用，被调用方能够使用对象，但不会共享所有权或延长生存期。

```c++
void use_shared_ptr_by_value(shared_ptr<int> sp);
void use_shared_ptr_by_reference(shared_ptr<int>& sp);
void use_shared_ptr_by_const_reference(const shared_ptr<int>& sp);
void use_raw_pointer(int* p);
void use_reference(int& r);

void test() {
    auto sp = make_shared<int>(5);
    use_shared_ptr_by_value(sp);       // 值传递复制，计数增加
    use_shared_ptr_by_reference(sp);   // 引用传递，计数不会增加
    use_shared_ptr_by_const_reference(sp);   // 常量引用传递，计数不会增加
    use_raw_pointer(sp.get());			// 传递基础指针，计数不会增加
    use_reference(*sp);				    // 传递对基础对象的引用，计数不会增加
    use_shared_ptr_by_value(move(sp));  // 移动语义，计数不会增加
}
```

> **weak_ptr**

`weak_ptr` 可以存储 `shared_ptr` 的基础对象，而不会增加引用计数。它无法阻止引用计数变为零，若 `shared_ptr` 被销毁，则 `if(weak_ptr.expired())` 返回 `true`。可以使用 `weak_ptr` 尝试获取 `shared_ptr` 的新副本。

```c++
int main()
{
    shared_ptr<int> sp = make_shared<int>(100);
    weak_ptr<int> wp = sp;
    sp.reset();
    if (wp.expired())
        cout << "对象已销毁" << endl;

    // 从 wp 到 sp 需要 lock()
    sp = make_shared<int>(200);
    wp = sp;
    shared_ptr<int> temp_sp = wp.lock();
	sp.reset();	
    if (temp_sp)                                   
        cout << "获取到的值：" << *temp_sp << endl; // 200
}
```

---
### 7. 数组

数组具有相同类型、连续存储元素的对象序列，无法创建引用或函数的数组。CV 限定仅作用于数组元素。数组名称被视为指向首元的指针。数组作为参数传递时，衰减为指针。

堆栈数组速度比堆数组更快。堆数组通常由 `new T[size]` 分配，并由 `delete[] T` 负责释放；堆栈数组在函数返回时自动清理。

```c++
int stack[] = { 1,2,3,4 };       // 堆栈数组 int[4]
int* heap = new int[size] {0};   // 堆数组 int[size]
//  ... use heap
delete[] heap;
// 多维数组
int arr[7];
int arr3x4[3][4];
int arr5x6x7[5][6][7];  
// 数组初始化
int a[6] = {};                  // 零初始化
int a1[5] = { 1,2,3 };          // 后续元素零初始化
int a2[] {6,7,8,9};             // 推断 int[4], 等号省略
int ma[][3] {{},{},{}};         // int[3][3]
int ma2[2][3] { 1,2,3,4,5,6 };    // 顺序初始化 
// 动态数组
int* p = new int[10];   // 内容未知
int* pa = new int[] {1, 2, 3, 4, 5};  // 推断, int[5]
int* pa2 = new int[5] {1, 2, 3};
int (*pm)[2] = new int[3][2];
int (*pm)[3][4] = new int[2][3][4];
```

>---
#### 7.1. 字符数组

字符数组可以由字符串值初始化，数组长度必须足够容纳字符串长度（包括 `\0`）。

```c
char hi[] = "Hi"; // char[3] {'\H','i','\0'}
char hi2[4] = "Hi"; // char[4] {'\H','i','\0','\0'}
wchar_t hi3[6] = L"Hello";  // len >= 6;

char err[5] = "Hello World";   // ERROR
```

---
### 8. 表达式

| Category     | Operators                                                                                            |
| :----------- | :--------------------------------------------------------------------------------------------------- |
| 基本表达式   | `x::y`,`x.y`,`x->y`,`f(x)`,`x[y]`,`x++`,`x--`,`typeid(x)`,`T{v}`,`T(x)`                              |
| cast 表达式  | `const_cast`,`dynamic_cast`,`reinterpret_cast`,`static_cast`                                         |
| 一元         | `sizeof(x \| T)`,`++x`,`--x`,`~x`,`!x`,`-x`,`+x`,`&x`,`*x`,`new`,`delete`,`(T)v`                     |
| 指向成员指针 | `x.*p`,`x->*p`                                                                                       |
| 乘法         | `x * y`,`x / y`,`x % y`                                                                              |
| 加法         | `x + y`,`x - y`                                                                                      |
| 移位         | `x << y`,`x >> y`                                                                                    |
| 关系         | `x > y`,`x < y`,`x >= y`,`x <= y`                                                                    |
| 相等         | `x == y`,`x != y`                                                                                    |
| 位运算       | `x & y`,`x ^ y`,`x \| y`                                                                             |
| 条件逻辑     | `x && y`,`x \|\| y`                                                                                  |
| 三目运算     | `cond ? <true> expr1 : <false> expr2`                                                                |
| 赋值         | `x = y`,`x *= y`,`x /= y`,`x %= y`,`x += y`,`x -= y`,`x <<= y`,`x >>= y`,`x &= y`,`x \|= y`,`x ^= y` |
| 异常         | `throw`                                                                                              |
| 逗号运算     | `(expr1, expr2 [, expr3, ..., exprN])` 返回 `exprN` 表达式结果                                       |
          
> **表达式值类别**

表达式具有类型和值类别两种特性，每个表达式只属于三种基本值（纯右值、亡值、左值）类别中的一种。值类别包括有：
- 泛左值 *glvalue*：求值时可确定某个对象或函数标识的表达式。
- 纯右值 *prvalue*：运算符操作数的值、初始化某个对象的结果对象。
- 亡值 *xvalue*：资源能够被重新使用的对象或位域的泛左值。
- 右值 *rvalue*：纯右值或亡值。
- 左值 *lvalue*：非亡值的泛左值。

```
         expression
         /        \
       glvalue   rvalue
       /     \  /     \
   lvalue   xvalue   pvalue
```

```cpp
#include <type_traits>
#include <utility>

template <class T> struct is_prvalue : std::true_type {};
template <class T> struct is_prvalue<T&> : std::false_type {};
template <class T> struct is_prvalue<T&&> : std::false_type {};

template <class T> struct is_lvalue : std::false_type {};
template <class T> struct is_lvalue<T&> : std::true_type {};
template <class T> struct is_lvalue<T&&> : std::false_type {};

template <class T> struct is_xvalue : std::false_type {};
template <class T> struct is_xvalue<T&> : std::false_type {};
template <class T> struct is_xvalue<T&&> : std::true_type {};

int main()
{
    int a{ 42 };
    int& b{ a };
    int&& r{ std::move(a) };

    // 表达式 `42` 是纯右值
    static_assert(is_prvalue<decltype((42))>::value);
    // 表达式 `a`,`b`,`r` 是左值
    static_assert(is_lvalue<decltype((a))>::value);
    static_assert(is_lvalue<decltype((b))>::value);
    static_assert(is_lvalue<decltype((r))>::value);
    // 表达式 `std::move(a)` 是亡值
    static_assert(is_xvalue<decltype((std::move(a)))>::value);
    // 变量 `r` 的类型是右值引用
    static_assert(std::is_rvalue_reference<decltype(r)>::value);
    // 变量 `b` 的类型是左值引用
    static_assert(std::is_lvalue_reference<decltype(b)>::value);
}
```

>---
#### 8.1. 类型转换：const_cast, static_cast, dynamic_cast, reinterpret_cast

> **显式类型转换**

```c++
// cast 形式
[CV] (<Type>) UnaryExpr;   
// 函数表示法
[CV] <Type> ( <init-list> );   
[CV] typename <Type-id> { <init-list> };

struct Point { int x, y; };
int main()
{
	auto p = Point(1, 2);
	auto p1 = Point(p);
	auto p2 = Point{ p };
	auto p3 = typename ::Point{ p };
}
```

> 用户定义转换函数

```cpp
operator <Type-id>                      // 隐式转换
explicit operator <Type-id>             // 显式转换
explicit (cond) operator <Type-id>      // 条件转换。true 为显式，false 为隐式

operator auto | decltype(auto)          // type-id 可以是 auto 占位
```

*Type-id* 不能是函数、数组、或出现运算符 `[]`（可以是数组指针别名）。


```c++
struct KM {
	long double distance;
public:
	KM() { distance = 0; };
	KM(long double dis) { distance = dis; };
	operator long double() { return distance; }   // 隐式类型转换
	string toString() { return to_string(distance) + "km"; }
	explicit operator string() { return this->toString(); }  // 显式类型转换

	KM operator=(int dis) { return KM{ (long double)dis }; }
	KM operator=(long double dis) { return KM{ dis }; }
	KM operator=(unsigned long long dis) { return KM{ (long double)dis }; }
};
KM operator ""km(long double dis) {   // 用户定义文本后缀
	return KM{ dis };
};
KM operator ""km(unsigned long long dis) {
	return KM{ (long double)dis };
};
KM operator ""km(long long dis) {
	return KM{ (long double)dis };
};
int main() {
	KM d1 = 10km;
	KM d2 = 3.14km;
	KM d3 = 10;
	KM d4 = 3.14f;
	long double d = d1;  // or (long double)d1;
	cout << d1.toString();  // 10.000000km
}
```

可以定义转换为自身类或引用、它的基类或引用、`void`，但它们无法作为转换序列的一部分执行。可以通过虚派发或成员函数语法调用：

```cpp
struct D;
struct B { virtual operator D() = 0; };
struct D : B {
	operator void() { cout << "void\n"; }
	operator D () override { return (cout << "D\n", D()); }
	operator B& () { return (cout << "B&\n", *this); }
};
int main() {
	D d;
	B& b1 = d;
	D d2 = b1;   // D
	d.operator void();  // 函数调用
	B& b2 = d.operator B & ();  // B&
}
```

> **const_cast**

`const_cast` 添加或移除（非函数）指针或引用的 CV 限定（编译时行为）。通过移除 `const` 限定的指针或引用修改 `const` 或 `constexpr` 原始值的行为未定义。

```c++
constexpr int cnum = 10086;  // 或 const int
const int& cr_cnum = cnum;
int& r_cnum = const_cast<int&>(cr_cnum);
r_cnum = 911;   // UB
assert(r_cnum == cnum); // failed; c_number = 10086
```

如果 `T` 是类类型或数组类型，它的纯右值可以通过 `const_cast<T&&>` 显式转换成 `T` 的亡值。

```cpp
typedef int(&&RA)[];
RA getPArr() {
	return { 1,2,3,4,5 };
}
RA&& ra = const_cast<int(&&)[]>(getPArr());
```

> **static_cast**

`static_cast` 用于基础类型之间转换、指针或引用在继承链中的向上 / 向下转换、转换构造函数或转换运算符、枚举与整数之间转换等，编译时检查。`static_cast<void>(expr)` 相当于 `(void)expr` 弃置。

```cpp
import std;
using namespace std;

struct B {
	int m = 42;
	const char* f() const { return "B\n"; }
};
struct D : B {
	const char* f() const { return "D\n"; }
};
enum class E { ONE = 1, TWO, THREE };
enum EU { ONE = 1, TWO, THREE };

int main()
{
// 静态向下转换
	D d; B& br = d; // 通过隐式转换向上转换
	D& another_d = static_cast<D&>(br); // 向下转换, B&(d) -> D&
// 左值到亡值
	vector<int> v0{ 1, 2, 3 };
	vector<int> v2 = static_cast<vector<int>&&>(v0);   // vector<int>(v0) -> vector<int>&& = 0
// 弃值表达式
	static_cast<void>(v2.size());  // void(v2.size())
// 初始化转换
	int n = static_cast<int>(3.14);   // float(3.14) -> int(3)
	vector<int> v = static_cast<vector<int>>(10);   // int(10) -> vector<int>(10)
// 隐式转换的逆转换
	void* nv = &n;   // int*
	int* ni = static_cast<int*>(nv);   // void*(n) -> int*
// 范围枚举到 int 或 float
	E e = E::TWO;
	int two = static_cast<int>(e);   // E::TWO -> int
// int 到枚举，枚举到另一枚举
	E e2 = static_cast<E>(two);	   // int(2) -> E::TWO 
	EU eu = static_cast<EU>(e2);   // E::TWO -> EU::TWO
// 指向成员指针的向上转换
	int D::* dpm = &D::m;
	int bm = br.*static_cast<int B::*>(dpm);  // D::*m -> B::*m
// void* 到任意对象指针
	void* vp = &e;
	vector<int>* p = static_cast<vector<int>*>(vp);  // void* ->  vector<int>*
}
```

向下转型存在​多态类型风险，可能会引发未定义行为（例如无效内存访问、栈被破坏等），优先使用 `dynamic_cast`。

```cpp
class B {
public:
	virtual void bf() {}
};
class D : public B {
public: 
	int v;
	virtual void vf() {};
	void f() {};
};
int main() {
	B b;
	D* pd = static_cast<D*>(&b);
	pd->f();		   // maybe ok
	// pd->vf();	   // err in run-time
	// pd->v = 100;    // err in run-time

	B* pb = new D;  
	if (D* pd2 = dynamic_cast<D*>(pb); pd2)   // B* -> D*
		pd2->vf();  
}
```

> **dynamic_cast**

`dynamic_cast` 用于多态类型的转换，运行时类型检查，如从多态基类的指针或引用到派生类指针或引用的向下转换、类层次之间的交叉转换（侧向转换）等。转换失败返回 `nullptr`；转换引用失败引发 `bad_cast`。

```cpp
// 多态类类型至少包含一个虚函数
struct V { virtual void f() { } };
struct A : virtual V { };
struct B : virtual V {
	void f() override { cout << "B.f\n"; };
};
struct D : A, B { };
int main() {
	D d; // 最终派生对象, 具有 A,B,V
	A& a = d; // 隐式向上转换，可以用 dynamic_cast，非必须
	D& new_d = dynamic_cast<D&>(a); // 向下转换, A&(d) -> D&
	B& new_b = dynamic_cast<B&>(a); // 侧向转换, A&(d) -> D& -> B&
}
```

如果有多重继承，可能会导致不明确。使用虚基类时，可能会导致更多不明确的情况。

```c++
class A { virtual void f(); };
class B : public A { virtual void f(); };
class C : public A { virtual void f(); };
class D : public B, public C { virtual void f(); };
int main() {
	D* pd = new D;
	A* pa = dynamic_cast<A*>(pd);   // 歧义转换失败 at runtime
	// `D*` 指针可以安全地转换为 `B*` 或 `C*`。
	// 但如果 `D*` 强制转换为 `A*`，将导致不明确的强制转换错误。
	// 若要解决此问题，可以执行两步明确的强制转换。
	B* pb = dynamic_cast<B*>(pd);   // D* -> B*
	A* pa2 = dynamic_cast<A*>(pb);  // B* -> A*
}
```

> **reinterpret_cast** 

`reinterpret_cast` 执行位模式重新解释，直接复制源类型的二进制位到目标类型。允许在任意指针或引用之间、整数与指针之间转换，无视类型关系。

```c++
struct Point2D { int x, y; };
struct Point3D { int x, y, z; };

int main() {
    Point2D* p2 = new Point2D{ 2, 2 };
    Point3D* p3 = new Point3D{ 3, 3, 3 };

    Point2D* p2D = reinterpret_cast<Point2D*>(p3);  // v3 -> v2; not safy
    p2D->x = 1; p2D->y = 99;  // v3 {1,99,3}

    Point3D* p3D = reinterpret_cast<Point3D*>(p2);  // v2 -> v3; not safy
    p3D->x = 9; p3D->y = 9;
    p3D->z = 10;  // maybe stackOverflow ; danger!

    // 未定义行为
    Point3D v3D = static_cast<Point3D>(*p3D);
    printf("v = {%d,%d,%d}", v3D.x, v3D.y, v3D.z); // not safe, {9,9,10}
    Point2D v2D = static_cast<Point2D>(*p2D);      // maybe safe, {1,99}
}
```


>---
#### 8.2. 内存分配

`new` 创建堆存储并初始化具有动态存储期的对象，这些对象的内存可能在未来的某个时刻由 `delete` 释放。

```cpp
[::] new ( placement-Params )? Typeid [initializer]
[::] delete ( placement-Params )? expr
[::] delete ( placement-Params )? [] expr

S* ps = new S;
int * p = new int[]{1,2,3};

delete ps;
delete[] p;
```

> 重载 `new / delete`

可创建自定义 `new / delete` 运算符，用于分配内存。

```cpp
struct allocteBlock {
private:
	static constexpr size_t blockSz = 8;
public:
	// User-defined operator new.
	void* operator new(size_t sz, size_t count) {
		if (count <= 0 || sz <= 0) return nullptr;
		auto p = calloc(count, blockSz);
		if (p)
			cout << "Memory block [" << p << "] allocated for "
			<< count << " bytes\n";
		return p;
	}
	// User-defined operator delete.
	void operator delete(void* pvMem) {
		if (!pvMem) return;
		cout << "Memory block [" << pvMem
			<< "] deallocated\n";
		free(pvMem);
	}
};
int main() {
	for (int i = 0; i < 10; ++i) {
		auto pMem = new(i) allocteBlock;
		delete pMem;
	}
}
```

> 可替代函数

可声明 `new`、`delete` 表达式的全局替代函数：

```cpp
operator new
operator new[]
operator delete
operator delete[]


void* operator new(size_t sz, int count) {
	return calloc(count, sz);
}
void operator delete(void* p, size_t) {
	free(p);
}

auto buf = new(10) double;
delete(buf);
```

>---
#### 8.3. typeid

```cpp
typeid ( <Type> | <expr> )
```

`typeid` 查询完整类型或表达式的基础类型信息，返回 `const std::type_info&`。

```cpp
static_assert(typeid(int) == typeid(int&));
static_assert(typeid(int) == typeid(int&&));
static_assert(typeid(int&&) == typeid(int&));
static_assert(typeid(int) == typeid(const int));
static_assert(typeid(volatile int&&) == typeid(const int));
		
cout << typeid(volatile int&&).name();   // int
```

多数情况下，`typeid` 是不求值表达式，对于多态类类型对象或表达式，`typeid` 对其求值并指代为实际动态类型的 `type_info`。

```cpp
struct Base { virtual void f() {} };
class DerivedA : public Base { virtual void f() override {} };
class DerivedB : public Base { virtual void f() override {} };
class Derived : public DerivedA, public DerivedB {};

int main() {
	Derived* d = new Derived;
	Base* da = static_cast<DerivedA*>(d);
	Base* db = static_cast<DerivedB*>(d);
	assert(typeid(*da) == typeid(*db));  // true;

	Base* ba = new DerivedA;
	Base* bb = new DerivedB;
	assert(typeid(*ba) == typeid(*bb));  // false;
}
```

---
### 9. 语句
#### 9.1. 空语句

```c++
while(cond)
    ;
```

>---
#### 9.2. 条件控制语句：if, switch

> **if-else**

```c++
if ( [init-expr;] condition ) 
    statements
[ else [if-clause] ]
    statements
```
```c++
if (auto x = getX(); x < 11)
    cout << "x < 11 is true!\n";  
else
    cout << "x < 11 is false!\n"; 
```

> **constexpr-if**

```cpp
if constexpr (<bool-constexpr>) 
    statements
[else [if ...]] 
```
```c++
template<typename T>
auto get_value(T t)
{
	if constexpr (std::is_pointer_v<T>)
		return *t; // 对 T = int* 推导返回类型为 int
	else
		return t;  // 对 T = int 推导返回类型为 int
}
```

> **consteval-if**

```cpp
if consteval {
    compile-time-statements  // 编译时计算分支
} [else {
    run-time-statements  // 运行时计算分支
} ] 
// or 
if !consteval { rt; } [ else { ct; } ] 
```
```cpp
constexpr int Square(int x) {
	if consteval {
		return x * x;
	}
	else {
		return (cout << "run-time\n", x * x);
	}
}
int main() {
	constexpr int cx = Square(10);
	int x = Square(10);   // run-time
}
```

> **switch-case**

在 `case` 中声明的变量属于 `switch` 语句范围，它们共享名称。避免空引用访问。

```c++
swtich ( [init-expr;] int-cond ){
    case cond1:
        statements
        [ break | return ];  // 允许贯穿，编译器警告
    case cond2: 
        statements
        [[fallthrough]];     // 抑制贯穿警告 
    [default: ]
}
```

>---
#### 9.3. 循环语句：while, do, for, for-range 

> while

```c++
while ( expression ) { 
   statements
};
```

> **do-while**

```c++
do {
    statements
} while ( expression );
```

> **for**

```c++
for ([init-expr]; [cond-expr]; [loop-expr] ) {
    statements
}  // for(;;) == while(true)
```

> **for-range**

```c++
for ( [init-expr;] elem-declaration : IteratorExpression ) {
    statements
}
```
```c++
int x[]{ 1,2,3,4,5,6,7,8,9,0 };
for (auto e : x)
    cout << e << ",";
for (typedef T = int; T e : x)
	// something
```

`for` 可以 *range* 的类型包括数组、拥有 `.begin()` 和 `.end()` 的容器或用户定义类型。

```c++
template <typename T>
class Iterator {  // 迭代器
private:
    T* iter;
public:
    Iterator(T* para, size_t n) { iter = para + n; }
    T& operator *() { return *iter; }
    bool operator != (const Iterator& that) { return this->iter != that.iter; }
    Iterator& operator++ () { ++iter; return *this; }
};
template <typename T>
class docker {  // 实现 begin 和 end
private:
    T* p;
    size_t size;
public :
    docker(size_t n, T arr[]) :size(n) { p = arr; }
    Iterator<T> begin() { return Iterator<T>(this->p, 0); }
    Iterator<T> end() { return Iterator<T>(this->p, size); }
};

int main(int argc, char* argv[])
{
    int x[10]{ 1,2,3,4,5,6,7,8,9,0 };
    docker<int> d(10, x);
    for (auto& i : d)
        cout << i << ",";
}
```

>---
#### 9.4. 跳转语句：break, continue, return, goto

```c++
int f(int i)
{
    int c = 0;
    while (i--) {
        if (i > 100) break;
        if (i % 2 == 0) continue;
        cout << i << ",";
        ++c;
    }
    cout << endl;
    return c;
}
int main() {
    auto c = f(10);    // 9,7,5,3,1
    printf("Iterate %d\n", c);  // 5
    c = f(1000);  // break
    printf("Iterate %d\n", c);  // 0
}
```

> **goto**

```c++
int main() {
    int a = 100;
    goto Test2;
    cout << "testing" << endl;
Test1:
    cerr << "At Test1 label." << endl;
    if(a < 10)
    {
    Test2:
        cerr << "At Test2 label." << endl;
    }
}  // At Test2 label.
```

---
### 10. 函数

函数声明可以作为非成员函数、友元函数、成员函数。返回类型不能是函数或数组。
- 非静态成员函数可以具有 CV 和 REF 限定。
- 非成员函数可以具有链接属性（`extern` or `static`）。
- `requires` 约束属于函数签名的一部分。`noexcept` 声明函数异常规范。
- 对于复杂类型，函数声明可以是尾随返回形式（`auto`）。
- `inline` 声明函数内联。

```cpp
FuncDecl     = Func ( Params ) [CV] [REF] [Except] [Attr]
TrailingFunc = auto FuncDecl -> <trailing>
TemplateFunc = template <Params> [RequiresClause] FuncDecl
```
```cpp
// 非成员函数
void func() noexcept;
static int func(int, ...);        // 变长参数, 内部链接  
auto func(int) -> int (*)(int);
struct S {
	// 成员函数
	virtual void func(int) const = 0;    // 纯虚函数
	void func(int, int) volatile&;		
	template <typename T>				 // 函数模板, 静态
	static auto func(T) -> decltype(new T{});
	friend void func(const S& s);    // 友元函数
};
```

非静态成员函数 `virtual` 声明（纯）虚函数，`override` 重写基类虚函数，`final` 声明密封。特殊成员函数或比较运算符函数可以是 `default` 预置或 `delete` 弃置。

```cpp
struct B1 {
	virtual void FuncA() = 0;   // 纯虚函数
	virtual void FuncB() {};    // 隐式 inline
};
struct B2 {
	void Func() = delete;       // 弃置
	B2() = default;				// 预置
	B2(B2&) = default;
};
struct D : B1, B2 {
	void B1::FuncA() override final {}  // 纯虚重写，密封
	void FuncB() override {}            // 重写
	bool operator <=> (const D& other) const = default;   // 默认比较运算符
    static inline void inFunc() {}		// 内联函数
    void FuncD() noexcept try {	}		// 函数 try 块
	catch (...) { }
};
```

类中定义的成员函数、弃置函数、`default` 或隐式生成的函数、首次声明 `constexpr` \ `consteval` 的函数、首次声明 `constexpr` 的静态数据成员等具有隐式 `inline`。

```cpp
struct Sample
{
	Sample();
	void Func();
	inline void InlineFunc();    // explicit inline
	// implicit inline
	Sample(const Sample&) = default;
	void implicitFunc() {}
	void* operator new(std::size_t) = delete;  
	Sample* operator &() = delete;  

};
Sample::Sample() {}  // not inline
inline void Sample::Func() {};   // explicit inline
void Sample::InlineFunc() {};    // inline
```

> **返回类型推导**

返回声明为类型占位符 `auto` / `decltype(auto)` 表示使用返回类型推导。虚函数和协程不支持返回类型推导。

```cpp
int x = 0;
auto func() { return x; }        // auto is int, reutun int
const auto& func() {return x; }  // auto is int, return const int& 

decltype(auto) funcA() { return x; }          // decltype(x) is int
decltype(auto) funcB() { return (x); }        // decltype((x)) is int&
decltype(auto) funcC() { return move(x); }    // decltype(move(x)) is int&&
```

> **参数列表**

非静态成员函数默认具有隐式 `this`。显式声明 `this` 参数的成员函数非虚非静态，可以是 lambda，无 CV 和引用限定。

```cpp
struct S {
    void F(void) const;
	void F(this S&);			// 显式对象参数
	void F(int x);				// 具名参数声明
	void F(double = 3.14);		// 默认参数
	void F(int, int*, int (*(*)(double))[3] = nullptr);   // 抽象参数声明
	void F(int x = 0, int y = 0, ...);       // 变参
    // void (*f)(this S) = [](this S s) {};
    template<typename... Args>
	void F(Args..., ...);		// 带形参包的变参函数模板
};
```

变长参数由 `<cstdarg>` 访问，同类型可由 `std::initializer_list` 访问。

```cpp
template<typename ... Ts>
void Iterator(Ts ... args) {
	for (std::initializer_list l{ args... }; auto e : l)
		cout << e << ",";
};

Iterator(1, 2, 3, 4, 5);
```

>---

#### 10.1. lambda 表达式

lambda 支持闭包，默认 `const` 限定，除非是 `mutable` 或显式 `this`。`mutable` 允许修改按复制捕获的对象、调用非 `const` 成员函数。`static` 具有空捕获规范。

```cpp
[Captures] [FrontAttr] (Params)? [[Specs][Except][BackAttr]] [-> Trailing][Requires] { Body };
[Captures] <TParams> [TRequires] [FrontAttr] (Params)? [[Specs][Except][BackAttr]] [-> Trailing][Requires] { Body };
	[Specs]: mutable|static, constexpr|consteval
```
```cpp
[] (int v) static constexpr -> int {return v;}    // 无捕获静态 lambda
[] <class C> -> decltype(auto) { return C{}; };   // lambda 模板
[=, &i] { i++; };           // 默认按复制捕获，引用捕获 i
[&, i] { return i * i; };   // 默认按引用捕获，复制捕获 i
```

> Captures：捕获

`[ ]` 无捕获；`[=]` 按值捕获；`[&]` 按引用捕获；非静态成员函数中可以捕获 `this`，`*this` 按值传递。lambda 闭包不会延长按引用捕获对象的生存期。
 
```cpp
struct S { void f(int i); };
void S::f(int i)
{
	[] static { };   // 无捕获，静态
	[&] {};          // 默认按引用捕获
	[&, i] {};       // 按引用捕获，但 i 按值捕获
	[&, this] {};    // 等价于 [&]
	[&, this, i] {}; // 等价于 [&, i]
	[=] {};          // 默认按复制捕获
	[=, &i] {};      // 按复制捕获，但 i 按引用捕获
	[=, *this] {};   // 按复制捕获外围的 S
	[=, this] {};    // 等价于 [=]
}
```

捕获可以具有初始化器，作用域为 lambda 表达式体，初始化时确定闭包捕获。

```cpp
int x = 4;  // global x
auto f = [&r = x, x = x + 1]() -> int {   // local x = global x + 1 = 5
    r += 2;         // r = & global x + 2
    return x * x;   // local x * local x
}; 
auto y = f(); // 更新 ::x 到 6 并初始化 y 为 25。
y = f();      // 更新 ::x 到 8, y 仍为 25。
```

如果 lambda 表达式作为默认实参，它不能捕获任何内容，除非所有捕获都带有初始化器，并满足参数约束条件。

```cpp
struct S
{
	int i = 1;
	// ERR：有捕获内容
	void g1(int = [i] { return i; }());
	void g2(int u, int = [i] { return 0; }()); 
	void g3(int = [=] { return i; }()); 
	void g4(int = [=] { return 0; }());      // ERR：非局部 lambda 具有默认捕获 
	void g5(int = [x = i] { return x; }());  // ERR：i 不能在默认实参中出现
	// OK：无捕获
	void g6(int = [] { return sizeof i; }()); 
	void g7(int = [x = 1] { return x; }()); 
	
}
```

> **闭包类型**

lambda 是纯右值的闭包类型，仅当捕获为空时该闭包类型是结构化类型。闭包类型具有以下成员：

```cpp
/* 辅助声明 */
using Params = void;
using Trailing = void;
using Body = void;
using Ta = int; using Tb = int; using Tc = int;
using Tx = int; using Ty = int; using Tz = int;
/* 闭包类型成员假定 */
using Lambda = decltype([](Params) /* Specs */ noexcept ->Trailing { Body; });  // (Params) Specs
struct LambdaType : Lambda {
	// Lambda::operator()(params): 仅当非 mutable 或没有显式对象形参时默认为 const
	/* Specs */ auto operator()(Params) const /* 默认行为 */ noexcept -> Trailing { Body; };
	// Lambda::operator returnType(*)(params): 仅当无捕获或没有显式对象形参时定义该用户转换函数
	using F = Trailing(*)(Params);
	constexpr operator F() const noexcept;  // 调用 Lambda::operator()
	// Lambda::Lambda()
	LambdaType() = default;    // 仅当无捕获时
	LambdaType(const LambdaType&) = default;
	LambdaType(LambdaType&&) = default;
	// Lambda::operator=(const Lambda&)
	LambdaType& operator=(const LambdaType&) = default;  // 无捕获时预置，其他情况弃置 delate
	LambdaType& operator=(LambdaType&&) = default;       // 仅当无捕获
	// Lambda::~Lambda()
	~LambdaType() = default;
	// Lambda::Captures: 复制捕获 [a,b,c]
	Ta a;
	Tb b;
	Tc c;
	// Lambda::Captures: 引用捕获 [&x,&y,&z]
	Tx& x;
	Ty& y;
	Tz& z;
};
```

类似的模板闭包类型：

```cpp
/* 辅助声明 */
using Params = void;
using Trailing = void;
using Body = void;
using Ta = int; using Tb = int; using Tc = int;
using Tx = int; using Ty = int; using Tz = int;
/* 闭包类型成员假定 */
using Lambda = decltype([]<class... TParams>(Params) /* Specs */ noexcept ->Trailing { Body; });
template<class ...TParams>
struct LambdaType : Lambda {
	// Lambda::operator()(params): 仅当非 mutable 或没有显式对象形参时默认为 const
	template<TParams...> /* Specs */ auto operator()(Params) const /* 默认行为 */ noexcept -> Trailing { Body; };
	// Lambda::operator returnType(*)(params): 仅当无捕获或没有显式对象形参时定义该用户转换函数
	template<TParams...params> using Fptr_t = Trailing(*)(Params);
	template<TParams...params> constexpr operator Fptr_t<params...>() const noexcept;
	// Lambda::Lambda()
	LambdaType() = default;    // 仅当无捕获时
	LambdaType(const LambdaType&) = default;
	LambdaType(LambdaType&&) = default;
	// Lambda::operator=(const Lambda&)
	LambdaType& operator=(const LambdaType&) = default;  // 无捕获时预置，其他情况弃置 delate
	LambdaType& operator=(LambdaType&&) = default;       // 仅当无捕获
	// Lambda::~Lambda()
	~LambdaType() = default;
	// Lambda::Captures: 复制捕获 [a,b,c]
	Ta a;
	Tb b;
	Tc c;
	// Lambda::Captures: 引用捕获 [&x,&y,&z]
	Tx& x;
	Ty& y;
	Tz& z;
};
```

>---
#### 10.2. 协程

协程的激活帧分为堆栈帧和协程帧：协程帧保留恢复协程状态的信息，堆栈帧在协程暂停时释放。

[**协程**](https://lewissbaker.github.io/2017/09/25/coroutine-theory) 将函数的调用和返回操作中的一些步骤拆分出一些额外的操作：*Suspend*、*Resume*、*Destroy*：
- *Suspend* 在 `co_await` 或 `co_yield` 处暂停协程并返回到调用方，且不会销毁协程的激活帧。当协程触及 *Suspend* 时将当前调用堆栈寄存器保存的任何值写入协程帧，然后保存一个值标记协程挂起点（由后续 *Resume* 或 *Destroy*）。协程返回调用方之前额外附加一个可以访问协程帧的句柄 *handle*，用于在稍后 *Resume* 或 *Destroy* 协程帧。
+ *Resume* 恢复暂停点协程的执行，并重新激活协程帧。*resumer* 调用 *handle* 的 `handle.resume()` 方法，重新分配一个堆栈帧、加载协程帧来跳转到挂起点。当下次挂起或运行完成时，调用 `handle.resume()` 将返回并恢复调用方的执行。
- *Destroy* 销毁激活帧并停止协程，释放协程激活帧资源。*Destroy* 重新激活协程激活帧并将执行转移至析构路径，即在挂起点调用所有局部变量的析构函数，最后释放协程帧使用的内存。*Destroy* 由 `handle.destroy()` 调用。
+ co-*Return* 将返回值存储在某个位置（可由协程自定义），然后析构所有的局部变量。在返回调用方可以执行一些额外的操作（例如执行某些操作发布返回值，或恢复正在等待返回值的另一个协程/异步）。

协程定义了 *Promise* 和 *Awaitable* 两种接口。
- *Promise* 指定协程本身行为的方法，例如定义调用协程时发生和返回、异常时发生的情况，或定义 `co_await` 和 `co_yield` 表达式的行为。
- *Awaitable* 指定控制 `co_await` 语义的方法表达，例如是否暂停当前协程、在当前协程暂停后执行一些逻辑以安排协程稍后恢复、以及在协程恢复后执行一些逻辑以生成 `co_await` 表达式的结果等。 
 
> [**Awaitable**](https://lewissbaker.github.io/2017/11/17/understanding-operator-co-await)

支持 `co_await` 运算的类型为 *Awaitable* 类型。*Awaiter* 类型实现了作为 `co_await` 表达式的一部分调用的三个方法：

```cpp
struct Awaiter {
	bool await_ready() noexcept;
	void/bool await_suspend(coroutine_handle<>) noexcept;  
	Result/void await_resume() noexcept;                   
};
```

*Promise* 中可以通过其 `await_transform` 方法更改 `co_await` 表达式的含义。假设等待协程的 *Promise* 对象 `promise` 具有类型 `P`，它包含一个 `await_transform` 成员函数，`co_await <expr>` 首先将 `<expr>` 传递到 `promise.await_transform(<expr>)` 中以获取 *Awaitable* 对象值 `awaitable`；若 `P` 类型不包含这样的成员，则直接将 `<expr>` 的结构作为 *Awaitable*。

然后若 `awaitable` 具有 `operator co_await()` 操作重载，则通过该函数获取 *Awaiter* 对象 `awaiter`；否则直接将 `awaitable` 用作 `awaiter`。

这些 `co_await` 规则可以抽象为：

```cpp
template<typename P, typename T>
decltype(auto) get_awaitable(P& promise, T&& expr)
{
  if constexpr (has_any_await_transform_member_v<P>)
    return promise.await_transform(static_cast<T&&>(expr));
  else
    return static_cast<T&&>(expr);
}

template<typename Awaitable>
decltype(auto) get_awaiter(Awaitable&& awaitable)
{
  if constexpr (has_member_operator_co_await_v<Awaitable>)
    return static_cast<Awaitable&&>(awaitable).operator co_await();
  else if constexpr (has_non_member_operator_co_await_v<Awaitable&&>)
    return operator co_await(static_cast<Awaitable&&>(awaitable));
  else
    return static_cast<Awaitable&&>(awaitable);
}
```

`co_await <expr>` 的语义可以抽象为：

```cpp
template<typename P, typename T>
decltype(auto) co_await_expr(P& promise, T&& expr) {
	auto&& value = expr;
	auto&& awaitable = get_awaitable(promise, static_cast<decltype(value)>(value));
	auto&& awaiter = get_awaiter(static_cast<decltype(awaitable)>(awaitable));
	if (!awaiter.await_ready())
	{
		using handle_t = std::coroutine_handle<P>;
		using await_suspend_result_t =
			decltype(awaiter.await_suspend(handle_t::from_promise(promise)));

		// <suspend-coroutine>

		if constexpr (std::is_void_v<await_suspend_result_t>)
		{	// 校验 awaiter.await_suspend 返回类型是 void 或 bool
			awaiter.await_suspend(handle_t::from_promise(promise));  // void await_suspend 无条件返回调用方
			// <return-to caller or resumer>
		}
		else 
		{
			static_assert(std::is_same_v<await_suspend_result_t, bool>,
				"await_suspend() must return 'void' or 'bool'.");
			if (awaiter.await_suspend(handle_t::from_promise(promise))) {  // false 将继续执行协程
				// <return-to caller or resumer>
			}
		}

		// <resume-point>
	}
	return awaiter.await_resume();
}
```

在 *`<suspend-coroutine>`* 暂停点，编译器生成一些代码保存协程的当前状态。例如存储 *`<resume-point>`* 的位置，转义协程调用栈寄存器中保存的任何值到协程帧内存中。

`awaiter.await_suspend()` 方法负责在操作完成后，安排协程在未来某个调用点恢复执行或销毁。`awaiter.await_ready()` 方法指示在已知操作已同步完成而无需挂起的情况下，避免 *`<suspend-coroutine>`* 操作产生的开销。

在 *`<return-to-caller-or-resumer>`* 点，执行返回到调用方并弹出本地堆栈帧，协程帧保持活动状态。当协程恢复时，执行在  *`<resume-point>`* 恢复，在调用 `awaiter.await_resume()` 方法之前获取操作的结果。

`awaiter.await_resume()` 返回值作为 `co_await` 表达式的结果。`await_resume()` 可以引发异常，异常会传播出 `co_await` 表达式。如果异常从 `await_suspend()` 发出，协程会自动恢复并直接传播出 `co_await` 表达式而不调用 `await_resume()`。

> **coroutine_handle**

`coroutine_handle` 类型标识协程帧的非拥有句柄，用于恢复协程执行或销毁协程栈；或访问协程的 *Promise* 对象。

```cpp
template<typename Promise>
struct coroutine_handle;

template<typename Promise>
struct coroutine_handle : coroutine_handle<void>
{
	Promise& promise() const;
	static coroutine_handle from_promise(Promise& promise);
	static coroutine_handle from_address(void* address);
};

template<>
struct coroutine_handle<void>
{
	bool done() const;
	void resume();
	void destroy();
	void* address() const;
	static coroutine_handle from_address(void* address);
};
```

`handle.resume()` 在 *`<resume-point>`* 点重新激活挂起的协程；`handle.destroy()` 会销毁协程帧，通常无需直接调用该方法，协程帧通常由某种调用协程的 *RAII* 类型所持有并调度，因此可能会出现双重销毁错误。`handle.promise()` 返回 *Promise* 对象引用，通常仅在设计 *Promise* 类型时有用；大多数情况可以选择 `coroutine_handle<void>` 作为 `awaiter.await_suspend()` 的参数类型。

`handle<P>.from_promise()` 允许从 `promise` 对象引用中重建协程 `handle`；必须确保类型 `P` 与 *Promise* 类型完全匹配，若从 *Promise* 类型为 *Derived* 的对象引用重建 `handle` 时，行为未定义。`coroutine_handle::from_address()` 或 `handle.address()` 将 `handle` 对象和 `void*` 之间互相转换；通常用于传递到 C 样式 API。 

协程恢复时首先调用 `awaiter.await_resume()` 获取结果，然后通常会立即销毁 `Awaiter` 对象（`await_suspend()` 调用的 `this` 指针）。然后协程可能会运行完成，进而销毁协程和 *Promise* 对象，所有这些操作都在 `await_suspend()` 返回之前完成。

因此在 `await_suspend()` 方法中，一旦协程可以在其他线程上恢复上，则需要确保避免访问 `this` 或协程的 `.promise()` 对象，它们可能已经被销毁。在作启动并计划协程恢复后，唯一可以安全访问的内容是 `await_suspend()` 中的局部变量。

使用协程时，可以利用协程帧中的局部变量在暂停时保持活动状态这一事实，避免为一些操作的状态（例如一些回调 API 或异步操作）分配堆分配。可以借助 *Awaiter* 对象在协程帧中借用内存，在 `co_await` 表达式的持续时间内存储每个操作的状态。协程帧可能仍会在堆上分配。但是分配后单个堆的协程帧可用来执行多个异步操作。

> **Promise** [^[↗]^](https://lewissbaker.github.io/2018/09/05/understanding-the-promise-type)

*Promise* 对象通过实现在协程执行期间的特定点调用的方法，定义和控制协程本身的行为。每次调用协程，都会在协程帧构造 *Promise* 对象实例，并在某些协程执行关键点生成对 *Promise* 对象上某些方法的调用。

协程函数的调用可以抽象为：

```cpp
SomeCoroutine() 
{ <body-statements> } ==> 
{
	co_await promise.initial_suspend();
	try
	{
		<body-statements>
	}
	catch (...)
	{
		promise.unhandled_exception();
	}
FinalSuspend:
	co_await promise.final_suspend();
}
```

协程开始执行 *`<body-statements>`* 之前，首先：
- 使用 `operator new`（可选）分配协程帧内存。异常导致分配协程帧失败的情况，在不允许异常的环境中，可选提供静态成员函数 `P::P get_return_object_on_allocation_failure()`，在 `operator new` 返回 `nullptr` 时立即调用该函数作为替代，而不是引发异常。
- 将任何函数参数复制到协程帧。
- 调用 `P` 类型 `promise` 对象的构造函数。允许在构造函数中访问复制后的参数；引发异常则销毁参数副本并释放协程帧，异常传播至调用方。
- 调用 `promise.get_return_object()` 获取在协程首次暂停或运行完成时返回给调用方的结果，并另存为局部变量。
- 执行 `co_await promise.initial_suspend()` 等待结果。这允许控制协程在执行协程正文 *`<body-statements>`* 之前是否可以暂停。若在 `.initial_suspend()` 点暂停，稍后可以通过 `handle.resume()` 恢复或 `handle.destroy()` 销毁。大多数协程函数中，`.initial_suspend()` 通常设计为返回 `suspend_always` 或 `suspend_never`。
- 当 `co_await promise.initial_suspend()` 恢复时，协程开始执行 *`<body-statements>`*。

当执行到达 `co_return;` 或 `co_return <expr>` 时（执行至协程末尾相当于在末尾有一个 `co_return;`），执行：
- 调用 `promise.return_void()` 或 `promise.return_value( <expr> )`。
- 以创建的相反顺序销毁所有具有自动存储期的变量。
- 执行 `co_await promise.final_suspend()` 等待结果。

若由于未处理的异常而离开 *`<body-statements>`*：
- 捕获异常并在 `catch` 块中调用 `promise.unhandled_exception()`。该方法的实现通常调用 `std::current_exception()` 来捕获异常的副本，以便在不同的上下文中重新引发；或直接 `throw;`。
- 执行 `co_await promise.final_suspend()` 等待结果。允许在协程返回之前在 `.final_suspend()` 执行一些额外的逻辑，例如发布结果、发出完成信号等。调用 `handle.resume()` 行为未定义，对于挂起的协程，唯一可选的操作是 `handle.destroy()`。

协程执行结束或异常终止后，协程帧被销毁：
- 调用 `promise.~promise()`（如果有或默认析构函数）。
- 调用函数参数的析构函数。
- 调用 `operator delete`（可选）删除协程帧内存。
- 执行返回调用方。

当执行首次到达 `co_await` 内的  点时，或者协程运行完成而不命中 <return-to-caller-or-resumer> point 的 Sym 中，则协程要么被挂起，要么被销毁，并且 return-object 之前 然后，从对 promise.get_return_object（） 的调用返回给协程的调用方。

当执行首次到达协程内 `co_await` 表达式中的 *`<return-to-caller-or-resumer>`* 点，或者协程运行至结束而未到达 *`<return-to-caller-or-resumer>`* 点时，协程被挂起或被销毁，然后先前从 `promise.get_return_object()` 结果对象将返回给协程的调用方。

`promise` 对象的类型是由模板 `std::coroutine_traits` 类根据协程的签名确定的。

```cpp
// 签名								模板 coroutine_traits 实例化
task<void> foo(int x);	            std::coroutine_traits<task<void>, int>::promise_type
task<void> Bar::foo(int x) const;	std::coroutine_traits<task<void>, const Bar&, int>::promise_type 
task<void> Bar::foo(int x) &&;	    std::coroutine_traits<task<void>, Bar&&, int>::promise_type
```

`template<typename RET, ...> coroutine_traits` 默认定义 `promise_type` 并通过 `RET` 查找 `promise_type` 的定义或 `typedef`。 

```cpp
template<typename RET, typename... ARGS>
struct coroutine_traits<RET, ARGS...> {
	using promise_type = typename RET::promise_type;
};

template<typename T>
struct task{
	using promise_type = task_promise<T>;
    ...
};
```

*Promise* 可以定义 `P::await_transform()` 方法来自定义 `co_await` 表达式的行为。`co_await <expr>` 被转换为 `co_await promise.await_transform(<expr>)`，因此可以启用某些非 *Awaitable* 的类型。也可以通过 `await_transform(T expr)` 重载声明为 `delete` 来禁止对类型 `T` 进行 `co_await`；

*Promise* 类型可以选择支持 `co_yield` 语义，需定义 `P::yield_value()` 方法。`co_yield <expr>` 被转换为 `co_await promise.yield_value(<expr>)`。

>---
#### 10.3. 运算符重载

```cpp
operator op     // op: + - * / % ^ & | ~ ! = < > += -= *= /= %= ^= &= |= 
                //     << >> >>= <<= == != <= >= <=> && || ++ -- , ->* -> () []
operator new | new[]
operetor delete | delete[]
operator co_await
```

`()`,`[]` 成员运算符可以声明为静态，其他静态声明的成员运算符仅能通过 `T::operaror op()` 调用。非成员函数的运算符操作数至少有一个具有类类型或枚举类型。

```cpp
struct S {
	int d;
	S operator + (int v) { return (d += v, *this); }
	S operator -(int v) { return (d -= v, *this); }
	S(int v) : d{ v } {}
};
// 非成员函数
S operator +(S s1, S s2) { return s1.d + s2.d; }

int main() {
	S s = 100;
	s = s + 1 - 2 + S(10);
}
```

`operator <=>` 三路比较运算符用于简化比较运算符的实现，编译器自动生成 `==`,`!=`,`<`,`>`,`<=`,`>=` 运算符重载。`<=>` 返回 `std::strong_ordering`、`std::weak_ordering` 或 `std::partial_ordering`。

```cpp
struct Point {
	int x, y;
	auto operator<=>(const Point&) const = default; // 默认实现, 按字典序
};

int main() {
	Point a{ 1, 2 }, b{ 3, 4 };
	cout << boolalpha << (a > b) << (a < b);   // false, true
}
```

将 `std::istream&` 或 `std::ostream&` 作为左侧参数的 `operator >>` 和 `operator <<` 的重载称为插入和提取运算符，必须将其实现为非成员函数或友元函数。

```cpp
std::ostream& operator<<(std::ostream& os, const T& obj) {
    // write obj to stream
    return os;
}
 
std::istream& operator>>(std::istream& is, T& obj) {
    // read obj from stream
    if (/* T could not be constructed */)
        is.setstate(std::ios::failbit);
    return is;
}
```

后缀递增递减通常根据前缀版本实现：

```cpp
struct Counter
{
	int count = 0;
	// ++Counter
	Counter& operator++() { return (++count, *this); }
	// Counter++
	Counter operator++(int)
	{
		Counter old = *this; // copy old value
		operator++();		 // prefix increment
		return old;          // return old value
	}
};
```

二元运算符通常实现为非成员，以保持对称性（*T op v*, *v op T*）；而成员函数仅支持 *T op v* 而不是 *v op T*；

```cpp
struct Point {
	int x, y;
	Point operator * (int m) { return Point(x * m, y * m); }
};

Point operator + (const Point& p1, const Point& p2) { return Point(p1.x + p2.x, p1.y + p2.y); }
int main() {
	Point p1{ 1,2 }, p2{ 3,4 };
	Point p = p1 + p2;   // T+v, v+T;
	Point p3 = p * 10;   // T*v; 
}
```

>---
#### 10.4. 用户定义文本函数

用户定义文本可以是整数、浮点数、字符或字符串用户定义。字符串文本函数需要两个参数：`const charType*` 和 `size_t`（允许单个参数 `const char *` 的原始字面量）。

```c++
using namespace std;
long double operator""_w(long double);
string operator""_w(const char16_t *, size_t);  // u"str"_w
unsigned operator""_w(const char *);	
int main()
{
    1.2_w;	  // operator ""_w(1.2L)
    u"one"_w; // operator ""_w(u"one", 3)
    12_w;	  // operator ""_w("12")
    "two"_w;  // error: 匹配 operator""_w(const char *, size_t)
}
```

---
### 11. 类类型

类类型为值类型 `class`、`struct`、`union`，其中 `union` 隐式密封。类成员可以声明静态或非静态的数据成员和成员函数、*typedefs* 或 *usings*、成员枚举、嵌套类和友元声明。具有内部定义的函数隐式内联，除非是命名模块导出（`export`）。

```cpp
class S {
private:
    int d1;            		       
    int a[10] = {1, 2}; 	       
    static const int d2 = 1;       // static 
protected: 
    virtual void f1(int) = 0;      // 纯虚函数
public:
    enum { NORTH, SOUTH, EAST, WEST };
    struct NestedS {
        std::string s;
    } d5, *d6;
    typedef NestedS value_type, *pointer_type;
	static constexpr decltype(auto) str = "hello";
};
```

`using` 声明可以将基类成员引入派生，或引入枚举成员。

```cpp
enum class color { red, orange, yellow };
class Base {
protected:
    int d;
};
class Derived : public Base {
public:
	using Base::d;    // make Base's protected member d a public member of Derived
	using Base::Base; // inherit all bases' constructors
	using enum color;
};
```

函数域中的局部类不包含静态数据成员，成员函数内部定义，无友元模板或友元函数，非闭包类型的局部类没有成员模板。

```cpp
int main()
{
	std::vector<int> v{ 3,4,5,7,2,8,2,4 };
	struct greater {
		bool operator()(int n, int m) { return n > m; }
	};
	std::sort(v.begin(), v.end(), greater()); // 降序
	for (int n : v)
		std::cout << n << ',';
}
```

非静态数据成员可以具有默认初始化器，静态数据成员仅声明（除非是 `constexpr`,`const`,`inline`）和外部定义。

```cpp
struct S {
	int v = 10086;
	const int cv = v;
	static int sv;  // 仅声明
	static const int scv = 10;
	static constexpr int d2 = 10;
};
int S::sv = 10010;  // 外部定义
```

嵌入匿名类只包含公共非静态数据成员。

```c++
struct S {
	union {
		int  ui;
		long ul;
	};
	class {
	public:
		int c;
	};
	struct {
		int s;
	};
};

S s{ 1,2,3 };
```

表达式 `this` 是一个 **纯右值** 表达式，表示为隐式对象形参的地址。CV 成员函数的 `this` 也是 CV 限定，且仅由对应 CV 对象调用；REF 成员函数仅限定隐式对象，不会限定 `this` 指针；构造函数和析构函数的 `this` 始终为 `T*`。可以 `delete this`，但确保 `this` 对象是由 `new` 分配的。

```cpp
struct Counter {
	static Counter* GetDeCounter(size_t count) {
		auto p = new Counter;
		p->v = count;
		return p;
	}
	void Inc() { ++v; }
	void Dec() {
		cout << v << endl;
		if (--v <= 0) 
			delete this;
	}
	~Counter() { cout << "delete this\n";; }
private:
	size_t v;
};
int main() {
	auto c = Counter::GetDeCounter(1);
	c->Dec();  // delete this
	// c 变为悬挂指针，对后续的调用行为未定义
	c->Dec();  // maybe 216544654648
	c->Dec();  //		216544654647
}
```

非虚非静态成员函数（无 CV 和 REF 限定）可以声明显式对象形参，允许推导类型（`this auto`）和值类别（`this T`）。***指向显式对象成员函数的指针是普通的函数指针，而不是成员指针***。

```cpp
struct S {
	int f(this auto);
	int g(this const S&, int, int);  // 显式对象形参
};
int main() {
	S s{};
	auto pf = &S::g;   // 非成员指针
	pf(s, 3, 4);	// ok
	(s.*pf)(3, 4);  // error: “pf” is not a pointer to member function
}
```

> *可变数据成员*

`mutable` 声明非静态、非常量和非引用的可变数据成员，它可以在 `const` 限定函数或 lambda 中修改。 

```c++
class X {
private:
    bool m_flag;
    mutable int m_accessCount;
public:
    bool GetFlag() const {
        m_accessCount++;   // ok
        // m_flag = true;  // err
        return m_flag;
    }
};
```

> **可访问性**

`class` 成员和基类默认 `private`。`struct` 成员和基类默认 `public`。`union` 成员默认 `public`。基类可声明 `public`、`protected`、`private`，限定继承成员在派生中的访问控制。

```cpp
struct A { int a; };
struct B : private A { int b = a; };   // a is private
struct C : public B {
	int c = a;  // err; 不可访问
};
```

> *注入类名* 

​*注入类名* ​是指类或类模板范围内，类名被视为公共成员名，可以直接使用而无需通过作用域解析（::）显式引用。注入类名可继承，非公共继承间接基类的注入类名可能在派生类中不可访问。

```cpp
int X;
struct X {
	X* p;    // OK, X is an injected-class-name
	::X* q;  // Error: name lookup finds a variable name, which hides the struct name
};

template<class T>
struct Y {
	Y* p;    // OK, Y is an injected-class-name
	Y<T>* q; // OK, Y is an injected-class-name, but Y<T> is not
};

struct A {};
struct B : private A {};
struct C : public B {
	A* p;    // Error: injected-class-name A is inaccessible
	::A* q;  // OK, does not use the injected-class-name
};
```

>---
#### 11.1. 构造函数

构造函数可以是 `friend`、`inline`、`constexpr`、`consteval` 或 `explicit`，可具有成员初始化列表。`constexpr` 构造函数使其类型成为字面类型。

在某些情况下，特殊成员函数可由编译器默认生成，它们是 *默认构造*、*复制构造*、*移动构造*、*复制赋值运算符*、*移动赋值运算符*、*析构函数*，可以显式声明弃置（`= delete`）以阻止默认生成。

```cpp
struct Point {
    constexpr Point() :x{}, y{} {};  // 带有成员初始化表达式的默认构造函数
    explicit Point(int x, int y) :x(x), y(y) {}; // 重载 
    int x, y;
};
constexpr Point origin{};            // 字面类型
```

> *转换构造函数*

*转换构造函数* 至少带有一个非默认参数。隐式声明或用户定义的非 `explicit` 复制构造和移动构造都是转换构造函数。

```cpp
struct S {
	explicit S(std::string) {}    // 仅显式转换
	S(int) {};   				  // 支持隐式转换
	S(int, int) {};
};
int main() {
	S s1 = 1;          // ok
	S s2(3);           // ok
	S s3 = { 5, 5 };   // ok
	S s4 = {};		   // err
	S s5 = "Hello";	   // err
	S s6 = (S)"Hello"; // ok，显式转换
}
```

转换运算符用于将类对象转换为目标类型：

```c++
struct money {
private:
    long double account;
public:
    money(long double n) : account{ n } {};
    operator long double() {   // 转换运算符
        return account;
    };
};
int main() {
    long double d0 = money(3.14);
    long double d1 = money(3.14L);
    long double d2 = money(10010);
    long double d3 = money(2.71f);
}
```

> *委托构造函数*

构造函数可以引用其他构造函数，不包含成员初始化器。

```c++
struct Point {
    Point(int x, int y, int z) :x{ x }, y{ y }, z{ z } {};
    Point(int x, int y) : Point(x, y, 0) {};
    Point() : Point(0, 0, 0) {};
    int x, y, z;
};
```

> *继承构造函数*

`using Base::Base` 将基类构造引入到派生。

```cpp
struct Base {
	int x, y;
	Base(int x, int y) : x{ x }, y{ y } {}
};
struct Derived : Base {
	using Base::Base;   // 继承 Base 的所有构造函数
	Base b = Derived { 1,2 };  // ok
};
```

> **复制构造函数**

复制构造函数从同类型对象复制成员值来初始化对象。对于指针类成员，默认复制构造只会复制指针值，因此可能需要用户定义 *分配指针内存* 操作。复制构造函数的首元参数为 `T&`，可以具有 CV 限定；若包含其他参数，需要具有默认实参。

```cpp
Point(Point& other);   				
Point(const Point& other);
Point(volatile Point& other);
Point(volatile const Point& other);
```

定义复制构造函数还需要定义相应的复制赋值运算符，无默认参数。

```c++
Point& operator=(Point& other);
Point& operator=(const Point& other);
Point& operator=(volatile Point& other);
Point& operator=(volatile const Point& other);
```

下面是一个可复制 *buffer*：

```cpp
struct buffer {
private:
    char *buf;
    size_t index, size;
public:
    buffer(size_t size) : size{size}, buf{new char[size]}, index{0} {}
    ~buffer() { delete[] buf; }
    buffer(const buffer &other) noexcept : size(other.size), index{other.index}, buf{nullptr} {
        buf = new char[size];
        copy(other.buf, other.buf + size, buf);
    };
    buffer &operator=(buffer &other) noexcept {
        if (this != &other) {
            delete[] buf;
            size = other.size;
            index = other.index;
            buf = new char[size];
            copy(other.buf, other.buf + size, buf);
        }
        return *this;
    }
    size_t Write(string str) {
        size_t n = str.length();
        size_t able = size - index - 1;
        if (str.length() >= able)
            n = able;
        copy(str.c_str(), str.c_str() + n, buf + index);
        index += n;
        buf[index] = '\n';
        return n;
    }
    operator string() {
        string s = string{buf, index};
        index = 0;
        buf[0] = '\0';
        return s;
    }
};
int main() {
    buffer b{512};
    b.Write("Hello World! ");
    buffer b2 = b;
    b2.Write("JimryYchao\n");
    cout << string(b2) << endl;  // Hello World! JimryYchao
}
```

> 移动构造函数

类似于复制构造函数，移动构造函数首元参数为 `T&&`，可以具有 CV 限定；若包含其他参数，需要具有默认实参。

```cpp
Point(Point&& other);   				
Point(const Point&& other);
Point(volatile Point&& other);
Point(volatile const Point&& other);
```

定义移动构造函数还需要定义相应的移动赋值运算符，无默认参数。

```c++
Point& operator=(Point&& other);
Point& operator=(const Point&& other);
Point& operator=(volatile Point&& other);
Point& operator=(volatile const Point&& other);
```

为 `buffer` 声明移动构造函数：

```cpp
struct buffer {
// ...
public:
	buffer(buffer&& other) noexcept : buf(nullptr) {
		cout << "move buffer\n";
		*this = std::move(other);
	};
	buffer& operator = (buffer&& other) noexcept {
		if (this != &other) {
			delete[] buf;
			size = other.size;
			index = other.index;
			buf = other.buf;
			other.buf = nullptr;
		}
		return *this;
	}
// ...
}
int main() {
    buffer b{ 512 };
    b.Write("Hello World! ");
    buffer b3 = std::move(b);
    b3.Write("CXX");
    cout << string(b3) << endl;  // Hello World! CXX
}
```

>---
#### 11.2. 析构函数

析构函数在对象超出范围或 `delete` 销毁对象指针时自动调用。析构函数不可继承，如果基类将其析构函数声明为 `virtual`，则派生始终会重写它。这使得可以通过指向基类的指针删除多态类型的动态分配对象。

纯虚析构函数强制类型为抽象类，必须具有定义。虚析构函数主要用于多态场景下派生对象的正确析构。当通过基类指针删除派生类对象时，虚析构函数确保调用派生析构再调用基类析构；其继承链的最终基类必须是虚析构函数。

```cpp
struct BaseA {
	virtual ~BaseA() = 0 { cout << "~BaseA\n"; }  // 纯虚析构函数，必须具有定义
};
struct BaseB : BaseA {
	~BaseB() { cout << "~BaseB\n"; }
};
struct Derived : BaseB {
public:
	~Derived() {
		/* 释放派生类资源 */
		cout << "~Derived\n";
	}
};
int main() {
	BaseA* obj = new Derived();
	delete obj; 
	// 正确调用 Derived::~Derived() → BaseB::~BaseB() → BaseA::~BaseA()
	// 若基类析构函数非虚，delete obj 仅调用 BaseA::~BaseA()，导致派生类资源泄漏
}
```

>---
#### 11.3. default, delete

`default` 预置或 `delete` 弃置特殊成员函数。如果某个类型未声明它本身，则编译器自动生成默认构造函数、复制构造和移动赋值运算符、移动构造和移动赋值运算符、析构函数。其中：
- 显式声明任何构造函数，则不生成默认构造函数。
- 显式声明虚析构函数，则不生成默认析构函数。
- 显式声明移动构造或移动赋值运算符，则不生成复制构造函数和复制赋值运算符。
- 显式声明复制构造、复制赋值运算符、移动构造、移动赋值运算符或析构函数，则不生成移动构造函数和移动赋值运算符。
- 显式声明复制构造或析构函数，则不生成复制赋值运算符。
- 显式声明复制赋值运算符或析构函数，则不生成复制构造函数。

如果基类不拥有供派生类可访问（`public`，`protected`）的默认构造函数，那么派生类无法自动生成它自己的默认构造函数。

对于创建不可移动、只能动态分配或无法动态分配的用户定义类型；可以通过 `default` 和 `delete` 方式进行设定。弃置复制构造和复制赋值运算符，可以使用户定义类型不可复制：

```c++
struct noncopyable {
	noncopyable() = default;   
	noncopyable(const noncopyable&) = delete;
  	noncopyable& operator=(const noncopyable&) = delete; 
};
```

> *defualt*

`default` 预置任何特殊成员函数以默认方式自动实现。

```c++
struct widget {
    widget() = default;
    widget& operator=(const widget&) = default;
};
```

> *delete*

`delete` 弃置特殊成员函数和普通成员函数以及非成员函数，以阻止定义或调用它们。已删除的函数仍参与重载决策。

```c++
struct widget {  // 无法 new 和 & 取址
    void* operator new(std::size_t) = delete;
    widget* operator &() = delete;   
};
widget w{};
widget* pw1 = new widget;  // ERR; 无法调用删除的函数；但是可以调用全局 new ；
widget* pw2 = &w  // ERR; address-of 被删除
widget* pw3 = ::new widget;  // ok
delete pw3;
```

若要限制发生隐式类型转换，确保仅发生对 `double` 重载的调用，可声明一个弃置模板版本：

```c++
template < typename T >
void call_with_true_double_only(T) = delete;   // 禁用任意类型隐式转换为 double
void call_with_true_double_only(double param) {  } // also define for const double, double&, etc. as needed.

call_with_true_double_only(3.1415);  // just only double
```

>---
#### 11.4. 继承

继承从现有类派生新类；可以是单一继承或多重继承，`final` 声明封装。同名成员（非 `override`）隐藏基类继承成员。	

```c++
class Derived : Base-Class, ... ;
	Base-Class = [virtual] [access-specifier] BaseClass
	access-specifier = public | protected | private
```

多重继承可以构建一个继承关系图，其中相同的基类是多个派生类的一部分。多重继承使得沿多个路径继承名称成为可能，沿这些路径的类成员名称不一定是唯一的。由于一个类可能多次成为派生类的间接基类，这些名称冲突存在 “多义性”。任何引用类成员的表达式必须采用明确的引用，可以通过限定名称方式消除多义性。

```c++
struct A {
    unsigned a;
    unsigned b();
};
struct B {
    unsigned a();  // class A also has a member "a"
    int b();       //  and a member "b".
    char c;
};
class C : public A, public B {};
int main(){
    C *pc = new C;
    // pc->b();     // 歧义，b 不明确
    pb->B::b();  // 限定名称访问
}
```

通过一个继承关系图到达多个名称（函数、对象或枚举器）是可能的。从派生指针或引用转换到基类指针或引用可能会导致歧义。 

```c++
class A { };
class B : public A {};
class C : public A {};
class D : public B, public C {};

A* pa = new D;  // 从 D 到 A 的转换是不明确的；A 无法辨别从 B 还是 C 传递
// 需要显式指定要使用的子对象
A* pab = (A*)(B*)(new D);  // D -> B -> A
A* pac = (A*)(C*)(new D);  // D -> C -> A
```

`virtual` 基类可以避免多重继承的类层次结构中出现多义性。虚拟基类的数据成员在多态继承路径上拥有唯一副本，由基类和派生类共享。

```c++
class A {  };
class B : public virtual A {};
class C : public virtual A {};
class D : public B, public C {};  // A 具有唯一副本
A* pa = new D; 
```

虚函数可以在派生类中覆写（`override` 和 `virtual` 可选），纯虚函数（`=0`）强制类型为抽象类型。更改重写函数的访问限定不会影响多态行为。`consteval` 虚函数不得重写非 `consteval` 虚函数。使用限定函数标识（`Base::VirFunc`）调用基类虚函数。

```c++
struct Account {
   Account( double d ) { _balance = d; }
   virtual ~Account() {}
   virtual double GetBalance() { return _balance; }
   virtual void PrintBalance() { cerr << "Error. Balance not available for base type." << endl; }
private:
    double _balance;
};
struct CheckingAccount : public Account {
   CheckingAccount(double d) : Account(d) {}
   void PrintBalance() { cout << "Checking account balance: " << GetBalance() << endl; }
};
struct SavingsAccount : public Account {
   SavingsAccount(double d) : Account(d) {}
   void PrintBalance() { cout << "Savings account balance: " << GetBalance(); }
};
int main() {
   CheckingAccount checking( 100.00 );
   SavingsAccount  savings( 1000.00 );
   Account *pAccount = &checking;
   pAccount->PrintBalance();  // call checking.PrintBalance
   pAccount = &savings;
   pAccount->PrintBalance();  // call savings.PrintBalance
}
```

可以提供纯虚函数的缺省定义（外部定义），纯虚析构函数（必须提供）内部定义。

```cpp
struct Abstract {
	virtual void f() = 0; // pure virtual
	virtual void g() {}   
	~Abstract() {
		g();           // OK: calls Abstract::g()
		// f();        // undefined behavior
		Abstract::f(); // OK: non-virtual call
	}
};
void Abstract::f() { std::cout << "A::f()\n"; }  // 缺省定义
struct Concrete : Abstract {
	void f() override {
		Abstract::f(); // OK: calls pure virtual function
	}
	void g() override {}
	~Concrete() {
		g(); // OK: calls Concrete::g()
		f(); // OK: calls Concrete::f()
	}
};
```

函数 `Derived::f` 重写 `Base::f`，它们的返回类型必须相同或是协变的。

```cpp
struct Base {
	virtual Base* P() { return this; };
	virtual Base& LR() { return *this; };
	virtual Base&& RR() { return Base{}; }
};
struct Derived : Base {
	Derived* P() override { return this; };        // 协变 D* -> B*
	Derived& LR() override { return *this; };      // 协变 D& -> B&
	Derived&& RR() override { return Derived{}; }; // 协变 D&& -> B&&
};
```


>---
#### 11.5. 友元

友元在类体中声明，以授予友元函数或友元类访问该类内部成员的权限。友元关系不具有继承和传递性，无访问限定符。友元类声明不能定义新类。

```cpp
class Product {
private:
	int secretCode;
	Product(int code) : secretCode{code} {} 
	friend class Factory;     // 友元声明
};
class Factory {      
public:
	Factory() = default;
	static Product CreateProduct();
};
// Factory 通过友元访问 Product 的私有构造函数
Product Factory::CreateProduct(){
	static int secretCode = 0;
	return Product(++secretCode);
}
```


>---
#### 11.6. 位域

类类型可包含位域成员，必须是整数类型或枚举。匿名位域用于填充；宽度为 0 的匿名位域强制边界对齐。

```c++
struct Date {               // offset
   unsigned nWeekDay  : 3;    // 0..2 
   unsigned nMonthDay : 6;    // 3..9 
   unsigned           : 0;    // 边界对齐
   unsigned nMonth    : 5;    // 32..36
   unsigned nYear     : 8;    // 37..45 
};
```

无法取值位域和非常量引用。常量引用会创建一个临时对象，用位域的值进行复制初始化，并且引用绑定到该临时对象。

```c++
Date d{ 1,1,1,1 };
const auto& r = d.nMonth;       // 临时对象，非指向位域
```

>---
#### 11.7. 联合体

联合体隐式密封，仅保存其非静态数据成员中的一个，不包含引用非静态数据成员。联合体的大小至少与最大的数据成员大小相同，最多一个成员可以具有初始化器。

```cpp
union U {
	std::int32_t n;     // occupies 4 bytes
	std::uint16_t s[2]; // occupies 4 bytes
	std::uint8_t c;     // occupies 1 byte
	unsigned int i1 : 4 = 1;   
};
```

联合具有非静态类类型成员时，编译器自动将任何非用户特殊成员函数标记为 `delete`。如果联合是 `class` 或 `struct` 中的匿名联合，则 `class` 或 `struct` 的任何非用户特殊成员函数标记为 `delete`。对于联合体的成员是具有用户定义构造和析构的类类型，当切换活动成员时，通常需要显式调用析构函数和 *new* 初始化。

```cpp
union S {
	std::string str;
	std::vector<int> vec;
	~S() { }
};          
int main() {
	S s = { "Hello, world" };
	std::cout << "s.str = " << s.str << '\n';
	s.str.~basic_string();   // 显式调用析构函数
	// 切换活动成员为 vec
	new (&s.vec) std::vector<int>;
	s.vec.push_back(10);
	std::cout << s.vec.size() << '\n';
	s.vec.~vector();
}
```

匿名联合体仅包含公共非静态数据成员，它的成员注入到封闭的作用域。命名空间域的联合体限定为 `static`。

```cpp
// 文件范围或命名空间范围的匿名 union
static union {
    short       iValue;
    long        lValue;
    double      dValue;
}
// 结构中的匿名 union
struct Input {
    WeatherDataType type;，
    union {
        TempData temp;
        WindData wind;
    };
};
```

---
### 12. 模板

模板可以是类模板、函数模板、别名模板、变量模板、概念约束等，以若干 *模板参数*（类型、常量或其他模板）参数化。特化支持显式提供：对类、函数、变量模板全特化，或对类和变量模板部分特化。模板声明可以包含约束。变量模板在类型中声明时仅支持静态。别名模版始终不进行推导。

```cpp
// 类模板
template <Params> [RequiresClause] ClassDecl;
// 函数模板
template <Params> [RequiresClause] FuncDecl;
// 别名模板
template <Params> using Name = Typeid;
// 变量模板
template <Params> [RequiresClause] ValueDecl;
// 概念
template <Params> concept ConceptName = ConstraintExpr;

template <typename T>      // 概念
concept C = requires{ T{}; };

template<typename... T>    // 类型模板
struct value_holder {

	template<T... Values>  // 嵌套类型模板
    struct apply {}; 
	
	template <C X>         // 静态数据成员模板
	static X Value = X{};		  
	
	void Func(C auto c);   // 成员函数模板
};

template <typename... T>   // 别名模板
using value_holder_R = value_holder<T...>&;    

template<class T>          // 变量模板
constexpr T pi = T(3.1415926535897932385L); 
```

> 简化函数模板

形参出现占位符 `auto` 或 `Concept auto` 时，该声明为一个函数模板，每个占位符向模板形参列表追加一个模板形参：

```cpp
template <typename> concept T = true;
template <typename> concept U = true;

void f1(auto);					 // template<class T> void f(T)
void f2(T auto);				 // template<T X> void f2(X)
void f3(U auto...);				 // template<U... Ts> void f3(Ts...) 
void f4(const T auto*, U auto&); // template<T X, U Y> void f4(const X*, Y&);
```

>---
#### 12.1. **模板形参**

```cpp
<Params> = 
// 常量模板形参
    <Typeid [name] [= Val]>      // 默认值
    <Typeid ... [name]>          // 形参包
// 类型模板形参
    <class|typename|ConceptName [name] [= Val]> 
    <class|typename|ConceptName ... [name]>
// 模板模板形参
    <template <Params> class|typename [name] [= Val]>
    <template <Params> class|typename ... [name]>
```

`Typeid` 可以是 *结构化类型*、包含占位符的类型、被推导类类型的占位符。

```cpp
// 常量模板形参
template<typename T, decltype(auto) sz = 0>  // 占位
    requires is_integral_v<T> && (sz >= 0)
struct Buffer { T buffer[sz]; };
Buffer<unsigned char, 128> buf{};

// 类型模板形参
template<class T = void>
concept Addable = requires(T t1, T t2) { t1 + t2; };
template <Addable A = int>
struct Sample {
	A value;
	Sample operator +(Sample  other) {
		return Sample{ value + other.value };
	}
};

// 模板模板形参
template<Addable A, template<Addable = A> typename AddableType = Sample>
AddableType<A> Add(AddableType<A> a1, AddableType<A> a2) {
	return a1 + a2;
}

int main() {
	buf.buffer[0] = 1;
	auto sum = Add(Sample{ 3.14 }, Sample{ 6.28 });  // Sample<double>
	cout << sum.value;  // 9.42
}
```


> **形参包**

变参模板中包含形参包（模板形参包、函数形参包）。

```cpp
TemplateParam ... [PackName]

// 变参类模板
template<class... Types>
struct Tuple {};  

// 变参函数模板
template<class... Types>
void Func(Types... args);   
int main() {
	Func();			  // Func<>
	Func(1);		  // Func<int>
	Func(1, 2.0);     // Func<int,double>
}
```

包模式展开：形参包被展开成多个模式实例。嵌套包展开从最内层开始。`sizeof	...(Pack)` 返回形参包长度。

```cpp
template<class... Us> void f(Us... pargs) {}
template<class... Ts> void g(Ts... args) {
	f(&args...); // &E1, &E2, ... , &En
}
// g(1,0.2,"a")
// &args... = &E1(1),&E2(0.2),&E3("a")
// Us...    = int *, double *, const char **

template<typename... Ts>
void func(Ts... args)
{
    const int size = sizeof...(args) + 2;
    int res[size] = {1, (args + 1)..., 2};   // a1+1,a2+1,...,
}
```

包展开可以在模板形参列表、基类列表、lambda 捕获子句中出现。

```cpp
template<class... Mixins>
class X : public Mixins... {  // public B1, public B2,...
public:
    X(const Mixins&... mixins) : Mixins(mixins)... {}  // 同时需要在构造函数中使用包展开
	// X(const B1& b1, const B2& b2,...) : B1(b1),B2(b2),... {}
};
```

包展开可以在关键词 `alignas` 所用的类型列表和表达式列表中使用。实例以空格分隔：

```cpp
template<class... T>
struct Align{
    alignas(T...) unsigned char buffer[128];
};
Align<int, short> a; // alignas(int) alignas(short) buffer
```

包可以由折叠表达式展开。`op` 支持 `+`, `-`, `*`, `/`, `%`, `^`, `&`, `|`, `=`, `<`, `>`, `<<`, `>>`, `+=`, `-=`, `*=`, `/=`, `%=`, `^=`, `&=` ,`|=`, `<<=`, `>>=`, `==`, `!=`, `<=`, `>=`, `&&`, `||`, `,`, `.*`, `->*`。对于一元展开空包，`&&` 返回 `true`，`||` 返回 `false`，`,` 返回 `void()`。

```cpp
Pack = {E1,E2,...,En};
(Pack op ...) = (E1 op (... op (En_1 op En)))      // 一元右折叠
(... op Pack) = (((E1 op E2 ) op ...) op En)       // 一元左折叠
(Pack op ... op V) = (E1 op (... op (En op V)))    // 二元右折叠
(V op ... op Pack) = (((V op E1) op ...) op En)    // 二元左折叠

template<typename... Args>
bool all(Args... args) { return (... && args); }
 
bool b = all(true, true, true, false);
// 在 all() 中，一元左折叠展开成
// return ((true && true) && true) && false;
```

>---
#### 12.2. 模板实例化

显式实例化定义其所指代的模板类型。`extern` 声明跳过隐式实例化步骤，表明该模板已在别处显式定义。模板显式特化在显式实例化之前出现，则显式实例化无效；反之实例化定义在特化之前，程序非良构。

```cpp
/* fileA.cpp */
template<typename T>
void Func(T v) { std::cout << v << std::endl; };   // 模板定义

template<> void Func<bool>(bool v) {     // Func<bool> 显式特化
	std::cout << std::boolalpha << v << std::endl;
};
template void Func<bool>(bool v);            // 已存在显式特化，显式实例化定义无效
template void Func<double>(double v);        // Func<double> 显式实例化定义
//template<> void Func<double>(double v) {}  // err 显式特化之前已有显式定义

/* fileB.cpp */
extern template void Func<bool>(bool);       // 引用外部显式特化定义
extern template void Func<double>(double);   // 引用外部显式实例化定义
extern template void Func<string>(string);   // 无效引用，未显式定义或特化
int main() {
	Func<int>(10086);    // 隐式实例化
	Func<double>(3.14);
	Func(1 > 2);         // false
	// Func<string>("hello");   // err, 未定义
}    
```

>---
#### 12.3. 模板特化   

全特化：允许对给定的模板实参集定制模板代码，显式特化可以在它的主模板的作用域中声明。

```cpp
template<typename T> 
struct A
{
	struct B {};      // 成员类
	template<class U> // 成员类模板
	struct C {};
};

template<> // 特化成员函数
struct A<int> { void f(int); };
void A<int>::f(int) { /* ... */ }   // 定义

template<> // 成员类的特化
struct A<char>::B { void f(); };
void A<char>::B::f() { /* ... */ }

template<> // 成员类模板的特化
template<class U>
struct A<char>::C { void f(); };
template<> template<class U>
void A<char>::C<U>::f() { /* ... */ }
```

部分特化：允许为给定的模板实参的类别定制类模板或变量模板。部分特化成员的模板形参列表和模板实参列表必须与部分特化的形参列表和实参列表相匹配。

```cpp
// 主模板
template<class... Ts>
struct S {
    template<Ts... Values>
    struct Inner {
        void print() { std::cout << "Generic case\n"; }
    };
};

// 特化外层模板 S，当第一个类型为 int 时
template<class... Rest>
struct S<int, Rest...> {
    // 内层模板的参数类型必须与外层模板的特化参数匹配
    template<int First, Rest... Values>
    struct Inner {
        void print() { std::cout << "First is int: " << First << "\n"; }
    };
};

// 进一步特化内层模板 Inner，当第一个非类型参数为 10 时
template<class... Rest>
template<Rest... Values>
struct S<int, Rest...>::Inner<10, Values...> {
    void print() { std::cout << "Specialization: First=10\n"; }
};

// 全特化
template<>
template<>
struct S<int, int, double>::Inner<1, 1, 3.14> {
    void print() { std::cout << "Full specialization: 1,1,3.14"; }
};

int main() {
    S<int, char>::Inner<10, 'a'> inner;
    inner.print();   // Specialization: First=10

    S<int, char>::Inner<42, 'a'> generic;
    generic.print(); // First is int: 42

    S<int, int, double>::Inner<1, 1, 3.14> full;
    full.print();    // Full specialization: 1,1,3.14
}
```

> *未知特化待决*

在模板定义内，某些名称被推导为属于某个未知特化：

```cpp
template<typename T>
struct Base {};
 
template<typename T>
struct Derived : Base<T> {
    void f() {
        // Derived<T> 指代当前实例化
        // 当前实例化没有 “unknown_type” 但有一个待决基类（Base<T>）
        // 因此 “unknown_type” 是未知特化的成员
        typename Derived<T>::unknown_type z;    // 待决名，假定 Base<T> 中有一个 unknown_type 成员
    }
};
 
template<>
struct Base<int> { // 它在此特化提供
    typedef int unknown_type; 
};
```

在模板（包括别名模版）的声明或定义中，非当前实例化成员且取决于某个模板形参的名称不会被认为是类型，除非使用 `typename` 声明，或它已经被设立为类型名（例如用 `typedef` 声明或通过用作基类名）。

```cpp
template<class T>
struct Iterator {};
int p = 1;

template<typename T>
void Func(const Iterator<T>& iter)
{
	// Iterator<T>::const_iterator 是待决名，
	typename Iterator<T>::const_iterator it = iter.begin();

	// 下列内容因为没有 “typename” 而会被解析成
	// 类型待决的成员变量 “const_iterator” 和 'p' 的乘法。
	Iterator<T>::const_iterator* p;

	typedef typename Iterator<T>::const_iterator iter_t;
	iter_t* p2; // iter_t 是待决名，但已知它是类型名
}

// template<>
// struct Iterator<int> {  
// 	struct const_iterator {};
// 	const_iterator begin() const& { return const_iterator{}; }
// };
int main() {
	Iterator<int> iter;  // 模板实例化失败，Iterator<int> 中没有 const_iterator 成员
	Func(iter);
}
```

>---
#### 12.4. 类模板推导指引

以尾随函数形式声明类类型模板的用户推导指引，不适用 `auto` 占位推断：

```cpp
[template <Params> [RequiresClause]]
[explicit] TemplateClassName( Params ) -> SampleTemplate [RequiresClause];
```

用户推导指引需要和类类型模板在同一语义作用域。

```cpp
template<class T>
struct S { S(T); };
S(char const*)->S<std::string>;
S s{ "hello" }; // 推导出 S<std::string>

template<class T>
struct container {
	container(T t) {}
	template<class Iter> container(Iter beg, Iter end);
};
template<class Iter>  // 额外的推导指引
container(Iter b, Iter e) -> container<typename std::iterator_traits<Iter>::value_type>;
vector v = { 1,2,3,4 };
container c{ v.begin(), v.end() };   // T is int
```

>---

#### 12.5. 概念与约束 

类模板、函数模板（lambda 模板）等可以与一项约束相关联。这类约束要求 `requires` 的具名集合称为概念。

```cpp
template <params>
concept ConceptName [Attr] = ConstraintExpr;
    ConstraintExpr = RequiresExpr | RequiresClause
```

概念可以作为模板类型形参声明、占位类型说明符或 `requires` 的复合要求。

```cpp
template<class T, class U>
concept Derived = std::is_base_of<U, T>::value;  // T 是 U 的基类
 
template<Derived<Base> T>
void f(T); // 隐式推断：T 被 Derived<T, Base> 约束
```

概念是 **要求** 的具名集合。要求由 `requires` 定义。不符合 `requires` 约束时返回 `false`。

```cpp
// requires 表达式
RequiresExpr = requires [( Params )] { 
    requirements = Sample | Type | Compound | Nested
        Sample   : expr;
        Type     : typename Identifier;
        Compound : { expr } [noexcept] [-> TypeConstraint | decltype(...)];
        Nested   : RequiresExpr | RequiresClause;
}; 
// requires 子句
RequiresClause = requires (ConceptName | ConstantExpr) [&& ConstraintExpr] [|| ConstraintExpr] // 联结

template<typename T> 
FuncDecl RequiresClause { ... };    // 函数声明符后面
template<typename T> RequiresClause // 或，模板形参列表后面
FuncDecl { ... };
```

Sample：简单要求断言表达式是有效的；

```cpp
template <class T, class U = T>
concept Swappable = requires(T && t, U && u) {
	swap(forward<T>(t), forward<U>(u));
	swap(forward<U>(u), forward<T>(t));
};
template <typename T>
concept Addable = requires(T a, T b) {
	a + b;
};
template <typename T>
void Add(T a, T b) requires Addable<T> {
	cout << a + b;
}  

Add(1,2);        // Add<int,int>
Add(1.2, 3.4);   // Add<double, double>
```

Type：类型要求断言 *标识符*（可有限定）类型是有效的：嵌套类型存在、类或别名模板特化有效。

```cpp
template<typename T> using Ref = T&;
template <typename T> class S {};
template<typename T>
concept C = requires(T t)
{
	typename T::inner; // 需要嵌套类型名
	typename S<T>;     // 需要类模板特化
	typename Ref<T>;   // 需要别名模板替换
};
// =================================================
template<class T, class U>
using CommonType = std::common_type_t<T, U>;
template<class T, class U>
concept Common = requires (T && t, U && u)
{
	typename CommonType<T, U>; // CommonType<T, U> 是合法的类型名
	{ CommonType<T, U>{std::forward<T>(t)} };
	{ CommonType<T, U>{std::forward<U>(u)} };
};
template<class T, class U> requires Common<T, U>
class Sample {};
template Sample<int, long>;       // CommonType is long
template Sample <float, double>;  // CommonType is double
template Sample <float, int>;     // err
template Sample <void*, int*>;    // CommonType is void*
```

Compound：复合要求断言表达式的属性。

```cpp
template<typename T>
concept C = requires(T x)
{
	// 表达式 *x 必须合法
	// 并且 类型 T::inner 必须存在
	// 并且 *x 的结果必须可以转换为 T::inner
	{ *x } -> std::convertible_to<typename T::inner>;
	// 表达式 x + 1 必须合法
	// 并且 std::same_as<decltype((x + 1)), int> 必须满足
	// 并且 x + 1 无异常抛出
	{ x + 1 } noexcept -> std::same_as<int>;
	// 表达式 x * 1 必须合法
	// 并且 它的结果必须可以转换为 T
	{ x * 1 } -> std::convertible_to<T>;
};
// =============================================
template <C T> T::inner F(T& t) { return *t; }
struct Sample {
	int value;
	typedef Sample* inner;		// T::inner
	inner operator* () { return this; }  // *T >> inner
	int operator +(int v) noexcept { return value += v; }	// x+1
	Sample operator *(int v) { return (this->value *= v, *this); }   // x*1
};
int main() {
	Sample s{};
	auto p = F<Sample>(s);
	(*p) + 10;  // value = 10
	(*p) * 2;	// value = 20
	cout << p->value;
}
```

Nested：嵌套要求可以根据本地形参指定其他约束。

```cpp
template<typename T>
concept Addable = requires (T a, T b) {
	a + b;
};
template<typename T, typename U>
concept Same = requires() {
	std::same_as<T, U>;
};
template<class T>
concept Semiregular = Addable<T> &&   // 联结
requires(T a, std::size_t n) {
	requires Same<T*, decltype(&a)>;       // 嵌套："Same<...> 求值为 true"
	{ a.~T() } noexcept;                   // 复合："a.~T()" 是不会抛出的合法表达式
	requires Same<T*, decltype(new T)>;    // 嵌套："Same<...> 求值为 true"
	requires Same<T*, decltype(new T[n])>; // 嵌套
	{ delete new T };                      // 复合
	{ delete new T[n] };                   // 复合
};
void F(Semiregular auto s) {  }

struct S {
	int v;
	S operator + (S s) { return S{ v + s.v }; }
	~S() noexcept {};
};
int main() {
	F(1);
	F(3.14);
	F(3.1415);
	F(S{ 10086 });
}
```

---
### 13. 异常处理

异常包括有编程逻辑错误和运行时错误。错误报告的管理方式是返回一个错误代码或设置一个全局变量，调用方选择性地检索该变量状态。例如 C `errno`。

`try-catch` 执行异常处理，未捕获时调用 `std::terminate` 终止程序。`throw` 可以抛出任何类型，但应引发派生自 `std::exception` 类型。对于每个可能引发或传播异常的函数，应提供三项异常保证之一：强保证、基本保证或 `noexcept` 保证。

`catch(...)` 程序块处理任意类型的异常。(MSVC 编译器) 当使用 `/EHa` 选项编译时，异常可包括 C 结构化异常和系统或应用程序生成的异步异常，例如内存保护、被零除和浮点冲突。

谨慎使用 `catch(...)`；除非 `catch` 明确知道特定异常。`catch(...)` 块一般用于在程序停止执行前记录错误和执行特殊的清理工作。在 `catch` 块中 `throw;` 表示 *rethrow*，异常对象是原始异常对象。

```cpp
template <typename T>
T Index(T arr[10], int i) throw(out_of_range) {  // 强保证
    if (i < 0 || i >= 10)
        throw std::out_of_range("index is out of range");   // 抛出值
    return arr[i];
}
int CatchIndex(int arr[]) noexcept {     // 无异常抛出
    try {
        auto v = Index(arr, -1);
    } catch (out_of_range& e) {          // 引用捕获
        cout << e.what() << endl;
        return -1;
    } catch (exception& e) {
        return -10;
    } catch (...) {   // 任意类型
        return -100;
    }
}
int main() {
    int arr[10] = { 0,1,2,3,4,5,6,7,8,9 };
    return CatchIndex(arr);
}
```

>---
#### 13.1. 未处理的异常

未匹配或未处理的异常，将调用预定义 `std::terminate()`（默认操作是 `abort()`）。可以在程序的任何点调用 `set_terminate()`。`terminate` 总是调用最后一次指定给 `set_terminate` 参数。

最好在 `set_terminate` 的 `terminate_handler` 处理程序中调用 `exit` 来终止程序或当前线程，否则返回调用方调用 `abort`。

```c++
using namespace std;
void term_func() {
   cout << "term_func was called by terminate." << endl;
   exit( -1 );
}
int main() {
   try {
      set_terminate( term_func );
      throw "Out of memory!"; // No catch handler for this exception
   }
   catch( int ) {
      cout << "Integer exception raised." << endl;
   }
}  // -1
```

>---
#### 13.2. 函数 try 块

```cpp
function try {
    // 函数体
} catch ( /*...*/ ) {
    // 处理异常
} 
constructor try [: init-list ] 
{ ... } catch { ... } 
```

函数 `try` 块是一类特殊的函数体。对于构造函数或析构函数 `try` 块，始终默认重新抛出异常；在这类函数 `try` 块中操作该对象的非静态成员或基类会导致未定义行为。

```cpp
struct Sample {
	int value;
	Sample(int value) try : value{ value } {
		throw this->value;
	} catch (exception& e) {
		cout << e.what();
	} catch (...) {
		cout << "catch unknown exception\n";
	}
};

int main(int, const char* argv[]) try {
	Sample s{ 10086 };
} catch (int v) {
	return v;  // 10086
}
// catch unknown exception
```

>---
#### 13.3. 异常规范与 noexcept

异常规范指示可由函数传播的异常类型的意图。`noexcept` 指定可以脱离函数的潜在异常集是否为空；`throw()` 是 `noexcept(true)` 的别名。

| 异常规范                                    | 含义                                                                                                                                                                                               |
| :------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `noexcept`<br>`noexcept(true)`<br>`throw()` | 函数不会引发异常，也不允许在其范围外传播异常。`noexcept` 和 `noexcept(true)` 是等效的。此规范声明的函数引发异常时，将直接调用 `std::terinate` 终止程序。并且不会保证调用任何范围内对象的析构函数 |
| `noexcept(false)`<br>`throw(...)`<br>无规范 | 函数可以引发任何类型的异常。                                                                                                                                                                       |
| `throw(type)`                               | C++14 之前表示函数可以引发指定 `type` 异常，之后编译器统一解释为 `noexcept(false)`                                                                                                               |

```c++
#include <type_traits>
template <typename T>
T copy_object(const T& obj) noexcept(std::is_pod<T>) // 当 T 是 POD 时，函数不会抛出异常
{
   // ...
}
```

> **noexcept(expr)**

`noexcept(expr)` 运算符进行编译时检查，在表达式声明不会抛出任何异常时返回 `true`。`expr` 是不求值操作数。

```cpp
void may_throw();
void no_throw() noexcept;
auto lmay_throw = []{};
auto lno_throw = []() noexcept {};

std::cout << std::boolalpha
    << "Will may_throw() throw an exception?" << !noexcept(may_throw()) << '\n'        // true 
    << "no_throw()?" << !noexcept(no_throw()) << '\n'          // false
    << "lmay_throw()?" << !noexcept(lmay_throw()) << '\n'      // true
    << "lno_throw()?" << !noexcept(lno_throw()) << '\n';       // false
std::cout << noexcept(1 / 0);   // true;
```

若 `expr` 具有类类型或它的数组类型，要求析构函数可访问且未 `delete`。

```cpp
struct T {
	~T() {} // 复制构造函数不会抛出异常
};
struct U {
	std::vector<int> v;
	~U() {} // 复制构造函数可能会抛出异常
};
struct V { std::vector<int> v; };
int main()
{
	T t; U u; V v;

	std::cout << std::boolalpha
		<< "~T() 可能会抛出异常吗？" << !noexcept(std::declval<T>().~T()) << '\n'      // false
		<< "T(T 右值) 可能会抛出异常吗？" << !noexcept(T(std::declval<T>())) << '\n'   // false
		<< "T(T 左值) 可能会抛出异常吗？" << !noexcept(T(t)) << '\n'                   // false
		<< "U(U 右值) 可能会抛出异常吗？" << !noexcept(U(std::declval<U>())) << '\n'   // true 
		<< "U(U 左值) 可能会抛出异常吗？" << !noexcept(U(u)) << '\n'                   // true
		<< "V(V 右值) 可能会抛出异常吗？" << !noexcept(V(std::declval<V>())) << '\n'   // false
		<< "V(V 左值) 可能会抛出异常吗？" << !noexcept(V(v)) << '\n';                  // true
}
```

---
### 14. 预处理指令

| Category       | Command                                                                  |
| :------------- | :----------------------------------------------------------------------- |
| 空指令         | `#`                                                                      |
| 条件包含       | `defined`,`__has_include`,`__has_cpp_attribute`                          |
| 条件控制       | `#if`,`#elif`,`#else`,`#endif`,`#ifdef`,`#ifndef`,`#elifdef`,`#elifndef` |
| 宏相关         | `#define`,`#undef`                                                       |
| 模块导出       | `export`,`module`                                                        |
| 源文件包含     | `#include`,`import`                                                      |
| 行信息         | `#line LINENUM [FILENAME]`                                               |
| 诊断           | `#error`,`#warning`                                                      |
| 编译器行为控制 | `#pragma`,`_Pragma`                                                      |

>---
#### 14.1. 宏定义

`#define` 创建宏定义或宏函数。编译器预处理阶段将源文件中每个宏标识符（类常量宏和类函数宏）的内容替换为对应的标记字符串。

```c++
#define identifier 									// 条件编译宏
#define identifier token-string						// 类常量宏
#define identifier(id0?, id1?, ...?) token-string?  // 类函数宏

#undef identifier    // 取消宏定义
```
```c++
#define DEBUG
#define MAX_BUFSZ  512
#define getRandom(min, max) \
    ((rand()%(int)(((max) + 1)-(min)))+ (min))

#if defined(DEBUG)    // 条件编译测试
    void foo(int){};
#else 
    void foo(char){};
#endif
```

> *字符串化运算符 `#`*

在仿函数宏中，`#` 字符串化 *替换列表* 对应的标识符。例如 `x` 是一个宏形参，`#x` 字符串化为 `"x"`。

```c
#include <stdio.h>
#define stringer( x ) printf_s( #x "\n" )
int main() {
   stringer( In quotes in the printf function call );
   stringer( "In quotes when printed to the screen" );
   stringer( "This: \"  prints an escaped double quote" );
   // 替换为
   printf_s( "In quotes in the printf function call" "\n" );
   printf_s( "\"In quotes when printed to the screen\"" "\n" );
   printf_s( "\"This: \\\" prints an escaped double quote\"" "\n" );
}
```

> *连接运算符 `##`*

仿函数宏中， `##` 用于连接两个 *token*。一些编译器允许 `##` 出现在逗号后和 `__VA_ARGS__` 前的扩展，在此情况下 `##` 在 `__VA_ARGS__` 非空时无效，但在 `__VA_ARGS__` 为空时移除逗号：这使得可以定义如 `fprintf (stderr, format, ##__VA_ARGS__)` 的宏。

```c
#include <stdio.h>
#define Printf(format, ...) fprintf(stderr, format, ## __VA_ARGS__)
#define XNAME(n) x ## n   // 表示将 x 与 n 组合成一个记号
#define PRINT_XN(n) printf("x" #n " = %d\n", x ## n);

int main(void)
{
    Printf("%d %d\n", 1, 2);   // fprintf(stderr, "%d %d\n", 1, 2);
    Printf("Hello World\n");   // fprintf(stderr, "Hello World\n");

    int XNAME(1) = 14; 	// int x1 = 14;
    int XNAME(2) = 20; 	// int x2 = 20;
    PRINT_XN(1); // printf("x1 = %d\n", x1);
    PRINT_XN(2); // printf("x2 = %d\n", x2);
    return 0;
}
```

>---
#### 14.2. 预定义宏

标准预定义标识符 `__func__`：

```c++
void example() {
    printf("%s\n", __func__);   // "example"
} 
```

编译器支持 ISO C99、C11、C17 和 ISO C++17 标准指定的以下预定义宏：

| macro                              | description                                                                                                                                                                                                 |
| :--------------------------------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `__cplusplus`                      | 当翻译单元编译为 C++ 时，定义为整数文本值。其他情况下则不定义。                                                                                                                                             |
| `__DATE__`                         | 当前源文件的编译日期。日期是 Mmm dd yyyy 格式的恒定长度字符串文本。月份名 Mmm 与 C 运行时库 (CRT) `asctime` 函数生成的缩写月份名相同。如果值小于 10，则日期 dd 的第一个字符为空格。任何情况下都会定义此宏。 |
| `__FILE__`                         | 当前源文件的名称。`__FILE__` 展开为字符型字符串文本。                                                                                                                                                       |
| `__LINE__`                         | 定义为当前源文件中的整数行号。可使用 `#line` 指令来更改此宏的值。                                                                                                                                           |
| `__STDC__`                         | 仅在编译为 C，它定义为 1。其他情况下则不定义。                                                                                                                                                              |
| `__STDC_HOSTED__`                  | 如果实现是托管实现并且支持整个必需的标准库，则定义为 1。其他情况下则定义为 0。                                                                                                                              |
| `__STDC_NO_ATOMICS__`              | 如果实现不支持可选的标准原子，则定义为 1。                                                                                                                                                                  |
| `__STDC_NO_COMPLEX__`              | 如果实现不支持可选的标准复数，则定义为 1。                                                                                                                                                                  |
| `__STDC_NO_THREADS__`              | 如果实现不支持可选的标准线程，则定义为 1。                                                                                                                                                                  |
| `__STDC_NO_VLA__`                  | 如果实现不支持可选的可变长度数组，则定义为 1。                                                                                                                                                              |
| `__STDC_VERSION__`                 | 标准 C 的 version。                                                                                                                                                                                         |
| `__STDCPP_DEFAULT_NEW_ALIGNMENT__` | 宏会扩展为 `size_t` 字面量，该字面量的对齐值由对非对齐感知的 `operator new` 的调用所保证。较大的对齐值传递到对齐感知重载，例如 `operator new(std::size_t, std::align_val_t)`。                              |
| `__STDCPP_THREADS__`               | 当且仅当程序可以有多个执行线程并编译为 C++ 时，定义为 1。其他情况下则不定义。                                                                                                                               |
| `__TIME__`                         | 预处理翻译单元的翻译时间。时间是 hh:mm:ss 格式的字符型字符串文本，与 CRT `asctime` 函数返回的时间相同。任何情况下都会定义此宏。                                                                             |

<!-- 
---
### 15. 附录 -->
---