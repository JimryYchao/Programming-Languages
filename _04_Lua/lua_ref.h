// Lua 5.4.8
#pragma once
#include <stdarg.h>
#include <stddef.h>

#define LUA_NUMBER float
#define LUA_INTEGER int
#define LUA_UNSIGNED unsigned
#define LUA_KCONTEXT ptrdiff_t
#define LUA_IDSIZE 60

/* 预编译代码标记 ('<esc>Lua') */
#define LUA_SIGNATURE "\x1bLua"

/* lua_pcall/lua_call 的多返回值选项 */
#define LUA_MULTRET (-1)

/* 伪索引（Pseudo-indices）(-LUAI_MAXSTACK 是最小有效索引；保留部分空间用于溢出检测) */
#define LUA_REGISTRYINDEX (-LUAI_MAXSTACK - 1000)     // 注册表伪索引
#define lua_upvalueindex(i) (LUA_REGISTRYINDEX - (i)) // 上值索引

/* 线程状态码 */
#define LUA_OK 0        // 成功
#define LUA_YIELD 1     // 协程挂起
#define LUA_ERRRUN 2    // 运行时错误
#define LUA_ERRSYNTAX 3 // 语法错误
#define LUA_ERRMEM 4    // 内存错误
#define LUA_ERRERR 5    // 错误处理函数自身出错

/* C 函数可用的最小 Lua 栈空间 */
#define LUA_MINSTACK 20

/* 注册表中的预定义索引 */
#define LUA_RIDX_MAINTHREAD 1 // 主线程引用
#define LUA_RIDX_GLOBALS 2    // 全局环境表
#define LUA_RIDX_LAST LUA_RIDX_GLOBALS

/* 类型预定义 */
typedef struct lua_State lua_State;                                          // Lua 状态机类型
typedef struct lua_Debug lua_Debug;                                          // 调试 API 使用的调试信息结构体
typedef LUA_NUMBER lua_Number;                                               // Lua 数值类型
typedef LUA_INTEGER lua_Integer;                                             // Lua 整数类型
typedef LUA_UNSIGNED lua_Unsigned;                                           // 无符号整数类型
typedef LUA_KCONTEXT lua_KContext;                                           // 延续函数上下文类型
typedef int (*lua_CFunction)(lua_State *L);                                  // 注册到 Lua 的 C 函数类型
typedef int (*lua_KFunction)(lua_State *L, int status, lua_KContext ctx);    // 延续函数类型
typedef const char *(*lua_Reader)(lua_State *L, void *ud, size_t *sz);       // 加载 Lua 代码块的读函数类型
typedef int (*lua_Writer)(lua_State *L, const void *p, size_t sz, void *ud); // 转储 Lua 代码块的写函数类型
typedef void *(*lua_Alloc)(void *ud, void *ptr, size_t osize, size_t nsize); // 内存分配函数类型
typedef void (*lua_WarnFunction)(void *ud, const char *msg, int tocont);     // 警告处理函数类型typedef
typedef void (*lua_Hook)(lua_State *L, lua_Debug *ar);                       // 调试器在特定事件中调用的钩子函数类型

/* ====================== 函数声明部分 ====================== */

/* 状态机管理 */
lua_State *lua_newstate(lua_Alloc f, void *ud);                // 创建新状态机
void lua_close(lua_State *L);                                  // 关闭状态机
lua_State *lua_newthread(lua_State *L);                        // 创建新协程
int lua_closethread(lua_State *L, lua_State *from);            // 关闭协程
lua_CFunction lua_atpanic(lua_State *L, lua_CFunction panicf); // 设置紧急函数
lua_Number lua_version(lua_State *L);                          // 获取版本号

