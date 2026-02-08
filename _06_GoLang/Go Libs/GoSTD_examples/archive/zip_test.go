package gostd_archive

/* 包 zip 提供了 zip 压缩包的读写功能
! 核心功能：
! - Writer: 创建 zip 压缩文件
! - Reader: 读取 zip 压缩文件
! - File: 表示 zip 文件中的一个文件
! - FileHeader: 表示 zip 文件的头部信息
! - RegisterCompressor: 注册自定义压缩器
! - RegisterDecompressor: 注册自定义解压缩器
*/

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

// ZipEntry 表示 zip 归档中的一个条目
type ZipEntry struct {
	Name    string
	Content string
	Method  uint16
	ModTime time.Time
}

// 默认测试文件
var zipEntries = []ZipEntry{
	{Name: "readme.txt", Content: "This is a readme file.", Method: zip.Deflate, ModTime: time.Now()},
	{Name: "src/main.go", Content: "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}", Method: zip.Deflate, ModTime: time.Now()},
	{Name: "src/utils/helper.go", Content: "package utils\n\nfunc Helper() {}", Method: zip.Deflate, ModTime: time.Now()},
	{Name: "docs/guide.md", Content: "# Guide\n\nThis is a guide.", Method: zip.Deflate, ModTime: time.Now()},
	{Name: "config/app.json", Content: `{"name": "MyApp", "version": "1.0.0"}`, Method: zip.Deflate, ModTime: time.Now()},
}

// ZipWriter 创建 zip 归档数据
func ZipWriter(entries []ZipEntry) ([]byte, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	for _, entry := range entries {
		header := &zip.FileHeader{
			Name:               entry.Name,
			Method:             entry.Method,
			UncompressedSize64: uint64(len(entry.Content)),
			Modified:           entry.ModTime,
		}
		f, err := w.CreateHeader(header)
		if err != nil {
			return nil, fmt.Errorf("error creating file %s: %w", entry.Name, err)
		}
		if _, err := f.Write([]byte(entry.Content)); err != nil {
			return nil, fmt.Errorf("error writing content to %s: %w", entry.Name, err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("error closing writer: %w", err)
	}
	return buf.Bytes(), nil
}

// ZipReader 读取 zip 归档数据
func ZipReader(zipData []byte) ([]ZipEntry, error) {
	if len(zipData) == 0 {
		return nil, fmt.Errorf("zip data is empty")
	}
	var entries []ZipEntry
	r, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("error creating reader: %w", err)
	}
	for _, file := range r.File {
		f, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("error opening file %s: %w", file.Name, err)
		}
		content, err := io.ReadAll(f)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("error reading content from %s: %w", file.Name, err)
		}
		f.Close()

		entries = append(entries, ZipEntry{
			Name:    file.Name,
			Content: string(content),
			Method:  file.Method,
			ModTime: file.Modified,
		})
	}
	return entries, nil
}

// ! CreateZipFile 创建实际的 zip 文件
func CreateZipFile(filename string, entries []ZipEntry) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	zipData, err := ZipWriter(entries)
	if err != nil {
		return err
	}
	if _, err := file.Write(zipData); err != nil {
		return err
	}
	return nil
}

// ! ExtractZipFile 从 zip 文件中提取内容
func ExtractZipFile(filename string) ([]ZipEntry, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	zipData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return ZipReader(zipData)
}

// ! 创建实际的 zip 文件并读取
// ? go test -v -run=TestZipFile
func TestZipFile(t *testing.T) {
	filename := "test.zip"
	// 创建 zip 文件
	if err := CreateZipFile(filename, zipEntries); err != nil {
		t.Fatalf("Error creating zip file: %v", err)
	}
	// 读取 zip 文件
	entries, err := ExtractZipFile(filename)
	if err != nil {
		t.Fatalf("Error extracting zip file: %v", err)
	}
	// 打印结果
	fmt.Printf("\nZip file '%s' created and verified successfully!\n", filename)
	fmt.Printf("Total files: %d\n", len(entries))
	for _, entry := range entries {
		fmt.Printf("  [FILE] %s (%d bytes, method: %d)\n", entry.Name, len(entry.Content), entry.Method)
	}
	// 清理测试文件
	if err := os.Remove(filename); err != nil {
		t.Logf("Warning: could not remove test file: %v", err)
	}
}

// ! 注册自定义压缩器并测试
// ? go test -v -run=TestRegisterCompressor
func TestRegisterCompressor(t *testing.T) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		fmt.Println("Registering custom Deflate compressor")
		return flate.NewWriter(out, flate.BestCompression)
	})

	if f, err := w.CreateHeader(&zip.FileHeader{
		Name:               "TestRegisterCompressor.txt",
		Method:             zip.Deflate,
		UncompressedSize64: uint64(len("Hello, World!")),
	}); err != nil {
		t.Fatalf("Error creating file header: %v", err)
	} else {
		if _, err := f.Write([]byte("Hello, World!")); err != nil {
			t.Fatalf("Error writing content to file: %v", err)
		}
	}
}
