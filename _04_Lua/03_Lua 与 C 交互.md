## Lua 与 C 交互

### C API

C API 用于读写 Lua 全局变量、调用 Lua 函数、运行 Lua 代码段、注册 C 函数等。Lua 标准库将所有状态保存在动态结构体 `lua_State`。

Lua 和 C 通过虚拟栈 stack 通信和数据交换。常见的栈操作：
- 压入元素 `lua_push*`。`lua_checkstack` 检查栈空闲空间；栈中元素按先后压入顺序从 1 开始索引，-1 表示栈顶元素。
- 查询元素 `lua_is*`。这类函数检查栈中某个元素是否可以转换为特定类型。`lua_type` 返回栈中元素的类型例如 `LUA_TNIL`、`LUA_TBOOLEAN`、`LUA_TNUMBER`、`LUA_TSTRING` 等。
- 获取元素 `lua_to*`。转换失败不会提示类型错误，`lua_tolstring` 和 `lua_tothread` 返回 `NULL`。

> *一个简易 Lua 解释器*
 
```c
#include <stdio.h>
#include "lua.h"                // CAPI 基础函数
#include "lauxlib.h"            // CAPI 辅助库函数
#include "lualib.h"             // Lua 标准库函数

int DoLuaBuff(char * buf, int bufsz)
{
    int error;
    lua_State *L = luaL_newstate();    // 创建一个 Lua 状态机
    luaL_openlibs(L);                  // 加载标准库
    while(fgets(buf, sizeof(bufsz), stdin) != NULL){
        error = luaL_loadstring(L,buf) || lua_pcall(L,0,0,0);   // 发生错误时将信息压入栈中
        if (error){
            fprintf_s(stderr, "%s\n", lua_tostring(L,-1));      // 获取错误信息
            lua_pop(L,1);              // 从 lua state 栈中弹出错误信息
        }
    }
    lua_close(L);
    return 0;
}
``` 

> *栈操作*

```c
// 遍历栈元素
void StackDump(lua_State *L)
{
    int top = lua_gettop(L);        // 获取栈的深度
    for (int i = 1; i <= top; i++)
    {
        int t = lua_type(L, i);
        switch (t)
        {
        case LUA_TSTRING:
            printf("string: '%s'", lua_tostring(L, i));
            break;
        case LUA_TBOOLEAN:
            printf("boolean: %s", (lua_toboolean(L, i) ? "true" : "false"));
            break;
        case LUA_TNUMBER:
            if (lua_isinteger(L, i))
                printf("integer: %lld", lua_tointeger(L, i));
            else
                printf("double: %g", lua_tonumber(L, i));
            break;
        default:
            printf("other: %s", lua_typename(L, t));
            break;
        }
    }
}
```

> *错误处理*

Lua 使用异常提示错误，C API 使用 `setjmp` 和 `longjmp` 模拟异常处理机制。也可以运行 `lua_pcall` 在保护模式中运行 C 代码。`lua_error` 或 `luaL_error` 

```c
static int F(lua_State *L)
{
    // code to run in protected mode
    return 0;
}
int secure_F(lua_State *L)
{
    lua_pushcfunction(L, F);
    int error = lua_pcall(L, 0, 0, 0);
    if (error)
        return lua_error(L);
    return 0;
}
```

> *内存分配*

```c
lua_State* luaL_newstate();                         // 使用默认 l_alloc 分配函数创建 luaState
lua_State *lua_newstate (lua_Alloc f, void *ud);    // 原始分配函数
typedef void * (*lua_Alloc) (       // 内存分配函数类型
    void *ud,                       // lua_newstate 不透明指针 ud
    void *ptr,                      // 指向正在分配/重新分配/释放的区块的指针
    size_t osize,                   // 块的原始大小或有关所分配内容的代码如 
                                        // LUA_TSTRING、LUA_TTABLE、LUA_TFUNCTION、LUA_TUSERDATA
                                        // 或 LUA_TTHREAD 的新对象 
    size_t nsize);                  // 块的新大小；0 类似于 free, 非 0 类似于 realloc

// 默认分配函数
static void *l_alloc(void *ud, void *ptr, size_t osize, size_t nsize)
{
    (void)ud;(void)osize; /* not used */
    if (nsize == 0) {
        free(ptr);
        return NULL;
    }
    else
        return realloc(ptr, nsize);
}

lua_Alloc lua_getallocf(lua_State *L, void **ud);           // 返回 L 的分配函数和不透明指针 ud
void lua_setallocf (lua_State *L, lua_Alloc f, void *ud);   // 更改 L 的分配函数和 ud，新分配函数有责任释放由前一个分配函数分配的块
```


