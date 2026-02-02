## Lua 规范摘要

### 1. 基本概念

#### 1.1. 类型系统

Lua 是一种动态类型语言，有 8 种内置基本类型：
- *nil*：空值，*table*，*function*，*thread*，*full userdata* 默认为 *nil*
- *number*： 64 位（默认）*integer* 和 *float*
- *string*：不可变 UTF-8 字符序列；
- *function*：Lua 或 C 函数；
- *userdata*：分为 *full userdata* (由 Lua 管理的占据内存的 C 对象) 和 *light userdata* (C 指针值)
- *thread*：表示独立的执行线程，用于实现协程。
- *table*：可用于表示数组、列表、集合、图、记录、字典、树等

*table*，*function*，*thread*，*full userdata* 值为对象，变量保存对它们的引用。`type(v)` 返回给定值的类型字符串。

>---
#### 1.2. 模块与环境

> `_ENV` 和 `_G`

非局部变量 `var` 的调用都转换为对 `_Env.var` 的调用。加载一个块时，`_Env` 首先初始化为全局环境 `_G` 并加载 Lua 标准库。全局变量通过 `_G` 表访问，`rawget` 执行表的原始访问：

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

`_ENV` 可以重新指向一个新环境，并丢失之前的状态。

```lua  
a = 1     
local newgt = {}    -- 创建新环境继承旧环境
setmetatable(newgt, { __index = _G })
_ENV = newgt

print(a, _G.a)  -- 1    1
a = 10
_G.a = 20
print(_ENV.a, _G.a)  -- 10   20 
```

> 模块

一个模块可以由 Lua 或 C 编写，通过 `require` 加载，任何的导出变量加载到本地 `_ENV` 中，仅加载一次并保存到 `package.loaded`。

```lua
mod = require(moduleName)
---------------------------------
local m = require "mod"			--> 引入 mod 模块
local f = require "mod".func    --> 引入 mod 模块中的 func 函数
local sub = require "mod.sub"	--> 引入 mod.sub 子模块
```
 
加载另一个模块时，模块的导出变量自动进入当前环境。可以利用 `_ENV` 的定界特性，将模块进行分离：

```lua
-- M1.lua
local M = {}
_ENV = M
function func()
    <body>
end
return M;

-- M2.lua
local M1 = require "M1"
M1.func()   
```


>---
#### 1.3. 编译与执行

Lua 可以在运行代码前执行预编译。
- `dofile(filename)` 用于直接执行 Lua 代码段。
- `loadfile` 与 `load` 执行预编译并返回一个函数；`assert(loadfile(...))` 检查编译错误。

```lua
i = 32
local i = 0
f = assert(load("i = i + 1; print(i)"))
g = function() i = i + 1; print(i) end
f()    -- 33
g()    -- 1
```

*luac* 程序生成预编译文件；`load` 和 `loadfile` 也可以接受预编译代码。

```shell
$ luac -o name.lc source.lua
```
```lua
assert(loadfile("name.lc"))()
```

>---
#### 1.4. 词法元素

> **标准关键字**

| description        | keywords                                             |
| :----------------- | :--------------------------------------------------- |
| 逻辑运算           | `and`,`or`,`not`                                     |
| 全局或局部变量声明 | `global`,`local`                                     |
| 局部变量常量限定   | `<const>`                                            |
| 局部变量关闭限定   | `<close>`                                            |
| 表函数自引用       | `self`                                               |
| 空值               | `nil`                                                |
| 布尔值             | `true`,`false`                                       |
| 函数声明           | `function-end`                                       |
| 跳转语句           | `break`,`goto`,`return`                              |
| 迭代语句           | `while-do`,`repeat-until`,`for-do`,`for-in-do`,`end` |
| 条件语句           | `if-then`,`else`,`elseif-then`,`end`                 |

> **运算符**

```lua
+    -   *   /   %   ^   #
&    ~   |   <<  >>  //  
==   ~=  <=  >=  <   >   =
(    )   {   }   [   ]   ::
;    :   ,   .   ..  ...
```

---
### 2. 类型和声明

Lua 中有 8 个基本类型：*nil*（空）、*boolean*（布尔）、*number*（数值）、*string*（字符串）、*userdata*（用户数据）、*function*（函数）、*thread*（线程） 和 *table*（表）。`local` 声明局部变量。`_` 表示弃元。

未初始化变量值为 `nil`，`nil` 可用于将变量或表键值置空，Lua 自动垃圾回收。真值相同的整数和浮点数被视为相等，`math.type(n)` 返回数值底层类型；浮点数支持 E 和 P 计数法。    


```lua
type(nil)       --> nil
type(true)      --> boolean
type(1234)      --> number
type("Hello")   --> string
type(io.stdin)  --> userdata
type(print)     --> function
type({})        --> table
type(type(X))   --> string
math.type(3)    --> integer
math.type(3.14) --> float
type(coroutine.create(type))   --> thread
type(coroutine.wrap(type))     --> function
```

