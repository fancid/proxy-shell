package main

import (
	"bufio"
	"fmt"
	"net/http"
)

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
	// TODO: Add Proxy Grabber
	// TODO: Test Proxy Grabber
	apiURL := "https://api.proxyscrape.com/v2/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all"

	fmt.Println("Fetching proxies...")
	proxies, err := fetch(apiURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print the list
	for i, proxy := range proxies {
		fmt.Printf("[%d] %s\n", i+1, proxy)
	}

	fmt.Printf("\nTotal proxies found: %d\n", len(proxies))
	// TODO: Add Proxy Checker
	// TODO: Test Proxy Checker
	// TODO: Add Proxy Router
	// TODO: Test Proxy Router
}
