package gostd

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"testing"
)

/*
! flag.Command 是默认的命令行标志集，从 os.Arg 解析过来，顶级函数是对 CommandLine 方法的包装
! flag.Usage() 将记录所有已定义的命令行标志的使用情况消息打印到 `flag.CommandLine` 的输出，默认为 os.Stderr
	当分析标志时, 发生错误时调用它。可以更改为指向自定义函数；默认调用 `PrintDefaults`。
	自定义函数可以选择退出程序；例如，命令行错误处理退出策略设置为 ExitOnError 时，退出就会发生
! PrintDefaults 默认在标准错误流上打印 flag 使用情况的信息；除非另有配置
*/

func TestFlagUsage(t *testing.T) {
	t.Run("Usage", func(t *testing.T) {
		flag.Usage()
	})

	t.Run("PrintDefaults", func(t *testing.T) {
		resetForTesting(flag.Usage, nil)
		setFlags(flag.CommandLine)
		flag.PrintDefaults()
	})
}

/* flag functions
! Args, Arg 返回命令行参数，Arg(0) 是处理标志后剩余的第一个参数；标志不存在时返回空
! NArg 是处理标志后剩余的参数数
! NFlag 返回已设置的命令行标志的数量
! Set 设置命名命令行标志 name 的值。
! Parse 解析 os.args[1:] 中的命令行标志; 必须在定义所有标志之后和程序访问标志之前调用;
! Parsed 报告是否已分析命令行标志。
! Set flags:
	types: Bool, Duration, Float64, Int, Int65, String, Uint, Uint64
	Var: Var, BoolVar, DurationVar, Float64Var, IntVar, Int64Var, StringVar, TextVar, UintVar, Uint64Var
! TextVar 使用指定的 name、默认值 value 和用法字符串 usage 定义标志。
	参数 p 必须是指向一个变量的指针，该变量将保存标志的值，并且 p 必须实现 encoding.TextUnmarshaler。
	如果使用了标志，标志值将被传递给 p 的 UnmarshalText 方法。默认值的类型必须与 p 的类型相同
! Func, BoolFunc 使用指定名称 name 和用法字符串 usage 定义一个标志，当看到该标志时都会使用调用 fn(var), 返回的错误被视为解析错误
! Visit 按字典顺序访问命令行标志，为每个标志调用 fn。它只访问那些已经设置的标志。
! VisitAll 按字典顺序访问命令行标志，并为每个标志调用 fn。它访问所有的旗帜，即使是那些没有设置
! flag.ErrorHandling 定义在分析失败时 FlagSet.Parse 的行为。如果分析失败，这些常量将使 FlagSet.Parse 按照所述方式运行
	ContinueOnError   返回一个错误描述
	ExitOnError		  调用 os.Exit(2) 或 -h/-help Exit(0)
	PanicOnError	  调用 panic(err)
! flag.Flag 代表一个标志的状态
! Lookup 返回指定名称 name 的标志的 *Flag。

! UnquoteUsage 从 *Flag 的用法字符串 usage 中提取一个带反引号的名称，并返回该名称和不带引号的 usage。
	给定 "a `name` to show"，它返回 ("name"，"a name to show"）。如果没有反引号，则名称是对标志值类型的推测
	如果标志为布尔值，则 name 为空字符串。
*/

func TestFlagFuncs(t *testing.T) {
	var fs = flag.NewFlagSet("userFlagSet", flag.ContinueOnError)
	args := setFlags(fs)
	if !fs.Parsed() {
		fs.Parse(args)
		logfln("total set %d flags at %s", fs.NFlag(), fs.Name())
	}
	// Visit
	fs.VisitAll(func(f *flag.Flag) {
		v := f.Value.String()
		if f.Value.String() != f.DefValue {
			f.Value.Set(f.DefValue)
			logfln("%s = %s, set default:%s", f.Name, v, f.Value)
		}
	})

	fs.Int("n_int", 10010, "new int value")
	// Lookup
	fs.Set("t_int", "10086")
	logfln("-t_int = %s", fs.Lookup("t_int").Value)

	// UnquoteUsage
	name, usage := flag.UnquoteUsage(fs.Lookup("t_ip"))
	logfln("name : %s, usage : %s", name, usage)
}

