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
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"time"
	"golang.org/x/net/html"

	"github.com/framps/golang_tutorial/sitemap/links"
)

const outputName = "genSitemap"
const lastSeenTimeout = time.Second * 1 // timeout for workers when there is no more work
const httpClientTimeout = 30 * time.Second // http get timeout
const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:105.0) Gecko/20100101 Firefox/105.0"

// const userAgent = "Mozilla/5.0 (X11; Linux x86_64; rv:106.0) Gecko/20100101 Firefox/106.0"

const cpyRght1 = "Copyright © 2017,2022 framp at linux-tips-and-tricks dot de"
const cpyRght2 = "Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan"

var (
	matchFile   *os.File // receives list of domain urls
	skippedFile *os.File // receives list of skipped urls and the skip reason
	errorFile   *os.File // receives list of pages unable to retrieve
	remoteFile  *os.File // receives list of remote links
	remoteFileNotfound  *os.File // receives list of remote links which cannot be resolved
)

var (
	debugFlag  *bool
	workerFlag *int
)

var crawled, notfound, matches, fails, skipped, remotes, errors int // counter for matched and skipped urls

var aborted = false

type linkRef struct {
	parent string
	link   string
}

func (l linkRef) String() string {
	return fmt.Sprintf("%s <- %s", l.link,l.parent)
}
// Extract makes an HTTP GET request to the specified URL, parses
// the response as HTML, and returns the links in the HTML document.
func Extract(url string) ([]string, error) {

	client := http.Client{
		Timeout: httpClientTimeout,
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("%d: %s (get error)", err, url)
	}

	var resp *http.Response
	if err == nil {
		req.Header.Set("User-Agent", userAgent)
		resp, err = client.Do(req)
	}

	if err != nil {
		return nil, fmt.Errorf("%d: %s (do error)", err, url)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("%d: %s (status code)",  resp.StatusCode, url)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %s (parse error)", err, url)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

// filter urls via a regex
func isValid(u *url.URL) bool {

	if len(u.Fragment) > 0 { // no fragment allowed
		return false
	}
	if len(u.Query()) > 0 { // no query allowed
		return false
	}
	re := regexp.MustCompile(`(?i).*(\.(htm(l)?|jp(e)?g|mp4|pdf|sql|sh|py|go|mp4|img|png))?$`)
	m := re.MatchString(u.Path)
	return m
}

// crawl urls
func crawl(nr int, parseURL linkRef, sourceURLs []string) []string {

	crawled++

	pu, err := url.Parse(parseURL.link)
	if err != nil {
		m := fmt.Sprintf("%s for %s (%2d) (parse1)\n", err, parseURL, nr)
		fails++
		errorFile.WriteString(m)
		return []string{}
	}

	for _, k := range sourceURLs {

		if parseURL.link != k {
			su, e := url.Parse(k)
			if e != nil {
				m := fmt.Sprintf("%s for %s (%2d) (parse2)\n",  err, k,nr)
				fails++
				errorFile.WriteString(m)
				return []string{}
			}
			if pu.Hostname() != su.Hostname() || pu.Scheme != su.Scheme {

				var m string
				client := http.Client{
					Timeout: httpClientTimeout,
				}

				req, err := http.NewRequest("GET", parseURL.link, nil)

				var resp *http.Response
				if err == nil {
					req.Header.Set("User-Agent", userAgent)
					resp, err = client.Do(req)
				}

				if err != nil {
					m = fmt.Sprintf("%s for %s (%2d)\n",  err, parseURL,nr)
					notfound++
					remoteFileNotfound.WriteString(m)
					return []string{}
				} else if resp.StatusCode != http.StatusOK {
					resp.Body.Close()
					m = fmt.Sprintf("%s for %s (%2d)\n", resp.Status, parseURL,nr)
					notfound++
					remoteFileNotfound.WriteString(m)
					return []string{}
				} else {
					m = fmt.Sprintf("%s (%2d)\n",  parseURL,nr)
				remotes++
				remoteFile.WriteString(m)
				return []string{}
			}
			}
			if !isValid(pu) {
				m := fmt.Sprintf("%s (%2d)\n", parseURL, nr)
				skipped++
				skippedFile.WriteString(m)
				return []string{}
			}
		}
	}

	if *debugFlag {
		fmt.Printf("%2d: --- Crawling %s\n", nr, parseURL)
	} else {
		// fmt.Printf(".")
		fmt.Printf("Pages crawled: %d\r",crawled)
	}

	list, err := links.Extract(parseURL.link)

	/*
		if *debugFlag {
			for _, u := range list {
				fmt.Printf("%2d: +++ Crawled %s\n", nr, u)
			}
		}
	*/

	if err != nil {
		m := fmt.Sprintf("%s for %s (%2d) (extract)\n", err, parseURL,nr)
		errors++
		errorFile.WriteString(m)
		return []string{}
	}

	var url = strings.TrimSuffix(parseURL.link, "/")
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

	remoteFileNotfound, e = os.Create(outputName + ".notfound")
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
			aborted=true
			activeWorkers.Wait()
			break
		}
	}()

	worklist := make(chan []linkRef)  // lists of URLs, may have duplicates
	unseenLinks := make(chan linkRef) // de-duplicated URLs

	sourceURLs := args // first arg are the domains to crawl

	// Add command-line arguments to worklist.
	go func() {
		var l []linkRef
		for i := range sourceURLs {
			l = append(l, linkRef{"", sourceURLs[i]})
		}
		worklist <- l
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
					var l []linkRef
					for _, ll := range foundLinks {
						l = append(l, linkRef{link.link, ll})
					}
					go func() {
						worklist <- l
					}()
				case <-time.After(lastSeenTimeout): // timer will expire if there is no more work to do
					if *debugFlag {
						fmt.Printf("%2d: --- Worker idle for %s and now terminating\n", nr, lastSeenTimeout)
					}
					activeWorkers.Done()
					return
				case <-abort:
					if *debugFlag {
						fmt.Printf("%2d: ??? Worker aborted and now terminating\n", nr)
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
				if !seen[link.link] {
					seen[link.link] = true
					unseenLinks <- link
				}
			}
		}
	}()

	fmt.Printf("%d workers now crawling ...\n", *workerFlag)
	activeWorkers.Wait()

	if ! aborted {
		elapsed := time.Since(start)
		fmt.Printf("\n%d pages crawled in %s on %s", crawled,elapsed,sourceURLs)

		fmt.Printf(`
Pages found: %d (-> %s)
Pages skipped: %d (-> %s)
Remote pages: %d (-> %s)
Remote pages not found: %d (-> %s)
Invalid URLs: %d (-> %s)
`,
			matches, matchFile.Name(),skipped, skippedFile.Name(),remotes, remoteFile.Name(),
			notfound, remoteFileNotfound.Name(),errors, errorFile.Name())
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
