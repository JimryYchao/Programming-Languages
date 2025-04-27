package gostd

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"
)

/*
! MapFS (map[string]*MapFile) 是内存中的一个用于测试的简单文件系统，表示为从路径名（Open 的参数）到它们所表示的文件或目录信息的映射
	文件系统操作直接从映射中读取，可以根据需要通过编辑映射来更改文件系统；
	文件系统操作不能与映射的更改并发运行；打开或读取一个目录需要迭代整个映射；
	methods：Glob, Open, ReadDir, ReadFile, Stat, Sub
! MapFile 描述为 MapFS 中的单个文件
	Data 	文件内容
	Mode 	info.Mode
	ModTime info.ModTime
	Sys		info.Sys
! TestFS 测试文件系统的实现，它在 fsys 中遍历整个文件树，打开并检查每个文件是否正确；fsys 的内容不能与 TestFS 同时更改。
*/

func TestMapFS(t *testing.T) {
	fsys, paths := getTmpMapFS(t, 3)

	if err := fstest.TestFS(fsys, paths...); err != nil {
		checkErr(err)
	}
}

func getTmpMapFS(t *testing.T, fsDepth int) (*readFS, []string) {
	mFs := make(readFS)
	paths := make([]string, 32)
	i := 0
	fsys, _ := getTmpFsys(t, fsDepth)
	newMapFile := func(path string, d fs.DirEntry) *fstest.MapFile {
		if i >= len(paths) {
			new := make([]string, len(paths)*2)
			copy(new, paths)
		}
		paths[i] = path
		i++

		info, _ := d.Info()
		var data []byte
		if !d.IsDir() {
			var w bytes.Buffer
			if f, err := fsys.Open(path); err == nil {
				if n, err := w.ReadFrom(f); err == nil {
					data = make([]byte, w.Len())
					copy(data, w.Bytes()[0:n])
				}
				f.Close()
			}
		}
		return &fstest.MapFile{Data: data, Mode: info.Mode(), ModTime: info.ModTime(), Sys: info.Sys()}
	}
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}
		mFs[path] = newMapFile(path, d)
		return nil
	})
	return &mFs, paths[:i]
}

type readFS fstest.MapFS

func (fsys *readFS) Open(name string) (fs.File, error) {
	f, ok := (*fsys)[name]
	if ok {
		if f.Mode.IsDir() {
			logfln("read dir(%s)", name)
		} else {
			logfln("read file(%s): %s", name, f.Data)
		}
	}
	return (fstest.MapFS(*fsys).Open(name))
}

func getTmpFsys(t *testing.T, fsDepth int) (fsys fs.FS, err error) {
	root := t.TempDir()
	firstDir := root
	for i := fsDepth; i >= 0; i-- {
		firstDir, err = mkTmpDir(firstDir, i)
	}
	return os.DirFS(root), err
}
func mkTmpDir(dir string, n int) (string, error) {
	i := n
	if n > 0 {
		n--
		if _, err := mkTmpDir(dir, n); err != nil {
			return "", err
		}
	}
	f, _ := os.CreateTemp(dir, fmt.Sprintf("f%d_", i))
	f.WriteString(f.Name())
	f.Close()
	return os.MkdirTemp(dir, fmt.Sprintf("d%d_", i))
}
