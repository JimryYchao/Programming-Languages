coroutine = {}

function coroutine.close(co) end       -- 关闭协程（Lua 5.4+）

function coroutine.create(f) end       -- 创建新协程

function coroutine.isyieldable(co) end -- 检查协程是否可让出

function coroutine.resume(co, val1, ...) end -- 开始/继续协程

function coroutine.running() end             -- 返回当前运行协程

function coroutine.status(co) end            -- 返回协程状态

function coroutine.wrap(f) end               -- 创建协程包装函数

function coroutine.yield(...) end            -- 挂起协程
