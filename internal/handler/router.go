package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(linkHandler *LinkHandler, pingHandler *PingHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(GzipMiddleware)
	router.Use(LoggingMiddleware)

	router.Post("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		linkHandler.CreateShortLinkHandle(w, r)
	}))

	router.Post("/api/shorten", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		linkHandler.CreateShortLinkJSONHandle(w, r)
	}))

	router.Post("/api/shorten/batch", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		linkHandler.CreateShortLinksBatchJSONHandle(w, r)
	}))

	router.Get("/{shortLink}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { linkHandler.GetShortLinkHandle(w, r) }))

	router.Get("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { pingHandler.PingHandle(w, r) }))

	return router
}
