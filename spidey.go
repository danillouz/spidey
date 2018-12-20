package spidey

import (
	"io"
	"net/http"
)

// FetcherService represents a service to fetch the child links of a URL.
type FetcherService interface {
	Fetch(url string) (*http.Response, error)
	FindLinks(r io.Reader, base, filter string) []string
}

// CacheService represents a service to keep track of visited links.
type CacheService interface {
	Add(url string)
	Check(url string) bool
}

// CrawlerService represents a service to crawl a webpage identified by an URL.
type CrawlerService interface {
	Crawl(url, filter string, depth int)
}
