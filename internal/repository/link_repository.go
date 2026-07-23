package repository

import (
	"github.com/SingletonVD/shortener/internal/model"
)

type LinkRepository interface {
	SaveIfAvailable(link model.ShortenedLink) bool
	FindFullLink(shortLink string) (string, bool)
}
