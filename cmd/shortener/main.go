package main

import (
	"net/http"
	"net/url"

	"github.com/SingletonVD/shortener/internal/handler"
	"github.com/SingletonVD/shortener/internal/repository"
)

const currentServerHost = "localhost:8080"

func run() error {
	storage := repository.MemStorage{Links: make(map[string]url.URL)}
	mux := http.NewServeMux()

	mux.Handle("POST /{$}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.CreateShortLinkHandle(&storage, currentServerHost, w, r)
	}))
	mux.Handle("GET /{shortLink}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handler.GetShortLinkHandle(&storage, w, r) }))

	return http.ListenAndServe(":8080", mux)
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
