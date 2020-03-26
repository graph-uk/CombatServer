package SignUpPage

import (
	"combat/Tests/Tests_shared/aas/browser"
	"combat/Tests/Tests_shared/aas/pages/webPage"
	//"runtime/debug"
	//"TestBot/tracer"
	//	"log"
)

type SignUpPage struct {
	webPage.WebPage
}

func (m *SignUpPage) Init(browser *browser.Browser, hostname string) {
	m.SetBrowser(browser)
	m.SetHostname(hostname)
	m.SetRelativeUrl("/signup")
}

func (m *SignUpPage) FillFirstName(text string) {
	firstNameInput, err := m.Browser.ElementFindByTagParamValue("input", "name", "firstname")
	if err != nil {
		panic("Cannot find element")
	}
	err = m.Browser.ElementTypeTo(firstNameInput, text)
	if err != nil {
		panic("Cannot type to element")
	}
}
