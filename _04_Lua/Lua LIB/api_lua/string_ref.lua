string = {}

function string.byte(s, i, j) end                 -- 返回字符的数字编码

function string.char(byte, ...) end               -- 从数字编码创建字符串

function string.dump(f, strip) end                -- 返回函数的二进制表示

function string.find(s, pattern, init, plain) end -- 查找模式首次出现位置

function string.format(s, ...) end                -- 格式化字符串

function string.gmatch(s, pattern, init) end      -- 返回模式匹配迭代器

function string.gsub(s, pattern, repl, n) end     -- 全局替换模式匹配

function string.len(s) end                        -- 返回字符串长度

function string.lower(s) end                      -- 转换为小写

function string.match(s, pattern, init) end       -- 返回模式匹配捕获

function string.pack(fmt, v1, v2, ...) end        -- 打包值为二进制串（Lua 5.3+）

function string.packsize(fmt) end                 -- 返回打包格式长度（Lua 5.3+）

function string.rep(s, n, sep) end                -- 重复字符串 n 次

function string.reverse(s) end                    -- 反转字符串

function string.sub(s, i, j) end                  -- 返回子字符串

function string.unpack(fmt, s, pos) end           -- 解包二进制串（Lua 5.3+）

function string.upper(s) end                      -- 转换为大写
