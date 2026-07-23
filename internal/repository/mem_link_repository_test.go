package repository

import (
	"testing"

	"github.com/SingletonVD/shortener/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestSaveIfAvailable(t *testing.T) {
	testCases := []struct {
		name            string
		link            model.ShortenedLink
		repositoryState map[string]string
		want            bool
	}{
		{
			name:            "Save new link",
			link:            model.ShortenedLink{Short: "short", FullLink: "long"},
			repositoryState: make(map[string]string),
			want:            true,
		},
		{
			name:            "Save new link with short link collision",
			link:            model.ShortenedLink{Short: "short", FullLink: "new long"},
			repositoryState: map[string]string{"short": "old long"},
			want:            false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo := NewMemLinkRepository()
			repo.links = testCase.repositoryState
			assert.Equal(t, testCase.want, repo.SaveIfAvailable(testCase.link))
		})
	}
}

func TestFindFullLink(t *testing.T) {
	type Want struct {
		found    bool
		fullLink string
	}

	testCases := []struct {
		name            string
		shortLink       string
		repositoryState map[string]string
		want            Want
	}{
		{
			name:            "Find existing link",
			shortLink:       "found",
			repositoryState: map[string]string{"found": "long"},
			want: Want{
				found:    true,
				fullLink: "long",
			},
		},
		{
			name:            "Not found link",
			shortLink:       "not found",
			repositoryState: map[string]string{"found": "long"},
			want: Want{
				found: false,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo := NewMemLinkRepository()
			repo.links = testCase.repositoryState
			fullLink, found := repo.FindFullLink(testCase.shortLink)
			assert.Equal(t, testCase.want.found, found)
			assert.Equal(t, testCase.want.fullLink, fullLink)
		})
	}
}
