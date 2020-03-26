package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"malibu/arrayUtils"
)

type Test struct {
	directory string
	name      string
	params    map[string]TestParameter
	tags      []string
}

type TestParameter struct {
	Name     string
	Type     string
	Variants []string
}

func (t *Test) addToGOPath(pathExtention string) []string {
	result := os.Environ()
	for curVarIndex, curVarValue := range result {
		if strings.HasPrefix(curVarValue, "GOPATH=") {
			result[curVarIndex] = result[curVarIndex] + string(os.PathListSeparator) + pathExtention
			return result
		}
	}
	return result
}

func (t *Test) LoadTagsAndParams() error {
	type UnmarshaledTestParams struct {
		Params []TestParameter
		Tags   []string
	}

	// get test's params in JSON
	rootTestsPath, _ := os.Getwd()
	rootTestsPath += string(os.PathSeparator) + ".." + string(os.PathSeparator) + ".."
	cmd := exec.Command("go", "run", t.directory+"/"+t.name+`/`+"main.go", "paramsJSON")
	cmd.Env = t.addToGOPath(rootTestsPath)
	var out, outErr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outErr
	cmd.Run()
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//println("LoadTagsAndParams: " + out.String())
	var TestParams UnmarshaledTestParams
	if err := json.Unmarshal(out.Bytes(), &TestParams); err != nil {
		log.Println("Cannot parse json for test: " + t.name)
		log.Println("JSON data: " + out.String())
		fmt.Println(outErr.String())
		panic(err)
	}
	t.tags = TestParams.Tags

	for _, curParameter := range TestParams.Params {
		t.params[curParameter.Name] = curParameter
	}

	return nil
}

func (t *Test) isCasesEqual(case1 []string, case2 []string) bool {
	if len(case1) != len(case2) {
		return false
	}

	result := true
	for _, curParameter := range case1 {
		parameterFound := false
		for _, curParameter2 := range case2 {
			if curParameter == curParameter2 {
				parameterFound = true
				break
			}
		}
		if !parameterFound {
			return false
		}
	}
	return result
}

func (t *Test) isCasePresented(allCases [][]string, aCase []string) bool {
	for _, curCase := range allCases {
		if t.isCasesEqual(curCase, aCase) {
			return true
		}
	}
	return false
}

func (t *Test) GetCasesByParameterCombinations(paramCombinations []*map[string]string) [][]string {
	var result [][]string

	for _, curCombination := range paramCombinations {
		curCombinationAccepted := true
		curCombinationCase := []string{t.name}
		for nameOfcurParamOfTest, curParamOfTest := range t.params {
			if curParamOfTest.Type == "EnumParam" {
				if !arrayUtils.StringInSlice((*curCombination)[nameOfcurParamOfTest], curParamOfTest.Variants) {
					curCombinationAccepted = false
					break
				}
			}
		}
		if curCombinationAccepted {
			for nameOfcurParamOfTest, _ := range t.params {
				curCombinationCase = append(curCombinationCase, "-"+nameOfcurParamOfTest+"="+(*curCombination)[nameOfcurParamOfTest])
			}
			if !t.isCasePresented(result, curCombinationCase) {
				result = append(result, curCombinationCase)
			}
		}
	}
	return result
}
