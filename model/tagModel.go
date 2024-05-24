package model

import "database/sql"

type Tag struct {
	ID        int
	Label     string
	DeletedAt sql.NullTime
}
