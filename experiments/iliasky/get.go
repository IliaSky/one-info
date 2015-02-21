package main

import (
	"fmt"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/peterbourgon/diskv"
	"io/ioutil"
	"net/http"
	"time"
)

func get(url string) {

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Could not get " + url + ": " + err.Error())
	}
	// fmt.Println(resp)
	fmt.Printf("%#v\n", resp.StatusCode)
	fmt.Printf("%#v\n", resp.Header["Etag"])
	fmt.Printf("%#v\n", resp.Header["Last-Modified"])
	_, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Could not read response body")
	}
	// fmt.Println(_)
}

func get2(c *http.Client, url string) {
	resp, err := c.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Could not get " + url + ": " + err.Error())
	}
	fmt.Printf("%#v\n", resp.StatusCode)
	fmt.Printf("%#v\n", resp.Header["X-From-Cache"])
	//
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

// req, err := http.NewRequest("GET", s.server.URL+"/lastmodified", nil)
// 	resp, err := s.client.Do(req)
// 	defer resp.Body.Close()
// 	c.Assert(err, IsNil)
// 	c.Assert(resp.Header.Get(XFromCache), Equals, "")

// 	resp2, err2 := s.client.Do(req)
//     func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error)

// CheckRedirect specifies the policy for handling redirects.
// If CheckRedirect is not nil, the client calls it before
// following an HTTP redirect. The arguments req and via are
// the upcoming request and the requests made already, oldest
// first. If CheckRedirect returns an error, the Client's Get
// method returns both the previous Response and
// CheckRedirect's error (wrapped in a url.Error) instead of
// issuing the Request req.
//
// If CheckRedirect is nil, the Client uses its default policy,
// which is to stop after 10 consecutive requests.
func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	fmt.Println("SASA?")
	fmt.Println(req)
	return nil
}

func main() {

	// cache := httpcache.NewMemoryCache()

	cache := diskcache.NewWithDiskv(diskv.New(diskv.Options{
		BasePath:     "./cache",
		CacheSizeMax: 10 * 1024 * 1024, // 10 MB
	}))

	c := &http.Client{
		// Transport:     httpcache.NewMemoryCacheTransport(),
		Transport:     httpcache.NewTransport(cache),
		CheckRedirect: redirectPolicyFunc,
	}
	// get("http://fmi.golang.bg/topics/293")
	// get("http://si1.free.bg/")
	// get("http://si1.free.bg/")
	// c := httpcache.NewTransport(cache).Client()
	// c := httpcache.NewMemoryCacheTransport().Client()
	// url := "http://si1.free.bg/"
	// url := "http://fmi.golang.bg/topics/293"
	// url := "http://sinoptik.bg/genrss.php?lid=100727011"
	url := "http://sinoptik.bg/rss/100727011"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Could not create request?")
	}
	// get2(c, "http://fmi.golang.bg/topics/293")
	// for i := 0; i < 9; i++ {
	// 	go get2(c, "http://si1.free.bg/")
	// }
	for i := 0; i < 9; i++ {
		time.Sleep(100 * time.Millisecond)

		go get3(c, req)
	}
	time.Sleep(0 * time.Second)

}
