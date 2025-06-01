## C++ 程序构建基块

### 1. 基本概念

#### 1.1. 类型系统

对象、引用、函数和表达式具有类型的性质。C++ 类型系统分为基础类型和复合类型。变量和成员函数声明可以包含 CV 限定：`const`、`volatile`、`const volatile`。

*基础类型* 包括有布尔、字符、整数、浮点数等类型；*复合类型* 包括有引用（左值 `&`、右值 `&&`）、指针、成员指针、数组、函数、枚举、类类型（`class`，`struct`，`union`）等类型。除函数类型、引用、`void` 之外的类型是对象类型。

*标量类型* 包括有数值、指针类、枚举、`nullptr_t` 等。*隐式生存期类型* 包括有标量类型、隐式生存期类类型、数组等。*可平凡复制类型* 包括有标量类型和可平凡复制类类型（数组）。*标准布局类型* 包括有标量类型和标准布局类类型（数组）。

> *隐式生存期类*

隐式生存期类满足：它是没有用户定义析构函数的聚合体、或至少具有一个平凡合格的构造函数和一个平凡未弃置析构函数的类类型。聚合体可以是数组类型或满足以下条件的类类型：
- 没有用户声明或继承的构造函数；
- 没有私有或受保护的直接非静态数据成员；
- 没有虚基类或私有受保护的直接基类；
- 没有虚成员函数；
- 没有默认成员初始化器。


> *可平凡复制类*

可平凡复制类满足：它至少包含有一个合格的复制构造函数、移动构造函数、复制赋值运算符、移动赋值运算符（全部要求平凡），和一个非弃置（`delete`）的平凡析构函数。
- 合格定义为该构造函数没有被弃置；它的所有关联约束（如果有）被满足；在所有满足关联约束的移动构造函数中，没有比其他都更受约束者。
- 平凡定义为类 `T` 的构造函数非用户提供，它没有虚成员函数和虚基类；它的每个直接基类的对应构造函数也是平凡的；类 `T` 的所有非静态的类类型成员的构造函数也是平凡的。

> *标准布局类*

标准布局类类型满足：没有非标准布局类类型或引用的非静态数据成员或数组成员；所有非静态数据成员具有相同的可访问性；没有非标准布局的基类；继承层级中仅有一个类具有非静态数据成员。

一个对象具有大小（`sizeof`）、对齐要求（`alignof`）、存储期、生存期、名称、类型、值、地址等信息。声明或继承至少一个虚函数的类类型对象是多态对象，这类对象包含一些额外的信息用于进行虚函数的调用。

> *结构化类型*

结构化类型可以具有 CV 限定，包括有数值类型、指向对象或函数的左值引用或指针、成员指针类型、枚举、`nullptr_t`、无闭包捕获的 Lambda、具有全部公开非 `mutable` 非静态数据成员（包括基类）且成员类型为结构化类型或数组的字面类类型。

>---
#### 1.2. 布局与 POD

布局是指类类型对象的成员在内存中的排列方式。如果当类或结构包含某些语言功能（如虚拟基类、虚拟函数、具有不同访问控制的成员）时，编译器可以自由选择布局，因此存储对象可能不会使用连续的内存区域。

例如，某个类具有虚拟函数，则该类的所有实例可能会共享单个虚拟函数表。由于布局未定义，无法将这类对象传递到其他语言（例如 C）编写的程序，因为它们可能是非连续的。

POD（简单旧数据）类型同时为平凡和标准布局，它的内存布局是连续的，可以对这些类型执行逐字节复制和二进制 I/O。标量类型（例如 `int`）也是 POD 类型。作为类的 POD 类型只能具有作为非静态数据成员的 POD 类型。

```cpp
#include <type_traits>
#include <iostream>
using namespace std;

struct B
{
protected:
	virtual void Foo() {}
};
// 非平凡非标准布局
struct A : B   
{
	int a;
	int b;
	void Foo() override {}
};
// 平凡非标准布局
struct C
{
	int a;
private:
	int b;   // 具有不同访问的非静态数据成员
};
// 标准布局非平凡
struct D
{
	int a;
	int b;
	D() {} // 具有用户定义否早函数
};
struct POD
{
	int a;
	int b;
};

int main()
{
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
#### 1.3. **对齐**

对齐表示类型不同对象所能分配连续相邻地址之间的字节数。`alignof(Type)` 获取类型的对齐要求，`alignas(Type | expr)` 设置类类型、类的非位域数据成员、变量（非函数参数）的对齐要求，对齐要求 = 0,1,2,4,8,16,... 不低于类型的自然对齐。

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

预处理指令 `#pragma pack` 设置当前翻译单元的编译器默认对齐方式：

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
#### 1.4. **作用域和存储期**

作用域（范围）包含有 *全局范围*、*文件范围*、*命名空间范围*、*类范围和枚举范围*、*函数域*、*块范围* 等。名称可以具有外部链接（`extern`）、内部链接（`static`）、模块链接或无链接（块范围）。全局范围成员可以通过 `::` 显式访问。顶层命名空间作用域（非模块声明）中为 `const` 且非 `extern` 的名称在具有内部链接。

对象的存储期定义包含该对象的存储的最短潜在生存期，包括有：*静态存储期*（`static`,`extern`）、*线程存储期*（`thread_local`）、*自动存储期*（局部非静态变量）、*动态存储期*（`new` 或显式内存分配的对象，异常对象）。应避免出现访问悬挂引用或悬挂指针（关联对象无效或内存位置被释放）。

```cpp
// 全局范围
#include "some.h"  // 嵌入头文件
import std.compat; // 导入模块

class Sample {};    // 类声明，外部链接，全局范围
static void Foo();  // 函数声明，内部链接，文件范围
const int V = 10086;  // 常量，内部链接，文件范围
namespace S {
	static const int* Foo() {  // 内部链接，命名空间范围
		int v = 10;				// 局部变量，动态存储期
		static int sv = v;		// 局部变量，静态存储期
		int* ptr = static_cast<int*>(std::malloc(sizeof(int) * sv));   // 动态存储期
		::Foo();   // :: 作用域解析
	}   
}
```

>---

#### 1.5. 声明与定义

程序的实体（*Entity*）包括值、对象、引用、结构化绑定、函数、枚举器、类型、类成员、位字段、模板、模板专用化、命名空间或包等。标识符名称在所属名称空间遵循单一定义原则，函数重载支持同名但签名不同。

变量是通过声明非静态数据成员或对象以外的引用来引入的。变量的名称（如果有）表示引用或对象。局部实体是具有自动存储期的变量、结构化绑定对应的变量，或 *this 对象。

声明可以（重新）将一个或多个名称和 / 或实体引入翻译单元，并指定这些名称的解释和语义属性。如果可以从实体或 *typedef-name* `X` 中获取另一个 `X` 声明，则该实体或 *typedef-name* `X` 的声明是 `X` 的重新声明。

在对象的定义中，该对象的类型不应该是不完整类型、抽象类类型或其数组（可能是多维的）。

```c++
// 定义声明
int a;                         // defines a
extern const int c = 1;        // defines c
int f(int x) { return x + a; } // defines f and defines x
struct S
{
    int a;
    int b;
}; // defines S, S::a, and S::b
struct X
{                 // defines X
    int x;        // defines non-static data member x
    static int y; // declares static data member y
    X() : x(0) {} // defines a constructor of X
};
int X::y = 1; // defines X::y
enum
{
    up,
    down
}; // defines up and down
namespace N
{
    int d;
} // defines N and N::d
namespace N1 = N; // defines N1
X anX;            // defines anX

// 仅声明
extern int a;       // declares a
extern const int c; // declares c
int f(int);         // declares f
struct S;           // declares S
typedef int Int;    // declares Int
extern X anotherX;  // declares anotherX
using N::d;         // declares d
```

在某些情况下，编译器会为一些定义声明隐式定义默认构造函数、赋值构造函数、移动构造函数、复制赋值运算符、移动赋值运算符或终结器。

```c++
struct C {
    std::string s; 
};
// 实现将隐式定义
struct C
{
    std::string s;
    C() : s() {}
    C(const C &x) : s(x.s) {}
    C(C &&x) : s(static_cast<std::string &&>(x.s)) {}
    // : s(std::move(x.s)) { }
    C &operator=(const C &x){
        s = x.s;
        return *this;
    }
    C &operator=(C &&x)
    {
        s = static_cast<std::string &&>(x.s);
        return *this;
    }
    // { s = std::move(x.s); return *this; }
    ~C() {}
};
```

>---
#### 1.6. 对象资源管理

C++ 没有自动回收垃圾，程序负责将所有已获取的资源返回到操作系统。未被释放未使用的资源是资源泄露（*leak*）。现代 C++ 通过声明堆栈上的对象，尽可能避免使用堆内存。

当对象初始化时，它会获取它拥有的资源，并且该对象负责在其析构函数中释放资源。在堆栈上声明拥有资源的对象本身。对象拥有资源的原则也称为 “资源获取即初始化” (RAII)。当拥有资源的堆栈对象超出范围时，会自动调用其析构函数。C++ 中的垃圾回收与对象生存期密切相关，是确定性的。资源始终在程序中的已知点释放。

```c++
struct buffer {
    friend void print(buffer&);
private:
    char* buf;
    size_t size;
public:
    buffer(size_t size);
    int writeString(string str);
    ~buffer();
};
buffer::buffer(size_t size) :size{ size } {
    buf = new char[size];
}
buffer::~buffer() {
    cout << "delete buffer" << endl;
    delete[] buf;
}
int buffer::writeString(string str) {
    int n = 0;
    try {
        for (char c : str) {
            if (n >= size - 1)
                break;
            buf[n++] = c;
        }
    }
    catch (exception e) {
        return -1;
    }
    buf[n] = '\0';
    return n;
}
void print(buffer& buf) {
    printf("%s\n", buf.buf);
}

void WriteToBuffer(std::string& str) {
    buffer b(512);
    int	n = b.writeString(str);
    if (n > 0)
        print(b);
}

int main()
{
    string s{ "Hello World" };
    WriteToBuffer(s);
}
// Hello World
// delete buffer
```

C++ 的设计可确保对象在超出范围时被销毁。也就是说，它们在块被退出时以与构造相反的顺序被摧毁。销毁对象时，将按特定顺序销毁其基项和成员。

可以使用智能指针处理对象所需内存资源的分配和删除。使用智能指针进行内存分配，可以消除内存泄漏的可能性。此模型适用于其他资源，例如文件句柄或套接字。


```c++
#include <memory>
class widget
{
private:
    std::unique_ptr<int[]> data;
public:
    widget(const int size) { data = std::make_unique<int[]>(size); }
    void do_something() {}
};

void usingWidget() {
    widget w(1000000);  // lifetime automatically tied to enclosing scope
                        // constructs w, including the w.data gadget member
    w.do_something();
} // automatic destruction and deallocation for w and w.data
```

>---
#### 1.7. 多线程与数据竞争

执行线程是程序中的控制流，从某个特定顶层函数调用（`thread`,`async`,`jthread` 等）开始。不同的执行线程始终可以同时访问（读和写）不同的内存位置，在不同线程上可能会出现读写冲突或写写冲突的数据竞争，除非是这两个冲突的求值是原子操作或顺序执行的。

对同一个除 `std::vector<bool>` 之外的所有标准库容器的不同元素对象的并行修改不会造成数据竞争。



>---
#### 1.8. 程序启动与终止

```c++
int main(); // 或
int main(int argc, char *argv[]);    // argv[0] 为程序名称
```

程序包含入口函数 `main`，动态链接库和静态库无需入口函数。


```c++
#include <iostream>

using namespace std;
int main( int argc,       // Number of strings in array argv
          char *argv[])   // Array of command-line argument strings
{
    // Display each command-line argument.
    cout << "\nCommand-line arguments:\n";
    for( int count = 0; count < argc; count++ )
         cout << "  argv[" << count << "]   "
                << argv[count] << "\n";
}
```

可以通过以下方式退出程序：

- 调用 `exit` 函数，终止程序并执行清理（如调用全局对象析构函数）。
- 调用 `abort` 函数。立即终止程序，跳过 `atexit` 机制。
- 从 `main` 执行 `return` 语句。

```c++
#include <cstdlib>

int main() {
    if (cond1){
        exit(EXIT_SUCCESS);  // or EXIT_FAILURE 
    }else if (cond2) {
        abort();
    }else 
        return 0;
}
```

`atexit` 为在程序终止之前执行的操作。`main` 以返回值作为参数调用 `exit` 并清理自动变量。直接调用 `exit` 时（不会销毁自动变量），销毁与当前线程关联的线程对象，然后（在调用指定给 `atexit` 的函数（如果有）之后销毁静态对象）。 

```c++
#include <cstdio>
class ShowData {
public:
    ShowData(const char* szDev) {
        errno_t err = fopen_s(&OutputDev, szDev, "w");
    }

    ~ShowData() { fclose(OutputDev); }
    void Disp(const char* szData) {
        fputs(szData, OutputDev);
    }
private:
    FILE* OutputDev;
};
// 静态
ShowData sd1 = "CON";
ShowData sd2 = "CON";

int main() {
    sd1.Disp("hello sd1\n");
    sd2.Disp("hello sd2\n");  // 相反顺序销毁并调用 ~ShowData()
}
// hello sd2  
// hello sd1
```

>---
#### 1.9. 词法元素

> **关键字**

```cpp
bool,false,true                 // 布尔类型
nullptr                         // 空指针常量
void                            // 无类型
float,double,long double        // 浮点数
signed,unsigned                 // 整数符号限定
char,short,int,long,long long   // 整数
char,char8_t,char16_t,char32_t,wchar_t   // 字符
namespace                       // 命名空间
enum                            // 枚举
class,struct,union              // 类类型
public,protected,private        // 访问性
virtual,final,override          // 继承与封装
this                            // 非成员函数隐式对象, 或显式对象声明
auto                // 类型推断占位符
decltype            // 生成表达式类型
const               // 常量限定
volatile            // 易变限定
register            // 弃用, 寄存器存储
extern              // 外部链接, 语言链接
static              // 内部链接, 成员静态声明
thread_local        // 线程存储声明
friend              // 友元声明
inline              // 内联声明
mutable             // 非静态数据 const 函数可修改修饰
operator            // 运算符声明
typedef,using       // 别名声明
using               // 导入声明            
module,export,import                        // 模块
template,tempname,class,concept,requires    // 模板
alignas,alignof     // 对齐
new,delete          // 内存分配和释放
default,delete      // 函数预置和弃置
sizeof              // 类型大小
typeid              // 类型信息
consteval           // 立即函数
constexpr           // 编译时计算
constinit           // 编译时变量初始化
const_cast          // 去 CV 限定转换
dynamic_cast        // 多态类型转换
reinterpret_cast    // 无视类型二进制转换
static_cast         // 编译时类型转换
explicit            // 显式转换函数
do-while,while,for              // 迭代语句
if-else,switch-case[-default]   // 条件分支语句
break,continue,return,goto      // 跳转语句
co_await,co_return,co_yield     // 协程
try,catch,throw,noexcept        // 异常处理和异常规范
static_assert                   // 静态断言
asm                             // 汇编
```

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
<%           {   
%>           }    
<:           [    
:>           ]    
%:           #    
%:%:         ##   

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
#### 1.10. 模块

