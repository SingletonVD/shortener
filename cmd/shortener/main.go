package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

type MemStorage struct {
	Links map[string]url.URL
	Lock  sync.Mutex
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const shortLinkLength = 8
const currentServerHost = "localhost:8080"

func generateShortLink() string {
	b := make([]byte, shortLinkLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (storage *MemStorage) createShortLink(fullLink url.URL) string {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	shortLink := generateShortLink()
	for _, found := storage.Links[shortLink]; found; {
		shortLink = generateShortLink()
	}
	storage.Links[shortLink] = fullLink

	return shortLink
}

func (storage *MemStorage) findLink(shortLink string) (url.URL, bool) {
	storage.Lock.Lock()
	defer storage.Lock.Unlock()

	fullLink, found := storage.Links[shortLink]
	return fullLink, found
}

func validateInputLink(inputLink string) (*url.URL, bool) {
	url, err := url.Parse(inputLink)

	if err != nil {
		return nil, false
	}

	linkValid := !(url.Host == "" || (url.Scheme != "http" && url.Scheme != "https"))

	return url, linkValid
}

func createShortLinkHandle(storage *MemStorage, w http.ResponseWriter, r *http.Request) {
	if (r.Header.Get("Content-Type")) != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	inputLink, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fullLink, valid := validateInputLink(string(inputLink))

	if !valid || fullLink == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLink := storage.createShortLink(*fullLink)

	w.Header().Add("Content-Type", "text/plain")
	fmt.Fprintf(w, "http://%s/%s", currentServerHost, shortLink)
	w.WriteHeader(http.StatusCreated)
}

func getShortLinkHandle(storage *MemStorage, w http.ResponseWriter, r *http.Request) {
	shortLink := r.PathValue("shortLink")
	fullLink, found := storage.findLink(shortLink)

	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Location", fullLink.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func run() error {
	storage := MemStorage{Links: make(map[string]url.URL)}
	mux := http.NewServeMux()

	mux.Handle("POST /{$}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { createShortLinkHandle(&storage, w, r) }))
	mux.Handle("GET /{shortLink}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { getShortLinkHandle(&storage, w, r) }))

	return http.ListenAndServe(":8080", mux)
}

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}
