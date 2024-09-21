package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	FirstName      string
	LastName       string
	Email          string
	HashedPassword string
	Created        time.Time
}

type UserModelInterface interface {
	Insert(firstName, lastName, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(firstName, lastName, email, password string) error {
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