> *浮点型与整型之间的转换*

数值相等的整数和浮点数类型之间可以互相转换：整数 与 0.0 加法转换为浮点数；浮点数（小数部分为 0）与 0 按位或转换为整数。

```lua
12345 + 0.0    --> 12345.0
12345.0 | 0    --> 12345
2^53           --> 9.007199254741e+15	(浮点型值)
2^53 | 0       --> 9007199254740922	(整型值)
-- error
3.2 | 0        --> 小数部分 > 0
2^64 | 0       --> 超出范围
```

>---
#### 2.1. string

Lua 字符串包括转义字符串 `"string"` 或 `'string'` 和原始字符串 `[=[string]=]`。`#str` 返回字符串长度，`x .. y` 拼接字符串，操作数可以是字符串或数值。

```lua
a = "a 'line\n'"
b = 'another "line\n"' .. "\255" .. "\xff" .. "\u{7fffffff}"
c = 1 .. "234"
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
\u{h…h}      --> 0 ~ 7fffffff
\[ \]        --> [ ]
```

`\z`：忽略任意空白字符直到第一个非空白字符。

```lua
"12345\z    6"  --> 123456
"12345\z
678"            --> 12345678
"123\z \r678"   --> 123\r678
```

>---
#### 2.2. table

表（table）可以表示数组、列表、符号表、集合、记录、图形、树等数据结构。`_G` 表用来存储全局变量和全局环境。键可以是除 `nil` 和 `NaN` 外的任何值，值为 `nil` 的任何键不被视为表的一部分。

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

列表的声明，索引从 1 开始。索引可以显式指定，计算列表边界时会忽略小于 1 的索引。列表出现连续 `nil` 值时，`#table` 返回连续空洞前的索引。 

```lua
-- 索引从 1 开始为每个列表元素建立索引关联
day = {"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
-- 显式索引
arr = {[1]=1,[2]=2,[4]=4,[5]=5,[7]=7,[9]=9}  -- #arr = 5
Arr = {[1]=1,[2]=2,[4]=4,[5]=5,[7]=7,[8]=8}  -- #Arr = 8
a1 = {1,2,3,4,5,6}           --> #a = 6
a2 = {1,2,3,nil,nil}         --> #a = 3
a3 = {1,nil,nil,nil,nil,2}   --> #a = 6
```

>---
#### 2.3. function

Lua 函数是第一类值。支持多值返回和闭包。

```lua
-- 函数声明
function Func( Params )
    body 
    [return r1,r2, ...]  -- 多值返回
end
Func = function( Params) body end

-- 闭包
function counter()
    local count = 0
    return function()
        count = count + 1
        return count
    end
end
local c = counter()
print(c())  -- 1
print(c())  -- 2
```

当函数只有一个参数且是字符串常量或表构造器时，`f()` 括号是可选的。

```lua
print "Hello"  -- print("Hello")
type {}        -- type({})
```

> *变长参数*

变长参数 `...` 可以由 *table* `{...}` 接收，或由多重赋值析构。

```lua
function Foo(a, ...)
    local t = { ... }
    -- v1,v2 = ...
    if (#t > 0) then
        for i = 1, #t do
            print(t[i])
        end
    end
    print("end :" .. a)
end
Foo(1, 2, 3, 4, 5, 6)
```

`select(n, ...)` 访问第 n 个变长参数；`select(#, ...)` 返回变长参数总数。

```lua
-- 打印奇数位的元素
function Foo(...)
    local len = select("#", ...)
    local i, a = 1, 0
    while i <= len do
        a = select(i, ...)
        print(a)
        i = i + 2
    end
end
Foo(1,2,3,4,5,6,7,8)	-- 1 3 5 7
```

> *表函数*

表函数通过两种方式调用：`t.fun()` 或 `t:fun()`。`t:fun` 将表自身作为第一个参数传入。`function T:Func(params)` 声明相当于 `function T.Func(self, params)`。

```lua
t = {1,2,3,4,5,6}
function t:Traverse()    -- t.Traverse(self) 
    for i =1,#self do
        print(self[i])
    end
end
t:Traverse()	-- or t.Traverse(t)
```

>---
#### 2.4. thread 与 coroutine

协程（*coroutine*）与线程（*thread*）类似：协程拥有自己的栈、局部变量和指令指针，与其他协程共享了全局变量和其他几乎一切资源。

`coroutine.create(f)` 返回一个 *thread* 协程，*thread* 协程有四种状态：*suspended*、*running*、*normal*、*dead*。创建协程时，协程不会自动运行而处于挂起状态。

```lua
co = coroutine.create(func)     
print(coroutine.status(co))     -- suspended
coroutine.resume(co [,params])  -- running
print(coroutine.status(co))     -- dead
```

`coroutine.wrap(f)` 返回一个 *function* 协程，无法获取此类协程的状态。*thread* 协程异常传播至 `coroutine.resume(co)` 返回中；*function* 协程异常直接导致程序错误。

