package mainPage

import (
	"malibu/Tests/Tests_shared/aas/browser"
	"malibu/Tests/Tests_shared/aas/pages/webPage"
	"runtime/debug"
	//"TestBot/tracer"
	//	"log"
)

type MainPage struct {
	webPage.WebPage
}

func (m *MainPage) Init(browser *browser.Browser, hostname string) {
	m.SetBrowser(browser)
	m.SetHostname(hostname)
	m.SetRelativeUrl("/")
}

func (m *MainPage) Header_ClickLogIn() error {
	loginLink, err := m.Browser.ElementFindByTagParamValue("a", "href", "/login")
	if err != nil {
		m.Browser.Log += "\tWarning: cannot click to element \r\n"
		return nil
	}
	m.Browser.ElementClick(loginLink)
	return nil
}

func (m *MainPage) Header_ClickSignUp() error {
	loginLink, err := m.Browser.ElementFindByTagParamValue("a", "href", "/signup")
	if err != nil {
		m.Browser.Log += "\tWarning: cannot click to element \r\n"
		return nil
	}
	m.Browser.ElementClick(loginLink)
	return nil
}

func (m *MainPage) Header_ClickResetPassword(isPanic bool) error {
	loginLink, err := m.Browser.ElementFindByTagParamValue("a", "href", "/account/reset-password")
	if err != nil {
		if isPanic {
			m.Browser.Log += string(debug.Stack())
			panic("Cannot find element")
		} else {
			m.Browser.Log += "\tWarning: cannot click to element \r\n"
			return err
		}
		//return nil
	}
	m.Browser.ElementClick(loginLink)
	return nil
}