>---
### 交互案例
#### C 调用 Lua

```lua
-- luaSrc.lua
function Add(x, y)
    return x + y
end
```
```c
void test_C_Call_Lua() {
	// load lua source file
	lua_State* L = luaL_newstate();
	if (luaL_loadfile(L, "luaSrc.lua") || lua_pcall(L, 0, 0, 0))
		luaL_error(L, "cannot load lua file: %s\n", lua_tostring(L, -1));

	// push lua function 
	lua_getglobal(L, "Add");
	lua_pushnumber(L, 3.14);
	lua_pushnumber(L, 9.81);

	// call Add and get result
	if (lua_pcall(L, 2, 1, 0) != LUA_OK)
		luaL_error(L, "error running function 'Add':%s", lua_tostring(L, -1));
	int isnum;
	double rt = lua_tonumberx(L, -1, &isnum);
	if (!isnum)
		luaL_error(L, "function 'Add' should return a integer");
	lua_pop(L, 1);

	// print result and close lua_State
	printf("in Lua, %g + %g = %g\n", 3.14, 9.81, rt);
	lua_close(L);
}
```

>---
#### Lua 调用 C

注册到 Lua 的 C 函数具有 `typedef int (*lua_CFunction)(lua_State *L)` 原型。在 Lua 中调用 C 函数，首先通过 `lua_pushcfunction` 注册该函数。在 C 函数返回后自动保存返回值并清空整个栈。lua 参数从左向右依次压入栈。

```c
// C 函数格式
int CFunction(lua_State *L) {
	// do with stack
	return rtValue_count;
}
```

```c
static int c_sin(lua_State* L) {
	int isnum;
	double d = lua_tonumberx(L, 1, &isnum);
	if (!isnum)
		lua_error(L);
	lua_pushnumber(L, sin(d));
	return 1;
}

void test_Lua_Call_C() {
	lua_State* L = luaL_newstate();
	luaL_openlibs(L);

	lua_pushcfunction(L, c_sin);
	lua_setglobal(L, "Sin");   
	// or lua_register(L, "Sin", c_sin);
	luaL_dostring(L, "return Sin(math.pi  / 6)");

	double rt = lua_tonumber(L, -1);
	printf("call c_sin: rt = %g\n", rt);
    lua_pop(L, 1);
	lua_close(L);
}
```

>---
#### C 函数延续

Lua 5.4 之前的 C 函数调用链中触发 *yield* 会因 longjmp 机制破坏 C 栈完整性而抛出 "attempt to yield across a C-call boundary"（例如 `pcall`）。使用延续实现 `pcall`：

```c
// lua < 5.4
// 恢复 C 函数调用时跳转至 finishpcall
static int finishpcall(lua_State* L, int status, intptr_t ctx) {
	printf("Call finishpcall\n");
	(void)ctx;
	status = (status != LUA_OK && status != LUA_YIELD);
	lua_pushboolean(L, (status == 0));
	lua_insert(L, 1);
	return lua_gettop(L);
}

static int luaB_pcall(lua_State* L) {
	printf("Call luaB_pcall\n");
	int status;
7	luaL_checkany(L, 1);
	status = lua_pcallk(L, lua_gettop(L) - 1, LUA_MULTRET, 0, 0, finishpcall);
	return finishpcall(L, status, 0);
}

void test_Coroutine_CFunc() {
	lua_State* L = luaL_newstate();
	luaL_openlibs(L);
	lua_pushcfunction(L, luaB_pcall);
	lua_setglobal(L, "pcall");

	luaL_dostring(L,
		"co = coroutine.create(function()			\
			print('yield co')						\
			print(pcall(coroutine.yield))			\
			print('resume co')						\
			end)									\
		coroutine.resume(co)");

	luaL_dostring(L, "coroutine.resume(co)");
	lua_close(L);
}
/*
yield co
Call luaB_pcall
Call finishpcall
true
resume co
*/
```

