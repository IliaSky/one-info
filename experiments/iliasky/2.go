package main

import (
	"fmt"
	"log"

	"github.com/PuerkitoBio/goquery"
)

func ExampleScrape() {
	doc, err := goquery.NewDocument("http://metalsucks.net")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".reviews-wrap article .review-rhs").Each(func(i int, s *goquery.Selection) {
		band := s.Find("h3").Text()
		title := s.Find("i").Text()
		fmt.Printf("Review %d: %s - %s\n", i, band, title)
	})
}

func ExampleScrape2() {
	doc, err := goquery.NewDocument("http://fmi.golang.bg/tasks")
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#tasks li").Each(func(i int, s *goquery.Selection) {
		// band := s.Find("h3").Text()
		// title := s.Find("i").Text()
		html, err := s.Html()
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("Review %d: %s %s\n", i, s.Text(), html)
	})
}

func main() {
	ExampleScrape2()
}