/* flag type
! NewFlagSet 返回一个指定名称和错误处理标志的新空标志集 FlagSet; 如果 name 不为空，它将被打印在默认 Usage 信息和错误信息中
! flag.FlagSet 表示一组已定义的标志，标志集的零值没有名称，并且具有 ContinueOnError 错误处理，标志集的名称必须唯一
non-top functions
	Name 返回标志集的名称。
	ErrorHandling 返回标志集的错误处理行为；
	OutPut 返回 usage 和 error 的输出；未设置时默认为 os.Stderr
	SetOutput 设置 usage 和 error 的输出。如果为 nil，则使用 os.Stderr。
	Init 为标志集设置名称 name 和错误处理属性。默认情况下，零 FlagSet 使用空名称和 ContinueOnError 错误处理策略。
	Parse 从参数列表中解析标志定义，该列表不应包括 command 名。必须在定义 FlagSet 中的所有标志之后且在程序访问标志之前调用。
		如果设置了 -help 或 -h 但未定义，则返回值为 ErrHelp。
! flag.Value 是存储在标志中的动态值的接口。
	如果 Value 有一个返回 true 的 `IsBoolFlag() bool` 方法，命令行解析器会使 -name 等效于 -name=true
	对于每个存在的标志，按命令行顺序调用 Set 一次。flag 包可以用一个零值的接收器调用 String 方法，比如一个 nil 指针
! flag.Getter 是一个接口，允许检索 Value 的内容。提供一个 Get 方法
	这个包提供的所有 Value 类型都满足 Getter 接口，除了 Func 使用的类型。
*/

func TestSetFlagSet(t *testing.T) {
	fs := flag.NewFlagSet("myFlags", flag.ContinueOnError)
	names := setFlags(fs)
	fs.Init("mFlags2", flag.ExitOnError)
	logfln("reInit a FlagSet(%s)", fs.Name())
	fs.SetOutput(os.Stdout)

	if fs.Parsed() {
		fs.Parse(names)
	}

	if g, ok := fs.Lookup("t_ip").Value.(flag.Getter); ok {
		logfln("get t_ip : %v, type : %s", g.Get(), reflect.TypeOf(g.Get()).Kind())
	}

}

var flagUsage = flag.Usage

func resetForTesting(usage func(), output io.Writer) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	if output == nil {
		output = os.Stderr
	}
	flag.CommandLine.SetOutput(output)
	flag.CommandLine.Usage = flagUsage
	flag.Usage = usage
}

// Declare a user-defined flag type.
type flagVar []string

func (f *flagVar) String() string {
	return fmt.Sprint([]string(*f))
}

func (f *flagVar) Set(value string) error {
	*f = append(*f, value)
	return nil
}

type VarStat struct {
	flag.Value
	name  string
	usage string
}

func setFlags(f *flag.FlagSet, argus ...VarStat) []string {
	var args = []string{"-t_bool", "-t_int", "-t_int64", "-t_uint", "-t_uint64", "-t_string", "-t_float64", "-t_duration", "-t_ip", "-t_userflag", "-t_func", "-t_boolfunc"}
	f.Bool("t_bool", true, "bool value")
	f.Int("t_int", 1, "int value")
	f.Int64("t_int64", 2, "int64 value")
	f.Uint("t_uint", 3, "uint value")
	f.Uint64("t_uint64", 4, "uint64 value")
	f.String("t_string", "5", "string value")
	f.Float64("t_float64", 6, "float64 value")
	f.Duration("t_duration", 7, "time.Duration value")
	f.TextVar(&net.IP{}, "t_ip", net.IPv4(192, 168, 0, 100), "`IP address` to parse")
	f.Var(&flagVar{"a", "b", "c"}, "t_userflag", "flagVar value")
	for _, v := range argus {
		f.Var(v.Value, v.name, v.usage)
		args = append(args, "-"+v.name)
	}
	f.Func("t_func", "func value", func(string) error { return nil })
	f.BoolFunc("t_boolfunc", "boolfunc value", func(string) error { return nil })
	return args
}
