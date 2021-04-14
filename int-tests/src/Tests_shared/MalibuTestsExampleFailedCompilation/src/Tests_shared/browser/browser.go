// provider of reliable versions of selenium operations.
// For example - selenium is not wait while element will be shown. Browser - wait it.
// Also browser try most actions twice, if fails.
// If you want use raw selenium functions, and process errors manual - call browser.Selenium.<...>
// All methods of browser return error, but not panic

package browser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

type Browser struct {
	Selenium           selenium.WebDriver
	WaitInterval       int
	minorFailPostDelay int
	browserCallLog     []*[]uintptr
}

func (b *Browser) setWaitInterval(interval int) error {
	err := b.Selenium.SetImplicitWaitTimeout(time.Second * time.Duration(interval))
	if err != nil {
		return err
	}
	b.WaitInterval = interval
	return nil
}

func ShuffleElements(src []selenium.WebElement) []selenium.WebElement {
	rand.Seed(time.Now().UnixNano())
	dest := make([]selenium.WebElement, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}
	return dest
}

//// existSeleniumSessionID - empty.
func NewBrowser(existSeleniumSessionID string) (*Browser, error) {
	caps := selenium.Capabilities{"browserName": "chrome"} // works chromeOpt
	gaBytes, err := ioutil.ReadFile(`GA.crx`)
	if err != nil {
		panic(`cannot read extension`)
	}
	caps["chromeOptions"] = selenium.ChromeOptions{
		//Args:       []string{`--start-maximized`},
		//specifyAgent, and allow login:pass@host.domain
		Args:       []string{`--user-agent=TestGraphAgent`, `--disable-blink-features=BlockCredentialedSubresources`},
		Extensions: []interface{}{gaBytes},
	}

	sel, err := selenium.NewExistRemote(caps, ``, existSeleniumSessionID)
	if err != nil {
		sel, err = selenium.NewRemote(caps, ``)
		if err != nil {
			return &Browser{nil, 5, 3, nil}, err
		}
	}

	sel.SetImplicitWaitTimeout(5)
	sel.ResizeWindow("", 1280, 1024)
	result := Browser{sel, 5, 3, nil}
	if err != nil {
		return &Browser{nil, 5, 3, nil}, err
	}

	return &result, err
}

func (b *Browser) ClickByXpath(xpath string) error {
	b.MakeScreenshot()
	element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
	if err == nil {
		err = element.Click()
	}
	if err != nil {
		time.Sleep(time.Duration(b.minorFailPostDelay) * time.Second)
		b.Selenium.ExecuteScript(`window.scrollTo(0, 0);`, nil) // scroll to top of page to prevent covering by floating header
		element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
		if err == nil {
			err = element.Click()
		}
		return err
	}
	return err
}

func (b *Browser) ClearByXpath(xpath string) error {
	b.MakeScreenshot()
	element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
	if err == nil {
		err = element.Clear()
	}
	if err != nil {
		time.Sleep(time.Duration(b.minorFailPostDelay) * time.Second)
		element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
		if err == nil {
			err = element.Clear()
		}
		return err
	}
	return err
}

func (b *Browser) ClearByXpathInvisible(xpath string) error {
	b.MakeScreenshot()
	element, err := b.findByXpathMayInvisible(xpath, b.WaitInterval)
	if err == nil {
		err = element.Clear()
	}
	if err != nil {
		time.Sleep(time.Duration(b.minorFailPostDelay) * time.Second)
		element, err := b.findByXpathMayInvisible(xpath, b.WaitInterval)
		if err == nil {
			err = element.Clear()
		}
		return err
	}
	return err
}

func (b *Browser) ClickRandomVisibleElementByXpath(xpath string) error {
	b.MakeScreenshot()
	elements, _ := b.Selenium.FindElements("xpath", xpath) //b.findByXpathDisplayed(xpath, b.waitInterval)
	elements = ShuffleElements(elements)
	for _, curElement := range elements {
		displayed, _ := curElement.IsDisplayed()
		if displayed {
			err := curElement.Click()
			return err
		}
	}
	return errors.New("Cannot find visible element in selection")
}

func (b *Browser) MouseMoveByXpath(xpath string) error {
	b.MakeScreenshot()
	element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
	if err == nil {
		err = element.MoveTo(1, 1)
	}
	if err != nil {
		time.Sleep(time.Duration(b.minorFailPostDelay) * time.Second)
		element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
		if err == nil {
			err = element.MoveTo(1, 1)
		}
		return err
	}
	return err
}

