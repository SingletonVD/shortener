package repository

import (
	"net/url"

	"github.com/SingletonVD/shortener/internal/model"
)

type LinksRepository interface {
	SaveIfAvailable(link model.ShortenedLink) bool
	FindFullLink(shortLink string) (url.URL, bool)
}
