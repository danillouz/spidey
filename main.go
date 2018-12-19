package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// Cache keeps track of visited links and implements a mutex.
type Cache struct {
	links map[string]bool
	mux   sync.Mutex
}

// Add caches a visited link to prevent loops and crawl more efficiently.
func (c *Cache) Add(l string) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.links[l]; !ok {
		c.links[l] = true
	}
}

// Link checks if a link has been visited already.
func (c *Cache) Link(l string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()

	_, ok := c.links[l]
	return ok
}

var visited = Cache{
	links: map[string]bool{},
}

func main() {
	urip := flag.String("uri", "", "The URI of the page to crawl. For example 'https://news.ycombinator.com/news'.")
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
			link, ok := getLink(t)
			if !ok {
				continue
			}

			go func() {
				if ok := visited.Link(link); ok {
					fmt.Println("\tcached", link)
					return
				}

				if valid := strings.HasPrefix(link, filter); valid {
					visited.Add(link)
					abs := getAbsLink(link, base)
					jobs <- abs
				}
			}()
		}
	}
}

func getLink(t html.Token) (link string, ok bool) {
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
