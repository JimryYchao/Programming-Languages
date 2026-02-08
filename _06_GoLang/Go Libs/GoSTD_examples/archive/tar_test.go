package gostd_archive

/* 包 tar 提供了对 tar 压缩包的读写功能
! 核心功能：
! - Writer: 创建 tar 压缩文件
! - Reader: 读取 tar 压缩文件
! - Header: 表示 tar 文件的头信息
*/

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TarEntry 表示 tar 压缩包中的一个条目
type TarEntry struct {
	Name    string
	Content string
	Mode    int64
	ModTime time.Time
	IsDir   bool
}

var tarEntries = []TarEntry{
	{Name: "readme.txt", Content: "This is a readme file.", Mode: 0644, ModTime: time.Now()},
	{Name: "src/main.go", Content: "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}", Mode: 0644, ModTime: time.Now()},
	{Name: "src/utils/helper.go", Content: "package utils\n\nfunc Helper() {}", Mode: 0644, ModTime: time.Now()},
	{Name: "docs/guide.md", Content: "# Guide\n\nThis is a guide.", Mode: 0644, ModTime: time.Now()},
	{Name: "config/app.json", Content: `{"name": "MyApp", "version": "1.0.0"}`, Mode: 0644, ModTime: time.Now()},
}

// ! TarWriter 创建 tar 压缩包数据
func TarWriter(entries []TarEntry) ([]byte, error) {
	var buf bytes.Buffer
	w := tar.NewWriter(&buf)
	for _, entry := range entries {
		// 目录 header
		dir := filepath.Dir(entry.Name)
		if dir != "." && dir != "/" {
			dirHeader := &tar.Header{
				Name:     dir + "/",
				Mode:     0755,
				ModTime:  entry.ModTime,
				Typeflag: tar.TypeDir,
			}
			if err := w.WriteHeader(dirHeader); err != nil {
				return nil, fmt.Errorf("error writing dir header for %s: %w", dir, err)
			}
		}
		// 文件 header
		header := &tar.Header{
			Name:     entry.Name,
			Size:     int64(len(entry.Content)),
			Mode:     entry.Mode,
			ModTime:  entry.ModTime,
			Typeflag: tar.TypeReg,
		}
		if entry.IsDir {
			header.Typeflag = tar.TypeDir
			header.Size = 0
		}

		if err := w.WriteHeader(header); err != nil {
			return nil, fmt.Errorf("error writing header for %s: %w", entry.Name, err)
		}
		if !entry.IsDir && len(entry.Content) > 0 {
			if _, err := w.Write([]byte(entry.Content)); err != nil {
				return nil, fmt.Errorf("error writing content for %s: %w", entry.Name, err)
			}
		}
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("error closing writer: %w", err)
	}
	return buf.Bytes(), nil
}

// ! TarReader 读取 tar 归档数据
func TarReader(tarData []byte) ([]TarEntry, error) {
	if len(tarData) == 0 {
		return nil, fmt.Errorf("tar data is empty")
	}
	var entries []TarEntry
	r := tar.NewReader(bytes.NewReader(tarData))
	for {
		header, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading header: %w", err)
		}
		entry := TarEntry{
			Name:    header.Name,
			Mode:    header.Mode,
			ModTime: header.ModTime,
			IsDir:   header.Typeflag == tar.TypeDir,
		}
		if !entry.IsDir && header.Size > 0 {
			content, err := io.ReadAll(r)
			if err != nil {
				return nil, fmt.Errorf("error reading content for %s: %w", header.Name, err)
			}
			entry.Content = string(content)
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// ! CreateTarFile 创建实际的 tar 文件
func CreateTarFile(filename string, entries []TarEntry) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	tarData, err := TarWriter(entries)
	if err != nil {
		return err
	}
	if _, err := file.Write(tarData); err != nil {
		return err
	}
	return nil
}

// ! ExtractTarFile 从 tar 文件中提取内容
func ExtractTarFile(filename string) ([]TarEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tarData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return TarReader(tarData)
}

// ! 创建实际的 tar 文件并读取
// ? go test -v -run=TestTarFile
func TestTarFile(t *testing.T) {
	filename := "test.tar"
	// 创建 tar 文件
	if err := CreateTarFile(filename, tarEntries); err != nil {
		t.Fatalf("Error creating tar file: %v", err)
	}
	// 读取 tar 文件
	entries, err := ExtractTarFile(filename)
	if err != nil {
		t.Fatalf("Error extracting tar file: %v", err)
	}
	// 打印结果
	fmt.Printf("\nTar file '%s' created and verified successfully!\n", filename)
	fmt.Printf("Total entries: %d\n", len(entries))
	for _, entry := range entries {
		if entry.IsDir {
			fmt.Printf("  [DIR]  %s\n", entry.Name)
		} else {
			fmt.Printf("  [FILE] %s (%d bytes)\n", entry.Name, len(entry.Content))
		}
	}
	// 清理测试文件
	if err := os.Remove(filename); err != nil {
		t.Logf("Warning: could not remove test file: %v", err)
	}
}
