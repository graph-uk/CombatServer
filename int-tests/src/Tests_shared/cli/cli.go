package cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//run command, if existStatus==0 - silent
//otherwise - print exit status, stdOut, stdErr, and exit process with code 1
func cmdSilent(dir string, env *[]string, command string, args ...string) {
	cmd := exec.Command(command, args...)
	if env != nil {
		cmd.Env = *env
	}
	if dir != `` {
		cmd.Dir = dir
	}
	var out, outErr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &outErr
	exitCode := cmd.Run()
	if exitCode != nil {
		log.Println(exitCode.Error())
		fmt.Println(out.String())
		fmt.Println(outErr.String())
		os.Exit(1)
	}
}

type cmdProcess struct {
	Cmd       *exec.Cmd
	Command   string
	StdOut    io.ReadCloser
	StdErr    io.ReadCloser
	StdErrBuf []byte
	StdOutBuf []byte
}

// this loop function - for separate concurrency go-routine.
// it is get text from console pipe.
// if command's buffer will overflow - command was paused untill we get this bytes
func (t *cmdProcess) refreshErrBufLoop() {
	buf := make([]byte, 512)
	for {
		len, err := t.StdErr.Read(buf)
		if err != nil {
			if err.Error() == `EOF` { // if the pipe closed (app is finished) - stop watching
				break
			} else {
				panic(err)
			}
		}
		if len > 0 {
			t.StdErrBuf = append(t.StdErrBuf, buf[:len]...)
		}
		if len == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// this function returns cut filepath on t.Command, and return short command
//D:\malibu_server_current\src\github.com\graph-uk\malibu-server\integration-tests\client\malibu-client.exe
//malibu-client.exe
func (t *cmdProcess) GetShortCommand() string {
	arr := strings.Split(t.Command, Sl) // split by '/' or '\'
	if len(arr) > 0 {
		return arr[len(arr)-1]
	} else {
		return `Cannot extract short command`
	}
}

// this loop function - for separate concurrency go-routine.
// it is get text from console pipe.
// if command's buffer will overflow - command was paused untill we get this bytes
func (t *cmdProcess) refreshOutBufLoop() {
	buf := make([]byte, 512)
	for {
		len, err := t.StdOut.Read(buf)
		if err != nil {
			if err.Error() == `EOF` { // if the pipe closed (app is finished) - stop watching
				break
			} else {
				panic(err)
			}
		}
		if len > 0 {
			t.StdOutBuf = append(t.StdOutBuf, buf[:len]...)
		}
		if len == 0 {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (t *cmdProcess) WaitingForStdErrContains(textPart string, timeout time.Duration) {
	startMoment := time.Now()
	log.Println(`AwaitErr - ` + t.GetShortCommand() + `: ` + textPart)
	for {
		if strings.Contains(string(t.StdErrBuf), textPart) {
			break
		}
		if startMoment.Add(timeout).Before(time.Now()) { // if timed out
			panic(`TimeoutErr - ` + t.GetShortCommand() + `: ` + textPart)
		}
		time.Sleep(time.Second)
	}
}

func (t *cmdProcess) WaitingForStdOutContains(textPart string, timeout time.Duration) {
	startMoment := time.Now()
	log.Println(`AwaitOut - ` + t.GetShortCommand() + `: ` + textPart)
	for {
		if strings.Contains(string(t.StdOutBuf), textPart) {
			break
		}
		if startMoment.Add(timeout).Before(time.Now()) { // if timed out
			panic(`TimeoutOut - ` + t.GetShortCommand() + `: ` + textPart)
		}
		time.Sleep(time.Second)
	}
}

func (t *cmdProcess) WaitingForExitWithCode(timeout time.Duration, expectedExitCode int) {
	log.Println(`AwaitExitWithExitCode ` + strconv.Itoa(expectedExitCode) + ` ` + t.GetShortCommand())

	ch := make(chan bool, 1)
	defer close(ch)

	go func() {
		t.Cmd.Wait()
		ch <- true
	}()

	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	select {
	case <-ch:
	case <-timer.C:
		panic(`TimeoutOut - Wait for exit with code ` + strconv.Itoa(expectedExitCode) + ` ` + t.GetShortCommand())
	}

	ws := t.Cmd.ProcessState.Sys().(syscall.WaitStatus)
	exitCode := ws.ExitStatus()
	if exitCode != expectedExitCode {
		panic(strconv.Itoa(expectedExitCode) + ` exit code expected, but the process is finished, with '` + strconv.Itoa(exitCode) + `' exit code. ` + t.GetShortCommand())
	}
}

func StartCmd(dir string, env *[]string, command string, args ...string) *cmdProcess {
	var res cmdProcess
	res.Command = command

	res.Cmd = exec.Command(command, args...)
	if env != nil {
		res.Cmd.Env = *env
	}
	if dir != `` {
		res.Cmd.Dir = dir
	}

	var err error
	res.StdErr, err = res.Cmd.StderrPipe()
	check(err)
	res.StdOut, err = res.Cmd.StdoutPipe()
	check(err)
	err = res.Cmd.Start()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	go res.refreshOutBufLoop() // stdout/stderr pipe-readers routines
	go res.refreshErrBufLoop()

	log.Println(`Started: ` + command)
	return &res
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func CopyFile(src, dst string, exec bool) {
	check(copyFile(src, dst))
	if exec && runtime.GOOS == "linux" {
		check(os.Chmod(dst, 0777))
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	Sl = string(os.PathSeparator)

	var err error
	curdir, err = os.Getwd()
	check(err)
}

var Sl, curdir string // system filepath separator (/ or \), dir which script started

func CopyDir(src, dst string) error {
	cmd := exec.Command(`xcopy`, `/s`, `/e`, `/c`, `/h`, `/k`, `/y`, src, dst+`\`)
	//log.Printf("Running cp -a")
	return cmd.Run()
}

func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func MkDir(path string, perm os.FileMode) {
	check(os.MkdirAll(path, perm))
}

func CreateFailTrigger(name string) {
	file := os.TempDir() + `\` + name
	f, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0666)
	check(err)
	check(f.Close())
}

func DeleteFailTrigger(name string) {
	file := os.TempDir() + `\` + name
	os.RemoveAll(file)
	//os.Remove(file)
}

func Pwd() string {
	dir, err := os.Getwd()
	check(err)
	return dir
}

// add item to environment variable with os-specified separator (like PATH, GOPATH, GOROOT)
// if value not exist - create new
func EnvExtend(env []string, name, value string) []string {
	for curVarIndex, curVarValue := range env {
		if strings.HasPrefix(curVarValue, name+`=`) {
			env[curVarIndex] = env[curVarIndex] + string(os.PathListSeparator) + value
			return env
		}
	}
	env = append(env, name+`=`+value)
	return env
}

// clear and rewrite exist value by new
// if value not exist - create new
func EnvRewrite(env []string, name, value string) []string {
	for curVarIndex, curVarValue := range env {
		if strings.HasPrefix(curVarValue, name+`=`) {
			env[curVarIndex] = name + `=` + value
			return env
		}
	}
	env = append(env, name+`=`+value)
	return env
}
