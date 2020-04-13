package MalibuWeb

import (
	"Tests_shared/MSTeams/Activities"
	"Tests_shared/browser"
	"runtime"
	"time"
)

type MalibuWeb struct {
	HostName   string
	Browser    *browser.Browser
	Activities Activities.Activities
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func NewMalibuWeb(hostname string, browser *browser.Browser) *MalibuWeb {
	return &MalibuWeb{hostname, browser, Activities.Activities{}}
}

func (t *MalibuWeb) OpenRootURL() *Activities.PageLogin {
	check(t.Browser.Selenium.Get(t.HostName + `/`))
	return t.Activities.NewPageLogin(t.Browser)
}