5.4 版本引入 `lua_KFunction` 作为协程回调函数，允许 C 函数在挂起后通过延续 *continuation* 恢复执行。当 C 函数调用 `lua_yieldk` 挂起时，可指定一个 `lua_KFunction` 作为恢复点，恢复时直接跳转到该函数，从而避免重新遍历调用栈。

```c
static int counter;
static void _C_DeCounter(lua_State* L, int n);

static int C_DeCounter_yield(lua_State* L, int status, lua_KContext ctx) {
	_C_DeCounter(L, counter);
}

static void _C_DeCounter(lua_State* L, int n) {
	if (!(n < 0)) {
		lua_pushinteger(L, n);
		counter = n - 1;
		lua_yieldk(L, 1, 0, C_DeCounter_yield);
	}
	else
		printf("counter end\n");
}
7
int C_DeCounter(lua_State* L) {
	counter = lua_tointeger(L, -1);
	lua_pop(L, 1);
	if (counter <= 0)
		return;
	_C_DeCounter(L, counter);
	return 0;
}

void test_Continue_CFunc() {
	lua_State* L = luaL_newstate();
	luaL_openlibs(L);
	lua_pushcfunction(L, C_DeCounter);
	lua_setglobal(L, "counter");

	luaL_dostring(L,
		"co = coroutine.create(function()			\
				counter(5)							\
			end)									\
		for i = 5, 0, -1 do						\
			print(coroutine.resume(co))				\
		end");
}
```

>---
#### 构造 C 模块

一个 C 模块只有一个用于打开库的公共函数 `luaopen_*`，其他函数为 `static`。库函数注册到 `luaL_Reg` 结构数组中。

```c
static int l_func(lua_State *L){
	// implementation
}

// 注册 C 库函数
static const struct luaL_Reg mylib[] = {
	{"func", l_func},
	{NULL, NULL}, // 哨兵, 标识数组结尾
};

// luaopen_moduleName
int luaopen_mylib(lua_State* L) {
	luaL_newlib(L, mylib);   // 为 mylib 构建一张表作为共同的环境
	return 1;
}
```

win 生成命令：

```powershell
$ gcc -shared -o mylib.dll mylib.c -I"C:\lua\include" -L"C:\lua\lib" -llua54  # -llua54 标识 lua54.lib
```

lua 中加载 C 库：

```lua
package.cpath = package.cpath .. ';/path/to/?.dll'
local mylib = require "mylib"  -- load luaopen_mylib
mylib.func();				   -- call C function
```

>---
#### 在 C 中保存状态

C 函数在调用结束后会清空 *stack*，C API 提供了注册表（*registry*）和上值（*upvalue*）来存储非局部数据。

注册表（registry）是一张只能被 C API 访问的全局 Lua 表，用来存储多个 C 模块间共享的数据，数据生命周期与 Lua 状态机绑定。注册表位于伪索引 `LUA_REGISTRYINDEX`，可由 `lua_getfield(L,LUA_REGISTRYINDEX,"key")` 获取注册表中键 `key` 的值。

`int ref = luaL_ref(L, LUA_REGISTRYINDEX)` 用于将当前栈顶的值弹出并创建引用，返回的引用值唯一，对于 `nil` 值不会创建新的引用；`luaL_unref` 释放引用和关联值。在注册表中通常不能使用数值类型的键，Lua 将其用作引用系统的保留字（例如 `lua_rawgeti(L,LUA_REGISTRYINDEX,ref)` 压栈引用键的值）。

```c
int pushstr(lua_State* L, const char* s) {
	lua_pushstring(L, s);
	return luaL_ref(L, LUA_REGISTRYINDEX);
}
int main() {
	lua_State* L = luaL_newstate();
	int ref = pushstr(L, "Hello World");
	printf("ref = %d\n", ref);
	lua_rawgeti(L, LUA_REGISTRYINDEX, ref);
	printf("type = %s, value = %s\n", luaL_typename(L, -1), lua_tostring(L, -1));
	// type = string, value = Hello World
	lua_pop(L, 1);
	luaL_unref(L, LUA_REGISTRYINDEX, ref);
	lua_rawgeti(L, LUA_REGISTRYINDEX, ref);
	printf("type = %s, value = %s\n", luaL_typename(L, -1), lua_tostring(L, -1));
	// type = string, value = 0

	lua_close(L);
}
```

