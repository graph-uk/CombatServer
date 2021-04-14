package CMD

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
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

//Run command. If existStatus==0 - silent
//Otherwise - print exit status, stdOut, stdErr, and panic
func CMDSync(dir string, env *[]string, command string, args ...string) string {
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
		fmt.Println(exitCode.Error())
		fmt.Println(out.String())
		fmt.Println(outErr.String())
		panic(`error on cmdSilent ` + command)
	}
	return out.String() + "\r\n" + outErr.String()
}
