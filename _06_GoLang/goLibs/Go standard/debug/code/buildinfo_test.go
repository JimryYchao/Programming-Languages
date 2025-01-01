package gostd

import (
	"debug/buildinfo"
	"debug/pe"
	"os"
	"testing"
)

/*
! BuildInfo 为 runtime/debug.BuildInfo 别名
! Read 返回嵌入在 Go 二进制文件中的构建信息。
! ReadFile 返回从给定路径 Go 二进制文件的构建信息。
*/

var exePath = "fortest/fortest.exe"

func TestBuildInfo(t *testing.T) {
	f, err := os.Open(exePath)
	if err != nil {
		t.Fatal(err)
	}
	binfo, err := buildinfo.Read(f)
	if err != nil {
		t.Fatal(err)
	}

	log(binfo)

	bpath, _ := buildinfo.ReadFile(exePath)
	log(bpath)
}

func TestPE(t *testing.T) {
	f, err := pe.Open(exePath)
	if err != nil {
		t.Fatal(err)
	}
	lobs, err := f.ImportedLibraries()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range lobs {
		logfln(v)
	}
}
