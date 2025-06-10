-- path, cpath
print(package.path)
print(package.cpath)

-- loaded, searchers
require "helloworld"   -- by package.searchers
package.loaded["helloworld"] = nil
os.execute("luac -o hello.lc helloworld.lua")
for key in pairs(package.loaded) do
    print(key)
end

-- searchpath
local file = package.searchpath("hello", package.path .. ";?.lc")
dofile(file)
os.remove(file)
