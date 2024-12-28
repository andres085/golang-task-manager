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
	GetUser(userId int) (*User, error)
	GetUserToInvite(email string, workspaceId int) (*User, error)
	AddUserToWorkspace(userId, workspaceId int) error
	GetWorkspaceUsers(workspaceId int) ([]UserWithRole, error)
	GetWorkspacesAsMemberCount(email string) (int, error)
	RemoveUserFromWorkspace(workspaceId, userId int) (int, error)
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

func (m *UserModel) GetUser(userId int) (*User, error) {
	stmt := "SELECT firstName, lastName FROM users where id = ?"

	var u User

	err := m.DB.QueryRow(stmt, userId).Scan(&u.FirstName, &u.LastName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &User{}, ErrNoRecord
		} else {
			return &User{}, err
		}
	}

	return &u, nil
}

func (m *UserModel) GetUserToInvite(email string, workspaceId int) (*User, error) {
	stmt := "SELECT u.* FROM users u LEFT JOIN users_workspaces uw ON u.id = uw.user_id AND uw.workspace_id = ? WHERE u.email = ? AND uw.workspace_id IS NULL"

	var u User

	err := m.DB.QueryRow(stmt, workspaceId, email).Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.HashedPassword, &u.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &User{}, ErrNoRecord
		} else {
			return &User{}, err
		}
	}

	return &u, nil
}

func (m *UserModel) GetWorkspacesAsMemberCount(email string) (int, error) {
	var totalWorkspaces int

	countStmt := "SELECT COUNT(*) FROM users u LEFT JOIN users_workspaces uw ON u.id = uw.user_id WHERE u.email = ? AND uw.`role` = 'MEMBER';"

	err := m.DB.QueryRow(countStmt, email).Scan(&totalWorkspaces)
	if err != nil {
		return 0, err
	}

	return totalWorkspaces, nil
}

type UserWithRole struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Role      string
}

func (m *UserModel) GetWorkspaceUsers(workspaceId int) ([]UserWithRole, error) {
	stmt := "SELECT u.id, u.firstName, u.lastName, u.email, uw.`role` FROM users u JOIN users_workspaces uw ON u.id = uw.user_id WHERE uw.workspace_id = ? ORDER BY role;"

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

func (m *UserModel) RemoveUserFromWorkspace(workspaceId, userId int) (int, error) {

	stmt := "DELETE FROM users_workspaces WHERE workspace_id = ? AND user_id = ? AND `role` != 'ADMIN';"

	result, err := m.DB.Exec(stmt, workspaceId, userId)
	if err != nil {
		return 0, err
	}

	var r int64
	r, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(r), nil
}