模块用于在翻译单元间共享声明和定义，它们可以在某些地方替代头文件的使用。编译器可以比头文件更快地处理模块文件。

```cpp
// helloworld.ixx
export module helloworld;  // 主模块分区声明
import <iostream>;         // 导入
export void Hello(){
    std::cout << "Hello C++23!" << endl;
}
// main.cpp
import helloworld;         // 导入模块
int main(){
    ::Hello();
}
```

一个完整的模块包含有一个主模块接口单元（*.ixx）和若干的模块分区接口单元，以及它们的模块实现单元。模块可以导出具有命名空间范围的顶级类型声明、模块分区、头文件或其他模块的内容。不会导出模块中出现或导入的任何宏定义。

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

| 语法                                                   | 描述                            |
| :----------------------------------------------------- | :------------------------------ |
| `export module <Name> [:<Partition>] [<Property>]`     | 声明模块接口单元                |
| `module <Name> [:<Partition>] [<Property>] `           | 声明模块实现单元                |
| `module;`                                              | 开始一个全局模块片段            |
| `module : private;`                                    | 开始一个私有模块片段            |
| `export <declaration>`,`export { <declaration-list> }` | 导出声明                        |
| `[export] import <Name>`                               | 导入（再 `export`）一个模块单元 |
| `[export] import <: Partition>`                        | 导入（再 `export`）一个模块分区 |
| `[export] import <Head>`                               | 导入（再 `export`）一个头文件   |

```cpp
// （每行表示一个单独的翻译单元）
export module A;    // 为具名模块 'A' 声明主模块接口单元
module A;           // 为具名模块 'A' 声明一个模块实现单元
module A;           // 为具名模块 'A' 声明另一个模块实现单元
export module A.B;  // 为具名模块 'A.B' 声明主模块接口单元
export module A:P1; // 为具名模块 'A' 声明一个模块分区接口单元 'P1'
module A:P1;        // 为具名模块 'A:P1' 声明一个模块分区实现单元
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

C++23 标准库引入了两个命名模块：`std` 和 `std.compat`；

- `import std` 导出 `std` 中定义的声明和名称，它还会导出 C 包装器标头的内容，例如 `<cstdio>` 和 `<cstdlib>`，提供类似 `std::printf()` 函数的内容。不会导出全局命名空间（如 `::printf()`）中定义的 C 函数。

- `import std.compat` 导出 `std` 中的所有内容，并添加 C 运行时全局命名空间，例如 `::printf`、`::fopen`、`::size_t`、`::strlen` 等。

```c++
import std.compat;
int main() {
    printf("Hello World! %s", "CPP!");
}
```

>---



---
### 2. 预处理指令

| Category       | Command                                                                  |
| :------------- | :----------------------------------------------------------------------- |
| 空指令         | `#`                                                                      |
| 条件包含       | `defined`,`__has_include`,`__has_cpp_attribute`                          |
| 条件控制       | `#if`,`#elif`,`#else`,`#endif`,`#ifdef`,`#ifndef`,`#elifdef`,`#elifndef` |
| 宏相关         | `#define`,`#undef`                                                       |
| 模块导出       | `export`,`module`                                                        |
| 源文件包含     | `#include`,`import`                                                      |
| 行信息         | `#line`                                                                  |
| 诊断           | `#error`,`#warning`                                                      |
| 编译器行为控制 | `#pragma`,`_Pragma`                                                      |

>---
#### 2.1. 宏定义

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
#### 2.2. 条件控制指令

```c++
#if constant-expr
#elif constant-expr
#ifdef identifier      // #if defined identifier
#ifndef identifier     // #if !defined identifier 
#elifdef identifier    // #elif defined identifier 
#elifndef identifier   // #elif !defined identifier 
#else
#endif
// 条件包含
defined ( identifier )               // 条件编译符号测试
__has_include ( heeder-name )        // 源文件包含测试
__has_cpp_attribute ( pre-tokens )   // 属性测试
```

`#ifdef`、`#ifndef`、`#elifdef` 和 `#elifndef` 指令以及 `defined` 条件包含应将 `__has_include` 和 `__has_cpp_attribute` 视为定义宏的名称。这些指令检查控制表达式结果是否为非零来选择条件编译的代码块。

```c++
#if __has_include(<optional>)
# include <optional>
#if __cpp_lib_optional >= 201603
# define have_optional 1
#endif
#elif __has_include(<experimental/optional>)
# include <experimental/optional>
#if __cpp_lib_experimental_optional >= 201411
# define have_optional 1
# define experimental_optional 1
#endif
#endif
#ifndef have_optional
# define have_optional 0
#endif

#if __has_cpp_attribute(acme::deprecated)
# define ATTR_DEPRECATED(msg) [[acme::deprecated(msg)]]
#else
# define ATTR_DEPRECATED(msg) [[deprecated(msg)]]
#endif
ATTR_DEPRECATED("This function is deprecated") void anvil();
```


>---
#### 2.3. 源文件包含

```c++
#include <header-name> 
#include "header-name" 
#include "source-name" 
```

`#include` 导入命名源文件，并将该位置替换为源文件的全部内容。

```c++
#if VERSION == 1
# define INCFILE "vers1.h"
#elif VERSION == 2
# define INCFILE "vers2.h" // and so on
#else
# define INCFILE "versN.h"
#endif
#include INCFILE
```

>---
#### 2.4. 诊断指令

```c++
#if !defined(__cplusplus)
    #error C++ compiler required.
#else 
#if _MSVC_STL_VERSION <= 140 
    #warning STL version should be 140 at least.
#endif
```

>---
#### 2.5. 行指令

```c++
#line digit-sequence ["filename"]
```

`#line` 指令更改源行号和文件名。编译器使用行号和文件名来确定预定义宏 `__FILE__` 和 `__LINE__` 的值。

```c++
#include <iostream>   // file：main.cpp
int main()
{
    printf( "This code is on line %d, in file %s\n", __LINE__, __FILE__ );
#line 10
    printf( "This code is on line %d, in file %s\n", __LINE__, __FILE__ );
#line 20 "hello.cpp"
    printf( "This code is on line %d, in file %s\n", __LINE__, __FILE__ );
    printf( "This code is on line %d, in file %s\n", __LINE__, __FILE__ );
}
/*
This code is on line 4, in file main.cpp
This code is on line 10, in file main.cpp
This code is on line 20, in file hello.cpp
This code is on line 21, in file hello.cpp
*/
```

>---
#### 2.6. 预定义宏

标准预定义标识符 `__func__`，封闭函数（用作 `char` 的函数局部 `static const` 数组）的未限定、未修饰名称。

```c++
void example()
{
    printf("%s\n", __func__);
} // prints "example"
```

编译器支持 ISO C99、C11、C17 和 ISO C++17 标准指定的以下预定义宏：

- `__cplusplus`：当翻译单元编译为 C++ 时，定义为整数文本值。其他情况下则不定义。
- `__DATE__`：当前源文件的编译日期。日期是 Mmm dd yyyy 格式的恒定长度字符串文本。月份名 Mmm 与 C 运行时库 (CRT) `asctime` 函数生成的缩写月份名相同。如果值小于 10，则日期 dd 的第一个字符为空格。任何情况下都会定义此宏。
- `__FILE__`：当前源文件的名称。`__FILE__` 展开为字符型字符串文本。
- `__LINE__`：定义为当前源文件中的整数行号。可使用 `#line` 指令来更改此宏的值。
- `__STDC__`：仅在编译为 C，它定义为 1。其他情况下则不定义。
`__STDC_HOSTED__`：如果实现是托管实现并且支持整个必需的标准库，则定义为 1。其他情况下则定义为 0。
- `__STDC_NO_ATOMICS__` 如果实现不支持可选的标准原子，则定义为 1。
- `__STDC_NO_COMPLEX__` 如果实现不支持可选的标准复数，则定义为 1。
- `__STDC_NO_THREADS__` 如果实现不支持可选的标准线程，则定义为 1。
- `__STDC_NO_VLA__` 如果实现不支持可选的可变长度数组，则定义为 1。
- `__STDC_VERSION__` 标准 C 的 version。 
- `__STDCPP_DEFAULT_NEW_ALIGNMENT__` 宏会扩展为 `size_t` 字面量，该字面量的对齐值由对非对齐感知的 `operator new` 的调用所保证。较大的对齐值传递到对齐感知重载，例如 `operator new(std::size_t, std::align_val_t)`。
- `__STDCPP_THREADS__` 当且仅当程序可以有多个执行线程并编译为 C++ 时，定义为 1。其他情况下则不定义。
- `__TIME__` 预处理翻译单元的翻译时间。时间是 hh:mm:ss 格式的字符型字符串文本，与 CRT `asctime` 函数返回的时间相同。任何情况下都会定义此宏。


---
### 3. 命名空间

命名空间是一个声明性区域，通过 `using` 指令导入命名空间（`using namespace std`）或导入类型成员（`using std::string`）；头文件中应始终使用完全限定的命名空间名称。

`main` 入口函数在全局命名空间中声明。所有 C++ 标准库类型和函数都在 `std` 命名空间或内部嵌套命名空间中声明。

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
	// using 声明
	using CustomData::ObjectManager;  // 导入 ObjectManager
	ObjectManager mgr;
	mgr.DoSomething();
	// using 指令
	using namespace CustomData;       // 导入 CustomData
	ObjectManager mgr;
	mgr.DoSomething();
	Func(mgr);
}
```

> *命名空间别名*

```c++
namespace a_very_long_namespace_name { class Foo {}; }
namespace AVLNN = a_very_long_namespace_name;
void Bar(AVLNN::Foo foo){ }
```

> *匿名命名空间*

匿名命名空间中的成员对外不可见，仅对其翻译单元范围可见，对外部文件不可见（相当于内部链接）。

```c++
// flieA.cpp
namespace Parent {
    namespace {
        void Myfunc() {
            std::cout << "Call Myfunc" << std::endl;
        }
    }
    void CallMyfunc() {
        Myfunc();  // 当前翻译单元可见
    }
}
// fileB.cpp
namespace Parent {
    void CallMyfuncOther() {
        // Parent::Myfunc 在当前翻译单元不可见
    }
}
int main() {
    Parent::CallMyfunc(); 
}
```

>---
#### 3.1. 内联命名空间

内联命名空间（`inline namesapce`）的成员被视为父级空间的成员。可以在内联命名空间中声明模板，然后在父命名空间中声明专用化：

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
    template<> class C<int> {};  // 专门化
    void new_ns::Func() {		 // 提供实现
        cout << 1;
    }
}
int main() {
    Parent::Func();   // 相当于 Parent::new_ns::Func();
}
```

可以将内联命名空间用作版本控制机制，以管理对库的公共接口的更改。

```c++
namespace Custom
{
    namespace v_10   
    {
        template <typename T>
        class Funcs  // 非内联；无法通过 Custom::Funcs 访问 v_10::Funcs
        {
        public:
            Funcs(void);            
            T Add(T a, T b);
            T Subtract(T a, T b);
            T Multiply(T a, T b);
            T Divide(T a, T b);  // 只能通过 Custom::v_10::Funcs::Divide 访问
        };
    }
    inline namespace v_20
    {
        template <typename T>
        class Funcs  // 内联，直接通过 Custom::Funcs 访问
        {
        public:
            Funcs(void);
            T Add(T a, T b);
            T Subtract(T a, T b);
            T Multiply(T a, T b);
            std::vector<double> Log(double);
            T Accumulate(std::vector<T> nums);
        };
    }
}
```

---
### 4. 声明

#### 4.1. typedef, using

`typedef` 在任意范围内创建任意既存类型的别名，具有外部链接。

```cpp
typedef struct {/* ... */ } S,*pS;  // 匿名类
typedef void Action();			// 函数
typedef int IntArray[];			// 数组
typedef int(&rFunc)();          // 函数引用
typedef void(*pFunc)(int);		// 函数指针
typedef int int_t, *intp_t, (&fp)(int, ulong), arr_t[10];  
```

`using` 创建指代已定义类型的别名。别名模板指代一族类型的别名。

```cpp
using S = struct { /* ... */ };
using Action = void();
using IntArray = int[];
using rFunc = int(&)();
using pFunc = int(*)();

template <typename T,typename U>
using Vec = vector<T,U>;   // 别名模板
```

>---
#### 4.2. auto, decltype

`auto` 用于类型推断声明或类型推断占位符。

```cpp
auto Hi = "Hello World";    // const char*
auto f(int v) { 	return v * v; }   // int
auto sum(auto typename v1, auto typename v2) { return v + v2; }   // template<class>
```

`decltype` 生成指定表达式的类型，包含 CV 限定和 REF 限定。

```cpp
template<class T> requires !(std::is_reference_v<T>)
using V = decltype(std::forward<T>(T{}));

V<const double> v = 10;   // const double&&
```

`decltype(auto)` 可用于函数占位符的类型推断返回。

```c++
template<typename T, typename U>
auto Func(T&& t, U&& u) -> decltype(auto)
        { return forward<T>(t) + forward<U>(u); };
// or
template<typename T, typename U>
decltype(auto) Func(T&& t, U&& u) { return forward<T>(t) + forward<U>(u); };
```

`func() -> decltype(expr)` 转发函数包装对其他函数的调用。转发函数的返回类型与包装函数的返回类型相同。

```c++
#include <iostream>
#include <string>
#include <utility>
#include <iomanip>

using namespace std;

template<typename T1, typename T2>
auto Plus(T1&& t1, T2&& t2) ->
decltype(forward<T1>(t1) + forward<T2>(t2)) {
    return forward<T1>(t1) + forward<T2>(t2);
}
class X
{
    friend X operator+(const X& x1, const X& x2)
    {
        return X(x1.m_data + x2.m_data);
    }

public:
    X(int data) : m_data(data) {}
    int Dump() const { return m_data; }
private:
    int m_data;
};

int main()
{
    // Integer
    int i = 4;
    cout << "Plus(i, 9) = " << Plus(i, 9) << endl;   // Plus(i, 9) = 13
    // Floating point
    float dx = 4.0;
    float dy = 9.5;
    cout << setprecision(3) <<
        "Plus(dx, dy) = " << Plus(dx, dy) << endl;  // Plus(dx, dy) = 13.5
    // String
    string hello = "Hello, ";
    string world = "world!";
    cout << Plus(hello, world) << endl;  // Hello, world!

    // Custom type
    X x1(20);
    X x2(22);
    X x3 = Plus(x1, x2);
    cout << "x3.Dump() = " <<
        x3.Dump() << endl;  // x3.Dump() = 42
}
```

