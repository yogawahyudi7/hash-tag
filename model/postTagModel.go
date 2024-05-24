package model

import "database/sql"

type PostTag struct {
	PostID    int
	TagID     int
	DeletedAt sql.NullTime
}
