package gostd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"testing"
)

var (
	dirPath  = "."
	fileName = "fs_testing.file"
)

/*
! fs.FS 提供对分层文件系统的访问。`FS` 接口是文件系统所需的最低实现。
	Open 以打开 `name` 指定的文件或目录，并返回一个 `File`。发生错误时返回一个 `*PathError` 类型的 `error`。
! fs.PathError 记录错误以及导致错误的操作和文件路径。
	Error 拼接 `e.Op + e.Path : e.Err`。
	Unwrap 返回 `e.Err`。
	Timeout 报告该错误是否是由于操作超时引发。
! ValidPath 报告给定的路径名格式是否有效，有效时可用于 `(FS).Open` 调用。它不检查路径是否存在。
! fs.File 提供对单个文件的访问。`File` 接口是文件所需的最低实现。目录文件也应实现 `fs.ReadDirFile`。
	Stat 返回该文件的 `FileInfo`。
	Read & Close 实现基本的 `io.Reader` 和 `io.Closer`。也可能实现 `io.ReaderAt` 和 `io.Seeker`。
! fs.FileInfo 描述文件信息。由 `Stat` 返回。
	Name 报告文件的名称。
	Size 报告常规文件的字节长度。其他文件类型依赖于系统。
	Mode 报告文件模式位 `FileMode`。
	ModTime 报告上次修改时间。
	IsDir 报告 `Mode().IsDir()`。
	Sys 报告底层数据源。

! fs.FileMode 表示文件的模式位和权限位。这些位在所有的系统上都有相同的定义。定义文件模式位是 `FileMode` 的最高有效位；最低 9 位是标准的 Unix 权限位。
	IsDir 描述是否为目录。
	IsRegular 描述是否为常规文件，即没有被设置模式类型位。
	Perm 返回 `fileMode` 的 Unix 权限位。
	String 字符串化 `FileMode` 的值。
	Type 返回 `fileMode & fs.ModeType` 中的类型位。

! FormatFileInfo 返回 `FileInfo` 格式化的字符串。
*/

