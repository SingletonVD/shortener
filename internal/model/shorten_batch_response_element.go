package model

type ShortenBatchResponseElement struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
