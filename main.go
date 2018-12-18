package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

var visited = map[string]bool{}

func main() {
	urip := flag.String("uri", "", "The URI of the page to crawl. For example ;'https://news.ycombinator.com/news'.")
	filterp := flag.String("filter", "", "Only these links will be crawled. For example '/news'.")
	flag.Parse()

	fmt.Println("uri", *urip)
	fmt.Println("filter", *filterp)

	jobs := make(chan string)
	go func() {
		jobs <- *urip
	}()

	for uri := range jobs {
		crawl(uri, jobs, *filterp)
	}
}

func crawl(uri string, jobs chan string, filter string) {
	fmt.Println("crawling", uri)

	resp, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	findLinks(resp.Body, jobs, uri, filter)
}

func findLinks(r io.Reader, jobs chan string, base, filter string) {
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

			go func() {
				if ok := visited[link]; ok {
					fmt.Println("cached", link)
					return
				}

				if valid := strings.HasPrefix(link, filter); valid {
					visited[link] = true
					abs := getAbsLink(link, base)
					jobs <- abs
				}
			}()
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

func getAbsLink(link, base string) string {
	l, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
	}

	b, err := url.Parse(base)
	if err != nil {
		log.Fatal(err)
	}

	abs := b.ResolveReference(l)
	return abs.String()
}
