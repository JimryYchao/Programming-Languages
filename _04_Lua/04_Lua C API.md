## Lua C API

- version = 5.5

>---

### C API

#### 1. 类型定义与常量

这些是 C API 中使用的核心类型，并非函数，但作为接口的重要组成部分。

```c
lua_Alloc;         // 内存分配器类型
lua_CFunction;     // C 函数类型
lua_Debug;        // 调试信息结构体
lua_Hook;         // 调试钩子类型
lua_Integer;      // Lua 整型（默认为 long long）
lua_KContext;     // 延续函数上下文类型
lua_KFunction;    // 延续函数类型
lua_Number;       // Lua 数字类型（默认为 double）
lua_Reader;       // 自定义读取器类型
lua_State;        // 虚拟机状态句柄
lua_Unsigned;     // 无符号整型
lua_WarnFunction; // 警告函数类型
lua_Writer;       // 自定义写入器类型
```

>---

#### 2. 状态生命周期管理

```c
lua_atpanic(lua_State* L, lua_CFunction panicf) -> lua_CFunction; // 设置 panic 函数
lua_close(lua_State* L);                             // 关闭状态
lua_closethread(lua_State* L, lua_State* L1);        // 关闭线程（协程）
lua_getextraspace(lua_State* L) -> void*;            // 获取状态附加内存空间
lua_newstate(lua_Alloc f, void* ud) -> lua_State*;   // 创建新状态（使用自定义分配器）
lua_version(lua_State* L) -> const char*;            // 返回 Lua 核心版本号
```

>---

#### 3. 栈操作

> 栈查询与移动

```c
lua_absindex(lua_State* L, int idx) -> int;       // 将相对/伪索引转换为绝对索引
lua_checkstack(lua_State* L, int n) -> int;       // 检查并扩容栈空间
lua_copy(lua_State* L, int fromidx, int toidx);   // 复制元素到另一位置
lua_gettop(lua_State* L) -> int;                  // 获取栈顶索引
lua_insert(lua_State* L, int idx);                // 将栈顶元素插入指定位置
lua_pop(lua_State* L, int n);                     // 弹出 n 个元素（宏）
lua_remove(lua_State* L, int idx);                // 移除指定索引元素
lua_replace(lua_State* L, int idx);                // 将栈顶元素替换到指定位置
lua_rotate(lua_State* L, int idx, int n);         // 旋转栈区间
lua_settop(lua_State* L, int idx);                // 设置栈顶（可弹出元素）
lua_xmove(lua_State* from, lua_State* to, int n); // 在不同状态间移动栈元素
```

> 压入元素

```c
lua_pushboolean(lua_State* L, bool b); // 压入布尔值
lua_pushcclosure(lua_State* L, lua_CFunction fn, int n); // 压入 C 闭包
lua_pushcfunction(lua_State* L, lua_CFunction fn); // 压入 C 函数
lua_pushexternalstring(lua_State* L, const char* s, size_t len, void (*dtor)(void*), void* ud); // 压入外部字符串（零拷贝）
lua_pushfstring(lua_State* L, const char* fmt, ...) -> const char*; // 格式化字符串并压入（参数列表）
lua_pushglobaltable(lua_State* L); // 压入全局环境表
lua_pushinteger(lua_State* L, lua_Integer n); // 压入整型
lua_pushlightuserdata(lua_State* L, void* p); // 压入轻量用户数据
lua_pushliteral(lua_State* L, const char* s); // 压入字符串字面量（宏）
lua_pushlstring(lua_State* L, const char* s, size_t len); // 压入指定长度字符串
lua_pushnil(lua_State* L); // 压入 nil
lua_pushnumber(lua_State* L, lua_Number n); // 压入数字
lua_pushstring(lua_State* L, const char* s); // 压入以零结尾的字符串
lua_pushthread(lua_State* L);                // 压入当前线程
lua_pushvalue(lua_State* L, int idx); // 压入指定索引的副本
lua_pushvfstring(lua_State* L, const char* fmt, va_list argp) -> const char*; // 格式化字符串并压入（可变参数）
```

> 关闭槽

```c
lua_closeslot(lua_State* L, int idx);        // 关闭指定槽（触发 to-be-closed 变量）
lua_toclose(lua_State* L, int idx) -> void*;               // 标记栈顶值为待关闭
```

>---

#### 4. 类型判断与转换

