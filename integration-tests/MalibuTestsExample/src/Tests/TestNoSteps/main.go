package main

import (
	"Tests_shared/aTest"
	"fmt"
	"time"
)

type theTest struct {
	aTest  aTest.ATest
	params struct {
		InternalIP aTest.StringParam
	}
	timestamp              string
	existSeleniumSessionID string
}

func createNewTest() *theTest {
	var result theTest
	result.aTest.DefaultParams = append(result.aTest.DefaultParams, "-InternalIP=192.168.1.1")
	result.aTest.Tags = append(result.aTest.Tags, "LiveMonitoring")
	result.aTest.FillParamsFromCLI(&result.params)
	result.timestamp = time.Now().Format("20060102150405")
	fmt.Println("Timestamp: " + result.timestamp)
	result.aTest.CreateOutputFolder()
	return &result
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	theTest := createNewTest()
	defer func() {
		if r := recover(); r != nil {
			aTest.PrintSourceAndContinuePanic(r)
		}
	}()
	fmt.Println(theTest.params.InternalIP.Value)
	fmt.Println(`This test has no steps, so it should not have "steps" tab`)

}