/* 栈操作 */
int lua_absindex(lua_State *L, int idx);                             // 转换为绝对索引
int lua_gettop(lua_State *L);                                        // 获取栈顶索引
void lua_settop(lua_State *L, int idx);                              // 设置栈顶
void lua_pushvalue(lua_State *L, int idx);                           // 压入栈中某位置的副本
void lua_rotate(lua_State *L, int idx, int n);                       // 在 idx 和栈顶之间旋转 n 个元素。n>0 沿栈顶方向；n<0 沿栈底方向
void lua_copy(lua_State *L, int fromidx, int toidx);                 // 复制栈元素
int lua_checkstack(lua_State *L, int n);                             // 检查栈空间
void lua_xmove(lua_State *from, lua_State *to, int n);               // 跨状态机移动值
#define lua_insert(L, idx) lua_rotate(L, (idx), 1)                   // 栈顶插入到目标位置
#define lua_remove(L, idx) (lua_rotate(L, (idx), -1), lua_pop(L, 1)) // 移除元素
#define lua_replace(L, idx) (lua_copy(L, -1, (idx)), lua_pop(L, 1))  // 栈顶替换到目标位置，并弹出栈顶

/* 基础类型标识和类型检查函数 */
#define LUA_TNONE (-1)                          // 无效类型
#define LUA_TNIL 0                              // nil
#define LUA_TBOOLEAN 1                          // 布尔
#define LUA_TLIGHTUSERDATA 2                    // 轻量用户数据
#define LUA_TNUMBER 3                           // 数字
#define LUA_TSTRING 4                           // 字符串
#define LUA_TTABLE 5                            // 表
#define LUA_TFUNCTION 6                         // 函数
#define LUA_TUSERDATA 7                         // 完整用户数据
#define LUA_TTHREAD 8                           // 协程
#define LUA_NUMTYPES 9                          // 类型总数
int lua_type(lua_State *L, int idx);            // 返回值的实际类型的 基础类型标识宏
const char *lua_typename(lua_State *L, int tp); // 基础类型标识宏 转字符串
int lua_isnumber(lua_State *L, int idx);        // 检查是否为数字或可转数字的字符串
int lua_isstring(lua_State *L, int idx);        // 检查是否为字符串或数字
int lua_iscfunction(lua_State *L, int idx);     // 检查是否为 C 函数
int lua_isinteger(lua_State *L, int idx);       // 检查是否为整数
int lua_isuserdata(lua_State *L, int idx);      // 检查是否为 userdata
#define LUA_NUMTAGS LUA_NUMTYPES
#define lua_isfunction(L, n) (lua_type(L, (n)) == LUA_TFUNCTION)
#define lua_istable(L, n) (lua_type(L, (n)) == LUA_TTABLE)
#define lua_islightuserdata(L, n) (lua_type(L, (n)) == LUA_TLIGHTUSERDATA)
#define lua_isnil(L, n) (lua_type(L, (n)) == LUA_TNIL)
#define lua_isboolean(L, n) (lua_type(L, (n)) == LUA_TBOOLEAN)
#define lua_isthread(L, n) (lua_type(L, (n)) == LUA_TTHREAD)
#define lua_isnone(L, n) (lua_type(L, (n)) == LUA_TNONE)
#define lua_isnoneornil(L, n) (lua_type(L, (n)) <= 0)

/* 类型转换函数 */
lua_Number lua_tonumberx(lua_State *L, int idx, int *isnum);   // 转换为 lua_Number
lua_Integer lua_tointegerx(lua_State *L, int idx, int *isnum); // 转换为 lua_Integer
int lua_toboolean(lua_State *L, int idx);                      // 转换为布尔值(0/1)
lua_Unsigned lua_rawlen(lua_State *L, int idx);                // 获取原始长度
lua_CFunction lua_tocfunction(lua_State *L, int idx);          // 提取 C 函数指针
void *lua_touserdata(lua_State *L, int idx);                   // 提取 userdata 指针
lua_State *lua_tothread(lua_State *L, int idx);                // 提取协程指针
const char *lua_tolstring(lua_State *L, int idx, size_t *len); // 转换为字符串
const void *lua_topointer(lua_State *L, int idx);              // 获取通用指针
int lua_numbertointeger(lua_Number n, lua_Integer *p);         // 尝试将 lua_Number 转换为 lua_Integer
#define lua_tointeger(L, i) lua_tointegerx(L, (i), NULL)
#define lua_tonumber(L, i) lua_tonumberx(L, (i), NULL)

