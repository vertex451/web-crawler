package crawler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"golang.org/x/net/html"
)

// Read IMPROVEMENTS.md for more context.

// WebCrawler defines a web crawler with a root URL and maximum depth.
type WebCrawler struct {
	visited     map[string]bool
	links       []string
	failedLinks []failedLinks
	history     []history
}

// history is used to store each crawl result
type history struct {
	rootURL     string
	maxDepth    int
	links       []string
	failedLinks []failedLinks
	start       time.Time
	finish      time.Time
}

type failedLinks struct {
	link string
	err  error
}

// NewWebCrawler creates a new WebCrawler instance.
func NewWebCrawler() *WebCrawler {
	return &WebCrawler{
		visited: make(map[string]bool),
	}
}

// Crawl initiates the crawling process.
func (c *WebCrawler) Crawl(ctx context.Context, rootURL string, maxDepth int) (result []string, err error) {
	historyEntry := history{
		start:    time.Now(),
		rootURL:  rootURL,
		maxDepth: maxDepth,
	}

	defer func() {
		c.history = append(c.history, historyEntry)
	}()

	if maxDepth == 0 {
		return nil, nil
	}

	c.crawlURL(ctx, rootURL, 0, maxDepth)
	historyEntry.links = c.links
	historyEntry.failedLinks = c.failedLinks
	historyEntry.finish = time.Now()

	c.links = nil // Reset links for the next crawl
	c.failedLinks = nil
	c.visited = map[string]bool{} // Reset visited for the next crawl

	// in case of at least one error, return the first error value to tell the user that something went wrong.
	if len(historyEntry.failedLinks) > 0 {
		return historyEntry.links, historyEntry.failedLinks[0].err
	}

	return historyEntry.links, nil
}

// crawlURL recursively crawls URLs up to the specified depth.
func (c *WebCrawler) crawlURL(ctx context.Context, url string, depth int, maxDepth int) {
	if depth > maxDepth || c.visited[url] {
		return
	}

	links, err := c.extractLinks(ctx, url)
	if err != nil {
		fmt.Printf("Failed to extracting links from %s: %v\n", url, err)
		c.failedLinks = append(c.failedLinks, failedLinks{
			link: url,
			err:  err,
		})
		return
	}

	c.visited[url] = true
	c.links = append(c.links, url)

	for _, link := range links {
		c.crawlURL(ctx, link, depth+1, maxDepth)
	}

	sort.Strings(c.links)
}

// extractLinks extracts links from an HTML page.
func (c *WebCrawler) extractLinks(ctx context.Context, urlStr string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	baseURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	tokenizer := html.NewTokenizer(resp.Body)
	var links []string

	for {
		tokenType := tokenizer.Next()
		switch tokenType { //nolint: exhaustive // (we don't need other token types)
		case html.ErrorToken:
			return links, nil
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data != "a" {
				continue
			}

			for _, attr := range token.Attr {
				if attr.Key != "href" {
					continue
				}

				linkURL, err := c.resolveURL(baseURL, attr.Val)
				if err == nil {
					links = append(links, linkURL.String())
				}
			}
		}
	}
}

// resolveURL resolves a URL against the base URL.
func (c *WebCrawler) resolveURL(baseURL *url.URL, link string) (*url.URL, error) {
	linkURL, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	if linkURL.IsAbs() {
		return linkURL, nil
	}

	return baseURL.ResolveReference(linkURL), nil
}
