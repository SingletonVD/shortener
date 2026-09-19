package memory

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
		userID          string
		repositoryState map[string]model.UserShortenedLink
		want            bool
	}{
		{
			name:            "Save new link",
			link:            model.ShortenedLink{Short: "short", FullLink: "long"},
			userID:          "1",
			repositoryState: make(map[string]model.UserShortenedLink),
			want:            true,
		},
		{
			name:   "Save new link with short link collision",
			link:   model.ShortenedLink{Short: "short", FullLink: "new long"},
			userID: "1",
			repositoryState: map[string]model.UserShortenedLink{
				"short": model.UserShortenedLink{
					ShortenedLink: model.ShortenedLink{
						Short:    "short",
						FullLink: "old long",
					},
					UserID: "1",
				},
			},
			want: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo := NewMemLinkRepository()
			repo.links = testCase.repositoryState
			saveResult, err := repo.SaveIfAvailable(context.TODO(), testCase.link, testCase.userID)
			require.NoError(t, err)
			assert.Equal(t, testCase.want, saveResult)
		})
	}
}

func TestSaveBatchIfAvailable(t *testing.T) {
	testCases := []struct {
		name            string
		links           []model.ShortenedLink
		userID          string
		repositoryState map[string]model.UserShortenedLink
		want            bool
	}{
		{
			name:            "Save new link",
			links:           []model.ShortenedLink{{Short: "short", FullLink: "long"}},
			userID:          "1",
			repositoryState: make(map[string]model.UserShortenedLink),
			want:            true,
		},
		{
			name:   "Save new link with short link collision",
			links:  []model.ShortenedLink{{Short: "short", FullLink: "new long"}, {Short: "short", FullLink: "long"}},
			userID: "1",
			repositoryState: map[string]model.UserShortenedLink{
				"short": model.UserShortenedLink{
					ShortenedLink: model.ShortenedLink{
						Short:    "short",
						FullLink: "old long",
					},
					UserID: "1",
				},
			},
			want: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo := NewMemLinkRepository()
			repo.links = testCase.repositoryState
			saveResult, err := repo.SaveBatchIfAvailable(context.TODO(), testCase.links, testCase.userID)
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
		repositoryState map[string]model.UserShortenedLink
		want            Want
	}{
		{
			name:      "Find existing link",
			shortLink: "found",
			repositoryState: map[string]model.UserShortenedLink{
				"found": {
					ShortenedLink: model.ShortenedLink{
						Short:    "found",
						FullLink: "long",
					},
					UserID: "1",
				},
			},
			want: Want{
				found:    true,
				fullLink: "long",
			},
		},
		{
			name:      "Not found link",
			shortLink: "not found",
			repositoryState: map[string]model.UserShortenedLink{
				"found": {
					ShortenedLink: model.ShortenedLink{
						Short:    "found",
						FullLink: "long",
					},
					UserID: "1",
				},
			},
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
