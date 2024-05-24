package model

import (
	"database/sql"
)

type User struct {
	ID        int
	Username  string
	Password  string
	IsAdmin   bool
	DeletedAt sql.NullTime
}
