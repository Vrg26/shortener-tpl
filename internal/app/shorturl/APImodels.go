package shorturl

type RequestURL struct {
	URL string `json:"url"`
}
type ResponseURL struct {
	Result string `json:"result"`
}
