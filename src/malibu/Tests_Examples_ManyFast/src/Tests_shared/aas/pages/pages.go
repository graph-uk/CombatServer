package pages

import (
	"malibu/Tests/Tests_shared/aas/browser"
	"malibu/Tests/Tests_shared/aas/pages/SignUpPage"
	"malibu/Tests/Tests_shared/aas/pages/mainPage"
)

type Pages struct {
	MainPage   mainPage.MainPage
	SignUpPage SignUpPage.SignUpPage
}

func (p *Pages) Init(browser *browser.Browser, hostname string) {
	p.MainPage.Init(browser, hostname)
	p.SignUpPage.Init(browser, hostname)
}
