package crawler

import (
	"fmt"
	"log"

	"github.com/danillouz/spidey"
)

// Ensure Crawler satisfies the CrawlerService interface.
var _ spidey.CrawlerService = &Crawler{}

// New returns a Crawler which can fetch and cache links.
func New(f spidey.FetcherService, c spidey.CacheService) *Crawler {
	return &Crawler{
		fetcher: f,
		cache:   c,
	}
}

// Crawler can fetch and cache links.
type Crawler struct {
	fetcher spidey.FetcherService
	cache   spidey.CacheService
}

// Crawl follows the links of a root URL to a certain depth.
// Child URLs can be filtered by providing a filter prefix.
func (c *Crawler) Crawl(url, filter string, depth int) {
	if depth <= 0 {
		fmt.Println("\tdone")
		return
	}

	fmt.Println("crawling", url)

	resp, err := c.fetcher.Fetch(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	c.cache.Add(url)

	links := c.fetcher.FindLinks(resp.Body, url, filter)
	total := len(links) - 1
	done := make(chan bool)

	for i, l := range links {
		fmt.Printf("\t>> %v/%v processing %v\n", i, total, l)

		if ok := c.cache.Check(l); ok {
			fmt.Println("\t!! cached", l)
			continue
		}

		go func(url string) {
			c.Crawl(url, filter, depth-1)
			done <- true
		}(l)
	}

	for i, l := range links {
		fmt.Printf("\t<< %v/%v waiting %v\n", i, total, l)
		<-done
	}
}
