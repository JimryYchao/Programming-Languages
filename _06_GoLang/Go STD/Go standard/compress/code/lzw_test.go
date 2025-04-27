package gostd

import (
	"compress/lzw"
	"io"
	"os"
	"testing"
)

/*
! NewReader 创建一个 zlib 的 io.ReadCloser，读取完成时调用方有责任调用其 Close；
	用于文字代码的位数 litWidth 必须在 [2,8] 范围内，它必须等于压缩期间使用的 litWidth
	ReadCloser 的底层类型 lzw.Reader: Close, Reset, Read
! NewWriter 创建一个新的 io.WriteCloser, 写入完成时调用方有责任调用其 Close；
	用于文字代码的位数 litWidth 必须在 [2，8] 范围内，通常为 8。输入字节必须小于 1<<litWidth
	WriteCloser 的底层类型是 lzw.Writer: Close, Reset, Read
! Order 指定 LZW 数据流中的位排序
	LSB 表示最低有效位优先，在 GIF 文件格式中使用
	MSB 表示最高有效位优先，在 TIFF 和 PDF 文件格式中使用
*/

// ? Compress
func TestLzwCompress(t *testing.T) {
	for _, f := range fileNames {
		lzwCompress(t, f)
	}
}

func lzwCompress(t *testing.T, file string) {
	os.Remove("testdata/" + file + ".lzw")
	f, err := os.Open("testdata/" + file)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	out, _ := os.Create("testdata/" + file + ".lzw")
	defer out.Close()
	zw := lzw.NewWriter(out, lzw.MSB, 8)
	defer zw.Close()
	if _, err = io.Copy(zw, f); err != nil {
		t.Fatal(err)
	}
}

// ? Decompress
func TestLzwDeCompress(t *testing.T) {
	lzwDeCompress(t, "compress.pdf.lzw")
	lzwDeCompress(t, "Isaac.Newton-Opticks.txt.lzw")
	lzwDeCompress(t, "e.txt.lzw")
}

func lzwDeCompress(t *testing.T, file string) {
	f, err := os.Open("testdata/" + file)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zw := lzw.NewReader(f, lzw.MSB, 8)
	defer zw.Close()
	if n, err := io.Copy(io.Discard, f); err != nil {
		t.Fatal(err)
	} else {
		_logfln("decompress %d bytes", n)
	}
}
