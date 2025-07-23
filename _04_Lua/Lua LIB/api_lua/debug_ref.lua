debug = {}

function debug.debug() end                               -- 进入交互调试模式

function debug.gethook(co) end                           -- 获取钩子设置

function debug.getinfo(thread, f, what) end              -- 获取函数信息

function debug.getlocal(thread, f, index) end            -- 获取局部变量

function debug.getmetatable(object) end                  -- 获取元表

function debug.getregistry() end                         -- 获取注册表

function debug.getupvalue(f, up) end                     -- 获取上值

function debug.getuservalue(u, n) end                    -- 获取用户值

function debug.sethook(thread, hook, mask, count) end    -- 设置钩子

function debug.setlocal(thread, level, index, value) end -- 设置局部变量

function debug.setmetatable(value, meta) end             -- 设置元表

function debug.setupvalue(f, up, value) end              -- 设置上值

function debug.setuservalue(udata, value, n) end         -- 设置用户值

function debug.traceback(thread, message, level) end     -- 获取调用栈回溯

function debug.upvalueid(f, n) end                       -- 返回上值唯一标识（Lua 5.2+）

function debug.upvaluejoin(f1, n1, f2, n2) end           -- 连接上值（Lua 5.2+）

---

function debug.getfenv(o) end            -- 获取对象环境（Lua 5.1）

function debug.setfenv(object, env) end  -- 设置对象环境（Lua 5.1）

function debug.setcstacklimit(limit) end -- 设置C栈限制（已废弃）
