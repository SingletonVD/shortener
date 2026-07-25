package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const currentServerHost = "localhost:8080"

func NewRouter(handler *LinkHandler) *chi.Mux {
	router := chi.NewRouter()
	// router.Use(middleware.AllowContentType("text/plain")) в задании указано возвращать 400 на ошибочные запросы, но тут вернется 415
	router.Post("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.CreateShortLinkHandle(currentServerHost, w, r)
	}))
	router.Get("/{shortLink}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handler.GetShortLinkHandle(w, r) }))
	return router
}
