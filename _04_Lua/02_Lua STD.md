## Lua STD

- version = 5.5
- [examples](./Lua%20Libs/LuaSTD_examples/README.md)

>---

### 1. 基础库 (base)

```lua
_G                            -- 全局环境表
_VERSION                      -- Lua 版本字符串
assert(v [, message])         -- 检查值是否为真，否则抛出错误
collectgarbage([opt [, arg]]) -- 控制垃圾收集器
dofile([filename])            -- 执行文件
error(message [, level])      -- 抛出错误
getmetatable(object)          -- 获取对象的元表
ipairs(t)                     -- 遍历表的数组部分
load(chunk [, chunkname [, mode [, env]]]) -- 加载代码块
loadfile([filename [, mode [, env]]]) -- 加载文件
next(table [, index])         -- 遍历表的下一个键值对
pairs(t)                      -- 遍历表的所有键值对
pcall(f [, arg1, ...])        -- 保护调用函数
print(...)                    -- 打印值
rawequal(v1, v2)              -- 原始相等性检查
rawget(table, index)          -- 原始获取表元素
rawlen(v)                     -- 原始获取长度
rawset(table, index, value)   -- 原始设置表元素
select(index, ...)            -- 返回可变参数中的元素
setmetatable(table, metatable) -- 设置表的元表
tonumber(e [, base])           -- 转换为数字
tostring(e)                   -- 转换为字符串
type(v)                       -- 返回值的类型
warn(...)                     -- 打印警告信息
xpcall(f, msgh [, arg1, ...]) -- 保护调用函数并捕获错误
```

>---

### 2. 协程库 (coroutine)

```lua
coroutine.close(co)            -- 关闭协程
coroutine.create(f)            -- 创建一个新的协程
coroutine.isyieldable()        -- 检查当前协程是否可以挂起
coroutine.resume(co, ...)      -- 恢复协程的执行
coroutine.running()            -- 返回当前正在运行的协程
coroutine.status(co)           -- 返回协程的状态
coroutine.wrap(f)              -- 创建一个包装了协程的函数
coroutine.yield(...)           -- 挂起协程的执行
```

>---

### 3. 包库 (package)

```lua
package.config                -- 包配置字符串
package.cpath                 -- C 库搜索路径
package.loaded                -- 已加载的包表
package.loadlib(libname, funcname) -- 加载动态链接库
package.path                  -- Lua 库搜索路径
package.preload               -- 预加载函数表
package.searchers             -- 模块搜索器表
package.searchpath(name, path [, sep [, rep]]) -- 搜索路径
```

>---

### 4. 字符串库 (string)

```lua
string.byte(s [, i [, j]])    -- 返回字符串中字符的 ASCII 值
string.char(...)              -- 将整数转换为字符串
string.dump(function [, strip]) -- 将函数转换为二进制字符串
string.find(s, pattern [, init [, plain]]) -- 在字符串中查找模式
string.format(format, ...)    -- 格式化字符串
string.gmatch(s, pattern)     -- 全局匹配模式
string.gsub(s, pattern, replacement [, n]) -- 全局替换
string.len(s)                 -- 返回字符串长度
string.lower(s)               -- 转换为小写
string.match(s, pattern [, init]) -- 匹配模式
string.pack(fmt, ...)         -- 将值打包为二进制字符串
string.packsize(fmt)          -- 计算打包格式的大小
string.rep(s, n [, sep])      -- 重复字符串
string.reverse(s)             -- 反转字符串
string.sub(s, i [, j])        -- 截取子串
string.unpack(fmt, s [, pos]) -- 从二进制字符串解包值
string.upper(s)               -- 转换为大写
```

>---
### 5. UTF-8 库 (utf8)

```lua
utf8.char(...)                -- 将码点转换为 UTF-8 字符串
utf8.charpattern              -- 匹配 UTF-8 字符的模式
utf8.codepoint(s [, i [, j]]) -- 返回字符串中字符的码点
utf8.codes(s)                 -- 迭代字符串中的码点
utf8.len(s [, i [, j]])       -- 计算 UTF-8 字符数
utf8.offset(s, n [, i])       -- 查找第 n 个字符的位置
```

>---

### 6. 表库 (table)

```lua
table.concat(list [, sep [, i [, j]]]) -- 连接表中的字符串
table.create(narr, nrec)         -- 创建指定大小的表
table.insert(list, [pos,] value) -- 插入元素到表
table.move(a1, f, e, t [, a2])   -- 移动表中的元素
table.pack(...)                  -- 将参数打包为表
table.remove(list [, pos])       -- 从表中移除元素
table.sort(list [, comp])        -- 排序表
table.unpack(list [, i [, j]])   -- 解包表为参数
```

