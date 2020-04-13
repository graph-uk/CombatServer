package malibutest

import (
	"Tests_shared/aTest"
	"encoding/base64"
	"fmt"

	"os"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type MalibuTest struct {
	ATest  aTest.ATest
	Params struct {
		HostName aTest.StringParam
	}
	//Browser                *browser.Browser
	Timestamp              string
	existSeleniumSessionID string
}

func (t *MalibuTest) startBrowser(filepath string) string {
	return strings.Replace(filepath, `/`, `\`, -1)
}

func (t *MalibuTest) pathToWindowsFormat(filepath string) string {
	return strings.Replace(filepath, `/`, `\`, -1)
}

func (t *MalibuTest) pathToLinuxFormat(filepath string) string {
	return strings.Replace(filepath, `\`, `/`, -1)
}

//decode parameter as base64, and store result to the same param
func decodeParam(str string) string {
	resBytes, err := base64.StdEncoding.DecodeString(str)
	check(err)
	return string(resBytes)
}

func NewMalibuTest() *MalibuTest {
	var result MalibuTest
	result.ATest.DefaultParams = append(result.ATest.DefaultParams, "-HostName="+os.Getenv(`MALIBU_TEST_HOSTNAME`))

	result.ATest.Tags = append(result.ATest.Tags, "NotForLive")
	result.ATest.FillParamsFromCLI(&result.Params)
	result.Timestamp = time.Now().Format("20060102150405")
	fmt.Println("Timestamp: " + result.Timestamp)

	fmt.Println("ExistSeleniumSessionID: " + result.existSeleniumSessionID)
	result.ATest.CreateOutputFolder()

	return &result
}
