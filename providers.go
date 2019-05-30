package soxy

import (
	"regexp"
	"strings"
	"sync"
)

// Providers maps the colloqial names of proxy providers to the function that returns their proxies.
var Providers = map[string]func() []Proxy{
	"freeproxylists": FreeProxyLists,
}

// FreeProxyLists returns all the HTTP proxies that it can find on the http://www.freeproxylists.com/ website.
func FreeProxyLists() (proxies []Proxy) {

	initialLinks := FindLinks("http://www.freeproxylists.com/socks.html", `^socks #\d+`)

	links := []string{}

	// Vomit, proxy regex
	proxyRegex := regexp.MustCompile(`&lt;td&gt;\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}&lt;\/td&gt;&lt;td&gt;\d+&lt;\/td&gt;`)
	// End Vomit

	for _, link := range initialLinks {
		parsedLink := "http://freeproxylists.com/load_socks_" + link[6:]
		links = append(links, parsedLink)
	}

	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()

			doc, err := GetURL(l)
			if err != nil {
				return
			}

			matches := proxyRegex.FindAllString(doc, -1)

			for i := range matches {
				proxies = append(
					proxies,
					NewProxy(
						matches[i][strings.Index(matches[i], "&lt;td&gt;") + len("&lt;td&gt;"):strings.Index(matches[i], "&lt;/td&gt;")],
						matches[i][strings.LastIndex(matches[i], "&lt;td&gt;") + len("&lt;td&gt;"):strings.LastIndex(matches[i], "&lt;/td&gt;")],
					))
			}
		}(link)
	}

	wg.Wait()

	return proxies
}