> 类型查询

```c
lua_type(lua_State* L, int idx) -> int;             // 返回元素类型
lua_typename(lua_State* L, int tp) -> const char*;  // 将类型编号转为字符串
```

> 类型判断

```c
lua_isboolean(lua_State* L, int idx) -> bool;       // 检查是否为布尔值
lua_iscfunction(lua_State* L, int idx) -> bool;     // 检查是否为 C 函数
lua_isfunction(lua_State* L, int idx) -> bool;      // 检查是否为函数
lua_isinteger(lua_State* L, int idx) -> bool;       // 检查是否为整数
lua_islightuserdata(lua_State* L, int idx) -> bool; // 检查是否为轻量用户数据
lua_isnil(lua_State* L, int idx) -> bool;           // 检查是否为 nil
lua_isnone(lua_State* L, int idx) -> bool;          // 检查是否为无值
lua_isnoneornil(lua_State* L, int idx) -> bool;     // 检查是否为无值或 nil
lua_isnumber(lua_State* L, int idx) -> bool;        // 检查是否为数字
lua_isstring(lua_State* L, int idx) -> bool;        // 检查是否为字符串
lua_istable(lua_State* L, int idx) -> bool;         // 检查是否为表
lua_isthread(lua_State* L, int idx) -> bool;        // 检查是否为线程
lua_isuserdata(lua_State* L, int idx) -> bool;      // 检查是否为用户数据
lua_isyieldable(lua_State* L) -> bool;              // 检查是否可以挂起
```

> 数值转换

```c
lua_numbertointeger(lua_Number n) -> lua_Integer;   // 将数字转换为整型（宏）
lua_stringtonumber(const char* s) -> lua_Number;    // 将字符串转换为数字并压入
lua_toboolean(lua_State* L, int idx) -> bool;              // 转为布尔值
lua_tocfunction(lua_State* L, int idx) -> lua_CFunction;   // 转 C 函数指针
lua_toclose(lua_State* L, int idx) -> void*;               // 标记栈顶值为待关闭
lua_tointeger(lua_State* L, int idx) -> lua_Integer;       // 转为整型
lua_tointegerx(lua_State* L, int idx, int* isint) -> lua_Integer; // 转为整型（带转换检查）
lua_tolstring(lua_State* L, int idx, size_t* len) -> const char*; // 转为字符串（可能引发 GC）
lua_tonumber(lua_State* L, int idx) -> lua_Number;         // 转为数字
lua_tonumberx(lua_State* L, int idx, int* isnum) -> lua_Number;   // 转为数字（带转换检查）
lua_topointer(lua_State* L, int idx) -> const void*;       // 转为指针（调试用）
lua_tostring(lua_State* L, int idx) -> const char*;        // 转为字符串
lua_tothread(lua_State* L, int idx) -> lua_State*;         // 转为线程指针
lua_touserdata(lua_State* L, int idx) -> void*;            // 转为用户数据指针
```

>---

#### 5. 表操作

```c
lua_createtable(lua_State* L, int narr, int nrec);  // 创建预分配大小的表
lua_getfield(lua_State* L, int idx, const char* k); // 根据字段名获取值
lua_geti(lua_State* L, int idx, lua_Integer n);     // 根据整数键获取值
lua_gettable(lua_State* L, int idx);                // 根据栈键获取值
lua_newtable(lua_State* L);                         // 创建新表
lua_next(lua_State* L, int idx) -> bool;            // 遍历表
lua_rawget(lua_State* L, int idx);                  // 原始访问（不触发元表）
lua_rawgeti(lua_State* L, int idx, lua_Integer n);  // 原始整数键访问
lua_rawgetp(lua_State* L, int idx, const void* p);  // 原始指针键访问（轻量用户数据作键）
lua_rawlen(lua_State* L, int idx) -> size_t;        // 获取原始长度（表或字符串）
lua_rawset(lua_State* L, int idx);                  // 原始设置
lua_rawseti(lua_State* L, int idx, lua_Integer n);  // 原始整数键设置
lua_rawsetp(lua_State* L, int idx, const void* p);  // 原始指针键设置
lua_setfield(lua_State* L, int idx, const char* k); // 设置表字段
lua_seti(lua_State* L, int idx, lua_Integer n);     // 设置整数键
lua_settable(lua_State* L, int idx);                // 设置表键值对
```

>---

#### 6. 元表与用户值

