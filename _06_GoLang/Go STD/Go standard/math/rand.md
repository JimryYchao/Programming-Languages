<a id="TOP"></a>

## Package rand

- [rand/v2](https://pkg.go.dev/math/rand/v2)
- [rand/v2 code](./code/randv2_test.go)


<div id="top" style="z-index:99999999;position:fixed;bottom:35px;right:50px;float:right">
	<a href="./code/rand_test.go" target="_blank"><img id="img-code" src="../_rsc/to-code.drawio.png" ></img></a>
	<a href="#TOP" ><img id="img-top" src="../_rsc/to-top.drawio.png" ></img></a>	
	<a href="https://pkg.go.dev/math/rand"  target="_blank"><img id="img-link" src="../_rsc/to-link.drawio.png" ></img></a>
	<a href="..\README.md"><img id="img-back" src="../_rsc/back.drawio.png"></img></a>
</div>

包 `rand` 实现了适合模拟等任务的伪随机数生成器，但不应用于安全敏感的工作。

随机数是由一个源 `Source` 产生的，通常包装在一个 `Rand`。这两种类型都应该由一个 *goroutine* 同时使用：在多个 *goroutine* 之间共享需要某种同步。

`Float64` 和 `Int`，对于多个 *goroutine* 并发使用是安全的。

这个包的输出可能很容易预测，不管它是如何设置种子。有关适用于安全敏感工作的随机数，请参阅 [`crypto/rand`](https://pkg.go.dev/crypto/rand)。

```go
func TestRand(t *testing.T) {
	answers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}
	for range 10 {
		fmt.Println("Magic 8-Ball says:", answers[rand.Intn(time.Now().Nanosecond())%len(answers)])
	}
}
```

---
<a id="exam" ><a>