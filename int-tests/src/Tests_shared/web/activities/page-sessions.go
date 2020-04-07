package activities

import (
	"Tests_shared/browser"
)

type PageLogin struct {
	Browser *browser.Browser
}

func (t *Activities) NewPageLogin(browser *browser.Browser) *PageLogin {
	return &PageLogin{browser}
}

func (t *PageLogin) FillEmail(value string) {
	el, err := t.Browser.FindByXpathMayInvisible(`//input[@type='email']`, 0)
	if err != nil {
		check(t.Browser.ClickByXpath(`//div[@id='otherTileText']`))
		el, err = t.Browser.FindByXpathMayInvisible(`//input[@type='email']`, 4)
		check(err)
	}

	check(el.SendKeys(value))
}

func (t *PageLogin) ClickNext() *PageLoginPassword {
	check(t.Browser.ClickByXpath(`//input[@type='submit']`))
	return &PageLoginPassword{t.Browser}
}
