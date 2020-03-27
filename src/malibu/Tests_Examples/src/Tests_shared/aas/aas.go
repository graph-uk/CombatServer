package aas

import (
	"malibu/Tests/Tests_shared/aas/browser"
	"malibu/Tests/Tests_shared/aas/pages"
	//	"errors"
	//	"fmt"
	//	"log"
	//	"strings"
	"time"
)

type AAS struct {
	Browser  browser.Browser
	Pages    pages.Pages
	Url      string
	loggedIn bool
}

func (a *AAS) Init(url string) {
	a.Url = url
	a.Browser.Init("D:\\Environment\\Utils\\ChromeWebDriver\\win32_2.14\\chromedriver.exe")
	a.Browser.StartChromeDriver()
	//a.Pages.Init(&a.Browser, a.Url)
	a.Pages.Init(&a.Browser, url)
}

func (a *AAS) SleepSeconds(delay time.Duration) {
	time.Sleep(delay * time.Second)
}
