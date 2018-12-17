package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/html"
)

// Links contains all links of a crawled page.
type Links []string

func main() {
	flag.Parse()
	urls := flag.Args()

	if len(urls) < 1 {
		log.Fatal("provide an URL, for example https://news.ycombinator.com/news")
	}

	url := urls[0]
	fmt.Printf("crawling %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	links := Links{}
	z := html.NewTokenizer(resp.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			fmt.Println("\nfound links:", links)
			return
		case tt == html.StartTagToken:
			t := z.Token()

			if t.Data == "a" {
				for _, attr := range t.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
						break
					}
				}
			}
		}
	}
}
