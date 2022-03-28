package shorturl

type RequestURL struct {
	URL string `json:"url"`
}
type RespResultURL struct {
	Result string `json:"result"`
}

type RespShortURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type RequestBatchURL struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBatchURL struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