```c
lua_getiuservalue(lua_State* L, int idx, int n) -> int;  // 获取用户数据关联值（按索引）
lua_getmetatable(lua_State* L, int idx) -> int;     // 获取元表
lua_setiuservalue(lua_State* L, int idx, int n) -> int;  // 设置用户数据关联值
lua_setmetatable(lua_State* L, int idx) -> int;     // 设置元表
```

>---

#### 7. 函数调用与协程

```c
lua_call(lua_State* L, int nargs, int nresults);        // 调用 Lua 函数
lua_callk(lua_State* L, int nargs, int nresults, lua_KContext ctx, lua_KFunction k); // 带延续的调用函数
lua_pcall(lua_State* L, int nargs, int nresults, int errfunc) -> int; // 保护模式调用
lua_pcallk(lua_State* L, int nargs, int nresults, int errfunc, lua_KContext ctx, lua_KFunction k) -> int; // 带延续的保护调用
lua_resume(lua_State* L, lua_State* from, int nargs) -> int; // 恢复协程
lua_status(lua_State* L) -> int;               // 获取协程状态
lua_yield(lua_State* L, int nresults) -> int;  // 挂起当前协程
lua_yieldk(lua_State* L, int nresults, lua_KContext ctx, lua_KFunction k) -> int; // 带延续的挂起协程
```

>---

#### 8. 错误处理与警告

```c
lua_error(lua_State* L) -> int; // 抛出错误
lua_setwarnf(lua_State* L, lua_WarnFunction warnf, void* ud); // 设置警告函数
lua_warning(lua_State* L, const char* msg, int tocont); // 发出警告
```

>---

#### 9. 垃圾回收

```c
lua_gc(lua_State* L, int what, int data) -> int;     // 控制垃圾回收器（GC 操作）
```

>---

#### 10. 调试接口

```c
lua_gethook(lua_State* L) -> lua_Hook;              // 获取当前钩子
lua_gethookcount(lua_State* L) -> int;              // 获取钩子计数值
lua_gethookmask(lua_State* L) -> int;               // 获取钩子掩码
lua_getinfo(lua_State* L, const char* what, lua_Debug* ar) -> int;       // 获取调试信息
lua_getlocal(lua_State* L, const lua_Debug* ar, int n) -> const char*;   // 获取局部变量
lua_getstack(lua_State* L, int level, lua_Debug* ar) -> int;             // 获取调用栈信息
lua_getupvalue(lua_State* L, int funcidx, int n) -> const char*;         // 获取上值
lua_sethook(lua_State* L, lua_Hook func, int mask, int count); // 设置钩子
lua_setlocal(lua_State* L, const lua_Debug* ar, int n) -> const char*;   // 设置局部变量
lua_setupvalue(lua_State* L, int funcidx, int n) -> const char*;         // 设置上值
lua_upvalueid(lua_State* L, int funcidx, int n) -> void*;                // 获取上值唯一 ID
lua_upvalueindex(int i) -> int;                         // 创建上值伪索引（宏）
lua_upvaluejoin(lua_State* L, int funcidx1, int n1, int funcidx2, int n2);  // 连接两个上值
```

>---

#### 11. 内存分配器

```c
lua_getallocf(lua_State* L, void** ud) -> lua_Alloc; // 获取内存分配器函数
lua_setallocf(lua_State* L, lua_Alloc f, void* ud);  // 设置内存分配器函数
```

>---

#### 12. 其他

```c
lua_arith(lua_State* L, int op) -> int;      // 执行算术运算
lua_compare(lua_State* L, int idx1, int idx2, int op) -> bool;           // 唯较两个值
lua_concat(lua_State* L, int n) -> int;      // 连接栈顶 n 个字符串
lua_dump(lua_State* L, lua_Writer writer, void* data, int strip) -> int; // 将函数转储为二进制块
lua_getglobal(lua_State* L, const char* name);  // 获取全局变量
lua_len(lua_State* L, int idx) -> size_t;          // 获取长度（触发长度元方法）
lua_load(lua_State* L, lua_Reader reader, void* data, const char* chunkname, const char* mode) -> int; // 加载 Lua 代码块
lua_newuserdatauv(lua_State* L, size_t sz, int nuvalue) -> void*;           // 创建带用户值的用户数据
lua_numbertocstring(lua_State* L, lua_Number n) -> const char*;          // 将数字转换为字符串
lua_rawequal(lua_State* L, int idx1, int idx2) -> bool; // 原始相等比较
lua_register(lua_State* L, const char* name, lua_CFunction f); // 将 C 函数注册为全局变量（宏）
lua_setglobal(lua_State* L, const char* name);  // 设置全局变量
```

