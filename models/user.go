package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"createdAt"`
}

type CreateUserRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (d *PostgresDatabase) CreateUser(userInfo CreateUserRequest) (int, error) {
	switch {
	case userInfo.Username == "":
		return 0, errors.New("username missing")
	case userInfo.FirstName == "":
		return 0, errors.New("first name missing")
	case userInfo.LastName == "":
		return 0, errors.New("last name missing")
	case userInfo.Password == "":
		return 0, errors.New("password missing")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), 12)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO users (username, first_name, last_name, hashed_password)
	VALUES($1,$2,$3,$4) RETURNING id`

	var id int
	err = d.db.QueryRow(query, userInfo.Username, userInfo.FirstName, userInfo.LastName, hashedPassword).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *PostgresDatabase) GetUsers() ([]*User, error) {
	rows, err := d.db.Query(`select * from users`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		user := new(User)
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.HashedPassword,
			&user.CreatedAt)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (d *PostgresDatabase) GetUserById(id int) (*User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	row := d.db.QueryRow(query, id)

	user := new(User)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.HashedPassword,
		&user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *PostgresDatabase) AuthenticateUser(username, password string) (int, error) {
	query := `SELECT id, hashed_password FROM users WHERE username = $1`
	var hashedPassword []byte
	var id int
	err := d.db.QueryRow(query, username).Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		// User not found.
		return 0, errors.New("user not found")
	} else if err != nil {
		// Other database-related errors.
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		// Password comparison failed.
		return 0, errors.New("authentication failed")
	}
	return id, nil
}
