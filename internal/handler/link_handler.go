package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/SingletonVD/shortener/internal/service"
	"github.com/SingletonVD/shortener/internal/validation"
	"github.com/go-chi/chi/v5"
)

type LinkHandler struct {
	linkService     *service.LinkService
	baseLinkAddress string
}

func NewLinkHandler(linkService *service.LinkService, baseLinkAddress string) *LinkHandler {
	return &LinkHandler{linkService: linkService, baseLinkAddress: baseLinkAddress}
}

func (handler *LinkHandler) CreateShortLinkHandle(w http.ResponseWriter, r *http.Request) {
	if (r.Header.Get("Content-Type")) != "text/plain" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	inputLink, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	valid := validation.ValidateRawLink(string(inputLink))

	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLink := handler.linkService.CreateShortLink(string(inputLink))

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", handler.baseLinkAddress, shortLink)
}

func (handler *LinkHandler) GetShortLinkHandle(w http.ResponseWriter, r *http.Request) {
	shortLink := chi.URLParam(r, "shortLink")
	fullLink, found := handler.linkService.FindFullLink(shortLink)

	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Location", fullLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
