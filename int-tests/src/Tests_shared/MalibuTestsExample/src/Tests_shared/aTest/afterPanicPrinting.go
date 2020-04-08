package aTest

import (
	"fmt"
	"io/ioutil"
	//	"os"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
)

func getMainFileNameAndLine(byteStack []byte) (string, int, error) {
	stack := string(byteStack)

	stackLines := strings.Split(stack, "\n")
	lastMainEntry := 0
	for i := len(stackLines) - 1; i > 0; i-- {
		match, _ := regexp.MatchString(`main.go:(\d*)\s`, stackLines[i])
		if match {
			lastMainEntry = i
			break
		}
	}

	mainStackLine := stackLines[lastMainEntry]
	mainStackLine = strings.Split(mainStackLine, " ")[0] //cut offset
	//	fmt.Println(len(stackLines))
	//	fmt.Println(strconv.Itoa(lastMainEntry))
	//	fmt.Println(stackLines[lastMainEntry])

	fnBegin := strings.Index(mainStackLine, "main.go:")

	mainFileName := mainStackLine[:fnBegin+7]
	//fmt.Println(mainFileName)
	stringLineNumber := mainStackLine[fnBegin+8:]
	//fmt.Println(stringLineNumber)
	lineNumber, err := strconv.Atoi(stringLineNumber)
	if err != nil {
		return "", 0, err
	}
	return mainFileName, lineNumber, nil
}

func getSrcLine(srcFileName string, line int) string {
	srcFileName = strings.TrimSpace(srcFileName)
	content, err := ioutil.ReadFile(srcFileName)
	if err != nil {
		fmt.Println(srcFileName)
		fmt.Println(err.Error())
		return "Cannot parse source line " + strconv.Itoa(line)
	}
	lines := strings.Split(string(content), "\n")
	return strings.TrimSpace(lines[line-1])
}

func getSrcBlock(srcFileName string, line int) string {
	srcFileName = strings.TrimSpace(srcFileName)
	content, err := ioutil.ReadFile(srcFileName)
	if err != nil {
		fmt.Println(srcFileName)
		fmt.Println(err.Error())
		return "Cannot parse source line " + strconv.Itoa(line)
	}
	lines := strings.Split(string(content), "\n")

	//fmt.Println(strconv.Itoa(line))

	// localize begin of block by empty line
	blockBegin := 0
	for i := line - 1; i > 0; i-- {
		curLine := lines[i]
		curLine = strings.TrimSpace(curLine)
		if curLine == "" {
			blockBegin = i + 1
			break
		}
		//		blockBegin = i
		//		break
	}

	// localize end of block by empty line
	blockEnd := len(lines) - 1
	for i := line - 1; i < len(lines); i++ {
		curLine := lines[i]
		curLine = strings.TrimSpace(curLine)
		if curLine == "" {
			blockEnd = i - 1
			break
		}
	}

	result := ""
	for i := blockBegin; i <= blockEnd; i++ {
		stringNumber := fmt.Sprintf("%-5d", i+1)

		arrow := "    " //arrow is enabled for error line only
		if i+1 == line {
			arrow = "===>"
		}
		result = result + stringNumber + arrow + strings.TrimSpace(lines[i]) + "\r\n"
	}
	result = `#############Code lines in main.go################` + "\r\n" + result

	return result
}

func GetClearedStack(originStack []byte) string {
	originStackLines := strings.Split(string(originStack), "\n")
	clearStack := ``
	for _, curLine := range originStackLines {
		if len(curLine) < 1 {
			continue
		}
		curLine := strings.TrimSpace(curLine)
		lineBegin := strings.Index(curLine, `src/`)
		lineEnd := strings.Index(curLine, `+`)
		if (lineBegin > -1) && (lineEnd > -1) {
			curLine = curLine[lineBegin:lineEnd]
			clearStack += curLine + "\r\n"
		}
	}
	return clearStack
}

func PrintSourceAndContinuePanic(r interface{}) {
	stack := debug.Stack()
	mainFileName, curLine, err := getMainFileNameAndLine(stack)
	if err != nil {
		fmt.Println("cannot extract main source line index")
		panic(r)
	}
	fmt.Println(getSrcBlock(mainFileName, curLine))

	fmt.Println(`#############Short stacktrace################`)
	fmt.Println(GetClearedStack(stack))
	fmt.Println(`#############Full stacktrace################`)
	//fmt.Println(string(stack))
	//os.Exit(1)
	panic(r)
}