```lua
-- 调用 thread 协程
local co = coroutine.create(f)
local ok [,rt] = coroutine.resume(co)
-- 调用 function 协程
local cf = coroutine.wrap(f)
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

>---
#### 2.5. 全局遍历和局部变量 

Lua 隐式启用 `global *` 全局模式，变量声明默认全局，`local` 声明局部变量。在某个范围内显式使用 `global` 则在当前范围内启用全局变量严格模式；`global *` 在当前范围内用于关闭全局变量严格模式。

```lua
global hi, print
hi = "World"
do
    local hi = "Hello "
    print(hi)  -- Hello
end
print(hi)  -- World

global *   -- 关闭全局严格模式
t = table.create(5, 10)  -- 隐式全局，预分配表空间
```

>---
#### 2.6. 局部属性

> *const*

`<const>` 属性赋予局部变量常量属性。

```lua
local t <const> ={ a = 1}
t = nil 	-- error
t.a = 2		-- ok
```

> *close*

局部变量（`<close>`）必须具有 `__close` 元方法或 `__close = false`。`<close>` 变量生存期结束时调用其关联元方法 `__close`。

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
### 3. 表达式

> **运算符与优先级**

| Category   | Operators                                                   |
| :--------- | :---------------------------------------------------------- |
| 幂运算     | `x ^ y`                                                     |
| 一元       | `-x`, `#string`, `#table`, `not x`, `~x`                    |
| 乘法       | `x * y`, `x / y`(浮点除法), `x // y`(向下取整除法), `x % y` |
| 加法       | `x + y`, `x - y`                                            |
| 字符串拼接 | `x .. y`                                                    |
| 移位       | `x << y`, `x >> y`                                          |
| 按位       | `x & y`, `x ~ y`(异或), `x \| y`                            |
| 关系       | `x < y`, `x > y`, `x <= y`, `x >= y`, `x == y`, `x ~= y`    |
| 逻辑       | `x and y`, `x or y`                                         |
| 赋值       | `x = y`                                                     |

> x % y

`x % y` 取模运算的定义 `x%y = x-((x//y)*y)`，表达式结果符号与第二操作数一致。可以利用取模运算保留浮点运算有效位。

```lua
-9 % 2      --> 1:  -9-(-5)*2 = 1
-9 % 2.0    --> 1.0
math.pi - math.pi % 0.0001	--> 3.1415
```

> and & or

`and` 与 `or` 支持短路原则。
- `and` 在第一个操作数为 `false` 或 `nil` 时返回第一个操作数，否则返回第二个操作数；
- `or` 在第一个操作数为 `false` 或 `nil` 时返回第二个操作数，否则返回第一个操作数。

```lua
10 or 20            --> 10
10 and 20           --> 20
false and nil       --> false
false or nil        --> nil
nil or "a"          --> "a"
```

可以利用 `and` 和 `or` 机制构造三目运算：

```lua
X and Y or Z
--> X == true --> Y
--> X == false --> Z
```

> << & >>

移位是逻辑移位，以 0 补齐空位

```lua
12(10) = 00001100(2)
00001100 >> 1 = 00000110  --> 12>>1 = 6
00001100 << 2 = 00110000  --> 12<<2 = 48

-18(10) = 11101110(2)     -- 补码表示
11101110 << 1 = 11011100  --> -18<<1 = -36
11101110 >> 2 = 01111011  --> -18>>2 = 123 (逻辑移位)
```

利用 floor 除法模拟（`num // (2^n)|0`）实现算术移位：当 n>0 表示算术右移；当 n<0 表示算术左移。

```lua
-- 负数的算术右移
    -10 >> 2 等价于 -10//2^2|0 --> -3
-- 算术左移
    -10 >> -2 == -10 << 2 == -10//(2^-2|0) --> -40
```

---
### 4. 语句

#### 4.1. 代码块：do-end

```lua
do
    <body>
end
```

>---
####  4.2. 条件控制：if

`false` 和 `nil` 值为假，`true` 和非 `nil` 值为真。

```lua
if <cond> then
    <body>
elseif <cond> then
    <body>
else
    <body>
end
```

>---
#### 4.3. 迭代语句：while, repeat, for

```lua
while <cond> do
    <body>
end
```
```lua
repeat
    <body>
until <cond>
```

> for

```lua
for i = Begin, End [, Step = 1] do
    <body>
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
        local var_1, ..., var_n = _f(_s, _var)
        _var = var_1
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
-- 相当于
do
    local _f, _s, _var = func, 5, 0
    while true do
        local i, v = _f(_s, _var)
        _var = i
        if _var == nil then break end
        print(i, v)
    end
end
```

> for-in pairs 

```lua
arr = {1,3,4,5,a="A",b="B"}
for k,v in pairs(arr) do   -- 遍历键值对 	
    print(k,v)
end
```

> for-in ipairs

