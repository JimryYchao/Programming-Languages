## Lua 规范摘要

### 基本概念

#### 类型系统

Lua 是一种动态类型语言，通常作为嵌入式语言在宿主语言上下文中被调用。有 8 种内置基本类型 ***nil***, ***boolean***, ***number***, ***string***, ***function***, ***userdata***, ***thread***, ***table***。
- *number* 包括 64 位大小的 *integer* 和 *float*，整数溢出则发生环绕；
- *string* 是不可变的 UTF-8 字符序列；
- *function* 表示 Lua 函数和 C 函数；
- *userdata* 用于将任意 C 数据的原始内存块存储在 Lua 变量中，分为 *full userdata* (由 Lua 管理的占据内存的对象) 和 *light userdata* (C 指针值)；它们没有预定义操作，只能通过 C API 创建或修改；
- *thread* 表示独立的执行线程，用于实现协程。
- *table* 实现关联数组，是 Lua 中唯一的数据结构化机制，可用于表示数组、列表、集合、图、记录、字典、树等。与 *nil* 关联的任何键不视为表的一部分。

*table*，*function*，*thread*，*full userdata* 值为对象，变量保存对它们的引用。`type()` 返回给定值的类型字符串。

>---
#### 环境

对任何自由变量 `var` 的调用都转换为对 `_Env.var` 的调用，每个 Lua 块的变量都在块的外部变量 `_Env` 范围内编译。当加载一个 Lua 块时，它的 `_Env` 首先初始化为全局环境 `_G`；然后所有的标准库都加载到全局环境中，可以使用 `load` 或 `loadfile` 加载具有不同环境的块。在 C 中必须先加载块，然后更改（`C::lua_setupvalue`）其第一个 *upvalue* 的值。 

全局变量都可以在 `_G` 表中由变量名为键进行索引或更改。测试一个变量是否存在，不能简单的和 `nil` 比较，访问 `nil` 值会引发一个错误；可以利用 `rawget`：

```lua
function isExist(varName)
    if rawget(_G, varName) == nil then
        return false
    else return true
    end
end

if isExist("varName") then
   -- do with _G["varName"] 
end
```

起初，`_G` 和 `_ENV` 指向同一个表，创建的全局变量均可通过两者进行访问；`_ENV` 可以指向一个新的用户定义环境，并丢失之前的状态。

```lua  
a = 1     
local newgt = {}    -- 创建新环境继承旧环境
setmetatable(newgt, { __index = _G })
_ENV = newgt

print(a, _G.a)  -- 1    1
a = 10
_G.a = 20
print(a, _G.a)  -- 10   20
```

当加载另一个模块时，被加载模块的全局变量会自动进入当前环境。可以利用 `_ENV` 的定界特性，将模块进行分离：

```lua
-- M1.lua
local M = {}
_ENV = M
function func()
    <code>
end
return M;

-- M2.lua
local M1 = require "M1"
M1.func()   
```





>---
#### 模块与包

一个模块（Module）可以是由 Lua 或 C 编写的一些代码，这些代码通过 `require` 加载并创建一个表，这个表中的导出变量被加载至本地的环境变量 `_ENV` 中。每个模块仅加载一次，`package.loaded` 保存 `require` 加载的模块。

```lua
mod = require(moduleName)
---------------------------------
local m = require "mod"			--> 引入 mod 模块
local f = require "mod".func    --> 引入 mod 模块中的 func 函数
local sub = require "mod.sub"	--> 引入 mod.sub 子模块
```

`require` 首先搜索 `package.path` 指定的路径检查模块是否存在，并通过 `loadfile` 对其进行加载。未找到时搜索 `package.cpath` 并通过 `package.loadlib` 进行加载。返回的加载函数具有返回值时，`require` 返回这个值并保存到 `package.loaded` 中。强制 `require` 加载同一模块两次，可以先将模块从 `package.loaded` 中移除（`package.loaded["mod"] = nil`）

```lua
local mod = require "Mod"    --> 首次加载
package.loaded["Mod"] = nil  --> 卸载 Mod
mod = require "Mod"          --> 再次加载
```

搜索路径是一组模板，其中的每项都指定了将模块名转换为文件名的方式。

```lua
package.path = [[?;?.lua;c:windows\?;/usr/local/lua/?/?.lua]]
local mod = require "sql"	-- 尝试搜索
    --[[
        spl
        sql.lua
        c:\windows\sql
        /usr/local/lua/sql/sql.lua
    ]]

package.cpath = [[.\"?.dll;C:\Program Files\Lua504\dll\?.dll]]
```

> *构建模块*

```lua
-- Mod.lua : 创建一个模块的一般方法
local M = {}
M.Add = function(c1, c2) end
M.Sub = function(c1, c2) end
M.Mul = function(c1, c2) end
M.Inv = function(c1, c2) end
M.Div = function(c1, c2) end
return M
-------------------------------
-- other.lua : 加载 Mod 模块
local mod = require "Mod"
```


>---
#### 编译与执行

Lua 作为解释型语言，可以在运行代码前执行预编译。`dofile(filename)` 是执行 Lua 代码段的主要方式之一。`loadfile` 与 `load` 执行预编译并返回一个函数，利用 `assert(loadfile(...))` 断言预编译是否发生错误。`load` 总是在全局环境中执行预编译，因此不涉及词法定界，常用于执行外部代码或动态生成的代码：

```lua
i = 32
local i = 0
f = assert(load("i = i + 1; print(i)"))
g = function() i = i + 1; print(i) end
f()    -- 33
g()    -- 1
```

可以利用 *luac* 程序生成 Lua 预编译文件。`load` 和 `loadfile` 也可以接受预编译代码。

```lua
-- shell
-- $ luac -o <name.lc> <source.lua>
assert(loadfile("name.lc"))()
```


>---
#### 错误处理

错误将中断程序正常流程。可以调用 `error(mess, level)` 显式抛出错误并沿着堆栈进行传播；`assert(exp, mess)` 执行断言并在 `exp` 为 `false` 或 `nil` 时抛出错误；函数 `pcall(func, args...)` 和 `xpcall(func, handle, args...)` 用于捕获错误，常用 `debug.debug` 和 `debug.traceback` 作为 `xpcall` 的消息处理函数。函数 `warn` 用于生成一条警告消息。

```lua
function panic(msg)
    error(msg, 2)
end

local ok, err = pcall(panic, "Lua panic");
if not ok then
    print("ERROR: " .. err)
end

function handle(errno)
    print("errno: " .. errno)
end
-- 同 pcall 并附带一个消息处理程序
ok, stat = xpcall(panic, handle, 1)
```


>---
#### 词法元素

> **关键字**

```lua
and,or,not             -- 逻辑运算
local                  -- 局部变量声明
<const>                -- 局部变量常量限定
<close>                -- 局部变量关闭限定
self                   -- 表函数自引用
_                      -- 弃元  
nil                    -- 空值
true,false             -- 布尔
function               -- 函数声明
break,goto,return      -- 跳转语句
for,in                 -- for
while-do,repeat-until  -- 循环语句
if,else,elseif         -- 条件语句 if
then                   -- 循环或条件语句的块开始
end                    -- 语句块结束
```

