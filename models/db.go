package models

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database interface {
	CreateUsersTable() error
	CreatePostsTable() error
	CreateCommentsTable() error
	CreateUser(CreateUserRequest) (int, error)
	GetUsers() ([]*User, error)
	GetUserById(int) (*User, error)
	AuthenticateUser(string, string) (int, error)
	CreatePost(NewPostRequest) (int, error)
	GetPosts() ([]*Post, error)
	GetPostById(int) (*Post, error)
	CreateComment(CreateCommentRequest) error
	GetCommentsByPostId(int) ([]*Comment, error)
}

type PostgresDatabase struct {
	db *sql.DB
}

func NewPostgresDatabase() (*PostgresDatabase, error) {
	connStr := "user=postgres dbname=blog password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// Ping the db to ensure connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	database := &PostgresDatabase{db}

	// Initialize database tables in the case they do not exist
	if err := database.CreateCommentsTable(); err != nil {
		return nil, err
	}
	if err := database.CreatePostsTable(); err != nil {
		return nil, err
	}
	if err := database.CreateUsersTable(); err != nil {
		return nil, err
	}

	return database, nil
}

func (d *PostgresDatabase) CreateUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		hashed_password VARCHAR(255),
		created timestamp default current_timestamp
	)`
	_, err := d.db.Exec(query)
	return err
}

func (d *PostgresDatabase) CreatePostsTable() error {
	query := `CREATE TABLE IF NOT EXISTS posts (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255),
        content TEXT,
        user_id INT,
        created TIMESTAMP DEFAULT current_timestamp
    )`
	_, err := d.db.Exec(query)
	return err
}

func (d *PostgresDatabase) CreateCommentsTable() error {
	query := `CREATE TABLE IF NOT EXISTS comments (
        id SERIAL PRIMARY KEY,
        post_id INT,
        content TEXT,
        user_id INT,
        created TIMESTAMP DEFAULT current_timestamp
    )`
	_, err := d.db.Exec(query)
	return err
}