>---
#### 4.3. extern "C", extern"C++"

语言规范提供在不同编程语言编写的程序单元之间进行连接的功能。每个函数类型、每个具有外部连接的函数名称以及每个具有外部连接的变量名称都具有一个名为语言连接的属性。语言连接封装了与另一种编程语言编写的程序单元链接所需的要求集：调用约定、名称修饰等。

```cpp
extern "C"|"C++" Decl | { Decl-List }
```

`"C"++` 为默认语言连接；`"C"` 允许与 C 函数链接并在 C++ 程序中定义可以从 C 翻译单元调用的函数。

```cpp
extern "C" int open(const char *path_name, int flags); // C function declaration
int main() {
    int fd = open("test.txt", 0); // calls a C function from a C++ program
}
// This C++ function can be called from C code
extern "C" void handler(int) {
    std::cout << "Callback invoked\n"; // It can use C++
}
```

当类成员、带有尾随 `requires` 子句的友元函数、或非静态成员函数出现在 `"C"` 语言块中时，其类型的连接仍然是 `"C++"` 但参数类型（如果有）仍然是 `"C"`。

```cpp
template<typename T>
struct A { struct B; };
template<typename T>
struct A<T>::B {
	extern "C" friend void f(B*) requires true {}   // C language linkage ignored
};
```


>---
#### 4.4. constexpr, consteval, constinit

`constexpr` 声明可以在编译时求出实体的值，这些实体（变量、函数等）用于编译时的常量表达式。`consteval` 指定函数为立即函数，对函数的每次调用都必须产生编译时常量。`constinit` 断言变量具有静态初始化，即零初始化和常量初始化；当声明的变量是引用时，`constinit` 等价于 `constexpr`；`extern constinit` 告知引用外部已初始化的变量。

```cpp
consteval int InitEval() {
	if consteval {
		return 10086;   // 编译时立即函数，只会返回 10086
	}
	else {
		return 10010;
	}
}
constexpr int InitExpr() {
	if consteval {
		return 10086;   // 编译时返回 10086
	}
	else {
		return 10010;   // 运行时返回 10010
	}
}
int main() {
	constinit static int init_v = InitExpr();   // 编译时求值，10086
	constexpr int expr_v = InitExpr();			// 编译时求值，10086
	int x = InitExpr();			// 运行时求值，10010
	int v = InitEval();			// 编译时求值，10086
}
```


>---
#### 4.5. Attribute

可以为类型、对象、代码等引入实现定义的属性。

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

### 5. 基础类型

| Types                                            | Description                                         |
| :----------------------------------------------- | :-------------------------------------------------- |
| `void`                                           | 无类型，声明非类型化非限定的一般指针 `void*`        |
| `bool`                                           | 布尔，真值 `true = 1` 和 `false = 0`                |
| `std::nullptr_t`,`nullptr`                       | 空指针字面值 `nullptr` 的类型                       |
| `char`,`wchar_t`,`char8_t`,`char16_t`,`char32_t` | 字符，`char` 和 `signed/unsigned char` 视为不同类型 |
| `float`,`double`,`long double`                   | 标准浮点类型                                        |
| `[signed]? char/short/int/long/long long`        | 有符号整数类型                                      |
| `unsigned char/short/int/long/long long`         | 无符号整数类型                                      |

> **bool**

可以将任何非零数值（包括无穷和 *NaN*）、非空指针、字符串的布尔表达式返回 `true`；空指针和零数值的布尔表达式返回 `false`。

```cpp
bool _;
_ = "Hi";       // true
_ = "";         // true
_ = 1;          // true
_ = 0.0/-0.0;   // true
_ = 0.0;        // false
```

> **整数字面值**

| Suffix     | description        |
| :--------- | :----------------- |
| `u`,`U`    | 无符号整数         |
| `l`, `L`   | `long/long long`   |
| `ll`, `LL` | `long long`        |
| `z`        | `size_t`,`ssize_t` |

```cpp
auto b = 0b1010'0101              // 二进制
auto o = 013245670;               // 八进制
auto i = 1234567890LL;            // 十进制
auto h = 0x12345'67890'abcdefLL;  // 十六进制

size_t z = -10010z;
ssize_t uz = 10086uz;
```

> **浮点数字面值**

| Suffix         | Description   |
| :------------- | :------------ |
| `f`, `F`       | `float`       |
| `l`, `L`       | `long double` |
| `f16`, `F16`   | `float16_t`   |
| `f32`, `F32`   | `float32_t`   |
| `f64`, `F64`   | `float64_t`   |
| `f128`, `F128` | `float128_t`  |
| `bf16`, `BF16` | `bfloat16_t`  |

```c++
float f         = 3.1415F;
double d        = 3.1415;
long double ld  = 3.1415L;
float16_t f16   = 3.1415f16;
float32_t f32   = 3.1415f32;
float64_t d2    = 3.1415f64;
float128_t f64  = 3.1415f128;
bfloat16_t bf16 = 3.1415bf16;  
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
char u2 = '\101';       // octal, 'A'; 1-3 位八进制，最大 \377
char u3 = '\x41';       // hexadecimal, 'A'
char u4 = '\u0041';     // \u hhhh 'A'
char u5 = '\U00000041'; // \U HHHHHHHH 'A'

// 转义字符
'\'' '\"' '\?' '\\' '\a' '\b' '\f' '\n' '\r' '\t' '\v'
```

>---
#### 5.1. 字符串字面值

字符串文本的种类、类型及其关联的字符编码由其编码前缀和字符序列决定。常规字符串文本和 UTF-8 字符串文本称为窄字符串文本。

```c++
auto s0 =   "hello";   // const char*
auto s1 = u8"hello";   // const char8_t* 
auto s2 =  L"hello";   // const wchar_t*
auto s3 =  u"hello";   // const char16_t*, encoded as UTF-16
auto s4 =  U"hello";   // const char32_t*, encoded as UTF-32
```

原始字符串以 `R"delimiter( char-sequence )delimiter"` 作为引导序列，其中 `delimiter` 最多包含 16 个字符；`R"((a|b))"` 等价于 `"(a|b)"`。

```c++
auto R0 =   R"("Hello \ world")";   // const char*
auto R1 = u8R"("Hello \ world")";   // const char8_t*
auto R2 =  LR"("Hello \ world")";   // const wchar_t*
auto R3 =  uR"("Hello \ world")";   // const char16_t*, encoded as UTF-16
auto R4 =  UR"("Hello \ world")";   // const char32_t*, encoded as UTF-32
```

`s` 后缀表示映射到 `std::string`。

```c++
auto S   =   "hello"s; // std::string
auto S8  = u8"hello"s; // std::u8string
auto SW  =  L"hello"s; // std::wstring
auto S16 =  u"hello"s; // std::u16string
auto S32 =  U"hello"s; // std::u32string

auto R   =   R"("Hello \ world")"s;  // std::string from a raw const char*
auto R8  = u8R"("Hello \ world")"s;  // std::u8string from a raw const char* 
auto RW  =  LR"("Hello \ world")"s;  // std::wstring from a raw const wchar_t*
auto R16 =  uR"("Hello \ world")"s;  // std::u16string from a raw const char16_t*, encoded as UTF-16
auto R32 =  UR"("Hello \ world")"s;  // std::u32string from a raw const char32_t*, encoded as UTF-32
```

两端字符串拼接要求编码前缀相同或其中一个没有编码前缀；任何其他组合格式不正确。

```c++
char str[] = "12" "34";  // "1234"
auto hi = u8"hello" " " "world"s;
auto err = U"hello" " " L"world"; // disagree on prefix
```

### 6. 枚举

枚举包含一组命名的关联常数项，具有基础整数类型（默认为 `int`），默认从 0 开始递增。枚举分为非区分范围枚举（`enum`，枚举项直接访问）和范围枚举（`enum class`，枚举项由枚举名称限定访问），

```c++
enum Identifier [: integerType ] { ... }                 // 无作用域枚举
enum  class | struct  Identifier [: integerType ] {...}  // 作用域枚举
```

无作用域枚举项可以隐式转换为整数；作用域枚举项需要 `static_cast` 强制转换。

```c++
enum Week : unsigned char { 
    Monday /*0*/, Tuesday /*1*/, 
    Wednesday = 10, Thursday /*11*/, Friday, Saturday, Sunday };
enum class Suit { Diamonds, Hearts, Clubs, Spades };

int main() {
    Week day = Monday;
    Suit kind = Suit::Hearts;

    int tue = Tuesday;
    int clubs = static_cast<int>( Suit::Clubs);
}
```

>--- 
#### 6.1. 笼统枚举

没有枚举项的枚举称为笼统枚举；无作用域枚举需要显式指定基础类型。

```cpp
enum E : <base-intType>;
enum class CE [: <base-intType>];
```

可以利用笼统枚举来声明一种新整数类型，例如 `std::byte`。

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
#### 6.2. using enum

`using enum E;` 声明将 `E` 枚举项作为成员引入到全局范围、命名空间、或类型中，在声明的范围内枚举名称可见。在块范围中为导入声明。

```cpp
enum class E { a, b, c, d };
using enum E;
struct S {
	using enum E;
};
namespace N {
	using enum E;
}
int main() {
	using enum E;
	auto e = a;
	auto e1 = ::b;  
	auto e2 = S::c;
	auto e2 = N::c;
}
```



---
### 7. 引用

`T&`,`T&&` 声明 `T` 类型的引用，声明即初始化。`T&` 创建左值引用，`T&&` 创建右值引用，编译器将已具名的右值引用视为左值。引用必须指代一个有效的对象或函数。引用可以转换为兼容的指针或其他左值引用。可以为临时右值创建 `const T&` 左值引用，要避免脱离函数域传递临时右值构造的 `const T&` 而出现悬挂引用。

```C++
// 直接绑定
int n = 10010;
int& lr = n;
int* pl = &lr; assert(&lr == pl);  // okay
int&& rrn = static_cast<int&&>(n); // rri 直接绑定到 n 
int&& rrm = std::move(lr);         // rrm 直接绑定到 n

// 间接绑定
int&& rr = 10086;
const int& cr = 110;             // 指代 110.0 的临时量
const int& lm = std::move(911);  // 指代 991.0 的临时量
double&& rrd = rr;               // 指代值 (double)rr 的临时量
rrd = 1000000000000;
cout << rr;   // 10086
```

`&&` 保留对右值表达式的引用；*Rvalue* 引用支持 “移动语义”。利用移动语义，可以将资源（如动态分配的内存）从一个对象转移到另一个对象，例如 `std::move` 某个对象转换为右值引用。也允许从临时对象转移资源。

```c
class Sample {};
int main() {
    Sample s = {};
    Sample& r = s;
    Sample&& rr = std::move(r); 
    Sample&& rs = std::move(s);   // s,r,rr,rs 引用同一对象。 
}
```

实现移动语义，通常可以向类提供 “移动构造函数”（`class(class&&)`），或提供移动赋值运算符（`operator=(class&&)`）。利用移动语义，可以直接移动对象而不必执行成本高昂的内存分配和复制操作。例如，`std::string` 实现了使用移动语义的操作：

```cpp
string s = string("h") + "e" + "ll" + "o";
```

>---
#### 7.1. 临时量的生存期

当引用绑定到临时对象或它的子对象，临时对象的生存期被延长以匹配引用的生存期。例外情景：
- `return` 语句中绑定到函数返回值的临时量不会被延续。这种 `return` 语句始终返回悬垂引用 (C++26 前)。
    ```cpp
    const std::string& getString() {
        return "hello";  // 临时字符串在 return 结束时销毁
    }
    int main() {
        const std::string& s = getString(); // s 是悬垂引用
        std::cout << s;                     // 未定义行为
    }
    ```
+ 在函数调用中绑定到函数形参的临时量，生存期仅到当次函数调用的全表达式结尾为止；如果函数返回一个生命长于全表达式的引用，那么它会成为悬垂引用。

    ```cpp
    const string& forwardRef(const string& param) {
        return param;  // 返回绑定到临时对象的引用
    }
    int main() {
        const string& ref = forwardRef("hello");  // 临时对象在全表达式后销毁
        std::cout << ref;                         // 未定义行为
    }
    ```

- 绑定到 `new` 表达式中所用的初始化器中的引用的临时量，存在到表达式结尾，而非是被初始化对象的存在期间。如果被初始化对象的生命长于全表达式，那么它的引用成员将成为悬垂引用。

    ```cpp
    struct Test {
        const string& ref;  // 引用成员
        Test(const string& r) : ref(r) {}  // 绑定临时对象
    };
    int main() {
        // 直接绑定临时对象（悬垂引用）
        Test* t1 = new Test("Hello");  
        std::cout << t1->ref;          // 未定义行为：访问已释放的内存

        // 通过函数返回临时对象（同样悬垂）
        auto create_temp = []() { return "10086"; };
        Test* t2 = new Test(create_temp());  // 临时对象在全表达式后销毁
        std::cout << t2->ref;                // 未定义行为
    }
    ```


+ 绑定到用直接初始化语法（圆括号），而非列表初始化语法（花括号）初始化的聚合体的引用元素中的引用的临时量，存在直至含该初始化器的全表达式末尾为止。

    ```cpp
    struct Test {
    	int v;
    	const string& ref;
    };

    int main() {

    	Test t1(1, "hello");
    	Test t2{ 2, "Hello" };

    	cout << t1.ref;   // 悬挂引用
    	cout << t2.ref;   // Hello
    }
    ```

---
### 8. 指针

指针可以指向类型化对象、函数或非类型指针 `void*`，指向类型成员声明成员指针。原始指针是指其生存期不受封装对象控制的指针。间接寻址无效指针值（导致内存损坏）或空指针值（导致异常），行为未定义。不存在指向引用或位域的指针。

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

`new` 运算返回一个堆分配对象的指针，当不再需要堆分配的对象时，必须对拥有对象的指针（或其副本）显式释放（`delete`），未能释放内存将会导致内存泄漏。。

```c++
MyClass* mc = new MyClass(); // allocate object on the heap
mc->print(); // access class member
delete mc; 	 // delete object (please don't forget!)
```

指针和数组密切相关。当数组按值传递给函数时，它将作为指向第一个元素的指针传递。
- `sizeof` 运算符返回数组的总大小（以字节为单位）；
- 当数组被传递给函数时，它会衰减为指针类型；
- 当 `sizeof` 运算符应用于指针时，它将返回指针大小。

`void*` 的指针仅指向原始内存位置。例如在 C++ 代码和 C 函数之间传递时需要使用 `void*`；将类型化指针强制转换为 `void*` 时，内容保持不变但类型信息丢失，无法执行递增或递减操作。


