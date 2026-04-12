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
	apiURL := []string{
		"https://api.proxyscrape.com/v2/?request=getproxies&protocol=http&timeout=10000&country=all&ssl=all&anonymity=all",
		"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt",
		"https://raw.githubusercontent.com/shiftytr/proxy-list/master/proxy.txt",
		"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt",
	}

	proxyMap := make(map[string]bool)

	for _, url := range apiURL {
		fmt.Printf("Fetching from: %s\n", url)

		results, err := fetch(url)
		if err != nil {
			fmt.Printf("  Skipping source due to error: %v\n", err)
			continue
		}

		// Add each found proxy to our map
		for _, p := range results {
			if p != "" { // Ensure we don't add empty strings
				proxyMap[p] = true
			}
		}
		fmt.Printf("  Successfully added %d proxies.\n", len(results))
	}

	// Print the final grouped list
	fmt.Println("\n--- Proxy List ---")
	for proxy := range proxyMap {
		fmt.Println(proxy)
	}

	fmt.Printf("\nProxies: %d\n", len(proxyMap))
	// TODO: Add Proxy Checker
	// TODO: Test Proxy Checker
	// TODO: Add Proxy Router
	// TODO: Test Proxy Router
}
