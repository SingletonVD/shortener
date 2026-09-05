package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/SingletonVD/shortener/internal/model"
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

	shortLink, err := handler.linkService.CreateShortLink(r.Context(), string(inputLink))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s/%s", handler.baseLinkAddress, shortLink)
}

func (handler *LinkHandler) CreateShortLinkJSONHandle(w http.ResponseWriter, r *http.Request) {
	if (r.Header.Get("Content-Type")) != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var shortenRequest model.ShortenRequest

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&shortenRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	valid := validation.ValidateRawLink(string(shortenRequest.URL))

	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLink, err := handler.linkService.CreateShortLink(r.Context(), string(shortenRequest.URL))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := fmt.Sprintf("%s/%s", handler.baseLinkAddress, shortLink)
	response := model.ShortenResponse{
		Result: result,
	}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

func (handler *LinkHandler) CreateShortLinksBatchJSONHandle(w http.ResponseWriter, r *http.Request) {
	if (r.Header.Get("Content-Type")) != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var shortenBatchRequest []model.ShortenBatchRequestElement

	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&shortenBatchRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, requestElement := range shortenBatchRequest {
		valid := validation.ValidateRawLink(string(requestElement.OriginalURL))

		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	originalLinksMap := make(map[string]string)
	for _, requestElement := range shortenBatchRequest {
		originalLinksMap[requestElement.CorrelationId] = requestElement.OriginalURL
	}

	shortenedLinksMap, err := handler.linkService.CreateShortLinksBatch(r.Context(), originalLinksMap)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := make([]model.ShortenBatchResponseElement, 0)

	for correlationId, shortenedLink := range shortenedLinksMap {
		shortURL := fmt.Sprintf("%s/%s", handler.baseLinkAddress, shortenedLink.Short)
		responseElement := model.ShortenBatchResponseElement{
			CorrelationId: correlationId,
			ShortURL:      shortURL,
		}
		response = append(response, responseElement)
	}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseJSON)
}

func (handler *LinkHandler) GetShortLinkHandle(w http.ResponseWriter, r *http.Request) {
	shortLink := chi.URLParam(r, "shortLink")
	link, err := handler.linkService.FindLink(r.Context(), shortLink)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if link == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Location", link.FullLink)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
