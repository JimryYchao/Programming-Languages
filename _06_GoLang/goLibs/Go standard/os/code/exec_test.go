package gostd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

/*
! LookPath 在 PATH 环境变量命名的目录中搜索可执行命名文件；LookPath 还使用 PATHEXT 环境变量来匹配合适的候选项;
! Command, CommandContext 返回使用给定参数执行 name 程序的 cmd 结构；它只设置返回结构中的 Path 和 Args
	如果 name 不包含路径分隔符，则 Command 将使用 LookPath 将 name 解析为完整路径（如果可能）
	否则，它直接使用 name 作为 Path；返回的 Cmd 的 Args 字段由命令名称后跟 arg 的元素构造
	CommandContext 使用 ctx 中断进程
! exec.Cmd 表示正在准备或运行的外部命令。调用 Cmd 的 Run、Output 或 CombinedOutput 方法后无法重用
	CombinedOutput 运行该命令并返回合并的标准输出和标准错误
	Environ 返回运行命令的环境的副本，该环境将按照当前配置进行运行。
	Output 运行命令并返回其标准输出
	Run 启动指定的命令并等待其完成
	Start 将启动指定的命令，但不会等待它完成。成功调用 Start 后，必须调用 Wait 方法才能释放关联的系统资源
	Wait 等待命令退出，并等待任何复制到 stdin 或从 stdout 或 stderr 复制完成
		如果 c.Stdin、c.Stdout 或 c.Stderr 中的任何一个不是 *os.File，Wait 还会等待相应的 I/O 循环复制到进程或从进程复制完成
	StdinPipe, StderrPipe, StdoutPipe 返回管道，该管道将在命令启动时连接到命令的标准输入/错误/输出
	String 返回 cmd 的可读描述
*/

func TestCmdSh(t *testing.T) {
	cmd := exec.Command("sh", "-c", "echo stdout; echo 1>&2 stderr")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}

func TestCmdEnviron(t *testing.T) {
	cmd := exec.Command("pwd")

	// Set Dir before calling cmd.Environ so that it will include an
	// updated PWD variable (on platforms where that is used).
	cmd.Dir = ".."
	cmd.Env = append(cmd.Environ(), "POSIXLY_CORRECT=1")

	out, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s\n", out)
}

func TestCmdOutput(t *testing.T) {
	cmd := exec.Command("./hello/hello.exe")
	var out []byte
	var err error
	done := make(chan bool)
	go func() {
		out, err = cmd.Output()
		if err != nil {
			log(err) // process killed
		}
		done <- true
	}()
	time.Sleep(1 * time.Second)
	cmd.Process.Signal(os.Kill)
	<-done
	fmt.Printf("The date is:\n%s", out)
}

func TestCmdStartWait(t *testing.T) {
	cmd := exec.Command("sleep", "5")
	err := cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	log("Waiting for command to finish...")
	err = cmd.Wait()
	logfln("Command finished with error: %v", err)
}

func TestCmdRun(t *testing.T) {
	cmd := exec.Command("./hello/hello.exe")
	done := make(chan error)
	go func() {
		err := cmd.Run()
		done <- err
	}()
	time.Sleep(3 * time.Second)
	cmd.Process.Signal(os.Kill)
	log(<-done)
}

func TestCmdPipe(t *testing.T) {
	cmd := exec.Command("echo", "-n", `{"Name": "Bob", "Age": 32}`)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	var person struct {
		Name string
		Age  int
	}
	if err := json.NewDecoder(stdout).Decode(&person); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%s is %d years old\n", person.Name, person.Age)
}

func TestGitPush(t *testing.T) {
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}

	if err = exec.Command(git, "add", ".").Run(); err != nil {
		t.Fatal(err)
	}

	if err = exec.Command(git, `commit`, `-m`, `go exec git commit`).Run(); err != nil {
		t.Fatal(err)
	}

	if err = exec.Command(git, `push`).Run(); err != nil {
		t.Fatal(err)
	}
}
