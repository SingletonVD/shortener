package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/SingletonVD/shortener/internal/model"
	"github.com/SingletonVD/shortener/internal/repository"
	"github.com/SingletonVD/shortener/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateShortLinkHandle(t *testing.T) {
	type Want struct {
		statusCode  int
		expectsBody bool
		contentType string
		bodyRegexp  string
	}

	baseLinkAddress := "http://localhost:8080"
	storage := repository.NewMemLinkRepository()
	service := service.NewLinkService(storage)
	handler := NewLinkHandler(service, baseLinkAddress)
	router := NewRouter(handler)

	testCases := []struct {
		name        string
		method      string
		path        string
		contentType string
		requestBody string
		want        Want
	}{
		{
			name:        "Positive case with getting short link",
			method:      http.MethodPost,
			path:        "/",
			contentType: "text/plain",
			requestBody: "https://practicum.yandex.ru",
			want: Want{
				statusCode:  http.StatusCreated,
				expectsBody: true,
				contentType: "text/plain",
				bodyRegexp:  `^http://localhost:8080/[a-zA-Z]{8}$`,
			},
		},
		{
			name:        "Bad request with wrong link",
			method:      http.MethodPost,
			path:        "/",
			contentType: "text/plain",
			requestBody: "/just/a/path",
			want: Want{
				statusCode:  http.StatusBadRequest,
				expectsBody: false,
			},
		},
		{
			name:        "Bad request with wrong content-type",
			method:      http.MethodPost,
			path:        "/",
			contentType: "application/json",
			requestBody: "https://practicum.yandex.ru",
			want: Want{
				statusCode:  http.StatusBadRequest,
				expectsBody: false,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(testCase.method, testCase.path, strings.NewReader(testCase.requestBody))
			request.Header.Set("Content-Type", testCase.contentType)
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)
			response := responseRecorder.Result()

			assert.Equal(t, testCase.want.statusCode, response.StatusCode)

			if testCase.want.expectsBody {
				assert.Equal(t, testCase.want.contentType, response.Header.Get("Content-Type"))

				defer response.Body.Close()
				body, err := io.ReadAll(response.Body)

				require.NoError(t, err)
				assert.Regexp(t, testCase.want.bodyRegexp, string(body))
			}
		})
	}
}

func TestCreateShortLinkJsonHandle(t *testing.T) {
	type Want struct {
		statusCode  int
		expectsBody bool
		contentType string
		bodyRegexp  string
	}

	baseLinkAddress := "http://localhost:8080"
	storage := repository.NewMemLinkRepository()
	service := service.NewLinkService(storage)
	handler := NewLinkHandler(service, baseLinkAddress)
	router := NewRouter(handler)

	testCases := []struct {
		name        string
		method      string
		path        string
		contentType string
		requestBody string
		want        Want
	}{
		{
			name:        "Positive case with getting short link",
			method:      http.MethodPost,
			path:        "/api/shorten",
			contentType: "application/json",
			requestBody: `{"url":"https://practicum.yandex.ru"}`,
			want: Want{
				statusCode:  http.StatusCreated,
				expectsBody: true,
				contentType: "application/json",
				bodyRegexp:  `^\{"result":"http://localhost:8080/[a-zA-Z]{8}"\}\n$`,
			},
		},
		{
			name:        "Bad request with wrong link",
			method:      http.MethodPost,
			path:        "/api/shorten",
			contentType: "application/json",
			requestBody: `{"url":"/just/a/path"}`,
			want: Want{
				statusCode:  http.StatusBadRequest,
				expectsBody: false,
			},
		},
		{
			name:        "Bad request with wrong content-type",
			method:      http.MethodPost,
			path:        "/api/shorten",
			contentType: "text/plain",
			requestBody: "https://practicum.yandex.ru",
			want: Want{
				statusCode:  http.StatusBadRequest,
				expectsBody: false,
			},
		},
		{
			name:        "Bad request with wrong body",
			method:      http.MethodPost,
			path:        "/api/shorten",
			contentType: "application/json",
			requestBody: "https://practicum.yandex.ru",
			want: Want{
				statusCode:  http.StatusBadRequest,
				expectsBody: false,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(testCase.method, testCase.path, strings.NewReader(testCase.requestBody))
			request.Header.Set("Content-Type", testCase.contentType)
			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, request)
			response := responseRecorder.Result()

			assert.Equal(t, testCase.want.statusCode, response.StatusCode)

			if testCase.want.expectsBody {
				assert.Equal(t, testCase.want.contentType, response.Header.Get("Content-Type"))

				defer response.Body.Close()
				body, err := io.ReadAll(response.Body)

				require.NoError(t, err)
				assert.Regexp(t, testCase.want.bodyRegexp, string(body))
			}
		})
	}
}

type FakeRepository struct {
	links map[string]string
}

func (repo *FakeRepository) SaveIfAvailable(link model.ShortenedLink) bool {
	return true
}

func (repo *FakeRepository) FindFullLink(shortLink string) (string, bool) {
	fullLink, found := repo.links[shortLink]
	return fullLink, found
}

func TestGetShortLinkHandle(t *testing.T) {

	type Want struct {
		statusCode       int
		expectRedirect   bool
		expectedRedirect string
	}

	testCases := []struct {
		name           string
		method         string
		shortLink      string
		path           string
		fakeRepository repository.LinkRepository
		want           Want
	}{
		{
			name:           "Positive with found link",
			method:         http.MethodGet,
			shortLink:      "yandexpr",
			path:           "/yandexpr",
			fakeRepository: &FakeRepository{map[string]string{"yandexpr": "https://practicum.yandex.ru"}},
			want: Want{
				statusCode:       http.StatusTemporaryRedirect,
				expectRedirect:   true,
				expectedRedirect: "https://practicum.yandex.ru",
			},
		},
		{
			name:           "Bad request with non existent link",
			method:         http.MethodGet,
			shortLink:      "yandexpr",
			path:           "/yandexpr",
			fakeRepository: &FakeRepository{make(map[string]string)},
			want: Want{
				statusCode:     http.StatusBadRequest,
				expectRedirect: false,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			baseLinkAddress := "http://localhost:8080"
			storage := testCase.fakeRepository
			service := service.NewLinkService(storage)
			handler := NewLinkHandler(service, baseLinkAddress)
			router := NewRouter(handler)

			request := httptest.NewRequest(testCase.method, testCase.path, nil)
			request.SetPathValue("shortLink", testCase.shortLink)
			responseRecorder := httptest.NewRecorder()

			router.ServeHTTP(responseRecorder, request)
			response := responseRecorder.Result()

			assert.Equal(t, testCase.want.statusCode, response.StatusCode)

			if testCase.want.expectRedirect {
				redirect := response.Header.Get("Location")
				assert.Equal(t, testCase.want.expectedRedirect, redirect)
			}
		})
	}
}
