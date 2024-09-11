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
	Insert(title, description string) (int, error)
	Get(id int) (Workspace, error)
	GetAll() ([]Workspace, error)
}

type WorkspaceModel struct {
	DB *sql.DB
}

func (m *WorkspaceModel) Insert(title, description string) (int, error) {
	stmt := `INSERT INTO workspaces (title, description,  created)  VALUES (?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, title, description)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
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

func (m *WorkspaceModel) GetAll() ([]Workspace, error) {
	stmt := `SELECT * FROM workspaces`

	rows, err := m.DB.Query(stmt)
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