>---
#### 8.1. 函数指针

函数指针由非成员函数或静态成员函数的地址初始化，存在函数到指针的隐式转换。解引用函数指针生成指向函数的左值。

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

	Foo* pf2 = &rf;	 // 转换引用为指针 
	pf2(4);

	Foo& rf2 = *pf;  // 解指针为引用
	rf2(5);

	void(&Tfoo)(double) = foo;   // foo<double>
	void(&Tfoo)(int) = foo;      // foo(int)  重载决策优先级更高
}
```

涉及函数指针的声明可以通过类型别名简化复杂抽象声明。

```cpp
void foo(int v);
using Fun = void(int);
typedef void DFun(int);   

// 等价声明
Fun& pf = foo;
DFun& df = foo;
void (*P)(int) = foo;
```

> **指针运算**

```c++
// address-of
int * pi = &i;
const type cp = &const_type;
volatile type vp = &vola_type;
// 算数
int * p = p + integer;
p += integer;
p++;
int * p = p - integer;
p -= integer;
p--;
ptrdiff_t diff = p1 - p2;
```

>---
#### 8.2. 成员指针

成员指针指向类类型的非静态成员。成员指针充当成员指针访问 `.*` 或 `->*` 的右操作数：

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

指向一个可访问且无歧义的非虚基类数据成员的指针，可以隐式转换成指向派生类的同一数据成员的指针。从指向派生类的数据或非虚函数的成员指针到指向无歧义（非虚）基类的成员指针，允许由 `static_cast` 和显式转换来进行，即使基类没有该成员。

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

	int B::* bM = &B::M;
	int D::* dM = bM;        // 隐式转换, int B::* -> int D::*
	cout << d.*dM << endl;   // 10086


	double D::* dN = &D::N;
	double B::* bN = static_cast<double B::*>(dN);
	cout << d.*bN << endl;   // 3.1415
	B b;
	cout << b.*bN << endl;   // 未定义行为

	void (D:: * dF)() = &D::F;
	void (B:: * bF)() = static_cast<void (B::*)()>(dF);
	(d.*bF)();    // F
	(b.*bF)();    // 未定义行为
}
```

指向基类的成员函数的指针可以隐式转换成指向派生类的同一成员函数的指针。

```cpp
struct B {  
	virtual void F() { cout << "B\n"; }
};
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
#### 8.3. 智能指针

智能指针（`<memory>`）用于确保程序不存在内存和资源泄漏且是异常安全的。

- `unique_ptr` 只允许基础指针的一个所有者。可以移到新所有者，但不会复制或共享。小巧高效且大小等同于一个指针，支持 *rvalue* 引用。
- `shared_ptr` 采用引用计数的智能指针。将一个原始指针分配给多个所有者时使用。直至所有 `shared_ptr` 所有者超出了范围或放弃所有权，才会删除原始指针。
- `weak_ptr` 结合 `shared_ptr` 使用的特例智能指针。`weak_ptr` 提供对一个或多个 `shared_ptr` 实例拥有的对象的访问，但不参与引用计数。

```c++
Point* p1 = new Point{};
unique_ptr<Point> p2(new Point{});
```

> **unique_ptr**

`unique_ptr` 不共享它的指针。它无法复制到其他 `unique_ptr`，无法通过值传递到函数，也无法用于需要副本的任何 C++ 标准库算法。只能移动 `unique_ptr`。创建 `unique_ptr` 实例并在函数之间传递实例：

```c++
struct Song {
    wstring artist;
    wstring title;
    Song(wstring artist, wstring title) : artist{ artist }, title{ title } {};
};
unique_ptr<Song> SongFactory(const std::wstring& artist, const std::wstring& title) {
    // Implicit move operation into the variable that stores the result.
    return make_unique<Song>(artist, title);
}
void MakeSongs() {
    // Create a new unique_ptr with a new object.
    auto song = make_unique<Song>(L"Mr. Children", L"Namonaki Uta");
    // Use the unique_ptr.
    vector<wstring> titles = { song->title };
    // Move raw pointer from one unique_ptr to another.
    unique_ptr<Song> song2 = std::move(song);
    // Obtain unique_ptr from function that returns by value.
    auto song3 = SongFactory(L"Michael Jackson", L"Beat It");
}
```

使用 `make_unique` 将 `unique_ptr` 创建到数组，但无法使用 `make_unique` 初始化数组元素：

```c++
// Create a unique_ptr to an array of 5 integers.
auto p = make_unique<int[]>(5);
// Initialize the array.
for (int i = 0; i < 5; ++i) {
    p[i] = i;
    wcout << p[i] << endl;
}
```

> **shared_ptr**

`shared_ptr` 为多个所有者需要管理对象生命周期的一种智能指针；`shared_ptr` 可复制，按值将其传入函数参数，然后将其分配给其他 `shared_ptr` 实例。 所有实例均指向同一个对象，并共享对一个 “控制块”（每当新的 `shared_ptr` 添加、超出范围或重置时增加和减少引用计数）的访问权限。当引用计数达到零时，控制块将删除内存资源和自身。

```c++
struct MediaAsset {
    virtual ~MediaAsset() = default; // make it polymorphic
};
struct Song : public MediaAsset {
    std::wstring artist;
    std::wstring title;
    Song(const std::wstring& artist_, const std::wstring& title_) :
        artist{ artist_ }, title{ title_ } {}
};
struct Photo : public MediaAsset {
    std::wstring date;
    std::wstring location;
    std::wstring subject;
    Photo(const std::wstring& date_,
          const std::wstring& location_,
          const std::wstring& subject_) 
        : date{ date_ }, location{ location_ }, subject{ subject_ } {}
};
```

第一次创建内存资源时，使用 `make_shared` 函数创建 `shared_ptr` 或显式 `new` 并传递给 `shared_ptr`：

```c++
// Use make_shared function when possible.
auto sp1 = make_shared<Song>(L"The Beatles", L"Im Happy Just to Dance With You");
// Ok, but slightly less efficient. 
// Note: Using new expression as constructor argument
// creates no named variable for other code to access.
shared_ptr<Song> sp2(new Song(L"Lady Gaga", L"Just Dance"));
// When initialization must be separate from declaration, e.g. class members, 
// initialize with nullptr to make your programming intent explicit.
shared_ptr<Song> sp5(nullptr);
//Equivalent to: shared_ptr<Song> sp5;
//...
sp5 = make_shared<Song>(L"Elton John", L"I'm Still Standing");
```

在容器中复制元素：

```c++
vector<shared_ptr<Song>> v {
  make_shared<Song>(L"Bob Dylan", L"The Times They Are A Changing"),
  make_shared<Song>(L"Aretha Franklin", L"Bridge Over Troubled Water"),
  make_shared<Song>(L"Thalía", L"Entre El Mar y Una Estrella")
};
vector<shared_ptr<Song>> v2;
remove_copy_if(v.begin(), v.end(), back_inserter(v2), [] (shared_ptr<Song> s) {
    return s->artist.compare(L"Bob Dylan") == 0;
});

for (const auto& s : v2) {
    wcout << s->artist << L":" << s->title << endl;
}
```

使用 `dynamic_pointer_cast`、`static_pointer_cast` 和 `const_pointer_cast` 来转换 `shared_ptr`：

```c++
vector<shared_ptr<MediaAsset>> assets {
  make_shared<Song>(L"Himesh Reshammiya", L"Tera Surroor"),
  make_shared<Song>(L"Penaz Masani", L"Tu Dil De De"),
  make_shared<Photo>(L"2011-04-06", L"Redmond, WA", L"Soccer field at Microsoft.")
};

vector<shared_ptr<MediaAsset>> photos;
copy_if(assets.begin(), assets.end(), back_inserter(photos), [] (shared_ptr<MediaAsset> p) -> bool
{
    // Use dynamic_pointer_cast to test whether element is a shared_ptr<Photo>.
    shared_ptr<Photo> temp = dynamic_pointer_cast<Photo>(p);
    return temp.get() != nullptr;
});

for (const auto&  p : photos)
{
    // We know that the photos vector contains only shared_ptr<Photo> objects, so use static_cast.
    wcout << "Photo location: " << (static_pointer_cast<Photo>(p))->location << endl;
}
```

按值传递 `shared_ptr`。这将调用复制构造函数，增加引用计数，并使被调用方成为所有者。按引用或常量引用传递 `shared_ptr`。引用计数不会增加，只要调用方不超出范围，被调用方就可以访问指针。传递基础指针或对基础对象的引用，被调用方能够使用对象，但不会共享所有权或延长生存期。

```c++
void use_shared_ptr_by_value(shared_ptr<int> sp);
void use_shared_ptr_by_reference(shared_ptr<int>& sp);
void use_shared_ptr_by_const_reference(const shared_ptr<int>& sp);
void use_raw_pointer(int* p);
void use_reference(int& r);

void test() {
    auto sp = make_shared<int>(5);
    // Pass the shared_ptr by value.
    // This invokes the copy constructor, increments the reference count, and makes the callee an owner.
    use_shared_ptr_by_value(sp);
    // Pass the shared_ptr by reference or const reference.
    // In this case, the reference count isn't incremented.
    use_shared_ptr_by_reference(sp);
    use_shared_ptr_by_const_reference(sp);
    // Pass the underlying pointer or a reference to the underlying object.
    use_raw_pointer(sp.get());
    use_reference(*sp);
    // Pass the shared_ptr by value.
    // This invokes the move constructor, which doesn't increment the reference count
    // but in fact transfers ownership to the callee.
    use_shared_ptr_by_value(move(sp));
}
```

> **weak_ptr** 

有时对象必须存储一种方法来访问 `shared_ptr` 的基础对象，而不会导致引用计数递增。例如在 `shared_ptr` 实例之间有循环引用时。`weak_ptr` 本身不参与引用计数，它无法阻止引用计数变为零。但是可以使用 `weak_ptr` 尝试获取初始化该副本的 `shared_ptr` 的新副本。若已删除内存，则 `weak_ptr` 的 `bool` 运算符 `if(weak_ptr)` 返回 `false`。

这些 `Controller` 对象表示计算机进程的一些方面，它们独立运行。每个控制器必须随时能够查询其他控制器的状态，每个控制器都包含一个专用 `vector<weak_ptr<Controller>>` 用于实现此目的。每个向量都包含一个循环引用，因此使用 `weak_ptr` 实例而不是 `shared_ptr`。

```c++
class Controller
{
public:
   int Num;
   wstring Status;
   vector<weak_ptr<Controller>> others;
   explicit Controller(int i) : Num(i), Status(L"On")  {
      wcout << L"Creating Controller" << Num << endl;
   }
   ~Controller() {
      wcout << L"Destroying Controller" << Num << endl;
   }
   // Demonstrates how to test whether the
   // pointed-to memory still exists or not.
   void CheckStatuses() const {
      for_each(others.begin(), others.end(), [](weak_ptr<Controller> wp) {
         auto p = wp.lock();
         if (p) 
            wcout << L"Status of " << p->Num << " = " << p->Status << endl;
         else
            wcout << L"Null object" << endl;
      });
   }
};
void RunTest() {
   vector<shared_ptr<Controller>> v {
       make_shared<Controller>(0),
       make_shared<Controller>(1),
       make_shared<Controller>(2),
       make_shared<Controller>(3),
       make_shared<Controller>(4),
   };
   // Each controller depends on all others not being deleted.
   // Give each controller a pointer to all the others.
   for (int i = 0; i < v.size(); ++i) {
      for_each(v.begin(), v.end(), [&v, i](shared_ptr<Controller> p) {
         if (p->Num != i) {
            v[i]->others.push_back(weak_ptr<Controller>(p));
            wcout << L"push_back to v[" << i << "]: " << p->Num << endl;
         }
      });
   }
   for_each(v.begin(), v.end(), [](shared_ptr<Controller> &p) {
      wcout << L"use_count = " << p.use_count() << endl;
      p->CheckStatuses();
   });
}

int main() {
   RunTest();
   wcout << L"Press any key" << endl;
   char ch;
   cin.getline(&ch, 1);
}
```


---
### 9. 数组

数组具有相同类型、连续存储数组元素的对象序列，不存在引用的数组或函数的数组。分配和访问堆栈数组的速度比堆数组更快，堆栈数组在编译时长度确定。堆数组通常由 `new T[size]` 分配，并由 `delete[] T` 负责释放。堆栈数组在函数返回时自动清理。CV 限定将作用于数组元素类型而不是数组本身。动态分配的数组大小可以为零，即 `new T[0]`。


```c++
int stack[] = { 1,2,3,4 };       // 堆栈数组 int[4]
int* heap = new int[size] {0};   // 堆数组 int[size]
//  ... use heap
delete[] heap;
```

存在从数组类型的左值和右值到指针类型的右值的隐式转换，它构造一个指向数组首元素的指针（`auto* p = arr` 等价于 `auto* p = &arr[0]`）。数组作为参数传递给函数时实际传递数组的首元指针，且不包括数组大小信息。

```c++
template <typename T>
void process(size_t len, T arr[])  // or T* arr
{
    for (size_t i = 0; i < len; ++i)
        // do something with arr[i]
        cout << arr[i] << endl;
}

int main() {
    int arr[4] = { 1, 2, 3, 4 };
    process(4, arr);
}
```

> **多维数组**

```c++
int arr[7]{};
int arr3x4[3][4]{};  	  
int arr5x6x7[5][6][7]{};  

// 数组退化至指针
int* parr = arr;  // ==>  parr = arr;
int (*parr2)[4] = arr2;
int (*parr3)[6][7] = arr3;
```

>---
#### 9.1. **数组初始化**

数组是一种聚合体，由聚合初始化器初始化：

```cpp
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

元素类型具有类构造函数的对象数组调用构造函数执行初始化，剩余元素由默认构造函数初始化。

```c++
class Point {
public:
    Point() = default;
    Point(int x, int y) : x(x), y(y) {};
    int x, y;
};
int main() {
    Point aPoint[3] = {
       Point(3, 3)     // 其他元素由默认构造函数初始化
    };
}
```

>---
#### 9.2. 字符数组

字符数组可以由相应类型的字符串字面值直接初始化，字符数组长度应大于字符串长度（包括隐式的空终止符），剩余元素零初始化。

```c
char hi[] = "Hi"; // char[3] {'\H','i','\0'}
char hi2[4] = "Hi"; // char[4] {'\H','i','\0','\0'}
wchar_t hi3[6] = { L"Hello" };  // len >= 6;
```



---
### 10. 表达式

C++ 语言包括所有 C 运算符并添加多个新的运算符。C++ 运算符的优先级和关联性：

