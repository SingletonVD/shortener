package service

import (
	"net/url"

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

func (service *LinkService) CreateShortLink(fullLink url.URL) string {
	shortLink := random.GenerateRandomString(shortLinkLength)

	shortenedLink := model.ShortenedLink{
		Short:    shortLink,
		FullLink: fullLink,
	}

	for saved := service.linkRepository.SaveIfAvailable(shortenedLink); !saved; {
		shortLink = random.GenerateRandomString(shortLinkLength)
		shortenedLink.Short = shortLink
	}

	return shortLink
}

func (service *LinkService) FindFullLink(shortLink string) (url.URL, bool) {
	return service.linkRepository.FindFullLink(shortLink)
}
