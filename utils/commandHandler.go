package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// CommandHandler ...
type CommandHandler struct {
}

// ExecuteCommand ...
func (t *CommandHandler) ExecuteCommand(command string, arguments []string, path string) (bytes.Buffer, error) {
	sl := string(os.PathSeparator)

	var outputBuffer bytes.Buffer
	var errorBuffer bytes.Buffer

	cmd := exec.Command(command, arguments...)
	cmd.Env = addToGOPath(path + sl + `..` + sl + `..`)

	cmd.Dir = path
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &errorBuffer

	exitStatus := cmd.Run()

	if exitStatus != nil {
		fmt.Println("Cannot extract cases")
		fmt.Println(exitStatus)
		fmt.Println(outputBuffer.String())
		fmt.Println(errorBuffer.String())
	}

	return outputBuffer, exitStatus
}

func addToGOPath(path string) []string {
	result := os.Environ()
	abslutePath, _ := filepath.Abs(path)
	for curVarIndex, curVarValue := range result {
		if strings.HasPrefix(curVarValue, "GOPATH=") {
			result[curVarIndex] = result[curVarIndex] + string(os.PathListSeparator) + abslutePath
			return result
		}
	}
	return result
}
