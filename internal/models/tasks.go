package models

import (
	"database/sql"
	"errors"
	"time"
)

type Task struct {
	ID          int
	Title       string
	Content     string
	Priority    string
	Created     time.Time
	Finished    time.Time
	WorkspaceId int
}

type TaskModelInterface interface {
	Insert(title, content, priority string, workspaceId int) (int, error)
	Get(id int) (Task, error)
	GetAll(workspaceId int) ([]Task, error)
	Update(id int, title, content, priority string) error
	Delete(id int) (int, error)
	ValidateOwnership(userId, taskId int) (bool, error)
}

type TaskModel struct {
	DB *sql.DB
}

func (m *TaskModel) Insert(title, content, priority string, workspaceId int) (int, error) {
	stmt := `INSERT INTO tasks (title, content, priority, created, finished, workspace_id)  VALUES (?, ?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL 2 WEEK), ?)`

	result, err := m.DB.Exec(stmt, title, content, priority, workspaceId)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *TaskModel) Get(id int) (Task, error) {
	stmt := `SELECT * FROM tasks WHERE id = ?`

	var t Task

	err := m.DB.QueryRow(stmt, id).Scan(&t.ID, &t.Title, &t.Content, &t.Priority, &t.Created, &t.Finished, &t.WorkspaceId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrNoRecord
		} else {
			return Task{}, err
		}
	}

	return t, nil
}

func (m *TaskModel) GetAll(workspaceId int) ([]Task, error) {
	stmt := `SELECT * FROM tasks where workspace_id = ?`

	rows, err := m.DB.Query(stmt, workspaceId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task

		err = rows.Scan(&t.ID, &t.Title, &t.Content, &t.Priority, &t.Created, &t.Finished, &t.WorkspaceId)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m *TaskModel) Update(id int, title, content, priority string) error {
	stmt := `UPDATE tasks SET title = ?, content = ?, priority = ? where id = ?`

	_, err := m.DB.Exec(stmt, title, content, priority, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *TaskModel) Delete(id int) (int, error) {
	stmt := `DELETE FROM tasks where id = ?`

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

func (m *TaskModel) ValidateOwnership(userId, taskId int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS (SELECT true FROM tasks JOIN users_workspaces uw ON tasks.workspace_id = uw.workspace_id WHERE tasks.id = ? AND uw.user_id = ?)"

	err := m.DB.QueryRow(stmt, taskId, userId).Scan(&exists)
	return exists, err
}
