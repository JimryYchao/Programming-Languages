## Python 内置函数

---
### 1. 数学运算

| 函数 | 说明 | 示例 |
|:------|:------|:------|
| `abs(x)` | 返回绝对值 | `abs(-5)` → `5` |
| `complex(real[,imag])` | 创建复数 | `complex(1,2)` → `(1+2j)` |
| `divmod(a,b)` | 返回(商,余数)元组 | `divmod(10,3)` → `(3,1)` |
| `max(iterable)` | 最大值 | `max([1,5,2])` → `5` |
| `min(iterable)` | 最小值 | `min([1,5,2])` → `1` |
| `pow(x,y[,z])` | 幂运算，可选取模 | `pow(2,10)` → `1024`<br>`pow(2,10,100)` → `24` |
| `round(x[,n])` | 四舍五入到n位小数 | `round(3.14159,2)` → `3.14` |
| `sum(iterable[,start])` | 求和 | `sum([1,2,3])` → `6` |

>---
### 2. 类型转换

| 函数 | 说明 | 示例 |
|:------|:------|:------|
| `bin(x)` | 转二进制字符串 | `bin(10)` → `'0b1010'` |
| `bool(x)` | 转布尔值 | `bool(0)` → `False` |
| `bytearray(x)` | 转可变字节 | `bytearray(5)` |
| `bytes(x)` | 转不可变字节 | `bytes("hi","utf-8")` |
| `chr(i)` | ASCII码转字符 | `chr(65)` → `'A'` |
| `dict(iterable)` | 转字典 | `dict([('a',1)])` → `{'a':1}` |
| `float(x)` | 转浮点数 | `float("3.14")` → `3.14` |
| `frozenset(iterable)` | 转不可变集合 | `frozenset([1,2,3])` |
| `hex(x)` | 转十六进制字符串 | `hex(255)` → `'0xff'` |
| `int(x)` | 转整数 | `int("123")` → `123` |
| `list(iterable)` | 转列表 | `list("abc")` → `['a','b','c']` |
| `oct(x)` | 转八进制字符串 | `oct(10)` → `'0o12'` |
| `ord(c)` | 字符转ASCII码 | `ord('A')` → `65` |
| `set(iterable)` | 转集合 | `set([1,2,2,3])` → `{1,2,3}` |
| `str(x)` | 转字符串 | `str(123)` → `"123"` |
| `tuple(iterable)` | 转元组 | `tuple([1,2,3])` → `(1,2,3)` |

>---
### 3. 容器操作

| 函数 | 说明 | 示例 |
|:------|:------|:------|   
| `all(iterable)` | 全为真返回True | `all([True,False])` → `False` |
| `any(iterable)` | 任一为真返回True | `any([False,True])` → `True` |
| `enumerate(iterable[,start])` | 枚举(索引,值)对 | `list(enumerate(['a','b']))` → `[(0,'a'),(1,'b')]` |
| `filter(func, iterable)` | 过滤序列 | `list(filter(lambda x:x>2,[1,2,3,4]))` → `[3,4]` |
| `len(s)` | 返回长度 | `len([1,2,3])` → `3` |
| `map(func, *iterables)` | 映射函数到序列 | `list(map(str,[1,2]))` → `['1','2']` |
| `range([start,]stop[,step])` | 生成整数序列 | `list(range(3))` → `[0,1,2]` |
| `reversed(seq)` | 反向迭代器 | `list(reversed([1,2,3]))` → `[3,2,1]` |
| `slice(start,stop[,step])` | 切片对象 | `s=slice(1,5,2)` 用于切片 |
| `sorted(iterable[,key][,reverse])` | 返回新排序列表 | `sorted([3,1,2])` → `[1,2,3]` |
| `zip(*iterables)` | 打包多个序列 | `list(zip([1,2],['a','b']))` → `[(1,'a'),(2,'b')]` |

>---
### 4. 对象操作

