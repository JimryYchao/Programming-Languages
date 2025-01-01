package gostd

import (
	"gostd/examples"
	"sync"
	"testing"
	"time"
)

/* time
! time.Time 表示具有纳秒精度的时间瞬间。应使用 Time 的值而不是指针；GobDecode、UnmarshalBinary、UnmarshalJSON 和 UnmarshalText 非并发安全
	Add, AddData 返回 t 加上一段时间后的 time
	After, Before, Compare, Equal 比较 t 和 u 值的前后
	Clock, Date, Weekday, YearDay, Year, Month, Day, Hour, Minute, Second, Nanosecond  返回 t 的相应时间值
	ISOWeek 返回 t 的年份和周数。周的范围从 1 到 53
	In 返回 t 同一时刻给定 loc 的时间值
	IsDST 报告 t 在其 loc 是否为夏令时
	IsZero 报告 t 是否为 January 1, year 1, 00:00:00 UTC
	Local, UTC 返回基于本地, UTC 的时间
	Round 返回将 t 舍入为 d 的最接近倍数的结果
	Truncate 返回将 t 向下舍入为 d 的倍数的结果
	Location 返回 t 关联的 loc 信息
	Zone 返回 t 的有效时区和及其自 UTC 以东的偏移量
	ZoneBounds 返回在时间 t 有效的时区的边界
	Sub 返回持续时间 t-u 的 duration

	GobEncode 实现了 gob.GobEncoder 接口
	GobDecode 实现了 gob.GobDecoder 接口

	MarshalBinary 实现 encoding.BinaryMarshaler 接口
	UnmarshalBinary 实现 encoding.BinaryUnmarshaler 接口

	MarshalJSON 实现了 json.Marshaler 接口, t 是 RFC 3339 格式
	UnmarshalJSON 实现了 json.Unmarshaler 接口

	MarshalText 实现 encoding.TextMarshaler 接口, t 是 RFC 3339 格式
	UnmarshalText 实现 encoding.TextUnmarshaler 接口
! Data 返回在给定位置区域 loc 的 Time
! Now 返回当前的本地时间
*/

func TestTime(t *testing.T) {
	log := func(name string, v any) {
		logfln("%s: %v", name, v)
	}
	tm := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	log("time", tm)
	log("Add+30d", tm.Add(30*24*60*60*time.Second))
	log("AddDate+1y1m1d", tm.AddDate(1, 1, 1))
	log("Compare(now)", tm.Compare(time.Now()))

	now := time.Now()
	log("now", now)
	h, m, s := now.Clock()
	logfln("Clock: %dh %dm %ds", h, m, s)
	y, mon, d := now.Date()
	logfln("Date: %dy %dm %dd", y, mon, d)
	log("Weekday", now.Weekday())
	log("YearDay", now.YearDay())
	y, w := now.ISOWeek()
	logfln("ISOWeek: %dy %dw", y, w)
	log("UTC", now.UTC())
	log("Location", now.Location())
	name, off := now.Zone()
	logfln("Zone: %s %d", name, off)
	start, end := now.ZoneBounds()
	logfln("ZoneBounds: %v %v", start, end)
	logfln("Sub(%v, %v) : %v", now, tm, now.Sub(tm))

	t.Run("Round & Truncate", func(t *testing.T) {
		for _, r := range ds {
			logfln("d.Round(%6s) = %s\n", r, now.Round(r).String())
		}

		for _, t := range ds {
			logfln("d.Truncate(%6s) = %s\n", t, now.Truncate(t).String())
		}
	})
}

/*
! Unix, UnixMicro, UnixMilli 返回自 1970/1/1 开始的本地位置的 time
! time.Unix, UnixMilli, UnixMicro, UnixNano 返回自 1970/1/1 的秒数, 毫秒数, 微秒数, 纳秒收集
*/

func TestTimeUnix(t *testing.T) {
	unix := func(sec, nsec int64) {
		logfln("Unix(%d, %d) = %v", sec, nsec, time.Unix(sec, nsec).UTC().Format(time.UnixDate))
	}
	unix(0, 0)
	unix(9999999999, 999_999_999)

	now := time.Now()
	logfln("time = (%v)\nUnix: %ds\nUnixMilli: %dms\nUnixMicro: %dμs\nUnixNano: %dns",
		now.Format(time.DateTime), now.Unix(), now.UnixMilli(), now.UnixMicro(), now.UnixNano())
}

