## Python 规范摘要

---

### 1. 基本概念


#### 1.1. 程序结构

Python 程序由代码块构成，每个名称绑定在对应代码块中，除非是 `nonlocal` 或 `global`。

```python
from typing import Optional,Generic,TypeVar

T = TypeVar("T") # 类型变量

global_var = "Hello,World!" # 全局变量

class Cls(Generic[T]): # 泛型类
    __class_var = "private"
    class_var = "public" 
    def __init__(self, value : T, __value : T): 
        self.public_var = value
        self.__private_var = __value
    def public_method(self) : 
        print(self.public_var, end=" ")
        print(self.__private_var)

c = Cls[str](global_var[0:5], global_var[6:])
c.public_method()
```

Python 默认源码编码为 utf-8。显式声明文件的编码 [codecs](https://docs.python.org/zh-cn/3/library/codecs.html#standard-encodings)：

```py
# -*- coding: __encoding__ -*- 
# 如
# -*- coding: cp1252 -*-
```

>---
#### 1.2. 运行时模型

当一个 Python 程序启动时，运行时执行模型大致为：

```txt
主机
  └── 进程
        └── Python 全局运行时 — 管理多个解释器
              └── Python 解释器 — 完整的 Python 运行时状态
                    └── Python 线程状态 — 线程专属数据
```

**全局运行时** 管理多个解释器的集合，处理它们之间的资源共享。**解释器** 保存 "Python 运行时"，包含 `sys.modules` 等持久状态。**线程状态** 包含线程专属数据如异常信息、调用栈、线程专属资源等。

每个线程状态始终关联到一个特定解释器和一个特定主机线程。同一主机线程可以关联多个线程状态，但同一时刻只能使用一个。线程状态之间隔离独立，不共享数据（可共享解释器及其中对象）。解释器之间相互隔离，各有自己的 `sys.modules`。


>---
#### 1.3. 数据模型

对象是 Python 对数据的抽象，所有数据都是由对象或对象间关系组成的。每个对象包含唯一标识号（或地址, `id(obj)`）、类型（`type(obj)`）和值。对象不会被显式销毁，当引用计数为 0 即无法访问时，自动垃圾回收。

数值、字符串、元组为不可变对象；列表、字典为可变对象。对于不可变对象，计算新值的操作可能会返回已有等价对象的引用。

```py
s = "hello"
s0 = "hello"
if s is s0:
    print(f"type = {type(s)}, id = {id(s)}, value = {s}")
# type = <class 'str'>, id = 2714914132240, value = hello
```

>---
#### 1.4. 词法元素

> **标准关键字和保留字**

| description    | keywords                                                                         |
| :------------- | :------------------------------------------------------------------------------- |
| 字面值         | `False`,`True`,`None`                                                            |
| 异常处理       | `raise`,`try`,`except`,`from`,`finally`,`else`                                   |
| 控制流         | `if`,`elif`,`else`,`for`,`while`,`break`,`continue`,`return`,`in`,`pass`,`yield` |
| 模块管理       | `import`,`from`, `lazy`                                                                |
| 声明           | `def`,`class`,`type`,`lambda`,`global`,`nonlocal`,`del`,`assert`,`as`     |
| 上下文管理     | `with`,`as`                                                                      |
| 异步与协程     | `async`,`await`                                                                  |
| 逻辑运算符     | `and`,`or`,`not`                                                                 |
| 成员与身份检查 | `in`,`not in`,`is`,`is not`                                                      |
| 模式匹配       | `match`,`case`                                                                   |

> **特殊标识符**

```py
# 输出: José
_*      # 不会被 from module import * 导入
_       # match case 模式中的 _ 通配符
__*__   # 系统定义名称 dunder
__*     # 类的私有名称
```

> **运算符和定界符**

```py
{ } [ ] ( )
, : ! ; = . @ ... ->
+ - * ** / // % & | ^ ~ << >>   
+= -= *= **= /= //= %= &= |= ^= <<= >>= @= :=     
< > <= >= == !=         
```

>---
#### 1.5. 包与模块

模块和包是 Python 组织代码的核心机制。每个 Python 文件是一个模块，包是一个特殊类型的模块。模块由 `import` 导入（或 `importlib.import_module()`, `__import__()`），`from ... import` 直接导入模块特定名称。

导入搜索首先查找 `sys.modules` 中是否存在模块缓存，找不到时唤起模块查找器与加载器，查找按导入路径搜索（内置模块 -> `sys.path` -> 包子模块）。

```python
import module               # 导入模块
from module import some_func, some_value # 导入特定成员
from package import *       # 导入所有公开成员（不推荐）
import module as m          # 别名
# 相对导入
from . import module           # 同级模块
from .module import some_func
from .. import parent_module   # 上级模块
from ..subpackage import mod   # 子包
# 显式惰性导入, 3.15+ 
lazy import large_module 
lazy from large_lib import some_func, some_value
# 兼容模式惰性导入
__lazy_module__ = ["large_module"]  
import lazy_module
```

包具有 `__path__` 属性。常规包通常以 `__init__.py` 文件的目录形式实现，该文件所定义的对象被绑定到该包的命名空间；命名空间包无需 `__init__.py` 文件。一个包的文件结构：

```txt
my_package/
├── __init__.py
├── module_a.py
├── module_b.py
└── sub_package/
    ├── __init__.py
    └── module_c.py
```

> 包和模块相关属性

| attr                     | description                             |
| ------------------------ | --------------------------------------- |
| `module.__name__`        | 模块名, 直接执行的模块为 `"__main__"`   |
| `module.__spec__`        | 模块规格                                |
| `module.__package__`     | 模块包名，顶级模块为 `''`               |
| `module.__path__`        | 模块路径                                |
| `module.__file__`        | 文件路径，可选                          |
| `module.__cached__`      | 模块已编译版本的路径，可选              |
| `module.__doc__`         | 模块文档字符串，可选                    |
| `module.__annotations__` | 变量标注，可选                          |
| `module.__annotate__`    | `annotate` 函数，可选                   |
| `module.__dict__`        | 模块命名空间                            |
| `package.__all__`        | 包公开成员列表，在 `__init__.py` 中定义 |

`__name__` 模块属性在直接运行时为 `"__main__"`，被导入时为模块名。

```python
# main.py
if __name__ == "__main__":
    main()  # 仅当直接运行此文件时执行
```

`__all__` 包公开成员列表，在 `__init__.py` 中定义。

```python
# __init__.py
__all__ = ['some_func','some_class','some_value']
```
> `__future__` 

`future` 语句允许在当前版本中使用未来版本的特性和语法。

```python
# 模块开头声明
from __future__ import [feature [as identifier]]*

# 在 Python2 中使用 Python3 的 print 函数
from __future__ import print_function

print("Hello", "World", sep="-")  # print
```


<!-- #### 1.3. 执行环境

Python 解释器执行字节码（`.pyc` 文件）。首次导入模块时编译为字节码，后续导入直接使用缓存。`python -m py_compile` 预编译模块。

```powershell
python script.py           # 直接运行脚本
python -c "print('Hello')"  # 执行单行代码
python -i script.py         # 运行后进入交互模式
python -m module_name      # 作为模块运行
```

>*** -->


>---
#### 1.6. 惰性求值

标注作用域中的标注、类型别名值、类型变量绑定/约束/默认值采用 **惰性求值**：

```py
x = 100
def foo(n : int) -> x + 5 :   # 标注
    return n
x = 200
print(foo.__annotations__) # {'n': <class 'int'>, 'return': 205}  惰性求值
x = 300
print(foo.__annotations__) # {'n': <class 'int'>, 'return': 205}

type Alias = 1/0           # 不会立即报错
Alias.__value__            # ZeroDivisionError 此时才触发
```

---
### 2. 内置类型

| type               | description                              |
| :----------------- | :--------------------------------------- |
| NoneType           | `None`                                   |
| NotImplementedType | `NotImplemented`                         |
| Ellipsis 占位      | `...`,`Ellipsis`                         |
| Number             | `int`,`bool`,`float`,`complex`           |
| Sequence           | `list`,`bytearray`,`str`,`bytes`,`tuple` |
| Set                | `set`,`frozenset`                        |
| Mapping            | `dict`,`frozendict`                                   |

>---
#### 2.1. Number

Number 包含 `int`,`float`,`complex`。`bool` 为 `int` 子类型。整数具有无限精度，浮点数通常由 C double 实现。

```py
# 整数
2_147_483_647
79228162514264337593543950336792281625     # 长度没有限制，仅受内存限制
0b_1100_0101     # 二进制, 0B
0o377            # 八进制, 0O
0xDead_Beef      # 十六进制, 0X
# 浮点数
3.14_15_92653
10.      # 10.0
.01      # 0.01
1e3      # E-计数法
f : float = "nan"   # float('nan')     
+/- "inf"
# 复数
1 + 2j
c =  1 + 2J
real = c.real             # 1.0
imag = c.imag             # 2.0
conjugate = c.conjugate() # 1-2j
```

> 布尔逻辑检测

任何对象都可以进行逻辑监测，一个对象被视为 `True`，除非对象所属的类上定义了 `__bool__()` 返回 `False` 或 `__len__()` 返回 `0`。被视为假值的内置对象包括 `None`、`False`、数值零、空序列 （`''`、`[]`、`{}`、`set()`、`range(0)`）。

```py
class C:
    def __bool__(self):
        return False
if not C():
    print("C() is false")
```

>---
#### 2.2. Sequence

序列表示以非负数为索引的有限有序集合，基本序列类型分为 `list`,`tuple`,`range`。不可变序列包含 `str`,`tuple`,`bytes`；可变序列包含 `list`,`bytearray` 或 `collections` 和 `array` 模块中定义的可变序列。可变序列不支持 `hash()` 操作。

> 通用序列操作

| description | operator                             |
| :---------- | :----------------------------------- |
| 成员检测    | `x in s`, `x not in s`               |
| 序列拼接    | `s1 + s2`, `s * n`(n 次自身拼接)     |
| 索引        | `s[i]`, `s[-i]`(从右索引)            |
| 切片        | `s[i:j]`, `s[i:j:step]`(step 为步长) |
| 内置函数    | `len(s)`, `max(s)`, `min(s)`         |

```py
arr = [1, 2, 3, 4, 5, 6, 7, 8, 9]
a = arr[::2]   # [1, 3, 5, 7, 9]， 从索引 0 开始，步长为 2
e = arr[len(arr)//2]   # 5
s1 = arr[:2]   # [1, 2]
s2 = arr[-2:]  # [8, 9]  
s = s1 + [e] + s2   # [1, 2, 5, 8, 9]

"ll" in "Hello"   # True

_list = [[1]] * -1   # <0 被视为 0，返回空序列 []
lists = [[]] * 3     # [[], [], []]，元素 [] 对象值不会被复制，仅被多次复制引用 
lists[0].append(2)   # [[2], [2], [2]]，全部引用 []

lists = [[] for i in range(3)]  # [] 重复创建 3 次，多维列表
lists[0].append(1)   # [[1], [], []]
```

> 可变序列还支持

| description                     | operator                                 |
| :------------------------------ | :--------------------------------------- |
| 切片替换                        | `s[i] = x`, `s[i:j] = t`(可迭代对象)     |
| 切片指定步长元素替换为 t 的元素 | `s[i:j:k] = t`(长度相同)                 |
| 元素删除                        | `del s[i]`, `del s[i:j]`, `del s[i:j:k]` |
| 切片自扩展                      | `s += t`, `s *= n`                       |

```py
l = [1,2,3,4,5,6,7,8,9]
l[1::2] =[6 for i in range(len(l[1::2]))]  # [1, 6, 3, 6, 5, 6, 7, 6, 9]
del l[1::2]      # [1, 3, 5, 7, 9]
l[1:-1] = [6,6]  # [1, 6, 6, 9]
l[1:-1] *= 2     # 对 l[1:-1] 范围的切片自扩展 2 倍
print(l)   # [1, 6, 6, 6, 6, 9]
```

> 拼接不可变序列

拼接不可变序列总是生成新对象，运行时开销 $O(n^2)$，可以使用线性开销 $O(n)$ 的替代方案：
- 字符串 str 采用 str.join() 或写入 `io.StringIO()`。
- 字节串 bytes 采用 bytes.join() 或写入 `io.BytesIO()`。
- 元组 tuple 扩展为列表。

```py
s = "Hello"
s = " ".join([s, "World"])   # "Hello World"

b = b"Hello"
s = b" ".join([b, b"World"]) # b"Hello World"

t = (1, 2, 3)
l = [i for i in t] + [4,5]   # [1, 2, 3, 4, 5]
```

>---
#### 2.3. list

list 为可变序列，常存放同类项的集合：

```py
l = []
l = [1, 2, 3, 4, 5]
l = [i for i in iterable]  # [i for i in range(1, 6)]
l = [i for i in iterable if cond]  # [i for i in range(1, 6) if i % 2 == 0]
l = list()   # []
l = list(iterable)   # list("abc") = ['a', 'b', 'c'], list((1, 2, 3)) = [1, 2, 3]
```

>---
#### 2.4. tuple

tuple 元组为不可变序列，常存放异构数据的多项集：

```py
t = ()      # 空元组
t = a,      # 或 t = (a,) 单项元组 
t = a,b,c   # t = (a, b, c)
t = tuple(iterable)   # tuple("abc") = ('a', 'b', 'c'), tuple([1, 2, 3]) = (1, 2, 3)
x,y,z = t   # 元组解包
```

>---
#### 2.5. range

range 表示不可变的数字序列，常用于循环：

```py
for i in range(10):   # 0 ~ 9
    pass

for i in range(1, 6):   # 1 ~ 5
    pass

for i in range(1,6,2):   # 1, 3, 5; step=2
    pass

for i in range(1,10) if i % 2 == 0:
    pass
```

>---
#### 2.6. str

str 字符串由 `"` 或 `'` 作为字面值引导，三引号表示为多行字符串。str 是 Unicode 码位（`0`~`0x10FFFF`）构成的不可变序列，内置函数 `ord()` 将字符（单字符字符串）转换为 Unicode 码位。

```py
"Hello" + 'World'      # 字符串拼接, 字面值可省略 +
"This string will not include \ 
backslashes or newline characters." 
s : str =R"This a raw string."         # 原始字符串
sbs : bytes =b"This a bytes string."   # 字节串
# 多行
"""Multi line"""
'''first line   
second line'''

# 转义序列
'\ ooo'     # 000 ~ 377 
'\x hh'
'\N{name}'  # 命名 Unicode 字符，'\N{SNAKE}' = '🐍'
'\u hhhh'
'\U hhhhhhhh'
'\\','\'','\"','\a','\b','\f','\n','\r','\t','\v'
```

> f-格式化字符串

`f"{}"` 标记替换字段，`{{ }}` 转义为 `{ }`。

```py
who = 'nobody'
nationality = 'Spanish'
f'{who.title()} expects the {nationality} Inquisition!'
```

替换字段还可以是 
- `=` 调试说明符
- 转换说明符：`!s`、`!r`、`!a` 表示 `str()`、`repr()`、`ascii()`。替换字段默认 `!s`；`=` 未使用格式说明符时默认 `!r`。
- 格式说明符：`expr:*[^,<,>][+][width][.precision]` 调用 `format()` 格式化；`*` 表示填充字符（默认空格）；`^` 居中对齐，`<` 左对齐，`>` 右对齐，`0` 零填充。`+` 显示数值符号位。

```py
from fractions import Fraction
one_third = Fraction(1, 3)
# 调试说明符 默认 !r
print(f"{one_third = }")    # 'one_third = Fraction(1, 3)'
# 格式说明符 
print(f'{one_third:7.3f}')  # '  0.333'
print(f'{10086:_^+10}')  # '__+10086__' 居中
# 转换说明符 
print(f'{one_third!s} is {one_third!r}') # '1/3 is Fraction(1, 3)'
```

格式字符串不用作 `__doc__` 文档字符串。

```py
def foo():
    f"Not a docstring"
print(foo.__doc__)   # None
```

> t-模板字符串

`t"{}"` 模板字符串与格式化字符串相同语法，延迟求值生成 `string.templatelib.Template` 对象。

```py
from html import escape
from string.templatelib import Template, Interpolation

def render_html(template: Template) -> str:
    """一个简单的HTML渲染器，它会自动转义所有动态内容"""   # doc string
    result_parts = []
    for item in template:
        if isinstance(item, str):
            # 静态文本部分，直接添加
            result_parts.append(item)
        elif isinstance(item, Interpolation):
            # 动态插值部分，进行HTML转义后再添加
            result_parts.append(escape(str(item.value)))
    return "".join(result_parts)

user_comment = "<script>alert('xss')</script>"
safe_template = t"<p>User says: {user_comment}</p>"
final_html = render_html(safe_template)
#<p>User says: &lt;script&gt;alert(&#x27;xss&#x27;)&lt;/script&gt;</p>
```

> 字符串插值运算符 %

`format % values` 转换标记符 `%` 为 `values` 中的对应值。values 是相同项数的元组或映射。

```py
print("%s has %d quote types" % ('Python', 2))   # 元组
print("%(language)s has %(number)d quote types" % 
    {'language': 'Python', 'number': 2})   # 映射
# Python has 2 quote types
```

转换标记符包含：

```py
'%<(mapKey)>?<flag>?<width>?<.precision>?<len_spec> type_spec'
''' 
flag : 
    '#'  替代形式
    '0'  零填充
    '-'  左对齐
    '+'  符号位
    ' '  正数留空
len_spec:  h, l/L  
type_spec:  
    %d,%i   有符号十进制整数  %u 等价于 %d
    %o      有符号八进制整数    
    %x,%X   有符号十六进制整数 
    %e,%E   科学计数法浮点数 
    %f,%F   浮点十进制
    %g,%G   浮点数，精度小于 -4 相当于 %e,%E
    %c      单个字符，整数或单字符串
    %r      repr() 字符串
    %s      str() 字符串
    %a      ascii() 字符串
    %%      %
# 替代形式说明:
    %#o             插入八进制数前缀 0o
    %#x,%X          插入十六进制数前缀 0x 或 0X
    %#f,%F,%e,%E    保留小数后 6 位
    %#g,%G          保留有效位 6 位
'''
print("%#g" % 31657646.14)    # 3.16576e+07
print("%#X" % 10086)          # 0X2766
```

> str.format()

```py
print("{} and {}".format("spam", "eggs"))   # 默认位置参数 0,1,...
print("{1} and {0}".format("spam", "eggs")) # 指定位置
print("{0} and {0}".format(("spam","eggs"))) # 传递元组
print("{0} and {eggs}".format("spam", eggs="eggs")) # 混合位置和关键字参数
print("{0[s]:s} and {0[e]:s}".format({"s": "spam", "e": "eggs"})) # 名称引用，传递字典
print("{s} and {e}".format(**{"s": "spam", "e": "eggs"})) # 字典解包
```

>---
#### 2.7. bytes, bytearray, memoryview

bytes 字节串为不可变数组，保存 8 位二进制字节（0~255 的整数）。

```py
b1 = b'hello'
b2 = "中文".encode("utf-8")
b3 = bytes([72, 101, 108, 108, 111])  # b'Hello'
b4 = b'hello' + b'world'
```

bytearray 字节数组为可变序列，由 `bytearray()` 创建。

```py
ba1 = bytearray(b'hello' + b'world')
ba1[0:6] = b'HELLO'   # b'HELLOworld'
ba2 = bytearray("中文", encoding="utf-8")
```

memoryview 内存视图为只读视图，引用底层内存，支持切片语法，可以操作支持缓冲区协议的对象如 bytes、bytearray、array 等。

```py
# 从支持缓冲区协议的对象创建
data = bytearray(b'Hello, World!')
view = memoryview(data)
print(view[7:12].tobytes())  # b'World'
# 修改原数据会影响视图（反之亦然） 
view[0] = 104  # b'hello, World!'
```

>---
#### 2.8. set, frozenset

set 为无序且不重复的可变可迭代集合。frozenset 为不可变可迭代冻结集合。

```python
fs = frozenset({1, 2, 3})

s = {1, 2, 3}
s.add(4)
s.remove(1)       # 元素不存在时抛出 KeyError
s.discard(99)     # 元素不存在时不报错
s1 | s2           # 并集
s1 & s2           # 交集
s1 - s2           # 差集
s1 ^ s2           # 对称差集
```


>---
#### 2.9. dict, frozendict

dict 字典为可变映射，键必须可哈希。字典保留插入顺序。frozenset 为不可变可迭代冻结字典。

```python
d = {'a': 1, 'b': 2}
d['a']           # 1
d.get('c', 0)    # 0（默认值）
d.keys()         # dict_keys(['a', 'b'])
d.values()       # dict_values([1, 2])
d.items()        # dict_items([('a', 1), ('b', 2)])
d.update({'c': 3})
```

---
### 3. 表达式

> 运算符与优先级

| Category           | Operators                                                                                                                                  |
| :----------------- | :----------------------------------------------------------------------------------------------------------------------------------------- |
| 基本表达式         | `(exprs...)`,`[exprs...]`,`{key:value...}`,`{exprs...}`,`x[i]`,`x[i:j]`,`x[i:j:step]`,`f(params)`,`x.attr`                                 |
| 异步表达式         | `await expr`                                                                                                                               |
| 幂运算             | `x ** y`                                                                                                                                   |
| 一元               | `+x`,`-x`,`~x`                                                                                                                             |
| 乘法和矩阵乘法     | `x * y`,`x / y`,`x // y`,`x % y`,`x @ y`                                                                                                   |
| 加法               | `x + y`, `x - y`                                                                                                                           |
| 移位               | `x << y`, `x >> y`                                                                                                                         |
| 位运算             | `x & y`,`x ^ y`,`x \| y`                                                                                                                   |
| 比较               | `x < y`,`x > y`,`x <= y`,`x >= y`,`x == y`,`x != y`                                                                                        |
| 成员检测与标识测试 | `a is b`,`a is not b`,`x in y`,`x not in y`                                                                                                |
| 布尔逻辑(短路)     | `not x`,`x and y`,`x or y`                                                                                                                 |
| 条件表达式         | `x if cond else y`                                                                                                                         |
| lambda             | `lambda args: expr`                                                                                                                        |
| 赋值表达式         | `x := y`,`x = y`,`x += y`,`x -= y`,`x *= y`,`x **= y`,`x /= y`,`x //= y`,`x %= y`,`x @= y`,`x &= y`,`x \|= y`,`x ^= y`,`x <<= y`,`x >>= y` |

> 链式比较

```py
x < y < z  # 等价于 (x < y) and (y < z)
```

> 海象赋值

`:=` 在表达式中赋值。

```py
x = a if (a := func(b)) > 0 else b

while (data:= input("Input : ")) != "#":
    print(data)
```

>---
#### 3.1. 打包与解包：*, **

*加星目标赋值* 对 iterable 部分元素打包为 list。

```py
a, *b = [1]  # 1, []
*a, = [1, 2, 3]       # [1, 2, 3]
[*a] = [1, 2, 3]      # [1, 2, 3]
a, *b = [1, 2, 3, 4, 5]     # 1, [2, 3, 4, 5]
*a, b = [1, 2, 3, 4, 5]     # [1, 2, 3, 4], 5
a, *b, c = [1, 2, 3, 4, 5]  # 1, [2, 3, 4], 5

d = {"a": 1, "b": 2, "c": 3, "d": 4}
k1, *otherKs = d   # 'a', ['b', 'c', 'd']
(k1,v1), *others = d.items() # ('a', 1), [('b', 2), ('c', 3), ('d', 4)]
*otherVs, v4 = d.values()    # [1, 2, 3], 4
```

`*` 和 `**` 解包 iterable。

```py
# 函数解包
def f(a, b, c): pass
f(*[1,2,3])
f(**{"a":1, "b":2, "c":3})  # f(a=1, b=2, c=3)  

# list/set/tuple 合并
list1 = [1, 2, 3]
list2 = [4, 5, 6]
combined = [*list1, *list2] # [1, 2, 3, 4, 5, 6]

# 字典合并
d1 = {"a": 1, "b": 2}
d2 = {"c": 3, "d": 4}
merged = {**d1, **d2}  # {'a': 1, 'b': 2, 'c': 3, 'd': 4}
keys = {*d1, *d2}      # {'a', 'b', 'c', 'd'}
values = {*d1.values(), *d2.values()}    # {1, 2, 3, 4}

# 组合为单一平面结构
lists = [[1,2], [4,5], [7,8]]
[*l for l in lists] # [1, 2, 4, 5, 7, 8], 相当于 [x for l in lists for x in l]
sets = [{1,2}, {2,3}, {3,4}]
[*s for s in sets] # [1, 2, 3, 4], 相当于 {x for s in sets for x in s}
dicts = [{'a': 1}, {'b': 2}]
[**d for d in dicts] # [{'a': 1}, {'b': 2}], 相当于 [k: v for d in dicts for k, v in d.items()]
gen = (*L for L in lists) # 生成器表达式
list(gen) # [1, 2, 4, 5, 7, 8]

```

>---
#### 3.2. 布尔逻辑：not, and, or

```py
not x       # x is false ? True : False
x or y      # x is true ? x : y
x and y     # x is true ? y : x
```

> 短路求值

```python
a = False and print("不会执行")  # 返回 False，不执行右侧
b = True or print("不会执行")    # 返回 True，不执行右侧
```

>---
#### 3.3. 推导式

```python
# 列表推导式
squares = [x ** 2 for x in range(10)]
evens = [x for x in range(20) if x % 2 == 0]
result = [x async for x in async_gen()]
# 集合推导式
unique_squares = {x ** 2 for x in [-1, 1, 2]}
# 字典推导式
word_lengths = {word: len(word) for word in ["hello", "world"]}
# 生成器推导式
gen = (x*10+y for x in range(10) for y in range(10))  # 惰性求值
next(gen) # 0
for i in gen:
    print(i)  # 1, 2, 3, ..., 99
```

---
### 4. 语句

#### 4.1. 空操作：`pass`

```py
def f(arg): 
    pass    # 空操作方法
class C:    
    pass    # 空操作类
```

>---

#### 4.2. 断言：assert

```py
assert expr [, expr2]
# 相当于
if __debug__
    if not expr: raise AssertionError(expr2)

assert 1 > 2, "Assert failed"
```

>---
#### 4.3. 作用域关联：global, nonlocal

> global

`global` 在模块、函数域或类定义中绑定全局变量名称。

```py
def increment():
    count = 10086     # 绑定为局部变量

def increment_global():
    global count      # 绑定为全局变量
    count = count + 1

count = 0
increment()         # count = 0
increment_global()  # count = 1
```

> nonlocal

`nonlocal` 在嵌套函数或类中引用外层函数域的变量，常用于函数闭包。

```py
def outer():
    def inner():
        nonlocal count  # 声明 count 为外层函数域的变量
        count += 1
    count = 0
    inner()
    print(count)  

outer()  # 1
```

>---
#### 4.4. 类型别名：type

```py
type Id [type_params] = expr

type Point = tuple[float, float]
type ListOrSet[T] = list[T] | set[T]
```

>---
#### 4.5. 删除：del

```py
del a [, b, ... ]
del list[i]
del list[i:j]
del list[i:j:step]
del m['key']
del obj.attr
```

>---
#### 4.6. 跳转语句：break, continue, return

> break, continue

在 `for`,`while` 中出现 `break`, `continue`。

```py
while condition:
    if break_condition:
        break
    if skip_condition:
        continue
else:
    # 循环正常结束时执行（未被break终止）
```

> return 

`return` 多值返回的实际对象为元组。

```py
def f():
    return 0,1,2    # return (0, 1, 2)
x,y,z = f()   # 元组解构
```


>---
#### 4.7. 条件语句：if, elif, else

```py
if cond1:
    stmt
(elif cond2:
    stmt)*
[else:
    stmt]
```

>---
#### 4.8. 迭代语句：for, while, async for

> for

```py
for target_list in iterable:
    stmt
[else:
    stmt]

for i in range(10):
    print(i) 
for i, value in enumerate(items):
    print(i, value)
for key, value in dictionary.items():
    print(key, value)
for a, b in zip(list1, list2):
    print(a, b)
print(i) # i 在离开循环后仍存在
```

> while

```py
while cond:
    stmt
[else:
    stmt]
```

> async for

```py
async for target in async_iterable:
    stmt1
[else:
    stmt2]
# 在语义上等价于：
aiter = async_iterable.__aiter__()
running = True
while running:
    try:
        target = await aiter.__anext__()
    except StopAsyncIteration:
        running = False
    else:
        stmt1
else: 
    stmt2
```


>---
#### 4.9. 上下文管理器：with, async with

（异步）上下文管理器（对象）定义执行 `[async] with` 语句的运行时上下文，常用于保存和恢复各种全局状态，锁定和释放资源，关闭文件等。

```py
[async] with e1 as t1 [, e2 as t2, ...]:
    stmt

with open("file.txt") as f:
    content = f.read()  # f.close() 自动关闭
with open("input") as f1, open("output", "w") as f2:
    f2.write(f1.read())
```

> 上下文管理器对象

```py
# 上下文管理器协议
class ManagedResource:
    def __enter__(self):
        return self
    def __exit__(self, exc_type, exc_value, traceback):
        # do some cleanup or exception-handle
        return False or None # 不抑制异常

# 异步上下文管理器协议
class AsyncManagedResource:
    async def __aenter__(self):
        [await] STMT
        return self
    async def __aexit__(self, exc_type, exc_value, traceback):
        # do some cleanup or exception-handle
        [await] STMT 
        return False or None # 不抑制异常
```

`with` 语句的执行过程：

```py
[async] with EXPRESSION as TARGET:
    STMT
# 在语义上等价于
manager = (EXPRESSION)
enter = manager.__enter__  # or __aenter__
exit = manager.__exit__    # or __aexit__
value = [await] enter()
hit_except = False
try:
    TARGET = value
    STMT
except:
    hit_except = True
    if not [await] exit(*sys.exc_info()):
        raise
finally:
    if not hit_except:
        [await] exit(None, None, None)
```


>---
#### 4.10. 模式匹配：match

```py
match subject_expr:
    case pattern <if guard>:  
        stmt
# 模式
    # 或模式
    case P1 | P2 <| P3 ...> : pass 
    # 绑定模式
    case P as v: pass
    # 字面值模式, number,str,None,True,False,...
    case literal : pass 
    # 值名称，匹配已存在的常量
    case name_or_attr: pass 
    # 组模式
    case ( P ) : pass
    # 序列模式，[a,*list,b],(a, b1 | b2)
    case [ sequence ] | ( sequence ) : pass  
    # 字典模式
    case { key_value <, key_value ...> <,**kwargs> } : pass
    # 类模式
    case ClassName(P | positionP | keywordP) : pass
    # 捕获模式，始终匹配，
    case name : pass
    # 通配符模式，始终匹配，不捕获，与捕获模式互斥
    case _ : pass     
```

> 类模式

```py
class Point:
    __match_args__ = ("x", "y")  # 位置参数

    def __init__(self, x, y, z=0):
        self.x = x
        self.y = y
        self.z = z

def match_point(point):
    match point:
        # 位置参数匹配 __match_args__
        case Point(0, 0):  
            print("Origin")
        case Point(0, y):
            print(f"Y-axis: {y}")
        case Point(x, 0):
            print(f"X-axis: {x}")
        # 关键字匹配, 不依赖 __match_args__
        case Point(x=1, y=2):
            print("Point(1, 2)")
        case Point(x=x, y=y, z=z) if z > 1:  # guard
            print(f"Point({x}, {y}, {z})")
        # 混合模式
        case Point(-10, -10, z=-10):
            print(f"Point({-10}, { -10}, {z})")
        case _:
            print("Unknown point")
```

---

### 5. 函数

```python
[@decorator]+
def funcName [typeParams] ([param [:p_annotation][= defaultValue]]+ ) [-> func_annotation]:
    [""" 函数文档字符串 __doc__ """]
    funcBody

# 默认参数
def add(a, b=10):
    return a + b
# 函数注解
def func(
    arg1: int,
    arg2: str = "default",
    *args: tuple,
    **kwargs: dict
) -> bool:
    pass
```

> 位置参数与关键字参数

位置参数按位置传递，关键字参数按名称传递。`/` 强制前置参数为位置参数，`*` 强制后置参数为关键字参数。

```py
def func(a, b, /, c=0, *, d=0): pass
func(1, 2, d=3)
func(a=1, b=2)     # ERR，a,b 是位置参数
func(1,2,3)    
func(1,2,3,4)      # ERR，d 是关键字参数
```

> 可变参数

```py
# *arg 打包为元组
# **kwargs 打包为字典
# 可变参数 *arg, **kwargs 位于任何其他参数之后
def func(*args, **kwargs): pass
func(1,2,3, a=1, b=2)  # 相当于 func((1,2,3), {'a':1, 'b':2})
```

> 函数闭包

```py
def make_adder(x):
    def adder(y):
        return x + y
    return adder
add_5 = make_adder(5)
add_5(3)  # 8
```


>***
#### 5.1. 标注声明

变量或函数可能带有标注，常用于类型提示。

```py
x : int = 1
def func(param : p_annotation) -> rt_annotation: pass
# 惰性求值
def func(p : x) -> int: pass
x = 2
print(func.__annotations__) # {'p': 2, 'return': <class 'int'>}
```

存在 `from __future__ import annotations` 时，标注将存储为字符串。

```py
from __future__ import annotations
def func(p : 1) -> int: pass
print(func.__annotations__) # {'p': '1', 'return': 'int'}
```

>---

#### 5.2. 装饰器函数

装饰器允许动态修改函数或类的行为，通常用于日志记录、性能分析、权限检查或函数缓存等。


```python
@decorator
def funcName(): pass  # 相当于 funcName = decorator(funcName)

@decorator(args)
def funcName(): pass  # 相当于 funcName = decorator(args)(funcName)

# 装饰器函数
def decorator(func):
    def wrapper(*args, **kwargs):
        print("前置处理 in decorator")  
        # 前置处理
        result = func(*args, **kwargs)  # 接收所有参数
        # 后置处理
        return result
    return wrapper

# 参数装饰器函数
def decorator_with_args(params):
    def inner_decorator(func):
        def wrapper(*args, **kwargs):
            print("前置处理 in decorator_with_args")
            # 前置处理
            result = func(*args, **kwargs)  # 接收所有参数
            # 后置处理
            return result
        return wrapper
    return inner_decorator

# 类装饰器
class decorator_class:
    def __init__(self, func):
        self.__func = func
    def __call__(self, *args, **kwargs):
        print("前置处理 in decorator_class")
        # 前置处理
        result = self.__func(*args, **kwargs)  # 接收所有参数
        # 后置处理
        return result

@decorator
@decorator_with_args("params") 
@decorator_class
def funcName(): pass
funcName() 
# 前置处理 in decorator
# 前置处理 in decorator_with_args
# 前置处理 in decorator_class
```

> **常用装饰器**

```python
@staticmethod      # 静态方法
@classmethod       # 类方法
@property          # 属性方法
@dataclass         # 数据类
```

>---
#### 5.3. 迭代器：iter, next, aiter, anext

```py
# 迭代器协议
class Iterator:
    """同步迭代器必须实现 __iter__ 和 __next__"""
    def __iter__(self):
        return self | iterable
    def __next__(self):
        # 返回下一个值，结束时抛出 StopIteration
        pass

# 异步迭代器协议  
class AsyncIterator:
    """异步迭代器必须实现 __aiter__ 和 __anext__"""
    def __aiter__(self):
        return self | iterable
    async def __anext__(self):
        # 异步返回下一个值，结束时抛出 StopAsyncIteration
        pass

for item in Iterator():
    print(item)
async for item in AsyncIterator():
    print(item)
```

>---
#### 5.4. 生成器：yield, yield from

`yield` 定义（异步）生成器函数，返回一个 generator 迭代器。`yield from` 传递一个 iterable，作为生成器函数的输出。generator 通过 `next()` 或 `for` 恢复执行。

```py
def gen():
    yield from [1,2,3,4,5,6,7,8,9]
# for generator    
for i in gen():
    print(i)
# next generator    
generator = gen()
print(next(generator)) # 1
print(next(generator)) # 2
# ... or StopIteration
```

生成器对象方法：

```py
generator.__next__()   # 相当于 next(generator)，恢复上一次 yield 执行
generator.send(value)  # 恢复执行并向生成器（非 yield from）send value，首次激活只能传递 None
generator.throw(value)
generator.throw(type[, value[, traceback]])  # yield 暂停处引发异常，并返回下一个值
generator.close()      # 关闭生成器，后续取值引发 StopIteration
```

> 异步生成器

`async def ... yield` 异步生成器（非 `yield from`）对象通常在协程函数的 `async for` 中使用。async_generator 通过 `await __anext__()` 或 `async for` 恢复执行。

当一个异步生成器因 `break` 或调用方任务被取消、异常等导致生成器退出时，调用方应调用 `aclose()` 显式关闭生成器，以将协程从事件循环中分离从而避免资源泄漏。

```py
import asyncio
# 异步生成器
async def async_gen(start, end):
    try:
        for i in range(start, end + 1):
            await asyncio.sleep(0.05)
            yield i
    except GeneratorExit:
        print("GeneratorExit caught!")  # aclose() 触发
        raise
    finally:
        print("Generator finally block")


# 消费异步生成器
async def async_consumer(start, end):
    try:
        async_generator = async_gen(start, end)
        async for num in async_generator:
            if num > 10:
                raise ValueError("Something wrong")
            print(num)
    except ValueError:
        print("Caught ValueError")
    finally:
        await async_generator.aclose()  # 需要手动关闭
        print("finally")

asyncio.run(async_consumer(0, 100))
# Caught ValueError
# GeneratorExit caught!
# Generator finally block
# finally
```

异步生成器对象方法：

```py
async_generator.__next__()  
async_generator.send(value)
async_generator.throw(value)
async_generator.throw(type[, value[, traceback]])
async_generator.aclose()  
```


>---
#### 5.5. lambda 表达式

lambda 创建匿名函数，函数不包含语句和标注

```py
lambda [parameter_list] : expression
# 相当于
def <lambda>(parameter_list):
    return expression

funcs = [lambda x: x + i for i in range(3)]  # 惰性求值 i = 2
print(funcs[0](5))  # 7, funcs[0] = lambda x: x + 2
print(funcs[1](5))  # 7, funcs[1] = lambda x: x + 2

funcs = [lambda x, v=i: x + v for i in range(3)] # 利用闭包
print(funcs[0](5))  # 5, funcs[0] = lambda x: x + 0
print(funcs[1](5))  # 6, funcs[1] = lambda x: x + 1
```

---
### 6. 类

类定义的属性顺序保存在 `cls.__dict__` 中。

```python
[@decorator]+
class ClassName [typeParams] [inheritance]:
    [""" 类文档字符串 __doc__ """]
    class_attr = "class_shared"  # 类属性

    def __init__(self, name: str):
        self.name = name     # 实例属性
        self._protected = 0  # 约定：受保护
        self.__private = 0   # 私有：名称重整为 __ClassName__private
    def instance_method(self):
        return f"Hello, {self.name}"
    @classmethod
    def class_method(cls):
        return cls.class_attr
    @staticmethod
    def static_method():
        return "Static"
```

>***

#### 6.1. 抽象与继承

```python
from abc import ABC, abstractmethod
class Base(ABC): 
    def __init__(self):
        self.base_attr = "base_shared"
    @abstractmethod
    def method(self):
        pass

class Derived(Base):
    def __init__(self):
        super().__init__()  # 调用父类构造器
    def method(self):  # 实现抽象方法
        print(self.base_attr)

# base = Base()   # 抽象类无法实例化
derived = Derived()
derived.method()  # base_shared
```

> 多继承

```py
class A: pass
class B: pass
class C(A, B): pass
# 方法解析顺序（MRO）, 继承链
C.__mro__   # (<class 'C'>, <class 'A'>, <class 'B'>, <class 'object'>)
```

>---
#### 6.2. 属性控制：property

```py
class Circle:
    def __init__(self, radius: float, pos: tuple = (0.0, 0.0)):
        self.__radius = radius
        self.__pos = pos
    
    @property   # 读写可删除
    def Radius(self):   # c.Radius
        return self.__radius
    @Radius.setter 
    def Radius(self, value: float): # c.Radius = value
        self.__radius = value
    @Radius.deleter
    def Radius(self):   # del c.Radius
        del self.__radius

    @property  # 只读
    def Pos(self):
        return self.__pos

c = Circle(5.0, (1.0, 2.0))
c.Radius = 10.0
print(c.Radius)
print(c.Pos)
del c.Radius 
```


>***
#### 6.3. 元类

Python 的类也是对象，由元类创建。`type` 为元类的默认实现。元类由 `metaclass` 关键字参数指定。

```py
class Person:
    def __init__(self, name: str):
        self.name = name
# 等价于
Person = type("Person", (), {"__init__": lambda self, name: self.name = name})
# type('className', (base_classes), {'attrs':value})
```

> 用户定义元类

```py
class SingletonMeta(type):
    __instances = {}
    def __call__(cls, *args, **kwargs):
        if cls not in cls.__instances:
            cls.__instances[cls] = super().__call__(*args, **kwargs)
        return cls.__instances[cls]
    
class DataBase(metaclass=SingletonMeta):
    def __init__(self):
        print("DataBase initialized")

db1 = DataBase()
db2 = DataBase()  #单例模式
assert db1 is db2 
```

---
### 7. 泛型

函数、类、类型别名包含类型参数以构造泛型，类型参数由 `__type_params__` 列出。类型参数分为：
- `typing.TypeVar` 由 `T`  引入表示一个单独类型，可以包含一个范围（`T:int`）或约束（`T:(str, bytes)`）。
- `typing.TypeVarTuple` 由 `*Ts` 表示一个类型元组。
- `typing.ParamSpec` 由 `**P` 表示一个可调用对象的形参。

```py
def overly_generic[
    SimpleTypeVar,         # 无默认值，其 __default__ 预设为 typing.NoDefault
    TypeVarWithDefault = int,    # 默认值
    TypeVarWithBound: int,       # 类型范围
    TypeVarWithConstraints: (str, bytes),    # 类型约束
    *SimpleTypeVarTuple = (int, float),
    **SimpleParamSpec = (str, bytearray),
](
    a: SimpleTypeVar,
    b: TypeVarWithDefault,
    c: TypeVarWithBound,
    d: Callable[SimpleParamSpec, TypeVarWithConstraints],
    *e: SimpleTypeVarTuple,
): pass
```

> 泛型声明

```py
# 泛型函数
[@decorator]
def GeneFunc[T](func: Callable[[T], Any] = some_default) -> T | None: pass
# 泛型类
[@decorator]
class GeneCls[T]: pass
class GeneDerived(Base[T], arg=T): pass
# 泛型类型别名
type ListOrSet[T] = list[T] | set[T]
```

---
### 8. 特殊属性和方法名称

#### 8.1. 用户定义函数 

> 只读属性

|name|description|
|:---|:----------|
|`function.__builtins__`| 指向保存函数内置命名空间的字典引用|
|`function.__globals__`| 函数定义所在模块的全局命名空间|
|`function.__closure__`| None 或单元格元组，即访问闭包中捕获的自由变量|

单元格元组的单元具有 `cell_contents` 属性，用于设置或获取单元的值。

```py
def outer(x):
    y = 1
    def inner():
        return x + y
    return inner

closure_func = outer(5)
closure_func.__closure__[0].cell_contents = 12345
closure_func.__closure__[1].cell_contents = 54321
print(closure_func())  # 66666
```

> 可写属性

|name|description|
|:---|:----------|
|`function.__doc__`|文档字符串|
|`function.__name__`| 函数名称|
|`function.__qualname__`| 完全限定名称|
|`function.__module__`| 所在模块|
|`function.__defaults__`| 默认参数|
|`function.__code__`| 已编译函数体的代码对象|
|`function.__dict__`| 对象字典|
|`function.__annotations__`| 标注|
|`function.__annotate__`| 函数的格式化标注|
|`function.__kwdefaults__`| 关键词参数默认值字典|
|`function.__type_params__`| 泛型函数类型参数元组|

> 实例方法

|name|description|
|:---|:----------|
|`method.__self__`| 方法绑定类实例对象|
|`method.__func__`| 指向原始函数对象|
|`method.__doc__`| 文档字符串|
|`method.__name__`| 方法名称|
|`method.__module__`| 所在模块|

>---
#### 8.2. 模块对象属性

|name|description|
|:---|:----------|
|`module.__name__`| 模块名称|
|`module.__spec__`| 模块与导入系统相关联的状态的记录|
|`module.__spec__.loader`| 模块加载器|
|`module.__spec__.parent`| 模块父模块规范|
|`module.__path__`| 模块路径|
|`module.__file__`| 模块文件路径，可选|
|`module.__spec__.cached`| 模块已编译文件缓冲，可选|
|`module.__doc__`| 文档字符串|
|`module.__annotate__`| 模块的格式化标注|
|`module.__dict__`| 模块命名空间，只读|
|`module.__lazy_module__`| 懒惰导入模块|

>---
#### 8.3. 自定义类

|name|description|
|:---|:----------|
|`class.__name__`| 类名称|
|`class.__qualname__`| 完全限定名称|
|`type.__module__`| 所在模块|
|`type.__dict__`| 类命名空间，只读|
|`type.__bases__`| 基类元组|
|`type.__base__`| 继承链中负责实例内存布局的单独基类|
|`type.__doc__`| 文档字符串|
|`type.__annotations__`| 类变量标注字典|
|`type.__annotate__()`| 类的格式化标注|
|`type.__type_params__`| 泛型类类型参数元组|
|`type.__static_attributes__`| 在类语句体中 `self.X` 赋值的属性名称元组|
|`type.__firstlineno__`| 类定义开始的行号|
|`type.__mro__`| 类解析阶段的基类元组|
|`type.mro()`| 可由元类重写，为其实例定制方法解析顺序。类实例化时调用并保存到 `__mro__`|
|`type.__subclasses__()`| 保存一个由指向其直接子类的弱引用组成的列表|
|`object.__class__`| 类实例所属类|
|`object.__dict__`| 实例属性字典，可写|

>---
#### 8.4. object 基本定制

类可以定义特殊名称的方法，用于自定义对象或模块的特定操作。`None` 弃置特殊方法。

| name                                          | description                                                              |
| :-------------------------------------------- | :----------------------------------------------------------------------- |
| `object.__new__(cls[, ...]) -> instance`      | 创建实例。                                                               |
| `object.__init__(self[, ...])`                | 初始化实例。                                                             |
| `object.__del__(self)`                        | 类似于析构函数。                                                   |
| `object.__repr__(self) -> str`                | `repr()` 调用。用于调试和打印。具有默认实现。                            |
| `object.__str__(self) -> str`                 | `str(obj)`,`__format__()`,`print()` 调用。`object` 默认实现调用 `repr`。 |
| `object.__bytes__(self) -> bytes`             | `bytes(obj)` 调用。object 不提供。                                       |
| `object.__format__(self, format_spec) -> str` | 格式化的字符串表示。                                                     |
| `object.__hash__(self) -> int`                | `hash()` 调用。需要定义 `__eq__`。                                       |
| `object.__bool__(self) -> bool`               | `bool()` 调用。真值检测。                                                |
| `object.__lt__(self, other) -> bool`          | `self < other`                                                           |
| `object.__le__(self, other) -> bool`          | `self <= other`                                                          |
| `object.__eq__(self, other) -> bool`          | `self == other`                                                          |
| `object.__ne__(self, other) -> bool`          | `self != other`                                                          |
| `object.__gt__(self, other) -> bool`          | `self > other`                                                           |
| `object.__ge__(self, other) -> bool`          | `self >= other`                                                          |

>---
#### 8.5. 自定义属性访问

自定义对类实例属性或模块属性的 `x.name` 的访问。

| name                                          | description                                                                                       |
| :-------------------------------------------- | :------------------------------------------------------------------------------------------------ |
| `object.__getattr__(self, name)`        | 在默认属性访问失败时调用。object 不提供。                                                         |
| `object.__getattribute__(self, name)`   | `obj.name`                                                                                        |
| `object.__setattr__(self, name, value)` | `obj.name = value`                                                                                |
| `object.__delattr__(self, name)`        | `del obj.name`                                                                                    |
| `object.__dir__(self)`                  | `dir()` 调用，返回一个可迭代对象。                                                                |
| `module.__getattr__(name)`              | 对模块属性的访问，当 `obj.getattribute` 未找到属性时调用 `mod.getattr` 在 `mod.__dict__` 中查找。 |
| `module.__dir__()`                      | 返回一个表示模块中可访问名称的字符串可迭代对象                                                    |
| `module.__class__`                      | 为模块对象的 `__class__` 设置一个 `ModuleType` 子类。                                             |
| `object.__slots__`                            | 显式数据成员并禁止实例创建 `__dict__` 和 `__weakref__`。|

> slots 实例属性内存优化

`__slots__` 显式数据成员并禁止实例创建 `__dict__` 和 `__weakref__`（除非数据成员中包含）。`__slots__` 返回非字符串的 Iterable。使用字典时，key 用作变量，value 用作变量描述，可通过 `help()` 查看。

类属性不能通过 `__slots__` 定义的实例变量设置默认值。

```py
class Person:
    __slots__ = {
        'name': '用户姓名',
        'age': '用户年龄',
        'email': '电子邮箱'
    }
    def __init__(self, name, age, email):
        self.name = name
        self.age = age
        self.email = email

print(help(Person))
# |  Data descriptors defined here:
# |  age 用户年龄
# |  email 电子邮箱
# |  name 用户姓名
```

>---
#### 8.6. 实现描述器

描述器（包含`__get__`, `__set__`, `__delete__` ）是 Python 实现属性访问控制的机制，表现为具有绑定行为的对象属性。定义 `__set__` / `__delete__` 的描述器是一个数据描述器。

| name | description |
| :--- | :---------- |
|`object.__get__(self, instance, owner=None)`| 获取所有者类或类实例的属性。|
|`object.__set__(self, instance, value)`| 设置 instance 所有者类的实例属性。|
|`object.__delete__(self, instance)`| 删除 instance 所有者类的实例属性。|
|`object.__objclass__`| 可选，在 `inspect` 中解读为此对象定义所在的类。|

属性访问的默认行为是从一个对象的 `__dict__` 中查找属性。`a.x` 查找从 `a.__dict__['x']` 开始，然后 `type(a).__dict__['x']`，接下来查找 `type(a)` 的上级基类（不包括元类）。当找到的值是定义了某个描述器方法的对象，解释器则重载默认行为并转向描述器方法。

```py
class Range:
    def __init__(self, type, start, end):
        self.type = type
        self.start = start
        self.end = end
    def __set__(self, instance, value):
        if type(value) is not self.type:
            raise TypeError(f"Expected {self.type}, got {type(value)}")
        if not self.start <= value <= self.end:
            raise ValueError(f"Value {value} is out of range {self.start} to {self.end}")
        instance.__dict__[self] = value
    def __get__(self, instance, owner):
        return instance.__dict__[self]

class Byte:
    Value = Range(int, 0, 255)
    def __init__(self, value: int):
        self.Value = value
    def __str__(self):
        return f"{self.Value}"

b = Byte(100)
print(b)  # 100
b.Value = 10086  # out of range
```

>---
#### 8.7. 元类操作

默认类的创建是由 `type()` 元类构建的，可通过在类定义行添加 `metaclass` 关键字参数来定制元类。

| name | description |
| :--- | :---------- |
|`object.__mro_entries__(self, bases) -> tuple`| 返回类元组，用于确定基类。|
|`object.__prepare__(name, bases, **kwargs) -> namespace`| 准备类命名空间。|
|`type.__instancecheck__(self, instance) -> bool` | `isinstance()` 调用。 |
|`type.__subclasscheck__(self, subclass) -> bool`| `issubclass()` 调用。 |

一个类定义的执行步骤：
- 解析 MRO (`object.__mro_entries__(self, bases)`, 返回类元组) 以确定基类；
- 准备类命名空间 (`namespace = metaclass.__prepare__(name, bases, **kwargs)`) 并传入到 `__new__`；
- 执行类体 (`exec(body, globals(), namespace)`)；
- 返回新类对象 (`cls = metaclass(name, bases, namespace, **kwargs)`)。

```py
class Substitute:
    def __init__(self, target):
        self.target = target
    def __mro_entries__(self, bases):
        print(f"__mro_entries__ called, self={self}, bases={bases}")
        return (self.target,)
class Meta(type): pass
class Target(metaclass=Meta):
    pass
# 创建类时，Substitute(Target) 不是 type 实例
class MyClass(Substitute(Target)): pass

print(MyClass.__bases__)  # (<class 'Target'>,)
print(MyClass.__mro__)    # MyClass, Target, object
print(type(MyClass))      # <class 'Meta'>
```

`__instancecheck__`, `__subclasscheck__` 查找基于类的类型（元类），不能作为类方法在实际的类中定义。

```py
class EnumMeta(type):
    def __subclasscheck__(cls, subclass):
        return type(cls) is type(subclass)
class Enum(metaclass=EnumMeta) : pass

class Color(Enum):
    RED = 1
    GREEN = 2
    BLUE = 3

assert issubclass(Color, Enum)
```

>---
#### 8.8. 定制类行为

| name | description |
| :--- | :---------- |
|`object.__init_subclass__(cls, **kwargs)`| 派生子类时调用。|
|`object.__set_name__(self, owner, name)`| 在所有者类 owner 被创建时自动调用。|
| `object.__call__(self[, args...]) -> any`            |`obj(args)` |
| `object.__class_getitem__(cls, key) -> type`  | 返回一个泛型类 `cls[key]` 的专门化。key 是一个类。 |
|`object.__match_args__`|`match` 中的位置参数匹配。 |

> 子类初始化: `__init_subclass__`

当一个类继承另一个类时，会在这个父类上调用 `__init_subclass__()` 方法。

```py
class Person:
    def __init_subclass__(cls, /, default_name, **kwargs):
        super().__init_subclass__(**kwargs)
        cls.default_name = default_name
class Student(Person, default_name="Alice"): pass

print(Student.default_name)  # Alice
```

> 类名自动设置: `__set_name__`

当一个类被定义时，`type.__new__()` 会扫描类变量并对其中带有 `__set_name__()` 钩子的对象执行回调。类定义被执行之后的类变量不会自动调用 `__set_name__()`。

```py
# 类型检查描述器
class Typed:
    def __init__(self, type):
        self.type = type
    def __set_name__(self, owner, name):
        self.__name__ = name       # 自动设置属性名
    def __get__(self, instance, owner):
        if instance is None:
            return self
        return instance.__dict__.get(self.__name__)
    def __set__(self, instance, value):
        if not isinstance(value, self.type):
            raise TypeError(f"Expected {self.type.__name__}, got {type(value).__name__}")
        instance.__dict__[self.__name__] = value

class Person:
    name = Typed(str)  # 转向描述器
    age = Typed(int)
    def __init__(self, name: str, age: int):
        self.name = name # name.__set__(self, value)
        self.age = age

p = Person("Alice", 30) # ok
p = Person(30, "Bob")   # TypeError
```

> 模拟泛型类型: `__class_getitem__`

```py
class collection(type):
    def __new__(mcls, name, bases, namespace):
        cls = super().__new__(mcls, name, bases, namespace)
        def __class_getitem__(cls, key):
            # 返回带 _type 属性的动态类
            return type(f"{cls.__name__}[{key.__name__}]", (cls,), {'_type': key})
        # 给类绑定 __class_getitem__ (classmethod)
        cls.__class_getitem__ = classmethod(__class_getitem__)
        return cls

class List(metaclass=collection):
    def __init__(self, *items):
        for item in items :
            if not isinstance(item, self._type):
                raise TypeError(f"Except {self._type.__name__}, but {type(item).__name__}") 
        self._data = list(items)
    def __repr__(self):
        return f"List({self._data})"
    def append(self, item):
        if not isinstance(item, self._type):
            raise TypeError(f"Except {self._type.__name__}, but {type(item).__name__}") 
        self._data.append(item)

l = List[int](1, 2, 3)           
l.append(4)   # List([1, 2, 3, 4])

l = List[int](1, 2, 3, "asd") # TypeError
```

> 匹配参数: `__match_args__`

```py
class Person:
    __match_args__ = ("name", "age")
    def __init__(self, name: str, age: int):
        self.name = name
        self.age = age

def match_person(p : Person):
    match p:
        case Person("Alice", age):
            print(f"Alice is {age} years old")
        case _:
            print("Unknown person")
```

>---
#### 8.9. 模拟容器类型

|name|description|
|:---|:----------|
|`object.__len__(self)->int`| `len(obj)`|
|`object.__length_hint__(self) -> int`| `operator.length_hint(obj)`|
|`object.__getitem__(self, subscript) -> any`| `self[subscript]`|
|`object.__setitem__(self, key, value)`| `self[key] = value`|
|`object.__delitem__(self, key)`| `del self[key]`|
|`object.__missing__(self, key) -> any`| `self[key]` 未找到时调用。|
|`object.__iter__(self) -> iterator`| 迭代器对象|
|`object.__reversed__(self) -> iterator`| 反向迭代器对象|
|`object.__contains__(self, item) -> bool`| `item in self`|

>---
#### 8.10. 模拟数字类型

> self op other

|name|description|
|:---|:----------|
|`object.__add__(self, other)`| `self + other`|
|`object.__sub__(self, other)`| `self - other`|
|`object.__mul__(self, other)`| `self * other`|
|`object.__matmul__(self, other)`| `self @ other`|
|`object.__truediv__(self, other)`| `self / other`|
|`object.__floordiv__(self, other)`| `self // other`|
|`object.__mod__(self, other)`| `self % other`|
|`object.__divmod__(self, other)`| `divmod(self, other)`|
|`object.__pow__(self, other)`| `self ** other`|
|`object.__pow__(self, other, modulo)`| `pow(self, other, modulo)`|
|`object.__lshift__(self, other)`| `self << other`|
|`object.__rshift__(self, other)`| `self >> other`|
|`object.__and__(self, other)`| `self & other`|
|`object.__xor__(self, other)`| `self ^ other`|
|`object.__or__(self, other)`| `self \| other`|

> other op self

|name|description|
|:---|:----------|
|`object.__radd__(self, other)`| `other + self`|
|`object.__rsub__(self, other)`| `other - self`|
|`object.__rmul__(self, other)`| `other * self`|
|`object.__rmatmul__(self, other)`| `other @ self`|
|`object.__rtruediv__(self, other)`| `other / self`|
|`object.__rfloordiv__(self, other)`| `other // self`|
|`object.__rmod__(self, other)`| `other % self`|
|`object.__rdivmod__(self, other)`| `divmod(other, self)`|
|`object.__rpow__(self, other)`| `other ** self`|
|`object.__rpow__(self, other, modulo)`| `pow(other, self, modulo)`|
|`object.__rlshift__(self, other)`| `other << self`|
|`object.__rrshift__(self, other)`| `other >> self`|
|`object.__rand__(self, other)`| `other & self`|
|`object.__rxor__(self, other)`| `other ^ self`|
|`object.__ror__(self, other)`| `other \| self`|

> self op= other

|name|description|
|:---|:----------|
|`object.__iadd__(self, other)`| `self += other`|
|`object.__isub__(self, other)`| `self -= other`|
|`object.__imul__(self, other)`| `self *= other`|
|`object.__imatmul__(self, other)`| `self @ other`|
|`object.__itruediv__(self, other)`| `self /= other`|
|`object.__ifloordiv__(self, other)`| `self //= other`|
|`object.__imod__(self, other)`| `self %= other`|
|`object.__ipow__(self, other)`| `self **= other`|
|`object.__ipow__(self, other, modulo)`| `self.__ipow__(other, modulo)`|
|`object.__ilshift__(self, other)`| `self <<= other`|
|`object.__irshift__(self, other)`| `self >>= other`|
|`object.__iand__(self, other)`| `self &= other`|
|`object.__ixor__(self, other)`| `self ^= other`|
|`object.__ior__(self, other)`| `self \|= other`|

> op self

|name|description|
|:---|:----------|
|`object.__neg__(self)`| `-self`|
|`object.__pos__(self)`| `+self`|
|`object.__abs__(self)`| `abs(self)`|
|`object.__invert__(self)`| `~self`|
|`object.__complex__(self)`| `complex(self)`|
|`object.__int__(self)`| `int(self)`|
|`object.__float__(self)`| `float(self)`|
|`object.__index__(self)`| `operator.index(self)`|
|`object.__round__(self[, ndigits])`| `round(self[, ndigits])`|
|`object.__trunc__(self)`| `math.trunc(self)`|
|`object.__floor__(self)`| `math.floor(self)`|
|`object.__ceil__(self)`| `math.ceil(self)`|

>---
#### 8.11. 模拟缓冲区类型

|name|description|
|:---|:----------|
|`object.__buffer__(self, flags)`| 当从 `self` 请求一个缓冲区时调用, `flags` 查阅 ` inspect.BufferFlags`|
|`object.__release_buffer__(self, buffer)`| 释放先前从 `self` 请求的缓冲区对象 |

```py
class SimpleBuffer:
    def __init__(self, size: int):
        self._size = size
    def __buffer__(self, flags: int) -> memoryview:
        self._data = bytearray(self._size)
        return memoryview(self._data)
    def __release_buffer__(self, buffer: memoryview):
        del self._data
        print(f"release buffer: {buffer.nbytes} bytes")

view = memoryview(SimpleBuffer(10)) # 获取缓冲区视图
view[0] = 65  # 写入 'A'
print(bytes(view[:5]))  # b'A\x00\x00\x00\x00'
del view  # release buffer: 10 bytes
```


>---
#### 8.12. 特殊协议

> 可等待对象

|name|description|
|:---|:----------|
|`object.__await__(self) -> Awaitable`| `await self` 返回一个可等待对象。|

> 迭代器

|name|description|
|:---|:----------|
|`object.__iter__(self) -> iterator`| 迭代器对象|
|`object.__next__(self) -> any`| 迭代器的下一个结果|
|`object.__aiter__(self) -> AsyncIterator`| 异步迭代器对象|
|`object.__anext__(self) -> any`| 异步迭代器的下一个结果|

> 上下文管理器

|name|description|
|:---|:----------|
|`object.__enter__(self)`| 进入 `with`，返回上下文管理器 |
|`object.__exit__(self, exc_type, exc_val, exc_tb)`| 退出 `with` 语句 |
|`object.__aenter__(self) -> awaitable`| 进入 `async with`，返回异步上下文管理器 |
|`object.__aexit__(self, exc_type, exc_val, exc_tb) -> awaitable`| 退出 `async with` 语句 |


>---
#### 8.13. 标注

|name|description|
|:---|:----------|
|`object.__annotations__`| 函数、类和模块的注解 |
|`object.__annotate__(format)`| 返回注解的格式化字典表示，`format` 是 `annotationlib.Format` 枚举成员 |

```py
from annotationlib import Format,get_annotations

class Example:
    a: int
    b: list[str]

# 获取完整求值的注解
print(get_annotations(Example,format=Format.VALUE))
# {'a': <class 'int'>, 'b': list[str]}
print(get_annotations(Example,format=Format.STRING))
# {'a': 'int', 'b': 'list[str]']}
```


>---

### 9. 异常处理

所有异常派生自 `BaseException`，用户定义异常从非致命异常 `Exception` 派生。`raise` 引发异常。异常被引发时会自动创建一个回溯对象关联到 `__traceback__` 属性，可以通过 `with_traceback(obj)` 设置用户回溯对象。

```py
# 引发异常
def raise_exception():
    raise [expr [from expr]]
    raise   # 在 except 子句重新引发当前处理的异常
    raise ValueError   # raise  ValueError() 的简化
    raise Exception("exception occurred").with_traceback(tracebackObj)   # 设置用户回溯对象
    raise Exception("exception occurred")  from other_exc  # 异常串联

def exception_handling():
    try:
        stmt 
    except ValueError as e:
        e.add_note("Add some information")      # 添加异常注释
    except RuntimeError, TypeError, NameError:  # 捕捉多个异常
        pass
    except SomeException: 
        tb = sys.exception().__traceback__
        raise OtherException().with_traceback(tb)  # 设置用户回溯对象
    except : # 捕捉任意异常
        raise    # 重新引发当前异常
    else : # 无异常时执行
        print("No exception")
    finally : 
        print("Do some cleanup")
```

`from` 用于异常串联，表明一个异常是另一个异常的直接后果。串联机制也会隐式发生，可以显式 `from None` 抑制串联：

```py
try:
    print(1 / 0)
except :
    raise RuntimeError("Something bad happened") # 隐式串联
# ZeroDivisionError: division by zero
# RuntimeError: Something bad happened  
    raise RuntimeError("Something bad happened") from None   # 抑制串联，仅抛出
# RuntimeError: Something bad happened 
```

`ExceptionGroup` 用于异常打包，`except*` 从分组中提取子异常，其他异常传播到其他子句。

```py
try:
    raise ExceptionGroup( "group 1", [ OSError(1), SystemError(2),
                ExceptionGroup( "group 2", [ OSError(3), RecursionError(4) ])])
except* OSError:
    print("There were OSErrors")
except* SystemError:
    print("There were SystemErrors")
'''
There were OSErrors
There were SystemErrors
  + Exception Group Traceback (most recent call last):
  |     raise ExceptionGroup( "group 1", [ OSError(1), SystemError(2),
  |                 ExceptionGroup( "group 2", [ OSError(3), RecursionError(4) ])])
  | ExceptionGroup: group 1 (1 sub-exception)
  +-+---------------- 1 ----------------
    | ExceptionGroup: group 2 (1 sub-exception)
    +-+---------------- 1 ----------------
      | RecursionError: 4
      +------------------------------------
'''
```

> 内置异常层次结构

```py
BaseException
├── BaseExceptionGroup
├── GeneratorExit
├── KeyboardInterrupt
├── SystemExit
└── Exception
    ├── ArithmeticError
    │   ├── FloatingPointError
    │   ├── OverflowError
    │   └── ZeroDivisionError
    ├── AssertionError
    ├── AttributeError
    ├── BufferError
    ├── EOFError
    ├── ExceptionGroup [BaseExceptionGroup]
    ├── ImportError
    │   └── ModuleNotFoundError
    ├── LookupError
    │   ├── IndexError
    │   └── KeyError
    ├── MemoryError
    ├── NameError
    │   └── UnboundLocalError
    ├── OSError   # EnvironmentError,IOError,WindowsError 别名
    │   ├── BlockingIOError
    │   ├── ChildProcessError
    │   ├── ConnectionError
    │   │   ├── BrokenPipeError
    │   │   ├── ConnectionAbortedError
    │   │   ├── ConnectionRefusedError
    │   │   └── ConnectionResetError
    │   ├── FileExistsError
    │   ├── FileNotFoundError
    │   ├── InterruptedError
    │   ├── IsADirectoryError
    │   ├── NotADirectoryError
    │   ├── PermissionError
    │   ├── ProcessLookupError
    │   └── TimeoutError
    ├── ReferenceError
    ├── RuntimeError
    │   ├── NotImplementedError
    │   ├── RecursionError
    │   └── PythonFinalizationError
    ├── StopAsyncIteration
    ├── StopIteration
    ├── SyntaxError
    │   └── IndentationError
    │       └── TabError
    ├── SystemError
    ├── TypeError
    ├── ValueError
    │   └── UnicodeError
    │       ├── UnicodeTranslateError
    │       ├── UnicodeDecodeError
    │       └── UnicodeEncodeError
    └── ...
```

> 自定义异常

```python
class MyError(Exception):
    def __init__(self, message, code=0):
        self.message = message
        self.code = code
        super().__init__(message)

raise MyError("Something went wrong", code=500)
```

>---
#### 9.1. 警告

`Warning` 派生自 `Exception`。调用 `warnings.warn()` 引发警告。

```python
import warnings

warnings.warn("deprecated", DeprecationWarning)
warnings.warn("deprecated", DeprecationWarning, stacklevel=2)
```

> **内置警告层次结构**

```py
Exception
└── Warning
    ├── BytesWarning        
    ├── DeprecationWarning  
    ├── EncodingWarning     
    ├── FutureWarning       
    ├── ImportWarning      # 导入模块时触发
    ├── PendingDeprecationWarning  # 即将废弃
    ├── ResourceWarning    # 资源使用相关
    ├── RuntimeWarning     # 可疑运行时特性
    ├── SyntaxWarning      # 疑问语法特性
    ├── UnicodeWarning     # Unicode 相关
    └── UserWarning       # 默认类别
```

> **警告过滤器**

警告过滤器控制是否忽略、显示或转为错误。可通过 `filterwarnings()` 或命令行 `-W` 选项配置。

| action      | 处置                         |
| ----------- | ---------------------------- |
| `"default"` | 为每个位置打印第一次匹配警告 |
| `"error"`   | 将匹配警告转换为异常         |
| `"ignore"`  | 从不打印匹配的警告           |
| `"always"`  | 总是打印匹配的警告           |
| `"module"`  | 为每个模块打印第一次匹配警告 |
| `"once"`    | 仅打印第一次出现的匹配警告   |

```python
import warnings

# 添加过滤器
warnings.filterwarnings("ignore", category=DeprecationWarning)
warnings.filterwarnings("error", message=".*deprecated.*", category=DeprecationWarning)

# 简单过滤器
warnings.simplefilter("ignore")
warnings.simplefilter("error", category=ResourceWarning)

# 重置过滤器
warnings.resetwarnings()
```

> **命令行选项**

```bash
python -W all                     # 显示所有警告
python -W ignore                  # 忽略所有警告
python -W error                   # 将所有警告转为错误
python -W error::DeprecationWarning    # DeprecationWarning 转为错误
python -W default::DeprecationWarning  # 显示 DeprecationWarning
```

> **catch_warnings 上下文管理器**

```python
import warnings

# 暂时抑制警告
with warnings.catch_warnings():
    warnings.simplefilter("ignore")
    deprecated_function()  # 不会显示警告

# 捕获警告以进行测试
with warnings.catch_warnings(record=True) as w:
    warnings.simplefilter("always")
    warnings.warn("deprecated", DeprecationWarning)
    assert len(w) == 1
    assert issubclass(w[-1].category, DeprecationWarning)
```

> **showwarning() 与 formatwarning()**

```python
import warnings

# 自定义警告显示
warnings.showwarning("message", UserWarning, "file.py", 10)

# 格式化警告信息
msg = warnings.formatwarning("message", UserWarning, "file.py", 10)
```

---

### 10. 异步编程

#### 10.1. 协程：async, await

协程（Coroutine）是 Python 中实现协作式多任务（异步）的核心机制，依靠事件循环调度，适用于 IO 密集型任务。`async def` 定义协程函数，调用时返回一个 coroutine 对象，属于 awaitable 对象。协程的执行通过调用 `__await__()` 并迭代其结果来控制。

一个协程函数可能包含 `await`, `async for`, `async with` 等语句。

```py
async def other_coro():
    print("do other coro")
    await asyncio.sleep(1)  # 模拟长时间操作
    print("other_coro complete")
    return "RESULT"

async def coro():
    print("do coro")
    task = asyncio.create_task(other_coro())  # 挂起 coro，创建一个异步任务
    print("back to coro")
    return await task  # 等待 other_coro 完成

print(asyncio.run(coro())) # RESULT
```

协程对象方法：

```py
coroutine.send(value)  # 开始或恢复协程的执行
coroutine.throw(value)
coroutine.throw(type[, value[, traceback]])  
coroutine.close()      
```

>---
#### 10.2. 可等待对象：awaitable

`await` 的操作数是 awaitable 对象，主要三种类型：
- Coroutine 可由其他 coroutine 中 `await`。
- Task 由 `asyncio.create_task()` 封装 coroutine 创建，自动调度执行。
- Future 表示一个异步操作的结果，通常由 `asyncio` 等库函数暴露。

> 协程对象

```py
class Coroutine:
    def __init__(self, action: callable):
        self.__action = action
        self.__gen = None
    def __await__(self):
        self.__gen = self.__generator()
        return self.__gen  # 返回一个可迭代对象
    def __generator(self):
        try: 
            yield None  
            if self.__action is not None:
                self.__action()
        except GeneratorExit:
            self.__closed = True
            print("Coroutine closed")
            raise
        except Exception as e:
            print(f"Coroutine error: {e}")
            raise

async def custom_coroutine():
    await Coroutine(lambda: print("hello", end=" "))
    await Coroutine(lambda: print("world")) 

asyncio.run(custom_coroutine())
```


>---
#### 10.3. 事件循环与并发

事件循环在单线程中按协作式多任务的方式调度协程。每个时刻只运行一个任务，当该任务 `await` 时，事件循环切换到另一个就绪的任务。
- `asyncio.run()` 是运行异步的高层级入口，用于创建一个新的事件循环，每个线程只能有一个事件循环。
- `asyncio.create_task()` 封装协程为 Task 并调度执行。Task 用于并发执行协程。
- `asyncio.gather()` 并发执行多个任务，等待所有任务完成。
- `asyncio.TaskGroup()` 用于管理多个任务，自动取消未完成的任务。
- `asyncio.wait()` 等待多个任务完成，返回完成的任务结果。
- `asyncio.as_completed()` 返回一个可迭代对象，按任务完成顺序返回完成的任务。

```python
asyncio.run(coro, *, debug=None, loop_factory=None) -> _Result
asyncio.create_task(coro, *, name=None, context=None, eager_start=None) -> task
asyncio.gather(*aws, return_exceptions=False) -> awaitable
asyncio.TaskGroup()
asyncio.wait(aws, *, timeout=None, return_when=ALL_COMPLETED) -> (done, pending)
asyncio.as_completed(aws, *, timeout=None) -> iterator
```
```py
import asyncio

async def coro(n):
    print(f"coro {n} start")
    await asyncio.sleep(1)
    print(f"coro {n} complete")
    return n
# gather
async def gather():
    print("gather 1,2 start")
    results = await asyncio.gather(coro(1), coro(2))
    for result in results:
        print(f"got {result}")
    print("gather 1,2 complete")
# wait
async def wait():
    print("wait 3,4 start")
    tasks = [
        asyncio.create_task(coro(3)),
        asyncio.create_task(coro(4))
    ]
    done, pending = await asyncio.wait(tasks, return_when=asyncio.ALL_COMPLETED)
    for task in done:
        print(f"got {task.result()}")
    print("wait 3,4 complete")
# as_completed
async def as_completed():
    print("as_completed 3,4 start")
    tasks = [
        asyncio.create_task(coro(5)),
        asyncio.create_task(coro(6))
    ]
    async for completed_task in asyncio.as_completed(tasks):
        result = await completed_task
        print(f"got {result}")
    print("as_completed 5,6 complete")
# TaskGroup
async def taskgroup():
    async with asyncio.TaskGroup() as tg:
        tg.create_task(gather())
        tg.create_task(wait())
        tg.create_task(as_completed())
    # async with 隐式等待所有任务完成
    print("taskgroup complete")

asyncio.run(taskgroup())
```

---

### 11. 附录

#### 11.1. 内置函数

| 分类        | 函数                                                                                                                                                                                                                    |
| :---------- | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 数学运算    | `abs()`、`complex()`、`divmod()`、`max()`、`min()`、`pow()`、`round()`、`sum()`                                                                                                                                         |
| 类型转换    | `bin()`、`bool()`、`bytearray()`、`bytes()`、`chr()`、`dict()`、`float()`、`frozenset()`、`hex()`、`int()`、`list()`、`oct()`、`ord()`、`set()`、`str()`、`tuple()`                                                     |
| 容器操作    | `all()`、`any()`、`enumerate()`、`filter()`、`len()`、`map()`、`range()`、`reversed()`、`slice()`、`sorted()`、`zip()`                                                                                                  |
| 对象操作    | `callable()`、`classmethod()`、`delattr()`、`dir()`、`getattr()`、`hasattr()`、`hash()`、`id()`、`isinstance()`、`issubclass()`、`object()`、`property()`、`setattr()`、`staticmethod()`、`super()`、`type()`、`vars()` |
| 输入 / 输出 | `ascii()`、`format()`、`input()`、`open()`、`print()`、`repr()`                                                                                                                                                         |
| 执行 / 编译 | `__import__()`、`compile()`、`eval()`、`exec()`                                                                                                                                                                         |
| 迭代器相关  | `aiter()`、`anext()`、`iter()`、`next()`                                                                                                                                                                                |
| 其他函数    | `breakpoint()`、`globals()`、`help()`、`locals()`、`memoryview()`、`sentinel()`                                                                                                                                                      |
---