package service

import (
	"context"

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

func (service *LinkService) FindLink(context context.Context, shortLink string) (*model.ShortenedLink, error) {
	return service.linkRepository.FindLink(context, shortLink)
}
