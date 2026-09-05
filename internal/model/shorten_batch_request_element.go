package model

type ShortenBatchRequestElement struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}