>---

### 7. 数学库 (math)

```lua
math.abs(x)                   -- 绝对值
math.acos(x)                  -- 反余弦
math.asin(x)                  -- 反正弦
math.atan(x [, y])            -- 反正切
math.ceil(x)                  -- 向上取整
math.cos(x)                   -- 余弦
math.deg(x)                   -- 弧度转角度
math.exp(x)                   -- 指数函数
math.floor(x)                 -- 向下取整
math.fmod(x, y)               -- 取模
math.frexp(x)                 -- 将数字分解为尾数和指数
math.huge                     -- 正无穷大
math.ldexp(m, e)            -- 将 m 乘以 2 的 e 次方
math.log(x [, base])          -- 对数
math.max(x, ...)              -- 最大值
math.maxinteger               -- 最大整数值
math.min(x, ...)              -- 最小值
math.mininteger               -- 最小整数值
math.modf(x)                  -- 分解为整数和小数部分
math.pi                       -- π 值
math.rad(x)                   -- 角度转弧度
math.random([m [, n]])        -- 随机数
math.randomseed(x)            -- 设置随机种子
math.sin(x)                   -- 正弦
math.sinh(x)                  -- 双曲正弦
math.sqrt(x)                  -- 平方根
math.tan(x)                   -- 正切
math.tanh(x)                  -- 双曲正切
math.tointeger(x)             -- 将值转换为整数
math.type(x)                  -- 返回值的类型
math.ult(m, n)                -- 无符号整数比较
```

>---

### 8. IO 库 (io)

```lua
io.stderr                     -- 标准错误文件句柄
io.stdin                      -- 标准输入文件句柄
io.stdout                     -- 标准输出文件句柄
io.close([file])              -- 关闭文件
io.flush()                    -- 刷新输出缓冲区
io.input([file])              -- 设置或获取标准输入
io.lines([filename [, ...]])  -- 迭代文件行
io.open(filename [, mode])    -- 打开文件
io.output([file])             -- 设置或获取标准输出
io.popen(prog [, mode])       -- 执行命令并返回文件句柄
io.read(...)                  -- 从标准输入读取
io.tmpfile()                  -- 创建临时文件
io.type(obj)                  -- 检查对象是否为文件句柄
io.write(...)                 -- 写入标准输出
file:close()                  -- 关闭文件
file:flush()                  -- 刷新文件缓冲区
file:lines(...)               -- 迭代文件行
file:read(...)                -- 从文件读取
file:seek([whence [, offset]]) -- 设置文件位置
file:setvbuf(mode [, size])   -- 设置缓冲区模式
file:write(...)               -- 写入文件
```

>---

### 9. 日期时间库 (os)

```lua
os.clock()                    -- 程序运行时间
os.date([format [, time]])    -- 格式化日期时间
os.difftime(t2, t1)           -- 计算时间差
os.execute([command])         -- 执行系统命令
os.exit([code [, close]])     -- 退出程序
os.getenv(varname)            -- 获取环境变量
os.remove(filename)           -- 删除文件
os.rename(oldname, newname)   -- 重命名文件
os.setlocale(locale [, category]) -- 设置区域设置
os.time([table])              -- 获取当前时间或转换表为时间戳
os.tmpname()                  -- 生成临时文件名
```

>---

### 10. 调试库 (debug)

```lua
debug.debug()                 -- 进入交互调试模式
debug.gethook([thread])       -- 获取钩子设置
debug.getinfo([thread,] f [, what]) -- 获取函数信息
debug.getlocal([thread,] f, local) -- 获取局部变量
debug.getmetatable(value)     -- 获取元表
debug.getregistry()           -- 获取注册表
debug.getupvalue(f, up)       -- 获取上值
debug.getuservalue(u, n)         -- 获取用户值
debug.sethook([thread,] hook, mask [, count]) -- 设置钩子
debug.setlocal([thread,] f, _local, value) -- 设置局部变量
debug.setmetatable(value, table) -- 设置元表
debug.setupvalue(f, up, value) -- 设置上值
debug.setuservalue(udata, value, n)  -- 设置用户值
debug.traceback([thread [, message [, level]]]) -- 获取调用栈跟踪
debug.upvalueid(f, n)         -- 获取上值的唯一标识符
debug.upvaluejoin(f1, n1, f2, n2) -- 连接两个函数的上值
```

---