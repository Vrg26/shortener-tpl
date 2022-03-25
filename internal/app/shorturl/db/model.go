package db

type ShortURL struct {
	ID        string `json:"id"`
	OriginURL string `json:"origin_url"`
	UserID    uint64 `json:"user_id"`
}
