package main

import (
	"flag"
	"log"

	"github.com/danillouz/spidey/pkg/cache"
	"github.com/danillouz/spidey/pkg/crawler"
	"github.com/danillouz/spidey/pkg/fetcher"
)

func main() {
	urip := flag.String("uri", "", "The URI of the web page to crawl. For example https://news.ycombinator.com/news.")
	filterp := flag.String("filter", "", "Only the links that have this prefix will be crawled. For example /news.")
	depthp := flag.Int("depth", 1, "How many levels of child links you'd like to crawl. For example 6")
	flag.Parse()

	if *urip == "" {
		log.Fatal("Provide a URI to crawl a webpage.")
	}

	f := fetcher.New()
	c := cache.New()
	cr := crawler.New(f, c)

	cr.Crawl(*urip, *filterp, *depthp)
}
