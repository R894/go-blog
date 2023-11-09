package models

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Database interface {
	// Initialize tables
	CreateUsersTable() error
	CreatePostsTable() error
	CreateCommentsTable() error

	// Posts
	CreatePost(NewPostRequest) (int, error)
	GetPosts() ([]*Post, error)
	GetPostById(int) (*Post, error)
	UpdatePostById(int, UpdatePostRequest) error
	DeletePostById(int) error

	// Comments
	CreateComment(CreateCommentRequest) error
	GetCommentsByPostId(int) ([]*Comment, error)

	// User related
	CreateUser(CreateUserRequest) (int, error)
	AuthenticateUser(string, string) (int, error)
	GetUsers() ([]*User, error)
	GetUserById(int) (*User, error)
}

type PostgresDatabase struct {
	db *sql.DB
}

func NewPostgresDatabase() (*PostgresDatabase, error) {
	dbUser, dbPassword, dbName := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", dbUser, dbName, dbPassword)
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
