local t = { 1, 5, 3, 2, 4, [7] = 7, [9] = 9, [100] = 100 }
print(table.concat(t, ",")) -- concat index 1~5

-- insert, concat, remove, move
table.insert(t, 3, 10086)  -- insert at 3
table.insert(t, 10010)     -- insert at #t+1
print(table.concat(t, ",")) 

table.remove(t, 6)         -- remove at 6
local t2 = { 1000, 99, -10, -20, 0 }
table.move(t, 1, 6, #t2+1, t2)   -- move 1-6 of t to #t2+1 of t2
print(table.concat(t2, ",")) 


-- pack, unpack, sort
local Sort1_100 = function(...)
    local t = table.pack(...)
    local j = 1
    for i = 1, #t do
        if t[j] > 100 or t[j] < 0 then
            table.remove(t, j)
            j = j - 1
        end
        j = j + 1
    end
    table.sort(t)
    return t
end
local t3 = Sort1_100(table.unpack(t2))
print(table.concat(t3, ","))