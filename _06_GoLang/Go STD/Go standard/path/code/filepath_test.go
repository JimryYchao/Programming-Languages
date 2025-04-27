package gostd

import (
	"io"
	"io/fs"
	"path/filepath"
	"testing"
)

/*
! Abs 返回 path 的绝对表示；不是绝对路径则与当前工作目录连接
! Base 返回 path 的最后一个元素。
! Clean 通过纯词法处理返回与 path 等效的最短路径名。
! Dir 返回 path 的最后一个元素（通常是 path 的目录）以外的所有元素
! EvalSymlinks 在计算任何符号链接后返回路径名
! Ext 返回 path 使用的文件扩展名
! FromSlash 返回将路径中的每个斜线（'/'）字符替换为分隔符的结果。多个斜线被多个分隔符替换。
! Glob 返回所有匹配 pattern 的文件的名称，如果没有匹配的文件，则返回 nil; 忽略文件系统错误，例如 I/O 错误（阅读目录）
! IsAbs 报告路径是否为绝对路径
! IsLocal 报告 path（仅使用词法分析）是否具有所有这些属性: 非绝对路径，非空，在 windows 上非保留名称，它位于计算路径所在目录的子树中
! Join 将任意数量的路径元素连接到单个路径中，并使用操作系统特定的分隔符将它们分隔开
! Match 报告 name 是否与 pattern 模式匹配; 在 Windows 上，禁用转义。而是将 \\ 视为路径分隔符。
! Rel 返回 targetpath 相对 basepath 的相对路径
! Split 在最后一个 Separator 之后立即拆分路径，将其分隔为目录和文件名组件
! SplitList 拆分由特定于操作系统的 ListSeparator 连接的路径列表，通常可以在 PATH 或 GOPATH 环境变量中找到
! ToSlash 返回将 path 中的每个分隔符替换为斜杠（“/”）字符的结果
! VolumeName 返回前导卷名。给定 C:\foo\bar 返回 C:。给定 \\host\share\foo 返回 \\host\share。在其他平台上，它返回 ""。
! Walk 遍历以 root 为根的文件树，为树中的每个文件或目录（包括 root）调用 fn; Walk 的效率不如 WalkDir
	WalkFunc 是 Walk 访问每个文件或目录所调用的函数的类型
! WalkDir 遍历以 root 为根的文件树，为树中的每个文件或目录（包括 root）调用 fn;
	与 io/fs.WalkDir 不同, WalkDir 调用 fn 时使用的路径使用了适合于操作系统的分隔符
*/

func TestFilepath(t *testing.T) {
	log("On Windows:")
	for _, p := range paths {
		logfln("path %s :", p)
		logfln("   > Base: %s", filepath.Base(p))
		logfln("   > Clean: %s", filepath.Clean(p))
		logfln("   > Dir: %s", filepath.Dir(p))
		logfln("   > Ext: %s", filepath.Ext(p))
		logfln("   > IsAbs: %t", filepath.IsAbs(p))
		dir, file := filepath.Split(p)
		logfln("   > Split: %s, %s", dir, file)
	}
}

func TestJoin(t *testing.T) {
	for i := 0; i < len(paths)-1; i++ {
		logfln("join [%s, %s] : %s", paths[i], paths[i+1], filepath.Join(paths[i], paths[i+1]))
	}
}

func TestMatch(t *testing.T) {
	log("On Windows: ")
	log(filepath.Match("/home/catch/*", "/home/catch/foo"))     // true
	log(filepath.Match("/home/catch/*", "/home/catch/foo/bar")) // true
	log(filepath.Match("/home/?opher", "/home/gopher"))         // true
	log(filepath.Match("/home/\\*", "/home/*"))                 // false
}

func TestSplitList(t *testing.T) {
	pathlist := "/bar/foo;/root;.;/dir/file.txt;foo/"

	log(filepath.SplitList(pathlist))
	// [/bar/foo /root . /dir/file.txt foo/]
}

func TestWalkDir(t *testing.T) {
	filepath.WalkDir(`..\`, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return io.EOF
		}
		if d.IsDir() {
			logfln("Dir : %s", filepath.Base(path))
		} else {
			logfln("File: %s", filepath.Base(path))
		}
		return nil
	})
}