/* 算术运算 */
#define LUA_OPADD 0                   // 加法 (排序与元方法相关)
#define LUA_OPSUB 1                   // 减法
#define LUA_OPMUL 2                   // 乘法
#define LUA_OPMOD 3                   // 取模
#define LUA_OPPOW 4                   // 幂运算
#define LUA_OPDIV 5                   // 除法
#define LUA_OPIDIV 6                  // 整除
#define LUA_OPBAND 7                  // 按位与
#define LUA_OPBOR 8                   // 按位或
#define LUA_OPBXOR 9                  // 按位异或
#define LUA_OPSHL 10                  // 左移
#define LUA_OPSHR 11                  // 右移
#define LUA_OPUNM 12                  // 取负
#define LUA_OPBNOT 13                 // 按位非
void lua_arith(lua_State *L, int op); // 执行算术运算

/* 比较运算 */
#define LUA_OPEQ 0                                         // 等于
#define LUA_OPLT 1                                         // 小于
#define LUA_OPLE 2                                         // 小于等于
int lua_rawequal(lua_State *L, int idx1, int idx2);        // 原始相等比较
int lua_compare(lua_State *L, int idx1, int idx2, int op); // 通用比较

/* 数据压栈 */
void lua_pushnil(lua_State *L);                                            // 压入 nil 值到栈顶
void lua_pushnumber(lua_State *L, lua_Number n);                           // 压入浮点数，lua_Number 通常为 double
void lua_pushinteger(lua_State *L, lua_Integer n);                         // 压入整数，lua_Integer 通常为 ptrdiff_t
void lua_pushboolean(lua_State *L, int b);                                 // 压入布尔值
const char *lua_pushlstring(lua_State *L, const char *s, size_t len);      // 压入指定长度的字符串
const char *lua_pushstring(lua_State *L, const char *s);                   // 压入以 \0 结尾的字符串
const char *lua_pushvfstring(lua_State *L, const char *fmt, va_list argp); // 压入格式化字符串 (va_list)
const char *lua_pushfstring(lua_State *L, const char *fmt, ...);           // 压入格式化字符串 (可变参数)
void lua_pushcclosure(lua_State *L, lua_CFunction fn, int n);              // 函数与闭包，压入 C 闭包并绑定 n 个上值
void lua_pushlightuserdata(lua_State *L, void *p);                         // 压入轻量 userdata，仅存储指针
int lua_pushthread(lua_State *L);                                          // 压入当前线程(协程)

/* 表操作 */
void lua_createtable(lua_State *L, int narr, int nrec);        // 创建新表
void *lua_newuserdatauv(lua_State *L, size_t sz, int nuvalue); // 创建用户数据
int lua_getmetatable(lua_State *L, int objindex);              // 获取元表
int lua_setmetatable(lua_State *L, int objindex);              // 设置元表
int lua_getiuservalue(lua_State *L, int idx, int n);           // 获取 uservalue
int lua_setiuservalue(lua_State *L, int idx, int n);           // 设置 uservalue
int lua_getglobal(lua_State *L, const char *name);             // 获取全局变量
void lua_setglobal(lua_State *L, const char *name);            // 设置全局变量
int lua_gettable(lua_State *L, int idx);                       // 通过键获取表值
void lua_settable(lua_State *L, int idx);                      // 设置表键值对
int lua_getfield(lua_State *L, int idx, const char *k);        // 获取字符串键值
void lua_setfield(lua_State *L, int idx, const char *k);       // 设置字符串键值
int lua_geti(lua_State *L, int idx, lua_Integer n);            // 获取整数键值
void lua_seti(lua_State *L, int idx, lua_Integer n);           // 设置整数键值
int lua_rawget(lua_State *L, int idx);                         // 原始获取（跳过元方法）
void lua_rawset(lua_State *L, int idx);
int lua_rawgeti(lua_State *L, int idx, lua_Integer n);
void lua_rawseti(lua_State *L, int idx, lua_Integer n);
int lua_rawgetp(lua_State *L, int idx, const void *p);
void lua_rawsetp(lua_State *L, int idx, const void *p);
#define lua_newuserdata(L, s) lua_newuserdatauv(L, s, 1)
#define lua_getuservalue(L, idx) lua_getiuservalue(L, idx, 1)
#define lua_setuservalue(L, idx) lua_setiuservalue(L, idx, 1)

