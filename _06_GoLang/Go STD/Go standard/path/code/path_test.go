package gostd

import (
	"path"
	"testing"
)

/*
! Base 返回 path 的最后一个元素；删除尾随的斜杠。空路径返回 "."; 路径完全由斜杠组成时返回 "/"
! Clean 通过纯词法处理返回与 path 等效的最短路径名
! Dir 返回 path 的最后一个元素（通常是 path 的目录）以外的所有元素。
! Ext 返回 path 使用的文件扩展名。
! IsAbs 报告 path 是否为绝对路径。
! Join 将任意数量的路径元素连接到单个路径中，并使用斜杠将它们分隔开
! Match 报告 name 是否与 pattern 匹配: { term }
	term:
		*         	matches any sequence of non-/ characters
		?         	matches any single non-/ character
		[ character-range ] or [^ character-range ]   	character class (must be non-empty)
		c           matches character c (c != '*', '?', '\\', '[')
		\\c      	matches character c
	character-range:
		c           matches character c (c != '\\', '-', ']')
		\\c      	matches character c
		lo ~ hi   	matches character c for lo <= c <= hi
! Split 在最后一个斜杠之后立即拆分 path，将其分隔为目录和文件名组件。
*/

var paths = []string{
	"", ".", "/",
	"root", "root/", "/root", "/root/",
	"root/file.txt",
	"root/dir", "root/dir/",
	"file.txt",
	"/root/dir/dir2/file.txt",
	"root/*/",
	"../../root",
	"./dir", "..",
}

func TestPath(t *testing.T) {
	checkPath := func(p string) {
		logfln("path %s :", p)
		logfln("   > Base: %s", path.Base(p))
		logfln("   > Clean: %s", path.Clean(p))
		logfln("   > Dir: %s", path.Dir(p))
		logfln("   > Ext: %s", path.Ext(p))
		logfln("   > IsAbs: %t", path.IsAbs(p))
		dir, file := path.Split(p)
		logfln("   > Split: %s, %s", dir, file)
	}

	for _, p := range paths {
		checkPath(p)
	}
}
