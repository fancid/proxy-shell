package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"h12.io/socks"
)

// Proxy struct holds address, protocol, and latency metadata
type Proxy struct {
	Address  string
	Protocol string
	Latency  time.Duration
}

func main() {
	fmt.Println("== Proxy Shell Generator ==")
	timeStart := time.Now()
	proxyList := fetchCycle()
	validProxies := checkList(proxyList)
	fmt.Printf("\nTotal Valid Proxies: %d\n", len(validProxies))
	rankedProxies := rankProxies(validProxies)
	fmt.Println("\nTop 10 Fastest Proxies:")
	for i, p := range rankedProxies {
		if i >= 10 {
			break
		}
		fmt.Printf("%d. %s (%s) - %.2f seconds\n", i+1, p.Address, p.Protocol, p.Latency.Seconds())
	}
	fmt.Printf("\nExecution Time: %.2f seconds\n", time.Since(timeStart).Seconds())
	// TODO: Add Proxy Router
	// TODO: Test Proxy Router
}

var apiURL = []string{
	"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt",
	"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks4.txt",
	"https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/socks5.txt",
	"https://raw.githubusercontent.com/hookzof/socks5_list/master/proxy.txt",
	"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt",
	"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks4.txt",
	"https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/socks5.txt",
	"https://raw.githubusercontent.com/mmpx12/proxy-list/master/http.txt",
	"https://raw.githubusercontent.com/mmpx12/proxy-list/master/https.txt",
	"https://raw.githubusercontent.com/mmpx12/proxy-list/master/socks4.txt",
	"https://raw.githubusercontent.com/mmpx12/proxy-list/master/socks5.txt",
	"https://raw.githubusercontent.com/clarketm/proxy-list/master/proxy-list-raw.txt",
	"https://raw.githubusercontent.com/shiftytr/proxy-list/master/proxy.txt",
	"https://raw.githubusercontent.com/sunny9577/proxy-scraper/master/proxies.txt",
	"https://raw.githubusercontent.com/roosterkid/openproxylist/main/HTTPS_RAW.txt",
	"https://raw.githubusercontent.com/roosterkid/openproxylist/main/SOCKS4_RAW.txt",
	"https://raw.githubusercontent.com/roosterkid/openproxylist/main/SOCKS5_RAW.txt",
	"https://raw.githubusercontent.com/OpsXCQ/proxy-list/master/list.txt",
	"https://raw.githubusercontent.com/rdavydov/proxy-list/master/proxies/http.txt",
	"https://raw.githubusercontent.com/rdavydov/proxy-list/master/proxies/socks4.txt",
	"https://raw.githubusercontent.com/rdavydov/proxy-list/master/proxies/socks5.txt",
	"https://raw.githubusercontent.com/Vann-Dev/proxy-list/main/proxies/http.txt",
	"https://raw.githubusercontent.com/Vann-Dev/proxy-list/main/proxies/https.txt",
	"https://raw.githubusercontent.com/Vann-Dev/proxy-list/main/proxies/socks4.txt",
	"https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/http.txt",
	"https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/https.txt",
	"https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/master/socks4.txt",
	"https://raw.githubusercontent.com/prxchk/proxy-list/main/http.txt",
	"https://raw.githubusercontent.com/prxchk/proxy-list/main/socks4.txt",
	"https://raw.githubusercontent.com/prxchk/proxy-list/main/socks5.txt",
	"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/http.txt",
	"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/https.txt",
	"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/socks4.txt",
	"https://raw.githubusercontent.com/ErcinDedeoglu/proxies/main/proxies/socks5.txt",
}

