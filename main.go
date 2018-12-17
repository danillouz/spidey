package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("body", string(body))
}
