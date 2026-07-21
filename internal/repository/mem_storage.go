package repository

import (
	"net/url"
	"sync"

	"github.com/SingletonVD/shortener/internal/random"
)

const shortLinkLength = 8

type MemStorage struct {
	Links map[string]url.URL
	Lock  sync.Mutex
}

func (storage *MemStorage) CreateShortLink(fullLink url.URL) string {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	shortLink := random.GenerateRandomString(shortLinkLength)
	for _, found := storage.Links[shortLink]; found; {
		shortLink = random.GenerateRandomString(shortLinkLength)
	}
	storage.Links[shortLink] = fullLink

	return shortLink
}

func (storage *MemStorage) FindLink(shortLink string) (url.URL, bool) {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	fullLink, found := storage.Links[shortLink]
	return fullLink, found
}
