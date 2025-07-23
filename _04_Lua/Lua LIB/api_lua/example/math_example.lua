-- 算术函数 abs, fmod, max, min, modf, ceil, floor
-- 圆相关   deg, rad
-- 三角函数 acos, asin, atan, cos, sin, tan
-- 幂函数   exp, sqrt,
-- 对数函数 log
-- 其他函数 tointeger, type

print(math.ceil(math.pi))    -- 4
print(math.floor(math.pi))   -- 3
print(math.fmod(3.1415, 2))  -- 1.1415
print(math.modf(3.1415))     -- 3   0.1415
print(math.deg(math.pi / 2)) -- 90.0
print(math.rad(60))          -- 1/3*pi

print(math.mininteger)       -- -9223372036854775808
print(math.maxinteger)       -- 9223372036854775807
print(math.pi)               -- 3.1415926535898
print(math.huge)             -- inf

-- random, randomseed
math.randomseed(10086, 10010)
print(math.random())             -- [0,1)
print(math.random(100))          -- [1,100]
print(math.random(10010, 10086)) -- [x,y]
