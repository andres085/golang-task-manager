package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	GetByEmail(email string) (*User, error)
	AddUserToWorkspace(userId, workspaceId int) error
	GetWorkspaceUsers(workspaceId int) ([]UserWithRole, error)
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(firstName, lastName, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (firstName, lastName, email, hashed_password, created) VALUES (?, ?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, firstName, lastName, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	stmt := "SELECT * FROM users WHERE email = ?"

	var u User

	err := m.DB.QueryRow(stmt, email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.HashedPassword, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &User{}, ErrNoRecord
		} else {
			return &User{}, err
		}
	}

	return &u, nil
}

type UserWithRole struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Role      string
}

func (m *UserModel) GetWorkspaceUsers(workspaceId int) ([]UserWithRole, error) {
	stmt := "SELECT u.id, u.firstName, u.lastName, u.email, uw.`role` FROM users u JOIN users_workspaces uw ON u.id = uw.user_id WHERE uw.workspace_id = ?"

	rows, err := m.DB.Query(stmt, workspaceId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []UserWithRole

	for rows.Next() {
		var u UserWithRole

		err = rows.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Role)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil

}

func (m *UserModel) AddUserToWorkspace(userId, workspaceId int) error {
	stmt := `INSERT INTO users_workspaces (user_id, workspace_id, role, created) VALUES (?, ?, "MEMBER", UTC_TIMESTAMP())`

	_, err := m.DB.Exec(stmt, userId, workspaceId)
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
