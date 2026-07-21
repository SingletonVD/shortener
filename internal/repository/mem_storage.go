package repository

import (
	"net/url"
	"sync"

	"github.com/SingletonVD/shortener/internal/model"
)

type MemStorage struct {
	Links map[string]url.URL
	Lock  sync.Mutex
}

func (storage *MemStorage) SaveIfAvailable(link model.ShortenedLink) bool {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	_, found := storage.Links[link.Short]
	if found {
		return false
	}

	storage.Links[link.Short] = link.FullLink

	return true
}

func (storage *MemStorage) FindFullLink(shortLink string) (url.URL, bool) {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	fullLink, found := storage.Links[shortLink]
	return fullLink, found
}
