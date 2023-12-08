package models

import (
	"database/sql"
	"fmt"
	"log"
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
	GetPosts(limit, offset int) ([]*Post, int, error)
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
	dbHost, dbPort, dbUser, dbPassword, dbName :=
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME")

	fmt.Printf("Connecting to DB with parameters:\nhost=%s port=%s user=%s password=%s dbname=%s\n", dbHost, dbPort, dbUser, dbPassword, dbName)
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Printf("Error opening database connection: %v\n", err)
		return nil, err
	}
	fmt.Println("success")
	// Ping the db to ensure connection
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging database: %v\n", err)
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
		username VARCHAR(255) UNIQUE NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		hashed_password VARCHAR(255) NOT NULL,
		created timestamp default current_timestamp
	)`
	_, err := d.db.Exec(query)
	return err
}

func (d *PostgresDatabase) CreatePostsTable() error {
	query := `CREATE TABLE IF NOT EXISTS posts (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        user_id INT NOT NULL,
		logo_path VARCHAR(255) DEFAULT 'none',
        created TIMESTAMP DEFAULT current_timestamp
    )`
	_, err := d.db.Exec(query)
	return err
}

func (d *PostgresDatabase) CreateCommentsTable() error {
	query := `CREATE TABLE IF NOT EXISTS comments (
        id SERIAL PRIMARY KEY,
        post_id INT NOT NULL,
        content TEXT NOT NULL,
        user_id INT NOT NULL,
        created TIMESTAMP DEFAULT current_timestamp
    )`
	_, err := d.db.Exec(query)
	return err
}
