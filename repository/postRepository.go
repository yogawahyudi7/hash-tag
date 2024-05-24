package repository

import (
	"database/sql"
	"fmt"

	"github.com/yogawahyudi7/social-media/helper"
	model "github.com/yogawahyudi7/social-media/model"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db}
}

func (r *PostRepository) FindAll() ([]model.Post, error) {

	rows, err := r.db.Query(
		`SELECT 
			p.* 
		FROM 
			post p 
		WHERE 
			deleted_at IS NULL`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Status, &post.PublishDate, &post.DeletedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostRepository) FindAllWithTag() ([]model.Post, error) {

	query := `
	SELECT 
		p.*, t.id, t.label 
	FROM 
		post p 
	LEFT JOIN 
		post_tag pt ON p.id = pt.post_id 
	LEFT JOIN 
		tag t ON pt.tag_id = t.id 
	WHERE 
		p.deleted_at IS NULL
	`

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	postsMap := make(map[int]*model.Post)
	for rows.Next() {
		var post model.Post

		var tagID sql.NullInt64
		var tagLabel sql.NullString
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Status, &post.PublishDate, &post.DeletedAt, &tagID, &tagLabel)
		if err != nil {
			return nil, err
		}

		if existingPost, ok := postsMap[post.ID]; ok {
			if tagID.Valid {
				existingPost.Tags = append(existingPost.Tags, model.Tag{ID: int(tagID.Int64), Label: tagLabel.String})
			}
		} else {
			if tagID.Valid {
				post.Tags = []model.Tag{{ID: int(tagID.Int64), Label: tagLabel.String}}
			}
			postsMap[post.ID] = &post
		}

	}

	var posts []model.Post
	for _, post := range postsMap {
		posts = append(posts, *post)
	}

	return posts, nil
}

func (r *PostRepository) FindByID(id int) (*model.Post, error) {
	var post model.Post

	query := `
	SELECT 
		p.*, t.id, t.label 
	FROM 
		post p 
	LEFT JOIN 
		post_tag pt ON p.id = pt.post_id 
	LEFT JOIN 
		tag t ON pt.tag_id = t.id 
	WHERE 
		p.deleted_at IS NULL AND p.id = $1
	`

	rows, err := r.db.Query(
		query,
		id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	postsMap := make(map[int]*model.Post)
	for rows.Next() {
		var tagID sql.NullInt64
		var tagLabel sql.NullString

		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Status, &post.PublishDate, &post.DeletedAt, &tagID, &tagLabel)
		if err != nil {
			return nil, err
		}

		if existingPost, ok := postsMap[post.ID]; ok {
			if tagID.Valid {
				existingPost.Tags = append(existingPost.Tags, model.Tag{ID: int(tagID.Int64), Label: tagLabel.String})
			}
		} else {
			if tagID.Valid {
				post.Tags = []model.Tag{{ID: int(tagID.Int64), Label: tagLabel.String}}
			}
			postsMap[post.ID] = &post
		}
	}
	return &post, nil
}