/* location
! time.Location 表示在一个地理区域中使用的时间偏移量的集合。location 用于在打印的时间值中提供时区，并用于涉及可能跨越夏令时边界的间隔的计算
! time.Local 表示系统的本地时区
! time.UTC 代表世界协调时间（UTC）。
! FixedZone 返回始终使用给定区域 name 和偏移量（自 UTC 以东的秒数）的 Location
! LoadLocation 返回具有给定名称的 Location。"" 或 "UTC" 返回 time.UTC；"Local" 返回 time.Local；给定的名称参考 IANA Time Zone database
! LoadLocationFromTZData 返回一个 Location，其给定名称是从 IANA 时区数据库格式化的数据中初始化的。数据应该是标准 IANA 时区文件的格式
*/

func TestTimeLocation(t *testing.T) {
	location, _ := time.LoadLocation("America/Los_Angeles")
	now := time.Now()
	log(now.In(time.Local))
	log(now.In(time.UTC))
	log(now.In(location))

	shanghaiLoc := time.FixedZone("Shanghai", 22*60+8*60*60) // UTC+8 + 22 min
	log(now.In(shanghaiLoc))
}

/* duration
! time.Duration 表示两个瞬间之间经过的时间
	Abs 返回 d 的绝对值
	Hours,Microseconds,Milliseconds,Minutes,Nanoseconds,Seconds 返回相应时间单位的持续时间
	Round 返回将 d 四舍五入到 m 的最接近倍数的结果
	String 返回表示持续时间的字符串，格式为 72h3m0.5s
	Truncate 返回将 d 向零舍入为 m 的倍数的结果
! Until 返回自 time.Now 到 t 持续时间
! Since 返回自 t 到 time.Now 经过的时间
*/

func TestDuration(t *testing.T) {
	start := time.Now()
	defer func() {
		logfln("test duration : %s", time.Since(start))
	}()
	d := time.Until(time.Date(2024, 1, 1, 1, 1, 1, 1, time.Local))
	d = d.Abs()

	logfln("duration(%q) : \nNanoseconds: %v: \nMicroseconds: %v: \nMilliseconds: %v: \nSeconds: %v: \nMinutes: %v: \nHours: %v",
		d, d.Nanoseconds(), d.Microseconds(), d.Milliseconds(), d.Seconds(), d.Minutes(), d.Hours())

	for _, r := range ds {
		logfln("d.Round(%6s) = %s\n", r, d.Round(r).String())
	}

	for _, t := range ds {
		logfln("d.Truncate(%6s) = %s\n", t, d.Truncate(t).String())
	}
}

/* Timer
! time.Timer 表示单个计时事件。当 Timer 超时触发时，当前时间将在 Timer.C 上发送，除非 Timer 是由 AfterFunc 创建的。必须使用 NewTimer 或 AfterFunc 创建 Timer。
	Reset 重置计时器持续时间 d；已停止时返回 false; 应始终在已停止或过期的通道上调用 Reset
	Stop 阻止 timer 触发，成功时返回 true，它不关闭 timer.C
! NewTimer 创建一个新的 Timer，它将在至少持续时间 d 之后在其通道上发送当前时间。
! AfterFunc 等待持续时间结束，并在自己的例程中调用 f；timer.Stop 来取消调用，timer.C 不可用
! After 等待时间结束，然后发送当前的时间。它相当于 NewTimer(d).C；底层的 timer 不会被 GC 收集
*/