> **操作符**

```lua
+    -   *   /   %   ^   #
&    ~   |   <<  >>  //  
==   ~=  <=  >=  <   >   =
(    )   {   }   [   ]   ::
;    :   ,   .   ..  ...
```

---
### 类型与声明

Lua 中有 8 个基本类型分别为：nil（空）、boolean（布尔）、number（数值）、string（字符串）、userdata（用户数据）、function（函数）、thread（线程） 和 table（表）。

简单值类型包含 *nil*，*boolean*，*number*，*string*。`local` 限定局部变量：
- 未定义变量值为 `nil`。`nil` 可用于将某个不再使用的变量或 *table* 的键置空，Lua 会自动回收该变量。
- `false` 和 `nil` 的布尔条件测试返回假，其他值返回真。
- 数值类型内置整数和浮点数，算数值相同的整数和浮点数被视为相等，`math.type(n)` 返回数值内部类型；浮点数支持 E 和 P 计数法。    


```lua
type(nil)       --> nil
type(true)      --> boolean
type(1234)      --> number
type("Hello")   --> string
type(io.stdin)  --> userdata
type(print)     --> function
type({})        --> table
type(type(X))   --> string
```

> *浮点型与整型之间的转换*

算数值相等的整数和浮点数之间可以互相转换：整数加 0.0 转换为浮点数；浮点数与 0 按位或运算转换为整数。

```lua
12345 + 0.0    --> 12345.0
12345.0 | 0    --> 12345
2^53           --> 9.007199254741e+15	(浮点型值)
2^53 | 0       --> 9007199254740922	(整型值)
-- error
3.2 | 0        --> 小数部分 > 0
2^64 | 0       --> 超出范围
```

> 局部变量与常量属性

`local` 声明与块范围关联的局部变量，局部变量无法通过 `_ENV.var` 访问，但可以作为模块的返回值；`<const>` 为局部变量赋予常量属性。

```lua
local m = {}
local C <const> = 10086
m.v = C
return m
```

>---
#### string

支持转义字符串 `"string"` 或 `'string'` 和原始字符串 `[=[string]=]`。`#str` 返回字符串字节数，`x .. y` 拼接字符串，操作数可以是字符串或数值。Lua 运行时提供了数值和字符串之间的自动转换。算术操作时，Lua 会预先尝试将字符串转换成数值类型。

```lua
a = "a 'line\n'"
b = 'another "line\n"' .. "\255" .. "\xff" .. "\u{7fffffff}"
m = [=[
first line\n
second line    
]=]
```

> *转义字符*

```lua
\a           --> 响铃
\b           --> 退格
\f           --> 换页
\n           --> 换行
\r           --> 回车
\t           --> 水平制表符
\v           --> 垂直制表符
\\           --> 反斜杠
\'           --> 单引号
\"           --> 双引号
\0           --> 空字符
\ddd         --> \000 ~ \255
\xhh         --> \x00 ~ \xff
\u{h…h}      --> \w{00000000} ~ \u{7fffffff}
\[,\]        --> 方括号
```

`\z`：忽略随后任意数目的空白字符直到第一个非空白字符。

```lua
"12345\z    6"  --> 123456
"12345\z
678"            --> 12345678
"123\z \r678"   --> 123\r678
```

>---
#### table

表（table）是 Lua 中唯一的数据结构机制，可以用来表示数组、列表、符号表、集合、记录、图形、树等数据结构。Lua 使用 `_G` 表用来存储全局变量。键可以是除 `nil` 和 `NaN` 外的任何值，值为 `nil` 任何键都不被视为表的一部分。

```lua
t1 = {
    K1 = "V1", ["K2"] = "V2",
    [1] = 1, [2] = 2, 3, 4, 5
}

a = {}			
a["K1"] = "V1"	
a.K2 = 999	
a[1] = "Great"	
a.__index = a	
```

表构造器用来创建并初始化表的表达式，对于列表式表元素的声明，从索引 1 开始为每个列表元素创建索引关联。

```lua
day = {"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
-- 索引从 1 开始为每个列表元素建立索引关联

print(day[4])	--> Wednesday
```

索引可以显式指定，在计算列表的边界时会忽略小于 1 的索引。显式声明的列表，`#table` 边界值受列表是否存在连续 `nil` 值的影响。列表出现连续 `nil` 值时，`#table` 返回连续空洞前的索引。

```lua
arr = {[1]=1,[2]=2,[4]=4,[5]=5,[7]=7,[9]=9}
-- #arr = 5
Arr = {[1]=1,[2]=2,[4]=4,[5]=5,[7]=7,[8]=8}
-- #Arr = 8
```

自动索引列表的边界返回不确定，列表末端的 `nil` 值会被忽略而不计入边界计算，对于中间存在空洞 `nil` 的列表而言，`#table` 是不可靠的。

```lua
a1 = {1,2,3,4,5,6}           --> #a = 6
a2 = {1,2,3,nil,nil}         --> #a = 3
a3 = {1,nil,nil,nil,nil,2}   --> #a = 6
a4 = {}
a4[1] = 1; a4[2] = 2; a4[4] = 4  --> #a = 4, 中间存在空洞
```


>---
#### function

Lua 函数是第一类值，函数定义实质上就是创建类型为 `function` 的值并赋值给变量

```lua
-- 函数声明
function Func( <params> )
    body | return multi -- 多值返回
end
Func = function Decl

-- 匿名函数
function GetCounter() 
    local v = 0
    return function()   -- 函数闭包
        v = v+1
        return v
    end
end
counter = GetCounter()
v = counter()  -- 1
```

当函数只有一个参数且该参数是字符串常量或表构造器时，函数调用表达式 `f()` 括号是可选的。

```lua
print "Hello"  -- print("Hello")
type {}        -- type({})
```

> *变长参数*

变长参数 `...` 在函数内部可以利用 *table* `{...}` 进行收集，或利用多重赋值按顺序提取；`table.pack` 检测参数中是否有 `nil`。

```lua
function Foo(a, ...)
    local t = { ... }
    if (#t > 0) then
        for i = 1, #t do
            print(t[i])
        end
    end
    print("end :" .. a)
end
Foo(1, 2, 3, 4, 5, 6)
```


另一种访问变长参数的方法是利用 `select(n,...)`，n 为数值时，返回第 n 个参数后的所有参数；n 是 `"#"` 时，返回额外参数的总数。

```lua
-- 打印奇数位的元素
function Foo(...)
    local len = select("#", ...)
    local i = 1
    local a = 0
    while i <= len do
        a = select(i, ...)
        print(a)
        i = i + 2
    end
end
Foo(1,2,3,4,5,6,7,8)	-- 1 3 5 7
```

> *表函数*

表调用自身的表函数成员时，可以通过两种方式调用：`table.fun()` 或 `table:fun()`。`table:fun` 默认将表自身作为第一个参数传递给表函数。

