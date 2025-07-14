## Lua 与 C 交互

### C API

C API 用于读写 Lua 全局变量、调用 Lua 函数、运行 Lua 代码段、注册 C 函数等。Lua 标准库将所有状态保存在动态结构体 `lua_State`。

Lua 和 C 通过虚拟栈 stack 通信和数据交换。常见的栈操作：
- 压入元素 `lua_push*`。`lua_checkstack` 检查栈空闲空间；栈中元素按先后压入顺序从 1 开始索引，-1 表示栈顶元素。
- 查询元素 `lua_is*`。这类函数检查栈中某个元素是否可以转换为特定类型。`lua_type` 返回栈中元素的类型例如 `LUA_TNIL`、`LUA_TBOOLEAN`、`LUA_TNUMBER`、`LUA_TSTRING` 等。
- 获取元素 `lua_to*`。转换失败不会提示类型错误，`lua_tolstring` 和 `lua_tothread` 返回 `NULL`。

> 一个简易 Lua 解释器
 
```c
#include <stdio.h>
#include "lua.h"                // CAPI 基础函数
#include "lauxlib.h"            // CAPI 辅助库函数
#include "lualib.h"             // Lua 标准库函数

int DoLuaBuff(char * buf, int bufsz)
{
    int error;
    lua_State *L = luaL_newstate();    // 创建一个 Lua state
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
    int top = lua_gettop(L); // 获取栈的深度
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
### 附录

#### CAPI Headers

- [lua.h](lua_ref.h) 定义 Lua 虚拟机（Lua State）的 C API 接口、基础数据类型和关键宏。

<!-- <details>
  <summary>展开查看代码</summary>
  <pre><code class="language-javascript">
  // 通过 JavaScript 动态加载文件内容
  fetch('./lua_ref.h').then(response => response.text()).then(text => {
    document.querySelector('lua_ref').innerText = text;
  });
  </code></pre>
</details> -->

<pre id="code-container">
    <code class="language-javascript">加载中...</code>
</pre>
<script>
    fetch('./lua_ref.h')
        .then(response => response.text())
        .then(text => {
            document.querySelector('#code-container code').textContent = text;
        });
</script>

---