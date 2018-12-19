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

	depth := 3
	crawl(*urip, *filterp, depth)
}

func crawl(uri, filter string, depth int) {
	if depth < 1 {
		fmt.Println("\tdone")
		return
	}

	fmt.Println("crawling", uri)
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	visited.Add(uri)
	links := findLinks(resp.Body, uri, filter)

	done := make(chan bool)
	for i, l := range links {
		fmt.Printf("\t>> %v/%v processing %v\n", i, len(links), l)

		if ok := visited.Link(l); ok {
			fmt.Println("\t!! cached", l)
			continue
		}

		go func(uri string) {
			crawl(uri, filter, depth-1)
			done <- true
		}(l)
	}

	for i, l := range links {
		fmt.Printf("\t<< %v/%v waiting %v\n", i, len(links), l)
		<-done
	}
}

func findLinks(r io.Reader, base, filter string) []string {
	links := []string{}
	z := html.NewTokenizer(r)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			return links
		case tt == html.StartTagToken:
			t := z.Token()
			link, ok := getLink(t)
			if !ok {
				continue
			}

			if valid := strings.HasPrefix(link, filter); valid {
				abs := getAbsLink(link, base)
				links = append(links, abs)
			}
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
