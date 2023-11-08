package models

import "time"

type Comment struct {
	ID        string    `json:"id"`
	PostId    int       `json:"post"`
	Content   string    `json:"content"`
	UserId    int       `json:"user"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateCommentRequest struct {
	PostId  int    `json:"post"`
	Content string `json:"content"`
	UserId  int    `json:"user"`
}

func (d *PostgresDatabase) GetCommentsByPostId(id int) ([]*Comment, error) {
	query := `SELECT * FROM comments WHERE post_id = $1`
	rows, err := d.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*Comment{}
	for rows.Next() {
		comment := new(Comment)
		err = rows.Scan(
			&comment.ID,
			&comment.PostId,
			&comment.Content,
			&comment.UserId,
			&comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, nil
}

func (d *PostgresDatabase) CreateComment(req CreateCommentRequest) error {
	query := `INSERT INTO comments (user_id, post_id, content) VALUES ($1, $2, $3)`
	_, err := d.db.Exec(query, req.UserId, req.PostId, req.Content)
	if err != nil {
		return err
	}
	return nil
}
