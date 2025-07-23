math = {}
math.huge = number                 -- 最大浮点数
math.maxinteger = integer          -- 最大整数（Lua 5.3+）
math.mininteger = integer          -- 最小整数（Lua 5.3+）
math.pi = number                   -- π 值

function math.abs(x) end           -- 绝对值

function math.acos(x) end          -- 反余弦

function math.asin(x) end          -- 反正弦

function math.atan(y, x) end       -- 反正切

function math.ceil(x) end          -- 向上取整

function math.cos(x) end           -- 余弦

function math.deg(x) end           -- 弧度转角度

function math.exp(x) end           -- e的x次方

function math.floor(x) end         -- 向下取整

function math.fmod(x, y) end       -- 取模

function math.log(x, base) end     -- 对数

function math.max(x, ...) end      -- 最大值

function math.min(x, ...) end      -- 最小值

function math.modf(x) end          -- 分离整数和小数

function math.rad(x) end           -- 角度转弧度

function math.random(m, n) end     -- 随机数

function math.randomseed(x, y) end -- 设置随机种子

function math.sin(x) end           -- 正弦

function math.sqrt(x) end          -- 平方根

function math.tan(x) end           -- 正切

function math.tointeger(x) end     -- 转换为整数（Lua 5.3+）

function math.type(x) end          -- 返回数字类型（Lua 5.3+）

function math.ult(m, n) end        -- 无符号小于比较（Lua 5.3+）

---

function math.tanh(x) end     -- 双曲正切（Lua 5.1）

function math.sinh(x) end     -- 双曲正弦（Lua 5.1）

function math.cosh(x) end     -- 双曲余弦（Lua 5.1）

function math.atan2(y, x) end -- 反正切（Lua 5.1）

function math.pow(x, y) end   -- x 的 y 次方（Lua 5.1）

function math.log10(x) end    -- 以 10 为底对数（Lua 5.1）

function math.ldexp(m, e) end -- m*(2^e)（Lua 5.1）

function math.frexp(x) end    -- 分解尾数和指数（Lua 5.1）
