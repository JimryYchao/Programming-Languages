table = {}

function table.concat(list, sep, i, j) end  -- 连接列表元素

function table.insert(list, pos, value) end -- 插入元素

function table.move(a1, f, e, t, a2) end    -- 移动元素（Lua 5.3+）

function table.pack(...) end                -- 打包参数为表（Lua 5.2+）

function table.remove(list, pos) end        -- 移除元素

function table.sort(list, comp) end         -- 排序列表

function table.unpack(list, i, j) end       -- 解包列表（Lua 5.2+）

---

function table.maxn(table) end -- 返回最大正索引（Lua 5.1）

---

function table.foreach(list, callback) end  -- 遍历表（已废弃）

function table.foreachi(list, callback) end -- 遍历数组（已废弃）

function table.getn(list) end               -- 返回表长度（已废弃）