上值使用类似 C 静态变量的机制，每次在 lua 中创建新 C 函数时，都可以将任意数目的上值与该函数关联。调用该函数时通过伪索引访问这些上值。这种关联方式利用了 C 函数闭包机制。

```c
static int Counter(lua_State* L) {
	int val = lua_tointeger(L, lua_upvalueindex(1));   // 第一个上值的伪索引
	lua_pushinteger(L, ++val);   // 新值
	lua_copy(L, -1, lua_upvalueindex(1));   // 更新上值
	return 1;   // 返回栈顶新值
}

int newCounter(lua_State* L) {
	lua_pushinteger(L, 0);
	lua_pushcclosure(L, &Counter, 1);  // 上值数量 1
	return 1;  // 返回一个 C 函数闭包
}

void main() {
	lua_State* L = luaL_newstate();
	luaL_openlibs(L);
	lua_register(L, "newCounter", newCounter);
	luaL_dostring(L, "counter = newCounter()");  

	luaL_dostring(L, "print(counter())");   // 1
	luaL_dostring(L, "print(counter())");   // 2
	luaL_dostring(L, "print(counter())");   // 3 
	luaL_dostring(L, "print(counter())");   // 4
}; 
```

>---
#### 跨状态机通信

Lua 状态机之间不能直接通信，需要借助 C 进行数据传递，例如字符串和数值；表需要序列化后才能传递。

```c
lua_pushstring(L2, lua_tostring(L1, -1));
```


实现一个多线程并发库（C++） [`lproc`](./Lua%20LIB/LuaWithC/Concurrency/concurrency.cpp)，为每个线程创建一个独立的 Lua 状态机。[Lua 方调用](./Lua%20LIB/LuaWithC/Concurrency/lproc.lua)：

```lua
local lproc = require "lproc" 
--[[ API
lproc.start(chunk)
lproc.send(channel, v1, v2, ...)
v1, v2 = lproc,receive(channel)
lproc.exit()
]]

lproc.start([[
for i = 1, 5 do
    lproc.send("mess_queue", "Mess_"..i)
end
lproc.send("mess_queue", nil) --结束信号
lproc.exit()
]])

lproc.start([[
while true do
    local mess = lproc.receive("mess_queue")
    if not mess then break end
    print("receive:", mess)
end
lproc.exit()
]])

--[[ output
receive:        Mess_1
receive:        Mess_2
receive:        Mess_3
receive:        Mess_4
receive:        Mess_5
]]
```

>---
### 附录
#### C API Headers

- [lua.h](./Lua%20LIB/api_c/lua_ref.h) C API 基础库，定义 Lua 虚拟机（Lua State）的 C API 接口、基础数据类型和关键宏。
- [lualib.h](./Lua%20LIB/api_c/lualib_ref.h) Lua 标准库。
- [lauxlib.h](./Lua%20LIB/api_c/lauxlib_ref.h) C API 辅助库，简化 C 模块的编写和注册。

>---
#### Source

> 虚拟机核心相关