| Category     | Operators                                                                                            |
| :----------- | :--------------------------------------------------------------------------------------------------- |
| 基本表达式   | `x::y`,`x.y`,`x->y`,`f(x)`,`x[y]`,`x++`,`x--`,`typeid(x)`,`T{v}`                                     |
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
| 逗号运算     | `(expr1, expr2 [, expr3, ... exprN])` 返回 `exprN` 表达式结果                                        |
          
> **表达式值类别**

表达式具有类型和值类别两种特性，每个表达式只属于三种基本值（纯右值、亡值、左值）类别中的一种。值类别包括有：
- 泛左值 *glvalue*：求值时可确定某个对象或函数标识的表达式。
- 纯右值 *prvalue*：运算符操作数的值、初始化某个对象的结果对象。
- 亡值 *xvalue*：资源能够被重新使用的对象或位域的泛左值。
- 左值 *lvalue*：非亡值的泛左值。
- 右值 *rvalue*：纯右值或亡值。

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
#### 10.1. 类型转换

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
	p = Point{ 1, 2 };
	p = Point{ .x = 1, .y = 2 };
	p = typename ::Point{ 1,2 };
}
```

> 用户定义转换函数

```cpp
operator <Type-id>                      // 隐式转换
explicit operator <Type-id>             // 显式转换
explicit (bool-expr) operator <Type-id> // 条件转换。true 为显式，false 为隐式
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
	explicit operator string() { return to_string(distance) + "km"; }  // 显式类型转换

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
int main() {
	KM d1 = 10km;    // int to KM
	KM d2 = 3.14km;  // double to KM
	KM d3 = 10;
	KM d4 = 3.14L;
	long double d = d1;  // or (long double)d1;
	cout << string(d1);  // 10.000000km
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
	d.operator void();  // void
	B& b2 = d.operator B & ();  // B&
}
```

> **const_cast**

`const_cast` 添加或移除任意对象类型指针、引用、数据成员指针的 CV 限定（编译时行为），结果指针或引用仍关联原始对象。原对象的 CV 限定不变，因此通过移除 `const` 限定的指针或引用修改 `const` 或 `constexpr` 变量值的行为未定义。`const_cast` 也可以将一个空指针值转换为其他目标类型的空指针值。

```c++
int num = 10010;
const int& cr_num = num;
int& r_num = const_cast<int&>(cr_num);
r_num = 110;
assert(r_num == num);  // assert success; both are 110

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

`static_cast` 用于基础类型之间转换、指针或引用在继承体系中的向上转型 / 向下转型、调用类的转换构造函数或转换运算符、枚举与整数之间转换等。`static_cast` 编译时检查源类型和目标类型是否满足合法转换规则，不进行运行时检查。`static_cast<void>(expr)` 为弃值表达式。

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

`dynamic_cast` 用于多态类型的转换并执行运行时类型检查，如从多态基类的指针或引用到派生类指针或引用的向下转换、类层次之间的交叉转换（侧向转换）等。转换指针失败返回 `nullptr`；转换引用失败引发 `bad_cast`。转换不移除 CV 限制。

```cpp
// 多态类类型至少包含一个虚函数
struct V { virtual void f() { } };
struct A : virtual V { };
struct B : virtual V
{
	B(V* v, A* a)
	{
		dynamic_cast<B*>(v); // 良好定义：V 是 B 的基类，V* -> B*
		dynamic_cast<B*>(a); // 未定义行为：A 不是 B 的基类
	}
	void f() override { cout << "B.f\n"; };
};
struct D : A, B {
	D() : B(static_cast<A*>(this), this) {}
};
int main()
{
	D d; // 最终派生对象, 具有 A,B,V
	A& a = d; // 隐式向上转换，可以用 dynamic_cast，非必须
	D& new_d = dynamic_cast<D&>(a); // 向下转换, A&(d) -> D&
	B& new_b = dynamic_cast<B&>(a); // 侧向转换, A&(d) -> D& -> B&
}
```

如果有多重继承，可能会导致不明确。使用虚拟基类时，可能会导致更多不明确的情况。

```c++
class A { virtual void f(); };
class B : public A { virtual void f(); };
class C : public A { virtual void f(); };
class D : public B, public C { virtual void f(); };

D* pd = new D;
A* pa = dynamic_cast<A*>(pd);   // ambiguous cast fails at runtime

// `D*` 指针可以安全地转换为 `B*` 或 `C*`。但如果 `D*` 强制转换为 `A*`，将导致不明确的强制转换错误。
// 若要解决此问题，可以执行两个明确的强制转换。
B* pb = dynamic_cast<B*>(pd);   // D* -> B*
A* pa2 = dynamic_cast<A*>(pb);  // B* -> A*
```

如果目标类型是 `void*`，则在运行时检查并确定转换对象的实际类型。结果是指向的源类型完整对象的指针。

```c++
class A {virtual void f();};
class B {virtual void f();};

A* pa = new A;
B* pb = new B;
void* pv = dynamic_cast<void*>(pa);  // pv now points to an object of type A
pv = dynamic_cast<void*>(pb);        // pv now points to an object of type B
```


> **reinterpret_cast** 

`reinterpret_cast` 执行位模式重新解释，直接复制源类型的二进制位到目标类型，在编译器完成转换。允许在任意指针或引用之间、整数与指针之间转换，无视类型关系。转换不移除 CV 限制。

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
#### 10.2. 内存分配

`new` 创建存储并初始化具有动态存储期的对象，这些对象的内存可能在未来的某个时刻由 `delete` 释放。

```cpp
[::] new ( placement-Params )? Typeid [initializer]
[::] delete ( placement-Params )? expr
[::] delete ( placement-Params )? [] expr

S* ps = new S;
int * p = new int[]{1,2,3};

delete ps;
delete[] p;
```

`expr` 由之前的 `new` 创建指向对象的指针或空指针，间接引用（`*p`）被释放的指针行为未定义。

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


> **可替代函数**

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
#### 10.3. sizeof

`sizeof` 查询对象或类型的字节大小，结果为 `size_t`。操作数不能是 `void`、函数、不完整类型、位域、动态数组、外部数组等。

```cpp
sizeof ( Type );
sizeof ( expr ) | expr;

auto Asz = sizeof(ClassA);
auto osz = sizeof obj;
auto len = sizeof arr / sizeof arr[0];
```

`sizeof...(P)` 查询包中的元素数目。

```cpp
template<typename... Ts>
void func(Ts... args)
{
	const int size = sizeof...(args) + 2;   // Ts 包元素数目
	int res[size] = { 1, args..., 2 };
	// 因为初始化器列表保证顺序，所以这可以用来对包的每个元素按顺序调用 cout：
	int dummy[sizeof...(Ts)] = { (std::cout << args << ",", 0)...};   // 折叠表达式
}
func(1, 2, 3, 4.0, 5.6); // 1,2,3,4,5.6,
```

>---
#### 10.4. typeid

```cpp
typeid ( <Type> | <expr> )
```

`typeid` 查询类型或表达式的类型信息，返回 `const std::type_info&`。*Type* 不能是不完整类类型；*expr* 不能是解空指针操作，指针需要指向有效对象。对于引用或 CV 限定，等同于返回其基础类型的类型信息。

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

>---
#### 10.5. 断言

可以在命名空间、类或块范围中执行 `static_assert( <const-expr> [, mess])` 编译时静态断言。`assert`、`_assert`、`_wassert` 执行运行时断言。

```c
// 命名空间范围断言
static_assert(sizeof(void*) == 8, "64-bit code generation is not supported.");
template <class CharT>
class b_string {
    // 类范围断言：模板参数是否为标准布局
    static_assert(std::is_standard_layout<CharT>::value, 
        "Template argument CharT must be a standard layout type in class template basic_string");
    // ...
};
struct NonPOD {  // 非标准布局
    NonPOD(const NonPOD&) {}
    virtual ~NonPOD() {}
};

int main()
{
    b_string<char> bs;
    b_string<NonPOD> vsp;  // 断言失败, 无法通过编译
}
```


---
### 11. 语句
#### 11.1. 空语句

```c++
while(cond)
    ;
```

>---
#### 11.2. 条件控制语句：if, switch

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

```c++
swtich ( [init-expr;] int-cond ){
    case cond1:
        statements
        [ break | return ];  // 允许贯穿，编译器警告
    case cond2: 
        statements
        [[fallthrough]];     // 贯穿属性，不警告 
    [default: ]
}
```
```c++
switch (1) {
case 1:
    cout << 1 << endl;
    [[fallthrough]];
case 2:
    cout << 2 << endl;
default:break; // warning: fallthrough
}
```

在 `case` 中声明的变量属于 `switch` 语句范围，它们共享名称。

```c++
/* 
switch (char szChEntered[] = "Character entered was:"; 'b')  // init at swtich condition; is ok
{   // const char* szChEntered; 
*/
switch ('b')
{
    // Declaration
    const char* szChEntered; // = "Declaration";  err
case 'a':
{
    // ok, reDeclaration. Local scope in block [case 'a'].
    const char* szChEntered = "Character entered was: ";
    cout << szChEntered << "a\n";
    break;
}
case 'b':
    // Error. Value of szChEntered undefined.
    cout << szChEntered << "b\n";
    break;
default:
    // Error. Value of szChEntered undefined.
    cout << szChEntered << "neither a nor b\n";
    break;
}
```

>---
#### 11.3. 循环语句：while, do, for, for-range 

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
```cpp
for (int a = 1, b = 10; a < b; a++, b--, cout << a * b << ',')
    ; // 18,24,28,30,30,
```

> **for-range**

```c++
for ( [init-expr;] elem-declaration : expression ) {
    statements
}
```
```c++
int x[10]{ 1,2,3,4,5,6,7,8,9,0 };
for (auto e : x)
    cout << e << ",";
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
#### 11.4. 跳转语句：break, continue, return, goto

`break` 终止最内层循环或 `switch`。`continue` 跳转至最内层循环的块末尾并开启下一次循环。`return` 终止函数执行并返回到调用方。

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

标签具有函数范围，在整个函数内可见，与声明位置无关。`goto` 不能跳过局部变量的初始化。

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
        cerr << "At Test1 label." << endl;
    }
}  // At Test2 label.
```

---
### 12. 函数类型

函数声明可以在任何作用域出现，可以作为非成员函数、友元函数、成员函数；函数定义仅在命名空间或类的作用域中出现。函数的返回类型不能是函数类型或数组类型，但可以是指针或引用。
- 非静态成员函数可以具有 CV 和引用（*ref*）限定。
- 非成员函数可以具有链接属性。
- `requires` 声明与模板化函数关联的约束，属于函数签名的一部分。
- `noexcept` 声明函数异常规范。
- 对于返回类型复杂时，函数声明可以是尾随返回形式（`auto`）。

```cpp
FuncDecl     = Func ( Params ) [CV] [REF] [Except] [Attr]
TrailingFunc = auto FuncDecl -> <trailing>
TemplateFunc = template <Params> [RequiresClause] FuncDecl
```
```cpp
void func() noexcept;
static int func(int, ...);        // 变参, 内部链接  
auto func(int) -> int (*)(int);
struct S {
	virtual void func(int) const = 0;    // 纯虚函数
	void func(int, int) volatile&;		
	template <typename T>				 // 函数模板, 静态
	static auto func(T) -> decltype(new T{});
};
```

> **返回类型推导**

返回声明为类型占位符 `auto` 表示使用返回类型推导；如果返回类型是 `decltype(auto)`，那么返回类型是将 `return expr` 的操作数到 `decltype( expr )` 中时所得到的类型。虚函数和协程不能使用返回类型推导。

```cpp
int x = 0;
auto func() { return x; }        // auto is int, reutun int
const auto& func() {return x; }  // auto is int, return const int& 

decltype(auto) funcA() { return x; }          // decltype(x) is int
decltype(auto) funcB() { return (x); }        // decltype((x)) is int&
decltype(auto) funcC() { return move(x); }    // decltype(move(x)) is int&&
```

>---
#### 12.1. **参数列表**

```cpp
Param = [Attr] ParamDecl [= Initializer]    
       | [Attr] this ParamDecl        // 显式对象形参       
       | void                 
```

具有 `this` 显式对象形参的成员函数非虚非静态，可以是 lambda，无 CV 和引用限定。显式对象形参只能是首个形参，无默认值。其他非静态成员函数具有隐式对象形参 `this`。尾随参数可以具有默认值，虚函数的重写函数不会从基类定义获得默认实参。

```cpp
struct S {
    void F(void) const;
	void F(this S&);			// 显式对象参数
	void F(int x);				// 具名参数声明
	void F(double = 3.14);		// 默认值
	void F(int, int*, int (*(*)(double))[3] = nullptr);   // 抽象参数声明
	void F(int x = 0, int y = 0, ...);       // 变参
    // void (*f)(this S) = [](this S s) {};
    template<typename... Args>
	void F(Args..., ...);		// 带形参包的变参函数模板
};
```

除了函数调用运算符和下标运算符外，运算符函数不能有默认实参。

```cpp
class C
{
    int operator++(int i = 0); // 非良构
    int operator[](int j = 0); // OK
    int operator()(int k = 0); // OK
};
```

变长参数由 `<cstdarg>` 访问；具有相同类型的变长实参也可以由 `std::initializer_list` 访问。

```cpp
template<typename ... Ts>
void Iterator(Ts ... args) {
	for (std::initializer_list l{ args... }; auto e : l)
		cout << e << ",";
};

Iterator(1, 2, 3, 4, 5);
```


>---
#### 12.2. **函数定义**

非静态成员函数可以定义虚说明（`override`、`final`），`virtual` 声明虚函数。非成员函数只能在命名空间范围定义，成员函数类范围定义默认具有 `inline` 内联。显式预置（`default`）的函数只能是特殊成员函数或比较运算符函数。任何弃置（`delete`）函数的使用都会导致程序无法编译。

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
    static inline void inFunc() {}		// 内联函数
    void FuncD() noexcept try {	}		// 函数 try 块
	catch (...) { }
};
```

> **inline**


在类类型定义中定义的成员函数、弃置函数、`default` 预置声明或隐式生成的函数、首次声明 `constexpr|consteval` 的函数、首次声明 `constexpr` 的静态数据成员等具有隐式 `inline`。内联函数或变量的定义必须在访问它的翻译单元中可访问。

多个翻译单元的相同内联函数必须具有相同的定义。在内联函数中，所有函数定义中的函数本地静态对象在所有翻译单元之间共享。

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



>---
#### 12.3. lambda 表达式

```cpp
[Captures] [FrontAttr] (Params)? [[Specs][Except][BackAttr]] [-> Trailing][Requires] { Body };
[Captures] <TParams> [TRequires] [FrontAttr] (Params)? [[Specs][Except][BackAttr]] [-> Trailing][Requires] { Body };