func (b *Browser) GetCookieValue(name string) (string, error) {
	cookies, err := b.Selenium.GetCookies()
	if err == nil {
		for _, curCookie := range cookies {
			if curCookie.Name == name {
				return curCookie.Value, nil
			}
		}
	} else {
		time.Sleep(time.Duration(b.minorFailPostDelay) * time.Second)
		for _, curCookie := range cookies {
			if curCookie.Name == name {
				return curCookie.Value, nil
			}
		}
	}
	return "", err
}

func (b *Browser) SubmitByXpath(xpath string) error {
	b.MakeScreenshot()
	element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
	if err == nil {
		err = element.Submit()
	}
	if err != nil {
		time.Sleep(time.Duration(b.minorFailPostDelay) * time.Second)
		element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
		if err == nil {
			err = element.Submit()
		}
		return err
	}
	return err
}

func (b *Browser) TypeByXpath(xpath string, text string) error {
	b.MakeScreenshot()
	element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
	if err == nil {
		err = element.SendKeys(text)
	}
	if err != nil {
		time.Sleep(time.Duration(b.minorFailPostDelay) * time.Second)
		element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
		if err == nil {
			err = element.SendKeys(text)
		}
		return err
	}
	return err
}

func (b *Browser) TypeByXpathMayInvisible(xpath string, text string) error {
	b.MakeScreenshot()
	element, err := b.findByXpathMayInvisible(xpath, b.WaitInterval)
	if err == nil {
		err = element.SendKeys(text)
	}
	if err != nil {
		time.Sleep(time.Duration(b.minorFailPostDelay) * time.Second)
		element, err := b.findByXpathMayInvisible(xpath, b.WaitInterval)
		if err == nil {
			err = element.SendKeys(text)
		}
		return err
	}
	return err
}

func (b *Browser) CheckUrlContains(urlPart string) error {
	secondsLfet := b.WaitInterval
	for secondsLfet > 0 {
		currentURL, err := b.Selenium.CurrentURL()
		if err != nil {
			return err
		}
		if strings.Contains(currentURL, urlPart) {
			return nil
		}
		time.Sleep(time.Second)
		secondsLfet--
	}
	return errors.New("Url is not contains '" + urlPart + "'")
}

// return nil if exist. Error if not exist.
func (b *Browser) CheckExistByXpath(xpath string) error {
	b.MakeScreenshot()
	_, err := b.findByXpathMayInvisible(xpath, b.WaitInterval)
	if err == nil {
		return nil
	}
	if err != nil {
		time.Sleep(time.Duration(b.minorFailPostDelay) * time.Second)
		_, err := b.findByXpathMayInvisible(xpath, b.WaitInterval)
		if err == nil {
			return nil
		}
		return err
	}
	return err
}

// return nil if exist. Error if not exist.
func (b *Browser) CheckNotExistByXpath(xpath string) error {
	b.MakeScreenshot()
	_, err := b.findByXpathMayInvisible(xpath, 2)
	if err == nil {
		return errors.New(`Element is exist, but shold not. XPath: ` + xpath)
	}
	return nil
}

// return nil if exist. Error if not exist.
func (b *Browser) CheckNotShownByXpath(xpath string) error {
	b.MakeScreenshot()
	element, err := b.findByXpathMayInvisible(xpath, 2)
	if err == nil {
		displayed, err := element.IsDisplayed()
		if err != nil {
			return errors.New(`Element is exist, but shold not. May be an error while check visibility. XPath: ` + xpath)
		}
		if displayed {
			return errors.New(`Element is shown, but shold not. XPath: ` + xpath)
		}
		//
	}
	return nil
}

// return nil if displayed. Error if not displayed.
func (b *Browser) CheckDisplayedByXpath(xpath string) error {
	b.MakeScreenshot()
	displayed := false
	element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
	if err == nil {
		displayed, err = element.IsDisplayed()
	}
	if err != nil {
		time.Sleep(time.Duration(b.minorFailPostDelay) * time.Second)
		element, err := b.FindByXpathDisplayed(xpath, b.WaitInterval)
		if err == nil {
			displayed, err = element.IsDisplayed()
			if err != nil {
				return errors.New("Cannot check is element displayed")
			} else {
				if displayed {
					return nil
				} else {
					return errors.New("Element is found but not displayed")
				}
			}
		}
		return err
	}
	return err
}

