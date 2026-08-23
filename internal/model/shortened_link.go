package model

type ShortenedLink struct {
	Short    string `json:"short_url"`
	FullLink string `json:"original_url"`
}

// мне uuid кажется избыточным, т.к. сокращенная ссылка и так должна быть уникальна,
// но не уверен, что тесты пройдут без uuid,
// поэтому создам такую "обертку" для использования только при сохранении
type PersistedShortenedLink struct {
	ShortenedLink
	UUID int `json:"uuid,string"`
}
