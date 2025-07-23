local bits = BitArray.new(5)

bits[1] = 3.14
bits[2] = true
bits[4] = nil
bits[5] = 15

for i = 1, #bits, 1 do
    print(bits[i])
end

print(bits);