// fetchCycle iterates over all api urls, calling fetch and inferProtocol
func fetchCycle() []Proxy {
	proxyMap := make(map[string]Proxy)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, url := range apiURL {
		wg.Add(1) // Increment Counter

		go func(url string) {
			defer wg.Done()

			fmt.Printf("[*] Fetching from: %s\n", url)

			results, err := fetch(url)
			if err != nil {
				fmt.Printf("[X] Skipping source due to error: %v\n", err)
				return
			}

			protocol := inferProtocol(url)

			mu.Lock()
			before := len(proxyMap)
			for _, p := range results {
				if p != "" {
					proxyMap[p] = Proxy{
						Address:  p,
						Protocol: protocol,
					}
				}
			}
			added := len(proxyMap) - before
			mu.Unlock()

			if added == 0 {
				fmt.Printf("[X] No Proxies Found in %s\n", url)
			} else {
				fmt.Printf("[✓] Successfully added %d proxies from %s \n", added, url)
			}
		}(url)
	}

	wg.Wait()

	fmt.Printf("\nTotal Proxies Fetched: %d\n", len(proxyMap))

	// Convert map to slice
	list := make([]Proxy, 0, len(proxyMap))
	for _, p := range proxyMap {
		list = append(list, p)
	}
	return list
}

// fetch parses the proxies from the url passed to it
func fetch(targetURL string) ([]string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Get(targetURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[X] Failed to fetch %s: status %d", targetURL, resp.StatusCode)
	}

	var proxies []string
	scanner := bufio.NewScanner(resp.Body)

	const maxCapacity = 1024 * 1024 // 1 MB buffer
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			proxies = append(proxies, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("[X] error reading response: %v", err)
	}

	return proxies, nil
}

// inferProtocol determines proxy protocol from the source URL
func inferProtocol(sourceURL string) string {
	lower := strings.ToLower(sourceURL)
	switch {
	case strings.Contains(lower, "socks5"):
		return "socks5"
	case strings.Contains(lower, "socks4"):
		return "socks4"
	case strings.Contains(lower, "https"):
		return "https"
	default:
		return "http"
	}
}

// checkList iterates through the proxylist to gather diagnostics on the proxies
func checkList(Proxies []Proxy) []Proxy {
	var (
		mu    sync.Mutex
		valid []Proxy
		wg    sync.WaitGroup
		sem   = make(chan struct{}, 1000)
	)

	for _, p := range Proxies {
		wg.Add(1)
		sem <- struct{}{}

		go func(p Proxy) {
			defer wg.Done()
			defer func() { <-sem }()

			if checkProxy(p) == 0 {
				return
			}

			var totalLatency time.Duration

			for i := 0; i < 5; i++ {
				latency := checkProxy(p)
				if latency == 0 {
					return
				}
				fmt.Printf("%s responded in %.2f seconds\n", p.Address, latency.Seconds())
				totalLatency += latency
			}

			p.Latency = totalLatency / 5

			mu.Lock()
			valid = append(valid, p)
			mu.Unlock()
		}(p)
	}

	wg.Wait()
	return valid
}

// checkProxy packages checkCurl and checkTCP together
func checkProxy(p Proxy) time.Duration {
	if !checkTCP(p.Address) {
		return 0
	}
	return checkCurl(p)
}

// checkTCP checks if the proxy can access TCP
func checkTCP(address string) bool {
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// checkCurl checks if the proxy can curl 8.8.8.8
func checkCurl(p Proxy) time.Duration {
	testURL := "http://httpbin.org/get"
	timeout := 3 * time.Second

	transport := &http.Transport{}

	switch p.Protocol {
	case "socks4", "socks5":
		proxyAddr := fmt.Sprintf("%s://%s", p.Protocol, p.Address)
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return socks.Dial(proxyAddr)(network, addr)
		}
	case "http", "https":
		proxyURL, _ := url.Parse(fmt.Sprintf("http://%s", p.Address))
		transport.Proxy = http.ProxyURL(proxyURL)

	default:
		return 0
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}

	start := time.Now()
	resp, err := client.Get(testURL)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return time.Since(start)
	}

	return 0
}

// rankProxies sorts the proxies by latency using merge sort
func rankProxies(proxies []Proxy) []Proxy {
	if len(proxies) <= 1 {
		return proxies
	}

	mid := len(proxies) / 2
	left := rankProxies(proxies[:mid])
	right := rankProxies(proxies[mid:])

	return merge(left, right)
}

func merge(left, right []Proxy) []Proxy {
	result := make([]Proxy, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		if left[i].Latency <= right[j].Latency {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}
