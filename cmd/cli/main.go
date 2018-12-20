package main

import (
	"flag"
	"fmt"

	"github.com/danillouz/spidey/pkg/cache"
	"github.com/danillouz/spidey/pkg/crawler"
	"github.com/danillouz/spidey/pkg/fetcher"
)

func main() {
	urip := flag.String("uri", "", "The URI of the page to crawl. For example 'https://news.ycombinator.com/news'.")
	filterp := flag.String("filter", "", "Only these links will be crawled. For example '/news'.")
	flag.Parse()

	fmt.Println("uri", *urip)
	fmt.Println("filter", *filterp)

	f := fetcher.New()
	c := cache.New()
	cr := crawler.New(f, c)

	depth := 3
	cr.Crawl(*urip, *filterp, depth)
}
