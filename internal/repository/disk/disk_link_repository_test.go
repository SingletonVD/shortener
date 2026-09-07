package disk

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/SingletonVD/shortener/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiskSaveIfAvailable(t *testing.T) {
	testCases := []struct {
		name        string
		link        model.ShortenedLink
		fileContent string
		want        bool
	}{
		{
			name:        "Save new link",
			link:        model.ShortenedLink{Short: "short", FullLink: "long"},
			fileContent: `[]`,
			want:        true,
		},
		{
			name:        "Save new link with short link collision",
			link:        model.ShortenedLink{Short: "short", FullLink: "new long"},
			fileContent: `[{"uuid":"1","short_url":"short","original_url":"http://yandex.ru"}]`,
			want:        false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFileName := filepath.Join(tmpDir, "disk_test_save.json")
			os.WriteFile(testFileName, []byte(testCase.fileContent), 0666)
			repo, err := NewDiskLinkRepository(testFileName)
			require.NoError(t, err)

			saveResult, err := repo.SaveIfAvailable(context.TODO(), testCase.link)
			require.NoError(t, err)
			assert.Equal(t, testCase.want, saveResult)
		})
	}
}

func TestDiskSaveBatchIfAvailable(t *testing.T) {
	testCases := []struct {
		name        string
		links       []model.ShortenedLink
		fileContent string
		want        bool
	}{
		{
			name:        "Save new link",
			links:       []model.ShortenedLink{{Short: "short", FullLink: "long"}},
			fileContent: `[]`,
			want:        true,
		},
		{
			name:        "Save new link with short link collision",
			links:       []model.ShortenedLink{{Short: "short", FullLink: "new long"}, {Short: "short", FullLink: "long"}},
			fileContent: `[{"uuid":"1","short_url":"short","original_url":"http://yandex.ru"}]`,
			want:        false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFileName := filepath.Join(tmpDir, "disk_test_save.json")
			os.WriteFile(testFileName, []byte(testCase.fileContent), 0666)
			repo, err := NewDiskLinkRepository(testFileName)
			require.NoError(t, err)

			saveResult, err := repo.SaveBatchIfAvailable(context.TODO(), testCase.links)
			require.NoError(t, err)
			assert.Equal(t, testCase.want, saveResult)
		})
	}
}

func TestDiskFindFullLink(t *testing.T) {
	type Want struct {
		found    bool
		fullLink string
	}

	testCases := []struct {
		name        string
		shortLink   string
		fileContent string
		want        Want
	}{
		{
			name:        "Find existing link",
			shortLink:   "4rSPg8ap",
			fileContent: `[{"uuid":"1","short_url":"4rSPg8ap","original_url":"http://yandex.ru"}]`,
			want: Want{
				found:    true,
				fullLink: "http://yandex.ru",
			},
		},
		{
			name:        "Not found link",
			shortLink:   "not found",
			fileContent: `[{"uuid":"1","short_url":"4rSPg8ap","original_url":"http://yandex.ru"}]`,
			want: Want{
				found: false,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFileName := filepath.Join(tmpDir, "disk_test_find.json")
			os.WriteFile(testFileName, []byte(testCase.fileContent), 0666)
			repo, err := NewDiskLinkRepository(testFileName)
			require.NoError(t, err)

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