```lua
arr = {1,3,4,5,a="A",b="B"}
for i,v in ipairs(arr) do	-- 遍历列表
    print(i,v)	-- 1,3,4,5
end
```

>---
#### 4.4. 跳转语句：goto, break, return

> `goto` 

`goto` 跳转到目标标签。

```lua
do 
    -- code
    goto OUT
    -- code
end
::OUT::
    -- code
```

> `break`

`break` 中断所属循环体的执行。

```lua
while cond do
    -- code
    break   -- 相当于 goto END
    -- code
end
::END::
```

> `return` 

`return` 终止函数执行并返回（多值）到调用方。文件范围 `return` 表示从模块返回值。 

```lua
local M =  _ENV
function Func()
    -- code
    return ok, rt   -- 多值返回
end
return M   -- 文件范围返回值
```

---
### 5. 元表与元方法

元表定义了其原始值允许的某些操作，可以设置元表中特定元方法来更改值的行为。*table* 和 *full userdata* 具有单独的元表，除字符串外其他类型默认没有元表。

```lua
local subTable, father = {}, {}
setmetatable(subTable, father)   -- 设置元表
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
__close     -- 待关闭变量 <close>
__mode      -- 弱表模式 "k","v","kv"
__metatable -- 元表
__pairs     -- 在 for-pairs 替选调用
```

通常将 `__le` 作为其他关系元方法（`__eq`,`__lt`）的基方法。

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
#### 5.1. __index, __newindex

读取 *table* 键值对时调用 `__index` 元方法，默认返回 `nil`。更新键值对时调用 `__newindex`。`rawget` 与 `rawset` 使用原始方式对 *table* 执行操作。

```lua
local Array = {} -- 固定长度数组
Array.New = function(...)
    local arr = { ... }
    local list = { Length = #arr, data = arr }
    setmetatable(list, Array)
    return list
end
function Array:__index(index)
    if index > self.Length or index < 1 then
        error("index is out of range.")
    end
    return self.data[index]
end

function Array:__newindex(index, value)
    if index > self.Length or index < 1 then
        error("index is out of range.")
    end
    self.data[index] = value
end

local a = Array.New(1, 2, 3, 4, 5, 6)
print(a.Length)

for i = 1, a.Length do
    print(a[i])
end
```

>---
#### 5.2. __call

`__call(args...)` 为表创建函数调用 `f()` 语法。

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
#### 5.3. __tostring, __name

函数 `print` 总是调用 `tostring` 格式化输出。`tostring` 首先会查找元方法 `__tostring`，否则查找 `__name` 作为替代。

```lua
local mt = {}
mt.__tostring = function(self)
    if #self == 0 then
        return tostring(self)
    end
    local s = '{' .. tostring(self[1])
    for i = 1, #self do
        if i ~= 1 then
            s = string.format(s .. ',' .. tostring(self[i]))
        end
    end
    return s .. '}'
end

local t = { 4, 10, 2 }
setmetatable(t, mt)
print(t) -- {4,10,2}
```

>---
#### 5.4. __gc

元表中拥有 `__gc` 字段时，子表在 `setmetatable` 时会被标记为可析构。可析构对象在垃圾回收阶段，垃圾收集器自动调用（仅一次）该对象的 `__gc` 元方法。

```lua
local mt = {}
mt.__gc = function ()
	print("call __gc ...")
end

local t = nil
t = setmetatable({}, mt)
t = nil  -- 置空
collectgarbage("collect")  -- 强制回收，自动调用 __gc
```

>---
#### 5.5. __close

可关闭对象离开其作用域时，将调用该对象的 `__close` 元方法。`__close` 可用于释放某些资源。

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

>---
#### 5.6. __mode

`__mode` 设置弱引用表的类型，`"k"` 键弱引用，`"v"` 值弱引用，`"kv"` 键值弱引用。对于一个弱引用表，只要一个弱引用的键或值被回收，就把这个键值对从表中删除。

```lua
local mt = { __mode = "v" } -- 值弱引用
local a = {}
local b = { key = a }
setmetatable(b, mt)
print(b.key)      -- table: 000001F0DD803580
a = nil
collectgarbage("collect")  -- 强制垃圾回收，a 被回收了
print(b.key)      -- nil
```

>---
#### 5.7. __metatable

`__metatable` 用于保护元表。`getmetatable` 调用子表时返回元表的 `__metatable`，而 `setmetatable` 更改此类子表时会引发错误。

```lua
local mt = { __metatable =  "metatable is protected" }
local t = setmetatable({}, mt)  -- t 受保护
print(getmetatable(t))          -- "metatable is protected"

local ok, err = pcall(setmetatable, t, {})
if not ok then
    print(err)  -- cannot change a protected metatable
end
```

>---
#### 5.8. __pairs

函数 `pairs` 对应元方法 `__pairs`，指定表在 `for-in-pairs` 的行为。

```lua
local mt = {}
local default = {
    name = "Lua",
    telephone = 123456,
    id = 7,
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
for k, v in pairs(test) do  -- test __pairs 
    print(v)
end
```

