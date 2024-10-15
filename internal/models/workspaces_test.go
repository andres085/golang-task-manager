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

	id, err := m.Insert("Test Workspace", "Test Workspace Description", 1)

	assert.Equal(t, id, 2)
	assert.NilError(t, err)
}

func TestWorkspacesGetAllMethod(t *testing.T) {
	db := newTestDB(t)

	m := WorkspaceModel{db}

	workspaces, err := m.GetAll(1)

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

func TestValidateOwnership(t *testing.T) {
	db := newTestDB(t)

	m := WorkspaceModel{db}

	isOwner, err := m.ValidateOwnership(1, 1)

	assert.Equal(t, isOwner, true)
	assert.NilError(t, err)
}