```lua
t = {1,2,3,4,5,6}
---@param t table
function t.Traverse(t)
    for i =1,#t do
        print(t[i])
    end
end

t:Traverse()	-- 传递自身作为首位参数到表函数中
```

外部声明表函数的通过 `:` 声明，表示默认将表自身作为第一个参数传入函数：

```lua
function T:Traverse()
-- 等价于
function T:Traverse(self)
```

>---
#### thread 与 coroutine

从多线程的角度看，协程（coroutine）与线程（thread）类似：协程是一系列的可执行语句，拥有自己的栈、局部变量和指令指针，同时协程又与其他协程共享了全局变量和其他几乎一切资源。

线程与协程的区别在于，多线程程序可以并行运行多个线程，协程需要彼此协作，任意指定的时刻只能有一个协程运行，正在运行的协程被挂起时其执行才会暂停。可以使用 `coroutine.create(function)` 创建一个新协程，返回一个 `thread` 类型。

一个协程有四种状态：挂起（suspended）、运行（running）、正常（normal）、死亡（dead）。`coroutine.status(co)` 查看协程对象状态。

当一个协程创建时，它不会自动运行而处于挂起（suspended）状态，利用 `coroutine.resume(co)` 用于启动或恢复一个协程，并改状态为运行（running）；若协程体运行之后就终止了，它的状态转变为死亡（dead）。当协程 A 唤醒协程 B 时（执行权移交给 B），A 会变成正常状态（normal），而协程 B 会变成运行（running）。

```lua
co = coroutine.create(<function(params)>)
print(coroutine.status(co))     -- suspended

coroutine.resume(co [,params])  -- running
print(coroutine.status(co))     -- dead
```

`coroutine.yield()` 可以让一个运行中的协程挂起，之后在 `resume` 后恢复运行，协程会继续执行直到遇到下一个 `yield` 或执行结束。`resume` 已经结束的协程（dead 状态）会返回 `false` 和一条信息（`"cannot resume dead coroutine"`）。协程中通过 `resume`，`yield` 在主函数与协程之间来交换数据。

```lua
co = coroutine.create(function()
    local count = 0
    while true do
        count = count + 1
        print("yield :", coroutine.yield(count))
    end
end)
print("resume:" , coroutine.resume(co, "hi"))
print("resume:" , coroutine.resume(co, 1, 1, 1))
print("resume:" , coroutine.resume(co, "hello"))
--[[
    resume::        true    1
    yield::         1       1       1
    resume::        true    2
    yield::         hello
    resume::        true    3
]]
```

`coroutine.wrap(f)` 返回一个 `function` 类型的协程，和 `create(co)` 构造的 `thread` 区别在于，`coroutine.yield` 或协程结束返回时，不会返回函数是否正常运行或恢复运行的状态，也无法获得 `function` 协程的状态。

`thread` 协程的主函数发生错误时不会终止程序，会将错误发送到 `resume` 的返回中；而 `function` 协程直接将导致程序错误

```lua
local co = coroutine.create(f)
local cf = coroutine.wrap(f)
-- 调用 thread 协程
coroutine.resume(co)
-- 调用 function 协程
cf()
```

> *将协程用作迭代器*

```lua
function permgen(a, n)
    n = n or #a
    if n <= 1 then
        coroutine.yield(a)
    else
        for i = 1, n do
            a[n], a[i] = a[i], a[n]
            permgen(a, n - 1)
            a[n], a[i] = a[i], a[n]
        end
    end
end

function printResult(a)
    for i = 1, #a do
        io.write(a[i], " ")
    end
    io.write('\n')
end

function permutations(a)
    local co = coroutine.create(function()
        permgen(a)
    end)
    return function()
        local code, res = coroutine.resume(co)
        return res
    end
end

for a in permutations { 1, 2, 3 } do
    printResult(a)
end
```

> *coroutine.running*

函数 `coroutine.running(co)` 返回正在运行的协程和一个 *boolean* 值，当正在运行的协程是主函数 `main` 时返回 `true`。

```lua
print(coroutine.running())
local co = coroutine.create(function()
    print("join in co")
    local c, ismain = coroutine.running()
    print(c, ismain, "\nco : " .. tostring(co))
end)
coroutine.resume(co)
--[[
    thread: 0000017C2872F458        true
    join in co
    thread: 0000017C287A43E8        false
    co : thread: 0000017C287A43E8
]]
```

> *coroutine.close*

函数 `coroutine.close(co)` 用于关闭待关闭的 suspended 或 dead 状态的协程并返回 `true`；关闭正在运行的协程会发生错误并返回 `false` 和错误信息。

```lua
local co = coroutine.create(function()
    print(coroutine.close(co)) -- cannot close a running coroutine
    print("join in co")
end)
coroutine.resume(co)
```

> *coroutine.isyieldable*

函数 `coroutine.isyieldable(co?)` 用于判断协程是否是可让步（yield）的。如果协程不是主线程，也不在不可让步 C 函数中，则该协程是可让步的。

```lua
-- main
print(coroutine.isyieldable())     -- false, 主线程

local co = coroutine.create(function()
    print("join in co")
end)
print(coroutine.isyieldable(co))   -- true
```

>---
#### 局部属性声明

> *const*

`<const>` 属性赋予局部变量常量属性，无法赋值操作，但是不影响作为常量表成员的任何操作。

```lua
local t <const> ={ a = 1}
t = nil 	-- error
t.a = 2		-- ok
```

> *close*

一个 `to-be-closed` 对象的行为类似于一个局部常量（无法重新定义），赋予 `<close>` 属性的值必须具有 `__close` 字段关联的元方法或 `false`。局部变量生存期结束时，将按声明的相反顺序依次调用 `__close`。

一个协程被挂起时且永远不会被恢复时，设定 `__close` 的对象值生存期被无限延长，因此它们也不会被关闭。可以调用 `coroutine.close(co)` 或调用 `__gc` 关闭这些变量。


```lua
function new_thing()
    local thing = {}
    setmetatable(thing, {
        __close = function()
            print("thing closed")
        end
    })
    return thing
end

do
    local x <close> = new_thing()
    print("use thing")
end
-- use thing
-- thing closed
```


---
### 表达式

| Category | Operators |
| :------- | :-------- |
|幂运算| `x ^ y` |
|一元| `-x` `#string` `#table` `not x` `~x` |
乘法|`x * y` `x / y`(浮点除法) `x // y`(向下取整除法) `x % y`|
加法 | `x + y` `x - y`
字符串拼接| `x .. y`|
移位 |`x << y` `x >> y`|
按位 |`x & y` `x ~ y`(异或) `x \| y`|
关系 |`x < y` `x > y` `x <= y` `x >= y` `x == y` `x ~= y`| 
逻辑 | `x and y` `x or y`|
赋值 | `x = y` |


`x // y` floor 触发对商向负无穷取整。

```Lua
3 // 2         --> 1
3.0 // 2       --> 1.0
-9 // 2        --> -5
1.5 // 0.5     --> 3.0
```