func TestValidPath(t *testing.T) {
	paths := []string{".", "x", "x/y", "", "..", "/", "x/", "/x", "x/y/", "/x/y",
		"./", "./x", "x/.", "x/./y", "../", "../x", "x/..", "x/../y", "x//y", `x\`, `x\y`, `x:y`, `\x`}

	for _, p := range paths {
		ok := fs.ValidPath(p)
		logf("ValidPath(%q) = %v", p, ok)
	}
}

func TestOpenFile(t *testing.T) {
	if fs.ValidPath(dirPath) { // check path format
		if file, err := os.DirFS(dirPath).Open(fileName); err == nil {
			file.Read(buf) // read file
			defer file.Close()
			logf("Read file : %s", buf)
		} else { // check err
			if pathErr, ok := err.(*fs.PathError); ok {
				checkErr(pathErr)
			} else {
				t.Fatal(err)
			}
		}
	}
}

func TestOpenFileWithErr(t *testing.T) {
	var perr *fs.PathError
	_, err := os.DirFS(t.TempDir()).Open("non-existent")
	if errors.As(err, &perr) {
		checkErr(perr)
	}
}

func TestFileInfo(t *testing.T) {
	file, _ := os.DirFS(".").Open(fileName)
	defer file.Close()
	info, _ := file.Stat()
	logf("Info of file : %s", fs.FormatFileInfo(info))
}

func TestFileMode(t *testing.T) {
	sys, _ := getTmpFsys(t, 3)
	entries, _ := fs.ReadDir(sys, ".")
	for _, e := range entries {
		fileMode := e.Type()
		logf("name: %s, isRegular: %v, perm: %v, type: %v, isDir: %v",
			e.Name(), fileMode.IsRegular(), fileMode.Perm(), fileMode.Type(), fileMode.IsDir())
	}
}

/* read dir
! fs.ReadDirFS 是由文件系统实现的接口，该文件系统提供 `ReadDir` 的优化实现。
! ReadDir 读取指定 `name` 的目录并返回按文件名排序的 `DirEntry` 条目列表。
! fs.DirEntry 是从目录中读取的条目。如使用 `ReadDir` 或 `ReadDirFile.ReadDir` 读取。
	Name 返回该条目的名称描述，可能是文件或是子目录名称。
	IsDir 报告该条目是否是目录。
	Type 返回该条目的类型位。
	Info 返回描述该条目的 `FileInfo`。
! FormatDirEntry 返回 `DirEntry` 格式化的字符串。
*/

func TestFS_ReadDirFS(t *testing.T) {
	if fsys, err := getTmpFsys(t, 3); err != nil {
		t.Fatal(err)
	} else {
		if dirfs, ok := fsys.(fs.ReadDirFS); ok {
			dirs, _ := dirfs.ReadDir(".")
			iterateDirEntries(dirs)
		}
	}
}

func TestReadDirOnly(t *testing.T) {
	fsys, err := getTmpFsys(t, 3)
	if err != nil {
		t.Fatal(err)
	}
	// 读取一个临时的文件树系统
	if dirs, err := fs.ReadDir(fsys, "."); err == nil {
		for _, dir := range dirs {
			if dir.IsDir() {
				info, _ := dir.Info()
				logf("- d: `%s`, type: `%s`, info: `%s`", dir.Name(), dir.Type(), fs.FormatFileInfo(info))
			}
		}
	}
}

func iterateDirEntries(dirs []fs.DirEntry) {
	if len(dirs) == 0 {
		return
	}
	for _, d := range dirs {
		logf(fs.FormatDirEntry(d))
	}
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

/*	read file
! fs.ReadFileFS 是由文件系统实现的接口，它提供 `ReadFile` 的优化实现。
	ReadFile 读取 `name` 文件并返回其内容，读取结束的 `io.EOF` 不被视为错误。
! ReadFile 从文件系统 `fsys` 中读取指定 `name` 的文件并返回其内容。
*/

func TestReadFile(t *testing.T) {
	fsys, _ := getTmpFsys(t, 3)
	dtris, _ := fs.ReadDir(fsys, ".")
	for _, d := range dtris {
		data, err := fs.ReadFile(fsys, d.Name()) // 读取非文件的目录文件将返回 err
		checkErr(err)                            //
		if err == nil {
			logf("Read file: %s", data)
		}
	}

	// read file err
	_, err := fs.ReadDir(fsys, "????")
	checkErr(err)
}

/*
! fs.StatFS 是一个带有 `Stat` 方法的文件系统。
! Stat 从文件系统 `fsys` 中返回描述 `name` 文件的 `FileInfo`。
! FileInfoToDirEntry 从 `FileInfo` 返回 `DirEntry`。
*/

func TestStatFS(t *testing.T) {
	fsys := os.DirFS(".")
	if statfs, ok := fsys.(fs.StatFS); ok {
		info, err := fs.Stat(statfs, "fs_testing.file")
		if err != nil {
			t.Fatal(err)
		}
		logf("FileInfo : %s", fs.FormatFileInfo(info))
		logf("DirEntry : %s", fs.FileInfoToDirEntry(info))
	}
}

/*
! fs.SubFS 是一个具有 `Sub` 方法的文件系统。
! Sub 返回一个对应于 `fsys` 系统下目录 `dir` 的子树文件系统。`Sub` 返回一个内置的 `SubFS`。
*/

func TestCheckFs(t *testing.T) {
	// 检查一个文件系统下的所有内容
	fsys, _ := getTmpFsys(t, 5)
	// fsys = os.DirFS(some path)   // 检查某个目录下的文件树
	dirTree, err := createTreeFromFS(fsys, ".")

	if err != nil {
		t.Fatal(err)
	}
	dirTree.walkDir(-1)
}

type dirTree struct {
	name      string
	subdirs   *[]dirTree
	files     *[]file
	dirCount  int
	fileCount int
}

// 0 表示当前目录，depth 表示文件深度，-1 时全部遍历
func (t *dirTree) walkDir(depth int) {
	walkDirHelper(t, "", depth+1)
}

func walkDirHelper(t *dirTree, indent string, depth int) {
	indent += "--"
	if t.fileCount > 0 {
		for _, f := range *t.files {
			logf(indent+" f: %s", f.name)
		}
	}
	if t.dirCount > 0 {
		depth--
		for _, dir := range *t.subdirs {
			logf(indent+" d: %s", dir.name)
			if depth != 0 {
				walkDirHelper(&dir, indent, depth)
			}
		}
	}
}

type file struct {
	name string
}
type subFSEntry struct {
	name  string
	subFS fs.SubFS
}

func checkSub(fsys fs.FS) ([]subFSEntry, []file, error) {
	dirs, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, nil, err
	}
	subs, dirCount := make([]subFSEntry, len(dirs)), 0
	files, fCount := make([]file, len(dirs)), 0
	var sub fs.FS
	for _, d := range dirs {
		if d.IsDir() {
			sub, err = fs.Sub(fsys, d.Name())
			subs[dirCount] = subFSEntry{d.Name(), sub.(fs.SubFS)}
			dirCount++
			if err != nil {
				return nil, nil, err
			}
		} else {
			files[fCount].name = d.Name()
			fCount++
		}
	}
	return subs[0:dirCount], files[0:fCount], nil
}

func createTreeFromFS(fsys fs.FS, name string) (*dirTree, error) {
	var tree dirTree
	tree.name = name
	subFSEntries, files, err := checkSub(fsys)
	if err != nil {
		return nil, err
	}
	if len(files) > 0 {
		tree.files = &files
		tree.fileCount = len(files)
	}
	if len(subFSEntries) == 0 {
		return &tree, nil
	}
	subT, subTrees := new(dirTree), make([]dirTree, len(subFSEntries))
	for _, _subEntry := range subFSEntries {
		subT, err = createTreeFromFS(_subEntry.subFS, _subEntry.name)
		if err != nil {
			return nil, err
		}
		subTrees[tree.dirCount] = *subT
		tree.dirCount++
	}
	subTrees = subTrees[0:tree.dirCount]
	tree.subdirs = &subTrees
	return &tree, nil
}

/*
! fs.GlobFs 是一个带有 `Glob` 方法的文件系统。
! Glob 返回指定文件系统 `fsys` 中所有匹配 `pattern` 的文件名称。
*/

func TestGlob(t *testing.T) {
	fsys := os.DirFS(".")
	if _, ok := fsys.(fs.GlobFS); ok {
		logf("os.DirFS return a fs.GlobFS")
	}

	matchs, err := fs.Glob(fsys, "*.go")
	if err != nil {
		checkErr(err)
		return
	}
	for _, m := range matchs {
		if file, err := fsys.Open(m); err != nil || file == nil {
			checkErr(err)
			continue
		} else {
			if data, err := io.ReadAll(file); err != nil {
				logf("Read file[%s] failed: %s", m, err)
			} else {
				logf("Read file[%s] successfully, and read %d bytes", m, len(data))
			}
			file.Close()
		}
	}
}

/*
! WalkDir 遍历以 `root` 为根的文件树，为树中的每个文件或目录（包括 `root`）调用 `WalkDirFunc fn`。
! fs.WalkDirFunc 是 `WalkDir` 调用的函数类型，用于访问每个文件或目录。`type WalkDirFunc func(path string, d DirEntry, err error) error`
*/

func TestWalkDir(t *testing.T) {
	fsys, _ := getTmpFsys(t, 3)
	walkHelper := func(path string, d fs.DirEntry, err error) error {
		fmt.Printf("walk to %s, ", d.Name())
		if err != nil {
			return err
		}
		if d.IsDir() {
			// do something in dir: d
			logf("do nothing with dir[%s]", path)
		} else {
			if file, e := fsys.Open(path); e != nil {
				return e
			} else {
				fmt.Print("read: ")
				readToStdout(file)
				file.Close()
			}
		}
		return nil
	}
	err := fs.WalkDir(fsys, ".", walkHelper)
	checkErr(err)
}
func TestWalkGoBin(t *testing.T) {
	recordD, err1 := os.OpenFile("gosrcDirs.md", os.O_CREATE|os.O_WRONLY, 0777)
	if err1 != nil {
		return
	}
	defer recordD.Close()

	recordF, err2 := os.OpenFile("gosrcFiles.md", os.O_CREATE|os.O_WRONLY, 0777)
	if err2 != nil {
		return
	}
	defer recordF.Close()

	fileSystem := os.DirFS("C:\\_Programme Environment_\\_Go_\\src")

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if d.IsDir() {
			io.WriteString(recordD, "- "+path+"\n")
		} else {
			io.WriteString(recordF, "- "+path+"\n")
		}
		return nil
	})
}

/*
! fs.ReadDirFile 是一个目录文件，其目录中的条目可以用其 `ReadDir` 方法读取。每个目录文件都应该实现这个接口。对于非目录文件，`ReadDir` 返回错误。
	ReadDir 读取目录的内容并返回，`n > 0` 时最多按顺序返回 `n` 个 `DirEntry`，抵达末尾时返回 `io.EOF`。`n < 0` 读取全部内容，抵达末尾返回 `err = nil`。
*/
//? go test -v -run=^TestReadDirFile$
func TestReadDirFile(t *testing.T) {
	f, _ := os.Open(".")
	fsRdf, ok := newReadDirFile(f)
	if !ok {
		t.Fail()
	}
	errCheck := func(err error) {
		t.Helper()
		if err != nil {
			if err == io.EOF {
				checkErr(err)
				return
			}
			t.Fatal(err)
		}
	}
	walkEntries := func(ds []fs.DirEntry, err error, n int64) {
		logf("Read %d entries at most", n)
		errCheck(err)
		for _, v := range ds {
			logf("- " + v.Name())
		}
	}

	ds, err := fsRdf.ReadDir(1)
	walkEntries(ds, err, 1)

	ds, err = fsRdf.ReadDir(3)
	walkEntries(ds, err, 3)

	ds, err = fsRdf.ReadDir(-1) // read all
	walkEntries(ds, err, -1)

	ds, err = fsRdf.ReadDir(1<<63 - 1) // read EOF
	walkEntries(ds, err, 1<<63-1)
}

type readDirFile struct {
	f *os.File
}

func newReadDirFile(f *os.File) (readDirFile, bool) {
	if info, err := f.Stat(); err != nil || !info.IsDir() {
		return readDirFile{}, false
	} else {
		return readDirFile{f}, true
	}
}
func (f *readDirFile) ReadDir(n int) ([]fs.DirEntry, error) {
	return f.f.ReadDir(n)
}
