package service

import (
	"context"
	"maps"
	"slices"

	"github.com/SingletonVD/shortener/internal/model"
	"github.com/SingletonVD/shortener/internal/random"
	"github.com/SingletonVD/shortener/internal/repository"
)

type LinkService struct {
	linkRepository repository.LinkRepository
}

func NewLinkService(linkRepository repository.LinkRepository) *LinkService {
	return &LinkService{linkRepository: linkRepository}
}

const (
	shortLinkLength = 8
)

func (service *LinkService) CreateShortLink(context context.Context, fullLink string) (string, error) {
	shortLink := random.GenerateRandomString(shortLinkLength)

	shortenedLink := model.ShortenedLink{
		Short:    shortLink,
		FullLink: fullLink,
	}

	for {
		saved, err := service.linkRepository.SaveIfAvailable(context, shortenedLink)

		if err != nil {
			return "", err
		}

		if saved {
			break
		}

		shortLink = random.GenerateRandomString(shortLinkLength)
		shortenedLink.Short = shortLink
	}

	return shortLink, nil
}

func (service *LinkService) CreateShortLinksBatch(context context.Context, fullLinksBatch map[string]string) (map[string]model.ShortenedLink, error) {
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

	for {

		saved, err := service.linkRepository.SaveBatchIfAvailable(context, shortenedLinks)

		if err != nil {
			return nil, err
		}

		if saved {
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

	return shortenedLinksMap, nil
}

func (service *LinkService) FindLink(context context.Context, shortLink string) (*model.ShortenedLink, error) {
	return service.linkRepository.FindLink(context, shortLink)
}
