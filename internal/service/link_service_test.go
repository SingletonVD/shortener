package service

import (
	"testing"

	"github.com/SingletonVD/shortener/internal/model"
	"github.com/stretchr/testify/assert"
)

type FakeRepository struct {
	ignoreNextSave bool
	saveCalls      int
	savedLink      string
}

func (repo *FakeRepository) SaveIfAvailable(link model.ShortenedLink) bool {
	repo.saveCalls += 1
	if repo.ignoreNextSave {
		repo.ignoreNextSave = false
		return false
	}
	repo.savedLink = link.FullLink
	return true
}

func (repo *FakeRepository) FindFullLink(shortLink string) (string, bool) {
	return "", false
}

func TestCreateShortLink(t *testing.T) {
	type Want struct {
		regexp    string
		saveCalls int
		savedLink string
	}

	testCases := []struct {
		name string
		link string
		repo FakeRepository
		want Want
	}{
		{
			name: "Save link when no collision",
			link: "https://practicum.yandex.ru",
			repo: FakeRepository{ignoreNextSave: false},
			want: Want{
				regexp:    `^[a-zA-Z]{8}$`,
				saveCalls: 1,
				savedLink: "https://practicum.yandex.ru",
			},
		},
		{
			name: "Save link when collision",
			link: "https://practicum.yandex.ru",
			repo: FakeRepository{ignoreNextSave: true},
			want: Want{
				regexp:    `^[a-zA-Z]{8}$`,
				saveCalls: 2,
				savedLink: "https://practicum.yandex.ru",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := NewLinkService(&testCase.repo)
			shortLink := service.CreateShortLink(testCase.link)

			assert.Regexp(t, testCase.want.regexp, shortLink)
			assert.Equal(t, testCase.want.saveCalls, testCase.repo.saveCalls)
			assert.Equal(t, testCase.want.savedLink, testCase.repo.savedLink)
		})
	}
}
