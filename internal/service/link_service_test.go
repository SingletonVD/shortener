package service

import (
	"context"
	"testing"

	"github.com/SingletonVD/shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FakeRepository struct {
	ignoreNextSave  bool
	saveCalls       int
	savedLink       string
	savedLinksBatch []string
}

func (repo *FakeRepository) SaveIfAvailable(_ context.Context, link model.ShortenedLink) (bool, error) {
	repo.saveCalls += 1
	if repo.ignoreNextSave {
		repo.ignoreNextSave = false
		return false, nil
	}
	repo.savedLink = link.FullLink
	return true, nil
}

func (repo *FakeRepository) SaveBatchIfAvailable(_ context.Context, linksBatch []model.ShortenedLink) (bool, error) {
	repo.saveCalls += 1
	if repo.ignoreNextSave {
		repo.ignoreNextSave = false
		return false, nil
	}

	for _, link := range linksBatch {
		repo.savedLinksBatch = append(repo.savedLinksBatch, link.FullLink)
	}
	return true, nil
}

func (repo *FakeRepository) FindLink(_ context.Context, _ string) (*model.ShortenedLink, error) {
	return nil, nil
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
			shortLink, err := service.CreateShortLink(context.TODO(), testCase.link)

			require.NoError(t, err)
			assert.Regexp(t, testCase.want.regexp, shortLink)
			assert.Equal(t, testCase.want.saveCalls, testCase.repo.saveCalls)
			assert.Equal(t, testCase.want.savedLink, testCase.repo.savedLink)
		})
	}
}

func TestCreateShortLinksBatch(t *testing.T) {
	type Want struct {
		regexp     string
		saveCalls  int
		savedLinks []string
	}

	testCases := []struct {
		name  string
		links map[string]string
		repo  FakeRepository
		want  Want
	}{
		{
			name:  "Save links when no collision",
			links: map[string]string{"1": "https://practicum.yandex.ru", "2": "https://ya.ru"},
			repo:  FakeRepository{ignoreNextSave: false},
			want: Want{
				regexp:     `^[a-zA-Z]{8}$`,
				saveCalls:  1,
				savedLinks: []string{"https://practicum.yandex.ru", "https://ya.ru"},
			},
		},
		{
			name:  "Save links when collision",
			links: map[string]string{"1": "https://practicum.yandex.ru", "2": "https://ya.ru"},
			repo:  FakeRepository{ignoreNextSave: true},
			want: Want{
				regexp:     `^[a-zA-Z]{8}$`,
				saveCalls:  2,
				savedLinks: []string{"https://practicum.yandex.ru", "https://ya.ru"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := NewLinkService(&testCase.repo)
			savedLinks, err := service.CreateShortLinksBatch(context.TODO(), testCase.links)

			require.NoError(t, err)

			savedFullLinks := make([]string, 0)
			for _, savedLink := range savedLinks {
				assert.Regexp(t, testCase.want.regexp, savedLink.Short)
				savedFullLinks = append(savedFullLinks, savedLink.FullLink)
			}
			assert.Equal(t, testCase.want.saveCalls, testCase.repo.saveCalls)
			assert.ElementsMatch(t, testCase.want.savedLinks, savedFullLinks)
		})
	}
}
