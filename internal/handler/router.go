package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *LinkHandler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(GzipMiddleware)
	router.Use(LoggingMiddleware)
	// router.Use(middleware.AllowContentType("text/plain")) в задании указано возвращать 400 на ошибочные запросы, но тут вернется 415

	router.Post("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.CreateShortLinkHandle(w, r)
	}))

	router.Post("/api/shorten", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.CreateShortLinkJsonHandle(w, r)
	}))

	router.Get("/{shortLink}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handler.GetShortLinkHandle(w, r) }))

	return router
}
