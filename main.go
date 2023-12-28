package main

import (
	"context"
	"crawler/internal/crawler-service"
	"flag"
	"fmt"
	"log"
	"time"
)

const CrawlingTimeoutSeconds = 300

// CrawlWebpage craws the given rootURL looking for <a href=""> tags
// that are targeting the current web page, either via an absolute url like http://mysite.com/mypath or by a relative url like /mypath
// and returns a sorted list of absolute urls  (eg: []string{"http://mysite.com/1","http://mysite.com/2"})
func CrawlWebpage(rootURL string, maxDepth int) ([]string, error) { //nolint: unparam
	cw := crawler.NewWebCrawler() // this must be in the main or dep-container function

	ctx, cancel := context.WithTimeout(context.Background(), CrawlingTimeoutSeconds*time.Second)
	defer cancel()

	links, err := cw.Crawl(ctx, rootURL, maxDepth)
	if err != nil {
		fmt.Println("failed to crawl all data", err)
	}

	// nil is returned for purpose, errors during crawling are expected,
	// this is not that critical and no need to stop execution in the main func
	return links, nil
}

// --- DO NOT MODIFY BELOW ---

func main() {
	const (
		defaultURL      = "https://www.example.com/"
		defaultMaxDepth = 3
	)
	urlFlag := flag.String("url", defaultURL, "the url that you want to crawl")
	maxDepth := flag.Int("depth", defaultMaxDepth, "the maximum number of links deep to traverse")
	flag.Parse()

	links, err := CrawlWebpage(*urlFlag, *maxDepth)
	if err != nil {
		log.Fatalln("ERROR:", err)
	}
	fmt.Println("Links")
	fmt.Println("-----")
	for i, l := range links {
		fmt.Printf("%03d. %s\n", i+1, l)
	}
	fmt.Println()
}
