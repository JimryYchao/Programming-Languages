local file = "helloworld.lua"

-- input, read, close
io.input(file)
io.stdout:write(io.read() .. "\n")
io.close()

-- open, file:setvbuf, write, flush, close
local tmp = io.open(file, "a+")
if tmp then
    tmp:setvbuf("full", 16)
    tmp:write("\n" .. [[print "Lua548"]])
    tmp:flush()
    tmp:close()
end

-- popen, lines
local r = io.popen("luac -o hello.lc helloworld.lua")
io.close(r)
dofile("hello.lc")
local firstLine
for line in io.lines(file, "l") do
    firstLine = line
    break
end

-- type, output, write 
if io.type(tmp) == "closed file" then
    tmp = io.output(io.open(file,"w+"))
    tmp:setvbuf("no")
    io.write(firstLine)
    io.close()
end
os.remove("hello.lc")

-- file:seek, read
tmp = io.tmpfile()
print(tmp)
print(io.type(tmp)) -- file
local _, err = tmp:write(io.lines("Lua LIB/io.lua", "a")())
tmp:flush()
if not err then
    tmp:seek("set")
    local line = ""
    repeat
        line = tmp:read("l")
        print(line)
    until not line
end
tmp:close()