func (r *PostRepository) Delete(id int) (rows int64, err error) {
	query := `UPDATE post SET deleted_at = NOW() WHERE id = $1`

	sql, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	rows, err = sql.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (r *PostRepository) HardDelete(id int) (rows int64, err error) {
	query := `DELETE FROM post WHERE id = $1`

	sql, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	rows, err = sql.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rows, nil
}

func (r *PostRepository) Create(post *model.Post, tags *[]model.Tag) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	tagIds := []int{}
	for _, tagValue := range *tags {
		tag := tagValue

		err := tx.QueryRow(
			`SELECT 
				t.id 
			FROM 
				tag t
			WHERE 
			t.label = $1`,
			tag.Label,
		).Scan(&tag.ID)

		if err != nil {
			if err == sql.ErrNoRows {
				tag.ID = 0
			} else {
				tx.Rollback()
				return 0, err
			}
		}

		if tag.ID == 0 {
			err := tx.QueryRow(
				`INSERT INTO 
					tag (label) 
				VALUES 
					($1) 
				RETURNING 
					id`,
				tag.Label,
			).Scan(&tag.ID)

			if err != nil {
				tx.Rollback()

				return 0, err
			}
		}

		tagIds = append(tagIds, tag.ID)
	}

	err = tx.QueryRow(
		`INSERT INTO 
			post (title, content, status) 
		VALUES 
			($1, $2, $3) 
		RETURNING 
			id`,
		post.Title,
		post.Content,
		post.Status,
	).Scan(&post.ID)

	if err != nil {
		tx.Rollback()
		return 0, err

	}

	for _, tagId := range tagIds {
		_, err := tx.Exec(
			`INSERT INTO
				post_tag (post_id, tag_id)
			VALUES
				($1, $2)`,
			post.ID,
			tagId,
		)

		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return post.ID, nil

}

func (r *PostRepository) Update(post *model.Post, tags *[]model.Tag) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	tagIds := []int{}
	for _, tagValue := range *tags {
		tag := tagValue

		err := tx.QueryRow(
			`SELECT 
				t.id 
			FROM 
				tag t
			WHERE 
			t.label = $1`,
			tag.Label,
		).Scan(&tag.ID)

		if err != nil {
			if err == sql.ErrNoRows {
				tag.ID = 0
			} else {
				tx.Rollback()
				return 0, err
			}
		}

		if tag.ID == 0 {
			err := tx.QueryRow(
				`INSERT INTO 
					tag (label) 
				VALUES 
					($1) 
				RETURNING 
					id`,
				tag.Label,
			).Scan(&tag.ID)

			if err != nil {
				tx.Rollback()

				return 0, err
			}
		}

		tagIds = append(tagIds, tag.ID)
	}

	sql, err := tx.Exec(
		`UPDATE 
			post
		SET 
			title = $1, content = $2, status = $3 
		WHERE 
			id = $4`,
		post.Title,
		post.Content,
		post.Status,
		post.ID,
	)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	row, err := sql.RowsAffected()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if row < 1 {
		tx.Rollback()
		return 0, nil
	}

	_, err = tx.Exec(
		`DELETE FROM 
			post_tag pt
		WHERE 
			pt.post_id = $1`,
		post.ID,
	)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, tagId := range tagIds {
		_, err := tx.Exec(
			`INSERT INTO
				post_tag (post_id, tag_id)
			VALUES
				($1, $2)`,
			post.ID,
			tagId,
		)

		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return post.ID, nil
}

func (r *PostRepository) FindByTag(tag string) ([]model.Post, error) {
	query := `
	SELECT 
		p.*, t.id, t.label 
	FROM 
		post p 
	LEFT JOIN 
		post_tag pt ON p.id = pt.post_id 
	LEFT JOIN 
		tag t ON pt.tag_id = t.id 
	WHERE 
		p.deleted_at IS NULL AND t.label = $1
	`

	rows, err := r.db.Query(query, tag)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	postsMap := make(map[int]*model.Post)
	for rows.Next() {
		var post model.Post

		var tagID sql.NullInt64
		var tagLabel sql.NullString
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Status, &post.PublishDate, &post.DeletedAt, &tagID, &tagLabel)
		if err != nil {
			return nil, err
		}

		if existingPost, ok := postsMap[post.ID]; ok {
			if tagID.Valid {
				existingPost.Tags = append(existingPost.Tags, model.Tag{ID: int(tagID.Int64), Label: tagLabel.String})
			}
		} else {
			if tagID.Valid {
				post.Tags = []model.Tag{{ID: int(tagID.Int64), Label: tagLabel.String}}
			}
			postsMap[post.ID] = &post
		}
	}

	var posts []model.Post
	for _, post := range postsMap {
		posts = append(posts, *post)
	}

	return posts, nil
}

func (r *PostRepository) FindIdByTag(tag string) ([]int, error) {
	query := `
	SELECT 
		p.id
	FROM 
		post p 
	LEFT JOIN 
		post_tag pt ON p.id = pt.post_id 
	LEFT JOIN 
		tag t ON pt.tag_id = t.id 
	WHERE 
		p.deleted_at IS NULL AND t.label = $1
	`

	rows, err := r.db.Query(query, tag)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var postIds []int
	for rows.Next() {
		var postId int
		err := rows.Scan(&postId)
		if err != nil {
			return nil, err
		}
		postIds = append(postIds, postId)
	}

	return postIds, nil
}

func (r *PostRepository) FindByIds(ids []int) ([]model.Post, error) {

	iQuery, iValue := helper.QueryIN(ids)

	query := fmt.Sprintf(`
	SELECT
		p.*, t.id, t.label
	FROM
		post p
	LEFT JOIN
		post_tag pt ON p.id = pt.post_id
	LEFT JOIN
		tag t ON pt.tag_id = t.id
	WHERE
		p.deleted_at IS NULL
	AND
		p.id IN (%v)`, iQuery)

	rows, err := r.db.Query(query, iValue...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	postsMap := make(map[int]*model.Post)
	for rows.Next() {
		var post model.Post

		var tagID sql.NullInt64
		var tagLabel sql.NullString
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Status, &post.PublishDate, &post.DeletedAt, &tagID, &tagLabel)
		if err != nil {
			return nil, err
		}

		if existingPost, ok := postsMap[post.ID]; ok {
			if tagID.Valid {
				existingPost.Tags = append(existingPost.Tags, model.Tag{ID: int(tagID.Int64), Label: tagLabel.String})
			}
		} else {
			if tagID.Valid {
				post.Tags = []model.Tag{{ID: int(tagID.Int64), Label: tagLabel.String}}
			}
			postsMap[post.ID] = &post
		}

	}

	var posts []model.Post
	for _, post := range postsMap {
		posts = append(posts, *post)
	}

	return posts, nil
}

func (r *PostRepository) PublishPost(id int) (int, error) {
	query := `UPDATE post SET status = 'publish' WHERE id = $1`

	sql, err := r.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	row, err := sql.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(row), nil
}