---
### 6. 面向对象编程

使用参数 `self` 是 Lua 面向对象语言的核心点。可以利用原型的概念实现面向对象编程，利用元表 `__index` （`setmetatable(A,{__index = B})`），让 B 成为 A 的一个原型。

```lua
local prototype = {}
prototype.__index = prototype
prototype.Func = function ()
    print("Call FuncA")
end 

function NewPrototype()
    return setmetatable({}, prototype)
end

local p = NewPrototype()
p.Func()
```

>---
####  6.1. 继承

利用 `__index` 和 `self` 机制实现继承。

```lua
local prototype = {}
function prototype:new(o)
    o = o or {}
    self.__index = self
    setmetatable(o, self)
    o.base = self   -- 关联超类
    return o
end

local o1 = prototype:new()
o1.hello = "world"
local o2 = o1:new()	-- 单一继承
print(o2.hello)	    -- world
```

多重继承意味着一个类可以具有多个超类，因此需要一个独立的方法（createClass）从一个类中创建子类，其参数为新类的所有超类。

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
### 7. 错误处理

错误将中断程序正常流程。
- `error(msg, level)` 抛出错误；
- `assert(exp, msg)` 执行断言；
- `warn(msg)` 生成警告；
- 安全调用函数 `pcall(func, args...)` 和 `xpcall(func, handle, args...)` 用于捕获错误，常用 `debug.debug` 和 `debug.traceback` 作为 `xpcall` 的消息处理函数。

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
ok, stat = xpcall(panic, handle, 1)  -- false, nil
```

---
### 8. 垃圾回收

Lua 使用自动内存管理。弱引用表（weak table）、析构器（finalizer）、函数 `collectgarbage` 是 Lua 中用来辅助垃圾收集器的主要机制。
  - 弱引用表允许收集 Lua 中还可以被程序访问的对象；
  - 析构器允许收集不在垃圾收集器直接控制下的外部对象；
  - 函数 `collectgarbage` 允许控制垃圾收集器的步长。

>---
#### 8.1. 弱引用表

弱引用是一种不在垃圾收集器考虑范围内的对象引用。对于一个弱引用表，只要一个弱引用键或值被回收，就将该键值对删除；弱引用表是由 `__mode` 设置，`"k"` 键弱引用，`"v"` 值弱引用，`"kv"` 键值弱引用。

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

一个具有弱引用键和强引用值的表是一个瞬表，瞬表键的可访问性控制着对应值的可访问性。对于瞬表的一对 `(k,v)`，指向 `v` 的引用只有当存在某些指向 `k` 的其他外部引用存在时才是强引用，否则即使 `v` 直接或间接地引用了 `k`，垃圾收集器都会将 `k` 从表中移除。

```lua
-- 瞬表
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

