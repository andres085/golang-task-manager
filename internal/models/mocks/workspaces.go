package mocks

import (
	"errors"

	"github.com/andres085/task_manager/internal/models"
)

var firstMockWorkspace = models.Workspace{
	ID:          1,
	Title:       "First Workspace",
	Description: "First workspace Description",
}

var secondMockWorkspace = models.Workspace{
	ID:          2,
	Title:       "Second Workspace",
	Description: "Second workspace Description",
}

type WorkspaceModel struct{}

func (t *WorkspaceModel) Insert(title, description string, userId int) (int, error) {
	return 2, nil
}

func (t *WorkspaceModel) Get(id int) (models.Workspace, error) {
	switch id {
	case 1:
		return firstMockWorkspace, nil
	default:
		return models.Workspace{}, models.ErrNoRecord
	}
}

func (m *WorkspaceModel) GetAll(userId int, role string) ([]models.Workspace, error) {
	return []models.Workspace{firstMockWorkspace, secondMockWorkspace}, nil
}

func (m *WorkspaceModel) Update(id int, title, description string) error {
	return nil
}

func (m *WorkspaceModel) Delete(id int) (int, error) {
	return 1, nil
}

func (m *WorkspaceModel) ValidateOwnership(userId, workspaceId int) (bool, error) {
	if userId == 1 && workspaceId == 1 {
		return true, nil
	}
	if userId == 1 && workspaceId == 0 {
		return false, errors.New("Internal Server Error")
	}
	return false, nil
}

func (m *WorkspaceModel) ValidateAdmin(userId, workspaceId int) (bool, error) {
	if userId == 1 && workspaceId == 1 {
		return true, nil
	}
	if userId == 1 && workspaceId == -1 {
		return false, errors.New("Internal server error")
	}
	return false, nil
}
