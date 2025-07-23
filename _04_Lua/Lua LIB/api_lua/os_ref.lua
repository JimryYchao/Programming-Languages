os = {}

function os.clock() end                     -- 返回 CPU 时间

function os.date(format, time) end          -- 格式化日期时间

function os.difftime(t2, t1) end            -- 计算时间差

function os.execute(command) end            -- 执行系统命令

function os.exit(code, close) end           -- 终止程序

function os.getenv(varname) end             -- 获取环境变量

function os.remove(filename) end            -- 删除文件

function os.rename(oldname, newname) end    -- 重命名文件

function os.setlocale(locale, category) end -- 设置区域设置

function os.time(date) end                  -- 获取当前时间或转换时间表

function os.tmpname() end                   -- 生成临时文件名
