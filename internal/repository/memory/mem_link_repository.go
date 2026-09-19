package memory

import (
	"context"
	"sync"

	"github.com/SingletonVD/shortener/internal/model"
)

type MemLinkRepository struct {
	links map[string]model.UserShortenedLink
	lock  sync.Mutex
}

func NewMemLinkRepository() *MemLinkRepository {
	return &MemLinkRepository{links: make(map[string]model.UserShortenedLink)}
}

func (storage *MemLinkRepository) SaveIfAvailable(_ context.Context, link model.ShortenedLink, userID string) (bool, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	_, found := storage.links[link.Short]
	if found {
		return false, nil
	}

	storage.links[link.Short] = model.UserShortenedLink{ShortenedLink: link, UserID: userID}

	return true, nil
}

func (storage *MemLinkRepository) SaveBatchIfAvailable(_ context.Context, links []model.ShortenedLink, userID string) (bool, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	for _, link := range links {
		_, found := storage.links[link.Short]
		if found {
			return false, nil
		}
	}

	for _, link := range links {
		storage.links[link.Short] = model.UserShortenedLink{ShortenedLink: link, UserID: userID}
	}

	return true, nil
}

func (storage *MemLinkRepository) FindLink(_ context.Context, shortLink string) (*model.ShortenedLink, error) {
	storage.lock.Lock()
	defer storage.lock.Unlock()

	userShortenedLink, found := storage.links[shortLink]

	if !found {
		return nil, nil
	}

	return &userShortenedLink.ShortenedLink, nil
}

func (storage *MemLinkRepository) Ping(ctx context.Context) error {
	return nil
}

func (storage *MemLinkRepository) GetUserLinks(ctx context.Context, userID string) ([]model.ShortenedLink, error) {
	// неоптимально, но не хотелось строить что-то типа индекса по пользователям для легаси реализации
	result := make([]model.ShortenedLink, 0, 0)

	for _, link := range storage.links {
		if link.UserID == userID {
			result = append(result, link.ShortenedLink)
		}
	}

	return result, nil
}
