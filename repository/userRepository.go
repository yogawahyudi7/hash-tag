package repository

import (
	"database/sql"

	model "github.com/yogawahyudi7/social-media/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Register(username, password string) (int, error) {
	userId := 0
	err := r.db.QueryRow(
		`INSERT INTO 
			"user" (username, password)
		VALUES 
			($1, $2)
		RETURNING
			id`,
		username, password,
	).Scan(&userId)

	if err != nil {
		return 0, err
	}

	if userId < 1 {
		return 0, nil
	}

	return userId, nil
}

func (r *UserRepository) FindUsername(username string) (*model.User, error) {

	var user model.User
	query := `
	SELECT 
		u.id, u.username, u.password, u.is_admin 
	FROM 
		"user" u
	WHERE
		u.deleted_at IS NULL
	AND
		username = $1
	`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.IsAdmin,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
