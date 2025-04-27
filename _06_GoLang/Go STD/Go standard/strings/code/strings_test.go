package gostd

import (
	"os"
	"strings"
	"testing"
	"unicode"
	"unsafe"
)

var content = `This is a content for strings test.
hello world
HELLO WORLD
你好，世界
abc123456789XYZ
`

func checkBool(b bool) {
	log(b)
}

// strings functions
func TestStringFunctions(t *testing.T) {
	// ! Clone 返回一个原字符串对象的副本
	t.Run("Clone", func(t *testing.T) {
		c := strings.Clone(content)
		checkBool(c == content)                                       // true
		checkBool(unsafe.StringData(c) == unsafe.StringData(content)) // false
	})

	// ! Compare 按字典顺序比较两个字符串
	t.Run("Compare", func(t *testing.T) {
		log(strings.Compare("a", "b")) // -1
		log(strings.Compare("a", "a")) // 0
		log(strings.Compare("b", "a")) // 1
	})
	// ! Count 统计 s 中的 substr 的非重叠实例数；substr 是 "" 时返回 1 + s 的 Unicode 码位数
	t.Run("Count", func(t *testing.T) {
		log(strings.Count("cheese", "e")) // 3
		log(strings.Count("five", ""))    // 1+4
	})

	// ! EqualFold 报告在简单 Unicode 大小写折叠下，解释为 UTF-8 字符串的 s 和  t 是否相等
	t.Run("EqualFold", func(t *testing.T) {
		checkBool(strings.EqualFold("Go", "go"))     // true
		checkBool(strings.EqualFold("ABcD", "abCD")) // true because comparison uses simple case-folding
		checkBool(strings.EqualFold("ß", "ss"))      // false because comparison does not use full case-folding

	})

	// ! Field 在连续空白前后拆分字符串并返回子字符串的切片组；FieldFunc 在满足 fn(rune) 的码位处前后拆分字符串并返回子字符串的切片组
	t.Run("Fields", func(t *testing.T) {
		logPerfln("Field: %s", strings.Fields(content)...)
		logPerfln("Field: %s", strings.FieldsFunc(content, func(r rune) bool { return r == '\n' })...)
	})

	// ! HasPrefix 报告是否以给定前缀开头；HasSuffix 报告是否以给定后缀结尾
	t.Run("has prefix or suffix", func(t *testing.T) {
		checkBool(strings.HasPrefix("Gopher", "Go")) // true
		checkBool(strings.HasPrefix("Gopher", "C"))  // false
		checkBool(strings.HasPrefix("Gopher", ""))   // true

		checkBool(strings.HasSuffix("Amigo", "go"))  // true
		checkBool(strings.HasSuffix("Amigo", "O"))   // false
		checkBool(strings.HasSuffix("Amigo", "Ami")) // false
		checkBool(strings.HasSuffix("Amigo", ""))    // true
	})
	// ! Join 将其第一个参数的元素连接起来以创建单个字符串。sep 用于间隔元素
	t.Run("Join", func(t *testing.T) {
		s := []string{"foo", "bar", "baz"}
		log(strings.Join(s, ", ")) // foo, bar, baz
	})

	// ! Map 返回字符串 s 的副本，其中所有字符都根据映射函数 `mapping` 进行了修改
	t.Run("Map", func(t *testing.T) {
		rot13 := func(r rune) rune {
			switch {
			case r >= 'A' && r <= 'Z':
				return 'A' + (r-'A'+13)%26
			case r >= 'a' && r <= 'z':
				return 'a' + (r-'a'+13)%26
			}
			return r
		}
		log(strings.Map(rot13, "'Twas brillig and the slithy gopher...")) // 'Gjnf oevyyvt naq gur fyvgul tbcure...
	})

	// ! Repeat 重复 count 个字符串 s 以组成一个新串
	log("ba" + strings.Repeat("na", 2)) // banana

	// ! Replace & ReplaceAll 返回一个替换 old 串为 new 串的字符串副本
	t.Run("Replace", func(t *testing.T) {
		log(strings.Replace("oink oink oink", "k", "ky", 2))      // oinky oinky oink
		log(strings.Replace("oink oink oink", "oink", "moo", -1)) // moo moo moo
		log(strings.ReplaceAll("oink oink oink", "oink", "moo"))  // moo moo moo
	})

	// ! ToLower, ToUpper, ToTitle
	t.Run("ToLower,ToUpper,ToTitle", func(t *testing.T) {
		log(strings.ToLower("Gopher"))                               // gopher
		log(strings.ToTitle("Gopher"))                               // gopher
		log(strings.ToUpper("Gopher"))                               // GOPHER
		log(strings.ToLowerSpecial(unicode.TurkishCase, "Önnek İş")) // önnek iş
		log(strings.ToTitleSpecial(unicode.TurkishCase, "Önnek İş")) // ÖRNEK İŞ
		log(strings.ToUpperSpecial(unicode.TurkishCase, "Önnek İş")) // ÖRNEK İŞ
	})

	// ! ToValidUTF8 返回字符串 s 的副本，其中每次运行无效的 UTF-8 字节序列都被 replacement 字符串（可以为空）替换。
	t.Run("ToValidUTF8", func(t *testing.T) {
		logfln("%s", strings.ToValidUTF8("abc", "\uFFFD"))            // abc
		logfln("%s", strings.ToValidUTF8("a\xffb\xC0\xAFc\xff", ",")) // a,b,c,
		logfln("%s", strings.ToValidUTF8("\xed\xa0\x80", "abc"))      // abc
	})
}

