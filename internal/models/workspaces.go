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
	GetAll() ([]Workspace, error)
	Update(id int, title, description string) error
	Delete(id int) (int, error)
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

func (m *WorkspaceModel) Update(id int, title, description string) error {
	stmt := `UPDATE workspaces SET title = ?, description = ? where id = ?`

	_, err := m.DB.Exec(stmt, title, description, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *WorkspaceModel) Delete(id int) (int, error) {

	tasksStmt := `SELECT * FROM tasks WHERE workspace_id = ?`

	rows, err := m.DB.Query(tasksStmt, id)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if rows.Next() {
		return m.DeleteWithTransaction(id)
	}

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

func (m *WorkspaceModel) DeleteWithTransaction(id int) (int, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	tasksStmt := `DELETE FROM tasks WHERE workspace_id = ?`

	_, err = tx.Exec(tasksStmt, id)
	if err != nil {
		return 0, err
	}

	stmt := `DELETE FROM workspaces WHERE id = ?`

	result, err := tx.Exec(stmt, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}