/* 代码加载与执行 */
void lua_callk(lua_State *L, int nargs, int nresults, lua_KContext ctx, lua_KFunction k);              // 调用函数（支持延续）
int lua_pcallk(lua_State *L, int nargs, int nresults, int errfunc, lua_KContext ctx, lua_KFunction k); // 保护模式调用
int lua_load(lua_State *L, lua_Reader reader, void *dt, const char *chunkname, const char *mode);      // 加载代码块
int lua_dump(lua_State *L, lua_Writer writer, void *data, int strip);                                  // 转储函数
#define lua_call(L, n, r) lua_callk(L, (n), (r), 0, NULL)
#define lua_pcall(L, n, r, f) lua_pcallk(L, (n), (r), (f), 0, NULL)

/* 协程控制 */
int lua_yieldk(lua_State *L, int nresults, lua_KContext ctx, lua_KFunction k); // 挂起协程
int lua_resume(lua_State *L, lua_State *from, int narg, int *nres);            // 恢复协程
int lua_status(lua_State *L);                                                  // 获取协程状态
int lua_isyieldable(lua_State *L);                                             // 检查是否可挂起
#define lua_yield(L, n) lua_yieldk(L, (n), 0, NULL)

/* 警告处理 */
void lua_setwarnf(lua_State *L, lua_WarnFunction f, void *ud); // 设置警告函数
void lua_warning(lua_State *L, const char *msg, int tocont);   // 触发警告

/* 垃圾回收控制 */
#define LUA_GCSTOP 0                     // 停止 GC
#define LUA_GCRESTART 1                  // 重启 GC
#define LUA_GCCOLLECT 2                  // 执行完整收集
#define LUA_GCCOUNT 3                    // 获取当前内存用量(KB)
#define LUA_GCCOUNTB 4                   // 获取内存用量余数(字节)
#define LUA_GCSTEP 5                     // 执行增量步骤
#define LUA_GCSETPAUSE 6                 // 设置暂停参数
#define LUA_GCSETSTEPMUL 7               // 设置步进乘数
#define LUA_GCISRUNNING 9                // 检查 GC 是否运行
#define LUA_GCGEN 10                     // 切换到分代模式
#define LUA_GCINC 11                     // 切换到增量模式
int lua_gc(lua_State *L, int what, ...); // 控制垃圾回收器

/* 杂项功能 */
int lua_error(lua_State *L);                             // 抛出错误
int lua_next(lua_State *L, int idx);                     // 遍历表
void lua_concat(lua_State *L, int n);                    // 连接栈顶字符串
void lua_len(lua_State *L, int idx);                     // 获取长度
size_t lua_stringtonumber(lua_State *L, const char *s);  // 字符串转数字
lua_Alloc lua_getallocf(lua_State *L, void **ud);        // 获取内存分配器
void lua_setallocf(lua_State *L, lua_Alloc f, void *ud); // 设置内存分配器
void lua_toclose(lua_State *L, int idx);                 // 标记待关闭值
void lua_closeslot(lua_State *L, int idx);               // 关闭指定槽位

/* ====================== 宏定义 ====================== */

