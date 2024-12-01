package models

import (
	"database/sql"
	"errors"
	"time"
)

type Workspace struct {
	ID          int
	Title       string
	Description string
	Created     time.Time
}

type WorkspaceModelInterface interface {
	Insert(title, description string, userId int) (int, error)
	Get(id int) (Workspace, error)
	GetAll(userId int, role string) ([]Workspace, error)
	Update(id int, title, description string) error
	Delete(id int) (int, error)
	ValidateOwnership(userId, workspaceId int) (bool, error)
	ValidateAdmin(userId, workspaceId int) (bool, error)
}

type WorkspaceModel struct {
	DB *sql.DB
}

func (m *WorkspaceModel) Insert(title, description string, userId int) (int, error) {
	stmt := `INSERT INTO workspaces (title, description, created)  VALUES (?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, title, description)
	if err != nil {
		return 0, err
	}

	workspaceId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	stmt = `INSERT INTO users_workspaces(user_id, workspace_id, role, created) VALUES (?, ?, ?, UTC_TIMESTAMP)`
	result, err = m.DB.Exec(stmt, userId, workspaceId, "ADMIN")
	if err != nil {
		return 0, err
	}

	return int(workspaceId), nil
}

func (m *WorkspaceModel) Get(id int) (Workspace, error) {
	stmt := `SELECT * FROM workspaces WHERE id = ?`

	var w Workspace

	err := m.DB.QueryRow(stmt, id).Scan(&w.ID, &w.Title, &w.Description, &w.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Workspace{}, ErrNoRecord
		} else {
			return Workspace{}, err
		}
	}

	return w, nil
}

func (m *WorkspaceModel) GetAll(userId int, role string) ([]Workspace, error) {
	stmt := `SELECT w.* FROM workspaces as w JOIN users_workspaces as uw ON w.id = uw.workspace_id WHERE uw.user_id = ? and uw.role = ?`

	rows, err := m.DB.Query(stmt, userId, role)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var workspaces []Workspace

	for rows.Next() {
		var w Workspace

		err = rows.Scan(&w.ID, &w.Title, &w.Description, &w.Created)
		if err != nil {
			return nil, err
		}

		workspaces = append(workspaces, w)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (m *WorkspaceModel) Update(id int, title, description string) error {
	stmt := `UPDATE workspaces SET title = ?, description = ? where id = ?`

	_, err := m.DB.Exec(stmt, title, description, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *WorkspaceModel) Delete(id int) (int, error) {
	stmt := `DELETE FROM workspaces where id = ?`

	result, err := m.DB.Exec(stmt, id)
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

func (m *WorkspaceModel) ValidateOwnership(userId, workspaceId int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users_workspaces WHERE user_id = ? AND workspace_id = ?)"

	err := m.DB.QueryRow(stmt, userId, workspaceId).Scan(&exists)
	return exists, err
}

func (m *WorkspaceModel) ValidateAdmin(userId, workspaceId int) (bool, error) {
	var isAdmin bool

	stmt := "SELECT EXISTS(SELECT true FROM users_workspaces WHERE user_id = ? AND workspace_id = ? AND role = 'ADMIN')"

	err := m.DB.QueryRow(stmt, userId, workspaceId).Scan(&isAdmin)
	return isAdmin, err
}
