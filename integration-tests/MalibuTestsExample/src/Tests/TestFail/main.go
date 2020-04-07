package main

import (
	"Tests_shared/aTest"
	"Tests_shared/fakescreenshots"
	"fmt"
	"strconv"
	"time"
	"os"
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

	f, _ := os.Create("out/SeleniumSessionID.txt")
    defer f.Close()

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

	panic(`Failed here for example`)
}
