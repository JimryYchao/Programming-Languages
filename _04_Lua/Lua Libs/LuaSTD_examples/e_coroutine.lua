-- wrap
local defer = {   -- 模拟 Go defer 语法
    __close = coroutine.wrap(function(self)
        for i = self.count, 1, -1 do
            self[i][1](table.unpack(self[i][2]))
            self[i] = nil
            self.count = self.count - 1
        end
    end),
    __call = function(self, func, ...)
        self.count = self.count + 1
        self[self.count] = { func, { ... } }
    end
}
defer.__index = defer
defer.New = function()
    return setmetatable({ count = 0 }, defer)
end

-- isyieldable, create, yield
function GetIteratorCO(table)
    if coroutine.isyieldable() then
        return nil
    end
    return coroutine.create(function()
        local de <close> = defer.New()
        de(function()
            print("coroutine closed")
        end)
        de(function()
            print("Hello World")
        end)

        for key, value in pairs(table) do
            coroutine.yield(key, value)
        end
    end)
end
-- resume
function GetKVFromCO(co)
    local rt = { coroutine.resume(co) }
    if rt[1] then
        return table.unpack(rt, 2)
    else
        return nil
    end
end

-- running, status, close
function close(co)
    local _, main = coroutine.running()
    if not main then
        return false
    end
    print(coroutine.status(co))
    return coroutine.close(co)
end

local mt = { 1, 2, 3, 4, 5, 6, 7, 8, 9 }
local co = GetIteratorCO(mt)

while true do
    local k, v = GetKVFromCO(co)
    if (v > 5) then
        close(co)
        break
    end
    print("Get KV = " .. k .. ", " .. v)
end
