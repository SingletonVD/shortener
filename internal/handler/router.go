package handler

import (
	"net/http"

	"github.com/SingletonVD/shortener/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(linkHandler *LinkHandler, pingHandler *PingHandler, authMiddleware *middleware.AuthMiddleware) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.GzipMiddleware)
	router.Use(middleware.LoggingMiddleware)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.IntrospectUserJWT)
		r.Post("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			linkHandler.CreateShortLinkHandle(w, r)
		}))

		r.Post("/api/shorten", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			linkHandler.CreateShortLinkJSONHandle(w, r)
		}))

		r.Post("/api/shorten/batch", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			linkHandler.CreateShortLinksBatchJSONHandle(w, r)
		}))

		r.Get("/api/user/urls", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			linkHandler.GetUserLinksHandle(w, r)
		}))
	})

	router.Get("/{shortLink}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { linkHandler.GetShortLinkHandle(w, r) }))
	router.Get("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { pingHandler.PingHandle(w, r) }))

	return router
}