`x % y` 取模运算的定义 ```x%y = x-((x//y)*y)```，其结果的符号与第二操作数一致。可以利用取模运算保留浮点运算有效位。

```lua
-9 % 2      --> 1:  -9-(-5)*2 = 1
-9 % 2.0    --> 1.0
math.pi - math.pi % 0.0001	--> 3.1415
```

`and` 与 `or` 支持短路原则。`and` 在第一个操作数为 `false` 或 `nil` 时返回第一个操作数，否则返回第二个操作数；`or` 在第一个操作数为 `false` 或 `nil` 时返回第二个操作数，否则返回第一个操作数。

```lua
10 or 20            --> 10
10 and 20           --> 20
nil or "a"          --> "a"
false and nil       --> false
false or nil        --> nil
```

可以利用 `and` 和 `or` 机制构造三目运算：

```lua
X and Y or Z
--> X == true --> Y
--> X == false --> Z
```

移位是逻辑移位，以 0 补齐空位

```lua
12(10) = 00001100(2)
00001100 >> 1 = 00000110  --> 12>>1 = 6
00001100 << 2 = 00110000  --> 12<<2 = 48

-18(10) = 11101110(2)     -- 补码表示
11101110 << 1 = 11011100  --> -18<<1 = -36
11101110 >> 2 = 01111011  --> -18>>2 = 123 (逻辑移位)
```

利用 floor 除法模拟实现算术移位，公式为 `num // (2^n)|0`，当 $n>0$ 表示算术右移；当 $n<0$ 表示算术左移。

```lua
-- 负数的算术右移
    -10 >> 2 等价于 -10//2^2|0 --> -3
-- 算术左移
    -10 >> -2 == -10 << 2 == -10//(2^-2|0) --> -40
```

---
### 语句

####  条件控制：if

`false` 和 `nil` 值为假，`true` 和非 `nil` 任意值为真。

```lua
if <condition> then
    <code>
elseif <condition> then
    <code>
else
    <code>
end
```

>---
#### 迭代语句：while, repeat, for

```lua
while <condition> do
    <code>
end
```
```lua
repeat
    <code>
until <condition>
```

> for

```lua
for i = Begin, End [, Step = 1] do
    <code>
end
```
```lua
arr = { 1,2,3,4,5 }
for i = 1, #arr do  
    print(arr[i]) 
end
```

> for-in

```lua
for	<varList> in <expList> do   -- varList 第一个值为 nil 循环停止
    <body>
end
-- 相当于
do
    local _f, _s, _var = explist   -- 返回一个迭代函数，不可变状态，控制变量初始值
    while true do
        local var_1, ... var_n = _f(_s, _var)
        _var = _var_1
        if _var == nil then break end
        <body>
    end
end
```
```lua
function func(maxCount,value) 
    if value < maxCount then
    value = value+1
    return value,value*2
    end
end
for i,v in func,5,0 do
    print(i,v)
end
```

> for-in pairs 

```lua
arr = {1,3,4,5,a="A",b="B"}
for k,v in pairs(arr) do	
    print(k,v)
end
```

> for-in ipairs

```lua
arr = {1,3,4,5,a="A",b="B"}
for i,v in ipairs(arr) do	-- 列表遍历
    print(i,v)	-- 1,3,4,5
end
```

>---
#### 跳转语句：break, return, goto

`break` 中断最内层循环语句；`goto` 无条件跳转当前范围的标签处；`return` 从当前范围返回零到多个值。


```lua
function Factorial(n)
    local rt = 1;
    ::start::
    if n == 0 then
        return rt
    else
        rt = rt * n
        n = n - 1
        goto start
    end
end

local i = 1
while (true) do
    if (i > 10) then
        break
    end
    print(i .. "  " .. Factorial(i))
    i = i + 1
end

return 0   -- 文件范围返回值
```

>---
#### 代码块：do-end

`do end` 可以在文件或函数域出现

```lua
do
    <code>
end
```

---
### 元表与元方法

元表定义了其原始值允许的某些操作，可以设置元表中特定元方法的字段来更改值行为。`getmetatable(t)` 获取父级元表。`rawget` 访问查询元表中的元方法。`setmetatable(t, metatable)` 替换表的元表，`metatable` 为 `nil` 时，表示删除 `t` 的元表。*table* 和 *full userdata* 具有单独的元表，除字符串外其他的值类型默认没有元表。

```lua
local subTable, father = {}, {}
setmetatable(subTable, father)
print(getmetatable(subTable) == father)   -- true
```

元表中的元方法可包括：

```lua
-- 算术运算
__unm       -- -x
__add       -- x + y
__sub       -- x - y
__mul       -- x * y
__div       -- x / y
__mod       -- x % y
__pow       -- x ^ y  
__idiv      -- x // y
__bnot      -- ~x
__band      -- x & y
__bor       -- x | y
__bxor      -- x ~ y
__shl       -- x << y
__shr       -- x >> y
__concat    -- x .. y
__len       -- #x
-- 关系运算
__eq        -- x == y
__lt        -- x < y
__le        -- x <= y
-- 表相关
__index     -- table[key]
__newindex  -- table[key] = value
-- 函数关联
__call      -- func(args)
__tostring  -- tostring 调用
__name      -- tostring 替选调用 
__gc        -- 终结器
__close     -- 待关闭变量
__mode      -- 弱表模式
__metatable -- 元表
__pairs     -- 在 for-pairs 替选调用
```

每种运算符都有一个对应的元方法。例如对两个表进行算术操作 `a + b` 时，首先查找第一个操作数的 `__add` 元方法并尝试调用；否则查找第二个操作数的 `__add`；否则抛出异常。

```lua
local mt = {}
function mt.__add(a, b)
    local set = {}
    local len = #a > #b and #a or #b
    for i = 1, len do
        set[i] = a[i]
    end
    for j = 1, #set do
        set[j] = a[j] + b[j]
    end
    return set
end

local t1 = { 1, 2, 3, 4, 5, 6 }
local t2 = { 6, 5, 4, 3, 2, 1 }
setmetatable(t1, mt)

local newt = t1 + t2
for i = 1, #newt do
    print(newt[i]) -- 7 7 7 7 7 7
end
```

通常会将 `a <= b` 作为其他关系元方法的基方法。

```lua
mt.__le = function(a, b)
    for k in pairs(a) do
        if not b[k] then return false end
    end
    return true
end
mt.__lt = function(a, b)
    return a <= b and not (b <= a)
end
mt.__eq = function(a, b)
    return a <= b and b <= a
end
```

>---
#### __index, __newindex

Lua 提供了一种改变表在访问和修改表中不存在字段这两种行为的方式。当访问不存在键时，这些访问首先查找 `__index` 元方法，默认返回 `nil`。`__newindex` 用于表的更新，当对不存在键赋值时，解释器会查找 `__newindex` 元方法。`rawget` 与 `rawset` 使用原始方式对 *table* 执行操作。

