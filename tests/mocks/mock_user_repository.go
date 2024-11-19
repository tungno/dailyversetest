// tests/mocks/mock_user_repository.go
package mocks

import (
	"context"
	"fmt"
	"proh2052-group6/pkg/models"
	"strings"
	"time"
)

type MockUserRepository struct {
	Users map[string]*models.User
}

func NewMockUserRepository(users map[string]*models.User) *MockUserRepository {
	return &MockUserRepository{Users: users}
}

func (mur *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if user, exists := mur.Users[email]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (mur *MockUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	for _, user := range mur.Users {
		if strings.ToLower(user.Username) == strings.ToLower(username) {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (mur *MockUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	if _, exists := mur.Users[user.Email]; exists {
		return fmt.Errorf("user already exists")
	}
	mur.Users[user.Email] = user
	return nil
}

func (mur *MockUserRepository) UpdateUser(ctx context.Context, email string, updates map[string]interface{}) error {
	user, exists := mur.Users[email]
	if !exists {
		return fmt.Errorf("user not found")
	}
	// Apply updates
	if otp, ok := updates["OTP"]; ok {
		user.OTP = otp.(string)
	}
	if otpExpiresAt, ok := updates["OTPExpiresAt"]; ok {
		user.OTPExpiresAt = otpExpiresAt.(time.Time)
	}
	if isVerified, ok := updates["IsVerified"]; ok {
		user.IsVerified = isVerified.(bool)
	}
	if password, ok := updates["Password"]; ok {
		user.Password = password.(string)
	}
	return nil
}

func (mur *MockUserRepository) SearchUsersByUsername(ctx context.Context, query string) ([]*models.User, error) {
	var users []*models.User
	queryLower := strings.ToLower(query)
	for _, user := range mur.Users {
		usernameLower := strings.ToLower(user.Username)
		if strings.HasPrefix(usernameLower, queryLower) {
			users = append(users, user)
		}
	}
	return users, nil
}
