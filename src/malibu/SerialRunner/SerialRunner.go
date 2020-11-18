package SerialRunner

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	TestName   string
	caseParams []string
	command    []string
	triesCount int
	isSucceed  bool
}

func addToGOPath(pathExtention string) []string {
	result := os.Environ()
	for curVarIndex, curVarValue := range result {
		if strings.HasPrefix(curVarValue, "GOPATH=") {
			result[curVarIndex] = result[curVarIndex] + string(os.PathListSeparator) + pathExtention
			return result
		}
	}
	return result
}

func addLeftTab(str string) string {
	result := ""
	strArray := strings.Split(str, "\n")
	for _, curStr := range strArray {
		result += "    " + strings.TrimSpace(curStr) + "\r\n"
	}
	result = "    " + strings.TrimSpace(result)
	return result
}

func isAllCasesRunned(testCases []testCase, maxTriesCount int) bool {
	result := true
	for _, curTestCase := range testCases {
		if curTestCase.isSucceed {
			continue
		} else {
			if maxTriesCount < curTestCase.triesCount {
				result = false
			}
		}
	}
	fmt.Println("allCasesRunned:", result)
	return result
}

func RunCasesSerial(cases [][]string, directory string) int {
	sl := string(os.PathSeparator)

	maxTriesCount := 3 // hardcoded 3 tries. May be custom values will be provided by CLI
	fmt.Println("Run cases.")
	var testCases []testCase

	// fill tries map
	for _, curCase := range cases {
		var curTestCase testCase
		curTestCase.TestName = curCase[0]
		curTestCase.caseParams = curCase[1:]

		curTestCase.isSucceed = false
		curTestCase.triesCount = 0
		curTestCase.command = []string{"run"}
		curTestCase.command = append(curTestCase.command, directory+sl+curTestCase.TestName+sl+"main.go")
		curTestCase.command = append(curTestCase.command, curTestCase.caseParams...)
		testCases = append(testCases, curTestCase)
	}

	for true {
		hasAnyTries := false
		for curCaseIndex, curCase := range testCases {
			if !curCase.isSucceed && curCase.triesCount < maxTriesCount {
				hasAnyTries = true
				testCases[curCaseIndex].triesCount++
				os.Chdir(directory + string(os.PathSeparator) + curCase.TestName)

				rootTestsPath, _ := os.Getwd()
				rootTestsPath += string(os.PathSeparator) + ".." + string(os.PathSeparator) + ".." + string(os.PathSeparator) + ".."

				cmd := exec.Command("go", curCase.command...)
				cmd.Env = addToGOPath(rootTestsPath)
				var out bytes.Buffer
				var outErr bytes.Buffer
				cmd.Stdout = &out
				cmd.Stderr = &outErr
				fmt.Println(curCase.TestName, " ", curCase.caseParams, " try:", testCases[curCaseIndex].triesCount)
				exitStatus := cmd.Run()

				if exitStatus != nil {
					fmt.Println(addLeftTab(exitStatus.Error()))
					fmt.Println(addLeftTab(out.String()))
					fmt.Println(addLeftTab(outErr.String()))
				} else {
					testCases[curCaseIndex].isSucceed = true
					fmt.Println(addLeftTab("          OK"))
				}
				fmt.Println()
			}
		}

		if !hasAnyTries {
			break
		}
	}

	FailedCasesCount := 0
	for _, curCase := range testCases {
		if !curCase.isSucceed {
			FailedCasesCount++
		}
	}

	fmt.Println("Total failed cases: ", FailedCasesCount)
	return FailedCasesCount
}