```lua
local Array = {}   -- 固定长度数组
Array.New = function(...)
    local arr = { ... }
    local list = { Length = #arr, array = arr }
    setmetatable(list, Array)
    return list
end
Array.__index = function(self, index)
    if index > self.Length or index < 1 then
        error("index is out of range.")
    end
    return self.array[index]
end
Array.__newindex = function(self, index, value)
    if index > self.Length or index < 0 then
        error("index is out of range.")
    end
    self.array[index] = value
end

local a = Array.New(1, 2, 3, 4, 5, 6)
print(a.Length)

for i = 1, a.Length do
    print(a[i])
end
```

>---
#### __call

`__call(args...)` 为表创建函数调用表达式语法。

```lua
local _Log = { "DEBUG", "INFO", "WARNING", "ERROR", "PANIC", }
function _Log:__call(mess, level)
    if not level then
        level = self.level
    end
    local pre = self[level]
    print((pre and (pre .. ": ") or "") .. mess);
end
function _Log:SetLevel(level)
    self.level = level
end
function NewLogger(level)
    local logger = { level = level }
    _Log.__index = _Log
    setmetatable(logger, _Log)
    return logger
end

Log = NewLogger(1)

Log("Hello")      -- DEBUG: Hello
Log:SetLevel(3)
Log("World")      -- WARNING: World
```

>---
#### __tostring, __name

函数 `print` 总是调用 `tostring` 来进行格式化输出，`tostring` 首先会查找对象是否拥有 `__tostring`；否则查找 `__name` 作为替代。

```lua
local mt = {}
print(mt) -- table: 00000111D4579E30
-- print 调用了表的元方法 __tostring
local t = { 4, 10, 2 }
mt.__tostring = function(self)
    if #self == 0 then
        return tostring(self)
    end
    local s = self[1]
    for i = 1, #self do
        if i ~= 1 then
            s = string.format(s .. ',' .. tostring(self[i]))
        end
    end
    return s
end
setmetatable(t, mt)
print(t) -- 4, 10, 2
t = {}
print(t) -- table: 000002A055677EE0
```

>---
#### __gc

元表中拥有 `__gc` 字段时，子表在 `setmetatable` 之后会被标记为可析构。如果设置元表时没有 `__gc` 字段，之后在元表中创建该字段，子表对象不会被标记为可析构。

标记为终结的对象在垃圾回收阶段，收集器会调用（仅一次）该对象的 `__gc` 元方法，并且只会被调用一次，即使在垃圾回收阶段被永久复苏。

```lua
local mt = {}
local t = nil
-- t = setmetatable({}, mt)
mt.__gc = function ()
	print("call __gc ...")
end
t = setmetatable({}, mt)
t = nil
collectgarbage("collect")
print("end ...")
--[[
call __gc ...
end ...
]]
```

>---
#### __close

当一个可关闭对象离开其作用域时，将调用该对象的 `__close` 元方法，这为手动释放一些资源提供了可行的方式。

```lua
function new_thing()
    local thing = {}
    setmetatable(thing, {
        __close = function()
            print("thing closed")
        end
    })
    return thing
end

do
    local x <close> = new_thing()
    print("use thing")
end
-- "thing closed" is printed here after "use thing"
```

>---
#### __mode

`__mode` 字段与弱引用相关，可以设置表为弱引用表，`__mode` 的值确定了弱引用表的类型，`"k"` 表示键弱引用，`"v"` 表示值弱引用，`"kv"` 表示键值弱引用。对于一个弱引用表，只要一个弱引用的键或值被回收，就把这个键值对从表中删除。

```lua
local mt = { __mode = "v" } -- 值弱引用
local a = {}
local b = { key = a }
setmetatable(b, mt)
print(b.key)      -- table: 000001F0DD803580
a = nil
print(b.key)      -- table: 000001F0DD803580
collectgarbage("collect")  -- 强制垃圾回收
print(b.key)      -- nil
```

>---
#### __metatable

函数 `setmetatable` 和 `getmetatable` 用到了元方法，用于保护元表。假设要保护集合，即使用户既看不到也不能修改集合的元表，可以在元表中设置 `__metatable` 字段，那么 `getmetatable` 会返回这个字段，`setmetatable` 会引发一个错误。

```lua
local mt = { __metatable = 0}
local t = setmetatable({}, mt)  -- 限制 t 受保护
print(getmetatable(t))          -- mt.__metatable = 0
local ok, err = pcall(setmetatable, t, {})
if not ok then
    print(err)  -- cannot change a protected metatable
end
```

>---
#### __pairs

函数 `pairs` 对应元方法 `__pairs`，指定表在 for 迭代器中的迭代行为。

```lua
local mt = {}
local default = {
    ["name"] = "Lua",
    ["telephone"] = 123456,
    ["id"] = 7,
}
mt.__index = function(tbl, key)
    local val = rawget(tbl, key)
    return val and val or default[key]
end

mt.__pairs = function(tbl, key)
    return function(t, k)
        local nk, nv = next(default, k)
        if nk then
            nv = t[nk]
        end
        return nk, nv
    end, tbl, nil
end

local test = setmetatable({ id = 8 }, mt)
-- test __pairs 
for k, v in pairs(test) do
    print(v)
end
-- test __index 
print(test.id)
print(test.name)
print(test.telephone)
```

---
### 面向对象编程

使用参数 `self` 是所有面向对象语言的核心点。避免使用全局变量进行操作，把操作限定给特定对象工作。可以利用表的与值无关的 `self` 表示作为操作的接受者，来避免将操作仅限定在特定的全局变量中才能工作。

```lua
function t.foo(self, arg)    -- self 声明
t.foo(t, arg)	-- 等价于 
t:foo(arg)		-- 表示隐藏 self
-- 即使 t 被置换成别名，该操作和表对象本身无关
```

Lua 可以利用原型的概念去实现面向对象编程，利用元表 `__index` 继承的方式（`setmetatable(A,{__index = B})`），让 B 成为 A 的一个原型。在此之后，A 就会在 B 中查找它没有的操作或字段。

```lua
local prototype = {}
prototype.__index = prototype
prototype.FuncA = function ()
    print("Call FuncA")
end 
prototype.FuncB = function ()
    print("Call FuncB")
end 
prototype.__close = function ()
    print("closed")
end

function NewPrototype()
    return setmetatable({}, prototype)
end

local p <close> = NewPrototype()
p.FuncA()
p.FuncB()
```

>---
####  继承

利用 `__index` 和 `self` 机制可以用来实现继承。

```lua
local metaclass = {}
---Create an object from metaclass
function metaclass:new(o)
    o = o or {}
    self.__index = self
    setmetatable(o, self)
    o.base = self			-- 关联超类
    return o
end

local o1 = metaclass:new()
o1.key1 = "hello"
local o2 = o1:new()	-- 单一继承
print(o2.key1)	-- hello
metaclass.key2 = 1
print(o2.key2)	-- 1
```

> 多重继承

多重继承意味着一个类可以具有多个超类，因此需要一个独立的方法（createClass）从一个类中创建子类，其参数为新类的所有超类

