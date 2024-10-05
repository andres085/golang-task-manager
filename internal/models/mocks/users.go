package mocks

import "github.com/andres085/task_manager/internal/models"

type UserModel struct{}
type UserWithRole struct{}

var firstMockUser = models.UserWithRole{
	ID:        1,
	FirstName: "Test",
	LastName:  "McTester",
	Email:     "testmctesterson@mail.com",
	Role:      "ADMIN",
}

var secondMockUser = models.UserWithRole{
	ID:        2,
	FirstName: "Pete",
	LastName:  "Peterson",
	Email:     "pete@mail.com",
	Role:      "MEMBER",
}

func (m *UserModel) Insert(firstName, lastName, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (m *UserModel) GetUserToInvite(email string, workspaceId int) (*models.User, error) {
	return &models.User{}, nil
}

func (m *UserModel) AddUserToWorkspace(userId, workspaceId int) error {
	return nil
}

func (m *UserModel) GetWorkspaceUsers(workspaceId int) ([]models.UserWithRole, error) {
	return []models.UserWithRole{firstMockUser, secondMockUser}, nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "pa$$word" {
		return 1, nil
	}

	return 0, models.ErrInvalidCredentials
}

func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) RemoveUserFromWorkspace(workspaceId, userId int) (int, error) {
	return 1, nil
}
