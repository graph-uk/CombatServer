package webPage

import (
	"malibu/Tests/Tests_shared/aas/browser"
	//"TestBot/aas"
	"errors"
	"log"
	"strings"
	"time"
)

type WebPage struct {
	Browser     *browser.Browser
	Hostname    string
	RelativeUrl string
}

func (w *WebPage) SetBrowser(browserLink *browser.Browser) {
	w.Browser = browserLink
}

func (w *WebPage) SetHostname(NewHostname string) {
	w.Hostname = NewHostname
}

func (w *WebPage) SetRelativeUrl(newUrl string) {
	w.RelativeUrl = newUrl
}

func (w *WebPage) Open() error {
	err := w.Browser.GetUrl(w.Hostname + w.RelativeUrl)
	log.Println(w.Hostname + w.RelativeUrl)
	if err != nil {
		log.Println(errors.New("Cannot upload part of archive"))
		log.Println(err)
		return err
	}
	time.Sleep(3 * time.Second)
	currentUrl, err := w.Browser.CurrentUrl()
	if err != nil {
		log.Println(errors.New("Cannot get current url"))
		log.Println(err)
		return err
	}
	if !strings.Contains(currentUrl, w.RelativeUrl) {
		return errors.New("Cannot open diskpage")
	}

	return nil
}