local k = "Hello"
local a = factory(k)
print(mem[k])  -- table[...]	
k = nil        -- 无论 v 是否引用 k 都会被移除
print(mem[k])  -- nil
```

>---
#### 8.2. 析构器

析构器 `__gc` 在对象被回收时调用。在每个回收周期，垃圾收集器在调用析构器之前清理弱引用表中的值，在调用析构器之后再清理键。

```lua
o = {x = "hi"}
setmetatable(o, {__gc = function(o) print(o.x) end})
o = nil
collectgarbage()    --> hi
```

当一个析构器被调用时，它的参数是正在被析构的对象。这个对象在析构期间重新变成活跃（临时复苏）。如果该对象在析构器返回后仍然可访问（析构时保存到全局变量时），这个对象就变成永久复苏。析构器仅调用一次。

由于复苏的存在，Lua 会在两个阶段中回收具有析构器的对象。当垃圾收集器首次发现某个析构器的对象不可达时，这个对象临时复苏并将其放入析构等待队列中，在析构器开始执行时标记为已被析构，当下一次垃圾收集器又发现这个对象不可达时，它就将这个对象删除。

程序结束时，Lua 会显式调用所有未被释放对象的析构器。可以利用这种特性实现某种形式的 `atexit()`。可以将这个特殊的表锚定在全局表中。

```lua
local t = {__gc = function()
    -- 'atexit' 的代码
    print("finishing Lua program")
end}
setmetatable(t, t)
_G["atexit"] = t
```

>---
#### 8.3. 垃圾收集器

Lua 执行自动内存管理，`string`、`table`、`userdata`、`function`、`thread`、`internal struct` 等都是自动管理的对象。

一直到 Lua5.0 使用的是 *标记清除式垃圾收集器*。这种垃圾收集器的特点是会时不时地停止主程序的运行来执行一次完整的垃圾收集周期：标记（mark）>> 清理（cleaning）>> 清除（sweep）>> -析构（finalization）。
  - 标记：把根节点集合（由 Lua 可以直接访问的对象组成）标记为活跃，这个集合只包括 C 注册表。保存在一个活跃对象中的对象是程序可达的，弱引用表中的元素不遵循这个规则。当所有可达对象被标记为活跃时，标记阶段完成。
  - 清理：Lua 主要处理析构器和弱引用表。首先 Lua 遍历所有被标记为需要进行析构但未被标记为活跃状态的对象重新标记为活跃（复苏），并被放在一个单独的列表中（析构阶段会用到）。然后 Lua 遍历弱引用表并从中移除键或值未被标记的元素。
  - 清除：遍历所有对象（Lua 会把所有创建的对象放在一个链表中），所有非活跃对象被回收，活跃对象被清理标记，进入下一个清除阶段。
  - 析构：Lua 调用清理阶段被分离出的对象的析构器。

Lua5.1 使用了 *增量式垃圾收集器*，也会像标记清除式一样执行相同的步骤。增量式不需要在垃圾收集期间停止主程序的运行。增量式与解释器一起交替运行，解释器可能会改变一个对象的可达性，为了保证收集的正确性，垃圾收集器中的有些操作具有发现危险改动和纠正涉及对象标记的内存屏障。

Lua5.2 引入了 *紧急垃圾收集*。当内存分配失败时，Lua 会强制进行一次完整的垃圾收集，并再次尝试分配。这些紧急情况可以发生在 Lua 进行内存分配的任意时刻，包括 Lua 处于不一致的代码执行状态时。这类垃圾收集动作不能运行析构器。

在 Lua5.4 之后，GC 可以在两种模式下工作，增量式或分代式。
- 增量模式：每个 GC 周期以小步骤执行标记、扫描和收集，并与程序的执行交替运行，收集器可以通过参数周期 *pause*（通过 `setpause`）、步长倍率 *stepmul*（通过 `setstepmul`）、步长 *stepsize*（通过 `step`）进行控制。

+ 分代模式：收集器经常进行次要收集，只遍历最近创建的对象，若小收集之后内存的使用仍高于限制，收集器将执行大收集（遍历所有对象）。分代模式使用两个参数，*次要收集频率* 和 *主要收集频率*。
  - 对于次要收集频率 x，在前一个主要收集后，当内存增长到比正在使用的内存大 x%，将执行次要收集。默认值为 20，最大值为 200
  - 对于主要收集频率 y，在前一个主要收集后，当内存增长到比使用的内存大 y% 时，将执行主要收集。默认值为 100（超过上一次收集后使用量的两倍），最大值为 1000

> 控制垃圾收集的步长

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
#### 8.4. 垃圾回收算法原理

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
### 9. 反射
#### 9.1. 自省机制

反射是程序用来检查和修改其自身某些部分的能力。Lua 支持的反射机制有：
- 环境允许运行时观察全局变量；
- 运行时检查和遍历未知数据结构（例如 `type` 和 `pairs`）；
- 允许程序在自身中追加代码或更新代码（`load` 和 `require`）。

但是 Lua 不能检查局部变量，开发人员不能跟踪代码的执行，函数也不知道调用方，`debug.library` 调试库填补了这些空缺。调试库由两类函数组成：自省函数（introspective function）和钩子（hook）
  - 自省函数允许检查一个正在运行的活动函数的栈、当前正在执行的代码行、局部变量的名称和值等；
  - 钩子函数允许跟踪一个程序的执行。

>---
#### 9.2. 自省函数

> *访问局部变量*

`debug.getlocal(fun|integer, index)` 检查任意活跃函数的局部变量，查询指定栈层次或函数的指定索引局部变量，返回变量名和值或 `nil`。Lua 按局部变量在函数中出现顺序对它们进行编号，编号只限于在函数当前作用域中活跃的变量（索引从 1 开始）。

`debug.setlocal(f|n,index,value)` 用于改变局部变量的值，函数返回被修改值的变量名。

```lua
local t = {}
print(debug.getlocal(1,1))
-- t	table: 000001E7622C4610

function foo(a,b)
    local i = 1
    print(debug.getlocal(1,3))   -- i	1
end
print(debug.getlocal(foo,2))     -- b
print(debug.getlocal(foo,3))     -- nil
```

`debug.getlocal(thread, fun, index)` 返回指定协程的栈层次局部变量信息。调试库中所有的自省函数都能够接受一个可选的协程作为第一个参数

```lua
function foo(a,b)
    local i = 1
    print(debug.getlocal(1,3))   -- i	1
end
local co = coroutine.create(foo)
print(debug.getlocal(co,foo,1))  -- a
print(debug.getlocal(co,foo,3))  -- nil
```

值为负的索引用于访问变长参数，-1 指向第一个额外参数，变量名称始终为 `(vararg)`

```lua
function foo(a, ...)
    local i = 1
    print(debug.getlocal(1, -2))
end