| 源文件                                          | 作用                                                                    | 备注                                                            | 对外接口前缀 |
| :---------------------------------------------- | :---------------------------------------------------------------------- | :-------------------------------------------------------------- | :----------- |
| [lapi.c](./Lua%20LIB/lua548/src/lapi.c)         | 提供 Lua 的 C API 实现，负责与宿主语言（如 ./Lua%20LIB/lua548/C++）交互 | 包含 `lua_push*`、`lua_call` 等函数，是 Lua 与 C 通信的核心接口 | `lua_`       |
| [lcode.c](./Lua%20LIB/lua548/src/lcode.c)       | 代码生成器，将语法树转换为虚拟机字节码                                  | 与 `lparser.c` 配合，生成优化后的指令序列                       | `luaK_`      |
| ​[lctype.c​](./Lua%20LIB/lua548/src/lctype.c)   | 实现 C 标准库 ctype 相关功能                                            | 提供高效的字符处理函数，避免依赖宿主环境的 C 库                 |
| [ldebug.c](./Lua%20LIB/lua548/src/ldebug.c)     | 调试接口，支持错误追踪、字节码验证等                                    | 通过 *CallInfo* 链表获取调用栈信息，用于错误回溯                | `luaG_`      |
| [ldo.c](./Lua%20LIB/lua548/src/ldo.c)           | 管理函数调用栈、异常处理和协程调度                                      | 核心机制包括 `setjmp/longjmp` 异常、*CallInfo* 链表             | `luaD_`      |
| [ldump.c](./Lua%20LIB/lua548/src/ldump.c)       | 序列化预编译的 Lua 字节码（如 luac 的输出）                             | 将字节码转换为二进制格式，支持跨平台加载                        |
| [lfunc.c](./Lua%20LIB/lua548/src/lfunc.c)       | 管理函数原型（*Proto*）和闭包（*Closure*）                              | 处理闭包的创建和调用，支持 *upvalue* 绑定                       | `luaF_`      |
| [lgc.c](./Lua%20LIB/lua548/src/lgc.c)           | 增量式垃圾回收器（GC），管理内存自动回收                                | 标记-清除算法，支持弱引用和分代回收                             | `luaC_`      |
| [llex.c](./Lua%20LIB/lua548/src/llex.c)         | 词法分析器，将源代码转换为 Token 流                                     | 与 `lparser.c` 协同工作，解析标识符、关键字等                   | `luaX_`      |
| [lmem.c](./Lua%20LIB/lua548/src/lmem.c)         | 内存管理接口，封装 `malloc/realloc`                                     | 提供内存分配和释放的统一接口，支持内存不足时的错误处理          | `luaM_`      |
| [lobject.c](./Lua%20LIB/lua548/src/lobject.c)   | 泛型对象操作函数（如类型转换、比较）                                    | 处理 TValue 结构（Lua 动态类型的核心表示）                      | `luaO_`      |
| [lopcodes.c](./Lua%20LIB/lua548/src/lopcodes.c) | 定义虚拟机字节码指令及其格式                                            | 包含指令编码和解码逻辑，与 `lvm.c` 紧密关联                     | `luaP_`      |
| [lparser.c](./Lua%20LIB/lua548/src/lparser.c)   | 递归下降解析器，将 Token 流转换为抽象语法树（AST）                      | 实现 Lua 语法的解析（如函数、表达式）                           | `luaY_`      |
| [lstate.c](./Lua%20LIB/lua548/src/lstate.c)     | 管理全局状态机（`global_State`）和线程状态（`lua_State`）               | 初始化虚拟机、维护注册表和字符串池                              | `luaE_`      |
| [lstring.c](./Lua%20LIB/lua548/src/lstring.c)   | 字符串池（String Interning），避免重复字符串的内存分配                  | 使用哈希表存储唯一字符串，优化内存使用                          | `luaS_`      |
| [ltable.c](./Lua%20LIB/lua548/src/ltable.c)     | ​	实现 Lua 的 *table* 类型，结合数组和哈希表                            | 动态调整数组和哈希部分的比例，优化访问效率                      | `luaH_`      |
| [ltm.c](./Lua%20LIB/lua548/src/ltm.c)           | ​	元方法（Metamethod）处理，支持运算符重载                              | 定义 `__add`、`__index` 等元方法的触发逻辑                      | `luaT_`      |
| [lundump.c](./Lua%20LIB/lua548/src/lundump.c)   | 虚拟机中负责加载预编译二进制代码                                        | 与 `ldump.c`（序列化字节码）配合完成 Lua 二进制文件的读写       | `luaU_`      |
| [lvm.c](./Lua%20LIB/lua548/src/lvm.c)           | 虚拟机执行核心，解释字节码并调度元方法                                  | 主循环 `luaV_execute` 解释指令，调用 `ltm.c` 的元方法           | `luaV_`      |
| [lzio.c](./Lua%20LIB/lua548/src/lzio.c)         | 输入流接口，支持从文件或内存中读取代码                                  | 为 `llex.c` 提供统一的字符流输入                                | `luaZ_`      |

> 内嵌库相关

