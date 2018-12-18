package main

import (
	"flag"
	"fmt"
	"io"
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

	linksp := &Links{}
	findLinks(linksp, resp.Body)
	fmt.Println("\nfound links:", *linksp)
}

func findLinks(l *Links, r io.Reader) {
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()
			ok, link := getLink(t)
			if !ok {
				continue
			}

			*l = append(*l, link)
		}
	}
}

func getLink(t html.Token) (ok bool, link string) {
	if t.Data == "a" {
		for _, attr := range t.Attr {
			if attr.Key == "href" {
				ok = true
				link = attr.Val
				break
			}
		}
	}

	return
}
