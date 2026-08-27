package repository

import (
	"sync"

	"github.com/SingletonVD/shortener/internal/model"
)

// теперь unused, но оставлю
type MemLinkRepository struct {
	links map[string]string
	lock  sync.Mutex
}

func NewMemLinkRepository() *MemLinkRepository {
	return &MemLinkRepository{links: make(map[string]string)}
}

func (storage *MemLinkRepository) SaveIfAvailable(link model.ShortenedLink) (bool, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	_, found := storage.links[link.Short]
	if found {
		return false, nil
	}

	storage.links[link.Short] = link.FullLink

	return true, nil
}

func (storage *MemLinkRepository) FindLink(shortLink string) (*model.ShortenedLink, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	fullLink, found := storage.links[shortLink]

	if !found {
		return nil, nil
	}

	return &model.ShortenedLink{
		Short:    shortLink,
		FullLink: fullLink,
	}, nil
}