>---

### C 辅助库 (auxiliary library)

#### 1. 类型与结构体

```c
luaL_Buffer     // 缓冲区结构体
luaL_Reg        // 函数注册结构体
luaL_Stream     // 流结构体（文件 I/O）
```

>---

#### 2. 状态创建与库加载

```c
luaL_newstate() -> lua_State*;                          // 创建新状态（使用标准分配器）
luaL_openlibs(lua_State* L);                            // 打开所有标准库
luaL_openselectedlibs(lua_State* L, int what);          // 选择性打开标准库（Lua 5.5 新增）
luaL_requiref(lua_State* L, const char* modname, lua_CFunction openf, int glb); // 加载模块并存入全局表
```

>---

#### 3. 错误处理与参数检查

```c
luaL_argcheck(lua_State* L, bool cond, int arg, const char* extramsg); // 检查条件，失败则报错
luaL_argerror(lua_State* L, int arg, const char* extramsg) -> int; // 报告参数错误
luaL_argexpected(lua_State* L, bool cond, int arg, const char* tname); // 报告期望类型错误
luaL_checkany(lua_State* L, int arg);                   // 确保参数存在
luaL_checkinteger(lua_State* L, int arg) -> lua_Integer; // 获取整数参数
luaL_checklstring(lua_State* L, int arg, size_t* l) -> const char*; // 获取字符串参数
luaL_checknumber(lua_State* L, int arg) -> lua_Number;  // 获取数字参数
luaL_checkoption(lua_State* L, int arg, const char* def, const char* const lst[]) -> int; // 检查选项字符串
luaL_checkstack(lua_State* L, int sz, const char* msg); // 检查栈空间
luaL_checkstring(lua_State* L, int arg) -> const char*; // 获取字符串参数
luaL_checktype(lua_State* L, int arg, int t);           // 检查参数类型
luaL_checkudata(lua_State* L, int arg, const char* tname) -> void*; // 检查用户数据类型
luaL_error(lua_State* L, const char* fmt, ...) -> int;  // 抛出错误（带格式化）
luaL_opt(lua_State* L, int arg, lua_Integer def) -> lua_Integer; // 可选参数处理（宏）
luaL_optinteger(lua_State* L, int arg, lua_Integer def) -> lua_Integer; // 可选整数参数
luaL_optlstring(lua_State* L, int arg, const char* def, size_t* l) -> const char*; // 可选字符串参数
luaL_optnumber(lua_State* L, int arg, lua_Number def) -> lua_Number; // 可选数字参数
luaL_optstring(lua_State* L, int arg, const char* def) -> const char*; // 可选字符串参数
luaL_typeerror(lua_State* L, int arg, const char* tname) -> int; // 报告类型错误
```

>---

#### 4. 缓冲区操作

```c
luaL_addchar(luaL_Buffer* B, char c);                   // 添加单个字符
luaL_addgsub(luaL_Buffer* B, const char* s, const char* p, const char* r); // 添加并执行替换
luaL_addlstring(luaL_Buffer* B, const char* s, size_t l); // 添加指定长度字符串
luaL_addsize(luaL_Buffer* B, size_t n);                 // 更新缓冲区大小
luaL_addstring(luaL_Buffer* B, const char* s);          // 添加零结尾字符串
luaL_addvalue(luaL_Buffer* B);                          // 将栈顶字符串添加到缓冲区
luaL_buffaddr(luaL_Buffer* B) -> char*;                 // 获取缓冲区地址
luaL_buffinitsize(lua_State* L, luaL_Buffer* B, size_t sz) -> char*; // 初始化并预分配大小
luaL_buffinit(lua_State* L, luaL_Buffer* B);            // 初始化缓冲区
luaL_bufflen(luaL_Buffer* B) -> size_t;                 // 获取缓冲区当前长度
luaL_buffsub(luaL_Buffer* B, int n);                    // 从缓冲区末尾减去长度
luaL_prepbuffer(luaL_Buffer* B) -> char*;               // 准备写入缓冲区
luaL_prepbuffsize(luaL_Buffer* B, size_t sz) -> char*;  // 准备指定大小的缓冲区
luaL_pushresult(luaL_Buffer* B);                        // 将缓冲区结果压栈
luaL_pushresultsize(luaL_Buffer* B, size_t sz);         // 设置大小并压栈
```