| 源文件                                          | 作用                                | 备注                                                               |
| :---------------------------------------------- | :---------------------------------- | :----------------------------------------------------------------- |
| [lauxlib.c](./Lua%20LIB/lua548/src/lauxlib.c)   | 辅助库函数，简化 C 模块的编写和注册 | 提供 `luaL_newlib`、`luaL_check*` 等工具函数，常用于扩展库开发     |
| [lbaselib.c](./Lua%20LIB/lua548/src/lbaselib.c) | Lua `base` 库实现                   | 包含标准库的初始化逻辑，如 `_G` 表的构建                           |
| [lcorolib.c](./Lua%20LIB/lua548/src/lcorolib.c) | Lua `coroutine` 库实现              | 基于 `ldo.c` 的栈管理机制实现协程功能                              |
| [ldblib.c](./Lua%20LIB/lua548/src/ldblib.c)     | Lua `debug` 库实现                  | 提供调试 *hook*、堆栈访问和错误追踪功能，与 `ldebug.c` 协作        |
| [liolib.c](./Lua%20LIB/lua548/src/liolib.c)     | Lua `io` 库实现                     | 封装文件读写和流操作，支持标准输入输出及文件处理                   |
| [lmathlib.c](./Lua%20LIB/lua548/src/lmathlib.c) | Lua `math` 库实现                   | 实现基础数学运算和随机数生成，依赖 C 标准库的数学函数              |
| [loadlib.c](./Lua%20LIB/lua548/src/loadlib.c)   | Lua `package` 库实现                | Lua 动态库加载器                                                   |
| [loslib.c](./Lua%20LIB/lua548/src/loslib.c)     | Lua `os` 库实现                     | 提供时间、进程管理等系统级功能，平台相关性较强                     |
| [lstrlib.c](./Lua%20LIB/lua548/src/lstrlib.c)   | Lua `string` 库实现                 | 实现字符串匹配、格式化等操作，优化性能                             |
| [ltablib.c](./Lua%20LIB/lua548/src/ltable.c)    | Lua `table` 库实现                  | 扩展 table 类型的常用方法，与 ltable.c 的核心实现分离              |
| [lutf8lib.c](./Lua%20LIB/lua548/src/lutf8lib.c) | Lua `utf8` 库实现                   | 处理 Unicode 字符串，适用于 Lua 5.3+ 版本                          |
| [linit.c](./Lua%20LIB/lua548/src/linit.c)       | 标准库初始化入口                    | 调用 `luaL_openlibs` 加载所有内置库（如 `lmathlib.c`、`liolib.c`） |


> 解析器、字节码相关

| 源文件                                  | 作用                                              | 备注                                                |
| :-------------------------------------- | :------------------------------------------------ | :-------------------------------------------------- |
| [lua.c](./Lua%20LIB/lua548/src/lua.c)   | 独立解释器，处理命令行交互和编译                  | 包含 `main` 函数，调用 `luaL_newstate` 初始化虚拟机 |
| [luac.c](./Lua%20LIB/lua548/src/luac.c) | 字节码编译器，将 Lua 源码编译为预编译的二进制文件 | 依赖 `ldump.c` 输出字节码，与 `lua.c` 共享核心模块  |


>---
#### Examples

- [C 调用 Lua](./Lua%20LIB/LuaWithC/code/C_Call_Lua.c)，[Lua source](./Lua%20LIB/LuaWithC/code/CCallLua.lua)
- [lua 调用 C 函数](./Lua%20LIB/LuaWithC/code/Lua_Call_C.c)
- [C 函数的协程调度模型](./Lua%20LIB/LuaWithC/code/Continue_CFunc.c)
- [用户数据类型](./Lua%20LIB/LuaWithC/code/BitArray.c)，定义一个布尔数组（*BitArray*）。
- [多线程并发 lproc 库设计](./Lua%20LIB/LuaWithC/code/concurrency.cpp)，[Lua source](./Lua%20LIB/LuaWithC/code/lproc.lua)
- [状态机内存分配限制](./Lua%20LIB/LuaWithC/code/memlimit.c)，[Lua source](./Lua%20LIB/LuaWithC/code/memlimit.lua)
- [C 库函数](./Lua%20LIB/CMudule/README.md)

---