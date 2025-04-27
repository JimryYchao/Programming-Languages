<a id="TOP"></a>

## Package big

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/big_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<!-- <a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	 -->
	<a href="https://pkg.go.dev/math/big"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `big` 实现了任意精度的算术（大数字）。支持：

```go
Int    // signed integers
Rat    // rational numbers
Float  // floating-point numbers

var i big.Int
var r = &big.Rat{}
var f = new(big.Float)

i = *big.NewInt(10010) // &i = 10010
f = big.NewFloat(0.1)  // f = 0.1
r = big.NewRat(1, 3)   // r = 1/3

var z1 = &big.Int{}
z1.SetUint64(123)                   // 123
z2 := new(big.Rat).SetFloat64(1.25) // 5/4
z3 := new(big.Float).SetInt(z1)     // 123.0
```

Setter, Numeric operations, Predicates 表示为以下形式的方法：

```go
func (z *T) SetV(v V) *T          // z = v
func (z *T) Unary(x *T) *T        // z = unary x
func (z *T) Binary(x, y *T) *T    // z = x binary y
func (x *T) Pred() P              // p = pred(x)

f := new(big.Float).SetInt64(10010)
fadd := f.Add(f, f)
r2, _ := f.Mul(f, big.NewFloat(0.666)).Rat(nil)
log(f, fadd, f.Prec(), r2)
```

算术表达式通常将结果值保存到方法的接收器，无论此前接收器的值为何：

```go
c.Add(a, b) 		// c = a + b
sum.Add(sum, b)     // sum += b
```

---
<a id="exam" ><a>

### Examples

- [some math calculation](./examples/bigExamples.go)

---