>---

#### 5. 元表与对象操作

```c
luaL_callmeta(lua_State* L, int obj, const char* e) -> int; // 调用元方法
luaL_checkudata(lua_State* L, int arg, const char* tname) -> void*; // 检查并返回用户数据
luaL_getmetatable(lua_State* L, const char* tname);     // 获取注册的元表
luaL_getmetafield(lua_State* L, int obj, const char* e) -> int; // 获取元表字段
luaL_len(lua_State* L, int idx) -> lua_Integer;         // 获取长度（支持元方法）
luaL_newmetatable(lua_State* L, const char* tname) -> int; // 创建元表并注册
luaL_setmetatable(lua_State* L, const char* tname);     // 设置元表
luaL_testudata(lua_State* L, int arg, const char* tname) -> void*; // 测试并返回用户数据
luaL_tolstring(lua_State* L, int idx, size_t* len) -> const char*; // 转换为字符串（支持元方法）
```

>---

#### 6. 函数注册与模块

```c
luaL_newlib(lua_State* L, const luaL_Reg l[]);          // 创建并注册函数表
luaL_newlibtable(lua_State* L, const luaL_Reg l[]);     // 创建预大小的函数表
luaL_register(lua_State* L, const char* libname, const luaL_Reg* l); // 注册库（已弃用）
luaL_requiref(lua_State* L, const char* modname, lua_CFunction openf, int glb); // 加载模块
luaL_setfuncs(lua_State* L, const luaL_Reg* l, int nup); // 设置函数到表中
```

>---

#### 7. 文件 I/O 辅助

```c
luaL_execresult(lua_State* L, int stat) -> int;         // 生成系统调用结果
luaL_fileresult(lua_State* L, int stat, const char* fname) -> int; // 生成文件操作结果
```

>---

#### 8. 字符串处理

```c
luaL_gsub(lua_State* L, const char* s, const char* p, const char* r) -> const char*; // 全局字符串替换
luaL_tolstring(lua_State* L, int idx, size_t* len) -> const char*; // 转为字符串
```

>---

#### 9. 引用系统

```c
luaL_ref(lua_State* L, int t) -> int;                   // 创建引用
luaL_unref(lua_State* L, int t, int ref);               // 释放引用
```

>---

#### 10. 调试与跟踪

```c
luaL_traceback(lua_State* L, lua_State* L1, const char* msg, int level); // 生成调用栈回溯
luaL_where(lua_State* L, int level);                    // 生成当前位置信息
```

>---

#### 11. 随机数种子

```c
luaL_makeseed(lua_State* L) -> unsigned int;            // 生成可靠的随机数种子
```

>---

#### 12. 其他

```c
luaL_alloc(void* ud, void* ptr, size_t osize, size_t nsize) -> void*; // 内存分配函数
luaL_checkversion(lua_State* L);                        // 检查库与核心版本兼容性
luaL_dofile(lua_State* L, const char* filename) -> int; // 执行文件（宏）
luaL_dostring(lua_State* L, const char* str) -> int;    // 执行字符串（宏）
luaL_execresult(lua_State* L, int stat) -> int;         // 处理系统调用结果
luaL_fileresult(lua_State* L, int stat, const char* fname) -> int; // 处理文件操作结果
luaL_getsubtable(lua_State* L, int idx, const char* fname) -> int; // 获取或创建子表
luaL_loadbuffer(lua_State* L, const char* buff, size_t sz, const char* name) -> int; // 加载缓冲区
luaL_loadbufferx(lua_State* L, const char* buff, size_t sz, const char* name, const char* mode) -> int; // 加载缓冲区（带模式）
luaL_loadfile(lua_State* L, const char* filename) -> int; // 加载文件
luaL_loadfilex(lua_State* L, const char* filename, const char* mode) -> int; // 加载文件（带模式）
luaL_loadstring(lua_State* L, const char* s) -> int;    // 加载字符串
luaL_pushfail(lua_State* L);                            // 压入失败标志（nil + 错误信息）
luaL_typename(lua_State* L, int idx) -> const char*;    // 获取类型名称
```

---