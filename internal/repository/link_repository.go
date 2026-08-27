package repository

import (
	"github.com/SingletonVD/shortener/internal/model"
)

type LinkRepository interface {
	SaveIfAvailable(link model.ShortenedLink) (bool, error)
	// здесь error пока не нужен, но сразу сделал задел на подключение БД, хоть и нарушаю YAGNI
	FindLink(shortLink string) (*model.ShortenedLink, error)
}
