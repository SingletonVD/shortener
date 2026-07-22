package repository

import (
	"net/url"
	"sync"

	"github.com/SingletonVD/shortener/internal/model"
)

type MemLinkRepository struct {
	links map[string]url.URL
	lock  sync.Mutex
}

func NewMemLinkRepository() *MemLinkRepository {
	return &MemLinkRepository{links: make(map[string]url.URL)}
}

func (storage *MemLinkRepository) SaveIfAvailable(link model.ShortenedLink) bool {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	_, found := storage.links[link.Short]
	if found {
		return false
	}

	storage.links[link.Short] = link.FullLink

	return true
}

func (storage *MemLinkRepository) FindFullLink(shortLink string) (url.URL, bool) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	fullLink, found := storage.links[shortLink]
	return fullLink, found
}
