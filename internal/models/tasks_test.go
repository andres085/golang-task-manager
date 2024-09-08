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