[] (int v) static constexpr -> int {return v;}    // 无捕获静态 lambda
[] <class C> -> decltype(auto) { return C{}; };   // lambda 模板
[=, &i] { i++; };           // 默认按复制捕获，引用捕获 i
[&, i] { return i * i; };   // 默认按引用捕获，复制捕获 i
```

lambda 表达式可以构造闭包，是能够捕获作用域中的变量的匿名函数对象。*Specs* 支持 `mutable|static`、`constexpr|consteval`。lambda 默认是 `const`，除非是 `mutable` 或具有显式对象形参 `this`。`mutable` 允许函数体修改按复制捕获的对象，以及调用非 `const` 成员函数。`static` 具有空捕获规范。

> **Captures**

`[ ]` 无捕获；`[=]` 默认为按值捕获，后继的简单捕获只能显式按引用 `&v` ；`[&]` 默认为按引用捕获，后继的简单捕获只能显式按值复制 `v`；在类非静态成员函数中可以将 `this` 传递给 *Capture*，`*this` 按值复制传递。当 `[=]` 或 `[&]` 总能

如果按引用隐式或显式捕获非引用实体，而在该实体的生存期结束之后调用闭包对象的 `lambda()`，会发生未定义行为。C++ 的闭包并不延长按引用捕获对象的生存期。
 
```cpp
struct S { void f(int i); };
void S::f(int i)
{
	[] static { };   // 无捕获，静态
	[&] {};          // 默认按引用捕获
	[&, i] {};       // 按引用捕获，但 i 按值捕获
	//[&, &i] {};    // 按引用捕获为默认时的按引用捕获
	[&, this] {};    // 等价于 [&]
	[&, this, i] {}; // 等价于 [&, i]
	[=] {};          // 默认按复制捕获
	[=, &i] {};      // 按复制捕获，但 i 按引用捕获
	[=, *this] {};   // 按复制捕获外围的 S
	[=, this] {};    // 等价于 [=]
}
```

捕获可以具有初始化器，其行为如同声明并显式捕获以 `auto` 声明并拥有相同初始化器的变量，该变量的作用域是 lambda 表达式体，但：
- 如果按复制捕获，引入闭包对象的非静态数据成员和该变量将被视为引用同一对象；换言之，源变量并不实际存在，而经由 `auto` 类型推导和初始化均应用到该非静态数据成员；
- 如果按引用捕获，那么引用变量的生存期在闭包对象的生存期结束时结束。
这可以用于，以 `x = std::move(x)` 这样的捕获符捕获仅可移动的类型。

```cpp
int x = 4;
 
auto y = [&r = x, x = x + 1]() -> int {
    r += 2;
    return x * x;
}(); // 更新 ::x 到 6 并初始化 y 为 25。
```

如果 lambda 表达式在默认实参中出现，那么它不能显式或隐式捕获任何内容，除非所有捕获都带有初始化器，并满足可以在默认实参中出现的表达式的约束条件。

```cpp
void f()
{
	int i = 1;

	void g1(int = [i] { return i; }()); // ERR：有捕获内容
	void g2(int = [i] { return 0; }()); // ERR：有捕获内容
	void g3(int = [=] { return i; }()); // ERR：有捕获内容

	void g4(int = [=] { return 0; }());       // OK：无捕获
	void g5(int = [] { return sizeof i; }()); // OK：无捕获

	void g6(int = [x = 1] { return x; }()); // OK：1 可以在默认实参中出现
	void g7(int = [x = i] { return x; }()); // ERR：i 不能在默认实参中出现
}
```

> **闭包类型**

lambda 是纯右值表达式的闭包类型，仅当捕获为空时该闭包类型是结构化类型。闭包类型具有以下成员：

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
#### 12.4. 协程

c++ 协程设计为无栈协程，通过返回到调用方来暂停执行。允许顺序代码异步执行。不使用可变参数、占位符返回类型、常规 `return` 语句。`constexpr/consteval`、构造函数、析构函数不能是协程函数。

可将将协程的激活帧视为堆栈帧和协程帧，协程帧保留重新激活协程函数状态的信息，因此该帧在协程暂停时仍然存在；而堆栈帧在协程暂停时释放。

协程 [^[↗]^](https://lewissbaker.github.io/2017/09/25/coroutine-theory) 将函数 *Call* 和 *Return* 操作中的一些步骤拆分出一些额外的操作：*Suspend*、*Resume*、*Destroy*：
- *Suspend* 暂停协程当前点的执行并将执行返回调用方，且不会销毁协程的激活帧。*Suspend* 点由 `co_await` 和 `co_yield` 标识。当协程触及 *Suspend* 时将协程当前调用堆栈寄存器保存的任何值写入协程帧。保存一个值指示协程在哪个挂起点暂停，并由后续的 *Resume* 执行恢复这些值或 *Destroy* 执行销毁这些值。协程在暂停点返回调用方之前额外附加一个可以访问协程帧的句柄  *handle*，以用于在稍后恢复或销毁协程帧。
+ *Resume* 恢复暂停点协程的执行，并重新激活协程的激活帧。*resumer* 通过调用协程帧句柄 *handle* 的 `void resume()` 方法，重新分配一个堆栈帧、加载协程帧来跳转移到上次协程的挂起点。当下次挂起或运行完成时，调用 `resume()` 将返回并恢复调用函数的执行。
- *Destroy* 销毁激活帧，且停止继续执行协程，所有存储协程激活帧的内存将被释放。*Destroy* 重新激活协程激活帧但不会将执行转移到协程主体，它将执行转移至一个另一条代码路径，该路径在协程的挂起点调用范围内所有局部变量的析构函数，最后释放协程帧使用的内存。*Destroy* 由 `handle.destroy()` 调用。
+ co-*Return* 将返回值存储在某个位置（可由协程自定义），然后析构所有的局部变量。在返回调用方可以执行一些额外的操作例如执行某些操作发布返回值或恢复正在等待返回值的另一个协程（异步）。与常规调用返回值不同的时，协程的返回操作可能在调用方恢复执行后很久才被执行。

协程定义了 *Promsie* 和 *Awaitable* 两种接口。*Promise* 指定自定义协程本身行为的方法，例如自定义调用协程时发生和返回、异常时发生的情况，并自定义协程中任何 `co_await` 和 `co_yield` 表达式的行为。*Awaitable* 指定控制 `co_await` 语义的方法表达，例如是否暂停当前协程，在当前协程暂停后执行一些逻辑以安排协程稍后恢复，以及在协程恢复后执行一些逻辑以生成 `co_await` 表达式的结果。 
 
> **Awaitable** [^[↗]^](https://lewissbaker.github.io/2017/11/17/understanding-operator-co-await)

支持 `co_await` 运算的类型称为 *Awaitable* 类型。*Awaiter* 类型实现了作为 `co_await` 表达式的一部分调用的三个方法：

```cpp
struct Awaiter {
	bool await_ready() noexcept;
	void await_suspend(coroutine_handle<>) noexcept;  // return void or bool
	Result await_resume() noexcept;                   // return Result or void
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
#### 12.5. 运算符重载

运算符函数是具有特殊名称的函数。

```cpp
operator op     // op: + - * / % ^ & | ~ ! = < > += -= *= /= %= ^= &= |= 
				//     << >> >>= <<= == != <= >= <=> && || ++ -- , ->* -> () []
operator new | new[]
operetor delete | delete[]
operator co_await
```

仅 `()`,`[]` 成员函数可以声明为静态，声明其他静态的成员运算符仅能通过 `T::operaror op()` 调用。非成员的重载运算符操作数至少有一个具有类类型或枚举类型。`.`, `.*`, `::`, `? :`, `#`, `##` 不支持重载。

```cpp
struct S {
	int d;
	S operator + (int v) { return (d += v, *this); }
	S operator -(int v) { return (d -= v, *this); }
	S(int v) : d{ v } {}
};
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

将 `std::istream&` 或 `std::ostream&` 作为左侧参数的 `operator >>` 和 `operator <<` 的重载称为插入和提取运算符，必须将其实现为非成员或友元函数。

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
	// prefix increment
	Counter& operator++() { return (++count, *this); }
	// postfix increment
	Counter operator++(int)
	{
		Counter old = *this; // copy old value
		operator++();		 // prefix increment
		return old;          // return old value
	}

	// prefix decrement
	Counter& operator--() { return (--count, *this); }
	// postfix decrement
	Counter operator--(int)
	{
		Counter old = *this; // copy old value
		operator--();		 // prefix decrement
		return old;			 // return old value
	}
};
```

二元运算符通常实现为非成员，以保持对称性（*T + v*, *v + T*）；而成员函数仅支持 *T + v* 而不是 *v + T*；

>---
#### 12.6. 用户定义文本函数

用户定义文本可以是整数、浮点数、字符或字符串用户定义，用户定义文本由一个用户定义后缀标识。用户定义文本被视为对文本运算符或文本运算符模板的调用。

```c++
using namespace std;
long double operator""_w(long double);
string operator""_w(const char16_t *, size_t);
unsigned operator""_w(const char *);
int main()
{
    1.2_w;	  // calls operator ""_w(1.2L)
    u"one"_w; // calls operator ""_w(u"one", 3)
    12_w;	  // calls operator ""_w("12")
    "two"_w;  // error: no applicable literal operator
}
```

---
### 13. 类类型

类类型为值类型，包含 `class`、`struct`、`union`，其中联合体类型隐式密封。类成员可以声明静态过非静态的数据成员和成员函数、成员 *typedefs* 或 *usings*、成员枚举、嵌套类和友元声明。类内部定义的成员函数隐式内联，除非是命名模块导出。

```cpp
class S {
private:
    int d1;            		       // non-static data member
    int a[10] = {1, 2}; 	       // non-static data member with initializer (C++11)
    static const int d2 = 1;       // static data member with initializer
protected: 
    virtual void f1(int) = 0;      // pure virtual member function
    std::string d3, *d4, f2(int);  // two data members and a member function
public:
    enum { NORTH, SOUTH, EAST, WEST };
    struct NestedS {
        std::string s;
    } d5, *d6;
    typedef NestedS value_type, *pointer_type;
	static constexpr decltype(auto) str = "hello";
};
```

`using` 声明可以将基类成员引入派生类定义，或将枚举成员引入类作用域。

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

可以函数块声明局部类，局部类不包含静态数据成员，成员函数类内部定义且无链接，无友元模板和友元函数，非闭包类型的局部类没有成员模板。

```cpp
int main()
{
	std::vector<int> v{ 3,4,5,7,2,8,2,4 };
	struct Local {
		bool operator()(int n, int m) { return n > m; }
	};
	std::sort(v.begin(), v.end(), Local()); // 降序
	for (int n : v)
		std::cout << n << ' ';
	std::cout << '\n';
}
```

类的非静态数据成员可以具有默认初始化器，静态数据成员仅声明（除非是 `constexpr`,`const`,`inline`）。当且仅当命名空间域的类本身具有外部链接时，类的静态数据成员具有外部链接。

```cpp
struct S {
	int v = 10086;
	const int cv = v;
	static int sv;
	static const int scv = 10;
	static constexpr int d2 = 10;
};
int S::sv = 10010;
```

嵌入的匿名类不能是静态，不包含非公共成员、静态成员、任何成员函数、构造函数和析构函数、不能作为参数传递、无法作为函数中返回值返回。

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

表达式 `this` 是一个 纯右值 表达式，其值是隐式对象形参的地址。CV 成员函数的 `this` 也是 CV `T*`，且仅由对应的 CV 对象调用；*ref* 限定的成员函数只会限定隐式对象，不会限定 `this` 指针；构造函数和析构函数的 `this` 始终为 `T*`。可以执行 `delete this;`，应确保对象是由 new 分配的。在 `delete this;` 返回后，使指向已释放对象的每个指针都无效，包括 `this` 指针本身。

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

非虚非静态成员函数可以声明显式对象形参，允许推导类型和值类别，无 CV 和引用限定。指向显式对象成员函数的指针是普通的函数指针，而不是成员指针。

```cpp
struct S {
	int g(this const S&, int, int);
};
int main() {
	S s{};
	auto pf = &S::g;
	pf(s, 3, 4);	// ok
	(s.*pf)(3, 4);  // error: “pg” is not a pointer to member function
}
```

> *可变数据成员*

`mutable` 声明一个类的非静态、非常量和非引用的可变数据成员，它可以在 `const` 限定函数或 `const` lambda 中进行修改。 

```c++
class X
{
public:
    bool GetFlag() const
    {
        m_accessCount++;
        // m_flag = true;  // err
        return m_flag;
    }
private:
    bool m_flag;
    mutable int m_accessCount;
};
```

> **可访问性**

`class` 默认对其成员及其基类具有 `private`。`struct` 默认对其成员及其基类具有 `public`。`union` 默认对其成员具有 `public`。成员或基类可声明 `public`、`protected`、`private`。限定基类访问隐式限定继承的基类成员的访问。

```cpp
struct A {
	int a;
};
struct B : private A { 
	int b = a;   // a is private
};
struct C : public B {
	int c = a;  // err; 不可访问
};
```

> *注入类名* 

在C++中，​注入类名（injected class name）​是指类或类模板的作用域内，类名被视为该类的公共成员名，可以直接使用而无需通过作用域解析运算符（::）显式引用。注入类名是可继承的。非公共继承的间接基类的注入类名可能在派生类中变得不可访问。

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
#### 13.1. 构造函数

构造函数允许的说明符是 `friend`、`inline`、`constexpr`、`consteval` 和 `explicit`。构造函数可选具有成员初始化表达式列表，该列表会在构造函数主体运行之前初始化类成员。`constexpr` 构造函数是其类型成为字面类型。

在某些情况下，特殊成员函数即使用户未定义也会由编译器定义。它们是 *默认构造函数*、*复制构造函数*、*移动构造函数*、*复制赋值运算符*、*移动赋值运算符*、*析构函数*，默认 `inline public`，可以显式声明弃置以阻止默认生成。特殊成员函数以及比较运算符是唯一可以默认化（`= default`）的函数。


```cpp
struct Point {
    constexpr Point() :x{}, y{} {};  // 带有成员初始化表达式的默认构造函数
    explicit Point(int x, int y) :x(x), y(y) {}; // 重载 
    int x, y;
};
constexpr Point origin{};            // 字面类型
```

> *转换构造函数*

没有 `explicit` 的构造函数且至少带有一个非默认参数被称为转换构造函数。隐式声明和用户定义的非显式复制构造函数和移动构造函数都是转换构造函数。

```cpp
struct S {
	explicit S() {}             // 仅显式转换
	explicit (true) S(int) {}   // 支持隐式转换
	S(int, int) {}
};
int main() {
	S s1 = 1;          // err
	S s2 = (S)1;       // ok
	S s3(3);           // ok
	S s4{ 4, 4 };      // ok
	S s5 = { 5, 5 };   // ok
	S s6;              // ok;
	S s7 = {};		   // err
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
int main()
{
    double d0 = money(3.14);
    double d1 = money(3.14L);
    double d2 = money(10010);
    double d3 = money(2.71f);
}
```

> *委托构造函数*

构造函数可以在成员初始化列表中使用其他构造函数。委托构造函数无法递归，不包含其他成员初始化器。

```c++
struct Point {
    Point(int x, int y, int z) :x{ x }, y{ y }, z{ z } {};
    Point(int x, int y) : Point(x, y, 0) {};
    Point() : Point(0, 0, 0) {};
    int x, y, z;
};
```

> *继承构造函数*

`using Base::Base` 声明基类时，将基类的构造函数引入基类中。

```cpp
struct Base {
	int x, y;
	Base(int x, int y) : x{ x }, y{ y } {}
};
struct DerivedA : Base {
	Base b = DerivedA{ 1,2 };   // err
};
struct DerivedB : Base {
	using Base::Base;   // 继承 Base 的所有构造函数
	Base b = DerivedB { 1,2 };  // ok
};
```

> **复制构造函数**

复制构造函数通过从相同类型的对象复制成员值来初始化对象。成员中存在指针时，自动生成复制构造只会复制指针值，因此可能需要自定义声明以分配指针内存。复制构造函数的首元参数为 `T&`，可以具有 CV 限定；若包含其他参数，需要具有默认实参。

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
	friend size_t Println(buffer&);
	char* buf;
	size_t index, size;
public:
	buffer(size_t size) :size{ size }, buf{ new char[size] }, index{ 0 } {}
	buffer(const buffer& other) noexcept :size(other.size), index{ other.index }, buf{ nullptr } {
		cout << "copy buffer" << endl;
		buf = new char[size];
		copy(other.buf, other.buf + size, buf);
	};
	buffer& operator = (buffer& other)noexcept {
		if (this != &other) {
			delete[] buf;
			size = other.size;
			index = other.index;
			buf = new char[size];
			copy(other.buf, other.buf + size, buf);
		}
		return *this;
	}
	size_t Write(string str);
	operator string () { return string{ buf, index }; }
	~buffer() { delete[] buf; }
};
size_t Println(buffer& bfr) {
	std::printf("%s\n", string{ bfr }.c_str());
	size_t n = bfr.index;
	bfr.index = 0;
	bfr.buf[0] = '\0';
	return n;
}
size_t buffer::Write(string str) {
	size_t n = str.length();
	size_t able = size - index - 1;
	if (str.length() >= able)
		n = able;
	copy(str.c_str(), str.c_str() + n, buf + index);
	index += n;
	buf[index] = '\n';
	return n;
}
int main() {
	buffer b{ 512 };
	b.Write("Hello World!\n");
	buffer b2 = b;
	b2.Write("JimryYchao\n");
	Println(b2);
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
    Println(b3);
}
```

>---
#### 13.2. 析构函数

析构函数在对象超出范围或通过调用 `delete` 显式销毁对象时自动调用。析构函数不可继承，如果基类将其析构函数声明为 `virtual`，则派生析构函数始终会重写它。这使得可以通过指向基类的指针删除多态类型的动态分配对象。

纯虚析构函数强制类型为抽象类，需要提供纯虚析构函数的定义。虚析构函数主要用于确保多态场景下派生类对象的正确析构，避免内存泄漏和资源未释放的问题。当通过基类指针删除派生类对象时，虚析构函数确保调用派生类的析构函数再调用基类的析构函数；要求其继承链的最终基类必须是虚析构函数。


```cpp
struct BaseA {
	virtual ~BaseA() { cout << "~BaseA\n"; }
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
	delete obj; // 正确调用 Derived::~Derived() → BaseB::~BaseB() → BaseA::~BaseA()
	// 若基类析构函数非虚，delete obj 仅调用 BaseA::~BaseA()，导致派生类资源泄漏
}
```

>---
#### 13.3. default, delete

`default` 函数和 `delete` 函数可以显式控制是否自动生成特殊成员函数。`delete` 函数防止所有类型的函数（特殊成员函数和普通成员函数以及非成员函数）的自变量中出现有问题的类型提升，这可能会导致意外的函数调用。

在 C++ 中，如果某个类型未声明它本身，则编译器将自动为该类型生成默认构造函数、复制构造函数、复制赋值运算符和析构函数、移动构造函数和移动赋值运算符。这些函数称为特殊成员函数，它们使 C++ 中的简单用户定义类型的行为如同 C 中的结构。其中：
- 显式声明了任何构造函数，则不会自动生成默认构造函数。
- 显式声明了虚拟析构函数，则不会自动生成默认析构函数。
- 显式声明了移动构造函数或移动赋值运算符，则不自动生成复制构造函数和成复制赋值运算符。
- 显式声明了复制构造函数、复制赋值运算符、移动构造函数、移动赋值运算符或析构函数，则不自动生成移动构造函数和移动赋值运算符。
- 显式声明了复制构造函数或析构函数，则弃用复制赋值运算符的自动生成。
- 显式声明了复制赋值运算符或析构函数，则弃用复制构造函数的自动生成。

如果基类不拥有可派生类调用的默认构造函数，例如没有 `public` 或 `protected` 的默认构造函数，那么派生类无法自动生成它自己的默认构造函数。

对于创建不可移动、只能动态分配或无法动态分配的用户定义类型；可以通过 `default` 和 `delete` 方式进行设定。弃置声明复制构造函数和复制赋值运算符，可以使用户定义类型不可复制：

```c++
struct noncopyable
{
	noncopyable() = default;   
	noncopyable(const noncopyable&) = delete;
  	noncopyable& operator=(const noncopyable&) = delete; 
};
```

> *defualt* *特殊成员函数*

可以默认设置任何特殊成员函数，以显式声明特殊成员函数使用默认实现、定义具有非公共访问限定符的特殊成员函数或恢复其他情况下被阻止其自动生成的特殊成员函数。通过对可内联的特殊成员函数设置为 `default` 而不是空函数体进行实现：

```c++
struct widget
{
  widget()=default;
  inline widget& operator=(const widget&);
};
inline widget& widget::operator=(const widget&) =default;
```

> *delete*

可以删除特殊成员函数和普通成员函数以及非成员函数，以阻止定义或调用它们。删除特殊成员函数，可以阻止编译器生成不需要的特殊成员函数。

```c++
struct widget
{
    // deleted operator new prevents widget from being dynamically allocated.
    void* operator new(std::size_t) = delete;
    widget* operator &() = delete;   // address-of 被删除
};
widget w{};
widget* pw = new widget;  // ERR; 无法调用删除的函数；但是可以调用全局 new ；
widget* pw = ::new widget;  // ok
widget* pw = &w  // ERR; address-of 被删除
```

删除普通成员函数或非成员函数可阻止有问题的类型提升导致调用意外函数。这可发挥作用的原因是，已删除的函数仍参与重载决策，并提供比提升类型之后可能调用的函数更好的匹配。函数调用将解析为更具体的但可删除的函数，并会导致编译器错误。

```c++
// deleted overload prevents call through type promotion of float to double from succeeding.
void call_with_true_double_only(float) =delete;
void call_with_true_double_only(double param) { return; }

