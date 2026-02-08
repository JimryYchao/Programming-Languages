package gostd_testing

/* 包 slogtest 提供了结构化日志测试工具, 主要用于测试 slog.Handler 的实现。
! 核心功能：
! - TestHandler: 测试 slog.Handler 的实现是否符合规范
*/

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
	"testing/slogtest"
)

// ! 使用 TestHandler 测试 slog.Handler 实现
// ? go test -v -run=TestSlogHandler
func TestSlogHandler(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewJSONHandler(&buf, nil)
	result := func() []map[string]any {
		var ms []map[string]any
		for _, line := range bytes.Split(buf.Bytes(), []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			var m map[string]any
			if err := json.Unmarshal(line, &m); err != nil {
				t.Fatal(err)
			}
			ms = append(ms, m)
		}
		return ms
	}

	// 测试处理器
	if err := slogtest.TestHandler(handler, result); err != nil {
		t.Fatalf("Handler test failed: %v", err)
	}
}

// ! 使用 Run 测试 slog.Handler 实现
// ? go test -v -run=TestSlogRun
func TestSlogRun(t *testing.T) {
	var buf bytes.Buffer
	handlerFunc := func(t *testing.T) slog.Handler {
		buf.Reset() // 每次测试都重置缓冲区
		return slog.NewJSONHandler(&buf, nil)
	}

	// 定义结果函数，用于解析日志输出
	resultFunc := func(t *testing.T) map[string]any {
		var m map[string]any
		lines := bytes.Split(buf.Bytes(), []byte{'\n'})
		for _, line := range lines {
			if len(line) > 0 {
				if err := json.Unmarshal(line, &m); err != nil {
					t.Fatal(err)
				}
				break
			}
		}
		return m
	}

	slogtest.Run(t, handlerFunc, resultFunc)
}