foo(1, 3, 5)    -- (vararg)		5
```

> *调用层次信息*

`debug.getinfo(function|num [, what?])` 返回与函数或栈层次的有关的一些数据的表。对于注册到 Lua 的 C 函数，只有字段 `what`、`name`、`namewhat`、`nups`、`func` 是有意义的。

<table>
    <tr>
        <th>表字段</th>
        <th>Lua 函数</th>
        <th>Lua 代码段</th>
        <th>C 函数</th>
    </tr>
    <tr>
        <td>source</td>
        <td>函数定义的位置</td>
        <td>load 返回定义字符串</td>
        <td>=[c]</td>
    </tr>
    <tr>
        <td>short_src</td>
        <td>source 的精简版本</td>
        <td>[string "code"]</td>
        <td>[c]</td>
    </tr>
    <tr>
        <td>linedefined</td>
        <td>函数定义的首行位置</td>
        <td>0</td>
        <td>-1</td>
    </tr>
    <tr>
        <td>lastlinedefined</td>
        <td>函数定义的末行位置</td>
        <td>0</td>
        <td>-1</td>
    </tr>
    <tr>
        <td>what</td>
        <td>函数类型 Lua</td>
        <td>main</td>
        <td>C</td>
    </tr>
    <tr>
        <td>name</td>
        <td colspan =3>返回函数名称的字段</td>
    </tr>
    <tr>
        <td>namewhat</td>
        <td colspan = 3>字段含义，可能是 global、local、method、field、""（空字符串）</td>
    </tr>
    <tr>
        <td>nups</td>
        <td colspan = 3>该函数上值的个数</td>
    </tr>
    <tr>
        <td>nparams</td>
        <td colspan = 3>函数参数的个数</td>
    </tr>
        <tr>
        <td>isvararg</td>
        <td colspan = 3>参数列表是否包含可变长参数</td>
    </tr>
    <tr>
        <td>activelines</td>
        <td colspan = 3>该函数所有活跃行的集合</td>
    </tr>
    <tr>
        <td>func</td>
        <td colspan = 3>该函数本身</td>
    </tr>
</table>

使用 `num` 作为参数调用 `getinfo` 返回有关相应栈层次上活跃函数的数据，0 表示 `getinfo` 自己。`num` 大于栈中活跃函数的数量时返回 `nil`。与栈层次相关的两个字段：`currentline` 表示当前该函数正在执行的代码所在的行；`istailcall` 表示函数是否是尾调用。`name` 只有在以一个数字为参数调用 ```getinfo``` 时才会起作用：

```lua
print(debug.getinfo(0).name)	-- getinfo
```

> *访问上值*

`debug.getupvalue(f,index)` 用于访问一个被 Lua 函数所使用的上值，Lua 按照函数引用上值的顺序对它们编号；`debug.setupvalue(f,index,value)` 用于更新上值，该函数返回被修改上值的变量名。

```lua
local up1 = 1
function foo()
    local a = up1
end
print(debug.getupvalue(foo, 1))     -- up1	1
print(debug.setupvalue(foo, 1, 2))  -- up1
print(up1)                          -- 2
```

>---
#### 9.3. 钩子函数

调试库中的钩子机制允许用户注册一个钩子函数，并且在 Lua 程序运行中某个特定事件发生时被调用：
  - 调用一个函数时产生的 *call* 事件
  - 函数返回时产生的 *return* 事件
  - 开始执行一行新代码时产生的 *line* 事件
  - 执行完指定数量的指令后产生的 *count* 事件

函数 `debug.sethook([thread,] hookf, mask [, count])` 中 `mask` 描述要监视事件的掩码，`count`（可选）描述以何种频度获取 `count` 事件。关闭钩子，只需不带任何参数的调用函数 `sethook`。

```lua
debug.sethook(print,"l")
-- Lua 发生 line 事件时会调用它，并输出解释器执行的每一行代码
```

>---
#### 9.4. 调优（profiler）

反射的一个常见用法是用于调优，即程序使用资源的行为分析。对于时间相关的调优最好使用 C 接口。开发一个性能调优工具来列出程序执行的每个函数的调用次数：

```lua
-- profiler.lua
local Counters = {}
local Names    = {}

-- 可以在函数活动时获取其名称
local function hook()
    local f = debug.getinfo(2, "f").func
    local count = Counters[f]
    if (count) == nil then
        Counters[f] = 1
        Names[f] = debug.getinfo(2, "Sn")
    else
        Counters[f] = count + 1
    end
end
local function getName(f)
    local n = Names[f]
    if (n.what == "C") then
        return n.name
    end
    local lc = string.format("[%s]:%d", n.short_src, n.linedefined)
    if n.what ~= "main" and n.namewhat ~= "" then
        return string.format("%s (%s)", lc, n.name)
    else
        return lc
    end
end
function Main()
    print("This is Main function...")
end
function funcA() end

-- 设置 call 事件的钩子函数
debug.sethook(hook, "c")    

Main()                      -- 运行主程序
for i = 1, 10 do
    funcA();
end

debug.sethook()             -- 关闭钩子

for func, count in pairs(Counters) do
    print(getName(func), count)
end