// ! Contains-like 类函数报告在字符串中是否包含符合给定要求的子串
func TestContainsLikes(t *testing.T) {
	//! Contains 报告 s 中是否包含子串 substr
	checkBool(strings.Contains(content, "hello")) // true
	//! ContainsAny 报告 s 中是否包含子串中任一的 Unicode 码位
	checkBool(strings.ContainsAny(content, "你好"[:4])) // true
	//! ContainsFunc 报告 s 中的字符是否有满足 f(c) 返回 true
	checkBool(strings.ContainsFunc(content, func(r rune) bool { return r == '好' })) // true
	//! ContainsRune 报告 s 中是否包含 Unicode 码位 r
	checkBool(strings.ContainsRune(content, '你')) // true
}

// ! Cut-likes 类函数在给定子串的周围剪切切片 s，并返回是否找到给定子串
func TestCutLikes(t *testing.T) {
	// ! Cut 在 sep 的第一个实例周围剪切 s
	t.Run("Cut", func(t *testing.T) {
		show := func(s, sep string) {
			before, after, found := strings.Cut(s, sep)
			logfln("Cut(%q, %q) = %q, %q, %v", s, sep, before, after, found)
		}
		show("Gopher", "Go")     // "", "pher", true
		show("Gopher", "ph")     // "Go", "er", true
		show("Gopher", "er")     // "Goph", "", true
		show("Gopher", "Badger") // "Gopher", "", false
		show("Gopher", "")       // "", "Gopher", true
	})

	// ! CutPrefix 返回不带 sep 作为前缀的后向子串
	t.Run("CutPrefix", func(t *testing.T) {
		show := func(s, sep string) {
			after, found := strings.CutPrefix(s, sep)
			logfln("CutPrefix(%q, %q) = %q, %v", s, sep, after, found)
		}
		show("Gopher", "Go") // "pher", true
		show("Gopher", "ph") // "Gopher", false
		show("Gopher", "er") // "Gopher", false
		show("Gopher", "")   // "Gopher", true
	})

	// ! CutSuffix 返回不包含 sep 作为后缀的前向子串
	t.Run("CutSuffix", func(t *testing.T) {
		show := func(s, sep string) {
			before, found := strings.CutSuffix(s, sep)
			logfln("CutSuffix(%q, %q) = %q, %v\n", s, sep, before, found)
		}
		show("Gopher", "Go") // "Gopher", false
		show("Gopher", "ph") // "Gopher", false
		show("Gopher", "er") // "Goph", true
		show("Gopher", "")   // "Gopher", true
	})

}

