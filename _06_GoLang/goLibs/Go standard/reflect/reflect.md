<a id="TOP"></a>

## Package reflect

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/reflect_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/reflect"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>


包 `reflect` 实现运行时反射，允许程序操作具有任意类型的对象。典型的用法是使用静态类型 `interface{}` 获取值，并通过调用 `TypeOf` 提取其动态类型信息。

计算中的反射是程序检查自身结构的能力，特别是通过类型；它是元编程的一种形式。

**反射第一定律**：反射从接口值到反射对象。

在基本级别上，反射只是一种检查存储在接口变量中的类型和值对的机制。`reflect` 中包含两种类型：`Value` 和 `Type`，它们提供给予对接口变量内容的访问权。`reflect.TypeOf` 和 `reflect.Value` 从接口值中检索 `Type` 和 `Value`。

```go
var x float64 = 3.1415
t := reflect.TypeOf(x)
v := reflect.ValueOf(x)

print(t, v)                   // float64 3.1415
print(t.Kind(), v.Kind())     // float64 float64
print(t.String(), v.String()) // float64 <float64 Value>
f = v.Float()               // f = x
```

**反射第二定律**：反射从反射对象到接口值。

Go 反射也会产生自己的逆反射。给定一个 `reflect.Value`，可以使用 `v.Interface` 方法恢复一个接口值；该方法将类型和值信息打包回接口表示并返回结果。

```go
f := v.Interface().(float64) 	 // f = x = 3.1415
fmt.Printf("%f", v.Interface())  // 3.141500
```

**反射第三定律**：若要修改反射对象，该值必须是可设置的。可设置性是反射 `Value` 的一个属性；`Value` 的 `CanSet` 方法报告了 `Value` 的可设置性，在不可设置的 `Value` 上调用 `Set` 方法是错误的：

```go
var x float64 = 3.1415
v := reflect.ValueOf(x)
print(v.CanSet())
// v.SetFloat(1.23456)  // Error: will panic

p := reflect.ValueOf(&x)
log(p.Type(), p.CanSet())               // *float64 false
log(p.Elem().Type(), p.Elem().CanSet()) // float64 true
p.Elem().SetFloat(1.23456)
log(p.Elem())  // 1.23456
log(x)         // 1.23456
```

---
<a id="exam" ><a>

### Examples

- [reflectHelper](reflect/helper/doc.go)

---