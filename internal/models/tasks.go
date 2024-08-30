package models

import (
	"database/sql"
	"errors"
	"time"
)

type Task struct {
	ID       int
	Title    string
	Content  string
	Priority string
	Created  time.Time
	Finished *time.Time
}

type TaskModel struct {
	DB *sql.DB
}

func (m *TaskModel) Insert(title, content, priority string) (int, error) {
	stmt := `INSERT INTO tasks (title, content, priority, created)  VALUES (?, ?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, title, content, priority)
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

	err := m.DB.QueryRow(stmt, id).Scan(&t.ID, &t.Title, &t.Content, &t.Priority, &t.Created, &t.Finished)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrNoRecord
		} else {
			return Task{}, err
		}
	}

	return t, nil
}

func (m *TaskModel) GetAll() ([]Task, error) {
	return nil, nil
}

func (m *TaskModel) Update(id int) (Task, error) {
	return Task{}, nil
}

func (m *TaskModel) Delete(id int) (int, error) {
	return 0, nil
}
