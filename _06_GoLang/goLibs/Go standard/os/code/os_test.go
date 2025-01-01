package gostd

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

/*
! os.Stdin, Stdout, Stderr 对应标准输入，标准输出，标准错误
! os.Args 保存命令行参数，从程序名称开始
! Hostname 返回内核报告的主机名。
! Executable 返回启动当前进程的可执行文件的路径名。主要用例是查找相对于可执行文件的资源。
! Exit 以给定的状态代码退出当前程序。
! Getgid 返回调用方的数字组 ID。win 返回 -1
! Getegid 返回调用方的数字有效组 ID。win 返回 -1
! Getuid 返回调用方的数字用户 ID。
! Geteuid 返回调用方的数字有效用户 ID。win 返回 -1
! Getpid 返回调用方的进程 ID。
! Getppid 返回调用方父级的进程 ID。
! Getgroups 返回调用方所属组的数字 ID 列表。Win 返回 syscall.EWINDOWS err
! Getpagesize 返回基础系统的内存页大小。
*/

func TestArgs(t *testing.T) {
	log(os.Executable())
	sr := bytes.NewBuffer(nil)
	for _, arg := range os.Args {
		sr.WriteString(arg + "\n")
		io.Copy(os.Stdout, sr)
	}

	logfln("gid  = %d", os.Getgid())
	logfln("egid = %d", os.Getegid())
	logfln("uid = %d", os.Getuid())
	logfln("euid = %d", os.Geteuid())
	gps, e := os.Getgroups()
	logfln("groups = %v, %s", gps, e)
	logfln("pagesize = %d", os.Getpagesize())
	logfln("pid = %d", os.Getpid())
	logfln("ppid = %d", os.Getppid())
	root, _ := os.Getwd()
	logfln("wd = %s", root)
	host, _ := os.Hostname()
	logfln("Hostname = %s", host)
}

/*
! IsPathSeparator 报告 c
! IsExist 报告 err 中是否已知文件或目录已存在
! IsNotExist 报告 err 中是否已知文件或目录不存在
! IsPermission 报告 err 中是否已知权限被拒绝
! IsTimeout 报告 err 中是否发生超时错误
*/

func TestIsErr(t *testing.T) {
	log(os.IsExist(os.ErrExist))
	log(os.IsNotExist(os.ErrNotExist))
	log(os.IsPermission(os.ErrPermission))
}

/*
! Getwd 返回与当前目录对应的 root 路径名
! Mkdir 创建一个具有指定名称和权限位的新目录（在 umask 之前）
! MkdirAll 创建一个名为 path 的目录以及任何必要的父目录
! MkdirTemp 在 dir 目录中创建一个新的临时目录，并返回新目录的路径名; 当不再需要目录时，调用方有责任删除该目录。
! TempDir 返回用于临时文件的默认目录; win %TMP%、%TEMP%、%USERPROFILE%...
! UserCacheDir 返回用于特定于用户的缓存数据的默认根目录; win %LocalAppData%
! UserConfigDir 返回用于用户特定配置数据的默认根目录; win %AppData%
! UserHomeDir 返回当前用户的主目录; win %USERPROFILE%
! Remove 删除 name 文件或目录
! RemoveAll 删除路径及其包含的任何子项。
! Rename 重命名（或移动）oldpath 到 newpath；newpath 已存在并且不是目录，Rename 将替换它
! SameFile 报告 fi1 和 fi2 是否描述同一文件
! Truncate 截断更改 name 文件的大小
*/

func TestMkDirTmp(t *testing.T) {
	for range 10 {
		logsDir, err := os.MkdirTemp(os.TempDir(), "*-logs")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(logsDir)
	}
	// 提前删除任何匹配 *-logs 模式的目录或文件
	globPattern := filepath.Join(os.TempDir(), "*-logs")
	matches, err := filepath.Glob(globPattern)
	if err != nil {
		t.Fatalf("Failed to match %q: %v", globPattern, err)
	}
	for _, match := range matches {
		if err := os.RemoveAll(match); err != nil {
			t.Logf("Failed to remove %q: %v", match, err)
		} else {
			logfln("remove : %s", match)
		}
	}
}

