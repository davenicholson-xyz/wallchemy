package download

import "sync"

type LinkManager struct {
	links []string
	mu    sync.RWMutex
}

func NewLinkManager() *LinkManager {
	return &LinkManager{
		links: make([]string, 0),
	}
}

func (lm *LinkManager) AddLinks(newLinks []string) {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lm.links = append(lm.links, newLinks...)
}

func (lm *LinkManager) GetLinks() []string {
	return lm.links
}

func (lm *LinkManager) Count() int {
	return len(lm.links)
}
