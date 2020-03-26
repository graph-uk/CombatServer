package main

import (
	"Tests_shared/aTest"
	"io/ioutil"
	"log"
	"os"
)

type theTest struct {
	aTest  aTest.ATest
	params struct {
		HostName         aTest.StringParam
		SessionTimestamp aTest.StringParam
		Locale           aTest.EnumParam
		Resolution       aTest.EnumParam
	}
}

func createNewTest() (*theTest, error) {
	var result theTest

	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "EN")
	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "RU")
	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "US")
	result.params.Resolution.AcceptedValues = append(result.params.Resolution.AcceptedValues, "DesktopView")
	result.params.Resolution.AcceptedValues = append(result.params.Resolution.AcceptedValues, "MobileView")

	//result.params.AdminName.Value = "TestDefaultValue"
	result.aTest.Tags = append(result.aTest.Tags, "NotForLive")
	result.aTest.Tags = append(result.aTest.Tags, "Lynx")
	result.aTest.Tags = append(result.aTest.Tags, "AlwaysFailedTest")

	result.aTest.FillParamsFromCLI(&result.params)
	result.aTest.CreateOutputFolder()
	return &result, nil
}

func main() {
	_, err := createNewTest()
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("out/log.txt", []byte("Fail"), 0777)
	os.Exit(12)
	log.Println("ok")
	return
}
