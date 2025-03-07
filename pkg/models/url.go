package models

type Url struct {
	UserId   int64  `json:"user_id"`
	ShortUrl string `json:"short_url"`
	LongUrl  string `json:"long_url"`
}

type AddUrl struct {
	UserId       int64  `json:"user_id"`
	ShortUrl     string `json:"short_url"`
	LongUrl      string `json:"long_url"`
	CustomAllias string `json:"custom_allias"`
}
