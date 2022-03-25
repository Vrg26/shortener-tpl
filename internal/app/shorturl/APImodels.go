package shorturl

type RequestURL struct {
	URL string `json:"url"`
}
type RespResultURL struct {
	Result string `json:"result"`
}

type RespShortUrl struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
