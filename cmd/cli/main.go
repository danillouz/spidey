package main

import (
	"flag"
	"log"

	"github.com/danillouz/spidey/pkg/cache"
	"github.com/danillouz/spidey/pkg/crawler"
	"github.com/danillouz/spidey/pkg/fetcher"
)

func main() {
	urip := flag.String("uri", "", "The URI of the page to crawl. For example 'https://news.ycombinator.com/news'.")
	filterp := flag.String("filter", "", "Only these links will be crawled. For example '/news'.")
	depthp := flag.Int("depth", 1, "How many levels of child links will be crawled. For example '6'")
	flag.Parse()

	if *urip == "" {
		log.Fatal("A URI must be provided in order to crawl a webpage.")
	}

	f := fetcher.New()
	c := cache.New()
	cr := crawler.New(f, c)

	cr.Crawl(*urip, *filterp, *depthp)
}
