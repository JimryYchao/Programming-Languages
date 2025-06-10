local str = [[Hello 世界. 你好, Lua]]

-- byte, char, len, reverse, rep, upper, lower
local t = {}
for _, b in ipairs({ string.byte(str, 1, string.len(str)) }) do
    t[#t + 1] = b
end
print(table.concat(t, ","))
print(string.char(table.unpack(t)))
print(string.reverse("Hello World"))
print(string.rep("😀", 10)) -- 重复
print(string.upper(str))
print(string.lower(str))

-- find, match, gsub, gmatch
print(string.find("hello world", "hello"))                   --> 1 5
print(string.match("Today is 1/1/1970", "%d+/%d+/%d+"))      --> 1/1/1970
print(string.match(str, "%bHL"))                             --> Hello 世界. 你好, L
print(string.gsub("Lua is cute", "cute", "great"))           --> Lua is great
for w in string.gmatch(str, "%a+") do
    print(w)                                                 -- Hello Lua
end
s = [[then he said: "it's all right"!]]                      --
qStart, quotePart, qEnd = string.match(s, "([\"'])(.-)(%1)") -- 相当于 ([\"'])(.-)([\"'])
print(qStart .. quotePart .. qEnd)

-- pack, unpack, packsize
local author = {
    name = "JimryYchao",
    tel = 123456789,
    age = 18,
    add = "666 East Changan Avenue, Beijing"
}
local packfmt = "!<c" .. #author.name .. " j" .. " j" .. "!<c" .. #author.add
print(#packfmt)   -- 输出一个 format 打包后的字节大小
local pack = string.pack("!<c" .. #packfmt .. packfmt,
    packfmt, author.name, author.tel, author.age, author.add)
print(string.packsize(packfmt)) -- 打包后总计是 58 字节

local fmt = string.unpack("!<c14", pack)
local _, name, tel, age, add = string.unpack("!<c14" .. packfmt, pack)
print(name, tel, age, add)

-- 输出有序列表中的前 n 个元素
function F(file, n)
    assert(tonumber(n) > 0)
    local f = io.input(file)
    local counter = {}
    for line in io.lines(file, "l") do
        for word in string.gmatch(line, "%w+") do
            counter[word] = (counter[word] or 0) + 1
        end
    end
    io.close(f)
    local words = {}
    for w in pairs(counter) do
        words[#words + 1] = w
    end
    table.sort(words, function(w1, w2) -- 按出现次数降序
        return counter[w1] > counter[w2] or counter[w1] == counter[w2] and w1 < w2
    end)
    for i = 1, n > #words and #words or n do
        io.output(io.stdout)
        io.write(words[i], "\t", counter[words[i]], "\n")
    end
    io.close(io.stdout)
end
-- dump 二进制返回一个函数
local fb = string.dump(F, true)
local f = load(fb, "F", "b")
if f then
    f("Lua Lib/string.lua", 10)
end


