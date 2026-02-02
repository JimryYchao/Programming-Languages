-- assert, error, pcall, xpcall, warn
(function(a)
    local xp        = function()
        local f = function()
            assert(a > 1, "a must large than 1")
        end
        local ok, err = pcall(f, a)
        if not ok then
            error(err)
        end
    end
    local ok, trace = xpcall(xp, debug.traceback, 1)
    if not ok and trace then
        warn("@on")
        warn(trace)
        warn("@off")
    end
end)(1);

-- dofile, loadfile, load, require
(function()
    dofile("helloworld.lua")
    loadfile("helloworld.lua")()
    load([[print "Hello World"]])()
    require "helloworld"
end)();

-- setmetatable, getmetatable, rawequal, rawlen, rawset, rawget, select, tonumber, tostring
(function(...)
    local mt = {
        __index    = function(table, index)
            return rawget(table, index)
        end,
        __newindex = function(table, index, value)
            if index then
                rawset(table, index, value)
            end
        end,
        __tostring = function(self)
            local str = "table[" .. rawlen(self) .. "] { "
            for _, v in ipairs(self) do
                str = str .. v .. ", "
            end
            str = str .. " }"
            return str
        end
    }
    local m = setmetatable({}, mt)
    print(rawequal(getmetatable(m), mt))
    local num
    for i = 1, select("#", ...) do
        num = select(i, ...)
        m[i] = tonumber(num)
    end
    print(m) -- 调用 tostring
end)("1", "0xff", "0x12345", "10086", "1e5", "3.1415");

-- pairs, ipairs, next, type
(function(table)
    for index, value in ipairs(table) do
        if type(value) == type(1) then
            table[index] = value * value
        end
    end
    for key, value in pairs(table) do
        print(key, value)
    end
    -- 仅迭代键值对部分
    local KVIterator = function(table)
        local it
        it = function(table, key)
            local k, v = key, table[key]
            ::next::
            k, v = next(table, k)
            if type(k) ~= type(1) or not k then
                return k, v
            else
                goto next
            end
        end
        return it, table, nil
    end
    print("KVIterator: ")
    for k, v in KVIterator(table) do
        print(k, v)
    end
end)({ 1, 2, 3, 4, 5, a = "a", b = "b", c = "c" });

-- collectgarbage
(function()
    local mt = {
        __mode = "v",
        __gc = function()
            print("collected")
        end
    }
    collectgarbage("collect")
    print("Current Memory = " .. collectgarbage("count"))
    local m = setmetatable({ insert = { "Hello World" } }, mt)
    print("Current Memory = " .. collectgarbage("count"))
    collectgarbage("collect") -- 弱表 m 的 insert 被回收
    print("Current Memory = " .. collectgarbage("count"))

    collectgarbage("stop")
    local a = ""
    for i = 1, 1000 do
        a = a .. "Hello World" .. i
    end
    print("Current Memory = " .. collectgarbage("count"))
    m = nil     -- 调用 __gc
    collectgarbage("collect") -- 手动 GC 不影响暂停状态

    print("Current Memory = " .. collectgarbage("count"))
    if not collectgarbage("isrunning") then
        collectgarbage("restart")
        print("GC restart")
    end
end)()
