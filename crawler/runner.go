package crawler

import "sync"

// Run starts the crawling.
func (c *Crawler) Run() map[string]string {
	c.Logger.Printf("Starting crawler with %d worker(s)", c.Workers)

	// Maps urls to titles.
	siteMap := make(map[string]string)

	// URLs produced by the workers.
	queue := make(chan []string)

	// Site map entries produced by the workers.
	entries := make(chan *siteMapEntry)

	// Control the number of concurrent page fetches.
	sem := make(chan struct{}, c.Workers)

	// Signal worker goroutines to exit.
	quit := make(chan struct{})

	seen := make(map[string]bool)

	go func() { queue <- []string{c.StartURL} }()

	// The termination condition if the page limit is not hit
	// is zero items in the queue.
	queueItems := 1

	var wg sync.WaitGroup

loop:
	for ; queueItems > 0; queueItems-- {
		urls := <-queue
		for _, url := range urls {
			if c.Limit > 0 && len(seen) == c.Limit {
				break loop
			}

			if !seen[url] {
				seen[url] = true
				queueItems++
				wg.Add(1)

				go func(url string) {
					c.Logger.Printf("Worker processing %q", url)

					sem <- struct{}{}
					entry, urls, err := c.processURL(url)
					<-sem

					if err != nil {
						c.Logger.Printf("Failed to process page: %v", err)
						// We can't just return from here without writing
						// to queue and calling Done() on the wg.
						urls = []string{}
					}

					// If we have exceeded the visited URLs limit,
					// we will block until we can read from quit.
					select {
					case queue <- urls:
					case <-quit:
					}

					go func() {
						defer wg.Done()
						if entry != nil {
							c.Logger.Printf("Publishing entry %s", entry.url)
							entries <- entry
						}
					}()
				}(url)
			}
		}
	}

	// If there are blocked workers, ensure they can continue.
	close(quit)

	go func() {
		// This needs to be done on its own goroutine, otherwise goroutines writing
		// to the entries channel would remain blocked.

		// When all entries are published, we can close the channel so that
		// we reach the final return point on the parent goroutine.
		wg.Wait()
		close(entries)
	}()

	c.Logger.Printf("Collecting results")
	for entry := range entries {
		siteMap[entry.url] = entry.title
	}

	return siteMap
}
