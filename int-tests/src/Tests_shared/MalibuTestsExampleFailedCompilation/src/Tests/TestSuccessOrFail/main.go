package main

import (
	"Tests_shared/aTest"
	"Tests_shared/fakescreenshots"
	"fmt"
	"os"
	"strconv"
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
	for i := 0; i < 100; i++ {
		if i < 10 {
			fakescreenshots.MakeFakeSetpArtifacts(strconv.Itoa(0) + strconv.Itoa(i))
			//time.Sleep(20 * time.Millisecond)
		} else {
			fakescreenshots.MakeFakeSetpArtifacts(strconv.Itoa(i))
			//time.Sleep(20 * time.Millisecond)
		}
	}
	//	fakescreenshots.MakeFakeSetpArtifacts(theTest.timestamp)
	//	time.Sleep(2 * time.Second)
	//	fakescreenshots.MakeFakeSetpArtifacts(time.Now().Format("20060102150405"))
	//	time.Sleep(2 * time.Second)
	//	fakescreenshots.MakeFakeSetpArtifacts(time.Now().Format("20060102150405"))

	//	file := os.TempDir() + `\testSuccessOrFailure.txt`
	//	_, err := os.Stat(file)
	//	if err == nil {

	//	} else if os.IsNotExist(err) {

	//		_, err := os.OpenFile(file, os.O_RDONLY|os.O_CREATE, 0666)
	//		fmt.Println(err.Error())
	//		panic("File not found")
	//	} else {
	//		fmt.Printf("file %s stat error: %v", file, err)
	//	}

	file := os.TempDir() + `\testSuccessOrFailure.txt`
	_, err := os.Stat(file)
	if err != nil {
		panic("Panicking")
	}
}