call_with_true_double_only(3.14f);  // err
call_with_true_double_only(100);  // ok; int cast to double 
```

若要限制发生隐式类型转换，确保仅发生对 `double` 类型的参数进行调用，可声明一个模板的已删除版本：

```c++
template < typename T >
void call_with_true_double_only(T) =delete; //prevent call through type promotion of any T to double from succeeding.
void call_with_true_double_only(double param) { return; } // also define for const double, double&, etc. as needed.

call_with_true_double_only(3.1415);  // just only double
```

>---
#### 13.4. 继承

继承从现有类派生新类；可以是单一继承，或多重继承，`final` 声明无法被继承。同名签名不同的函数声明隐藏基类成员。	

```c++
class Derived : Base-Class, ... ;
Base-Class = [virtual] [access-specifier] BaseClass
```

在多重继承中，可以构建一个继承关系图，其中相同的基类是多个派生类的一部分。多重继承使得沿多个路径继承名称成为可能。沿这些路径的类成员名称不一定是唯一的。由于一个类可能多次成为派生类的间接基类，这些名称冲突存在 “多义性”。任何引用类成员的表达式必须采用明确的引用，可以通过限定名称及其类名消除多义性。

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
    pc->b();     // b 不明确
    pb->B::b();  // ok
}
```

通过一个继承关系图到达多个名称（函数、对象或枚举器）是可能的。这种情况被视为与非虚拟基类一起使用时目的不明确，从派生类型指针或引用转换到基类指针或引用的显式或隐式转换可能会导致多义性。 

```c++
class A { };
class B : public A {};
class C : public A {};
class D : public B, public C {};

A* pa = new D;  // 从 D 到 A 的转换是不明确的；A 无法辨别从 B 还是 C 传递
// 需要显式指定要使用的子对象
A* pab = (A*)(B*)(new D);
A* pac = (A*)(C*)(new D);
```

`virtual` 基类可以避免多重继承的类层次结构中出现多义性。虚拟基类的数据成员在多态继承路径上拥有唯一副本，由基类和派生类共享。

```c++
class A {  };
class B : public virtual A {};
class C : public virtual A {};
class D : public B, public C {};
A* pa = new D; 
```

虚函数可以在派生类中重新定义（`override` 和 `virtual` 重写时可选），纯虚函数强制类型为抽象类型。更改重写函数的访问限定不会影响多态行为。`consteval` 虚函数不得重写非 `consteval` 虚函数，也不得被其重写。使用限定函数标识（`Base::VirFunc`）调用基类虚函数。

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
   // Create objects of type CheckingAccount and SavingsAccount.
   CheckingAccount checking( 100.00 );
   SavingsAccount  savings( 1000.00 );
   // Call PrintBalance using a pointer to Account.
   Account *pAccount = &checking;
   pAccount->PrintBalance();  // call checking.PrintBalance
   // Call PrintBalance using a pointer to Account.
   pAccount = &savings;
   pAccount->PrintBalance();  // call savings.PrintBalance
}
```

可以提供纯虚函数的定义（纯虚析构函数必须提供），除纯虚析构函数外，其他纯虚函数的定义必须在类外部提供，派生类的成员函数可以自由地使用限定函数标识调用抽象基类的纯虚函数。

```cpp
struct Abstract {
	virtual void f() = 0; // pure virtual
	virtual void g() {}   // non-pure virtual
	~Abstract() {
		g();           // OK: calls Abstract::g()
		// f();        // undefined behavior
		Abstract::f(); // OK: non-virtual call
	}
};
// definition of the pure virtual function
void Abstract::f() { std::cout << "A::f()\n"; }
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


如果使用限定名查找选择函数（即，如果函数名称出现在作用域解析运算符 :: 的右侧），则会抑制虚函数调用。

```cpp
struct Base {
	virtual void f() { std::cout << "base\n"; }
};
struct Derived : Base {
	void f() override { std::cout << "derived\n"; }
};
int main()
{
	Base b;
	Derived d;
	// virtual function call through reference
	Base& br = b; // the type of br is Base&
	Base& dr = d; // the type of dr is Base& as well
	br.f(); // prints "base"
	dr.f(); // prints "derived"
	// non-virtual function call
	br.Base::f(); // prints "base"
	dr.Base::f(); // prints "base"
}
```

