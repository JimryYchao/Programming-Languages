<a id="TOP"></a>

## Package unsafe

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/unsafe_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<!-- <a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	 -->
	<a href="https://pkg.go.dev/unsafe"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `unsafe` 包含绕过 Go 程序类型安全的操作。

---
### Pointer

`Pointer` 表示指向任何类型的指针，指针允许程序绕过类型系统读写任意内存。
- 任何类型的指针值都可以转换 `Pointer`
- `Pointer` 可以转换为任何类型的指针值
- `uintptr` 类型可以转换为 `Pointer` 类型
- `Pointer` 可以转换为 `uintptr` 类型

假设 `T2` 不大于 `T1`，并且两者共享等效的存储器布局，则这种转换允许将一种类型的数据重新解释为另一种类型的数据。一个例子是 `math.Float64bits` 的实现：

```go
func Float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}
```

将 `Pointer` 转换为 `uintptr` 会产生所指向的值的内存地址，作为整数，通常用来打印。`uintptr` 回转为 `Pointer` 通常是无效的。无法保证 `uintptr` 代表地址位置的值是否被回收或者被移动。

可以利用 `uintptr` 进行指针算术并回转为 `Pointer`，例如回转时添加一个偏移量实现一个 `Pointer` 的前进，最常见的用法是访问结构中的字段或数组的元素：

```go
p = unsafe.Pointer(uintptr(p) + offset)

// equivalent to f := unsafe.Pointer(&s.f)
f := unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Offsetof(s.f))

// equivalent to e := unsafe.Pointer(&x[i])
e := unsafe.Pointer(uintptr(unsafe.Pointer(&x[0])) + i*unsafe.Sizeof(x[0]))
```

将 `Pointer` 前进到其原始分配的末尾是无效的，它必须指向一个已分配的对象

```go
// INVALID: end points outside allocated space.
var s thing
end = unsafe.Pointer(uintptr(unsafe.Pointer(&s)) + unsafe.Sizeof(s))

// INVALID: end points outside allocated space.
b := make([]byte, n)
end = unsafe.Pointer(uintptr(unsafe.Pointer(&b[0])) + uintptr(n))
```

调用 `syscall.Syscall` 时将指针转换为 `uintptr`。包 `syscall` 中的 `Syscall` 函数将它们的 `uintptr` 参数直接传递给操作系统，然后操作系统可能会根据调用的细节将其中一些参数重新解释为指针。

```go
yscall.Syscall(SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(p)), uintptr(n))
```

`reflect` 的 `Pointer` 和 `UnsafeAddr.Value` 方法返回类型 `uintptr` 而不是 `unsafe.Pointer` 用于防止调用者在未首先导入 “`unsafe`” 的情况下将结果更改为任意类型。必须在调用后立即转换为指针：

```go
p := (*int)(unsafe.Pointer(reflect.ValueOf(new(int)).Pointer()))
```

在转换之前存储结果是无效的：

```go
// INVALID: uintptr cannot be stored in variable
// before implicit conversion back to Pointer during system call.
u := uintptr(unsafe.Pointer(p))
syscall.Syscall(SYS_READ, uintptr(fd), u, uintptr(n))

u := reflect.ValueOf(new(int)).Pointer()
p := (*int)(unsafe.Pointer(u))
```

`reflect` 的数据结构 `SliceHeader` 和 `StringHeader` 将字段 `Data` 声明为 `uintptr`，以防止调用者在没有首先导入 “`unsafe`” 的情况下将结果更改为任意类型。这意味着 `SliceHeader` 和 `StringHeader` 仅在解释实际切片或字符串值的内容时有效：

```go
var s string
hdr := (*reflect.StringHeader)(unsafe.Pointer(&s)) 
hdr.Data = uintptr(unsafe.Pointer(p))              
hdr.Len = n

// INVALID: a directly-declared header will not hold Data as a reference.
var hdr reflect.StringHeader
hdr.Data = uintptr(unsafe.Pointer(p))
hdr.Len = n
s := *(*string)(unsafe.Pointer(&hdr)) // p possibly already lost 
```

---
<a id="exam" ><a>