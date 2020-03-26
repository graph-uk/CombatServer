package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"os"

	"malibu/CLIParser"

	"malibu/arrayUtils"
)

// This is the base struct contain all required in all test fields
type TestManager struct {
	tests                 map[string]*Test
	parametersFromCLI     map[string]string
	testParametersFromCLI map[string]string
	//testMergedParameters TestParameter
}

// Parse all parameters from CLI. Fill default values if needed.
func (t *TestManager) parseAllCLIParameters() {
	t.parametersFromCLI = CLIParser.ParseAllCLIFlags()
	if _, ok := t.parametersFromCLI["name"]; !ok {
		t.parametersFromCLI["name"] = ""
	}

	if _, ok := t.parametersFromCLI["tag"]; !ok {
		t.parametersFromCLI["tag"] = ""
	}
}

// Parse test parameters from CLI (except for action, name, tag).
func (t *TestManager) parseTestCLIParameters() {
	t.testParametersFromCLI = make(map[string]string)
	for curParamName, curParamValue := range t.parametersFromCLI {
		if curParamName != "name" && curParamName != "tag" {
			t.testParametersFromCLI[curParamName] = curParamValue
		}
	}
}

// parse parameters from CLI, load parameters of each test, and filter tests by CLI parameters (name,tag)
func (t *TestManager) Init(directory string) error {
	t.parseAllCLIParameters()
	t.parseTestCLIParameters()
	t.selectAllTests(directory)
	t.filterTestsByName()
	t.filterTestsByTag()
	return nil
}

//Select all tests in the directory, load that's parameters, and collect it to t.tests
func (t *TestManager) selectAllTests(directory string) error {
	// clear test list
	t.tests = make(map[string]*Test)

	// read test's directory
	testsFileList, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Println("Error: cannot list directory: " + directory)
		log.Fatal(err)
	}

	// check that no files in the test's directory
	for _, curTestFile := range testsFileList {
		if !curTestFile.IsDir() {
			log.Fatal("File " + curTestFile.Name() + " in tests directory: " + directory + ". There is should exist folders only.")
		}
	}

	// create new items in t.tests,
	for _, curTestFile := range testsFileList {
		t.tests[curTestFile.Name()] = &Test{
			directory: directory,
			name:      curTestFile.Name(),
			params:    map[string]TestParameter{},
			tags:      []string{},
		}
	}

	// load allowed parameters of each test
	for _, curTest := range t.tests {
		curTest.LoadTagsAndParams()
	}
	return nil
}

// Saves in t.tests the tests with a suitable name only
func (t *TestManager) filterTestsByName() error {
	name := t.parametersFromCLI["name"]
	for curTestName, _ := range t.tests {
		match, err := regexp.MatchString(name, curTestName)
		if err != nil {
			log.Fatal("Incorrect regexp in name parameter")
		}
		if !match {
			delete(t.tests, curTestName)
		}
	}
	return nil
}

// Saves in t.tests the tests with a suitable tag only
func (t *TestManager) filterTestsByTag() error {
	tag := t.parametersFromCLI["tag"]
	for curTestName, curTest := range t.tests {
		tagFound := false
		for _, curTag := range curTest.tags {
			match, err := regexp.MatchString(tag, curTag)
			if err != nil {
				log.Fatal("Incorrect regexp in name parameter")
			}
			if match {
				tagFound = true
				break
			}
		}
		if !tagFound {
			delete(t.tests, curTestName)
		}
	}
	return nil
}

// Print to STDOUT list of tests ordered by name
func (t *TestManager) PrintListOrderedByNames() error {
	for _, curTest := range t.tests {
		fmt.Println(curTest.name)
		fmt.Println("-------------------------------------------------")

		for _, curParam := range curTest.params {
			fmt.Printf("%-20s %-20s", curParam.Name, curParam.Type)
			if curParam.Type == "EnumParam" {
				for _, curEnumVariant := range curParam.Variants {
					fmt.Print(curEnumVariant + " ")
				}
			}
			fmt.Println()
		}
	}
	return nil
}

// Print to STDOUT list of tests ordered by tag
func (t *TestManager) PrintListOrderedByTag() error {
	var allTags map[string][]string
	allTags = make(map[string][]string)

	for _, curTest := range t.tests {
		for _, curTag := range curTest.tags {
			allTags[curTag] = append(allTags[curTag], curTest.name)
		}
	}
	for curTagKey, curTagTests := range allTags {
		fmt.Printf("%s(%d)\r\n", curTagKey, len(curTagTests))
		for _, curTagTest := range curTagTests {
			fmt.Println(curTagTest)
		}

		fmt.Println()
	}
	return nil
}

