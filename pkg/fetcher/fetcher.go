package fetcher

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/danillouz/spidey"
	"golang.org/x/net/html"
)

// Ensure Fetcher satisfies the FetcherService interface.
var _ spidey.FetcherService = &Fetcher{}

// New return a Fetcher which retrieves the content of a web page.
func New() *Fetcher {
	return &Fetcher{}
}

// Fetcher can retrieve the contents of a web page and its links.
type Fetcher struct {
}

// Fetch the content of a web page by URL.
func (f *Fetcher) Fetch(url string) (*http.Response, error) {
	return http.Get(url)
}

// FindLinks finds the links of a webpage using an optional filter.
// Relative links will be transformed to absolute links, by using an URL base.
func (f *Fetcher) FindLinks(r io.Reader, base, filter string) []string {
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
