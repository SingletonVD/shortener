package handler

import (
	"github.com/SingletonVD/shortener/internal/handler/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(linkHandler *LinkHandler, pingHandler *PingHandler, authMiddleware *middleware.AuthMiddleware) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.GzipMiddleware)
	router.Use(middleware.LoggingMiddleware)

	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.IntrospectUserJWT)
		r.Post("/", linkHandler.CreateShortLinkHandle)
		r.Post("/api/shorten", linkHandler.CreateShortLinkJSONHandle)
		r.Post("/api/shorten/batch", linkHandler.CreateShortLinksBatchJSONHandle)
		r.Get("/api/user/urls", linkHandler.GetUserLinksHandle)
		r.Delete("/api/user/urls", linkHandler.DeleteLinks)
	})

	router.Get("/{shortLink}", linkHandler.GetShortLinkHandle)
	router.Get("/ping", pingHandler.PingHandle)

	return router
}
