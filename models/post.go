package models

import (
	"errors"
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	UserId    int       `json:"-"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}

type NewPostRequest struct {
	Title   string `json:"title"`
	UserId  int    `json:"-"`
	Content string `json:"content"`
}

type UpdatePostRequest struct {
	Title   string `json:"title"`
	UserId  int    `json:"-"`
	Content string `json:"content"`
}

func (d *PostgresDatabase) GetPosts() ([]*Post, error) {
	query := `SELECT * FROM posts`
	rows, err := d.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := []*Post{}
	for rows.Next() {
		post := new(Post)
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserId,
			&post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (d *PostgresDatabase) GetPostById(id int) (*Post, error) {
	query := `SELECT * FROM posts WHERE id = $1`
	row := d.db.QueryRow(query, id)

	post := new(Post)
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserId,
		&post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (d *PostgresDatabase) CreatePost(p NewPostRequest) (int, error) {
	query := `INSERT INTO posts (title, user_id, content) VALUES($1,$2,$3) RETURNING id`
	switch {
	case p.Title == "":
		return 0, errors.New("title missing")
	case p.Content == "":
		return 0, errors.New("content missing")
	case p.UserId <= 0:
		return 0, errors.New("invalid user id")
	}
	var postID int
	err := d.db.QueryRow(query, p.Title, p.UserId, p.Content).Scan(&postID)
	if err != nil {
		return 0, err
	}
	return postID, nil
}

func (d *PostgresDatabase) UpdatePostById(id int, req UpdatePostRequest) error {
	query := `UPDATE posts SET title = $1, content = $2 WHERE id = $3`
	_, err := d.db.Exec(query, req.Title, req.Content, req.UserId)

	if err != nil {
		return err
	}
	return nil
}

func (d *PostgresDatabase) DeletePostById(id int) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := d.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