--[[
    sethook 1
    print   1
    [profiler.lua]:28 (Main)    1
    [profiler.lua]:31 (funcA)   10
]]
```

>---
#### 9.5. 沙盒（sandbox）

> 一个使用钩子的简单沙盒

该程序把钩子设置为监听 count 事件，使得 Lua 每执行 100 条指令就调用一次钩子函数；钩子函数只是一个递增计数器，并检查其是否超过了设定的阙值：

```lua
local steplimit = 1000     -- 最大能执行的 steps
local count = 0

local function step()
    count = count + 1
    if count > steplimit then
        error("script uses too much CPU")
    end
end
debug.sethook(step, "clr", 100) -- 设置钩子
```

> 控制内存使用

```lua
local function checkmem()
    if collectgarbage("count") > memlimit then
        error("script uses too much memory")
    end
end

local count = 0
local function step()
    checkmem()
    count = count + 1
    if count > steplimit then
        error("script uses to much CPU")
    end
end

debug.sethook(step, "clr", 100)

--[[
local s = "123456789012345"
for i = 1, 36 do
    s = s .. s
end
]]
```

> 使用钩子阻止对未授权函数的访问

```lua
local debug = require "debug"
local steplimit = 1000
local count = 0

-- 设置授权的函数
local validfunc = {
    [print] = true,
    [string.upper] = true,
    [string.lower] = true,
    [string.format] = false,
    --...    -- 其他授权函数
}

local function hook(event)
    print(tostring(event))
    if event == "call" then
        local info = debug.getinfo(2, "fn")
        if not validfunc[info.func] then
            error("calling bad function: " .. (info.name or "?"))
        end
    end
    count = count + 1
    if (count > steplimit) then
        error("script uses too much CPU")
    end
end

debug.sethook(hook, "clr", 100)
print(string.upper("hello"))
-- HELLO
print(string.format("%s %d", "foo", 1))
-- calling bad function: format
```

---
### 10. 附录


#### 10.1. 格式化输出

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

#### 10.2. 日期与时间

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
#### 10.3. 模式匹配

`string.find` 在目标字符串中搜索指定的模式并返回第一个匹配到的开始与结束位置的索引（字节）。
`string.match` 则直接返回符合模式匹配的子串。
`string.gsub` 将匹配到的子串替换为给定子串。

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

`%f[set]` 表示前向匹配，该模式只有在后一个字符位于 *set* 内而前一个字符不在范围内时匹配一个空字符串。前向匹配把目标字符串中第一个字符前和最后一个字符后的位置当成空字符串。

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
q, quotePart = string.match(s, "([\"'])(.-)%1")  -- 即 "([\"'])(.-)([\"'])"
q           --> "
quotePart   --> it's all right
```

案例：读取文本，并记录每个单词出现的次数。

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
#### 10.4. 打包与解包 

`string.pack` 和 `string.unpack` 用于在二进制和基本类型值（数值，字符串）之间进行转换。第一个参数为格式化字符串：

```lua
<           -- 小端编码
>           -- 大端编码
=           -- 本地编码，默认
![n]        -- 设置最大对齐为 n，默认为本地对齐;n = 0,1,2,4,8,...
b           -- signed char
B           -- unsigned char
h           -- signed short
H           -- unsigned short
l           -- signed long
L           -- unsigned long
j           -- lua_Integer
J           -- lua_Unsigned
n           -- Lua_Number
T           -- size_t
i[n]        -- n 字节的有符号整数，默认为本地大小
I[n]        -- n 字节的无符号整数，默认为本地大小
f           -- float
d           -- double
cn          -- 固定长度 n 的字符串
z           -- 以 \0 终结的字符串
s[n]        -- 前置长度 n 字节的字符串，默认为 size_t
x           -- 1 字节填充 (0)
Xop         -- 根据转换项 op 对齐的空项
' '        -- 忽略格式项，分隔各个选项
-- 对于 ![n],s[n],i[n],I[n],  n 可选 [1,16]
```
```lua
local fmt = "<I4i4ffs2"
--    "<"   : 小端序
--    "I4"  : 无符号4字节整数 (id)
--    "i4"  : 有符号4字节整数 (health)
--    "f"   : 单精度浮点数 (position.x)
--    "f"   : 单精度浮点数 (position.y)
--    "s2"  : 字符串，前2字节存储长度 (name)

-- 打包数据
local id = 1001
local health = 850
local pos_x = 1024.5
local pos_y = 2048.25
local name = "Astronaut"
local packed_data = string.pack(fmt, id, health, pos_x, pos_y, name)
print("data length:", #packed_data)  -- 27

-- 解包数据
local _id, _health, _x, _y, _name, next_pos = string.unpack(fmt, packed_data)
print("ID         :", _id)        -- 1001
print("Health     :", _health)    -- 850
print("Position X :", _x)         -- 1024.5 (浮点数可能有微小误差)
print("Position Y :", _y)         -- 2048.25
print("Name       :", _name)      -- Astronaut
print("next_pos   :", next_pos)   -- 28
```

---