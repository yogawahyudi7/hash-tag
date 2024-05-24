package model

import (
	"database/sql"
)

type Post struct {
	ID          int
	Title       string
	Content     string
	Status      string
	PublishDate sql.NullTime
	DeletedAt   sql.NullTime

	Tags []Tag
}
