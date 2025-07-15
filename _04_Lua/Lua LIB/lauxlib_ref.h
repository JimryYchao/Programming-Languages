// Lua 5.4.8
#pragma once
#include <stddef.h>
#include <stdio.h>
#include "lua.h"

/* 全局表名称 */
#define LUA_GNAME "_G"

/* 前向声明缓冲区结构 */
typedef struct luaL_Buffer luaL_Buffer;

/* 文件加载错误码 */
#define LUA_ERRFILE (LUA_ERRERR + 1)

/* 注册表中已加载模块表的键名 */
#define LUA_LOADED_TABLE "_LOADED"

/* 注册表中预加载加载器表的键名 */
#define LUA_PRELOAD_TABLE "_PRELOAD"

/* 版本检查 */
#define LUAL_NUMSIZES (sizeof(lua_Integer) * 16 + sizeof(lua_Number))
void luaL_checkversion_(lua_State *L, lua_Number ver, size_t sz); // 检查核心和库版本是否兼容
#define luaL_checkversion(L) luaL_checkversion_(L, LUA_VERSION_NUM, LUAL_NUMSIZES)

/* 元表操作 */
int luaL_getmetafield(lua_State *L, int obj, const char *e);    // 获取对象的元表字段
int luaL_callmeta(lua_State *L, int obj, const char *e);        // 调用对象的元方法
const char *luaL_tolstring(lua_State *L, int idx, size_t *len); // 将索引处的值转换为字符串

/* 参数检查 */
int luaL_argerror(lua_State *L, int arg, const char *extramsg);                 // 生成参数错误
int luaL_typeerror(lua_State *L, int arg, const char *tname);                   // 生成类型错误
const char *luaL_checklstring(lua_State *L, int arg, size_t *l);                // 检查字符串参数
const char *luaL_optlstring(lua_State *L, int arg, const char *def, size_t *l); // 获取可选字符串参数
lua_Number luaL_checknumber(lua_State *L, int arg);                             // 检查数字参数
lua_Number luaL_optnumber(lua_State *L, int arg, lua_Number def);               // 获取可选数字参数
lua_Integer luaL_checkinteger(lua_State *L, int arg);                           // 检查整数参数
lua_Integer luaL_optinteger(lua_State *L, int arg, lua_Integer def);            // 获取可选整数参数
void luaL_checkstack(lua_State *L, int sz, const char *msg);                    // 检查栈空间
void luaL_checktype(lua_State *L, int arg, int t);                              // 检查参数类型
void luaL_checkany(lua_State *L, int arg);                                      // 检查参数是否存在
#define luaL_checkstring(L, n) (luaL_checklstring(L, (n), NULL))
#define luaL_optstring(L, n, d) (luaL_optlstring(L, (n), (d), NULL))

/* 用户数据操作 */
int luaL_newmetatable(lua_State *L, const char *tname);         // 创建新元表
void luaL_setmetatable(lua_State *L, const char *tname);        // 设置元表
void *luaL_testudata(lua_State *L, int ud, const char *tname);  // 测试用户数据
void *luaL_checkudata(lua_State *L, int ud, const char *tname); // 检查用户数据

/* 错误处理 */
void luaL_where(lua_State *L, int lvl);             // 添加错误位置信息
int luaL_error(lua_State *L, const char *fmt, ...); // 抛出格式化错误

/* 选项检查 */
int luaL_checkoption(lua_State *L, int arg, const char *def, const char *const lst[]); // 检查选项参数

/* 文件和执行结果处理 */
int luaL_fileresult(lua_State *L, int stat, const char *fname); // 处理文件操作结果
int luaL_execresult(lua_State *L, int stat);                    // 处理执行结果

/* 引用系统 */
#define LUA_NOREF (-2)                         /* 无效引用 */
#define LUA_REFNIL (-1)                        /* nil 引用 */
int luaL_ref(lua_State *L, int t);             // 创建引用
void luaL_unref(lua_State *L, int t, int ref); // 释放引用

/* 代码加载 */
int luaL_loadfilex(lua_State *L, const char *filename, const char *mode);                            // 从文件加载代码
int luaL_loadbufferx(lua_State *L, const char *buff, size_t sz, const char *name, const char *mode); // 从缓冲区加载代码
int luaL_loadstring(lua_State *L, const char *s);                                                    // 从字符串加载代码
#define luaL_loadfile(L, f) luaL_loadfilex(L, f, NULL)

/* 状态管理 */
lua_State *luaL_newstate(void);              // 创建新的 Lua 状态机
lua_Integer luaL_len(lua_State *L, int idx); // 获取表或字符串的长度

/* 字符串处理 */
void luaL_addgsub(luaL_Buffer *B, const char *s, const char *p, const char *r);   // 缓冲区字符串替换
const char *luaL_gsub(lua_State *L, const char *s, const char *p, const char *r); // 全局字符串替换

/* 函数注册 */
typedef struct luaL_Reg
{
    const char *name;   /* 函数名称 */
    lua_CFunction func; /* 对应的 C 函数 */
} luaL_Reg;
void luaL_setfuncs(lua_State *L, const luaL_Reg *l, int nup);   // 注册函数表
int luaL_getsubtable(lua_State *L, int idx, const char *fname); // 获取或创建子表

/* 调试支持 */
void luaL_traceback(lua_State *L, lua_State *L1, const char *msg, int level); // 生成调用栈信息

/* 模块加载 */
void luaL_requiref(lua_State *L, const char *modname, lua_CFunction openf, int glb); // 加载并注册模块