```lua
--- 在表 plist 的列表中查找 k
local function search(k, plist)
    for i = 1, #plist, 1 do
        local v = plist[i][k]
        if v then return v end
    end
end
--- 多重继承
function createClass(...)
    local c = {}
    local parent = { ... }
    setmetatable(c,
        { __index = function(t, k)
            return search(k, parent)
        end })

    function c:new()	-- 继承模式
        o = o or {}
        self.__index = self
        setmetatable(o, self)
        return o
    end
    return c
end
```

---
### 垃圾回收

Lua 使用自动内存管理。一般而言，垃圾收集器是不可见的。弱引用表（weak table）、析构器（finalizer）、函数 `collectgarbage` 是 Lua 中用来辅助垃圾收集器的主要机制。
  - 弱引用表允许收集 Lua 中还可以被程序访问的对象；
  - 析构器允许收集不在垃圾收集器直接控制下的外部对象；
  - 函数 `collectgarbage` 允许控制垃圾收集器的步长。

>---
#### 弱引用表

弱引用是一种不在垃圾收集器考虑范围内的对象引用。对于一个弱引用表，只要一个弱引用的键或值被回收，就把这个键值对从表中删除；一个表是否为弱引用表是由 `__mode` 字段决定，`"k"` 表示这个表的键是弱引用的，`"v"` 表示值是弱引用的，`"kv"` 表整个键值对是弱引用。只有对象可以从弱引用表中移除，值类型数值、布尔、字符串等不可回收，除非它们相关联的对象（可以是表、函数等）被回收。

```lua
local mt = { __mode = "v" } -- 值弱引用
local a = {}
local b = { key = a }
setmetatable(b, mt)
print(b.key)      -- table: 000001F0DD803580
a = nil
print(b.key)      -- table: 000001F0DD803580
collectgarbage()  -- 强制垃圾回收
print(b.key)      -- nil
```

垃圾回收中一种棘手的情况是，一个具有弱引用键的表中的值又引用了对应的键，从而弱引用键总是会被关联。一个典型的示例是常量函数工厂（参数是一个对象，返回值是一个被调用时返回传入对象的函数）

```lua
t[o] = function return (function() return o end) end
-- 值中引用了键本身
```

从弱引用表的概念中，每一个函数都指向其对应的键对象，因此对于每一个键来说都存在一个强引用。所以即使有弱引用的键，这些对象也不会被回收。一个具有弱引用键和强引用值的表是一个瞬表，在一个瞬表中，一个键的可访问性控制着对应值的可访问性。使用记忆技术的常量函数工厂，表 `mem` 中与一个对象关联的值回指了自己的键（对象本身）。

```lua
do
    local mem = {}
    setmetatable(mem, {__mode = "k"})
    function factory(o)
        local res = mem[o]
        if not res then
            res = (function() return o end)
            mem[o] = res	-- 对象和常量函数关联
        end
        return res
    end
end
```

对于瞬表的一对 `(k,v)`，指向 `v` 的引用只有当存在某些指向 `k` 的其他外部引用存在时才是强引用，否则即使 `v` 直接或间接地引用了 `k`，垃圾收集器最终会收集 `k` 并把元素从表中移除。

>---
#### 析构器

垃圾收集器的目标是回收对象，析构器在对象即将被回收时调用。在每个垃圾收集周期，垃圾收集器会在调用析构器之前清理弱引用表中的值，在调用析构器之后再清理键。

Lua 通过元方法 `__gc` 实现析构器。通过给对象设置一个具有非空 `__gc` 元方法的元表，就可以把一个对象标记为可析构。析构器不能产生或运行垃圾收集器，因为它们可以在不可预知的时间内运行。运行析构器时的任何错误都会产生警告，且错误不会传播。

```lua
o = {x = "hi"}
setmetatable(o, {__gc = function(o) print(o.x) end})
o = nil
collectgarbage()    --> hi
```

有关析构器的一个特性是复苏（resurrection）。当一个析构器被调用时，它的参数是正在被析构的对象。这个对象在析构期间重新变成活跃（临时复苏），如果该对象在析构器返回后仍然可访问（析构时保存到全局变量时），这个对象就变成永久复苏。

由于复苏的存在，Lua 会在两个阶段中回收具有析构器的对象。当垃圾收集器首次发现某个析构器的对象不可达时，垃圾收集器把这个对象复苏并将其放入等待被析构的队列中，一旦析构器开始执行，Lua 就将该对象标记为已被析构，当下一次垃圾收集器又发现这个对象不可达时，它就将这个对象删除。

每一个对象的析构器都只会精确地运行一次。程序结束时，Lua 也会显式的调用所有未被释放对象的析构器。因此可以利用在程序结束时的这种特性实现某种形式的 `atexit()`。可以将这个特殊的表锚定在全局表中。

```lua
local t = {__gc = function()
    -- 'atexit' 的代码
    print("finishing Lua program")
end}
setmetatable(t, t)
_G["atexit"] = t
```

>---
#### 垃圾收集器

Lua 执行自动内存管理，通过运行一个垃圾收集器来收集所有的死对象来自动管理内存，其中 `string`、`table`、`userdata`、`function`、`thread`、`internal struct` 等都是自动管理的对象。

一直到 Lua5.0 使用的都是标记清除式垃圾收集器。这种垃圾收集器的特点是会时不时地停止主程序的运行来执行一次完整的垃圾收集周期：标记（mark）、清理（cleaning）、清除（sweep）、析构（finalization）。
  - 标记：把根节点集合（由 Lua 可以直接访问的对象组成）标记为活跃，这个集合只包括 C 注册表。保存在一个活跃对象中的对象是程序可达的，弱引用表中的元素不遵循这个规则。当所有可达对象被标记为活跃时，标记阶段完成。
  - 清理：Lua 主要处理析构器和弱引用表。首先 Lua 遍历所有被标记为需要进行析构但未被标记为活跃状态的对象重新标记为活跃（复苏），并被放在一个单独的列表中（析构阶段会用到）。然后 Lua 遍历弱引用表并从中移除键或值未被标记的元素。
  - 清除：遍历所有对象（Lua 会把所有创建的对象放在一个链表中），所有非活跃对象被回收，活跃对象被清理标记，进入下一个清除阶段。
  - 析构：Lua 调用清理阶段被分离出的对象的析构器。

Lua5.1 使用了增量式垃圾收集器，也会像标记清除式一样执行相同的步骤。不同的是，增量式不需要在垃圾收集期间停止主程序的运行。增量式与解释器一起交替运行，解释器可能会改变一个对象的可达性，为了保证收集的正确性，垃圾收集器中的有些操作具有发现危险改动和纠正涉及对象标记的内存屏障。

Lua5.2 引入了紧急垃圾收集。当内存分配失败时，Lua 会强制进行一次完整的垃圾收集，并再次尝试分配。这些紧急情况可以发生在 Lua 进行内存分配的任意时刻，包括 Lua 处于不一致的代码执行状态时。这类型的垃圾收集动作不能运行析构器。