函数 `Derived::f` 重写了函数 `Base::f`，它们的返回类型必须相同或为协变的。

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
int main() {
	Derived d;
	Base& br = d;
	Base* p = br.P();     // calls Derived::P() and converts the result to B*
	Base& l = br.LR();    // calls Derived::LR() and converts the result to B&
	Base&& r = br.RR();   // calls Derived::RR() and converts the result to B&&
}
```


>---
#### 13.5. 友元声明

友元声明出现在类体中，并授予函数或另一个类访问声明友元声明的类的私有和受保护成员的权限。友元函数非成员函数并导出到封闭非类范围。友元关系不具有继承和传递性。访问说明符对友元声明的含义没有影响，友元类声明不能定义新类（例如 `friend class X {};` 是错误的）

```cpp
class Product {
private:
	int secretCode;
	Product(int code) : secretCode{code} {} 
	friend class Factory;     // 友元声明，类声明
};
class Factory {      // 定义
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
#### 13.6. 位域

类和结构可包含比整型类型占用更少存储空间的位域成员，它必须是整数类型或枚举。匿名位域用于填充；宽度为 0 的匿名位域强制将下一个位域与下一个类型边界对齐。

```c++
struct Date {
   unsigned nWeekDay  : 3;    // 0..7   (3 bits)
   unsigned nMonthDay : 6;    // 0..31  (6 bits)
   unsigned           : 0;    // Force alignment to next boundary.
   unsigned nMonth    : 5;    // 0..12  (5 bits)
   unsigned nYear     : 8;    // 0..100 (8 bits)
};
```

无法创建指向位域的指针和非常量引用。当从位域初始化常量引用时，会创建一个临时对象（其类型是位域的类型），用位域的值进行复制初始化，并且引用绑定到该临时对象。

```c++
Date d{ 1,1,1,1 };
const auto& r = d.nMonth;       // 临时对象
unsigned int& pr = const_cast<unsigned int&> (r);
pr = 2;
cout << pr << endl;				// 2
cout << d.nMonth << endl;		// 1
cout << r << endl;				// 2
```

>---
#### 13.7. 联合体

联合体隐式密封无基类，对象仅能保存其非静态数据成员中的一个，不包含引用类型非静态数据成员，可以包含位域。成员默认 `public`。联合体的大小至少与容纳其最大的数据成员所需的大小相同。最多一个变体成员可以具有默认初始化器。

```cpp
union U {
	std::int32_t n;     // occupies 4 bytes
	std::uint16_t s[2]; // occupies 4 bytes
	std::uint8_t c;     // occupies 1 byte
	unsigned int i1 : 4 = 1;   
};
```

联合具有非静态类类型成员时，编译器自动将非用户提供的任何特殊成员函数标记为 `delete`。如果联合是 `class` 或 `struct` 中的匿名联合，则 `class` 或 `struct` 的非用户提供的任何特殊成员函数都会被标记为 `delete`。如果联合体的成员是具有用户定义构造函数和析构函数的类，当切换活动成员时，通常需要显式析构函数和 *placement new*。

```cpp
union S {
	std::string str;
	std::vector<int> vec;
	~S() { }
};          
int main() {
	S s = { "Hello, world" };
	// reading from s.vec is undefined behavior
	std::cout << "s.str = " << s.str << '\n';
	s.str.~basic_string();   // 显式调用析构函数
	new (&s.vec) std::vector<int>;
	// now, s.vec is the active member of the union
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
### 14. 模板

模板可以定义类模板、函数模板、别名模板、变量模板、概念约束等实体，以若干的 *模板参数*（类型、常量或其他模板）参数化。通过提供具体模板实参以特化，特化支持显式提供：对类、函数、变量模板全特化，或对类和变量模板部分特化。模板声明可以包含约束。变量模板在类型中仅支持静态。别名模版始终不进行推导。

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

	template <C T>         // 静态数据成员模板
	static T Value = T{};		  

	void Func(C auto c);   // 成员函数模板
};

template <typename... T>   // 别名模板
using value_holder_R = value_holder<T...>&;    

template<class T>          // 变量模板
constexpr T pi = T(3.1415926535897932385L); 
```

> 简化函数模板

当函数声明或函数模板声明的形参列表中出现占位符类型（`auto` 或 `Concept auto`）时，该声明为一个函数模板，并且每个占位符向模板形参列表追加一个虚设的模板形参：

```cpp
template <typename> concept T = true;
template <typename> concept U = true;

void f1(auto);					 // template<class T> void f(T)
void f2(T auto);				 // template<T X> void f2(X)
void f3(U auto...);				 // template<U... Ts> void f3(Ts...) 
void f4(const T auto*, U auto&); // template<T X, U Y> void f4(const X*, Y&);

template<class X, T Y>
void g(X x, Y y, U auto z);      // template<class X, T Y, U Z> void g(X x, Y y, Z z);
```

>---
#### 14.1. **模板形参**

```cpp
<Params> = 
// 常量模板形参
    <Typeid [name] [= Val]>      // 默认值
    <Typeid ... [name]>          // 常量形参包
// 类型模板形参
    <class | typename | ConceptName [name] [= Val]> 
    <class | typename | ConceptName ... [name]>   // 模板形参包
// 模板模板形参
    <template <Params> class | typename [name] [= Val]>
    <template <Params> class | typename ... [name]>
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
Sample<int> s{};   // Sample<int,int>

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
Tuple<> t0;           // Types<>
Tuple<int> t1;        // Types<int>
Tuple<int, float> t2; // Types<int,float>
// 变参函数模板
template<class... Types>
void Func(Types... args);   
int main() {
	Func();			  // Func<>
	Func(1);		  // Func<int>
	Func(1, 2.0);     // Func<int,double>
}
```

包模式展开：形参包被展开成零个或更多个逗号分隔的模式实例，每个实例元素按顺序被替换成包中的各个元素。同一模式中出现的两个形参包长度必须相同。嵌套包模式展开从最内层最先开始。`size...(Pack)` 返回形参包长度。

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
    int res[size] = {1, args..., 2};
    // 因为初始化器列表保证顺序，所以这可以用来对包的每个元素按顺序调用函数：
    int dummy[sizeof...(Ts)] = {(std::cout << args, 0)...};
}
```

包展开可以在模板形参列表中出现、或用于指定类声明的基类列表、可以在 lambda 捕获子句中出现。

```cpp
template<class... Mixins>
class X : public Mixins... {
public:
    X(const Mixins&... mixins) : Mixins(mixins)... {}  // 同时需要在构造函数中使用包展开
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


包还可以由折叠表达式展开。`op` 支持 `+`, `-`, `*`, `/`, `%`, `^`, `&`, `|`, `=`, `<`, `>`, `<<`, `>>`, `+=`, `-=`, `*=`, `/=`, `%=`, `^=`, `&=` ,`|=`, `<<=`, `>>=`, `==`, `!=`, `<=`, `>=`, `&&`, `||`, `,`, `.*`, `->*`。其中对于一元展开空包，`&&` 返回 `true`，`||` 返回 `false`，`,` 返回 `void()`。

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
#### 14.2. 模板实例化

模板的某个特定类型尚未被显式实例化时，由编译器隐式实例化。显式实例化定义强制实例化其所指代的模板类型，要求在此之前必须存在它的模板定义。显式实例化声明（`extern`）跳过隐式实例化步骤，表明该模板已在别处显式实例化定义。若同一组模板实参的显式特化在显式实例化定义之前出现，则该显式实例化无效（反之定义在特化之前声明，程序非良构）。

```cpp
/* fileA.cpp */
template<typename T>
void Func(T v) { std::cout << v << std::endl; };   // 模板定义

template<> void Func<bool>(bool v) {     // Func<bool> 显式特化
	std::cout << std::boolalpha << v << std::endl;
};
template void Func<bool>(bool v);            // 无效显式实例化定义
template void Func<double>(double v);        // Func<double> 显式实例化定义
//template<> void Func<double>(double v) {}  // err 显式特化之前已有显式定义

/* fileB.cpp */
extern template void Func<bool>(bool);       // 引用外部显式特化定义
extern template void Func<double>(double);   // 引用外部显式实例化定义
extern template void Func<string>(string);   // 无效显式实例化引用声明，无显式实例化定义
int main() {
	Func<int>(10086);    // 隐式实例化
	Func<double>(3.14);
	Func(1 > 2);         // false
	// Func<string>("hello");   // err, 未定义
}    
```

>---
#### 14.3. 模板特化   

全特化：允许对给定的模板实参集定制模板代码，显式特化可以在它的主模板的作用域中声明。

```cpp
// 显式全特化
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

在模板定义内，某些名称被推导为属于某个未知特化，特别是，
- 在 :: 左侧出现了并非当前实例化成员的待决类型的名字的限定名。
- 限定符是当前实例化，且无法在当前实例化或它的任何非待决基类中找到，并存在待决基类的限定名。
- 类成员访问表达式中的成员名（`x.y` 或 `xp->y` 中的 `y`），如果对象表达式（`x` 或 `*xp`）的类型是待决类型且非当前实例化。
- 类成员访问表达式中的成员名（`x.y` 或 `xp->y` 中的 `y`），如果对象表达式（`x` 或 `*xp`）的类型是当前实例化，且在当前实例化或任何其非待决基类中找不到该名字，并存在待决基类。

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

在模板（包括别名模版）的声明或定义中，不是当前实例化的成员且取决于某个模板形参的名字不会被认为是类型，除非使用 `typename` 或它已经被设立为类型名（例如用 `typedef` 声明或通过用作基类名）。

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
#### 14.4. 用户定义模板推导

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
struct container
{
	container(T t) {}
	template<class Iter> container(Iter beg, Iter end);
};
template<class Iter>  // 额外的推导指引
container(Iter b, Iter e) -> container<typename std::iterator_traits<Iter>::value_type>;
vector v = { 1,2,3,4 };
container c{ v.begin(), v.end() };   // T is int
```

>---
#### 14.5. SFINAE

“替换失败不是错误” (Substitution Failure Is Not An Error)：当模板形参在替换成显式指定的类型或推导出的类型失败时，从重载集中丢弃这个特化，而非导致编译失败。

试图创建 void 的数组，引用的数组，函数的数组，负大小的数组，非整型大小的数组，或者零大小的数组。

```cpp
template<int I>
void div(char(*)[I % 2 == 0] = 0) {
    // 当 I 是偶数时选择这个重载
}
template<int I>
void div(char(*)[I % 2 == 1] = 0) {
    // 当 I 是奇数时选择这个重载
}
```

试图在作用域解析运算符 `::` 左侧使用类和枚举以外的类型。

```cpp
template<class T>
int f(typename T::B*);
template<class T>
int f(T);

int i = f<int>(0); // 使用第二个重载
```

>---
#### 14.6. 概念与约束 

类模板、函数模板（包括泛型 lambda）等可以与一项约束相关联。这类约束要求（`requires`）的具名集合称为概念。

```cpp
template <params>
concept ConceptName [Attr] = ConstraintExpr;
    ConstraintExpr = RequiresExpr | RequiresClause
```

概念可以作为模板类型形参声明、占位类型说明符或 `requires` 的复合要求。

```cpp
template<class T, class U>
concept Derived = std::is_base_of<U, T>::value;
 
template<Derived<Base> T>
void f(T); // 隐式推断：T 被 Derived<T, Base> 约束
```

> **要求**

概念是要求的具名集合。要求由 `requires` 表达式定义，编译时计算。任意一项不符合 *requirements* 时，`requires` 返回 `false`。

```cpp
// requires 表达式
RequiresExpr = requires [( Params )] { 
    requirements = Sample | Type | Compound | Nested
        Sample   : expr;
        Type     : typename Identifier;
        Compound : { expr } [noexcept] [-> TypeConstraint | decltype(...)];
        nested   : RequiresExpr | RequiresClause;
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

Type：类型要求断言 *标识符*（可有限定）知名的类型是有效的：嵌套类型存在、类或别名模板特化有效。

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
template <C T> T::inner F(T& t) { return *t; }

struct Sample {
	int value;
	typedef Sample* inner;
	inner operator* () { return this; }
	int operator +(int v) noexcept { return value += v; }
	Sample operator *(int v) { return (this->value *= v, *this); }
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
concept Semiregular = Addable<T> &&
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
### 15. 异常处理

异常包括有编程逻辑错误和运行时错误。错误报告的管理方式是返回一个错误代码或设置一个全局变量，调用方选择性地检索该变量状态。例如 C `errno`。

`try-catch` 执行异常处理，未捕获时调用 `std::terminate` 终止程序。`throw` 可以抛出任何类型，但应引发直接或间接派生自 `std::exception` 的类型。

对于每个可能引发或传播异常的函数，应提供三项异常保证之一：强保证、基本保证或 nothrow `noexcept` 保证。通过 `throw` 值传递引发异常，通过 `catch` 引用传递捕获异常。

```cpp
template <typename T>
T Index(T arr[10], int i) throw(out_of_range) {  // 强保证
    if (i < 0 || i >= 10)
        throw std::out_of_range("index is out of range");   // 抛出值
    return arr[i];
}

int CatchIndex(int arr[]) noexcept {   // 无异常保证

    try {
        auto v = Index(arr, -1);
    } catch (out_of_range& e) {          // 引用捕获
        cout << e.what() << endl;
        return -1;
    } catch (exception& e) {
        return -10;
    } catch (...) {
        return -100;
    }
}

int main() {
    int arr[10] = { 0,1,2,3,4,5,6,7,8,9 };
    return CatchIndex(arr);
}
```

>---
#### 15.1. try-catch, throw

`try` 执行受保护代码。`throw` 引发异常或 *rethrow*。`catch` 执行异常处理程序。`catch(...)` 程序块处理任意类型的异常。(MSVC 编译器) 当使用 `/EHa` 选项编译时，异常可包括 C 结构化异常和系统或应用程序生成的异步异常，例如内存保护、被零除和浮点冲突。

谨慎使用 `catch(...)`；除非 `catch` 明确知道特定异常。`catch(...)` 块一般用于在程序停止执行前记录错误和执行特殊的清理工作。在 `catch` 块中 `throw;` 语句重新引发当前正在处理的异常。*rethrow* 的异常对象是原始异常对象。

> *堆栈展开示例*

```c++
// MyProgram.cpp
import Example;
import std;
using namespace std;

class MyException :exception {
public:
    MyException(const char* const mess) {
        _mess = mess;
    };
    const char* const What() {
        return _mess;
    }
private: const char* _mess;
};
class Dummy {
public:
    Dummy(string s) : MyName(s) { printMsg("Created Dummy:"); }
    Dummy(const Dummy& other) : MyName(other.MyName) { printMsg("Copy created Dummy:"); }
    string MyName;
    ~Dummy() { printMsg("Destroyed Dummy:"); }
private:
    void printMsg(string s) { cout << s << MyName << endl; }
};
void B(Dummy d, int i) {    
    d.MyName = " B";
    throw MyException("Throw in B");
    cout << "Exiting FunctionB" << endl;
}
void A(Dummy d, int i) {
    d.MyName = " A";
    B(d, i + 1);
    cout << "Exiting FunctionA" << endl;
}
int main() {
    try {
        Dummy d(" M");
        A(d, 1);
    }
    catch (MyException& e) {
        cout << "Caught:" << e.What() << endl;
    }
    catch (exception) {
        throw;  // 重新引发
    }
    cout << "Exiting main." << endl;
}
/* Output:
    Created Dummy: M
    Copy created Dummy: M
    Copy created Dummy: A
    Destroyed Dummy: B
    Destroyed Dummy: A
    Destroyed Dummy: M
    Caught:Throw in B
    Exiting main.
*/
```


>---
#### 15.2. 未处理的异常

如果无法找到当前异常的匹配处理程序（或 `catch(...)` 处理程序），则调用预定义的 `terminate` 运行时函数（默认操作是 `abort()`）。可以在程序的任何点调用 `set_terminate()`。`terminate` 总是调用最后一次指定给 `set_terminate` 参数。

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
   return 0;
}
```

>---
#### 15.3. 函数 try 块

```cpp
function try {
    // 函数体
} catch ( /*...*/ ) {
    // 处理异常
} 
constructor try [: initer-list ] 
{ ... } catch { ... } 
```

函数 `try` 块是一类特殊的函数体。若 `try` 块或构造函数初始化器抛出异常，则在函数 `try` 块后的 `catch` 处理块中处理异常。对于构造函数或析构函数 `try` 块，始终默认重新抛出异常；在这类函数 `try` 块中涉及该对象的非静态成员或基类会导致未定义行为。

```cpp
struct Sample {
	int value;
	Sample(int value) try : value{ value } {
		throw this->value;
	}
	catch (exception& e) {
		cout << e.what();
	}
	catch (...) {
		cout << "catch unknown exception\n";
	}
};

int main(int, const char* argv[]) try
{
	Sample s{ 10086 };
}
catch (int v) {
	return v;  // 10086
}
```

>---
#### 15.4. 异常规范与 noexcept

异常规范指示可由函数传播的异常类型的意图，表明指定函数可以或不可以因异常退出。`noexcept` 指定可以脱离函数的潜在异常集是否为空；`throw()` 是 `noexcept(true)` 的别名。

| 异常规范                                    | 含义                                                                                                                                                                                               |
| :------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `noexcept`<br>`noexcept(true)`<br>`throw()` | 函数不会引发异常，也不允许在其范围外传播异常。`noexcept` 和 `noexcept(true)` 是等效的。此规范声明的函数引发异常时，将直接调用 `std::terinate` 终止程序。并且不会保证将调用任何范围内对象的析构函数 |
| `noexcept(false)`<br>`throw(...)`<br>无规范 | 函数可以引发任何类型的异常。                                                                                                                                                                       |
| `throw(type)`                               | C++14 之前表示函数可以引发 `type` 类型的异常，之后编译器将其解释为 `noexcept(false)`                                                                                                               |

当要复制的对象是普通的旧数据类型 (POD) 时，可将复制其自变量的函数模板声明为 `noexcept`。

```c++
#include <type_traits>
template <typename T>
T copy_object(const T& obj) noexcept(std::is_pod<T>)
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

若 `expr` 具有类类型或它的数组类型，要求类型的析构函数可访问且未 `delete`。

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
<!-- ### 附录

#### 具名要求 -->