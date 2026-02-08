package gostd_testing

/* 包 fstest 提供了文件系统测试工具
! fstest 包主要用于测试 io/fs.FS 接口的实现
! 核心功能：
! - MapFS: 内存文件系统实现
! - TestFS: 测试文件系统
*/

import (
	"fmt"
	"io/fs"
	"testing"
	"testing/fstest"
)

// ! 使用 MapFS 创建内存文件系统
// ? go test -v -run=ExampleMapFS
func ExampleMapFS() {
	// 创建一个内存文件系统
	memfs := fstest.MapFS{
		"hello.txt": {
			Data: []byte("Hello, World!"),
		},
		"dir/": {},
		"dir/file.txt": {
			Data: []byte("File in directory"),
		},
	}
	// 读取文件内容
	data, err := fs.ReadFile(memfs, "hello.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Println("Content of hello.txt:", string(data))

	// 读取目录内容
	entries, err := fs.ReadDir(memfs, "dir")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	fmt.Println("Entries in dir:")
	for _, entry := range entries {
		if entry.Name() != "." {
			fmt.Printf("  %s (is dir: %v)\n", entry.Name(), entry.IsDir())
		}
	}

	// output:
	// Content of hello.txt: Hello, World!
	// Entries in dir:
	//   file.txt (is dir: false)
}

// ! 使用 TestFS 测试文件系统
// ? go test -v -run=TestFS
func TestFS(t *testing.T) {
	// 创建测试文件系统
	testFS := fstest.MapFS{
		"hello.txt": {
			Data: []byte("Hello, World!"),
		},
		"file1.txt": {
			Data: []byte("Content 1"),
		},
		"subdir/file2.txt": {
			Data: []byte("Content 2"),
		},
	}

	// 测试文件系统接口，至少包含 expected ...
	fstest.TestFS(testFS, "file1.txt", "subdir/file2.txt")
}
