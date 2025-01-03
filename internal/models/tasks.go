package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID          int
	Title       string
	Content     string
	Priority    string
	Created     time.Time
	Finished    *time.Time
	WorkspaceId int
	UserId      int
	Status      string
}

type TaskModelInterface interface {
	Insert(title, content, priority string, workspaceId, userId int) (int, error)
	Get(id int) (Task, error)
	GetAll(workspaceId, limit, offset int, title, priority, status, sort string) ([]Task, error)
	GetTotalTasks(workspaceId int, title, priority, status string) (int, error)
	Update(id int, title, content, priority string, userId int, status string) error
	Delete(id int) (int, error)
	ValidateOwnership(userId, taskId int) (bool, error)
	ValidateAdmin(userId, taskId int) (bool, error)
}

type TaskModel struct {
	DB *sql.DB
}

func (m *TaskModel) Insert(title, content, priority string, workspaceId, userId int) (int, error) {
	stmt := `INSERT INTO tasks (title, content, priority, created, workspace_id, user_id)  VALUES (?, ?, ?, UTC_TIMESTAMP(), ?, ?)`

	result, err := m.DB.Exec(stmt, title, content, priority, workspaceId, userId)
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

	err := m.DB.QueryRow(stmt, id).Scan(&t.ID, &t.Title, &t.Content, &t.Priority, &t.Created, &t.Finished, &t.WorkspaceId, &t.UserId, &t.Status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrNoRecord
		} else {
			return Task{}, err
		}
	}

	return t, nil
}

func (m *TaskModel) GetAll(workspaceId, limit, offset int, title, priority, status, sort string) ([]Task, error) {

	stmt := `SELECT * FROM tasks where workspace_id = ?`

	conditions := map[string]interface{}{
		"workspaceId": workspaceId,
		"title":       title,
		"priority":    priority,
		"status":      status,
		"sort":        sort,
		"limit":       limit,
		"offset":      offset,
	}

	preparedStmt, args := prepareStmt(stmt, conditions)

	rows, err := m.DB.Query(preparedStmt, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task

		err = rows.Scan(&t.ID, &t.Title, &t.Content, &t.Priority, &t.Created, &t.Finished, &t.WorkspaceId, &t.UserId, &t.Status)
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

func (m *TaskModel) GetTotalTasks(workspaceId int, title, priority, status string) (int, error) {
	var totalTasks int

	countStmt := `SELECT COUNT(*) FROM tasks WHERE workspace_id = ? `

	args := []interface{}{workspaceId}

	if title != "" {
		countStmt += "AND title LIKE ? "
		args = append(args, "%"+title+"%")
	}
	if priority != "" {
		countStmt += "AND priority = ? "
		args = append(args, priority)
	}
	if status != "" {
		countStmt += "AND status = ? "
		args = append(args, status)
	}

	err := m.DB.QueryRow(countStmt, args...).Scan(&totalTasks)
	if err != nil {
		return 0, err
	}

	return totalTasks, nil
}

func (m *TaskModel) Update(id int, title, content, priority string, userId int, status string) error {
	var finished *time.Time

	if status == "Completed" {
		now := time.Now()
		finished = &now
	} else {
		finished = nil
	}

	stmt := `UPDATE tasks SET title = ?, content = ?, priority = ?, user_id = ?, status = ?, finished = ? where id = ?`

	_, err := m.DB.Exec(stmt, title, content, priority, userId, status, finished, id)
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

func (m *TaskModel) ValidateAdmin(userId, taskId int) (bool, error) {
	var isAdmin bool

	stmt := "SELECT EXISTS (SELECT true FROM tasks JOIN users_workspaces uw ON tasks.workspace_id = uw.workspace_id WHERE tasks.id = ? AND uw.user_id = ? AND uw.role = 'ADMIN')"

	err := m.DB.QueryRow(stmt, taskId, userId).Scan(&isAdmin)
	return isAdmin, err
}

func prepareStmt(baseStmt string, conditions map[string]interface{}) (string, []interface{}) {
	workspaceId := conditions["workspaceId"]
	args := []interface{}{workspaceId}

	for column, value := range conditions {

		if column == "title" && value != "" {
			baseStmt += fmt.Sprintf(" AND %s LIKE ? ", column)
			args = append(args, fmt.Sprintf(`%%%s%%`, value))
		}
		if column == "priority" && value != "" {
			baseStmt += fmt.Sprintf(" AND %s = ? ", column)
			args = append(args, value)
		}
		if column == "status" && value != "" {
			baseStmt += fmt.Sprintf(" AND %s = ? ", column)
			args = append(args, value)
		}
	}

	sort := conditions["sort"].(string)
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	baseStmt += " ORDER BY created " + sort + " LIMIT ? OFFSET ?"
	args = append(args, conditions["limit"], conditions["offset"])

	return baseStmt, args
}