在 Lua5.4 之后，GC 可以在两种模式下工作，增量式或分代式。
- 增量模式：每个 GC 周期以小步骤执行标记、扫描和收集，并与程序的执行交替运行，收集器可以通过参数周期 *pause*（通过 `setpause`）、步长倍率 *stepmul*（通过 `setstepmul`）、步长 *stepsize*（通过 `step`）进行控制。

- 分代模式：收集器经常进行次要收集，只遍历最近创建的对象，若小收集之后内存的使用仍高于限制，收集器将执行大收集（遍历所有对象）。代入模式使用两个参数，次要收集频率和主要收集频率
  - 对于次要收集频率 x，在前一个主要收集后，当内存增长到比正在使用的内存大 x%，将执行次要收集。默认值为 20，最大值为 200
  - 对于主要收集频率 y，在前一个主要收集后，当内存增长到比使用的内存大 y% 时，将执行主要收集。默认值为 100（超过上一次收集后使用量的两倍），最大值为 1000

> *控制垃圾收集的步长*

在函数 `collectgarbage(opt, ...)` 中，`opt` 提供了一个可选参数用来说明收集器进行何种操作：
  - `"stop"`：停止垃圾收集器，直到使用选项 `"restart"` 调用 `collectgarbage`。
  - `"restart"`：重启垃圾收集器。
  - `"collect"`：执行一次完整的垃圾收集，回收和所有不可达的对象（默认）。
  - `"step"`：执行某些垃圾收集工作，第二个参数 `data:integer` 指明工作量，默认值为 13。
  - `"count"`：以 KB 为单位返回当前已用的内存数，包括了尚未被回收的死对象。
  - `"setpause"`：设置收集器的 *pause* 参数（间歇率），```data``` 以百分比给出设定的新值，data = 100，参数设定为 1（100%）。
  - ```"setstepmul"```：设置收集器的 *stepmul* 参数（步进倍率），```data``` 为百分比单位，默认值为 100%。

任何垃圾收集器都是使用 CPU 时间换内存空间。不消耗 CPU 时间是以巨大的内存消耗为代价的，而程序能够使用尽可能少的内存，都是以巨大的 CPU 消耗为代价的。*pause* 和 *stepmul* 用于平衡这两个极端：
  - *pause* 用于控制垃圾收集器在一次收集完成后等待多久再开始新的一次收集（0 表示在上一次垃圾回收结束后立即开始新的收集，200% 使得等待时间翻倍）。该值等于小于 100 意味着收集器不会等待便开始一个新的循环，200 表示等待正在使用的内存翻倍时才开始一个新的周期；最大值为 1000。
  - *stepmul* 控制对于每分配 1KB 内存，垃圾收集器应该进行多少工作，即收集器相对于内存分配的速度。小于 100 时会使收集器过于缓慢，并可能导致收集器永远无法完成一个周期。默认值为 100，最大值为 1000。
  - *step* 控制每个增量步长的大小，较大的值会使收集器的行为像标记清除式垃圾收集器。默认为 13，即大约 8kb 的步长，解释器每分配大约 8kb 内存就进行一次收集。

>---
#### 垃圾回收算法原理

Lua 采用了标记清除式（Mark and Sweep）GC 算法。
- 标记：每次执行 GC 时，先以若干根节点开始，逐个把直接或间接和它们相关的节点都做上标记；
- 清除：当标记完成后，遍历整个对象链表，把被标记为需要删除的节点一一删除即可。

> *颜色标记*

Lua 用白、灰、黑三色来标记一个对象的可回收状态，其中白色又分为 **白1**，**白2**。

- **白色**：可回收状态（未标记待访问，标记待回收）：
  - 如果该对象未被 GC 标记过则此时白色代表当前对象为待访问状态。
  - 新创建的对象的初始状态就应该被设定为白色，因为该对象还没有被 GC 标记到，所以保持初始状态颜色不变，仍然为白色。
  - 如果该对象在 GC 标记阶段结束后，仍然为白色则此时白色代表当前对象为可回收状态。但其实本质上白色的设定就是为了标识可回收。

+ **灰色**：中间状态（已访问待标记）：当前对象为待标记状态。当前对象已经被 GC 访问过，但是该对象引用的其他对象还没有被标记。

- **黑色**：不可回收状态（标记不可回收）：当前对象为已标记状态。当前对象已经被 GC 访问过，并且对象引用的其他对象也被标记了。（表示有引用关联）。

白色分为 **白1**，**白2**：

- 在 GC 标记阶段结束而清除阶段尚未开始时，如果新建一个对象，由于其未被发现引用关系，原则上应该被标记为白色，于是之后的清除阶段就会按照白色被清除的规则将新建的对象清除。这是不合理的。
- 于是 Lua 用两种白色进行标识，如果发生上述情况，Lua 依然会将新建对象标识为白色，不过是 “当前白”（比如 *白1*）。而 Lua 在清扫阶段只会清扫 “旧白”（比如 *白2*）。
- 在清扫结束之后，则会更新 “当前白”，即将 *白2* 作为当前白。下一轮 GC 将会清扫作为 “旧白” 的 *白1* 标识对象。

```
  new  ——> 白  ——> Free 
         ↙  ↖
       灰 ———> 黑 
```

> *垃圾回收详细过程*

```
新建可回收对象，并将其标记为白色 
              ↓          ↑ 
 yes ——  满足 GC 条件 —— no
  ↓
从根对象开始标记，将白色置为灰色，并加入到灰色链表中
                      ↓
  yes ———————— 灰色链表是否为空 —— no
   |                ↑             ↓
   |     从灰色链表中取出一个对象将其标记为黑色，并遍历和这个对象相关联的其他对象
   ↓
最后对灰色链表进行一次清除且保证是原子操作
                   ↓
根据不同类型的对象进行分步回收。回收中遍历不同类型对象的存储链表。
                   ↓
该对象存储链表是否达到链尾 —— no
      ↑                     ↓
将对象置为白 ← no — 逐个判断对象是否为白 — yes → 释放对象所占用的空间
```

---
### 附录
#### 格式化输出

`string.format()` 生成格式化字符串。

```lua
-- 格式转换符
%c     -- 接受整数, ASCII 字符
%a+    -- 接受数值, p 计数法
%d,%i  -- 接受整数, 有符号整数
%o     -- 接受整数, 八进制数
%u     -- 接受整数, 无符号整数
%x,%X  -- 接受整数, 十六进制数
%e,%E  -- 接受数值, 科学记数法
%f     -- 接受数值, 浮点数格式
%g,%G  -- 接受数值, %e/%E 或 %f 中较短的一种格式
%q     -- 接受字符串, 转化为 Lua 格式
%s     -- 接受字符串,
-- 修饰符
+      -- 显示数值符号位
0      -- 宽度前导零在后面指定了字串宽度时占位用. 不填时的默认占位符是空格|
-      -- 左对齐, 默认右对齐
n      -- 整数, 指定宽度
.n     -- 小数位数, 或字串裁切
```

>---

#### 日期与时间

Lua 使用两种方式表示日期和时间：*integer*（从 `Jan 01,1970,0:00 UTC` 开始后的秒数）和 *data* 表（`date = {year,month,day,hour,min,sec,wday,yday,isdst}`）。`os.time()` 以整数形式返回当前日期和时间。

