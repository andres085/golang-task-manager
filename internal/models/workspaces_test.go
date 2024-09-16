package models

import (
	"testing"

	"github.com/andres085/task_manager/internal/assert"
)

func TestWorkspacesGetMethod(t *testing.T) {
	workspaceId := 1

	db := newTestDB(t)

	m := WorkspaceModel{db}

	workspace, err := m.Get(workspaceId)

	assert.Equal(t, workspace.ID, workspaceId)
	assert.NilError(t, err)
}

func TestWorkspacesInsertMethod(t *testing.T) {
	db := newTestDB(t)

	m := WorkspaceModel{db}

	id, err := m.Insert("Test Workspace", "Test Workspace Description")

	assert.Equal(t, id, 2)
	assert.NilError(t, err)
}

func TestWorkspacesGetAllMethod(t *testing.T) {
	db := newTestDB(t)

	m := WorkspaceModel{db}

	workspaces, err := m.GetAll()

	assert.Equal(t, len(workspaces), 1)
	assert.NilError(t, err)
}

func TestWorkspacesUpdateMethod(t *testing.T) {
	db := newTestDB(t)

	m := WorkspaceModel{db}

	newTitle := "Updated Title"
	err := m.Update(1, newTitle, "Test Task Description")

	assert.NilError(t, err)

	updatedWorkspace, err := m.Get(1)

	assert.Equal(t, updatedWorkspace.Title, newTitle)
}

func TestWorkspacesDeleteMethod(t *testing.T) {
	db := newTestDB(t)

	m := WorkspaceModel{db}

	row, err := m.Delete(1)

	assert.Equal(t, row, 1)
	assert.NilError(t, err)
}

func TestWorkspacesDeleteWithTransacionMethod(t *testing.T) {
	db := newTestDB(t)

	m := WorkspaceModel{db}
	m2 := TaskModel{db}

	row, err := m.Delete(1)
	tasks, err := m2.GetAll(1)

	assert.Equal(t, len(tasks), 0)
	assert.Equal(t, row, 1)
	assert.NilError(t, err)
}
