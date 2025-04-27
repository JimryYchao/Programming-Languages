<a id="TOP"></a>

## Package fmt

<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/fmt_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/fmt" ><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `fmt` 使用类似于 C 的 `printf` 和 `scanf` 的函数实现格式化 I/O。“*verbs*” 格式是从 C 派生的。

---
### verbs

```go
%v		// value 的默认格式; 对于 struct, %+v 添加字段名
%#v		// value 的 go 语法表示
%T		// type of value 的 go 语法表示
%% 		// %

// bool 
%t		// boolean 的字面值
// integer
%b		// 二进制整数
%c		// Unicode code 对应的字符
%d		// 十进制整数
%o		// 八进制整数
%O		// 0o 前缀的八进制整数
%q		// 安全转义的单引号字符，例如 65 转义为 'A'，无效整数码位改为 U+FFFD
%x,%X	// 十六进制整数
%U		// 如 Unicode `U+1234` 形式 
// floating-point & complex
%b		// 无小数 p 计数法, -123p-45
%e,%E	// e 计数法
%f,%F	// 浮点数
%g,%G	// 择优选择 %e/%E 还是 %f/%F
%x,%X	// 十六进制 p 计数法

// strings & []byte
%s		// 字符串或字节切片的无解释字节
%q		// 安全转义的双引号字符串
%x,%X	// 逆向转义为 \xhh，每个字节两个字符，例如 "我" > e68891

// pointers
%p		// 十六进制地址值，切片返回 &s[0] 的地址，指针用于整数 verbs 时被视为一般整数；例如 0xc0000dc4e0
```

对于 `%v` 的默认格式为：

```go
bool					%v >> %t
int,int8,...			%v >> %d
uint,uint8,...			%v >> %d, %#v >> %#x
float32,complex64,...	%v >> %g
string					%s
chan					%p
pointer					%p
struct					{field0 field1 ...}
array,slice				[elem0 elem1 ...]
map						map[k1:v1 k2:v2 ...]
compound pointers		&{}, &[], &map[] 
```

宽度和精度由 $width.precision$ 表示，宽度表示要输出的最小字符（`rune`）数，必要时空格填充。浮点的精度默认为 6；对于 `string` 和 `bytestring`，精度限制了格式化输入的长度，以 `rune` 为单位；`%x, %X` 以字节为单位。对于复数，宽度和精度独立应用到它的两个分量。

其他的一些标志位：

```go
'+'			// 对于数值打印符号位；%+q 保证至输出 ASCII
'-'			// 空格填充在右侧
'#'			/* %#b,%#o,%#x,%#X		打印前导 0b,0,0x,0X
			   %#p					取消打印前导 0x		
			   %#q					打印原始字符串（若支持）
			   %#f,%#e,%#g			打印小数点，不删除 %g 的零 
			   %#U					U+0078 'x', 同时打印对应的可打印字符  */
' '			// % d 整数预留符号位空格；% x 字节之间放置空格
'0'			// 0 代替空格作为宽度的填充，对于字符串类型将忽略
```

`[n]` 表示显示参数索引，`("%[2]d %[1]d\n", 11, 22)` 将产生 `"22 11"`。显式索引会影响后续的 *verbs*，后续的索引依次递增，`("%d %d %#[1]x %#x", 16, 17)` 将产生 `"16 17 0x10 0x11"`

如果操作数是接口值，则使用内部具体值。除了 `%T, %p`，实现某些接口的操作数应用一些特殊规则：
1. 如果操作数是一个 `reflect.Value`，则操作数将被它所保存的具体值替换，并且打印将继续执行下一个规则。 
2. 如果一个操作数实现了 `Formatter` 接口，它将被调用，*verbs* 和标志的解释由该实现控制。
3. 如果使用 `%#v`，并且操作数实现了 `GoStringer` 接口，则将调用该接口。如果格式（对于 `Println` 等，隐式为 `%v`）对于字符串（`%s %q %x %X`）有效，或者是 `%v` 但不是 `%#v`，则适用以下两个规则：
4. 如果一个操作数实现了 `error` 接口，则调用 `Error` 以将对象转换为字符串，然后根据 *verb*（如果有的话）的要求进行格式化。
5. 如果一个操作数实现了方法 `String() string`，则调用 `String()` 以将对象转换为字符串，然后根据 *verb*（如果有）的要求进行格式化。


---
### Printers

 `Printf-like` 函数使用格式化 `format` 输出。
 `Print-like` 函数在两个非字符串参数之间添加一个空格。对每个参数使用 `%v`。
 `Println-like` 函数在操作数之间空格，并追加一个换行符。

