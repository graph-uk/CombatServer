package aTest

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// This is the base struct contain all required in all test fields
type ATest struct {
	Tags          []string
	DefaultParams []string
}

// This param able to contain value only in accepted values.
type EnumParam struct {
	Value          string
	AcceptedValues []string
}

// This param able to contain any string values.
type StringParam struct {
	Value string
}

// Return name of the test. (filename without extension)
func (a *ATest) GetName() string {
	return strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
}

// Fill test parameters passed from CLI. If it is clearly passed - print manual and exit.
func (a *ATest) FillParamsFromCLI(params interface{}) error {
	if len(os.Args) < 2 {
		os.Args = append([]string{os.Args[0]}, a.DefaultParams...)
	}
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "params":
			fmt.Println("")
			a.PrintParamsInPlaintext(params)
			fmt.Println("")
			a.PrintTagsInPlaintext()
			fmt.Println("")
			os.Exit(1)
		case "paramsJSON":
			a.PrintParamsAndTagsInJSON(params, a.Tags)
			os.Exit(1)
		default:
			a.FillParamsFromCLI2(params)
		}
	} else {
		a.PrintManual()
		os.Exit(1)
	}
	return nil
}

// Fill test parameters from CLI.
func (a *ATest) FillParamsFromCLI2(params interface{}) error {
	s := reflect.ValueOf(params)

	// parse CLI flags to flagMap
	var flagMap map[string]*string
	flagMap = make(map[string]*string)
	for curFieldIndex := 0; curFieldIndex < s.Elem().NumField(); curFieldIndex++ {
		flagMap[s.Elem().Type().Field(curFieldIndex).Name] = flag.String(s.Elem().Type().Field(curFieldIndex).Name, "", "")
	}
	flag.Parse()

	// set parameters from flagMap
	for curFieldIndex := 0; curFieldIndex < s.Elem().NumField(); curFieldIndex++ {
		if s.Elem().Field(curFieldIndex).Type().Name() == "StringParam" {
			s.Elem().Field(curFieldIndex).FieldByName("Value").Set(reflect.ValueOf(*flagMap[s.Elem().Type().Field(curFieldIndex).Name]))
		}
		if s.Elem().Field(curFieldIndex).Type().Name() == "EnumParam" {
			AcceptedValues := s.Elem().Field(curFieldIndex).FieldByName("AcceptedValues")
			ParamValueAccepted := false
			for i := 0; i < AcceptedValues.Len(); i++ {
				if AcceptedValues.Index(i).String() == *flagMap[s.Elem().Type().Field(curFieldIndex).Name] {
					ParamValueAccepted = true
					break
				}
			}
			if ParamValueAccepted {
				s.Elem().Field(curFieldIndex).FieldByName("Value").Set(reflect.ValueOf(*flagMap[s.Elem().Type().Field(curFieldIndex).Name]))
			} else {
				fmt.Println("Incorrect value " + *flagMap[s.Elem().Type().Field(curFieldIndex).Name] + " for parameter " + s.Elem().Type().Field(curFieldIndex).Name)
				fmt.Println("The parameter accept the only following values:")
				for i := 0; i < AcceptedValues.Len(); i++ {
					fmt.Println(AcceptedValues.Index(i).String())
				}
				os.Exit(1)
			}
		}
	}

	// Print params table
	fmt.Println("##########Starting tests with parameters#############")
	for curFieldIndex := 0; curFieldIndex < s.Elem().NumField(); curFieldIndex++ {
		//		fmt.Printf("%-20s \"%s\"\n", s.Elem().Type().Field(curFieldIndex).Name, *flagMap[s.Elem().Type().Field(curFieldIndex).Name])
	}
	fmt.Println()

	return nil
}

// Print user manual to console.
func (a *ATest) PrintManual() {
	fmt.Println("Incorrect parameter.")
	fmt.Println(a.GetName() + " params                                    - Show list of accepted tags/params in human format")
	fmt.Println(a.GetName() + " paramsJSON                                - Show list of accepted tags/params in JSON format")
	fmt.Println(a.GetName() + " -param1=\"value1\" -param2=\"value2\"...      - Run test with parameters")
}

// Print accepted parameters as human-readable format
func (a *ATest) PrintParamsInPlaintext(params interface{}) {
	fmt.Println("Parameters:")
	fmt.Println("Name                 Type                  Variants")
	fmt.Println("---------------------------------------------------")
	s := reflect.ValueOf(params).Elem()
	for curFieldIndex := 0; curFieldIndex < s.NumField(); curFieldIndex++ {
		curField := s.Field(curFieldIndex)
		if curField.Type().Name() == "StringParam" {
			fmt.Printf("%-20s%-10s\n", s.Type().Field(curFieldIndex).Name, " StringParam")
		}
		if curField.Type().Name() == "EnumParam" {
			//fmt.Print(s.Type().Field(curFieldIndex).Name + " EnumParam ")
			fmt.Printf("%-20s%-10s", s.Type().Field(curFieldIndex).Name, " EnumParam             ")
			for EnumIndex := 0; EnumIndex < curField.Field(1).Len(); EnumIndex++ {
				fmt.Print(curField.Field(1).Index(EnumIndex).String() + " ")
			}
			fmt.Println("")
		}
	}
}

func (a *ATest) PrintTagsInPlaintext() {
	fmt.Println("Tags:")
	for _, curTag := range a.Tags {
		fmt.Print(curTag + " ")
	}
}

func (a *ATest) CreateOutputFolder() {
	os.RemoveAll("out")
	err := os.Mkdir("out", 0777)
	if err != nil {
		fmt.Println("Cannot create output folder")
	}
}

//Function for tracing test step by step to check states manually (in CMS, API, etc)
func (a *ATest) WaitForReturnKey() {
	fmt.Println(`Test paused. Press return to continue`)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
