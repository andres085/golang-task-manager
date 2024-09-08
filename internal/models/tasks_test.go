package models

import (
	"testing"

	"github.com/andres085/task_manager/internal/assert"
)

func TestGetMethod(t *testing.T) {
	taskId := 1

	db := newTestDB(t)

	m := TaskModel{db}

	task, err := m.Get(taskId)

	assert.Equal(t, task.ID, taskId)
	assert.NilError(t, err)
}

func TestInsertMethod(t *testing.T) {
	db := newTestDB(t)

	m := TaskModel{db}

	id, err := m.Insert("Test Task", "Test Task Body", "HIGH")

	assert.Equal(t, id, 4)
	assert.NilError(t, err)
}

func TestGetAllMethod(t *testing.T) {
	db := newTestDB(t)

	m := TaskModel{db}

	tasks, err := m.GetAll()

	assert.Equal(t, len(tasks), 3)
	assert.NilError(t, err)
}
