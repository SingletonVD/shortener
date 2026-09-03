package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(linkHandler *LinkHandler, pingHandler *PingHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(GzipMiddleware)
	router.Use(LoggingMiddleware)
	// router.Use(middleware.AllowContentType("text/plain")) в задании указано возвращать 400 на ошибочные запросы, но тут вернется 415

	router.Post("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		linkHandler.CreateShortLinkHandle(w, r)
	}))

	router.Post("/api/shorten", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		linkHandler.CreateShortLinkJsonHandle(w, r)
	}))

	router.Get("/{shortLink}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { linkHandler.GetShortLinkHandle(w, r) }))

	router.Get("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { pingHandler.PingHandle(w, r) }))

	return router
}
