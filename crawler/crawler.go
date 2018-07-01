// Package crawler implements a simple, concurrent web crawler.
package crawler

import (
	"net/http"
)

// HTTPClient is used to retrieve URLs.
type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

// Logger defines what will be used for logging.
type Logger interface {
	Printf(format string, v ...interface{})
}

// siteMapEntry holds the result for one page.
type siteMapEntry struct {
	url, title string
}

// Crawler holds the state of the web crawler.
type Crawler struct {
	Workers  int
	StartURL string

	// Max number of pages to fetch. Zero means no limit.
	Limit int

	Client HTTPClient
	Logger Logger
}