// Print to STDOUT list of tests ordered by parameter
func (t *TestManager) PrintListOrderedByParameter() error {
	var allParametersTests map[string][]string
	allParametersTests = make(map[string][]string)

	var allParametersVariants map[string][]string
	allParametersVariants = make(map[string][]string)

	for _, curTest := range t.tests {
		for _, curParameter := range curTest.params {
			allParametersTests[curParameter.Name] = append(allParametersTests[curParameter.Name], curTest.name)
			if curParameter.Type == "EnumParam" {
				for _, curVariant := range curParameter.Variants {
					if !arrayUtils.StringInSlice(curVariant, allParametersVariants[curParameter.Name]) {
						allParametersVariants[curParameter.Name] = append(allParametersVariants[curParameter.Name], curVariant)
					}
				}
			}
		}
	}

	for curParameterKey, curParameter := range allParametersTests {
		fmt.Print(curParameterKey)
		if len(allParametersVariants[curParameterKey]) > 1 {
			fmt.Print("(")
			for curVariantKey, curVariant := range allParametersVariants[curParameterKey] {
				fmt.Print(curVariant)
				if curVariantKey < len(allParametersVariants[curParameterKey])-1 {
					fmt.Print(",")
				}
			}
			fmt.Print(")")
		}
		fmt.Println()
		fmt.Println("-------------------------------------------------")
		for _, curParameterTest := range curParameter {
			if t.tests[curParameterTest].params[curParameterKey].Type == "EnumParam" {
				fmt.Println(curParameterTest, t.tests[curParameterTest].params[curParameterKey].Variants)
			} else {
				fmt.Println(curParameterTest)
			}

		}
		fmt.Println()
	}
	return nil
}

// return all params with all variants for each
func (t *TestManager) getAllTestParamsWithVariants() map[string][]string {
	var allParameters map[string][]string
	allParameters = make(map[string][]string)
	// collect all params with all variants for each
	for _, curTest := range t.tests {
		for _, curParameter := range curTest.params {
			if curParameter.Type == "EnumParam" { // if parameter's type is enum - get single parameter from CLI
				for _, curVariant := range curParameter.Variants {
					if !arrayUtils.StringInSlice(curVariant, allParameters[curParameter.Name]) {
						allParameters[curParameter.Name] = append(allParameters[curParameter.Name], curVariant)
					}
				}
			} else { // if parameter's type is string - get single parameter from CLI
				allParameters[curParameter.Name] = []string{t.parametersFromCLI[curParameter.Name]}
			}
		}
	}
	return allParameters
}

func (t *TestManager) filterParametersCombinationsByGlobalParams(paramCombinations []*map[string]string) []*map[string]string {
	var result []*map[string]string
	//fmt.Println(t.testParametersFromCLI)
	//fmt.Println("")
	for _, curCombine := range paramCombinations {
		//fmt.Print(*curCombine)
		combineAccepted := true
		for curParamName, curParamValue := range *curCombine {
			if _, ok := t.testParametersFromCLI[curParamName]; ok { // if parameter found in CLI - check that it is accepted.
				if !arrayUtils.StringInSlice(curParamValue, CLIParser.GetAllVariantsOfFlagSeparatedBy(t.testParametersFromCLI[curParamName], ",")) {
					combineAccepted = false
					break
				}
			}
		}
		if combineAccepted {
			//result= append(result,curCombine)
			//fmt.Print(" accepted")
			result = append(result, curCombine)
		}
		//fmt.Println()

	}
	//os.Exit(0)
	return result
}

func (t *TestManager) getAllCases(ParamCombinations []*map[string]string) [][]string {
	var result [][]string
	for _, curTest := range t.tests {
		casesOfCurTest := curTest.GetCasesByParameterCombinations(ParamCombinations)
		//fmt.Println(casesOfCurTest)
		result = append(result, casesOfCurTest...)

		//os.Exit(0)
		//result = append(result, curTest.GetCasesByParameterCombination(ParamCombinations))
	}
	return result
}

func (t *TestManager) isRequiredParametersPresented() bool {
	var allParamsRequiredToTesting map[string]string
	allParamsRequiredToTesting = make(map[string]string)

	for _, curTest := range t.tests {
		for _, curParameter := range curTest.params {
			if curParameter.Type == "StringParam" {
				allParamsRequiredToTesting[curParameter.Name] = ""
			}
		}
	}

	allRequiredFlagsPresented := true
	for curParameterKey, _ := range allParamsRequiredToTesting {
		if _, ok := t.parametersFromCLI[curParameterKey]; !ok {
			println("Flag \"" + curParameterKey + "\" is required.")
			allRequiredFlagsPresented = false
		}
	}

	if !allRequiredFlagsPresented {
		return false
	} else {
		return true
	}
}

// return array with cases are allowed for this parameters combination
func (t *TestManager) AllCases() [][]string {
	if !t.isRequiredParametersPresented() {
		os.Exit(1)
	}
	allParameters := t.getAllTestParamsWithVariants()
	allParametersCombinations := getAllParamsCombinations(allParameters)
	allParametersCombinations = t.filterParametersCombinationsByGlobalParams(allParametersCombinations)

	return t.getAllCases(allParametersCombinations)
}

// Print to STDOUT all cases are allowed for this parameters combination
func (t *TestManager) PrintCases() error {
	allCases := t.AllCases()
	for _, curCase := range allCases {
		for _, curElement := range curCase {
			fmt.Print(curElement, " ")
		}
		fmt.Println()
	}
	return nil
}