// ! Index-like 返回符合给定条件的首个实例的索引位置，不满足时返回 -1。相应的 LastIndex-like 函数返回最后一个实例的索引位置
func TestIndexLikes(t *testing.T) {
	// ! Index 返回 s 中 substr 的第一个实例的索引;
	log(strings.Index("chicken", "ken")) // 4
	// ! IndexAny 返回来自 s 中任何 Unicode 码位第一个实例的索引
	log(strings.IndexAny("chicken", "aeiouy")) // 2
	// ! IndexByte 返回字节 c 在 s 中的第一个实例的索引
	log(strings.IndexByte("golang", 'x')) // -1
	// ! IndexFunc 返回 s 中第一个满足 f(c) 的 Unicode 码位的索引
	log(strings.IndexFunc("Hello, 世界", func(r rune) bool { return unicode.Is(unicode.Han, r) })) // 7
	// ! IndexRune 返回s 中 Unicode 码位 r 的第一个实例的索引
	log(strings.IndexRune("chicken", 'k')) // 4
}

// ! Split-like 字符串切片
func TestSplitLikes(t *testing.T) {
	// ! Split 切片所有由 sep 分隔的子字符串
	logfln("%q\n", strings.Split("a,b,c,,,", ",")) // ["a" "b" "c" "" "" ""]
	// ! SplitN 最多 n-1 次切片由 sep 分隔的子字符串，返回 n 个子串；n=-1 与 Split 相同
	logfln("%q\n", strings.SplitN("a,b,c,,,", ",", 3)) // ["a" "b" "c,,,"]
	// ! SplitAfter 切片 sep 之后的子字符串
	logfln("%q\n", strings.SplitAfter("a,b,c,,,", ",")) // ["a," "b," "c," "," "," ""]
	// ! SplitAfterN 最多 n-1 次切片 sep 之后的子字符串，返回 n 个子串；n=-1 与 SplitAfter 相同
	logfln("%q\n", strings.SplitAfterN("a,b,c,,,", ",", 3)) // ["a," "b," "c,,,"]
}

// ! Trim-like 删除满足给定条件的字符串片段
func TestTrimLikes(t *testing.T) {
	// ! Trim 删除 cutset 中所有的前导和尾随 Unicode 码位；TrimLeft 只删除前导部分；TrimRight 只删除尾随部分；
	// ! TrimFunc 删除了所有满足 f(c) 的前导和尾随 Unicode 码位 c；TrimLeftFunc 只删除前导部分，TrimRightFunc 只删除尾随部分
	// ! TrimSpace 删除前导和尾随的所有空白；
	// ! TrimPrefix 删除前缀；TrimSuffix 删除后缀；

	log(strings.Trim("¡¡¡Hello, Gophers!!!", "!¡")) // Hello, Gophers

	log(strings.TrimFunc("¡¡¡Hello, Gophers!!!", func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) })) // Hello, Gophers
}

// ! strings.Builder 使用 Builder.Write 方法高效地构建字符串。以最大限度地减少内存复制。
func TestStringBuilder(t *testing.T) {
	sb := strings.Builder{}

	if sb.Cap() < 32 {
		sb.Grow(32)
	}
	sb.Write([]byte{'a', 'b', 'c'})
	sb.WriteByte('d')
	sb.WriteRune('我')
	sb.WriteString("1234")

	logfln("len:%d, string:%s", sb.Len(), sb.String())
	sb.Reset()

	sb.WriteString(content)
	logfln("len:%d, string:%s", sb.Len(), sb.String())
}

// ! strings.Reader 实现 io.Reader,io.ReaderAt,io.ByteReader,io.ByteScanner,io.RuneReader,io.RuneScanner,io.Seeker,io.WriterTo接口以从 s 中读取
func TestStringReader(t *testing.T) {
	sr := strings.NewReader("Hello World")
	for sr.Len() > 0 {
		if b, err := sr.ReadByte(); err == nil {
			logfln("read byte %#v %[1]c", b)
		}
	}
	sr.Reset(content)
	sr.WriteTo(os.Stdout)
}

// ! strings.Replacer 用替换项替换字符串列表。它对于多个 goroutine 并发使用是安全的。
func TestStringReplacer(t *testing.T) {
	r := strings.NewReplacer("<", "&lt;", ">", "&gt;")
	log(r.Replace("This is <b>HTML</b>!"))

	r = strings.NewReplacer(" ", "_", "，", ",")
	r.WriteString(os.Stdout, content)
}
