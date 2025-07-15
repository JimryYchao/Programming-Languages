local str = [[Hello 世界. 你好, Lua]]

-- len, offset, codepoint
print(#str, utf8.len(str, 1, #str)) -- 字节数 #str = 25, 字符数 len = 17
print("offset(你) = " .. utf8.offset(str, 11))     -- 14
print("codepoint(好) = " .. utf8.codepoint("好"))  -- 

-- codes, char
for _, c in utf8.codes(str) do
    if c ~= 32  then
        print(string.format("codepoint [0x%08X] = '%s'" ,c, utf8.char(c)))
    end
end