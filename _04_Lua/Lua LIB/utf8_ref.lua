utf8 = {}
utf8.charpattern = "[\0-\x7F\xC2-\xFD][\x80-\xBF]*"                     -- UTF-8 字节序列匹配模式

function utf8.char(code, ...) end         -- 从 Unicode 编码创建字符串

function utf8.codes(s, lax) end           -- 返回 UTF-8 字符迭代器

function utf8.codepoint(s, i, j, lax) end -- 返回 Unicode 编码值

function utf8.len(s, i, j, lax) end       -- 返回 UTF-8 字符数

function utf8.offset(s, n, i) end         -- 返回第 n 个字符的字节位置