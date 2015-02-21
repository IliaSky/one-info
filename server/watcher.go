package server

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

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

func (w *Watcher) Init() {
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
	for _, page := range w.Pages {
		page.RootUrl = w.RootUrl
		if page.interval == 0 {
			page.interval = time.Hour
		}
		ticker := time.NewTicker(page.interval)
		req, err := http.NewRequest("GET", page.RootUrl+page.Url, nil)
		if err != nil {
			page.skip = true
			page.Error = err.Error()
			fmt.Println(err)
			continue // skip the invalid urls
		}
		page.req = req
		w.wg.Add(1)
		go func() {
			page.checkForUpdates(w.client)

			for {
				select {
				case <-ticker.C:
					page.checkForUpdates(w.client)
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
func (w *Watcher) Update(since time.Time) map[string]interface{} {
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
	return data
	// return toJSON(data)
}

func (w *Watcher) Stop() {
	close(w.quit)
	w.wg.Wait()
	return
}

// func jsonError(err error) string {
// 	return toJSON(map[string]string{
// 		"Error": err.Error(),
// 	})
// }
// func jsonDataAndError(data interface{}, err error) {

// }

func toJSON(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Println(err.Error())
	}
	return string(bytes)
}
