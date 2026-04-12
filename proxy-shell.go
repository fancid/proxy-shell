package main

import (
	"bufio"
	"fmt"
	"net/http"
	"sync"
)

func fetchCycle() []string {
	apiURL := []string{
		"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt",
		"https://raw.githubusercontent.com/shiftytr/proxy-list/master/proxy.txt",
		"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt",
	}

	proxyMap := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, url := range apiURL {
		wg.Add(1) // Increment Count

		go func(url string) {
			defer wg.Done()

			var i int = 0
			fmt.Printf("Fetching from: %s\n", url)

			results, err := fetch(url)
			if err != nil {
				fmt.Printf("  Skipping source due to error: %v\n", err)
				return
			}

			mu.Lock()
			for _, p := range results {
				if p != "" {
					proxyMap[p] = true
					i++
				}
			}
			mu.Unlock()

			if i == 0 {
				fmt.Printf("  No Proxies Found in %s\n", url)
			} else {
				fmt.Printf("  Successfully added %d proxies. \n", i)
			}
		}(url)
	}

	wg.Wait()

	fmt.Printf("\nTotal Proxies Fetched: %d\n", len(proxyMap))

	var list []string
	for p := range proxyMap {
		list = append(list, p)
	}
	return list
}

func fetch(targetURL string) ([]string, error) {
	resp, err := http.Get(targetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var proxies []string
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	return proxies, nil
}

func main() {
	proxies := fetchCycle()

	// TODO: Add Proxy Checker
	// TODO: Test Proxy Checker
	// TODO: Add Proxy Router
	// TODO: Test Proxy Router
}
