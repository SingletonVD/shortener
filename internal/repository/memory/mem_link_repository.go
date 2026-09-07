package memory

import (
	"context"
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

func (storage *MemLinkRepository) SaveIfAvailable(_ context.Context, link model.ShortenedLink) (bool, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	_, found := storage.links[link.Short]
	if found {
		return false, nil
	}

	storage.links[link.Short] = link.FullLink

	return true, nil
}

func (storage *MemLinkRepository) SaveBatchIfAvailable(_ context.Context, links []model.ShortenedLink) (bool, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	for _, link := range links {
		_, found := storage.links[link.Short]
		if found {
			return false, nil
		}
	}

	for _, link := range links {
		storage.links[link.Short] = link.FullLink
	}

	return true, nil
}

func (storage *MemLinkRepository) FindLink(_ context.Context, shortLink string) (*model.ShortenedLink, error) {
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

func (storage *MemLinkRepository) Ping(ctx context.Context) error {
	return nil
}
