memlimit.set(10240 * 1024)  -- 设置内存上限为 10MB（单位：字节）

-- 测试内存分配
local function test_memAlloc()
    local data = {}
    print("开始分配内存...")
    for i = 1, 100000 do
        data[i] = string.rep("A", 1024) -- 每次分配1KB
        if i % 1000 == 0 then
            print("已分配:", i * 1024, "字节")
        end
    end
end
print(666)

-- 捕获内存不足错误
local status, err = pcall(test_memAlloc)
if not status then
    print("错误: ", err) -- 输出：错误: not enough memory (或 NULL 分配失败)
    print("当前内存占用: ", collectgarbage("count") * 1024, "字节") 
end

-- 强制触发垃圾回收
collectgarbage("collect")
print("GC 后内存占用: ", collectgarbage("count") * 1024, "字节")

