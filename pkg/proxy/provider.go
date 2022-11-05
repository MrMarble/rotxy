package proxy

import "math/rand"

type ProxyIterator func() string

// Provider is a iterator builder that returns a new iterator using the given algorithm.
func Provider(algorithm string, proxyList *ProxyList) ProxyIterator {
	switch algorithm {
	case "random":
		return Random(proxyList)
	case "round-robin":
		return RoundRobin(proxyList)
	default:
		return Random(proxyList)
	}
}

// Random returns a new iterator that returns a random proxy from the given list.
func Random(proxiList *ProxyList) func() string {
	return func() string {
		return proxiList.Proxies[rand.Intn(len(proxiList.Proxies))]
	}
}

// RoundRobin returns a new iterator that returns a proxy from the given list in a round-robin fashion.
func RoundRobin(proxiList *ProxyList) func() string {
	var index int
	return func() string {
		proxy := proxiList.Proxies[index]
		index = (index + 1) % len(proxiList.Proxies)
		return proxy
	}
}