#define lua_getextraspace(L) ((void *)((char *)(L) - LUA_EXTRASPACE))                      // 获取状态机额外空间指针
#define lua_pop(L, n) lua_settop(L, -(n) - 1)                                              // 弹出 n 个值
#define lua_newtable(L) lua_createtable(L, 0, 0)                                           // 创建空表
#define lua_register(L, n, f) (lua_pushcfunction(L, (f)), lua_setglobal(L, (n)))           // 注册 C 函数为全局变量
#define lua_pushcfunction(L, f) lua_pushcclosure(L, (f), 0)                                // 压入 C 函数（无上值）
#define lua_pushliteral(L, s) lua_pushstring(L, "" s)                                      // 压入字面量字符串
#define lua_pushglobaltable(L) ((void)lua_rawgeti(L, LUA_REGISTRYINDEX, LUA_RIDX_GLOBALS)) // 获取全局环境表
#define lua_tostring(L, i) lua_tolstring(L, (i), NULL)                                     // 转换为字符串（无长度检查）

/* ====================== 调试 API ====================== */

/* 事件类型 */
#define LUA_HOOKCALL 0     // 函数调用时触发
#define LUA_HOOKRET 1      // 函数返回时触发
#define LUA_HOOKLINE 2     // 执行新行时触发
#define LUA_HOOKCOUNT 3    // 指令计数达到阈值时触发
#define LUA_HOOKTAILCALL 4 // 尾调用时触发

/* 事件掩码 */
#define LUA_MASKCALL (1 << LUA_HOOKCALL)
#define LUA_MASKRET (1 << LUA_HOOKRET)
#define LUA_MASKLINE (1 << LUA_HOOKLINE)
#define LUA_MASKCOUNT (1 << LUA_HOOKCOUNT)

/* 调试函数 */
int lua_getstack(lua_State *L, int level, lua_Debug *ar);                 // 获取调用栈信息
int lua_getinfo(lua_State *L, const char *what, lua_Debug *ar);           // 获取调试信息
const char *lua_getlocal(lua_State *L, const lua_Debug *ar, int n);       // 获取局部变量
const char *lua_setlocal(lua_State *L, const lua_Debug *ar, int n);       // 设置局部变量
const char *lua_getupvalue(lua_State *L, int funcindex, int n);           // 获取上值
const char *lua_setupvalue(lua_State *L, int funcindex, int n);           // 设置上值
void *lua_upvalueid(lua_State *L, int fidx, int n);                       // 获取上值唯一标识
void lua_upvaluejoin(lua_State *L, int fidx1, int n1, int fidx2, int n2); // 关联上值
void lua_sethook(lua_State *L, lua_Hook func, int mask, int count);       // 设置钩子
lua_Hook lua_gethook(lua_State *L);                                       // 获取当前钩子
int lua_gethookmask(lua_State *L);                                        // 获取钩子掩码
int lua_gethookcount(lua_State *L);                                       // 获取钩子计数器
int lua_setcstacklimit(lua_State *L, unsigned int limit);                 // 设置 C 栈限制

/* 调试信息结构体 */
struct lua_Debug
{
  int event;                  // 事件类型
  const char *name;           // (n) 函数名或变量名
  const char *namewhat;       // (n) 名称类型：'global'|'local'|'field'|'method'
  const char *what;           // (S) 函数类型：'Lua'|'C'|'main'|'tail'
  const char *source;         // (S) 源码路径
  size_t srclen;              // (S) 源码长度
  int currentline;            // (l) 当前行号
  int linedefined;            // (S) 函数定义起始行
  int lastlinedefined;        // (S) 函数定义结束行
  unsigned char nups;         // (u) 上值数量
  unsigned char nparams;      // (u) 参数数量
  char isvararg;              // (u) 是否可变参数
  char istailcall;            // (t) 是否为尾调用
  unsigned short ftransfer;   // (r) 第一个转移值的索引
  unsigned short ntransfer;   // (r) 转移值数量
  char short_src[LUA_IDSIZE]; // (S) 源码缩写（用于错误消息）
  /* 私有字段 */
  struct CallInfo *i_ci; // 当前执行的函数
};