package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func getService(name string) (func(w *WatchPage, resp *http.Response) (string, error), bool) {
	services := map[string]func(*WatchPage, *http.Response) (string, error){
		"Web": ServiceWeb,
		"Rss": ServiceRSS,
	}
	result, ok := services[name]
	return result, ok
}

func ServiceWeb(w *WatchPage, resp *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", errors.New("Could not create document from response: " + err.Error())
	}
	fmt.Println("Got result")
	result := toJSON(doc.Find(w.Settings["Selector"]).Map(func(i int, s *goquery.Selection) string {
		html, _ := s.Html()
		return html
	}))
	return result, nil
}

func ServiceRSS(w *WatchPage, resp *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", errors.New("Could not create document from response: " + err.Error())
	}

	result := toJSON(doc.Find("item").Map(func(i int, s *goquery.Selection) string {
		html, _ := s.Html()
		return html
	}))
	return result, nil
}
