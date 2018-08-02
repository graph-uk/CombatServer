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
	"strings"
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

func (t *cmdProcess) refreshErrBuf() {
	buf := make([]byte, 5120)
	n, err := io.ReadAtLeast(t.StdErr, buf, 1)
	check(err)
	t.StdErrBuf = append(t.StdErrBuf, buf[:n]...)
}

func (t *cmdProcess) refreshOutBuf() {
	buf := make([]byte, 512)
	n, err := io.ReadAtLeast(t.StdOut, buf, 1)
	check(err)
	t.StdOutBuf = append(t.StdOutBuf, buf[:n]...)
}

func (t *cmdProcess) WaitingForErrBufContains(textPart string, timeout time.Duration) {
	startMoment := time.Now()
	log.Println(`AwaitErr - ` + t.Command + `: ` + textPart)
	for {
		t.refreshErrBuf()
		if strings.Contains(string(t.StdErrBuf), textPart) {
			//log.Println(`FoundErr - ` + t.Command + `: ` + textPart)
			break
		}
		if startMoment.Add(timeout).Before(time.Now()) { // if timed out
			panic(`TimeoutErr - ` + t.Command + `: ` + textPart)
		}
		time.Sleep(time.Second)
	}
}

func (t *cmdProcess) WaitingForOutBufContains(textPart string, timeout time.Duration) {
	startMoment := time.Now()
	log.Println(`AwaitOut - ` + t.Command + `: ` + textPart)
	for {
		t.refreshOutBuf()
		if strings.Contains(string(t.StdOutBuf), textPart) {
			//log.Println(`FoundOut - ` + t.Command + `: ` + textPart)
			break
		}
		if startMoment.Add(timeout).Before(time.Now()) { // if timed out
			panic(`TimeoutOut - ` + t.Command + `: ` + textPart)
		}
		time.Sleep(time.Second)
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

func main() {
	//Re-create (clear) folders for test binaries
	os.RemoveAll(`server`)
	os.RemoveAll(`client`)
	os.RemoveAll(`worker`)
	check(os.MkdirAll(`server`, 0777))
	check(os.MkdirAll(`client`, 0777))
	check(os.MkdirAll(`worker`, 0777))

	//Copy compiled binaries to correspond test folders
	check(CopyFile(`..`+sl+`..`+sl+`combat-server`+sl+`combat-server.exe`, `server`+sl+`combat-server.exe`))
	check(CopyFile(`..`+sl+`..`+sl+`combat-client`+sl+`combat-client.exe`, `client`+sl+`combat-client.exe`))
	check(CopyFile(`..`+sl+`..`+sl+`combat-worker`+sl+`combat-worker.exe`, `worker`+sl+`combat-worker.exe`))

	//Configure environment variable for the server and workers
	env := envRewrite(os.Environ(), `GOPATH`, curdir+sl+`CombatTestsExample`+sl)
	env = envRewrite(env, `GOROOT`, curdir+sl+`..`+sl+`..`+sl+`..`+sl+`..`+sl+`..`+sl+`combat-dev-utils`+sl+`combat-dev-go`)
	env = envAdd(env, `PATH`, curdir+sl+`..`+sl+`..`+sl+`combat`)
	env = envAdd(env, `PATH`, curdir+sl+`..`+sl+`..`+sl+`..`+sl+`..`+sl+`..`+sl+`combat-dev-utils`+sl+`combat-dev-go`+sl+`bin`)

	//	fmt.Println(env)
	//	return

	//run server, client worker. Kill before quit.
	server := startCmd(curdir+sl+`server`, &env, `.`+sl+`combat-server.exe`)
	defer server.Cmd.Process.Kill()
	client := startCmd(curdir+sl+`CombatTestsExample`+sl+`src`+sl+`Tests`, nil, curdir+sl+`client`+sl+`combat-client.exe`, `http://localhost:9090`, `40`, `-InternalIP=192.168.1.1`)
	defer client.Cmd.Process.Kill()
	worker := startCmd(curdir+sl+`worker`, &env, `.`+sl+`combat-worker.exe`, `http://localhost:9090`)
	defer worker.Cmd.Process.Kill()

	time.Sleep(10 * time.Second)

	//Check server's output
	//server.WaitingForOutBufContains(`config.json is not found. Default config will be created`, 10*time.Second)
	//	server.WaitingForOutBufContains(`Serving combat tests at port: 9090...`, 10*time.Second)
	//	server.WaitingForOutBufContains(`Create new session: `, 10*time.Second)
	//	server.WaitingForOutBufContains(` 40 -InternalIP=192.168.1.1`, 10*time.Second)
	//	server.WaitingForOutBufContains(`Explored 2 cases for session: `, 10*time.Second)
	//	server.WaitingForOutBufContains(`Get a job (CasesRun) for case: `, 10*time.Second)
	//	server.WaitingForOutBufContains(`Provide result for case: `, 10*time.Second)

	//	//Check worker's output
	//	worker.WaitingForOutBufContains(`getJob - idle`, time.Minute)
	//	worker.WaitingForOutBufContains(`getJob - RunCase`, time.Minute)
	//	worker.WaitingForOutBufContains(`CaseRunning TestFail -InternalIP=192.168.1.1`, time.Minute)
	//	worker.WaitingForOutBufContains(`Run case... Fail`, time.Minute)
	//	worker.WaitingForOutBufContains(`CaseRunning TestSuccess -InternalIP=192.168.1.1`, time.Minute)
	//	//worker.WaitingForOutBufContains(`Failed here for example`, time.Minute)

	//	//Check server's output
	//	client.WaitingForOutBufContains(`Cleanup tests`, 10*time.Second)
	//	client.WaitingForOutBufContains(`Packing tests`, 10*time.Second)
	//	client.WaitingForOutBufContains(`Uploading session`, 10*time.Second)
	//	client.WaitingForOutBufContains(`Session status: http://localhost:9090/sessions/`, 10*time.Second)
	//	client.WaitingForOutBufContains(`Cases exploring`, 10*time.Second)
	//	client.WaitingForOutBufContains(`Testing (0/2)`, 10*time.Second)
	//client.WaitingForOutBufContains(`Testing (1/2)`, 10*time.Second)
	// client.WaitingForOutBufContains(`Finished with `, 60*time.Second)
	// client.WaitingForOutBufContains(`More info at: `, 60*time.Second)
	// client.WaitingForOutBufContains(`Time of testing: `, 60*time.Second)

	//	time.Sleep(10 * time.Second)
	//	client.refreshOutBuf()
	//	fmt.Println(string(client.StdOutBuf))

	time.Sleep(3 * time.Minute)
}
