package gostd

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"testing/slogtest"
)

/*
! TestHandler 测试一个 slog.Handler, 当发现任何错误时, 将它们组合成一个 joinError 返回.
	TestHandler 将给定的 Handler 安装到 slog.Logger 上, 并对 Logger 的输出方法进行多次调用, Handler 的级别应至少为 InfoLevel
	它返回一个嵌套的 map[string]any, 对于每个 Logger 的输出方法都有一个; 映射的键值对应于 Handler 输出的键值; 输出的每个组表示为自己的嵌套映射 map[group]submap;
	标准键应使用 slog.TimeKey, slog.LevelKey, slog.MessageKey

	如果 Handler 输出 JSON，则使用 `map[string]any` 调用 encoding/json.Unmarshal 将创建正确的数据结构。
	如果 Handler 故意丢弃了一个被测试检查过的属性，那么 results 函数应该检查它是否不存在，并将它添加到它返回的 map 中。
*/

func TestTestHandler(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, nil)

	results := func() []map[string]any {
		var ms []map[string]any
		for _, line := range bytes.Split(buf.Bytes(), []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			var m map[string]any
			if err := json.Unmarshal(line, &m); err != nil {
				t.Fatal(err) // In a real test, use t.Fatal.
			}
			ms = append(ms, m)
		}
		return ms
	}
	err := slogtest.TestHandler(h, results)

	buf.WriteTo(os.Stdout) // check test output

	if err != nil {
		t.Fatal(err)
	}
}

/*
! Run 与 TestHandler 在相同的测试用例上使用 Handler, 但在每个子测试中运行每个用例
	对于每个子测试用例,它首先调用 newHandler 函数来构造一个 slog.Handler; 然后运行测试用例; 最后使用 result 获取结果
	测试用例失败, 它将调用 t.Error
*/

func TestRun(t *testing.T) {
	var buf bytes.Buffer

	newHandler := func(*testing.T) slog.Handler {
		buf.Reset()
		return slog.NewJSONHandler(&buf, nil)
	}
	result := func(t *testing.T) map[string]any {
		m := map[string]any{}
		if err := json.Unmarshal(buf.Bytes(), &m); err != nil {
			t.Fatal(err)
		}
		return m
	}

	slogtest.Run(t, newHandler, result)

	buf.WriteTo(os.Stdout)
}
