package mocks

import (
	"time"

	"github.com/andres085/task_manager/internal/models"
)

var firstMockTask = models.Task{
	ID:       1,
	Title:    "First Test Task",
	Content:  "First Test Task Content",
	Priority: "LOW",
	Created:  time.Now(),
	Finished: nil,
	Status:   "To Do",
}

var secondMockTask = models.Task{
	ID:       2,
	Title:    "Second Test Task",
	Content:  "Second Test Task Content",
	Priority: "MEDIUM",
	Created:  time.Now(),
	Finished: nil,
	Status:   "To Do",
}

type TaskModel struct{}

func (t *TaskModel) Insert(title, content, priority string, workspaceId, userId int) (int, error) {
	return 2, nil
}

func (t *TaskModel) Get(id int) (models.Task, error) {
	switch id {
	case 1:
		return firstMockTask, nil
	default:
		return models.Task{}, models.ErrNoRecord
	}
}

func (m *TaskModel) GetAll(id, limit, offset int, query string) ([]models.Task, error) {
	return []models.Task{firstMockTask, secondMockTask}, nil
}

func (m *TaskModel) GetTotalTasks(workspaceId int) (int, error) {
	return 0, nil
}

func (m *TaskModel) Update(id int, title, content, priority string, userId int, status string) error {
	return nil
}

func (m *TaskModel) Delete(id int) (int, error) {
	return 1, nil
}

func (m *TaskModel) ValidateOwnership(userId, taskId int) (bool, error) {
	if userId == 1 && taskId == 1 {
		return true, nil
	}
	return false, nil
}

func (m *TaskModel) ValidateAdmin(userId, taskId int) (bool, error) {
	if userId == 1 && taskId == 1 {
		return true, nil
	}
	return false, nil
}