```go
// ? go test -v -run=^TestPrinters$
func TestPrinters(t *testing.T) {
	printErrNew := func(mess string) error {
		return fmt.Errorf("printer error: %s", mess)
	}
	dir := t.TempDir()
	file, err := os.CreateTemp(dir, "fmt")
	if err != nil {
		t.Fatal(printErrNew(err.Error()))
	}
	defer file.Close()
	var WriteFormat = "Write: %s"
	var printArgs = []any{1, "H", true, "E", 'A', "L",
		struct{ name string }{"World"}, "L", []byte("!!!"), "O", [...]int{1, 2, 3, 4}, []int{55, 6, 7, 8},
		map[string]int{"a": 1, "b": 2}, &struct{ n int }{100},
		func(a, b int) int { return a + b }, make(chan bool), new(fmt.Formatter)}

	fmt.Printf("\nfmt.Sprintf: \n%s", fmt.Sprintf(WriteFormat, "test Printf\n"))
	fmt.Printf("\nfmt.Sprint: \n%s", fmt.Sprint(printArgs...))
	fmt.Printf("\nfmt.Sprintln: \n%s", fmt.Sprintln(printArgs...))

	fmt.Fprintf(file, WriteFormat, "test Fprintf\n")
	io.Copy(os.Stdout, file)
}
```

---
### Scanners

`Scanf-like` 函数根据格式字符串解析参数。
`Scan-like` 函数将输入的换行符视为空格。
`Scanln-like` 函数在扫描到换行符或 EOF 后停止。
	
*Scanners* 的 *verbs* 中不包括 `%T`、`%p`、`#`、`+`；`\r\n` 的意思与 `\n` 相同；由 *verbs* 处理的输入是隐式空白分隔的：除了 `%c` 其他 *verbs* 的实现将忽略任何的前导空白；`%s`（和读取字符串时的 `%v`）在第一个空白之后或换行符处停止使用输入。

宽度可以在输入文本中解释，但没有用于以精度扫描的语法（没有 `%5.2f`，只有 `%5f`）。如果提供了 *width*，则它将在前导空格被修剪后应用，并指定要读取的 *rune* 的最大数量。例如 `Sscanf(" 1234567 ", "%5s%d", &s, &i)` 中 `s="12345", i=67`；`Sscanf(" 12 34 567 ", "%5s%d", &s, &i)` 中 `s="12", i=34`。

在所有的 *Scanners* 函数中，如果一个操作数实现了 `Scan` 方法（它实现了 `fmt.Scanner` 接口），那么该方法将被用于扫描该操作数的文本。要扫描的所有参数必须是指向基本类型的指针或 `Scanner` 接口的实现。

`Fscan` 等可以读取它们返回的输入后的一个字符（`rune`），这意味着调用扫描例程的循环可能会跳过一些输入。这通常只在输入值之间没有空格时才有问题。如果提供给 `Fscan` 的读取器实现了 `fmt.ScanState.ReadRune`，则该方法将用于读取字符。如果读取器也实现 `fmt.ScanState.UnreadRune`，则该方法将用于保存字符，并且后续调用不会丢失数据。若要将 `ReadRune` 和 `UnreadRune` 方法附加到没有此功能的读取器，使用 `bufio.NewReader`。

```go
//? go test -v -run=^TestScanners$
func TestScanners(t *testing.T) {
	var text string = "5 true gophers\naaaaa"
	t.Run("Scan-like", func(t *testing.T) {
		var (
			b     bool
			n     int
			s, s0 string
		)
		if c, err := fmt.Sscan(text, &n, &b, &s, &s0); err == nil {
			fmt.Printf("Read %d : %s", c, fmt.Sprintln(n, b, s, s0))
		}
	})
	t.Run("Scanln-like", func(t *testing.T) {
		var (
			b bool
			n int
			s string
		)
		if c, err := fmt.Sscanln(text, &n, &b, &s); err == nil {
			fmt.Printf("Read %d : %s", c, fmt.Sprintln(n, b, s))
		}
	})

	t.Run("Scanf-like", func(t *testing.T) {
		var (
			r rune
			b bool
			s string
			h int
		)
		if c, err := fmt.Sscanf(text, "%c %t %s\n%x", &r, &b, &s, &h); err == nil {
			fmt.Printf("Read %d : %s", c, fmt.Sprintln(r, b, s, h))
		}
	})
}
```

---