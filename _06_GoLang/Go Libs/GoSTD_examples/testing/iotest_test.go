package gostd_testing

/* 包 iotest 提供了 I/O 测试工具
! iotest 包主要用于测试 I/O 相关的代码
! 核心功能：
! - OneByteReader: 每次只读一个字节的 Reader
! - HalfReader: 每次只读取一半数据的 Reader
! - TimeoutReader: 模拟超时的 Reader
! - ErrReader: 模拟错误的 Reader
! - TruncateReader: 模拟截断数据的 Reader
! - WriteLogger: 记录写入操作的 Writer
*/

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"testing/iotest"
)

// ! 使用 HalfReader 测试
// ? go test -v -run=TestHalfReader
func TestHalfReader(t *testing.T) {
	data := []byte("Hello, World!")
	r := iotest.HalfReader(bytes.NewReader(data))
	// 读取数据（多次读取）
	buf := make([]byte, 0, len(data))
	for {
		b := make([]byte, len(data))
		if n, err := r.Read(b); err != nil { // 读取一半
			if err == io.EOF {
				break
			}
			t.Fatalf("Error reading: %v", err)
		} else if n > 0 {
			t.Log(string(b[:n]))
			buf = append(buf, b[:n]...)
		}
	}
	if !bytes.Equal(buf, data) {
		t.Fatalf("Expected %q, got %q", data, buf)
	}
}

// ! 使用 ErrReader 测试
// ? go test -v -run=TestErrReader
func TestErrReader(t *testing.T) {
	r := iotest.ErrReader(errors.New("custom error"))
	n, err := r.Read(nil)
	fmt.Printf("n:   %d\nerr: %q\n", n, err)
}
