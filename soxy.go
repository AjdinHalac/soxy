package soxy

import (
	"github.com/pkg/errors"
	"math/rand"
)

type Soxy struct {
	proxies []Proxy
}

func NewSoxy() (soxy Soxy) {
	soxy = Soxy{}
	soxy.Refresh()
	return
}

// Refresh will load all proxies from all providers.
func (b *Soxy) Refresh() {
	for provider := range Providers {
		b.fetchProvider(provider)
	}
}

// Random fetches a random proxy.
func (b *Soxy) Random() Proxy {
	chosen := b.proxies[rand.Intn(len(b.proxies)-1)]
	chosen.Used = true
	return chosen
}

// Unused returns first unused proxy.
func (b *Soxy) Unused() (chosen Proxy,err error) {
	for _, proxy := range b.proxies {
		if !proxy.Used {
			proxy.Used = true
			chosen = proxy
			return
		}
	}
	err =  errors.New("No unused proxies left.")
	return
}

/* Loading & Managing proxies */

// add adds a list of proxies.
func (b *Soxy) add(proxies []Proxy) {
	b.proxies = UniqueProxies(append(b.proxies, proxies...))
}

// fetchProvider fetches all the proxies from a given provider.
func (b *Soxy) fetchProvider(provider string) {
	proxies := Providers[provider]()

	b.add(proxies)
}

