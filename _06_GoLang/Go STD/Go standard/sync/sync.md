<a id="TOP"></a>

## Package sync

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/sync_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/sync"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `sync` 提供基本的同步原语，如互斥锁。除了 `Once` 和 `WaitGroup` 类型之外，大多数类型都是供低级库例程使用的。更高级别的同步最好通过通道和通信来完成。不应复制包含此包中定义的类型的值。

---
<a id="exam" ><a>

### Examples

- [SyncPool](./examples/syncPool.go)

- [Wrong locker race](./examples/race_test.go)

---