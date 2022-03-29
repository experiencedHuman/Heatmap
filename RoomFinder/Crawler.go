package RoomFinder

import (
	"fmt"
	"sync"
)

type mutexCache struct {
	mux   sync.Mutex
	store map[string]bool
}

func (cache *mutexCache) setVisited(name string) bool {
	cache.mux.Lock()
	defer cache.mux.Unlock()

	if cache.store[name] {
		return true
	}

	cache.store[name] = true

	return false
}

var cacheInstance = mutexCache{store: make(map[string]bool)}

// Fetcher interface
type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func crawlInner(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup) {
	defer wg.Done()
	if depth <= 0 {
		return
	}

	if cacheInstance.setVisited(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		wg.Add(1)
		go crawlInner(u, depth-1, fetcher, wg)
	}
	return
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(1)

	go crawlInner(url, depth, fetcher, waitGroup)

	waitGroup.Wait()
}
