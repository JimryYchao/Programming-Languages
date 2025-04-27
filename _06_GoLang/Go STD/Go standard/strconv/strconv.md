<a id="TOP"></a>

## Package strconv

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/strconv_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<!-- <a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	 -->
	<a href="https://pkg.go.dev/strconv"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `strconv` 实现了基本数据类型的字符串表示形式之间的转换。

最常见的数字转换是 `Atoi(string to int)` 和 `Itoa(int to string)`。

```go
i, err := strconv.Atoi("-42")
s := strconv.Itoa(-42)
```

`ParseT` 将字符串转换为值；`FormatT` 将值转换为字符串；`AppendT` 则将格式化的值追加到目标切片：

```go
b, err := strconv.ParseBool("true")
f, err := strconv.ParseFloat("3.1415", 64)
i, err := strconv.ParseInt("-42", 10, 64)
u, err := strconv.ParseUint("42", 10, 32)

s := strconv.FormatBool(true)
s := strconv.FormatFloat(3.1415, 'E', -1, 64)
s := strconv.FormatInt(-42, 16)
s := strconv.FormatUint(42, 10)
```

`Quote` 和 `QuoteToASCII` 将字符串转换为带引号的 Go 字符串文字。后者通过使用 `\u` 转义任何非 ASCII Unicode 来保证结果是 ASCII 字符串；`Unquote` 和 `UnquoteChar` 为逆函数：

```go
q := strconv.Quote("Hello, 世界")
q := strconv.QuoteToASCII("Hello, 世界")
```

---
<a id="exam" ><a>

### Examples

- [file StringWriter](code/stringWriter.go)

---