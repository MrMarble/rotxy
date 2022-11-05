package proxy

import (
	"net/http"
	"net/url"
	"time"
)

// Check makes a CONNECT request to the proxy and filters out the ones that
// don't work.
func Check(proxies []string, timeout time.Duration, tls bool) []string {
	newProxies := []string{}
	ch := make(chan string)
	client := &http.Client{Timeout: timeout}
	testURL := "https://httpbin.org"
	if !tls {
		testURL = "http://httpbin.org"
	}

	for _, proxy := range proxies {
		go check(proxy, ch, client, testURL)
	}

	for range proxies {
		proxy := <-ch
		if proxy != "" {
			newProxies = append(newProxies, proxy)
		}
	}

	return newProxies
}

func check(proxy string, ch chan string, client *http.Client, testURL string) {
	proxyURL, err := url.Parse("http://" + proxy)
	if err != nil {
		ch <- ""
		return
	}
	client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	req, err := http.NewRequest(http.MethodConnect, testURL, nil)
	if err != nil {
		ch <- ""
		return
	}
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != 200 {
		ch <- ""
		return
	}
	ch <- proxy
}
