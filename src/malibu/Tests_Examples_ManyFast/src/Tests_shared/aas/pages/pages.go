package pages

import (
	"combat/Tests/Tests_shared/aas/browser"
	"combat/Tests/Tests_shared/aas/pages/SignUpPage"
	"combat/Tests/Tests_shared/aas/pages/mainPage"
)

type Pages struct {
	MainPage   mainPage.MainPage
	SignUpPage SignUpPage.SignUpPage
}

func (p *Pages) Init(browser *browser.Browser, hostname string) {
	p.MainPage.Init(browser, hostname)
	p.SignUpPage.Init(browser, hostname)
}
