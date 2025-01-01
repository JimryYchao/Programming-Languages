package gostd

import (
	"compress/zlib"
	"io"
	"os"
	"testing"
)

/* reader
decompress
! NewReader 从 r 构造一个 ReaderCloser 解压缩器，可能会从 r 中读取更多的数据；返回的 io.ReadCloser 也实现了 Resetter
! NewReaderDict 类似于 NewReader，但使用预设的 dict 作为支持
compress
! NewWriter 返回一个在给定级别压缩数据的新 Writer;
! NewWriterLevel 类似于 NewWriter，但使用指定的压缩级别 [-2,9]
! NewWriterLevelDict 类似于 NewWriter 但使用预设的 dict, level, 写入 w 的压缩数据只能由使用相同 dict 初始化的 reader 解压缩
*/

// ? Compress
func TestZlibCompress(t *testing.T) {
	for _, f := range fileNames {
		zlibCompress(t, f)
	}
}

func zlibCompress(t *testing.T, file string) {
	os.Remove("testdata/" + file + ".zl")
	in, err := os.Open("testdata/" + file)
	if err != nil {
		t.Fatal(err)
	}
	defer in.Close()
	out, err := os.Create("testdata/" + file + ".zl")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	zw, _ := zlib.NewWriterLevel(out, zlib.DefaultCompression)

	if _, err = io.Copy(zw, in); err != nil {
		zw.Close()
		t.Fatal(err)
	}
	zw.Close()
}

// ? Decompress
func TestZlibReader(t *testing.T) {
	f, err := os.Open("testdata/e.txt.zl")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	zr, err := zlib.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	if bs, err := io.ReadAll(zr); err != nil {
		t.Fatal(err)
	} else {
		zr.Close()
		os.Stdout.Write(bs)
	}
}
