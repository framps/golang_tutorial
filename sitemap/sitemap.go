// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
//
// Sample code how to crawl a website and generate a list of urls. Can be used
// to generate Google sitemaps.
//
// See github.com/framps/golang_tutorial for latest code
//
// This code is based on https://github.com/adonovan/gopl.io/blob/master/ch8/crawl3/findlinks.go
// and was enhanced by a parse termination condition and parsing only happens on the base url

// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 243.

// Crawl3 crawls web links starting with the command-line arguments.
//
// This version uses bounded parallelism.
// For simplicity, it does not address the termination problem.
//
package main

import (
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/framps/golang_tutorial/sitemap/links"
)

const maxWorkers = 20
const lastSeenTimeout = time.Second * 3 // timeout for workers when there is no more work

var (
	matchFile  *os.File
	rejectFile *os.File
)

var matches, fails int

// filter urls via a regex
func isValid(u *url.URL) bool {

	if len(u.Fragment) > 0 { // no fragment allowed
		return false
	}
	if len(u.Query()) > 0 { // no query allowed
		return false
	}
	re := regexp.MustCompile(`(?i).*(\.(html|jpg|jpeg|mp4|htm|pdf|sql))?$`)
	m := re.MatchString(u.Path)
	if m {
		matchFile.WriteString(u.String() + "\n")
		matches++
	}
	return m
}

// crawl urls
func crawl(nr int, parseURL string, sourceURLs []string) []string {

	fmt.Printf("%2d: Crawling %s\n", nr, parseURL)

	pu, err := url.Parse(parseURL)
	if err != nil {
		m := fmt.Sprintf("%2d: ??? URL parse error %s for %s\n", nr, err, parseURL)
		fails++
		rejectFile.WriteString(m)
		return []string{}
	}

	for _, k := range sourceURLs {
		if parseURL != k {
			su, e := url.Parse(k)
			if e != nil {
				m := fmt.Sprintf("%2d: ??? URL parse error %s for %s\n", nr, err, k)
				fails++
				rejectFile.WriteString(m)
				return []string{}
			}
			if !strings.Contains(pu.String(), su.Hostname()) {
				m := fmt.Sprintf("%2d: --- Remote URL %s\n", nr, parseURL)
				fails++
				rejectFile.WriteString(m)
				return []string{}
			}
			if !isValid(pu) {
				m := fmt.Sprintf("%2d: --- No match %s\n", nr, parseURL)
				fails++
				rejectFile.WriteString(m)
				return []string{}
			}
		}
	}
	list, err := links.Extract(parseURL)

	if err != nil {
		m := fmt.Sprintf("%2d: ??? Extract error %s for %s\n", nr, err, parseURL)
		fails++
		rejectFile.WriteString(m)
		return []string{}
	}
	return list
}

// parse a website and create a file with all same domain url links. Create a file which will log skipped urls and the skip reason

func main() {

	var activeWorkers sync.WaitGroup // waitgroup for active workers

	matchFile, err := os.Create("sitemap.match")
	if err != nil {
		panic(err)
	}
	defer func() {
		fmt.Println("Closing matchfile")
		matchFile.Close()
	}()

	rejectFile, err = os.Create("sitemap.reject")
	if err != nil {
		panic(err)
	}
	defer func() {
		fmt.Println("Closing rejectfile")
		rejectFile.Close()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for _ = range signalChan {
			fmt.Println("\nReceived an interrupt, stopping ...")
			rejectFile.Close()
			matchFile.Close()
			fmt.Printf("Found pages: %d\nSkipped pages: %d\n", matches, fails)
			os.Exit(1)
		}
	}()

	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	sourceURLs := os.Args[1:] // first arg are the domains to crawl

	// Add command-line argument to worklist.
	go func() {
		worklist <- os.Args[1:]
	}()

	activeWorkers.Add(maxWorkers)
	// Create crawler goroutines to fetch each unseen link.
	for i := 0; i < maxWorkers; i++ {
		go func(nr int) {
			for {
				select {
				case link := <-unseenLinks:
					foundLinks := crawl(nr, link, sourceURLs)
					go func() {
						worklist <- foundLinks
					}()
				case <-time.After(lastSeenTimeout): // timer will expire if there is no more work to do
					fmt.Printf("%2d: Worker idle for %s and now terminating\n", nr, lastSeenTimeout)
					activeWorkers.Done()
				}
			}
		}(i)
	}

	// de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	go func() {
		seen := make(map[string]bool)
		for list := range worklist {
			for _, link := range list {
				if !seen[link] {
					seen[link] = true
					unseenLinks <- link
				}
			}
		}
	}()

	fmt.Printf("Waiting for %d workers to finish ...\n", maxWorkers)
	activeWorkers.Wait()
	fmt.Printf("Found pages: %d\nSkipped pages: %d\n", matches, fails)
}
