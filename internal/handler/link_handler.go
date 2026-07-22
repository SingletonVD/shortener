package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/SingletonVD/shortener/internal/service"
)

type LinkHandler struct {
	linkService *service.LinkService
}

func NewLinkHandler(linkService *service.LinkService) *LinkHandler {
	return &LinkHandler{linkService: linkService}
}

func validateInputLink(inputLink string) (*url.URL, bool) {
	url, err := url.Parse(inputLink)

	if err != nil {
		return nil, false
	}

	linkValid := !(url.Host == "" || (url.Scheme != "http" && url.Scheme != "https"))

	return url, linkValid
}

func (handler *LinkHandler) CreateShortLinkHandle(currentServerHost string, w http.ResponseWriter, r *http.Request) {
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

	fullLink, valid := validateInputLink(string(inputLink))

	if !valid || fullLink == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLink := handler.linkService.CreateShortLink(*fullLink)

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "http://%s/%s", currentServerHost, shortLink)
}

func (handler *LinkHandler) GetShortLinkHandle(w http.ResponseWriter, r *http.Request) {
	shortLink := r.PathValue("shortLink")
	fullLink, found := handler.linkService.FindFullLink(shortLink)

	if !found {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Location", fullLink.String())
	w.WriteHeader(http.StatusTemporaryRedirect)
}