func TestWindowDirs(t *testing.T) {
	logfln("TempDir : %s", os.TempDir())
	p, _ := os.UserCacheDir()
	logfln("UserCacheDir : %s", p)
	p, _ = os.UserConfigDir()
	logfln("UserConfigDir : %s", p)
	p, _ = os.UserHomeDir()
	logfln("UserHomeDir : %s", p)
}

/*
! Symlink 创建 newname 作为指向 oldname 的符号链接
! Link 将 newname 创建为指向 oldname 文件的硬链接
! Readlink 返回 name 符号链接的目标
*/

func Test(t *testing.T) {
	path := "tmp/rolfl/symexample"
	target := "symtarget.txt"
	os.MkdirAll(path, 0755)
	os.WriteFile(filepath.Join(path, "symtarget.txt"), []byte("Hello\n"), 0644)
	symlink := filepath.Join(path, "symlink")
	os.Symlink(target, symlink) // NO PERMISSION
}

/*
! Chdir 将当前工作目录更改为指定 dir
! Chmod 将 name 文件的模式改为指定的 mode
! Chown 更改 name 文件的 uid 和 gid；-1 表示不更改
! Lchown 更改 name 文件的数字 uid 和 gid; win 返回 syscall.EWINDOWS err
! Chtimes 更改 name 文件的访问和修改时间
*/

func TestCh(t *testing.T) {
	tmp, _ := os.MkdirTemp("", "tmp")
	os.Chdir(tmp)
	defer os.Remove(tmp)

	for range 5 {
		f, _ := os.CreateTemp(".", "*-logs")
		defer f.Close()
	}
	dirs, _ := os.ReadDir(tmp)
	for _, d := range dirs {
		log(d.Name())
	}
}

/*
! Environ 返回表示环境的字符串副本 "key=value"
! Setenv 设置由 key 命名的环境变量的值。
! Getenv 检索由 key 命名的环境变量的值。
! Unsetenv 取消设置 key 环境变量
! LookupEnv 检索由 key 命名的环境变量的值
! Clearenv 删除所有环境变量。
! Expand 根据映射函数 map 替换字符串中的 ${var} 或 $var。
! ExpandEnv 根据当前环境变量的值替换字符串中的 ${var} 或 $var
*/

func TestEnvs(t *testing.T) {
	t.Run("GetSet", func(t *testing.T) {
		os.Setenv("NAME", "JimryYchao")
		os.Setenv("HI", "Hello World")
		gopath := os.Getenv("GOPATH")
		logfln("%s\nGOPATH = %s", os.ExpandEnv("${NAME}: $HI"), gopath)
	})

	// os.Clearenv()
	t.Run("Environ", func(t *testing.T) {
		var envs map[string]string = make(map[string]string)
		for _, e := range os.Environ() {
			r := strings.Split(e, "=")
			if len(r) > 1 {
				envs[r[0]] = r[1]
			}
		}
		for _, p := range strings.Split(envs["PATH"], ";") {
			if p != "" {
				logfln("%s", p)
			}
		}
	})

	t.Run("Expand", func(t *testing.T) {
		log(os.Expand("${AUTHOR}, $TIME_NOW", func(s string) string {
			switch s {
			case "AUTHOR":
				return "JimryYchao"
			case "TIME_NOW":
				return time.Now().UTC().String()
			default:
				return ""
			}
		}))
	})
}

/*
! Pipe 返回一对连接的 *File; 从 w 写入，从 R 读取。返回文件的基础 Windows 句柄被标记为可由子进程继承。
! ReadFile 读取 name 文件并返回内容。
! WriteFile 将数据写入 name 文件，并在必要时创建它
! DirFS 返回一个文件系统（一个 fs.FS） 表示根目录目录中的文件树。结果实现 io/fs.StatFS、io/fs.ReadFileFS 和 io/fs.ReadDirFS
! ReadDir 读取 name 目录，返回按文件名排序的所有目录条目 DirEntry
	os.DirEntry 是从目录读取的条目（使用 ReadDir 函数或 File 的 ReadDir 方法）
! Stat 返回描述 name 文件的 FileInfo
! os.File 表示打开的文件描述符。
	Create 创建或截断 name 文件。如果文件已存在，则将其截断; 文件不存在，则使用模式 0666（在 umask 之前）创建该文件。
	CreateTemp 在目录目录中创建一个新的临时文件
	NewFile 返回一个具有给定文件描述符和名称的新文件
	Open 将打开 name 文件以供读取
	OpenFile “打开” 或 “创建” name 文件

	ChDir, Chmod, Chown, Truncate
	Name, ReadDir, Readdir, Readdirnames
	Stat, Lstat 返回 os.FileInfo 描述一个文件
	Close, Read, ReadAt, ReadFrom, Seek, Write, WriteAt, WriteString, WriteTo
	Fd 返回 windows 文件句柄（描述符）
	Sync 同步到稳定存储
	SetDeadline, SetReadDeadline, SetWriteDeadline 设置文件的读取和写入截止时间; 普通文件一般不支持，但 Pipe 支持。
		超过截止时间后，可以通过在将来设置截止时间来刷新连接
	SyscallConn 返回原始文件。这将实现 Syscall.Conn 接口
*/

