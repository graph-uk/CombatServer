package malibutest

import (
	"Tests_shared/aTest"
	//"Tests_shared/browser"
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

// func (t *MalibuTest) loadSeleniumSession() {
// 	t.existSeleniumSessionID = ``
// 	f, err := os.Open(`out` + string(os.PathSeparator) + `SeleniumSessionID.txt`)
// 	if err == nil {
// 		buf := bytes.NewBuffer(nil)
// 		_, err := io.Copy(buf, f) // Error handling elided for brevity.
// 		f.Close()
// 		if err == nil {
// 			t.existSeleniumSessionID = strings.TrimSpace(string(buf.Bytes()))
// 		}
// 	}
// }

// func (t *MalibuTest) saveSeleniumSession() {
// 	SeleniumSessionID := []byte(t.Browser.Selenium.SessionID())
// 	err := ioutil.WriteFile(`out`+string(os.PathSeparator)+`SeleniumSessionID.txt`, SeleniumSessionID, 0644)
// 	if err != nil {
// 		fmt.Println(`Cannot write selenium session ID to file. Error: ` + err.Error())
// 	}
// }

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

	//result.loadSeleniumSession()
	fmt.Println("ExistSeleniumSessionID: " + result.existSeleniumSessionID)
	result.ATest.CreateOutputFolder()
	//result.Browser, err = browser.NewBrowser(result.existSeleniumSessionID) // attach to exist session, or create new if not exist.
	//check(err)
	//result.saveSeleniumSession()

	// result.Params.CSIPassword.Value = decodeParam(result.Params.CSIPassword.Value)
	// result.Params.CSIPassword1.Value = decodeParam(result.Params.CSIPassword1.Value)
	// result.Params.CSIPassword2.Value = decodeParam(result.Params.CSIPassword2.Value)
	// result.Params.CSIPassword3.Value = decodeParam(result.Params.CSIPassword3.Value)
	// result.Params.CSIPassword4.Value = decodeParam(result.Params.CSIPassword4.Value)
	// result.Params.CSIPassword5.Value = decodeParam(result.Params.CSIPassword5.Value)
	// result.Params.CSIPassword6.Value = decodeParam(result.Params.CSIPassword6.Value)

	return &result
}