func (b *Browser) FindByXpathDisplayed(xpath string, timeoutSeconds int) (selenium.WebElement, error) {
	displayed := false
	element, err := b.Selenium.FindElement("xpath", xpath)
	if err == nil {
		displayed, err = element.IsDisplayed()
	}
	if err != nil || !displayed {
		if timeoutSeconds < 1 {
			if err != nil {
				return nil, err
			} else {
				return nil, errors.New("Element is found but not displayed")
			}
		}
		time.Sleep(1 * time.Second)
		element, err := b.FindByXpathDisplayed(xpath, timeoutSeconds-1)
		if err != nil {
			return nil, err
		}
		return element, nil
	}
	return element, err
}

func (b *Browser) findByXpathMayInvisible(xpath string, timeoutSeconds int) (selenium.WebElement, error) {
	element, err := b.Selenium.FindElement("xpath", xpath)
	if err == nil {
		return element, nil
	} else {
		if timeoutSeconds < 1 {
			if err != nil {
				return nil, err
			}
		}
		time.Sleep(1 * time.Second)
		element, err := b.findByXpathMayInvisible(xpath, timeoutSeconds-1)
		if err != nil {
			return nil, err
		}
		return element, nil
	}
}

func (b *Browser) Close() error {
	b.MakeScreenshot()
	err := b.Selenium.Quit()
	if err != nil {
		return err
	}
	return nil
}

func (b *Browser) MakeScreenshot() {
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	if b.Selenium != nil {
		screenshot, err := b.Selenium.Screenshot()
		if err != nil {
			fmt.Println("cannot make screenshot")
		}
		ioutil.WriteFile("out"+string(os.PathSeparator)+timestamp+".png", screenshot, 0777)

		pageSrc, err := b.Selenium.PageSource()
		if err != nil {
			fmt.Println("cannot get page source")
		}
		ioutil.WriteFile("out"+string(os.PathSeparator)+timestamp+".html", []byte(pageSrc), 0777)

		pageURL, err := b.Selenium.CurrentURL()
		if err != nil {
			fmt.Println("cannot get current URL")
		}
		ioutil.WriteFile("out"+string(os.PathSeparator)+timestamp+".txt", []byte(pageURL), 0777)

		// extract and format errors and warnings from browser console
		log, err := b.Selenium.Log(`browser`)
		if err != nil {
			fmt.Println("cannot get browser log")
			fmt.Println(err.Error())
			return
		}
		consoleMessages := `Timestamp | Level | Message` + "\r\n"
		consoleMessages += `----------------------------------------------------------------------------` + "\r\n"
		for _, curMessage := range log {
			consoleMessages += strconv.Itoa(curMessage.Timestamp) + ` ` + curMessage.Level + ` ` + curMessage.Message + "\r\n"
		}
		ioutil.WriteFile("out"+string(os.PathSeparator)+timestamp+".console.log", []byte(consoleMessages), 0777)

	}
}

func (b *Browser) GetParameterValueByXpath(xpath, paramName string) (string, error) {
	b.MakeScreenshot()
	element, err := b.findByXpathMayInvisible(xpath, b.WaitInterval)
	if err == nil {
		value, err := element.GetAttribute(paramName)
		if err == nil {
			return value, nil
		} else {
			return "", err
		}
	}
	if err != nil {
		time.Sleep(time.Duration(b.minorFailPostDelay) * time.Second)
		element, err := b.findByXpathMayInvisible(xpath, b.WaitInterval)
		if err == nil {
			value, err := element.GetAttribute(paramName)
			if err == nil {
				return value, nil
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	}
	return "", err
}

func (b *Browser) ClearCookies() error {
	b.MakeScreenshot()
	err := b.Selenium.DeleteAllCookies()
	if err != nil {
		return err
	}
	err = b.RunScript(`window.localStorage.clear();`)
	if err != nil {
		return err
	}
	err = b.RunScript(`window.sessionStorage.clear();`)
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	return nil
}

func (b *Browser) RefreshPage() error {
	b.MakeScreenshot()
	err := b.Selenium.Refresh()
	if err != nil {
		return err
	}
	return nil
}

func (b *Browser) RunScript(script string) error {
	b.MakeScreenshot()
	_, err := b.Selenium.ExecuteScript(script, nil)
	if err != nil {
		return err
	}
	return nil
}
