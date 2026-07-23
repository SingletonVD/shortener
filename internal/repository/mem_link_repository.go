package repository

import (
	"sync"

	"github.com/SingletonVD/shortener/internal/model"
)

type MemLinkRepository struct {
	links map[string]string
	lock  sync.Mutex
}

func NewMemLinkRepository() *MemLinkRepository {
	return &MemLinkRepository{links: make(map[string]string)}
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

func (storage *MemLinkRepository) FindFullLink(shortLink string) (string, bool) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	fullLink, found := storage.links[shortLink]
	return fullLink, found
}
