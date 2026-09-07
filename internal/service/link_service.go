package service

import (
	"context"
	"errors"
	"maps"
	"slices"

	"github.com/SingletonVD/shortener/internal/model"
	"github.com/SingletonVD/shortener/internal/random"
)

type LinkRepository interface {
	SaveIfAvailable(ctx context.Context, link model.ShortenedLink) (bool, error)
	SaveBatchIfAvailable(ctx context.Context, links []model.ShortenedLink) (bool, error)
	FindLink(ctx context.Context, shortLink string) (*model.ShortenedLink, error)
}

type LinkService struct {
	linkRepository LinkRepository
}

func NewLinkService(linkRepository LinkRepository) *LinkService {
	return &LinkService{linkRepository: linkRepository}
}

const (
	shortLinkLength = 8
	retryLimit      = 3
)

var ErrUniqueShortLinkViolation = errors.New("generated short link violates unique constraint")

func (service *LinkService) CreateShortLink(ctx context.Context, fullLink string) (string, error) {
	shortLink := random.GenerateRandomString(shortLinkLength)

	shortenedLink := model.ShortenedLink{
		Short:    shortLink,
		FullLink: fullLink,
	}

	returnError := ErrUniqueShortLinkViolation

	for range retryLimit {
		saved, err := service.linkRepository.SaveIfAvailable(ctx, shortenedLink)

		if err != nil {
			return "", err
		}

		if saved {
			returnError = nil
			break
		}

		shortLink = random.GenerateRandomString(shortLinkLength)
		shortenedLink.Short = shortLink
	}

	return shortLink, returnError
}

func (service *LinkService) CreateShortLinksBatch(ctx context.Context, fullLinksBatch map[string]string) (map[string]model.ShortenedLink, error) {
	if len(fullLinksBatch) == 0 {
		return make(map[string]model.ShortenedLink), nil
	}

	shortenedLinksMap := make(map[string]model.ShortenedLink)

	for correlationID, fullLink := range fullLinksBatch {
		shortLink := random.GenerateRandomString(shortLinkLength)

		shortenedLinksMap[correlationID] = model.ShortenedLink{
			Short:    shortLink,
			FullLink: fullLink,
		}
	}

	shortenedLinks := slices.Collect(maps.Values(shortenedLinksMap))
	returnError := ErrUniqueShortLinkViolation

	for range retryLimit {
		saved, err := service.linkRepository.SaveBatchIfAvailable(ctx, shortenedLinks)

		if err != nil {
			return nil, err
		}

		if saved {
			returnError = nil
			break
		}

		// смелое предположение, что энтропии хватит, чтобы следующая генерация в обозримом будущем обошлась без коллизий
		for correlationID, shortenedLink := range shortenedLinksMap {
			shortLink := random.GenerateRandomString(shortLinkLength)

			shortenedLinksMap[correlationID] = model.ShortenedLink{
				Short:    shortLink,
				FullLink: shortenedLink.FullLink,
			}
		}
	}

	return shortenedLinksMap, returnError
}

func (service *LinkService) FindLink(ctx context.Context, shortLink string) (*model.ShortenedLink, error) {
	return service.linkRepository.FindLink(ctx, shortLink)
}
