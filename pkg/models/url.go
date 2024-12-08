package models

import "github.com/gocql/gocql"

type Url struct {
	Id       gocql.UUID
	UserId   int
	ShortUrl string
	LongUrl  string
}
