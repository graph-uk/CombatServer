//This file is download or/and update dependencies, and build binaries.
//Recommended to set env "GOPATH" to THIS_FILE_PATH/..
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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
//D:\combat_server_current\src\github.com\graph-uk\combat-server\integration-tests\client\combat-client.exe
//combat-client.exe
func (t *cmdProcess) GetShortCommand() string {
	arr := strings.Split(t.Command, sl) // split by '/' or '\'
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

func startCmd(dir string, env *[]string, command string, args ...string) *cmdProcess {
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

// add value to environment using separator (like PATH, GOPATH, GOROOT)
// if value not exist - create new
func envAdd(env []string, name, value string) []string {
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
func envRewrite(env []string, name, value string) []string {
	for curVarIndex, curVarValue := range env {
		if strings.HasPrefix(curVarValue, name+`=`) {
			env[curVarIndex] = name + `=` + value
			return env
		}
	}
	env = append(env, name+`=`+value)
	return env
}

// Copy the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func CopyFile(src, dst string) error {
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	sl = string(os.PathSeparator)

	var err error
	curdir, err = os.Getwd()
	check(err)
}

var sl, curdir string // system filepath separator (/ or \), dir which script started

func CopyDir(src, dst string) error {
	cmd := exec.Command(`xcopy`, `/s`, `/e`, `/c`, `/h`, `/k`, `/y`, src, dst+`\`)
	//log.Printf("Running cp -a")
	return cmd.Run()
}
func createFailTrigger() {
	file := os.TempDir() + `\testSuccessOrFailure.txt`
	f, _ := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0666)
	f.Close()
}

func deleteFailTrigger() {
	file := os.TempDir() + `\testSuccessOrFailure.txt`
	//	var err error
	os.RemoveAll(file)
	//fmt.Println(err.Error())
	os.Remove(file)
	//fmt.Println(err.Error())
}

func main() {
	//Re-create (clear) folders for test binaries
	os.RemoveAll(`server`)
	os.RemoveAll(`client`)
	os.RemoveAll(`worker`)
	check(os.MkdirAll(`server`, 0777))
	check(os.MkdirAll(`client`, 0777))
	check(os.MkdirAll(`worker`, 0777))
	createFailTrigger()

	//Copy compiled binaries to correspond test folders
	check(CopyFile(`..`+sl+`..`+sl+`combat-server`+sl+`combat-server.exe`, `server`+sl+`combat-server.exe`))
	check(CopyFile(`config.json`, `server`+sl+`config.json`))
	check(CopyFile(`..`+sl+`..`+sl+`combat-client`+sl+`combat-client.exe`, `client`+sl+`combat-client.exe`))
	check(CopyFile(`..`+sl+`..`+sl+`combat-worker`+sl+`combat-worker.exe`, `worker`+sl+`combat-worker.exe`))
	//check(CopyDir(`..`+sl+`..`+sl+`combat-server`+sl+`server`, `server`+sl+`server`))
	//check(CopyDir(`..`+sl+`..`+sl+`combat-server`+sl+`assets`, `server`+sl+`assets`))

	//Configure environment variable for the server and workers
	env := envRewrite(os.Environ(), `GOPATH`, curdir+sl+`CombatTestsExample`+sl)
	env = envRewrite(env, `GOROOT`, curdir+sl+`..`+sl+`..`+sl+`..`+sl+`..`+sl+`..`+sl+`combat-dev-utils`+sl+`combat-dev-go`)
	env = envAdd(env, `PATH`, curdir+sl+`..`+sl+`..`+sl+`combat`)
	env = envAdd(env, `PATH`, curdir+sl+`..`+sl+`..`+sl+`..`+sl+`..`+sl+`..`+sl+`combat-dev-utils`+sl+`combat-dev-go`+sl+`bin`)

	//fmt.Println(env)
	//return

	//run server, client worker. Kill before quit.
	server := startCmd(curdir+sl+`server`, &env, `.`+sl+`combat-server.exe`)
	client := startCmd(curdir+sl+`CombatTestsExample`+sl+`src`+sl+`Tests`, nil, curdir+sl+`client`+sl+`combat-client.exe`, `http://localhost:3133`, `./../..`, `40`, `-InternalIP=192.168.1.1`)

	worker := startCmd(curdir+sl+`worker`, &env, `.`+sl+`combat-worker.exe`, `http://localhost:3133`)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(`----------------------------------------Server stdout-----------------------------------------`)
			fmt.Println(string(server.StdOutBuf))
			fmt.Println(`----------------------------------------Server stderr-----------------------------------------`)
			fmt.Println(string(server.StdErrBuf))
			fmt.Println(`----------------------------------------Client stdout-----------------------------------------`)
			fmt.Println(string(client.StdOutBuf))
			fmt.Println(`----------------------------------------Client stderr-----------------------------------------`)
			fmt.Println(string(client.StdErrBuf))
			fmt.Println(`----------------------------------------Worker stdout-----------------------------------------`)
			fmt.Println(string(worker.StdOutBuf))
			fmt.Println(`----------------------------------------Worker stderr-----------------------------------------`)
			fmt.Println(string(worker.StdErrBuf))
		}
	}()

	//defer server.Cmd.Process.Kill()
	defer client.Cmd.Process.Kill()
	//defer worker.Cmd.Process.Kill()

	//time.Sleep(10 * time.Second)

	//Check server's output
	//server.WaitingForStdOutContains(`config.json is not found. Default config will be created`, 10*time.Second)
	server.WaitingForStdOutContains(`http server started on`, 10*time.Second)
	server.WaitingForStdOutContains(`Created:  _data/sessions`, 10*time.Second)
	server.WaitingForStdOutContains(`TestFail -InternalIP=192.168.1.1`, 40*time.Second)
	server.WaitingForStdOutContains(`TestSuccess -InternalIP=192.168.1.1`, 40*time.Second)
	server.WaitingForStdOutContains(`Try status: Failed`, 200*time.Second)

	//Check worker's output
	worker.WaitingForStdOutContains(`CaseRunning TestSuccess -InternalIP=192.168.1.1`, 2*time.Minute)
	worker.WaitingForStdOutContains(`Run case... Ok`, time.Minute)
	worker.WaitingForStdOutContains(`CaseRunning TestFail -InternalIP=192.168.1.1`, 2*time.Minute)
	worker.WaitingForStdOutContains(`Run case... Fail`, time.Minute)
	worker.WaitingForStdOutContains(`Failed here for example`, time.Minute)

	//Check client's output
	client.WaitingForStdOutContains(`Cleanup tests`, 10*time.Second)
	client.WaitingForStdOutContains(`Packing tests`, 10*time.Second)
	client.WaitingForStdOutContains(`Uploading session`, 10*time.Second)
	client.WaitingForStdOutContains(` - Pending`, 10*time.Second)
	client.WaitingForStdOutContains(`Case exploring`, time.Minute)
	client.WaitingForStdOutContains(` - Processing`, 40*time.Second)
	client.WaitingForStdOutContains(`Processed 0 of 4 tests`, 40*time.Second)
	client.WaitingForStdOutContains(`Time of testing`, 400*time.Second)
	//Time of testing
	//panic(`test`)

	deleteFailTrigger()
	client = startCmd(curdir+sl+`CombatTestsExample`+sl+`src`+sl+`Tests`, nil, curdir+sl+`client`+sl+`combat-client.exe`, `http://localhost:3133`, `./../..`, `40`, `-InternalIP=192.168.1.1`)
	client.WaitingForStdOutContains(`Time of testing`, 400*time.Second)

	//client.WaitingForExitWithCode(40*time.Second, 0)
	//time.Sleep(20*time.Second)

	//server.WaitingForStdOutContains(`Slack alert sent. Response:  ok`, 40*time.Second)

	//panic(`test`)
	log.Println(`The test finished succeed.`)
}
