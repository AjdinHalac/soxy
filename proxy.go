package soxy

// Proxy represents a URL that can be used.
type Proxy struct {
	Host     string
	Port     string
	Used     bool
}

// NewProxy is a convenience function for generating a new proxy struct.
func NewProxy(Host string, Port string) Proxy {
	return Proxy{Host: Host, Port: Port, Used: false}
}

// UniqueProxies returns a list of proxies with all duplicates removed.
// Instead of just using the plain Proxy type, it only checks if the URL matches.
// This is because a proxy that is identical to another from a different provider will not match because the provider will be different, although fundementally it is the same proxy.
func UniqueProxies(proxies []Proxy) []Proxy {
	keys := make(map[string]bool)
	list := []Proxy{}

	for _, entry := range proxies {
		if _, value := keys[entry.Host]; !value {
			keys[entry.Host] = true
			list = append(list, entry)
		}
	}
	return list
}