func TestPipe(t *testing.T) {
	r, w, _ := os.Pipe()
	rinfo, _ := r.Stat()
	winfo, _ := w.Stat()
	log("Pipe:", rinfo.Name(), winfo.Name())

	err := w.SetDeadline(time.Now().Add(100 * time.Millisecond))
	log(err)
	ctx, _ := context.WithTimeout(context.Background(), 200*time.Millisecond)
	go func() {
		for {
			select {
			case <-ctx.Done():
				w.Close()
				return
			default:
				_, err := w.WriteString("Hello World\n")
				time.Sleep(10 * time.Millisecond)
				if err != nil {
					w.Close()
					return
				}
			}
		}
	}()

out:
	for {
		select {
		case <-ctx.Done():
			r.Close()
			break out
		default:
			_, err := io.Copy(os.Stdout, r)
			if err != nil {
				r.Close()
				break out
			}
		}
	}
}

/*
! FindProcess 通过其 pid 查找正在运行的进程
! StartProcess 使用由 name、argv 和 attr 指定的程序、参数和属性启动一个新进程
! os.ProcAttr 包含将应用于 StartProcess 启动的新进程的属性
! os.Process 存储有关 StartProcess 创建的进程的信息
	Kill 会导致进程立即退出。Kill 不会等到进程实际退出。这只会终止进程本身，而不会终止它可能已经启动的任何其他进程。
	Release 会释放与 Process p 关联的任何资源，使其将来无法使用。仅当 Wait 不是时，才需要调用 Release。
	Signal 向进程发送信号。未在 Windows 上实现发送中断。
		os.Signal: Interrupt, Kill
	Wait 等待进程退出，然后返回描述其状态和错误（如果有）的 ProcessState; 在大多数操作系统上，进程必须是当前进程的子进程，否则将返回错误。
! ProcessState 存储有关进程的信息，如 Wait 报告的那样
	ExitCode 退出进程的退出代码
	Exited 报告程序是否已退出
	Pid 返回已退出进程的进程 ID
	Success 报告程序是否成功退出
	Sys 返回有关进程与系统相关的退出信息
	SysUsage 返回有关退出进程的与系统相关的资源使用情况信息
	SystemTime 返回退出进程及其子进程的系统 CPU 时间
	UserTime 返回退出进程及其子进程的用户 CPU 时间
*/

func TestProcess(t *testing.T) {
	cmddir, cmdbase := filepath.Split("hello/hello.exe")
	attr := &os.ProcAttr{Dir: cmddir, Files: []*os.File{nil, os.Stdout, os.Stderr}}

	p, err := os.StartProcess(cmdbase, nil, attr)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(3 * time.Second)
	if err = p.Signal(os.Kill); err != nil {
		p.Kill()
		t.Fatal(err)
	}

	if stat, err := p.Wait(); err == nil {
		logfln("exit code : %d", stat.ExitCode())
		logfln("is exited : %t", stat.Exited())
		logfln("pid : %d", stat.Pid())
		logfln("Sys : %v", stat.Sys())
		usag := stat.SysUsage().(*syscall.Rusage)
		logfln("Sys Usage : create at %v, exit at %v", usag.CreationTime, usag.ExitTime)
		logfln("SysTime : %v", stat.SystemTime())
		logfln("UserTime: %v", stat.UserTime())
	}
}
