package models

import (
	"database/sql"
	"fmt"
	"fourth-week/bcryptPassword"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Password string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(name, username, password string) (*User, error) {
	username = strings.ToLower(username)
	hash, _ := bcryptPassword.HashPassword(password)
	PasswordHash := string(hash)

	user := User{
		Name:     name,
		Username: username,
		Password: PasswordHash,
	}

	row := us.DB.QueryRow(`
	INSERT INTO users(name, username, password) VALUES ($1, $2, $3);
	`, name, username, PasswordHash)

	err := row.Scan(&user.ID)

	if err != nil {
		return nil, fmt.Errorf("Create user: %w", err)
	}

	return &user, nil
}

func (us *UserService) Authenticate(username, password string) (*User, error) {
	username = strings.ToLower(username)
	user := User{
		Username: username,
	}

	row := us.DB.QueryRow(`
		SELECT id, password FROM users WHERE username = $1
	`, username)

	err := row.Scan(&user.ID, &user.Password)

	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	match := bcryptPassword.CheckPasswordHash(password, user.Password)
	if match {
		return &user, nil
	} else {
		return nil, fmt.Errorf("authenticate: %w", err)
	}
}