| 函数 | 说明 | 示例 |
|:------|:------|:------|
| `callable(obj)` | 判断是否可调用 | `callable(len)` → `True` |
| `classmethod(func)` | 类方法装饰器 | 类中使用 |
| `delattr(obj, name)` | 删除属性 | `delattr(obj,'name')` |
| `dir([obj])` | 返回属性列表 | `dir([])` 列表方法 |
| `getattr(obj, name[,default])` | 获取属性 | `getattr(obj,'name','default')` |
| `hasattr(obj, name)` | 检查属性 | `hasattr(obj,'name')` → `True/False` |
| `hash(obj)` | 返回哈希值 | `hash("hello")` |
| `id(obj)` | 返回对象唯一ID | `id("hello")` |
| `isinstance(obj, classinfo)` | 判断类型 | `isinstance(5, int)` → `True` |
| `issubclass(cls, classinfo)` | 判断子类 | `issubclass(bool, int)` → `True` |
| `object()` | 空对象基类 | `obj = object()` |
| `property([fget[,fset[,fdel[,doc]]]])` | 属性装饰器 | 类中使用 |
| `setattr(obj, name, value)` | 设置属性 | `setattr(obj,'name','value')` |
| `staticmethod(func)` | 静态方法装饰器 | 类中使用 |
| `super([type[,object-or-type]])` | 调用父类方法 | `super().__init__()` |
| `type(obj)` | 返回对象类型 | `type(123)` → `<class 'int'>` |
| `vars([obj])` | 返回`__dict__`属性 | `vars()` 局部变量字典 |

>---
### 5. 输入 / 输出

| 函数 | 说明 | 示例 |
|:------|:------|:------|
| `ascii(obj)` | ASCII表示 | `ascii('中文')` → `"'\\u4e2d\\u6587'"` |
| `format(value[,format_spec])` | 格式化值 | `format(3.14,'.1f')` → `'3.1'` |
| `input([prompt])` | 获取用户输入 | `name = input("Name: ")` |
| `open(file, mode='r')` | 打开文件 | `f = open('file.txt','r')` |
| `print(*objects, sep=' ', end='\n')` | 打印输出 | `print("Hello")` |
| `repr(obj)` | 官方字符串表示 | `repr("Hi")` → `"'Hi'"` |

>---
### 6. 执行 / 编译

| 函数 | 说明 | 示例 |
|:------|:------|:------|
| `__import__(name)` | 动态导入模块 | `math = __import__('math')` |
| `compile(source, filename, mode)` | 编译代码为对象 | `code = compile("2+3","","eval")` |
| `eval(expression[,globals[,locals]])` | 计算表达式 | `eval("2+3")` → `5` |
| `exec(object[,globals[,locals]])` | 执行代码 | `exec("a=5")` |


>---
### 7. 迭代器相关

| 函数 | 说明 | 示例 |
|:------|:------|:------|
| `aiter(async_iterable)` | 返回异步迭代器（3.10+） | `ait = aiter(async_gen)` |
| `anext(async_iterator[,default])` | 异步取下一个值（3.10+） | `await anext(ait)` |
| `iter(iterable)` | 返回迭代器 | `it = iter([1,2,3])` |
| `next(iterator[,default])` | 取下一个值 | `next(it)` → `1` |

>---
### 8. 其他函数

| 函数 | 说明 | 示例 |
|:------|:------|:------|
| `breakpoint()` | 调试器入口点 | `breakpoint()` 进入调试器 |
| `globals()` | 返回全局变量字典 | `globals()` 全局作用域 |
| `help([object])` | 获取帮助信息 | `help(print)` 交互式帮助 |
| `locals()` | 返回局部变量字典 | `locals()` 当前作用域 |
| `memoryview(obj)` | 内存视图 | `mv = memoryview(b'hello')` |
| `sentinel(name)` | 哨兵对象(3.15+) | `Singleton = sentinel("Singleton")` |

---