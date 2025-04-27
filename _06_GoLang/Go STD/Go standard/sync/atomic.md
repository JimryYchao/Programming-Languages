<a id="TOP"></a>

## Package atomic

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/atomic_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/sync/atomic"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `atomic` 提供了用于实现同步算法的低级原子内存原语。

这些函数需要非常小心才能正确使用。除了特殊的低级应用程序外，最好使用通道或 `sync` 包的工具来完成同步。**通过通信来共享内存；不要通过共享记忆来通信。**

`AddT` 函数实现的 add 操作在原子上等同于：

```go
*addr += delta
return *addr
```

`SwapT` 函数实现的交换操作在原子上等同于：

```go
old, *addr = *addr, new
return old
```

`CompareAndSwapT` 函数实现的 *compare-and-swap* 操作在原子上等效于：

```go
if *addr == old {
	*addr = new
	return true
}
return false
```

由 `LoadT` 和 `StoreT` 函数实现的加载和存储操作是 `return *addr`和 `*addr = value` 的原子等价物。

---
<a id="exam" ><a>