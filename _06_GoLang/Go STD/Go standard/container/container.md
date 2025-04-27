<a id="TOP"></a>

## container

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<!-- <a href="./code/XXX_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a> -->
	<!-- <a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	 -->
	<a href="https://pkg.go.dev/container"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `heap` 为实现 `heap.Interface` 的任何类型提供堆操作。堆是一棵树，它的属性是每个节点都是它的子树中的最小值节点。树中的最小元素是根，索引为 0。堆是实现优先级队列的常用方法。要构建一个优先级队列，使用（负）优先级实现 `Heap` 接口作为 `Less` 方法的排序，`Push` 添加项目，而 `Pop` 从队列中删除最高优先级的项目。[[↗]](code/heap_test.go)

包 `list` 实现了双向链表。[[↗]](code/list_test.go)

包 `ring` 实现了循环列表上的操作。[[↗]](code/ring_test.go)

---
<a id="exam" ><a>

### Examples

- [heap: PriorityQueue](examples/priorityQueue_test.go)

---