```lua
print(os.time())	--> 1669608000 (11/28/2022, 12:00 UTC)
print(os.time({year = 2022,month = 11,day = 28}))
```

`os.date(format, time?)` 在一定程度上是函数 `os.time` 的反函数，返回一个包含日期及时刻的字符串或表。格式化方法取决于所给字符串 `format`：

```lua
*t             -- 返回 table 格式
%a,%A          -- 周几简写(Wed)，全写(Wednesday)
%b,%B          -- 月份简写(Sep)，全写(September)
%c             -- 日期和时间，(09/12/22 23:01:01)
%d             -- 一个月的第几天，(16) , [01,31]
%H             -- 24小时制的小时数，(23), [00,23]
%I             -- 12小时制的小时数，(11), [01,12]
%j             -- 一年中的第几天，(259), [001,365]
%m             -- 月份
%M             -- 分钟
%p             -- "am" 或 "pm"
%S             -- 秒数
%w             -- 星期(0)	Sunday-Saturday(0-6)
%W             -- 一年的第几周 [00,53]
%x             -- 日期，(09/16/99)
%X             -- 时间，(23:48:10)
%y             -- 表示年份的两位数，(22)，[00,99]
%Y             -- 完整的年份，(2022)
%z             -- 时区，(-0300)
%%             -- 百分号
```

`!%<signal>` 叹号开头，`os.date` 会以 UTC 格式进行解析。

```lua
local now = os.time()
print(os.date("%c", now)) 
print(os.date("!%c", now)) 
-- Tue Jun 10 02:46:23 2025
-- Mon Jun  9 18:46:23 2025
```

>---
#### 模式匹配

`string.find` 在目标字符串中搜索指定的模式并返回第一个匹配到的开始与结束位置的索引（字节）。`string.match` 则直接返回符合模式匹配的子串。`string.gsub` 将匹配到的子串替换为给定子串。

```lua
print(string.find("hello world", "hello"))                  --> 1 5
print(string.match("Today is 1/1/1970", "%d+/%d+/%d+"))     --> 1/1/1970
print(string.gsub("Lua is cute", "cute", "great"))          --> Lua is great
```

`string.gmatch` 返回一个迭代器，连续调用时依次返回符合模式的子串。

```lua
s = "some string"
words = {}
for w in string.gmatch(s, "%a+") do
    words[#words + 1] = w    -- words = {"some", "string"}
end
```

> **模式**

```lua
.          -- 任意字符
%a         -- 字母
%c         -- 控制字符
%d         -- 数字
%g         -- 除空格外的可打印字符
%l         -- 小写字母
%p         -- 标点符号
%s         -- 空白字符
%u         -- 大写字母
%w         -- 字母和数字
%x         -- 十六进制数字
-- 以上字符的大写形式表示该字符分类的补集
%M         -- M 表示对魔法字符的转义，可以是 ().%+-*?[]^$，例如 %% 表示 %
-- 魔法字符
[]         -- char-set 字符集	[0123456] 表示 0-6
-          -- 连接字符集范围 [0-6]
^          -- 字符集的补集 [^\n] 表示换行符外的字符，%S 等价于 [^%s]
+          -- 重复至少一次，%d+ 表示匹配一个或多个数字
*          -- 重复最少零次
-          -- 重复最少零次（最小匹配）
?          -- 出现零次或一次
^,$        -- 锚定目标字符串开头(^)或结尾($)，^ 表示从字符串开头开始查找
%b xy      -- 匹配成对的字符串，x 为起始，y 为结束。例如 %b() 返回包含 () 内中间的字符串
```

`%f[set]`前向匹配，该模式只有在后一个字符位于 *set* 内而前一个字符不在范围内时匹配一个空字符串。前向匹配把目标字符串中第一个字符前和最后一个字符后的位置当成空字符串。

```lua
s = "the anthem is the theme"
print(string.gsub(s, "%f[%w]the%f[%W]", "one"))
    --> one anthem is one theme
```

> **捕获**

捕获（capture）机制允许根据一个模式从目标字符中抽出与该模式匹配的内容用于后续用途，可以通过把模式中需要捕获的部分放到 `()` 中来指定捕获。函数 `string.match` 会将所有捕获到的值最为单独的结果返回。

```lua
pair = "name = Anna"
key, value = string.match(pair, "(%a+)%s*=%s*(%a+)")
print(key, value)	--> name  Anna
```

空白捕捉 `()` 用于锚定捕获模式在目标字符串中的位置。

```lua
print(string.match("hello", "()ll()"))	-->  3  5 (和 string.find 有所区别)
```

`% num` 表示匹配第 *num* 个捕获的副本，`%0` 表示整个匹配。

```lua
s = [[then he said: "it's all right"!]]
q, quotePart = string.match(s, "([\"'])(.-)%1")
q           --> "
quotePart   --> it's all right
```

读取文本，并记录每个单词出现的次数。

```lua
function F(file, n)    -- 输出有序列表中的前 n 个元素
    assert(tonumber(n) > 0)
    local f = io.input(file)
    local counter = {}
    for line in io.lines(file,"l") do
        for word in string.gmatch(line, "%w+") do
            counter[word] = (counter[word] or 0) + 1
        end
    end
    io.close(f)
    local words = {}
    for w in pairs(counter) do
        words[#words + 1] = w
    end
    table.sort(words, function(w1, w2)  -- 按出现次数降序
        return counter[w1] > counter[w2] or counter[w1] == counter[w2] and w1 < w2
    end)
    for i = 1, n > #words and #words or n do
        io.output(io.stdout)
        io.write(words[i], "\t", counter[words[i]], "\n")
    end
    io.close(io.stdout)
end


F("filename", 10)
```

>---
#### 打包与解包

`string.pack` 和 `string.unpack` 用于在二进制和基本类型值（数值，字符串）之间进行转换。第一个参数为格式化字符串：

```lua
<           -- 小端编码
>           -- 大端编码
=           -- 本地环境的原生大小端编码
![n]        -- 将最大对齐设为 n，默认为本地
b           -- 有符号字节 char
B           -- 无符号字节 unsigned char
h           -- short
H           -- unsigned short
l           -- long
L           -- unsigned long
j           -- Lua 有符号 lua_Integer
J           -- Lua 无符号 lua_Unsigned
n           -- Lua_Number
T           -- size_t
i[n]        -- n 字节的有符号整数，默认为本地
I[n]        -- n 字节的无符号整数，默认为本地
f           -- float
d           -- double
cn          -- 固定长度 n 的字符串
z           -- 以 \0 终结的字符串
s[n]        -- 由表示长度的 n 字节大小整数（默认是 size_t）打头的字符串
x           -- 1 字节填充 (0)
Xop         -- 根据转换项 op 对齐的空对象
空格        -- 忽略格式项，分隔各个选项
-- [n] 一个可选整数, [1,16]
```

`string.packsize` 返回 `pack` 结果的长度，该结果不包含 `s` 和 `z` 选项。

```lua

```

---