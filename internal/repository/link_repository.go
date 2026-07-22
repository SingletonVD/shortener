package repository

import (
	"net/url"

	"github.com/SingletonVD/shortener/internal/model"
)

type LinkRepository interface {
	SaveIfAvailable(link model.ShortenedLink) bool
	FindFullLink(shortLink string) (url.URL, bool)
}
