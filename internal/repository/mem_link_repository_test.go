package repository

import (
	"context"
	"testing"

	"github.com/SingletonVD/shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			saveResult, err := repo.SaveIfAvailable(context.TODO(), testCase.link)
			require.NoError(t, err)
			assert.Equal(t, testCase.want, saveResult)
		})
	}
}

func TestSaveBatchIfAvailable(t *testing.T) {
	testCases := []struct {
		name            string
		links           []model.ShortenedLink
		repositoryState map[string]string
		want            bool
	}{
		{
			name:            "Save new link",
			links:           []model.ShortenedLink{{Short: "short", FullLink: "long"}},
			repositoryState: make(map[string]string),
			want:            true,
		},
		{
			name:            "Save new link with short link collision",
			links:           []model.ShortenedLink{{Short: "short", FullLink: "new long"}, {Short: "short", FullLink: "long"}},
			repositoryState: map[string]string{"short": "old long"},
			want:            false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo := NewMemLinkRepository()
			repo.links = testCase.repositoryState
			saveResult, err := repo.SaveBatchIfAvailable(context.TODO(), testCase.links)
			require.NoError(t, err)
			assert.Equal(t, testCase.want, saveResult)
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
			link, err := repo.FindLink(context.TODO(), testCase.shortLink)
			require.NoError(t, err)

			if testCase.want.found {
				require.NotNil(t, link)
				assert.Equal(t, testCase.want.fullLink, link.FullLink)
			} else {
				assert.Nil(t, link)
			}
		})
	}
}
