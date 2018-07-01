package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/anr/crawler/crawler"
)

func main() {
	startURL := flag.String("start_url", "", "starting point")
	workers := flag.Int("workers", 1, "number of concurrent workers")
	timeout := flag.Int("timeout", 5, "timeout in seconds")
	limit := flag.Int("limit", 10, "max number of pages to visit")

	flag.Parse()

	logger := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

	if *startURL == "" {
		logger.Printf("Please provide a start url")
		os.Exit(1)
	}

	c := &crawler.Crawler{
		StartURL: *startURL,
		Workers:  *workers,
		Limit:    *limit,
		Client:   &http.Client{Timeout: time.Duration(*timeout) * time.Second},
		Logger:   logger,
	}

	start := time.Now()
	siteMap := c.Run()
	logger.Printf("Execution took %s", time.Since(start).String())

	fmt.Println("Site map")
	for k, v := range siteMap {
		fmt.Printf("title: %s\nurl: %s\n\n", v, k)
	}
}