func TestTimer(t *testing.T) {
	log(time.Now())
	t.Run("Sleep", func(t *testing.T) {
		time.Sleep(100 * time.Millisecond)
		log("sleep 100 milliseconds")
	})

	t.Run("Timer", func(t *testing.T) {
		now := time.Now()
		timer := time.NewTimer(100 * time.Millisecond)

		logfln("timer passed %d milliseconds", (<-timer.C).Sub(now).Milliseconds()) // 耗尽 timer.C

		timer.Reset(1)
		var wg sync.WaitGroup
		done := make(chan bool)
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-timer.C: // 阻塞
				log("timer fired")
			case <-done:
				log("timer stopped")
			}
		}()

		// <-time.After(10)  // 将阻塞 C
		if !timer.Stop() {
			<-timer.C // 耗尽 C
		} else {
			done <- true //
		}
		wg.Wait()
	})

	t.Run("AfterFunc", func(t *testing.T) {
		var wg sync.WaitGroup
		i := 10
		var afterFunc = func(d time.Duration, f func()) {
			wg.Add(1)
			time.AfterFunc(d, f)
		}

		var f func()
		f = func() {
			defer wg.Done()
			i--
			if i >= 0 {
				afterFunc(time.Millisecond*time.Duration(i), f)
			}
			log("call f")
		}
		afterFunc(0, f)
		wg.Wait()
	})

	t.Run("After", func(t *testing.T) {
		// 提前终止 AfterFunc
		timer := time.AfterFunc(3*time.Second, func() { log("Hello") })
		log(<-time.After(100 * time.Millisecond))
		log(timer.Stop())
	})
}

/* Tick
! time.Ticker 持有一个通道，该通道每隔一段时间发送一次当前时间
	Reset 重新设置间隔时间 duration
	Stop 停止 ticker 但不会关闭通道
! NewTicker 构造一个 Ticker
! Tick 包装一个 NewTicker(d).C；底层的 ticker 不会被 GC 收集，因此如果没有关闭 C 的方法，ticker 将泄露
*/

func TestTicker(t *testing.T) {
	ticker := examples.TickFunc(100*time.Millisecond, func() {
		log("do something")
	})

	time.Sleep(1 * time.Second)
	ticker.Reset(200*time.Millisecond, func() {
		log("do other things")
	})
	time.Sleep(1 * time.Second)
	ticker.Stop()

	// reset after stop
	ticker.Reset(100*time.Millisecond, func() {
		log("reset ticker")
	})
	time.Sleep(1 * time.Second)
	ticker.Stop()
}

/* Format
! Time.Parse-Functions
	String 返回 t 的格式化字符串的时间文本
	GoString 实现 fmt.GoStringer, %#v
	Format 返回给定 layout 布局的格式化字符串的时间文本
	AppendFormat 类似于 Format, 但将文本附加到 b 并返回扩展的缓冲区
! Parse 解析格式化字符串并返回它所表示的时间
! ParseInLocation 类似于 Parse；但若没有时区信息时，将依照给定 location 进行解释
! ParseDuration 解析一段持续时间的字符串
*/

func TestTimeFormat(t *testing.T) {
	format := func(t time.Time, layouts []struct {
		name   string
		layout string
	}) {
		logfln("%q Format", t)
		for _, l := range layouts {
			logfln("%-11s : %s", l.name, t.Format(l.layout))
		}
	}
	format(time.Now(), layouts)
}

func TestParse(t *testing.T) {
	parse := func(tm string, layouts []struct {
		name   string
		layout string
	}) {
		logfln("%s Parse", tm)
		for _, l := range layouts {
			t, err := time.Parse(l.layout, tm)
			if err == nil {
				logfln("%-11s : %v", l.name, t)
			}
		}
		log(" ")
	}
	for _, tm := range tms {
		parse(tm, layouts)
	}
}

func TestParseLocation(t *testing.T) {
	parse := func(tm string, loc *time.Location, layouts []struct {
		name   string
		layout string
	}) {
		for _, l := range layouts {
			t, err := time.ParseInLocation(l.layout, tm, loc)
			if err == nil {
				logfln("%-11s : %-30v | %s", l.name, t, loc)
			}
		}
	}
	for _, ts := range tms {
		logfln("%s ParseInLocation", ts)
		for _, l := range locations {
			if loc, err := time.LoadLocation(l); err == nil {
				parse(ts, loc, layouts)
			}
		}
		log(" ")
	}
}

func TestParseDuration(t *testing.T) {
	for _, s := range parseDuration {
		if d, err := time.ParseDuration(s); err == nil {
			logfln("%30s >> %q", s, d)
		}
	}
}
