arg = {}                                       -- 独立版 Lua 的启动参数表

_G = {}                                        -- 全局环境表

_VERSION = "Lua 5.4"                           -- Lua 版本号

function assert(v, message, ...) end           -- 断言检查，失败时抛出错误

function collectgarbage(opt, ...) end          -- 控制垃圾收集器行为

function dofile(filename) end                  -- 执行指定文件内容

function error(message, level) end             -- 抛出错误

function getmetatable(object) end              -- 获取对象元表

function ipairs(t) end                         -- 返回数组迭代器

function load(chunk, chunkname, mode, env) end -- 加载代码块

function loadfile(filename, mode, env) end     -- 从文件加载代码块

function next(table, index) end                -- 遍历表中的键值对

function pairs(t) end                          -- 返回表迭代器

function pcall(f, arg1, ...) end               -- 保护模式调用函数

function print(...) end                        -- 打印输出到 stdout

function rawequal(v1, v2) end                  -- 原始相等比较

function rawget(table, index) end              -- 原始表访问

function rawlen(v) end                         -- 原始长度获取

function rawset(table, index, value) end       -- 原始表赋值

function select(index, ...) end                -- 选择参数子集

function setmetatable(table, metatable) end    -- 设置表元表

function tonumber(e) end                       -- 转换为数字

function tostring(v) end                       -- 转换为字符串

function type(v) end                           -- 返回类型名称

function warn(message, ...) end                -- 发出警告（Lua 5.4+）

function xpcall(f, msgh, arg1, ...) end        -- 带错误处理的保护调用

---

function unpack(list, i, j) end          -- 解包列表元素（Lua 5.1）

function getfenv(f) end                  -- 获取函数环境（Lua 5.1）

function loadstring(text, chunkname) end -- 从字符串加载代码（Lua 5.1）

function module(name, ...) end           -- 创建模块（Lua 5.1）

function setfenv(f, table) end           -- 设置函数环境（Lua 5.1）
