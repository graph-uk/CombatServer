package browser

import (
	"log"
	"time"

	"github.com/fedesog/webdriver"
)

type Browser struct {
	session      *webdriver.Session
	driverHandle *webdriver.ChromeDriver
	driverPath   string
	Log          string
}

func (b *Browser) Init(driverPath string) error {
	b.driverPath = driverPath
	return nil
}

func (b *Browser) StartChromeDriver() error {
	b.driverHandle = webdriver.NewChromeDriver(b.driverPath)

	err := b.driverHandle.Start()
	if err != nil {
		log.Println(err)
		return err
	}
	//	desired := webdriver.Capabilities{"Platform": "Windows"}
	//	required := webdriver.Capabilities{}

	//	b.session, err = b.driverHandle.NewSession(desired, required)
	//	if err != nil {
	//		log.Println(err)
	//		return err
	//	}

	//	err = b.session.SetTimeoutsAsyncScript(10000)
	//	if err != nil {
	//		log.Println(err)
	//		return err
	//	}

	//	err = b.session.SetTimeoutsImplicitWait(10000)
	//	if err != nil {
	//		log.Println(err)
	//		return err
	//	}
	return err
}

func (b *Browser) StartSession() error {
	desired := webdriver.Capabilities{"Platform": "Windows"}
	required := webdriver.Capabilities{}

	var err error
	b.session, err = b.driverHandle.NewSession(desired, required)
	if err != nil {
		log.Println(err)
		return err
	}

	err = b.session.SetTimeoutsAsyncScript(10000)
	if err != nil {
		log.Println(err)
		return err
	}

	err = b.session.SetTimeoutsImplicitWait(10000)
	if err != nil {
		log.Println(err)
		return err
	}

	w, _ := b.session.WindowHandle()

	var windowPosition webdriver.Position
	windowPosition.X = 0
	windowPosition.Y = 0

	var windowSize webdriver.Size
	windowSize.Width = 1280
	windowSize.Height = 800

	w.SetPosition(windowPosition)
	w.SetSize(windowSize)
	return err
}

func (b *Browser) StopSession() error {
	if b.session != nil {
		err := b.session.Delete()
		if err != nil {
			log.Println(err)
			return err
		}
		return err
	} else {
		return nil
	}
}

func (b *Browser) StopChromeDriver() error {
	if b.driverHandle != nil {
		err := b.driverHandle.Stop()
		if err != nil {
			log.Println(err)
			return err
		}
		return err
	} else {
		return nil
	}
}

func (b *Browser) RestartIfDied() error {
	//fmt.Println("GetUrl")

	var err error
	if b.session != nil {
		_, err = b.session.GetUrl()
	}

	if err != nil || b.session == nil {
		err = b.StartSession()
		if err != nil {
			panic("Cannot restart browser session")
		}
		time.Sleep(2 * time.Second)
		//log.Println(err)
	}
	return err
}

func (b *Browser) GetUrl(url string) error {
	//fmt.Println("GetUrl")
	err := b.session.Url(url)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (b *Browser) Refresh() error {
	err := b.session.Refresh()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (b *Browser) CurrentUrl() (string, error) {
	//fmt.Println("GetUrl")
	currentUrl, err := b.session.GetUrl()
	if err != nil {
		log.Println(err)
	}
	return currentUrl, err
}

func (b *Browser) ElementClick(element webdriver.WebElement) error {
	err := b.session.MoveTo(element, 5, 5)
	if err != nil {
		log.Println("CannotMoveTo")
		log.Println(err)
		return err
	}
	time.Sleep(1 * time.Second)
	err = element.Click()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (b *Browser) ElementIsDisplayed(element webdriver.WebElement) (bool, error) {
	displayed, err := element.IsDisplayed()
	if err != nil {
		log.Println(err)
	}
	return displayed, err
}

func (b *Browser) ElementSubmit(element webdriver.WebElement) error {
	err := element.Submit()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (b *Browser) ElementSendFile(element webdriver.WebElement, value string) error {
	err := element.SendKeys(value)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(1 * time.Second)
	return err
}

func (b *Browser) ElementTypeTo(element webdriver.WebElement, value string) error {
	err := element.SendKeys(value)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(1 * time.Second)
	return err
}

func (b *Browser) ElementMoveTo(element webdriver.WebElement, x, y int) error {
	err := b.session.MoveTo(element, x, y)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (b *Browser) ElementFindByID(id string) (webdriver.WebElement, error) {
	el, err := b.session.FindElement("id", id)
	if err != nil {
		log.Println(err)
	}
	return el, err
}

func (b *Browser) ElementFindByClass(class string) (webdriver.WebElement, error) {
	el, err := b.ElementFindByXPath("//*[@class[contains(.,'" + class + "')]]")
	if err != nil {
		log.Println(err)
	}
	return el, err
}

func (b *Browser) ElementFindByTagParamValue(tag string, param string, value string) (webdriver.WebElement, error) {
	el, err := b.ElementFindByXPath("//" + tag + "[@" + param + "[contains(.,'" + value + "')]]")
	return el, err
}

func (b *Browser) ElementFindByXPath(xpath string) (webdriver.WebElement, error) {
	el, err := b.session.FindElement("xpath", xpath)
	if err != nil {
		log.Println(err)
	}
	return el, err
}

func (b *Browser) ElementsFindByXPath(xpath string) ([]webdriver.WebElement, error) {
	elArray, err := b.session.FindElements("xpath", xpath)
	if err != nil {
		log.Println(err)
	}
	return elArray, err
}

func (b *Browser) ClearCookiesAndStorage() error {
	err := b.session.DeleteCookies()
	if err != nil {
		log.Println(err)
		return err
	}
	err = b.session.SessionStorageClear()
	if err != nil {
		log.Println(err)
		return err
	}
	time.Sleep(1 * time.Second)
	return err
}
