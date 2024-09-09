package models

import (
	"database/sql"
	"time"
)

type Workspace struct {
	ID          int
	Title       string
	Description string
	Created     time.Time
}

type WorkspaceModelInterface interface {
	GetAll() ([]Workspace, error)
}

type WorkspaceModel struct {
	DB *sql.DB
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
