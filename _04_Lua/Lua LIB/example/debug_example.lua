-- 交互模式
-- debug.debug()
-- getuservalue, setuservalue

-- getmetatable, setmetatable, sethook, gethook
debug.sethook(function ()
    print("Call Function")
end, "c")
print(debug.getmetatable("string"))
local mt <close> = debug.setmetatable({}, {__close = function ()
    print("debug test over")
end})
print(debug.gethook())
debug.sethook()

-- getregistry
local reg = debug.getregistry() -- lua 注册表
for k, v in pairs(reg[2]) do
    print(k, v)
end

-- getupvalue, setupvalue, getlocal, setlocal, upvalueid, upvaluejoin, getinfo
function F(x, y, level)
    if x < y then
        x = x + 1
        y = y - 1
        print("x = " .. x .. ", y = " .. y)
        local g
        g = function()                          -- 闭包变量顺序 _ENV,x,y,g
            print("x * y = " .. (x * y))
            local _, v = debug.getupvalue(g, 3) -- 获取闭包 y -> y
            print(debug.upvalueid(g, 3))        -- 上值 y id
            debug.setupvalue(g, 3, v - 1)       -- 设置闭包 y = v - 1
            print(debug.getupvalue(g, 3))
        end
        g()
        print("闭包 g 内更改捕获对象 y = " .. y)

        local i = 0
        local f
        f = function(V)                   -- 上层函数变量顺序 x,y,level,g,f
            i = 0
            debug.upvaluejoin(f, 1, g, 3) -- 令 f[v] = g[x]
            print("i = g[y] = " .. i)
            print("当前栈上 1 层 x = " .. x)
            local _, v = debug.getlocal(2, 1) -- 获取栈上层 up[2][1] x -> v
            debug.setlocal(2, 1, v + 1)       -- 设置栈上次 up[2][1] x = v + 1
            print("更改栈上 1 层 x = " .. x)
        end
        f()
        level = level + 1
        return F(x, y, level)
    end
    print(debug.traceback("F 递归结束", level))
end
local tf = debug.getinfo(F)
tf.func(1, 20, 1)