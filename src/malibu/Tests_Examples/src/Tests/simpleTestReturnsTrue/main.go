package main

import (
	"Tests_shared/aTest"
	"io/ioutil"
	"log"
)

type theTest struct {
	aTest  aTest.ATest
	params struct {
		HostName         aTest.StringParam
		SessionTimestamp aTest.StringParam
		Locale           aTest.EnumParam
		AdminName        aTest.StringParam
	}
}

func createNewTest() (*theTest, error) {
	var result theTest

	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "EN")
	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "RU")
	result.params.AdminName.Value = "TestDefaultValue"
	result.aTest.Tags = append(result.aTest.Tags, "NotForLive")

	result.aTest.FillParamsFromCLI(&result.params)
	result.aTest.CreateOutputFolder()
	return &result, nil
}

func main() {
	_, err := createNewTest()
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./out/log.txt", []byte("Ok"), 0777)
	if err != nil {
		println(err.Error())
	}

	log.Println("ok")
	return
}
