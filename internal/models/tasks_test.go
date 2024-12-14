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

	id, err := m.Insert("Test Task", "Test Task Body", "HIGH", 1, 1)

	assert.Equal(t, id, 4)
	assert.NilError(t, err)
}

func TestGetAllMethod(t *testing.T) {
	db := newTestDB(t)

	m := TaskModel{db}
	id := 1

	tasks, err := m.GetAll(id, 10, 0)

	assert.Equal(t, len(tasks), 3)
	assert.NilError(t, err)
}

func TestUpdateMethod(t *testing.T) {
	db := newTestDB(t)

	m := TaskModel{db}

	newTitle := "Updated Title"
	err := m.Update(1, newTitle, "Test Task Body", "HIGH", 1, "Completed")

	assert.NilError(t, err)

	updatedTask, err := m.Get(1)

	assert.Equal(t, updatedTask.Title, newTitle)
}

func TestDeleteMethod(t *testing.T) {
	db := newTestDB(t)

	m := TaskModel{db}

	row, err := m.Delete(1)

	assert.Equal(t, row, 1)
	assert.NilError(t, err)
}

func TestValidateTaskOwnership(t *testing.T) {
	db := newTestDB(t)

	m := TaskModel{db}

	isOwner, err := m.ValidateOwnership(1, 1)

	assert.Equal(t, isOwner, true)
	assert.NilError(t, err)
}

func TestValidateTaskAdmin(t *testing.T) {
	db := newTestDB(t)

	m := TaskModel{db}

	isAdmin, err := m.ValidateAdmin(1, 1)

	assert.Equal(t, isAdmin, true)
	assert.NilError(t, err)
}
