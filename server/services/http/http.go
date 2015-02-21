package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/peterbourgon/diskv"
)

type Watcher struct {
	Pages   []*WatchPage
	Name    string
	RootUrl string
	client  *http.Client
	quit    chan struct{}
	wg      sync.WaitGroup
}

func (w *Watcher) init() {
	w.quit = make(chan struct{})

	cache := diskcache.NewWithDiskv(diskv.New(diskv.Options{
		BasePath:     "./cache/" + w.Name,
		CacheSizeMax: 10 * 1024 * 1024, // 10 MB
	}))

	w.client = &http.Client{
		Transport: httpcache.NewTransport(cache),
		// CheckRedirect: redirectPolicyFunc,
	}

	// req, err := http.NewRequest("GET", url, nil)
	// w = &Watcher{
	// 	Pages:  pages,
	// 	Name:   name,
	// 	client: client,
	// 	quit:   quit,
	// }
	for _, page := range w.Pages {
		ticker := time.NewTicker(page.interval)
		req, err := http.NewRequest("GET", page.Url, nil)
		if err != nil {
			fmt.Print(err)
			continue
		}
		page.req = req
		w.wg.Add(1)
		go func() {
			for {
				select {
				case <-ticker.C:
					get3(w.client, page.req)
					// do stuff
				case <-w.quit:
					ticker.Stop()
					w.wg.Done()
					return
				}
			}
		}()
	}
	return
}

func (w *Watcher) Load() string {
	return toJSON(w)
}
func (w *Watcher) Update(since time.Time) string {
	data := map[string]interface{}{
		"Name":    w.Name,
		"RootUrl": w.RootUrl,
	}
	pages := []*WatchPage{}
	for _, page := range w.Pages {
		if page.UpdatedAt.After(since) {
			pages = append(pages, page)
		}
	}
	data["Pages"] = pages
	return toJSON(data)
}
func jsonError(err error) string {
	return toJSON(map[string]string{
		"Error": err.Error(),
	})
}
func jsonDataAndError(data interface{}, err error) {

}

type WatchPage struct {
	Name     string
	RootUrl  string
	Url      string
	Selector string
	Filter   string

	Error     string
	Value     string
	UpdatedAt time.Time

	req      *http.Request
	interval time.Duration

	// auth     bool
	// safeHtml bool

	// filterFunc func(*Selection) bool
	// mapperFunc func(*Selection) string
}

// func NewWatchPage(json string) *WatchPage {

// 	return &WatchPage{
// 		RootUrl: "",
// 		Url:     "",
// 		Name: "",
// 		Selector: "",
// 		Filter:   "",

// 		// Value   string
// 		// UpdatedAt time.Time
// 	}

// }
func (w *WatchPage) checkForUpdates(client *http.Client) {
	resp, err := client.Do(w.req)
	if err != nil {
		w.Error = "Could not complete request: " + err.Error()
		return
	} else {
		defer resp.Body.Close()

		// fmt.Printf("%#v\n", resp.StatusCode)
		if resp.Header.Get("X-From-Cache") != "1" {
			doc, err := goquery.NewDocumentFromResponse(resp)
			if err != nil {
				w.Error = "Could not create document from response: " + err.Error()
				return
			}

			result := toJSON(doc.Find(w.Selector).Map(func(i int, s *goquery.Selection) string {
				html, _ := s.Html()
				return html
			}))
			if result != w.Value {
				w.Value = result
				w.UpdatedAt = time.Now()
			}
		}
	}
}
func toJSON(v interface{}) string {
	bytes, err := json.Marshal(v)
	log.Println(err.Error())
	return string(bytes)
}

func ParseServices() {
	type Services struct {
		Http []*Watcher
		Rss  []*Watcher
	}
	var services Services
	bytes, _ := ioutil.ReadFile("../config/services.json")
	fmt.Println(string(bytes))
	// return string(bytes), err
	err := json.Unmarshal(bytes, &services)
	if err != nil {
		fmt.Println("error:", err)
	}
	a, _ := json.MarshalIndent(services, "", "  ")
	fmt.Println(string(a))
	// fmt.Printf("%+v", services)
}

func registerServices() {

}

// func NewWatchPage() *WatchPage {

// }

func get2(url, selector string) {
	doc, err := goquery.NewDocument("http://fmi.golang.bg/tasks")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		// band := s.Find("h3").Text()
		// title := s.Find("i").Text()
		html, err := s.Html()
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("Review %d: %s %s\n", i, s.Text(), html)
	})
}

func get3(c *http.Client, req *http.Request) {
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("Could not get " + ": " + err.Error())
	} else {
		defer resp.Body.Close()
		fmt.Printf("%#v\n", resp.StatusCode)
		fmt.Printf("%#v\n", resp.Header.Get("X-From-Cache"))
	}
}
