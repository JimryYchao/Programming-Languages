-- clock, tmpname, rename
local start = os.clock()
local filename = os.tmpname()
local oldname = filename
filename = filename .. ".lua"
os.rename(oldname, filename)
print(filename)

-- exit, time, getenv, setlocale, date, difftime
local file = io.open(filename, "w")
if not file then
    os.exit()
end
file:write([=[
do
    local now = os.time()
    print("COMPUTERNAME=" .. (os.getenv("COMPUTERNAME") or ""))
    os.setlocale("C","time")
    print("Current Time [setlocale(C )]: " .. os.date("%c", now))
    os.setlocale("zh","time")
    print("Current Time [setlocale(zh)]: " .. os.date("%c", now))
    print "hello Lua548"
    print("difftime(now, 2000/1/1 00:00:00) = " ..
        string.format("%.fs", os.difftime(now, os.time(
        { year = 2000, month = 1, day = 1, hour = 0, min = 0, sec = 0 }))))
    print("Current Time in US: " .. os.date("%c", now))
end
]=])
file:close()

-- execute, remove
os.execute("lua " .. filename)
os.remove(filename)
print(string.format("Elapsed time: %fs", os.clock() - start))
