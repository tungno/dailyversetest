// tests/mocks/mock_user_repository.go
package mocks

import (
	"context"
	"errors"
	"proh2052-group6/pkg/models"
)

type MockUserRepository struct {
	Users map[string]*models.User
}

func NewMockUserRepository(users map[string]*models.User) *MockUserRepository {
	return &MockUserRepository{Users: users}
}

func (mur *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, exists := mur.Users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (mur *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	for _, user := range mur.Users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, errors.New("user not found")
}
