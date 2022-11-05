package proxy

import (
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mrmarble/rotxy/pkg/utils"
)

type ProxyList struct {
	mu      sync.Mutex
	Proxies []string
	sources []string
}

var InternalProxyList = []string{
	"https://www.my-proxy.com/free-proxy-list.html",
	"https://www.my-proxy.com/free-proxy-list-2.html",
	"https://www.my-proxy.com/free-proxy-list-3.html",
	"https://www.my-proxy.com/free-proxy-list-4.html",
	"https://www.my-proxy.com/free-proxy-list-5.html",
	"https://www.my-proxy.com/free-proxy-list-6.html",
	"https://www.my-proxy.com/free-proxy-list-7.html",
	"https://www.my-proxy.com/free-proxy-list-8.html",
	"https://www.my-proxy.com/free-proxy-list-9.html",
	"https://www.my-proxy.com/free-proxy-list-10.html",
	"https://www.my-proxy.com/free-elite-proxy.html",
	"https://spys.me/proxy.txt",
}

var ProxyRegex = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d{1,5}`)

func NewProxyList(sources []string) *ProxyList {
	return &ProxyList{
		sources: sources,
	}
}

func (p *ProxyList) Update() {
	wg := sync.WaitGroup{}
	p.Proxies = []string{}
	for _, source := range p.sources {
		wg.Add(1)
		if strings.HasPrefix(source, "http") {
			go p.updateFromURL(source, &wg)
		} else {
			go p.updateFromFile(source, &wg)
		}
	}
	wg.Wait()
	p.Proxies = RemoveDuplicates(filter(p.Proxies))
}

func (p *ProxyList) Prune(cycles uint, timeout time.Duration, tls bool) {
	checked := []string{}
	for i := 0; i < int(cycles); i++ {
		tmp := Check(utils.FilterBy(p.Proxies, checked), timeout, tls)
		checked = append(checked, tmp...)
	}
	p.Proxies = checked
}

func (p *ProxyList) updateFromFile(source string, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.ReadFile(source)
	if err != nil {
		panic(err)
	}
	proxies := ProxyRegex.FindAllString(string(file), -1)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Proxies = append(p.Proxies, proxies...)
}

func (p *ProxyList) updateFromURL(source string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(source)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	proxies := ProxyRegex.FindAllString(string(html), -1)

	p.mu.Lock()
	defer p.mu.Unlock()
	p.Proxies = append(p.Proxies, proxies...)
}

func filter(proxyList []string) []string {
	newProxyList := []string{}
	for _, proxy := range proxyList {
		if proxy == "" {
			continue
		}

		// Remove country code.
		proxy = strings.Split(proxy, "#")[0]

		newProxyList = append(newProxyList, proxy)
	}
	return newProxyList
}

func RemoveDuplicates(proxyList []string) []string {
	proxyMap := make(map[string]struct{})
	for _, proxy := range proxyList {
		proxyMap[proxy] = struct{}{}
	}
	return MapToSlice(proxyMap)
}

func MapToSlice(proxyMap map[string]struct{}) []string {
	proxyList := []string{}
	for proxy := range proxyMap {
		proxyList = append(proxyList, proxy)
	}
	return proxyList
}
