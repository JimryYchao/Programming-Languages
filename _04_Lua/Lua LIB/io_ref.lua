---@class iolib
---@field stdin  file*
---@field stdout file*
---@field stderr file*
io = {}


function io.close(file) end          -- 关闭文件

function io.flush() end              -- 刷新输出缓冲

function io.input(file) end          -- 设置默认输入文件

function io.lines(filename, ...) end -- 返回文件行迭代器

function io.open(filename, mode) end -- 打开文件

function io.output(file) end         -- 设置默认输出文件

function io.popen(prog, mode) end    -- 启动进程

function io.read(...) end            -- 读取输入

function io.tmpfile() end            -- 创建临时文件

function io.type(file) end           -- 检查文件句柄类型

function io.write(...) end           -- 写入输出

file = {}

function file:close() end              -- 关闭文件

function file:flush() end              -- 刷新文件缓冲

function file:lines(...) end           -- 返回文件行迭代器

function file:read(...) end            -- 读取文件

function file:seek(whence, offset) end -- 设置/获取文件位置

function file:setvbuf(mode, size) end  -- 设置缓冲模式

function file:write(...) end           -- 写入文件
