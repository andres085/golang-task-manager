package models

import (
	"testing"

	"github.com/andres085/task_manager/internal/assert"
)

func TestUserInsertMethod(t *testing.T) {
	db := newTestDB(t)

	m := UserModel{db}
	err := m.Insert("Test", "McTester", "test@mail.com", "pa$$word")

	assert.NilError(t, err)
}

func TestUserAuthenticateMethod(t *testing.T) {
	db := newTestDB(t)

	email := "test@mail.com"
	password := "pa$$word"

	m := UserModel{db}
	err := m.Insert("Test", "McTester", email, password)
	if err != nil {
		t.Fatal(err)
	}

	row, err := m.Authenticate(email, password)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, row, 3)
}