/*
** ===============================================================
** 实用宏定义
** ===============================================================
*/

#define luaL_typename(L, i) lua_typename(L, lua_type(L, (i)))                                    // 类型名称宏
#define luaL_newlibtable(L, l) lua_createtable(L, 0, sizeof(l) / sizeof((l)[0]) - 1)             // 创建新库表
#define luaL_newlib(L, l) (luaL_checkversion(L), luaL_newlibtable(L, l), luaL_setfuncs(L, l, 0)) // 创建并注册新库
#define luaL_dofile(L, fn) (luaL_loadfile(L, fn) || lua_pcall(L, 0, LUA_MULTRET, 0))             // 执行文件
#define luaL_dostring(L, s) (luaL_loadstring(L, s) || lua_pcall(L, 0, LUA_MULTRET, 0))           // 执行字符串
#define luaL_getmetatable(L, n) (lua_getfield(L, LUA_REGISTRYINDEX, (n)))                        // 元表操作简化宏
#define luaL_opt(L, f, n, d) (lua_isnoneornil(L, (n)) ? (d) : f(L, (n)))                         // 可选参数处理宏
#define luaL_loadbuffer(L, s, sz, n) luaL_loadbufferx(L, s, sz, n, NULL)                         // 缓冲区加载简化宏
#define luaL_intop(op, v1, v2) ((lua_Integer)((lua_Unsigned)(v1)op(lua_Unsigned)(v2)))           // 整数运算宏（使用环绕语义）
#define luaL_pushfail(L) lua_pushnil(L)                                                          // 推送失败值
/* 参数检查宏 */
#define luaL_argcheck(L, cond, arg, extramsg) ((void)(luai_likely(cond) || luaL_argerror(L, (arg), (extramsg))))
#define luaL_argexpected(L, cond, arg, tname) ((void)(luai_likely(cond) || luaL_typeerror(L, (arg), (tname))))

/*
** =======================================================
** 缓冲区操作
** =======================================================
*/

struct luaL_Buffer // 缓冲区结构
{
    char *b;      /* 缓冲区地址 */
    size_t size;  /* 缓冲区大小 */
    size_t n;     /* 缓冲区中字符数 */
    lua_State *L; /* 关联的 Lua 状态 */
    union
    {
        LUAI_MAXALIGN;           /* 确保最大对齐 */
        char b[LUAL_BUFFERSIZE]; /* 初始缓冲区 */
    } init;
};

/* 缓冲区信息宏 */
#define luaL_bufflen(bf) ((bf)->n)
#define luaL_buffaddr(bf) ((bf)->b)

/* 缓冲区操作宏 */
#define luaL_addchar(B, c) ((void)((B)->n < (B)->size || luaL_prepbuffsize((B), 1)), ((B)->b[(B)->n++] = (c)))
#define luaL_addsize(B, s) ((B)->n += (s))
#define luaL_buffsub(B, s) ((B)->n -= (s))

/* 缓冲区函数 */
void luaL_buffinit(lua_State *L, luaL_Buffer *B);                 // 初始化缓冲区
char *luaL_prepbuffsize(luaL_Buffer *B, size_t sz);               // 准备缓冲区空间
void luaL_addlstring(luaL_Buffer *B, const char *s, size_t l);    // 添加字符串到缓冲区
void luaL_addstring(luaL_Buffer *B, const char *s);               // 添加字符串到缓冲区
void luaL_addvalue(luaL_Buffer *B);                               // 添加栈顶值到缓冲区
void luaL_pushresult(luaL_Buffer *B);                             // 将缓冲区内容推入栈
void luaL_pushresultsize(luaL_Buffer *B, size_t sz);              // 将缓冲区内容推入栈并设置大小
char *luaL_buffinitsize(lua_State *L, luaL_Buffer *B, size_t sz); // 初始化缓冲区并准备空间
#define luaL_prepbuffer(B) luaL_prepbuffsize(B, LUAL_BUFFERSIZE)

/*
** =======================================================
** 文件流处理
** =======================================================
*/

#define LUA_FILEHANDLE "FILE*" // 文件流类型标识
typedef struct luaL_Stream     // 文件流结构
{
    FILE *f;              /* 文件指针(NULL表示未完全创建的流) */
    lua_CFunction closef; /* 关闭函数(NULL表示已关闭的流) */
} luaL_Stream;

/*
** ===================================================================
** 消息和错误报告的基础抽象层
** ===================================================================
*/

#define lua_writestring(s, l) fwrite((s), sizeof(char), (l), stdout)           // 打印字符串
#define lua_writeline() (lua_writestring("\n", 1), fflush(stdout))             // 打印换行并刷新输出
#define lua_writestringerror(s, p) (fprintf(stderr, (s), (p)), fflush(stderr)) // 打印错误消息

/*
** =============================================================
** 兼容性定义（已弃用的类型转换）
** =============================================================
*/
#if defined(LUA_COMPAT_APIINTCASTS)
#define luaL_checkunsigned(L, a) ((lua_Unsigned)luaL_checkinteger(L, a))
#define luaL_optunsigned(L, a, d) ((lua_Unsigned)luaL_optinteger(L, a, (lua_Integer)(d)))
#define luaL_checkint(L, n) ((int)luaL_checkinteger(L, (n)))
#define luaL_optint(L, n, d) ((int)luaL_optinteger(L, (n), (d)))
#define luaL_checklong(L, n) ((long)luaL_checkinteger(L, (n)))
#define luaL_optlong(L, n, d) ((long)luaL_optinteger(L, (n), (d)))
#endif