package models

import "github.com/gocql/gocql"

type Url struct {
	Id       gocql.UUID `json:"id"`
	UserId   int        `json:"user_id"`
	ShortUrl string     `json:"short_url"`
	LongUrl  string     `json:"long_url"`
}
