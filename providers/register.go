package providers

import "sync"

var (
	providers     = make(map[string]Provider)
	providersLock sync.RWMutex
)

func RegisterProvider(p Provider) {
	providersLock.Lock()
	defer providersLock.Unlock()
	providers[p.Name()] = p
}

func GetProvider(name string) (Provider, bool) {
	providersLock.RLock()
	defer providersLock.RUnlock()
	p, exists := providers[name]
	return p, exists
}

func init() {
	RegisterProvider(&WallhavenProvider{})
}
