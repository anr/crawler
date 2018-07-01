package crawler

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// absURL receives an absolute base URL and a candidate URL, possibly relative.
// It returns an absolute URL and whether they have the same Host.
func absURL(baseURL, candidateURL string) (string, bool, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", false, fmt.Errorf("couldn't parse base url %q: %v", baseURL, err)
	}

	candidate, err := url.Parse(candidateURL)
	if err != nil {
		return "", false, fmt.Errorf("couldn't parse candidate url %q: %v", candidateURL, err)
	}

	// Make the candidate URL absolute if needed.
	candidate = base.ResolveReference(candidate)

	if candidate.Host != base.Host {
		return candidate.String(), false, nil
	}

	return candidate.String(), true, nil
}

// processURL fetches a URL and extracts all its links. It returns the page's
// entry for the site map and the internal links found.
func (c *Crawler) processURL(myURL string) (*siteMapEntry, []string, error) {
	resp, err := c.Client.Get(myURL)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get page %q: %v", myURL, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("Bad status for %q: %s", myURL, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to read body for %q: %v", myURL, err)
	}

	entry := &siteMapEntry{
		url:   myURL,
		title: doc.Find("title").Text(),
	}

	links := make([]string, 0)
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		item, _ := s.Attr("href")
		u, ok, err := absURL(myURL, item)
		if err != nil {
			c.Logger.Printf("Skipping url %q: %v", item, err)
			return
		}
		if ok {
			links = append(links, u)
		}
	})

	return entry, links, nil
}
