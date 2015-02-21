package server

import (
	"fmt"
	"net/http"
	"time"
)

type WatchPage struct {
	Name     string
	RootUrl  string
	Url      string
	Type     string
	Settings map[string]string
	Selector string
	Filter   string

	Error     string
	Value     string
	UpdatedAt time.Time

	skip     bool
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
		fmt.Println(w.Url)
		defer resp.Body.Close()

		// fmt.Printf("%#v\n", resp.StatusCode)
		if resp.Header.Get("X-From-Cache") != "1" {

			service, ok := getService(w.Type)
			if !ok {
				w.Error = fmt.Sprintf("Service \"%s\" does not exist", w.Type)
				return
			}

			result, err := service(w, resp)
			if err != nil {
				w.Error = err.Error()
				return
			}

			if result != w.Value {
				w.Value = result
				w.UpdatedAt = time.Now()
			}
		}
	}
}
