package db

type ShortURL struct {
	ID            string `json:"id"`
	OriginURL     string `json:"origin_url"`
	UserID        uint32 `json:"user_id"`
	CorrelationID string
}
