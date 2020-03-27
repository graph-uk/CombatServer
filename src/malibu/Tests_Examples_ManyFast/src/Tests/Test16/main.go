package main

import (
	"Tests_shared/aTest"
	"io/ioutil"
	"log"
)

type theTest struct {
	aTest  aTest.ATest
	params struct {
		Locale  aTest.EnumParam
		Locale2 aTest.EnumParam
		Locale3 aTest.EnumParam
	}
}

func createNewTest() (*theTest, error) {
	var result theTest
	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "EN")
	result.params.Locale.AcceptedValues = append(result.params.Locale.AcceptedValues, "RU")
	result.params.Locale2.AcceptedValues = append(result.params.Locale2.AcceptedValues, "qw")
	result.params.Locale2.AcceptedValues = append(result.params.Locale2.AcceptedValues, "we")
	result.params.Locale3.AcceptedValues = append(result.params.Locale3.AcceptedValues, "er")
	result.params.Locale3.AcceptedValues = append(result.params.Locale3.AcceptedValues, "rt")
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

	err = ioutil.WriteFile("./out/log.txt", []byte("fail"), 0777)
	if err != nil {
		println(err.Error())
	}
	panic("sdf")

	log.Println("ok")
	return
}
