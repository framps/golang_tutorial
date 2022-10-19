// Copyright (C) 2017 framp at linux-tips-and-tricks dot de
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
//
// Sample code how to crawl a website and generate a list of urls. Can be used
// to generate Google sitemaps.
//
// See github.com/framps/golang_tutorial for latest code
//
// This code is based on https://github.com/adonovan/gopl.io/blob/master/ch8/crawl3/findlinks.go
// and was enhanced by
// 1. a parse termination condition
// 2. remote urls are skipped
// 3. additional messages for better understanding of the crawl logic

// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 243.

// Crawl3 crawls web links starting with the command-line arguments.

package main

import (
	"flag"
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

const outputName = "genSitemap"
const lastSeenTimeout = time.Second * 3 // timeout for workers when there is no more work

const cpyRght1 = "Copyright © 2017 framp at linux-tips-and-tricks dot de"
const cpyRght2 = "Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan"

var (
	matchFile   *os.File // receives list of domain urls
	skippedFile *os.File // receives list of skipped urls and the skip reason
	errorFile   *os.File // receives list of pages unable to retrieve
	remoteFile  *os.File // receives list of remote links
)

var (
	debugFlag  *bool
	workerFlag *int
)

var matches, fails, skipped, remotes, errors int // counter for matched and skipped urls

// filter urls via a regex
func isValid(u *url.URL) bool {

	if len(u.Fragment) > 0 { // no fragment allowed
		return false
	}
	if len(u.Query()) > 0 { // no query allowed
		return false
	}
	re := regexp.MustCompile(`(?i).*(\.(htm(l)?|jp(e)?g|mp4|pdf|sql|sh|py|go))?$`)
	m := re.MatchString(u.Path)
	return m
}

// crawl urls
func crawl(nr int, parseURL string, sourceURLs []string) []string {

	pu, err := url.Parse(parseURL)
	if err != nil {
		m := fmt.Sprintf("%2d: ??? URL parse error %s for %s\n", nr, err, parseURL)
		fails++
		errorFile.WriteString(m)
		return []string{}
	}

	for _, k := range sourceURLs {

		if parseURL != k {
			su, e := url.Parse(k)
			if e != nil {
				m := fmt.Sprintf("%2d: ??? URL parse error %s for %s\n", nr, err, k)
				fails++
				errorFile.WriteString(m)
				return []string{}
			}
			if pu.Hostname() != su.Hostname() || pu.Scheme != su.Scheme {
				m := fmt.Sprintf("%2d: --- Remote URL %s\n", nr, parseURL)
				remotes++
				remoteFile.WriteString(m)
				return []string{}
			}
			if !isValid(pu) {
				m := fmt.Sprintf("%2d: --- No match %s\n", nr, parseURL)
				skipped++
				skippedFile.WriteString(m)
				return []string{}
			}
		}
	}

	if *debugFlag {
		fmt.Printf("%2d: Crawling %s\n", nr, parseURL)
	} else {
		fmt.Printf(".")
	}

	list, err := links.Extract(parseURL)

	if err != nil {
		m := fmt.Sprintf("%2d: ??? Extract error %s for %s\n", nr, err, parseURL)
		errors++
		errorFile.WriteString(m)
		return []string{}
	}

	var url = strings.TrimSuffix(parseURL, "/")
	matchFile.WriteString(url + "\n")
	matches++

	return list
}

// parse a website and create a file with all same domain url links. Create a file which will log skipped urls and the skip reason

func main() {

	fmt.Println(cpyRght1)
	fmt.Println(cpyRght2)
	fmt.Println()

	debugFlag = flag.Bool("debug", false, "Debug mode")
	workerFlag = flag.Int("worker", 20, "Number of workers")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Missing URL to parse")
		os.Exit(1)
	}

	fmt.Printf("Crawling %s\n", args[0])

	var activeWorkers sync.WaitGroup // waitgroup for active workers
	abort := make(chan bool, 1)      // channel to signal abort to worker

	var e error
	matchFile, e = os.Create(outputName + ".match")
	if e != nil {
		panic(e)
	}
	defer func() {
		matchFile.Close()
	}()

	skippedFile, e = os.Create(outputName + ".skipped")
	if e != nil {
		panic(e)
	}
	defer func() {
		skippedFile.Close()
	}()

	errorFile, e = os.Create(outputName + ".error")
	if e != nil {
		panic(e)
	}
	defer func() {
		errorFile.Close()
	}()

	remoteFile, e = os.Create(outputName + ".remote")
	if e != nil {
		panic(e)
	}
	defer func() {
		remoteFile.Close()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			close(abort)
			fmt.Println("\nReceived an interrupt, stopping ...")
			activeWorkers.Wait()
		}
	}()

	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	sourceURLs := args // first arg are the domains to crawl

	// Add command-line arguments to worklist.
	go func() {
		worklist <- os.Args[1:]
	}()

	start := time.Now()

	activeWorkers.Add(*workerFlag)
	// Create crawler goroutines to fetch each unseen link.
	for i := 0; i < *workerFlag; i++ {
		go func(nr int, abort chan bool) {
			for {
				select {
				case link := <-unseenLinks:
					foundLinks := crawl(nr, link, sourceURLs)
					go func() {
						worklist <- foundLinks
					}()
				case <-time.After(lastSeenTimeout): // timer will expire if there is no more work to do
					if *debugFlag {
						fmt.Printf("%2d: Worker idle for %s and now terminating\n", nr, lastSeenTimeout)
					}
					activeWorkers.Done()
					return
				case <-abort:
					if *debugFlag {
						fmt.Printf("%2d: Worker aborted and now terminating\n", nr)
					}
					activeWorkers.Done()
					return
				}
			}
		}(i, abort)
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

	fmt.Printf("Waiting for %d workers to finish ...\n", *workerFlag)
	activeWorkers.Wait()
	fmt.Printf("\nPages found: %d\nPages skipped: %d\nRemote pages: %d\nInvalid URLs %d\n", matches, skipped, remotes, errors)

	elapsed := time.Since(start)
	fmt.Printf("Sitemap creation for %s took %s\n", sourceURLs, elapsed)

}
