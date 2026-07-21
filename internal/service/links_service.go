package service

import (
	"net/url"

	"github.com/SingletonVD/shortener/internal/model"
	"github.com/SingletonVD/shortener/internal/random"
	"github.com/SingletonVD/shortener/internal/repository"
)

type LinksService struct {
	LinksRepository repository.LinksRepository
}

const (
	shortLinkLength = 8
)

func (service *LinksService) CreateShortLink(fullLink url.URL) string {
	shortLink := random.GenerateRandomString(shortLinkLength)

	shortenedLink := model.ShortenedLink{
		Short:    shortLink,
		FullLink: fullLink,
	}

	for saved := service.LinksRepository.SaveIfAvailable(shortenedLink); !saved; {
		shortLink = random.GenerateRandomString(shortLinkLength)
		shortenedLink.Short = shortLink
	}

	return shortLink
}

func (service *LinksService) FindFullLink(shortLink string) (url.URL, bool) {
	return service.LinksRepository.FindFullLink(shortLink)
}
