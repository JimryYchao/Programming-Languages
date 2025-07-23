lproc.start([[
print("start")
for i = 1, 5 do
    lproc.send("mess_queue", "Mess_"..i, i)
end
lproc.send("mess_queue", nil) --结束信号
lproc.exit()
]])

lproc.start([[
while true do
    local mess, i = lproc.receive("mess_queue")
    if not mess then break end
    print("receive:", mess, i)
end
lproc.exit()
]])

